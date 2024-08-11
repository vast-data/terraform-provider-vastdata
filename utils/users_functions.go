package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var user_lists_attributes []string = []string{"gids", "groups", "s3_policies_ids"}
var user_boolean_attributes []string = []string{"s3_superuser", "allow_delete_bucket", "allow_create_bucket"}

func UserBeforePatchFunc(m map[string]interface{}, i interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	FieldsUpdate(ctx, user_lists_attributes, d, &m)
	FieldsUpdate(ctx, user_boolean_attributes, d, &m)
	return m, nil
}
