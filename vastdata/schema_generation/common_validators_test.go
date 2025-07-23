// Copyright (c) HashiCorp, Inc.

package schema_generation

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestGracePeriodValidator(t *testing.T) {
	t.Parallel()

	validatorFn := GracePeriodValidator()

	cases := []struct {
		name     string
		input    string
		wantErr  bool
		errorMsg string
	}{
		{"Valid Old Format", "01:02:03", false, ""},
		{"Valid New Format", "1 12:34:56", false, ""},
		{"Valid New Format Multi-Digit Day", "10 23:59:59", false, ""},

		{"Invalid Empty String", "", true, "Expected HH:MM:SS"},
		{"Invalid Format", "abc", true, "Expected HH:MM:SS"},
		{"Invalid Time Format in New Style", "2 99:00:00", true, "Expected HH:MM:SS"},
		{"Zero Days Not Allowed", "0 01:00:00", true, "Expected a positive integer"},
		{"Negative Days Not Allowed", "-1 01:00:00", true, "Expected a positive integer"},
		{"Invalid Day Token", "one 01:00:00", true, "Expected a positive integer"},
		{"Extra Space", "1 12:00:00 extra", true, "Expected HH:MM:SS"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: types.StringValue(tc.input),
				Path:        path.Root("grace_period"),
			}

			var resp validator.StringResponse
			validatorFn.ValidateString(context.Background(), req, &resp)

			if tc.wantErr {
				require.True(t, resp.Diagnostics.HasError(), "expected error but got none")
				if tc.errorMsg != "" {
					require.Contains(t, resp.Diagnostics.Errors()[0].Detail(), tc.errorMsg)
				}
			} else {
				require.False(t, resp.Diagnostics.HasError(), "expected no error but got one")
			}
		})
	}
}
