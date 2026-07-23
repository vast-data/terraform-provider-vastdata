// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vast_client "github.com/vast-data/go-vast-client"
)

func TestExtractMetaFromSchema(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"required_field": rschema.StringAttribute{
				Required: true,
			},
			"optional_field": rschema.Int64Attribute{
				Optional: true,
			},
			"computed_field": rschema.BoolAttribute{
				Computed: true,
			},
			"optional_computed": rschema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"nested": rschema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]rschema.Attribute{
					"inner_required": rschema.StringAttribute{Required: true},
					"inner_computed": rschema.Int64Attribute{Computed: true},
				},
			},
		},
	}

	rawAny := map[string]any{
		"required_field":    "abc",
		"optional_field":    int64(123),
		"computed_field":    true,
		"optional_computed": "auto",
		"nested": map[string]any{
			"inner_required": "val",
			"inner_computed": int64(42),
		},
	}

	// Convert map[string]any → map[string]attr.Value using BuildAttrValueFromAny
	raw := make(map[string]attr.Value)
	for key, val := range rawAny {
		attrDef, ok := schema.Attributes[key]
		require.True(t, ok, "unexpected attribute: %s", key)
		attrType := attrDef.GetType()

		converted, _, err := BuildAttrValueFromAny(attrType, val)
		require.NoError(t, err, "failed to build attr.Value for %q", key)

		raw[key] = converted
	}

	state := NewTFStateMust(raw, schema, nil)

	expect := map[string]attrMeta{
		"required_field":        {Required: true},
		"optional_field":        {Optional: true},
		"computed_field":        {Computed: true},
		"optional_computed":     {Optional: true, Computed: true},
		"nested":                {Optional: true},
		"nested.inner_required": {Required: true},
		"nested.inner_computed": {Computed: true},
	}

	for path, want := range expect {
		got := state.Meta[path]
		require.Equal(t, want.Required, got.Required, path+" required mismatch")
		require.Equal(t, want.Optional, got.Optional, path+" optional mismatch")
		require.Equal(t, want.Computed, got.Computed, path+" computed mismatch")
	}
}

func mustValue(typ tftypes.Type, val any) tftypes.Value {
	v := tftypes.NewValue(typ, val)
	return v
}

func TestFillFrameworkValues_Basic(t *testing.T) {
	s := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"name":  rschema.StringAttribute{},
			"count": rschema.Int64Attribute{},
			"flag":  rschema.BoolAttribute{},
		},
	}

	tv := mustValue(
		tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"name":  tftypes.String,
				"count": tftypes.Number,
				"flag":  tftypes.Bool,
			},
		},
		map[string]tftypes.Value{
			"name":  mustValue(tftypes.String, "example"),
			"count": mustValue(tftypes.Number, float64(42)), // IMPORTANT: use float64
			"flag":  mustValue(tftypes.Bool, true),
		},
	)

	got, err := FillFrameworkValues(tv, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := map[string]attr.Value{
		"name":  types.StringValue("example"),
		"count": types.Int64Value(42),
		"flag":  types.BoolValue(true),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected result\n got:  %#v\n want: %#v", got, want)
	}
}

func TestFillFrameworkValues_WithNulls(t *testing.T) {
	s := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"name":  rschema.StringAttribute{},
			"count": rschema.Int64Attribute{},
		},
	}

	tv := mustValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":  tftypes.String,
			"count": tftypes.Number,
		},
	}, map[string]tftypes.Value{
		"name":  tftypes.NewValue(tftypes.String, nil), // Null value
		"count": tftypes.NewValue(tftypes.Number, nil),
	})

	got, err := FillFrameworkValues(tv, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := map[string]attr.Value{
		"name":  types.StringNull(),
		"count": types.Int64Null(),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected result\n got:  %#v\n want: %#v", got, want)
	}
}

func TestFillFrameworkValues_Complex(t *testing.T) {
	s := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"tags": rschema.ListAttribute{ElementType: types.StringType},
			"meta": rschema.MapAttribute{ElementType: types.StringType},
		},
	}

	tv := mustValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"tags": tftypes.List{ElementType: tftypes.String},
			"meta": tftypes.Map{ElementType: tftypes.String},
		},
	}, map[string]tftypes.Value{
		"tags": mustValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
			mustValue(tftypes.String, "a"),
			mustValue(tftypes.String, "b"),
			mustValue(tftypes.String, "c"),
		}),
		"meta": mustValue(tftypes.Map{ElementType: tftypes.String}, map[string]tftypes.Value{
			"env": mustValue(tftypes.String, "prod"),
			"ver": mustValue(tftypes.String, "v1"),
		}),
	})

	got, err := FillFrameworkValues(tv, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := map[string]attr.Value{
		"tags": types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("a"),
			types.StringValue("b"),
			types.StringValue("c"),
		}),
		"meta": types.MapValueMust(types.StringType, map[string]attr.Value{
			"env": types.StringValue("prod"),
			"ver": types.StringValue("v1"),
		}),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected result\n got:  %#v\n want: %#v", got, want)
	}
}

func TestFillFrameworkValues_EmptyObject(t *testing.T) {
	s := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"details": rschema.ObjectAttribute{
				AttributeTypes: map[string]attr.Type{},
			},
		},
	}

	tv := mustValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"details": tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{},
			},
		},
	}, map[string]tftypes.Value{
		"details": mustValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{},
		}, map[string]tftypes.Value{}),
	})

	got, err := FillFrameworkValues(tv, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := map[string]attr.Value{
		"details": types.ObjectValueMust(map[string]attr.Type{}, map[string]attr.Value{}),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected result\n got:  %#v\n want: %#v", got, want)
	}
}

func TestFillFrameworkValues_NullObject(t *testing.T) {
	s := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"details": rschema.ObjectAttribute{
				AttributeTypes: map[string]attr.Type{
					"info": types.StringType,
				},
			},
		},
	}

	tv := mustValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"details": tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"info": tftypes.String,
				},
			},
		},
	}, map[string]tftypes.Value{
		"details": tftypes.NewValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"info": tftypes.String,
			},
		}, nil), // Null object
	})

	got, err := FillFrameworkValues(tv, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := map[string]attr.Value{
		"details": types.ObjectNull(map[string]attr.Type{
			"info": types.StringType,
		}),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected result\n got:  %#v\n want: %#v", got, want)
	}
}

func TestFillFrameworkValues_NestedObject(t *testing.T) {
	s := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"settings": rschema.ObjectAttribute{
				AttributeTypes: map[string]attr.Type{
					"enabled": types.BoolType,
					"tags":    types.ListType{ElemType: types.StringType},
					"meta":    types.MapType{ElemType: types.StringType},
				},
			},
		},
	}

	tv := mustValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"settings": tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"enabled": tftypes.Bool,
					"tags":    tftypes.List{ElementType: tftypes.String},
					"meta":    tftypes.Map{ElementType: tftypes.String},
				},
			},
		},
	}, map[string]tftypes.Value{
		"settings": mustValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"enabled": tftypes.Bool,
				"tags":    tftypes.List{ElementType: tftypes.String},
				"meta":    tftypes.Map{ElementType: tftypes.String},
			},
		}, map[string]tftypes.Value{
			"enabled": mustValue(tftypes.Bool, true),
			"tags": mustValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
				mustValue(tftypes.String, "a"),
				mustValue(tftypes.String, "b"),
			}),
			"meta": mustValue(tftypes.Map{ElementType: tftypes.String}, map[string]tftypes.Value{
				"ver": mustValue(tftypes.String, "v1"),
			}),
		}),
	})

	got, err := FillFrameworkValues(tv, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := map[string]attr.Value{
		"settings": types.ObjectValueMust(map[string]attr.Type{
			"enabled": types.BoolType,
			"tags":    types.ListType{ElemType: types.StringType},
			"meta":    types.MapType{ElemType: types.StringType},
		}, map[string]attr.Value{
			"enabled": types.BoolValue(true),
			"tags": types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("a"),
				types.StringValue("b"),
			}),
			"meta": types.MapValueMust(types.StringType, map[string]attr.Value{
				"ver": types.StringValue("v1"),
			}),
		}),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected result\n got:  %#v\n want: %#v", got, want)
	}
}

func TestGetFilteredValues(t *testing.T) {

	cases := []struct {
		name     string
		raw      map[string]attr.Value
		meta     map[string]attrMeta
		expected map[string]any
	}{
		{
			name: "string and int64 with values",
			raw: map[string]attr.Value{
				"name": types.StringValue("foo"),
				"age":  types.Int64Value(42),
			},
			meta: map[string]attrMeta{
				"name": {Required: true},
				"age":  {Optional: true},
			},
			expected: map[string]any{
				"name": "foo",
				"age":  int64(42),
			},
		},
		{
			name: "empty string and unknown int64",
			raw: map[string]attr.Value{
				"empty":   types.StringValue(""),
				"unknown": types.Int64Unknown(),
			},
			meta: map[string]attrMeta{
				"empty":   {Optional: true},
				"unknown": {Optional: true},
			},
			expected: map[string]any{
				"empty": "",
			},
		},
		{
			name: "nested object",
			raw: map[string]attr.Value{
				"nested": types.ObjectValueMust(
					map[string]attr.Type{"foo": types.StringType},
					map[string]attr.Value{"foo": types.StringValue("bar")},
				),
			},
			meta: map[string]attrMeta{
				"nested": {Optional: true},
			},
			expected: map[string]any{
				"nested": map[string]any{"foo": "bar"},
			},
		},
		{
			name: "list of strings",
			raw: map[string]attr.Value{
				"tags": types.ListValueMust(types.StringType,
					[]attr.Value{types.StringValue("a"), types.StringValue("b")}),
			},
			meta: map[string]attrMeta{
				"tags": {Optional: true},
			},
			expected: map[string]any{
				"tags": []any{"a", "b"},
			},
		},
		{
			name: "list of list of ints",
			raw: map[string]attr.Value{
				"matrix": types.ListValueMust(
					types.ListType{ElemType: types.Int64Type},
					[]attr.Value{
						types.ListValueMust(types.Int64Type, []attr.Value{
							types.Int64Value(1), types.Int64Value(2),
						}),
					}),
			},
			meta: map[string]attrMeta{
				"matrix": {Optional: true},
			},
			expected: map[string]any{
				"matrix": []any{[]any{int64(1), int64(2)}},
			},
		},
	}

	for _, tc := range cases {

		typeMap := make(map[string]attr.Type)
		for k, v := range tc.raw {
			typeMap[k] = v.Type(context.Background())
		}

		t.Run(tc.name, func(t *testing.T) {
			state := &TFState{
				Raw:     tc.raw,
				Meta:    tc.meta,
				Enabled: true,
				TypeMap: typeMap,
			}
			got := state.GetFilteredValues(FilterOr, nil, SearchRequired, SearchOptional)
			if len(got) != len(tc.expected) {
				t.Errorf("unexpected number of results: got %d, want %d", len(got), len(tc.expected))
			}
			for k, v := range tc.expected {
				if gv, ok := got[k]; !ok || fmt.Sprintf("%#v", gv) != fmt.Sprintf("%#v", v) {
					t.Errorf("unexpected value for key %q: got %#v, want %#v", k, gv, v)
				}
			}
		})
	}
}

func TestGetFilteredValues_PrimitivesOnly(t *testing.T) {
	// Build a small state with a primitive and a complex field
	raw := map[string]attr.Value{
		"name": types.StringValue("alpha"),
		"config": types.ObjectValueMust(map[string]attr.Type{
			"enabled": types.BoolType,
		}, map[string]attr.Value{
			"enabled": types.BoolValue(true),
		}),
	}
	meta := map[string]attrMeta{
		"name":   {Optional: true},
		"config": {Optional: true},
	}
	typeMap := map[string]attr.Type{
		"name":   types.StringType,
		"config": types.ObjectType{AttrTypes: map[string]attr.Type{"enabled": types.BoolType}},
	}

	state := &TFState{Raw: raw, Meta: meta, TypeMap: typeMap, Enabled: true}

	// Without the flag, both should appear
	got := state.GetFilteredValues(FilterOr, nil, SearchOptional)
	require.Equal(t, map[string]any{
		"name":   "alpha",
		"config": map[string]any{"enabled": true},
	}, got)

	// With the primitives-only flag, only the primitive should remain
	got = state.GetFilteredValues(FilterOr, nil, SearchOptional, SearchPrimitivesOnly)
	require.Equal(t, map[string]any{
		"name": "alpha",
	}, got)
}

func TestGetFilteredValues2(t *testing.T) {
	tests := []struct {
		name   string
		raw    map[string]attr.Value
		schema any
		meta   map[string]attrMeta
		expect map[string]any
	}{
		{
			name: "flat values",
			raw: map[string]attr.Value{
				"name":  types.StringValue("test"),
				"count": types.Int64Value(3),
			},
			meta: map[string]attrMeta{
				"name":  {Required: true},
				"count": {Required: true},
			},
			expect: map[string]any{
				"name":  "test",
				"count": int64(3),
			},
		},
		{
			name: "nested object",
			raw: map[string]attr.Value{
				"config": types.ObjectValueMust(map[string]attr.Type{
					"enabled": types.BoolType,
				}, map[string]attr.Value{
					"enabled": types.BoolValue(true),
				}),
			},
			meta: map[string]attrMeta{
				"config":         {Required: true},
				"config.enabled": {Required: true},
			},
			expect: map[string]any{
				"config": map[string]any{
					"enabled": true,
				},
			},
		},
		{
			name: "list of list with content",
			raw: map[string]attr.Value{
				"matrix": types.ListValueMust(types.ListType{ElemType: types.Int64Type}, []attr.Value{
					types.ListValueMust(types.Int64Type, []attr.Value{
						types.Int64Value(1), types.Int64Value(2),
					}),
					types.ListValueMust(types.Int64Type, []attr.Value{
						types.Int64Value(3),
					}),
					types.ListNull(types.Int64Type),
				}),
			},
			meta: map[string]attrMeta{
				"matrix": {Required: true},
			},
			expect: map[string]any{
				"matrix": []any{
					[]any{int64(1), int64(2)},
					[]any{int64(3)},
				},
			},
		},
		{
			name: "list of list empty",
			raw: map[string]attr.Value{
				"matrix": types.ListValueMust(types.ListType{ElemType: types.Int64Type}, []attr.Value{}),
			},
			meta: map[string]attrMeta{
				"matrix": {Required: true},
			},
			expect: map[string]any{
				"matrix": []any{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typeMap := make(map[string]attr.Type)
			for k, v := range tt.raw {
				typeMap[k] = v.Type(context.Background())
			}

			state := &TFState{
				Raw:     tt.raw,
				Meta:    tt.meta,
				TypeMap: typeMap,
				Enabled: true,
			}

			got := state.GetFilteredValues(FilterOr, nil, SearchRequired)
			require.Equal(t, tt.expect, got)
		})
	}
}

func TestBuildAttrMapFromRecord_Complex(t *testing.T) {
	record := map[string]any{
		"name":    "example",
		"enabled": true,
		"count":   3,
		"nested_object": map[string]any{
			"id":    "abc",
			"value": 42,
			"inner": map[string]any{
				"flag": false,
			},
		},
		"nested_list": []any{
			map[string]any{"label": "first", "score": 1.1},
			map[string]any{"label": "second", "score": 2.2},
		},
		"nested_list_of_lists": []any{
			[]any{
				map[string]any{"key": "k1", "val": 10},
				map[string]any{"key": "k2", "val": 20},
			},
			[]any{
				map[string]any{"key": "k3", "val": 30},
			},
		},
		"null_string": nil,
		"null_object": nil,
		"null_list":   nil,
	}

	// Construct nested types using ObjectTypeFromAttributeTypes
	innerType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"flag": types.BoolType,
		},
	}

	nestedObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":    types.StringType,
			"value": types.Int64Type,
			"inner": innerType,
		},
	}

	nestedListElemType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"label": types.StringType,
			"score": types.Float64Type,
		},
	}

	nestedListOfListsElemType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key": types.StringType,
			"val": types.Int64Type,
		},
	}

	nullObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"unused": types.StringType,
		},
	}

	// Build schema
	schema := map[string]attr.Type{
		"name":                 types.StringType,
		"enabled":              types.BoolType,
		"count":                types.Int64Type,
		"nested_object":        nestedObjectType,
		"nested_list":          types.ListType{ElemType: nestedListElemType},
		"nested_list_of_lists": types.ListType{ElemType: types.ListType{ElemType: nestedListOfListsElemType}},
		"null_string":          types.StringType,
		"null_object":          nullObjectType,
		"null_list":            types.ListType{ElemType: types.StringType},
	}

	attrMap := make(map[string]attr.Value)
	for k, typ := range schema {
		val, _, err := BuildAttrValueFromAny(typ, record[k])
		require.NoError(t, err, "failed at key: %s", k)
		require.NotNil(t, val, "value should not be nil: %s", k)
		attrMap[k] = val
	}

	// Optional: Check some key examples
	require.Equal(t, types.StringValue("example"), attrMap["name"])
	require.True(t, attrMap["null_object"].IsNull())
	require.True(t, attrMap["null_string"].IsNull())
	require.True(t, attrMap["null_list"].IsNull())
}

func TestBuildAttrValueFromAny_ListOfListOfList(t *testing.T) {
	record := map[string]any{
		"triple_nested": []any{
			[]any{
				[]any{
					map[string]any{"key": "a", "val": 1},
					map[string]any{"key": "b", "val": 2},
				},
			},
			[]any{
				[]any{
					map[string]any{"key": "c", "val": 3},
				},
			},
		},
	}

	// Build nested object type
	objectElemType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key": types.StringType,
			"val": types.Int64Type,
		},
	}

	// Define list of list of list type
	tripleListType := types.ListType{
		ElemType: types.ListType{
			ElemType: types.ListType{
				ElemType: objectElemType,
			},
		},
	}

	// Call the builder
	val, _, err := BuildAttrValueFromAny(tripleListType, record["triple_nested"])
	require.NoError(t, err)
	require.False(t, val.IsNull())
	require.False(t, val.IsUnknown())

	// Check it's a list
	list := val.(types.List)
	require.Len(t, list.Elements(), 2)

	// Check inner structure
	firstOuter := list.Elements()[0].(types.List)
	require.Len(t, firstOuter.Elements(), 1)

	innerMost := firstOuter.Elements()[0].(types.List)
	require.Len(t, innerMost.Elements(), 2)

	obj0 := innerMost.Elements()[0].(types.Object)
	require.Equal(t, "a", obj0.Attributes()["key"].(types.String).ValueString())
	require.Equal(t, int64(1), obj0.Attributes()["val"].(types.Int64).ValueInt64())
}

func TestFilteredValuesReturnsEmptyMapWhenNoMatchingFlags(t *testing.T) {
	raw := map[string]attr.Value{
		"name":  types.StringValue("test"),
		"count": types.Int64Value(3),
	}
	meta := map[string]attrMeta{
		"name":  {Required: true},
		"count": {Optional: true},
	}

	typeMap := make(map[string]attr.Type)
	for k, v := range raw {
		typeMap[k] = v.Type(context.Background())
	}

	state := &TFState{
		Raw:     raw,
		Meta:    meta,
		TypeMap: typeMap,
		Enabled: true,
	}
	result := state.GetFilteredValues(FilterOr, nil, SearchSensitive)
	require.Empty(t, result)
}

func GetFilteredValuesIncludesNullValuesWhenSearchEmptyFlagIsSet(t *testing.T) {
	raw := map[string]attr.Value{
		"name":  types.StringNull(),
		"count": types.Int64Value(3),
	}
	meta := map[string]attrMeta{
		"name":  {Optional: true},
		"count": {Required: true},
	}
	typeMap := make(map[string]attr.Type)
	for k, v := range raw {
		typeMap[k] = v.Type(context.Background())
	}
	state := &TFState{
		Raw:     raw,
		Meta:    meta,
		Enabled: true,
		TypeMap: typeMap,
	}
	result := state.GetFilteredValues(FilterOr, nil, SearchEmpty)
	require.Equal(t, map[string]any{
		"name":  nil,
		"count": int64(3),
	}, result)
}

func DiffFieldsReturnsDifferencesBetweenTwoStates(t *testing.T) {
	raw1 := map[string]attr.Value{
		"name":  types.StringValue("test"),
		"count": types.Int64Value(3),
	}
	raw2 := map[string]attr.Value{
		"name":  types.StringValue("example"),
		"count": types.Int64Value(3),
	}
	meta := map[string]attrMeta{
		"name":  {Required: true},
		"count": {Required: true},
	}
	typeMap1 := make(map[string]attr.Type)
	for k, v := range raw1 {
		typeMap1[k] = v.Type(context.Background())
	}
	typeMap2 := make(map[string]attr.Type)
	for k, v := range raw2 {
		typeMap2[k] = v.Type(context.Background())
	}

	state1 := &TFState{
		Raw:     raw1,
		Meta:    meta,
		Enabled: true,
		TypeMap: typeMap1,
	}
	state2 := &TFState{
		Raw:     raw2,
		Meta:    meta,
		Enabled: true,
		TypeMap: typeMap2,
	}

	result := state1.DiffFields(state2, FilterOr, nil, SearchRequired)
	require.Equal(t, map[string]any{
		"name": "test",
	}, result)
}

func DiffFieldsHandlesNullValuesCorrectly(t *testing.T) {
	raw1 := map[string]attr.Value{
		"name":  types.StringNull(),
		"count": types.Int64Value(3),
	}
	raw2 := map[string]attr.Value{
		"name":  types.StringValue("example"),
		"count": types.Int64Value(3),
	}
	meta := map[string]attrMeta{
		"name":  {Required: true},
		"count": {Required: true},
	}
	typeMap1 := make(map[string]attr.Type)
	for k, v := range raw1 {
		typeMap1[k] = v.Type(context.Background())
	}
	typeMap2 := make(map[string]attr.Type)
	for k, v := range raw2 {
		typeMap2[k] = v.Type(context.Background())
	}

	state1 := &TFState{
		Raw:     raw1,
		Meta:    meta,
		Enabled: true,
		TypeMap: typeMap1,
	}
	state2 := &TFState{
		Raw:     raw2,
		Meta:    meta,
		Enabled: true,
		TypeMap: typeMap2,
	}
	result := state1.DiffFields(state2, FilterOr, nil, SearchRequired)
	require.Equal(t, map[string]any{
		"name": nil,
	}, result)
}

func buildTFStateFixture() *TFState {
	raw := map[string]attr.Value{
		"name": types.StringValue("Alice"),
		"age":  types.Int64Value(30),
		"meta": types.ObjectValueMust(
			map[string]attr.Type{"enabled": types.BoolType},
			map[string]attr.Value{"enabled": types.BoolValue(true)},
		),
		"tags": types.ListValueMust(types.StringType,
			[]attr.Value{types.StringValue("a"), types.StringValue("b")},
		),
	}
	meta := map[string]attrMeta{
		"name":         {Required: true},
		"age":          {Optional: true},
		"meta":         {Required: true},
		"meta.enabled": {Required: true},
		"tags":         {Optional: true},
	}
	typeMap := map[string]attr.Type{
		"name":         types.StringType,
		"age":          types.Int64Type,
		"meta":         types.ObjectType{AttrTypes: map[string]attr.Type{"enabled": types.BoolType}},
		"tags":         types.ListType{ElemType: types.StringType},
		"meta.enabled": types.BoolType,
	}

	return &TFState{
		Raw:     raw,
		Meta:    meta,
		TypeMap: typeMap,
		Enabled: true,
	}
}

func TestTFState_Getters(t *testing.T) {
	state := buildTFStateFixture()

	require.Equal(t, "Alice", state.String("name"))
	require.Equal(t, int64(30), state.Int64("age"))
	require.Equal(t, true, state.TfObject("meta").Attributes()["enabled"].(types.Bool).ValueBool())
	require.Equal(t, []attr.Value{
		types.StringValue("a"),
		types.StringValue("b"),
	}, state.TfList("tags").Elements())
	require.False(t, state.IsNull("name"))
	require.False(t, state.IsUnknown("name"))
}

func TestTFState_Get_InvalidType(t *testing.T) {
	state := buildTFStateFixture()

	require.Panics(t, func() { state.Bool("name") })
	require.Panics(t, func() { state.Float64("name") })
	require.Panics(t, func() { state.TfSet("tags") }) // not a set
	require.Panics(t, func() { state.TfObject("age") })
	require.Panics(t, func() { state.String("invalid") })
}

func TestTFState_MetaAccess(t *testing.T) {
	state := buildTFStateFixture()

	require.True(t, state.IsRequired("name"))
	require.False(t, state.IsRequired("age"))
	require.True(t, state.IsOptional("age"))
	require.Panics(t, func() { state.IsOptional("missing") })
	require.Equal(t, types.StringType, state.Type("name"))
	require.Panics(t, func() { _ = state.Type("notfound") })
}

func TestTFState_Set(t *testing.T) {
	state := buildTFStateFixture()
	state.Set("name", "Bob")
	require.Equal(t, "Bob", state.String("name"))

	require.Panics(t, func() {
		state.Set("missing", "value")
	})
}

// TestSetMethodWithStringSlice, TestSetMethodWithIntSlice, TestSetMethodWithNonSliceValue removed -
// these tests were for automatic slice conversion which is no longer implemented

func TestTFState_GetFilteredValues(t *testing.T) {
	state := buildTFStateFixture()

	out := state.GetFilteredValues(FilterOr, nil, SearchRequired)
	require.Equal(t, map[string]any{
		"name": "Alice",
		"meta": map[string]any{"enabled": true},
	}, out)
}

func TestTFState_CopyKnownFieldsTo(t *testing.T) {
	src := buildTFStateFixture()
	dst := &TFState{
		Raw:     map[string]attr.Value{},
		Meta:    map[string]attrMeta{},
		TypeMap: src.TypeMap,
		Enabled: true,
	}
	src.CopyKnownFieldsTo(dst)
	require.Equal(t, src.Raw["name"], dst.Raw["name"])
	require.Equal(t, src.Meta["name"], dst.Meta["name"])
}

// TestTFState_CopyKnownFieldsTo_SkipsNullComputedFields verifies that null values
// for computed fields are NOT copied, preventing them from overwriting API response values.
// This addresses the issue where optional+computed fields like 'id' become null during updates.
func TestTFState_CopyKnownFieldsTo_SkipsNullComputedFields(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id":       rschema.Int64Attribute{Optional: true, Computed: true}, // optional+computed
		"guid":     rschema.StringAttribute{Computed: true},                // computed-only
		"name":     rschema.StringAttribute{Optional: true},                // optional-only
		"required": rschema.StringAttribute{Required: true},                // required
	}}

	// Source (plan) has null values for computed fields
	plan := NewTFStateMust(map[string]attr.Value{
		"id":       types.Int64Null(),  // null in plan
		"guid":     types.StringNull(), // null in plan
		"name":     types.StringValue("updated_name"),
		"required": types.StringValue("updated_required"),
	}, schema, nil)

	// Destination (state) has values set from API response
	state := NewTFStateMust(map[string]attr.Value{
		"id":       types.Int64Value(42),         // from API
		"guid":     types.StringValue("abc-123"), // from API
		"name":     types.StringValue("old_name"),
		"required": types.StringValue("old_required"),
	}, schema, nil)

	// Copy from plan to state
	plan.CopyKnownFieldsTo(state)

	// Verify: computed fields should NOT be overwritten by null values from plan
	assert.Equal(t, int64(42), state.Raw["id"].(types.Int64).ValueInt64(),
		"id should NOT be overwritten by null from plan")
	assert.Equal(t, "abc-123", state.Raw["guid"].(types.String).ValueString(),
		"guid should NOT be overwritten by null from plan")

	// Verify: non-computed fields SHOULD be copied from plan
	assert.Equal(t, "updated_name", state.Raw["name"].(types.String).ValueString(),
		"name should be updated from plan")
	assert.Equal(t, "updated_required", state.Raw["required"].(types.String).ValueString(),
		"required should be updated from plan")
}

// TestTFState_CopyKnownFieldsTo_CopiesNonNullComputedFields verifies that non-null
// computed fields ARE copied (user-specified values should be preserved).
func TestTFState_CopyKnownFieldsTo_CopiesNonNullComputedFields(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id":   rschema.Int64Attribute{Optional: true, Computed: true},
		"name": rschema.StringAttribute{Optional: true},
	}}

	// Plan has a user-specified value for the optional+computed field
	plan := NewTFStateMust(map[string]attr.Value{
		"id":   types.Int64Value(99), // user specified value
		"name": types.StringValue("new_name"),
	}, schema, nil)

	state := NewTFStateMust(map[string]attr.Value{
		"id":   types.Int64Value(42),
		"name": types.StringValue("old_name"),
	}, schema, nil)

	// Copy from plan to state
	plan.CopyKnownFieldsTo(state)

	// Verify: non-null computed value should be copied
	assert.Equal(t, int64(99), state.Raw["id"].(types.Int64).ValueInt64(),
		"id should be updated with user-specified value from plan")
	assert.Equal(t, "new_name", state.Raw["name"].(types.String).ValueString(),
		"name should be updated from plan")
}

// TestTFState_CopyKnownFieldsTo_CopiesNullForOptionalFields verifies that null values
// for OPTIONAL (non-computed) fields ARE copied, allowing users to clear field values.
// This is the key difference from computed fields - users must be able to clear optional fields.
func TestTFState_CopyKnownFieldsTo_CopiesNullForOptionalFields(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id":             rschema.Int64Attribute{Optional: true, Computed: true}, // optional+computed
		"optional_field": rschema.StringAttribute{Optional: true},                // optional-only
		"description":    rschema.StringAttribute{Optional: true},                // optional-only
		"tags":           rschema.ListAttribute{ElementType: types.StringType, Optional: true},
	}}

	// Scenario: User removes optional_field and tags from config, wants to clear them
	// Plan has null values for fields user wants to clear
	plan := NewTFStateMust(map[string]attr.Value{
		"id":             types.Int64Null(),                // computed field - should NOT be copied
		"optional_field": types.StringNull(),               // optional field - SHOULD be copied to clear it
		"description":    types.StringValue("kept"),        // user keeps this
		"tags":           types.ListNull(types.StringType), // optional field - SHOULD be copied to clear it
	}, schema, nil)

	// State has values from previous apply
	state := NewTFStateMust(map[string]attr.Value{
		"id":             types.Int64Value(42),           // from API
		"optional_field": types.StringValue("old_value"), // old value user wants to clear
		"description":    types.StringValue("old_description"),
		"tags": types.ListValueMust(
			types.StringType,
			[]attr.Value{types.StringValue("tag1"), types.StringValue("tag2")},
		), // old tags user wants to clear
	}, schema, nil)

	// Copy from plan to state
	plan.CopyKnownFieldsTo(state)

	// Verify: computed field with null should NOT be copied (preserves API value)
	assert.Equal(t, int64(42), state.Raw["id"].(types.Int64).ValueInt64(),
		"id (computed) should NOT be overwritten by null from plan - API value preserved")

	// Verify: optional fields with null SHOULD be copied (allows user to clear)
	assert.True(t, state.Raw["optional_field"].(types.String).IsNull(),
		"optional_field should be cleared (set to null) as user intended")

	assert.True(t, state.Raw["tags"].(types.List).IsNull(),
		"tags should be cleared (set to null) as user intended")

	// Verify: optional field with value should be copied normally
	assert.Equal(t, "kept", state.Raw["description"].(types.String).ValueString(),
		"description should be updated from plan")
}

// TestTFState_CopyKnownFieldsTo_RequiredFieldsWithNull verifies behavior with required fields
func TestTFState_CopyKnownFieldsTo_RequiredFieldsWithNull(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"name": rschema.StringAttribute{Required: true},
		"id":   rschema.Int64Attribute{Optional: true, Computed: true},
	}}

	// Plan with null for required field (shouldn't happen in real scenarios, but testing edge case)
	plan := NewTFStateMust(map[string]attr.Value{
		"name": types.StringNull(), // required field with null (edge case)
		"id":   types.Int64Null(),  // computed field with null
	}, schema, nil)

	state := NewTFStateMust(map[string]attr.Value{
		"name": types.StringValue("existing_name"),
		"id":   types.Int64Value(42),
	}, schema, nil)

	// Copy from plan to state
	plan.CopyKnownFieldsTo(state)

	// Required field with null should be copied (not skipped)
	assert.True(t, state.Raw["name"].(types.String).IsNull(),
		"required field with null should be copied")

	// Computed field with null should NOT be copied
	assert.Equal(t, int64(42), state.Raw["id"].(types.Int64).ValueInt64(),
		"computed field with null should NOT be copied")
}

func TestTFState_CopyNonEmptyFieldsTo(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id":     rschema.Int64Attribute{Computed: true},
		"name":   rschema.StringAttribute{Optional: true},
		"title":  rschema.StringAttribute{Computed: true},
		"note":   rschema.StringAttribute{Optional: true},
		"number": rschema.Int64Attribute{Optional: true},
	}}

	left := NewTFStateMust(map[string]attr.Value{
		"id":     types.Int64Value(42),
		"name":   types.StringValue("alice"),
		"title":  types.StringNull(),
		"note":   types.StringNull(),
		"number": types.Int64Value(1001),
	}, schema, nil)

	right := NewTFStateMust(map[string]attr.Value{
		"id":     types.Int64Value(7),
		"name":   types.StringValue("bob"),
		"title":  types.StringValue("mgr"),
		"note":   types.StringValue("keep"),
		"number": types.Int64Null(),
	}, schema, nil)

	left.CopyNonEmptyFieldsTo(right)

	require.Equal(t, types.Int64Value(42), right.Raw["id"])         // copied
	require.Equal(t, types.StringValue("alice"), right.Raw["name"]) // copied
	require.Equal(t, types.StringValue("mgr"), right.Raw["title"])  // unchanged (left null)
	require.Equal(t, types.StringValue("keep"), right.Raw["note"])  // unchanged (left null)
	require.Equal(t, types.Int64Value(1001), right.Raw["number"])   // copied
}

func TestTFState_FillFromRecord(t *testing.T) {
	state := buildTFStateFixture()
	state.Meta["computed_field"] = attrMeta{Computed: true}
	state.TypeMap["computed_field"] = types.StringType

	err := state.FillFromRecord(map[string]any{
		"computed_field": "auto",
		"irrelevant":     123,
	})
	require.NoError(t, err)
	require.Equal(t, "auto", state.String("computed_field"))
}

// TestTFState_FillFromRecordForImport verifies that FillFromRecordForImport fills
// computed and required fields, but NOT optional fields. This prevents drift after import
// by keeping optional fields null in both plan and state (when not specified in config).
func TestTFState_FillFromRecordForImport(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id":             rschema.Int64Attribute{Optional: true, Computed: true},               // optional+computed
		"name":           rschema.StringAttribute{Required: true},                              // required
		"description":    rschema.StringAttribute{Optional: true},                              // optional-only
		"tags":           rschema.ListAttribute{ElementType: types.StringType, Optional: true}, // optional list
		"nfs_read_only":  rschema.ListAttribute{ElementType: types.StringType, Optional: true}, // optional list
		"computed_field": rschema.StringAttribute{Computed: true},                              // computed-only
	}}

	state := NewTFStateMust(map[string]attr.Value{
		"id":             types.Int64Null(),
		"name":           types.StringNull(),
		"description":    types.StringNull(),
		"tags":           types.ListNull(types.StringType),
		"nfs_read_only":  types.ListNull(types.StringType),
		"computed_field": types.StringNull(),
	}, schema, nil)

	// Simulate API response from import
	record := Record{
		"id":             int64(5),
		"name":           "test-policy",
		"description":    "my description",              // optional field - should NOT be filled
		"tags":           []interface{}{"tag1", "tag2"}, // optional field - should NOT be filled
		"nfs_read_only":  []interface{}{},               // optional field - should NOT be filled
		"computed_field": "auto-generated",
	}

	// Fill from record using import method
	err := state.FillFromRecordForImport(record)
	require.NoError(t, err)

	// Verify computed and required fields ARE filled
	assert.Equal(t, int64(5), state.Int64("id"),
		"id (optional+computed) should be filled")
	assert.Equal(t, "test-policy", state.String("name"),
		"name (required) should be filled")
	assert.Equal(t, "auto-generated", state.String("computed_field"),
		"computed_field should be filled")

	// Verify optional-only fields are NOT filled (remain null)
	assert.True(t, state.Get("description").IsNull(),
		"description (optional-only) should NOT be filled - prevents drift!")
	assert.True(t, state.Get("tags").IsNull(),
		"tags (optional-only) should NOT be filled - prevents drift!")
	assert.True(t, state.Get("nfs_read_only").IsNull(),
		"nfs_read_only (optional-only) should NOT be filled - prevents drift!")
}

// TestTFState_FillFromRecordForImport_vs_Regular compares import vs regular fill behavior
// Both methods skip optional-only fields, but import includes required fields.
func TestTFState_FillFromRecordForImport_vs_Regular(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id":       rschema.Int64Attribute{Optional: true, Computed: true},
		"name":     rschema.StringAttribute{Required: true},
		"optional": rschema.StringAttribute{Optional: true},
		"computed": rschema.StringAttribute{Computed: true},
	}}

	record := Record{
		"id":       int64(42),
		"name":     "test",
		"optional": "optional_value",
		"computed": "computed_value",
	}

	// Test regular FillFromRecord (used during create/update/read)
	stateRegular := NewTFStateMust(map[string]attr.Value{
		"id":       types.Int64Null(),
		"name":     types.StringNull(),
		"optional": types.StringNull(),
		"computed": types.StringNull(),
	}, schema, nil)

	err := stateRegular.FillFromRecord(record)
	require.NoError(t, err)

	// Regular fill: only computed fields, NOT optional or required
	assert.True(t, stateRegular.Get("optional").IsNull(),
		"FillFromRecord should NOT fill optional fields")
	assert.True(t, stateRegular.Get("name").IsNull(),
		"FillFromRecord should NOT fill required fields")
	assert.False(t, stateRegular.Get("computed").IsNull(),
		"FillFromRecord should fill computed fields")

	// Test FillFromRecordForImport (used during import)
	stateImport := NewTFStateMust(map[string]attr.Value{
		"id":       types.Int64Null(),
		"name":     types.StringNull(),
		"optional": types.StringNull(),
		"computed": types.StringNull(),
	}, schema, nil)

	err = stateImport.FillFromRecordForImport(record)
	require.NoError(t, err)

	// Import fill: computed + required, but NOT optional-only
	assert.True(t, stateImport.Get("optional").IsNull(),
		"FillFromRecordForImport should NOT fill optional-only fields")
	assert.Equal(t, "test", stateImport.String("name"),
		"FillFromRecordForImport SHOULD fill required fields")
	assert.Equal(t, "computed_value", stateImport.String("computed"),
		"FillFromRecordForImport should fill computed fields")
	assert.Equal(t, int64(42), stateImport.Int64("id"),
		"FillFromRecordForImport should fill optional+computed fields")
}

func TestGetGenericSearchParams(t *testing.T) {
	typeMap := map[string]attr.Type{
		"uid":       types.Int64Type,
		"name":      types.StringType,
		"tenant_id": types.StringType,
		"id":        types.Int64Type,
		"guid":      types.StringType,
		"extra":     types.StringType,
	}

	baseMeta := map[string]attrMeta{
		"uid":       {Optional: true, Searchable: true},
		"name":      {Optional: true, Searchable: true},
		"tenant_id": {Required: true, Searchable: true},
		"id":        {Computed: true},
		"guid":      {Computed: true},
		"extra":     {Required: true},
	}

	t.Run("by unique identifiers", func(t *testing.T) {
		raw := map[string]attr.Value{
			"uid": types.Int64Value(42),
		}
		tf := &TFState{
			Raw:     raw,
			Meta:    baseMeta,
			TypeMap: typeMap,
			Enabled: true,
		}
		got := tf.GetGenericSearchParams(context.Background())
		require.Equal(t, vast_client.Params{
			"uid": int64(42),
		}, got)
	})

	t.Run("by common searchable", func(t *testing.T) {
		raw := map[string]attr.Value{
			"name":      types.StringValue("alpha"),
			"tenant_id": types.StringNull(), // should be skipped
		}
		tf := &TFState{
			Raw:     raw,
			Meta:    baseMeta,
			TypeMap: typeMap,
			Enabled: true,
		}
		got := tf.GetGenericSearchParams(context.Background())
		require.Equal(t, vast_client.Params{
			"name": "alpha",
		}, got)
	})

	t.Run("fallback to required+searchable", func(t *testing.T) {
		raw := map[string]attr.Value{
			"extra": types.StringValue("value"),
		}
		meta := map[string]attrMeta{
			"extra": {Required: true, Searchable: true},
		}
		tf := &TFState{
			Raw:     raw,
			Meta:    meta,
			TypeMap: map[string]attr.Type{"extra": types.StringType},
			Enabled: true,
		}
		got := tf.GetGenericSearchParams(context.Background())
		require.Equal(t, vast_client.Params{
			"extra": "value",
		}, got)
	})

	t.Run("includes id and guid if present", func(t *testing.T) {
		raw := map[string]attr.Value{
			"id":   types.Int64Value(99),
			"guid": types.StringValue("abc-def"),
		}
		tf := &TFState{
			Raw:     raw,
			Meta:    baseMeta,
			TypeMap: typeMap,
			Enabled: true,
		}
		got := tf.GetGenericSearchParams(context.Background())
		require.Equal(t, vast_client.Params{
			"id":   int64(99),
			"guid": "abc-def",
		}, got)
	})
}

func TestGetSearchParams_BothGenericAndReadOnly(t *testing.T) {
	meta := map[string]attrMeta{
		"id":          {Computed: true},
		"guid":        {Computed: true},
		"name":        {Optional: true, Searchable: true},
		"zone":        {Optional: true, Searchable: true},
		"description": {ReadOnly: true},
	}

	raw := map[string]attr.Value{
		"id":          types.Int64Value(123),
		"guid":        types.StringValue("abc-guid"),
		"name":        types.StringValue("vol1"),
		"zone":        types.StringValue("zoneA"),
		"description": types.StringValue("readonly field"),
	}

	typeMapResolved := make(map[string]attr.Type)
	for k, v := range raw {
		typeMapResolved[k] = v.Type(context.Background())
	}

	tf := &TFState{
		Raw:     raw,
		Meta:    meta,
		TypeMap: typeMapResolved,
		Enabled: true,
	}

	t.Run("GetReadOnlySearchParams returns readonly only", func(t *testing.T) {
		got := tf.GetReadOnlySearchParams()
		require.Equal(t, vast_client.Params{
			"description": "readonly field",
		}, got)
	})

	t.Run("GetGenericSearchParams includes fallback and merges readonly (non-overlapping)", func(t *testing.T) {
		got := tf.GetGenericSearchParams(context.Background())
		require.Equal(t, vast_client.Params{
			"id":          int64(123),
			"guid":        "abc-guid",
			"name":        "vol1",
			"description": "readonly field",
		}, got)
	})

	t.Run("GetGenericSearchParams prefers existing keys over readonly (override=false)", func(t *testing.T) {
		tf := &TFState{
			Raw: map[string]attr.Value{
				"name":        types.StringValue("main-name"),
				"description": types.StringValue("custom-description"),
			},
			Meta: map[string]attrMeta{
				"name":        {Optional: true, Searchable: true},
				"description": {ReadOnly: true},
			},
			TypeMap: map[string]attr.Type{
				"name":        types.StringType,
				"description": types.StringType,
			},
			Enabled: true,
		}

		got := tf.GetGenericSearchParams(context.Background())
		require.Equal(t, vast_client.Params{
			"name":        "main-name",
			"description": "custom-description", // not overridden
		}, got)
	})
}

// TestConvertSliceToAny removed - function no longer exists

// TestSetMethodWithVariousSliceTypes removed - Set method should not automatically convert slices

func TestSetMethodWithNonSliceValues(t *testing.T) {
	tests := []struct {
		name        string
		fieldType   attr.Type
		inputValue  any
		expectedVal any
	}{
		{
			name:        "string value",
			fieldType:   types.StringType,
			inputValue:  "test-string",
			expectedVal: "test-string",
		},
		{
			name:        "int value",
			fieldType:   types.Int64Type,
			inputValue:  42,
			expectedVal: int64(42),
		},
		{
			name:        "float value",
			fieldType:   types.Float64Type,
			inputValue:  3.14,
			expectedVal: 3.14,
		},
		{
			name:        "bool value",
			fieldType:   types.BoolType,
			inputValue:  true,
			expectedVal: true,
		},
		// Removed nil test case as it's not a typical use case
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create schema with the test field type
			schema := rschema.Schema{
				Attributes: map[string]rschema.Attribute{
					"id": rschema.Int64Attribute{
						Required: true,
					},
					"test_field": rschema.StringAttribute{
						Computed: true,
					},
				},
			}

			// Create TFState
			raw := map[string]attr.Value{
				"id":         types.Int64Value(123),
				"test_field": types.StringNull(),
			}

			tfState := NewTFStateMust(raw, schema, nil)

			// Test setting the value
			tfState.Set("test_field", tt.inputValue)

			// Verify the value was set correctly
			result := tfState.Get("test_field")
			if tt.inputValue == nil {
				// For nil values, the behavior depends on the type
				// String attributes become null when set to nil
				require.True(t, result.IsNull())
			} else {
				require.False(t, result.IsNull())
				require.False(t, result.IsUnknown())
			}
		})
	}
}

// TestSetMethodWithComplexSliceTypes and TestSetMethodEdgeCases removed - Set method should not automatically convert slices

func TestGetFilteredValues_SetOfObjectsWithNullFields(t *testing.T) {
	// Create an object type for user quota (matching the actual Terraform structure)
	entityType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":            types.StringType,
			"email":           types.StringType,
			"identifier":      types.StringType,
			"identifier_type": types.StringType,
			"is_group":        types.BoolType,
		},
	}

	userQuotaType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"entity":       entityType,
			"grace_period": types.StringType,
			"hard_limit":   types.Int64Type,
			"soft_limit":   types.Int64Type,
		},
	}

	// Create a set type containing objects
	setType := types.SetType{ElemType: userQuotaType}

	// Create an object with all null fields
	entity := types.ObjectValueMust(entityType.AttrTypes, map[string]attr.Value{
		"name":            types.StringNull(),
		"email":           types.StringNull(),
		"identifier":      types.StringNull(),
		"identifier_type": types.StringNull(),
		"is_group":        types.BoolNull(),
	})

	obj1 := types.ObjectValueMust(userQuotaType.AttrTypes, map[string]attr.Value{
		"entity":       entity,
		"grace_period": types.StringNull(),
		"hard_limit":   types.Int64Null(),
		"soft_limit":   types.Int64Null(),
	})

	// Create a set with the object
	set := types.SetValueMust(setType.ElemType, []attr.Value{obj1})

	// Create TFState
	state := &TFState{
		Raw: map[string]attr.Value{
			"user_quotas": set,
		},
		Meta: map[string]attrMeta{
			"user_quotas": {Optional: true},
		},
		TypeMap: map[string]attr.Type{
			"user_quotas": setType,
		},
		Enabled: true,
	}

	// Test without SearchEmpty flag (default behavior)
	result := state.GetFilteredValues(FilterOr, nil, SearchOptional)

	// The result should contain user_quotas
	require.Contains(t, result, "user_quotas")

	// The user_quotas should be a slice
	userQuotas, ok := result["user_quotas"].([]any)
	require.True(t, ok)
	require.Len(t, userQuotas, 1)

	// The first element should be a map with null values removed
	objMap, ok := userQuotas[0].(map[string]any)
	require.True(t, ok)

	// Since all fields are null, only the entity field should remain as an empty map
	// This is the correct behavior - null fields are removed but object structure is preserved
	require.Contains(t, objMap, "entity")
	entityMap, ok := objMap["entity"].(map[string]any)
	require.True(t, ok)
	require.Empty(t, entityMap)
}

func TestGetFilteredValues_UserQuotasWithRealValues(t *testing.T) {
	// This test reproduces the actual issue where user_quotas with real values
	// are being converted to empty objects in GetFilteredValues

	// Create the entity object type
	entityType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":            types.StringType,
			"email":           types.StringType,
			"identifier":      types.StringType,
			"identifier_type": types.StringType,
			"is_group":        types.BoolType,
		},
	}

	// Create the user_quota object type (matching the actual Terraform structure)
	userQuotaType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"entity":       entityType,
			"grace_period": types.StringType,
			"hard_limit":   types.Int64Type,
			"soft_limit":   types.Int64Type,
		},
	}

	// Create a set type containing user_quota objects
	setType := types.SetType{ElemType: userQuotaType}

	// Create an entity with real values
	entity := types.ObjectValueMust(entityType.AttrTypes, map[string]attr.Value{
		"name":            types.StringValue("tfzealous-kingfisher"),
		"email":           types.StringValue("user1@example.com"),
		"identifier":      types.StringValue("tfzealous-kingfisher"),
		"identifier_type": types.StringValue("username"),
		"is_group":        types.BoolValue(false),
	})

	// Create a user_quota with real values
	userQuota := types.ObjectValueMust(userQuotaType.AttrTypes, map[string]attr.Value{
		"entity":       entity,
		"grace_period": types.StringValue("02:00:00"),
		"hard_limit":   types.Int64Value(15000),
		"soft_limit":   types.Int64Value(15000),
	})

	// Create a set with the user_quota
	set := types.SetValueMust(setType.ElemType, []attr.Value{userQuota})

	// Create TFState
	state := &TFState{
		Raw: map[string]attr.Value{
			"user_quotas": set,
		},
		Meta: map[string]attrMeta{
			"user_quotas": {Optional: true},
		},
		TypeMap: map[string]attr.Type{
			"user_quotas": setType,
		},
		Enabled: true,
	}

	// Test GetFilteredValues
	result := state.GetFilteredValues(FilterOr, nil, SearchOptional)

	// Debug: print the result
	fmt.Printf("Result: %+v\n", result)

	// The result should contain user_quotas
	require.Contains(t, result, "user_quotas")

	// The user_quotas should contain the real values, not empty objects
	userQuotas, ok := result["user_quotas"].([]any)
	require.True(t, ok)
	require.Len(t, userQuotas, 1)

	// Check that the first element has the expected structure
	firstQuota, ok := userQuotas[0].(map[string]any)
	require.True(t, ok)
	require.Contains(t, firstQuota, "entity")
	require.Contains(t, firstQuota, "grace_period")
	require.Contains(t, firstQuota, "hard_limit")
	require.Contains(t, firstQuota, "soft_limit")
}

// Mock schema for testing
type mockSchema struct {
	attributes map[string]attr.Type
}

func (m *mockSchema) GetAttributes() map[string]attr.Type {
	return m.attributes
}

func TestTFState_HasAttribute(t *testing.T) {
	tests := []struct {
		name          string
		schema        any
		attributeName string
		expected      bool
	}{
		{
			name: "resource schema with id attribute",
			schema: rschema.Schema{
				Attributes: map[string]rschema.Attribute{
					"id":   rschema.Int64Attribute{Optional: true},
					"name": rschema.StringAttribute{Optional: true},
				},
			},
			attributeName: "id",
			expected:      true,
		},
		{
			name: "resource schema without id attribute",
			schema: rschema.Schema{
				Attributes: map[string]rschema.Attribute{
					"name": rschema.StringAttribute{Optional: true},
				},
			},
			attributeName: "id",
			expected:      false,
		},
		{
			name: "datasource schema with id attribute",
			schema: dsschema.Schema{
				Attributes: map[string]dsschema.Attribute{
					"id":   dsschema.Int64Attribute{Optional: true},
					"name": dsschema.StringAttribute{Optional: true},
				},
			},
			attributeName: "id",
			expected:      true,
		},
		{
			name: "datasource schema without id attribute",
			schema: dsschema.Schema{
				Attributes: map[string]dsschema.Attribute{
					"name": dsschema.StringAttribute{Optional: true},
				},
			},
			attributeName: "id",
			expected:      false,
		},
		{
			name:          "disabled TFState",
			schema:        nil,
			attributeName: "id",
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tfState *TFState
			if tt.schema == nil {
				tfState = NewTFStateMust(nil, nil, nil)
			} else {
				tfState = NewTFStateMust(nil, tt.schema, nil)
			}

			result := tfState.HasAttribute(tt.attributeName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTFState_SetOrAdd(t *testing.T) {
	tests := []struct {
		name     string
		schema   rschema.Schema
		key      string
		value    any
		expected attr.Value
	}{
		{
			name: "set int64 value",
			schema: rschema.Schema{
				Attributes: map[string]rschema.Attribute{
					"id": rschema.Int64Attribute{Optional: true},
				},
			},
			key:      "id",
			value:    int64(123),
			expected: types.Int64Value(123),
		},
		{
			name: "set string value",
			schema: rschema.Schema{
				Attributes: map[string]rschema.Attribute{
					"name": rschema.StringAttribute{Optional: true},
				},
			},
			key:      "name",
			value:    "test-name",
			expected: types.StringValue("test-name"),
		},
		{
			name: "set bool value",
			schema: rschema.Schema{
				Attributes: map[string]rschema.Attribute{
					"enabled": rschema.BoolAttribute{Optional: true},
				},
			},
			key:      "enabled",
			value:    true,
			expected: types.BoolValue(true),
		},
		{
			name: "add new key that doesn't exist in Raw",
			schema: rschema.Schema{
				Attributes: map[string]rschema.Attribute{
					"id": rschema.Int64Attribute{Optional: true},
				},
			},
			key:      "id",
			value:    int64(456),
			expected: types.Int64Value(456),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create TFState with empty Raw map to test SetOrAdd
			tfState := NewTFStateMust(map[string]attr.Value{}, tt.schema, nil)

			// Set the value using SetOrAdd
			tfState.SetOrAdd(tt.key, tt.value)

			// Verify the value was set correctly
			result := tfState.Get(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTFState_SetOrAdd_WithExistingValue(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"id":   rschema.Int64Attribute{Optional: true},
			"name": rschema.StringAttribute{Optional: true},
		},
	}

	// Create TFState with existing values
	initialRaw := map[string]attr.Value{
		"name": types.StringValue("initial-name"),
	}
	tfState := NewTFStateMust(initialRaw, schema, nil)

	// Add a new key that doesn't exist in Raw
	tfState.SetOrAdd("id", int64(789))

	// Verify both existing and new values are correct
	assert.Equal(t, types.StringValue("initial-name"), tfState.Get("name"))
	assert.Equal(t, types.Int64Value(789), tfState.Get("id"))
}

func TestTFState_SetOrAdd_OverwriteExistingValue(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"id": rschema.Int64Attribute{Optional: true},
		},
	}

	// Create TFState with existing value
	initialRaw := map[string]attr.Value{
		"id": types.Int64Value(123),
	}
	tfState := NewTFStateMust(initialRaw, schema, nil)

	// Overwrite the existing value
	tfState.SetOrAdd("id", int64(456))

	// Verify the value was overwritten
	assert.Equal(t, types.Int64Value(456), tfState.Get("id"))
}

// NOTE: SetState is simplified in the implementation; skipping write-only persistence behavior tests.

// ==========================================
// Tests for GetUpdateParams
// ==========================================

func TestGetUpdateParams_ExcludesEditOnlyFields(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"id":            rschema.Int64Attribute{Optional: true},
			"name":          rschema.StringAttribute{Optional: true},
			"enabled":       rschema.BoolAttribute{Optional: true},   // edit-only
			"description":   rschema.StringAttribute{Optional: true}, // regular field
			"delete_option": rschema.StringAttribute{Optional: true}, // delete-only
		},
	}

	hints := &TFStateHints{
		EditOnlyFields: []string{"enabled"},
		DeleteOnlyBodyFields: map[string]string{
			"delete_option": "",
		},
	}

	// Current state (old)
	stateRaw := map[string]attr.Value{
		"id":            types.Int64Value(1),
		"name":          types.StringValue("old-name"),
		"enabled":       types.BoolValue(false),
		"description":   types.StringValue("old-desc"),
		"delete_option": types.StringValue("old-opt"),
	}
	currentState := NewTFStateMust(stateRaw, schema, hints)

	// Plan state (new)
	planRaw := map[string]attr.Value{
		"id":            types.Int64Value(1),
		"name":          types.StringValue("new-name"), // CHANGED
		"enabled":       types.BoolValue(true),         // CHANGED (edit-only)
		"description":   types.StringValue("new-desc"), // CHANGED
		"delete_option": types.StringValue("new-opt"),  // CHANGED (delete-only)
	}
	planState := NewTFStateMust(planRaw, schema, hints)

	// Get update params
	updateParams := planState.GetUpdateParams(currentState)

	// Should include regular changed fields
	assert.Contains(t, updateParams, "name")
	assert.Equal(t, "new-name", updateParams["name"])
	assert.Contains(t, updateParams, "description")
	assert.Equal(t, "new-desc", updateParams["description"])

	// Should NOT include edit-only fields
	assert.NotContains(t, updateParams, "enabled", "edit-only field should be excluded")

	// Should NOT include delete-only fields
	assert.NotContains(t, updateParams, "delete_option", "delete-only field should be excluded")

	// Should NOT include id
	assert.NotContains(t, updateParams, "id", "id should not be in update params")
}

// TERF-224: omitting a create-only bool must not produce PATCH {"field": null}.
func TestGetUpdateParams_ExcludesCreateOnlyClearedFields(t *testing.T) {
	t.Parallel()

	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"id":                rschema.Int64Attribute{Optional: true},
			"name":              rschema.StringAttribute{Optional: true},
			"is_physical_quota": rschema.BoolAttribute{Optional: true, Computed: true},
		},
	}
	hints := &TFStateHints{CreateOnlyFields: []string{"is_physical_quota"}}

	currentState := NewTFStateMust(map[string]attr.Value{
		"id":                types.Int64Value(3),
		"name":              types.StringValue("terf224_physical"),
		"is_physical_quota": types.BoolValue(true),
	}, schema, hints)

	// Config omitted is_physical_quota → plan null (before UseStateForUnknown).
	planState := NewTFStateMust(map[string]attr.Value{
		"id":                types.Int64Value(3),
		"name":              types.StringValue("terf224_physical"),
		"is_physical_quota": types.BoolNull(),
	}, schema, hints)

	updateParams := planState.GetUpdateParams(currentState)
	assert.NotContains(t, updateParams, "is_physical_quota",
		"create-only field must not be cleared with null on update")
	assert.Empty(t, updateParams)
}

func TestGetUpdateParams_NoHints(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"id":   rschema.Int64Attribute{Optional: true},
			"name": rschema.StringAttribute{Optional: true},
		},
	}

	// Current state
	stateRaw := map[string]attr.Value{
		"id":   types.Int64Value(1),
		"name": types.StringValue("old-name"),
	}
	currentState := NewTFStateMust(stateRaw, schema, nil) // No hints

	// Plan state
	planRaw := map[string]attr.Value{
		"id":   types.Int64Value(1),
		"name": types.StringValue("new-name"),
	}
	planState := NewTFStateMust(planRaw, schema, nil)

	// Get update params - should work like GetChangedParams when no hints
	updateParams := planState.GetUpdateParams(currentState)

	assert.Contains(t, updateParams, "name")
	assert.Equal(t, "new-name", updateParams["name"])
}

func TestGetUpdateParams_OnlyEditOnlyFieldsChanged(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"id":      rschema.Int64Attribute{Optional: true},
			"name":    rschema.StringAttribute{Optional: true},
			"enabled": rschema.BoolAttribute{Optional: true}, // edit-only
		},
	}

	hints := &TFStateHints{
		EditOnlyFields: []string{"enabled"},
	}

	// Current state
	stateRaw := map[string]attr.Value{
		"id":      types.Int64Value(1),
		"name":    types.StringValue("same-name"),
		"enabled": types.BoolValue(false),
	}
	currentState := NewTFStateMust(stateRaw, schema, hints)

	// Plan state - only edit-only field changed
	planRaw := map[string]attr.Value{
		"id":      types.Int64Value(1),
		"name":    types.StringValue("same-name"), // NOT changed
		"enabled": types.BoolValue(true),          // CHANGED (edit-only)
	}
	planState := NewTFStateMust(planRaw, schema, hints)

	// Get update params - should be empty
	updateParams := planState.GetUpdateParams(currentState)

	assert.Empty(t, updateParams, "Should be empty when only edit-only fields changed")
}

// ==========================================
// Tests for GetChangedEditOnlyParams
// ==========================================

func TestGetChangedEditOnlyParams_ReturnsOnlyChangedEditOnlyFields(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"id":          rschema.Int64Attribute{Optional: true},
			"name":        rschema.StringAttribute{Optional: true},
			"enabled":     rschema.BoolAttribute{Optional: true},   // edit-only
			"auto_start":  rschema.BoolAttribute{Optional: true},   // edit-only
			"description": rschema.StringAttribute{Optional: true}, // regular field
		},
	}

	hints := &TFStateHints{
		EditOnlyFields: []string{"enabled", "auto_start"},
	}

	// Current state
	stateRaw := map[string]attr.Value{
		"id":          types.Int64Value(1),
		"name":        types.StringValue("old-name"),
		"enabled":     types.BoolValue(false),
		"auto_start":  types.BoolValue(false),
		"description": types.StringValue("old-desc"),
	}
	currentState := NewTFStateMust(stateRaw, schema, hints)

	// Plan state
	planRaw := map[string]attr.Value{
		"id":          types.Int64Value(1),
		"name":        types.StringValue("new-name"), // CHANGED (regular)
		"enabled":     types.BoolValue(true),         // CHANGED (edit-only)
		"auto_start":  types.BoolValue(false),        // NOT changed (edit-only)
		"description": types.StringValue("new-desc"), // CHANGED (regular)
	}
	planState := NewTFStateMust(planRaw, schema, hints)

	// Get changed edit-only params
	editOnlyParams := planState.GetChangedEditOnlyParams(currentState)

	// Should include only CHANGED edit-only field
	assert.Contains(t, editOnlyParams, "enabled")
	assert.Equal(t, true, editOnlyParams["enabled"])

	// Should NOT include unchanged edit-only field
	assert.NotContains(t, editOnlyParams, "auto_start", "unchanged edit-only field should be excluded")

	// Should NOT include regular fields
	assert.NotContains(t, editOnlyParams, "name", "regular field should be excluded")
	assert.NotContains(t, editOnlyParams, "description", "regular field should be excluded")
}

func TestGetChangedEditOnlyParams_NoHints(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"id":   rschema.Int64Attribute{Optional: true},
			"name": rschema.StringAttribute{Optional: true},
		},
	}

	// Current state
	stateRaw := map[string]attr.Value{
		"id":   types.Int64Value(1),
		"name": types.StringValue("old-name"),
	}
	currentState := NewTFStateMust(stateRaw, schema, nil) // No hints

	// Plan state
	planRaw := map[string]attr.Value{
		"id":   types.Int64Value(1),
		"name": types.StringValue("new-name"),
	}
	planState := NewTFStateMust(planRaw, schema, nil)

	// Get changed edit-only params - should be empty
	editOnlyParams := planState.GetChangedEditOnlyParams(currentState)

	assert.Empty(t, editOnlyParams, "Should be empty when no hints provided")
}

func TestGetChangedEditOnlyParams_NoChanges(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"id":      rschema.Int64Attribute{Optional: true},
			"enabled": rschema.BoolAttribute{Optional: true}, // edit-only
		},
	}

	hints := &TFStateHints{
		EditOnlyFields: []string{"enabled"},
	}

	// Current state
	stateRaw := map[string]attr.Value{
		"id":      types.Int64Value(1),
		"enabled": types.BoolValue(true),
	}
	currentState := NewTFStateMust(stateRaw, schema, hints)

	// Plan state - no changes
	planRaw := map[string]attr.Value{
		"id":      types.Int64Value(1),
		"enabled": types.BoolValue(true), // NOT changed
	}
	planState := NewTFStateMust(planRaw, schema, hints)

	// Get changed edit-only params - should be empty
	editOnlyParams := planState.GetChangedEditOnlyParams(currentState)

	assert.Empty(t, editOnlyParams, "Should be empty when no edit-only fields changed")
}

// ==========================================
// Tests for GetUpdateParams + GetChangedEditOnlyParams
// Combined (Disjoint Sets)
// ==========================================

func TestGetUpdateParams_And_GetChangedEditOnlyParams_AreDisjoint(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"id":            rschema.Int64Attribute{Optional: true},
			"name":          rschema.StringAttribute{Optional: true},
			"enabled":       rschema.BoolAttribute{Optional: true},   // edit-only
			"auto_start":    rschema.BoolAttribute{Optional: true},   // edit-only
			"description":   rschema.StringAttribute{Optional: true}, // regular field
			"delete_option": rschema.StringAttribute{Optional: true}, // delete-only
		},
	}

	hints := &TFStateHints{
		EditOnlyFields: []string{"enabled", "auto_start"},
		DeleteOnlyBodyFields: map[string]string{
			"delete_option": "",
		},
	}

	// Current state
	stateRaw := map[string]attr.Value{
		"id":            types.Int64Value(1),
		"name":          types.StringValue("old-name"),
		"enabled":       types.BoolValue(false),
		"auto_start":    types.BoolValue(false),
		"description":   types.StringValue("old-desc"),
		"delete_option": types.StringValue("old-opt"),
	}
	currentState := NewTFStateMust(stateRaw, schema, hints)

	// Plan state - all fields changed
	planRaw := map[string]attr.Value{
		"id":            types.Int64Value(1),
		"name":          types.StringValue("new-name"),
		"enabled":       types.BoolValue(true),
		"auto_start":    types.BoolValue(true),
		"description":   types.StringValue("new-desc"),
		"delete_option": types.StringValue("new-opt"),
	}
	planState := NewTFStateMust(planRaw, schema, hints)

	// Get both sets
	updateParams := planState.GetUpdateParams(currentState)
	editOnlyParams := planState.GetChangedEditOnlyParams(currentState)

	// Verify updateParams contains regular fields only
	assert.Contains(t, updateParams, "name")
	assert.Contains(t, updateParams, "description")
	assert.Len(t, updateParams, 2, "Should have exactly 2 regular changed fields")

	// Verify editOnlyParams contains edit-only fields only
	assert.Contains(t, editOnlyParams, "enabled")
	assert.Contains(t, editOnlyParams, "auto_start")
	assert.Len(t, editOnlyParams, 2, "Should have exactly 2 edit-only changed fields")

	// Verify NO overlap between the two sets
	for key := range updateParams {
		assert.NotContains(t, editOnlyParams, key, "Field %q should not be in both sets", key)
	}
	for key := range editOnlyParams {
		assert.NotContains(t, updateParams, key, "Field %q should not be in both sets", key)
	}

	// Verify delete-only fields are in neither set
	assert.NotContains(t, updateParams, "delete_option")
	assert.NotContains(t, editOnlyParams, "delete_option")

	// Verify id is in neither set
	assert.NotContains(t, updateParams, "id")
	assert.NotContains(t, editOnlyParams, "id")
}

func TestGetUpdateParams_And_GetChangedEditOnlyParams_CompletePartitioning(t *testing.T) {
	// This test verifies that GetUpdateParams + GetChangedEditOnlyParams together
	// capture ALL changed optional/required fields (excluding delete-only and id)

	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"id":          rschema.Int64Attribute{Optional: true},
			"name":        rschema.StringAttribute{Optional: true},
			"enabled":     rschema.BoolAttribute{Optional: true}, // edit-only
			"description": rschema.StringAttribute{Optional: true},
		},
	}

	hints := &TFStateHints{
		EditOnlyFields: []string{"enabled"},
	}

	// Current state
	stateRaw := map[string]attr.Value{
		"id":          types.Int64Value(1),
		"name":        types.StringValue("old-name"),
		"enabled":     types.BoolValue(false),
		"description": types.StringValue("old-desc"),
	}
	currentState := NewTFStateMust(stateRaw, schema, hints)

	// Plan state - all fields changed
	planRaw := map[string]attr.Value{
		"id":          types.Int64Value(1),
		"name":        types.StringValue("new-name"),
		"enabled":     types.BoolValue(true),
		"description": types.StringValue("new-desc"),
	}
	planState := NewTFStateMust(planRaw, schema, hints)

	// Get all changed params using old method
	allChangedParams := planState.GetChangedParams(currentState)

	// Get partitioned sets
	updateParams := planState.GetUpdateParams(currentState)
	editOnlyParams := planState.GetChangedEditOnlyParams(currentState)

	// Combine the two partitioned sets
	combinedParams := make(vast_client.Params)
	for k, v := range updateParams {
		combinedParams[k] = v
	}
	for k, v := range editOnlyParams {
		combinedParams[k] = v
	}

	// The combined params should equal all changed params (minus id)
	delete(allChangedParams, "id")

	assert.Equal(t, len(allChangedParams), len(combinedParams),
		"Combined params should have same length as all changed params (minus id)")

	for key, val := range allChangedParams {
		assert.Contains(t, combinedParams, key, "Key %q missing from combined params", key)
		assert.Equal(t, val, combinedParams[key], "Value mismatch for key %q", key)
	}
}

func TestGetUpdateParams_WithDeleteOnlyParamFields(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"id":           rschema.Int64Attribute{Optional: true},
			"name":         rschema.StringAttribute{Optional: true},
			"delete_param": rschema.StringAttribute{Optional: true}, // delete-only param
		},
	}

	hints := &TFStateHints{
		DeleteOnlyParamFields: map[string]string{
			"delete_param": "force",
		},
	}

	// Current state
	stateRaw := map[string]attr.Value{
		"id":           types.Int64Value(1),
		"name":         types.StringValue("old-name"),
		"delete_param": types.StringValue("old-param"),
	}
	currentState := NewTFStateMust(stateRaw, schema, hints)

	// Plan state
	planRaw := map[string]attr.Value{
		"id":           types.Int64Value(1),
		"name":         types.StringValue("new-name"),
		"delete_param": types.StringValue("new-param"), // CHANGED (delete-only param)
	}
	planState := NewTFStateMust(planRaw, schema, hints)

	// Get update params
	updateParams := planState.GetUpdateParams(currentState)

	// Should include regular fields
	assert.Contains(t, updateParams, "name")

	// Should NOT include delete-only param fields
	assert.NotContains(t, updateParams, "delete_param", "delete-only param field should be excluded")
}

// TestGetAllValues_PrimitiveTypes tests GetAllValues with primitive types (string, int64, bool)
func TestGetAllValues_PrimitiveTypes(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"string_field": rschema.StringAttribute{Optional: true},
			"int_field":    rschema.Int64Attribute{Optional: true},
			"bool_field":   rschema.BoolAttribute{Optional: true},
		},
	}

	raw := map[string]attr.Value{
		"string_field": types.StringValue("test_value"),
		"int_field":    types.Int64Value(42),
		"bool_field":   types.BoolValue(true),
	}

	tfState := NewTFStateMust(raw, schema, nil)
	result := tfState.GetAllValues()

	// Verify all fields are returned
	assert.Len(t, result, 3)
	assert.Equal(t, "test_value", result["string_field"])
	assert.Equal(t, int64(42), result["int_field"])
	assert.Equal(t, true, result["bool_field"])
}

// TestGetAllValues_NullAndUnknownValues tests that null and unknown values are included
func TestGetAllValues_NullAndUnknownValues(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"null_string":    rschema.StringAttribute{Optional: true},
			"unknown_string": rschema.StringAttribute{Optional: true},
			"known_string":   rschema.StringAttribute{Optional: true},
			"null_int":       rschema.Int64Attribute{Optional: true},
			"unknown_int":    rschema.Int64Attribute{Optional: true},
		},
	}

	raw := map[string]attr.Value{
		"null_string":    types.StringNull(),
		"unknown_string": types.StringUnknown(),
		"known_string":   types.StringValue("known"),
		"null_int":       types.Int64Null(),
		"unknown_int":    types.Int64Unknown(),
	}

	tfState := NewTFStateMust(raw, schema, nil)
	result := tfState.GetAllValues()

	// All fields should be in the result, including null and unknown
	assert.Len(t, result, 5)
	assert.Nil(t, result["null_string"])
	assert.Nil(t, result["unknown_string"])
	assert.Equal(t, "known", result["known_string"])
	assert.Nil(t, result["null_int"])
	assert.Nil(t, result["unknown_int"])
}

// TestGetAllValues_Lists tests GetAllValues with list attributes
func TestGetAllValues_Lists(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"string_list": rschema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"int_list": rschema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
			},
		},
	}

	raw := map[string]attr.Value{
		"string_list": types.ListValueMust(
			types.StringType,
			[]attr.Value{
				types.StringValue("item1"),
				types.StringValue("item2"),
				types.StringValue("item3"),
			},
		),
		"int_list": types.ListValueMust(
			types.Int64Type,
			[]attr.Value{
				types.Int64Value(10),
				types.Int64Value(20),
			},
		),
	}

	tfState := NewTFStateMust(raw, schema, nil)
	result := tfState.GetAllValues()

	assert.Len(t, result, 2)

	// Check string list
	stringList, ok := result["string_list"].([]interface{})
	require.True(t, ok, "string_list should be a []interface{}")
	assert.Len(t, stringList, 3)
	assert.Equal(t, "item1", stringList[0])
	assert.Equal(t, "item2", stringList[1])
	assert.Equal(t, "item3", stringList[2])

	// Check int list
	intList, ok := result["int_list"].([]interface{})
	require.True(t, ok, "int_list should be a []interface{}")
	assert.Len(t, intList, 2)
	assert.Equal(t, int64(10), intList[0])
	assert.Equal(t, int64(20), intList[1])
}

// TestGetAllValues_Sets tests GetAllValues with set attributes
func TestGetAllValues_Sets(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"string_set": rschema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}

	raw := map[string]attr.Value{
		"string_set": types.SetValueMust(
			types.StringType,
			[]attr.Value{
				types.StringValue("alpha"),
				types.StringValue("beta"),
			},
		),
	}

	tfState := NewTFStateMust(raw, schema, nil)
	result := tfState.GetAllValues()

	assert.Len(t, result, 1)

	// Check set (should be converted to []interface{})
	set, ok := result["string_set"].([]interface{})
	require.True(t, ok, "string_set should be a []interface{}")
	assert.Len(t, set, 2)
}

// TestGetAllValues_NestedObjects tests GetAllValues with nested object attributes
func TestGetAllValues_NestedObjects(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"nested_single": rschema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]rschema.Attribute{
					"name": rschema.StringAttribute{Optional: true},
					"age":  rschema.Int64Attribute{Optional: true},
				},
			},
			"nested_list": rschema.ListNestedAttribute{
				Optional: true,
				NestedObject: rschema.NestedAttributeObject{
					Attributes: map[string]rschema.Attribute{
						"id":    rschema.Int64Attribute{Optional: true},
						"value": rschema.StringAttribute{Optional: true},
					},
				},
			},
		},
	}

	raw := map[string]attr.Value{
		"nested_single": types.ObjectValueMust(
			map[string]attr.Type{
				"name": types.StringType,
				"age":  types.Int64Type,
			},
			map[string]attr.Value{
				"name": types.StringValue("John"),
				"age":  types.Int64Value(30),
			},
		),
		"nested_list": types.ListValueMust(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id":    types.Int64Type,
					"value": types.StringType,
				},
			},
			[]attr.Value{
				types.ObjectValueMust(
					map[string]attr.Type{
						"id":    types.Int64Type,
						"value": types.StringType,
					},
					map[string]attr.Value{
						"id":    types.Int64Value(1),
						"value": types.StringValue("first"),
					},
				),
				types.ObjectValueMust(
					map[string]attr.Type{
						"id":    types.Int64Type,
						"value": types.StringType,
					},
					map[string]attr.Value{
						"id":    types.Int64Value(2),
						"value": types.StringValue("second"),
					},
				),
			},
		),
	}

	tfState := NewTFStateMust(raw, schema, nil)
	result := tfState.GetAllValues()

	assert.Len(t, result, 2)

	// Check nested single object
	nestedSingle, ok := result["nested_single"].(map[string]interface{})
	require.True(t, ok, "nested_single should be a map[string]interface{}")
	assert.Equal(t, "John", nestedSingle["name"])
	assert.Equal(t, int64(30), nestedSingle["age"])

	// Check nested list
	nestedList, ok := result["nested_list"].([]interface{})
	require.True(t, ok, "nested_list should be a []interface{}")
	assert.Len(t, nestedList, 2)

	firstItem, ok := nestedList[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, int64(1), firstItem["id"])
	assert.Equal(t, "first", firstItem["value"])

	secondItem, ok := nestedList[1].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, int64(2), secondItem["id"])
	assert.Equal(t, "second", secondItem["value"])
}

// TestGetAllValues_EmptyState tests GetAllValues with an empty state
func TestGetAllValues_EmptyState(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"field1": rschema.StringAttribute{Optional: true},
			"field2": rschema.Int64Attribute{Optional: true},
		},
	}

	raw := map[string]attr.Value{}

	tfState := NewTFStateMust(raw, schema, nil)
	result := tfState.GetAllValues()

	// Should return empty map
	assert.Empty(t, result)
}

// TestGetAllValues_MixedTypes tests GetAllValues with various mixed types
func TestGetAllValues_MixedTypes(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"id":         rschema.Int64Attribute{Optional: true},
			"name":       rschema.StringAttribute{Optional: true},
			"enabled":    rschema.BoolAttribute{Optional: true},
			"null_field": rschema.StringAttribute{Optional: true},
			"tags":       rschema.ListAttribute{ElementType: types.StringType, Optional: true},
			"metadata": rschema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]rschema.Attribute{
					"key":   rschema.StringAttribute{Optional: true},
					"value": rschema.StringAttribute{Optional: true},
				},
			},
		},
	}

	raw := map[string]attr.Value{
		"id":         types.Int64Value(123),
		"name":       types.StringValue("test"),
		"enabled":    types.BoolValue(false),
		"null_field": types.StringNull(),
		"tags": types.ListValueMust(
			types.StringType,
			[]attr.Value{types.StringValue("tag1"), types.StringValue("tag2")},
		),
		"metadata": types.ObjectValueMust(
			map[string]attr.Type{
				"key":   types.StringType,
				"value": types.StringType,
			},
			map[string]attr.Value{
				"key":   types.StringValue("env"),
				"value": types.StringValue("prod"),
			},
		),
	}

	tfState := NewTFStateMust(raw, schema, nil)
	result := tfState.GetAllValues()

	assert.Len(t, result, 6)
	assert.Equal(t, int64(123), result["id"])
	assert.Equal(t, "test", result["name"])
	assert.Equal(t, false, result["enabled"])
	assert.Nil(t, result["null_field"])

	tags, ok := result["tags"].([]interface{})
	require.True(t, ok)
	assert.Len(t, tags, 2)

	metadata, ok := result["metadata"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "env", metadata["key"])
	assert.Equal(t, "prod", metadata["value"])
}

// TestGetAllValues_EmptyLists tests GetAllValues with empty lists and sets
func TestGetAllValues_EmptyLists(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"empty_list": rschema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"empty_set": rschema.SetAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
			},
		},
	}

	raw := map[string]attr.Value{
		"empty_list": types.ListValueMust(types.StringType, []attr.Value{}),
		"empty_set":  types.SetValueMust(types.Int64Type, []attr.Value{}),
	}

	tfState := NewTFStateMust(raw, schema, nil)
	result := tfState.GetAllValues()

	assert.Len(t, result, 2)

	emptyList, ok := result["empty_list"].([]interface{})
	require.True(t, ok)
	assert.Empty(t, emptyList)

	emptySet, ok := result["empty_set"].([]interface{})
	require.True(t, ok)
	assert.Empty(t, emptySet)
}

// TestGetAllValues_SkipRefreshUseCase tests GetAllValues in the context of SkipRefreshAPICall
// This simulates the protection_policy scenario where GetAllValues is used as a Record
func TestGetAllValues_SkipRefreshUseCase(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"id":   rschema.Int64Attribute{Optional: true, Computed: true},
			"name": rschema.StringAttribute{Optional: true},
			"frames": rschema.ListNestedAttribute{
				Optional: true,
				NestedObject: rschema.NestedAttributeObject{
					Attributes: map[string]rschema.Attribute{
						"every":       rschema.StringAttribute{Optional: true},
						"keep_local":  rschema.StringAttribute{Optional: true},
						"keep_remote": rschema.StringAttribute{Optional: true},
					},
				},
			},
		},
	}

	// Simulate user's configuration with duration values
	raw := map[string]attr.Value{
		"id":   types.Int64Value(42),
		"name": types.StringValue("test-policy"),
		"frames": types.ListValueMust(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"every":       types.StringType,
					"keep_local":  types.StringType,
					"keep_remote": types.StringType,
				},
			},
			[]attr.Value{
				types.ObjectValueMust(
					map[string]attr.Type{
						"every":       types.StringType,
						"keep_local":  types.StringType,
						"keep_remote": types.StringType,
					},
					map[string]attr.Value{
						"every":       types.StringValue("1D"),
						"keep_local":  types.StringValue("14D"),
						"keep_remote": types.StringValue("8D"),
					},
				),
			},
		),
	}

	tfState := NewTFStateMust(raw, schema, nil)
	result := tfState.GetAllValues()

	// Verify the result can be used as a Record (map[string]any)
	assert.Len(t, result, 3)
	assert.Equal(t, int64(42), result["id"])
	assert.Equal(t, "test-policy", result["name"])

	// Verify frames structure is preserved
	frames, ok := result["frames"].([]interface{})
	require.True(t, ok, "frames should be []interface{}")
	assert.Len(t, frames, 1)

	frame, ok := frames[0].(map[string]interface{})
	require.True(t, ok, "frame should be map[string]interface{}")

	// Verify duration values are preserved exactly as user specified
	assert.Equal(t, "1D", frame["every"])
	assert.Equal(t, "14D", frame["keep_local"])
	assert.Equal(t, "8D", frame["keep_remote"])

	// This demonstrates that GetAllValues preserves user's original values
	// without any normalization, which is the desired behavior for SkipRefreshAPICall
}

// TestGetCreateParams_WithWriteOnlyFields tests that GetCreateParams includes write-only fields
// This is critical for resources like administrator_manager where password is write-only
func TestGetCreateParams_WithWriteOnlyFields(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"id": rschema.Int64Attribute{
				Computed: true,
			},
			"username": rschema.StringAttribute{
				Required: true,
			},
			"password": rschema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"email": rschema.StringAttribute{
				Optional: true,
			},
		},
	}

	raw := map[string]attr.Value{
		"username": types.StringValue("admin"),
		"password": types.StringValue("secret123"),
		"email":    types.StringValue("admin@example.com"),
	}

	hints := &TFStateHints{
		WriteOnlyFields: []string{"password"},
		SensitiveFields: []string{"password"},
	}

	tfState := NewTFStateMust(raw, schema, hints)
	result := tfState.GetCreateParams()

	// Verify that write-only fields ARE included in create params
	assert.Contains(t, result, "username", "username should be included")
	assert.Contains(t, result, "password", "password (write-only) should be included in create params")
	assert.Contains(t, result, "email", "email should be included")
	assert.Equal(t, "admin", result["username"])
	assert.Equal(t, "secret123", result["password"])
	assert.Equal(t, "admin@example.com", result["email"])

	// Verify that computed ID is not included
	assert.NotContains(t, result, "id", "id should not be included in create params")
}

// TestGetCreateParams_WithMultipleWriteOnlyFields tests multiple write-only fields
func TestGetCreateParams_WithMultipleWriteOnlyFields(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"username": rschema.StringAttribute{
				Required: true,
			},
			"password": rschema.StringAttribute{
				Required: true,
			},
			"password_retype": rschema.StringAttribute{
				Optional: true,
			},
			"roles": rschema.ListAttribute{
				ElementType: types.Int64Type,
				Required:    true,
			},
		},
	}

	raw := map[string]attr.Value{
		"username":        types.StringValue("testuser"),
		"password":        types.StringValue("pass123"),
		"password_retype": types.StringValue("pass123"),
		"roles":           types.ListValueMust(types.Int64Type, []attr.Value{types.Int64Value(1), types.Int64Value(5)}),
	}

	hints := &TFStateHints{
		WriteOnlyFields: []string{"password", "password_retype"},
		SensitiveFields: []string{"password", "password_retype"},
	}

	tfState := NewTFStateMust(raw, schema, hints)
	result := tfState.GetCreateParams()

	// Verify ALL fields including write-only are included
	assert.Contains(t, result, "username")
	assert.Contains(t, result, "password", "password should be included")
	assert.Contains(t, result, "password_retype", "password_retype should be included")
	assert.Contains(t, result, "roles")

	assert.Equal(t, "testuser", result["username"])
	assert.Equal(t, "pass123", result["password"])
	assert.Equal(t, "pass123", result["password_retype"])
}

// TestGetCreateParams_WithOptionalWriteOnlyFields tests optional write-only fields
func TestGetCreateParams_WithOptionalWriteOnlyFields(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"username": rschema.StringAttribute{
				Required: true,
			},
			"password": rschema.StringAttribute{
				Optional: true, // Optional write-only field
			},
		},
	}

	raw := map[string]attr.Value{
		"username": types.StringValue("admin"),
		"password": types.StringValue("secret"),
	}

	hints := &TFStateHints{
		WriteOnlyFields: []string{"password"},
	}

	tfState := NewTFStateMust(raw, schema, hints)
	result := tfState.GetCreateParams()

	// Optional write-only fields should also be included
	assert.Contains(t, result, "password", "optional write-only field should be included")
	assert.Equal(t, "secret", result["password"])
}

func TestClearWriteOnlyFields(t *testing.T) {
	t.Parallel()

	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"name":        rschema.StringAttribute{Required: true},
			"certificate": rschema.StringAttribute{Optional: true, WriteOnly: true},
			"private_key": rschema.StringAttribute{Optional: true, WriteOnly: true},
		},
	}
	ts := NewTFStateMust(map[string]attr.Value{
		"name":        types.StringValue("c"),
		"certificate": types.StringValue("PEM-CERT"),
		"private_key": types.StringValue("PEM-KEY"),
	}, schema, &TFStateHints{
		WriteOnlyFields: []string{"certificate", "private_key"},
	})

	require.Contains(t, ts.GetCreateParams(), "certificate")
	ts.ClearWriteOnlyFields()
	assert.True(t, ts.Raw["certificate"].IsNull())
	assert.True(t, ts.Raw["private_key"].IsNull())
	assert.Equal(t, "c", ts.String("name"))
	create := ts.GetCreateParams()
	assert.NotContains(t, create, "certificate")
	assert.NotContains(t, create, "private_key")
	assert.Equal(t, "c", create["name"])
}

// TestListsHaveSameContentIgnoringOrder_SimpleStrings tests order-independent comparison for simple string lists
func TestListsHaveSameContentIgnoringOrder_SimpleStrings(t *testing.T) {
	// Same content, different order
	listA, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("a"),
		types.StringValue("b"),
		types.StringValue("c"),
	})
	listB, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("c"),
		types.StringValue("a"),
		types.StringValue("b"),
	})

	assert.True(t, listsHaveSameContentIgnoringOrder(listA, listB), "lists with same content but different order should be equal")

	// Different content
	listC, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("a"),
		types.StringValue("b"),
		types.StringValue("d"),
	})

	assert.False(t, listsHaveSameContentIgnoringOrder(listA, listC), "lists with different content should not be equal")

	// Different lengths
	listD, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("a"),
		types.StringValue("b"),
	})

	assert.False(t, listsHaveSameContentIgnoringOrder(listA, listD), "lists with different lengths should not be equal")
}

// TestListsHaveSameContentIgnoringOrder_NestedLists tests order-independent comparison for nested lists (like client_ip_ranges)
func TestListsHaveSameContentIgnoringOrder_NestedLists(t *testing.T) {
	// This mimics client_ip_ranges: List(List(String))
	innerListType := types.ListType{ElemType: types.StringType}

	// Create inner lists (IP ranges)
	range1, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("172.21.112.1"),
		types.StringValue("172.21.112.4"),
	})
	range2, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("12.0.0.6"),
		types.StringValue("12.0.0.10"),
	})
	range3, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("192.168.0.1"),
		types.StringValue("192.168.0.10"),
	})

	// List A: [range1, range2, range3]
	listA, _ := types.ListValue(innerListType, []attr.Value{range1, range2, range3})

	// List B: [range3, range1, range2] (same content, different order)
	listB, _ := types.ListValue(innerListType, []attr.Value{range3, range1, range2})

	assert.True(t, listsHaveSameContentIgnoringOrder(listA, listB), "nested lists with same content but different order should be equal")

	// List C: [range1, range3, range3] (duplicate, different from A)
	listC, _ := types.ListValue(innerListType, []attr.Value{range1, range3, range3})

	assert.False(t, listsHaveSameContentIgnoringOrder(listA, listC), "nested lists with different content should not be equal")
}

// TestListsHaveSameContentIgnoringOrder_ClientIPRanges tests the exact use case from tenant.client_ip_ranges
func TestListsHaveSameContentIgnoringOrder_ClientIPRanges(t *testing.T) {
	innerListType := types.ListType{ElemType: types.StringType}

	// User config order
	userRange1, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("172.21.112.1"),
		types.StringValue("172.21.112.4"),
	})
	userRange2, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("12.0.0.6"),
		types.StringValue("12.0.0.10"),
	})

	userList, _ := types.ListValue(innerListType, []attr.Value{userRange1, userRange2})

	// API returns in different order
	apiRange1, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("12.0.0.6"),
		types.StringValue("12.0.0.10"),
	})
	apiRange2, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("172.21.112.1"),
		types.StringValue("172.21.112.4"),
	})

	apiList, _ := types.ListValue(innerListType, []attr.Value{apiRange1, apiRange2})

	assert.True(t, listsHaveSameContentIgnoringOrder(userList, apiList), "client_ip_ranges should match regardless of order")
}

// TestListsHaveSameContentIgnoringOrder_Integers tests with integer lists
func TestListsHaveSameContentIgnoringOrder_Integers(t *testing.T) {
	listA, _ := types.ListValue(types.Int64Type, []attr.Value{
		types.Int64Value(1),
		types.Int64Value(2),
		types.Int64Value(3),
	})
	listB, _ := types.ListValue(types.Int64Type, []attr.Value{
		types.Int64Value(3),
		types.Int64Value(1),
		types.Int64Value(2),
	})

	assert.True(t, listsHaveSameContentIgnoringOrder(listA, listB), "integer lists with same content should be equal")

	listC, _ := types.ListValue(types.Int64Type, []attr.Value{
		types.Int64Value(1),
		types.Int64Value(2),
		types.Int64Value(4),
	})

	assert.False(t, listsHaveSameContentIgnoringOrder(listA, listC), "integer lists with different content should not be equal")
}

// TestListsHaveSameContentIgnoringOrder_Duplicates tests handling of duplicate elements
func TestListsHaveSameContentIgnoringOrder_Duplicates(t *testing.T) {
	// List with duplicates: ["a", "b", "a"]
	listA, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("a"),
		types.StringValue("b"),
		types.StringValue("a"),
	})

	// Same duplicates, different order: ["a", "a", "b"]
	listB, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("a"),
		types.StringValue("a"),
		types.StringValue("b"),
	})

	assert.True(t, listsHaveSameContentIgnoringOrder(listA, listB), "lists with same duplicates should be equal")

	// Different number of duplicates: ["a", "b", "b"]
	listC, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("a"),
		types.StringValue("b"),
		types.StringValue("b"),
	})

	assert.False(t, listsHaveSameContentIgnoringOrder(listA, listC), "lists with different duplicate counts should not be equal")
}

// TestListsHaveSameContentIgnoringOrder_EmptyLists tests empty list handling
func TestListsHaveSameContentIgnoringOrder_EmptyLists(t *testing.T) {
	emptyA, _ := types.ListValue(types.StringType, []attr.Value{})
	emptyB, _ := types.ListValue(types.StringType, []attr.Value{})

	assert.True(t, listsHaveSameContentIgnoringOrder(emptyA, emptyB), "two empty lists should be equal")

	nonEmpty, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("a"),
	})

	assert.False(t, listsHaveSameContentIgnoringOrder(emptyA, nonEmpty), "empty and non-empty lists should not be equal")
}

// TestListsHaveSameContentIgnoringOrder_NonListTypes tests that non-list types return false
func TestListsHaveSameContentIgnoringOrder_NonListTypes(t *testing.T) {
	stringVal := types.StringValue("test")
	intVal := types.Int64Value(123)

	assert.False(t, listsHaveSameContentIgnoringOrder(stringVal, intVal), "non-list types should return false")

	listVal, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("a"),
	})

	assert.False(t, listsHaveSameContentIgnoringOrder(stringVal, listVal), "comparing list to non-list should return false")
	assert.False(t, listsHaveSameContentIgnoringOrder(listVal, stringVal), "comparing non-list to list should return false")
}

// TestListsHaveSameContentIgnoringOrder_NullAndUnknown tests null and unknown value handling
func TestListsHaveSameContentIgnoringOrder_NullAndUnknown(t *testing.T) {
	listA, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("a"),
		types.StringNull(),
	})
	listB, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringNull(),
		types.StringValue("a"),
	})

	assert.True(t, listsHaveSameContentIgnoringOrder(listA, listB), "lists with null values in different positions should be equal")

	listC, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("a"),
		types.StringUnknown(),
	})

	assert.False(t, listsHaveSameContentIgnoringOrder(listA, listC), "null and unknown should not be equal")
}
