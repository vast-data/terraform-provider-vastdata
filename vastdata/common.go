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
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	planmodifiers "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	vast_client "github.com/vast-data/go-vast-client"
	"github.com/vast-data/go-vast-client/core"
	"github.com/vast-data/go-vast-client/resources/untyped"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/client"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/logging"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/schema_generation"
)

const (
	customRawKey = "@raw" // used to store raw string values in Record
)

var notImportable = false

// ProviderData contains the client and configuration passed from the provider
type ProviderData struct {
	Client      *VMSRest
	MigrateMode bool // When true, Create reads from cluster instead of creating, Delete is no-op
}

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
	isNotFoundErr  = vast_client.IsNotFoundErr
	ignoreNotFound = vast_client.IgnoreNotFound
	isApiError     = vast_client.IsApiError
)

// asApiError returns the *ApiError in err's chain, if any. API errors are
// often wrapped (for example by WaitAPICondition when an async task fails),
// so a direct err.(*ApiError) assertion can panic even when isApiError is true.
func asApiError(err error) (*ApiError, bool) {
	var apiErr *ApiError
	if errors.As(err, &apiErr) {
		return apiErr, true
	}
	return nil, false
}

// ignoreStatusCodes returns nil when err is an API error (possibly wrapped)
// with one of the given status codes, and err otherwise.
func ignoreStatusCodes(err error, codes ...int) error {
	if expectStatusCodes(err, codes...) {
		return nil
	}
	return err
}

// expectStatusCodes reports whether err is an API error (possibly wrapped)
// with one of the given status codes.
func expectStatusCodes(err error, codes ...int) bool {
	apiErr, ok := asApiError(err)
	if !ok {
		return false
	}
	for _, code := range codes {
		if apiErr.StatusCode == code {
			return true
		}
	}
	return false
}

type ForceCleanState struct{}

func (ForceCleanState) Error() string {
	return "force-clean-state"
}

type CustomImportOnly struct{}

func (CustomImportOnly) Error() string {
	return "custom-import-only"
}

// cached decorates fn so it runs once. The cache object lives in the returned
// closure and is named after fn; the name comes from the function itself.
func cached[T any](fn func() T) func() T {
	c := &onceCache[T]{name: callableName(fn)}
	if c.name == "" {
		panic("cached: cannot determine wrapped function name")
	}
	return func() T {
		c.once.Do(func() { c.value = fn() })
		return c.value
	}
}

// onceCache is the value closed over by cached.
type onceCache[T any] struct {
	name  string
	once  sync.Once
	value T
}

// callableName returns the wrapped function's name, without the package path.
func callableName(fn any) string {
	v := reflect.ValueOf(fn)
	if !v.IsValid() || v.Kind() != reflect.Func || v.IsNil() {
		return ""
	}
	pc := runtime.FuncForPC(v.Pointer())
	if pc == nil {
		return ""
	}
	name := pc.Name()
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[i+1:]
	}
	return name
}

// ignoreResourceGone treats client NotFoundError and HTTP 404 as success (resource already gone).
func ignoreResourceGone(err error) error {
	if err == nil || isNotFoundErr(err) {
		return nil
	}
	return ignoreStatusCodes(err, http.StatusNotFound)
}

func notFoundErr(resource string, query any) error {
	return &core.NotFoundError{Resource: resource, Query: fmt.Sprint(query)}
}

// sleepCtx waits for d or returns early if ctx is cancelled.
func sleepCtx(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// waitCondition is called once per wait cycle.
type waitCondition func(ctx context.Context) (done bool, err error)

type waitWithTimeoutConfig struct {
	// Timeout is the maximum wall-clock time to wait. Required when Cycles <= 0.
	Timeout time.Duration
	// Sleep is the delay after a false condition before the next cycle.
	// Defaults to 5s when <= 0.
	Sleep time.Duration
	// Cycles is the maximum number of condition evaluations. 0 means unlimited
	// (Timeout must be > 0).
	Cycles int
	// TimeoutMessage is included in the default timeout error when OnTimeout is nil.
	TimeoutMessage string
	// OnTimeout, if set, builds the error returned when Timeout/Cycles is exhausted.
	OnTimeout func() error
}

// errWaitTimeout marks exhaustion of waitWithTimeout's Timeout or Cycles budget.
var errWaitTimeout = errors.New("wait timed out")

// waitWithTimeout repeatedly invokes cond until it returns done=true, an error,
// or the configured Timeout/Cycles budget is exhausted. Sleep is applied between
// unsuccessful attempts (not after success or fatal error).
func waitWithTimeout(ctx context.Context, cfg waitWithTimeoutConfig, cond waitCondition) error {
	if cond == nil {
		return fmt.Errorf("waitWithTimeout: condition is required")
	}
	if cfg.Timeout <= 0 && cfg.Cycles <= 0 {
		return fmt.Errorf("waitWithTimeout: Timeout or Cycles is required")
	}
	if cfg.Sleep <= 0 {
		cfg.Sleep = 5 * time.Second
	}

	var deadline time.Time
	if cfg.Timeout > 0 {
		deadline = time.Now().Add(cfg.Timeout)
	}

	timeoutErr := func() error {
		if cfg.OnTimeout != nil {
			return cfg.OnTimeout()
		}
		if cfg.TimeoutMessage != "" {
			return fmt.Errorf("%w: %s", errWaitTimeout, cfg.TimeoutMessage)
		}
		if cfg.Cycles > 0 && cfg.Timeout <= 0 {
			return fmt.Errorf("%w: condition not met after %d cycles", errWaitTimeout, cfg.Cycles)
		}
		return fmt.Errorf("%w after %v", errWaitTimeout, cfg.Timeout)
	}

	for attempt := 1; ; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if !deadline.IsZero() && time.Now().After(deadline) {
			return timeoutErr()
		}
		if cfg.Cycles > 0 && attempt > cfg.Cycles {
			return timeoutErr()
		}

		done, err := cond(ctx)
		if err != nil {
			return err
		}
		if done {
			return nil
		}

		if cfg.Cycles > 0 && attempt >= cfg.Cycles {
			return timeoutErr()
		}
		if !deadline.IsZero() && time.Now().After(deadline) {
			return timeoutErr()
		}
		if err := sleepCtx(ctx, cfg.Sleep); err != nil {
			return err
		}
	}
}

// requirePathGUID returns a non-empty trimmed "guid" from tfstate for APIs that
// address resources as /{collection}/{guid}. resourceName is used in errors.
func requirePathGUID(tfstate is.ResourceState, resourceName string) (string, error) {
	if tfstate == nil || is.IsNilState(tfstate) || !tfstate.IsKnownAndNotNull("guid") {
		return "", fmt.Errorf("%s: guid is required", resourceName)
	}
	guid := strings.TrimSpace(tfstate.String("guid"))
	if guid == "" {
		return "", fmt.Errorf("%s: guid is empty", resourceName)
	}
	return guid, nil
}

// pathGUIDOrPlan returns guid from state, falling back to plan (create→update edge).
func pathGUIDOrPlan(state, plan is.ResourceState, resourceName string) (string, error) {
	guid, err := requirePathGUID(state, resourceName)
	if err == nil {
		return guid, nil
	}
	if plan != nil && !is.IsNilState(plan) {
		if g, planErr := requirePathGUID(plan, resourceName); planErr == nil {
			return g, nil
		}
	}
	return "", err
}

// nestedMap walks nested map/Record keys and returns the final map, or nil.
func nestedMap(rec Record, keys ...string) map[string]any {
	cur := any(rec)
	for _, key := range keys {
		m, ok := cur.(map[string]any)
		if !ok {
			if r, ok := cur.(Record); ok {
				m = map[string]any(r)
			} else {
				return nil
			}
		}
		cur = m[key]
	}
	if m, ok := cur.(map[string]any); ok {
		return m
	}
	if r, ok := cur.(Record); ok {
		return map[string]any(r)
	}
	return nil
}

// nestedString walks nested map/Record keys and returns a trimmed string, or "".
func nestedString(rec Record, keys ...string) string {
	cur := any(rec)
	for _, key := range keys {
		m, ok := cur.(map[string]any)
		if !ok {
			if r, ok := cur.(Record); ok {
				m = map[string]any(r)
			} else {
				return ""
			}
		}
		cur = m[key]
	}
	if cur == nil {
		return ""
	}
	s := strings.TrimSpace(fmt.Sprint(cur))
	if s == "<nil>" {
		return ""
	}
	return s
}

// withTransformedQueryParams applies TransformQueryParams when manager implements
// it (Create/Read/Update/Delete). Non-implementers keep query as-is (nil → empty).
func withTransformedQueryParams(ctx context.Context, manager any, managerName string, query params, rest *VMSRest) params {
	if query == nil {
		query = params{}
	}
	if t, ok := manager.(TransformQueryParams); ok {
		logging.Debug(ctx, fmt.Sprintf("TransformQueryParams[%s]: do.", managerName))
		return t.TransformQueryParams(query, rest)
	}
	return query
}

// filterRecordSetByAttrs returns records whose attributes match all entries in want.
// An empty want returns records unchanged.
func filterRecordSetByAttrs(records RecordSet, want map[string]any) RecordSet {
	if len(want) == 0 {
		return records
	}
	out := make(RecordSet, 0, len(records))
	for _, rec := range records {
		if recordMatchesAttrs(rec, want) {
			out = append(out, rec)
		}
	}
	return out
}

func recordMatchesAttrs(rec Record, want map[string]any) bool {
	for k, v := range want {
		rv, ok := rec[k]
		if !ok || fmt.Sprint(rv) != fmt.Sprint(v) {
			return false
		}
	}
	return true
}

// PopulateIDFieldsFromNestedObjects examines all keys in tfstate that end with "_id"
// and attempts to populate them from nested objects in the record.
//
// For example:
//   - If tfstate has key "local_provider_id" but record doesn't have it
//   - Look for record["local_provider"] (object with "id" and "name")
//   - Extract record["local_provider"]["id"] and set record["local_provider_id"]
//
// This handles the common API pattern where responses contain nested objects
// like {tenant: {id: 1, name: "foo"}} but the schema expects tenant_id: 1
func PopulateIDFieldsFromNestedObjects(ctx context.Context, tfstate is.ResourceState, record Record) {
	if record == nil || tfstate == nil {
		return
	}

	// Iterate over all tfstate keys
	for key := range tfstate.RawValues() {
		// Check if key ends with "_id"
		if !strings.HasSuffix(key, "_id") {
			continue
		}

		// Check if this field already exists in the record
		if _, exists := record[key]; exists {
			continue
		}

		// Derive the parent object name by removing "_id" suffix
		// e.g., "local_provider_id" -> "local_provider"
		parentKey := strings.TrimSuffix(key, "_id")

		// Check if parent object exists in record
		parentValue, parentExists := record[parentKey]
		if !parentExists || parentValue == nil {
			continue
		}

		// Try to cast parent to map[string]any
		parentMap, ok := parentValue.(map[string]any)
		if !ok {
			continue
		}

		// Try to extract "id" from the parent object
		if id, idExists := parentMap["id"]; idExists && id != nil {
			record[key] = id
			logging.Debug(ctx, fmt.Sprintf(
				"PopulateIDFieldsFromNestedObjects: extracted %s=%v from %s.id",
				key, id, parentKey,
			))
		}
	}
}

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

// Common modifiers
const (
	ModifierForceNew             = schema_generation.ModifierForceNew
	ModifierForceNewIfConfigured = schema_generation.ModifierForceNewIfConfigured
)

// ----------------------------------
// MISC
// ----------------------------------

// isRedacted reports whether v is an API placeholder for secret material
// (e.g. "******") rather than a real value. Use when GET responses redact
// fields that must not be written into Terraform state.
func isRedacted(v any) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	switch strings.TrimSpace(s) {
	case "******", "********", "[REDACTED]", "redacted":
		return true
	default:
		return false
	}
}

func withContext(ctx context.Context, method string, managerName string, fn func(ctx context.Context)) {
	ctx = client.ContextWithRequestID(ctx)
	logging.Debug(ctx, fmt.Sprintf("◉ %s[%s] start", method, managerName))
	defer logging.Debug(ctx, fmt.Sprintf("◉ %s[%s] end", method, managerName))
	fn(ctx)
}

// isZeroValue reports whether v is the zero value for its type
// (e.g. "" for strings, 0 for numbers, false for bools, nil pointer, empty slice/map).
func isZeroValue(v any) bool {
	if v == nil {
		return true
	}
	return reflect.ValueOf(v).IsZero()
}

// deleteZeroValues removes entries from params whose keys are listed in
// fields and whose values are the zero/empty value for their type
// (e.g. "" for strings, 0 for numbers, false for bools, empty slices/maps).
func deleteZeroValues(params map[string]any, fields []string) {
	for _, key := range fields {
		if val, ok := params[key]; ok && isZeroValue(val) {
			delete(params, key)
		}
	}
}

// ----------------------------------
// Common Resource Validators
// ----------------------------------

// validateOneOf ensures that at most one of the given fields is set (i.e., known and not null).
// Returns an error if more than one field is set.
func validateOneOf(tf is.ResourceState, fields ...string) error {
	var setFields []string
	var fieldDeclared bool

	for _, field := range fields {
		if tf.IsKnownAndNotNull(field) {
			fieldDeclared = true
			setFields = append(setFields, field)
		}
	}
	if !fieldDeclared {
		return fmt.Errorf("one of %q must be set", fields)
	}

	if len(setFields) > 1 {
		return fmt.Errorf("only one of %q can be set, but multiple were provided: %v", fields, setFields)
	}

	return nil
}

// validateAllOf ensures that all the given fields are set (i.e., known and not null).
// Returns an error listing the fields that are missing.
func validateAllOf(tf is.ResourceState, fields ...string) error {
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
func validateNoneOf(tf is.ResourceState, fields ...string) error {
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

// ValidateFieldIsOneOf ensures that the specified field in allows only a limited set of valid values.
func ValidateFieldIsOneOf[T string | int64 | float64](tfState is.ResourceState, field string, validValues ...T) error {
	if !tfState.IsKnownAndNotNull(field) {
		return nil // Field is not set, no validation needed
	}

	var value T
	switch any(value).(type) {
	case string:
		value = any(tfState.String(field)).(T)
	case int64:
		value = any(tfState.Int64(field)).(T)
	case float64:
		value = any(tfState.Float64(field)).(T)
	default:
		return fmt.Errorf("unsupported type for validation: %T", value)
	}

	if slices.Contains(validValues, value) {
		return nil // Valid value found
	}

	return fmt.Errorf("field %q must be one of %v, but got '%v'", field, validValues, value)
}

// ValidateFieldIsNoneOf ensures that the specified field in allows none of the invalid values.
func ValidateFieldIsNoneOf[T string | int64 | float64](tfState is.ResourceState, field string, invalidValues ...T) error {
	if !tfState.IsKnownAndNotNull(field) {
		return nil // Field is not set, no validation needed
	}

	var value T
	switch any(value).(type) {
	case string:
		value = any(tfState.String(field)).(T)
	case int64:
		value = any(tfState.Int64(field)).(T)
	case float64:
		value = any(tfState.Float64(field)).(T)
	default:
		return fmt.Errorf("unsupported type for validation: %T", value)
	}

	if slices.Contains(invalidValues, value) {
		return fmt.Errorf("field %q must not be one of %v, but got '%v'", field, invalidValues, value)
	}

	return nil
}

func ensureNotChanged(tfState is.ResourceState, planTfState is.ResourceState, fields ...string) error {
	for _, field := range fields {
		if planTfState.IsUnknown(field) {
			continue
		}
		if !tfState.Get(field).Equal(planTfState.Get(field)) {
			return fmt.Errorf(
				"field %q is immutable, current value: %v, planned value: %v",
				field,
				tfState.Get(field),
				planTfState.Get(field),
			)
		}
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
func getSearchParams(ctx context.Context, tfState, planTfState is.ResourceState) params {
	searchParams := tfState.GetGenericSearchParams(ctx)
	if !is.IsNilState(planTfState) {
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
// pathIDField returns the record attribute used as the URL path id for Get/Update/Delete.
// Empty PathIDField means the default VMS "id" field.
func pathIDField(hints *is.TFStateHints) string {
	if hints != nil && hints.PathIDField != "" {
		return hints.PathIDField
	}
	return "id"
}

func getRecordBySearchParams(ctx context.Context, api VastResourceAPIWithContext, tfState is.ResourceState, planTfState is.ResourceState, managerName, op string) (DisplayableRecord, error) {
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

	idField := pathIDField(tfState.StateHints())
	pathID, pathIDExists := searchParams[idField]
	guid, guidExists := searchParams["guid"]
	// When PathIDField is "guid", do not fall back to numeric "id" for GetById —
	// DataEngine URLs expect {guid} and numeric ids are precision-unsafe.
	useDefaultIDPath := idField == "id"

	// Query keys that belong on GetById (e.g. ?tenant_name=), not in the path.
	idQuery := params{}
	if tn, ok := searchParams["tenant_name"]; ok {
		idQuery["tenant_name"] = tn
	}

	// Attempt 1: get resource by path id (id or guid).
	if pathIDExists {
		logging.Debug(ctx, fmt.Sprintf("%s[%s]: found %s = %v.", op, managerName, idField, pathID))
		if len(idQuery) > 0 {
			record, err = api.GetByIdWithContext(ctx, pathID, idQuery)
		} else {
			record, err = api.GetByIdWithContext(ctx, pathID)
		}
	}
	// Attempt 2: get resource by GUID (VMS-style fallback only).
	// Listed: some endpoints ignore the guid query parameter and
	// return every record, which Get rejects as "too many records". The guid is compared
	// here, and a list with no match is a 404 so the caller treats it as not found.
	if useDefaultIDPath && (!pathIDExists || expectStatusCodes(err, http.StatusNotFound)) && guidExists {
		logging.Debug(ctx, fmt.Sprintf("%s[%s]: found GUID = %v.", op, managerName, guid))
		var records RecordSet
		records, err = api.ListWithContext(ctx, params{"guid": guid})
		if err == nil {
			record, err = recordByGUID(records, guid)
		}
	}
	// Attempt 3: if path id (and guid for default id path) are not set or prior reads failed, use search parameters.
	triedDirect := pathIDExists || (useDefaultIDPath && guidExists)
	if !triedDirect || expectStatusCodes(err, http.StatusNotFound) {
		searchParams.Without("id", "guid")
		if len(searchParams) > 0 {
			if !triedDirect {
				logging.Debug(
					ctx,
					fmt.Sprintf("%s[%s]: no path id set, using search parameters %v.",
						op, managerName, searchParams,
					),
				)
			} else {
				logging.Debug(
					ctx,
					fmt.Sprintf("%s[%s]: %s %v not found, using search parameters %v.",
						op, managerName, idField, pathID, searchParams,
					),
				)
			}
			record, err = api.GetWithContext(ctx, searchParams)
		}
	}

	return record, err

}

// recordByGUID returns the record whose guid equals want. No match is a 404 ApiError, which
// the callers of getRecordBySearchParams already treat as "not found".
func recordByGUID(records RecordSet, want any) (Record, error) {
	wanted := fmt.Sprint(want)
	for _, rec := range records {
		if rec == nil {
			continue
		}
		got, ok := rec["guid"]
		if ok && fmt.Sprint(got) == wanted {
			return rec, nil
		}
	}
	return nil, &ApiError{
		StatusCode: http.StatusNotFound,
		Body:       fmt.Sprintf("no record with guid %s", wanted),
	}
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
func deleteRecordBySearchParams(ctx context.Context, api VastResourceAPIWithContext, tfState is.ResourceState, managerName, op string) (vast_client.Record, error) {
	var err error
	var result vast_client.Record
	searchParams := getSearchParams(ctx, tfState, nil)
	if len(searchParams) == 0 {
		return nil, fmt.Errorf(
			"%s[%s]: no search parameters provided for %q resource."+
				" Verify presence of required fields or add searchable hints to resource",
			op,
			managerName,
			managerName,
		)
	}

	var deleteBodyParams map[string]any
	var deleteQueryParams map[string]any
	if len(tfState.StateHints().DeleteOnlyBodyFields) > 0 {
		deleteBodyParams = tfState.GetDeleteOnlyBodyParams()
	}
	if len(tfState.StateHints().DeleteOnlyParamFields) > 0 {
		deleteQueryParams = tfState.GetDeleteOnlyQueryParams()
	}

	idField := pathIDField(tfState.StateHints())
	if tn, ok := searchParams["tenant_name"]; ok {
		if deleteQueryParams == nil {
			deleteQueryParams = params{}
		}
		q := params(deleteQueryParams)
		q.Update(params{"tenant_name": tn}, true)
		deleteQueryParams = q
	}
	if id, ok := searchParams[idField]; ok {
		logging.Debug(ctx, fmt.Sprintf("%s[%s]: found %s = %v.", op, managerName, idField, id))
		result, err = api.DeleteByIdWithContext(ctx, id, deleteQueryParams, deleteBodyParams)
	} else {
		logging.Debug(ctx, fmt.Sprintf("%s[%s]: no %s found, using search params.", op, managerName, idField))
		result, err = api.DeleteWithContext(ctx, searchParams, deleteQueryParams, deleteBodyParams)
	}
	return result, ignoreStatusCodes(err, http.StatusNotFound)

}

// ----------------------------------
// Transformations
// ----------------------------------

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

// parseImportId parses the import ID into the TFState attributes.
func parseImportId(importID string, tfState is.ResourceState) error {
	// Use default import implementation
	hints := tfState.StateHints()

	if strings.Contains(importID, "=") {
		// Parse key=value pairs regardless of hints
		sep := ","
		if strings.Contains(importID, ";") {
			sep = ";"
		}
		parts := strings.Split(importID, sep)
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			kv := strings.SplitN(p, "=", 2)

			if len(kv) != 2 {
				return fmt.Errorf("segment %q is not in key=value form", p)
			}
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])
			if !tfState.HasAttribute(key) {
				return fmt.Errorf("field %q is not present in the resource schema", key)
			}
			t := tfState.Type(key)
			switch {
			case t.Equal(types.Int64Type):
				n, convErr := strconv.ParseInt(val, 10, 64)
				if convErr != nil {
					return fmt.Errorf("field %q contains invalid value %q: %w", key, val, convErr)
				}
				tfState.SetOrAdd(key, types.Int64Value(n))
			case t.Equal(types.BoolType):
				bv := strings.EqualFold(val, "true") || val == "1"
				tfState.SetOrAdd(key, types.BoolValue(bv))
			case t.Equal(types.StringType):
				tfState.SetOrAdd(key, types.StringValue(val))
			default:
				// store as string for unsupported types
				tfState.SetOrAdd(key, types.StringValue(val))
			}
		}
	} else if hints != nil && len(hints.ImportFields) > 0 && strings.Contains(importID, "|") {
		// Ordered values mode via hints
		if err := parseAndApplyCompositeImport(importID, hints.ImportFields, tfState, func(k string, v attr.Value) {
			tfState.SetOrAdd(k, v)
		}); err != nil {
			return err
		}
	} else {
		// Treat as single ID token
		idField := "id"
		if !tfState.HasAttribute(idField) {
			return fmt.Errorf("field %q is not present in the resource schema", idField)
		}
		idType := tfState.Type(idField)
		if idType.Equal(types.Int64Type) {
			idInt64, convErr := strconv.ParseInt(importID, 10, 64)
			if convErr != nil {
				return fmt.Errorf("field %q contains invalid value %q: %w", idField, importID, convErr)
			}
			tfState.SetOrAdd(idField, types.Int64Value(idInt64))
		} else if idType.Equal(types.StringType) {
			tfState.SetOrAdd(idField, types.StringValue(importID))
		} else if idType.Equal(types.BoolType) {
			bv := strings.EqualFold(importID, "true") || importID == "1"
			tfState.SetOrAdd(idField, types.BoolValue(bv))
		} else {
			return fmt.Errorf("field %q is not present in the resource schema", idField)
		}
	}
	return nil
}

// deleteDuplicates removes duplicate elements from a slice, preserving order.
// Works with any comparable type (string, int, float64, etc.).
func deleteDuplicates[T comparable](s []T) []T {
	seen := make(map[T]struct{}, len(s))
	result := make([]T, 0, len(s))
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}

// allPrimitiveScalars reports whether every element is nil or a scalar
// (string, bool, int64, float64). Used when normalizing list comparisons.
func allPrimitiveScalars(vals []any) bool {
	for _, v := range vals {
		switch v.(type) {
		case nil, string, bool, int64, float64:
		default:
			return false
		}
	}
	return true
}

// ----------------------------------
// Write-only field fingerprints
// ----------------------------------
//
// Helpers for tracking changes to write-only / sensitive string fields that
// Terraform does not store in state. Pattern:
//
//  1. Keep the secret as WriteOnly (+ Sensitive).
//  2. Add a computed "<field>_md5" attribute (see ContentMD5SchemaAttribute).
//  3. List the field in WriteOnlyMD5Fields — ResourceLifecycle.PrepareDesiredState
//     seeds <field>_md5 via ApplyWriteOnlyMD5FromHints (TF ModifyPlan and Crossplane).
//  4. For non-ContentMD5 fingerprints (e.g. DER cert md5), implement PrepareDesiredState.
//  5. On read/transform, call SetRecordContentMD5 when the API still returns the
//     secret (or clear the md5 when it does not) so state stays in sync.
//
// Used by vastdata_license (replace on md5 change) and
// vastdata_administrator_manager (update on md5 change via WriteOnlyMD5Fields).

// ContentMD5Hex returns the lowercase hex MD5 of s.
//
//nolint:gosec // plan/state equality fingerprint, not a security control
func ContentMD5Hex(s string) string {
	return BytesMD5Hex([]byte(s))
}

// BytesMD5Hex returns the lowercase hex MD5 of b.
//
//nolint:gosec // plan/state equality fingerprint, not a security control
func BytesMD5Hex(b []byte) string {
	sum := md5.Sum(b)
	return hex.EncodeToString(sum[:])
}

// ContentMD5SchemaAttribute builds a computed string attribute that stores an
// MD5 fingerprint of a write-only field for change detection.
// If requiresReplace is true, a fingerprint change forces resource replace
// (use when the API has no PATCH for that field).
func ContentMD5SchemaAttribute(sourceField string, requiresReplace bool) rschema.StringAttribute {
	mods := []planmodifiers.String{
		stringplanmodifier.UseStateForUnknown(),
	}
	if requiresReplace {
		mods = append([]planmodifiers.String{stringplanmodifier.RequiresReplace()}, mods...)
	}
	return rschema.StringAttribute{
		Computed: true,
		Description: fmt.Sprintf(
			"MD5 of %s (lowercase hex). Computed from the configured value at plan time so changes to the write-only field are detected. Not settable.",
			sourceField,
		),
		PlanModifiers: mods,
	}
}

// SetStateContentMD5 sets md5Field on st from ContentMD5Hex(sourceField).
// Empty/null source clears md5Field (TERF-320). Shared by TF and Crossplane.
func SetStateContentMD5(st is.ResourceState, sourceField, md5Field string) {
	if st == nil || !st.HasAttribute(md5Field) {
		return
	}
	raw := st.RawValues()
	v, ok := raw[sourceField]
	if !ok || v == nil || v.IsNull() || v.IsUnknown() {
		raw[md5Field] = types.StringNull()
		return
	}
	s, ok := v.(types.String)
	if !ok || s.ValueString() == "" {
		raw[md5Field] = types.StringNull()
		return
	}
	st.Set(md5Field, ContentMD5Hex(s.ValueString()))
}

// ApplyWriteOnlyMD5FromHints seeds every WriteOnlyMD5Fields companion <field>_md5
// from the secret currently on st (ContentMD5Hex). Custom PrepareDesiredState
// hooks may overwrite specific companions afterward (e.g. DER cert md5).
func ApplyWriteOnlyMD5FromHints(st is.ResourceState) {
	if st == nil {
		return
	}
	hints := st.StateHints()
	if hints == nil || len(hints.WriteOnlyMD5Fields) == 0 {
		return
	}
	for _, field := range hints.WriteOnlyMD5Fields {
		SetStateContentMD5(st, field, field+"_md5")
	}
}

// SetPlanContentMD5 reads sourceField from config and sets md5Field on the plan
// to ContentMD5Hex(value). Prefer ApplyWriteOnlyMD5FromHints / PrepareDesiredState;
// kept for tests and any remaining direct callers.
func SetPlanContentMD5(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
	sourceField, md5Field string,
) {
	if req.Plan.Raw.IsNull() {
		return
	}
	var src types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root(sourceField), &src)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if src.IsUnknown() {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root(md5Field), types.StringUnknown())...)
		return
	}
	if src.IsNull() || src.ValueString() == "" {
		// Clear fingerprint left by UseStateForUnknown when the write-only
		// optional secret is omitted from config.
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root(md5Field), types.StringNull())...)
		return
	}
	resp.Diagnostics.Append(resp.Plan.SetAttribute(
		ctx,
		path.Root(md5Field),
		types.StringValue(ContentMD5Hex(src.ValueString())),
	)...)
}

// SetRecordContentMD5 sets record[md5Field] to ContentMD5Hex of record[sourceField]
// when sourceField is a non-empty string. Leaves md5Field untouched otherwise.
func SetRecordContentMD5(record Record, sourceField, md5Field string) {
	if record == nil {
		return
	}
	if key, ok := record[sourceField].(string); ok && key != "" {
		record[md5Field] = ContentMD5Hex(key)
	}
}

// SetPlanContentMD5Fields calls SetPlanContentMD5 for each source→md5 pair.
// pairs must be even: source1, md5_1, source2, md5_2, ...
func SetPlanContentMD5Fields(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
	pairs ...string,
) {
	for i := 0; i+1 < len(pairs); i += 2 {
		SetPlanContentMD5(ctx, req, resp, pairs[i], pairs[i+1])
	}
}

// DeleteRecordFields removes keys from an API record before it is written to state.
func DeleteRecordFields(record Record, fields ...string) {
	if record == nil {
		return
	}
	for _, f := range fields {
		delete(record, f)
	}
}

// ----------------------------------
// Async tasks
// ----------------------------------

func handleMaybeAsyncTask(ctx context.Context, rest *VMSRest, record Record, timeout *time.Duration) error {
	waitTimeout := 10 * time.Minute
	if timeout != nil {
		waitTimeout = *timeout
	}
	_, err := untyped.MaybeWaitAsyncResultWithContext(ctx, record, rest, waitTimeout)
	return err
}

func resolveRecordAfterAsyncTask(ctx context.Context, manager ResourceManager, rest *VMSRest, record Record, managerName string) (Record, error) {
	resolver, ok := manager.(ResolveRecordAfterAsyncTask)
	if !ok {
		return record, nil
	}
	logging.Debug(ctx, fmt.Sprintf("ResolveRecordAfterAsyncTask[%s]: do.", managerName))
	return resolver.ResolveRecordAfterAsyncTask(ctx, rest, record)
}

const (
	defaultRetryTimes        = 5
	defaultRetrySleepSeconds = 10
)

// ---------------------------------------
//      HELPER FUNCTIONS
// ---------------------------------------

// parseAndApplyCompositeImport parses a composite import string into values for the given fields
// and invokes set(key, value) with properly typed Terraform attr values.
// Accepted formats:
//   - key=value pairs separated by ',' (e.g., "gid=1001,tenant_id=22,context=ad")
//   - ordered values separated by '|' matching fields order
func parseAndApplyCompositeImport(importID string, fields []string, tfState is.ResourceState, set func(string, attr.Value)) error {
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
		rest                         = r.providerData.Client
		managerName                  = r.managerName
		tfState                      = manager.TfState()
		planTfState is.ResourceState = nil
		api                          = manager.API(rest)
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
		visualization := is.BuildResourceAttributesString(tfState.SchemaDefinition().(rschema.Schema).Attributes, false, 2)
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

// WithRetry wraps fn with retry logic driven by expr.
// operation is used for log messages (e.g. "CreateResource", "DeleteResource").
//
// A retry is attempted when:
//  1. The call returns an *ApiError whose StatusCode is listed in expr.StatusCodes, AND
//  2. If expr.BodyContains is non-empty, at least one substring is present in the response body.
//
// Any other error type is returned immediately without retrying.
func WithRetry[T any](
	expr *is.RetryExpression,
	operation, managerName string,
	fn func(context.Context, *VMSRest) (T, error),
) func(context.Context, *VMSRest) (T, error) {
	return func(ctx context.Context, rest *VMSRest) (T, error) {
		return retryOnExpression(ctx, expr, operation, managerName, func() (T, error) {
			return fn(ctx, rest)
		})
	}
}

// retryOnExpression retries fn when it returns an error matching expr.
func retryOnExpression[T any](
	ctx context.Context,
	expr *is.RetryExpression,
	operation string,
	scope string,
	fn func() (T, error),
) (T, error) {
	var zero T
	maxAttempts := expr.Times
	if maxAttempts <= 0 {
		maxAttempts = defaultRetryTimes
	}
	sleepSecs := expr.SleepSeconds
	if sleepSecs <= 0 {
		sleepSecs = defaultRetrySleepSeconds
	}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		result, err := fn()
		if err == nil {
			if attempt > 1 {
				logging.Info(ctx, fmt.Sprintf(
					"%s[%s]: succeeded on attempt %d/%d.",
					operation, scope, attempt, maxAttempts,
				))
			}
			return result, nil
		}

		lastErr = err
		if !shouldRetry(expr, err) {
			return zero, err
		}

		logging.Warn(ctx, fmt.Sprintf(
			"%s[%s]: attempt %d/%d failed (%s), retrying in %ds...",
			operation, scope, attempt, maxAttempts, err.Error(), sleepSecs,
		))

		if attempt < maxAttempts {
			time.Sleep(time.Duration(sleepSecs) * time.Second)
		}
	}

	return zero, fmt.Errorf(
		"%s[%s]: all %d attempts failed, last error: %w",
		operation, scope, maxAttempts, lastErr,
	)
}

// shouldRetry returns true when err satisfies the retry conditions defined in expr.
func shouldRetry(expr *is.RetryExpression, err error) bool {
	if !expectStatusCodes(err, expr.StatusCodes...) {
		return false
	}

	if len(expr.BodyContains) == 0 {
		return true
	}
	var apiErr *ApiError
	if !errors.As(err, &apiErr) {
		return true
	}
	for _, pattern := range expr.BodyContains {
		if strings.Contains(apiErr.Body, pattern) {
			return true
		}
	}
	return false
}

// mergeSubResources calls GetSubResources (if implemented) and merges the
// returned Record into the main record before FillFromRecord is called.
// This keeps sub-resource data out of tfstate manipulation: the merged record
// is the single source used to populate state.
// manager may be a ResourceManager or a DataSourceManager — only the
// GetSubResources interface is required.
//
// If at least one sub-resource hint declares MinVastVersion, the cluster version
// is retrieved once (using the process-level cache) and passed to GetSubResources.
// The implementation is then responsible for using clusterVersion to gate fetches.
// When no hint uses MinVastVersion, nil is passed — implementations must tolerate nil.
type hasTFState interface {
	TfState() is.ResourceState
}

func mergeSubResources(ctx context.Context, manager any, rest *VMSRest, record DisplayableRecord, managerName string) error {
	if record == nil {
		return nil
	}

	var hints *is.TFStateHints
	if hts, ok := manager.(hasTFState); ok {
		hints = hts.TfState().StateHints()
		if hints == nil || len(hints.SubResources) == 0 {
			return nil
		}
	}

	imp, ok := manager.(GetSubResources)
	if !ok {
		return nil
	}

	// Fetch cluster version exactly once if any hint requires it.
	// The result is passed into GetSubResources so implementations
	// don't need to call GetCachedClusterVersion themselves.
	var clusterVersion *version.Version
	if hints != nil {
		for _, sr := range hints.SubResources {
			if sr.MinVastVersion != nil {
				v, err := GetCachedClusterVersion(ctx, rest)
				if err != nil {
					logging.Warn(ctx, fmt.Sprintf(
						"mergeSubResources[%s]: cannot retrieve cluster version for MinVastVersion check: %v",
						managerName, err,
					))
				} else {
					clusterVersion = v
				}
				break
			}
		}
	}

	logging.Debug(ctx, fmt.Sprintf("GetSubResources[%s]: do.", managerName))
	sub, err := imp.GetSubResources(ctx, rest, record.(Record), clusterVersion)
	if err != nil {
		return err
	}
	if sub != nil {
		for k, v := range sub {
			record.(Record)[k] = v
		}
		return nil
	}

	// GetSubResources returned nil — the trigger is off or fetch was skipped.
	if rm, ok := manager.(ResourceManager); ok {
		if h := rm.TfState().StateHints(); h != nil {
			for _, sr := range h.SubResources {
				if sr.SchemaKey != "" {
					// Nested sub-resource: null the top-level key.
					record.(Record)[sr.SchemaKey] = nil
				} else {
					// Flat sub-resource: null each declared attribute individually.
					for k := range sr.SchemaAttributes {
						record.(Record)[k] = nil
					}
				}
			}
		}
	}
	return nil
}
