// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var TenantMetricLabelsSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"tenants/metric_labels",
	http.MethodGet,
	"tenants/metric_labels",
)

type TenantMetricLabels struct {
	tfstate *is.TFState
}

func (m *TenantMetricLabels) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &TenantMetricLabels{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: TenantMetricLabelsSchemaRef,
			CommonModifiersMapping: map[string]string{
				"key":           ModifierForceNew,
				"default_value": ModifierForceNew,
				"description":   ModifierForceNew,
			},
		},
	)}
}

func (m *TenantMetricLabels) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &TenantMetricLabels{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: TenantMetricLabelsSchemaRef,
		},
	)}
}

func (m *TenantMetricLabels) TfState() *is.TFState {
	return m.tfstate
}

func (m *TenantMetricLabels) API(rest *VMSRest) VastResourceAPIWithContext {
	return nil
}

// CreateResource calls POST /tenants/metric_labels/ directly.
// The default CRUD cannot be used because API() returns rest.Tenants which
// routes to /tenants/, not /tenants/metric_labels/.
func (m *TenantMetricLabels) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	body := params{}
	ts.SetToMapIfAvailable(body, "key", "default_value", "description")
	return rest.Tenants.TenantMetricLabelsWithContext_POST(ctx, body)
}

// ReadResource fetches by ID via GET /tenants/metric_labels/{id}/.
// The /tenants/metric_labels/ list endpoint supports no attribute filters,
// so lookup must always be done by ID.
func (m *TenantMetricLabels) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.readByID(ctx, rest)
}

// ReadDatasource looks up by ID when known, otherwise scans the full list by key.
func (m *TenantMetricLabels) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	if m.tfstate.IsKnownAndNotNull("id") {
		return m.readByID(ctx, rest)
	}
	return m.readByKey(ctx, rest)
}

func (m *TenantMetricLabels) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	return m.readByID(ctx, rest)
}

// DeleteResource calls DELETE /tenants/metric_labels/{id}/.
func (m *TenantMetricLabels) DeleteResource(ctx context.Context, rest *VMSRest) error {
	id := m.tfstate.Int64("id")
	if id == 0 {
		return fmt.Errorf("'id' is required to delete a metric label")
	}
	return rest.Tenants.TenantMetricLabelsWithContext_DELETE(ctx, id)
}

func (m *TenantMetricLabels) readByID(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	id := m.tfstate.Int64("id")
	if id == 0 {
		return nil, fmt.Errorf("'id' is required to read a metric label")
	}
	return rest.Tenants.TenantMetricLabelsByIdWithContext_GET(ctx, id)
}

// readByKey lists all metric labels and returns the one matching 'key'.
// Used only by the data source when 'id' is not provided.
func (m *TenantMetricLabels) readByKey(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	key := m.tfstate.String("key")
	if key == "" {
		return nil, fmt.Errorf("either 'id' or 'key' must be specified to look up a metric label")
	}
	records, err := rest.Tenants.TenantMetricLabelsListWithContext_GET(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list metric labels: %w", err)
	}
	for _, record := range records {
		if k, ok := record["key"].(string); ok && k == key {
			return record, nil
		}
	}
	return nil, fmt.Errorf("metric label with key %q not found", key)
}
