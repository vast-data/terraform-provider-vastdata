package mappings

import (
	"fmt"
	"reflect"
	"strconv"
)

func GetType(v interface{}) reflect.Type {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func toInt(val any) (int64, error) {
	var idInt int64
	switch v := val.(type) {
	case int64:
		idInt = v
	case float64:
		idInt = int64(v)
	case int:
		idInt = int64(v)
	default:
		return 0, fmt.Errorf("unexpected type for id field: %T", v)
	}
	return idInt, nil
}

func isNil(val any) bool {
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

func must[T any](val T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("must: %v", err))
	}
	return val
}
func toFloat(v any) (float64, error) {
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
		return 0, fmt.Errorf("unsupported type %T for float64", v)
	}
}
