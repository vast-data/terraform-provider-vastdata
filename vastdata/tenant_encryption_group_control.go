// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	planmodifiers "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

type TenantEncryptionGroupControl struct {
	tfstate *is.TFState
}

func (m *TenantEncryptionGroupControl) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &TenantEncryptionGroupControl{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			TFStateHintsForCustom: &is.TFStateHintsForCustom{
				Description: "Control operations on tenant encryption groups. This resource allows you to perform actions like revoke, deactivate, reinstate, or rotate keys on tenant encryption groups.",
				SchemaAttributes: map[string]any{
					"id": rschema.Int64Attribute{
						Required:    true,
						Description: "ID of the tenant to control encryption groups for.",
						PlanModifiers: []planmodifiers.Int64{
							int64planmodifier.RequiresReplace(),
						},
					},
					"action": rschema.StringAttribute{
						Required:    true,
						Description: "Action to perform on the tenant encryption group. Valid values: revoke, deactivate, reinstate, rotate_key.",
						PlanModifiers: []planmodifiers.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
		},
	)}
}

func (m *TenantEncryptionGroupControl) TfState() *is.TFState {
	return m.tfstate
}

func (m *TenantEncryptionGroupControl) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Tenants
}

func (m *TenantEncryptionGroupControl) ValidateResourceConfig(context.Context) error {
	return ValidateFieldIsOneOf(m.tfstate, "action", "revoke", "deactivate", "reinstate", "rotate_key")
}

func (m *TenantEncryptionGroupControl) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.performAction(ctx, rest)
}

func (m *TenantEncryptionGroupControl) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	return m.performAction(ctx, rest)
}

func (m *TenantEncryptionGroupControl) DeleteResource(ctx context.Context, rest *VMSRest) error {
	// No-op for tenant encryption group control, as it cannot be deleted.
	return nil
}

// performAction executes the specified action on the tenant encryption group
func (m *TenantEncryptionGroupControl) performAction(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	var (
		tenantId = m.tfstate.Int64("id")
		action   = m.tfstate.String("action")
		err      error
	)

	if tenantId == 0 {
		return nil, fmt.Errorf("failed to get tenant ID: tenant ID is empty")
	}

	if action == "" {
		return nil, fmt.Errorf("failed to get action: action is empty")
	}

	switch action {
	case "revoke":
		_, err = rest.Tenants.RevokeEncryptionGroupWithContext(ctx, tenantId)
	case "deactivate":
		_, err = rest.Tenants.DeactivateEncryptionGroupWithContext(ctx, tenantId)
	case "reinstate":
		_, err = rest.Tenants.ReinstateEncryptionGroupWithContext(ctx, tenantId)
	case "rotate_key":
		_, err = rest.Tenants.RotateEncryptionGroupKeyWithContext(ctx, tenantId)
	default:
		return nil, fmt.Errorf("invalid action '%s'. Valid actions are: revoke, deactivate, reinstate, rotate_key", action)
	}

	return nil, err
}
