// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	vast_client "github.com/vast-data/go-vast-client"
)

type Record = vast_client.Record

func GetType(v any) reflect.Type {
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func ToInt(val any) (int64, error) {
	var idInt int64
	switch v := val.(type) {
	case int64:
		idInt = v
	case float64:
		idInt = int64(v)
	case int:
		idInt = int64(v)
	case string:
		// Parse string as int64 for import scenarios
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse string %q as int64: %w", v, err)
		}
		idInt = parsed
	default:
		return 0, fmt.Errorf("unexpected type %T for int. Value = %v", v, val)
	}
	return idInt, nil
}

func IsNil(val any) bool {
	if val == nil {
		return true
	}
	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map:
		return rv.IsNil()
	}
	return false
}

func Must[T any](val T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("must: %v", err))
	}
	return val
}
func ToFloat(v any) (float64, error) {
	switch n := v.(type) {
	case float64:
		return n, nil
	case float32:
		return float64(n), nil
	case int, int64, int32:
		return float64(reflect.ValueOf(v).Int()), nil
	case string:
		return strconv.ParseFloat(n, 64)
	default:
		return 0, fmt.Errorf("unsupported type %T for float64. Value = %v", v, v)
	}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func SnakeCaseName(v any) string {
	t := GetType(v)
	name := t.Name()

	// Convert PascalCase to snake_case
	snake := matchFirstCap.ReplaceAllString(name, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func contains[T comparable](list []T, key T) bool {
	for _, item := range list {
		if item == key {
			return true
		}
	}
	return false
}

// isPrimitiveType reports whether the given Terraform framework type is a
// primitive type we consider for simple filtering/searching.
// For our purposes, primitives are limited to: string, bool, int64, and numeric (float/number).
// Container and complex types (list, set, map, object, tuple, dynamic) are treated as non-primitive.
func isPrimitiveType(t attr.Type) bool {
	switch t.(type) {
	case basetypes.StringType, basetypes.BoolType, basetypes.Int64Type, basetypes.Float64Type, basetypes.NumberType, basetypes.Int32Type, basetypes.Float32Type:
		return true
	default:
		return false
	}
}

var pathPartRegex = regexp.MustCompile(`([^[.\]]+)|\[(\d+)\]`)

func parsePath(path string) []any {
	var parts []any
	matches := pathPartRegex.FindAllStringSubmatch(path, -1)
	for _, m := range matches {
		if m[1] != "" {
			parts = append(parts, m[1])
		} else if m[2] != "" {
			idx, _ := strconv.Atoi(m[2])
			parts = append(parts, idx)
		}
	}
	return parts
}

func getElements(v attr.Value) []attr.Value {
	switch val := v.(type) {
	case types.List:
		return val.Elements()
	case types.Set:
		return val.Elements()
	default:
		panic(fmt.Sprintf("type %T does not support indexing", v))
	}
}

func partsToStrings(parts []any) []string {
	var out []string
	for _, p := range parts {
		switch v := p.(type) {
		case string:
			out = append(out, v)
		case int:
			out = append(out, fmt.Sprintf("[%d]", v))
		}
	}
	return out
}

// RemoveNilValues removes nil values recursively from maps and slices.
// It is safe to use before sending request bodies to avoid serializing nulls.
// Note: this is a public wrapper to allow usage from provider code.
func RemoveNilValues(v any) any {
	switch val := v.(type) {
	case map[string]any:
		cleaned := make(map[string]any)
		for k, v2 := range val {
			if v2 == nil {
				continue
			}
			cleanedVal := RemoveNilValues(v2)
			if cleanedVal != nil {
				cleaned[k] = cleanedVal
			}
		}
		return cleaned

	case []any:
		var cleaned []any
		for _, item := range val {
			cleanedItem := RemoveNilValues(item)
			if cleanedItem != nil {
				cleaned = append(cleaned, cleanedItem)
			}
		}
		if cleaned == nil {
			return []any{}
		}
		return cleaned

	default:
		return v
	}
}

type Set[T comparable] struct {
	m map[T]struct{}
}

// NewSet creates a new set from a slice. If the slice is nil, the set is empty.
func NewSet[T comparable](items []T) *Set[T] {
	s := &Set[T]{m: make(map[T]struct{})}
	for _, item := range items {
		s.m[item] = struct{}{}
	}
	return s
}

// Add inserts the element into the set.
// Returns true if the element was added (i.e., it wasn't already present).
func (s *Set[T]) Add(item T) bool {
	if _, exists := s.m[item]; exists {
		return false
	}
	s.m[item] = struct{}{}
	return true
}

// Remove deletes the element from the set.
// Returns true if the element existed and was removed.
func (s *Set[T]) Remove(item T) bool {
	if _, exists := s.m[item]; exists {
		delete(s.m, item)
		return true
	}
	return false
}

// Contains checks if the item is present in the set.
func (s *Set[T]) Contains(item T) bool {
	_, exists := s.m[item]
	return exists
}

// ToSlice returns all elements in the set as a slice.
func (s *Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s.m))
	for k := range s.m {
		result = append(result, k)
	}
	return result
}

// Len returns the number of elements in the set.
func (s *Set[T]) Len() int {
	return len(s.m)
}

// Clear removes all elements from the set.
func (s *Set[T]) Clear() {
	s.m = make(map[T]struct{})
}

// NewSetFromAny creates a set from any input value, handling []interface{} and []T.
// Supports conversion to T: int64, float64, string.
func NewSetFromAny[T comparable](input any) (*Set[T], error) {
	switch casted := input.(type) {
	case nil:
		return NewSet[T](nil), nil
	case []T:
		return NewSet[T](casted), nil
	case []interface{}:
		converted := make([]T, 0, len(casted))
		for _, item := range casted {
			val, err := convertTo[T](item)
			if err != nil {
				return nil, err
			}
			converted = append(converted, val)
		}
		return NewSet[T](converted), nil
	default:
		return nil, fmt.Errorf("unsupported input type for set: %T", input)
	}
}

// convertTo converts an interface{} to T, supporting int64, float64, string.
func convertTo[T comparable](v any) (T, error) {
	var zero T
	switch any(zero).(type) {
	case int64:
		switch val := v.(type) {
		case int:
			return any(int64(val)).(T), nil
		case int64:
			return any(val).(T), nil
		case float64:
			return any(int64(val)).(T), nil
		default:
			return zero, fmt.Errorf("cannot convert %T to int64", v)
		}
	case float64:
		switch val := v.(type) {
		case float64:
			return any(val).(T), nil
		case int:
			return any(float64(val)).(T), nil
		case int64:
			return any(float64(val)).(T), nil
		default:
			return zero, fmt.Errorf("cannot convert %T to float64", v)
		}
	case string:
		if str, ok := v.(string); ok {
			return any(str).(T), nil
		}
		return zero, fmt.Errorf("cannot convert %T to string", v)
	default:
		return zero, fmt.Errorf("unsupported target type: %v", reflect.TypeOf(zero))
	}
}
