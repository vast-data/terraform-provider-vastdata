// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/schema_generation"
)

type TenantMetricLabelValuesBulk struct {
	tfstate *is.TFState
}

var tenantMetricLabelValuesBulkSchemaAttributes = map[string]any{
	"tenant_id": rschema.Int64Attribute{
		Required:    true,
		Description: "The ID of the tenant whose metric label values are managed.",
	},
	"values": rschema.MapAttribute{
		ElementType: types.StringType,
		Required:    true,
		Description: "Dictionary of metric label keys to values. Keys must already exist as vastdata_tenant_metric_labels.",
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
			CommonModifiersMapping: map[string]string{
				"tenant_id": schema_generation.ModifierForceNew,
			},
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

func normalizeBulkLabelValuesRecord(r Record) (map[string]any, error) {
	if r == nil {
		return map[string]any{}, nil
	}

	out := make(map[string]any, len(r))
	for key, value := range r {
		if strings.HasPrefix(key, "@") {
			continue // go-vast-client internal metadata (e.g. @resourceType)
		}
		str, err := metricLabelScalarString(key, value)
		if err != nil {
			return nil, err
		}
		out[key] = str
	}
	return out, nil
}

// metricLabelScalarString enforces the TenantMetricLabelValues contract:
// a flat dictionary of label keys to string values (additionalProperties: string).
func metricLabelScalarString(key string, value any) (string, error) {
	s, ok := value.(string)
	if !ok {
		return "", fmt.Errorf(
			"metric label %q: expected string value, got %T (%v)",
			key, value, value,
		)
	}
	return s, nil
}

func validateBulkMetricLabelKeys(ctx context.Context, rest *VMSRest, values map[string]any) error {
	if len(values) == 0 {
		return nil
	}

	labels, err := rest.Tenants.TenantMetricLabelsListWithContext_GET(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to list tenant metric labels: %w", err)
	}

	registered := make(map[string]struct{}, len(labels))
	for _, label := range labels {
		key, _ := label["key"].(string)
		if key != "" {
			registered[key] = struct{}{}
		}
	}

	var unknown []string
	for key := range values {
		if _, ok := registered[key]; !ok {
			unknown = append(unknown, key)
		}
	}
	if len(unknown) == 0 {
		return nil
	}
	sort.Strings(unknown)
	return fmt.Errorf(
		"values contains unregistered metric label keys %v: create vastdata_tenant_metric_labels resources for these keys first",
		unknown,
	)
}

func (m *TenantMetricLabelValuesBulk) applyBulkResponse(tenantID int64, rec Record) error {
	values, err := normalizeBulkLabelValuesRecord(rec)
	if err != nil {
		return fmt.Errorf("invalid bulk metric label values response: %w", err)
	}
	m.tfstate.Set("id", tenantID)
	m.tfstate.Set("values", values)
	return nil
}

func (m *TenantMetricLabelValuesBulk) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	return ensureTenantMetricLabelValuesBulkUpdatedWith(ctx, m, ts, ts, rest)
}

func (m *TenantMetricLabelValuesBulk) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	tenantID := m.tfstate.Int64("tenant_id")
	if tenantID == 0 {
		tenantID = m.tfstate.Int64("id")
	}
	if tenantID == 0 {
		return nil, fmt.Errorf("tenant_id is required")
	}

	rec, err := rest.Tenants.TenantBulkWithContext_GET(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if err := m.applyBulkResponse(tenantID, rec); err != nil {
		return nil, err
	}
	return nil, nil
}

func (m *TenantMetricLabelValuesBulk) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	stateTs := m.tfstate
	planTs := plan.(*TenantMetricLabelValuesBulk).TfState()
	return ensureTenantMetricLabelValuesBulkUpdatedWith(ctx, m, stateTs, planTs, rest)
}

func ensureTenantMetricLabelValuesBulkUpdatedWith(
	ctx context.Context,
	m *TenantMetricLabelValuesBulk,
	stateTs, fieldsTs *is.TFState,
	rest *VMSRest,
) (DisplayableRecord, error) {
	tenantID := stateTs.Int64("tenant_id")
	body := fieldsTs.ToMap("values")
	if err := validateBulkMetricLabelKeys(ctx, rest, body); err != nil {
		return nil, err
	}
	rec, err := rest.Tenants.TenantBulkWithContext_POST(ctx, tenantID, body)
	if err != nil {
		return nil, err
	}
	if err := m.applyBulkResponse(tenantID, rec); err != nil {
		return nil, err
	}
	return nil, nil
}

func (m *TenantMetricLabelValuesBulk) DeleteResource(ctx context.Context, rest *VMSRest) error {
	tenantID := m.tfstate.Int64("tenant_id")
	if tenantID == 0 {
		tenantID = m.tfstate.Int64("id")
	}

	_, err := rest.Tenants.TenantBulkWithContext_POST(ctx, tenantID, nil)
	return err
}

func (m *TenantMetricLabelValuesBulk) ImportResourceState(req resource.ImportStateRequest, ctx context.Context, rest *VMSRest) error {
	if err := parseImportId(req.ID, m.tfstate); err != nil {
		return err
	}
	tenantID := m.tfstate.Int64("id")
	if tenantID == 0 {
		return fmt.Errorf("import id must be a valid tenant id")
	}
	m.tfstate.Set("tenant_id", tenantID)
	if _, err := m.ReadResource(ctx, rest); err != nil {
		return err
	}
	return CustomImportOnly{}
}

func (m *TenantMetricLabelValuesBulk) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	tenantID := m.tfstate.Int64("tenant_id")

	rec, err := rest.Tenants.TenantBulkWithContext_GET(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	values, err := normalizeBulkLabelValuesRecord(rec)
	if err != nil {
		return nil, fmt.Errorf("invalid bulk metric label values response: %w", err)
	}
	m.tfstate.Set("values", values)
	return nil, nil
}
