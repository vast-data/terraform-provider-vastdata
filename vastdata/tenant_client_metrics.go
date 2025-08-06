// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"

	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var TenantClientMetricsSchemaRef = is.NewSchemaReference(
	http.MethodPatch,
	"tenants/{id}/client_metrics",
	http.MethodGet,
	"tenants/{id}/client_metrics",
)

type TenantClientMetrics struct {
	tfstate *is.TFState
}

func (m *TenantClientMetrics) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &TenantClientMetrics{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: TenantClientMetricsSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"tenant_id": rschema.Int64Attribute{
					Required:    true,
					Description: "ID of the tenant to manage client metrics for.",
				},
			},
		},
	)}
}

func (m *TenantClientMetrics) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &TenantClientMetrics{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: TenantClientMetricsSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"tenant_id": dschema.Int64Attribute{
					Required:    true,
					Description: "ID of the tenant to manage client metrics for.",
				},
			},
		}),
	}
}

func (m *TenantClientMetrics) TfState() *is.TFState {
	return m.tfstate
}

func (m *TenantClientMetrics) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Tenants
}

func (m *TenantClientMetrics) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	tenantId := m.tfstate.Int64("tenant_id")
	return rest.Tenants.GetClientMetricsWithContext(ctx, tenantId)
}

func (m *TenantClientMetrics) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	return ensureTenantClientMetricsUpdatedWith(ctx, ts, ts, rest)
}

func (m *TenantClientMetrics) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	stateTs := m.tfstate
	planTs := plan.(*TenantClientMetrics).TfState()
	return ensureTenantClientMetricsUpdatedWith(ctx, stateTs, planTs, rest)
}

func (m *TenantClientMetrics) DeleteResource(ctx context.Context, rest *VMSRest) error {
	// No-op for tenant client metrics, as it cannot be deleted.
	return nil
}

// ensureTenantClientMetricsUpdatedWith verifies if the given TenantClientMetrics (looked up by state)
// needs to be updated with new fields and performs the update if necessary.
//
// This is used in both CreateResource and UpdateResource for TenantClientMetrics.
func ensureTenantClientMetricsUpdatedWith(ctx context.Context, stateTs, fieldsTs *is.TFState, rest *VMSRest) (DisplayableRecord, error) {
	// Get tenant ID from tfstate
	tenantId := stateTs.Int64("tenant_id")
	if tenantId == 0 {
		return nil, fmt.Errorf("failed to get tenant ID: tenant ID is empty")
	}

	// Create params with the client metrics configuration
	params := params{}
	if ok := fieldsTs.SetToMapIfAvailable(
		params,
		"config",
		"user_defined_columns",
	); ok {
		// Use the custom API method to update client metrics
		return rest.Tenants.UpdateClientMetricsWithContext(ctx, tenantId, params)
	}

	// If no fields to update, just get the current client metrics
	return rest.Tenants.GetClientMetricsWithContext(ctx, tenantId)
}
