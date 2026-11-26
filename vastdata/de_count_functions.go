// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"maps"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/vast-data/go-vast-client/core"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/schema_generation"
)

var DeCountFunctionsSchemaRef = schema_generation.NewTelemetrySchemaReference(
	"",
	"",
	http.MethodGet,
	"telemetries/count-functions",
)

// DeCountFunctions is the read-only data source for GET /serverless/telemetries/count-functions.
type DeCountFunctions struct {
	tfstate is.ResourceState
}

func (m *DeCountFunctions) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &DeCountFunctions{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: DeCountFunctionsSchemaRef,
			// Terraform reserves the root attribute name "count"
			//for meta-arguments so here we rename attribute.
			ExcludedSchemaFields: []string{"count"},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameDatasourceAttribute(),
				"function_count": dschema.Int64Attribute{
					Computed:    true,
					Description: "The number of functions in the database that started between start_time and end_time.",
				},
			},
		},
	)}
}

func (m *DeCountFunctions) TfState() is.ResourceState { return m.tfstate }

func (m *DeCountFunctions) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.DataEngine.Telemetry
}

func (m *DeCountFunctions) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	query := deTelemetryQueryParams(m.tfstate, rest)
	return core.Request[core.Record](ctx, rest.DataEngine.Telemetry, http.MethodGet, "telemetries/count-functions", query, nil)
}

func (m *DeCountFunctions) TransformResponseRecord(response Record) Record {
	if response != nil {
		if count, ok := response["count"]; ok {
			out := Record{}
			maps.Copy(out, response)
			out["function_count"] = count
			response = out
		}
	}
	return ensureTenantNameInRecord(response, m.tfstate, nil)
}
