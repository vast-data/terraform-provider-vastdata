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

func getNonLocalGroupId(gid int, tenantId int, context string) string {
	return fmt.Sprintf("%v-%v-%v", gid, tenantId, context)
}

func decomposeNonLocalGroupId(id string) (int, int, string, error) {
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

func NonLocalGroupCreateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	client := _client.(*vastclient.VMSSession)
	attributes, err := getAttributesAsString([]string{"path"}, attr)
	if err != nil {
		return nil, err
	}
	if _, ok := data["s3_policies_ids"]; !ok {
		data["s3_policies_ids"] = []int{}
	}
	buffer, marshallingError := json.Marshal(data)
	if marshallingError != nil {
		return nil, marshallingError
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling POST to path \"%v\"", attr))
	response, err := client.Patch(ctx, (*attributes)["path"], bytes.NewReader(buffer), map[string]string{})
	if err != nil {
		return response, err
	}
	responseBody := map[string]interface{}{}
	err = UnmarshalBodyToMap(response, &responseBody)
	if err != nil {
		return nil, err
	}
	gid := data["gid"].(int)
	tenantId := data["tenant_id"].(int)
	contextValue := data["context"].(string)
	responseBody["id"] = getNonLocalGroupId(gid, data["tenant_id"].(int), contextValue)
	responseBody["tenant_id"] = tenantId
	responseBody["sid"] = fmt.Sprintf("%v", responseBody["sid"]) // Force string (can be int)
	responseBody["groupname"] = responseBody["name"]
	responseBody["context"] = contextValue
	return FakeHttpResponse(response, responseBody)
}

func NonLocalGroupGetFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	client := _client.(*vastclient.VMSSession)
	attributes, err := getAttributesAsString([]string{"path"}, attr)
	if err != nil {
		return nil, err
	}
	path := (*attributes)["path"]
	gid := d.Get("gid").(int)
	tenantId := d.Get("tenant_id").(int)
	contextValue := d.Get("context").(string)
	query := fmt.Sprintf("gid=%v&tenant_id=%v&context=%v", gid, tenantId, contextValue)
	tflog.Debug(ctx, fmt.Sprintf("Calling GET to path \"%v\" , with Query %v", path, query))
	response, err := client.Get(ctx, path, query, headers)
	if err != nil {
		return nil, err
	}

	responseBody := map[string]interface{}{}
	err = UnmarshalBodyToMap(response, &responseBody)
	if err != nil {
		return nil, err
	}
	responseBody["id"] = getNonLocalGroupId(gid, tenantId, contextValue)
	responseBody["tenant_id"] = tenantId
	responseBody["sid"] = fmt.Sprintf("%v", responseBody["sid"]) // Force string (can be int)
	responseBody["groupname"] = responseBody["name"]
	responseBody["context"] = contextValue
	return FakeHttpResponse(response, responseBody)
}

func NonLocalGroupUpdateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	id := attr["id"].(string)
	gid, tenantId, contextValue, err := decomposeNonLocalGroupId(id)
	if err != nil {
		return nil, err
	}
	data["gid"] = gid
	data["tenant_id"] = tenantId
	data["context"] = contextValue
	return NonLocalGroupCreateFunc(ctx, _client, attr, data, headers)
}

func NonLocalGroupDeleteFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	tflog.Info(ctx, "Doing nothing. We cannot delete non-local group.")
	return nil, nil
}

func mimicListResponseForSingularNonLocalGroup(ctx context.Context, response *http.Response) (*http.Response, error) {
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
	gid := int((*unmarshalledBody)["gid"].(float64))
	tenantId, _ := strconv.Atoi(response.Request.URL.Query().Get("tenant_id"))
	contextValue := response.Request.URL.Query().Get("context")
	id := getNonLocalGroupId(gid, tenantId, contextValue)
	(*unmarshalledBody)["id"] = id
	(*unmarshalledBody)["tenant_id"] = tenantId
	(*unmarshalledBody)["groupname"] = (*unmarshalledBody)["name"]
	(*unmarshalledBody)["context"] = contextValue
	(*unmarshalledBody)["sid"] = fmt.Sprintf("%v", (*unmarshalledBody)["sid"]) // Force string (can be int)
	var list []*map[string]interface{}
	list = append(list, unmarshalledBody)
	return FakeHttpResponseAny(response, list)
}

func NonLocalGroupImportFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, g GetFuncType) (*http.Response, error) {
	defaultResponse, err := DefaultImportFunc(ctx, _client, attr, d, g)
	if err != nil {
		return nil, err
	}
	mimickedResponse, err := mimicListResponseForSingularNonLocalGroup(ctx, defaultResponse)
	if err != nil {
		return nil, err
	}
	return mimickedResponse, nil
}

func NonLocalGroupProcessingFunc(ctx context.Context, response *http.Response) ([]byte, error) {
	mimickedResponse, err := mimicListResponseForSingularNonLocalGroup(ctx, response)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(mimickedResponse.Body)
	tflog.Debug(ctx, fmt.Sprintf("HTTP Response body %s", string(body)))
	return body, err
}
