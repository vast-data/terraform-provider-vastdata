package utils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CloneTypeLocal(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	clone_type, exists := d.GetOkExists("clone_type")
	tflog.Debug(ctx, fmt.Sprintf("[CloneTypeLocal] - Clone Type: %, Exists: %v", clone_type, exists))
	if !exists {
		return m, nil
	}
	if fmt.Sprintf("%v", clone_type) == "LOCAL" {
		tflog.Debug(ctx, "[CloneTypeLocal] Clone Type is local , this means we should remove keep-remote from every frame if exists")
		frames, exists := m["frames"]
		if exists {
			for _, fr := range frames.([]interface{}) {
				q := fr.(map[string]interface{})
				_, exists := q["keep-remote"]
				if exists {
					delete(q, "keep-remote")
				}
			}
		}
	}
	return m, nil
}
