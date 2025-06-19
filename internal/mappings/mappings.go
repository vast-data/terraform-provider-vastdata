package mappings

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	vast_client "github.com/vast-data/go-vast-client"
	"reflect"
	"strings"

	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SchemaContext int

const (
	SchemaForDataSource SchemaContext = iota
	SchemaForResource
)

type SearchContext int

const (
	SearchAny SearchContext = iota
	SearchRequited
	SearchOptional
	SearchSensitive
	SearchSearchable
)

func (sc SchemaContext) toSchemaTag() string {
	switch sc {
	case SchemaForDataSource:
		return "datasource"
	case SchemaForResource:
		return "resource"
	default:
		panic("unknown schema context")
	}
}

func (sc SchemaContext) toSchema() any {
	switch sc {
	case SchemaForDataSource:
		return datasource_schema.Schema{}
	case SchemaForResource:
		return resource_schema.Schema{}
	default:
		panic("unknown schema context")
	}
}

type attrMeta struct {
	Required   bool
	Optional   bool
	Computed   bool
	Sensitive  bool
	SearchAble bool
}

func (meta attrMeta) satisfySearchContext(search SearchContext) bool {
	switch search {
	case SearchAny:
		return meta.Required || meta.Optional
	case SearchRequited:
		return meta.Required
	case SearchOptional:
		return meta.Optional
	case SearchSensitive:
		return meta.Sensitive
	case SearchSearchable:
		return meta.SearchAble
	default:
		panic("unknown search context")
	}
}

func parseSchemaTag(tag string) attrMeta {
	meta := attrMeta{
		Computed: true, // default
	}

	tokens := strings.Split(tag, ",")
	for _, tok := range tokens {
		switch strings.TrimSpace(tok) {
		case "required":
			meta.Required = true
			meta.Computed = false
		case "optional":
			meta.Optional = true
		case "computed":
			meta.Computed = true
		case "searchable":
			meta.SearchAble = true
		case "sensitive":
			meta.Sensitive = true
		}
	}

	if meta.Required && meta.Optional {
		panic("attribute cannot be both required and optional")
	}
	return meta
}

func GenerateSchemaFromStruct(rs interface{}, kind SchemaContext) any {
	t := GetType(rs)
	switch kind {
	case SchemaForDataSource:
		return datasource_schema.Schema{
			Attributes: generateDatasourceAttributes(t, kind),
		}
	case SchemaForResource:
		return resource_schema.Schema{
			Attributes: generateResourceAttributes(t, kind),
		}
	default:
		panic("unknown schema context")
	}
}

func parseElementType(tag string) attr.Type {
	switch strings.TrimSpace(tag) {
	case "string":
		return types.StringType
	case "int64":
		return types.Int64Type
	case "float64":
		return types.Float64Type
	case "bool":
		return types.BoolType
	default:
		panic(fmt.Sprintf("unsupported element_type: %q", tag))
	}
}

func generateDatasourceAttributes(t reflect.Type, kind SchemaContext) map[string]datasource_schema.Attribute {
	attrs := make(map[string]datasource_schema.Attribute)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.PkgPath != "" {
			continue
		}

		tag := field.Tag.Get("tfsdk")
		tagName := strings.Split(tag, ",")[0]
		if tagName == "" || tagName == "-" {
			panic(fmt.Sprintf("field '%s' missing valid tfsdk tag", field.Name))
		}

		schemaContextTag := kind.toSchemaTag()
		meta := parseSchemaTag(field.Tag.Get(schemaContextTag))
		ft := field.Type

		switch ft {
		case reflect.TypeOf(types.String{}):
			attrs[tagName] = makeDatasourceAttr(datasource_schema.StringAttribute{}, meta)
		case reflect.TypeOf(types.Int64{}):
			attrs[tagName] = makeDatasourceAttr(datasource_schema.Int64Attribute{}, meta)
		case reflect.TypeOf(types.Float64{}):
			attrs[tagName] = makeDatasourceAttr(datasource_schema.Float64Attribute{}, meta)
		case reflect.TypeOf(types.Bool{}):
			attrs[tagName] = makeDatasourceAttr(datasource_schema.BoolAttribute{}, meta)
		case reflect.TypeOf(types.List{}):
			elementTag := field.Tag.Get("element_type")
			if elementTag == "" {
				panic(fmt.Sprintf("field '%s' is types.List but missing `element_type` tag", field.Name))
			}
			attrs[tagName] = makeDatasourceAttr(datasource_schema.ListAttribute{
				ElementType: parseElementType(elementTag),
			}, meta)
		case reflect.TypeOf(types.Set{}):
			elementTag := field.Tag.Get("element_type")
			if elementTag == "" {
				panic(fmt.Sprintf("field '%s' is types.Set but missing `element_type` tag", field.Name))
			}
			attrs[tagName] = makeDatasourceAttr(datasource_schema.SetAttribute{
				ElementType: parseElementType(elementTag),
			}, meta)

		case reflect.TypeOf(types.Map{}):
			elementTag := field.Tag.Get("element_type")
			if elementTag == "" {
				panic(fmt.Sprintf("field '%s' is types.Map but missing `element_type` tag", field.Name))
			}
			attrs[tagName] = makeDatasourceAttr(datasource_schema.MapAttribute{
				ElementType: parseElementType(elementTag),
			}, meta)
		default:
			if ft.Kind() == reflect.Struct {
				attrs[tagName] = makeDatasourceAttr(datasource_schema.SingleNestedAttribute{
					Attributes: generateDatasourceAttributes(ft, kind),
				}, meta)
			} else {
				panic(fmt.Sprintf("field '%s' has unsupported type: %s", field.Name, ft.String()))
			}
		}
	}
	return attrs

}

func generateResourceAttributes(t reflect.Type, kind SchemaContext) map[string]resource_schema.Attribute {
	attrs := make(map[string]resource_schema.Attribute)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.PkgPath != "" {
			continue
		}

		tag := field.Tag.Get("tfsdk")
		tagName := strings.Split(tag, ",")[0]
		if tagName == "" || tagName == "-" {
			panic(fmt.Sprintf("field '%s' missing valid tfsdk tag", field.Name))
		}

		schemaContextTag := kind.toSchemaTag()
		meta := parseSchemaTag(field.Tag.Get(schemaContextTag))
		ft := field.Type

		switch ft {
		case reflect.TypeOf(types.String{}):
			attrs[tagName] = makeResourceAttr(resource_schema.StringAttribute{}, meta)
		case reflect.TypeOf(types.Int64{}):
			attrs[tagName] = makeResourceAttr(resource_schema.Int64Attribute{}, meta)
		case reflect.TypeOf(types.Float64{}):
			attrs[tagName] = makeResourceAttr(resource_schema.Float64Attribute{}, meta)
		case reflect.TypeOf(types.Bool{}):
			attrs[tagName] = makeResourceAttr(resource_schema.BoolAttribute{}, meta)
		case reflect.TypeOf(types.List{}):
			elementTag := field.Tag.Get("element_type")
			if elementTag == "" {
				panic(fmt.Sprintf("field '%s' is types.List but missing `element_type` tag", field.Name))
			}
			attrs[tagName] = makeResourceAttr(resource_schema.ListAttribute{
				ElementType: parseElementType(elementTag),
			}, meta)
		case reflect.TypeOf(types.Set{}):
			elementTag := field.Tag.Get("element_type")
			if elementTag == "" {
				panic(fmt.Sprintf("field '%s' is types.Set but missing `element_type` tag", field.Name))
			}
			attrs[tagName] = makeResourceAttr(resource_schema.SetAttribute{
				ElementType: parseElementType(elementTag),
			}, meta)

		case reflect.TypeOf(types.Map{}):
			elementTag := field.Tag.Get("element_type")
			if elementTag == "" {
				panic(fmt.Sprintf("field '%s' is types.Map but missing `element_type` tag", field.Name))
			}
			attrs[tagName] = makeResourceAttr(resource_schema.MapAttribute{
				ElementType: parseElementType(elementTag),
			}, meta)
		default:
			if ft.Kind() == reflect.Struct {
				attrs[tagName] = makeResourceAttr(resource_schema.SingleNestedAttribute{
					Attributes: generateResourceAttributes(ft, kind),
				}, meta)
			} else {
				panic(fmt.Sprintf("field '%s' has unsupported type: %s", field.Name, ft.String()))
			}
		}
	}
	return attrs

}

func makeDatasourceAttr(attr datasource_schema.Attribute, meta attrMeta) datasource_schema.Attribute {
	switch a := attr.(type) {
	case datasource_schema.StringAttribute:
		a.Required = meta.Required
		a.Optional = meta.Optional
		a.Computed = meta.Computed
		a.Sensitive = meta.Sensitive
		return a
	case datasource_schema.Int64Attribute:
		a.Required = meta.Required
		a.Optional = meta.Optional
		a.Computed = meta.Computed
		a.Sensitive = meta.Sensitive
		return a
	case datasource_schema.Float64Attribute:
		a.Required = meta.Required
		a.Optional = meta.Optional
		a.Computed = meta.Computed
		a.Sensitive = meta.Sensitive
		return a
	case datasource_schema.BoolAttribute:
		a.Required = meta.Required
		a.Optional = meta.Optional
		a.Computed = meta.Computed
		a.Sensitive = meta.Sensitive
		return a
	case datasource_schema.ListAttribute:
		a.Required = meta.Required
		a.Optional = meta.Optional
		a.Computed = meta.Computed
		a.Sensitive = meta.Sensitive
		return a
	case datasource_schema.SetAttribute:
		a.Required = meta.Required
		a.Optional = meta.Optional
		a.Computed = meta.Computed
		a.Sensitive = meta.Sensitive
		return a
	case datasource_schema.ListNestedAttribute:
		a.Required = meta.Required
		a.Optional = meta.Optional
		a.Computed = meta.Computed
		a.Sensitive = meta.Sensitive
		return a
	case datasource_schema.SetNestedAttribute:
		a.Required = meta.Required
		a.Optional = meta.Optional
		a.Computed = meta.Computed
		a.Sensitive = meta.Sensitive
		return a
	case datasource_schema.SingleNestedAttribute:
		a.Required = meta.Required
		a.Optional = meta.Optional
		a.Computed = meta.Computed
		a.Sensitive = meta.Sensitive
		return a
	default:
		panic(fmt.Sprintf("unsupported datasource schema attribute type: %T", attr))
	}
}

func makeResourceAttr(attr resource_schema.Attribute, meta attrMeta) resource_schema.Attribute {
	switch a := attr.(type) {
	case resource_schema.StringAttribute:
		a.Required = meta.Required
		a.Optional = meta.Optional
		a.Computed = meta.Computed
		a.Sensitive = meta.Sensitive
		return a
	case resource_schema.Int64Attribute:
		a.Required = meta.Required
		a.Optional = meta.Optional
		a.Computed = meta.Computed
		a.Sensitive = meta.Sensitive
		return a
	case resource_schema.Float64Attribute:
		a.Required = meta.Required
		a.Optional = meta.Optional
		a.Computed = meta.Computed
		a.Sensitive = meta.Sensitive
		return a
	case resource_schema.BoolAttribute:
		a.Required = meta.Required
		a.Optional = meta.Optional
		a.Computed = meta.Computed
		a.Sensitive = meta.Sensitive
		return a
	case resource_schema.ListAttribute:
		a.Required = meta.Required
		a.Optional = meta.Optional
		a.Computed = meta.Computed
		a.Sensitive = meta.Sensitive
		return a
	case resource_schema.SetAttribute:
		a.Required = meta.Required
		a.Optional = meta.Optional
		a.Computed = meta.Computed
		a.Sensitive = meta.Sensitive
		return a
	case resource_schema.ListNestedAttribute:
		a.Required = meta.Required
		a.Optional = meta.Optional
		a.Computed = meta.Computed
		a.Sensitive = meta.Sensitive
		return a
	case resource_schema.SetNestedAttribute:
		a.Required = meta.Required
		a.Optional = meta.Optional
		a.Computed = meta.Computed
		a.Sensitive = meta.Sensitive
		return a
	case resource_schema.SingleNestedAttribute:
		a.Required = meta.Required
		a.Optional = meta.Optional
		a.Computed = meta.Computed
		a.Sensitive = meta.Sensitive
		return a
	default:
		panic(fmt.Sprintf("unsupported resource schema attribute type: %T", attr))
	}
}

// GetNotEmptyFields returns a map of attribute names to values for all fields in the input struct
// that are marked as `required` or `optional` in the struct tag corresponding to the given SchemaContext,
// and are not null or unknown.
//
// It is useful for extracting user-supplied input values before calling API clients or performing validation.
//
// Supported values must implement the Terraform Plugin Framework value interface with IsNull() and IsUnknown() methods.
// ExtractNonNullConfiguredFields returns a map of attribute names to concrete Go values
// for fields tagged as required or optional and that are not null or unknown.
func GetNotEmptyFields(input any, kind SchemaContext, search SearchContext) (map[string]any, error) {
	result := make(map[string]any)

	t, v, err := resolveStructValue(input)
	if err != nil {
		return nil, err
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %T", input)
	}

	schemaContextTag := kind.toSchemaTag()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue // unexported
		}

		tfsdkTag := field.Tag.Get("tfsdk")
		tagName := strings.Split(tfsdkTag, ",")[0]
		if tagName == "" || tagName == "-" {
			continue
		}

		meta := parseSchemaTag(field.Tag.Get(schemaContextTag))
		if !meta.satisfySearchContext(search) {
			continue
		}

		fieldVal := v.Field(i)
		if fieldVal.Kind() != reflect.Struct {
			continue
		}

		isNullMethod := fieldVal.MethodByName("IsNull")
		isUnknownMethod := fieldVal.MethodByName("IsUnknown")
		if !isNullMethod.IsValid() || !isUnknownMethod.IsValid() {
			continue
		}
		if isNullMethod.Call(nil)[0].Bool() || isUnknownMethod.Call(nil)[0].Bool() {
			continue
		}

		// Extract actual Go value from Terraform type
		switch fv := fieldVal.Interface().(type) {
		case types.String:
			result[tagName] = fv.ValueString()
		case types.Int64:
			result[tagName] = fv.ValueInt64()
		case types.Bool:
			result[tagName] = fv.ValueBool()
		case types.Float64:
			result[tagName] = fv.ValueFloat64()
		case types.List:
			var unpacked []any
			for _, el := range fv.Elements() {
				switch e := el.(type) {
				case types.String:
					unpacked = append(unpacked, e.ValueString())
				case types.Int64:
					unpacked = append(unpacked, e.ValueInt64())
				case types.Float64:
					unpacked = append(unpacked, e.ValueFloat64())
				case types.Bool:
					unpacked = append(unpacked, e.ValueBool())
				default:
					panic(fmt.Sprintf("unsupported element type %T in list for field %s", e, tagName))
				}
			}
			result[tagName] = unpacked
		case types.Set:
			var unpacked []any
			for _, el := range fv.Elements() {
				switch e := el.(type) {
				case types.String:
					unpacked = append(unpacked, e.ValueString())
				case types.Int64:
					unpacked = append(unpacked, e.ValueInt64())
				case types.Float64:
					unpacked = append(unpacked, e.ValueFloat64())
				case types.Bool:
					unpacked = append(unpacked, e.ValueBool())
				default:
					panic(fmt.Sprintf("unsupported element type %T in set for field %s", e, tagName))
				}
			}
			result[tagName] = unpacked
		case types.Map:
			unpacked := make(map[string]any)
			for k, v := range fv.Elements() {
				switch e := v.(type) {
				case types.String:
					unpacked[k] = e.ValueString()
				case types.Int64:
					unpacked[k] = e.ValueInt64()
				case types.Float64:
					unpacked[k] = e.ValueFloat64()
				case types.Bool:
					unpacked[k] = e.ValueBool()
				default:
					panic(fmt.Sprintf("unsupported element type %T in map for field %s", e, tagName))
				}
			}
			result[tagName] = unpacked

		default:
			panic(fmt.Sprintf("unsupported type %T for field %s", fv, tagName))
		}
	}

	return result, nil
}

// GetIdPtr extracts the value of the field tagged with `tfsdk:"id"` as a pointer to int64.
// It returns nil if the field is not found or is null/unknown.
func GetIdPtr(input any) (*int64, error) {
	t, v, err := resolveStructValue(input)
	if err != nil {
		return nil, err
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %T", input)
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue // skip unexported
		}

		tfsdkTag := field.Tag.Get("tfsdk")
		tagName := strings.Split(tfsdkTag, ",")[0]
		if tagName != "id" {
			continue
		}

		fieldVal := v.Field(i)
		if fieldVal.Kind() != reflect.Struct {
			continue
		}

		isNullMethod := fieldVal.MethodByName("IsNull")
		isUnknownMethod := fieldVal.MethodByName("IsUnknown")
		if !isNullMethod.IsValid() || !isUnknownMethod.IsValid() {
			continue
		}
		if isNullMethod.Call(nil)[0].Bool() || isUnknownMethod.Call(nil)[0].Bool() {
			return nil, nil
		}

		if idVal, ok := fieldVal.Interface().(types.Int64); ok {
			val := idVal.ValueInt64()
			return &val, nil
		}
	}

	return nil, nil
}

// GetOptionalFieldNames returns a list of field names tagged as "optional"
// for the given SchemaContext (e.g., "datasource" or "resource").
func GetOptionalFieldNames(input interface{}, kind SchemaContext) ([]string, error) {
	return getFieldNamesByTag(input, kind, "optional")
}

// GetRequiredFieldNames returns a list of field names tagged as "required"
// for the given SchemaContext (e.g., "datasource" or "resource").
func GetRequiredFieldNames(input interface{}, kind SchemaContext) ([]string, error) {
	return getFieldNamesByTag(input, kind, "required")
}

// HasRequiredFields returns true if there is at least one field tagged as "required"
// for the given SchemaContext in the provided struct.
func HasRequiredFields(input interface{}, kind SchemaContext) (bool, error) {
	requiredFieldNames, err := getFieldNamesByTag(input, kind, "required")
	if err != nil {
		return false, err
	}
	return len(requiredFieldNames) > 0, nil
}

// HasOptionalFields returns true if there is at least one field tagged as "optional"
// for the given SchemaContext in the provided struct.
func HasOptionalFields(input interface{}, kind SchemaContext) (bool, error) {
	optionalFieldNames, err := getFieldNamesByTag(input, kind, "optional")
	if err != nil {
		return false, err
	}
	return len(optionalFieldNames) > 0, nil
}

// getFieldNamesByTag is an internal helper that returns field names by tag (e.g. "optional", "required").
func getFieldNamesByTag(input interface{}, kind SchemaContext, targetTag string) ([]string, error) {
	t := reflect.TypeOf(input)
	if t.Kind() != reflect.Struct {
		t = reflect.Indirect(reflect.ValueOf(input)).Type()
		if t.Kind() != reflect.Struct {
			return nil, fmt.Errorf("expected struct but got %s", t.Kind())
		}
	}

	schemaContextTag := kind.toSchemaTag()
	var result []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.PkgPath != "" {
			continue // unexported
		}

		tfTag := field.Tag.Get("tfsdk")
		name := strings.Split(tfTag, ",")[0]
		if name == "" || name == "-" {
			continue
		}

		schemaTag := field.Tag.Get(schemaContextTag)
		for _, token := range strings.Split(schemaTag, ",") {
			if strings.TrimSpace(token) == targetTag {
				result = append(result, name)
				break
			}
		}
	}

	return result, nil
}

// GetNotEmptyOptionalFieldNames returns a list of "optional" field names that are not null for the given schema context.
func GetNotEmptyOptionalFieldNames(input interface{}, kind SchemaContext) ([]string, error) {
	return getNotEmptyFieldNamesByTag(input, kind, "optional")
}

// GetNotEmptyRequiredFieldNames returns a list of "required" field names that are not null for the given schema context.
func GetNotEmptyRequiredFieldNames(input interface{}, kind SchemaContext) ([]string, error) {
	return getNotEmptyFieldNamesByTag(input, kind, "required")
}

// HasNotEmptyOptionalFields returns true if there is at least one field
// tagged as "optional" and its value is not null for the given SchemaContext.
func HasNotEmptyOptionalFields(input interface{}, kind SchemaContext) (bool, error) {
	fields, err := getNotEmptyFieldNamesByTag(input, kind, "optional")
	if err != nil {
		return false, err
	}
	return len(fields) > 0, nil
}

// HasNotEmptyRequiredFields returns true if there is at least one field
// tagged as "required" and its value is not null for the given SchemaContext.
func HasNotEmptyRequiredFields(input interface{}, kind SchemaContext) (bool, error) {
	fields, err := getNotEmptyFieldNamesByTag(input, kind, "required")
	if err != nil {
		return false, err
	}
	return len(fields) > 0, nil
}

// HasAnyNotEmptyFields returns true if the input struct contains at least one field
// that is either tagged as "required" or "optional" and is not null, based on the given SchemaContext.
func HasAnyNotEmptyFields(input interface{}, kind SchemaContext) (bool, error) {
	required, err := HasNotEmptyRequiredFields(input, kind)
	if err != nil {
		return false, err
	}
	if required {
		return true, nil
	}

	optional, err := HasNotEmptyOptionalFields(input, kind)
	if err != nil {
		return false, err
	}
	return optional, nil
}

// getNotEmptyFieldNamesByTag returns field names that match the given tag (e.g., "optional", "required") and are not null.
func getNotEmptyFieldNamesByTag(input interface{}, kind SchemaContext, targetTag string) ([]string, error) {
	t, v, err := resolveStructValue(input)
	if err != nil {
		return nil, err
	}

	schemaContextTag := kind.toSchemaTag()
	var result []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.PkgPath != "" {
			continue
		}

		tfTag := field.Tag.Get("tfsdk")
		name := strings.Split(tfTag, ",")[0]
		if name == "" || name == "-" {
			continue
		}

		// Match tag
		schemaTag := field.Tag.Get(schemaContextTag)
		hasTag := false
		for _, token := range strings.Split(schemaTag, ",") {
			if strings.TrimSpace(token) == targetTag {
				hasTag = true
				break
			}
		}
		if !hasTag {
			continue
		}

		// Check IsNull
		val := v.Field(i)
		if val.Kind() == reflect.Struct {
			isNull := val.MethodByName("IsNull")
			if isNull.IsValid() {
				res := isNull.Call(nil)
				if len(res) == 1 && !res[0].Bool() {
					result = append(result, name)
				}
			}
		}
	}

	return result, nil
}

// resolveStructValue dereferences a pointer to a struct if needed and returns the struct's reflect.Type and reflect.Value.
// It returns an error if the input is not a struct or a pointer to a struct.
func resolveStructValue(input any) (reflect.Type, reflect.Value, error) {
	v := reflect.ValueOf(input)
	if !v.IsValid() {
		return nil, reflect.Value{}, fmt.Errorf("input is invalid")
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, reflect.Value{}, fmt.Errorf("input is a nil pointer")
		}
		v = v.Elem()
	}

	t := v.Type()
	if t.Kind() != reflect.Struct {
		return nil, reflect.Value{}, fmt.Errorf("expected struct but got %s", t.Kind())
	}

	return t, v, nil
}

func FillFromRecord(r vast_client.Record, container any) error {
	v := reflect.ValueOf(container).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("tfsdk")
		if tag == "" || tag == "-" {
			continue
		}
		key := strings.Split(tag, ",")[0]
		val, ok := r[key]

		fieldValue := v.Field(i)
		if !fieldValue.CanSet() {
			continue
		}

		switch field.Type {
		case reflect.TypeOf(types.String{}):
			if !ok || isNil(val) {
				fieldValue.Set(reflect.ValueOf(types.StringNull()))
			} else {
				fieldValue.Set(reflect.ValueOf(types.StringValue(fmt.Sprintf("%v", val))))
			}

		case reflect.TypeOf(types.Int64{}):
			if !ok || isNil(val) {
				fieldValue.Set(reflect.ValueOf(types.Int64Null()))
			} else {
				intVal, err := toInt(val)
				if err != nil {
					return fmt.Errorf("error converting %q to int64: %w", key, err)
				}
				fieldValue.Set(reflect.ValueOf(types.Int64Value(intVal)))
			}

		case reflect.TypeOf(types.Float64{}):
			if !ok || isNil(val) {
				fieldValue.Set(reflect.ValueOf(types.Float64Null()))
			} else {
				floatVal, err := toFloat(val)
				if err != nil {
					return fmt.Errorf("error converting %q to float64: %w", key, err)
				}
				fieldValue.Set(reflect.ValueOf(types.Float64Value(floatVal)))
			}

		case reflect.TypeOf(types.Bool{}):
			if !ok || isNil(val) {
				fieldValue.Set(reflect.ValueOf(types.BoolNull()))
			} else {
				boolVal, isBool := val.(bool)
				if !isBool {
					return fmt.Errorf("value for %q is not bool", key)
				}
				fieldValue.Set(reflect.ValueOf(types.BoolValue(boolVal)))
			}

		case reflect.TypeOf(types.List{}):
			elementTag := field.Tag.Get("element_type")
			if elementTag == "" {
				return fmt.Errorf("field %q is types.List but missing `element_type` tag", field.Name)
			}
			elemType := parseElementType(elementTag)

			if !ok || isNil(val) {
				fieldValue.Set(reflect.ValueOf(types.ListNull(elemType)))
				continue
			}

			var elements []attr.Value
			switch elemType {
			case types.StringType:
				raw, ok := val.([]any)
				if !ok {
					rawStrings, ok := val.([]string)
					if !ok {
						return fmt.Errorf("field %q expects []string", key)
					}
					for _, s := range rawStrings {
						elements = append(elements, types.StringValue(s))
					}
				} else {
					for _, s := range raw {
						strVal := fmt.Sprintf("%v", s)
						elements = append(elements, types.StringValue(strVal))
					}
				}

			case types.Int64Type:
				raw, ok := val.([]any)
				if !ok {
					rawInts, ok := val.([]int64)
					if !ok {
						return fmt.Errorf("field %q expects []int64", key)
					}
					for _, i := range rawInts {
						elements = append(elements, types.Int64Value(i))
					}
				} else {
					for _, n := range raw {
						intVal, err := toInt(n)
						if err != nil {
							return fmt.Errorf("error converting %q element to int64: %w", key, err)
						}
						elements = append(elements, types.Int64Value(intVal))
					}
				}

			default:
				return fmt.Errorf("unsupported element_type %q in list for field %q", elementTag, key)
			}

			if len(elements) == 0 {
				fieldValue.Set(reflect.ValueOf(types.ListNull(elemType)))
			} else {
				fieldValue.Set(reflect.ValueOf(types.ListValueMust(elemType, elements)))
			}

		case reflect.TypeOf(types.Set{}):
			elementTag := field.Tag.Get("element_type")
			if elementTag == "" {
				return fmt.Errorf("field %q is types.Set but missing `element_type` tag", field.Name)
			}
			elemType := parseElementType(elementTag)

			if !ok || isNil(val) {
				fieldValue.Set(reflect.ValueOf(types.SetNull(elemType)))
				continue
			}

			var elements []attr.Value
			switch elemType {
			case types.StringType:
				raw, ok := val.([]any)
				if !ok {
					rawStrings, ok := val.([]string)
					if !ok {
						return fmt.Errorf("field %q expects []string", key)
					}
					for _, s := range rawStrings {
						elements = append(elements, types.StringValue(s))
					}
				} else {
					for _, s := range raw {
						strVal := fmt.Sprintf("%v", s)
						elements = append(elements, types.StringValue(strVal))
					}
				}

			case types.Int64Type:
				raw, ok := val.([]any)
				if !ok {
					rawInts, ok := val.([]int64)
					if !ok {
						return fmt.Errorf("field %q expects []int64", key)
					}
					for _, i := range rawInts {
						elements = append(elements, types.Int64Value(i))
					}
				} else {
					for _, n := range raw {
						intVal, err := toInt(n)
						if err != nil {
							return fmt.Errorf("error converting %q element to int64: %w", key, err)
						}
						elements = append(elements, types.Int64Value(intVal))
					}
				}

			default:
				return fmt.Errorf("unsupported element_type %q in set for field %q", elementTag, key)
			}

			if len(elements) == 0 {
				fieldValue.Set(reflect.ValueOf(types.SetNull(elemType)))
			} else {
				fieldValue.Set(reflect.ValueOf(types.SetValueMust(elemType, elements)))
			}

		default:
			// Handle []struct{} (nested slices of structs)
			if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct {
				if !ok || isNil(val) {
					fieldValue.Set(reflect.Zero(field.Type))
					continue
				}

				rawList, ok := val.([]any)
				if !ok {
					return fmt.Errorf("value for %q is not a slice", key)
				}

				slice := reflect.MakeSlice(field.Type, len(rawList), len(rawList))
				for j, item := range rawList {
					itemMap, ok := item.(map[string]any)
					if !ok {
						return fmt.Errorf("element %d in %q is not a map", j, key)
					}
					elemPtr := slice.Index(j).Addr().Interface()
					if err := FillFromRecord(vast_client.Record(itemMap), elemPtr); err != nil {
						return fmt.Errorf("error filling nested struct in %q: %w", key, err)
					}
				}
				fieldValue.Set(slice)
			} else {
				panic(fmt.Sprintf("unsupported field type %s for field %q", field.Type, key))
			}
		}
	}

	return nil
}
