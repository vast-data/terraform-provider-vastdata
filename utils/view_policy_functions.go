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
	return DefaultUpdateFunc(ctx, _client, attr, data, d, headers)
}
