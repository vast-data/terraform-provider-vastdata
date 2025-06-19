package vastdata

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/vast-data/terraform-provider-vastdata/internal/mappings"
)

type Resource struct {
	newState func() ResourceState
	client   *VMSRest
}

func RegisterResourceFor(resourceStateFactory func() ResourceState) func() resource.Resource {
	return func() resource.Resource {
		return &Resource{
			newState: resourceStateFactory,
		}
	}
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Info(ctx, fmt.Sprintf(">> %s[%p] start", "Metadata", r))
	defer tflog.Info(ctx, fmt.Sprintf(">> %s[%p] end", "Metadata", r))

	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, entryName(r.newState()))
}

func (r *Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Info(ctx, fmt.Sprintf(">> %s[%p] start", "Schema", r))
	defer tflog.Info(ctx, fmt.Sprintf(">> %s[%p] end", "Schema", r))

	resp.Schema = mappings.GenerateSchemaFromStruct(r.newState(), mappings.SchemaForResource).(schema.Schema)
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Info(ctx, fmt.Sprintf(">> %s[%p] start", "Configure", r))
	defer tflog.Info(ctx, fmt.Sprintf(">> %s[%p] end", "Configure", r))

	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*VMSRest)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource configure type",
			fmt.Sprintf("Expected *VMSRest, got: %T", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, fmt.Sprintf(">> %s[%p] start", "ImportState", r))
	defer tflog.Info(ctx, fmt.Sprintf(">> %s[%p] end", "ImportState", r))

	var (
		rest  = r.client
		state = r.newState()
	)

	if err := state.PrepareImportResourceState(ctx, rest, req); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error preparing import for %q resource.", entryName(state)),
			err.Error(),
		)
		return
	}

	record, err := state.ImportResourceState(ctx, rest, req)
	if record == nil && err == nil {
		resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error importing %q resource.", entryName(state)),
			err.Error(),
		)
		return
	}
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, fmt.Sprintf(">> %s[%p] start", "Create", r))
	defer tflog.Info(ctx, fmt.Sprintf(">> %s[%p] end", "Create", r))

	var (
		rest   = r.client
		state  = r.newState()
		record DisplayableRecord
		err    error
	)

	if diags := req.Plan.Get(ctx, state); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if err = state.PrepareCreateResource(ctx, rest, req); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error preparing %q resource for creation.", entryName(state)),
			err.Error(),
		)
		return
	}

	if !r.checkNonEmptyFields(state, resp.Diagnostics) {
		return
	}

	record, err = state.CreateResource(ctx, rest, req)
	if record == nil && err == nil {
		// Delegate to the default create implementation
		tflog.Debug(ctx, fmt.Sprintf("CreateResource[%s]: use default implementation.", entryName(state)))
		searchParams := r.getSearchAbleParams(state)
		updateParams := r.getAllNonEmptyParams(state)
		record, err = state.GetRestResource(rest).GetWithContext(ctx, searchParams)
		if err != nil {
			// Something not expected happened, we should not continue.
			if !isNotFoundErr(err) {
				resp.Diagnostics.AddError(
					fmt.Sprintf("error reading %q resource.", entryName(state)),
					err.Error(),
				)
				return
			} else {
				tflog.Debug(ctx, fmt.Sprintf("CreateResource[%s]: no existing resource found, proceeding to create.", entryName(state)))
				// Transactional error handling
				defer func() {
					if resp.Diagnostics.HasError() {
						tflog.Error(ctx, fmt.Sprintf("CreateResource[%s]: error occurred: %s", entryName(state), resp.Diagnostics.Errors()))
						deleteParams := r.getSearchAbleParams(state)
						if len(deleteParams) > 0 {
							state.GetRestResource(rest).DeleteWithContext(ctx, deleteParams, nil)
						}
					}
				}()

				record, err = state.GetRestResource(rest).CreateWithContext(ctx, updateParams)
			}
		} else {
			tflog.Debug(ctx, fmt.Sprintf("CreateResource[%s]: found existing resource - taking management.", entryName(state)))
			id := record.(Record).RecordID()
			record, err = state.GetRestResource(rest).UpdateWithContext(ctx, id, updateParams)
		}
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error creating %q resource.", entryName(state)),
			err.Error(),
		)
		return
	}
	if err = mappings.FillFromRecord(record.(Record), state); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error filling %q resource.", entryName(state)),
			err.Error(),
		)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, fmt.Sprintf(">> %s[%p] start", "Read", r))
	defer tflog.Info(ctx, fmt.Sprintf(">> %s[%p] end", "Read", r))

	var (
		rest   = r.client
		state  = r.newState()
		record DisplayableRecord
		err    error
	)

	if diags := req.State.Get(ctx, state); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if err = state.PrepareReadResource(ctx, rest, req); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error preparing %q resource for reading.", entryName(state)),
			err.Error(),
		)
		return
	}

	if !r.checkNonEmptyFields(state, resp.Diagnostics) {
		return
	}

	record, err = state.ReadResource(ctx, rest, req)
	if record == nil && err == nil {
		// Delegate to the default read implementation
		tflog.Debug(ctx, fmt.Sprintf("ReadResource[%s]: use default implementation.", entryName(state)))
		id := must(mappings.GetIdPtr(state))
		if id != nil {
			tflog.Debug(ctx, fmt.Sprintf("ReadResource[%s]: found id %v.", entryName(state), *id))
			// If ID is provided, read by ID.
			record, err = state.GetRestResource(rest).GetByIdWithContext(ctx, *id)
		} else {
			tflog.Debug(ctx, fmt.Sprintf("ReadResource[%s]: no id found, using search params.", entryName(state)))
			// If no ID is provided, use the search parameters to find the resource.
			// This is useful for resources that do not have a unique ID or when the ID is not known.
			requestParams := r.getSearchAbleParams(state)
			record, err = state.GetRestResource(rest).GetWithContext(ctx, requestParams)
		}
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error reading %q resource.", entryName(state)),
			err.Error(),
		)
		return

	}
	if record == nil {
		tflog.Warn(ctx,
			fmt.Sprintf("ReadResource[%s]: no such resource.", entryName(state)),
		)
		resp.State.RemoveResource(ctx)
		return
	}
	if err = mappings.FillFromRecord(record.(Record), state); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error filling %q resource.", entryName(state)),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, fmt.Sprintf(">> %s[%p] start", "Update", r))
	defer tflog.Info(ctx, fmt.Sprintf(">> %s[%p] end", "Update", r))

	var (
		rest   = r.client
		state  = r.newState()
		record DisplayableRecord
		err    error
	)

	if diags := req.Plan.Get(ctx, state); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if err = state.PrepareUpdateResource(ctx, rest, req); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error preparing %q resource for update.", entryName(state)),
			err.Error(),
		)
		return
	}

	if !r.checkNonEmptyFields(state, resp.Diagnostics) {
		return
	}

	record, err = state.UpdateResource(ctx, rest, req)
	if record == nil && err == nil {
		// Delegate to the default update implementation
		tflog.Debug(ctx, fmt.Sprintf("UpdateResource[%s]: use default implementation.", entryName(state)))
		searchParams := r.getSearchAbleParams(state)
		record, err = state.GetRestResource(rest).GetWithContext(ctx, searchParams)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("error reading %q resource before update.", entryName(state)),
				err.Error(),
			)
			return
		}
		id := record.(Record).RecordID()
		updateParams := r.getAllNonEmptyParams(state)
		record, err = state.GetRestResource(rest).UpdateWithContext(ctx, id, updateParams)
	}

	if record == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("unexpected error: record is still nil for etry %q.", entryName(state)),
			"Please report this issue to the provider developers.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error updating %q resource.", entryName(state)),
			err.Error(),
		)
		return
	}

	if err = mappings.FillFromRecord(record.(Record), state); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error filling %q resource.", entryName(state)),
			err.Error(),
		)
		return

	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, fmt.Sprintf(">> %s[%p] start", "Delete", r))
	defer tflog.Info(ctx, fmt.Sprintf(">> %s[%p] end", "Delete", r))

	var (
		rest  = r.client
		state = r.newState()
	)

	if diags := req.State.Get(ctx, state); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if err := state.PrepareDeleteResource(ctx, rest, req); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error preparing %q resource for deletion.", entryName(state)),
			err.Error(),
		)
		return
	}

	if !r.checkNonEmptyFields(state, resp.Diagnostics) {
		return
	}

	record, err := state.DeleteResource(ctx, rest, req)
	if record == nil && err == nil {
		// Delegate to the default delete implementation
		tflog.Debug(ctx, fmt.Sprintf("DeleteResource[%s]: use default implementation.", entryName(state)))
		id := must(mappings.GetIdPtr(state))
		if id != nil {
			tflog.Debug(ctx, fmt.Sprintf("DeleteResource[%s]: found id %v.", entryName(state), id))
			// If ID is provided, delete by ID.
			record, err = state.GetRestResource(rest).DeleteByIdWithContext(ctx, *id, nil)
		} else {
			tflog.Debug(ctx, fmt.Sprintf("DeleteResource[%s]: no id found, using search params.", entryName(state)))
			// If no ID is provided, use the search parameters to find the resource.
			// This is useful for resources that do not have a unique ID or when the ID is not known.
			deleteParams := r.getSearchAbleParams(state)
			if len(deleteParams) == 0 {
				resp.Diagnostics.AddError(
					fmt.Sprintf("error deleting %q resource.", entryName(state)),
					"At least one field must be set in the resource to delete.",
				)
				return
			}
			record, err = state.GetRestResource(rest).DeleteWithContext(ctx, deleteParams, nil)
		}
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error deleting %q resource.", entryName(state)),
			err.Error(),
		)
		return
	}
}

// getSearchAbleParams extracts common search parameters from the resource state.
// The typical pattern is to search by "name", or by the pair ("name", "tenant_id")
// if both are marked as required.
//
// If a different search strategy is needed, it should be implemented directly in:
//   - ImportResourceState
//   - CreateResource
//   - ReadResource
//   - UpdateResource
//   - DeleteResource
func (r *Resource) getSearchAbleParams(state ResourceState) params {
	defaultSearchableParams := must(mappings.GetNotEmptyFields(state, mappings.SchemaForResource, mappings.SearchSearchable))
	if len(defaultSearchableParams) > 0 {
		// "searchable" is present in "resource" tags. These are default for search.
		return defaultSearchableParams
	}

	// "searchable" is not present in "resource" tags.
	var searchParams = make(params)

	allRequiredParams := must(mappings.GetNotEmptyFields(state, mappings.SchemaForResource, mappings.SearchRequited))
	if name, ok := allRequiredParams["name"]; ok {
		// If name is required, we can use it as a search parameter.
		searchParams["name"] = name
	}
	if tenantId, ok := allRequiredParams["tenant_id"]; ok {
		// If tenant_id is required, we can use it as a search parameter.
		searchParams["tenant_id"] = tenantId
	}

	return searchParams
}

func (r *Resource) getAllNonEmptyParams(state ResourceState) params {
	return must(mappings.GetNotEmptyFields(state, mappings.SchemaForResource, mappings.SearchAny))
}

func (r *Resource) checkNonEmptyFields(state ResourceState, dg diag.Diagnostics) bool {
	ok := true
	if !must(mappings.HasAnyNotEmptyFields(state, mappings.SchemaForResource)) {
		ok = false
		availableFields := must(mappings.GetOptionalFieldNames(state, mappings.SchemaForResource))
		dg.AddAttributeError(
			path.Root(""),
			fmt.Sprintf("empty %q resource.", entryName(state)),
			fmt.Sprintf(
				"At least one field must be set in the resource."+
					" Please set at least one of the following fields: %+v.",
				availableFields,
			),
		)
	}
	return ok
}
