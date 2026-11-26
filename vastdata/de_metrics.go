// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/vast-data/go-vast-client/core"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/schema_generation"
)

var DeMetricsSchemaRef = schema_generation.NewTelemetrySchemaReference(
	"",
	"",
	http.MethodGet,
	"telemetries/metrics",
)

// DeMetrics is the read-only data source for GET /serverless/telemetries/metrics.
type DeMetrics struct {
	tfstate is.ResourceState
}

func (m *DeMetrics) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &DeMetrics{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: DeMetricsSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameDatasourceAttribute(),
			},
		},
	)}
}

func (m *DeMetrics) TfState() is.ResourceState { return m.tfstate }

func (m *DeMetrics) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.DataEngine.Metrics
}

func (m *DeMetrics) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	query := deTelemetryQueryParams(m.tfstate, rest)
	return core.Request[core.Record](ctx, rest.DataEngine.Metrics, http.MethodGet, "telemetries/metrics", query, nil)
}

func (m *DeMetrics) TransformResponseRecord(response Record) Record {
	return ensureTenantNameInRecord(response, m.tfstate, nil)
}
