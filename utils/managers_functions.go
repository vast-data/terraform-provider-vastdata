package utils

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	metadata "github.com/vast-data/terraform-provider-vastdata/metadata"
)

func passwordTosha256(s string) string {
	c := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", c)
}

func backwordCompatability(ctx context.Context, m map[string]interface{}) map[string]interface{} {
	// older versions will not allow managers creation with unkown attributes
	c := metadata.GetClusterVersion()
	full_featured_manager_version, _ := version.NewVersion("5.2.0")
	tflog.Debug(ctx, fmt.Sprintf("[backwordCompatability] Cluster Version Found %v", c.String()))
	if c.LessThan(full_featured_manager_version) {
		tflog.Debug(ctx, fmt.Sprintf("[backwordCompatability] Cluster version is lower than %v , removing attributes", full_featured_manager_version.String()))
		for _, t := range []string{"is_temporary_password", "password_expiration_disabled"} {
			tflog.Debug(ctx, fmt.Sprintf("[backwordCompatability] Attribute %v will be removed", t))
			_, exists := m[t]
			tflog.Debug(ctx, fmt.Sprintf("[backwordCompatability] Attribute %v found and will be removed", t))
			if exists {
				delete(m, t)
			}
		}

	}
	return m

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

func ManagerCreateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	var password string = ""
	tflog.Debug(ctx, fmt.Sprintf("[ManagerCreateFunc] Creating manager with data %v", data))
	tflog.Debug(ctx, fmt.Sprintf("[ManagerCreateFunc] Alligning data to match older versions if neede"))
	data = backwordCompatability(ctx, data)
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

	_, password_exists := b["password"]
	if password_exists {
		password = fmt.Sprintf("%v", b["password"])
	}
	b["password"] = passwordTosha256(password)
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
	data = backwordCompatability(ctx, data)
	r, err := DefaultUpdateFunc(ctx, _client, attr, data, d, headers)
	if err != nil {
		return r, err
	}
	return FakeHttpResponse(r, data)

}

func convertRoleIdsMapToList(ctx context.Context, b map[string]interface{}) (map[string]interface{}, error) {
	//Get Role IDs
	roles_ids := []interface{}{}
	roles, roles_exists := b["roles"]
	if roles_exists {
		tflog.Debug(ctx, fmt.Sprintf("[convertRoleIdsMapToList] Roles returned from VMS %v", roles))
		roles_list, is_list := roles.([]interface{})
		if !is_list {
			return b, fmt.Errorf("[convertRoleIdsMapToList] Roles return from VMS are not at list format")
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
	tflog.Debug(ctx, fmt.Sprintf("[convertRoleIdsMapToList] Roles IDs found %v", roles_ids))
	b["roles"] = roles_ids
	return b, nil

}

func collectRolesIds(ctx context.Context, r *http.Response, m map[string]interface{}) (map[string]interface{}, error) {
	b := map[string]interface{}{}
	e := UnmarshalBodyToMap(r, &b)
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

func convertPermissionsToPermissonsList(n map[string]interface{}) map[string]interface{} {
	var _permissions []interface{}
	permissions, permissions_exists := n["permissions"]
	if permissions_exists {
		m, is_permissions_map := permissions.(map[string]interface{})
		if is_permissions_map {
			roles, has_roles := m["role"]
			if has_roles {
				r, is_roles_list := roles.([]interface{})
				if is_roles_list {
					_permissions = r
				}
			}
		}
	}
	n["permissions"] = _permissions
	n["permissions_list"] = _permissions
	return n

}

func ManagersGetFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	tflog.Debug(ctx, fmt.Sprintf("[ManagersGetFunc] Calling Default Get Func"))
	r, e := DefaultGetFunc(ctx, _client, attr, d, headers)
	if e != nil {
		return r, e
	}
	m := map[string]interface{}{}
	//	p := passwordTosha256(fmt.Sprintf("%v", d.Get("password")))

	b, e := collectRolesIds(ctx, r, m)
	if e != nil {
		return r, e
	}

	b, e = collectPermissions(ctx, r, b)
	if e != nil {
		return r, e
	}
	b["password"] = fmt.Sprintf("%v", d.Get("password"))
	tflog.Debug(ctx, fmt.Sprintf("[ManagersGetFunc] Password Found %v", b["password"]))
	t, _ := json.MarshalIndent(b, "", "   ")
	tflog.Debug(ctx, fmt.Sprintf("[ManagersGetFunc] Manager returned: %v", string(t)))
	return FakeHttpResponse(r, b)
}

func collect_permissions_list(i interface{}, ctx context.Context, d *schema.ResourceData) error {
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

func ManagerAfterReadFunc(i interface{}, ctx context.Context, d *schema.ResourceData) error {
	return collect_permissions_list(i, ctx, d)
}

func ManagerPasswordChangedDiffSupress(k, oldValue, newValue string, d *schema.ResourceData) bool {
	o := fmt.Sprintf("%v", oldValue)
	n := passwordTosha256(fmt.Sprintf("%v", d.Get("password")))
	return o == n
}

func ManagerImportGetFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	tflog.Debug(ctx, fmt.Sprintf("[ManagerImportGetFunc] Calling Default get func with attributes %v", attr))
	r, e := DefaultGetFunc(ctx, _client, attr, d, headers)
	// We asume that the result is a list
	if e != nil {
		return r, e
	}

	m := []map[string]interface{}{}
	e = UnmarshalBodyToMapsList(r, &m)
	tflog.Debug(ctx, fmt.Sprintf("[ManagerImportGetFunc] Umarshel response %v", m))
	if len(m) < 1 {
		return r, fmt.Errorf("[ManagerImportGetFunc] The value returned for data %v is not a map", attr)
	}
	i := m[0]
	b, e := convertRoleIdsMapToList(ctx, i)
	if e != nil {
		return r, e
	}
	tflog.Debug(ctx, fmt.Sprintf("[ManagerImportGetFunc] Collected Roles IDs %v", b))
	b, e = collectPermissions(ctx, r, b)
	if e != nil {
		return r, e
	}
	tflog.Debug(ctx, fmt.Sprintf("[ManagerImportGetFunc] Collected Permissions %v", b))

	//	b = convertPermissionsToPermissonsList(b)
	tflog.Debug(ctx, fmt.Sprintf("[ManagerImportGetFunc] Converted Response returned %v", b))
	return FakeHttpResponseAny(r, []map[string]interface{}{b})
}

type ManagerImporter struct {
}

func (i *ManagerImporter) GetDoc() []string {
	return []string{"<guid>", "<Username>"}
}

func (i *ManagerImporter) GetFunc() GetFuncType {
	return i.getFunc
}

func (i *ManagerImporter) genQuery(s string) (string, error) {
	values := url.Values{}
	// we check if GUID was provided
	_, err := guid.FromString(s)
	if err == nil {
		//The given string is a GUID, query string will return guid=<GUID>
		values.Add("guid", s)
		return values.Encode(), nil
	}
	//If we got here than what ever was provided was not a valid guid so we would use username fields
	values.Add("username", s)
	return values.Encode(), nil
}

func (i *ManagerImporter) getFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	s := fmt.Sprintf("%v", d.Id())
	query, err := i.genQuery(s)
	if err != nil {
		return nil, err
	}
	attr["query"] = query
	return ManagerImportGetFunc(ctx, _client, attr, d, headers)
}
