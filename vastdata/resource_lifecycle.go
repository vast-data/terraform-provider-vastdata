// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/logging"
)

// ResourceLifecycle ----------------
// Framework-agnostic lifecycle
// ----------------------------------
//
// ResourceLifecycle drives a ResourceManager's CRUD hooks (PrepareX / CreateResource /
// ReadResource / UpdateResource / DeleteResource / AfterX, LookupForCreate,
// TransformRequestBody, GetSubResources, UpgradeResourceSchema, ...) and the default
// API-backed behaviour used when a manager does not implement a hook.
//
// It deliberately knows nothing about Terraform: no plan, config, state requests,
// diagnostics, or tfsdk types. Inputs are managers whose ResourceState is already seeded.
// Outputs are plain errors, and the manager's ResourceState is updated in place so the
// caller can persist it through an is.StateSink.
//
// One lifecycle serves one resource kind for the whole session of its runtime: Terraform
// builds it when the resource is configured, Crossplane when it connects. It drives Read
// and whichever of Create, Update and Delete follow. The ResourceCallbacks it embeds are
// therefore shared by every operation it runs (see components.go, "Resource callbacks").
type ResourceLifecycle struct {
	ResourceCallbacks

	// Rest is the connected VAST client.
	Rest *VMSRest
	// Name is the manager's snake_case name, used in errors and logs.
	Name string
}

// NewResourceLifecycle returns a lifecycle bound to a client and resource name.
func NewResourceLifecycle(rest *VMSRest, name string) *ResourceLifecycle {
	return &ResourceLifecycle{Rest: rest, Name: name}
}

func (l *ResourceLifecycle) stageErr(stage string, err error) error {
	return fmt.Errorf("%s[%s]: %w", stage, l.Name, err)
}

func asRecord(d DisplayableRecord) Record {
	if d == nil {
		return nil
	}
	r, _ := d.(Record)
	return r
}

// checkSearchable fails when the state has no field the default CRUD can search by.
func (l *ResourceLifecycle) checkSearchable(ctx context.Context, state is.ResourceState) error {
	if len(getSearchParams(ctx, state, nil)) > 0 {
		return nil
	}
	available := state.GetFilteredValues(is.FilterOr, nil, is.SearchEmpty, is.SearchRequired, is.SearchSearchable)
	keys := make([]string, 0, len(available))
	for k := range available {
		keys = append(keys, k)
	}
	if len(keys) > 0 {
		return fmt.Errorf("no searchable fields provided for %q resource: at least one of the following fields must be set: %v", l.Name, keys)
	}
	return fmt.Errorf("no searchable fields provided for %q resource: none of the fields are marked as required or searchable in the OpenAPI schema", l.Name)
}

// postProcessRecord applies the steps shared by create and update once the API returned
// a record: async-task wait and the resolve-after-async hook.
func (l *ResourceLifecycle) postProcessRecord(
	ctx context.Context,
	op string,
	manager ResourceManager,
	record Record,
) (Record, error) {
	rest := l.Rest
	var asyncTimeout *time.Duration
	if h := manager.TfState().StateHints(); h != nil {
		asyncTimeout = h.AsyncTaskTimeout
	}
	if err := handleMaybeAsyncTask(ctx, rest, record, asyncTimeout); err != nil {
		return nil, l.stageErr("AsyncTask - "+op, err)
	}
	resolved, err := resolveRecordAfterAsyncTask(ctx, manager, rest, record, l.Name)
	if err != nil {
		return nil, l.stageErr("ResolveRecordAfterAsyncTask - "+op, err)
	}
	record = resolved
	return record, nil
}

// Create creates the external resource described by manager's state, adopting an existing
// one when LookupForCreate / the default search finds it. On success the manager's state is
// refreshed from the API record. If a resource was created and a later step fails, it is
// rolled back (same as Terraform's transaction delete).
func (l *ResourceLifecycle) Create(ctx context.Context, manager ResourceManager) (err error) {
	var (
		rest    = l.Rest
		api     = manager.API(rest)
		state   = manager.TfState()
		record  DisplayableRecord
		created bool
		// Original (desired) values; write-only values are still present here.
		desired = state.Copy()
	)
	defer func() {
		if err != nil && created {
			l.rollbackCreate(ctx, manager)
		}
	}()

	if err = l.runPrepareCreateResource(ctx, manager); err != nil {
		return l.stageErr("PrepareCreateResource", err)
	}

	var createRetryOn *is.RetryExpression
	if h := state.StateHints(); h != nil && h.RetryOn != nil {
		createRetryOn = h.RetryOn.Create
	}

	var customCreate bool
	record, customCreate, err = l.runCreateResource(ctx, manager, createRetryOn)
	if !customCreate {
		if err = l.checkSearchable(ctx, state); err != nil {
			return err
		}
		createParams := state.GetCreateParams()
		createParams = runTransformRequestBody(manager, createParams)

		var (
			lookupErr    error
			customLookup bool
		)
		record, customLookup, lookupErr = l.runLookupForCreate(ctx, manager)
		if !customLookup {
			record, lookupErr = getRecordBySearchParams(ctx, api, state, nil, l.Name, "Create")
		}

		switch {
		case lookupErr != nil && !isNotFoundErr(lookupErr):
			// Something unexpected happened, we must not continue.
			return l.stageErr("read before create", lookupErr)

		case lookupErr != nil:
			// Nothing exists yet: create it.
			created = true
			defaultCreate := func(ctx context.Context, _ *VMSRest) (DisplayableRecord, error) {
				return api.CreateWithContext(ctx, createParams)
			}
			if createRetryOn != nil {
				defaultCreate = WithRetry(createRetryOn, "CreateResource", l.Name, defaultCreate)
			}
			record, err = defaultCreate(ctx, rest)

		default:
			// Found an existing resource: take management of it.
			existing := asRecord(record)
			existing = runTransformResponseRecord(manager, existing)
			existing = runNormalizeRecordForCreateAdopt(manager, existing)
			record = existing
			// NOTE: the default implementation only works for resources with an 'id' field.
			// Resources without one must implement CreateResource.
			if diff := diffMap(createParams, existing); len(diff) > 0 {
				id, exists := existing["id"]
				if !exists {
					return l.stageErr("adopt", fmt.Errorf("record does not have 'id' field; implement CreateResource for this resource"))
				}
				// Send only the difference between the existing record and what we want.
				record, err = api.UpdateWithContext(ctx, id, diff)
			}
		}
	}
	if err != nil {
		return l.stageErr("create", err)
	}

	rec := asRecord(record)
	if rec != nil && rec.Empty() {
		rec = nil
	}

	if rec != nil {
		if rec, err = l.postProcessRecord(ctx, "create", manager, rec); err != nil {
			return err
		}

		_, ownAfterCreate := manager.(AfterCreateResource)
		if err = l.runAfterCreateResource(ctx, manager, rec); err != nil {
			return l.stageErr("AfterCreateResource", err)
		} else if !ownAfterCreate {
			if h := state.StateHints(); h != nil && len(h.EditOnlyFields) > 0 {
				// Fields that cannot be set on creation (e.g. "enabled") are applied via update.
				if updateParams := state.GetReadEditOnlyParams(); len(updateParams) > 0 {
					id, exists := rec["id"]
					if !exists {
						return l.stageErr("create edit-only fields", fmt.Errorf("record does not have 'id' field"))
					}
					if _, err = api.UpdateWithContext(ctx, id, updateParams); err != nil {
						return l.stageErr("create edit-only fields", err)
					}
					for k, v := range updateParams {
						rec[k] = v
					}
				}
			}
		}

		rec = runTransformResponseRecord(manager, rec)
		if err = mergeSubResources(ctx, manager, rest, rec, l.Name); err != nil {
			return l.stageErr("GetSubResources", err)
		}
		if err = state.FillFromRecord(rec); err != nil {
			return l.stageErr("fill state", err)
		}
	}

	// Align the desired values with what the backend actually returned.
	desired.CopyNonEmptyFieldsTo(state)
	state.ClearWriteOnlyFields()
	return nil
}

func (l *ResourceLifecycle) rollbackCreate(ctx context.Context, manager ResourceManager) {
	rest := l.Rest
	if handled, _ := l.runDeleteResource(ctx, manager); !handled {
		if _, err := deleteRecordBySearchParams(ctx, manager.API(rest), manager.TfState(), l.Name, "TransactionDelete"); err != nil {
			logging.Warn(ctx, fmt.Sprintf("TransactionDelete[%s]: error deleting resource after error: %s", l.Name, err))
		}
	}
	_ = l.runAfterDeleteResource(ctx, manager)
}

// Read reads the external resource and refreshes manager's state from it.
// It reports exists=false (and no error) when the resource is gone.
func (l *ResourceLifecycle) Read(ctx context.Context, manager ResourceManager) (exists bool, err error) {
	var (
		rest   = l.Rest
		state  = manager.TfState()
		record DisplayableRecord
	)

	if err = l.runPrepareReadResource(ctx, manager); err != nil {
		return false, l.stageErr("PrepareReadResource", err)
	}

	var customRead bool
	record, customRead, err = l.runReadResource(ctx, manager)
	if !customRead {
		if err = l.checkSearchable(ctx, state); err != nil {
			return false, err
		}
		if h := state.StateHints(); h != nil && h.SkipRefreshAPICall {
			record = Record(state.GetAllValues())
		} else {
			record, err = getRecordBySearchParams(ctx, manager.API(rest), state, nil, l.Name, "Read")
		}
	}
	if err != nil {
		if isNotFoundErr(err) || expectStatusCodes(err, http.StatusNotFound) || errors.As(err, &ForceCleanState{}) {
			return false, nil
		}
		return false, l.stageErr("read", err)
	}

	if rec := asRecord(record); rec != nil {
		rec = runTransformResponseRecord(manager, rec)
		PopulateIDFieldsFromNestedObjects(ctx, state, rec)
		if err = mergeSubResources(ctx, manager, rest, rec, l.Name); err != nil {
			return false, l.stageErr("GetSubResources", err)
		}
		if err = state.FillFromRecordIncludingRequired(rec, true); err != nil {
			return false, l.stageErr("fill state", err)
		}
	}

	if err = l.runAfterReadResource(ctx, manager); err != nil {
		return false, l.stageErr("AfterReadResource", err)
	}
	return true, nil
}

// Update applies the difference between current (what exists) and desired (what is wanted)
// and refreshes current's state from the API. Both managers must be of the same resource type.
func (l *ResourceLifecycle) Update(ctx context.Context, current, desired ResourceManager) (err error) {
	var (
		rest         = l.Rest
		api          = current.API(rest)
		currentState = current.TfState()
		desiredState = desired.TfState()
		record       DisplayableRecord
	)

	if err = l.runPrepareUpdateResource(ctx, current, desired); err != nil {
		return l.stageErr("PrepareUpdateResource", err)
	}

	var customUpdate bool
	record, customUpdate, err = l.runUpdateResource(ctx, current, desired)
	if !customUpdate {
		if err = l.checkSearchable(ctx, desiredState); err != nil {
			return err
		}
		if record, err = getRecordBySearchParams(ctx, api, currentState, desiredState, l.Name, "Update"); err != nil {
			return l.stageErr("update", err)
		}
		// NOTE: the default implementation only works for resources with an 'id' field.
		id, exists := asRecord(record)["id"]
		if !exists {
			return l.stageErr("update", fmt.Errorf("record does not have 'id' field; implement UpdateResource for this resource"))
		}
		updateParams := desiredState.GetUpdateParams(currentState)
		updateParams = runTransformRequestBody(current, updateParams)
		delete(updateParams, "id") // The ID must never be updated.
		if len(updateParams) > 0 {
			record, err = api.UpdateWithContext(ctx, id, updateParams)
		}
	}
	if err != nil {
		return l.stageErr("update", err)
	}

	rec := asRecord(record)
	if rec != nil && rec.Empty() {
		rec = nil
	}

	if rec != nil {
		if rec, err = l.postProcessRecord(ctx, "update", current, rec); err != nil {
			return err
		}
		rec = runTransformResponseRecord(current, rec)
		// Merge sub-resources using the desired manager so trigger fields reflect the
		// wanted state, not the (potentially stale) current one.
		if err = mergeSubResources(ctx, desired, rest, rec, l.Name); err != nil {
			return l.stageErr("GetSubResources", err)
		}
		if err = currentState.FillFromRecord(rec); err != nil {
			return l.stageErr("fill state", err)
		}

		_, ownAfterUpdate := current.(AfterUpdateResource)
		if err = l.runAfterUpdateResource(ctx, current, desired, rec); err != nil {
			return l.stageErr("AfterUpdateResource", err)
		} else if !ownAfterUpdate {
			if h := currentState.StateHints(); h != nil && len(h.EditOnlyFields) > 0 {
				// If AfterUpdateResource is implemented it is assumed to handle edit-only fields itself.
				if editOnlyParams := desiredState.GetChangedEditOnlyParams(currentState); len(editOnlyParams) > 0 {
					id, exists := rec["id"]
					if !exists {
						return l.stageErr("update edit-only fields", fmt.Errorf("record does not have 'id' field"))
					}
					editRecord, editErr := api.UpdateWithContext(ctx, id, editOnlyParams)
					if editErr != nil {
						return l.stageErr("update edit-only fields", editErr)
					}
					if er := asRecord(editRecord); er != nil && !er.Empty() {
						maps.Copy(rec, er)
					} else {
						maps.Copy(rec, editOnlyParams)
					}
				}
			}
		}
	}

	// Copy wanted values onto the refreshed state.
	desiredState.CopyKnownFieldsTo(currentState)
	currentState.ClearWriteOnlyFields()
	return nil
}

// Delete removes the external resource described by manager's state.
func (l *ResourceLifecycle) Delete(ctx context.Context, manager ResourceManager) (err error) {
	var (
		rest  = l.Rest
		state = manager.TfState()
	)

	if err = l.runPrepareDeleteResource(ctx, manager); err != nil {
		return l.stageErr("PrepareDeleteResource", err)
	}

	var deleteRetryOn *is.RetryExpression
	if h := state.StateHints(); h != nil && h.RetryOn != nil {
		deleteRetryOn = h.RetryOn.Delete
	}

	deleteFn := func(ctx context.Context, rest *VMSRest) error {
		if handled, err := l.runDeleteResource(ctx, manager); handled {
			return err
		}
		if err := l.checkSearchable(ctx, state); err != nil {
			return err
		}
		record, delErr := deleteRecordBySearchParams(ctx, manager.API(rest), state, l.Name, "Delete")
		if delErr != nil {
			return delErr
		}
		if record != nil {
			var asyncTimeout *time.Duration
			if h := state.StateHints(); h != nil {
				asyncTimeout = h.AsyncTaskTimeout
			}
			return handleMaybeAsyncTask(ctx, rest, record, asyncTimeout)
		}
		return nil
	}

	if deleteRetryOn != nil {
		_, err = WithRetry(deleteRetryOn, "DeleteResource", l.Name, func(ctx context.Context, rest *VMSRest) (struct{}, error) {
			return struct{}{}, deleteFn(ctx, rest)
		})(ctx, rest)
	} else {
		err = deleteFn(ctx, rest)
	}
	if err != nil {
		return l.stageErr("delete", err)
	}

	if err = l.runAfterDeleteResource(ctx, manager); err != nil {
		return l.stageErr("AfterDeleteResource", err)
	}
	return nil
}

// -----------------
// Stage runners for the regular hooks
// -----------------
//
// Each runner hides the type assertion for one optional hook of a manager. The CRUD ones
// report whether the manager implements the hook (handled): when it does not, the caller
// falls back to the default API-backed behavior. The stages with callbacks are further down,
// next to ResourceCallbacks.

// runCreateResource runs the manager's CreateResource, retrying it with retryOn when set.
func (l *ResourceLifecycle) runCreateResource(ctx context.Context, m ResourceManager, retryOn *is.RetryExpression) (DisplayableRecord, bool, error) {
	imp, ok := m.(CreateResource)
	if !ok {
		return nil, false, nil
	}
	createFn := imp.CreateResource
	if retryOn != nil {
		createFn = WithRetry(retryOn, "CreateResource", l.Name, createFn)
	}
	record, err := createFn(ctx, l.Rest)
	return record, true, err
}

// runLookupForCreate runs the manager's LookupForCreate, which finds an existing resource
// to adopt before creating.
func (l *ResourceLifecycle) runLookupForCreate(ctx context.Context, m ResourceManager) (DisplayableRecord, bool, error) {
	imp, ok := m.(LookupForCreate)
	if !ok {
		return nil, false, nil
	}
	record, err := imp.LookupForCreate(ctx, l.Rest)
	return record, true, err
}

// runReadResource runs the manager's ReadResource.
func (l *ResourceLifecycle) runReadResource(ctx context.Context, m ResourceManager) (DisplayableRecord, bool, error) {
	imp, ok := m.(ReadResource)
	if !ok {
		return nil, false, nil
	}
	record, err := imp.ReadResource(ctx, l.Rest)
	return record, true, err
}

// runUpdateResource runs current's UpdateResource with desired as the plan. Both managers
// must implement it.
func (l *ResourceLifecycle) runUpdateResource(ctx context.Context, current, desired ResourceManager) (DisplayableRecord, bool, error) {
	imp, ok := current.(UpdateResource)
	if !ok {
		return nil, false, nil
	}
	plan, ok := desired.(UpdateResource)
	if !ok {
		return nil, true, fmt.Errorf("desired manager %T does not implement UpdateResource", desired)
	}
	record, err := imp.UpdateResource(ctx, plan, l.Rest)
	return record, true, err
}

// runDeleteResource runs the manager's DeleteResource.
func (l *ResourceLifecycle) runDeleteResource(ctx context.Context, m ResourceManager) (bool, error) {
	imp, ok := m.(DeleteResource)
	if !ok {
		return false, nil
	}
	return true, imp.DeleteResource(ctx, l.Rest)
}

// runPrepareImportResourceState runs the manager's PrepareImportResourceState, if any.
func (l *ResourceLifecycle) runPrepareImportResourceState(ctx context.Context, m ResourceManager) error {
	if imp, ok := m.(PrepareImportResourceState); ok {
		return imp.PrepareImportResourceState(ctx, l.Rest)
	}
	return nil
}

// runImportResourceState runs the manager's ImportResourceState and reports whether it
// has one. Without one, the caller imports with the default import ID parsing.
func (l *ResourceLifecycle) runImportResourceState(ctx context.Context, req resource.ImportStateRequest, m ResourceManager) (bool, error) {
	imp, ok := m.(ImportResourceState)
	if !ok {
		return false, nil
	}
	return true, imp.ImportResourceState(req, ctx, l.Rest)
}

// runAfterImportResourceState runs the manager's AfterImportResourceState, if any.
func (l *ResourceLifecycle) runAfterImportResourceState(ctx context.Context, m ResourceManager) error {
	if imp, ok := m.(AfterImportResourceState); ok {
		return imp.AfterImportResourceState(ctx, l.Rest)
	}
	return nil
}

// -----------------
// Transforms
// -----------------
//
// Plain functions of the manager: the datasource read path uses them as well and has no
// lifecycle. Without the hook the value is returned as it is.

// runTransformRequestBody lets the manager rewrite the body of a create or update request.
func runTransformRequestBody(m any, body params) params {
	if imp, ok := m.(TransformRequestBody); ok {
		return imp.TransformRequestBody(body)
	}
	return body
}

// runTransformResponseRecord lets the manager rewrite a record the API returned.
func runTransformResponseRecord(m any, record Record) Record {
	if imp, ok := m.(TransformResponseRecord); ok {
		return imp.TransformResponseRecord(record)
	}
	return record
}

// runNormalizeRecordForCreateAdopt lets the manager normalize an existing record before
// create adopts it.
func runNormalizeRecordForCreateAdopt(m any, record Record) Record {
	if imp, ok := m.(NormalizeRecordForCreateAdopt); ok {
		return imp.NormalizeRecordForCreateAdopt(record)
	}
	return record
}

// -----------------
// Resource callbacks
// -----------------

// ResourceCallbacks stores the callbacks of every stage. ResourceLifecycle embeds it, which
// makes the lifecycle a ResourceCallbackExecutor. The zero value has no callbacks.
type ResourceCallbacks struct {
	prepareCreate []func(context.Context, *VMSRest) error
	afterCreate   []func(context.Context, *VMSRest, Record) error
	prepareRead   []func(context.Context, *VMSRest) error
	afterRead     []func(context.Context, *VMSRest) error
	prepareUpdate []func(context.Context, ResourceManager, *VMSRest) error
	afterUpdate   []func(context.Context, ResourceManager, *VMSRest, Record) error
	prepareDelete []func(context.Context, *VMSRest) error
	afterDelete   []func(context.Context, *VMSRest) error
	upgradeSchema []func(ctx context.Context, fromVersion int64, ts *is.TFState) error
}

var _ ResourceCallbackExecutor = (*ResourceCallbacks)(nil)

// runAll calls fn for every callback and returns the first error.
func runAll[F any](callbacks []F, fn func(F) error) error {
	for _, cb := range callbacks {
		if err := fn(cb); err != nil {
			return err
		}
	}
	return nil
}

func (c *ResourceCallbacks) AddPrepareCreateResourceCallback(fn func(context.Context, *VMSRest) error) {
	c.prepareCreate = append(c.prepareCreate, fn)
}

func (c *ResourceCallbacks) ExecutePrepareCreateResourceCallbacks(ctx context.Context, rest *VMSRest) error {
	return runAll(c.prepareCreate, func(fn func(context.Context, *VMSRest) error) error {
		return fn(ctx, rest)
	})
}

func (c *ResourceCallbacks) AddAfterCreateResourceCallback(fn func(context.Context, *VMSRest, Record) error) {
	c.afterCreate = append(c.afterCreate, fn)
}

func (c *ResourceCallbacks) ExecuteAfterCreateResourceCallbacks(ctx context.Context, rest *VMSRest, record Record) error {
	return runAll(c.afterCreate, func(fn func(context.Context, *VMSRest, Record) error) error {
		return fn(ctx, rest, record)
	})
}

func (c *ResourceCallbacks) AddPrepareReadResourceCallback(fn func(context.Context, *VMSRest) error) {
	c.prepareRead = append(c.prepareRead, fn)
}

func (c *ResourceCallbacks) ExecutePrepareReadResourceCallbacks(ctx context.Context, rest *VMSRest) error {
	return runAll(c.prepareRead, func(fn func(context.Context, *VMSRest) error) error {
		return fn(ctx, rest)
	})
}

func (c *ResourceCallbacks) AddAfterReadResourceCallback(fn func(context.Context, *VMSRest) error) {
	c.afterRead = append(c.afterRead, fn)
}

func (c *ResourceCallbacks) ExecuteAfterReadResourceCallbacks(ctx context.Context, rest *VMSRest) error {
	return runAll(c.afterRead, func(fn func(context.Context, *VMSRest) error) error {
		return fn(ctx, rest)
	})
}

func (c *ResourceCallbacks) AddPrepareUpdateResourceCallback(fn func(context.Context, ResourceManager, *VMSRest) error) {
	c.prepareUpdate = append(c.prepareUpdate, fn)
}

func (c *ResourceCallbacks) ExecutePrepareUpdateResourceCallbacks(ctx context.Context, desired ResourceManager, rest *VMSRest) error {
	return runAll(c.prepareUpdate, func(fn func(context.Context, ResourceManager, *VMSRest) error) error {
		return fn(ctx, desired, rest)
	})
}

func (c *ResourceCallbacks) AddAfterUpdateResourceCallback(fn func(context.Context, ResourceManager, *VMSRest, Record) error) {
	c.afterUpdate = append(c.afterUpdate, fn)
}

func (c *ResourceCallbacks) ExecuteAfterUpdateResourceCallbacks(ctx context.Context, desired ResourceManager, rest *VMSRest, record Record) error {
	return runAll(c.afterUpdate, func(fn func(context.Context, ResourceManager, *VMSRest, Record) error) error {
		return fn(ctx, desired, rest, record)
	})
}

func (c *ResourceCallbacks) AddPrepareDeleteResourceCallback(fn func(context.Context, *VMSRest) error) {
	c.prepareDelete = append(c.prepareDelete, fn)
}

func (c *ResourceCallbacks) ExecutePrepareDeleteResourceCallbacks(ctx context.Context, rest *VMSRest) error {
	return runAll(c.prepareDelete, func(fn func(context.Context, *VMSRest) error) error {
		return fn(ctx, rest)
	})
}

func (c *ResourceCallbacks) AddAfterDeleteResourceCallback(fn func(context.Context, *VMSRest) error) {
	c.afterDelete = append(c.afterDelete, fn)
}

func (c *ResourceCallbacks) ExecuteAfterDeleteResourceCallbacks(ctx context.Context, rest *VMSRest) error {
	return runAll(c.afterDelete, func(fn func(context.Context, *VMSRest) error) error {
		return fn(ctx, rest)
	})
}

func (c *ResourceCallbacks) AddUpgradeResourceSchemaCallback(fn func(ctx context.Context, fromVersion int64, ts *is.TFState) error) {
	c.upgradeSchema = append(c.upgradeSchema, fn)
}

func (c *ResourceCallbacks) ExecuteUpgradeResourceSchemaCallbacks(ctx context.Context, fromVersion int64, ts *is.TFState) error {
	return runAll(c.upgradeSchema, func(fn func(context.Context, int64, *is.TFState) error) error {
		return fn(ctx, fromVersion, ts)
	})
}

// UpgradeResourceSchema runs the manager's UpgradeResourceSchema hook.
func (l *ResourceLifecycle) UpgradeResourceSchema(
	ctx context.Context,
	fromVersion int64,
	manager ResourceManager,
	ts *is.TFState,
) error {
	if err := l.runUpgradeResourceSchema(ctx, fromVersion, manager, ts); err != nil {
		return l.stageErr("UpgradeResourceSchema", err)
	}
	return nil
}

// -----------------
// Stage runners: the manager's own hook first, then the lifecycle's callbacks.
// A failing hook is returned immediately, and no callback of that stage runs.
// -----------------

func (l *ResourceLifecycle) runPrepareCreateResource(ctx context.Context, m ResourceManager) error {
	if imp, ok := m.(PrepareCreateResource); ok {
		if err := imp.PrepareCreateResource(ctx, l.Rest); err != nil {
			return err
		}
	}
	return l.ExecutePrepareCreateResourceCallbacks(ctx, l.Rest)
}

func (l *ResourceLifecycle) runAfterCreateResource(ctx context.Context, m ResourceManager, record Record) error {
	if imp, ok := m.(AfterCreateResource); ok {
		if err := imp.AfterCreateResource(ctx, l.Rest, record); err != nil {
			return err
		}
	}
	return l.ExecuteAfterCreateResourceCallbacks(ctx, l.Rest, record)
}

func (l *ResourceLifecycle) runPrepareReadResource(ctx context.Context, m ResourceManager) error {
	if imp, ok := m.(PrepareReadResource); ok {
		if err := imp.PrepareReadResource(ctx, l.Rest); err != nil {
			return err
		}
	}
	return l.ExecutePrepareReadResourceCallbacks(ctx, l.Rest)
}

func (l *ResourceLifecycle) runAfterReadResource(ctx context.Context, m ResourceManager) error {
	if imp, ok := m.(AfterReadResource); ok {
		if err := imp.AfterReadResource(ctx, l.Rest); err != nil {
			return err
		}
	}
	return l.ExecuteAfterReadResourceCallbacks(ctx, l.Rest)
}

// runPrepareUpdateResource runs current's hook and the callbacks with desired as the plan.
func (l *ResourceLifecycle) runPrepareUpdateResource(ctx context.Context, current, desired ResourceManager) error {
	if imp, ok := current.(PrepareUpdateResource); ok {
		plan, ok := desired.(PrepareUpdateResource)
		if !ok {
			return errors.New("desired manager does not implement PrepareUpdateResource")
		}
		if err := imp.PrepareUpdateResource(ctx, plan, l.Rest); err != nil {
			return err
		}
	}
	return l.ExecutePrepareUpdateResourceCallbacks(ctx, desired, l.Rest)
}

// runAfterUpdateResource runs current's hook and the callbacks with desired as the plan.
func (l *ResourceLifecycle) runAfterUpdateResource(ctx context.Context, current, desired ResourceManager, record Record) error {
	if imp, ok := current.(AfterUpdateResource); ok {
		plan, ok := desired.(AfterUpdateResource)
		if !ok {
			return errors.New("desired manager does not implement AfterUpdateResource")
		}
		if err := imp.AfterUpdateResource(ctx, plan, l.Rest, record); err != nil {
			return err
		}
	}
	return l.ExecuteAfterUpdateResourceCallbacks(ctx, desired, l.Rest, record)
}

func (l *ResourceLifecycle) runPrepareDeleteResource(ctx context.Context, m ResourceManager) error {
	if imp, ok := m.(PrepareDeleteResource); ok {
		if err := imp.PrepareDeleteResource(ctx, l.Rest); err != nil {
			return err
		}
	}
	return l.ExecutePrepareDeleteResourceCallbacks(ctx, l.Rest)
}

func (l *ResourceLifecycle) runAfterDeleteResource(ctx context.Context, m ResourceManager) error {
	if imp, ok := m.(AfterDeleteResource); ok {
		if err := imp.AfterDeleteResource(ctx, l.Rest); err != nil {
			return err
		}
	}
	return l.ExecuteAfterDeleteResourceCallbacks(ctx, l.Rest)
}

func (l *ResourceLifecycle) runUpgradeResourceSchema(
	ctx context.Context,
	fromVersion int64,
	m ResourceManager,
	ts *is.TFState,
) error {
	if imp, ok := m.(UpgradeResourceSchema); ok {
		if err := imp.UpgradeResourceSchema(ctx, fromVersion, ts); err != nil {
			return err
		}
	}
	return l.ExecuteUpgradeResourceSchemaCallbacks(ctx, fromVersion, ts)
}
