// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	planmodifiers "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vast-data/go-vast-client/core"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

type TenantMetricLabelValuesBulk struct {
	tfstate *is.TFState
}

var tenantMetricLabelValuesBulkSchemaAttributes = map[string]any{
	"tenant_id": rschema.Int64Attribute{
		Required:    true,
		Description: "The ID of the tenant whose metric label values are managed.",
		PlanModifiers: []planmodifiers.Int64{
			int64planmodifier.RequiresReplace(),
		},
	},
	"values": rschema.MapAttribute{
		ElementType: types.StringType,
		Required:    true,
		Description: "Dictionary of metric label keys to values. Replaces all existing values for the tenant on create and update.",
	},
	"id": rschema.Int64Attribute{
		Computed:    true,
		Description: "Terraform identifier for this resource, equal to tenant_id.",
	},
}

var tenantMetricLabelValuesBulkDatasourceAttributes = map[string]any{
	"tenant_id": dschema.Int64Attribute{
		Required:    true,
		Description: "The ID of the tenant whose metric label values are read.",
	},
	"values": dschema.MapAttribute{
		ElementType: types.StringType,
		Computed:    true,
		Description: "Dictionary of metric label keys to values.",
	},
}

func (m *TenantMetricLabelValuesBulk) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &TenantMetricLabelValuesBulk{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			TFStateHintsForCustom: &is.TFStateHintsForCustom{
				Description:      "Bulk management of tenant metric label values. This resource replaces all metric label values for a tenant with the provided dictionary.",
				SchemaAttributes: tenantMetricLabelValuesBulkSchemaAttributes,
			},
		},
	)}
}

func (m *TenantMetricLabelValuesBulk) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &TenantMetricLabelValuesBulk{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			TFStateHintsForCustom: &is.TFStateHintsForCustom{
				Description:      "Read all tenant metric label values as a dictionary of label keys to values.",
				SchemaAttributes: tenantMetricLabelValuesBulkDatasourceAttributes,
			},
		},
	)}
}

func (m *TenantMetricLabelValuesBulk) TfState() *is.TFState {
	return m.tfstate
}

func (m *TenantMetricLabelValuesBulk) API(rest *VMSRest) VastResourceAPIWithContext {
	return nil
}

func tenantMetricLabelValuesBulkPath(tenantID int64) string {
	return core.BuildResourcePathWithID("tenants", tenantID, "metric_label_values", "bulk")
}

func (m *TenantMetricLabelValuesBulk) bulkGet(ctx context.Context, rest *VMSRest, tenantID int64) (Record, error) {
	path := tenantMetricLabelValuesBulkPath(tenantID)
	return core.Request[Record](ctx, rest.Tenants, http.MethodGet, path, nil, nil)
}

func (m *TenantMetricLabelValuesBulk) bulkPost(ctx context.Context, rest *VMSRest, tenantID int64, body params) (Record, error) {
	path := tenantMetricLabelValuesBulkPath(tenantID)
	return core.Request[Record](ctx, rest.Tenants, http.MethodPost, path, nil, body)
}

func recordToStringMap(r Record) map[string]any {
	if r == nil {
		return map[string]any{}
	}
	out := make(map[string]any, len(r))
	for k, v := range r {
		out[k] = fmt.Sprint(v)
	}
	return out
}

func (m *TenantMetricLabelValuesBulk) applyBulkResponse(tenantID int64, rec Record) {
	m.tfstate.Set("id", tenantID)
	m.tfstate.Set("values", recordToStringMap(rec))
}

func (m *TenantMetricLabelValuesBulk) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	tenantID := m.tfstate.Int64("tenant_id")
	if tenantID == 0 {
		return nil, fmt.Errorf("tenant_id is required")
	}

	body := params(m.tfstate.ToMap("values"))
	rec, err := m.bulkPost(ctx, rest, tenantID, body)
	if err != nil {
		return nil, err
	}
	m.applyBulkResponse(tenantID, rec)
	return nil, nil
}

func (m *TenantMetricLabelValuesBulk) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	tenantID := m.tfstate.Int64("tenant_id")
	if tenantID == 0 {
		return nil, fmt.Errorf("tenant_id is required")
	}

	rec, err := m.bulkGet(ctx, rest, tenantID)
	if err != nil {
		return nil, err
	}
	m.applyBulkResponse(tenantID, rec)
	return nil, nil
}

func (m *TenantMetricLabelValuesBulk) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	planTs := plan.(*TenantMetricLabelValuesBulk).TfState()
	tenantID := planTs.Int64("tenant_id")
	if tenantID == 0 {
		return nil, fmt.Errorf("tenant_id is required")
	}

	body := params(planTs.ToMap("values"))
	rec, err := m.bulkPost(ctx, rest, tenantID, body)
	if err != nil {
		return nil, err
	}
	m.applyBulkResponse(tenantID, rec)
	return nil, nil
}

func (m *TenantMetricLabelValuesBulk) DeleteResource(ctx context.Context, rest *VMSRest) error {
	tenantID := m.tfstate.Int64("tenant_id")
	if tenantID == 0 {
		return fmt.Errorf("tenant_id is required")
	}

	_, err := m.bulkPost(ctx, rest, tenantID, params{})
	return err
}

func (m *TenantMetricLabelValuesBulk) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	tenantID := m.tfstate.Int64("tenant_id")
	if tenantID == 0 {
		return nil, fmt.Errorf("tenant_id is required")
	}

	rec, err := m.bulkGet(ctx, rest, tenantID)
	if err != nil {
		return nil, err
	}
	m.tfstate.Set("values", recordToStringMap(rec))
	return nil, nil
}
