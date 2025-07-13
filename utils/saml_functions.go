package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vast-data/terraform-provider-vastdata/vast-client"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func getSamlId(vmsId int, idpName string) string {
	return fmt.Sprintf("%v-%v", vmsId, idpName)
}

func getPathAndQuery(path string, vmsId int, idpName string) (string, string) {
	pathWithVmsId := fmt.Sprintf(path, vmsId)
	query := fmt.Sprintf("idp_name=%v", idpName)
	return pathWithVmsId, query
}

func getKeyFromMap(dict map[string]interface{}, fallback string) string {
	for key := range dict {
		return key
	}
	return fallback
}

func unmarshallBody(response *http.Response, vmsId int, idpName string, fallbackIdpEntityid string) (*map[string]interface{}, error) {
	unmarshalledBody := map[string]interface{}{}
	err := UnmarshalBodyToMap(response, &unmarshalledBody)
	if err != nil {
		return nil, err
	}
	id := getSamlId(vmsId, idpName)

	unmarshalledBody["id"] = id
	_metadata := unmarshalledBody["metadata"].(map[string]interface{})
	if remote, ok := _metadata["remote"].([]interface{}); ok && len(remote) > 0 {
		if first, ok := remote[0].(map[string]interface{}); ok {
			if url, ok := first["url"].(string); ok {
				unmarshalledBody["idp_metadata_url"] = url
			}
		}
	}
	_idp := unmarshalledBody["idp"].(map[string]interface{})
	_spSettings := unmarshalledBody["sp_settings"].(map[string]interface{})
	unmarshalledBody["idp_name"] = idpName
	unmarshalledBody["vms_id"] = vmsId
	unmarshalledBody["idp_entityid"] = getKeyFromMap(_idp, fallbackIdpEntityid)
	unmarshalledBody["encrypt_assertion"] = _spSettings["encrypt_assertion"]
	unmarshalledBody["force_authn"] = _spSettings["force_authn"]
	unmarshalledBody["want_assertions_or_response_signed"] = _spSettings["want_assertions_or_response_signed"]
	return &unmarshalledBody, nil
}

func SamlCreateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	client := _client.(*vast_client.VMSSession)
	payload := map[string]interface{}{}
	payload["saml_settings"] = data
	marshalledBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	path := attr["path"].(string)
	vmsId := data["vms_id"].(int)
	idpName := data["idp_name"].(string)
	pathWithVmsId, query := getPathAndQuery(path, vmsId, idpName)

	tflog.Debug(ctx, fmt.Sprintf("Calling POST to path \"%v\"", pathWithVmsId))
	response, err := client.Post(ctx, pathWithVmsId, query, bytes.NewReader(marshalledBody), map[string]string{})
	if err != nil {
		return nil, err
	}
	data["id"] = getSamlId(vmsId, idpName)
	return FakeHttpResponse(response, data)
}

func SamlGetFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	client := _client.(*vast_client.VMSSession)

	path := attr["path"].(string)
	vmsId := d.Get("vms_id").(int)
	idpName := d.Get("idp_name").(string)
	idpMetadata := d.Get("idp_metadata").(string)
	signingCert := d.Get("signing_cert").(string)
	signingKey := d.Get("signing_key").(string)
	encryptionSamlCrt := d.Get("encryption_saml_crt").(string)
	encryptionSamlKey := d.Get("encryption_saml_key").(string)
	idpEntityid := d.Get("idp_entityid").(string)

	pathWithVmsId, query := getPathAndQuery(path, vmsId, idpName)

	tflog.Debug(ctx, fmt.Sprintf("Calling GET to path \"%v\" , with Query %v", path, query))
	response, err := client.Get(ctx, pathWithVmsId, query, headers)
	if err != nil {
		return nil, err
	}

	unmarshalledBody, err := unmarshallBody(response, vmsId, idpName, idpEntityid)
	if err != nil {
		return nil, err
	}
	(*unmarshalledBody)["idp_metadata"] = idpMetadata
	(*unmarshalledBody)["signing_cert"] = signingCert
	(*unmarshalledBody)["signing_key"] = signingKey
	(*unmarshalledBody)["encryption_saml_crt"] = encryptionSamlCrt
	(*unmarshalledBody)["encryption_saml_key"] = encryptionSamlKey

	return FakeHttpResponse(response, *unmarshalledBody)
}

func SamlBeforeDeleteFunc(ctx context.Context, d *schema.ResourceData, m interface{}) (io.Reader, error) {
	data := map[string]interface{}{"vms_id": d.Get("vms_id").(int), "idp_name": d.Get("idp_name").(string)}
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func SamlDeleteFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	client := _client.(*vast_client.VMSSession)

	path := attr["path"].(string)
	vmsId := int(data["vms_id"].(float64))
	idpName := data["idp_name"].(string)

	pathWithVmsId, query := getPathAndQuery(path, vmsId, idpName)

	tflog.Debug(ctx, fmt.Sprintf("Calling Delete to path \"%v\"", pathWithVmsId))
	return client.Delete(ctx, pathWithVmsId, query, nil, map[string]string{})
}

func SamlUpdateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	client := _client.(*vast_client.VMSSession)
	path := attr["path"].(string)
	vmsId := d.Get("vms_id").(int)
	idpName := d.Get("idp_name").(string)

	pathWithVmsId, query := getPathAndQuery(path, vmsId, idpName)

	payload := map[string]interface{}{}
	payload["saml_settings"] = data
	marshalledBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling PATCH to path \"%v\"", pathWithVmsId))
	tflog.Debug(ctx, fmt.Sprintf("Calling PATCH with payload: %v", string(marshalledBody)))
	response, err := client.Post(ctx, pathWithVmsId, query, bytes.NewReader(marshalledBody), map[string]string{})
	if err != nil {
		return nil, err
	}
	data["id"] = getSamlId(vmsId, idpName)
	return FakeHttpResponse(response, data)
}

func SamlImportFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, g GetFuncType) (*http.Response, error) {
	client := _client.(*vast_client.VMSSession)

	path := attr["path"].(string)
	passedId := d.Id()
	splitId := strings.Split(passedId, "|")
	if len(splitId) != 2 {
		return nil, fmt.Errorf("import function must contain exactly one '|' separator")
	}
	vmsId, err := strconv.Atoi(splitId[0])
	if err != nil {
		return nil, err
	}
	idpName := splitId[1]

	pathWithVmsId, query := getPathAndQuery(path, vmsId, idpName)

	tflog.Debug(ctx, fmt.Sprintf("Calling GET to path \"%v\" , with Query %v", path, query))
	response, err := client.Get(ctx, pathWithVmsId, query, map[string]string{})
	if err != nil {
		return nil, err
	}

	unmarshalledBody, err := unmarshallBody(response, vmsId, idpName, "")
	if err != nil {
		return nil, err
	}
	var list []*map[string]interface{}
	list = append(list, unmarshalledBody)
	return FakeHttpResponseAny(response, list)
}

func SamlProcessingFunc(ctx context.Context, response *http.Response, d *schema.ResourceData) ([]byte, error) {
	vmsId := d.Get("vms_id").(int)
	idpName := d.Get("idp_name").(string)
	unmarshalledBody, err := unmarshallBody(response, vmsId, idpName, "")
	if err != nil {
		return nil, err
	}
	var list []*map[string]interface{}
	list = append(list, unmarshalledBody)
	fakeResponse, err := FakeHttpResponseAny(response, list)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(fakeResponse.Body)
	tflog.Debug(ctx, fmt.Sprintf("HTTP Response body %s", string(body)))
	return body, err
}
