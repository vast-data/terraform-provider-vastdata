// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/logging"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/schema_generation"
)

// Operation constants for resource operations
const (
	OpMetadata       = "Metadata"
	OpSchema         = "Schema"
	OpConfigure      = "Configure"
	OpImportState    = "ImportState"
	OpUpgradeState   = "UpgradeState"
	OpCreate         = "Create"
	OpRead           = "Read"
	OpUpdate         = "Update"
	OpDelete         = "Delete"
	OpValidateConfig = "ValidateConfig"
	OpModifyPlan     = "ModifyPlan"
)

type Resource struct {
	newManager   ResourceFactoryFn
	providerData *ProviderData
	managerName  string
	// lifecycle drives the CRUD operations of this resource kind. It is built in Configure.
	// The framework creates a new Resource and calls Configure for every RPC, so it lives for
	// one RPC, not for the provider process: add callbacks right after it is built.
	lifecycle *ResourceLifecycle
}

func (r *Resource) EmptyManager() ResourceManager {
	return r.newManager(nil, nil)
}

func (r *Resource) ManagerWithSchemaOnly(ctx context.Context) (ResourceManager, error) {
	return newManagerWithSchemaOnly(ctx, r.newManager)
}

// newManagerWithSchemaOnly builds a manager whose state carries the generated resource
// schema and a zeroed (all-null) attribute map matching it. It does not depend on any
// Terraform request/response types, so Terraform and Crossplane share it.
func newManagerWithSchemaOnly(ctx context.Context, newManager ResourceFactoryFn) (ResourceManager, error) {
	// Get schema for the resource to create a proper manager
	emptyManager := newManager(nil, nil)
	hints := emptyManager.TfState().StateHints()
	schema, err := schema_generation.GetResourceSchema(ctx, hints)
	if err != nil {
		return nil, err
	}
	if u, ok := emptyManager.(UpgradeResourceSchema); ok {
		if v := u.SchemaVersion(); v > schema.Version {
			schema.Version = v
		}
	}
	// Create a new manager with the schema and empty Raw filled according to schema types
	// Build a zeroed attr map matching the schema so the state has all keys with Null values
	zeroRaw := make(map[string]attr.Value)
	switch sch := any(*schema).(type) {
	case rschema.Schema:
		for k, a := range sch.Attributes {
			zeroRaw[k], _, _ = is.BuildAttrValueFromAny(a.GetType(), nil)
		}
	default:
		// Fallback to passing nil Raw if schema kind unexpected
	}
	return newManager(zeroRaw, *schema), nil
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

// overlayWriteOnlyFromConfig copies write-only attribute values from Config into
// manager's tfstate. Write-only values are always null in Plan/State and are only
// available on Config (Terraform >= 1.11). Without this, create/update would omit
// secrets such as certificate PEMs from the API body (TERF-268).
func (r *Resource) overlayWriteOnlyFromConfig(manager ResourceManager, config tfsdk.Config) {
	ts := manager.TfState()
	if ts == nil || !ts.IsEnabled() || ts.StateHints() == nil || len(ts.StateHints().WriteOnlyFields) == 0 {
		return
	}
	cfgTs := r.NewManager(config).TfState()
	for _, key := range ts.StateHints().WriteOnlyFields {
		v, ok := cfgTs.RawValues()[key]
		if !ok || v == nil || v.IsNull() || v.IsUnknown() {
			continue
		}
		ts.RawValues()[key] = v
	}
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
		if r.providerData != nil && r.providerData.MigrateMode {
			r.createMigrateImpl(ctx, req, resp)
		} else {
			r.createImpl(ctx, req, resp)
		}
	})
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	withContext(ctx, OpRead, r.managerName, func(ctx context.Context) {
		r.readImpl(ctx, req, resp)
	})
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	withContext(ctx, OpUpdate, r.managerName, func(ctx context.Context) {
		if r.providerData != nil && r.providerData.MigrateMode {
			// Non-importable resources (e.g. vastdata_user_key) are preserved
			// from the old state verbatim.  Their schema may differ from the new
			// provider, causing a schema-reconciliation update.  Resources that
			// implement MigrateModePassThroughUpdate opt in to allowing that
			// update instead of being blocked.
			if _, ok := r.EmptyManager().(MigrateModePassThroughUpdate); !ok {
				resp.Diagnostics.AddError(
					fmt.Sprintf("MigrateMode[%s]: update operation blocked.", r.managerName),
					"VASTDATA_MIGRATE_MODE is intended for populating an EMPTY state from existing infrastructure. "+
						"An update was triggered, which means the state already contains this resource. "+
						"Please start with an empty state (remove terraform.tfstate) and run again, "+
						"or unset VASTDATA_MIGRATE_MODE to operate normally.",
				)
				return
			}
			logging.Warn(ctx, fmt.Sprintf("MigrateMode[%s]: allowing schema-reconciliation update for non-importable resource.", r.managerName))
		}
		r.updateImpl(ctx, req, resp)
	})
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	withContext(ctx, OpDelete, r.managerName, func(ctx context.Context) {
		if r.providerData != nil && r.providerData.MigrateMode {
			resp.Diagnostics.AddError(
				fmt.Sprintf("MigrateMode[%s]: delete operation blocked.", r.managerName),
				"VASTDATA_MIGRATE_MODE is intended for populating an EMPTY state from existing infrastructure. "+
					"A delete was triggered, which means the state already contains this resource. "+
					"Please start with an empty state (remove terraform.tfstate) and run again, "+
					"or unset VASTDATA_MIGRATE_MODE to operate normally.",
			)
			return
		}
		r.deleteImpl(ctx, req, resp)
	})
}

func (r *Resource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	withContext(ctx, OpValidateConfig, r.managerName, func(ctx context.Context) {
		r.validateConfigImpl(ctx, req, resp)
	})
}

func (r *Resource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	withContext(ctx, OpModifyPlan, r.managerName, func(ctx context.Context) {
		r.modifyPlanImpl(ctx, req, resp)
	})
}

// UpgradeState registers schema state upgraders when the resource uses WriteOnly
// / *_md5 hints and/or implements UpgradeResourceSchema. Prior versions
// 0..Version-1 each get an upgrader that decodes into the current schema, runs
// shared WriteOnly migration when applicable, then ResourceLifecycle.UpgradeResourceSchema
// (manager hook + external callbacks).
func (r *Resource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	var out map[int64]resource.StateUpgrader
	withContext(ctx, OpUpgradeState, r.managerName, func(ctx context.Context) {
		out = r.upgradeStateImpl(ctx)
	})
	return out
}

func (r *Resource) modifyPlanImpl(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}
	manager := r.NewManager(req.Plan)
	imp, ok := manager.(ModifyPlanResource)
	if !ok {
		return
	}
	logging.Debug(ctx, fmt.Sprintf("ModifyPlanResource[%s]: do.", r.managerName))
	if err := imp.ModifyPlan(ctx, req, resp); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("ModifyPlanResource[%q]", r.managerName),
			err.Error(),
		)
	}
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

	resp.Schema = manager.TfState().SchemaDefinition().(rschema.Schema)
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
	r.lifecycle = NewResourceLifecycle(providerData.Client, r.managerName)
}

func (r *Resource) importStateImpl(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	manager, _ := r.ManagerWithSchemaOnly(ctx)
	managerName := r.managerName
	tfState := manager.TfState()
	hints := tfState.StateHints()

	// Check importable flag (defaults to true) before touching provider/client.
	if hints != nil && hints.Importable != nil && !*hints.Importable {
		resp.Diagnostics.AddError(
			fmt.Sprintf("ImportState[%q]: import not supported.", managerName),
			"This resource has been marked as not importable.",
		)
		return
	}

	var (
		rest = r.providerData.Client
		err  error
	)

	importID := req.ID
	if strings.TrimSpace(importID) == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("ImportState[%s]: missing import ID.", managerName),
			fmt.Sprintf("An import ID or key=value list is required for importing the %q resource.", managerName),
		)
		return
	}

	logging.Debug(ctx, fmt.Sprintf("PrepareImportResourceState[%s]: do.", managerName))
	if err = r.lifecycle.runPrepareImportResourceState(ctx, manager); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("PrepareImportResourceState[%q]", managerName),
			err.Error(),
		)
		return
	}

	logging.Debug(
		ctx,
		fmt.Sprintf("ImportState[%q] - state:\n%s\n", managerName, tfState.Pretty()),
	)

	if customImport, err := r.lifecycle.runImportResourceState(ctx, req, manager); customImport {
		logging.Debug(ctx, fmt.Sprintf("ImportResourceState[%s]: do.", managerName))
		if err != nil {
			if errors.As(err, &CustomImportOnly{}) {
				logging.Debug(
					ctx,
					fmt.Sprintf("ImportResourceState[%s]: custom import only - stopping.", managerName),
				)
				// Finally, persist state
				if err = is.NewTFSdkSink(&resp.State).WriteState(ctx, tfState); err != nil {
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
		logging.Debug(ctx, fmt.Sprintf("ImportState[%s]: use default import implementation.", managerName))
		if err := parseImportId(importID, tfState); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("ImportState[%s]: invalid import ID.", managerName),
				err.Error(),
			)
			return
		}
	}

	// Before reading, allow resource to prepare read (same as in readImpl)
	logging.Debug(ctx, fmt.Sprintf("PrepareReadResource[%s]: do.", managerName))
	if err = r.lifecycle.runPrepareReadResource(ctx, manager); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("PrepareReadResource[%q]", managerName),
			err.Error(),
		)
		return
	}

	// Populate full state from backend using the same logic as read (without SDK response wrapper)
	var (
		record     DisplayableRecord
		customRead bool
	)
	record, customRead, err = r.lifecycle.runReadResource(ctx, manager)
	if customRead {
		logging.Debug(ctx, fmt.Sprintf("ImportState[%s]: read resource after key parsing with ReadResource.", managerName))
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
		logging.Debug(ctx, fmt.Sprintf("ImportState[%s]: default read after key parsing.", managerName))
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
		record = runTransformResponseRecord(manager, record.(Record))

		// Automatically populate _id fields from nested objects
		// e.g., extract local_provider_id from local_provider.id
		PopulateIDFieldsFromNestedObjects(ctx, tfState, record.(Record))

		if err = mergeSubResources(ctx, manager, rest, record, managerName); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("GetSubResources[%q]", managerName), err.Error())
			return
		}

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
	logging.Debug(ctx, fmt.Sprintf("AfterReadResource[%s]: do.", managerName))
	if err = r.lifecycle.runAfterReadResource(ctx, manager); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("AfterReadResource[%q]", managerName),
			err.Error(),
		)
		return
	}

	// Finally, persist state
	if err = is.NewTFSdkSink(&resp.State).WriteState(ctx, tfState); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("ImportState[%s]: failed to set state.", managerName),
			fmt.Sprintf("Failed to set the imported state: %s", err.Error()),
		)
		return
	}

	logging.Debug(ctx, fmt.Sprintf("AfterImportResourceState[%s]: do.", managerName))
	if err = r.lifecycle.runAfterImportResourceState(ctx, manager); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("AfterImportResourceState[%q]", managerName),
			err.Error(),
		)
		return
	}

}

func (r *Resource) createImpl(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	manager := r.NewManager(req.Plan)
	// Write-only attrs are null in Plan; load them from Config for the API request.
	r.overlayWriteOnlyFromConfig(manager, req.Config)

	if err := r.lifecycle.Create(ctx, manager); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error creating %q resource.", r.managerName),
			err.Error(),
		)
		return
	}
	if err := is.NewTFSdkSink(&resp.State).WriteState(ctx, manager.TfState()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error setting internalstate for %q resource.", r.managerName),
			err.Error(),
		)
	}
}

func (r *Resource) createMigrateImpl(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var (
		manager     = r.NewManager(req.Plan)
		managerName = r.managerName
		tfState     = manager.TfState()
		record      DisplayableRecord
		err         error
	)

	logging.Warn(ctx, fmt.Sprintf("MigrateMode[%s]: reading existing resource instead of creating.", managerName))

	// Prepare read hooks (same as in readImpl/importStateImpl)
	logging.Debug(ctx, fmt.Sprintf("MigrateMode/PrepareReadResource[%s]: do.", managerName))
	if err = r.lifecycle.runPrepareReadResource(ctx, manager); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("MigrateMode/PrepareReadResource[%q]", managerName),
			err.Error(),
		)
		return
	}

	// Read the resource using the same logic as import
	var customRead bool
	record, customRead, err = r.lifecycle.runReadResource(ctx, manager)
	if !customRead {
		logging.Debug(ctx, fmt.Sprintf("MigrateMode[%s]: using default read by search params.", managerName))
		record, err = r.getRecordBySearchParams(ctx, manager, nil, "MigrateCreate")
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("MigrateMode[%s]: failed to read existing resource from cluster.", managerName),
			fmt.Sprintf("Ensure the resource exists on the VAST cluster. Error: %s", err.Error()),
		)
		return
	}

	if record != nil {
		// Transform response if needed
		record = runTransformResponseRecord(manager, record.(Record))

		// Populate _id fields from nested objects (same as import)
		PopulateIDFieldsFromNestedObjects(ctx, tfState, record.(Record))

		// Fill state from the record (using import-style fill which populates computed + required)
		if err = tfState.FillFromRecordForImport(record.(Record)); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("MigrateMode[%s]: error filling state from cluster record.", managerName),
				err.Error(),
			)
			return
		}
	} else {
		resp.Diagnostics.AddError(
			fmt.Sprintf("MigrateMode[%s]: resource not found on cluster.", managerName),
			"The resource was not found on the VAST cluster. Ensure the resource configuration matches an existing resource.",
		)
		return
	}

	// After read hook
	logging.Debug(ctx, fmt.Sprintf("MigrateMode/AfterReadResource[%s]: do.", managerName))
	if err = r.lifecycle.runAfterReadResource(ctx, manager); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("MigrateMode/AfterReadResource[%q]", managerName),
			err.Error(),
		)
		return
	}

	// Persist state
	if err = is.NewTFSdkSink(&resp.State).WriteState(ctx, tfState); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("MigrateMode[%s]: failed to set state.", managerName),
			err.Error(),
		)
		return
	}

	logging.Warn(ctx, fmt.Sprintf("MigrateMode[%s]: successfully populated state from existing cluster resource.", managerName))
}

func (r *Resource) readImpl(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	manager := r.NewManager(req.State)
	exists, err := r.lifecycle.Read(ctx, manager)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Read[%s]", r.managerName), err.Error())
		logging.Error(ctx, fmt.Sprintf("Read[%s]: error reading resource: %s", r.managerName, err))
		return
	}
	if !exists {
		logging.Warn(ctx, fmt.Sprintf("Read[%s]: no such resource.", r.managerName))
		resp.State.RemoveResource(ctx)
		return
	}
	if err := is.NewTFSdkSink(&resp.State).WriteState(ctx, manager.TfState()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Read[%s]: error setting internalstate.", r.managerName),
			err.Error(),
		)
	}
}

func (r *Resource) updateImpl(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	current := r.NewManager(req.State)
	desired := r.NewManager(req.Plan)
	// Write-only attrs are null in Plan; load them from Config for the API request.
	r.overlayWriteOnlyFromConfig(desired, req.Config)

	if err := r.lifecycle.Update(ctx, current, desired); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Update[%s].", r.managerName), err.Error())
		return
	}
	if err := is.NewTFSdkSink(&resp.State).WriteState(ctx, current.TfState()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Update[%s]: error setting internalstate.", r.managerName),
			err.Error(),
		)
	}
}

func (r *Resource) deleteImpl(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	manager := r.NewManager(req.State)
	if err := r.lifecycle.Delete(ctx, manager); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Delete[%s]", r.managerName), err.Error())
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

func (r *Resource) upgradeStateImpl(ctx context.Context) map[int64]resource.StateUpgrader {
	empty := r.EmptyManager()
	hints := empty.TfState().StateHints()
	writeOnly := hints != nil && (len(hints.WriteOnlyFields) > 0 || len(hints.WriteOnlyMD5Fields) > 0)

	version := int64(0)
	if writeOnly {
		version = 1
	}
	if u, ok := empty.(UpgradeResourceSchema); ok {
		if v := u.SchemaVersion(); v > version {
			version = v
		}
	}
	if version == 0 {
		return nil
	}

	out := make(map[int64]resource.StateUpgrader, version)
	for from := int64(0); from < version; from++ {
		fromVersion := from
		out[fromVersion] = resource.StateUpgrader{
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				withContext(ctx, OpUpgradeState, r.managerName, func(ctx context.Context) {
					r.upgradeStateFromImpl(ctx, fromVersion, writeOnly, req, resp)
				})
			},
		}
	}
	return out
}

func (r *Resource) upgradeStateFromImpl(
	ctx context.Context,
	fromVersion int64,
	writeOnly bool,
	req resource.UpgradeStateRequest,
	resp *resource.UpgradeStateResponse,
) {
	if req.RawState == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("UpgradeState[%s]", r.managerName),
			"missing prior RawState",
		)
		return
	}

	manager, err := r.ManagerWithSchemaOnly(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("UpgradeState[%s]: schema", r.managerName),
			err.Error(),
		)
		return
	}

	sch := manager.TfState().SchemaDefinition().(rschema.Schema)
	hints := manager.TfState().StateHints()

	rawVal, err := req.RawState.UnmarshalWithOpts(sch.Type().TerraformType(ctx), tfprotov6.UnmarshalOpts{
		ValueFromJSONOpts: tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true},
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("UpgradeState[%s]: unmarshal prior state", r.managerName),
			err.Error(),
		)
		return
	}

	attrMap, err := is.FillFrameworkValues(rawVal, sch)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("UpgradeState[%s]: decode prior state", r.managerName),
			err.Error(),
		)
		return
	}

	// Concrete *TFState: SeedWriteOnlyMD5FromSecrets is upgrade-only and not
	// part of ResourceState (managers/lifecycle never need it).
	ts := is.NewTFStateMust(attrMap, sch, hints)
	if writeOnly {
		ts.SeedWriteOnlyMD5FromSecrets()
		ts.ClearWriteOnlyFields()
	}

	lc := r.lifecycle
	if lc == nil {
		// UpgradeState can run before Configure in some hosts/tests; still run
		// the manager hook via a throwaway lifecycle (no external callbacks).
		lc = NewResourceLifecycle(nil, r.managerName)
	}
	if err := lc.UpgradeResourceSchema(ctx, fromVersion, manager, ts); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("UpgradeState[%s]", r.managerName),
			err.Error(),
		)
		return
	}

	objType, ok := sch.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("UpgradeState[%s]", r.managerName),
			fmt.Sprintf("expected object schema type, got %T", sch.Type().TerraformType(ctx)),
		)
		return
	}

	vals := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, typ := range objType.AttributeTypes {
		v, present := ts.RawValues()[name]
		if !present || v == nil {
			vals[name] = tftypes.NewValue(typ, nil)
			continue
		}
		tv, convErr := v.ToTerraformValue(ctx)
		if convErr != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("UpgradeState[%s]: encode %q", r.managerName, name),
				convErr.Error(),
			)
			return
		}
		vals[name] = tv
	}

	resp.State = tfsdk.State{
		Schema: sch,
		Raw:    tftypes.NewValue(objType, vals),
	}
}
