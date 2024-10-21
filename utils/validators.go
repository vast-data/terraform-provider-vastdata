package utils

import (
	//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	grace_period_format = "15:04:05"
	start_at_format     = "2006-01-02 15:04:05"
)

func validTime(s string) error {
	_, err := time.Parse(grace_period_format, s)
	return err

}

func GracePeriodFormatValidation(i interface{}, c cty.Path) diag.Diagnostics {
	//Old sytle quota is composed on of only time like 12:01:23 , while new style is at the format of DD HH:MM:SS
	var diags diag.Diagnostics
	gp := i.(string)
	s := strings.SplitN(gp, " ", 2)
	if len(s) == 2 {
		// in this case the format must be <N HH:MM:SS> where N is an interger larget than 0
		n, err := strconv.Atoi(s[0])
		if err != nil {
			return diag.FromErr(fmt.Errorf("Given number of days %v is not an integer", s[0]))
		}
		if n <= 0 {
			return diag.Errorf("The number of days must be larger than 0")
		}
		err = validTime(s[1])
		if err != nil {
			return diag.FromErr(fmt.Errorf("Wrong time format %v provided  %v", s[1], err))
		}

	} else if len(s) == 1 {
		err := validTime(s[0])
		if err != nil {
			return diag.FromErr(fmt.Errorf("Wrong time format %v provided  %v", s[0], err))
		}

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
