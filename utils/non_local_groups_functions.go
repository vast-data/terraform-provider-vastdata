package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vastclient "github.com/vast-data/terraform-provider-vastdata/vast-client"
	"net/http"
	"strconv"
)

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
	responseBody["id"] = strconv.Itoa(gid)
	responseBody["sid"] = fmt.Sprintf("%v", responseBody["sid"]) // Force string (can be int)
	responseBody["groupname"] = responseBody["name"]
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
	query := fmt.Sprintf("gid=%v", gid)
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
	responseBody["id"] = strconv.Itoa(gid)
	responseBody["sid"] = fmt.Sprintf("%v", responseBody["sid"]) // Force string (can be int)
	responseBody["groupname"] = responseBody["name"]
	return FakeHttpResponse(response, responseBody)
}

func NonLocalGroupUpdateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	id := attr["id"].(string)
	gid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	data["gid"] = gid
	return NonLocalGroupCreateFunc(ctx, _client, attr, data, headers)
}

func NonLocalGroupDeleteFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	tflog.Info(ctx, "Doing nothing. We cannot delete non-local group.")
	return nil, nil
}
