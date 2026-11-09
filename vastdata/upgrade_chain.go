// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

// UpgradeState migrates pre-WriteOnly state (schema v0) to schema v1: seed
// <field>_md5 from leftover secrets, then clear WriteOnlyFields so Terraform
// does not reject non-null write-only values during state upgrade.
func (r *Resource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	hints := r.EmptyManager().TfState().StateHints()
	if hints == nil || len(hints.WriteOnlyFields) == 0 {
		return nil
	}
	return map[int64]resource.StateUpgrader{
		0: {StateUpgrader: r.upgradeWriteOnlyStateFromV0},
	}
}

func (r *Resource) upgradeWriteOnlyStateFromV0(
	ctx context.Context,
	req resource.UpgradeStateRequest,
	resp *resource.UpgradeStateResponse,
) {
	if req.RawState == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("UpgradeState[%s]", r.managerName),
			"missing prior RawState",
		)
		return
	}

	manager, err := r.ManagerWithSchemaOnly(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("UpgradeState[%s]: schema", r.managerName),
			err.Error(),
		)
		return
	}

	sch := manager.TfState().SchemaDefinition().(rschema.Schema)
	hints := manager.TfState().StateHints()

	rawVal, err := req.RawState.UnmarshalWithOpts(sch.Type().TerraformType(ctx), tfprotov6.UnmarshalOpts{
		ValueFromJSONOpts: tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true},
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("UpgradeState[%s]: unmarshal prior state", r.managerName),
			err.Error(),
		)
		return
	}

	attrMap, err := is.FillFrameworkValues(rawVal, sch)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("UpgradeState[%s]: decode prior state", r.managerName),
			err.Error(),
		)
		return
	}

	// Concrete *TFState: SeedWriteOnlyMD5FromSecrets is upgrade-only and not
	// part of ResourceState (managers/lifecycle never need it).
	ts := is.NewTFStateMust(attrMap, sch, hints)
	ts.SeedWriteOnlyMD5FromSecrets()
	ts.ClearWriteOnlyFields()

	objType, ok := sch.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("UpgradeState[%s]", r.managerName),
			fmt.Sprintf("expected object schema type, got %T", sch.Type().TerraformType(ctx)),
		)
		return
	}

	vals := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, typ := range objType.AttributeTypes {
		v, present := ts.RawValues()[name]
		if !present || v == nil {
			vals[name] = tftypes.NewValue(typ, nil)
			continue
		}
		tv, convErr := v.ToTerraformValue(ctx)
		if convErr != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("UpgradeState[%s]: encode %q", r.managerName, name),
				convErr.Error(),
			)
			return
		}
		vals[name] = tv
	}

	resp.State = tfsdk.State{
		Schema: sch,
		Raw:    tftypes.NewValue(objType, vals),
	}
}
