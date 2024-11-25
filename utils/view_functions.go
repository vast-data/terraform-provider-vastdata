package utils

import (
	"fmt"
	"os"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func BucketLoggingDiff(k, oldValue, newValue string, d *schema.ResourceData) bool {
	old, new := d.GetChange(k)
	e, ex := d.GetOkExists(k)
	b := fmt.Sprintf("Old: %v , New :%v Exists: %v, E:%v", old, new, ex, e)
	os.WriteFile("/tmp/bl", []byte(b), 0644)
	if new == nil {
		return true
	}
	return reflect.DeepEqual(old, new)
}
