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

type EncryptionGroupControl struct {
	tfstate *is.TFState
}

func (m *EncryptionGroupControl) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &EncryptionGroupControl{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			TFStateHintsForCustom: &is.TFStateHintsForCustom{
				Description: "Control operations on encryption groups. This resource allows you to perform actions like revoke, deactivate, reinstate, or rotate keys on encryption groups.",
				SchemaAttributes: map[string]any{
					"id": rschema.Int64Attribute{
						Required:    true,
						Description: "ID of the encryption group to control.",
						PlanModifiers: []planmodifiers.Int64{
							int64planmodifier.RequiresReplace(),
						},
					},
					"action": rschema.StringAttribute{
						Required:    true,
						Description: "Action to perform on the encryption group. Valid values: revoke, deactivate, reinstate, rotate_key.",
						PlanModifiers: []planmodifiers.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
		},
	)}
}

func (m *EncryptionGroupControl) TfState() *is.TFState {
	return m.tfstate
}

func (m *EncryptionGroupControl) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.EncryptionGroups
}

func (m *EncryptionGroupControl) ValidateResourceConfig(context.Context) error {
	return ValidateFieldIsOneOf(m.tfstate, "action", "revoke", "deactivate", "reinstate", "rotate_key")
}

func (m *EncryptionGroupControl) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.performAction(ctx, rest)
}

func (m *EncryptionGroupControl) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	return m.performAction(ctx, rest)
}

func (m *EncryptionGroupControl) DeleteResource(ctx context.Context, rest *VMSRest) error {
	// No-op for encryption group control, as it cannot be deleted.
	return nil
}

// performAction executes the specified action on the encryption group
func (m *EncryptionGroupControl) performAction(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	var (
		encryptionGroupId = m.tfstate.Int64("id")
		action            = m.tfstate.String("action")
		err               error
	)

	if encryptionGroupId == 0 {
		return nil, fmt.Errorf("failed to get encryption group ID: ID is empty")
	}

	if action == "" {
		return nil, fmt.Errorf("failed to get action: action is empty")
	}

	switch action {
	case "revoke":
		_, err = rest.EncryptionGroups.RevokeEncryptionGroupWithContext(ctx, encryptionGroupId)
	case "deactivate":
		_, err = rest.EncryptionGroups.DeactivateEncryptionGroupWithContext(ctx, encryptionGroupId)
	case "reinstate":
		_, err = rest.EncryptionGroups.ReinstateEncryptionGroupWithContext(ctx, encryptionGroupId)
	case "rotate_key":
		_, err = rest.EncryptionGroups.RotateEncryptionGroupKeyWithContext(ctx, encryptionGroupId)
	default:
		return nil, fmt.Errorf("invalid action '%s'. Valid actions are: revoke, deactivate, reinstate, rotate_key", action)
	}

	return nil, err
}
