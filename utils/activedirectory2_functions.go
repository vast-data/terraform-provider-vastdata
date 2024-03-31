package utils

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		  "machine_account_name": "cluster123",
		  .
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

	err = UnmarshelBodyToMap(response, &temp)
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
