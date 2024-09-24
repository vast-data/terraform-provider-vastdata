package utils

import (
	//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"errors"
	"fmt"
	"regexp"
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

func ProtectionPolicyTimeIntervalValidation(i interface{}, c cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	matched, _ := regexp.MatchString(`^[1-9]+[0-9]{0,}[DWHMYsm]{1}`, i.(string))
	if !matched {
		return diag.FromErr(errors.New(fmt.Sprintf("The value given does not match the format <integer><time period>  time period can be D - Days ,W - Weeks ,s - Seconds ,m - Minutes, H - Hours, M - Months, Y - Years , Ex 1D = 1 Day")))
	}
	return diags
}

func ValidateRetention(i interface{}, c cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	b, e := regexp.MatchString(`^[1-9](\d)*[h|m|d|y]$`, i.(string))
	if !b {
		return diag.FromErr(e)
	}
	return diags

}

func ValidateStringListMembers(s []string, allow_duplicates bool) schema.SchemaValidateDiagFunc {
	return func(i interface{}, c cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		l, is_string_list := i.([]string)
		if !is_string_list {
			return diag.FromErr(fmt.Errorf("%v is not a list of strings", i))
		}
		m := map[string]int{}
		for _, j := range s {
			m[j] = 0
		}
		for _, t := range l {
			_, exists := m[t]
			if !exists {
				return diag.FromErr(fmt.Errorf("%v is not a valid value , valid values are, %v", t, s))
			}
		}
		if !allow_duplicates {
			// in the case that we are not allowing duplicates and we are sure that all the values in the given list are valid,
			// if we try to pop more than one time the same value it means it is duplicate.
			for _, t := range l {
				_, exit := m[t]
				if !exit {
					return diag.FromErr(fmt.Errorf("Duplicate value %s was found , duplicates are not allowed", t))
				}
				delete(m, t)
			}

		}
		return diags
	}
}

func ValidateAbeProtocols(i interface{}, c cty.Path) diag.Diagnostics {
	return ValidateStringListMembers([]string{"NFS", "SMB", "NFS4", "S3"}, false)(i, c)
}
