package vast_versions

import (
	"reflect"
	"regexp"
	"strings"

	version "github.com/hashicorp/go-version"
	"github.com/vast-data/terraform-provider-vastdata/metadata"
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

func GetTFformatName(s string) string {
	re := regexp.MustCompile("([A-Z])")
	t := re.ReplaceAllString(s, "_${1}")
	return strings.ToLower(strings.TrimLeft(t, "_"))
}

func TypeAttributeExists(r reflect.Type, a string) bool {
	for _, f := range reflect.VisibleFields(r) {
		if GetTFformatName(f.Name) == a {
			return true
		}
	}
	return false
}

func VersionsSupportingAttributes(struct_name string, attribute_name string) []string {
	//return a list of all versions (as long at the are larger or equals to the min version ) that hold this attribute.
	l := []string{}
	min_version := metadata.GetMinVersion()
	for k, _ := range vast_versions {
		v, e := version.NewVersion(k)
		if e != nil {
			continue
		}
		if v.Compare(&min_version) < 0 {
			//if found version is smaller than minimal version we skip validation
			continue
		}
		t, u := GetVersionedType(k, struct_name)
		if !u {
			continue
		}
		if TypeAttributeExists(t, attribute_name) {
			l = append(l, k)
		}
	}
	return l
}
