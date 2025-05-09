package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
)

type ResponseConversionFunc func(map[string]interface{}, interface{}, context.Context, *schema.ResourceData) (map[string]interface{}, error)

func EntityMergeToUserQuotas(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	/*This function should handle the case of the Quota object where sending is defferant than reading sturctue
	  to move the fields from the entity object into the user quotas
	*/
	for _, key := range []string{"user_quotas", "group_quotas"} {
		quotas, exists := m[key]
		if exists {
			old_quotas := quotas.([]interface{})
			new_quotas := []map[string]interface{}{}
			for _, quota := range old_quotas {
				new_quota := make(map[string]interface{})
				_quota := quota.(map[string]interface{})
				entity, entity_exists := _quota["entity"]
				if entity_exists {
					for k, v := range entity.(map[string]interface{}) {
						new_quota[k] = v
					}
				}
				for k, v := range _quota {
					if k == "entity" {
						continue
					}
					new_quota[k] = v
				}

				new_quotas = append(new_quotas, new_quota)
			}
			m[key] = new_quotas
		}
	}
	return m, nil
}

func EnabledMustBeSet(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	m["enabled"] = d.Get("enabled")
	return m, nil
}

func list_snapshoted_paths_remote(remote_target_guid, remote_tenant_guid string, client interface{}, ctx context.Context) ([]map[string]interface{}, error) {
	m := []map[string]interface{}{}
	values := url.Values{}
	if remote_target_guid != "" {
		values.Add("remote_target_guid", remote_target_guid)
	}
	if remote_tenant_guid != "" {
		values.Add("remote_tenant_guid", remote_tenant_guid)
	}
	c := client.(*vast_client.VMSSession)

	r, err := c.Get(ctx, "/api/latest/clusters/list_snapshoted_paths_remote/", values.Encode(), map[string]string{})
	if err != nil {
		return m, err
	}
	b, e := io.ReadAll(r.Body)

	if e != nil {
		return m, err
	}

	json.Unmarshal(b, &m)
	tflog.Debug(ctx, fmt.Sprintf("Paths found remotly %v", m))
	return m, nil

}

func list_clone_snapshoted_paths_remote(remote_target_guid, handle string, client interface{}, ctx context.Context) ([]map[string]interface{}, error) {
	m := []map[string]interface{}{}
	values := url.Values{}
	if remote_target_guid != "" {
		values.Add("remote_target_guid", remote_target_guid)
	}
	if handle != "" {
		values.Add("handle", handle)
	}
	values.Add("start_snapshot_id", "0")
	c := client.(*vast_client.VMSSession)

	r, err := c.Get(ctx, "/api/latest/clusters/list_clone_snapshoted_paths_remote/", values.Encode(), map[string]string{})
	if err != nil {
		return m, err
	}
	b, e := io.ReadAll(r.Body)

	if e != nil {
		return m, err
	}

	json.Unmarshal(b, &m)
	return m, nil

}
func get_snapshot_handle(remote_tenant_guid, remote_target_guid, path string, client interface{}, ctx context.Context) (string, error) {
	m, err := list_snapshoted_paths_remote(remote_target_guid, remote_tenant_guid, client, ctx)
	if err != nil {
		return "", err
	}
	for _, o := range m {
		k, v := o["name"]
		if !v {
			continue
		}
		if strings.TrimSuffix(k.(string), "/") == strings.TrimSuffix(path, "/") {
			h, e := o["handle"]
			if !e {
				return "", errors.New(fmt.Sprintf("Could not find handle at %v", o))
			}
			return fmt.Sprintf("%v", h), nil
		}
	}
	return "", errors.New("Could not find path")
}

func get_snapshot_clone_id(handle, remote_target_guid, snapshot_name string, client interface{}, ctx context.Context) (interface{}, error) {
	m, err := list_clone_snapshoted_paths_remote(remote_target_guid, handle, client, ctx)
	if err != nil {
		return 0, err
	}
	for _, v := range m {
		h, e := v["name"]
		if !e {
			continue
		}
		if fmt.Sprintf("%v", h) == snapshot_name {
			clone_id, e := v["clone_id"]
			if !e {
				continue
			}
			return clone_id, nil
		}
	}
	return uint64(0), errors.New(fmt.Sprintf("Could not find a snapshot with the name %s for the handle %s", snapshot_name, handle))
}

func AddStreamInfo(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	tflog.Debug(ctx, fmt.Sprintf("Data Before Processing: %v ", m))
	remote_tenant_guid := m["owner_tenant"].(map[string]interface{})["guid"].(string)
	remote_target_guid := m["remote_target_guid"].(string)
	remote_target_path := m["remote_target_path"].(string)
	snapshot_name := m["owner_root_snapshot"].(map[string]interface{})["name"].(string)
	client := i.(*vast_client.VMSSession)
	handle, err := get_snapshot_handle(remote_tenant_guid, remote_target_guid, remote_target_path, client, ctx)
	if err != nil {
		return m, err
	}
	clone_id, err := get_snapshot_clone_id(handle, remote_target_guid, snapshot_name, client, ctx)
	if err != nil {
		return m, err
	}
	m["owner_root_snapshot"].(map[string]interface{})["clone_id"] = clone_id
	m["owner_root_snapshot"].(map[string]interface{})["parent_handle_ehandle"] = strings.Split(handle, "_")[0]

	delete(m, "remote_target_path")
	delete(m, "remote_target_guid")

	return m, nil
}

func UpdateStreamInfo(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	//The only update possiable is enable/disable
	v, exists := m["enabled"]
	if exists {
		return map[string]interface{}{"enabled": v}, nil
	}
	return map[string]interface{}{"enabled": false}, nil

}

func has_bucket_logging(ctx context.Context, d *schema.ResourceData) bool {
	l, exists := d.GetOkExists("bucket_logging")
	tflog.Debug(ctx, fmt.Sprintf("[has_bucket_logging] bucket_logging %v", l))
	if !exists {
		tflog.Debug(ctx, "bucket_logging is not defined")
		return false
	}
	o, is_list := l.([]interface{})
	if is_list && len(o) > 0 {
		m, is_map := o[0].(map[string]interface{})
		if is_map {
			dest_id, has_dest_id := m["destination_id"]
			if has_dest_id && fmt.Sprintf("%v", dest_id) != "0" {
				return true
			} else {
				tflog.Debug(ctx, fmt.Sprintf("bucket_logging first array element , does not have a key called destination_id or it's value is 0 %v", m))
				return false
			}
		} else {
			tflog.Debug(ctx, fmt.Sprintf("bucket_logging first element is not a map from the type of map[string]interface{}, %v", l))
			return false
		}
	} else {
		tflog.Debug(ctx, fmt.Sprintf("bucket_logging is not from the type of []interface{} but %v, or has 0 length", l))
		return false

	}
	return false
}

func AlwaysSendCreateDir(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	create_dir := d.Get("create_dir")
	m["create_dir"] = create_dir
	//Due to some VMS issue it return a broken configurations of bucket logging , this causes the update to fail
	//in the case of shared ACL set , but the acl is missing ,we must set it to an empty list
	bucket_logging_configured := has_bucket_logging(ctx, d)
	tflog.Debug(ctx, fmt.Sprintf("[AlwaysSendCreateDir] has_bucket_logging bucket: %v Bucket logging configured: %v", d.Get("bucket"), bucket_logging_configured))
	if !bucket_logging_configured {
		_, e := m["bucket_logging"]
		if e {
			tflog.Debug(ctx, fmt.Sprintf("[AlwaysSendCreateDir] bucket_logging is not configured setting it to nil"))
			m["bucket_logging"] = nil
		}

	}
	share_acl, share_acl_exists := m["share_acl"]
	if share_acl_exists {
		qos_policy_id, qos_policy_id_exists := d.GetOkExists("qos_policy_id")
		if qos_policy_id_exists && qos_policy_id != 0 {
			m["qos_policy_id"] = qos_policy_id
		} else {
			m["qos_policy_id"] = nil
		}

		_share_acl := share_acl.(map[string]interface{})
		_, acl_exists := _share_acl["acl"]
		if !acl_exists {
			o, n := d.GetChange("enabled")
			tflog.Debug(ctx, fmt.Sprintf("VIEW: share_acl->acl doe not exists, creating empty acl with enabled value of: Old %v, New: %v", o, n))
			_share_acl["acl"] = []interface{}{}
			if d.Get("enabled") == nil {
				_share_acl["enabled"] = false
			}
		}
	}
	return m, nil
}
func AlwaysStoreCreateDir(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	_create_dir, exists := d.GetOkExists("create_dir")
	o, n := d.GetChange("create_dir")
	tflog.Debug(ctx, fmt.Sprintf("The Value of Create Dir Obtained:  Exists(%v),Value(%v),Changed(%v),Old(%v),New(%v)", exists, _create_dir, d.HasChange("create_dir"), o, n))
	if !exists {
		d.Set("create_dir", false)
		tflog.Debug(ctx, fmt.Sprintf("CREATE DIR: Was not found and set to false"))
	}

	d.Set("create_dir", n)
	tflog.Debug(ctx, fmt.Sprintf("CREATE DIR: Value Found %v", _create_dir))

	return m, nil
}

type SchemaManipulationFunc func(interface{}, context.Context, *schema.ResourceData) error

func KeepCreateDirState(i interface{}, ctx context.Context, d *schema.ResourceData) error {
	_, exists := d.GetOkExists("create_dir")
	o, n := d.GetChange("create_dir")
	if has_bucket_logging(ctx, d) {

	}
	tflog.Debug(ctx, fmt.Sprintf("OLD: %v, NEW: %v", o, n))
	if !exists {
		return nil
	}
	if !d.HasChange("create_dir") {
		d.Set("create_dir", n)
	}
	return nil
}

type PreDeleteFunc func(context.Context, *schema.ResourceData, interface{}) (io.Reader, error)

func AlwaysSkipDeleteLdap(ctx context.Context, d *schema.ResourceData, m interface{}) (io.Reader, error) {
	data := map[string]interface{}{"skip_ldap": true}
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
