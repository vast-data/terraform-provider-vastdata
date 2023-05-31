package utils

import (
	//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
