// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestGetType_ReturnsElemTypeForDoublePointer(t *testing.T) {
	type foo struct{}
	var v **foo
	typ := GetType(v)
	require.Equal(t, "foo", typ.Name())
}

func TestToInt_HandlesNegativeAndZero(t *testing.T) {
	i, err := ToInt(int64(-5))
	require.NoError(t, err)
	require.Equal(t, int64(-5), i)

	i, err = ToInt(0)
	require.NoError(t, err)
	require.Equal(t, int64(0), i)
}

func TestIsNil_HandlesNonNilSliceAndMap(t *testing.T) {
	s := []int{1}
	m := map[string]int{"a": 1}
	require.False(t, IsNil(s))
	require.False(t, IsNil(m))
}

func TestMust_PanicsWithErrorMessage(t *testing.T) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
		require.Contains(t, r.(string), "must:")
	}()
	Must(0, assertError())
}

func TestToFloat_ParsesStringError(t *testing.T) {
	_, err := ToFloat("not-a-number")
	require.Error(t, err)
}

func TestSnakeCaseName_HandlesSingleWord(t *testing.T) {
	type Foo struct{}
	s := SnakeCaseName(Foo{})
	require.Equal(t, "foo", s)
}

func TestContains_EmptyList(t *testing.T) {
	require.False(t, contains([]int{}, 1))
}

func TestParsePath_EmptyString(t *testing.T) {
	parts := parsePath("")
	require.Empty(t, parts)
}

func TestGetElements_PanicsOnIntType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic")
		}
	}()
	getElements(types.Int64Value(1))
}

func TestPartsToStrings_EmptyParts(t *testing.T) {
	strs := partsToStrings([]any{})
	require.Empty(t, strs)
}

func TestRemoveNilValues_NestedNilSlice(t *testing.T) {
	s := []any{nil, []any{nil, 1}, 2}
	cleaned := RemoveNilValues(s).([]any)
	require.Equal(t, []any{[]any{1}, 2}, cleaned)
}

func TestRemoveNilValues_NestedNilMap(t *testing.T) {
	m := map[string]any{"a": map[string]any{"b": nil, "c": 3}, "d": nil}
	cleaned := RemoveNilValues(m).(map[string]any)
	require.Equal(t, map[string]any{"a": map[string]any{"c": 3}}, cleaned)
}

func TestRemoveNilValues_EmptyMapAndSlice(t *testing.T) {
	require.Equal(t, map[string]any{}, RemoveNilValues(map[string]any{}))
	require.Equal(t, []any{}, RemoveNilValues([]any{}))
}

func TestGetType(t *testing.T) {
	var v *string
	require.Equal(t, reflect.TypeOf(""), GetType(v))
	require.Equal(t, reflect.TypeOf(0), GetType(42))
}

func TestToInt(t *testing.T) {
	n, err := ToInt(42)
	require.NoError(t, err)
	require.Equal(t, int64(42), n)

	n, err = ToInt(float64(3.14))
	require.NoError(t, err)
	require.Equal(t, int64(3), n)

	_, err = ToInt("not-an-int")
	require.Error(t, err)
}

func TestToFloat(t *testing.T) {
	f, err := ToFloat(1)
	require.NoError(t, err)
	require.Equal(t, float64(1), f)

	f, err = ToFloat("3.14")
	require.NoError(t, err)
	require.Equal(t, 3.14, f)

	_, err = ToFloat(struct{}{})
	require.Error(t, err)
}

func TestIsNil(t *testing.T) {
	var m map[string]any
	var s []string
	var p *int
	require.True(t, IsNil(nil))
	require.True(t, IsNil(m))
	require.True(t, IsNil(s))
	require.True(t, IsNil(p))
	require.False(t, IsNil(0))
}

func TestSnakeCaseName(t *testing.T) {
	type HelloWorld struct{}
	name := SnakeCaseName(HelloWorld{})
	require.Equal(t, "hello_world", name)
}

func TestContains(t *testing.T) {
	require.True(t, contains([]string{"a", "b"}, "a"))
	require.False(t, contains([]int{1, 2}, 3))
}

func TestParsePath(t *testing.T) {
	path := "spec[0].items[2].name"
	parts := parsePath(path)
	expected := []any{"spec", 0, "items", 2, "name"}
	require.Equal(t, expected, parts)
}

func TestPartsToStrings(t *testing.T) {
	input := []any{"spec", 0, "data"}
	out := partsToStrings(input)
	require.Equal(t, []string{"spec", "[0]", "data"}, out)
}

func TestRemoveNilValues(t *testing.T) {
	input := map[string]any{
		"a": nil,
		"b": 1,
		"c": map[string]any{
			"d": nil,
			"e": "value",
		},
		"f": []any{nil, "x"},
	}
	cleaned := RemoveNilValues(input).(map[string]any)

	require.NotContains(t, cleaned, "a")
	require.Contains(t, cleaned, "b")
	require.Equal(t, "value", cleaned["c"].(map[string]any)["e"])
	require.Equal(t, []any{"x"}, cleaned["f"])
}

func TestRemoveNilValues_NestedObjectStructure(t *testing.T) {
	input := map[string]any{
		"config": map[string]any{
			"bucket_name":        "vastdb-metrics",
			"bucket_owner":       "metrics-user",
			"enabled":            true,
			"max_capacity_mb":    int64(2048),
			"retention_time_sec": int64(172800),
		},
		"user_defined_columns": []any{
			map[string]any{
				"field": map[string]any{
					"column_type": "integer",
					"key_type":    nil,
					"value_type":  nil,
				},
				"name": "ENV_ACCESS_COUNT",
			},
			map[string]any{
				"field": map[string]any{
					"column_type": "string",
					"key_type":    nil,
					"value_type":  nil,
				},
				"name": "ENV_USER_ID",
			},
		},
	}

	cleaned := RemoveNilValues(input).(map[string]any)

	udc := cleaned["user_defined_columns"].([]any)
	for _, e := range udc {
		obj := e.(map[string]any)
		field := obj["field"].(map[string]any)
		if _, ok := field["key_type"]; ok {
			t.Fatalf("key_type should have been removed, found: %v", field["key_type"])
		}
		if _, ok := field["value_type"]; ok {
			t.Fatalf("value_type should have been removed, found: %v", field["value_type"])
		}
		if _, ok := field["column_type"]; !ok {
			t.Fatalf("column_type should remain in field")
		}
	}
}

func TestSet_InstantiateFromNil(t *testing.T) {
	var initial []int = nil
	s := NewSet(initial)
	require.Equal(t, 0, s.Len())
	require.False(t, s.Contains(1))
}

func TestSet_IntOperations(t *testing.T) {
	s := NewSet([]int{1, 2, 3})
	require.Equal(t, 3, s.Len())

	require.False(t, s.Add(2)) // already exists
	require.True(t, s.Add(4))  // new
	require.True(t, s.Contains(4))
	require.True(t, s.Remove(1))
	require.False(t, s.Contains(1))
	require.False(t, s.Remove(10)) // doesn't exist

	out := s.ToSlice()
	sort.Ints(out)
	require.Equal(t, []int{2, 3, 4}, out)
}

func TestSet_StringOperations(t *testing.T) {
	s := NewSet([]string{"a", "b"})
	require.True(t, s.Contains("a"))
	require.False(t, s.Contains("z"))

	require.True(t, s.Add("z"))
	require.True(t, s.Remove("b"))
	require.Equal(t, 2, s.Len())

	require.ElementsMatch(t, []string{"a", "z"}, s.ToSlice())
}

func TestSet_Clear(t *testing.T) {
	s := NewSet([]string{"x", "y"})
	require.Equal(t, 2, s.Len())

	s.Clear()
	require.Equal(t, 0, s.Len())
	require.False(t, s.Contains("x"))
	require.Empty(t, s.ToSlice())
}

func TestSet_ToSlice_IsCopy(t *testing.T) {
	s := NewSet([]int{1, 2})
	sl := s.ToSlice()
	sl = append(sl, 99) // should not affect internal set
	require.False(t, s.Contains(99))
	require.Equal(t, 2, s.Len())
}

func TestSet_NewSetFromAny_Int64(t *testing.T) {
	input := []any{1, int64(2), float64(3)}
	s, err := NewSetFromAny[int64](input)
	require.NoError(t, err)
	require.ElementsMatch(t, []int64{1, 2, 3}, s.ToSlice())
}

func TestSet_NewSetFromAny_Float64(t *testing.T) {
	input := []any{1, float64(2.5), int64(3)}
	s, err := NewSetFromAny[float64](input)
	require.NoError(t, err)
	require.ElementsMatch(t, []float64{1.0, 2.5, 3.0}, s.ToSlice())
}

func TestSet_NewSetFromAny_String(t *testing.T) {
	input := []any{"foo", "bar"}
	s, err := NewSetFromAny[string](input)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"foo", "bar"}, s.ToSlice())
}

func TestSet_NewSetFromAny_Nil(t *testing.T) {
	s, err := NewSetFromAny[int64](nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	require.Empty(t, s.ToSlice())
}

func TestSet_NewSetFromAny_AlreadyTypedSlice(t *testing.T) {
	s, err := NewSetFromAny[int64]([]int64{10, 20})
	require.NoError(t, err)
	require.ElementsMatch(t, []int64{10, 20}, s.ToSlice())
}

func TestSet_NewSetFromAny_InvalidType(t *testing.T) {
	_, err := NewSetFromAny[int64]("not a slice")
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported input type")
}

func TestSet_NewSetFromAny_InvalidElementType(t *testing.T) {
	input := []any{1, "not-an-int"}
	_, err := NewSetFromAny[int64](input)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot convert")
}

func Test_isPrimitiveType(t *testing.T) {
	t.Run("primitive types", func(t *testing.T) {
		require.True(t, isPrimitiveType(types.StringType))
		require.True(t, isPrimitiveType(types.BoolType))
		require.True(t, isPrimitiveType(types.Int64Type))
		require.True(t, isPrimitiveType(types.Float64Type))
		require.True(t, isPrimitiveType(types.NumberType))
	})

	t.Run("non primitive types", func(t *testing.T) {
		require.False(t, isPrimitiveType(types.ListType{ElemType: types.StringType}))
		require.False(t, isPrimitiveType(types.SetType{ElemType: types.StringType}))
		require.False(t, isPrimitiveType(types.MapType{ElemType: types.StringType}))
		require.False(t, isPrimitiveType(types.ObjectType{AttrTypes: map[string]attr.Type{"a": types.StringType}}))
	})
}
