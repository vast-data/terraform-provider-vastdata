package utils

import (
	//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"fmt"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	grace_period_format = "15:04:05"
	start_at_format     = "2006-01-02 15:04:05"
)

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

func SnapshotExpirationFormatValidation(i interface{}, c cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	_, e := time.Parse(time.RFC3339, i.(string))
	if e != nil {
		return diag.FromErr(e)
	}
	return diags

}

func ProtectionPolicyStartAt(i interface{}, c cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	_, e := time.Parse(start_at_format, i.(string))
	if e != nil {
		return diag.FromErr(e)
	}
	return diags

}
