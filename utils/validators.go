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

func ValidateManagerPassword(i interface{}, c cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	s := fmt.Sprintf("%v", i)
	if len(s) < 8 {
		return diag.Errorf("Password length should be larger than 8 characters")
	}
	for k, v := range map[string]string{
		`[A-Z]`:                                "Password should contain at least one upper case character",
		`[a-z]`:                                "Password should contain at least one lower case character",
		`[0-9]`:                                "Password should contain at least one number",
		`[\!\@\#\$\%\^\&\*\(\)\~\{\}\[\]\.\?]`: "Password should contain at least one special character"} {
		o, _ := regexp.Match(k, []byte(s))
		if !o {
			return diag.Errorf(v)
		}

	}
	return diags
}
