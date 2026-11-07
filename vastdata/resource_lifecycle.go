// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/logging"
)

// ResourceLifecycle ----------------------------------
// Framework-agnostic lifecycle
// ----------------------------------
//
// ResourceLifecycle drives a ResourceManager's CRUD hooks (PrepareX / CreateResource /
// ReadResource / UpdateResource / DeleteResource / AfterX, LookupForCreate,
// TransformRequestBody, GetSubResources, ...) and the default API-backed behaviour used
// when a manager does not implement a hook.
//
// It deliberately knows nothing about Terraform: no plan, config, state requests,
// diagnostics, or tfsdk types. Inputs are managers whose ResourceState is already seeded.
// Outputs are plain errors, and the manager's ResourceState is updated in place so the
// caller can persist it through an is.StateSink.
type ResourceLifecycle struct {
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

	if err = forEachHook(manager, func(imp PrepareCreateResource) error {
		return imp.PrepareCreateResource(ctx, rest)
	}); err != nil {
		return l.stageErr("PrepareCreateResource", err)
	}

	var createRetryOn *is.RetryExpression
	if h := state.StateHints(); h != nil && h.RetryOn != nil {
		createRetryOn = h.RetryOn.Create
	}

	if imp, ok := findHook[CreateResource](manager); ok {
		createFn := imp.CreateResource
		if createRetryOn != nil {
			createFn = WithRetry(createRetryOn, "CreateResource", l.Name, imp.CreateResource)
		}
		record, err = createFn(ctx, rest)
	} else {
		if err = l.checkSearchable(ctx, state); err != nil {
			return err
		}
		createParams := state.GetCreateParams()
		if transformer, ok := findHook[TransformRequestBody](manager); ok {
			createParams = transformer.TransformRequestBody(createParams)
		}

		var lookupErr error
		if lookup, ok := findHook[LookupForCreate](manager); ok {
			record, lookupErr = lookup.LookupForCreate(ctx, rest)
		} else {
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
			if transformer, ok := findHook[TransformResponseRecord](manager); ok {
				existing = transformer.TransformResponseRecord(existing)
			}
			if normalizer, ok := findHook[NormalizeRecordForCreateAdopt](manager); ok {
				existing = normalizer.NormalizeRecordForCreateAdopt(existing)
			}
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

		if err = forEachHook(manager, func(imp AfterCreateResource) error {
			return imp.AfterCreateResource(ctx, rest, rec)
		}); err != nil {
			return l.stageErr("AfterCreateResource", err)
		} else if !hasOwnHook[AfterCreateResource](manager) {
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

		if transformer, ok := findHook[TransformResponseRecord](manager); ok {
			rec = transformer.TransformResponseRecord(rec)
		}
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
	if imp, ok := findHook[DeleteResource](manager); ok {
		_ = imp.DeleteResource(ctx, rest)
	} else if _, err := deleteRecordBySearchParams(ctx, manager.API(rest), manager.TfState(), l.Name, "TransactionDelete"); err != nil {
		logging.Warn(ctx, fmt.Sprintf("TransactionDelete[%s]: error deleting resource after error: %s", l.Name, err))
	}
	_ = forEachHook(manager, func(imp AfterDeleteResource) error {
		return imp.AfterDeleteResource(ctx, rest)
	})
}

// Read reads the external resource and refreshes manager's state from it.
// It reports exists=false (and no error) when the resource is gone.
func (l *ResourceLifecycle) Read(ctx context.Context, manager ResourceManager) (exists bool, err error) {
	var (
		rest   = l.Rest
		state  = manager.TfState()
		record DisplayableRecord
	)

	if err = forEachHook(manager, func(imp PrepareReadResource) error {
		return imp.PrepareReadResource(ctx, rest)
	}); err != nil {
		return false, l.stageErr("PrepareReadResource", err)
	}

	if imp, ok := findHook[ReadResource](manager); ok {
		record, err = imp.ReadResource(ctx, rest)
	} else {
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
		if transformer, ok := findHook[TransformResponseRecord](manager); ok {
			rec = transformer.TransformResponseRecord(rec)
		}
		PopulateIDFieldsFromNestedObjects(ctx, state, rec)
		if err = mergeSubResources(ctx, manager, rest, rec, l.Name); err != nil {
			return false, l.stageErr("GetSubResources", err)
		}
		if err = state.FillFromRecordIncludingRequired(rec, true); err != nil {
			return false, l.stageErr("fill state", err)
		}
	}

	if err = forEachHook(manager, func(imp AfterReadResource) error {
		return imp.AfterReadResource(ctx, rest)
	}); err != nil {
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

	desiredPrep, desiredPrepOK := findHook[PrepareUpdateResource](desired)
	if err = forEachHook(current, func(imp PrepareUpdateResource) error {
		if !desiredPrepOK {
			return fmt.Errorf("desired manager %T does not implement PrepareUpdateResource", desired)
		}
		return imp.PrepareUpdateResource(ctx, desiredPrep, rest)
	}); err != nil {
		return l.stageErr("PrepareUpdateResource", err)
	}

	if imp, ok := findHook[UpdateResource](current); ok {
		desiredUpd, ok := findHook[UpdateResource](desired)
		if !ok {
			return l.stageErr("UpdateResource", fmt.Errorf("desired manager %T does not implement UpdateResource", desired))
		}
		record, err = imp.UpdateResource(ctx, desiredUpd, rest)
	} else {
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
		if transformer, ok := findHook[TransformRequestBody](current); ok {
			updateParams = transformer.TransformRequestBody(updateParams)
		}
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
		if transformer, ok := findHook[TransformResponseRecord](current); ok {
			rec = transformer.TransformResponseRecord(rec)
		}
		// Merge sub-resources using the desired manager so trigger fields reflect the
		// wanted state, not the (potentially stale) current one.
		if err = mergeSubResources(ctx, desired, rest, rec, l.Name); err != nil {
			return l.stageErr("GetSubResources", err)
		}
		if err = currentState.FillFromRecord(rec); err != nil {
			return l.stageErr("fill state", err)
		}

		desiredAfter, desiredAfterOK := findHook[AfterUpdateResource](desired)
		if err = forEachHook(current, func(imp AfterUpdateResource) error {
			if !desiredAfterOK {
				return fmt.Errorf("desired manager %T does not implement AfterUpdateResource", desired)
			}
			return imp.AfterUpdateResource(ctx, desiredAfter, rest, rec)
		}); err != nil {
			return l.stageErr("AfterUpdateResource", err)
		} else if !hasOwnHook[AfterUpdateResource](current) {
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
						for k, v := range er {
							rec[k] = v
						}
					} else {
						for k, v := range editOnlyParams {
							rec[k] = v
						}
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

	if err = forEachHook(manager, func(imp PrepareDeleteResource) error {
		return imp.PrepareDeleteResource(ctx, rest)
	}); err != nil {
		return l.stageErr("PrepareDeleteResource", err)
	}

	var deleteRetryOn *is.RetryExpression
	if h := state.StateHints(); h != nil && h.RetryOn != nil {
		deleteRetryOn = h.RetryOn.Delete
	}

	deleteFn := func(ctx context.Context, rest *VMSRest) error {
		if imp, ok := findHook[DeleteResource](manager); ok {
			return imp.DeleteResource(ctx, rest)
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

	if err = forEachHook(manager, func(imp AfterDeleteResource) error {
		return imp.AfterDeleteResource(ctx, rest)
	}); err != nil {
		return l.stageErr("AfterDeleteResource", err)
	}
	return nil
}
