package main

import (
	"bytes"
	"fmt"
	"text/template"

	codegen_configs "github.com/vast-data/terraform-provider-vastdata/codegen_tools/configs"
)

func BuildDataSourceTemplateHeader(r codegen_configs.ResourceTemplateV2) string {
	var b bytes.Buffer
	header :=
		`package datasources

import (
        "encoding/json"
        "fmt"
        "context"
        "strconv"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
        "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
        vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
        api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
        "github.com/hashicorp/terraform-plugin-log/tflog"
        "net/url"
        utils "github.com/vast-data/terraform-provider-vastdata/utils"
       
        
)

func DataSource{{ .ResourceName }}() *schema.Resource {
     return  &schema.Resource{
     ReadContext: dataSource{{ .ResourceName }}Read,
     Description: {{getBT}}{{ .ResourceDocumantation }}{{getBT}},
     Schema: map[string]*schema.Schema{
`

	t := template.Must(template.New("header").Funcs(funcMap).Parse(header))
	err := t.Execute(&b, r)
	if err != nil {
		fmt.Println(err)
	}
	return b.String()
}

func ResourceTemplateToTerrafromElem(r codegen_configs.ResourceElem, indent int) string {
	var b bytes.Buffer
	r.Indent = indent
	tmpl := `     
             {{ $I:=.Indent}}
             {{ $name:=.Attributes.name}}
	     {{indent $I " "}}"{{ .Attributes.name }}": &schema.Schema{
	     {{indent $I " "}}   Type: 	  schema.{{ .Attributes.type }},
	     {{indent $I " "}}   Computed: {{ .Attributes.computed }},
             {{indent $I " "}}   Required: {{ .Attributes.required }},
             {{indent $I " "}}   Optional: {{ .Attributes.optional }},
             {{indent $I " "}}   Description: {{getBT}}{{ GetSchemaProperyDocument .Attributes.name }}{{getBT}},
             {{ if and (eq .Attributes.length "1") (eq .Attributes.list_type "simple") (eq .Attributes.type "TypeList")}}
                {{indent $I " "}}Elem: &schema.Schema{
                {{indent $I " "}}            Type: schema.Type{{.Attributes.set_type}},                               
                                   },              
             {{ end }}
             {{ if and (eq .Attributes.type "TypeList") (ne .Attributes.length "1") }}
                {{ $f:=.Attributes }}
                {{indent $I " "}}Elem: &schema.Resource{
                {{indent $I " "}}    Schema: map[string]*schema.Schema{
                {{ if or (eq .Attributes.set_type "Int") (eq .Attributes.set_type "String") (eq .Attributes.set_type "Float") }}
                {{ range $t:= split .Attributes.names "," }}
                {{indent $I " "}}      "{{$t}}": &schema.Schema{
                {{indent $I " "}}               Type:  schema.Type{{ $f.set_type }},
                {{indent $I " "}}               Computed:  true,
                {{indent $I " "}}               Description: "{{GetFakeFieldDescription $name $t}}",
                {{indent $I " "}}         },
                {{end}}
                {{else}} 
                {{ BuildTemplateFromModelName .Attributes.set_type ( AddInt $I 7) }}{{end}}
              {{indent $I " "}}},
             {{indent $I " "}}},{{ end }}
	    {{indent $I " "}}},
`
	//A dirty workaround to avoid loop decleration
	_, exists := funcMap["BuildTemplateFromModelName"]
	if !exists {
		funcMap["BuildTemplateFromModelName"] = BuildTemplateFromModelName
	}
	localFuncMap := template.FuncMap{}

	for k, v := range funcMap {
		localFuncMap[k] = v
	}
	localFuncMap["GetFakeFieldDescription"] = r.Parent.GetFakeFieldDescription
	localFuncMap["GetSchemaProperyDocument"] = r.Parent.GetSchemaProperyDocument

	t := template.Must(template.New("tf").Funcs(localFuncMap).Parse(tmpl))

	err := t.Execute(&b, r)
	if err != nil {
		fmt.Println(err)
	}
	return b.String()

}

func BuildTemplate(r codegen_configs.ResourceTemplateV2, indent int) string {
	out := ""
	for _, f := range r.Fields {
		out += ResourceTemplateToTerrafromElem(f, indent)
	}
	return out
}

func BuildTemplateFromModelName(n string, indent int) string {
	model, _ := datasources_templates_map[n]
	return BuildTemplate(model, indent)
}

func BuildDataSourcesList(datasources_templates []codegen_configs.ResourceTemplateV2) string {
	var b bytes.Buffer
	datasources :=
		`package datasources

import (
 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
        
)

var DataSources map[string]*schema.Resource = map[string]*schema.Resource{
{{range . }}
{{- if .Generate }}     "{{.DataSourceName}}":DataSource{{ .ResourceName -}}(),
{{ end -}}
{{end -}}
}
`

	t := template.Must(template.New("datasources").Parse(datasources))
	err := t.Execute(&b, datasources_templates)
	if err != nil {
		fmt.Println(err)
	}
	return b.String()

}

func BuildDataSourceTemplateReadFunction(r codegen_configs.ResourceTemplateV2) string {
	var b bytes.Buffer
	read_function := `
func dataSource{{ .ResourceName }}Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
     var diags diag.Diagnostics
     {{ $cbr:="}" }}
     {{ $cbl:="{" }}
     client:=m.(vast_client.JwtSession)
     values := url.Values{}
     {{ range $i,$v := .RequiredIdentifierFields.ToArray }} 
     {{$v}}:=d.Get("{{$v}}") 
     values.Add("{{$v}}",fmt.Sprintf("%v",{{$v}}))
     {{ end }}
     {{ range $i,$v := .OptionalIdentifierFields.ToArray }} 
     if d.HasChanges("{{$v}}") {
         {{$v}}:=d.Get("{{$v}}")
         tflog.Debug(ctx,"Using optional attribute \"{{$v}}\"")
         values.Add("{{$v}}",fmt.Sprintf("%v",{{$v}}))
     }
     {{ end }}
     response,err:=client.Get(ctx,utils.GenPath("{{.Path}}"),values.Encode(), map[string]string{})
     tflog.Info(ctx,response.Request.URL.String())
     if err!=nil {
        diags = append(diags, diag.Diagnostic {
		Severity: diag.Error,
		Summary:  "Error occured while obtaining data from the vastdata cluster",
		Detail:   err.Error(),
		})
       return diags

     }
     resource_l:=[]api_latest.{{.ResourceName}}{}
     {{ if ne .ResponseProcessingFunc "" }}
     body,err:=utils.{{.ResponseProcessingFunc}}(ctx,response)
     {{else }}
     body,err:=utils.DefaultProcessingFunc(ctx,response)
     {{end -}}
     if err!=nil {
         diags = append(diags, diag.Diagnostic {
		Severity: diag.Error,
		Summary:  "Error occured reading data recived from VastData cluster",
		Detail:   err.Error(),
		})
       return diags

     }
     {{ if .ResponseGetByURL }} 
     body,err = utils.ResponseGetByURL(ctx,body,client)
     if err!=nil {
         diags = append(diags, diag.Diagnostic {
		Severity: diag.Error,
		Summary:  "Error occured reading urls from response",
		Detail:   err.Error(),
		})
       return diags
     }
     {{end -}}
     err=json.Unmarshal(body,&resource_l)
     if err!=nil {
                diags = append(diags, diag.Diagnostic {
		Severity: diag.Error,
		Summary:  "Error occured while parsing data recived from VastData cluster",
		Detail:   err.Error(),
		})
       return diags

     }
     if len(resource_l) == 0 {
         d.SetId("")
         diags = append(diags, diag.Diagnostic {
		Severity: diag.Error,
		Summary:  "Could not find a resource that matches those attributes",
		Detail:   "Could not find a resource that matches those attributes",
		})
         return diags         
     }
     if len(resource_l) > 1 {
         d.SetId("")
         diags = append(diags, diag.Diagnostic {
		Severity: diag.Error,
		Summary:  "Multiple results returned, you might want to add more attributes to get a specific resource",
		Detail:   "Multiple results returned, you might want to add more attributes to get a specific resource",
		})
         return diags         
     }

     resource:=resource_l[0]

     {{range .Fields}}
     tflog.Info(ctx,fmt.Sprintf("%v - %v","{{.Attributes.modelName}}",resource.{{.Attributes.modelName}}))
     {{if eq .Attributes.type "TypeList"}}
     {{ if and (eq .Attributes.type "TypeList") (eq .Attributes.set_type "String") -}}
     {{ if eq .Attributes.length "1" }}
     err=d.Set("{{.Attributes.name}}",utils.FlattenListOfPrimitives(&resource.{{.Attributes.modelName}}))
     {{else}}
     err=d.Set("{{.Attributes.name}}",utils.FlattenListOfStringsList(&resource.{{.Attributes.modelName}} ,[]string{{$cbl}} {{ToListOfStrings .Attributes.names}} {{$cbr}}))
     {{end}}
     {{end -}}
     {{if and (eq .Attributes.type "TypeList") (eq .Attributes.set_type "Int") -}}
     {{ if eq .Attributes.length "1" }}
     err=d.Set("{{.Attributes.name}}",utils.FlattenListOfPrimitives(&resource.{{.Attributes.modelName}} ))
     {{end -}}
     {{end -}}
     {{if and (eq .Attributes.type "TypeList") (ne .Attributes.set_type "String" ) (ne .Attributes.set_type "Int") -}}
     {{if (eq .Attributes.set_access "Object") }} 
     err=d.Set("{{.Attributes.name}}",utils.FlattenModelAsList(ctx,*resource.{{.Attributes.modelName}}))
     {{end -}}
     {{if (eq .Attributes.set_access "List") }} 
     err=d.Set("{{.Attributes.name}}",utils.FlattenListOfModelsToList(ctx,resource.{{.Attributes.modelName}}))
     {{end -}}
     {{if (eq .Attributes.set_access "Pointer") }}
     tflog.Debug(ctx,fmt.Sprintf("Found a pointer object %v", resource.{{.Attributes.modelName}})) 
     err=d.Set("{{.Attributes.name}}",utils.FlattenModelAsList(ctx,resource.{{.Attributes.modelName}}))
     {{end -}}
     {{end -}}
     {{else}}
     err=d.Set("{{.Attributes.name}}",resource.{{.Attributes.modelName}})
     {{end}}
     if err!=nil {
          diags = append(diags, diag.Diagnostic {
		Severity: diag.Error,
		Summary:  "Error occured setting value to \"{{.Attributes.name}}\"",
		Detail:   err.Error(),
		})
          }

     {{ end }}
     Id:=(int64)(resource.Id)
     d.SetId(strconv.FormatInt(Id,10))
     return diags
}
`
	_, exists := funcMap["BuildTemplateFromModelName"]
	if !exists {
		funcMap["BuildTemplateFromModelName"] = BuildTemplateFromModelName
	}

	t := template.Must(template.New("read_function").Funcs(funcMap).Parse(read_function))
	err := t.Execute(&b, r)
	if err != nil {
		fmt.Println(err)
	}
	return b.String()

}
func GenDataSourceTemplate(r codegen_configs.ResourceTemplateV2) string {
	func_footer := `     },
   }
}
`
	return BuildDataSourceTemplateHeader(r) + BuildTemplate(r, 0) + func_footer + BuildDataSourceTemplateReadFunction(r)
}
