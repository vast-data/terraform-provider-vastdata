package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"sync"

	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
)

var mu sync.Mutex

func createUserKeyPath(path, user_id interface{}) string {
	return fmt.Sprintf("%v/%v/access_keys/", path, user_id)
}

func getID(user_id, access_key interface{}) string {
	return fmt.Sprintf("%v-%v", user_id, access_key)
}

func encryptSecretToken(token, pgp_public_key interface{}) (string, error) {
	//Encypt the secret with the PGP public key
	//for now it does nothing
	_token := fmt.Sprintf("%v", token)
	_pgp_public_key := fmt.Sprintf("%v", pgp_public_key)
	return helper.EncryptMessageArmored(_pgp_public_key, _token)
}

func readFromResource(resource api_latest.UserKey) map[string]interface{} {
	o := map[string]interface{}{}
	//Read values from struct
	o["pgp_public_key"] = resource.PgpPublicKey
	o["secret_key"] = resource.SecretKey
	o["encrypted_secret_key"] = resource.EncryptedSecretKey
	return o
}

func genNewHttpResponse(ctx context.Context, _client interface{}, path string, h *http.Response, d map[string]interface{}) (*http.Response, error) {
	response_body, _ := io.ReadAll(h.Body)
	i := map[string]interface{}{}

	uerr := json.Unmarshal(response_body, &i)
	if uerr != nil {
		return nil, uerr
	}
	tflog.Debug(ctx, fmt.Sprintf("USERKEY: Date returned %v", i))
	_enabled := false
	enabled, enabled_exists := d["enabled"] //If it can not be found it is a false
	tflog.Debug(ctx, fmt.Sprintf("USERKEY: Data recived %v", d))
	if enabled_exists {
		_enabled = enabled.(bool)
	}
	disableKeyIfneeded(ctx, _client, path, fmt.Sprintf("%v", i["access_key"]), _enabled)
	i["id"] = getID(d["user_id"], i["access_key"])
	i["user_id"] = d["user_id"]
	i["pgp_public_key"] = d["pgp_public_key"]
	pgp_public_key, exists := d["pgp_public_key"]
	if exists {
		tflog.Debug(ctx, fmt.Sprintf("USERKEY: PGP public key found , encrypting secret"))
		e, err := encryptSecretToken(i["secret_key"], pgp_public_key)
		if err != nil {
			return nil, err
		}
		i["secret_key"] = nil
		i["encrypted_secret_key"] = e
	} else {
		tflog.Warn(ctx, fmt.Sprintf("USERKEY: No PGP public key was found secret_key will hold the plain text file"))
		i["encrypted_secret_key"] = nil
	}
	tflog.Debug(ctx, fmt.Sprintf("USERKEY: Date written back %v", i))
	b, merr := json.Marshal(i)
	if merr != nil {
		return nil, merr
	}
	tflog.Warn(ctx, fmt.Sprintf("USERKEY: JSON return %v", string(b)))
	return &http.Response{Body: io.NopCloser(bytes.NewReader(b))}, nil
}

func disableKeyIfneeded(ctx context.Context, _client interface{}, path, access_key string, enabled bool) (*http.Response, error) {
	/*You can not create a disabled key, when using he Vast WebUI it is not a problem since it is not even an option to define new key as disabled.
	  However using terraform it is possiable to define enabled=false when definning an UserKey.
	  But it can not be done since the Vast API does not supports it, there for this function which should be called only upon creation and as much as possiable
	  close to the key creation will disable the key.
	*/
	tflog.Debug(ctx, fmt.Sprintf("USERKEY: Checking to see if we should disable access_key:%v for path:%v, enabled:%v", access_key, path, enabled))
	if enabled {
		tflog.Debug(ctx, "access_key is set to be enabled, nothign to do")
		return nil, nil
	}
	client := _client.(vast_client.JwtSession)
	payload := map[string]interface{}{}
	payload["access_key"] = access_key
	payload["enabled"] = enabled

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	tflog.Debug(ctx, fmt.Sprintf("USERKEY: Disabling key with paylod %v", string(b)))
	return client.Patch(ctx, path, "", bytes.NewReader(b), map[string]string{})
}

func CreateUserKeyFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	mu.Lock() //Multiple creations are not working well , we need to make sure only one key is created at one time
	defer mu.Unlock()
	client := _client.(vast_client.JwtSession)
	attributes, err := getAttributesAsString([]string{"path"}, attr)
	if err != nil {
		return nil, err
	}
	path := createUserKeyPath((*attributes)["path"], data["user_id"])
	tflog.Debug(ctx, fmt.Sprintf("Creating UserKey for user with id:%s , path:%s", data["user_id"], path))
	//The API needs to send an empty POST object
	response, response_err := client.Post(ctx, path, bytes.NewReader([]byte("{}")), map[string]string{})
	if response_err != nil {
		return nil, response_err
	}
	return genNewHttpResponse(ctx, _client, path, response, data)
}

func DeleteUserKeyFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	client := _client.(vast_client.JwtSession)
	attributes, err := getAttributesAsString([]string{"path", "id"}, attr)
	if err != nil {
		return nil, err
	}
	s := strings.SplitN((*attributes)["id"], "-", 2)
	key_path := createUserKeyPath((*attributes)["path"], s[0])
	payload := map[string]string{"access_key": s[1]}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return client.Delete(ctx, key_path, "", bytes.NewReader(b), headers)

}

func UpdateUserKeyFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	client := _client.(vast_client.JwtSession)
	payload := map[string]interface{}{"access_key": d.Get("access_key"), "enabled": d.Get("enabled")}
	attributes, err := getAttributesAsString([]string{"path"}, attr)
	key_path := createUserKeyPath((*attributes)["path"], d.Get("user_id"))
	tflog.Debug(ctx, fmt.Sprintf("USERKEY: Update path %v", key_path))
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	tflog.Debug(ctx, fmt.Sprintf("USERKEY: Update payload", string(b)))
	return client.Patch(ctx, key_path, "", bytes.NewReader(b), headers)

}

func GetUserKeyFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	//There is not GET for a key we will have to iterate over all user keys to find this specific key
	client := _client.(vast_client.JwtSession)
	attributes, err := getAttributesAsString([]string{"path", "id"}, attr)
	resource := ctx.Value(ContextKey("resource"))
	tflog.Debug(ctx, fmt.Sprintf("USERKEY: Resource %v found", resource))
	if resource != nil {
		r := readFromResource(resource.(api_latest.UserKey))
		d.Set("pgp_public_key", r["pgp_public_key"])
		d.Set("secret_key", r["secret_key"])
		d.Set("encrypted_secret_key", r["encrypted_secret_key"])
	}
	if err != nil {
		return nil, err
	}
	s := strings.SplitN((*attributes)["id"], "-", 2)
	i := map[string]interface{}{}
	access_key := s[1]
	path := fmt.Sprintf("%v/%v", (*attributes)["path"], s[0])
	tflog.Debug(ctx, fmt.Sprintf("USERKEY: Calling Path %v to get user details", path))
	response, rerr := client.Get(ctx, path, "", headers)
	if rerr != nil {
		return nil, rerr
	}
	body, berr := io.ReadAll(response.Body)
	if berr != nil {
		return nil, berr
	}

	err = json.Unmarshal(body, &i)
	if err != nil {
		return nil, err
	}
	tflog.Debug(ctx, fmt.Sprintf("USERKEY: Reponse unmarsheled: %v", i))
	response_payload := map[string]interface{}{}
	//Basic Setups
	response_payload["user_id"] = d.Get("user_id")
	response_payload["pgp_public_key"] = d.Get("pgp_public_key")
	response_payload["secret_key"] = d.Get("secret_key")
	response_payload["encrypted_secret_key"] = d.Get("encrypted_secret_key")
	access_keys, exists := i["access_keys"]
	tflog.Debug(ctx, fmt.Sprintf("USERKEY: %v User Access Key Found: %v", i["id"], access_keys))
	if exists {
		for _, l := range access_keys.([]interface{}) {
			v := l.(map[string]interface{})
			key, key_exists := v["key"]
			if !key_exists {
				key = ""
			}
			enabled, enabled_exists := v["status"]
			if !enabled_exists {
				enabled = ""
			}
			//we assume that the the respone is []string where fisrt string element is always the access key and the second is enabled/disabled
			if key == access_key {
				response_payload["access_key"] = key
				if enabled == "enabled" {
					response_payload["enabled"] = true
				} else {

					response_payload["enabled"] = false
				}

			}

		}

	}
	tflog.Debug(ctx, fmt.Sprintf("USERKEY: Response payload returned: %v", response_payload))
	b, err := json.Marshal(response_payload)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		Body:    io.NopCloser(bytes.NewReader(b)),
		Request: response.Request}, nil

}

func AddLostDataBackToUserKey(i interface{}, ctx context.Context, d *schema.ResourceData) error {
	o, n := d.GetChange("secret_key")
	tflog.Debug(ctx, fmt.Sprintf("USERKEY: Old secret ID %v, New secret Id %v", o, n))
	return nil
}
