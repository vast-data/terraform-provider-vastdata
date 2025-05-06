package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
	"net/http"
	"strconv"
	"strings"
)

func CreateNonLocalUserKeyFunc(ctx context.Context, _client any, attr map[string]any, data map[string]any, headers map[string]string) (*http.Response, error) {
	client := _client.(*vast_client.VMSSession)
	attributes, err := getAttributesAsString([]string{"path"}, attr)
	if err != nil {
		return nil, err
	}
	path := (*attributes)["path"]
	uid := data["uid"]
	tenantId := data["tenant_id"]
	if pgpPublicKey, ok := data["pgp_public_key"]; ok {
		if _, err := encryptSecretToken("Testing value", pgpPublicKey); err != nil {
			tflog.Error(ctx, "Public GPG key is incorrect.")
			return nil, err
		}
	}

	enabled, ok := data["enabled"]
	if !ok {
		enabled = true
	}
	buffer, marshallingError := json.Marshal(data)
	if marshallingError != nil {
		return nil, marshallingError
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling POST to path \"%v\"", attr))
	response, err := client.Post(ctx, path, bytes.NewReader(buffer), map[string]string{})
	if err != nil {
		return response, err
	}
	responseBody := map[string]any{}
	err = UnmarshalBodyToMap(response, &responseBody)
	if err != nil {
		return nil, err
	}
	accessKey := responseBody["access_key"]
	secretKey := responseBody["secret_key"]
	responseBody["id"] = fmt.Sprintf("%v-%v", uid, accessKey)
	responseBody["uid"] = uid
	responseBody["enabled"] = enabled
	responseBody["tenant_id"] = tenantId
	responseBody["pgp_public_key"] = data["pgp_public_key"]

	if !(enabled.(bool)) {
		payload := map[string]any{"uid": uid, "access_key": accessKey, "enabled": false, "tenant_id": tenantId}
		pBuffer, pMarshallingError := json.Marshal(payload)
		if pMarshallingError != nil {
			return nil, pMarshallingError
		}
		tflog.Debug(ctx, fmt.Sprintf("Disabling NonLocalUserKey %v", accessKey))
		presponse, perr := client.Patch(ctx, path, bytes.NewReader(pBuffer), map[string]string{})
		if perr != nil {
			return presponse, perr
		}
	}
	if pgpKey, ok := data["pgp_public_key"]; ok {
		tflog.Debug(ctx, fmt.Sprintf("NONLOCALUSERKEY: PGP public key found , encrypting secret"))
		encryptedKey, err := encryptSecretToken(secretKey, pgpKey)
		if err != nil {
			return nil, err
		}
		responseBody["secret_key"] = nil
		responseBody["encrypted_secret_key"] = encryptedKey
	} else {
		tflog.Warn(ctx, fmt.Sprintf("NONLOCALUSERKEY: No PGP public key was found secret_key will hold the plain text file"))
		responseBody["encrypted_secret_key"] = nil
	}
	return FakeHttpResponse(response, responseBody)
}

func UpdateNonLocalUserKeyFunc(ctx context.Context, _client any, attr map[string]any, data map[string]any, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	client := _client.(*vast_client.VMSSession)
	attributes, err := getAttributesAsString([]string{"path", "id"}, attr)
	if err != nil {
		return nil, err
	}
	path := (*attributes)["path"]
	s := strings.SplitN((*attributes)["id"], "-", 2)
	uid, err := strconv.Atoi(s[0])
	if err != nil {
		return nil, err
	}
	accessKey := s[1]
	enabled := d.Get("enabled").(bool)
	tenantId := d.Get("tenant_id").(int)
	payload := map[string]any{"uid": uid, "access_key": accessKey, "enabled": enabled, "tenant_id": tenantId}
	buffer, marshallingError := json.Marshal(payload)
	if marshallingError != nil {
		return nil, marshallingError
	}
	tflog.Debug(ctx, fmt.Sprintf("Updating NonLocalUserKey %v with data %v", accessKey, payload))
	response, err := client.Patch(ctx, path, bytes.NewReader(buffer), map[string]string{})
	if err != nil {
		return response, err
	}
	responseBody := payload
	responseBody["id"] = accessKey
	return FakeHttpResponse(response, responseBody)
}

func GetNonLocalUserKeyFunc(ctx context.Context, _client any, attr map[string]any, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	//There is no GET for a key we will have to iterate over all user keys to find this specific key
	client := _client.(*vast_client.VMSSession)
	attributes, err := getAttributesAsString([]string{"id"}, attr)
	if err != nil {
		return nil, err
	}
	resource := ctx.Value(ContextKey("resource"))
	tflog.Debug(ctx, fmt.Sprintf("NonLocalUserKey: Resource %v found", resource))
	if resource != nil {
		r := resource.(api_latest.NonLocalUserKey)
		d.Set("secret_key", r.SecretKey)
		d.Set("tenant_id", r.TenantId)
		d.Set("pgp_public_key", r.PgpPublicKey)
		d.Set("encrypted_secret_key", r.EncryptedSecretKey)
	}
	s := strings.SplitN((*attributes)["id"], "-", 2)
	uid, err := strconv.Atoi(s[0])
	if err != nil {
		return nil, err
	}
	accessKey := s[1]
	tenantId := d.Get("tenant_id").(int)
	secretKey := d.Get("secret_key").(string)
	pgpKey := d.Get("pgp_public_key")
	encryptedKey := d.Get("encrypted_secret_key")
	query := fmt.Sprintf("uid=%v", uid)
	path := GenPath("users/query")
	tflog.Debug(ctx, fmt.Sprintf("Calling GET for uid %v to get user detail", uid))
	response, err := client.Get(ctx, path, query, headers)
	if err != nil {
		return nil, err
	}
	responseBody := map[string]any{}
	err = UnmarshalBodyToMap(response, &responseBody)
	if err != nil {
		return nil, err
	}
	tflog.Debug(ctx, fmt.Sprintf("NonLocalUserKey: Reponse: %v", responseBody))
	accessKeys, exists := responseBody["access_keys"]
	if exists {
		for _, l := range accessKeys.([]any) {
			v := l.(map[string]any)
			key, keyExists := v["key"]
			if !keyExists {
				key = ""
			}
			enabled, enabledExists := v["status"]
			if !enabledExists {
				enabled = ""
			}
			if key == accessKey {
				tflog.Debug(ctx, fmt.Sprintf("NonLocalUserKey: key found: %v", accessKey))
				responseBody["id"] = accessKey
				responseBody["uid"] = uid
				responseBody["access_key"] = accessKey
				responseBody["secret_key"] = secretKey
				responseBody["tenant_id"] = tenantId
				responseBody["pgp_public_key"] = pgpKey
				responseBody["encrypted_secret_key"] = encryptedKey
				if enabled == "enabled" {
					responseBody["enabled"] = true
				} else {
					responseBody["enabled"] = false
				}
			}
		}
	}
	return FakeHttpResponse(response, responseBody)
}

func DeleteNonLocalUserKeyFunc(ctx context.Context, _client any, attr map[string]any, data map[string]any, headers map[string]string) (*http.Response, error) {
	client := _client.(*vast_client.VMSSession)
	attributes, err := getAttributesAsString([]string{"path", "id"}, attr)
	if err != nil {
		return nil, err
	}
	s := strings.SplitN((*attributes)["id"], "-", 2)
	uid, err := strconv.Atoi(s[0])
	if err != nil {
		return nil, err
	}
	accessKey := s[1]
	path := (*attributes)["path"]
	payload := map[string]any{"access_key": accessKey, "uid": uid}
	buffer, marshallingError := json.Marshal(payload)
	if marshallingError != nil {
		return nil, marshallingError
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling DELETE for %v", accessKey))
	return client.Delete(ctx, path, "", bytes.NewReader(buffer), headers)
}
