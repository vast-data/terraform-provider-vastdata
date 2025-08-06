// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"net/http"

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
	return m.updateClientMetrics(ctx, rest)
}

func (m *TenantClientMetrics) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	return m.updateClientMetrics(ctx, rest)
}

func (m *TenantClientMetrics) DeleteResource(ctx context.Context, rest *VMSRest) error {
	// No-op for tenant client metrics, as it cannot be deleted.
	return nil
}

// updateClientMetrics updates the client metrics settings for a tenant
func (m *TenantClientMetrics) updateClientMetrics(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	tenantId := m.tfstate.Int64("tenant_id")
	if tenantId == 0 {
		return nil, fmt.Errorf("failed to get tenant ID: tenant ID is empty")
	}

	// Create params with the client metrics configuration
	params := params{}
	if ok := m.tfstate.SetToMapIfAvailable(
		params,
		"config",
		"user_defined_columns",
	); ok {
		// Use the custom API method to update client metrics
		record, err := rest.Tenants.UpdateClientMetricsWithContext(ctx, tenantId, params)
		if err != nil {
			return nil, fmt.Errorf("failed to update client metrics: %w", err)
		}
		return record, nil
	}

	// If no fields to update, just get the current client metrics
	record, err := rest.Tenants.GetClientMetricsWithContext(ctx, tenantId)
	if err != nil {
		return nil, fmt.Errorf("failed to get client metrics: %w", err)
	}

	return record, nil
}
