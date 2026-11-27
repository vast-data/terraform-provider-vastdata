// Copyright (c) HashiCorp, Inc.

// Package planner runs Terraform's plan for vastdata resources in-process, for runtimes that
// use this provider as a library instead of through Terraform Core.
//
// The provider half of a plan — defaults, plan modifiers (UseStateForUnknown,
// RequiresReplace), ModifyPlan / PrepareDesiredState, write-only and *_md5 handling — runs
// unchanged through the provider's own protocol server. The Core half is small and lives
// here: validate the config, merge it with prior state into the proposed new state
// (ProposedNew), keep the replacements whose value really changes, and decide the action from
// planned vs prior state. The result is the same change Terraform would report for
// `terraform plan`.
package planner

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/provider"
)

// Action is what Terraform would do to an existing resource.
type Action string

const (
	NoOp    Action = "no-op"
	Update  Action = "update"
	Replace Action = "replace"
	Create  Action = "create"
)

// Result is the outcome of planning one resource: the planned state, and how it differs from
// prior.
//
// The planned state is what the provider's Create or Update receives as req.Plan. Like there,
// write-only values are null: is.OverlayWriteOnly takes them from the configuration.
type Result struct {
	is.ResourceState

	Action Action
	// Changed lists top-level attributes whose planned value differs from prior: a different
	// value, or a prior value that becomes known only after apply. A changed nested value may
	// still hold parts that are known only after apply.
	Changed []string
	// Unknown lists top-level attributes whose whole planned value is known only after apply,
	// including those without a prior value.
	Unknown []string
	// RequiresReplace lists top-level attributes that force replacement.
	RequiresReplace []string

	planned tftypes.Value
}

// DiagnosticsError carries the error diagnostics the provider returned.
type DiagnosticsError struct {
	Stage       string
	Diagnostics []*tfprotov6.Diagnostic
}

func (e *DiagnosticsError) Error() string {
	var parts []string
	for _, d := range e.Diagnostics {
		if d.Severity != tfprotov6.DiagnosticSeverityError {
			continue
		}
		msg := d.Summary
		if d.Detail != "" {
			msg += ": " + d.Detail
		}
		if d.Attribute != nil && len(d.Attribute.Steps()) > 0 {
			msg = d.Attribute.String() + ": " + msg
		}
		parts = append(parts, msg)
	}
	return fmt.Sprintf("%s: %s", e.Stage, strings.Join(parts, "; "))
}

// Planner holds one in-process provider server. It is safe for concurrent use, like the
// provider process Terraform Core talks to.
type Planner struct {
	server tfprotov6.ProviderServer
}

// New builds a planner around the provider's protocol server. The server is not configured
// with a client: planning never calls the VAST API.
func New(version string) (*Planner, error) {
	srv, err := providerserver.NewProtocol6WithError(provider.New(version)())()
	if err != nil {
		return nil, fmt.Errorf("start embedded provider server: %w", err)
	}
	return &Planner{server: srv}, nil
}

var defaultPlanner = sync.OnceValues(func() (*Planner, error) { return New("embedded") })

// Plan plans with a shared planner. See (*Planner).Plan.
func Plan(ctx context.Context, name string, prior, config is.ResourceState) (*Result, error) {
	p, err := defaultPlanner()
	if err != nil {
		return nil, err
	}
	return p.Plan(ctx, name, prior, config)
}

// Plan returns what Terraform would plan for resource `name` (registry name, e.g. "vip_pool").
//
// prior is the state after refresh (what Terraform keeps in state); a nil or empty prior
// plans a creation. config holds only what the user configured, including write-only values.
func (p *Planner) Plan(ctx context.Context, name string, prior, config is.ResourceState) (*Result, error) {
	if is.IsNilState(config) || config.IsEmpty() {
		return nil, errors.New("config is required")
	}
	sch, ok := config.SchemaDefinition().(rschema.Schema)
	if !ok {
		return nil, fmt.Errorf("resource %q: config state has no resource schema", name)
	}
	if len(sch.Blocks) > 0 {
		return nil, fmt.Errorf("resource %q: nested blocks are not supported by the planner", name)
	}
	typeName := "vastdata_" + name
	objType := sch.Type().TerraformType(ctx)

	configVal, err := toTerraformObject(ctx, sch, config)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	if err := checkConfigShape(sch.Attributes, configVal, tftypes.NewAttributePath()); err != nil {
		return nil, err
	}
	priorVal := tftypes.NewValue(objType, nil)
	if !is.IsNilState(prior) && !prior.IsEmpty() {
		if priorVal, err = toTerraformObject(ctx, sch, prior); err != nil {
			return nil, fmt.Errorf("prior state: %w", err)
		}
	}

	configDV, err := tfprotov6.NewDynamicValue(objType, configVal)
	if err != nil {
		return nil, err
	}
	vresp, err := p.server.ValidateResourceConfig(ctx, &tfprotov6.ValidateResourceConfigRequest{
		TypeName:           typeName,
		Config:             &configDV,
		ClientCapabilities: &tfprotov6.ValidateResourceConfigClientCapabilities{WriteOnlyAttributesAllowed: true},
	})
	if err != nil {
		return nil, err
	}
	if hasErrors(vresp.Diagnostics) {
		return nil, &DiagnosticsError{Stage: "ValidateResourceConfig", Diagnostics: vresp.Diagnostics}
	}

	proposed, err := proposedNew(sch.Attributes, priorVal, configVal)
	if err != nil {
		return nil, fmt.Errorf("proposed new state: %w", err)
	}
	priorDV, err := tfprotov6.NewDynamicValue(objType, priorVal)
	if err != nil {
		return nil, err
	}
	proposedDV, err := tfprotov6.NewDynamicValue(objType, proposed)
	if err != nil {
		return nil, err
	}
	presp, err := p.server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
		TypeName:           typeName,
		PriorState:         &priorDV,
		ProposedNewState:   &proposedDV,
		Config:             &configDV,
		ClientCapabilities: &tfprotov6.PlanResourceChangeClientCapabilities{},
	})
	if err != nil {
		return nil, err
	}
	if hasErrors(presp.Diagnostics) {
		return nil, &DiagnosticsError{Stage: "PlanResourceChange", Diagnostics: presp.Diagnostics}
	}
	if presp.PlannedState == nil {
		return nil, errors.New("PlanResourceChange returned no planned state")
	}
	planned, err := presp.PlannedState.Unmarshal(objType)
	if err != nil {
		return nil, fmt.Errorf("decode planned state: %w", err)
	}

	res := &Result{RequiresReplace: topLevelNames(requiresReplace(presp.RequiresReplace, priorVal, planned)), planned: planned}
	if res.Changed, res.Unknown, err = changedAttributes(planned, priorVal); err != nil {
		return nil, err
	}
	res.Action = decideAction(planned, priorVal, res.RequiresReplace)
	raw, err := is.FillFrameworkValues(planned, sch)
	if err != nil {
		return nil, fmt.Errorf("planned state: %w", err)
	}
	res.ResourceState = is.NewResourceState(raw, sch, config.StateHints())
	return res, nil
}

// decideAction is Terraform Core's choice of action from planned vs prior state.
func decideAction(planned, prior tftypes.Value, requiresReplace []string) Action {
	switch {
	case prior.IsNull():
		return Create
	case planned.Equal(prior):
		return NoOp
	case len(requiresReplace) > 0:
		return Replace
	default:
		return Update
	}
}

// requiresReplace keeps the paths the provider asked to replace for whose planned value
// differs from prior, as Terraform Core does. A RequiresReplace modifier can flag an
// attribute whose value a later modifier (UseStateForUnknown) restores to the prior one.
func requiresReplace(paths []*tftypes.AttributePath, prior, planned tftypes.Value) []*tftypes.AttributePath {
	if prior.IsNull() {
		return nil
	}
	var out []*tftypes.AttributePath
	for _, p := range paths {
		pv, pok := valueAt(prior, p)
		nv, nok := valueAt(planned, p)
		switch {
		case !pok && !nok:
			continue
		case !pok:
			pv = tftypes.NewValue(nv.Type(), nil)
		case !nok:
			nv = tftypes.NewValue(pv.Type(), nil)
		}
		if !nv.IsFullyKnown() || !nv.Equal(pv) {
			out = append(out, p)
		}
	}
	return out
}

// valueAt returns the value at path p of v. A path through a null value has none.
func valueAt(v tftypes.Value, p *tftypes.AttributePath) (tftypes.Value, bool) {
	got, _, err := tftypes.WalkAttributePath(v, p)
	if err != nil {
		return tftypes.Value{}, false
	}
	tv, ok := got.(tftypes.Value)
	return tv, ok
}

// toTerraformObject converts a resource state into the object value of schema sch.
func toTerraformObject(ctx context.Context, sch rschema.Schema, st is.ResourceState) (tftypes.Value, error) {
	raw := st.RawValues()
	vals := make(map[string]tftypes.Value, len(sch.Attributes))
	for k, a := range sch.Attributes {
		v, ok := raw[k]
		if !ok || v == nil {
			vals[k] = tftypes.NewValue(a.GetType().TerraformType(ctx), nil)
			continue
		}
		tv, err := v.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.Value{}, fmt.Errorf("attribute %q: %w", k, err)
		}
		vals[k] = tv
	}
	return tftypes.NewValue(sch.Type().TerraformType(ctx), vals), nil
}

// changedAttributes compares planned and prior top-level attributes.
func changedAttributes(planned, prior tftypes.Value) (changed, unknown []string, err error) {
	if planned.IsNull() || !planned.IsKnown() {
		return nil, nil, nil
	}
	pv := map[string]tftypes.Value{}
	if err := planned.As(&pv); err != nil {
		return nil, nil, err
	}
	qv := map[string]tftypes.Value{}
	if !prior.IsNull() {
		if err := prior.As(&qv); err != nil {
			return nil, nil, err
		}
	}
	for k, v := range pv {
		q, ok := qv[k]
		if ok && v.Equal(q) {
			continue
		}
		if !v.IsKnown() {
			unknown = append(unknown, k)
			if !ok || q.IsNull() {
				continue
			}
		}
		changed = append(changed, k)
	}
	sort.Strings(changed)
	sort.Strings(unknown)
	return changed, unknown, nil
}

// checkConfigShape is the schema check Core does while decoding configuration: a required
// attribute must be set ("Missing required argument") and a computed-only attribute must
// not be ("Invalid Configuration for Read-Only Attribute"). The framework's
// ValidateResourceConfig does not repeat either check.
func checkConfigShape(attrs map[string]rschema.Attribute, config tftypes.Value, at *tftypes.AttributePath) error {
	if config.IsNull() || !config.IsKnown() {
		return nil
	}
	vals := map[string]tftypes.Value{}
	if err := config.As(&vals); err != nil {
		return err
	}
	var missing []string
	for k, a := range attrs {
		v, ok := vals[k]
		p := at.WithAttributeName(k)
		if !ok || v.IsNull() {
			if a.IsRequired() {
				missing = append(missing, p.String())
			}
			continue
		}
		if a.IsComputed() && !a.IsOptional() && !a.IsRequired() {
			return fmt.Errorf("%s: invalid configuration for read-only attribute: the value is decided by the provider", p)
		}
		if err := forEachNestedObject(a, v, func(nested map[string]rschema.Attribute, obj tftypes.Value) error {
			return checkConfigShape(nested, obj, p)
		}); err != nil {
			return err
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("missing required argument: %s", strings.Join(missing, ", "))
	}
	return nil
}

func topLevelNames(paths []*tftypes.AttributePath) []string {
	seen := map[string]bool{}
	var out []string
	for _, p := range paths {
		if p == nil || len(p.Steps()) == 0 {
			continue
		}
		if n, ok := p.Steps()[0].(tftypes.AttributeName); ok && !seen[string(n)] {
			seen[string(n)] = true
			out = append(out, string(n))
		}
	}
	sort.Strings(out)
	return out
}

func hasErrors(diags []*tfprotov6.Diagnostic) bool {
	for _, d := range diags {
		if d.Severity == tfprotov6.DiagnosticSeverityError {
			return true
		}
	}
	return false
}
