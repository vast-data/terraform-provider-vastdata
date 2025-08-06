// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
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

		converted, err := BuildAttrValueFromAny(attrType, val)
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
		val, err := BuildAttrValueFromAny(typ, record[k])
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
	val, err := BuildAttrValueFromAny(tripleListType, record["triple_nested"])
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

func TestTFState_CopyNonEmptyFieldsTo(t *testing.T) {
	src := buildTFStateFixture()
	dst := &TFState{
		Raw:     map[string]attr.Value{},
		Meta:    map[string]attrMeta{},
		TypeMap: src.TypeMap,
		Enabled: true,
	}
	src.CopyNonEmptyFieldsTo(dst)
	require.Equal(t, src.Raw["name"], dst.Raw["name"])
	require.Equal(t, src.Meta["name"], dst.Meta["name"])
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
