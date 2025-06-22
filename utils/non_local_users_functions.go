package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vastclient "github.com/vast-data/terraform-provider-vastdata/vast-client"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var nonLocalUserListsAttributes = []string{"s3_policies_ids"}
var nonLocalUserBooleanAttributes = []string{"allow_delete_bucket", "allow_create_bucket"}

func getNonLocalUserId(uid int, tenantId int, contextValue string) string {
	return fmt.Sprintf("%v-%v-%v", uid, tenantId, contextValue)
}

func decomposeNonLocalUserId(id string) (int, int, string, error) {
	split := strings.Split(id, "-")
	if len(split) != 3 {
		return 0, 0, "", fmt.Errorf("invalid NonLocalUser ID: %s", id)
	}
	uid, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, 0, "", err
	}
	tenantId, err := strconv.Atoi(split[1])
	if err != nil {
		return 0, 0, "", err
	}
	contextValue := split[2]
	return uid, tenantId, contextValue, nil
}

func NonLocalUserBeforePatchFunc(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	FieldsUpdate(ctx, nonLocalUserListsAttributes, d, &m)
	FieldsUpdate(ctx, nonLocalUserBooleanAttributes, d, &m)
	return m, nil
}

func NonLocalUserCreateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	client := _client.(*vastclient.VMSSession)
	attributes, err := getAttributesAsString([]string{"path"}, attr)
	if err != nil {
		return nil, err
	}
	buffer, marshallingError := json.Marshal(data)
	if marshallingError != nil {
		return nil, marshallingError
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling POST to path \"%v\"", attr))
	response, err := client.Patch(ctx, (*attributes)["path"], "", bytes.NewReader(buffer), map[string]string{})
	if err != nil {
		return response, err
	}
	unmarshalledBody := map[string]interface{}{}
	err = UnmarshalBodyToMap(response, &unmarshalledBody)
	if err != nil {
		return nil, err
	}
	uid := data["uid"].(int)
	tenantId := data["tenant_id"].(int)
	contextValue := data["context"].(string)
	id := getNonLocalUserId(uid, tenantId, contextValue)
	unmarshalledBody["id"] = id
	unmarshalledBody["tenant_id"] = tenantId
	unmarshalledBody["username"] = unmarshalledBody["name"]
	unmarshalledBody["context"] = contextValue
	return FakeHttpResponse(response, unmarshalledBody)
}

func NonLocalUserGetFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	client := _client.(*vastclient.VMSSession)
	attributes, err := getAttributesAsString([]string{"path"}, attr)
	if err != nil {
		return nil, err
	}
	path := (*attributes)["path"]
	uid := d.Get("uid").(int)
	tenantId := d.Get("tenant_id").(int)
	contextValue := d.Get("context").(string)
	query := fmt.Sprintf("uid=%v&tenant_id=%v", uid, tenantId)
	if contextValue != "" {
		query = fmt.Sprintf("%v&context=%v", query, contextValue)
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling GET to path \"%v\" , with Query %v", path, query))
	response, err := client.Get(ctx, path, query, headers)
	if err != nil {
		return nil, err
	}

	unmarshalledBody := map[string]interface{}{}
	err = UnmarshalBodyToMap(response, &unmarshalledBody)
	if err != nil {
		return nil, err
	}
	id := getNonLocalUserId(uid, tenantId, contextValue)
	unmarshalledBody["id"] = id
	unmarshalledBody["tenant_id"] = tenantId
	unmarshalledBody["username"] = unmarshalledBody["name"]
	unmarshalledBody["context"] = contextValue
	return FakeHttpResponse(response, unmarshalledBody)
}

func NonLocalUserUpdateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	id := attr["id"].(string)
	uid, tenantId, contextValue, err := decomposeNonLocalUserId(id)
	if err != nil {
		return nil, err
	}
	data["uid"] = uid
	data["tenant_id"] = tenantId
	data["context"] = contextValue
	return NonLocalUserCreateFunc(ctx, _client, attr, data, headers)
}

func NonLocalUserDeleteFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	tflog.Info(ctx, "Doing nothing. We cannot delete non-local user.")
	return nil, nil
}

func mimicListResponseForSingularNonLocalUser(ctx context.Context, response *http.Response) (*http.Response, error) {
	unmarshalledBody := new(map[string]interface{})
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, unmarshalledBody)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Resonse From Cluster %v", string(body)))
		return nil, err
	}
	uid := int((*unmarshalledBody)["uid"].(float64))
	tenantId, _ := strconv.Atoi(response.Request.URL.Query().Get("tenant_id"))
	contextValue := response.Request.URL.Query().Get("context")
	id := getNonLocalUserId(uid, tenantId, contextValue)
	(*unmarshalledBody)["id"] = id
	(*unmarshalledBody)["tenant_id"] = tenantId
	(*unmarshalledBody)["username"] = (*unmarshalledBody)["name"]
	(*unmarshalledBody)["context"] = contextValue
	var list []*map[string]interface{}
	list = append(list, unmarshalledBody)
	return FakeHttpResponseAny(response, list)
}

func NonLocalUserImportFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, g GetFuncType) (*http.Response, error) {
	defaultResponse, err := DefaultImportFunc(ctx, _client, attr, d, g)
	if err != nil {
		return nil, err
	}
	mimickedResponse, err := mimicListResponseForSingularNonLocalUser(ctx, defaultResponse)
	if err != nil {
		return nil, err
	}
	return mimickedResponse, nil
}

func NonLocalUserProcessingFunc(ctx context.Context, response *http.Response) ([]byte, error) {
	mimickedResponse, err := mimicListResponseForSingularNonLocalUser(ctx, response)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(mimickedResponse.Body)
	tflog.Debug(ctx, fmt.Sprintf("HTTP Response body %s", string(body)))
	return body, err
}
