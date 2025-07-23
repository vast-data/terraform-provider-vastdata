// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestBuildAttrValueFromAny_ArbitraryToString(t *testing.T) {
	v := struct{ A string }{A: "val"}
	val, err := BuildAttrValueFromAny(types.StringType, v)
	require.NoError(t, err)
	require.Equal(t, "{val}", val.(types.String).ValueString())
}

func TestBuildAttrValueFromAny_NullValues(t *testing.T) {
	val, err := BuildAttrValueFromAny(types.StringType, nil)
	require.NoError(t, err)
	require.True(t, val.IsNull())

	val, err = BuildAttrValueFromAny(types.Int64Type, nil)
	require.NoError(t, err)
	require.True(t, val.IsNull())

	val, err = BuildAttrValueFromAny(types.ListType{ElemType: types.StringType}, nil)
	require.NoError(t, err)
	require.True(t, val.IsNull())
}

func TestBuildAttrValueFromAny_ListAndSet(t *testing.T) {
	listVal, err := BuildAttrValueFromAny(types.ListType{ElemType: types.StringType}, []any{"a", "b"})
	require.NoError(t, err)
	require.Equal(t, types.ListValueMust(types.StringType, []attr.Value{types.StringValue("a"), types.StringValue("b")}), listVal)

	setVal, err := BuildAttrValueFromAny(types.SetType{ElemType: types.Int64Type}, []any{int64(1), int64(2)})
	require.NoError(t, err)
	require.Equal(t, types.SetValueMust(types.Int64Type, []attr.Value{types.Int64Value(1), types.Int64Value(2)}), setVal)
}

func TestBuildAttrValueFromAny_ObjectMissingField(t *testing.T) {
	objType := types.ObjectType{AttrTypes: map[string]attr.Type{"foo": types.StringType, "bar": types.Int64Type}}
	val, err := BuildAttrValueFromAny(objType, map[string]any{"foo": "baz"})
	require.NoError(t, err)
	require.True(t, val.(types.Object).Attributes()["bar"].IsNull())
}

func TestTfTypeToAttrType_UnsupportedType(t *testing.T) {
	_, err := tfTypeToAttrType(types.ObjectType{AttrTypes: map[string]attr.Type{}}, tftypes.NewValue(tftypes.String, "foo"))
	require.Error(t, err)
}

func TestConvertAttrValueToRaw_NullAndUnknown(t *testing.T) {
	require.Nil(t, ConvertAttrValueToRaw(types.StringNull(), types.StringType))
	require.Nil(t, ConvertAttrValueToRaw(types.Int64Unknown(), types.Int64Type))
}

func TestConvertAttrValueToRaw_ListSetMapObject(t *testing.T) {
	list := types.ListValueMust(types.StringType, []attr.Value{types.StringValue("a"), types.StringValue("b")})
	require.Equal(t, []any{"a", "b"}, ConvertAttrValueToRaw(list, types.ListType{ElemType: types.StringType}))

	set := types.SetValueMust(types.Int64Type, []attr.Value{types.Int64Value(1), types.Int64Value(2)})
	require.ElementsMatch(t, []any{int64(1), int64(2)}, ConvertAttrValueToRaw(set, types.SetType{ElemType: types.Int64Type}))

	m := types.MapValueMust(types.StringType, map[string]attr.Value{"k": types.StringValue("v")})
	require.Equal(t, map[string]any{"k": "v"}, ConvertAttrValueToRaw(m, types.MapType{ElemType: types.StringType}))

	objType := types.ObjectType{AttrTypes: map[string]attr.Type{"foo": types.StringType}}
	obj := types.ObjectValueMust(objType.AttributeTypes(), map[string]attr.Value{"foo": types.StringValue("bar")})
	require.Equal(t, map[string]any{"foo": "bar"}, ConvertAttrValueToRaw(obj, objType))
}

func TestGetNestedAttributes_UnknownType(t *testing.T) {
	require.Nil(t, getNestedAttributes(123))
}

func TestConvertAttrMap_Empty(t *testing.T) {
	require.Empty(t, convertAttrMap(map[string]int{}))
}

func TestBuildAttrValueFromAny_Primitives(t *testing.T) {
	cases := []struct {
		name string
		typ  attr.Type
		val  any
		want attr.Value
	}{
		{"string", types.StringType, "foo", types.StringValue("foo")},
		{"int64", types.Int64Type, int64(42), types.Int64Value(42)},
		{"float64", types.Float64Type, 3.14, types.Float64Value(3.14)},
		{"bool", types.BoolType, true, types.BoolValue(true)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := BuildAttrValueFromAny(tc.typ, tc.val)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestBuildAttrValueFromAny_Null(t *testing.T) {
	val, err := BuildAttrValueFromAny(types.StringType, nil)
	require.NoError(t, err)
	require.True(t, val.IsNull())
}

func TestBuildAttrValueFromAny_List(t *testing.T) {
	listType := types.ListType{ElemType: types.StringType}
	input := []any{"a", "b"}
	val, err := BuildAttrValueFromAny(listType, input)
	require.NoError(t, err)

	list := val.(types.List)
	require.Equal(t, 2, len(list.Elements()))
}

func TestConvertAttrValueToRaw_Object(t *testing.T) {
	val := types.ObjectValueMust(map[string]attr.Type{
		"name": types.StringType,
	}, map[string]attr.Value{
		"name": types.StringValue("Alice"),
	})

	raw := ConvertAttrValueToRaw(val, types.ObjectType{AttrTypes: map[string]attr.Type{
		"name": types.StringType,
	}})

	require.Equal(t, map[string]any{"name": "Alice"}, raw)
}

func TestTfTypeToAttrType_Object(t *testing.T) {
	objType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name": types.StringType,
			"age":  types.Int64Type,
		},
	}

	tfVal := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name": tftypes.String,
			"age":  tftypes.Number,
		},
	}, map[string]tftypes.Value{
		"name": tftypes.NewValue(tftypes.String, "Alice"),
		"age":  tftypes.NewValue(tftypes.Number, big.NewFloat(30)),
	})

	val, err := tfTypeToAttrType(objType, tfVal)
	require.NoError(t, err)

	obj := val.(types.Object)
	require.Equal(t, "Alice", obj.Attributes()["name"].(types.String).ValueString())
	require.Equal(t, int64(30), obj.Attributes()["age"].(types.Int64).ValueInt64())
}
