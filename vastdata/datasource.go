package vastdata

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"

	"github.com/vast-data/terraform-provider-vastdata/internal/mappings"
)

type Datasource struct {
	newState func() DataSourceState
	client   *VMSRest
}

func RegisterDataSourceFor(datasourceStateFactory func() DataSourceState) func() datasource.DataSource {
	return func() datasource.DataSource {
		return &Datasource{
			newState: datasourceStateFactory,
		}
	}
}

// Configure adds the provider configured client to the data source.
func (d *Datasource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*VMSRest)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource configure type",
			fmt.Sprintf("Expected *vast_client.VMSRest, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *Datasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Info(ctx, fmt.Sprintf(">> %s[%p] start", "Metadata", d))
	defer tflog.Info(ctx, fmt.Sprintf(">> %s[%p] end", "Metadata", d))
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, entryName(d.newState()))
}

func (d *Datasource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Info(ctx, fmt.Sprintf(">> %s[%p] start", "Schema", d))
	defer tflog.Info(ctx, fmt.Sprintf(">> %s[%p] end", "Schema", d))
	resp.Schema = mappings.GenerateSchemaFromStruct(d.newState(), mappings.SchemaForDataSource).(schema.Schema)
}

// Read refreshes the Terraform state with the latest data.
func (d *Datasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, fmt.Sprintf(">> %s[%p] start", "Delete", d))
	defer tflog.Info(ctx, fmt.Sprintf(">> %s[%p] end", "Delete", d))

	var (
		state  = d.newState()
		rest   = d.client
		record DisplayableRecord
		err    error
	)

	diags := req.Config.Get(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err = state.PrepareReadDatasource(ctx, rest, req); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error preparing %q datasource for reading.", entryName(state)),
			err.Error(),
		)
		return
	}

	if !d.checkNonEmptyFields(state, resp.Diagnostics) {
		return
	}

	record, err = state.ReadDatasource(ctx, rest, req)
	if record == nil && err == nil {
		// Delegate to the default read implementation
		tflog.Debug(ctx, fmt.Sprintf("ReadDatasource[%s]: use default implementation.", entryName(state)))
		requestParams := d.getAllNonEmptyParams(state)
		record, err = state.GetRestResource(rest).GetWithContext(ctx, requestParams)
	}

	if record, ok := record.(Record); !ok {
		// For now only single Record type is supported. Update if needed (RecordSet support).
		resp.Diagnostics.AddError(
			fmt.Sprintf("invalid %q datasource.", entryName(state)),
			fmt.Sprintf(
				"Expected a Record type, got: %T."+
					" Please report this issue to the provider developers.",
				record,
			),
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error reading %q datasource.", entryName(state)),
			err.Error(),
		)
		return
	}
	if err = mappings.FillFromRecord(record.(Record), state); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error filling  %q datasource.", entryName(state)),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

}

func entryName(entry VastResourceGetter) string {
	return strings.ToLower(mappings.GetType(entry).Name())
}

func (d *Datasource) getAllNonEmptyParams(state DataSourceState) params {
	return must(mappings.GetNotEmptyFields(state, mappings.SchemaForDataSource, mappings.SearchAny))
}

func (d *Datasource) checkNonEmptyFields(state DataSourceState, dg diag.Diagnostics) bool {
	ok := true
	if !must(mappings.HasAnyNotEmptyFields(state, mappings.SchemaForDataSource)) {
		ok = false
		availableFields := must(mappings.GetOptionalFieldNames(state, mappings.SchemaForDataSource))
		dg.AddAttributeError(
			path.Root(""),
			fmt.Sprintf("empty %q datasource.", entryName(state)),
			fmt.Sprintf(
				"At least one field must be set in the datasource."+
					" Please set at least one of the following fields: %+v.",
				availableFields,
			),
		)
	}
	return ok
}
