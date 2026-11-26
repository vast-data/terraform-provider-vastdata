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

var DeLogsSchemaRef = schema_generation.NewTelemetrySchemaReference(
	"",
	"",
	http.MethodGet,
	"telemetries/logs",
)

// DeLogs is the read-only data source for GET /serverless/telemetries/logs.
type DeLogs struct {
	tfstate is.ResourceState
}

func (m *DeLogs) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &DeLogs{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:              DeLogsSchemaRef,
			AllowEmptySearchParams: true,
			PreserveOrderFields:    []string{"data"},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameDatasourceAttribute(),
			},
		},
	)}
}

func (m *DeLogs) TfState() is.ResourceState { return m.tfstate }

func (m *DeLogs) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.DataEngine.Logs
}

func (m *DeLogs) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	query := deTelemetryQueryParams(m.tfstate, rest)
	return core.Request[core.Record](ctx, rest.DataEngine.Logs, http.MethodGet, "telemetries/logs", query, nil)
}

func (m *DeLogs) TransformResponseRecord(response Record) Record {
	return ensureTenantNameInRecord(response, m.tfstate, nil)
}
