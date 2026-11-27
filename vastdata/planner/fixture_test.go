// Copyright (c) HashiCorp, Inc.

package planner

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"

	vsd "github.com/vast-data/terraform-provider-vastdata/vastdata"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

// validValues are values for attributes that a validator or PrepareDesiredState checks in
// code rather than through a schema validator synth can probe.
var validValues = map[string]map[string]func(variant int) any{
	"compute_cluster_control":         {"action": oneOf("start", "stop")},
	"encryption_group_control":        {"action": oneOf("revoke", "deactivate")},
	"tenant_encryption_group_control": {"action": oneOf("revoke", "deactivate")},
	"nfs_certificate":                 {"certificate": pemCertificate},
	"s3_certificate":                  {"certificate": pemCertificate},
	"vms_certificate":                 {"certificate": pemCertificate},
}

func oneOf(values ...string) func(int) any {
	return func(variant int) any { return values[variant%len(values)] }
}

// pemCertificate returns a self-signed certificate; certificates of different variants differ.
func pemCertificate(variant int) any {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(int64(variant + 1)),
		Subject:      pkix.Name{CommonName: fmt.Sprintf("planner-%d.test", variant)},
		NotBefore:    time.Unix(1_700_000_000, 0),
		NotAfter:     time.Unix(1_900_000_000, 0),
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
}

// fixture plans one resource from raw Terraform values, the way Terraform Core holds them.
type fixture struct {
	t     *testing.T
	name  string
	sch   rschema.Schema
	typ   tftypes.Object
	hints *is.TFStateHints
}

func newFixture(t *testing.T, name string) *fixture {
	t.Helper()
	m, err := vsd.NewResourceManagerByName(context.Background(), name)
	require.NoError(t, err)
	sch, ok := m.TfState().SchemaDefinition().(rschema.Schema)
	require.True(t, ok)
	hints := m.TfState().StateHints()
	if hints == nil {
		hints = &is.TFStateHints{}
	}
	return &fixture{
		t: t, name: name, sch: sch, hints: hints,
		typ: sch.Type().TerraformType(context.Background()).(tftypes.Object),
	}
}

// object builds a resource object; attributes missing from vals are null.
func (f *fixture) object(vals map[string]tftypes.Value) tftypes.Value {
	out := make(map[string]tftypes.Value, len(f.typ.AttributeTypes))
	for k, ty := range f.typ.AttributeTypes {
		if v, ok := vals[k]; ok {
			out[k] = v
		} else {
			out[k] = tftypes.NewValue(ty, nil)
		}
	}
	return tftypes.NewValue(f.typ, out)
}

func (f *fixture) null() tftypes.Value { return tftypes.NewValue(f.typ, nil) }

// state turns a resource object into the provider's resource state.
func (f *fixture) state(v tftypes.Value) is.ResourceState {
	f.t.Helper()
	raw, err := is.FillFrameworkValues(v, f.sch)
	require.NoError(f.t, err)
	return vsd.ResourceManagerPrototypes()[f.name].NewResourceManager(raw, f.sch).TfState()
}

func (f *fixture) plan(prior tftypes.Value, config map[string]tftypes.Value) (*Result, error) {
	var priorSt is.ResourceState
	if !prior.IsNull() {
		priorSt = f.state(prior)
	}
	return Plan(context.Background(), f.name, priorSt, f.state(f.object(config)))
}

func (f *fixture) mustPlan(prior tftypes.Value, config map[string]tftypes.Value) *Result {
	f.t.Helper()
	res, err := f.plan(prior, config)
	require.NoError(f.t, err)
	return res
}

// apply plans the creation of config and returns the state after apply: values the plan
// leaves unknown are filled in the way a server fills them.
func (f *fixture) apply(config map[string]tftypes.Value) (tftypes.Value, error) {
	res, err := f.plan(f.null(), config)
	if err != nil {
		return tftypes.Value{}, err
	}
	if res.Action != Create {
		return tftypes.Value{}, fmt.Errorf("planning without prior state: action %s", res.Action)
	}
	return applied(res.planned)
}

func (f *fixture) mustApply(config map[string]tftypes.Value) tftypes.Value {
	f.t.Helper()
	v, err := f.apply(config)
	require.NoError(f.t, err)
	return v
}

// attr returns the value of top-level attribute k.
func attrOf(t *testing.T, obj tftypes.Value, k string) tftypes.Value {
	t.Helper()
	vals, err := objectValues(obj)
	require.NoError(t, err)
	v, ok := vals[k]
	require.True(t, ok, "no attribute %q", k)
	return v
}

// with sets top-level attribute k of obj.
func with(t *testing.T, obj tftypes.Value, k string, v tftypes.Value) tftypes.Value {
	t.Helper()
	vals, err := objectValues(obj)
	require.NoError(t, err)
	vals[k] = v
	return tftypes.NewValue(obj.Type(), vals)
}

// minimalConfig sets every required attribute.
func (f *fixture) minimalConfig() map[string]tftypes.Value {
	out := map[string]tftypes.Value{}
	for k, a := range f.sch.Attributes {
		if a.IsRequired() {
			out[k] = f.value(k, 0)
		}
	}
	return out
}

// value returns a configurable value for top-level attribute k. Values of different
// variants differ.
func (f *fixture) value(k string, variant int) tftypes.Value {
	ty := f.typ.AttributeTypes[k]
	if fn, ok := validValues[f.name][k]; ok {
		return tftypes.NewValue(ty, fn(variant))
	}
	return synth(f.sch.Attributes[k], ty, k, variant)
}

func (f *fixture) isWriteOnly(k string) bool {
	if f.sch.Attributes[k].IsWriteOnly() {
		return true
	}
	for _, w := range f.hints.WriteOnlyFields {
		if w == k {
			return true
		}
	}
	return false
}

func copyConfig(c map[string]tftypes.Value) map[string]tftypes.Value {
	out := make(map[string]tftypes.Value, len(c))
	for k, v := range c {
		out[k] = v
	}
	return out
}

// synth returns a known value for attribute a. Nested objects get their required
// attributes, or their first optional one when none is required.
func synth(a rschema.Attribute, ty tftypes.Type, seed string, variant int) tftypes.Value {
	attrs, mode, ok := nestedOf(a)
	if !ok {
		switch sa := a.(type) {
		case rschema.StringAttribute:
			if len(sa.Validators) > 0 {
				return tftypes.NewValue(ty, validString(sa.Validators, seed, variant))
			}
		case rschema.Int64Attribute:
			if len(sa.Validators) > 0 {
				return tftypes.NewValue(ty, big.NewFloat(float64(validInt64(sa.Validators, variant))))
			}
		}
		return synthType(ty, seed, variant)
	}
	switch mode {
	case nestingSingle:
		return synthObject(attrs, ty.(tftypes.Object), seed, variant)
	case nestingList:
		elem := synthObject(attrs, ty.(tftypes.List).ElementType.(tftypes.Object), seed, variant)
		return tftypes.NewValue(ty, []tftypes.Value{elem})
	case nestingSet:
		elem := synthObject(attrs, ty.(tftypes.Set).ElementType.(tftypes.Object), seed, variant)
		return tftypes.NewValue(ty, []tftypes.Value{elem})
	default:
		elem := synthObject(attrs, ty.(tftypes.Map).ElementType.(tftypes.Object), seed, variant)
		return tftypes.NewValue(ty, map[string]tftypes.Value{"k": elem})
	}
}

func synthObject(attrs map[string]rschema.Attribute, ty tftypes.Object, seed string, variant int) tftypes.Value {
	var required, optional []string
	for k, a := range attrs {
		switch {
		case a.IsRequired():
			required = append(required, k)
		case a.IsOptional():
			optional = append(optional, k)
		}
	}
	sort.Strings(optional)
	set := required
	if len(set) == 0 && len(optional) > 0 {
		set = optional[:1]
	}
	vals := make(map[string]tftypes.Value, len(ty.AttributeTypes))
	for k, at := range ty.AttributeTypes {
		vals[k] = tftypes.NewValue(at, nil)
	}
	for _, k := range set {
		vals[k] = synth(attrs[k], ty.AttributeTypes[k], seed+"."+k, variant)
	}
	return tftypes.NewValue(ty, vals)
}

func synthType(ty tftypes.Type, seed string, variant int) tftypes.Value {
	switch {
	case ty.Is(tftypes.String):
		return tftypes.NewValue(ty, fmt.Sprintf("%s-%d", seed, variant))
	case ty.Is(tftypes.Number):
		return tftypes.NewValue(ty, big.NewFloat(float64(1+variant)))
	case ty.Is(tftypes.Bool):
		return tftypes.NewValue(ty, variant%2 == 0)
	}
	switch tt := ty.(type) {
	case tftypes.List:
		return tftypes.NewValue(ty, []tftypes.Value{synthType(tt.ElementType, seed, variant)})
	case tftypes.Set:
		return tftypes.NewValue(ty, []tftypes.Value{synthType(tt.ElementType, seed, variant)})
	case tftypes.Map:
		return tftypes.NewValue(ty, map[string]tftypes.Value{"k": synthType(tt.ElementType, seed, variant)})
	case tftypes.Object:
		vals := make(map[string]tftypes.Value, len(tt.AttributeTypes))
		for k, at := range tt.AttributeTypes {
			vals[k] = synthType(at, seed+"."+k, variant)
		}
		return tftypes.NewValue(ty, vals)
	}
	panic(fmt.Sprintf("synthType: unsupported type %s", ty))
}

var quoted = regexp.MustCompile(`"([^"]*)"`)

// validString returns a value every validator accepts: an enumeration value a validator
// lists, or one of a few common formats (paths, durations, timestamps).
func validString(validators []validator.String, seed string, variant int) string {
	var candidates []string
	for _, v := range validators {
		desc := v.Description(context.Background())
		if strings.Contains(desc, "one of") {
			for _, m := range quoted.FindAllStringSubmatch(desc, -1) {
				candidates = append(candidates, m[1])
			}
		}
	}
	if n := len(candidates); n > 0 {
		candidates = append(candidates[variant%n:], candidates[:variant%n]...)
	}
	n := 1 + variant
	candidates = append(candidates,
		fmt.Sprintf("%s-%d", seed, variant),
		fmt.Sprintf("/%s-%d/", seed, variant),
		fmt.Sprintf("2025-07-06T13:45:%02dZ", n),
		fmt.Sprintf("%dD", n), fmt.Sprintf("%dW", n), fmt.Sprintf("%dh", n), fmt.Sprintf("%dy", n),
		fmt.Sprintf("%dd", n), fmt.Sprintf("%dm", n), fmt.Sprintf("0%d:00:00", n),
	)
	for _, c := range candidates {
		ok := true
		for _, v := range validators {
			var resp validator.StringResponse
			v.ValidateString(context.Background(), validator.StringRequest{ConfigValue: types.StringValue(c)}, &resp)
			if resp.Diagnostics.HasError() {
				ok = false
				break
			}
		}
		if ok {
			return c
		}
	}
	return fmt.Sprintf("%s-%d", seed, variant)
}

func validInt64(validators []validator.Int64, variant int) int64 {
	for _, c := range []int64{int64(1 + variant), int64(10 + variant), int64(100 + variant), int64(variant), int64(1024 + variant)} {
		ok := true
		for _, v := range validators {
			var resp validator.Int64Response
			v.ValidateInt64(context.Background(), validator.Int64Request{ConfigValue: types.Int64Value(c)}, &resp)
			if resp.Diagnostics.HasError() {
				ok = false
				break
			}
		}
		if ok {
			return c
		}
	}
	return int64(1 + variant)
}

// applied fills what a plan leaves unknown the way a server does: primitives get a value,
// collections come back empty and objects with null attributes.
func applied(planned tftypes.Value) (tftypes.Value, error) {
	return tftypes.Transform(planned, func(_ *tftypes.AttributePath, v tftypes.Value) (tftypes.Value, error) {
		if v.IsKnown() {
			return v, nil
		}
		switch tt := v.Type().(type) {
		case tftypes.List:
			return tftypes.NewValue(tt, []tftypes.Value{}), nil
		case tftypes.Set:
			return tftypes.NewValue(tt, []tftypes.Value{}), nil
		case tftypes.Map:
			return tftypes.NewValue(tt, map[string]tftypes.Value{}), nil
		case tftypes.Object:
			return emptyObject(tt), nil
		}
		return synthType(v.Type(), "server", 7), nil
	})
}

// rejectedByValidation reports whether err is the provider refusing a configuration.
func rejectedByValidation(err error) bool {
	var de *DiagnosticsError
	return errors.As(err, &de) && de.Stage == "ValidateResourceConfig"
}

func sortedAttributeNames(attrs map[string]rschema.Attribute) []string {
	out := make([]string, 0, len(attrs))
	for k := range attrs {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func sortedResourceNames() []string {
	names := append([]string(nil), vsd.ResourceNames()...)
	sort.Strings(names)
	return names
}
