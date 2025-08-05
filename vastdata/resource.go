// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/schema_generation"
	"net/http"
)

type Resource struct {
	newManager  ResourceFactoryFn
	client      *VMSRest
	managerName string
}

func (r *Resource) EmptyManager() ResourceManager {
	return r.newManager(nil, nil)
}

func (r *Resource) NewManager(state any) ResourceManager {
	var (
		out    map[string]attr.Value
		schema any
		err    error
	)
	switch v := state.(type) {
	case tfsdk.Plan:
		schema = v.Schema
		out, err = is.FillFrameworkValues(v.Raw, schema)
	case tfsdk.State:
		schema = v.Schema
		out, err = is.FillFrameworkValues(v.Raw, schema)
	case tfsdk.Config:
		schema = v.Schema
		out, err = is.FillFrameworkValues(v.Raw, schema)
	default:
		panic(fmt.Sprintf("unsupported type: %T", v))
	}
	if err != nil {
		panic(fmt.Sprintf("error filling resource: %s", err))
	}
	return r.newManager(out, schema)
}

// ----------------------------------------
//      RESOURCE INTERFACE IMPLEMENTATION
// ----------------------------------------

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	withContext(ctx, "Metadata", r.managerName, func(ctx context.Context) {
		r.metadataImpl(ctx, req, resp)
	})
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	withContext(ctx, "Schema", r.managerName, func(ctx context.Context) {
		r.schemaImpl(ctx, req, resp)
	})
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	withContext(ctx, "Configure", r.managerName, func(ctx context.Context) {
		r.configureImpl(ctx, req, resp)
	})
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	withContext(ctx, "ImportState", r.managerName, func(ctx context.Context) {
		r.importStateImpl(ctx, req, resp)
	})
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	withContext(ctx, "Create", r.managerName, func(ctx context.Context) {
		r.createImpl(ctx, req, resp)
	})
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	withContext(ctx, "Read", r.managerName, func(ctx context.Context) {
		r.readImpl(ctx, req, resp)
	})
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	withContext(ctx, "Update", r.managerName, func(ctx context.Context) {
		r.updateImpl(ctx, req, resp)
	})
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	withContext(ctx, "Delete", r.managerName, func(ctx context.Context) {
		r.deleteImpl(ctx, req, resp)
	})
}

func (r *Resource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	withContext(ctx, "ValidateConfig", r.managerName, func(ctx context.Context) {
		r.validateConfigImpl(ctx, req, resp)
	})

}

// ----------------------------------------

func (r *Resource) metadataImpl(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, r.managerName)
}

func (r *Resource) schemaImpl(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	manager := r.EmptyManager()
	hints := manager.TfState().Hints

	schema, err := schema_generation.GetResourceSchema(ctx, hints)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error fetching OpenAPI schema for %q resource.", r.managerName),
			err.Error(),
		)
		return
	}
	resp.Schema = *schema
}

func (r *Resource) configureImpl(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*VMSRest)
}

func (r *Resource) importStateImpl(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var (
		rest        = r.client
		manager     = r.EmptyManager()
		managerName = r.managerName
		tfState     = manager.TfState()
		err         error
	)

	if imp, ok := manager.(PrepareImportResourceState); ok {
		tflog.Debug(ctx, fmt.Sprintf("PrepareImportResourceState[%s]: do.", managerName))
		if err = imp.PrepareImportResourceState(ctx, rest); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("PrepareImportResourceState[%q]", managerName),
				err.Error(),
			)
			return
		}
	}

	tflog.Debug(
		ctx,
		fmt.Sprintf("ImportState[%q] - state:\n%s\n", managerName, tfState.Pretty()),
	)

	if imp, ok := manager.(ImportResourceState); ok {
		tflog.Debug(ctx, fmt.Sprintf("ImportResourceState[%s]: do.", managerName))
		_, err = imp.ImportResourceState(ctx, rest)
	} else {
		resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error importing %q resource.", managerName),
			err.Error(),
		)
		return
	}

	if imp, ok := manager.(AfterImportResourceState); ok {
		tflog.Debug(ctx, fmt.Sprintf("AfterImportResourceState[%s]: do.", managerName))
		if err = imp.AfterImportResourceState(ctx, rest); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("AfterImportResourceState[%q]", managerName),
				err.Error(),
			)
			return
		}
	}

}

func (r *Resource) createImpl(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var (
		rest              = r.client
		manager           = r.NewManager(req.Plan)
		api               = manager.API(rest)
		managerName       = r.managerName
		record            DisplayableRecord
		tfState           = manager.TfState()
		tsStateCopy       = tfState.Copy() // Original vlaues from plan.
		err               error
		transactionDelete = func() {
			if resp.Diagnostics.HasError() {
				if imp, ok := manager.(DeleteResource); ok {
					imp.DeleteResource(ctx, rest)
				} else {
					if err = r.deleteRecordBySearchParams(ctx, manager, "TransactionDelete"); err != nil {
						tflog.Warn(
							ctx,
							fmt.Sprintf("TransactionDelete[%s]:"+
								" error deleting resource after error: %s",
								managerName,
								err.Error(),
							),
						)
					}
				}
			}
		}
	)

	if !r.checkNonEmptyFields(ctx, manager, &resp.Diagnostics) {
		return
	}

	if imp, ok := manager.(PrepareCreateResource); ok {
		tflog.Debug(ctx, fmt.Sprintf("PrepareCreateResource[%s]: do.", managerName))
		if err = imp.PrepareCreateResource(ctx, rest); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("PrepareCreateResource[%q]", managerName),
				err.Error(),
			)
			return
		}
	}

	tflog.Debug(
		ctx,
		fmt.Sprintf("Create[%q] - plan:\n%s\n", managerName, tfState.Pretty()),
	)

	if imp, ok := manager.(CreateResource); ok {
		tflog.Debug(ctx, fmt.Sprintf("CreateResource[%s]: do.", managerName))
		record, err = imp.CreateResource(ctx, rest)
	} else {
		// Delegate to the default create implementation
		tflog.Debug(ctx, fmt.Sprintf("Create[%s]: use default implementation.", managerName))
		// Get all params required + optional for creation.
		createParams := tfState.GetCreateParams()
		if transformer, ok := manager.(TransformRequestBody); ok {
			tflog.Debug(ctx, fmt.Sprintf("TransformRequestBody[%s]: do.", managerName))
			createParams = transformer.TransformRequestBody(createParams)
		}

		record, err = r.getRecordBySearchParams(ctx, manager, nil, "Create")
		if err != nil {
			// Something not expected happened, we should not continue.
			if !isNotFoundErr(err) {
				resp.Diagnostics.AddError(
					fmt.Sprintf("error reading %q resource.", managerName),
					err.Error(),
				)
				return
			} else {
				tflog.Debug(
					ctx, fmt.Sprintf("Create[%s]: no existing resource found, proceeding to create.", managerName),
				)

				defer transactionDelete()
				if record, err = api.CreateWithContext(ctx, createParams); err == nil {
					r.checkIntegrity(ctx, record.(Record), createParams)
				}
			}
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Create[%s]: found existing resource - taking management.", managerName))
			if transformer, ok := manager.(TransformResponseRecord); ok {
				tflog.Debug(ctx, fmt.Sprintf("TransformResponseRecord[%s]: do.", managerName))
				record = transformer.TransformResponseRecord(record.(Record))
			}
			// !NOTE: default implementation works only for resources with 'id' field.
			// For other resources please implement CreateResource to avoid entering this branch.
			createParamsDiff := diffMap(createParams, record.(Record))
			if len(createParamsDiff) > 0 {
				id, exists := record.(Record)["id"]
				if !exists {
					panic(fmt.Sprintf("Create[%s]: record does not have 'id' field.", managerName))
				}
				// Send only difference between current record from vast and createParams.
				if record, err = api.UpdateWithContext(ctx, id, createParamsDiff); err == nil {
					r.checkIntegrity(ctx, record.(Record), createParamsDiff)
				}
			} else {
				tflog.Debug(ctx, fmt.Sprintf("Create[%s]: no changes detected.", managerName))
			}
		}
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error creating %q resource.", managerName),
			err.Error(),
		)
		return
	}

	if imp, ok := manager.(AfterCreateResource); ok {
		tflog.Debug(ctx, fmt.Sprintf("AfterCreateResource[%s]: do.", managerName))
		if err = imp.AfterCreateResource(ctx, rest, record.(Record)); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("AfterCreateResource[%q]", managerName),
				err.Error(),
			)
			return
		}
	} else if len(tfState.Hints.EditOnlyFields) > 0 {
		// Update fields on resource that cannot be set on creation. For instance "enabled" field for some resources.
		updateParams := tfState.GetReadEditOnlyParams()
		if len(updateParams) > 0 {
			tflog.Debug(ctx, fmt.Sprintf("Create[%s]: Update 'EditOnly' fields.", managerName))
			id, exists := record.(Record)["id"]
			if !exists {
				panic(fmt.Sprintf("Create[%s]: record does not have 'id' field.", managerName))
			}
			_, err = api.UpdateWithContext(ctx, id, updateParams)
			for k, v := range updateParams {
				record.(Record)[k] = v // Update record with new values.
			}
		}
	}

	if record != nil {
		if transformer, ok := manager.(TransformResponseRecord); ok {
			tflog.Debug(ctx, fmt.Sprintf("TransformResponseRecord[%s]: do.", managerName))
			record = transformer.TransformResponseRecord(record.(Record))
		}

		// In particular scenarios we might want to populate all internalstate in custom handler.
		// In this case we might want to return nil to avoid this population.
		if err = tfState.FillFromRecord(record.(Record)); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("error filling %q resource.", managerName),
				err.Error(),
			)
		}
	} else {
		tflog.Debug(
			ctx,
			fmt.Sprintf("CreateResource[%s]: no record returned, skipping internalstate update.",
				managerName,
			),
		)
	}

	// Need to align original plan state with the current state in case response returned inconsistent data.
	// IOW you set fieldA as optional to valueA but backend returned valueB.
	tsStateCopy.CopyNonEmptyFieldsTo(tfState)

	if err = tfState.SetState(ctx, &resp.State); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error setting internalstate for %q resource.", managerName),
			err.Error(),
		)
		return
	}
}

func (r *Resource) readImpl(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var (
		rest        = r.client
		manager     = r.NewManager(req.State)
		managerName = r.managerName
		tfState     = manager.TfState()
		record      DisplayableRecord
		err         error
	)

	if imp, ok := manager.(PrepareReadResource); ok {
		tflog.Debug(ctx, fmt.Sprintf("PrepareReadResource[%s]: do.", managerName))
		if err = imp.PrepareReadResource(ctx, rest); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("PrepareReadResource[%q]", managerName),
				err.Error(),
			)
			return
		}
	}

	tflog.Debug(
		ctx,
		fmt.Sprintf("Read[%q] - state:\n%s\n", managerName, tfState.Pretty()),
	)

	if !r.checkNonEmptyFields(ctx, manager, &resp.Diagnostics) {
		return
	}

	if imp, ok := manager.(ReadResource); ok {
		tflog.Debug(ctx, fmt.Sprintf("ReadResource[%s]: do.", managerName))
		record, err = imp.ReadResource(ctx, rest)
	} else {
		// Delegate to the default read implementation
		tflog.Debug(ctx, fmt.Sprintf("Read[%s]: use default implementation.", managerName))
		record, err = r.getRecordBySearchParams(ctx, manager, nil, "Read")
	}

	if err != nil {
		if isNotFoundErr(err) || expectStatusCodes(err, http.StatusNotFound) || errors.As(err, &ForceCleanState{}) {
			// Ignore not found errors.
			// The next terraform plan will recreate the resource.
			tflog.Warn(ctx,
				fmt.Sprintf("Read[%s]: no such resource. err: %s", managerName, err),
			)
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Read[%s]", managerName),
				err.Error(),
			)
			tflog.Error(ctx, fmt.Sprintf("Read[%s]: error reading resource: %s", managerName, err))
		}
		return
	}

	if record != nil {
		if transformer, ok := manager.(TransformResponseRecord); ok {
			tflog.Debug(ctx, fmt.Sprintf("TransformResponseRecord[%s]: do.", managerName))
			record = transformer.TransformResponseRecord(record.(Record))
		}

		// In particular scenarios we might want to populate all internalstate in custom handler.
		// In this case we might want to return nil to avoid this population.
		if err = tfState.FillFromRecord(record.(Record)); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Read[%s] error filling resource.", managerName),
				err.Error(),
			)
			return
		}
	} else {
		tflog.Debug(ctx, fmt.Sprintf("Read[%s]: no record found, skipping internalstate update.", managerName))
	}

	if imp, ok := manager.(AfterReadResource); ok {
		tflog.Debug(ctx, fmt.Sprintf("AfterReadResource[%s]: do.", managerName))
		if err = imp.AfterReadResource(ctx, rest); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("AfterReadResource[%q]", managerName),
				err.Error(),
			)
			return
		}
	}

	if err = tfState.SetState(ctx, &resp.State); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Read[%s]: error setting internalstate.", managerName),
			err.Error(),
		)
		return
	}

}

func (r *Resource) updateImpl(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		rest        = r.client
		stateManger = r.NewManager(req.State)
		planManager = r.NewManager(req.Plan) // Planned changes to the resource (only diff fields)
		tfState     = stateManger.TfState()
		planTfState = planManager.TfState()
		managerName = r.managerName
		api         = stateManger.API(rest)
		record      DisplayableRecord
		err         error
	)

	if imp, ok := stateManger.(PrepareUpdateResource); ok {
		tflog.Debug(ctx, fmt.Sprintf("PrepareUpdateResource[%s]: do.", managerName))
		if err = imp.PrepareUpdateResource(ctx, planManager.(PrepareUpdateResource), rest); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("PrepareUpdateResource[%q]", managerName),
				err.Error(),
			)
			return
		}
	}

	tflog.Debug(
		ctx,
		fmt.Sprintf("Update[%q] - plan:\n%s\n", managerName, planTfState.Pretty()),
	)

	if !r.checkNonEmptyFields(ctx, planManager, &resp.Diagnostics) {
		return
	}

	if imp, ok := stateManger.(UpdateResource); ok {
		tflog.Debug(ctx, fmt.Sprintf("UpdateResource[%s]: do.", managerName))
		record, err = imp.UpdateResource(ctx, planManager.(UpdateResource), rest)
	} else {
		// Delegate to the default update implementation
		tflog.Debug(ctx, fmt.Sprintf("Update[%s]: use default implementation.", managerName))
		record, err = r.getRecordBySearchParams(ctx, stateManger, planManager, "Update")
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Update[%s]", managerName),
				err.Error(),
			)
			return
		}
		// !NOTE: default implementation works only for resources with 'id' field.
		// For other resources please implement CreateResource to avoid entering this branch.
		id, exists := record.(Record)["id"]
		if !exists {
			panic(fmt.Sprintf("Update[%s]: record does not have 'id' field.", managerName))
		}
		updateParams := planTfState.DiffFields(tfState, is.FilterOr, nil, is.SearchOptional, is.SearchRequired)
		if transformer, ok := stateManger.(TransformRequestBody); ok {
			tflog.Debug(ctx, fmt.Sprintf("TransformRequestBody[%s]: do.", managerName))
			updateParams = transformer.TransformRequestBody(updateParams)
		}

		delete(updateParams, "id") // Remove ID from update parameters, as it should not be updated.
		if record, err = api.UpdateWithContext(ctx, id, updateParams); err == nil {
			r.checkIntegrity(ctx, record.(Record), updateParams)
		}
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Update[%s].", managerName),
			err.Error(),
		)
		return
	}

	if record != nil {
		if transformer, ok := stateManger.(TransformResponseRecord); ok {
			tflog.Debug(ctx, fmt.Sprintf("TransformResponseRecord[%s]: do.", managerName))
			record = transformer.TransformResponseRecord(record.(Record))
		}

		// In particular scenarios we might want to populate all internalstate in custom handler.
		// In this case we might want to return nil to avoid this population.
		if err = tfState.FillFromRecord(record.(Record)); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Update[%s]: error filling resource.", managerName),
				err.Error(),
			)
			return
		}
	} else {
		tflog.Debug(ctx, fmt.Sprintf("Update[%s]: no record returned, skipping internalstate update.", managerName))
	}

	if imp, ok := stateManger.(AfterUpdateResource); ok {
		tflog.Debug(ctx, fmt.Sprintf("AfterUpdateResource[%s]: do.", managerName))
		if err = imp.AfterUpdateResource(ctx, planManager.(AfterUpdateResource), rest, record.(Record)); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("AfterUpdateResource[%q]", managerName),
				err.Error(),
			)
			return
		}
	}

	// Copy changes from plan to state.
	planTfState.CopyNonEmptyFieldsTo(tfState)

	if err = tfState.SetState(ctx, &resp.State); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Update[%s]: error setting internalstate.", managerName),
			err.Error(),
		)
		return
	}
}

func (r *Resource) deleteImpl(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var (
		rest        = r.client
		manager     = r.NewManager(req.State)
		tfState     = manager.TfState()
		managerName = r.managerName
		err         error
	)

	if imp, ok := manager.(PrepareDeleteResource); ok {
		tflog.Debug(ctx, fmt.Sprintf("PrepareDeleteResource[%s]: do.", managerName))
		if err = imp.PrepareDeleteResource(ctx, rest); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("PrepareDeleteResource[%q]", managerName),
				err.Error(),
			)
			return
		}
	}

	tflog.Debug(
		ctx,
		fmt.Sprintf("Delete[%q] - state:\n%s\n", managerName, tfState.Pretty()),
	)

	if !r.checkNonEmptyFields(ctx, manager, &resp.Diagnostics) {
		return
	}

	if imp, ok := manager.(DeleteResource); ok {
		tflog.Debug(ctx, fmt.Sprintf("DeleteResource[%s]: do.", managerName))
		err = imp.DeleteResource(ctx, rest)
	} else {
		// Delegate to the default delete implementation
		tflog.Debug(ctx, fmt.Sprintf("Delete[%s]: use default implementation.", managerName))
		err = r.deleteRecordBySearchParams(ctx, manager, "Delete")
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Delete[%s]", managerName),
			err.Error(),
		)
		return
	}

	if imp, ok := manager.(AfterDeleteResource); ok {
		tflog.Debug(ctx, fmt.Sprintf("AfterDeleteResource[%s]: do.", managerName))
		if err = imp.AfterDeleteResource(ctx, rest); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("AfterDeleteResource[%q]", managerName),
				err.Error(),
			)
			return
		}
	}

}

func (r *Resource) validateConfigImpl(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	manager := r.NewManager(req.Config)
	if imp, ok := manager.(ValidateResourceConfig); ok {
		if err := imp.ValidateResourceConfig(ctx); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("ValidateResourceConfig[%q]", r.managerName),
				err.Error(),
			)
			return
		}
	}

}

// ---------------------------------------
//      HELPER FUNCTIONS
// ---------------------------------------

func (r *Resource) getRecordBySearchParams(ctx context.Context, manager, planManager ResourceManager, op string) (DisplayableRecord, error) {
	var (
		rest                    = r.client
		managerName             = r.managerName
		tfState                 = manager.TfState()
		planTfState *is.TFState = nil
		api                     = manager.API(rest)
	)
	if planManager != nil {
		planTfState = planManager.TfState()
	}

	return getRecordBySearchParams(ctx, api, tfState, planTfState, managerName, op)

}

func (r *Resource) deleteRecordBySearchParams(ctx context.Context, manager ResourceManager, op string) error {
	var (
		rest        = r.client
		managerName = r.managerName
		tfState     = manager.TfState()
		api         = manager.API(rest)
	)
	return deleteRecordBySearchParams(ctx, api, tfState, managerName, op)
}

func (r *Resource) checkNonEmptyFields(ctx context.Context, manager ResourceManager, dg *diag.Diagnostics) bool {
	ok := true
	tfState := manager.TfState()
	searchParams := getSearchParams(ctx, tfState, nil)
	if len(searchParams) == 0 {
		ok = false
		availableParams := tfState.GetFilteredValues(
			is.FilterOr,
			nil,
			is.SearchEmpty,
			is.SearchRequired,
			is.SearchSearchable,
		)
		keys := make([]string, 0, len(availableParams))
		for k := range availableParams {
			keys = append(keys, k)
		}
		visualization := is.BuildResourceAttributesString(tfState.Schema.(rschema.Schema).Attributes, false, 2)
		managerName := r.managerName
		if len(keys) > 0 {
			dg.AddAttributeError(
				path.Root(""),
				fmt.Sprintf("No searchable fields provided for %q resource.", managerName),
				fmt.Sprintf(
					"Schema:\n%s\nAt least one of the following fields must be set: %+v",
					visualization,
					keys,
				),
			)
		} else {
			dg.AddAttributeError(
				path.Root(""),
				fmt.Sprintf("No searchable fields provided for %q resource.", managerName),
				fmt.Sprintf(
					"Schema:\n%s\nNone of the fields are marked as required or searchable in the OpenAPI schema for %q.\n"+
						"This typically indicates a missing or incomplete schema definition.",
					visualization,
					managerName,
				),
			)
		}
	}
	return ok
}

func (r *Resource) checkIntegrity(
	ctx context.Context,
	record Record,
	expected map[string]any,
) {
	// Wait consistency of the response.
	// Mismatched fields are not error state.

	var mismatches []string

	for key, expectedVal := range expected {
		actualVal, ok := record[key]
		if !ok {
			continue // treat missing as consistent
		}

		expectedVal = normalizeNumber(expectedVal)
		actualVal = normalizeNumber(actualVal)

		diff, panicked := safeDeepEqual(expectedVal, actualVal)
		if panicked {
			tflog.Debug(ctx, fmt.Sprintf("skipping deep.Equal (unhashable) for key: %q", key))
			continue
		}
		if diff != nil {
			mismatches = append(mismatches, diff...)
		}
	}
	if len(mismatches) > 0 {
		// Most likely it is VMS bug or OpenAPI schema mismatch.
		//IOW field that should not be marked as available for request are actually marked. (RequestBody model)
		tflog.Warn(
			ctx,
			"Record integrity check failed.",
		)
		for _, mismatch := range mismatches {
			tflog.Warn(ctx, mismatch)
		}
	} else {
		tflog.Debug(
			ctx,
			"Record integrity check passed.",
		)
	}
}
