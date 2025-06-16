package utils

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	metadata "github.com/vast-data/terraform-provider-vastdata/metadata"
)

func GetListDimentionAndType(t reflect.Type) (int, reflect.Type) {
	d := 0
	for t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		t = t.Elem()
		d++
	}
	return d, t
}

func IsPrimitive(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.String:
		return true
	}
	return false
}

func MakeStructMap(t reflect.Type, m *map[string]interface{}) {
	for _, f := range reflect.VisibleFields(t) {
		j := f.Tag.Get("json")
		names := strings.Split(j, ",")
		if len(names) == 0 {
			// if we dont have a json tag we give up
			continue
		}
		switch f.Type.Kind() {
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.String:
			(*m)[names[0]] = f.Type.String()
		case reflect.Struct, reflect.Pointer:
			t := f.Type
			if f.Type.Kind() == reflect.Pointer {
				t = f.Type.Elem()
			}
			if t.Kind() == reflect.Struct {
				n := map[string]interface{}{}
				(*m)[names[0]] = n
				MakeStructMap(t, &n)

			}
		case reflect.Array, reflect.Slice:
			dim, tp := GetListDimentionAndType(f.Type)
			if (dim >= 1 && dim <= 2) && IsPrimitive(tp) {
				(*m)[names[0]] = f.Type.Elem().String()

			}
			if dim > 1 {
				//we can opnly deal with a single dimention list from now on
				continue
			}

			if tp.Kind() == reflect.Pointer {
				//we deal only with the direct type not a pointer
				tp = tp.Elem()
			}
			if tp.Kind() != reflect.Struct {
				//we only handle structs list
				continue
			}
			n := map[string]interface{}{}
			(*m)[names[0]] = n
			MakeStructMap(tp, &n)

		}

	}

}

func isValidMap(i interface{}) bool {
	return reflect.TypeOf(i) == reflect.TypeOf(map[string]interface{}{})
}

func toValidMap(i interface{}) *map[string]interface{} {
	if isValidMap(i) {
		t := i.(map[string]interface{})
		return &t
	}
	return nil
}

func StructMapDiff(struct_map map[string]interface{}, data_map map[string]interface{}, prefix string, diffs *[]string) {

	for key, value := range data_map {
		_, exists := struct_map[key]
		if !exists {

			*diffs = append(*diffs, prefix+"."+key)
		} else {
			if value == nil || IsPrimitive(reflect.TypeOf(value)) {
				continue
			}
			switch reflect.TypeOf(value) {
			case reflect.TypeOf(map[string]interface{}{}):
				StructMapDiff(struct_map[key].(map[string]interface{}), data_map[key].(map[string]interface{}), prefix+"."+key, diffs)
			case reflect.TypeOf([]interface{}{}):
				sm := toValidMap(struct_map[key])
				if sm == nil {
					/*
					   this means it cannot be converted we will have to skip
					   but we should add it as we can not verify it matches
					*/
					*diffs = append(*diffs, prefix+"."+key)
					continue

				}
				_v := value.([]interface{})
				for i, _ := range _v {
					vs := toValidMap(_v[i])
					if vs == nil {
						/*
						   this means it cannot be converted we will have to skip
						   but we should add it as we can not verify it matches
						*/
						*diffs = append(*diffs, prefix+"."+key+"."+fmt.Sprintf("%d", i))
						continue
					}
					StructMapDiff(*sm, *vs, prefix+"."+key+"."+fmt.Sprintf("%d", i), diffs)
				}

			}

		}
	}

}

func VastVersionsWarn(ctx context.Context) int {
	//Version checking//
	version_compare := metadata.ClusterVersionCompare()
	tflog.Debug(ctx, fmt.Sprintf("Version Compare %v", version_compare))
	if metadata.IsLowerThanMinVersion() {
		tflog.Warn(ctx, fmt.Sprintf("Cluster Version is lower than the minimum provider version (%s<%s), strict validation (strict_version_validation=True) is not supported", metadata.ClusterVersionString(), metadata.GetMinVersion()))
	} else if version_compare == metadata.CLUSTER_VERSION_GRATER {
		tflog.Warn(ctx, fmt.Sprintf("Cluster Version is greater than provider (%s>%s)  build version, please consider upgrading", metadata.ClusterVersionString(), metadata.BuildVersionString()))
	} else if version_compare == metadata.CLUSTER_VERSION_LOWER {
		tflog.Warn(ctx, fmt.Sprintf("Cluster Version is lower than the provider (%s<%s) this might result in resouce creation/update faliure", metadata.ClusterVersionString(), metadata.BuildVersionString()))
	} else {
		tflog.Debug(ctx, fmt.Sprintf("Cluster Version is matching build version"))
	}
	return version_compare
}

func VersionMatch(t reflect.Type, data_map map[string]interface{}) error {
	//Return an error if
	struct_map := map[string]interface{}{}
	MakeStructMap(t, &struct_map)
	//	map[string]interface{}, prefix string, diffs *[]string
	diffs := []string{}
	StructMapDiff(struct_map, data_map, "", &diffs)
	if len(diffs) > 0 {
		sort.Strings(diffs)
		msg := fmt.Sprintf("The following fields found which do not match the Cluster version %s \n", metadata.ClusterVersionString())
		for _, fld := range diffs {
			msg += fld + "\n"
		}
		return errors.New(msg)

	}
	return nil
}
