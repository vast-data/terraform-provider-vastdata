package utils

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
)

var forbiddenKeys map[string]interface{} = map[string]interface{}{
	"guid":                nil,
	"id":                  nil,
	"active_directory_id": nil,
	"title":               nil,
	"state":               nil,
	"active_directory":    nil,
}

/*
	When performing get action to https://<cluster ip>/api/<version>/activedirectory/<id> the data returned is not flattend
	some attributes such as "searchbase", port ..... are returned under an object ldap.

Ex:

		{ "id": 1,
		  "machine_account_name": "cluster123",		  .
		  .
		  .
		  "ldap": {
		            "bindpw": "<password>",
		            "use_tls": false
		            .
		            .
		          }

	   To cope with this since the POST request is flatten , but the GET is not we will flatten attributes as part of this GET request
*/
func ActiveDirectory2GetFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	var temp map[string]interface{} = map[string]interface{}{}
	var temp2 map[string]interface{} = map[string]interface{}{}

	response, err := DefaultGetFunc(ctx, _client, attr, d, headers)
	if err != nil {
		return response, err
	}

	err = UnmarshalBodyToMap(response, &temp)
	if err != nil {
		return response, err
	}
	//Copy all key:value pairs whcih are not ldap
	for k, v := range temp {
		if k != "ldap" {
			temp2[k] = v
		}
	}
	//Copy selected Values under the ldap object (if such exists)
	ldap, exists := temp["ldap"]
	if exists {
		//ldap is an object this means map[string]interface{}
		_ldap := ldap.(map[string]interface{})
		for k, v := range _ldap {
			_, is_forbidden := forbiddenKeys[k]
			if is_forbidden {
				continue
			}
			temp2[k] = v
		}
	}
	return FakeHttpResponse(response, temp2)
}

func ActiveDirectory2DeleteFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	client := _client.(*vast_client.VMSSession)
	attributes, err := getAttributesAsString([]string{"path", "id"}, attr)
	if err != nil {
		return nil, err
	}
	delete_path := fmt.Sprintf("%v%v", (*attributes)["path"], (*attributes)["id"])
	query := ""
	_query := getAttributeOrDefault("query", nil, attr)
	if _query != nil {
		query = *_query
	}
	b := []byte("{}")
	r := bytes.NewReader(b)
	tflog.Debug(ctx, fmt.Sprintf("Calling Delete to path \"%v\"", delete_path))
	return client.Delete(ctx, delete_path, query, r, map[string]string{})

}

type ImportActiveDirectory2ByHttpFields struct {
	*ImportByHttpFields
}

func (i *ImportActiveDirectory2ByHttpFields) getFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	s := fmt.Sprintf("%v", d.Id())
	query, err := i.genQuery(s)
	if err != nil {
		return nil, err
	}
	attr["query"] = query
	return ActiveDirectory2GetFunc(ctx, _client, attr, d, headers)
}

func NewImportActiveDirectory2ByHttpFields(disable_guid_import bool, fields []HttpFieldTuple) *ImportActiveDirectory2ByHttpFields {
	return &ImportActiveDirectory2ByHttpFields{ImportByHttpFields: NewImportByHttpFields(disable_guid_import, fields)}

}
