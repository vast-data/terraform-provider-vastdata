package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"

	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/pb33f/libopenapi"
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

var re *regexp.Regexp = regexp.MustCompile("\\[\\]")

var datasources_templates_map map[string]ResourceTemplateV2

var resources_templates_map map[string]ResourceTemplateV2

var spingFuncMap template.FuncMap = sprig.FuncMap()

var funcMap template.FuncMap = template.FuncMap{
	"upper":           strings.ToUpper,
	"split":           strings.Split,
	"ToListOfStrings": ToListOfStrings,
	"AddInt":          AddInt,
	"indent":          spingFuncMap["indent"],
	"replace":         strings.Replace,
	"replaceAll":      strings.ReplaceAll,
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

type ResourceTemplateV2 struct {
	ResourceName             string
	Fields                   []ResourceElem
	Path                     *string
	Model                    interface{}
	DestFile                 *string
	IgnoreFields             *StringSet
	RequiredIdentifierFields *StringSet
	OptionalIdentifierFields *StringSet
	ListsNamesMap            map[string][]string
	Generate                 bool
	DataSourceName           string
	ResponseProcessingFunc   string
	ResponseGetByURL         bool
	IgnoreUpdates            *StringSet
	TfNameToModelName        map[string]string
}

func (r *ResourceTemplateV2) ConvertTfNameToModelName(tf_name string) string {
	name, exists := r.TfNameToModelName[tf_name]
	if exists {
		return name
	}
	return ""
}

func (r *ResourceTemplateV2) X() []string {
	return []string{}
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

func ProcessResourceTemplate(R *ResourceTemplateV2) {
	var elem_type string
	TfNameToModelName := map[string]string{}
	Fields := []ResourceElem{}
	r := R.Model
	t := reflect.TypeOf(r)
	for _, e := range reflect.VisibleFields(t) {
		if R.IgnoreFields.In(e.Name) {
			continue
		}
		m := map[string]string{"name": GetTFformatName(e.Name), "modelName": e.Name}
		TfNameToModelName[GetTFformatName(e.Name)] = e.Name
		m["max_items"] = "0"
		if R.IgnoreUpdates == nil {
			R.IgnoreUpdates = NewStringSet()
		}
		if R.IgnoreUpdates.In(m["name"]) {
			m["ignore_update"] = "true"
		} else {
			m["ignore_update"] = "false"
		}

		if R.RequiredIdentifierFields.In(m["name"]) {
			m["computed"] = "false"
			m["required"] = "true"
			m["optional"] = "false"
		} else if R.OptionalIdentifierFields.In(m["name"]) {
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
			R.SetupListProperties(elem_type, &m)
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
		Fields = append(Fields, NewResourceElem(m))

	}
	R.Fields = Fields
	R.TfNameToModelName = TfNameToModelName

}

type ResourceElem struct {
	Attributes   map[string]string
	ResourceElem *ResourceElem
	Indent       int
}

func (r *ResourceElem) IsReferance() bool {
	if r.ResourceElem != nil {
		return true
	}
	return false
}

func NewResourceElem(m map[string]string) ResourceElem {
	return ResourceElem{
		Attributes:   m,
		ResourceElem: nil,
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
       version_{{ replaceAll $k "." "_" }} "github.com/vast-data/terraform-provider-vastdata.git/codegen/{{$k}}"
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
	gen_datasources()
	gen_resources()
	BuildVersionsRefs()
	build_resources_tests()
	build_datasources_tests()

}
