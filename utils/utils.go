package utils

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetTFformatName(s string) string {
	re := regexp.MustCompile("([A-Z])")
	t := re.ReplaceAllString(s, "_${1}")
	return strings.ToLower(strings.TrimLeft(t, "_"))
}

func FieldsUpdate(ctx context.Context, fields []string, d *schema.ResourceData, m *map[string]interface{}) {
	/*While most of the functionality of this function is done by functions which are designed to pupulate
	  data in preperation to POST/PATCH requests , when it comes to booleans there is a problem with TF assuming false is zero value
	  and json not populating them

	  This function will  force the value to the map so it will be sent.

	  It is advisable to use this function mostly for booleans
	*/

	for _, i := range fields {
		_, e := d.GetOkExists(i)
		if !e {
			continue
		}
		old, new := d.GetChange(i)
		tflog.Debug(ctx, fmt.Sprintf("[FieldsUpdate] - Old:%v , New:%v", old, new))
		if old == nil || new == nil { // creation time can generate nil especially if no default is provided
			(*m)[i] = false

		} else {
			(*m)[i] = new
		}
	}
}
