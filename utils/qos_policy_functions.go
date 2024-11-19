package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	metadata "github.com/vast-data/terraform-provider-vastdata/metadata"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
)

func always_send_policy_type(data *map[string]interface{}, d *schema.ResourceData) {
	policy_type := d.Get("policy_type")
	(*data)["policy_type"] = policy_type
}

func always_send_mode(data *map[string]interface{}, d *schema.ResourceData) {
	mode := d.Get("mode")
	(*data)["mode"] = mode
}

func set_use_total_limits(m *map[string]interface{}, ctx context.Context, d *schema.ResourceData) {
	policy_type, policy_type_exists := d.GetOkExists("policy_type")
	if policy_type_exists && fmt.Sprintf("%v", policy_type) == "USER" { // When policy type is USER we should not send use_totla_limits
		return
	}
	_, static_total_limits_exists := d.GetOkExists("static_total_limits")
	_, capacity_total_limits_exsts := d.GetOkExists("capacity_total_limits")
	tflog.Debug(ctx, fmt.Sprintf("[set_use_total_limits] , capacity_total_limits_exists: %v, static_total_limits_exists:%v", capacity_total_limits_exsts, static_total_limits_exists))
	(*m)["use_total_limits"] = static_total_limits_exists || capacity_total_limits_exsts
}

func policy_type_validation(ctx context.Context, data map[string]interface{}) error {
	v := metadata.GetClusterVersion()
	min, _ := version.NewVersion("5.1.0")
	if v.GreaterThanOrEqual(min) {
		tflog.Debug(ctx, fmt.Sprintf("[policy_type_validation] Cluster Version >= 5.1.0 going verifying QOS policy with the following data: %v", data))
		policy_type, exists := data["policy_type"]
		if !exists {
			return fmt.Errorf("policy_type , must be provided")
		}
		_, mode_exists := data["mode"]
		if fmt.Sprintf("%v", policy_type) != "USER" && !mode_exists {
			return fmt.Errorf("When policy_type is not USER , mode should be provided")
		}

	}
	return nil
}

func populate_users(ctx context.Context, users_identifiers []interface{}, c interface{}) ([]interface{}, error) {
	client := c.(vast_client.JwtSession)
	users := []interface{}{}
	for _, ind := range users_identifiers {
		i := fmt.Sprintf("%v", ind)
		u := map[string]interface{}{}
		user := map[string]interface{}{}
		r, e := client.Get(ctx, GenPath(fmt.Sprintf("%v/%v", "users", i)), "", map[string]string{})
		if e != nil {
			return nil, e
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(body, &u)
		user["fqdn"] = ""
		user["is_sid"] = true
		sid, sid_exists := u["sid"]
		if !sid_exists {
			sid = ""
		}
		user["sid_str"] = sid
		user["label"] = fmt.Sprintf("%v (%v)", u["name"], u["name"])
		user["login_name"] = u["name"]
		user["value"] = u["name"]
		user["name"] = u["name"]
		user["identifier_type"] = "sid_str"
		user["identifier_value"] = sid
		users = append(users, user)
	}
	return users, nil
}

func SetQoSPolicyUseTotalLimits(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	set_use_total_limits(&m, ctx, d)
	return m, nil
}

func QosCreateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	if e := policy_type_validation(ctx, data); e != nil {
		return nil, e
	}
	attached_users, exists := data["attached_users_identifiers"]
	if exists {
		p, e := populate_users(ctx, attached_users.([]interface{}), _client)
		if e != nil {
			return nil, e
		}
		data["attached_users"] = p
	}
	return DefaultCreateFunc(ctx, _client, attr, data, headers)
}

func QosUpdateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	always_send_policy_type(&data, d)
	if fmt.Sprintf("%v", data["policy_type"]) != "USER" {
		always_send_mode(&data, d)
	}
	if e := policy_type_validation(ctx, data); e != nil {
		return nil, e
	}

	attached_users, exists := data["attached_users_identifiers"]
	if exists {
		p, e := populate_users(ctx, attached_users.([]interface{}), _client)
		if e != nil {
			return nil, e
		}
		data["attached_users"] = p
	}
	return DefaultUpdateFunc(ctx, _client, attr, data, d, headers)
}

func QosAfterReadFunc(i interface{}, ctx context.Context, d *schema.ResourceData) error {
	tflog.Debug(ctx, fmt.Sprintf("[QosAfterReadFunc] - Start"))
	attached_users, exists := d.GetOkExists("attached_users")
	client := i.(vast_client.JwtSession)
	attached_users_identifiers := []interface{}{}
	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("[QosAfterReadFunc] - attached_users was not found"))
		return nil
	}
	tflog.Debug(ctx, fmt.Sprintf("[QosAfterReadFunc] - Attached users: %v", attached_users))
	m, valid := attached_users.([]interface{})
	if !valid {
		tflog.Debug(ctx, fmt.Sprintf("[QosAfterReadFunc] - Attached users: %v is not valid []interface{}", attached_users))
		return nil
	}
	for _, o := range m {
		t, is_map := o.(map[string]interface{})
		if !is_map {
			tflog.Debug(ctx, fmt.Sprintf("[QosAfterReadFunc] - Converting user data %v to map[string]interface{} failed", t))
			continue
		}
		tflog.Debug(ctx, fmt.Sprintf("[QosAfterReadFunc] - Converting user data %v", t))
		sid, e := t["identifier_value"]
		if !e {
			tflog.Debug(ctx, fmt.Sprintf("[QosAfterReadFunc] - Converting user data %v failed nod sid found", t))
			continue
		}
		tflog.Debug(ctx, fmt.Sprintf("[QosAfterReadFunc] - Querying for sid %v", sid))
		u := url.Values{}
		u.Add("fields", "id")
		u.Add("sid", fmt.Sprintf("%v", sid))
		response, err := client.Get(ctx, GenPath("users"), u.Encode(), map[string]string{})
		if err != nil {
			tflog.Debug(ctx, fmt.Sprintf("[QosAfterReadFunc] - Querying for user with sid %v failed , reason %v", sid, err))
			continue
		}
		j := []map[string]interface{}{}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			tflog.Debug(ctx, fmt.Sprintf("[QosAfterReadFunc] - Converting user data %v failed to read reponse body %v", t, e))
			continue
		}
		uerr := json.Unmarshal(body, &j)
		if uerr != nil {
			tflog.Debug(ctx, fmt.Sprintf("[QosAfterReadFunc] - Converting user data %v failed to unmarshal response body %v", t, uerr))
			continue
		}
		if len(j) == 0 {
			tflog.Debug(ctx, fmt.Sprintf("[QosAfterReadFunc] - Converting user data %v failed not elements returned from response body", t))
			continue
		}
		q := j[0]
		id, id_exists := q["id"]
		if !id_exists {
			tflog.Debug(ctx, fmt.Sprintf("[QosAfterReadFunc] - Converting user data %v failed , no id was found", t))
			continue
		}
		attached_users_identifiers = append(attached_users_identifiers, id)

	}
	tflog.Debug(ctx, fmt.Sprintf("[QosAfterReadFunc] Setting attached_users_identifiers to %v", attached_users_identifiers))
	d.Set("attached_users_identifiers", attached_users_identifiers)
	return nil
}
