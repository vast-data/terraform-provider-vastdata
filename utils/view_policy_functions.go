package utils

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	metadata "github.com/vast-data/terraform-provider-vastdata/metadata"
)

var permissions_attributes = []string{"nfs_all_squash", "nfs_root_squash", "nfs_read_write", "nfs_read_only", "s3_read_only", "s3_read_write", "smb_read_only", "smb_read_write", "nfs_no_squash"}
var viewpolicy_boolean_attributes = []string{"smb_is_ca", "enable_access_to_snapshot_dir_in_subdirs", "enable_visibility_of_snapshot_dir", "nfs_enforce_tls", "nfs_case_insensitive", "nfs_posix_acl"}
var min_vippool_permission_version, _ = version.NewVersion("5.1.0")

func _convert_vip_pools_to_permission_per_vip_pool(i interface{}, m *map[string]interface{}) {
	t := map[string]string{}
	l, is_list := i.([]interface{})
	if is_list {
		for _, v := range l {
			t[fmt.Sprintf("%v", v)] = "RW"
		}
	}
	if len(t) > 0 {
		(*m)["permission_per_vip_pool"] = t
	}

}

func __vippool_permission_convert(i interface{}) map[string]string {
	permission_per_vip_pool := map[string]string{}
	perms, is_interface_list := i.([]interface{})
	if is_interface_list && len(perms) > 0 {
		for _, p := range perms {
			q, is_interface_map := p.(map[string]interface{})
			if is_interface_map {
				vippool_id, e1 := q["vippool_id"]
				vippool_permissions, e2 := q["vippool_permissions"]
				if e1 && e2 {
					permission_per_vip_pool[fmt.Sprintf("%v", vippool_id)] = fmt.Sprintf("%v", vippool_permissions)
				}

			}
		}
	}
	return permission_per_vip_pool
}

func _vippool_permission_convert(ctx context.Context, i interface{}, m *map[string]interface{}) {
	tflog.Debug(ctx, fmt.Sprintf("[ViewPolicyConvertVippoolPermissions] - VipPool Permissions : %v ", i))
	permission_per_vip_pool := __vippool_permission_convert(i)
	if len(permission_per_vip_pool) > 0 {
		(*m)["permission_per_vip_pool"] = permission_per_vip_pool
	}
}

func vippool_permission_convert_for_update(ctx context.Context, d *schema.ResourceData, m *map[string]interface{}) {
	tflog.Debug(ctx, fmt.Sprintf("[ViewPolicyUpdate] - VipPool Permissions : %v", d.Get("vippool_permissions")))
	i, e := d.GetOkExists("vippool_permissions")
	if !e {
		return
	}
	_vippool_permission_convert(ctx, i, m)
}

func vippool_permission_convert_for_create(ctx context.Context, m *map[string]interface{}) {
	i, exists := (*m)["vippool_permissions"]
	tflog.Debug(ctx, fmt.Sprintf("[ViewPolicyCreate] - VipPool Permissions : %v , Exist: %v", i, exists))
	if !exists {
		return
	}
	_vippool_permission_convert(ctx, i, m)

}
func setupS3SpecialCharsSupport(ctx context.Context, v string, m *map[string]interface{}) {
	(*m)["s3_special_chars_support"] = v
}

func VippoolPermissionsIdsDiffSupress(k, oldValue, newValue string, d *schema.ResourceData) bool {

	oldData, newData := d.GetChange("vippool_permissions")
	if oldData == nil || newData == nil { // if any of them is nil it means new data was set so there can be no diff
		return false
	}
	o := __vippool_permission_convert(oldData)
	n := __vippool_permission_convert(newData)
	return reflect.DeepEqual(o, n)
}

func checkAuthProviders(ctx context.Context, data map[string]interface{}) (string, error) {
	use_auth_providers, exists := data["use_auth_provider"]
	_use_auth_providers := strings.ToLower(fmt.Sprintf("%v", use_auth_providers))
	auth_source, auth_source_exists := data["auth_source"]
	_auth_source := fmt.Sprintf("%v", auth_source)
	tflog.Debug(ctx, fmt.Sprintf("Evaluating the usage of auth providers, use_auth_providers:%v , auth_source:%v", _use_auth_providers, _auth_source))
	if exists && _use_auth_providers == "false" {
		/*it means false we would not want it set so if it was defined we returne the defined value if not we return RPC
		  But if use_auth_providers is set but the value of auth_source is anything but PROVIDERS we return an error
		*/
		if !auth_source_exists || _auth_source == "nil" {
			return "RPC", nil
		}
		if _auth_source != "RPC" {
			return "", fmt.Errorf("When use_auth_providers is set to false auth_source must be set to either \"RPC\"")
		}
		return fmt.Sprintf("%v", auth_source), nil

	}
	if exists && _use_auth_providers == "true" {
		if !auth_source_exists || auth_source == nil {
			return "PROVIDERS", nil
		}
		if _auth_source != "PROVIDERS" && _auth_source != "RPC_AND_PROVIDERS" {
			return "", fmt.Errorf("When use_auth_providers is set to \"true\" auth_source must be set to PROVIDERS or compleatly removed")
		}
	}
	if auth_source_exists {
		return _auth_source, nil
	}
	//If we got to this part is means that use_auth_provider and auth_source are not provided which means we only left with RPC
	return "RPC", nil
}

func ViewPolicyPermissionsSetup(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	for _, v := range permissions_attributes {
		q, k := d.GetOk(v)
		tflog.Debug(ctx, fmt.Sprintf("Data recived for ViewPolicy : %v Before Creation %v:%v", v, q, k))
		if !k {
			/*
					If k is false it means one of 2 things
					1. No value was given and this means to use the Zero value ([]).
					2. That the user provided the Zero Value ([]).
				        In any case we want the Zero Value otherwise it means that some value was provided by the user which is not Zero.
			*/
			m[v] = []string{}
		}

	}

	//	tflog.Debug(ctx, fmt.Sprintf("Data recived for ViewPolicy Before Creation nfs_read_write %v", d.Get("nfs_read_write")))
	//	tflog.Debug(ctx, fmt.Sprintf("Data recived for ViewPolicy Before Creation nfs_all_squash %v", d.Get("nfs_all_squash")))
	FieldsUpdate(ctx, viewpolicy_boolean_attributes, d, &m)
	return m, nil
}

func ViewPolicyCreateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	auth_provider, err := checkAuthProviders(ctx, data)
	if err != nil {
		return nil, err
	}
	data["auth_provider"] = auth_provider
	flavor, flavor_exists := data["flavor"]
	if flavor_exists && (fmt.Sprintf("%v", flavor) == "S3_NATIVE") {
		z, e := data["s3_special_chars_support"]
		if !e {
			z = "false"
		}
		setupS3SpecialCharsSupport(ctx, fmt.Sprintf("%v", z), &data)

	}
	vippool, vippoolexists := data["vip_pools"]
	ver := metadata.GetClusterVersion()
	if vippoolexists && ver.GreaterThanOrEqual(min_vippool_permission_version) {
		_convert_vip_pools_to_permission_per_vip_pool(vippool, &data)
	}
	vippool_permission_convert_for_create(ctx, &data)
	return DefaultCreateFunc(ctx, _client, attr, data, headers)
}

func ViewPolicyUpdateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	use_auth_provider, exists := d.GetOk("use_auth_provider")
	auth_source, auth_source_exists := d.GetOk("auth_source")
	tflog.Debug(ctx, fmt.Sprintf("Updating View Policy with with use_auth_provider = %v & exists = %v", use_auth_provider, exists))
	tflog.Debug(ctx, fmt.Sprintf("Updating View Policy with with auth_source = %v & exists = %v", auth_source, auth_source_exists))
	if exists {
		data["use_auth_provider"] = use_auth_provider
	} else {
		data["use_auth_provider"] = false
	}

	if auth_source_exists {
		data["auth_source"] = auth_source
	} else {
		data["auth_source"] = nil
	}

	auth_provider, err := checkAuthProviders(ctx, data)
	if err != nil {
		return nil, err
	}
	data["auth_provider"] = auth_provider
	zero := []string{}

	for _, v := range permissions_attributes {
		q, f := d.GetChange(v)
		k := d.HasChange(v)
		tflog.Debug(ctx, fmt.Sprintf("ViewPolicy attribute: %v Has Change: %v , Change: %v <==> %v", v, k, q, f))
		if k {
			i := f.([]interface{})
			if len(i) == 0 {
				tflog.Debug(ctx, fmt.Sprintf("ViewPolicy attribute : %v is zero value and will be set", v))
				data[v] = zero
			}
		}

	}
	flavor, flavor_exists := data["flavor"]
	if flavor_exists && (fmt.Sprintf("%v", flavor) == "S3_NATIVE") {
		z, e := data["s3_special_chars_support"]
		if !e {
			z = "false"
		}
		setupS3SpecialCharsSupport(ctx, fmt.Sprintf("%v", z), &data)

	}
	vippool, vippoolexists := data["vip_pools"]
	tflog.Debug(ctx, fmt.Sprintf("[ViewPolicyUpdateFunc] - Vippool Exist : %v , vippool %v", vippoolexists, vippool))
	ver := metadata.GetClusterVersion()
	if vippoolexists && ver.GreaterThanOrEqual(min_vippool_permission_version) {
		tflog.Debug(ctx, fmt.Sprintf("[ViewPolicyUpdateFunc] - Converting vip_pools %v to permission_per_vip_pool format", vippool))
		_convert_vip_pools_to_permission_per_vip_pool(vippool, &data)
		pr, _ := data["permission_per_vip_pool"]
		tflog.Debug(ctx, fmt.Sprintf("[ViewPolicyUpdateFunc] - converted vip_pools %v to permission_per_vip_pool format, result :%v", vippool, pr))

	}
	old, new := d.GetChange("vip_pools")
	tflog.Debug(ctx, fmt.Sprintf("[ViewPolicyUpdateFunc] - vip_pools changed from %v To %v", old, new)) // this is a special case where the vip_pools was defined but set to []
	if (!vippoolexists) && (!reflect.DeepEqual(old, new)) {
		vippoolexists = true
		data["permission_per_vip_pool"] = map[string]string{}
		data["vip_pools"] = []int{}
	}

	if !vippoolexists {
		vippool_permission_convert_for_update(ctx, d, &data)
	}

	return DefaultUpdateFunc(ctx, _client, attr, data, d, headers)
}

func ViewPolicyGetFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	response, err := DefaultGetFunc(ctx, _client, attr, d, headers)
	if err != nil {
		return response, err
	}
	u := map[string]interface{}{}
	err = UnmarshalBodyToMap(response, &u)
	if err != nil {
		return response, err
	}
	l := []map[string]interface{}{}
	r := []int{}
	i, e := u["permission_per_vip_pool"]
	if e {
		for k, v := range i.(map[string]interface{}) {
			l = append(l, map[string]interface{}{"vippool_id": k, "vippool_permissions": v})
			n, converr := strconv.Atoi(k)
			if converr != nil {
				tflog.Debug(ctx, fmt.Sprintf("[ViewPolicyGetFunc] - Can not convert %v to integer, error: %v", k, converr))
			} else {
				r = append(r, n)
			}
		}
		u["vippool_permissions"] = l
		u["vip_pools"] = r

	}

	return FakeHttpResponse(response, u)
}

func ViewPolicyBeforePatchFunc(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	FieldsUpdate(ctx, viewpolicy_boolean_attributes, d, &m)
	return m, nil
}
