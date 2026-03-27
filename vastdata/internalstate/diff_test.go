// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestGetType_ReturnsElemTypeForPointer(t *testing.T) {
	type foo struct{}
	var v *foo
	typ := GetType(v)
	require.Equal(t, "foo", typ.Name())
}

func TestGetType_ReturnsTypeForNonPointer(t *testing.T) {
	type bar struct{}
	var v bar
	typ := GetType(v)
	require.Equal(t, "bar", typ.Name())
}

func TestToInt_HandlesIntTypes(t *testing.T) {
	i, err := ToInt(int64(42))
	require.NoError(t, err)
	require.Equal(t, int64(42), i)

	i, err = ToInt(int(7))
	require.NoError(t, err)
	require.Equal(t, int64(7), i)

	i, err = ToInt(float64(3.0))
	require.NoError(t, err)
	require.Equal(t, int64(3), i)
}

func TestToInt_ReturnsErrorForUnsupportedType(t *testing.T) {
	_, err := ToInt("foo")
	require.Error(t, err)
}

func TestIsNil_HandlesNilAndNonNil(t *testing.T) {
	var p *int
	require.True(t, IsNil(nil))
	require.True(t, IsNil(p))
	require.False(t, IsNil(5))
	require.False(t, IsNil(&struct{}{}))
}

func TestIsNil_HandlesNilSliceAndMap(t *testing.T) {
	var s []int
	var m map[string]int
	require.True(t, IsNil(s))
	require.True(t, IsNil(m))
}

func TestMust_PanicsOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic")
		}
	}()
	Must(0, assertError())
}

func assertError() error {
	return &customError{}
}

type customError struct{}

func (e *customError) Error() string { return "err" }

func TestMust_ReturnsValueOnNoError(t *testing.T) {
	val := Must(42, nil)
	require.Equal(t, 42, val)
}

func TestToFloat_HandlesVariousTypes(t *testing.T) {
	f, err := ToFloat(float64(1.5))
	require.NoError(t, err)
	require.Equal(t, 1.5, f)

	f, err = ToFloat(float32(2.5))
	require.NoError(t, err)
	require.Equal(t, 2.5, f)

	f, err = ToFloat(int64(3))
	require.NoError(t, err)
	require.Equal(t, 3.0, f)

	f, err = ToFloat("4.2")
	require.NoError(t, err)
	require.Equal(t, 4.2, f)
}

func TestToFloat_UnsupportedType(t *testing.T) {
	_, err := ToFloat([]int{1})
	require.Error(t, err)
}

func TestSnakeCaseName_ConvertsPascalToSnake(t *testing.T) {
	type FooBarBaz struct{}
	s := SnakeCaseName(FooBarBaz{})
	require.Equal(t, "foo_bar_baz", s)
}

func TestContains_FindsElement(t *testing.T) {
	require.True(t, contains([]int{1, 2, 3}, 2))
	require.False(t, contains([]string{"a", "b"}, "c"))
}

func TestParsePath_ParsesComplexPath(t *testing.T) {
	parts := parsePath("foo.bar[2].baz[10]")
	require.Equal(t, []any{"foo", "bar", 2, "baz", 10}, parts)
}

func TestGetElements_ListAndSet(t *testing.T) {
	l := types.ListValueMust(types.StringType, []attr.Value{types.StringValue("a")})
	require.Len(t, getElements(l), 1)
	s := types.SetValueMust(types.StringType, []attr.Value{types.StringValue("b")})
	require.Len(t, getElements(s), 1)
}

func TestGetElements_PanicsOnUnsupportedType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic")
		}
	}()
	getElements(types.StringValue("foo"))
}

func TestPartsToStrings_ConvertsParts(t *testing.T) {
	parts := []any{"foo", 2, "bar"}
	strs := partsToStrings(parts)
	require.Equal(t, []string{"foo", "[2]", "bar"}, strs)
}

func TestRemoveNilValues_RemovesNilFromMapAndSlice(t *testing.T) {
	m := map[string]any{"a": 1, "b": nil, "c": map[string]any{"d": nil, "e": 2}}
	cleaned := RemoveNilValues(m).(map[string]any)
	require.Equal(t, map[string]any{"a": 1, "c": map[string]any{"e": 2}}, cleaned)

	s := []any{1, nil, 2}
	cleanedSlice := RemoveNilValues(s).([]any)
	require.Equal(t, []any{1, 2}, cleanedSlice)
}

func TestRemoveNilValues_ReturnsDefaultForOtherTypes(t *testing.T) {
	require.Equal(t, 5, RemoveNilValues(5))
	require.Equal(t, "foo", RemoveNilValues("foo"))
}

func TestEqual_Primitives(t *testing.T) {
	ok, diff := Equal("a", "a")
	require.True(t, ok)
	require.Empty(t, diff)

	ok, diff = Equal(123, 456)
	require.False(t, ok)
	require.Contains(t, diff, "mismatch")
}

func TestEqual_StringSlices_OrderInsensitive(t *testing.T) {
	a := []string{"b", "a", "c"}
	b := []string{"c", "b", "a"}
	ok, diff := Equal(a, b)
	require.True(t, ok, diff)
}

func TestEqual_NestedSlices_OrderInsensitive(t *testing.T) {
	a := [][]string{{"b", "a"}, {"z"}}
	b := [][]string{{"a", "b"}, {"z"}}
	ok, diff := Equal(a, b)
	require.True(t, ok, diff)
}

func TestDiffMap_Basic(t *testing.T) {
	map1 := map[string]any{
		"a": "val1",
		"b": []string{"a", "b"},
		"c": 1,
	}
	map2 := map[string]any{
		"a": "val1",
		"b": []string{"b", "a"},
	}

	diff := DiffMap(map1, map2)
	require.Contains(t, diff, "c")
	require.Equal(t, 1, diff["c"])
}

func TestDiffMap_Empty(t *testing.T) {
	map1 := map[string]any{
		"a": "x",
		"b": []string{"1", "2"},
	}
	map2 := map[string]any{
		"a": "x",
		"b": []string{"2", "1"},
	}

	diff := DiffMap(map1, map2)
	require.Empty(t, diff)
}

// ---------------------------------------------------------------------------
// DiffFields tests
// ---------------------------------------------------------------------------

// buildSimpleOptionalSchema returns a schema with the given optional string fields.
func buildSimpleOptionalSchema(fields ...string) rschema.Schema {
	attrs := make(map[string]rschema.Attribute, len(fields))
	for _, f := range fields {
		attrs[f] = rschema.StringAttribute{Optional: true}
	}
	return rschema.Schema{Attributes: attrs}
}

func TestDiffFields_NonNullToNull_IncludedAsNil(t *testing.T) {
	// Regression test for: setting an optional field to null should produce a
	// nil entry in the diff so the PATCH body sends `"qos_policy": null` and
	// the API clears the field.
	schema := buildSimpleOptionalSchema("name", "qos_policy")

	currentState := NewTFStateMust(map[string]attr.Value{
		"name":       types.StringValue("my_view"),
		"qos_policy": types.StringValue("my_qos_policy"),
	}, schema, nil)

	planState := NewTFStateMust(map[string]attr.Value{
		"name":       types.StringValue("my_view"),
		"qos_policy": types.StringNull(), // user sets qos_policy = null
	}, schema, nil)

	diff := planState.DiffFields(currentState, FilterOr, nil, SearchOptional)

	require.Contains(t, diff, "qos_policy", "cleared field must appear in diff")
	require.Nil(t, diff["qos_policy"], "cleared field value must be nil so it serialises as JSON null")
	require.NotContains(t, diff, "name", "unchanged field must not appear in diff")
}

func TestDiffFields_NullToNull_NotIncluded(t *testing.T) {
	// A field that is null in both plan and current state must not appear in the diff.
	schema := buildSimpleOptionalSchema("name", "qos_policy")

	currentState := NewTFStateMust(map[string]attr.Value{
		"name":       types.StringValue("my_view"),
		"qos_policy": types.StringNull(),
	}, schema, nil)

	planState := NewTFStateMust(map[string]attr.Value{
		"name":       types.StringValue("my_view"),
		"qos_policy": types.StringNull(),
	}, schema, nil)

	diff := planState.DiffFields(currentState, FilterOr, nil, SearchOptional)

	require.NotContains(t, diff, "qos_policy", "field null in both states must not appear in diff")
	require.Empty(t, diff)
}

func TestDiffFields_NullToValue_Included(t *testing.T) {
	// A field transitioning from null to a concrete value must appear in the diff.
	schema := buildSimpleOptionalSchema("name", "qos_policy")

	currentState := NewTFStateMust(map[string]attr.Value{
		"name":       types.StringValue("my_view"),
		"qos_policy": types.StringNull(),
	}, schema, nil)

	planState := NewTFStateMust(map[string]attr.Value{
		"name":       types.StringValue("my_view"),
		"qos_policy": types.StringValue("new_policy"),
	}, schema, nil)

	diff := planState.DiffFields(currentState, FilterOr, nil, SearchOptional)

	require.Contains(t, diff, "qos_policy")
	require.Equal(t, "new_policy", diff["qos_policy"])
}

func TestDiffFields_ValueToValue_Changed_Included(t *testing.T) {
	// A field changing from one non-null value to another must appear in the diff.
	schema := buildSimpleOptionalSchema("name", "qos_policy")

	currentState := NewTFStateMust(map[string]attr.Value{
		"name":       types.StringValue("my_view"),
		"qos_policy": types.StringValue("old_policy"),
	}, schema, nil)

	planState := NewTFStateMust(map[string]attr.Value{
		"name":       types.StringValue("my_view"),
		"qos_policy": types.StringValue("new_policy"),
	}, schema, nil)

	diff := planState.DiffFields(currentState, FilterOr, nil, SearchOptional)

	require.Contains(t, diff, "qos_policy")
	require.Equal(t, "new_policy", diff["qos_policy"])
	require.NotContains(t, diff, "name")
}

func TestDiffFields_ComputedField_ClearedNotIncluded(t *testing.T) {
	// A computed-only field set to null in the plan must NOT appear in the diff
	// because computed fields are managed by the provider/API, not the user.
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"name":   rschema.StringAttribute{Optional: true},
		"id":     rschema.Int64Attribute{Computed: true},
		"qos_id": rschema.Int64Attribute{Computed: true},
	}}

	currentState := NewTFStateMust(map[string]attr.Value{
		"name":   types.StringValue("my_view"),
		"id":     types.Int64Value(1),
		"qos_id": types.Int64Value(42),
	}, schema, nil)

	planState := NewTFStateMust(map[string]attr.Value{
		"name":   types.StringValue("my_view"),
		"id":     types.Int64Null(), // computed - null in plan is normal
		"qos_id": types.Int64Null(), // computed - null in plan is normal
	}, schema, nil)

	diff := planState.DiffFields(currentState, FilterOr, nil, SearchOptional)

	require.NotContains(t, diff, "id", "computed field must not be included in update diff")
	require.NotContains(t, diff, "qos_id", "computed field must not be included in update diff")
}
