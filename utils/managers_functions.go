package utils

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
	data = addRealmPermissions(ctx, data)
	delete(data, "realms_permissions")
	r, e := DefaultCreateFunc(ctx, _client, attr, data, headers)
	if e != nil {
		return nil, e
	}

	b, e := convertManagerResponseToTfModel(ctx, r, data)
	if e != nil {
		return r, e
	}
	b = addRealmPermissions(ctx, b)
	return FakeHttpResponse(r, b)

}

// func ManagerUpdateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
// 	tflog.Debug(ctx, fmt.Sprintf("[ManagerUpdateFunc] Data recive %v , data map %v", d.Get("realms_permissions"), data))
// 	l, e := d.GetOkExists("permissions_list")
// 	if e {
// 		p, z := l.([]interface{})
// 		if z {
// 			data["permissions_list"] = p
// 		}
// 	}
// 	a, b := d.GetOkExists("realms_permissions")
// 	if b {
// 		p, z := a.([]interface{})
// 		if z {
// 			data["realms_permissions"] = p
// 		}
// 	}
// 	return FakeHttpResponse(r, b)

// }

func ManagersGetFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	tflog.Debug(ctx, fmt.Sprintf("[ManagersGetFunc] Calling Default Get Func"))
	r, e := DefaultGetFunc(ctx, _client, attr, d, headers)
	if e != nil {
		return r, e
	}
	m := map[string]interface{}{}
	p := fmt.Sprintf("%v", d.Get("password"))
	m["password"] = p
	b, e := convertManagerResponseToTfModel(ctx, r, m)
	if e != nil {
		return r, e
	}
	j, has_permissions_list := m["permissions_list"]
	var realms_permissions []map[string]interface{}
	if has_permissions_list {
		tflog.Debug(ctx, fmt.Sprintf("[ManagersGetFunc] Processing permissions list to collect permissions list: %v", j))
		realms_permissions = readReamlsPermissions(ctx, j)

	} else {

		realms_permissions = []map[string]interface{}{}
	}
	m["realms_permissions"] = realms_permissions
	tflog.Debug(ctx, fmt.Sprintf("[ManagersGetFunc] Realms permissions collected: %v", m["realms_permissions"]))
	return FakeHttpResponse(r, b)
}

func ManagerPasswordDiffSupress(k, oldValue, newValue string, d *schema.ResourceData) bool {
	return true
}
