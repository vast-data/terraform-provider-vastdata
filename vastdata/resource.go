// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/schema_generation"
)

// Operation constants for resource operations
const (
	OpMetadata       = "Metadata"
	OpSchema         = "Schema"
	OpConfigure      = "Configure"
	OpImportState    = "ImportState"
	OpCreate         = "Create"
	OpRead           = "Read"
	OpUpdate         = "Update"
	OpDelete         = "Delete"
	OpValidateConfig = "ValidateConfig"
)

type Resource struct {
	newManager   ResourceFactoryFn
	providerData *ProviderData
	managerName  string
}

func (r *Resource) EmptyManager() ResourceManager {
	return r.newManager(nil, nil)
}

func (r *Resource) ManagerWithSchemaOnly(ctx context.Context) (ResourceManager, error) {
	// Get schema for the resource to create a proper manager
	emptyManager := r.newManager(nil, nil)
	hints := emptyManager.TfState().Hints
	schema, err := schema_generation.GetResourceSchema(ctx, hints)
	if err != nil {
		return nil, err
	}
	// Create a new manager with the schema and empty Raw filled according to schema types
	// Build a zeroed attr map matching the schema so TFState has all keys with Null values
	zeroRaw := make(map[string]attr.Value)
	switch sch := any(*schema).(type) {
	case rschema.Schema:
		for k, a := range sch.Attributes {
			zeroRaw[k], _, _ = is.BuildAttrValueFromAny(a.GetType(), nil)
		}
	default:
		// Fallback to passing nil Raw if schema kind unexpected
	}
	return r.newManager(zeroRaw, *schema), nil
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
	withContext(ctx, OpMetadata, r.managerName, func(ctx context.Context) {
		r.metadataImpl(ctx, req, resp)
	})
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	withContext(ctx, OpSchema, r.managerName, func(ctx context.Context) {
		r.schemaImpl(ctx, req, resp)
	})
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	withContext(ctx, OpConfigure, r.managerName, func(ctx context.Context) {
		r.configureImpl(ctx, req, resp)
	})
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	withContext(ctx, OpImportState, r.managerName, func(ctx context.Context) {
		r.importStateImpl(ctx, req, resp)
	})
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	withContext(ctx, OpCreate, r.managerName, func(ctx context.Context) {
		r.createImpl(ctx, req, resp)
	})
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	withContext(ctx, OpRead, r.managerName, func(ctx context.Context) {
		r.readImpl(ctx, req, resp)
	})
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	withContext(ctx, OpUpdate, r.managerName, func(ctx context.Context) {
		r.updateImpl(ctx, req, resp)
	})
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	withContext(ctx, OpDelete, r.managerName, func(ctx context.Context) {
		r.deleteImpl(ctx, req, resp)
	})
}

func (r *Resource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	withContext(ctx, OpValidateConfig, r.managerName, func(ctx context.Context) {
		r.validateConfigImpl(ctx, req, resp)
	})
}

// ----------------------------------------

func (r *Resource) metadataImpl(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, r.managerName)
}

func (r *Resource) schemaImpl(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	manager, err := r.ManagerWithSchemaOnly(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error fetching OpenAPI schema for %q resource.", r.managerName),
			err.Error(),
		)
		return
	}

	resp.Schema = manager.TfState().Schema.(rschema.Schema)
}

func (r *Resource) configureImpl(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Check if the provider data is provided
	if req.ProviderData == nil {
		return
	}

	// Extract the provider data
	providerData, ok := req.ProviderData.(*ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.providerData = providerData
}

func (r *Resource) importStateImpl(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var (
		rest        = r.providerData.Client
		manager, _  = r.ManagerWithSchemaOnly(ctx)
		managerName = r.managerName
		tfState     = manager.TfState()
		hints       = tfState.Hints
		err         error
	)

	// Check importable flag (defaults to true)
	if hints != nil && hints.Importable != nil && !*hints.Importable {
		resp.Diagnostics.AddError(
			fmt.Sprintf("ImportState[%q]: import not supported.", managerName),
			"This resource has been marked as not importable.",
		)
		return
	}

	importID := req.ID
	if strings.TrimSpace(importID) == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("ImportState[%s]: missing import ID.", managerName),
			fmt.Sprintf("An import ID or key=value list is required for importing the %q resource.", managerName),
		)
		return
	}

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
		if err := imp.ImportResourceState(req, ctx, rest); err != nil {
			if errors.As(err, &CustomImportOnly{}) {
				tflog.Debug(
					ctx,
					fmt.Sprintf("ImportResourceState[%s]: custom import only - stopping.", managerName),
				)
				// Finally, persist state
				if err = tfState.SetState(ctx, &resp.State); err != nil {
					resp.Diagnostics.AddError(
						fmt.Sprintf("ImportState[%s]: failed to set state.", managerName),
						fmt.Sprintf("Failed to set the imported state: %s", err.Error()),
					)
				}
				return

			}

			resp.Diagnostics.AddError(
				fmt.Sprintf("error importing %q resource.", managerName),
				err.Error(),
			)
			return
		}

	} else {
		// Use default import implementation
		tflog.Debug(ctx, fmt.Sprintf("ImportState[%s]: use default import implementation.", managerName))
		if err := parseImportId(importID, tfState); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("ImportState[%s]: invalid import ID.", managerName),
				err.Error(),
			)
			return
		}
	}

	// Before reading, allow resource to prepare read (same as in readImpl)
	if prep, ok := manager.(PrepareReadResource); ok {
		tflog.Debug(ctx, fmt.Sprintf("PrepareReadResource[%s]: do.", managerName))
		if err = prep.PrepareReadResource(ctx, rest); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("PrepareReadResource[%q]", managerName),
				err.Error(),
			)
			return
		}
	}

	// Populate full state from backend using the same logic as read (without SDK response wrapper)
	var record DisplayableRecord
	if reader, ok := manager.(ReadResource); ok {
		tflog.Debug(ctx, fmt.Sprintf("ImportState[%s]: reading resource after key parsing.", managerName))
		record, err = reader.ReadResource(ctx, rest)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("ImportState[%s]: failed to read imported resource.", managerName),
				err.Error(),
			)
			return
		}
		if record == nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("ImportState[%s]: resource import is not supported for this type.", managerName),
				"ReadResource returned nil",
			)
			return
		}
	} else {
		tflog.Debug(ctx, fmt.Sprintf("ImportState[%s]: default read after key parsing.", managerName))
		record, err = r.getRecordBySearchParams(ctx, manager, nil, "Import")
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("ImportState[%s]: failed to read imported resource.", managerName),
				err.Error(),
			)
			return
		}
	}

	if record != nil {
		if transformer, ok := manager.(TransformResponseRecord); ok {
			tflog.Debug(ctx, fmt.Sprintf("TransformResponseRecord[%s]: do.", managerName))
			record = transformer.TransformResponseRecord(record.(Record))
		}

		// Automatically populate _id fields from nested objects
		// e.g., extract local_provider_id from local_provider.id
		PopulateIDFieldsFromNestedObjects(ctx, tfState, record.(Record))

		// On import, populate computed and required fields (but NOT optional fields)
		if err = tfState.FillFromRecordForImport(record.(Record)); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("ImportState[%s]: error filling state.", managerName),
				err.Error(),
			)
			return
		}
	}

	// After read hook (same as in readImpl)
	if aft, ok := manager.(AfterReadResource); ok {
		tflog.Debug(ctx, fmt.Sprintf("AfterReadResource[%s]: do.", managerName))
		if err = aft.AfterReadResource(ctx, rest); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("AfterReadResource[%q]", managerName),
				err.Error(),
			)
			return
		}
	}

	// Finally, persist state
	if err = tfState.SetState(ctx, &resp.State); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("ImportState[%s]: failed to set state.", managerName),
			fmt.Sprintf("Failed to set the imported state: %s", err.Error()),
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
		rest              = r.providerData.Client
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
					if _, err = r.deleteRecordBySearchParams(ctx, manager, "TransactionDelete"); err != nil {
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
					panic(fmt.Sprintf("Create[%s]: record does not have 'id' field. Record: %s", managerName, record.(Record).PrettyJson("  ")))
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

	// Ensure record is not empty, otherwise set to nil.
	if record, ok := record.(Record); ok && record.Empty() {
		record = nil
	}

	if record != nil {
		// In case record is AsyncTask
		if err := handleMaybeAsyncTask(ctx, rest, record.(Record), 10*time.Minute); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("AsyncTask - create[%s].", managerName),
				err.Error(),
			)
			return
		}

		// Handle AfterCreateResource hook
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
					panic(fmt.Sprintf("Create[%s]: record does not have 'id' field. Record: %s", managerName, record.(Record).PrettyJson("  ")))
				}
				_, err = api.UpdateWithContext(ctx, id, updateParams)
				for k, v := range updateParams {
					record.(Record)[k] = v // Update record with new values.
				}
			}
		}

		// Handle response transformation
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
		rest        = r.providerData.Client
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
		// Check if we should skip API calls and use current tfstate
		if tfState.Hints.SkipRefreshAPICall {
			tflog.Debug(ctx, fmt.Sprintf("Read[%s]: skip API call, using current tfstate as record.", managerName))
			// Convert current tfstate to map[string]any and use as record
			record = Record(tfState.GetAllValues())
		} else {
			// Delegate to the default read implementation
			tflog.Debug(ctx, fmt.Sprintf("Read[%s]: use default implementation.", managerName))
			record, err = r.getRecordBySearchParams(ctx, manager, nil, "Read")
		}
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

		// Automatically populate _id fields from nested objects
		// e.g., extract local_provider_id from local_provider.id
		PopulateIDFieldsFromNestedObjects(ctx, tfState, record.(Record))

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
		rest        = r.providerData.Client
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
			panic(fmt.Sprintf("Update[%s]: record does not have 'id' field. Record: %s", managerName, record.(Record).PrettyJson("  ")))
		}

		// Get update params excluding edit-only and delete-only fields
		updateParams := planTfState.GetUpdateParams(tfState)
		if transformer, ok := stateManger.(TransformRequestBody); ok {
			tflog.Debug(ctx, fmt.Sprintf("TransformRequestBody[%s]: do.", managerName))
			updateParams = transformer.TransformRequestBody(updateParams)
		}

		delete(updateParams, "id") // Remove ID from update parameters, as it should not be updated.
		if len(updateParams) == 0 {
			tflog.Debug(
				ctx, fmt.Sprintf("Update[%s]: no update request parameters. Skipping update...", managerName),
			)
		} else {
			if record, err = api.UpdateWithContext(ctx, id, updateParams); err == nil {
				r.checkIntegrity(ctx, record.(Record), updateParams)
			}
		}
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Update[%s].", managerName),
			err.Error(),
		)
		return
	}

	// Ensure record is not empty, otherwise set to nil.
	if record, ok := record.(Record); ok && record.Empty() {
		record = nil
	}

	if record != nil {
		// In case record is AsyncTask
		if err := handleMaybeAsyncTask(ctx, rest, record.(Record), 10*time.Minute); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("AsyncTask - update[%s].", managerName),
				err.Error(),
			)
			return
		}

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

		if imp, ok := stateManger.(AfterUpdateResource); ok {
			tflog.Debug(ctx, fmt.Sprintf("AfterUpdateResource[%s]: do.", managerName))
			if err = imp.AfterUpdateResource(ctx, planManager.(AfterUpdateResource), rest, record.(Record)); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("AfterUpdateResource[%q]", managerName),
					err.Error(),
				)
				return
			}
		} else if len(tfState.Hints.EditOnlyFields) > 0 {
			// Only handle edit-only fields automatically if AfterUpdateResource is not implemented
			// If AfterUpdateResource is implemented, it's assumed to handle edit-only fields itself
			editOnlyParams := planTfState.GetChangedEditOnlyParams(tfState)
			if len(editOnlyParams) > 0 {
				tflog.Debug(ctx, fmt.Sprintf("Update[%s]: Updating 'EditOnly' fields.", managerName))
				id, exists := record.(Record)["id"]
				if !exists {
					resp.Diagnostics.AddError(
						fmt.Sprintf("Update[%s]: cannot update edit-only fields.", managerName),
						"Record does not have 'id' field",
					)
					return
				}
				var editRecord DisplayableRecord
				if editRecord, err = api.UpdateWithContext(ctx, id, editOnlyParams); err != nil {
					resp.Diagnostics.AddError(
						fmt.Sprintf("Update[%s]: error updating edit-only fields.", managerName),
						err.Error(),
					)
					return
				}
				// Update the record with edit-only values
				if editRecord != nil && !editRecord.(Record).Empty() {
					for k, v := range editRecord.(Record) {
						record.(Record)[k] = v
					}
				} else {
					for k, v := range editOnlyParams {
						record.(Record)[k] = v
					}
				}
			}
		}

	} else {
		tflog.Debug(ctx, fmt.Sprintf("Update[%s]: no record returned, skipping internalstate update.", managerName))
	}

	// Copy changes from plan to state.
	planTfState.CopyKnownFieldsTo(tfState)

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
		rest        = r.providerData.Client
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
		record, err := r.deleteRecordBySearchParams(ctx, manager, "Delete")
		if err == nil && record != nil {
			// In case record is AsyncTask
			if err := handleMaybeAsyncTask(ctx, rest, record, 10*time.Minute); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("AsyncTask - delete[%s].", managerName),
					err.Error(),
				)
				return
			}
		}
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

// parseAndApplyCompositeImport parses a composite import string into values for the given fields
// and invokes set(key, value) with properly typed Terraform attr values.
// Accepted formats:
//   - key=value pairs separated by ',' (e.g., "gid=1001,tenant_id=22,context=ad")
//   - ordered values separated by '|' matching fields order
func parseAndApplyCompositeImport(importID string, fields []string, tfState *is.TFState, set func(string, attr.Value)) error {
	kv := make(map[string]string)
	s := strings.TrimSpace(importID)
	if s == "" {
		return fmt.Errorf("empty import id")
	}

	if strings.Contains(s, "=") {
		// Handle key=value comma-separated format
		sep := ","
		if strings.Contains(s, ";") {
			sep = ";"
		}
		parts := strings.Split(s, sep)
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			kvPair := strings.SplitN(p, "=", 2)
			if len(kvPair) != 2 {
				return fmt.Errorf("segment %q is not in key=value form", p)
			}
			key := strings.TrimSpace(kvPair[0])
			val := strings.TrimSpace(kvPair[1])
			kv[key] = val
		}
	} else if strings.Contains(s, "|") {
		// Handle ordered pipe-separated values
		values := strings.Split(s, "|")
		if len(values) != len(fields) {
			return fmt.Errorf("expected %d values for fields %v, got %d", len(fields), fields, len(values))
		}
		for i, f := range fields {
			kv[f] = strings.TrimSpace(values[i])
		}
	} else {
		// Handle single token - treat as "id" field
		if tfState.HasAttribute("id") {
			kv["id"] = s
		} else {
			return fmt.Errorf("single token provided but no 'id' field available in schema")
		}
	}

	// Set parsed fields with type coercion and schema validation
	for f, val := range kv {
		if !tfState.HasAttribute(f) {
			return fmt.Errorf("field %q is not present in the schema", f)
		}
		t := tfState.Type(f)
		switch {
		case t.Equal(types.Int64Type):
			n, perr := strconv.ParseInt(val, 10, 64)
			if perr != nil {
				return fmt.Errorf("invalid int64 for %q: %v", f, perr)
			}
			set(f, types.Int64Value(n))
		case t.Equal(types.StringType):
			set(f, types.StringValue(val))
		case t.Equal(types.BoolType):
			bv := strings.EqualFold(val, "true") || val == "1"
			set(f, types.BoolValue(bv))
		default:
			// fallback: attempt passthrough as string
			set(f, types.StringValue(val))
		}
	}
	return nil
}

func (r *Resource) getRecordBySearchParams(ctx context.Context, manager, planManager ResourceManager, op string) (DisplayableRecord, error) {
	var (
		rest                    = r.providerData.Client
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

func (r *Resource) deleteRecordBySearchParams(ctx context.Context, manager ResourceManager, op string) (Record, error) {
	var (
		rest        = r.providerData.Client
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
