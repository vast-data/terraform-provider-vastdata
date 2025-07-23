// Copyright (c) HashiCorp, Inc.

package schema_generation

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	ValidatorPathStartsEndsWithSlash = "path_starts_ends_with_slash"
	ValidatorPathStartsWithSlash     = "path_starts_with_slash"
	ValidatorGracePeriodFormat       = "grace_period_format"
	ValidatorRFC3339Format           = "rfc3339_format"
	ValidatorIntervalFormat          = "interval_format"  // e.g., 5D, 1W
	ValidatorRetentionFormat         = "retention_format" // e.g., 12h, 3y
)

var commonStringValidators = map[string][]validator.String{
	ValidatorPathStartsEndsWithSlash: {
		stringvalidator.RegexMatches(
			regexp.MustCompile(`^/([^/].*/)?$`),
			"must start and end with '/'",
		),
	},
	ValidatorPathStartsWithSlash: {
		stringvalidator.RegexMatches(
			regexp.MustCompile(`^/.*`),
			"must start with '/'",
		),
	},
	ValidatorGracePeriodFormat: {
		GracePeriodValidator(),
	},
	ValidatorRFC3339Format: {
		stringvalidator.RegexMatches(
			regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
			"must be a valid RFC3339 timestamp (e.g., 2025-07-06T13:45:00Z)",
		),
	},
	ValidatorIntervalFormat: {
		stringvalidator.RegexMatches(
			regexp.MustCompile(`^[1-9][0-9]*[DWHMYsm]$`),
			"must match format <integer><unit> (D/W/H/M/Y/s/m), e.g. 1D or 5H",
		),
	},
	ValidatorRetentionFormat: {
		stringvalidator.RegexMatches(
			regexp.MustCompile(`^[1-9][0-9]*[hmdy]$`),
			"must match format <integer><unit> (h/m/d/y), e.g. 7d or 12h",
		),
	},
}

var commonIntValidators = map[string][]validator.Int64{}
var commonFloatValidators = map[string][]validator.Float64{}

// -------------------------
// GracePeriodValidator
// -------------------------

const gracePeriodTimeFormat = "15:04:05"

type gracePeriodValidator struct{}

func (v gracePeriodValidator) Description(ctx context.Context) string {
	return "String must be in HH:MM:SS or <days> HH:MM:SS format."
}

func (v gracePeriodValidator) MarkdownDescription(ctx context.Context) string {
	return "String must be in `HH:MM:SS` or `<days> HH:MM:SS` format, where `<days>` is a positive integer."
}

func (v gracePeriodValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	val := req.ConfigValue
	if val.IsNull() || val.IsUnknown() {
		return
	}

	s := val.ValueString()
	parts := strings.SplitN(s, " ", 2)

	if len(parts) == 2 {
		days, err := strconv.Atoi(parts[0])
		if err != nil || days <= 0 {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Days in Grace Period",
				fmt.Sprintf("Expected a positive integer for days, got: %q", parts[0]),
			)
			return
		}
		if err := validateHHMMSS(parts[1]); err != nil {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Time Format",
				fmt.Sprintf("Expected HH:MM:SS, got %q: %s", parts[1], err),
			)
		}
	} else if len(parts) == 1 {
		if err := validateHHMMSS(parts[0]); err != nil {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Time Format",
				fmt.Sprintf("Expected HH:MM:SS, got %q: %s", parts[0], err),
			)
		}
	}
}

func validateHHMMSS(s string) error {
	_, err := time.Parse(gracePeriodTimeFormat, s)
	return err
}

func GracePeriodValidator() validator.String {
	return gracePeriodValidator{}
}
