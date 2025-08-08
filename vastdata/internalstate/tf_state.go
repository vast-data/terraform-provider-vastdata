// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	vast_client "github.com/vast-data/go-vast-client"

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
	"name", "path", "tenant_id", "tenant_name", "bucket", "gid", "uid",
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
	for _, ex := range fs.Exclude {
		if ex == field {
			return false
		}
	}
	if len(fs.Include) > 0 {
		for _, in := range fs.Include {
			if in == field {
				return true
			}
		}
		return false
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
		default:
			panic(fmt.Sprintf("unknown search context flag: %d", f))
		}
	}

	switch comb {
	case FilterOr:
		for _, f := range flags {
			if matches(f) {
				return true
			}
		}
		return false

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

func NewTFStateMust(raw map[string]attr.Value, schema any, hints *TFStateHints) *TFState {
	var kind SchemaContext
	if schema == nil {
		kind = SchemaTypeUnknown
	} else {
		switch schema.(type) {
		case rschema.Schema:
			kind = SchemaForResource
		case dsschema.Schema:
			kind = SchemaForDataSource
		default:
			panic(fmt.Sprintf("unknown schema type: %T", schema))
		}
	}
	return Must(NewTFState(raw, schema, kind, hints))
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
func (s *TFState) Copy() *TFState {
	s.assertEnabled()

	rawCopy := make(map[string]attr.Value, len(s.Raw))
	for k, v := range s.Raw {
		rawCopy[k] = v
	}

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

	val := Must(BuildAttrValueFromAny(destType, value))
	s.Raw[key] = val
}

// SetOrAdd sets a value in the state, adding the key to Raw if it doesn't exist
// This is useful for import operations where we need to set values that aren't in the initial state
func (s *TFState) SetOrAdd(key string, value any) {
	s.assertEnabled()

	// Get the type for this key from the schema
	destType := s.Type(key)

	// Convert the value to the appropriate type
	val := Must(BuildAttrValueFromAny(destType, value))

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
	return removeNilValues(converted).([]any)
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
	return removeNilValues(converted).(map[string]any)
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
	if record == nil {
		return errors.New("record is nil")
	}
	for key, rawVal := range record {
		typ, ok := s.TypeMap[key]
		if !ok {
			continue // silently skip unknown keys
		}
		if !s.IsComputed(key) {
			// Set only computed fields
			continue
		}
		val, err := BuildAttrValueFromAny(typ, rawVal)
		if err != nil {
			return fmt.Errorf(
				"FillFromRecord for %q failed: %w\nInspected object:\n%v",
				key, err, record.PrettyJson("     "),
			)
		}
		s.Raw[key] = val
	}
	return nil
}

func (s *TFState) SetState(ctx context.Context, state *tfsdk.State) error {
	for k, v := range s.Raw {
		attrPath := path.Root(k)

		tflog.Debug(ctx, fmt.Sprintf(
			"SetState: key=%q type=%T value=%v",
			k, v, v,
		))

		if diags := state.SetAttribute(ctx, attrPath, v); diags.HasError() {
			return fmt.Errorf("set attribute %q: %s", k, diags.Errors())
		}
	}
	return nil
}

// CopyNonEmptyFieldsTo copies only non-null and known fields from this TFState
// to another, along with their associated attribute metadata.
func (s *TFState) CopyNonEmptyFieldsTo(other *TFState) {
	s.assertEnabled()
	other.assertEnabled()

	for k, v := range s.Raw {
		if v.IsNull() || v.IsUnknown() {
			continue // skip null or unknown values
		}
		other.Raw[k] = v
		other.Meta[k] = s.Meta[k]
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
		result = removeNilValues(result).(map[string]any)
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
	other *TFState,
	comb FilterTagCombination,
	fields []string,
	flags ...FieldFilterFlag,

) map[string]any {
	s.assertEnabled()
	other.assertEnabled()

	searchEmpty := contains(flags, SearchEmpty)

	fieldMap := make(map[string]struct{})
	if fields != nil {
		for _, field := range fields {
			fieldMap[field] = struct{}{}
		}
	}

	diff := make(map[string]any)

	for k, v := range s.Raw {
		if fields != nil {
			if _, ok := fieldMap[k]; !ok {
				continue // skip fields not in the specified list
			}
		}

		valType := s.Type(k)
		meta, ok := s.Meta[k]

		if v.IsNull() || v.IsUnknown() {
			continue // skip null or unknown values
		}

		if v.IsNull() || v.IsUnknown() {
			if !searchEmpty {
				continue
			} else {
				diff[k] = ConvertAttrValueToRaw(v, valType)
				continue
			}
		}

		if !ok || !meta.satisfyFieldFilterFlag(comb, flags...) {
			continue
		}

		if otherVal, ok := other.Raw[k]; !ok || otherVal.IsNull() || otherVal.IsUnknown() || !v.Equal(otherVal) {
			diff[k] = ConvertAttrValueToRaw(v, valType)
		}
	}

	if !searchEmpty {
		diff = removeNilValues(diff).(map[string]any)
	}
	return diff
}

// GetGenericSearchParams returns a map of search parameters for the resource
// Such cases for search parameters are super common for VAST.
func (s *TFState) GetGenericSearchParams(ctx context.Context) vast_client.Params {
	var exclude []string
	if s.Hints != nil {
		exclude = append(exclude, s.Hints.EditOnlyFields...)   // Edit only fields should not be set on creation.
		exclude = append(exclude, s.Hints.DeleteOnlyFields...) // Delete only fields should not be set on creation.
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
		tflog.Debug(ctx, "++ 'optional(common search scope)'")
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
		tflog.Debug(ctx, "++ 'required+searchable'")
		searchParams = sp
	}
	// Try go get resource ID and GUID from current state
	if idFromState, ok := s.Raw["id"]; ok && !idFromState.IsNull() && !idFromState.IsUnknown() {
		tflog.Debug(ctx, "++ 'search by ID'")
		// Convert the ID to raw value for search params
		searchParams["id"] = ConvertAttrValueToRaw(idFromState, s.Type("id"))
	}

	if guid, ok := s.Raw["guid"]; ok && !guid.IsNull() && !guid.IsUnknown() {
		tflog.Debug(ctx, "++ 'search by GUID'")
		searchParams["guid"] = ConvertAttrValueToRaw(guid, s.Type("guid"))
	}

	searchParams.Update(s.GetReadOnlySearchParams(), false)

	if len(searchParams) == 0 {
		// Still no search params, try to get all non-primitive fields
		tflog.Debug(ctx, "++ 'search by all non-primitive fields'")
		if sp := s.GetFilteredValues(
			FilterOr,
			nil,
			SearchOptional,
			SearchPrimitivesOnly,
		); len(sp) > 0 {
			searchParams = sp
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

// GetDeleteOnlyParams returns a map of parameters used exclusively for delete operations (delete-only).
// These fields are not used during normal lifecycle operations, but may be required for safe deletion.
func (s *TFState) GetDeleteOnlyParams() vast_client.Params {
	searchParams := make(vast_client.Params)
	if s.Hints != nil && len(s.Hints.DeleteOnlyFields) > 0 {
		searchParams.Update(s.GetFilteredValues(
			FilterOr,
			&FieldSet{
				Include: s.Hints.DeleteOnlyFields,
			},
			SearchOptional,
		), true)
	}
	return searchParams
}

// GetCreateParams returns a map of parameters used for resource creation
func (s *TFState) GetCreateParams() vast_client.Params {
	// Get all params required + optional for creation.
	var exclude []string
	if s.Hints != nil {
		exclude = append(exclude, s.Hints.EditOnlyFields...)   // Edit only fields should not be set on creation.
		exclude = append(exclude, s.Hints.DeleteOnlyFields...) // Delete only fields should not be set on creation.
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
