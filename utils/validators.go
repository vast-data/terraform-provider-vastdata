package utils

import (
	//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"fmt"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const grace_period_format = "15:04:05"

func GracePeriodFormatValidation(i interface{}, c cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	_, e := time.Parse(grace_period_format, i.(string))
	if e != nil {
		return diag.FromErr(e)
	}
	return diags

}

func OneOf(l []string) schema.SchemaValidateDiagFunc {
	return func(i interface{}, c cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		for _, n := range l {
			if n == i.(string) {
				return diags
			}
		}
		return diag.Errorf(fmt.Sprintf("Wrong Value Provided %v, Allowed values are %v", i, l))
	}
}
