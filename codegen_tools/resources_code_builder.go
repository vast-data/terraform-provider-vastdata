package main

import (
	"bytes"
	"fmt"
	"text/template"

	codegen_configs "github.com/vast-data/terraform-provider-vastdata/codegen_tools/configs"
)

func BuildResourcesList(resources_templates []codegen_configs.ResourceTemplateV2) string {
	var b bytes.Buffer
	resources :=
		`package resources

import (
 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
       
)

var Resources map[string]*schema.Resource = map[string]*schema.Resource{
{{range . }}
{{- if .Generate }}     "{{.DataSourceName}}":Resource{{ .ResourceName -}}(),
{{ end -}}
{{end -}}
}
`

	t := template.Must(template.New("resources").Parse(resources))
	err := t.Execute(&b, resources_templates)
	if err != nil {
		fmt.Println(err)
	}
	return b.String()

}

func BuildResourceTemplateHeader(r codegen_configs.ResourceTemplateV2) string {
	var b bytes.Buffer
	header :=
		`package resources

import (
        "io"
        "reflect"
        "encoding/json"
        "fmt"
        "context"
        {{ if not .DisableImport }}
//        "net/url"
        "errors"
        {{ end }}
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
        "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
        vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
        api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
        "github.com/hashicorp/terraform-plugin-log/tflog"
        utils "github.com/vast-data/terraform-provider-vastdata/utils"
    	metadata "github.com/vast-data/terraform-provider-vastdata/metadata"
        vast_versions  "github.com/vast-data/terraform-provider-vastdata/vast_versions"
        codegen_configs "github.com/vast-data/terraform-provider-vastdata/codegen_tools/configs"
        
)

func Resource{{ .ResourceName }}() *schema.Resource {
     return  &schema.Resource{
     ReadContext: resource{{ .ResourceName }}Read,
     DeleteContext: resource{{ .ResourceName }}Delete,
     CreateContext: resource{{ .ResourceName }}Create,
     UpdateContext: resource{{ .ResourceName }}Update,
     {{ if not .DisableImport }}
     Importer: &schema.ResourceImporter{
                      StateContext: resource{{ .ResourceName }}Importer,
     },
     {{ end }}
     {{ if .Timeouts }}
     Timeouts: codegen_configs.GetResourceByName("{{ .ResourceName }}").GetResourceTimeouts(),
     {{ end }}
     Description: {{getBT}}{{ .ResourceDocumantation }}{{getBT}},
     Schema: getResource{{ .ResourceName }}Schema(),
   }
}

func getResource{{ .ResourceName }}Schema() map[string]*schema.Schema {
     return map[string]*schema.Schema{


`

	t := template.Must(template.New("header").Funcs(funcMap).Parse(header))
	err := t.Execute(&b, r)
	if err != nil {
		fmt.Println(err)
	}
	return b.String()
}

func BuildResourceTemplateReadFunction(r codegen_configs.ResourceTemplateV2) string {
	var b bytes.Buffer
	read_function := `
var {{ .ResourceName }}NamesMapping = map[string][]string{
     {{range .Fields}}
     {{- if and (eq .Attributes.type "TypeList") ( eq .Attributes.list_type "simple") (eq .Attributes.length "2" ) -}}
         "{{.Attributes.name}}": []string{"{{ replaceAll .Attributes.names ","  "\",\""}}"},
     {{- end -}}
     {{end}}
     }


func Resource{{ .ResourceName }}ReadStructIntoSchema(ctx context.Context, resource api_latest.{{.ResourceName}} ,d *schema.ResourceData) diag.Diagnostics {
     var diags diag.Diagnostics
     var err error
     {{ $cbr:="}" }}
     {{ $cbl:="{" }}
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
     err=d.Set("{{.Attributes.name}}",utils.FlattenListOfPrimitives(&resource.{{.Attributes.modelName}}))
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
		Summary:  "Error occurred setting value to \"{{.Attributes.name}}\"",
		Detail:   err.Error(),
		})
          }

     {{ end }}
     return diags


}
func resource{{ .ResourceName }}Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
     var diags diag.Diagnostics
     {{ $cbr:="}" }}
     {{ $cbl:="{" }}
     client:=m.(*vast_client.VMSSession)
     resourceConfig := codegen_configs.GetResourceByName("{{ .ResourceName }}")
     attrs:=map[string]interface{}{"path":utils.GenPath("{{.Path}}"),"id":d.Id()}
     tflog.Debug(ctx,fmt.Sprintf("[resource{{ .ResourceName }}Read] Calling Get Function : %v for resource {{ .ResourceName }}",utils.GetFuncName(resourceConfig.GetFunc)))  
     response,err:=resourceConfig.GetFunc(ctx,client,attrs,d,map[string]string{})
     utils.VastVersionsWarn(ctx)
	
	var body []byte
	var resource api_latest.{{.ResourceName}}
	if err != nil && response != nil && response.StatusCode == 404 && !resourceConfig.DisableFallbackRequest {
		var fallbackErr error
		body, fallbackErr = utils.HandleFallback(ctx, client, attrs, d, resourceConfig.IdFunc)
		if fallbackErr != nil {
			errorMessage := fmt.Sprintf("Initial request failed:\n%v\nFallback request also failed:\n%v", err.Error(), fallbackErr.Error())
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error occurred while obtaining data from the VAST Data cluster",
				Detail:   errorMessage,
			})
			return diags
		}
	} else if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred while obtaining data from the VAST Data cluster",
			Detail:   err.Error(),
		})
		return diags
	} else {
		tflog.Info(ctx, response.Request.URL.String())
		body, err = resourceConfig.ResponseProcessingFunc(ctx, response, d)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error occurred reading data received from VAST Data cluster",
				Detail:   err.Error(),
			})
			return diags
		}
	}
     err=json.Unmarshal(body,&resource)
     if err!=nil {
                diags = append(diags, diag.Diagnostic {
		Severity: diag.Error,
		Summary:  "Error occurred while parsing data received from VAST Data cluster",
		Detail:   err.Error(),
		})
       return diags

     }
 diags = Resource{{ .ResourceName }}ReadStructIntoSchema(ctx, resource ,d )
 {{ if .AfterReadFunc }}
 var after_read_error error
 after_read_error=resourceConfig.AfterReadFunc(client,ctx,d)
 if after_read_error!=nil {
    return diag.FromErr(after_read_error)
 }
 {{end}}
 return diags 
}

func resource{{ .ResourceName }}Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
     var diags diag.Diagnostics
     client:=m.(*vast_client.VMSSession)
     resourceConfig := codegen_configs.GetResourceByName("{{ .ResourceName }}")
     attrs:=map[string]interface{}{"path":utils.GenPath("{{.Path}}"),"id":d.Id()}
     {{ if .BeforeDeleteFunc  }}
     data,beforeDeleteError:=resourceConfig.BeforeDeleteFunc(ctx,d,m)
     if beforeDeleteError!=nil {
        return diag.FromErr(beforeDeleteError)
     }
     unmarshalledData := map[string]interface{}{}
     _data,_:=io.ReadAll(data)
     err := json.Unmarshal(_data, &unmarshalledData)
     if err!=nil {
          diags = append(diags, diag.Diagnostic {
  		Severity: diag.Error,
  		Summary:  "Failed to unmarshall json data",
  		Detail:   err.Error(),
  		})
          return diags
      }
     response,err:=resourceConfig.DeleteFunc(ctx,client,attrs,unmarshalledData,map[string]string{});
     {{else}}
     response,err:=resourceConfig.DeleteFunc(ctx,client,attrs,nil,map[string]string{});
     {{end}}
     tflog.Info(ctx,fmt.Sprintf("Removing Resource"))
     if response != nil {
		 tflog.Info(ctx,response.Request.URL.String())
		 tflog.Info(ctx,utils.GetResponseBodyAsStr(response))
	}

     if err!=nil {
        diags = append(diags, diag.Diagnostic {
		Severity: diag.Error,
		Summary:  "Error occurred while deleting a resource from the VAST Data cluster",
		Detail:   err.Error(),
		})

     }

 return diags

}

func resource{{ .ResourceName }}Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    namesMapping := utils.ContextKey("names_mapping")
    newCtx := context.WithValue(ctx, namesMapping, {{ .ResourceName }}NamesMapping)
    var diags diag.Diagnostics
    data := make(map[string]interface{})
    client:=m.(*vast_client.VMSSession)
    resourceConfig := codegen_configs.GetResourceByName("{{ .ResourceName }}")
    tflog.Info(ctx,fmt.Sprintf("Creating Resource {{.ResourceName}}"))
    reflect{{.ResourceName}} := reflect.TypeOf((*api_latest.{{.ResourceName}})(nil))
    utils.PopulateResourceMap(newCtx, reflect{{.ResourceName}}.Elem(),d, &data,"",false)
    {{ if  .BeforePostFunc  }}
    var before_post_error error
    data,before_post_error=resourceConfig.BeforePostFunc(data,client,ctx,d)
    if before_post_error!=nil {
       return diag.FromErr(before_post_error)
    }
    {{end}}
    versionsEqual:=utils.VastVersionsWarn(ctx)
   
    if versionsEqual!= metadata.CLUSTER_VERSION_EQUALS {
          clusterVersion:=metadata.ClusterVersionString()
          t,typeExists:=vast_versions.GetVersionedType(clusterVersion,"{{.ResourceName}}")
          if typeExists {
          versionError:=utils.VersionMatch(t,data) 
          if versionError!=nil {
               tflog.Warn(ctx,versionError.Error())
               versionValidationMode,versionValidationModeExists:=metadata.GetClusterConfig("version_validation_mode")
               tflog.Warn(ctx,fmt.Sprintf("Version Validation Mode Detected %s",versionValidationMode))
               if versionValidationModeExists && versionValidationMode=="strict" {
		    diags = append(diags, diag.Diagnostic {
			    Severity: diag.Error,
			    Summary:  "Cluster Version & Build Version Are Too Different",
			    Detail:   versionError.Error(),
			    })
		    return diags                           
                    }                       
             }
          } else {
             tflog.Warn(ctx,fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly","{{ .ResourceName }}",clusterVersion))
          }
    }     
    tflog.Debug(ctx,fmt.Sprintf("Data %v" , data))    
    b,err:=json.MarshalIndent(data,"","   ")
    if err!=nil {
        diags = append(diags, diag.Diagnostic {
		Severity: diag.Error,
		Summary:  "Could have not generate request json",
		Detail:   err.Error(),
		})
        return diags
    }
    tflog.Debug(ctx,fmt.Sprintf("Request json created %v", string(b)))
    attrs:=map[string]interface{}{"path":utils.GenPath("{{.Path}}")}
    response ,createErr:=resourceConfig.CreateFunc(ctx,client,attrs,data,map[string]string{});
    tflog.Info(ctx,fmt.Sprintf("Server Error for  {{.ResourceName}} %v" , createErr))
    
    if createErr != nil {
            errorMessage:=fmt.Sprintf("server response:\n%v\nUnderlying error:\n%v", utils.GetResponseBodyAsStr(response), createErr.Error()) 
            diags = append(diags, diag.Diagnostic {
		Severity: diag.Error,
		Summary:  "Object Creation Failed",
		Detail:   errorMessage,
		})
        return diags
     }
   responseBody,_:=io.ReadAll(response.Body)
   tflog.Debug(ctx,fmt.Sprintf("Object created, server response %v", string(responseBody)))
   resource:=api_latest.{{.ResourceName}}{}
   err=json.Unmarshal(responseBody,&resource)
   if err!=nil {
        diags = append(diags, diag.Diagnostic {
		Severity: diag.Error,
		Summary:  "Failed to convert response body into {{.ResourceName}}",
		Detail:   err.Error(),
		})
        return diags
    }
   
   err=resourceConfig.IdFunc(ctx,client,resource.Id,d)
   if err!=nil {
        diags = append(diags, diag.Diagnostic {
		Severity: diag.Error,
		Summary:  "Failed to set Id",
		Detail:   err.Error(),
		})
        return diags
    }
   ctxWithResource:=context.WithValue(ctx, utils.ContextKey("resource"), resource)
   resource{{ .ResourceName }}Read(ctxWithResource,d,m)
    {{ if .BeforeCreateFunc }}
    var beforeCreateErr error
    _,beforeCreateErr=resourceConfig.BeforeCreateFunc(data,client,ctx,d)
    if beforeCreateErr!=nil {
       return diag.FromErr(beforeCreateErr)
    }

    {{end}}
   return diags
}

func resource{{ .ResourceName }}Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    namesMapping := utils.ContextKey("names_mapping")
    newCtx := context.WithValue(ctx, namesMapping, {{ .ResourceName }}NamesMapping)
    var diags diag.Diagnostics
    data := make(map[string]interface{})
    versionsEqual:=utils.VastVersionsWarn(ctx)
    resourceConfig := codegen_configs.GetResourceByName("{{ .ResourceName }}")
    if versionsEqual!= metadata.CLUSTER_VERSION_EQUALS {
          clusterVersion:=metadata.ClusterVersionString()
          t,typeExists:=vast_versions.GetVersionedType(clusterVersion,"{{.ResourceName}}")
          if typeExists {
          versionError:=utils.VersionMatch(t,data) 
          if versionError!=nil {
               tflog.Warn(ctx,versionError.Error())
               versionValidationMode,versionValidationModeExists:=metadata.GetClusterConfig("version_validation_mode")
               tflog.Warn(ctx,fmt.Sprintf("Version Validation Mode Detected %s",versionValidationMode))
               if versionValidationModeExists && versionValidationMode=="strict" {
		    diags = append(diags, diag.Diagnostic {
			    Severity: diag.Error,
			    Summary:  "Cluster Version & Build Version Are Too Different",
			    Detail:   versionError.Error(),
			    })
		    return diags                           
                    }                       
             }
          } else {
             tflog.Warn(ctx,fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly","{{ .ResourceName }}",clusterVersion))
          }
    }     

    client:=m.(*vast_client.VMSSession)
    tflog.Info(ctx,fmt.Sprintf("Updating Resource {{.ResourceName}}"))
    reflect{{.ResourceName}} := reflect.TypeOf((*api_latest.{{.ResourceName}})(nil))
    utils.PopulateResourceMap(newCtx, reflect{{.ResourceName}}.Elem(),d, &data,"",false)
    {{ if .BeforePatchFunc }}
    var beforePatchError error
    data,beforePatchError=resourceConfig.BeforePatchFunc(data,client,ctx,d)
    if beforePatchError!=nil {
       return diag.FromErr(beforePatchError)
    }

    {{end}}
    tflog.Debug(ctx,fmt.Sprintf("Data %v" , data ))    
    b,err:=json.MarshalIndent(data,"","   ")
    if err!=nil {
        diags = append(diags, diag.Diagnostic {
		Severity: diag.Error,
		Summary:  "Could have not generate request json",
		Detail:   err.Error(),
		})
        return diags
    }
    tflog.Debug(ctx,fmt.Sprintf("Request json created %v", string(b)))
    attrs:=map[string]interface{}{"path":utils.GenPath("{{.Path}}"),"id":d.Id()}
    response ,patchErr := resourceConfig.UpdateFunc(ctx,client,attrs,data,d,map[string]string{})
    tflog.Info(ctx,fmt.Sprintf("Server Error for  {{.ResourceName}} %v" , patchErr))
    if patchErr != nil {
            errorMessage:=fmt.Sprintf("server response:\n%v\nUnderlying error:\n%v", utils.GetResponseBodyAsStr(response), patchErr.Error()) 
            diags = append(diags, diag.Diagnostic {
		Severity: diag.Error,
		Summary:  "Object Creation Failed",
		Detail:   errorMessage,
		})
        return diags
     }
   resource{{ .ResourceName }}Read(ctx,d,m)
   {{ if .AfterPatchFunc }}
   var afterPatchError error
   data,afterPatchError=resourceConfig.AfterPatchFunc(data,client,ctx,d)
   if afterPatchError!=nil {
      return diag.FromErr(afterPatchError)
   }
   {{end}}

   return diags



}
{{ if not .DisableImport }}
func resource{{ .ResourceName }}Importer(ctx context.Context, d *schema.ResourceData, m interface{})  ([]*schema.ResourceData, error) {

    var result []*schema.ResourceData
    client := m.(*vast_client.VMSSession)
    resourceConfig := codegen_configs.GetResourceByName("{{ .ResourceName }}")
    attrs:=map[string]interface{}{"path":utils.GenPath("{{.Path}}")}
    response,err:=resourceConfig.ImportFunc(ctx,client,attrs,d,resourceConfig.Importer.GetFunc())

    if err != nil {
	    return result, err
    }
     
    var resourceList []api_latest.{{.ResourceName}}
    body,err:=resourceConfig.ResponseProcessingFunc(ctx,response, d)

    if err!=nil {
       return result, err
    }
     {{ if .ResponseGetByURL }} 
     body,err = utils.ResponseGetByURL(ctx,body,client)
     if err!=nil {
       return result, err
      }
     {{end -}}

    err=json.Unmarshal(body,&resourceList)
    if err!=nil {
       return result,err
    }

     if len(resourceList) == 0 {
        return result,errors.New("cluster returned 0 elements matching provided guid")
     }
     
    resource:=resourceList[0]
    idErr:=resourceConfig.IdFunc(ctx,client,resource.Id,d)
    if idErr!=nil {
	 return result,idErr
     }

     diags := Resource{{.ResourceName}}ReadStructIntoSchema(ctx, resource, d)
     if diags.HasError() {
         allErrors:="Errors occurred while importing:\n"
         for _,dig := range diags {
           allErrors+=fmt.Sprintf("Summary:%s\nDetails:%s\n",dig.Summary,dig.Detail)
         }
         return result,errors.New(allErrors)
     }
     result=append(result,d)
     
     return result, err

}
{{ end }}
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

func BuildResourceTemplateFromModelName(n string, indent int) string {
	model, _ := resources_templates_map[n]
	return BuildResourceTemplate(model, indent)
}

func ResourceBuildTemplateToTerrafromElem(r codegen_configs.ResourceElem, indent int) string {
	var b bytes.Buffer
	r.Indent = indent
	tmpl := `     
             {{ $I:=.Indent}}
             {{ $name:=.Attributes.name}}
	     {{indent $I " "}}"{{ .Attributes.name }}": {
	     {{indent $I " "}}   Type: 	  schema.{{ .Attributes.type }},
	     {{indent $I " "}}   ConflictsWith: codegen_configs.GetResourceByName("{{GetResourceName}}").GetConflictingFields("{{.Attributes.name}}"),
             {{if eq .Attributes.ignore_update "true" }}
             {{indent $I " "}}   DiffSuppressOnRefresh: false,
             DiffSuppressFunc: utils.DoNothingOnUpdate(),
             {{ end }}
             {{if AttributeHasDiffFunc .Attributes.name  }}
             {{indent $I " "}}   DiffSuppressOnRefresh: false,
             DiffSuppressFunc: codegen_configs.GetResourceByName("{{GetResourceName}}").GetAttributeDiffFunc("{{.Attributes.name}}"),
             {{ end }}
	     {{- if eq .Attributes.required "true" -}}
	     {{indent $I " "}}   Required: true,
	     {{ else -}}
	     {{indent $I " "}}   Computed: {{.Attributes.computed}},
	     {{indent $I " "}}   Optional: {{.Attributes.optional}},
	     {{indent $I " "}}   Sensitive: {{.Attributes.sensitive}},
             {{indent $I " "}} 
             {{- end }}
             {{-  if  .Attributes.validator_func  }}
             {{indent $I " "}}  ValidateDiagFunc: {{.Attributes.validator_func}},
             {{ else }}
             {{-  if  .Attributes.enum }}
             {{indent $I " "}}  ValidateDiagFunc: utils.OneOf({{.Attributes.enum}}),
             {{ end -}}
             {{indent $I " "}}   Description: {{getBT}}{{ GetSchemaProperyDocument .Attributes.name }}{{getBT}},
	     {{ end -}}
             {{- $default_value:=GetSchemaProperyDefault .Attributes.name -}}
             {{  if  ne $default_value "" }}
             {{indent $I " "}}  Default: {{ GetSchemaProperyDefault .Attributes.name }},
             {{ end -}}
             {{- if eq  .Attributes.force_new "true" }}
	     {{- indent $I " "}}   ForceNew: {{.Attributes.force_new}},                          
             {{ end -}}
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
                {{indent $I " "}}               Optional:  true,
                {{indent $I " "}}               Description: "{{GetFakeFieldDescription $name $t}}",
                {{indent $I " "}}         },
                {{end}}
                {{else}} 
                {{ BuildResourceTemplateFromModelName .Attributes.set_type ( AddInt $I 7) }}{{end}}
              {{indent $I " "}}},
             {{indent $I " "}}},{{ end }}
	    {{indent $I " "}}},
`
	//A dirty workaround to avoid loop decleration
	_, exists := funcMap["BuildResourceTemplateFromModelName"]
	if !exists {
		funcMap["BuildResourceTemplateFromModelName"] = BuildResourceTemplateFromModelName
	}
	//Create a local copy to have local only changes
	localFuncMap := template.FuncMap{}

	for k, v := range funcMap {
		localFuncMap[k] = v
	}
	localFuncMap["GetFakeFieldDescription"] = r.Parent.GetFakeFieldDescription
	localFuncMap["GetSchemaProperyDocument"] = r.Parent.GetSchemaProperyDocument
	localFuncMap["GetSchemaProperyDefault"] = r.Parent.GetSchemaProperyDefault
	localFuncMap["HasValidatorFunc"] = r.Parent.HasValidatorFunc
	localFuncMap["GetValidatorFunc"] = r.Parent.GetValidatorFunc
	localFuncMap["AttributeHasDiffFunc"] = r.Parent.AttributeHasDiffFunc
	localFuncMap["GetAttributeDiffFunc"] = r.Parent.GetAttributeDiffFunc
	localFuncMap["GetResourceName"] = func() string {
		return r.Parent.ResourceName
	}
	t := template.Must(template.New("tf").Funcs(localFuncMap).Parse(tmpl))
	err := t.Execute(&b, r)
	if err != nil {
		fmt.Println(err)
	}
	return b.String()

}

func BuildResourceTemplate(r codegen_configs.ResourceTemplateV2, indent int) string {
	out := ""
	for _, f := range r.Fields {
		out += ResourceBuildTemplateToTerrafromElem(f, indent)
	}
	return out
}

func GenResourceTemplate(r codegen_configs.ResourceTemplateV2) string {
	func_footer := `     
   }
}
`
	return BuildResourceTemplateHeader(r) + BuildResourceTemplate(r, 0) + func_footer + BuildResourceTemplateReadFunction(r)
}
