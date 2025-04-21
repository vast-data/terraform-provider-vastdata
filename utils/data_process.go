package utils

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/constraints"
)

type ContextKey string

/*
	func GetTFformatName(s string) string {
		re := regexp.MustCompile("([A-Z])")
		t := re.ReplaceAllString(s, "_${1}")
		return strings.ToLower(strings.TrimLeft(t, "_"))
	}
*/
func ReverseTFName(s string) string {
	result := ""
	re := regexp.MustCompile("[a-z0-9]+_")
	for _, f := range re.FindAllStringSubmatch(s, -1) {
		t := strings.TrimRight(f[0], "_")
		z := strings.Split(t, "")
		z[0] = strings.ToUpper(z[0])
		result += strings.Join(z, "")
	}
	tail := strings.Split(re.ReplaceAllString(s, ""), "")
	tail[0] = strings.ToUpper(tail[0])
	return result + strings.Join(tail, "")
}

func FlattenListOfPrimitives[T constraints.Integer | constraints.Float | string | bool | constraints.Signed | constraints.Unsigned](elements *[]T) []interface{} {
	if elements != nil {
		flattenList := make([]interface{}, len(*elements), len(*elements))
		for i, v := range *elements {
			flattenList[i] = v
		}
		return flattenList
	}
	return make([]interface{}, 0)
}

func FlattenListOfStrings(elements *[]string, field_name string) []interface{} {
	if elements != nil {
		flattenStrings := make([]interface{}, len(*elements), len(*elements))
		for i, s := range *elements {
			m := make(map[string]interface{})
			m[field_name] = s
			flattenStrings[i] = m
		}
		return flattenStrings
	}
	return make([]interface{}, 0)
}

func FlattenListOfIntegers[T int | int32 | int64](elements *[]T, field_name string) []interface{} {
	if elements != nil {
		flattenStrings := make([]interface{}, len(*elements), len(*elements))
		for i, s := range *elements {
			m := make(map[string]interface{})
			m[field_name] = s
			flattenStrings[i] = m
		}
		return flattenStrings
	}
	return make([]interface{}, 0)
}

func FlattenListOfStringsList(elements *[][]string, fields []string) []interface{} {
	if elements != nil {
		flattenStrings := make([]interface{}, len(*elements), len(*elements))
		for i, s := range *elements {
			m := make(map[string]interface{})
			for j, f := range s {
				m[fields[j]] = f
			}
			flattenStrings[i] = m
		}
		return flattenStrings
	}
	return make([]interface{}, 0)
}

func BuildModelNameMapping(ctx context.Context, model interface{}) map[string]string {
	tflog.Debug(ctx, fmt.Sprintf("Building Model Name Mapping for %v", model))
	model_reflection := reflect.TypeOf(model)
	//Build Element names list map
	model_name_mapping := map[string]string{}
	for _, e := range reflect.VisibleFields(model_reflection) {
		model_name_mapping[e.Name] = GetTFformatName(e.Name)
	}
	return model_name_mapping
}

func FlattenModelToList(ctx context.Context, element interface{}) []interface{} {
	tflog.Info(ctx, fmt.Sprintf("%v", element))
	if element != nil {
		name_mapping := BuildModelNameMapping(ctx, element)
		flattenModel := make([]interface{}, 1, 1)
		m := make(map[string]interface{})
		model_value := reflect.ValueOf(element)
		for k, v := range name_mapping {
			value := model_value.FieldByName(k)
			switch value.Type().String() {
			case "int", "int32", "int64", "int8", "int16":
				m[v] = value.Int()
			case "string":
				m[v] = value.String()
			case "bool":
				m[v] = value.Bool()
			case "uint", "uint32", "uint64", "uint8", "uint16":
				m[v] = value.Uint()
			case "float32", "float64":
				m[v] = value.Float()
			}

			tflog.Info(ctx, fmt.Sprintf("%v", model_value.FieldByName(k).Type().String()))
			//m[v]=model_value.FieldByName(k)

		}
		flattenModel[0] = m
		return flattenModel
	}

	return make([]interface{}, 0)
}

func ModelToMap(ctx context.Context, model interface{}) map[string]interface{} {
	/*
	   Get A model and return it as a map of map[string]interface{}
	   where the string is the fieled name and the Value is the actual Value
	*/
	tflog.Debug(ctx, fmt.Sprintf("Processing %v", model))
	m := make(map[string]interface{})
	model_value := reflect.ValueOf(model)
	if model_value.Type().Kind() == reflect.Pointer && model_value.IsNil() {
		return m
	}
	name_mapping := BuildModelNameMapping(ctx, model)
	for k, v := range name_mapping {
		value := model_value.FieldByName(k)
		switch value.Type().String() {
		case "int", "int32", "int64", "int8", "int16":
			m[v] = value.Int()
		case "string":
			m[v] = value.String()
		case "bool":
			m[v] = value.Bool()
		case "uint", "uint32", "uint64", "uint8", "uint16":
			m[v] = value.Uint()
		case "float32", "float64":
			m[v] = value.Float()
		default:
			//If we got here it means that this is not a primitive but a complex type
			kind := value.Type().Kind()
			tflog.Debug(ctx, fmt.Sprintf("Model Kind %v", kind))
			if kind == reflect.Pointer && value.IsNil() {
				tflog.Debug(ctx, fmt.Sprintf("Nil Pointer processing return empty mapping"))
				m[v] = make([]interface{}, 0)
			} else if kind == reflect.Slice || kind == reflect.Array {
				//Handling Slices / Arrays
				tflog.Debug(ctx, fmt.Sprintf("Processing list of models"))
				l := []interface{}{}
				for i := 0; i < value.Len(); i++ {
					l = append(l, ModelToMap(ctx, value.Index(i).Interface()))
				}
				m[v] = l
			} else if kind == reflect.Pointer && !value.IsNil() {
				//Handling of pointers
				i := reflect.Indirect(value)
				tflog.Debug(ctx, fmt.Sprintf("Processing a pointer to %v", i))
				o := make([]interface{}, 1, 1)
				o[0] = ModelToMap(ctx, i.Interface())
				m[v] = o
			} else if kind == reflect.Struct {
				tflog.Debug(ctx, fmt.Sprintf("Processing struct %v", model))
				//Handlig Structs (we assume those are models)
				o := make([]interface{}, 1, 1)
				o[0] = ModelToMap(ctx, value.Interface())
				m[v] = o
			} else {
				panic(fmt.Sprintf("Can not Hanlde values from the type of %s", kind))
			}

		}

	}
	return m
}
func FlattenModelAsList(ctx context.Context, model interface{}) []interface{} {
	tflog.Debug(ctx, fmt.Sprintf("Flatenning model %v", model))
	model_value := reflect.ValueOf(model)
	tflog.Debug(ctx, fmt.Sprintf("Flatenning model %v,kind: %v", model, model_value.Type().Kind()))
	if model_value.Type().Kind() == reflect.Pointer && !model_value.IsNil() {
		tflog.Debug(ctx, fmt.Sprintf("Recived a non-nil pointer, will use the pointed value"))
		model = reflect.Indirect(model_value).Interface()
	}
	if model != nil {
		flattenModel := make([]interface{}, 1, 1)
		flattenModel[0] = ModelToMap(ctx, model)
		return flattenModel
	}
	return make([]interface{}, 0)
}

func FlattenListOfModelsToList(ctx context.Context, elements interface{}) []interface{} {
	tflog.Debug(ctx, fmt.Sprintf("Flattening List of models %v", elements))
	valueof := reflect.ValueOf(elements)
	if !(valueof.Type().Kind() == reflect.Array || valueof.Type().Kind() == reflect.Slice) {
		panic(fmt.Sprintf("elements must from the Kind of \"%s\" or \"&s\" , given type is \"%s\"", reflect.Array, reflect.Slice, valueof.Type().Kind()))
	}
	if elements != nil {
		flattenModel := make([]interface{}, valueof.Len(), valueof.Len())
		for i := 0; i < valueof.Len(); i++ {
			flattenModel[i] = ModelToMap(ctx, valueof.Index(i).Interface())
		}
		tflog.Debug(ctx, fmt.Sprintf("Flattened list of models list %v", flattenModel))
		return flattenModel

	}

	return make([]interface{}, 0)
}

func ReadSingleleValue(element interface{}) interface{} {
	return element
}

func ReadSingleList(elements *[]interface{}, field_name string) []interface{} {
	if elements != nil {
		response := make([]interface{}, len(*elements), len(*elements))
		for i, e := range *elements {
			t := e.(map[string]interface{})
			response[i] = t[field_name]
		}
		return response
	}
	return make([]interface{}, 0, 0)
}

func ReadDoubleList(elements *[]interface{}, fields_names []string) [][]interface{} {
	if elements != nil {
		response := make([][]interface{}, len(*elements), len(*elements))
		for i, e := range *elements {
			t := e.(map[string]interface{})
			n := make([]interface{}, len(fields_names), len(fields_names))
			for d, k := range fields_names {
				n[d] = t[k]
			}

			response[i] = n
		}
		return response
	}
	return make([][]interface{}, 0, 0)
}

func GetJsonTag(t reflect.StructTag) *string {
	tag, exists := t.Lookup("json")
	if !exists {
		return nil
	}
	s := strings.Split(tag, ",")
	if len(s) == 0 {
		return nil
	}
	return &s[0]

}

func PopulateResourceMap(ctx context.Context, t reflect.Type, d *schema.ResourceData, m *map[string]interface{}, prefix string, ignore_changes bool) {
	tflog.Debug(ctx, fmt.Sprintf("Populating Resource Map with %v", t))
	for _, fld := range reflect.VisibleFields(t) {
		tf_name := GetTFformatName(fld.Name)
		tflog.Debug(ctx, fmt.Sprintf("Processing Field %v of %v with tag: %v", fld, t, fld.Tag))
		tag := GetJsonTag(fld.Tag)
		if tag == nil {
			continue
		}
		full_tag := tf_name
		if prefix != "" {
			full_tag = prefix + "." + tf_name
		}
		tflog.Debug(ctx, fmt.Sprintf("Full Tag Name %s", full_tag))
		value, value_exists := d.GetOk(full_tag)
		if !value_exists {
			if fld.Type.Kind() == reflect.Slice {
				switch fld.Type.String() {
				case "[]string", "[]int", "[]int32", "[]int64", "[]float", "[]float32", "[]float64":
					(*m)[*tag] = value
				}
			}
			continue
		}
		tflog.Debug(ctx, fmt.Sprintf("Tag %s Exists , %v", full_tag, value))
		//Do not parse values values which have not been changed
		if !ignore_changes {
			changed := d.HasChange(full_tag)
			if !changed {
				continue
			}
		}
		if value_exists {
			tflog.Debug(ctx, fmt.Sprintf("Value of field %s", full_tag))
		}
		switch fld.Type.Kind() {
		case reflect.Struct, reflect.Pointer:
			tflog.Debug(ctx, fmt.Sprintf("Handling Struct/Pointer %v", fld))
			s := fld.Type.Elem()
			n := map[string]interface{}{}
			(*m)[*tag] = n
			PopulateResourceMap(ctx, s, d, &n, full_tag+".0", true)
		case reflect.String, reflect.Int, reflect.Int32, reflect.Float64, reflect.Int64, reflect.Bool:
			tflog.Debug(ctx, fmt.Sprintf("Converting string/int/bool %v", fld))
			(*m)[*tag] = value
		case reflect.Array, reflect.Slice:
			tflog.Debug(ctx, fmt.Sprintf("Converting array %v", fld))
			switch fld.Type.String() {
			case "[]string", "[]int", "[]int32", "[]int64", "[]float", "[]float32", "[]float64":
				(*m)[*tag] = value
			case "[][]string", "[][]int", "[][]int32", "[][]int64", "[][]float", "[][]float32", "[][]float64":
				out := [][]interface{}{}
				names := ctx.Value(ContextKey("names_mapping")).(map[string][]string)
				for _, v := range value.([]interface{}) {
					keys, exists := names[full_tag]
					if !exists {
						continue
					}
					t := v.(map[string]interface{})
					i := []interface{}{}
					for _, o := range keys {
						f, e := t[o]
						if e {
							i = append(i, f)
						}
					}
					out = append(out, i)
				}
				(*m)[*tag] = out
			default:
				/*If we got there than we assume the following
				  1. This is a single Array
				  2. It is made of struct
				*/
				list_type := fld.Type.Elem()
				//It must be a single list struct
				if list_type.Kind() != reflect.Struct {
					tflog.Debug(ctx, fmt.Sprintf("%s does not seems to be a slice of structs , skipping ", fld.Type.String()))
					continue

				}
				l := len(value.([]interface{}))
				o := make([]interface{}, l, l)
				(*m)[*tag] = o
				for i, _ := range value.([]interface{}) {
					n := map[string]interface{}{}
					o[i] = n
					PopulateResourceMap(ctx, list_type, d, &n, full_tag+"."+strconv.Itoa(i), true)

				}

			}

		}
	}
}
