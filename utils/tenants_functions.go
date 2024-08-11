package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
)

var tenant_booleans []string = []string{"use_smb_privileged_user", "use_smb_privileged_group", "smb_privileged_group_full_access", "is_nfsv42_supported", "allow_locked_users", "allow_disabled_users"}

type GenericInt64ID struct {
	Id int64 `json:"id,omitempty"`
}

func ConvertVippoolToIDs(i interface{}, ctx context.Context, d *schema.ResourceData) error {
	v, e := d.GetOkExists("vippool_ids")
	tflog.Debug(ctx, fmt.Sprintf("[ConvertVippoolToIDs] - Tenant %v , vipool_id: %v", d.Get("name"), v))

	vippool_ids := []int64{}
	if e && (len(v.([]interface{})) > 0) {
		return nil
	}
	vippool_names, vippool_names_exists := d.GetOkExists("vippool_names")
	tflog.Debug(ctx, fmt.Sprintf("[ConvertVippoolToIDs] - Tenant %v , vippool_names: %v", d.Get("name"), vippool_names))
	if !vippool_names_exists {
		return nil
	}
	client := i.(vast_client.JwtSession)
	for _, vp := range vippool_names.([]interface{}) {
		vps := fmt.Sprintf("name=%v", vp)
		tflog.Debug(ctx, fmt.Sprintf("[ConvertVippoolToIDs] - Tenant %v , converting name: %v to id", d.Get("name"), vps))
		h, err := client.Get(ctx, GenPath("/vippools"), vps, map[string]string{})
		if err != nil {
			continue
		}
		response_body, _ := io.ReadAll(h.Body)
		g := []GenericInt64ID{}
		err = json.Unmarshal(response_body, &g)
		if err != nil {
			continue
		}
		if len(g) != 0 {
			vippool_ids = append(vippool_ids, g[0].Id)
		}

	}
	d.Set("vippool_ids", vippool_ids)
	return nil
}

func TenantBeforePostFunc(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	FieldsUpdate(ctx, tenant_booleans, d, &m)
	return m, nil
}

func TenantBeforePatchFunc(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	FieldsUpdate(ctx, tenant_booleans, d, &m)
	return m, nil
}
