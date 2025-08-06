// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var UserTenantDataSchemaRef = is.NewSchemaReference(
	http.MethodPatch,
	"users/{id}/tenant_data",
	http.MethodGet,
	"users/{id}/tenant_data",
)

type UserTenantData struct {
	tfstate *is.TFState
}

func (m *UserTenantData) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &UserTenantData{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: UserTenantDataSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"user_id": rschema.Int64Attribute{
					Required:    true,
					Description: "The ID of the user to manage tenant data for.",
				},
				"s3_policies_ids": rschema.SetAttribute{
					ElementType: types.Int64Type,
					Optional:    true,
					Computed:    true,
					Description: "IDs of S3 policies to attach to the user.",
				},
			},
		},
	)}
}

func (m *UserTenantData) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &UserTenantData{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: UserTenantDataSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"user_id": dschema.Int64Attribute{
					Required:    true,
					Description: "The ID of the user to manage tenant data for.",
				},
			},
		},
	)}
}

func (m *UserTenantData) TfState() *is.TFState {
	return m.tfstate
}

func (m *UserTenantData) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Users
}

func (m *UserTenantData) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	userID := m.tfstate.Int64("user_id")
	return rest.Users.GetTenantDataWithContext(ctx, userID)
}

func (m *UserTenantData) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.ReadDatasource(ctx, rest)
}

func (m *UserTenantData) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	return ensureUserTenantDataUpdatedWith(ctx, ts, ts, rest)
}

func (m *UserTenantData) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	stateTs := m.tfstate
	planTs := plan.(*UserTenantData).TfState()
	return ensureUserTenantDataUpdatedWith(ctx, stateTs, planTs, rest)
}

func (m *UserTenantData) DeleteResource(ctx context.Context, rest *VMSRest) error {
	// No-op for user tenant data, as it cannot be deleted.
	return nil
}

// ensureUserTenantDataUpdatedWith verifies if the given UserTenantData (looked up by state)
// needs to be updated with new fields and performs the update if necessary.
//
// This is used in both CreateResource and UpdateResource for UserTenantData.
func ensureUserTenantDataUpdatedWith(ctx context.Context, stateTs, fieldsTs *is.TFState, rest *VMSRest) (DisplayableRecord, error) {
	// Get user ID from tfstate
	userID := stateTs.Int64("user_id")
	if userID == 0 {
		return nil, fmt.Errorf("failed to get user ID: user ID is empty")
	}

	// Create params with the tenant data fields
	params := params{}
	if ok := fieldsTs.SetToMapIfAvailable(
		params,
		"tenant_id",
		"allow_create_bucket",
		"allow_delete_bucket",
		"s3_superuser",
		"s3_policies_ids",
	); ok {
		// Use the custom API method to update tenant data
		record, err := rest.Users.UpdateTenantDataWithContext(ctx, userID, params)
		if err != nil {
			return nil, fmt.Errorf("failed to update tenant data: %w", err)
		}
		return record, nil
	}

	// If no fields to update, just get the current tenant data
	record, err := rest.Users.GetTenantDataWithContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant data: %w", err)
	}

	return record, nil
}
