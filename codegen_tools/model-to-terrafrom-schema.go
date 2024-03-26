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
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/pb33f/libopenapi"

	base "github.com/pb33f/libopenapi/datamodel/high/base"
	codegen_configs "github.com/vast-data/terraform-provider-vastdata/codegen_tools/configs"
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

var apiSchemasMap map[string]*base.SchemaProxy = make(map[string]*base.SchemaProxy)

var re *regexp.Regexp = regexp.MustCompile("\\[\\]")

var datasources_templates_map map[string]codegen_configs.ResourceTemplateV2

var resources_templates_map map[string]codegen_configs.ResourceTemplateV2

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
	"funcName":        codegen_configs.FuncName,
}

func AsStringListDefentions(s []string) string {
	a := "[]string{"
	for _, i := range s {
		a += "\"" + i + "\","
	}
	a += "}"
	return a
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

func ProcessResourceTemplate(R *codegen_configs.ResourceTemplateV2) {
	var elem_type string
	TfNameToModelName := map[string]string{}
	Fields := []codegen_configs.ResourceElem{}
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
			l := codegen_configs.ListAsStringsList(enum)
			m["enum"] = AsStringListDefentions(l)
		}
		m["max_items"] = "0"
		if R.IgnoreUpdates == nil {
			R.IgnoreUpdates = codegen_configs.NewStringSet()
		}
		if R.ForceNewFields == nil {
			R.ForceNewFields = codegen_configs.NewStringSet()
		}
		if R.ComputedFields == nil {
			R.ComputedFields = codegen_configs.NewStringSet()
		}
		if R.ForceNewFields.In(m["name"]) {
			m["force_new"] = "true"
		} else {
			m["force_new"] = "false"
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
		} else if R.GetSchemaProperyDefault(m["name"]) != "" {
			m["required"] = "false"
			m["computed"] = "false"
			m["optional"] = "true"
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
		Fields = append(Fields, codegen_configs.NewResourceElem(m, R))

	}
	R.Fields = Fields
	R.TfNameToModelName = TfNameToModelName
	R.ApiSchema = getSchemaProxy(strings.ToLower(R.ResourceName))
	R.ResourceDocumantation = R.GetSchemaDocumentation()
	R.SetFunctions()
}

/*
Return a map of resource templates based on a list
The keys will be the names of the Resource Names
*/
func TemplatesListToMap(datasources_templates []codegen_configs.ResourceTemplateV2) map[string]codegen_configs.ResourceTemplateV2 {
	m := make(map[string]codegen_configs.ResourceTemplateV2)
	for _, r := range datasources_templates {
		m[r.ResourceName] = r
	}
	return m
}

func WriteDataSourceCodeCodeFile(base_path string, r codegen_configs.ResourceTemplateV2) {
	path := filepath.Join(base_path, *r.DestFile)
	f, _ := os.Create(path)
	defer f.Close()
	f.WriteString(GenDataSourceTemplate(r))
}

func WriteResourceCodeCodeFile(base_path string, r codegen_configs.ResourceTemplateV2) {
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
	for i, _ := range codegen_configs.DatasourcesTemplates {
		d := &codegen_configs.DatasourcesTemplates[i]
		d.IsDataSource = true
		ProcessResourceTemplate(&codegen_configs.DatasourcesTemplates[i])

	}
	datasources_templates_map = TemplatesListToMap(codegen_configs.DatasourcesTemplates)
	for _, resource_template := range codegen_configs.DatasourcesTemplates {
		if resource_template.Generate {
			WriteDataSourceCodeCodeFile(base_path, resource_template)
		}
	}
	WriteStringToFile(base_path, "datasources.go", BuildDataSourcesList(codegen_configs.DatasourcesTemplates))
}

func gen_resources() {
	base_path := "../resources/"
	for i, _ := range codegen_configs.ResourcesTemplates {
		d := &codegen_configs.ResourcesTemplates[i]
		d.IsDataSource = false
		ProcessResourceTemplate(&codegen_configs.ResourcesTemplates[i])

	}
	resources_templates_map = TemplatesListToMap(codegen_configs.ResourcesTemplates)
	for _, resource_template := range codegen_configs.ResourcesTemplates {
		if resource_template.Generate {
			WriteResourceCodeCodeFile(base_path, resource_template)
		}
	}
	WriteStringToFile(base_path, "resources.go", BuildResourcesList(codegen_configs.ResourcesTemplates))
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

func gen_import_command(path string) {
	for _, r := range codegen_configs.ResourcesTemplates {
		if !r.DisableImport && r.Generate {
			fmt.Printf("Creating import file for Resource: %v ,Resource Name:%v Doc:%v\n", r.ResourceName, r.DataSourceName, r.Importer.GetDoc())
			base_path := filepath.Join(path, r.DataSourceName)
			import_string := ""
			for _, d := range r.Importer.GetDoc() {
				import_string += fmt.Sprintf("terraform import %v.example %v\n", r.DataSourceName, d)
			}
			WriteStringToFile(base_path, "import.sh", import_string)
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
	gen_import_command("../examples/resources/")

}
