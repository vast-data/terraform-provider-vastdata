// Copyright (c) HashiCorp, Inc.

// This file defines commonly used aliases for external libraries (e.g., vast_client),
// along with shared helper functions and constants used throughout the Terraform provider.
//
// Contents include:
// - Type aliases for REST client and schema utilities
// - Logging wrappers for consistent diagnostic output
// - Validation-related constants used in schema generation
// - General-purpose utility functions (e.g., deep equality comparison, type normalization)
//
// Intended to reduce duplication and centralize shared logic across resources and internal tooling.

package provider

import (
	"context"
	"fmt"
	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	vast_client "github.com/vast-data/go-vast-client"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/client"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/schema_generation"
	"net/http"
	"strings"
)

// Rest Client
type (
	VastResourceAPI            = vast_client.VastResourceAPI
	VastResourceAPIWithContext = vast_client.VastResourceAPIWithContext
	DisplayableRecord          = vast_client.DisplayableRecord
	Renderable                 = vast_client.Renderable
	Record                     = vast_client.Record
	RecordSet                  = vast_client.RecordSet
	params                     = vast_client.Params
	VMSRest                    = vast_client.VMSRest
	RestFn                     func(context.Context, params) (Record, error)
	ApiError                   = vast_client.ApiError
)

var (
	isNotFoundErr     = vast_client.IsNotFoundErr
	ignoreNotFound    = vast_client.IgnoreNotFound
	ignoreStatusCodes = vast_client.IgnoreStatusCodes
	expectStatusCodes = vast_client.ExpectStatusCodes
	isApiError        = vast_client.IsApiError
)

// Helpers
var (
	toInt   func(val any) (int64, error)             = is.ToInt
	toFloat func(val any) (float64, error)           = is.ToFloat
	diffMap func(a, b map[string]any) map[string]any = is.DiffMap
)

// Common validators
const (
	ValidatorPathStartsWithSlash     = schema_generation.ValidatorPathStartsWithSlash
	ValidatorPathStartsEndsWithSlash = schema_generation.ValidatorPathStartsEndsWithSlash
	ValidatorGracePeriodFormat       = schema_generation.ValidatorGracePeriodFormat
	ValidatorRFC3339Format           = schema_generation.ValidatorRFC3339Format
	ValidatorIntervalFormat          = schema_generation.ValidatorIntervalFormat
	ValidatorRetentionFormat         = schema_generation.ValidatorRetentionFormat
)

// MISC
func withContext(ctx context.Context, method string, managerName string, fn func(ctx context.Context)) {
	ctx = client.ContextWithRequestID(ctx)
	tflog.Debug(ctx, fmt.Sprintf("◉ %s[%s] start", method, managerName))
	defer tflog.Debug(ctx, fmt.Sprintf("◉ %s[%s] end", method, managerName))
	fn(ctx)
}

func safeDeepEqual(expected, actual any) (diff []string, panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()
	diff = deep.Equal(expected, actual, deep.FLAG_IGNORE_SLICE_ORDER)
	return diff, false
}

// ----------------------------------
// Common Resource Validators
// ----------------------------------

// validateOneOf ensures that at most one of the given fields is set (i.e., known and not null).
// Returns an error if more than one field is set.
func validateOneOf(tf *is.TFState, fields ...string) error {
	var setFields []string

	for _, field := range fields {
		if tf.IsKnownAndNotNull(field) {
			setFields = append(setFields, field)
		}
	}

	if len(setFields) > 1 {
		return fmt.Errorf("only one of %q can be set, but multiple were provided: %v", fields, setFields)
	}

	return nil
}

// validateAllOf ensures that all the given fields are set (i.e., known and not null).
// Returns an error listing the fields that are missing.
func validateAllOf(tf *is.TFState, fields ...string) error {
	var missing []string

	for _, field := range fields {
		if !tf.IsKnownAndNotNull(field) {
			missing = append(missing, field)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("all of %q must be set, but missing: %v", fields, missing)
	}

	return nil
}

// validateNoneOf ensures that none of the given fields are set (i.e., all are null or unknown).
// Returns an error listing the fields that were incorrectly set.
func validateNoneOf(tf *is.TFState, fields ...string) error {
	var setFields []string

	for _, field := range fields {
		if tf.IsKnownAndNotNull(field) {
			setFields = append(setFields, field)
		}
	}

	if len(setFields) > 0 {
		return fmt.Errorf("none of %q should be set, but found: %v", fields, setFields)
	}

	return nil
}

// ----------------------------------
// Query Resources and Search Parameters
// ----------------------------------

// getSearchParams extracts generic search parameters from the given TF state and optional Plan state.
// It ensures that at least one searchable parameter is present, otherwise returns an error.
// This is typically used to build queries for lookup operations.
//
// Returns:
//   - map[string]any: a map of search fields extracted from state (e.g., id, guid, name, etc.)
//   - error: if no searchable fields are found
func getSearchParams(ctx context.Context, tfState, planTfState *is.TFState) params {
	searchParams := tfState.GetGenericSearchParams(ctx)
	if planTfState != nil {
		if sp := planTfState.GetReadOnlySearchParams(); len(sp) > 0 {
			for k, v := range sp {
				searchParams[k] = v
			}
		}
	}
	return searchParams

}

// getRecordBySearchParams attempts to retrieve a record using parameters extracted from TF state.
// It performs lookup in the following order:
//  1. By ID (if present)
//  2. By GUID (if ID not present or returns 404)
//  3. By other search parameters (if ID and GUID are missing or both fail)
//
// Parameters:
//   - ctx: request context
//   - api: resource client interface
//   - tfState: Terraform state wrapper
//   - planTfState: Terraform state wrapper for plan (optional)
//   - managerName: name of resource manager for logging
//   - op: operation name (e.g. "read", "update") for logging
//
// Returns:
//   - DisplayableRecord: the found record
//   - error: if not found or all lookups fail
func getRecordBySearchParams(ctx context.Context, api VastResourceAPIWithContext, tfState *is.TFState, planTfState *is.TFState, managerName, op string) (DisplayableRecord, error) {
	var (
		record DisplayableRecord
		err    error
	)
	searchParams := getSearchParams(ctx, tfState, planTfState)
	if len(searchParams) == 0 {
		return nil, fmt.Errorf(
			"%s[%s]: no search parameters provided for %q resource."+
				" Verify presence of required fields or add searchable hints to resource",
			op,
			managerName,
			managerName,
		)
	}

	id, idExists := searchParams["id"]
	guid, guidExists := searchParams["guid"]

	// Attempt 1: get resource by ID.
	if idExists {
		tflog.Debug(ctx, fmt.Sprintf("%s[%s]: found ID = %v.", op, managerName, id))
		// If the ID is set, we assume it's a direct call by ID
		record, err = api.GetByIdWithContext(ctx, is.Must(toInt(id)))
	}
	// Attempt 2: get resource by GUID.
	if (!idExists || expectStatusCodes(err, http.StatusNotFound)) && guidExists {
		tflog.Debug(ctx, fmt.Sprintf("%s[%s]: found GUID = %v.", op, managerName, guid))
		record, err = api.GetWithContext(ctx, params{"guid": guid})
	}
	// Attempt 3: if both ID and GUID are not set or both reads failed, we use the search parameters.
	if (!idExists && !guidExists) || expectStatusCodes(err, http.StatusNotFound) {
		// If neither ID nor GUID is set or both reads failed, we use the search parameters
		searchParams.Without("id", "guid")
		if len(searchParams) > 0 {
			if !idExists && !guidExists {
				tflog.Debug(
					ctx,
					fmt.Sprintf("%s[%s]: no ID or GUID set, using search parameters %v.",
						op, managerName, searchParams,
					),
				)
			} else {
				tflog.Debug(
					ctx,
					fmt.Sprintf("%s[%s]: ID %v or GUID %v not found, using search parameters %v.",
						op, managerName, id, guid, searchParams,
					),
				)
			}
			record, err = api.GetWithContext(ctx, searchParams)
		}
	}

	return record, err

}

// deleteRecordBySearchParams attempts to delete a resource by ID (preferred) or by generic search parameters.
// It also supports appending additional deletion parameters using TFState hints (DeleteOnlyFields).
//
// Parameters:
//   - ctx: request context
//   - api: resource client interface
//   - tfState: Terraform state wrapper
//   - managerName: name of resource manager for logging
//   - op: operation name (e.g. "delete") for logging
//
// Returns:
//   - error: any error that occurred during deletion, excluding 404s (they are ignored)
func deleteRecordBySearchParams(ctx context.Context, api VastResourceAPIWithContext, tfState *is.TFState, managerName, op string) error {
	var err error
	searchParams := getSearchParams(ctx, tfState, nil)
	if len(searchParams) == 0 {
		return fmt.Errorf(
			"%s[%s]: no search parameters provided for %q resource."+
				" Verify presence of required fields or add searchable hints to resource",
			op,
			managerName,
			managerName,
		)
	}

	var deleteParams map[string]any
	if len(tfState.Hints.DeleteOnlyFields) > 0 {
		deleteParams = tfState.GetDeleteOnlyParams()
	}

	if id, ok := searchParams["id"]; ok {
		tflog.Debug(ctx, fmt.Sprintf("%s[%s]: found ID = %v.", op, managerName, id))
		// If the ID is set, we assume it's a direct call by ID
		id := is.Must(toInt(id))
		_, err = api.DeleteByIdWithContext(ctx, id, deleteParams)
	} else {
		tflog.Debug(ctx, fmt.Sprintf("%s[%s]: no ID found, using search params.", op, managerName))
		_, err = api.DeleteWithContext(ctx, searchParams, deleteParams)
	}
	return ignoreStatusCodes(err, http.StatusNotFound)

}

// ----------------------------------
// Transformations
// ----------------------------------

// normalizeNumber attempts to convert float64 values to int64 if they are whole numbers.
// This is useful for decoding JSON where numbers are unmarshaled as float64,
// but Terraform schema (or other systems) expect integer types.
// It also recursively normalizes values in slices and maps.
func normalizeNumber(v any) any {
	switch n := v.(type) {
	case float64:
		if float64(int64(n)) == n {
			return int64(n)
		}
	case []any:
		for i, elem := range n {
			n[i] = normalizeNumber(elem)
		}
		return n
	case map[string]any:
		for k, v := range n {
			n[k] = normalizeNumber(v)
		}
		return n
	}
	return v
}

// convertMapKeysRecursive recursively applies a key transformation function to all
// keys in a map[string]any and its nested structures. It also traverses lists containing maps.
// This is useful for converting between naming conventions like snake_case and dash-case.
func convertMapKeysRecursive(val any, convertFn func(string) string) any {
	switch v := val.(type) {
	case map[string]any:
		newMap := make(map[string]any, len(v))
		for key, val := range v {
			newMap[convertFn(key)] = convertMapKeysRecursive(val, convertFn)
		}
		return newMap
	case []any:
		for i, item := range v {
			v[i] = convertMapKeysRecursive(item, convertFn)
		}
		return v
	default:
		return val
	}
}

// underscoreToDash converts a string from snake_case to dash-case.
// Example: "start_at" → "start-at"
func underscoreToDash(s string) string {
	return strings.ReplaceAll(s, "_", "-")
}

// dashToUnderscore converts a string from dash-case to snake_case.
// Example: "start-at" → "start_at"
func dashToUnderscore(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}
