package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SkipDeleteAsUserSelected(ctx context.Context, d *schema.ResourceData, m interface{}) (io.Reader, error) {
	data := map[string]interface{}{}
	skip_delete := d.Get("skip_ldap")
	if skip_delete == nil {
		data["skip_ldap"] = false
	} else {
		data["skip_ldap"] = skip_delete
	}

	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
