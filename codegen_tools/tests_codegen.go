package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"text/template"

	codegen_configs "github.com/vast-data/terraform-provider-vastdata/codegen_tools/configs"
	"github.com/vast-data/terraform-provider-vastdata/utils"
)

func fillArr(i interface{}, dim int) interface{} {
	if dim == 0 {
		return 10
	}
	if reflect.TypeOf(i).Kind() == reflect.Slice {
		i = append(i.([]interface{}), fillArr([]interface{}{}, dim-1))
	}
	return i
}

func GenIntegersList(size, sub_size int) [][]int {
	l := make([][]int, size, size)
	for i := 0; i < size; i++ {
		m := make([]int, sub_size, sub_size)
		for j := 0; j < sub_size; j++ {
			m[j] = j
		}
		l[i] = m
	}
	return l

}

func GenStringsList(size, sub_size int) [][]string {
	l := make([][]string, size, size)
	for i := 0; i < size; i++ {
		m := make([]string, sub_size, sub_size)
		for j := 0; j < sub_size; j++ {
			m[j] = fmt.Sprintf("string-%v", j)
		}
		l[i] = m
	}
	return l

}

func ShouldSkip(resource, field_name string, resources_map map[string]codegen_configs.ResourceTemplateV2) bool {
	r, exists := resources_map[resource]
	if !exists {
		return true
	}
	if r.IgnoreFields.In(field_name) {
		return true
	}
	return false
}

func StructAsFilledMap(t reflect.Type, m *map[string]interface{}, resources_map map[string]codegen_configs.ResourceTemplateV2) {
	for _, fld := range reflect.VisibleFields(t) {
		if ShouldSkip(t.Name(), fld.Name, resources_map) {
			continue
		}
		tag := utils.GetJsonTag(fld.Tag)
		if tag == nil {
			continue
		}
		//we should skip tag=type as this is a very special case where for safty reason type is converted to Type_ <type>  but the tag of the json it type .... this causes some problems when testing
		//it is better to skip it
		if *tag == "type" {
			continue
		}
		if utils.IsPrimitive(fld.Type) {
			switch fld.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				(*m)[*tag] = 100
			case reflect.String:
				(*m)[*tag] = "string"
			case reflect.Float32, reflect.Float64:
				(*m)[*tag] = 10.5
			}
		} else {
			switch fld.Type.Kind() {
			case reflect.Struct, reflect.Pointer:
				n := map[string]interface{}{}
				(*m)[*tag] = n
				StructAsFilledMap(fld.Type.Elem(), &n, resources_map)
			case reflect.Array, reflect.Slice:
				dim, tp := utils.GetListDimentionAndType(fld.Type)
				if (dim >= 1) && (dim <= 2) && utils.IsPrimitive(tp) {
					//we only handle up to 2 dimentional lists of primitives
					if dim == 1 {
						switch tp.Kind() {
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
							(*m)[*tag] = []interface{}{1, 2, 3, 4, 5, 6}
						case reflect.String:
							(*m)[*tag] = []interface{}{"A", "B", "C", "D", "E"}

						}

					} else { //dim must be 2
						o, exists := resources_map[t.Name()]
						if !exists {
							panic(fmt.Sprintf("Unable to find struct data for %s", t.Name()))

						}
						var names_length int
						for _, attr_info := range o.Fields {
							if attr_info.Attributes["name"] == GetTFformatName(fld.Name) {
								names_length = len(strings.Split(attr_info.Attributes["names"], ","))
								break
							}

						}
						switch tp.Kind() {
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
							(*m)[*tag] = GenIntegersList(10, names_length)
						case reflect.String:
							(*m)[*tag] = GenStringsList(10, names_length)

						}

					}

				} else if !utils.IsPrimitive(tp) && dim == 1 {
					a := make([]interface{}, 1, 1)
					n := map[string]interface{}{}
					StructAsFilledMap(tp, &n, resources_map)
					a[0] = n
					(*m)[*tag] = a

				}

			}
		}
	}
}

func getTestJson(resource_name string, mp map[string]codegen_configs.ResourceTemplateV2) string {
	m := map[string]interface{}{}
	u, _ := mp[resource_name]
	StructAsFilledMap(reflect.TypeOf(u.Model), &m, mp)
	x, _ := json.MarshalIndent(m, "", "   ")
	return string(x)
}

func GenResourceTestCode(resource_name string) string {
	resource, exists := resources_templates_map[resource_name]
	var b bytes.Buffer
	if !exists {
		panic(fmt.Sprintf("Unable to find resource withthe name of %s in the resources map", resource_name))
	}

	test_code := `
package resources_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/vast-data/terraform-provider-vastdata/resources"
	utils "github.com/vast-data/terraform-provider-vastdata/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
	"github.com/hashicorp/terraform-plugin-log/tflogtest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var _ = Describe(" {{.RESOURCE.ResourceName}}", func() {
	var ReadContext schema.ReadContextFunc
	var DeleteContext schema.DeleteContextFunc
	var CreateContext schema.CreateContextFunc
	var UpdateContext schema.UpdateContextFunc
	var Importer schema.ResourceImporter
	//	var ResourceSchema map[string]*schema.Schema
	//An empty resource data to be populated per test
	var {{.RESOURCE.ResourceName}}ResourceData *schema.ResourceData
	var model_json = {{ .BT }}
                         {{.TEST_JSON}}
                         {{ .BT }}
	var server *ghttp.Server
	var client vast_client.JwtSession
	 {{.RESOURCE.ResourceName}}Resource := resources.Resource{{.RESOURCE.ResourceName}}()
	ReadContext = {{.RESOURCE.ResourceName}}Resource.ReadContext
	DeleteContext = {{.RESOURCE.ResourceName}}Resource.DeleteContext
	CreateContext = {{.RESOURCE.ResourceName}}Resource.CreateContext
	UpdateContext = {{.RESOURCE.ResourceName}}Resource.UpdateContext
	Importer = *{{.RESOURCE.ResourceName}}Resource.Importer
	//	ResourceSchema = {{.RESOURCE.ResourceName}}Resource.Schema

	BeforeEach(func() {
		{{.RESOURCE.ResourceName}}ResourceData = {{.RESOURCE.ResourceName}}Resource.TestResourceData()
		{{.RESOURCE.ResourceName}}ResourceData.SetId("100")
		server = ghttp.NewTLSServer()
		host_port := strings.Split(server.Addr(), ":")
		host := host_port[0]
		_port := host_port[1]
		port, _ := strconv.ParseUint(_port, 10, 64)
		client = vast_client.NewJwtSession(host, "user", "pwd", port, true)
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/api/token/"),
			ghttp.VerifyJSON("{\"username\":\"user\",\"password\":\"pwd\"}"),
			ghttp.RespondWith(200, {{ .BT }}{"access":"femcew2d332f2e2e322e2qqw#2","":"32dm0932kde,ml;sd,s;l,322332"}{{ .BT }}),
		))

	},
	)
	Describe("Validating Resource Read Context", func() {
		Context("Read Data into a ResourceData", func() {
			It("Resource:{{.RESOURCE.ResourceName}} ,Reads Data", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "{{.RESOURCE.Path}}100"),
					ghttp.RespondWith(200, model_json),
				),
				)
				e := client.Start()
				Expect(e).To(BeNil())
				o := new(map[string]interface{})
				json.Unmarshal([]byte(model_json), o)

				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				d := ReadContext(ctx, {{.RESOURCE.ResourceName}}ResourceData, client)
				Expect(d).To(BeNil())
				attributes := {{.RESOURCE.ResourceName}}ResourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
	Describe("Validating Resource Delete Context", func() {
		Context("Delete A resource", func() {
			It("Resource:{{.RESOURCE.ResourceName}} ,Deletes the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "{{.RESOURCE.Path}}100//"),
					ghttp.RespondWith(200, "DELETED"),
				),
				)
				e := client.Start()
				Expect(e).To(BeNil())
				d := DeleteContext(ctx, {{.RESOURCE.ResourceName}}ResourceData, client)
				Expect(d).To(BeNil())
			})
		})
	},
	)
	Describe("Validating Resource Creation Context", func() {
		Context("Create A resource", func() {
			It("Resource:{{.RESOURCE.ResourceName}} ,Creates the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "{{.RESOURCE.Path}}"),
					ghttp.RespondWith(200, model_json),
				),
				)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "{{.RESOURCE.Path}}0"), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, model_json),
				),
				)

				e := client.Start()
				Expect(e).To(BeNil())
				resource := api_latest.{{.RESOURCE.ResourceName}}{}
				json.Unmarshal([]byte(model_json), &resource)
				resources.Resource{{.RESOURCE.ResourceName}}ReadStructIntoSchema(ctx, resource, {{.RESOURCE.ResourceName}}ResourceData)
				{{.RESOURCE.ResourceName}}ResourceData.SetId("100")
				d := CreateContext(ctx, {{.RESOURCE.ResourceName}}ResourceData, client)
				Expect(d).To(BeNil())

			})
		})
	},
	)
	Describe("Validating Resource Update Context", func() {
		Context("Update A resource", func() {
			It("Resource:{{.RESOURCE.ResourceName}} ,Update the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var new_guid = "11111-11111-11111-11111-11111"

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "{{.RESOURCE.Path}}"),
					ghttp.RespondWith(200, model_json),
				),
				)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "{{.RESOURCE.Path}}0"), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, model_json),
				),
				)

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("PATCH", "{{.RESOURCE.Path}}/0/"), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, model_json),
				),
				)

				e := client.Start()
				Expect(e).To(BeNil())
				//First we create a resource than we change it and see if it was updated
				resource := api_latest.{{.RESOURCE.ResourceName}}{}
				json.Unmarshal([]byte(model_json), &resource)
				resources.Resource{{.RESOURCE.ResourceName}}ReadStructIntoSchema(ctx, resource, {{.RESOURCE.ResourceName}}ResourceData)
				d := CreateContext(ctx, {{.RESOURCE.ResourceName}}ResourceData, client)
				Expect(d).To(BeNil())
				//We update the guid as this is a fieled that always exists
				resource.Guid = new_guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "{{.RESOURCE.Path}}0"), //the new_guid is returned and it should change the value of the resource
					ghttp.RespondWith(200, string(b)),
				),
				)

				d = UpdateContext(ctx, {{.RESOURCE.ResourceName}}ResourceData, client)
				Expect(d).To(BeNil())
				Expect({{.RESOURCE.ResourceName}}ResourceData.Get("guid")).To(Equal(new_guid))

			})
		})
	},
	)
	Describe("Validating Resource Importer", func() {
		Context("Import A resource", func() {
			It("Resource:{{.RESOURCE.ResourceName}} ,Imports the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var guid = "11111-11111-11111-11111-11111"
                            
				resource := api_latest.{{.RESOURCE.ResourceName}}{}
				json.Unmarshal([]byte(model_json), &resource)
				{{.RESOURCE.ResourceName}}ResourceData.SetId(guid)
				resource.Guid = guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())
                                {{ if  .RESOURCE.ResponseGetByURL }}
                                request_url := {{.BT}}[{"url":"https://{{.BT}} + server.Addr() + {{.BT}}{{.RESOURCE.Path}}100"}]{{.BT}}
                                server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "{{.RESOURCE.Path}}", fmt.Sprintf("guid=%s", guid)), 
					ghttp.RespondWith(200, request_url),
				),
				)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "{{.RESOURCE.Path}}100"),
					ghttp.RespondWith(200, string(b)),
				),
				)
                                
                                {{else}}


                                
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "{{.RESOURCE.Path}}", fmt.Sprintf("guid=%s", guid)), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, {{ .BT }}[{{ .BT }}+string(b)+{{ .BT }}]{{ .BT }}),
				),
				)
                                
                                {{end}}
				ctx = tflogtest.RootLogger(ctx, &output)
				e := client.Start()
				Expect(e).To(BeNil())
				Importer.StateContext(ctx, {{.RESOURCE.ResourceName}}ResourceData, client)
				//				Expect(d).To(BeNil())

				o := new(map[string]interface{})
				json.Unmarshal(b, o)
				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				attributes := {{.RESOURCE.ResourceName}}ResourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
})
`

	m := map[string]interface{}{"BT": "`", "TEST_JSON": getTestJson(resource_name, resources_templates_map), "RESOURCE": resource}

	t := template.Must(template.New("test_code").Parse(test_code))
	err := t.Execute(&b, m)
	if err != nil {
		fmt.Println(err)
	}
	return b.String()

}

func GenDataSourceTestCode(datasource_name string) string {
	resource, exists := datasources_templates_map[datasource_name]
	var b bytes.Buffer
	var funcMap template.FuncMap = template.FuncMap{}

	if !exists {
		panic(fmt.Sprintf("Unable to find a datasource with the name of %s in the resources map", datasource_name))
	}

	test_code := `
package datasources_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
        "net/url"

	"github.com/vast-data/terraform-provider-vastdata/datasources"
	utils "github.com/vast-data/terraform-provider-vastdata/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
	"github.com/hashicorp/terraform-plugin-log/tflogtest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var _ = Describe(" {{.DATASOURCE.ResourceName}}", func() {
	var ReadContext schema.ReadContextFunc
	//An empty resource data to be populated per test
	var {{.DATASOURCE.ResourceName}}DataSourceData *schema.ResourceData
	var model_json = {{ .BT }}
                         {{.TEST_JSON}}
                         {{ .BT }}
	var server *ghttp.Server
	var client vast_client.JwtSession
	{{.DATASOURCE.ResourceName}}DataSource := datasources.DataSource{{.DATASOURCE.ResourceName}}()
	ReadContext = {{.DATASOURCE.ResourceName}}DataSource.ReadContext

	BeforeEach(func() {
		{{.DATASOURCE.ResourceName}}DataSourceData = {{.DATASOURCE.ResourceName}}DataSource.TestResourceData()
		{{.DATASOURCE.ResourceName}}DataSourceData.SetId("100")
		server = ghttp.NewTLSServer()
		host_port := strings.Split(server.Addr(), ":")
		host := host_port[0]
		_port := host_port[1]
		port, _ := strconv.ParseUint(_port, 10, 64)
		client = vast_client.NewJwtSession(host, "user", "pwd", port, true)
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/api/token/"),
			ghttp.VerifyJSON("{\"username\":\"user\",\"password\":\"pwd\"}"),
			ghttp.RespondWith(200, {{ .BT }}{"access":"femcew2d332f2e2e322e2qqw#2","":"32dm0932kde,ml;sd,s;l,322332"}{{ .BT }}),
		))

	},
	)
	Describe("Validating Datasource Read", func() {
		Context("Read A datasource", func() {
			It("Datasource:{{.DATASOURCE.ResourceName}} ,reads", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var guid = "11111-11111-11111-11111-11111"
                            
				resource := api_latest.{{.DATASOURCE.ResourceName}}{}
				json.Unmarshal([]byte(model_json), &resource)
				{{.DATASOURCE.ResourceName}}DataSourceData.SetId(guid)
				resource.Guid = guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())
                                values := url.Values{}
                                {{ $ds:=.DATASOURCE}}
                                {{ range $i,$v := .DATASOURCE.RequiredIdentifierFields.ToArray }}
                                values.Add("{{$v}}", fmt.Sprintf("%v" ,resource.{{ ConvertTfNameToModelName $v}}))
                                {{$ds.ResourceName}}DataSourceData.Set("{{$v}}" ,resource.{{ ConvertTfNameToModelName $v}})
                                {{end}}
                                {{ if  .DATASOURCE.ResponseGetByURL }}
                                request_url := {{.BT}}[{"url":"https://{{.BT}} + server.Addr() + {{.BT}}{{.DATASOURCE.Path}}100"}]{{.BT}}
                                server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "{{.DATASOURCE.Path}}", values.Encode()), 
					ghttp.RespondWith(200, request_url),
				),
				)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "{{.DATASOURCE.Path}}100"),
					ghttp.RespondWith(200, string(b)),
				),
				)

                                {{else}}

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "{{.DATASOURCE.Path}}", values.Encode()), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, {{ .BT }}[{{ .BT }}+string(b)+{{ .BT }}]{{ .BT }}),
				),
				)
                               
                                {{end}}
				ctx = tflogtest.RootLogger(ctx, &output)
				e := client.Start()
				Expect(e).To(BeNil())
				d:=ReadContext(ctx, {{.DATASOURCE.ResourceName}}DataSourceData, client)
				Expect(d).To(BeNil())

				o := new(map[string]interface{})
				json.Unmarshal(b, o)
				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				attributes := {{.DATASOURCE.ResourceName}}DataSourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
})

`

	funcMap["ConvertTfNameToModelName"] = resource.ConvertTfNameToModelName
	m := map[string]interface{}{"BT": "`", "TEST_JSON": getTestJson(datasource_name, datasources_templates_map), "DATASOURCE": resource}

	t := template.Must(template.New("test_code").Funcs(funcMap).Parse(test_code))
	err := t.Execute(&b, m)
	if err != nil {
		fmt.Println(err)
	}
	return b.String()

}
