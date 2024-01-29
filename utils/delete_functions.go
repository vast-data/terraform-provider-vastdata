package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
)

type HttpDeleteFunc func(context.Context, interface{}, string, string, map[string]string, *schema.ResourceData) (*http.Response, error)

func DeleteVastDatabaseSchema(ctx context.Context, m interface{}, path string, query string, headers map[string]string, d *schema.ResourceData) (*http.Response, error) {
	var data map[string]interface{} = map[string]interface{}{}
	client := m.(vast_client.JwtSession)
	data["database_name"] = d.Get("database_name")
	data["name"] = d.Get("name")
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return client.Delete(ctx, path, "", bytes.NewReader(b), headers)

}

func generateTableDeleteData(ctx context.Context, d *schema.ResourceData) ([]byte, error) {
	data := map[string]interface{}{}

	s := strings.SplitN(fmt.Sprintf("%v", d.Get("schema_identifier")), "/", 2)
	data["database_name"] = s[0]
	data["schema_name"] = s[1]
	data["name"] = d.Get("name")

	return json.MarshalIndent(data, "", " ")
}

func DeleteVastDatabaseTable(ctx context.Context, m interface{}, path string, query string, headers map[string]string, d *schema.ResourceData) (*http.Response, error) {
	data, err := generateTableDeleteData(ctx, d)
	if err != nil {
		return nil, err
	}
	client := m.(vast_client.JwtSession)
	return client.Delete(ctx, path, query, bytes.NewReader(data), headers)

}
