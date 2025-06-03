package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
)

func ProtectedPathDeleteFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	//First we call delete than we wait up to 10 min till the protected path is deleted
	response, err := DefaultDeleteFunc(ctx, _client, attr, data, headers)
	if err != nil {
		return response, err
	}
	//Now we wait for the protected path deletion
	client := _client.(*vast_client.VMSSession)
	attributes, err := getAttributesAsString([]string{"path", "id"}, attr)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("%v%v", (*attributes)["path"], (*attributes)["id"])
	deleted := false

	for i := time.Minute * 10; i > 0; i = (i - 10*time.Second) {
		tflog.Debug(ctx, fmt.Sprintf("Wating for protectedpath: %v to be deleted", path))
		time.Sleep(time.Second * 10)
		r, _ := client.Get(ctx, path, "", map[string]string{})
		switch r.StatusCode {
		case 404:
			tflog.Debug(ctx, fmt.Sprintf("Protected Path %v Deleted", path))
			deleted = true
			break
		default:
			tflog.Debug(ctx, fmt.Sprintf("Protected Path %v is still being deleted", path))

		}
		if deleted {
			break
		}

	}
	return response, nil

}

func ProtectedPathBeforePostFunc(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	FieldsUpdate(ctx, []string{"enabled"}, d, &m)
	return m, nil
}

func ProtectedPathBeforePatchFunc(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	FieldsUpdate(ctx, []string{"enabled"}, d, &m)
	return m, nil
}

//In case of creating sending enabled=false will have no affect, still this indicates that the user intended to so we patch the newly created protected_path

func ProtectedPathBeforeCreateFunc(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	enabled, exists := m["enabled"]
	if !exists {
		return m, nil
	}

	client := i.(*vast_client.VMSSession)
	tflog.Debug(ctx, fmt.Sprintf("[ProtectedPathBeforeCreateFunc] Setting the value of enabled to: %v ", enabled))
	id := fmt.Sprintf("%v", d.Id())
	z := map[string]interface{}{"enabled": enabled}
	b, _ := json.Marshal(z)
	_, err := client.Patch(ctx, GenPath(fmt.Sprintf("%v/%v", "protectedpaths", id)), bytes.NewReader(b), map[string]string{})
	if err != nil {
		return m, err
	}
	d.Set("enabled", enabled)
	return m, nil
}
