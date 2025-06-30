package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
)

func S3PolicyReSendEnable(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	enabled := m["enabled"]
	client := i.(*vast_client.VMSSession)
	tflog.Debug(ctx, fmt.Sprintf("[S3PolicyReSendEnable] %v ", enabled))
	id := fmt.Sprintf("%v", d.Id())
	z := map[string]interface{}{"enabled": enabled}
	b, _ := json.Marshal(z)
	_, err := client.Patch(ctx, GenPath(fmt.Sprintf("%v/%v", "s3policies", id)), "", bytes.NewReader(b), map[string]string{})
	if err != nil {
		return m, err
	}
	d.Set("enabled", enabled)
	return m, nil
}
