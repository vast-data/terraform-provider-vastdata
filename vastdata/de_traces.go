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

var DeTracesSchemaRef = schema_generation.NewTelemetrySchemaReference(
	"",
	"",
	http.MethodGet,
	"telemetries/traces",
)

// DeTraces is the read-only data source for GET /serverless/telemetries/traces.
type DeTraces struct {
	tfstate is.ResourceState
}

func (m *DeTraces) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &DeTraces{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:              DeTracesSchemaRef,
			AllowEmptySearchParams: true,
			PreserveOrderFields:    []string{"data"},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameDatasourceAttribute(),
			},
		},
	)}
}

func (m *DeTraces) TfState() is.ResourceState { return m.tfstate }

func (m *DeTraces) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.DataEngine.Traces
}

func (m *DeTraces) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	query := deTelemetryQueryParams(m.tfstate, rest)
	return core.Request[core.Record](ctx, rest.DataEngine.Traces, http.MethodGet, "telemetries/traces", query, nil)
}

func (m *DeTraces) TransformResponseRecord(response Record) Record {
	return ensureTenantNameInRecord(response, m.tfstate, nil)
}
