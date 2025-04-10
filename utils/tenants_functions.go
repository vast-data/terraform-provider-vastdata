package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
)

var tenant_booleans []string = []string{"use_smb_privileged_user", "use_smb_privileged_group", "smb_privileged_group_full_access", "is_nfsv42_supported", "allow_locked_users", "allow_disabled_users"}
var tenant_lists []string = []string{"vippool_ids"}

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
	id := d.Id()
	u := url.Values{"tenant_id": {id}}
	client := i.(*vast_client.VMSSession)
	h, err := client.Get(ctx, GenPath("/vippools"), u.Encode(), map[string]string{})
	if err != nil {
		return err
	}
	response_body, _ := io.ReadAll(h.Body)
	g := []GenericInt64ID{}
	err = json.Unmarshal(response_body, &g)
	if err != nil {
		return err
	}
	for _, i := range g {
		vippool_ids = append(vippool_ids, i.Id)
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
	FieldsUpdate(ctx, tenant_lists, d, &m)
	return m, nil
}
