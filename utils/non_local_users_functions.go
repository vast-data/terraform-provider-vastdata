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
	"strings"
)

var nonLocalUserListsAttributes = []string{"s3_policies_ids"}
var nonLocalUserBooleanAttributes = []string{"allow_delete_bucket", "allow_create_bucket"}

func getNonLocalUserId(uid int, tenantId int) string {
	return fmt.Sprintf("%v-%v", uid, tenantId)
}

func getNonLocalUserUidAndTenantId(id string) (int, int, error) {
	split := strings.Split(id, "-")
	if len(split) != 2 {
		return 0, 0, fmt.Errorf("invalid NonLocalUser ID: %s", id)
	}
	uid, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, 0, err
	}
	tenantId, err := strconv.Atoi(split[1])
	if err != nil {
		return 0, 0, err
	}
	return uid, tenantId, nil
}

func NonLocalUserBeforePatchFunc(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	FieldsUpdate(ctx, nonLocalUserListsAttributes, d, &m)
	FieldsUpdate(ctx, nonLocalUserBooleanAttributes, d, &m)
	return m, nil
}

func NonLocalUserCreateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	client := _client.(vastclient.JwtSession)
	attributes, err := getAttributesAsString([]string{"path"}, attr)
	if err != nil {
		return nil, err
	}
	buffer, marshallingError := json.Marshal(data)
	if marshallingError != nil {
		return nil, marshallingError
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling POST to path \"%v\"", attr))
	response, err := client.Patch(ctx, (*attributes)["path"], "application/json", bytes.NewReader(buffer), map[string]string{})
	if err != nil {
		return nil, err
	}
	responseBody := map[string]interface{}{}
	err = UnmarshalBodyToMap(response, &responseBody)
	if err != nil {
		return nil, err
	}
	uid := data["uid"].(int)
	tenantId := data["tenant_id"].(int)
	responseBody["id"] = getNonLocalUserId(uid, tenantId)
	return FakeHttpResponse(response, responseBody)
}

func NonLocalUserGetFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	client := _client.(vastclient.JwtSession)
	attributes, err := getAttributesAsString([]string{"path"}, attr)
	if err != nil {
		return nil, err
	}
	path := (*attributes)["path"]
	uid := d.Get("uid").(int)
	tenantId := d.Get("tenant_id").(int)
	query := fmt.Sprintf("uid=%v&tenantId=%v", uid, tenantId)
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
	responseBody["tenant_id"] = tenantId // tenant_id is missing from the response
	responseBody["id"] = getNonLocalUserId(uid, tenantId)
	return FakeHttpResponse(response, responseBody)
}

func NonLocalUserUpdateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	id := attr["id"].(string)
	uid, tenantId, err := getNonLocalUserUidAndTenantId(id)
	if err != nil {
		return nil, err
	}
	data["uid"] = uid
	data["tenant_id"] = tenantId
	return NonLocalUserCreateFunc(ctx, _client, attr, data, headers)
}

func NonLocalUserDeleteFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	tflog.Info(ctx, "Doing nothing. We cannot delete non-local user.")
	return nil, nil
}
