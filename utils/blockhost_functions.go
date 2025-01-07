package utils

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func BlockHostCreateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	data, err := listToKVmap(ctx, data, "blockhost_tags")
	if err != nil {
		return nil, err
	}
	return DefaultCreateFunc(ctx, _client, attr, data, headers)
}

func BlockHostUpdateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	data, err := listToKVmap(ctx, data, "blockhost_tags")
	if err != nil {
		return nil, err
	}
	block_host_ids, exists := d.GetOkExists("block_host_ids")
	if !exists {
		data["block_host_ids"] = []int64{}
	} else {
		data["block_host_ids"] = block_host_ids
	}

	return DefaultUpdateFunc(ctx, _client, attr, data, d, headers)
}

func BlockHostGetFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	response, err := DefaultGetFunc(ctx, _client, attr, d, headers)
	if err != nil {
		return response, err
	}
	data_with_tags, tags_err := kvMapToList(ctx, response, "blockhost_tags")
	if tags_err != nil {
		return response, tags_err
	}
	tflog.Debug(ctx, fmt.Sprintf("[BlockHostGetFunc] Data With Tags Returned %v", data_with_tags))
	return FakeHttpResponse(response, data_with_tags)

}
