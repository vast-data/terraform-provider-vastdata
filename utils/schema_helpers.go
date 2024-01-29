package utils

//This package will hold helper functions to be used when definning a schema

import (
	"encoding/json"
	"reflect"

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

func JsonStructureCompare(k, old, new string, d *schema.ResourceData) bool {
	if old == "" && new != "" {
		//Empty oldValue while The new Value is not empty means
		//that this is probably creation time
		//So we need this one for creation , but not for update
		return false
	}
	var x1 interface{}
	var x2 interface{}
	json.Unmarshal([]byte(old), x1)
	json.Unmarshal([]byte(new), x2)
	return reflect.DeepEqual(x1, x2)
}
