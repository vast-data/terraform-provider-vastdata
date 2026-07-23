// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"
	"time"

	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var QuotaGroupSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"quotagroups",
	http.MethodGet,
	"quotagroups",
)

// quotaGroupActionFields are provider-only fields not returned by the API.
// They are preserved in Terraform state so the user's configuration is not
// overwritten by API reads.
var quotaGroupActionFields = []string{
	"quotas_ids",
	"refresh_user_quotas",
	"reset_grace_period",
}

type QuotaGroup struct {
	tfstate *is.TFState
}

func (m *QuotaGroup) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &QuotaGroup{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: QuotaGroupSchemaRef,
			PreserveUserValueFields: append([]string{
				"default_group_quota",
				"default_user_quota",
				"grace_period",
				"group_quotas",
				"user_quotas",
			}, quotaGroupActionFields...),
			PreserveUserValueFieldsWhenApiReturnsNull: []string{
				"hard_limit",
				"soft_limit",
				"hard_limit_inodes",
				"soft_limit_inodes",
			},
			CreateOnlyFields: []string{"is_physical_quota"},
			AdditionalSchemaAttributes: map[string]any{
				// quotas_ids drives the assign_quotas sub-action.
				// Send the desired list of quota IDs; the API response shows them
				// as the nested "quotas" array (name+id objects).
				"quotas_ids": rschema.ListAttribute{
					ElementType: types.Int64Type,
					Optional:    true,
					Description: "List of quota IDs to assign to this quota group. " +
						"Calls PATCH /quotagroups/{id}/assign_quotas/ after create or update.",
				},
				// refresh_user_quotas triggers PATCH /quotagroups/{id}/refresh_user_quotas/.
				"refresh_user_quotas": rschema.BoolAttribute{
					Optional: true,
					Description: "When true, refreshes all user quota usage counters for this group " +
						"after every create or update operation.",
				},
				// reset_grace_period triggers PATCH /quotagroups/{id}/reset_grace_period/.
				"reset_grace_period": rschema.BoolAttribute{
					Optional: true,
					Description: "When true, resets the grace period countdown for this quota group " +
						"after every create or update operation.",
				},
			},
		},
	)}
}

func (m *QuotaGroup) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &QuotaGroup{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: QuotaGroupSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"quotas_ids": dschema.ListAttribute{
					ElementType: types.Int64Type,
					Optional:    true,
					Description: "List of quota IDs assigned to this quota group.",
				},
				"refresh_user_quotas": dschema.BoolAttribute{
					Optional:    true,
					Description: "Trigger field for refreshing user quota counters.",
				},
				"reset_grace_period": dschema.BoolAttribute{
					Optional:    true,
					Description: "Trigger field for resetting the grace period countdown.",
				},
			},
		}),
	}
}

func (m *QuotaGroup) TfState() *is.TFState {
	return m.tfstate
}

func (m *QuotaGroup) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.QuotaGroups
}

func (m *QuotaGroup) PrepareUpdateResource(_ context.Context, plan PrepareUpdateResource, _ *VMSRest) error {
	planTs := plan.(*QuotaGroup).tfstate
	if !planTs.IsKnownAndNotNull("is_physical_quota") {
		return nil
	}
	return ensureNotChanged(m.tfstate, planTs, "is_physical_quota")
}

// AfterCreateResource calls action sub-endpoints after a quota group is created.
func (m *QuotaGroup) AfterCreateResource(ctx context.Context, rest *VMSRest, record Record) error {
	return m.runActions(ctx, rest, m.tfstate, record)
}

// AfterUpdateResource calls action sub-endpoints after a quota group is updated.
func (m *QuotaGroup) AfterUpdateResource(ctx context.Context, plan AfterUpdateResource, rest *VMSRest, record Record) error {
	planTs := plan.(*QuotaGroup).TfState()
	return m.runActions(ctx, rest, planTs, record)
}

// TransformRequestBody strips provider-only action trigger fields before sending to the API.
// The VAST API does not accept quotas_ids, refresh_user_quotas, or reset_grace_period in
// its POST/PATCH body — they are handled separately via AfterCreateResource/AfterUpdateResource.
func (m *QuotaGroup) TransformRequestBody(body params) params {
	delete(body, "quotas_ids")
	delete(body, "refresh_user_quotas")
	delete(body, "reset_grace_period")
	return body
}

// runActions calls assign_quotas, refresh_user_quotas, and reset_grace_period
// endpoints when their corresponding trigger fields are set.
func (m *QuotaGroup) runActions(ctx context.Context, rest *VMSRest, ts *is.TFState, record Record) error {
	id := record.RecordID()

	// Assign quotas when quotas_ids is provided.
	if ts.IsKnownAndNotNull("quotas_ids") {
		quotaIds := ts.ToSlice("quotas_ids")
		body := params{"quotas_ids": quotaIds}
		if _, err := rest.QuotaGroups.QuotaGroupAssignQuotasWithContext_PATCH(ctx, id, body, 5*time.Minute); err != nil {
			return fmt.Errorf("failed to assign quotas to quota group %v: %w", id, err)
		}
	}

	// Refresh user quota usage counters.
	if ts.IsKnownAndNotNull("refresh_user_quotas") && ts.Bool("refresh_user_quotas") {
		if err := rest.QuotaGroups.QuotaGroupRefreshUserQuotasWithContext_PATCH(ctx, id, nil); err != nil {
			return fmt.Errorf("failed to refresh user quotas for quota group %v: %w", id, err)
		}
	}

	// Reset the grace period countdown.
	if ts.IsKnownAndNotNull("reset_grace_period") && ts.Bool("reset_grace_period") {
		if err := rest.QuotaGroups.QuotaGroupResetGracePeriodWithContext_PATCH(ctx, id, nil); err != nil {
			return fmt.Errorf("failed to reset grace period for quota group %v: %w", id, err)
		}
	}

	return nil
}
