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

var DeTraceTreeSchemaRef = schema_generation.NewTelemetrySchemaReference(
	"",
	"",
	http.MethodGet,
	"telemetries/trace-tree",
)

// DeTraceTree is the read-only data source for GET /serverless/telemetries/trace-tree.
type DeTraceTree struct {
	tfstate is.ResourceState
}

func (m *DeTraceTree) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &DeTraceTree{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: DeTraceTreeSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameDatasourceAttribute(),
			},
		},
	)}
}

func (m *DeTraceTree) TfState() is.ResourceState { return m.tfstate }

func (m *DeTraceTree) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.DataEngine.Telemetry
}

func (m *DeTraceTree) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	query := deTelemetryQueryParams(m.tfstate, rest)
	return core.Request[core.Record](ctx, rest.DataEngine.Telemetry, http.MethodGet, "telemetries/trace-tree", query, nil)
}

func (m *DeTraceTree) TransformResponseRecord(response Record) Record {
	return ensureTenantNameInRecord(response, m.tfstate, nil)
}
