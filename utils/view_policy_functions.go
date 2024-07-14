package utils

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
	for _, v := range []string{"nfs_all_squash", "nfs_root_squash", "nfs_read_write", "nfs_read_only", "s3_read_only", "s3_read_write", "smb_read_only", "smb_read_write"} {
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
	return m, nil
}

func ViewPolicyCreateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	auth_provider, err := checkAuthProviders(ctx, data)
	if err != nil {
		return nil, err
	}
	data["auth_provider"] = auth_provider
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

	for _, v := range []string{"nfs_all_squash", "nfs_root_squash", "nfs_read_write", "nfs_read_only", "s3_read_only", "s3_read_write", "smb_read_only", "smb_read_write"} {
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

	return DefaultUpdateFunc(ctx, _client, attr, data, d, headers)
}
