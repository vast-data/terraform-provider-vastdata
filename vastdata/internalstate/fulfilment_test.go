// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

func TestBuildAttrValueFromAny_ArbitraryToString(t *testing.T) {
	v := struct{ A string }{A: "val"}
	val, _, err := BuildAttrValueFromAny(types.StringType, v)
	require.NoError(t, err)
	require.Equal(t, "{val}", val.(types.String).ValueString())
}

func TestBuildAttrValueFromAny_NullValues(t *testing.T) {
	val, _, err := BuildAttrValueFromAny(types.StringType, nil)
	require.NoError(t, err)
	require.True(t, val.IsNull())

	val, _, err = BuildAttrValueFromAny(types.Int64Type, nil)
	require.NoError(t, err)
	require.True(t, val.IsNull())

	val, _, err = BuildAttrValueFromAny(types.ListType{ElemType: types.StringType}, nil)
	require.NoError(t, err)
	require.True(t, val.IsNull())
}

func TestBuildAttrValueFromAny_ListAndSet(t *testing.T) {
	listVal, _, err := BuildAttrValueFromAny(types.ListType{ElemType: types.StringType}, []any{"a", "b"})
	require.NoError(t, err)
	require.Equal(t, types.ListValueMust(types.StringType, []attr.Value{types.StringValue("a"), types.StringValue("b")}), listVal)

	setVal, _, err := BuildAttrValueFromAny(types.SetType{ElemType: types.Int64Type}, []any{int64(1), int64(2)})
	require.NoError(t, err)
	require.Equal(t, types.SetValueMust(types.Int64Type, []attr.Value{types.Int64Value(1), types.Int64Value(2)}), setVal)
}

func TestBuildAttrValueFromAny_ObjectMissingField(t *testing.T) {
	objType := types.ObjectType{AttrTypes: map[string]attr.Type{"foo": types.StringType, "bar": types.Int64Type}}
	val, _, err := BuildAttrValueFromAny(objType, map[string]any{"foo": "baz"})
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
			got, _, err := BuildAttrValueFromAny(tc.typ, tc.val)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestBuildAttrValueFromAny_Null(t *testing.T) {
	val, _, err := BuildAttrValueFromAny(types.StringType, nil)
	require.NoError(t, err)
	require.True(t, val.IsNull())
}

func TestBuildAttrValueFromAny_List(t *testing.T) {
	listType := types.ListType{ElemType: types.StringType}
	input := []any{"a", "b"}
	val, _, err := BuildAttrValueFromAny(listType, input)
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

func TestConvertAttrValueToRaw_SetOfObjects(t *testing.T) {
	// Create an object type for user quota
	objType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"entity_identifier": types.StringType,
			"guid":              types.StringType,
			"id":                types.Int64Type,
			"identifier":        types.StringType,
			"identifier_type":   types.StringType,
			"is_group":          types.BoolType,
			"name":              types.StringType,
			"path":              types.StringType,
		},
	}

	// Create a set type containing objects
	setType := types.SetType{ElemType: objType}

	// Create an object with some values
	obj1 := types.ObjectValueMust(objType.AttrTypes, map[string]attr.Value{
		"entity_identifier": types.StringValue("user1"),
		"guid":              types.StringValue("guid1"),
		"id":                types.Int64Value(1),
		"identifier":        types.StringValue("user1"),
		"identifier_type":   types.StringValue("name"),
		"is_group":          types.BoolValue(false),
		"name":              types.StringValue("User One"),
		"path":              types.StringValue("/path1"),
	})

	// Create a set with the object
	set := types.SetValueMust(setType.ElemType, []attr.Value{obj1})

	// Convert to raw
	raw := ConvertAttrValueToRaw(set, setType)

	// Should be a slice of maps
	require.IsType(t, []any{}, raw)
	result := raw.([]any)
	require.Len(t, result, 1)

	// The first element should be a map
	require.IsType(t, map[string]any{}, result[0])
	objMap := result[0].(map[string]any)

	// Check that the object fields are preserved
	require.Equal(t, "user1", objMap["entity_identifier"])
	require.Equal(t, "guid1", objMap["guid"])
	require.Equal(t, int64(1), objMap["id"])
	require.Equal(t, "user1", objMap["identifier"])
	require.Equal(t, "name", objMap["identifier_type"])
	require.Equal(t, false, objMap["is_group"])
	require.Equal(t, "User One", objMap["name"])
	require.Equal(t, "/path1", objMap["path"])
}

func TestConvertAttrValueToRaw_SetOfEmptyObjects(t *testing.T) {
	// Create an object type for user quota
	objType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"entity_identifier": types.StringType,
			"guid":              types.StringType,
			"id":                types.Int64Type,
			"identifier":        types.StringType,
			"identifier_type":   types.StringType,
			"is_group":          types.BoolType,
			"name":              types.StringType,
			"path":              types.StringType,
		},
	}

	// Create a set type containing objects
	setType := types.SetType{ElemType: objType}

	// Create an empty object (all fields null)
	obj1 := types.ObjectValueMust(objType.AttrTypes, map[string]attr.Value{
		"entity_identifier": types.StringNull(),
		"guid":              types.StringNull(),
		"id":                types.Int64Null(),
		"identifier":        types.StringNull(),
		"identifier_type":   types.StringNull(),
		"is_group":          types.BoolNull(),
		"name":              types.StringNull(),
		"path":              types.StringNull(),
	})

	// Create a set with the empty object
	set := types.SetValueMust(setType.ElemType, []attr.Value{obj1})

	// Convert to raw
	raw := ConvertAttrValueToRaw(set, setType)

	// Should be a slice of maps
	require.IsType(t, []any{}, raw)
	result := raw.([]any)
	require.Len(t, result, 1)

	// The first element should be a map
	require.IsType(t, map[string]any{}, result[0])
	objMap := result[0].(map[string]any)

	// Check that the object fields are preserved (even if null)
	require.Contains(t, objMap, "entity_identifier")
	require.Contains(t, objMap, "guid")
	require.Contains(t, objMap, "id")
	require.Contains(t, objMap, "identifier")
	require.Contains(t, objMap, "identifier_type")
	require.Contains(t, objMap, "is_group")
	require.Contains(t, objMap, "name")
	require.Contains(t, objMap, "path")

	// All values should be nil for null fields
	require.Nil(t, objMap["entity_identifier"])
	require.Nil(t, objMap["guid"])
	require.Nil(t, objMap["id"])
	require.Nil(t, objMap["identifier"])
	require.Nil(t, objMap["identifier_type"])
	require.Nil(t, objMap["is_group"])
	require.Nil(t, objMap["name"])
	require.Nil(t, objMap["path"])
}

func TestConvertAttrValueToRaw_SetOfObjectsWithNullFields(t *testing.T) {
	// Create an object type for user quota
	objType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"entity_identifier": types.StringType,
			"guid":              types.StringType,
			"id":                types.Int64Type,
			"identifier":        types.StringType,
			"identifier_type":   types.StringType,
			"is_group":          types.BoolType,
			"name":              types.StringType,
			"path":              types.StringType,
		},
	}

	// Create a set type containing objects
	setType := types.SetType{ElemType: objType}

	// Create an object with all null fields
	obj1 := types.ObjectValueMust(objType.AttrTypes, map[string]attr.Value{
		"entity_identifier": types.StringNull(),
		"guid":              types.StringNull(),
		"id":                types.Int64Null(),
		"identifier":        types.StringNull(),
		"identifier_type":   types.StringNull(),
		"is_group":          types.BoolNull(),
		"name":              types.StringNull(),
		"path":              types.StringNull(),
	})

	// Create a set with the object
	set := types.SetValueMust(setType.ElemType, []attr.Value{obj1})

	// Convert to raw
	raw := ConvertAttrValueToRaw(set, setType)

	// Should be a slice of maps
	require.IsType(t, []any{}, raw)
	result := raw.([]any)
	require.Len(t, result, 1)

	// The first element should be a map with all null values
	require.IsType(t, map[string]any{}, result[0])
	objMap := result[0].(map[string]any)

	// Check that the object fields are preserved (even if null)
	require.Contains(t, objMap, "entity_identifier")
	require.Contains(t, objMap, "guid")
	require.Contains(t, objMap, "id")
	require.Contains(t, objMap, "identifier")
	require.Contains(t, objMap, "identifier_type")
	require.Contains(t, objMap, "is_group")
	require.Contains(t, objMap, "name")
	require.Contains(t, objMap, "path")

	// All values should be nil for null fields
	require.Nil(t, objMap["entity_identifier"])
	require.Nil(t, objMap["guid"])
	require.Nil(t, objMap["id"])
	require.Nil(t, objMap["identifier"])
	require.Nil(t, objMap["identifier_type"])
	require.Nil(t, objMap["is_group"])
	require.Nil(t, objMap["name"])
	require.Nil(t, objMap["path"])
}

func TestConvertAttrValueToRaw_UserQuotasExample(t *testing.T) {
	// This test demonstrates the fix for the original issue where user_quotas
	// was being converted to [{}] instead of preserving the object structure

	// Create an object type for user quota (matching the real schema)
	objType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"entity_identifier": types.StringType,
			"guid":              types.StringType,
			"id":                types.Int64Type,
			"identifier":        types.StringType,
			"identifier_type":   types.StringType,
			"is_group":          types.BoolType,
			"name":              types.StringType,
			"path":              types.StringType,
		},
	}

	// Create a set type containing objects
	setType := types.SetType{ElemType: objType}

	// Create an object with all null fields (like in the original issue)
	obj1 := types.ObjectValueMust(objType.AttrTypes, map[string]attr.Value{
		"entity_identifier": types.StringNull(),
		"guid":              types.StringNull(),
		"id":                types.Int64Null(),
		"identifier":        types.StringNull(),
		"identifier_type":   types.StringNull(),
		"is_group":          types.BoolNull(),
		"name":              types.StringNull(),
		"path":              types.StringNull(),
	})

	// Create a set with the object
	set := types.SetValueMust(setType.ElemType, []attr.Value{obj1})

	// Convert to raw
	raw := ConvertAttrValueToRaw(set, setType)

	// Should be a slice of maps
	require.IsType(t, []any{}, raw)
	result := raw.([]any)
	require.Len(t, result, 1)

	// The first element should be a map with all null values preserved
	require.IsType(t, map[string]any{}, result[0])
	objMap := result[0].(map[string]any)

	// Check that the object fields are preserved (even if null)
	require.Contains(t, objMap, "entity_identifier")
	require.Contains(t, objMap, "guid")
	require.Contains(t, objMap, "id")
	require.Contains(t, objMap, "identifier")
	require.Contains(t, objMap, "identifier_type")
	require.Contains(t, objMap, "is_group")
	require.Contains(t, objMap, "name")
	require.Contains(t, objMap, "path")

	// All values should be nil for null fields
	require.Nil(t, objMap["entity_identifier"])
	require.Nil(t, objMap["guid"])
	require.Nil(t, objMap["id"])
	require.Nil(t, objMap["identifier"])
	require.Nil(t, objMap["identifier_type"])
	require.Nil(t, objMap["is_group"])
	require.Nil(t, objMap["name"])
	require.Nil(t, objMap["path"])
}

func TestBuildAttrValueFromAny_EmptyString(t *testing.T) {
	// Test soft skip behavior in FillFromRecord - unparseable strings return 0 with success=false
	// Empty string for Int64 returns 0 with success=false (soft skip)
	val, success, err := BuildAttrValueFromAny(types.Int64Type, "")
	require.NotNil(t, val)
	require.False(t, success, "empty string should not be considered successful parse")
	require.Error(t, err, "empty string should produce an error")
	require.Equal(t, int64(0), val.(types.Int64).ValueInt64(), "empty string should default to 0")

	// Test with valid string
	val, success, err = BuildAttrValueFromAny(types.Int64Type, "42")
	require.NotNil(t, val)
	require.True(t, success, "valid number string should be successful")
	require.NoError(t, err)
	require.Equal(t, int64(42), val.(types.Int64).ValueInt64())

	// Test with valid int64
	val, success, err = BuildAttrValueFromAny(types.Int64Type, int64(100))
	require.NotNil(t, val)
	require.True(t, success)
	require.NoError(t, err)
	require.Equal(t, int64(100), val.(types.Int64).ValueInt64())

	// Test with invalid string - soft skip behavior (success=false)
	val, success, err = BuildAttrValueFromAny(types.Int64Type, "not-a-number")
	require.NotNil(t, val)
	require.False(t, success, "unparseable string should not be considered successful parse")
	require.Error(t, err, "unparseable string should produce an error")
	require.Equal(t, int64(0), val.(types.Int64).ValueInt64())

	// Test with "UNKNOWN" string - soft skip behavior (the original issue)
	val, success, err = BuildAttrValueFromAny(types.Int64Type, "UNKNOWN")
	require.NotNil(t, val)
	require.False(t, success, "UNKNOWN string should not be considered successful parse")
	require.Error(t, err, "UNKNOWN should produce an error")
	require.Equal(t, int64(0), val.(types.Int64).ValueInt64(), "UNKNOWN should default to 0")

	// Test with "UNKNOWN" for Float64Type as well - soft skip behavior
	val, success, err = BuildAttrValueFromAny(types.Float64Type, "UNKNOWN")
	require.NotNil(t, val)
	require.False(t, success, "UNKNOWN string should not be considered successful parse")
	require.Error(t, err, "UNKNOWN should produce an error")
	require.Equal(t, 0.0, val.(types.Float64).ValueFloat64(), "UNKNOWN should default to 0.0")
}

func TestBuildAttrValueFromAny_OtherTypes(t *testing.T) {
	// Test that other types still work with success indicator

	// String type - always successful
	val, success, err := BuildAttrValueFromAny(types.StringType, "hello")
	require.NoError(t, err)
	require.True(t, success)
	require.Equal(t, "hello", val.(types.String).ValueString())

	// Bool type - successful
	val, success, err = BuildAttrValueFromAny(types.BoolType, true)
	require.NoError(t, err)
	require.True(t, success)
	require.Equal(t, true, val.(types.Bool).ValueBool())

	// Float type - successful
	val, success, err = BuildAttrValueFromAny(types.Float64Type, 3.14)
	require.NoError(t, err)
	require.True(t, success)
	require.Equal(t, 3.14, val.(types.Float64).ValueFloat64())

	// Null value - successful
	val, success, err = BuildAttrValueFromAny(types.Int64Type, nil)
	require.NoError(t, err)
	require.True(t, success)
	require.True(t, val.IsNull())
}

func TestBuildAttrValueFromAny_NestedStructures(t *testing.T) {
	// Test list with empty string in Int64 element
	// The list itself should return true (successful), but the Int64 element inside will have returned false
	listType := types.ListType{ElemType: types.Int64Type}
	val, success, err := BuildAttrValueFromAny(listType, []any{int64(1), "", int64(3)})
	require.NotNil(t, val)
	require.True(t, success, "list itself succeeds even if element had issue")
	require.NoError(t, err, "list should be created successfully")

	list := val.(types.List)
	elems := list.Elements()
	require.Equal(t, 3, len(elems))
	require.Equal(t, int64(1), elems[0].(types.Int64).ValueInt64())
	require.Equal(t, int64(0), elems[1].(types.Int64).ValueInt64(), "empty string should default to 0")
	require.Equal(t, int64(3), elems[2].(types.Int64).ValueInt64())
}

func TestConvertAttrValueToRaw_UserQuotasWithRealValues(t *testing.T) {
	// This test reproduces the actual issue where user_quotas with real values
	// are being converted to empty objects

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

	// Create the user_quota object type
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

	// Convert to raw
	raw := ConvertAttrValueToRaw(set, setType)

	// Should be a slice of maps
	require.IsType(t, []any{}, raw)
	result := raw.([]any)
	require.Len(t, result, 1)

	// The first element should be a map with the actual values
	require.IsType(t, map[string]any{}, result[0])
	objMap := result[0].(map[string]any)

	// Check that the user_quota fields are preserved
	require.Contains(t, objMap, "entity")
	require.Contains(t, objMap, "grace_period")
	require.Contains(t, objMap, "hard_limit")
	require.Contains(t, objMap, "soft_limit")

	// Check the actual values
	require.Equal(t, "02:00:00", objMap["grace_period"])
	require.Equal(t, int64(15000), objMap["hard_limit"])
	require.Equal(t, int64(15000), objMap["soft_limit"])

	// Check the entity values
	entityMap, ok := objMap["entity"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "tfzealous-kingfisher", entityMap["name"])
	require.Equal(t, "user1@example.com", entityMap["email"])
	require.Equal(t, "tfzealous-kingfisher", entityMap["identifier"])
	require.Equal(t, "username", entityMap["identifier_type"])
	require.Equal(t, false, entityMap["is_group"])
}
