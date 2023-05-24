package utils

import (
	"fmt"
	"reflect"
)

//This file should hold tools that can help us in unittesting

func MapToTFSchema(m map[string]interface{}, n *map[string]string, prefix string) {
	_prefix := ""
	if prefix == "" {
		_prefix = ""
	} else {
		_prefix = prefix + "."
	}
	for k, v := range m {
		if IsPrimitive(reflect.TypeOf(v)) {
			(*n)[prefix+k] = fmt.Sprintf("%v", v)
		} else if reflect.TypeOf(v).Kind() == reflect.Array || reflect.TypeOf(v).Kind() == reflect.Slice {
			(*n)[prefix+k+".#"] = fmt.Sprintf("%v", len(v.([]interface{})))
			for i, j := range v.([]interface{}) {
				if IsPrimitive(reflect.TypeOf(j)) {
					(*n)[_prefix+k+fmt.Sprintf(".%v", i)] = fmt.Sprintf("%v", j)
				}

			}
		}

	}
}
