package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"runtime"
	"sort"

	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pb33f/libopenapi"
	"github.com/vast-data/terraform-provider-vastdata/utils"

	base "github.com/pb33f/libopenapi/datamodel/high/base"
)

func ToListOfStrings(s string) string {
	o := []string{}
	for _, t := range strings.Split(s, ",") {
		o = append(o, fmt.Sprintf("\"%s\"", t))
	}
	return strings.Join(o, ",")
}

func AddInt(x, y int) int {
	return x + y
}

func GetTypeName(i interface{}) string {
	return reflect.TypeOf(i).Name()

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

var apiSchemasMap map[string]*base.SchemaProxy = make(map[string]*base.SchemaProxy)

var re *regexp.Regexp = regexp.MustCompile("\\[\\]")

var datasources_templates_map map[string]ResourceTemplateV2

var resources_templates_map map[string]ResourceTemplateV2

var spingFuncMap template.FuncMap = sprig.FuncMap()

func GetBT() string {
	return "`"
}

var funcMap template.FuncMap = template.FuncMap{
	"upper":           strings.ToUpper,
	"split":           strings.Split,
	"ToListOfStrings": ToListOfStrings,
	"AddInt":          AddInt,
	"indent":          spingFuncMap["indent"],
	"replace":         strings.Replace,
	"replaceAll":      strings.ReplaceAll,
	"getBT":           GetBT,
	"getTypeName":     GetTypeName,
	"funcName":        FuncName,
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

func AsStringListDefentions(s []string) string {
	a := "[]string{"
	for _, i := range s {
		a += "\"" + i + "\","
	}
	a += "}"
	return a
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
	ComputedFields           *StringSet
	ListsNamesMap            map[string][]string
	Generate                 bool
	DataSourceName           string
	ResponseProcessingFunc   string
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
	FieldsValidators         map[string]schema.SchemaValidateDiagFunc
	SensitiveFields          *StringSet
	IsDataSource             bool
	BeforeCreateFunc         utils.ResponseConversionFunc
}

func (r *ResourceTemplateV2) HasProperty(property string) bool {
	if r.ApiSchema != nil {
		schema := r.ApiSchema.Schema()
		_, exists := schema.Properties[property]
		return exists
	}
	return false
}

func (r *ResourceTemplateV2) GetEnum(property string) []interface{} {
	if r.HasProperty(property) {
		schema := r.ApiSchema.Schema()
		property_schema := schema.Properties[property].Schema()
		if property_schema == nil {
			return []interface{}{}
		}
		return property_schema.Enum

	}
	return []interface{}{}
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
		return r.ApiSchema.Schema().Description
	}
	return ""
}

func (r *ResourceTemplateV2) GetSchemaProperyDocument(property string) string {
	if r.ApiSchema != nil {
		schema := r.ApiSchema.Schema()
		property_schema, exists := schema.Properties[property]
		if !exists {
			return ""
		}
		enum := r.GetEnum(property)
		out := property_schema.Schema().Description
		if len(enum) > 0 {
			out += fmt.Sprintf(" Allowed Values are %v", ListAsStringsList(enum))
		}
		return out
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

func GetTFformatName(s string) string {
	re := regexp.MustCompile("([A-Z])")
	t := re.ReplaceAllString(s, "_${1}")
	return strings.ToLower(strings.TrimLeft(t, "_"))
}
func LoadApi(filename string) (libopenapi.Document, error) {
	_api, _ := os.ReadFile(filename)
	doc, err := libopenapi.NewDocument(_api)
	return doc, err

}

func aPISchemasPopulate(doc libopenapi.Document) {
	v3, _ := doc.BuildV3Model()
	//All is converted to lower to preven conversion names with swagger
	for n, s := range v3.Model.Components.Schemas {
		apiSchemasMap[strings.ToLower(n)] = s
	}

}
func getSchemaProxy(name string) *base.SchemaProxy {
	name = strings.ToLower(name)
	p, exists := apiSchemasMap[name]
	if exists {
		return p
	}
	return nil
}

func ProcessResourceTemplate(R *ResourceTemplateV2) {
	var elem_type string
	TfNameToModelName := map[string]string{}
	Fields := []ResourceElem{}
	r := R.Model
	R.ApiSchema = getSchemaProxy(strings.ToLower(R.ResourceName))
	t := reflect.TypeOf(r)
	for _, e := range reflect.VisibleFields(t) {
		if R.IgnoreFields.In(e.Name) {
			continue
		}
		m := map[string]string{"name": GetTFformatName(e.Name), "modelName": e.Name}
		TfNameToModelName[GetTFformatName(e.Name)] = e.Name
		m["validator_func"] = ""
		if R.HasValidatorFunc(m["name"]) {
			m["validator_func"] = R.GetValidatorFunc(m["name"])
		}
		m["enum"] = ""
		enum := R.GetEnum(m["name"])
		if len(enum) > 0 && (!R.HasValidatorFunc(m["name"])) {
			fmt.Println(m["name"])
			l := ListAsStringsList(enum)
			m["enum"] = AsStringListDefentions(l)
		}
		m["max_items"] = "0"
		if R.IgnoreUpdates == nil {
			R.IgnoreUpdates = NewStringSet()
		}
		if R.ComputedFields == nil {
			R.ComputedFields = NewStringSet()
		}
		if R.IgnoreUpdates.In(m["name"]) {
			m["ignore_update"] = "true"
		} else {
			m["ignore_update"] = "false"
		}
		m["sensitive"] = "false"
		if R.SensitiveFields != nil && R.SensitiveFields.In(m["name"]) {
			m["sensitive"] = "true"
		}

		if R.RequiredIdentifierFields.In(m["name"]) {
			m["computed"] = "false"
			m["required"] = "true"
			m["optional"] = "false"
		} else if R.OptionalIdentifierFields.In(m["name"]) {
			m["computed"] = "true"
			m["required"] = "false"
			m["optional"] = "true"
		} else if R.ComputedFields.In(m["name"]) {
			m["computed"] = "true"
			m["required"] = "false"
			m["optional"] = "false"
		} else if m["name"] == "guid" {
			m["computed"] = "true"
			m["required"] = "false"
			m["optional"] = "false"
		} else if !R.IsDataSource {
			m["computed"] = "true"
			m["required"] = "false"
			m["optional"] = "true"
		} else {
			m["computed"] = "true"
			m["required"] = "false"
			m["optional"] = "false"
		}
		elem_type = e.Type.String()
		m["elem_type"] = elem_type
		switch elem_type {
		case "int", "int32", "int64", "int8", "int16":
			m["type"] = "TypeInt"
		case "float", "float32", "float64":
			m["type"] = "TypeFloat"
		case "string":
			m["type"] = "TypeString"
		case "bool":
			m["type"] = "TypeBool"
		case "[]string", "[]int", "[]int32", "[]int64", "[]float", "[]float32", "[]float64", "[][]string", "[][]int", "[][]int32", "[][]int64", "[][]float", "[][]float32", "[][]float64":
			if R.HasFakeField(m["name"]) {
				R.SetupFakeFields(elem_type, R.ListFields[m["name"]], &m)
			} else {
				R.SetupListProperties(elem_type, &m)
			}
		default:
			//If we got here than we are not dealing with primitive / primitives array
			m["type"] = "TypeList"
			m["list_type"] = "complex"
			m["set_access"] = "Object"

			if e.Type.Kind() == reflect.Slice {
				m["set_type"] = e.Type.Elem().Name()
				m["set_access"] = "List"
			} else if e.Type.Kind() == reflect.Pointer {
				m["set_type"] = e.Type.Elem().Name()
				m["set_access"] = "Pointer"

			}
			m["set_type"] = e.Type.Elem().Name()
		}
		Fields = append(Fields, NewResourceElem(m, R))

	}
	R.Fields = Fields
	R.TfNameToModelName = TfNameToModelName
	R.ApiSchema = getSchemaProxy(strings.ToLower(R.ResourceName))
	R.ResourceDocumantation = R.GetSchemaDocumentation()
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

/*
Return a map of resource templates based on a list
The keys will be the names of the Resource Names
*/
func TemplatesListToMap(datasources_templates []ResourceTemplateV2) map[string]ResourceTemplateV2 {
	m := make(map[string]ResourceTemplateV2)
	for _, r := range datasources_templates {
		m[r.ResourceName] = r
	}
	return m
}

func ToStringPointer(s string) *string {
	return &s
}

func WriteDataSourceCodeCodeFile(base_path string, r ResourceTemplateV2) {
	path := filepath.Join(base_path, *r.DestFile)
	f, _ := os.Create(path)
	defer f.Close()
	f.WriteString(GenDataSourceTemplate(r))
}

func WriteResourceCodeCodeFile(base_path string, r ResourceTemplateV2) {
	path := filepath.Join(base_path, *r.DestFile)
	f, _ := os.Create(path)
	defer f.Close()
	f.WriteString(GenResourceTemplate(r))
}

func WriteStringToFile(base_path, filename, content string) {
	path := filepath.Join(base_path, filename)
	f, _ := os.Create(path)
	defer f.Close()
	f.WriteString(content)
}
func gen_datasources() {
	base_path := "../datasources/"
	for i, _ := range datasources_templates {
		d := &datasources_templates[i]
		d.IsDataSource = true
		ProcessResourceTemplate(&datasources_templates[i])

	}
	datasources_templates_map = TemplatesListToMap(datasources_templates)
	for _, resource_template := range datasources_templates {
		if resource_template.Generate {
			WriteDataSourceCodeCodeFile(base_path, resource_template)
		}
	}
	WriteStringToFile(base_path, "datasources.go", BuildDataSourcesList(datasources_templates))
}

func gen_resources() {
	base_path := "../resources/"
	for i, _ := range resources_templates {
		d := &resources_templates[i]
		d.IsDataSource = false
		ProcessResourceTemplate(&resources_templates[i])

	}
	resources_templates_map = TemplatesListToMap(resources_templates)
	for _, resource_template := range resources_templates {
		if resource_template.Generate {
			WriteResourceCodeCodeFile(base_path, resource_template)
		}
	}
	WriteStringToFile(base_path, "resources.go", BuildResourcesList(resources_templates))
}

func scanPackageForStructs(package_path string) []string {
	//This function will go over a go package and list all structs which have been defined by it
	structs := []string{}
	t := token.NewFileSet()
	pkgs, _ := parser.ParseDir(t, package_path, nil, parser.ParseComments)
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				if tDecl, ok := decl.(*ast.GenDecl); ok && tDecl.Tok == token.TYPE {
					for _, spec := range tDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							if _, ok := typeSpec.Type.(*ast.StructType); ok {
								structs = append(structs, typeSpec.Name.Name)
							}
						}

					}
				}

			}

		}
	}
	sort.Strings(structs) // Sort in order to keep ordered for later usage of file creation
	return structs
}

func MakeNillPointersToTypes(struct_names map[string][]string) string {
	var b bytes.Buffer
	tmp := `
package vast_versions 

import (
      "reflect"
       {{ range $k,$v :=  . -}}
       version_{{ replaceAll $k "." "_" }} "github.com/vast-data/terraform-provider-vastdata/codegen/{{$k}}"
       {{ end }}
      
)

var vast_versions map[string]map[string]reflect.Type = map[string]map[string]reflect.Type{
        {{ range $k,$v :=  . -}}
           "{{$k}}":map[string]reflect.Type{
            {{ $ver:= replaceAll $k "." "_" }}
            {{- range $i := $v -}}
           "{{$i}}":reflect.TypeOf((*version_{{$ver}}.{{$i}})(nil)).Elem(),
           {{ end -}}
          },
        {{ end -}}
      }
`
	t := template.Must(template.New("pointers_map").Funcs(funcMap).Parse(tmp))
	err := t.Execute(&b, struct_names)
	if err != nil {
		fmt.Println(err)
	}
	return b.String()
}

func BuildVersionsRefs() {
	versions := os.Getenv("BUILD_VERSIONS")
	m := map[string][]string{}
	if versions != "" {
		for _, ver := range strings.Split(versions, " ") {
			m[ver] = scanPackageForStructs(fmt.Sprintf("../codegen/%s", ver))

		}

	}
	WriteStringToFile("../vast_versions/", "versions_map.go", MakeNillPointersToTypes(m))

}

func build_resources_tests() {
	for resource_name, resource := range resources_templates_map {
		if resource.Generate {
			test_code := GenResourceTestCode(resource_name)
			filename := *resource.DestFile
			filename = strings.Replace(filename, ".go", "_test.go", 1)
			fmt.Printf("Writing test file %s for resource %s\n", filename, resource_name)
			WriteStringToFile("../resources", filename, test_code)
		}

	}

}

func build_datasources_tests() {
	for datasource_name, datasource := range datasources_templates_map {
		if datasource.Generate {
			test_code := GenDataSourceTestCode(datasource_name)
			filename := *datasource.DestFile
			filename = strings.Replace(filename, ".go", "_test.go", 1)
			fmt.Printf("Writing test file %s for datasource %s\n", filename, datasource_name)
			WriteStringToFile("../datasources", filename, test_code)
		}

	}

}

func main() {
	doc, err := LoadApi("../codegen/latest/api.yaml")
	if err != nil {
		panic(err.Error())
	}
	aPISchemasPopulate(doc)
	gen_datasources()
	gen_resources()
	BuildVersionsRefs()
	build_resources_tests()
	build_datasources_tests()

}
