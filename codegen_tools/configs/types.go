package configs

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	base "github.com/pb33f/libopenapi/datamodel/high/base"
	utils "github.com/vast-data/terraform-provider-vastdata/utils"
	vast_versions "github.com/vast-data/terraform-provider-vastdata/vast_versions"
)

func ToStringPointer(s string) *string {
	return &s
}

func GetFuncRunTimeTypeName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

}

func LastSplit(s, sp string) string {
	i := strings.Split(s, sp)
	l := len(i)
	if l > 0 {
		return i[l-1]
	}
	return ""
}

func FuncName(i interface{}) string {
	n := GetFuncRunTimeTypeName(i)
	return LastSplit(n, "/")
}

type StringSet struct {
	_set map[string]struct{}
}

// Add a string key to the set
func (s *StringSet) Add(k string) {
	s._set[k] = struct{}{}
}

// Check if a string key is in the set
func (s *StringSet) In(k string) bool {
	_, exists := s._set[k]
	return exists
}

// Remove a string key from set
func (s *StringSet) Del(k string) {
	if s.In(k) {
		delete(s._set, k)
	}
}

func (s *StringSet) ToArray() []string {
	arr := []string{}
	for k, _ := range s._set {
		arr = append(arr, k)
	}
	return arr
}

func (s *StringSet) QueryStringFmt() string {
	arr := s.ToArray()
	q1 := ""
	q2 := ""
	for _, a := range arr {
		q1 += a + "=%v&"
		q2 += "," + a
	}
	return "\"" + strings.TrimRight(q1, "&") + "\"" + q2
}

// Get a new StringSet that can be populated from the beginning with strings
func NewStringSet(s ...string) *StringSet {
	new_string_set := StringSet{_set: make(map[string]struct{})}
	for _, elem := range s {
		new_string_set.Add(elem)

	}
	return &new_string_set
}

func ListAsStringsList(s []interface{}) []string {
	//convert any list to a list of strings with the value as the string for each element
	n := make([]string, len(s), len(s))
	for i, j := range s {
		n[i] = fmt.Sprintf("%v", j)
	}
	return n
}

type ResourceElem struct {
	Attributes   map[string]string
	ResourceElem *ResourceElem
	Indent       int
	Parent       *ResourceTemplateV2
}

func (r *ResourceElem) IsReferance() bool {
	if r.ResourceElem != nil {
		return true
	}
	return false
}

func NewResourceElem(m map[string]string, p *ResourceTemplateV2) ResourceElem {
	return ResourceElem{
		Attributes:   m,
		ResourceElem: nil,
		Parent:       p,
	}

}

type FakeField struct {
	Name        string
	Description string
}

type ResourceTemplateV2 struct {
	ResourceName             string
	Fields                   []ResourceElem
	Path                     *string
	Model                    interface{}
	DestFile                 *string
	IgnoreFields             *StringSet
	RequiredIdentifierFields *StringSet
	OptionalIdentifierFields *StringSet
	PathIdentifierFields     *StringSet
	ComputedFields           *StringSet
	ForceNewFields           *StringSet
	ConflictingFields        map[string][]string
	ListsNamesMap            map[string][]string
	Generate                 bool
	DisableImport            bool
	DataSourceName           string
	ResponseProcessingFunc   utils.ResponseProcessingFunc
	ResponseGetByURL         bool
	IgnoreUpdates            *StringSet
	TfNameToModelName        map[string]string
	ListFields               map[string][]FakeField
	ApiSchema                *base.SchemaProxy
	ResourceDocumantation    string
	BeforePostFunc           utils.ResponseConversionFunc
	BeforePatchFunc          utils.ResponseConversionFunc
	AfterPostFunc            utils.ResponseConversionFunc
	AfterPatchFunc           utils.ResponseConversionFunc
	AfterReadFunc            utils.SchemaManipulationFunc
	BeforeDeleteFunc         utils.PreDeleteFunc
	FieldsValidators         map[string]schema.SchemaValidateDiagFunc
	SensitiveFields          *StringSet
	IsDataSource             bool
	BeforeCreateFunc         utils.ResponseConversionFunc
	CreateFunc               utils.CreateFuncType
	UpdateFunc               utils.UpdateFuncType
	DeleteFunc               utils.DeleteFuncType
	GetFunc                  utils.GetFuncType
	IdFunc                   utils.IdFuncType
	ImportFunc               utils.ImportFunc
	Importer                 utils.ImportInterface
	AttributesDiffFuncs      map[string]schema.SchemaDiffSuppressFunc
	Timeouts                 *schema.ResourceTimeout
	DisableAutoValidator     *StringSet
	DisableFallbackRequest   bool
}

func (r *ResourceTemplateV2) AutomaticValidationIsDisabled(s string) bool {
	if r.DisableAutoValidator == nil {
		return false
	}
	return r.DisableAutoValidator.In(s)
}

func (r *ResourceTemplateV2) GetResourceTimeouts() *schema.ResourceTimeout {
	return r.Timeouts
}

func (r *ResourceTemplateV2) GetAttributeDiffFunc(attr string) schema.SchemaDiffSuppressFunc {
	return r.AttributesDiffFuncs[attr]
}

func (r *ResourceTemplateV2) AttributeHasDiffFunc(attr string) bool {
	_, exists := r.AttributesDiffFuncs[attr]
	return exists
}

func (r *ResourceTemplateV2) SetFunctions() {
	if r.CreateFunc == nil {
		r.CreateFunc = utils.DefaultCreateFunc
	}
	if r.DeleteFunc == nil {
		r.DeleteFunc = utils.DefaultDeleteFunc
	}
	if r.UpdateFunc == nil {
		r.UpdateFunc = utils.DefaultUpdateFunc
	}
	if r.GetFunc == nil {
		r.GetFunc = utils.DefaultGetFunc
	}
	if r.IdFunc == nil {
		r.IdFunc = utils.DefaultIdFunc
	}
	if r.ImportFunc == nil {
		r.ImportFunc = utils.DefaultImportFunc
	}
	if r.Importer == nil {
		r.Importer = utils.GetDefaultImporter()
	}
	if r.ResponseProcessingFunc == nil {
		r.ResponseProcessingFunc = utils.DefaultProcessingFunc
	}
}

func (r *ResourceTemplateV2) GetProperty(property string) *base.SchemaProxy {
	if r.ApiSchema != nil {
		schema := r.ApiSchema.Schema()
		p := schema.Properties
		for pair := p.First(); pair != nil; pair = pair.Next() {
			if pair.Key() == property {
				return pair.Value()
			}
		}

	}
	return nil
}

func (r *ResourceTemplateV2) GetEnum(property string) []interface{} {
	e := []interface{}{}
	property_schema_proxy := r.GetProperty(property)
	if property_schema_proxy == nil {
		return e
	}
	property_schema := property_schema_proxy.Schema()
	enum := property_schema.Enum
	if len(enum) == 0 {
		return e
	}
	for _, r := range enum {
		e = append(e, r.Value)
	}
	return e
}

func (r *ResourceTemplateV2) GetSchemaProperySupportedVersions(property string) string {
	l := vast_versions.VersionsSupportingAttributes(r.ResourceName, property)
	if len(l) == 0 {
		return "Valid for versions: Unkown"
	}
	e := "Valid for versions: "
	for _, s := range l {
		e = e + fmt.Sprintf("%v,", s)
	}
	return strings.TrimSuffix(e, ",")
}

func (r ResourceTemplateV2) HasValidatorFunc(s string) bool {
	_, exists := r.FieldsValidators[s]
	return exists
}

func (r ResourceTemplateV2) GetValidatorFunc(s string) string {
	if r.HasValidatorFunc(s) {
		f, _ := r.FieldsValidators[s]
		return FuncName(f)
	}
	return ""

}

func (r *ResourceTemplateV2) GetSchemaDocumentation() string {
	if r.ApiSchema != nil {
		return fmt.Sprintf("%v", r.ApiSchema.Schema().Description)
	}
	return ""
}

func (r *ResourceTemplateV2) GetSchemaProperyDocument(property string) string {
	if r.ApiSchema != nil {
		property_schema_proxy := r.GetProperty(property)
		if property_schema_proxy == nil {
			return ""
		}
		property_schema := property_schema_proxy.Schema()
		enum := r.GetEnum(property)
		out := fmt.Sprintf("%v", property_schema.Description)
		if len(enum) > 0 {
			out += fmt.Sprintf(" Allowed Values are %v", ListAsStringsList(enum))
		}
		out = fmt.Sprintf("(%s) %s", r.GetSchemaProperySupportedVersions(property), out)
		return out
	}
	return ""
}

func (r *ResourceTemplateV2) GetSchemaProperyDefault(property string) string {
	if r.ApiSchema != nil {
		property_schema_proxy := r.GetProperty(property)
		if property_schema_proxy == nil {
			return ""
		}
		property_schema := property_schema_proxy.Schema()
		_default := property_schema.Default
		if _default == nil {
			return ""
		}
		if property_schema.Type[0] == "string" {
			return fmt.Sprintf(`"%v"`, _default.Value)
		}
		return fmt.Sprintf("%v", _default.Value)
	}

	return ""
}

func (r *ResourceTemplateV2) ConvertTfNameToModelName(tf_name string) string {
	name, exists := r.TfNameToModelName[tf_name]
	if exists {
		return name
	}
	return ""
}

func (r *ResourceTemplateV2) GetFakeFieldDescription(fieldName string, fieldListName string) string {
	// Get the description of fake filed X item in list Y , if either X or Y are not founf "" is returned
	f, exists := r.ListFields[fieldName]
	if !exists {
		return ""
	}
	for _, fk := range f {
		if fk.Name == fieldListName {
			return fk.Description
		}
	}

	return ""

}

func (r *ResourceTemplateV2) HasFakeField(s string) bool {
	_, exists := r.ListFields[s]
	return exists
}

func (r *ResourceTemplateV2) SetupFakeFields(s string, f []FakeField, m *map[string]string) {
	re := regexp.MustCompile("\\[\\]")
	t := re.FindAllString(s, -1)
	(*m)["type"] = "TypeList"
	(*m)["list_type"] = "simple"
	(*m)["length"] = strconv.Itoa(len(t))

	n := make([]string, len(f), len(f))
	for i, j := range f {
		n[i] = j.Name
	}
	(*m)["names"] = strings.Join(n, ",")
	if strings.Contains(s, "string") {
		(*m)["set_type"] = "String"
	} else if strings.Contains(s, "int") {
		(*m)["set_type"] = "Int"
	} else if strings.Contains(s, "Float") {
		(*m)["set_type"] = "Float"
	}

}

func (r *ResourceTemplateV2) SetupListProperties(s string, m *map[string]string) {
	re := regexp.MustCompile("\\[\\]")
	t := re.FindAllString(s, -1)
	(*m)["type"] = "TypeList"
	(*m)["list_type"] = "simple"
	(*m)["length"] = strconv.Itoa(len(t))
	l, exists := r.ListsNamesMap[(*m)["name"]]
	if exists {
		(*m)["names"] = strings.Join(l, ",")
	} else {
		x := []string{}
		for i := 0; i < len(t); i++ {

			x = append(x, fmt.Sprintf("elem%d", i))
		}
		(*m)["names"] = strings.Join(x, ",")
	}
	if strings.Contains(s, "string") {
		(*m)["set_type"] = "String"
	} else if strings.Contains(s, "int") {
		(*m)["set_type"] = "Int"
	} else if strings.Contains(s, "Float") {
		(*m)["set_type"] = "Float"
	}

}

func (r *ResourceTemplateV2) GetConflictingFields(name string) []string {
	i, e := r.ConflictingFields[name]
	if e {
		return i
	}
	return []string{}
}
