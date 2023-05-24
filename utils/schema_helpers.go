package utils

//This package will hold helper functions to be used when definning a schema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DoNothingOnUpdate() schema.SchemaDiffSuppressFunc {
	return func(k, oldValue, newValue string, d *schema.ResourceData) bool {
		if oldValue == "" && newValue != "" {
			//Empty oldValue while The new Value is not empty means
			//that this is probably creation time
			//So we need this one for creation , but not for update
			return false
		}
		return true
	}
}
