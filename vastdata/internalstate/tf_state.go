// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	vast_client "github.com/vast-data/go-vast-client"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/logging"

	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// FieldFilterFlag defines filtering criteria used during schema introspection.
// These flags are used to identify fields with specific traits, such as being
// required, optional, sensitive, or marked as searchable.
type FieldFilterFlag int

const (
	SearchRequired FieldFilterFlag = iota
	SearchOptional
	SearchSensitive
	SearchSearchable
	SearchComputed
	SearchWriteOnly
	SearchReadonly // Indicates that field is only for read operations, not for create/update
	SearchNotRequired
	SearchNotOptional
	SearchNotSensitive
	SearchNotSearchable
	SearchNotWriteOnly
	SearchEmpty
	SearchPrimitivesOnly // Indicates that only primitive types (string, int, bool) should be considered
)

// CommonSearchableFields defines a set of standard fields that are commonly
// searchable across most VAST API resources. These fields are typically used
// in query parameters.
// But using just one of them sometimes is not enough to identify a record.
// For instance many subsystems (views with BLOCK protocol) can have the same name but different tenants.
var CommonSearchableFields = []string{
	"name", "path", "bucket", "gid", "uid",
}

// AdditionalSearchableFields If exists these fields are also used for search.
// Should be always added to search fields
var AdditionalSearchableFields = []string{
	"guid", "id", "tenant_id", "tenant_name",
}

type FilterTagCombination int

const (
	FilterOr FilterTagCombination = iota
	FilterAnd
)

type FieldSet struct {
	Include []string // fields to include (whitelist)
	Exclude []string // fields to exclude (blacklist)
}

func (fs *FieldSet) shouldInclude(field string) bool {
	if slices.Contains(fs.Exclude, field) {
		return false
	}
	if len(fs.Include) > 0 {
		return slices.Contains(fs.Include, field)
	}
	return true
}

type attrMeta struct {
	Required   bool
	Optional   bool
	Computed   bool
	Sensitive  bool
	Searchable bool
	WriteOnly  bool
	ReadOnly   bool // Indicates that field is only for read operations, not for create/update
	EditOnly   bool // Indicates that field is only for edit operations, not for create/read
	DeleteOnly bool // Indicates that field is only for delete operations, not for create/read/edit
}

// satisfyFieldFilterFlag returns true if the attribute matches the given flags,
// based on the combination mode (AND or OR).
func (meta attrMeta) satisfyFieldFilterFlag(comb FilterTagCombination, flags ...FieldFilterFlag) bool {
	matches := func(f FieldFilterFlag) bool {
		switch f {
		case SearchRequired:
			return meta.Required
		case SearchOptional:
			return meta.Optional
		case SearchSensitive:
			return meta.Sensitive
		case SearchComputed:
			return meta.Computed
		case SearchWriteOnly:
			return meta.WriteOnly
		case SearchReadonly:
			return meta.ReadOnly
		case SearchSearchable:
			return meta.Searchable
		case SearchNotRequired:
			return !meta.Required
		case SearchNotOptional:
			return !meta.Optional
		case SearchNotSensitive:
			return !meta.Sensitive
		case SearchNotWriteOnly:
			return !meta.WriteOnly
		case SearchNotSearchable:
			return !meta.Searchable
		case SearchEmpty:
			return false
		case SearchPrimitivesOnly:
			// This flag is enforced via type checks in callers (e.g., GetFilteredValues).
			// It has no bearing on attribute metadata, so treat as neutral true here.
			return true
		default:
			panic(fmt.Sprintf("unknown search context flag: %d", f))
		}
	}

	switch comb {
	case FilterOr:
		return slices.ContainsFunc(flags, matches)

	case FilterAnd:
		for _, f := range flags {
			if !matches(f) {
				return false
			}
		}
		return true

	default:
		panic(fmt.Sprintf("unknown filter combination mode: %d", comb))
	}
}

type TFState struct {
	Raw     map[string]attr.Value
	Schema  any           // datasource_schema.Schema or resource_schema.Schema
	Kind    SchemaContext // indicates schema kind
	Meta    map[string]attrMeta
	TypeMap map[string]attr.Type // for schema introspection
	Hints   *TFStateHints
	Enabled bool
}

func NewTFState(raw map[string]attr.Value, schema any, kind SchemaContext, hints *TFStateHints) (*TFState, error) {
	if raw == nil {
		raw = make(map[string]attr.Value)
	}
	if schema == nil {
		return &TFState{
			Raw:     raw,
			Schema:  nil,
			Kind:    kind,
			Meta:    map[string]attrMeta{},
			Hints:   hints,
			Enabled: false,
		}, nil
	}

	meta, err := extractMetaFromSchema(schema, kind, hints)
	if err != nil {
		return nil, err
	}
	typeMap, err := extractTypesFromSchema(schema, kind)
	if err != nil {
		return nil, fmt.Errorf("extract types from schema: %w", err)
	}
	return &TFState{
		Raw:     raw,
		Schema:  schema,
		Kind:    kind,
		Meta:    meta,
		TypeMap: typeMap,
		Hints:   hints,
		Enabled: true,
	}, nil
}

// NewTFStateFromSchema builds a TFState detecting the schema kind (resource/datasource)
// from the schema type. It reports problems as errors, which makes it suitable for
// StateFactory implementations; NewTFStateMust is the panicking variant.
func NewTFStateFromSchema(raw map[string]attr.Value, schema any, hints *TFStateHints) (*TFState, error) {
	kind := SchemaTypeUnknown
	switch schema.(type) {
	case nil:
	case rschema.Schema:
		kind = SchemaForResource
	case dsschema.Schema:
		kind = SchemaForDataSource
	default:
		return nil, fmt.Errorf("unknown schema type: %T", schema)
	}
	return NewTFState(raw, schema, kind, hints)
}

// NewTFStateMust is NewTFStateFromSchema that panics on error.
func NewTFStateMust(raw map[string]attr.Value, schema any, hints *TFStateHints) *TFState {
	return Must(NewTFStateFromSchema(raw, schema, hints))
}

func (s *TFState) assertEnabled() {
	if !s.Enabled {
		panic("TFState is disabled; schema and meta not available")
	}
}

// Copy creates a shallow copy of the TFState.
// - The Raw map is deeply copied (per key).
// - Schema, Kind, Meta, TypeMap, Hints are reused by reference.
// This is useful when you want to modify Raw independently but keep schema and metadata intact.
func (s *TFState) Copy() ResourceState {
	return s.CopyTFState()
}

// CopyTFState is the concrete-typed implementation of Copy. Implementations that embed
// *TFState use it to re-wrap the copy in their own type when overriding Copy.
func (s *TFState) CopyTFState() *TFState {
	s.assertEnabled()

	rawCopy := make(map[string]attr.Value, len(s.Raw))
	maps.Copy(rawCopy, s.Raw)

	return &TFState{
		Raw:     rawCopy,
		Schema:  s.Schema,
		Kind:    s.Kind,
		Meta:    s.Meta,
		TypeMap: s.TypeMap,
		Hints:   s.Hints,
		Enabled: s.Enabled,
	}
}

func (s *TFState) Pretty() string {
	s.assertEnabled()
	return prettyWithMeta(s.Raw, s.Meta)
}

func (s *TFState) convPanic(msg string) {
	panic(fmt.Sprintf(
		"conversion error: %s\n\nSchema:\n%s",
		msg,
		schemaVisualization(s.Schema, s.Kind, false),
	))
}

// --- Helper Functions ---

// HasAttribute checks if a schema has a specific attribute
func (s *TFState) HasAttribute(attributeName string) bool {
	if !s.Enabled {
		return false
	}

	if schema, ok := s.Schema.(rschema.Schema); ok {
		_, exists := schema.Attributes[attributeName]
		return exists
	}

	if schema, ok := s.Schema.(dsschema.Schema); ok {
		_, exists := schema.Attributes[attributeName]
		return exists
	}

	return false
}

// --- Getters ---

func (s *TFState) String(path string) string {
	v := s.Get(path)
	str, ok := v.(types.String)
	if !ok {
		s.convPanic(fmt.Sprintf("not a sting at %q", path))
	}
	return str.ValueString()
}

func (s *TFState) Bool(path string) bool {
	v := s.Get(path)
	b, ok := v.(types.Bool)
	if !ok {
		s.convPanic(fmt.Sprintf("not a bool at %q", path))
	}
	return b.ValueBool()
}

func (s *TFState) Int64(path string) int64 {
	v := s.Get(path)

	switch val := v.(type) {
	case types.Int64:
		return val.ValueInt64()

	case types.Float64:
		// Convert float64 to int64 only if it's a whole number
		f := val.ValueFloat64()
		i := int64(f)
		if float64(i) != f {
			s.convPanic(fmt.Sprintf("cannot convert non-integer float64 to int64 at %q", path))
		}
		return i
	}
	s.convPanic(fmt.Sprintf("expected types.Int64 or types.Float64 at at %q", path))
	return 0
}

func (s *TFState) Float64(path string) float64 {
	v := s.Get(path)
	f, ok := v.(types.Float64)
	if !ok {
		s.convPanic(fmt.Sprintf("not a float64 at %q", path))
	}
	return f.ValueFloat64()
}

func (s *TFState) TfObject(path string) types.Object {
	v := s.Get(path)
	o, ok := v.(types.Object)
	if !ok {
		s.convPanic(fmt.Sprintf("not a types.Object at %q", path))
	}
	return o
}

func (s *TFState) TfList(path string) types.List {
	v := s.Get(path)
	l, ok := v.(types.List)
	if !ok {
		s.convPanic(fmt.Sprintf("not a types.List at %q", path))
	}
	return l
}

func (s *TFState) TfSet(path string) types.Set {
	v := s.Get(path)
	set, ok := v.(types.Set)
	if !ok {
		s.convPanic(fmt.Sprintf("not a types.Set at %q", path))
	}
	return set
}

func (s *TFState) IsNull(path string) bool {
	return s.Get(path).IsNull()
}

func (s *TFState) IsUnknown(path string) bool {
	return s.Get(path).IsUnknown()
}

func (s *TFState) IsKnownAndNotNull(path string) bool {
	return !s.IsNull(path) && !s.IsUnknown(path)
}

// StringWithFallback gets a string value from the current TFState, falling back to another TFState if not found
// Returns the value and a boolean indicating if it was found in either state
func (s *TFState) StringWithFallback(fallback ResourceState, path string) (string, bool) {
	if s.IsKnownAndNotNull(path) {
		return s.String(path), true
	}
	if !IsNilState(fallback) && fallback.IsKnownAndNotNull(path) {
		return fallback.String(path), true
	}
	return "", false
}

// BoolWithFallback gets a bool value from the current TFState, falling back to another TFState if not found
// Returns the value and a boolean indicating if it was found in either state
func (s *TFState) BoolWithFallback(fallback ResourceState, path string) (bool, bool) {
	if s.IsKnownAndNotNull(path) {
		return s.Bool(path), true
	}
	if !IsNilState(fallback) && fallback.IsKnownAndNotNull(path) {
		return fallback.Bool(path), true
	}
	return false, false
}

// Int64WithFallback gets an int64 value from the current TFState, falling back to another TFState if not found
// Returns the value and a boolean indicating if it was found in either state
func (s *TFState) Int64WithFallback(fallback ResourceState, path string) (int64, bool) {
	if s.IsKnownAndNotNull(path) {
		return s.Int64(path), true
	}
	if !IsNilState(fallback) && fallback.IsKnownAndNotNull(path) {
		return fallback.Int64(path), true
	}
	return 0, false
}

// Float64WithFallback gets a float64 value from the current TFState, falling back to another TFState if not found
// Returns the value and a boolean indicating if it was found in either state
func (s *TFState) Float64WithFallback(fallback ResourceState, path string) (float64, bool) {
	if s.IsKnownAndNotNull(path) {
		return s.Float64(path), true
	}
	if !IsNilState(fallback) && fallback.IsKnownAndNotNull(path) {
		return fallback.Float64(path), true
	}
	return 0.0, false
}

func (s *TFState) Get(path string) attr.Value {
	s.assertEnabled()
	parts := parsePath(path)

	var current attr.Value
	currentMap := s.Raw

	for i, part := range parts {
		switch key := part.(type) {
		case string:
			val, ok := currentMap[key]
			if !ok {
				s.convPanic(fmt.Sprintf("key %q not found at path %q", key, path))
			}
			current = val

			if obj, ok := current.(types.Object); ok {
				currentMap = obj.Attributes()
			} else if i < len(parts)-1 {
				s.convPanic(fmt.Sprintf(
					"unexpected non-object at %q (%T)", strings.Join(partsToStrings(parts[:i+1]), "."), current),
				)
			}
		case int:
			elems := getElements(current)
			if key < 0 || key >= len(elems) {
				s.convPanic(fmt.Sprintf("index [%d] out of bounds at %q", key, path))
			}
			current = elems[key]
		}
	}
	return current
}

// --- Setters ---

func (s *TFState) Set(key string, value any) {
	s.assertEnabled()
	_, ok := s.Raw[key]
	if !ok {
		s.convPanic(fmt.Sprintf("Set: key %q not found in state", key))
	}
	destType := s.Type(key)

	attrVal, _, err := BuildAttrValueFromAny(destType, value)
	val := Must(attrVal, err)
	s.Raw[key] = val
}

// SetOrAdd sets a value in the state, adding the key to Raw if it doesn't exist
// This is useful for import operations where we need to set values that aren't in the initial state
func (s *TFState) SetOrAdd(key string, value any) {
	s.assertEnabled()

	// Get the type for this key from the schema
	destType := s.Type(key)

	// Convert the value to the appropriate type
	attrVal, _, err := BuildAttrValueFromAny(destType, value)
	val := Must(attrVal, err)

	// Set the value in Raw (this will add the key if it doesn't exist)
	s.Raw[key] = val
}

// --- Converters ---

func (s *TFState) ToSlice(path string) []any {
	l := s.Get(path)
	if l.IsNull() || l.IsUnknown() {
		return []any{}
	}
	converted, ok := ConvertAttrValueToRaw(l, s.Type(path)).([]any)
	if !ok {
		s.convPanic(fmt.Sprintf("ToSlice: expected a slice at %q, got %T", path, l))
	}
	return RemoveNilValues(converted).([]any)
}

func (s *TFState) ToMap(path string) map[string]any {
	v := s.Get(path)
	if v.IsNull() || v.IsUnknown() {
		return map[string]any{} // return empty map for null or unknown values
	}
	converted, ok := ConvertAttrValueToRaw(v, s.Type(path)).(map[string]any)
	if !ok {
		s.convPanic(fmt.Sprintf("ToMap: expected map[string]any at %q, got %T", path, v))
	}
	return RemoveNilValues(converted).(map[string]any)
}

// SetToMapIfAvailable sets the value at the given path to the provided map
// Helps to avoid checking for null or unknown values before setting
// Returns true if at least one field was set.
func (s *TFState) SetToMapIfAvailable(m map[string]any, path ...string) bool {
	updated := false
	for _, p := range path {
		if !s.IsKnownAndNotNull(p) {
			continue // do not set if the value is unknown or null
		}
		m[p] = ConvertAttrValueToRaw(s.Get(p), s.Type(p))
		updated = true
	}
	return updated
}

// SetIfAvailable is similar to SetToMapIfAvailable but creates and returns a new map
// containing only the fields from the given paths that are known and not null.
// This avoids manually checking values before assignment.
// Returns the new map and a boolean indicating whether any fields were set.
func (s *TFState) SetIfAvailable(path ...string) (map[string]any, bool) {
	m := make(map[string]any, len(path))
	updated := false
	for _, p := range path {
		if !s.IsKnownAndNotNull(p) {
			continue // do not set if the value is unknown or null
		}
		m[p] = ConvertAttrValueToRaw(s.Get(p), s.Type(p))
		updated = true
	}
	return m, updated
}

// --- Attribute meta access via Meta map ---

func (s *TFState) IsRequired(path string) bool {
	s.assertEnabled()
	meta, ok := s.Meta[path]
	if !ok {
		s.convPanic(fmt.Sprintf("IsRequired: path %q not found in meta", path))
	}
	return meta.Required
}

func (s *TFState) IsOptional(path string) bool {
	s.assertEnabled()
	meta, ok := s.Meta[path]
	if !ok {
		s.convPanic(fmt.Sprintf("IsOptional: path %q not found in meta", path))
	}
	return meta.Optional
}

func (s *TFState) IsComputed(path string) bool {
	s.assertEnabled()
	meta, ok := s.Meta[path]
	if !ok {
		s.convPanic(fmt.Sprintf("IsComputed: path %q not found in meta", path))
	}
	return meta.Computed
}

func (s *TFState) Type(path string) attr.Type {
	typ, ok := s.TypeMap[path]
	if !ok {
		s.convPanic(fmt.Sprintf("Type: path %q not found in meta", path))
	}
	return typ
}

func (s *TFState) FillFromRecord(record Record) error {
	return s.FillFromRecordIncludingRequired(record, false)
}

// FillFromRecordIncludingRequired populates TF state from backend record.
// Behavior:
// - Always sets attributes marked Computed.
// - If includeRequired is true, also sets attributes marked Required (even if not Computed).
// - Optional non-computed attributes are not set.
func (s *TFState) FillFromRecordIncludingRequired(record Record, includeRequired bool) error {
	return s.fillFromRecordInternal(record, includeRequired, false)
}

// FillFromRecordForImport populates TF state from backend record during import operations.
// It fills computed and required fields, but NOT optional fields.
// This prevents drift after import by ensuring optional fields remain null in both plan and state
// (since they're not specified in config).
func (s *TFState) FillFromRecordForImport(record Record) error {
	return s.fillFromRecordInternal(record, true, false)
}

// fillFromRecordInternal is the internal implementation for filling state from records.
// Parameters:
// - includeRequired: if true, sets required fields (even if not computed)
// - includeOptional: if true, sets optional fields (even if not computed or required)
func (s *TFState) fillFromRecordInternal(record Record, includeRequired bool, includeOptional bool) error {
	if record == nil {
		return errors.New("record is nil")
	}

	hints := s.Hints
	for key, rawVal := range record {
		typ, ok := s.TypeMap[key]
		if !ok {
			continue // silently skip unknown keys
		}
		// Decide whether to set this field
		if !s.IsComputed(key) {
			// Not computed: check if we should include it
			shouldInclude := false
			if includeRequired && s.IsRequired(key) {
				shouldInclude = true
			}
			if includeOptional && s.IsOptional(key) {
				shouldInclude = true
			}
			if !shouldInclude {
				continue
			}
		}
		val, success, err := BuildAttrValueFromAny(typ, rawVal)
		if !success {
			// Soft error - parsing failed (e.g., empty string for int64)
			continue
		}
		if err != nil && val == nil {
			// Hard error - no value could be produced
			return fmt.Errorf(
				"FillFromRecord for %q failed: %w\nInspected object:\n%v",
				key, err, record.PrettyJson("     "),
			)
		}

		if hints != nil {
			if existing, ok := s.Raw[key]; ok && !existing.IsNull() && !existing.IsUnknown() {
				// For user-configurable fields (Optional or Required), preserve the user's declared value
				// according to provided hint flags. Computed-only fields should always come from API.
				if s.IsOptional(key) || s.IsRequired(key) {
					// PreserveUserValueFields — immutable: once a user sets a value
					// in the .tf config, it is NEVER overwritten by the API response.
					if contains(hints.PreserveUserValueFields, key) {
						continue
					}

					// PreserveUserValueFieldsWhenApiReturnsNull — protect from null:
					// the user's value CAN be updated when the API returns a real,
					// non-null value that differs.  But if the API returned
					// null/unknown we refuse to reset the field — keep the user's
					// declared value.
					if contains(hints.PreserveUserValueFieldsWhenApiReturnsNull, key) {
						// a nil attr.Value.  Treat a Go-nil val the same as a
						// logical null — preserve the user value.
						if val == nil || val.IsNull() || val.IsUnknown() {
							continue
						}
						// Non-null API value → let it fall through to s.Raw[key] = val
					}
				}

				// For PreserveOrderFields, check if the content is the same but order is different.
				// If so, preserve the user's order to avoid perpetual drift.
				if contains(hints.PreserveOrderFields, key) {
					// Check if both values are lists with the same content (ignoring order)
					if listsHaveSameContentIgnoringOrder(existing, val) {
						// Content is identical, preserve user's order
						continue
					}
				}

				// IPSetEquivalenceFields: compare by expanding both sides to their
				// full constituent IP sets.  If the sets are identical the API has
				// merely re-formatted the addresses (e.g. collapsed adjacent IPs
				// into a CIDR) — keep the user's original form so Terraform does
				// not generate a spurious diff.
				if contains(hints.IPSetEquivalenceFields, key) {
					if IPSetsEquivalentAttrs(existing, val) {
						continue
					}
				}
			}
		}

		s.Raw[key] = val
	}
	return nil
}

// SetState is Terraform-specific and intentionally NOT part of ResourceState.
// Orchestration code persists state through a StateSink (see TFSdkSink) instead.
func (s *TFState) SetState(ctx context.Context, state *tfsdk.State) error {
	return writeAttributesToTFSDK(ctx, s.Raw, state)
}

// CopyNonEmptyFieldsTo copies only non-null and known fields from this TFState
// to another, along with their associated attribute metadata.
func (s *TFState) CopyNonEmptyFieldsTo(other ResourceState) {
	s.assertEnabled()
	assertStateEnabled(other)

	otherRaw := other.RawValues()
	otherCore := coreOf(other)
	for k, v := range s.Raw {
		if v.IsNull() || v.IsUnknown() {
			continue // skip null or unknown values
		}
		otherRaw[k] = v
		if otherCore != nil {
			otherCore.Meta[k] = s.Meta[k]
		}
	}
}

// CopyKnownFieldsTo copies only known fields from this TFState
// to another, along with their associated attribute metadata.
// It skips unknown values and null values for computed fields (which should come from API responses).
func (s *TFState) CopyKnownFieldsTo(other ResourceState) {
	s.assertEnabled()
	assertStateEnabled(other)

	otherRaw := other.RawValues()
	otherCore := coreOf(other)
	for k, v := range s.Raw {
		if v.IsUnknown() {
			continue // skip unknown values
		}

		// Skip null values for computed fields to avoid overwriting values from API responses
		// Example: During update, id=5 is set by FillFromRecord from API, but plan has id=null
		// We should NOT overwrite the API value with null from the plan.
		// Exception: <field>_md5 companions for WriteOnlyMD5Fields — a null plan
		// fingerprint means the write-only secret was removed from config (TERF-320).
		meta, ok := s.Meta[k]
		if ok && meta.Computed && v.IsNull() {
			if s.Hints == nil || !isWriteOnlyMD5Companion(s.Hints, k) {
				continue
			}
		}
		// Create-only: omit from config must not wipe the value from Terraform state.
		if v.IsNull() && s.Hints != nil && slices.Contains(s.Hints.CreateOnlyFields, k) {
			continue
		}

		otherRaw[k] = v
		if otherCore != nil {
			otherCore.Meta[k] = s.Meta[k]
		}
	}
}

// GetFilteredValues returns a map of state values filtered by the provided field filter flags.
//
// By default, it returns only non-null and known values (i.e., populated fields).
// If the SearchEmpty flag is included, it also includes fields that are null or unknown.
//
// Each field is included only if it satisfies the given metadata flags (such as SearchRequired, SearchOptional, etc.),
// and its value matches the presence condition (non-empty or empty) based on the SearchEmpty flag.
//
// Parameters:
//   - comb: whether to use OR or AND logic when evaluating multiple flags.
//   - flags: a list of FieldFilterFlag values used to determine which fields to include.
//
// Returns:
//   - map[string]any: a map of field names to their converted raw values, filtered based on metadata and value presence.
func (s *TFState) GetFilteredValues(
	comb FilterTagCombination,
	fieldSet *FieldSet,
	flags ...FieldFilterFlag,
) map[string]any {
	s.assertEnabled()

	searchEmpty := contains(flags, SearchEmpty)
	primitivesOnly := contains(flags, SearchPrimitivesOnly)

	result := make(map[string]any)

	for k, v := range s.Raw {
		if fieldSet != nil && !fieldSet.shouldInclude(k) {
			continue
		}

		valType := s.Type(k)

		if primitivesOnly && !isPrimitiveType(valType) {
			continue
		}

		meta, ok := s.Meta[k]
		if !ok || !meta.satisfyFieldFilterFlag(comb, flags...) {
			continue
		}

		if v.IsNull() || v.IsUnknown() {
			if searchEmpty {
				result[k] = ConvertAttrValueToRaw(v, valType)
			}
			continue
		}

		raw := ConvertAttrValueToRaw(v, valType)
		if raw != nil {
			result[k] = raw
		}
	}

	if !searchEmpty {
		result = RemoveNilValues(result).(map[string]any)
	}
	return result
}

func (s *TFState) GetAllValues() map[string]any {
	s.assertEnabled()

	result := make(map[string]any, len(s.Raw))
	for k, v := range s.Raw {
		typ := s.Type(k)
		result[k] = ConvertAttrValueToRaw(v, typ)
	}
	return result
}

func (s *TFState) DiffFields(
	other ResourceState,
	comb FilterTagCombination,
	fields []string,
	flags ...FieldFilterFlag,

) map[string]any {
	s.assertEnabled()
	assertStateEnabled(other)

	otherRaw := other.RawValues()
	searchEmpty := contains(flags, SearchEmpty)

	fieldMap := make(map[string]struct{})
	if fields != nil {
		for _, field := range fields {
			fieldMap[field] = struct{}{}
		}
	}

	diff := make(map[string]any)
	// clearedFields collects fields explicitly set to null in the plan that previously
	// had a non-null value in the current state. These must be sent as JSON null in the
	// PATCH body so the API removes the value; they are merged back after RemoveNilValues
	// so they are not accidentally dropped.
	clearedFields := make(map[string]any)

	for k, v := range s.Raw {
		if fields != nil {
			if _, ok := fieldMap[k]; !ok {
				continue // skip fields not in the specified list
			}
		}

		valType := s.Type(k)
		meta, ok := s.Meta[k]

		if v.IsNull() || v.IsUnknown() {
			if !searchEmpty && ok && meta.satisfyFieldFilterFlag(comb, flags...) {
				// Create-only fields must not be cleared via PATCH when omitted from config.
				if s.Hints != nil && slices.Contains(s.Hints.CreateOnlyFields, k) {
					continue
				}
				// Detect a non-null → null transition: the plan clears this field.
				if otherVal, exists := otherRaw[k]; exists && !otherVal.IsNull() && !otherVal.IsUnknown() {
					clearedFields[k] = nil
				}
			}
			continue
		}

		if !ok || !meta.satisfyFieldFilterFlag(comb, flags...) {
			continue
		}

		if otherVal, ok := otherRaw[k]; !ok || otherVal.IsNull() || otherVal.IsUnknown() || !v.Equal(otherVal) {
			diff[k] = ConvertAttrValueToRaw(v, valType)
		}
	}

	if !searchEmpty {
		diff = RemoveNilValues(diff).(map[string]any)
	}
	// Merge cleared fields last so RemoveNilValues does not strip the intentional nulls.
	for k, v := range clearedFields {
		diff[k] = v
	}

	return diff
}

// GetGenericSearchParams returns a map of search parameters for the resource
// Such cases for search parameters are super common for VAST.
func (s *TFState) GetGenericSearchParams(ctx context.Context) vast_client.Params {
	var exclude []string
	if s.Hints != nil {
		exclude = append(exclude, s.Hints.EditOnlyFields...)                                   // Edit only fields should not be set on creation.
		exclude = append(exclude, slices.Collect(maps.Keys(s.Hints.DeleteOnlyBodyFields))...)  // Delete only fields should not be set on creation.
		exclude = append(exclude, slices.Collect(maps.Keys(s.Hints.DeleteOnlyParamFields))...) // Delete only fields should not be set on creation.
		exclude = append(exclude, s.Hints.SensitiveFields...)                                  // Sensitive fields must never appear in URL query params.
	}

	searchParams := make(vast_client.Params)
	if sp := s.GetFilteredValues(
		FilterOr,
		&FieldSet{
			Include: CommonSearchableFields,
		},
		SearchRequired,
		SearchOptional,
		SearchSearchable,
	); len(sp) > 0 {
		// Last attempt: get all optional fields and perform search by all of them
		logging.Debug(ctx, "++ 'optional(common search scope)'")
		searchParams = sp
	} else if sp := s.GetFilteredValues(
		FilterOr,
		&FieldSet{
			Exclude: exclude,
		},
		SearchRequired,
		SearchSearchable,
	); len(sp) > 0 {
		// Get all params required + searchable for search.
		logging.Debug(ctx, "++ 'required+searchable'")
		searchParams = sp
	} else if sp := s.GetFilteredValues(
		FilterOr,
		nil,
		SearchOptional,
		SearchPrimitivesOnly,
	); len(sp) > 0 {
		// Get all optional fields that are primitive types (string, int, bool) for search.
		logging.Debug(ctx, "++ 'search by all non-primitive fields'")
		searchParams = sp
	}

	for _, field := range AdditionalSearchableFields {
		// Try go get resource ID and GUID and other additional fields from current state
		if val, ok := s.Raw[field]; ok && !val.IsNull() && !val.IsUnknown() {
			logging.Debug(ctx, fmt.Sprintf("++ 'search by %s'", field))
			searchParams[field] = ConvertAttrValueToRaw(val, s.Type(field))
		}
	}

	// Always append SearchFilterFields when set — these narrow the search to the
	// correct scope (AND semantics) regardless of which primary search branch fired.
	if s.Hints != nil {
		for _, field := range s.Hints.SearchFilterFields {
			if val, ok := s.Raw[field]; ok && !val.IsNull() && !val.IsUnknown() {
				logging.Debug(ctx, fmt.Sprintf("++ 'search filter field %s'", field))
				searchParams[field] = ConvertAttrValueToRaw(val, s.Type(field))
			}
		}
	}

	searchParams.Update(s.GetReadOnlySearchParams(), false)

	// Sensitive fields must never appear in URL query parameters.
	if s.Hints != nil {
		for _, field := range s.Hints.SensitiveFields {
			delete(searchParams, field)
		}
	}

	return searchParams

}

// GetReadOnlySearchParams returns a map of search parameters for the resource
// where only 'SearchReadonly' fields are included.
func (s *TFState) GetReadOnlySearchParams() vast_client.Params {
	searchParams := make(vast_client.Params)
	searchParams.Update(s.GetFilteredValues(
		FilterOr,
		nil,
		SearchReadonly,
	), true)

	return searchParams

}

// GetReadEditOnlyParams returns a map of parameters used exclusively for create/update (edit-only) operations.
// These are fields that are not part of search or identification but are required for modifying a resource.
func (s *TFState) GetReadEditOnlyParams() vast_client.Params {
	searchParams := make(vast_client.Params)
	if s.Hints != nil && len(s.Hints.EditOnlyFields) > 0 {
		searchParams.Update(s.GetFilteredValues(
			FilterOr,
			&FieldSet{
				Include: s.Hints.EditOnlyFields,
			},
			SearchOptional,
		), true)
	}
	return searchParams
}

// GetChangedEditOnlyParams returns edit-only fields that have changed between this TFState (plan)
// and another TFState (current state).
//
// This is used in UpdateResource to handle edit-only fields separately from the main update.
// Edit-only fields are those that cannot be set during creation but can be modified via update operations.
func (s *TFState) GetChangedEditOnlyParams(otherState ResourceState) vast_client.Params {
	searchParams := make(vast_client.Params)
	if s.Hints != nil && len(s.Hints.EditOnlyFields) > 0 {
		// Get all changed fields
		diffParams := s.DiffFields(otherState, FilterOr, nil, SearchOptional)

		// Keep only edit-only fields
		for key := range diffParams {
			if !slices.Contains(s.Hints.EditOnlyFields, key) {
				delete(diffParams, key)
			}
		}

		searchParams.Update(diffParams, true)
	}
	return searchParams
}

// GetChangedParams returns a map of parameters that differ between this TFState (plan)
// and another TFState (current state), returning only the changed fields.
// This is commonly used in UpdateResource methods to send only modified fields to the API.
func (s *TFState) GetChangedParams(otherState ResourceState) vast_client.Params {
	diffParams := s.DiffFields(otherState, FilterOr, nil, SearchOptional, SearchRequired)
	return diffParams
}

// GetUpdateParams returns a map of changed parameters suitable for update operations,
// excluding edit-only fields, create-only fields, and delete-only fields.
//
// This method:
//  1. Gets all changed fields between plan and current state (like GetChangedParams)
//  2. Excludes EditOnlyFields - fields that should only be set via separate edit operations
//  3. Excludes CreateOnlyFields - fields that may be set on create but must never be PATCH'd
//  4. Excludes DeleteOnlyBodyFields - fields only used during delete operations
//  5. Excludes DeleteOnlyParamFields - query params only used during delete operations
//  6. When WriteOnlyMD5Fields is set, omits WriteOnlyFields unless a <field>_md5
//     fingerprint changed (avoids re-sending secrets after overlayWriteOnlyFromConfig)
//
// EditOnlyFields should be handled separately after the main update operation.
func (s *TFState) GetUpdateParams(otherState ResourceState) vast_client.Params {
	// Build exclusion list
	var exclude []string
	if s.Hints != nil {
		exclude = append(exclude, s.Hints.EditOnlyFields...)                                   // Edit only fields should be updated separately
		exclude = append(exclude, s.Hints.CreateOnlyFields...)                                 // Create-only fields must never be PATCH'd
		exclude = append(exclude, slices.Collect(maps.Keys(s.Hints.DeleteOnlyBodyFields))...)  // Delete only fields should not be in update
		exclude = append(exclude, slices.Collect(maps.Keys(s.Hints.DeleteOnlyParamFields))...) // Delete only fields should not be in update
		exclude = append(exclude, s.subResourceKeys()...)                                      // Sub-resource fields are managed via hooks, not the main API.
	}

	// Get all changed params
	diffParams := s.DiffFields(otherState, FilterOr, nil, SearchOptional, SearchRequired)

	// Remove excluded fields from the diff
	for _, excludeField := range exclude {
		delete(diffParams, excludeField)
	}

	// Write-only secrets are overlaid from Config onto the plan, so DiffFields
	// always sees plan=value vs state=null. Only send them when their MD5
	// fingerprint actually changed.
	if s.Hints != nil && len(s.Hints.WriteOnlyMD5Fields) > 0 && !s.writeOnlyMD5Changed(otherState) {
		for _, field := range s.Hints.WriteOnlyFields {
			delete(diffParams, field)
		}
	}

	return diffParams
}

// writeOnlyMD5Changed reports whether any WriteOnlyMD5Fields fingerprint
// (<field>_md5) differs between this plan state and otherState.
//
// Upgrade path: when state has no <field>_md5 yet (provider just added the
// companion) but still holds the pre-WriteOnly secret, compare MD5(state secret)
// to the plan fingerprint. Matching means the password did not change — only
// seed md5 into state; do not re-send the secret. Import (no state secret, no
// md5) still counts as changed so the configured secret is applied once.
func (s *TFState) writeOnlyMD5Changed(otherState ResourceState) bool {
	if s.Hints == nil || len(s.Hints.WriteOnlyMD5Fields) == 0 {
		return false
	}
	otherRaw := otherState.RawValues()
	for _, field := range s.Hints.WriteOnlyMD5Fields {
		md5Field := field + "_md5"
		planMD5, planOK := attrValueString(s.Raw[md5Field])
		if !planOK || planMD5 == "" {
			// Cleared optional write-only secret (TERF-320): do not treat as a
			// secret re-send. Resources that need side effects on clear (e.g.
			// tls_certificate CRL delete) handle that in AfterUpdateResource.
			continue
		}
		if stateMD5, ok := attrValueString(otherRaw[md5Field]); ok {
			if stateMD5 == planMD5 {
				continue
			}
			return true
		}
		// State md5 missing: prefer hashing a leftover state secret (upgrade).
		if stateSecret, ok := attrValueString(otherRaw[field]); ok && stateSecret != "" {
			if contentMD5Hex(stateSecret) == planMD5 {
				continue
			}
			return true
		}
		// No state md5 and no state secret (e.g. import) → treat as new secret.
		return true
	}
	return false
}

func attrValueString(v attr.Value) (string, bool) {
	if v == nil || v.IsNull() || v.IsUnknown() {
		return "", false
	}
	s, ok := v.(types.String)
	if !ok {
		return "", false
	}
	return s.ValueString(), true
}

//nolint:gosec // plan/state equality fingerprint, not a security control
func contentMD5Hex(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

// isWriteOnlyMD5Companion reports whether attrName is <field>_md5 for a field
// listed in hints.WriteOnlyMD5Fields.
func isWriteOnlyMD5Companion(hints *TFStateHints, attrName string) bool {
	if hints == nil || !strings.HasSuffix(attrName, "_md5") {
		return false
	}
	field := strings.TrimSuffix(attrName, "_md5")
	return slices.Contains(hints.WriteOnlyMD5Fields, field)
}

// GetDeleteOnlyBodyParams returns a map of parameters used exclusively for delete operations (delete-only).
// These fields are not used during normal lifecycle operations, but may be required for safe deletion.
func (s *TFState) GetDeleteOnlyBodyParams() vast_client.Params {
	searchParams := make(vast_client.Params)
	if s.Hints != nil && len(s.Hints.DeleteOnlyBodyFields) > 0 {
		deleteOnlyKeys := slices.Collect(maps.Keys(s.Hints.DeleteOnlyBodyFields))
		searchParams.Update(s.GetFilteredValues(
			FilterOr,
			&FieldSet{
				Include: deleteOnlyKeys,
			},
			SearchOptional,
		), true)
		for _, key := range deleteOnlyKeys {
			apiName, ok := s.Hints.DeleteOnlyBodyFields[key]
			if !ok || apiName == "" {
				delete(searchParams, key)
				continue
			}
			if v, present := searchParams[key]; present {
				if apiName != key {
					searchParams[apiName] = v // remap to API name
					delete(searchParams, key)
				}
			}
		}
	}
	return searchParams
}

// GetDeleteOnlyQueryParams returns a map of parameters used exclusively for delete operations (delete-only).
// These fields are not used during normal lifecycle operations, but may be required for safe deletion.
func (s *TFState) GetDeleteOnlyQueryParams() vast_client.Params {
	searchParams := make(vast_client.Params)
	if s.Hints != nil && len(s.Hints.DeleteOnlyParamFields) > 0 {
		deleteOnlyKeys := slices.Collect(maps.Keys(s.Hints.DeleteOnlyParamFields))

		searchParams.Update(s.GetFilteredValues(
			FilterOr,
			&FieldSet{
				Include: deleteOnlyKeys,
			},
			SearchOptional,
		), true)
		for _, key := range deleteOnlyKeys {
			apiName, ok := s.Hints.DeleteOnlyParamFields[key]
			if !ok || apiName == "" {
				delete(searchParams, key)
				continue
			}
			if v, present := searchParams[key]; present {
				if apiName != key {
					searchParams[apiName] = v // remap to API name
					delete(searchParams, key)
				}
			}
		}
	}
	return searchParams
}

// subResourceKeys returns the top-level field names contributed by all declared
// SubResources hints. These fields are managed exclusively via lifecycle hooks
// (AfterCreateResource / AfterUpdateResource) and must never appear in the
// main resource create or update body sent to the API.
func (s *TFState) subResourceKeys() []string {
	if s.Hints == nil {
		return nil
	}
	var keys []string
	for _, sr := range s.Hints.SubResources {
		if sr.FieldTrigger != "" {
			keys = append(keys, sr.FieldTrigger)
		}
		if sr.SchemaKey != "" {
			keys = append(keys, sr.SchemaKey)
		} else {
			for k := range sr.SchemaAttributes {
				keys = append(keys, k)
			}
		}
	}
	return keys
}

// GetCreateParams returns a map of parameters used for resource creation
func (s *TFState) GetCreateParams() vast_client.Params {
	// Get all params required + optional for creation.
	var exclude []string
	if s.Hints != nil {
		exclude = append(exclude, s.Hints.EditOnlyFields...)                                   // Edit only fields should not be set on creation.
		exclude = append(exclude, slices.Collect(maps.Keys(s.Hints.DeleteOnlyBodyFields))...)  // Delete only fields should not be set on creation.
		exclude = append(exclude, slices.Collect(maps.Keys(s.Hints.DeleteOnlyParamFields))...) // Delete only fields should not be set on creation.
		exclude = append(exclude, s.subResourceKeys()...)                                      // Sub-resource fields are managed via hooks, not the main API.
	}

	createParams := s.GetFilteredValues(
		FilterOr,
		&FieldSet{
			Exclude: exclude,
		},
		SearchRequired,
		SearchOptional,
	)
	delete(createParams, "id") // Remove ID from update parameters, as it should not be set during creation.
	return createParams
}

// ClearWriteOnlyFields nulls all WriteOnlyFields in Raw so they are never
// persisted to Terraform state (write-only values must always be null in state).
func (s *TFState) ClearWriteOnlyFields() {
	if s == nil || !s.Enabled || s.Hints == nil {
		return
	}
	for _, key := range s.Hints.WriteOnlyFields {
		typ, ok := s.TypeMap[key]
		if !ok {
			delete(s.Raw, key)
			continue
		}
		nullVal, _, err := BuildAttrValueFromAny(typ, nil)
		if err != nil {
			delete(s.Raw, key)
			continue
		}
		s.Raw[key] = nullVal
	}
}

// SeedWriteOnlyMD5FromSecrets fills empty <field>_md5 companions from leftover
// WriteOnly secret values. Used only by schema upgrades (concrete *TFState);
// intentionally not on ResourceState.
func (s *TFState) SeedWriteOnlyMD5FromSecrets() {
	if s == nil || !s.Enabled || s.Hints == nil {
		return
	}
	for _, field := range s.Hints.WriteOnlyMD5Fields {
		md5Field := field + "_md5"
		if existing, ok := attrValueString(s.Raw[md5Field]); ok && existing != "" {
			continue
		}
		secret, ok := attrValueString(s.Raw[field])
		if !ok || secret == "" {
			continue
		}
		s.Raw[md5Field] = types.StringValue(contentMD5Hex(secret))
	}
}

// listsHaveSameContentIgnoringOrder checks if two attr.Value instances are lists
// with the same content, regardless of order. This is used to prevent drift when
// the API returns list elements in a different order than the user's configuration.
func listsHaveSameContentIgnoringOrder(a, b attr.Value) bool {
	// Both must be lists
	listA, okA := a.(types.List)
	listB, okB := b.(types.List)
	if !okA || !okB {
		return false
	}

	// Must have the same length
	elementsA := listA.Elements()
	elementsB := listB.Elements()

	if len(elementsA) != len(elementsB) {
		return false
	}

	// Create a slice to track which elements in B have been matched
	matchedB := make([]bool, len(elementsB))

	// For each element in A, find a matching element in B (that hasn't been matched yet)
	for _, elemA := range elementsA {
		found := false
		for j, elemB := range elementsB {
			if matchedB[j] {
				continue // Already matched
			}
			if elemA.Equal(elemB) {
				matchedB[j] = true
				found = true
				break
			}
		}
		if !found {
			return false // Element in A not found in B
		}
	}

	return true
}

// ============================================================================
// Resource state abstraction (shared by Terraform and Crossplane)
// ============================================================================
//
// Resource managers (User, VipPool, License, ...) only ever talk to ResourceState.
// The attribute bag logic (typed getters/setters, FillFromRecord, request params
// builders, diffing) is identical for every runtime; what differs per runtime is:
//
//   - how the bag is seeded (Terraform plan/state vs Crossplane CR spec),
//   - which values may be unknown (Terraform plans have unknowns, CRs do not),
//   - where the result is persisted (tfsdk.State vs CR status) - see StateSink.
//
// Terraform uses *TFState directly. Other runtimes embed *TFState and override
// only what differs (see crossplane-provider-vastdata/internal/state).

// ResourceState is the state abstraction managers depend on.
type ResourceState interface {
	// --- Introspection ---

	// IsEnabled reports whether a schema is attached (false for schema-less placeholder states).
	IsEnabled() bool
	// StateHints returns the manager hints (may be nil for placeholder states).
	StateHints() *TFStateHints
	// SchemaDefinition returns the schema the state was built from
	// (resource_schema.Schema or datasource_schema.Schema).
	SchemaDefinition() any
	// RawValues returns the LIVE attribute map. Mutating it mutates the state.
	RawValues() map[string]attr.Value
	HasAttribute(attributeName string) bool
	Type(path string) attr.Type
	IsRequired(path string) bool
	IsOptional(path string) bool
	IsComputed(path string) bool
	Pretty() string

	// --- Getters ---

	String(path string) string
	Bool(path string) bool
	Int64(path string) int64
	Float64(path string) float64
	TfObject(path string) types.Object
	TfList(path string) types.List
	TfSet(path string) types.Set
	IsNull(path string) bool
	IsUnknown(path string) bool
	IsKnownAndNotNull(path string) bool
	StringWithFallback(fallback ResourceState, path string) (string, bool)
	BoolWithFallback(fallback ResourceState, path string) (bool, bool)
	Int64WithFallback(fallback ResourceState, path string) (int64, bool)
	Float64WithFallback(fallback ResourceState, path string) (float64, bool)
	Get(path string) attr.Value
	ToSlice(path string) []any
	ToMap(path string) map[string]any
	SetToMapIfAvailable(m map[string]any, path ...string) bool
	SetIfAvailable(path ...string) (map[string]any, bool)
	GetAllValues() map[string]any
	GetFilteredValues(comb FilterTagCombination, fieldSet *FieldSet, flags ...FieldFilterFlag) map[string]any

	// --- Mutators ---

	Set(key string, value any)
	SetOrAdd(key string, value any)
	FillFromRecord(record Record) error
	FillFromRecordIncludingRequired(record Record, includeRequired bool) error
	FillFromRecordForImport(record Record) error
	ClearWriteOnlyFields()
	CopyNonEmptyFieldsTo(other ResourceState)
	CopyKnownFieldsTo(other ResourceState)
	// Copy returns an independent copy of the same implementation type.
	Copy() ResourceState

	// --- Request parameter builders ---

	DiffFields(other ResourceState, comb FilterTagCombination, fields []string, flags ...FieldFilterFlag) map[string]any
	GetGenericSearchParams(ctx context.Context) vast_client.Params
	GetReadOnlySearchParams() vast_client.Params
	GetReadEditOnlyParams() vast_client.Params
	GetCreateParams() vast_client.Params
	GetUpdateParams(other ResourceState) vast_client.Params
	GetChangedParams(other ResourceState) vast_client.Params
	GetChangedEditOnlyParams(other ResourceState) vast_client.Params
	GetDeleteOnlyBodyParams() vast_client.Params
	GetDeleteOnlyQueryParams() vast_client.Params
}

var _ ResourceState = (*TFState)(nil)

// IsEnabled implements ResourceState.
func (s *TFState) IsEnabled() bool { return s.Enabled }

// StateHints implements ResourceState.
func (s *TFState) StateHints() *TFStateHints { return s.Hints }

// SchemaDefinition implements ResourceState.
func (s *TFState) SchemaDefinition() any { return s.Schema }

// RawValues implements ResourceState. It returns the live attribute map.
func (s *TFState) RawValues() map[string]attr.Value { return s.Raw }

// TFStateCore returns the underlying *TFState. It is deliberately not part of
// ResourceState: it is promoted through embedding so wrapper implementations
// remain usable wherever the concrete attribute bag (Meta/TypeMap) is needed.
func (s *TFState) TFStateCore() *TFState { return s }

// coreOf returns the underlying *TFState of any state built on top of one, or nil.
func coreOf(s ResourceState) *TFState {
	if c, ok := s.(interface{ TFStateCore() *TFState }); ok {
		return c.TFStateCore()
	}
	return nil
}

// IsNilState reports whether s is nil, including a typed-nil pointer boxed in the interface.
func IsNilState(s ResourceState) bool {
	if s == nil {
		return true
	}
	v := reflect.ValueOf(s)
	return v.Kind() == reflect.Ptr && v.IsNil()
}

func assertStateEnabled(s ResourceState) {
	if IsNilState(s) || !s.IsEnabled() {
		panic("ResourceState is nil or disabled; schema and meta not available")
	}
}

// --- State sinks ---

// StateSink persists a ResourceState into the runtime's backing store.
// Terraform: a *tfsdk.State. Crossplane: the managed resource status.
type StateSink interface {
	WriteState(ctx context.Context, s ResourceState) error
}

// TFSdkSink writes a ResourceState into a Terraform Plugin Framework state.
type TFSdkSink struct {
	State *tfsdk.State
}

// NewTFSdkSink returns a StateSink backed by the given Terraform state.
func NewTFSdkSink(state *tfsdk.State) StateSink {
	return TFSdkSink{State: state}
}

// WriteState implements StateSink.
func (t TFSdkSink) WriteState(ctx context.Context, s ResourceState) error {
	return writeAttributesToTFSDK(ctx, s.RawValues(), t.State)
}

func writeAttributesToTFSDK(ctx context.Context, raw map[string]attr.Value, state *tfsdk.State) error {
	for k, v := range raw {
		attrPath := path.Root(k)

		logging.Debug(ctx, fmt.Sprintf(
			"SetState: key=%q type=%T value=%v",
			k, v, v,
		))

		if diags := state.SetAttribute(ctx, attrPath, v); diags.HasError() {
			return fmt.Errorf("set attribute %q: %s", k, diags.Errors())
		}
	}
	return nil
}

// --- State factory ---

// StateFactory builds the ResourceState implementation used by resource managers.
type StateFactory func(raw map[string]attr.Value, schema any, hints *TFStateHints) (ResourceState, error)

var stateFactory StateFactory

// SetStateFactory overrides how NewResourceState builds states. It is meant to be
// called once at process start-up (e.g. by the Crossplane provider's main) before
// any manager is created; passing nil restores the default Terraform behaviour.
// It is process-wide and not safe to call concurrently with NewResourceState.
func SetStateFactory(f StateFactory) {
	stateFactory = f
}

// NewResourceState is the single construction point managers use for their state.
// By default it returns a *TFState; see SetStateFactory.
func NewResourceState(raw map[string]attr.Value, schema any, hints *TFStateHints) ResourceState {
	if stateFactory == nil {
		return NewTFStateMust(raw, schema, hints)
	}
	return Must(stateFactory(raw, schema, hints))
}
