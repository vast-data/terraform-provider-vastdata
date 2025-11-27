// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/schema_generation"
)

type Datasource struct {
	newManager   DatasourceFactoryFn
	providerData *ProviderData
	managerName  string
}

func (d *Datasource) EmptyManager() DataSourceManager {
	return d.newManager(nil, nil)
}

func (d *Datasource) ManagerWithSchemaOnly(ctx context.Context) (DataSourceManager, error) {
	// Get schema for the datasource to create a proper manager
	emptyManager := d.newManager(nil, nil)
	hints := emptyManager.TfState().Hints
	schema, err := schema_generation.GetDatasourceSchema(ctx, hints)
	if err != nil {
		return nil, err
	}
	// Create a new manager with the schema and empty Raw filled according to schema types
	// Build a zeroed attr map matching the schema so TFState has all keys with Null values
	zeroRaw := make(map[string]attr.Value)
	switch sch := any(*schema).(type) {
	case dschema.Schema:
		for k, a := range sch.Attributes {
			zeroRaw[k], _, _ = is.BuildAttrValueFromAny(a.GetType(), nil)
		}
	default:
		// Fallback to passing nil Raw if schema kind unexpected
	}
	return d.newManager(zeroRaw, *schema), nil
}

func (d *Datasource) NewManager(config tfsdk.Config) DataSourceManager {
	out, err := is.FillFrameworkValues(config.Raw, config.Schema)
	if err != nil {
		panic(fmt.Sprintf("error filling datasource: %s", err))
	}
	return d.newManager(out, config.Schema)
}

// ----------------------------------------
//      RESOURCE INTERFACE IMPLEMENTATION
// ----------------------------------------

func (d *Datasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	withContext(ctx, OpMetadata, d.managerName, func(ctx context.Context) {
		d.metadataImpl(ctx, req, resp)
	})
}

func (d *Datasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	withContext(ctx, OpSchema, d.managerName, func(ctx context.Context) {
		d.schemaImpl(ctx, req, resp)
	})
}

func (d *Datasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	withContext(ctx, OpConfigure, d.managerName, func(ctx context.Context) {
		d.configureImpl(ctx, req, resp)
	})
}

func (d *Datasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	withContext(ctx, OpRead, d.managerName, func(ctx context.Context) {
		d.readImpl(ctx, req, resp)
	})
}

// ----------------------------------------

func (d *Datasource) metadataImpl(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, d.managerName)
}

func (d *Datasource) schemaImpl(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	manager, err := d.ManagerWithSchemaOnly(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error getting schema for %q datasource.", d.managerName),
			err.Error(),
		)
		return
	}

	resp.Schema = manager.TfState().Schema.(dschema.Schema)
}

func (d *Datasource) configureImpl(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Check if the provider data is provided
	if req.ProviderData == nil {
		return
	}

	// Extract the provider data
	providerData, ok := req.ProviderData.(*ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Datasource Configure Type",
			fmt.Sprintf("Expected *ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.providerData = providerData
}

func (d *Datasource) readImpl(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		manager     = d.NewManager(req.Config)
		rest        = d.providerData.Client
		managerName = d.managerName
		tfState     = manager.TfState()
		record      DisplayableRecord
		err         error
	)

	if imp, ok := manager.(PrepareReadDatasource); ok {
		tflog.Debug(ctx, fmt.Sprintf("PrepareReadDatasource[%s]: do.", managerName))
		if err = imp.PrepareReadDatasource(ctx, rest); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("PrepareReadDatasource[%q]", managerName),
				err.Error(),
			)
			return
		}
	}

	tflog.Debug(
		ctx,
		fmt.Sprintf("Read[%q] - state:\n%s\n", managerName, tfState.Pretty()),
	)

	if !d.checkNonEmptyFields(ctx, manager, &resp.Diagnostics) {
		return
	}

	if imp, ok := manager.(ReadDatasource); ok {
		tflog.Debug(ctx, fmt.Sprintf("ReadDatasource[%s]: do.", managerName))
		record, err = imp.ReadDatasource(ctx, rest)
	} else {
		// Delegate to the default read implementation
		tflog.Debug(ctx, fmt.Sprintf("Read[%s]: use default implementation.", managerName))
		record, err = d.getRecordBySearchParams(ctx, manager, "Read")
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Read[%s]: error reading datasource.", managerName),
			err.Error(),
		)
		return
	}

	if record != nil {
		if transformer, ok := manager.(TransformResponseRecord); ok {
			tflog.Debug(ctx, fmt.Sprintf("TransformResponseRecord[%s]: do.", managerName))
			record = transformer.TransformResponseRecord(record.(Record))
		}

		// Automatically populate _id fields from nested objects
		// e.g., extract local_provider_id from local_provider.id
		PopulateIDFieldsFromNestedObjects(ctx, tfState, record.(Record))

		if err = tfState.FillFromRecord(record.(Record)); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Read[%s]: error filling datasource.", managerName),
				err.Error(),
			)
			return
		}
	} else {
		tflog.Debug(ctx, fmt.Sprintf("Read[%s]: no record returned, skipping state update.", managerName))
	}

	if imp, ok := manager.(AfterReadDatasource); ok {
		tflog.Debug(ctx, fmt.Sprintf("AfterReadDatasource[%s]: do.", managerName))
		if err = imp.AfterReadDatasource(ctx, rest); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("AfterReadDatasource[%q]", managerName),
				err.Error(),
			)
			return
		}

	}

	if err = tfState.SetState(ctx, &resp.State); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Read[%s]: error setting state.", managerName),
			err.Error(),
		)
		return
	}
}

// ---------------------------------------
//      HELPER FUNCTIONS
// ---------------------------------------

func (d *Datasource) getRecordBySearchParams(ctx context.Context, manager DataSourceManager, op string) (DisplayableRecord, error) {
	var (
		rest        = d.providerData.Client
		managerName = d.managerName
		tfState     = manager.TfState()
		api         = manager.API(rest)
	)
	return getRecordBySearchParams(ctx, api, tfState, nil, managerName, op)

}

func (d *Datasource) checkNonEmptyFields(ctx context.Context, manager DataSourceManager, dg *diag.Diagnostics) bool {
	ok := true
	tfState := manager.TfState()
	searchParams := getSearchParams(ctx, tfState, nil)
	if len(searchParams) == 0 {
		ok = false
		availableParams := manager.TfState().GetFilteredValues(
			is.FilterOr,
			nil,
			is.SearchEmpty,
			is.SearchOptional,
			is.SearchRequired,
			is.SearchPrimitivesOnly,
		)
		keys := make([]string, 0, len(availableParams))
		for k := range availableParams {
			keys = append(keys, k)
		}
		visualization := is.BuildDataSourceAttributesString(tfState.Schema.(dschema.Schema).Attributes, false, 2)

		dg.AddAttributeError(
			path.Root(""),
			fmt.Sprintf("empty %q datasource.", d.managerName),
			fmt.Sprintf(
				"Schema:\n%s\nAt least one field must be set: %+v",
				visualization, keys,
			),
		)
	}
	return ok
}
