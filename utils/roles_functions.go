package utils

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var builin_permissions = map[string]interface{}{"create_support": nil, "create_settings": nil, "create_security": nil, "create_monitoring": nil, "create_logical": nil, "create_hardware": nil, "create_events": nil, "create_database": nil, "create_applications": nil, "view_support": nil, "view_settings": nil, "view_security": nil, "view_monitoring": nil, "view_logical": nil, "view_hardware": nil, "view_events": nil, "view_applications": nil, "view_database": nil, "edit_support": nil, "edit_settings": nil, "edit_security": nil, "edit_monitoring": nil, "edit_logical": nil, "edit_hardware": nil, "edit_events": nil, "edit_database": nil, "edit_applications": nil, "delete_support": nil, "delete_settings": nil, "delete_security": nil, "delete_monitoring": nil, "delete_logical": nil, "delete_hardware": nil, "delete_events": nil, "delete_applications": nil, "delete_database": nil}

func realms_permissions_to_list(m interface{}) []interface{} {
	s := []interface{}{}
	l, is_list := m.([]interface{})
	if !is_list {
		return s
	}
	for _, r := range l {
		o, is_map := r.(map[string]interface{})
		if !is_map {
			continue
		}
		_realm_name, exists := o["realm_name"]
		if !exists {
			continue
		}
		realm_name := fmt.Sprintf("%v", _realm_name)
		for _, x := range []string{"create", "view", "realm_name", "edit"} {
			v, exists := o[x]
			if !exists {
				continue
			}
			b, e := v.(bool)
			if !e {
				continue
			}
			if b {
				s = append(s, fmt.Sprintf("%v_%v", x, realm_name))
			}
		}

	}
	return s
}

func removeRealmsPermissions(ctx context.Context, i interface{}) []interface{} {
	p := []interface{}{}
	l, is_list := i.([]interface{})
	if !is_list {
		tflog.Debug(ctx, fmt.Sprintf("[removeRealmsPermissions] the value give %v given is not a list", i))
		return p
	}
	for _, v := range l {
		_, e := builin_permissions[fmt.Sprintf("%v", v)]
		if e {
			p = append(p, e)
			tflog.Debug(ctx, fmt.Sprintf("[removeRealmsPermissions] permissions %v which belongs to builtin permissions found and will be added", e))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("[removeRealmsPermissions] permissions %v which does not belongs to builtin permissions (probably Realm) found and will be added", e))
		}
	}
	return p

}

func addRealmPermissions(ctx context.Context, data map[string]interface{}) map[string]interface{} {
	realms_permissions, exists := data["realms_permissions"]

	if exists {
		realms_permissions_list := realms_permissions_to_list(realms_permissions)
		if len(realms_permissions_list) > 0 {

			tflog.Debug(ctx, fmt.Sprintf("[RoleCreateFunc] Realms permissions found %v", realms_permissions))
			permissions_list, e := data["permissions_list"]
			if !e {
				data["permissions_list"] = realms_permissions_list
			}
			_permissions_list, e := permissions_list.([]interface{})
			if e {
				_permissions_list = append(_permissions_list, realms_permissions_list...)
				data["permissions_list"] = _permissions_list
				tflog.Debug(ctx, fmt.Sprintf("[RoleCreateFunc] Realms permissions appended %v", _permissions_list))
			}
		}
	}
	return data

}

func readReamlsPermissions(ctx context.Context, permissions interface{}) []map[string]interface{} {
	realms_permissions := map[string]map[string]interface{}{}
	l, is_list := permissions.([]interface{})
	if !is_list {
		return []map[string]interface{}{}
	}
	for _, k := range l {
		_, e := builin_permissions[fmt.Sprintf("%v", k)]
		if e { //this is the case of built in permissions which means that this is not related to a realm permissions
			continue
		}
		realm_permissions := strings.SplitN(fmt.Sprintf("%v", k), "_", 2)
		tflog.Debug(ctx, fmt.Sprintf("[readReamlsPermissions] Checking realm permission %v", realm_permissions))
		realm := realm_permissions[1]
		permisson := realm_permissions[0]
		q, e := realms_permissions[realm]
		if !e {
			realms_permissions[realm] = map[string]interface{}{permisson: true, "realm_name": realm}
		} else {
			q[permisson] = true
		}

	}
	//here we fixthe realms permissions to have all the permissions sed
	for _, v := range realms_permissions {
		for _, i := range []string{"create", "delete", "edit", "view"} {
			_, e := v[i]
			if !e {
				v[i] = false
			}
		}
	}
	t := []map[string]interface{}{}
	for _, v := range realms_permissions {
		t = append(t, v)
	}
	return t
}

func RoleCreateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	tflog.Debug(ctx, fmt.Sprintf("[RoleCreateFunc] Data recive ", data))
	data = addRealmPermissions(ctx, data)
	return DefaultCreateFunc(ctx, _client, attr, data, headers)
}

func RoleUpdateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	tflog.Debug(ctx, fmt.Sprintf("[RoleUpdateFunc] Data recive %v , data map %v", d.Get("realms_permissions"), data))
	l, e := d.GetOkExists("permissions_list")
	if e {
		p, z := l.([]interface{})
		if z {
			data["permissions_list"] = p
		}
	}
	a, b := d.GetOkExists("realms_permissions")
	if b {
		p, z := a.([]interface{})
		if z {
			data["realms_permissions"] = p
		}
	}
	data = addRealmPermissions(ctx, data)
	tflog.Debug(ctx, fmt.Sprintf("[RoleUpdateFunc] Updated recive %v ", data))
	return DefaultUpdateFunc(ctx, _client, attr, data, d, headers)
}

func RoleAfterReadFunc(i interface{}, ctx context.Context, d *schema.ResourceData) error {
	permissions, exists := d.GetOkExists("permissions")
	tflog.Debug(ctx, fmt.Sprintf("[RoleAfterReadFunc] permissions: %v, permissions_list: %v", permissions, d.Get("permissions_list")))
	if exists {
		d.Set("permissions_list", removeRealmsPermissions(ctx, permissions))
		realm_permissions := readReamlsPermissions(ctx, permissions)
		tflog.Debug(ctx, fmt.Sprintf("[RoleAfterReadFunc] Data recive %v , permissions %v, permissions_list %v, has change, %v", realm_permissions, permissions, d.Get("permissions_list"), d.HasChange("permissions_list")))
		if len(realm_permissions) > 0 {
			d.Set("realms_permissions", realm_permissions)
		}
	}

	return nil

}
