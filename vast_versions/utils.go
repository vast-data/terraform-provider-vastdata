package vast_versions

import (
	"reflect"
)

func GetVersionedType(version, struct_name string) (reflect.Type, bool) {
	/*
	   This function will return the type of a struct with the name struct_name whcih is matching a version.
	   If it does not exists the function will return TypeOf(nil) ,false
	*/
	structs_map, versions_map_exists := vast_versions[version]
	if !versions_map_exists {
		return reflect.TypeOf(nil), false
	}
	struct_type, struct_type_exist := structs_map[struct_name]
	if !struct_type_exist {
		return reflect.TypeOf(nil), false

	}
	return struct_type, true

}
