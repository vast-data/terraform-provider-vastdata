package utils

import (
	"context"
	"fmt"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
	"reflect"
	"regexp"
	"runtime"
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
	/*While most of the functionality of this function is done by functions which are designed to populate
	  data in preparation to POST/PATCH requests , when it comes to booleans there is a problem with TF assuming false is zero value
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

func GetFuncName(i interface{}) string {
	//Return the function name
	p := runtime.FuncForPC(reflect.ValueOf(i).Pointer())
	if p != nil {
		return p.Name()
	}
	return "Unknown"

}

func HandleFallback(ctx context.Context, client *vast_client.VMSSession, attrs map[string]interface{}, d *schema.ResourceData, idFunc IdFuncType) ([]byte, error) {
	response, fallbackErr := DefaultGetByGUIDFunc(ctx, client, attrs, d, map[string]string{})
	if fallbackErr != nil {
		return nil, fallbackErr
	}
	var id string
	body, id, fallbackErr := GetBodyBytesAndId(response)
	if fallbackErr != nil {
		return nil, fallbackErr
	}
	fallbackErr = idFunc(ctx, nil, id, d)
	if fallbackErr != nil {
		return nil, fallbackErr
	}
	return body, nil
}
