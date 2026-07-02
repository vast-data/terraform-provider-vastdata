// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/vast-data/go-vast-client/core"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

// TenantMetricLabelValuesSchemaRef uses the per-tenant metric label value endpoints.
var TenantMetricLabelValuesSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"tenants/{tenant_id}/metric_label_values",
	http.MethodGet,
	"tenants/{tenant_id}/metric_label_values",
)

type TenantMetricLabelValues struct {
	tfstate *is.TFState
}

func (m *TenantMetricLabelValues) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &TenantMetricLabelValues{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: TenantMetricLabelValuesSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"tenant_id": rschema.Int64Attribute{
					Required:    true,
					Description: "The ID of the tenant that owns this metric label value.",
				},
			},
			CommonModifiersMapping: map[string]string{
				"tenant_id": ModifierForceNew,
				"label_id":  ModifierForceNew,
			},
		},
	)}
}

func (m *TenantMetricLabelValues) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &TenantMetricLabelValues{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: TenantMetricLabelValuesSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"tenant_id": dschema.Int64Attribute{
					Required:    true,
					Description: "The ID of the tenant that owns this metric label value.",
				},
			},
		},
	)}
}

func (m *TenantMetricLabelValues) TfState() *is.TFState {
	return m.tfstate
}

func (m *TenantMetricLabelValues) API(rest *VMSRest) VastResourceAPIWithContext {
	return nil
}

func tenantMetricLabelValuesCollectionPath(tenantID int64) string {
	return core.BuildResourcePathWithID("tenants", tenantID, "metric_label_values")
}

func tenantMetricLabelValueItemPath(tenantID, valueID int64) string {
	return core.BuildResourcePathWithID(tenantMetricLabelValuesCollectionPath(tenantID), valueID)
}

func (m *TenantMetricLabelValues) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	tenantID := ts.Int64("tenant_id")
	if tenantID == 0 {
		return nil, fmt.Errorf("tenant_id is required")
	}

	body := params{}
	ts.SetToMapIfAvailable(body, "label_id", "value")
	path := tenantMetricLabelValuesCollectionPath(tenantID)
	return core.Request[Record](ctx, rest.Tenants, http.MethodPost, path, nil, body)
}

func (m *TenantMetricLabelValues) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	tenantID := ts.Int64("tenant_id")
	if tenantID == 0 {
		return nil, fmt.Errorf("tenant_id is required")
	}
	id := ts.Int64("id")
	if id == 0 {
		return nil, fmt.Errorf("id is required for read")
	}
	path := tenantMetricLabelValueItemPath(tenantID, id)
	return core.Request[Record](ctx, rest.Tenants, http.MethodGet, path, nil, nil)
}

func (m *TenantMetricLabelValues) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	planTs := plan.(*TenantMetricLabelValues).TfState()
	tenantID := planTs.Int64("tenant_id")
	id := m.tfstate.Int64("id")

	body := params{}
	planTs.SetToMapIfAvailable(body, "value")
	path := tenantMetricLabelValueItemPath(tenantID, id)
	return core.Request[Record](ctx, rest.Tenants, http.MethodPatch, path, nil, body)
}

func (m *TenantMetricLabelValues) DeleteResource(ctx context.Context, rest *VMSRest) error {
	ts := m.tfstate
	tenantID := ts.Int64("tenant_id")
	id := ts.Int64("id")
	path := tenantMetricLabelValueItemPath(tenantID, id)
	_, err := core.Request[Record](ctx, rest.Tenants, http.MethodDelete, path, nil, nil)
	return err
}

func (m *TenantMetricLabelValues) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	tenantID := ts.Int64("tenant_id")
	if tenantID == 0 {
		return nil, fmt.Errorf("tenant_id is required")
	}

	path := tenantMetricLabelValuesCollectionPath(tenantID)
	records, err := core.Request[RecordSet](ctx, rest.Tenants, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list metric label values for tenant %d: %w", tenantID, err)
	}

	labelID := ts.Int64("label_id")
	for _, record := range records {
		if labelID != 0 {
			label, _ := record["label"].(map[string]any)
			if label == nil {
				continue
			}
			recLabelID, ok := labelIDAsInt64(label["id"])
			if ok && recLabelID == labelID {
				return record, nil
			}
		}
	}

	return nil, fmt.Errorf("no metric label value found for tenant %d with label_id %d", tenantID, labelID)
}

func labelIDAsInt64(v any) (int64, bool) {
	switch val := v.(type) {
	case int64:
		return val, true
	case float64:
		return int64(val), true
	case int:
		return int64(val), true
	}
	return 0, false
}
