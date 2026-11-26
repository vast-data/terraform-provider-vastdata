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

var DeMetricNamesSchemaRef = schema_generation.NewTelemetrySchemaReference(
	"",
	"",
	http.MethodGet,
	"telemetries/metrics/names",
)

// DeMetricNames is the read-only data source for GET /serverless/telemetries/metrics/names.
type DeMetricNames struct {
	tfstate is.ResourceState
}

func (m *DeMetricNames) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &DeMetricNames{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:              DeMetricNamesSchemaRef,
			AllowEmptySearchParams: true,
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameDatasourceAttribute(),
			},
		},
	)}
}

func (m *DeMetricNames) TfState() is.ResourceState { return m.tfstate }

func (m *DeMetricNames) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.DataEngine.Metrics
}

func (m *DeMetricNames) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	query := deTelemetryQueryParams(m.tfstate, rest)
	return core.Request[core.Record](ctx, rest.DataEngine.Metrics, http.MethodGet, "telemetries/metrics/names", query, nil)
}

func (m *DeMetricNames) TransformResponseRecord(response Record) Record {
	return ensureTenantNameInRecord(response, m.tfstate, nil)
}
