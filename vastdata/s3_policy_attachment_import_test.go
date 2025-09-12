// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestS3PolicyAttachment_ImportResourceState(t *testing.T) {
	tests := []struct {
		name           string
		importID       string
		expectError    bool
		errorMsg       string
		expectedValues map[string]interface{}
	}{
		{
			name:        "valid_import_with_gid_and_s3_policy_id",
			importID:    "s3_policy_id=1,gid=1000",
			expectError: false,
			expectedValues: map[string]interface{}{
				"s3_policy_id": int64(1),
				"gid":          int64(1000),
			},
		},
		{
			name:        "valid_import_with_uid_and_s3_policy_id",
			importID:    "s3_policy_id=2,uid=2000",
			expectError: false,
			expectedValues: map[string]interface{}{
				"s3_policy_id": int64(2),
				"uid":          int64(2000),
			},
		},
		{
			name:        "valid_import_with_context_and_tenant_id",
			importID:    "s3_policy_id=1,gid=1000,context=ldap,tenant_id=5",
			expectError: false,
			expectedValues: map[string]interface{}{
				"s3_policy_id": int64(1),
				"gid":          int64(1000),
				"context":      "ldap",
				"tenant_id":    int64(5),
			},
		},
		{
			name:        "valid_import_with_s3_policy_guid",
			importID:    "s3_policy_guid=policy-guid-123,uid=3000",
			expectError: false,
			expectedValues: map[string]interface{}{
				"s3_policy_guid": "policy-guid-123",
				"uid":            int64(3000),
			},
		},
		{
			name:        "invalid_import_missing_policy",
			importID:    "gid=1000",
			expectError: true,
			errorMsg:    "one of",
		},
		{
			name:        "invalid_import_missing_user_group",
			importID:    "s3_policy_id=1",
			expectError: true,
			errorMsg:    "one of",
		},
		{
			name:        "invalid_import_both_user_and_group",
			importID:    "s3_policy_id=1,gid=1000,uid=2000",
			expectError: true,
			errorMsg:    "only one of",
		},
		{
			name:        "invalid_import_both_policy_id_and_guid",
			importID:    "s3_policy_id=1,s3_policy_guid=policy-guid-123,gid=1000",
			expectError: true,
			errorMsg:    "only one of",
		},
		{
			name:        "invalid_import_format",
			importID:    "invalid_format",
			expectError: true,
			errorMsg:    "not present in the resource schema",
		},
		{
			name:        "invalid_field_name",
			importID:    "invalid_field=123,gid=1000",
			expectError: true,
			errorMsg:    "not present in the resource schema",
		},
		{
			name:        "invalid_integer_value",
			importID:    "s3_policy_id=not_a_number,gid=1000",
			expectError: true,
			errorMsg:    "invalid value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock S3PolicyAttachment with empty state
			attachment := createTestS3PolicyAttachment(map[string]attr.Value{})

			// Test the import parsing logic by calling parseImportId directly
			parseErr := parseImportId(tt.importID, attachment.tfstate)

			// Test validation after parsing (if parsing succeeded)
			var validationErr error
			if parseErr == nil {
				validationErr = attachment.validateS3PolicyAttachmentConfig()
			}

			// Determine if we expect an error from either parsing or validation
			err := parseErr
			if err == nil {
				err = validationErr
			}

			if tt.expectError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				return
			}

			require.NoError(t, parseErr, "Parsing should succeed")
			require.NoError(t, validationErr, "Validation should succeed")

			// Verify the expected values were set correctly
			for key, expectedValue := range tt.expectedValues {
				assert.True(t, attachment.tfstate.IsKnownAndNotNull(key), "Field %s should be set", key)

				switch v := expectedValue.(type) {
				case int64:
					assert.Equal(t, v, attachment.tfstate.Int64(key), "Field %s should have correct int64 value", key)
				case string:
					assert.Equal(t, v, attachment.tfstate.String(key), "Field %s should have correct string value", key)
				}
			}
		})
	}
}
