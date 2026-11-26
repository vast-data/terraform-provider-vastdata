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

var DeSpanLogsSchemaRef = schema_generation.NewTelemetrySchemaReference(
	"",
	"",
	http.MethodGet,
	"telemetries/span-logs",
)

// DeSpanLogs is the read-only data source for GET /serverless/telemetries/span-logs.
type DeSpanLogs struct {
	tfstate is.ResourceState
}

func (m *DeSpanLogs) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &DeSpanLogs{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: DeSpanLogsSchemaRef,
			// Backend query_span_logs raises 400 if trace_id is omitted, even though
			// the OpenAPI parameter ref does not mark it required.
			RequiredSchemaFields: []string{"trace_id"},
			PreserveOrderFields:  []string{"data"},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameDatasourceAttribute(),
			},
		},
	)}
}

func (m *DeSpanLogs) TfState() is.ResourceState { return m.tfstate }

func (m *DeSpanLogs) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.DataEngine.Telemetry
}

func (m *DeSpanLogs) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	query := deTelemetryQueryParams(m.tfstate, rest)
	return core.Request[core.Record](ctx, rest.DataEngine.Telemetry, http.MethodGet, "telemetries/span-logs", query, nil)
}

func (m *DeSpanLogs) TransformResponseRecord(response Record) Record {
	return ensureTenantNameInRecord(response, m.tfstate, nil)
}
