package utils

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
)

type HttpReadFunc func(context.Context, interface{}, string, string, map[string]string, *schema.ResourceData) (*http.Response, error)

func ReadVastDataDatabseSchema(ctx context.Context, m interface{}, path string, query string, headers map[string]string, d *schema.ResourceData) (*http.Response, error) {
	client := m.(vast_client.JwtSession)
	qs := url.Values{}
	qs.Add("database_name", fmt.Sprintf("%v", d.Get("database_name")))
	qs.Add("schema", fmt.Sprintf("%v", d.Get("name")))
	return client.Get(ctx, path, qs.Encode(), headers)
}

func ReadVastDataDatabseTable(ctx context.Context, m interface{}, path string, query string, headers map[string]string, d *schema.ResourceData) (*http.Response, error) {
	schema_identifier := fmt.Sprintf("%v", d.Get("schema_identifier"))
	s := strings.SplitN(schema_identifier, "/", 2)
	database_name := s[0]
	schema_name := s[1]
	client := m.(vast_client.JwtSession)
	qs := url.Values{}

	qs.Add("table_name", fmt.Sprintf("%v", d.Get("name")))
	qs.Add("database_name", database_name)
	qs.Add("schema_name", schema_name)
	return client.Get(ctx, "/api/latest/columns/", qs.Encode(), headers)
}
