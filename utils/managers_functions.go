package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func filterNativePermissions(ctx context.Context, i interface{}) []interface{} {
	p := []interface{}{}
	l, is_list := i.([]interface{})
	if !is_list {
		tflog.Debug(ctx, fmt.Sprintf("[filterNativePermissions] the value give %v given is not a list", i))
		return p
	}
	for _, v := range l {
		_, e := builin_permissions[fmt.Sprintf("%v", v)]
		if e {
			p = append(p, v)
			tflog.Debug(ctx, fmt.Sprintf("[filterNativePermissions] permissions %v which belongs to builtin permissions found and will be added", v))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("[filterNativePermissions] permissions %v which does not belongs to builtin permissions (probably Realm) found and will not be added", v))
		}
	}
	return p

}

func collectPermissions(ctx context.Context, r *http.Response, m map[string]interface{}) (map[string]interface{}, error) {
	permissions, permissions_exists := m["permissions"]
	if permissions_exists {
		p, is_map := permissions.(map[string]interface{})
		if !is_map {
			return m, fmt.Errorf("[collectPermissions] permissions return from VMS are not at map format")
		}
		manager, has_manager := p["manager"]
		if has_manager {
			m["permissions_list"] = manager
			tflog.Debug(ctx, fmt.Sprintf("[collectPermissions] permissions.manager found %v and will be used to replace permissions_list", manager))
		}
		delete(m, "permissions")
	}
	return m, nil
}

func convertManagerResponseToTfModel(ctx context.Context, r *http.Response, m map[string]interface{}) (map[string]interface{}, error) {
	b := map[string]interface{}{}
	e := UnmarshelBodyToMap(r, &b)
	if e != nil {
		return b, e
	}
	//Get Role IDs
	roles_ids := []interface{}{}
	roles, roles_exists := b["roles"]
	if roles_exists {
		tflog.Debug(ctx, fmt.Sprintf("[convertManagerResponseToTfModel] Roles returned from VMS %v", roles))
		roles_list, is_list := roles.([]interface{})
		if !is_list {
			return b, fmt.Errorf("[convertManagerResponseToTfModel] Roles return from VMS are not at list format")
		}
		for _, r := range roles_list {
			q, is_map := r.(map[string]interface{})
			if !is_map {
				continue
			}
			id, id_exists := q["id"]
			if id_exists {
				roles_ids = append(roles_ids, id)
			}
		}
	}
	tflog.Debug(ctx, fmt.Sprintf("[convertManagerResponseToTfModel] Roles IDs found %v", roles_ids))

	b["roles"] = roles_ids
	//Get Permissions
	permissions, permissions_exists := b["permissions"]
	if permissions_exists {
		p, is_map := permissions.(map[string]interface{})
		if !is_map {
			return b, fmt.Errorf("[convertManagerResponseToTfModel] permissions return from VMS are not at map format")
		}
		role, has_role := p["role"]
		if has_role {
			b["permissions"] = role
			tflog.Debug(ctx, fmt.Sprintf("[convertManagerResponseToTfModel] permissions.role found %v and will be used to replace permissions", role))
		}
		manager, has_manager := p["manager"]
		if has_manager {
			b["permissions_list"] = manager
			tflog.Debug(ctx, fmt.Sprintf("[convertManagerResponseToTfModel] permissions.manager found %v and will be used to replace permissions_list", role))
		}

	}
	b["password"] = m["password"]
	tflog.Debug(ctx, fmt.Sprintf("[convertManagerResponseToTfModel] Data that will be returned %v", b))
	return b, nil
}

func ManagerCreateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	tflog.Debug(ctx, fmt.Sprintf("[ManagerCreateFunc] Creating manager with data %v", data))
	r, e := DefaultCreateFunc(ctx, _client, attr, data, headers)
	if e != nil {
		return nil, e
	}

	b, e := collectRolesIds(ctx, r, data)
	if e != nil {
		return r, e
	}
	b, e = collectPermissions(ctx, r, b)
	if e != nil {
		return r, e
	}
	tflog.Debug(ctx, fmt.Sprintf("[ManagerCreateFunc] Converted respose data  %v", b))
	return FakeHttpResponse(r, b)

}

func ManagerUpdateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	tflog.Debug(ctx, fmt.Sprintf("[ManagerUpdateFunc] Update Func Called with data:%v", data))
	l, e := d.GetOkExists("permissions_list")
	if e {
		p, z := l.([]interface{})
		if z {
			data["permissions_list"] = p
		}
	}
	r, err := DefaultUpdateFunc(ctx, _client, attr, data, d, headers)
	if err != nil {
		return r, err
	}
	return FakeHttpResponse(r, data)

}

// func ManagersGetFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
// 	tflog.Debug(ctx, fmt.Sprintf("[ManagersGetFunc] Calling Default Get Func"))
// 	r, e := DefaultGetFunc(ctx, _client, attr, d, headers)
// 	if e != nil {
// 		return r, e
// 	}
// 	m := map[string]interface{}{}
// 	p := fmt.Sprintf("%v", d.Get("password"))
// 	m["password"] = p
// 	b, e := convertManagerResponseToTfModel(ctx, r, m)
// 	if e != nil {
// 		return r, e
// 	}
// 	j, has_permissions_list := b["permissions_list"]
// 	var realms_permissions []map[string]interface{}
// 	if has_permissions_list {
// 		tflog.Debug(ctx, fmt.Sprintf("[ManagersGetFunc] Processing permissions list to collect permissions list: %v", j))
// 		realms_permissions = readReamlsPermissions(ctx, j)

// 	} else {

// 		realms_permissions = []map[string]interface{}{}
// 	}
// 	b["realms_permissions"] = realms_permissions
// 	tflog.Debug(ctx, fmt.Sprintf("[ManagersGetFunc] Removing Realms permissions from permissions list: %v", j))
// 	tflog.Debug(ctx, fmt.Sprintf("[ManagersGetFunc] Realms permissions removed from permissions list: %v", filterNativePermissions(ctx, j)))
// 	b["permissions_list"] = filterNativePermissions(ctx, j)
// 	tflog.Debug(ctx, fmt.Sprintf("[ManagersGetFunc] Realms permissions collected: %v", b["realms_permissions"]))
// 	t, _ := json.MarshalIndent(b, "", "   ")
// 	tflog.Debug(ctx, fmt.Sprintf("[ManagersGetFunc] Manager returned: %v", string(t)))

// 	return FakeHttpResponse(r, b)
// }

func collectRolesIds(ctx context.Context, r *http.Response, m map[string]interface{}) (map[string]interface{}, error) {
	b := map[string]interface{}{}
	e := UnmarshelBodyToMap(r, &b)
	if e != nil {
		return b, e
	}
	//Get Role IDs
	roles_ids := []interface{}{}
	roles, roles_exists := b["roles"]
	if roles_exists {
		tflog.Debug(ctx, fmt.Sprintf("[convertManagerResponseToTfModel] Roles returned from VMS %v", roles))
		roles_list, is_list := roles.([]interface{})
		if !is_list {
			return b, fmt.Errorf("[convertManagerResponseToTfModel] Roles return from VMS are not at list format")
		}
		for _, r := range roles_list {
			q, is_map := r.(map[string]interface{})
			if !is_map {
				continue
			}
			id, id_exists := q["id"]
			if id_exists {
				roles_ids = append(roles_ids, id)
			}
		}
	}
	tflog.Debug(ctx, fmt.Sprintf("[collectRolesIds] Roles IDs found %v", roles_ids))
	b["roles"] = roles_ids
	return b, nil
}

func ManagersGetFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	tflog.Debug(ctx, fmt.Sprintf("[ManagersGetFunc] Calling Default Get Func"))
	r, e := DefaultGetFunc(ctx, _client, attr, d, headers)
	if e != nil {
		return r, e
	}
	m := map[string]interface{}{}
	p := fmt.Sprintf("%v", d.Get("password"))
	m["password"] = p
	b, e := collectRolesIds(ctx, r, m)
	if e != nil {
		return r, e
	}

	b, e = collectPermissions(ctx, r, b)
	if e != nil {
		return r, e
	}
	t, _ := json.MarshalIndent(b, "", "   ")
	tflog.Debug(ctx, fmt.Sprintf("[ManagersGetFunc] Manager returned: %v", string(t)))
	return FakeHttpResponse(r, b)

}

func ManagerAfterReadFunc(i interface{}, ctx context.Context, d *schema.ResourceData) error {
	permissions, permissions_exists := d.GetOkExists("permissions")
	if permissions_exists {
		m, is_permissions_map := permissions.(map[string]interface{})
		if is_permissions_map {
			roles, has_roles := m["role"]
			if has_roles {
				r, is_roles_list := roles.([]interface{})
				if is_roles_list {
					e := d.Set("permissions_list", r)
					if e != nil {
						return e
					}
					return nil
				}
			}
		}
	}
	return nil
}
