package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	//        "net/url"
	"errors"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	codegen_configs "github.com/vast-data/terraform-provider-vastdata/codegen_tools/configs"
	metadata "github.com/vast-data/terraform-provider-vastdata/metadata"
	utils "github.com/vast-data/terraform-provider-vastdata/utils"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
	vast_versions "github.com/vast-data/terraform-provider-vastdata/vast_versions"
)

func ResourceSaml() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceSamlRead,
		DeleteContext: resourceSamlDelete,
		CreateContext: resourceSamlCreate,
		UpdateContext: resourceSamlUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSamlImporter,
		},

		Description: ``,
		Schema:      getResourceSamlSchema(),
	}
}

func getResourceSamlSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"vms_id": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Saml").GetConflictingFields("vms_id"),

			Required:    true,
			Description: `(Valid for versions: 5.1.0,5.2.0) VMS ID`,
		},

		"idp_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Saml").GetConflictingFields("idp_name"),

			Required:    true,
			Description: `(Valid for versions: 5.1.0,5.2.0) SAML IDP name`,
		},

		"encrypt_assertion": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Saml").GetConflictingFields("encrypt_assertion"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Set to true if the IdP encrypts the assertion. If true, an encryption certificate and key must be uploaded. Use encryption_saml_crt and encryption_saml_key to provide the required certificate and key. Default: false. Set to false to disable encryption.`,
		},

		"encryption_saml_crt": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Saml").GetConflictingFields("encryption_saml_crt"),

			Computed:    true,
			Optional:    true,
			Sensitive:   true,
			Description: `(Valid for versions: 5.1.0,5.2.0) Specifies the encryption certificate file content to upload. Required if encrypt_assertion is true.`,
		},

		"encryption_saml_key": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Saml").GetConflictingFields("encryption_saml_key"),

			Computed:    true,
			Optional:    true,
			Sensitive:   true,
			Description: `(Valid for versions: 5.1.0,5.2.0) Specifies the encryption key file content to upload. Required if encrypt_assertion is true.`,
		},

		"force_authn": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Saml").GetConflictingFields("force_authn"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Set to true to force authentication with the IDP even if there is an active session with the IdP for the user. Default: false.`,
		},

		"idp_entityid": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Saml").GetConflictingFields("idp_entityid"),

			Required:    true,
			Description: `(Valid for versions: 5.1.0,5.2.0) A unique identifier for the IdP instance`,
		},

		"idp_metadata": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Saml").GetConflictingFields("idp_metadata"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Use local metadata. Supply local metadata XML.`,
		},

		"idp_metadata_url": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Saml").GetConflictingFields("idp_metadata_url"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Use metadata located at specified remote URL. For example: 'https://dev-12914105.okta.com/app/exke7ia133bKXWP2g5d7/sso/saml/metadata'`,
		},

		"signing_cert": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Saml").GetConflictingFields("signing_cert"),

			Computed:    true,
			Optional:    true,
			Sensitive:   true,
			Description: `(Valid for versions: 5.1.0,5.2.0) Specifies the certificate file content to use for requiring signed responses from the IdP. Required if want_assertions_or_response_signed is true.`,
		},

		"signing_key": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Saml").GetConflictingFields("signing_key"),

			Computed:    true,
			Optional:    true,
			Sensitive:   true,
			Description: `(Valid for versions: 5.1.0,5.2.0) Specifies the key file content to use for requiring signed responses from the IdP. Required if want_assertions_or_response_signed is true.`,
		},

		"want_assertions_or_response_signed": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Saml").GetConflictingFields("want_assertions_or_response_signed"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Set to true to require a signed response or assertion from the IdP. VMS then fails the user authentication if an unsigned response is received. If true, upload a certificate and key. Use signing_cert and signing_key to provide certificate and key. Default: false.`,
		},
	}
}

var SamlNamesMapping = map[string][]string{}

func ResourceSamlReadStructIntoSchema(ctx context.Context, resource api_latest.Saml, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "VmsId", resource.VmsId))

	err = d.Set("vms_id", resource.VmsId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"vms_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IdpName", resource.IdpName))

	err = d.Set("idp_name", resource.IdpName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"idp_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "EncryptAssertion", resource.EncryptAssertion))

	err = d.Set("encrypt_assertion", resource.EncryptAssertion)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"encrypt_assertion\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "EncryptionSamlCrt", resource.EncryptionSamlCrt))

	err = d.Set("encryption_saml_crt", resource.EncryptionSamlCrt)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"encryption_saml_crt\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "EncryptionSamlKey", resource.EncryptionSamlKey))

	err = d.Set("encryption_saml_key", resource.EncryptionSamlKey)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"encryption_saml_key\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ForceAuthn", resource.ForceAuthn))

	err = d.Set("force_authn", resource.ForceAuthn)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"force_authn\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IdpEntityid", resource.IdpEntityid))

	err = d.Set("idp_entityid", resource.IdpEntityid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"idp_entityid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IdpMetadata", resource.IdpMetadata))

	err = d.Set("idp_metadata", resource.IdpMetadata)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"idp_metadata\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IdpMetadataUrl", resource.IdpMetadataUrl))

	err = d.Set("idp_metadata_url", resource.IdpMetadataUrl)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"idp_metadata_url\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SigningCert", resource.SigningCert))

	err = d.Set("signing_cert", resource.SigningCert)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"signing_cert\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SigningKey", resource.SigningKey))

	err = d.Set("signing_key", resource.SigningKey)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"signing_key\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "WantAssertionsOrResponseSigned", resource.WantAssertionsOrResponseSigned))

	err = d.Set("want_assertions_or_response_signed", resource.WantAssertionsOrResponseSigned)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"want_assertions_or_response_signed\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceSamlRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Saml")
	attrs := map[string]interface{}{"path": utils.GenPath("vms/%v/saml_config"), "id": d.Id()}
	tflog.Debug(ctx, fmt.Sprintf("[resourceSamlRead] Calling Get Function : %v for resource Saml", utils.GetFuncName(resourceConfig.GetFunc)))
	response, err := resourceConfig.GetFunc(ctx, client, attrs, d, map[string]string{})
	utils.VastVersionsWarn(ctx)

	var body []byte
	var resource api_latest.Saml
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
	err = json.Unmarshal(body, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred while parsing data received from VAST Data cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	diags = ResourceSamlReadStructIntoSchema(ctx, resource, d)

	return diags
}

func resourceSamlDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Saml")
	attrs := map[string]interface{}{"path": utils.GenPath("vms/%v/saml_config"), "id": d.Id()}

	data, beforeDeleteError := resourceConfig.BeforeDeleteFunc(ctx, d, m)
	if beforeDeleteError != nil {
		return diag.FromErr(beforeDeleteError)
	}
	unmarshalledData := map[string]interface{}{}
	_data, _ := io.ReadAll(data)
	err := json.Unmarshal(_data, &unmarshalledData)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to unmarshall json data",
			Detail:   err.Error(),
		})
		return diags
	}
	response, err := resourceConfig.DeleteFunc(ctx, client, attrs, unmarshalledData, map[string]string{})

	tflog.Info(ctx, fmt.Sprintf("Removing Resource"))
	if response != nil {
		tflog.Info(ctx, response.Request.URL.String())
		tflog.Info(ctx, utils.GetResponseBodyAsStr(response))
	}

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred while deleting a resource from the VAST Data cluster",
			Detail:   err.Error(),
		})

	}

	return diags

}

func resourceSamlCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("names_mapping")
	newCtx := context.WithValue(ctx, namesMapping, SamlNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Saml")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource Saml"))
	reflectSaml := reflect.TypeOf((*api_latest.Saml)(nil))
	utils.PopulateResourceMap(newCtx, reflectSaml.Elem(), d, &data, "", false)

	versionsEqual := utils.VastVersionsWarn(ctx)

	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "Saml")
		if typeExists {
			versionError := utils.VersionMatch(t, data)
			if versionError != nil {
				tflog.Warn(ctx, versionError.Error())
				versionValidationMode, versionValidationModeExists := metadata.GetClusterConfig("version_validation_mode")
				tflog.Warn(ctx, fmt.Sprintf("Version Validation Mode Detected %s", versionValidationMode))
				if versionValidationModeExists && versionValidationMode == "strict" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Cluster Version & Build Version Are Too Different",
						Detail:   versionError.Error(),
					})
					return diags
				}
			}
		} else {
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "Saml", clusterVersion))
		}
	}
	tflog.Debug(ctx, fmt.Sprintf("Data %v", data))
	b, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could have not generate request json",
			Detail:   err.Error(),
		})
		return diags
	}
	tflog.Debug(ctx, fmt.Sprintf("Request json created %v", string(b)))
	attrs := map[string]interface{}{"path": utils.GenPath("vms/%v/saml_config")}
	response, createErr := resourceConfig.CreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  Saml %v", createErr))

	if createErr != nil {
		errorMessage := fmt.Sprintf("server response:\n%v\nUnderlying error:\n%v", utils.GetResponseBodyAsStr(response), createErr.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   errorMessage,
		})
		return diags
	}
	responseBody, _ := io.ReadAll(response.Body)
	tflog.Debug(ctx, fmt.Sprintf("Object created, server response %v", string(responseBody)))
	resource := api_latest.Saml{}
	err = json.Unmarshal(responseBody, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into Saml",
			Detail:   err.Error(),
		})
		return diags
	}

	err = resourceConfig.IdFunc(ctx, client, resource.Id, d)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to set Id",
			Detail:   err.Error(),
		})
		return diags
	}
	ctxWithResource := context.WithValue(ctx, utils.ContextKey("resource"), resource)
	resourceSamlRead(ctxWithResource, d, m)

	return diags
}

func resourceSamlUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("names_mapping")
	newCtx := context.WithValue(ctx, namesMapping, SamlNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	versionsEqual := utils.VastVersionsWarn(ctx)
	resourceConfig := codegen_configs.GetResourceByName("Saml")
	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "Saml")
		if typeExists {
			versionError := utils.VersionMatch(t, data)
			if versionError != nil {
				tflog.Warn(ctx, versionError.Error())
				versionValidationMode, versionValidationModeExists := metadata.GetClusterConfig("version_validation_mode")
				tflog.Warn(ctx, fmt.Sprintf("Version Validation Mode Detected %s", versionValidationMode))
				if versionValidationModeExists && versionValidationMode == "strict" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Cluster Version & Build Version Are Too Different",
						Detail:   versionError.Error(),
					})
					return diags
				}
			}
		} else {
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "Saml", clusterVersion))
		}
	}

	client := m.(*vast_client.VMSSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource Saml"))
	reflectSaml := reflect.TypeOf((*api_latest.Saml)(nil))
	utils.PopulateResourceMap(newCtx, reflectSaml.Elem(), d, &data, "", false)

	tflog.Debug(ctx, fmt.Sprintf("Data %v", data))
	b, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could have not generate request json",
			Detail:   err.Error(),
		})
		return diags
	}
	tflog.Debug(ctx, fmt.Sprintf("Request json created %v", string(b)))
	attrs := map[string]interface{}{"path": utils.GenPath("vms/%v/saml_config"), "id": d.Id()}
	response, patchErr := resourceConfig.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  Saml %v", patchErr))
	if patchErr != nil {
		errorMessage := fmt.Sprintf("server response:\n%v\nUnderlying error:\n%v", utils.GetResponseBodyAsStr(response), patchErr.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   errorMessage,
		})
		return diags
	}
	resourceSamlRead(ctx, d, m)

	return diags

}

func resourceSamlImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	var result []*schema.ResourceData
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Saml")
	attrs := map[string]interface{}{"path": utils.GenPath("vms/%v/saml_config")}
	response, err := resourceConfig.ImportFunc(ctx, client, attrs, d, resourceConfig.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	var resourceList []api_latest.Saml
	body, err := resourceConfig.ResponseProcessingFunc(ctx, response, d)

	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &resourceList)
	if err != nil {
		return result, err
	}

	if len(resourceList) == 0 {
		return result, errors.New("cluster returned 0 elements matching provided guid")
	}

	resource := resourceList[0]
	idErr := resourceConfig.IdFunc(ctx, client, resource.Id, d)
	if idErr != nil {
		return result, idErr
	}

	diags := ResourceSamlReadStructIntoSchema(ctx, resource, d)
	if diags.HasError() {
		allErrors := "Errors occurred while importing:\n"
		for _, dig := range diags {
			allErrors += fmt.Sprintf("Summary:%s\nDetails:%s\n", dig.Summary, dig.Detail)
		}
		return result, errors.New(allErrors)
	}
	result = append(result, d)

	return result, err

}
