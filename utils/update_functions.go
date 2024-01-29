package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
)

type HttpUpdateFunc func(context.Context, interface{}, string, string, map[string]string, *schema.ResourceData) (*http.Response, error)

func UpdateVastDatabaseSchemaName(ctx context.Context, m interface{}, path string, content_type string, headers map[string]string, d *schema.ResourceData) (*http.Response, error) {
	var new_data map[string]interface{} = map[string]interface{}{}
	client := m.(vast_client.JwtSession)
	new_data["database_name"] = d.Get("database_name")
	old, new := d.GetChange("name")
	new_data["name"] = old
	new_data["new_schema_name"] = new
	b, err := json.Marshal(new_data)
	if err != nil {
		return nil, err
	}
	return client.Patch(ctx, fmt.Sprintf("%s/rename", path), "content_type", bytes.NewReader(b), map[string]string{})
}
