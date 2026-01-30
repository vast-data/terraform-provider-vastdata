// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

// Mock implementations are commented out due to interface compatibility issues.
// The tests below focus on validation logic which is the core functionality.

// Helper to create S3PolicyAttachment with test data
func createTestS3PolicyAttachment(rawValues map[string]attr.Value) *S3PolicyAttachment {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"group_sid":      rschema.StringAttribute{Optional: true},
			"groupname":      rschema.StringAttribute{Optional: true},
			"gid":            rschema.Int64Attribute{Optional: true},
			"uid":            rschema.Int64Attribute{Optional: true},
			"sid":            rschema.StringAttribute{Optional: true},
			"username":       rschema.StringAttribute{Optional: true},
			"s3_policy_id":   rschema.Int64Attribute{Optional: true},
			"s3_policy_guid": rschema.StringAttribute{Optional: true, Computed: true},
			"ignore_present": rschema.BoolAttribute{Optional: true, Computed: true},
			"context":        rschema.StringAttribute{Optional: true},
			"tenant_id":      rschema.Int64Attribute{Optional: true},
		},
	}

	// Ensure all schema fields have values (null if not provided)
	fullRawValues := map[string]attr.Value{
		"group_sid":      types.StringNull(),
		"groupname":      types.StringNull(),
		"gid":            types.Int64Null(),
		"uid":            types.Int64Null(),
		"sid":            types.StringNull(),
		"username":       types.StringNull(),
		"s3_policy_id":   types.Int64Null(),
		"s3_policy_guid": types.StringNull(),
		"ignore_present": types.BoolNull(),
		"context":        types.StringNull(),
		"tenant_id":      types.Int64Null(),
	}

	// Override with provided values
	for k, v := range rawValues {
		fullRawValues[k] = v
	}

	return &S3PolicyAttachment{
		tfstate: is.NewTFStateMust(fullRawValues, schema, nil),
	}
}

func TestS3PolicyAttachment_validateS3PolicyAttachmentConfig(t *testing.T) {
	tests := []struct {
		name      string
		rawValues map[string]attr.Value
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid_with_gid_and_s3_policy_id",
			rawValues: map[string]attr.Value{
				"gid":          types.Int64Value(1001),
				"s3_policy_id": types.Int64Value(100),
			},
			expectErr: false,
		},
		{
			name: "valid_with_groupname_and_s3_policy_id",
			rawValues: map[string]attr.Value{
				"groupname":    types.StringValue("test-group"),
				"s3_policy_id": types.Int64Value(100),
			},
			expectErr: false,
		},
		{
			name: "valid_with_uid_and_s3_policy_guid",
			rawValues: map[string]attr.Value{
				"uid":            types.Int64Value(2001),
				"s3_policy_guid": types.StringValue("policy-guid-123"),
			},
			expectErr: false,
		},
		{
			name: "valid_with_username_and_s3_policy_guid",
			rawValues: map[string]attr.Value{
				"username":       types.StringValue("test-user"),
				"s3_policy_guid": types.StringValue("policy-guid-123"),
			},
			expectErr: false,
		},
		{
			name: "invalid_missing_user_group",
			rawValues: map[string]attr.Value{
				"s3_policy_id": types.Int64Value(100),
			},
			expectErr: true,
			errMsg:    "one of",
		},
		{
			name: "invalid_mixing_user_and_group_fields",
			rawValues: map[string]attr.Value{
				"gid":          types.Int64Value(1001),
				"uid":          types.Int64Value(2001),
				"s3_policy_id": types.Int64Value(100),
			},
			expectErr: true,
			errMsg:    "only one of",
		},
		{
			name: "invalid_mixing_gid_and_username",
			rawValues: map[string]attr.Value{
				"gid":          types.Int64Value(1001),
				"username":     types.StringValue("test-user"),
				"s3_policy_id": types.Int64Value(100),
			},
			expectErr: true,
			errMsg:    "only one of",
		},
		{
			name: "invalid_mixing_groupname_and_uid",
			rawValues: map[string]attr.Value{
				"groupname":    types.StringValue("test-group"),
				"uid":          types.Int64Value(2001),
				"s3_policy_id": types.Int64Value(100),
			},
			expectErr: true,
			errMsg:    "only one of",
		},
		{
			name: "invalid_missing_policy",
			rawValues: map[string]attr.Value{
				"gid": types.Int64Value(1001),
			},
			expectErr: true,
			errMsg:    "one of",
		},
		{
			name: "invalid_both_policy_id_and_guid",
			rawValues: map[string]attr.Value{
				"gid":            types.Int64Value(1001),
				"s3_policy_id":   types.Int64Value(100),
				"s3_policy_guid": types.StringValue("policy-guid-123"),
			},
			expectErr: true,
			errMsg:    "only one of",
		},
		{
			name: "valid_with_group_sid",
			rawValues: map[string]attr.Value{
				"group_sid":    types.StringValue("S-1-5-21-3775954470-3969645310-4227734646-72157"),
				"s3_policy_id": types.Int64Value(100),
			},
			expectErr: false,
		},
		{
			name: "valid_with_sid_for_user",
			rawValues: map[string]attr.Value{
				"sid":          types.StringValue("S-1-5-21-1234567890-1234567890-1234567890-2000"),
				"s3_policy_id": types.Int64Value(100),
			},
			expectErr: false,
		},
		{
			name: "invalid_mixing_group_sid_and_gid",
			rawValues: map[string]attr.Value{
				"group_sid":    types.StringValue("S-1-5-21-1234567890-1234567890-1234567890-1000"),
				"gid":          types.Int64Value(1001),
				"s3_policy_id": types.Int64Value(100),
			},
			expectErr: true,
			errMsg:    "only one of",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attachment := createTestS3PolicyAttachment(tt.rawValues)
			err := attachment.validateS3PolicyAttachmentConfig()

			if tt.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

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
			name:        "valid_import_with_groupname",
			importID:    "s3_policy_id=1,groupname=test-group",
			expectError: false,
			expectedValues: map[string]interface{}{
				"s3_policy_id": int64(1),
				"groupname":    "test-group",
			},
		},
		{
			name:        "valid_import_with_username",
			importID:    "s3_policy_id=2,username=test-user",
			expectError: false,
			expectedValues: map[string]interface{}{
				"s3_policy_id": int64(2),
				"username":     "test-user",
			},
		},
		{
			name:        "valid_import_with_username_and_context",
			importID:    "s3_policy_id=1,username=test-user,context=ldap,tenant_id=5",
			expectError: false,
			expectedValues: map[string]interface{}{
				"s3_policy_id": int64(1),
				"username":     "test-user",
				"context":      "ldap",
				"tenant_id":    int64(5),
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
			name:        "invalid_import_mixing_gid_and_username",
			importID:    "s3_policy_id=1,gid=1000,username=test-user",
			expectError: true,
			errorMsg:    "only one of",
		},
		{
			name:        "invalid_import_mixing_groupname_and_uid",
			importID:    "s3_policy_id=1,groupname=test-group,uid=2000",
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

func TestS3PolicyAttachment_getSearchParamsFromState(t *testing.T) {
	tests := []struct {
		name            string
		rawValues       map[string]attr.Value
		expectedKey     string
		expectedContext string
		checkValueFunc  func(*testing.T, interface{})
	}{
		{
			name: "group_sid",
			rawValues: map[string]attr.Value{
				"group_sid": types.StringValue("S-1-5-21-3775954470-3969645310-4227734646-72157"),
				"context":   types.StringValue("ad"),
			},
			expectedKey:     "sid",
			expectedContext: "group",
			checkValueFunc: func(t *testing.T, val interface{}) {
				assert.Equal(t, "S-1-5-21-3775954470-3969645310-4227734646-72157", val)
			},
		},
		{
			name: "groupname",
			rawValues: map[string]attr.Value{
				"groupname": types.StringValue("test-group"),
			},
			expectedKey:     "groupname",
			expectedContext: "group",
			checkValueFunc: func(t *testing.T, val interface{}) {
				assert.Equal(t, "test-group", val)
			},
		},
		{
			name: "sid_for_user",
			rawValues: map[string]attr.Value{
				"sid": types.StringValue("S-1-5-21-1234567890-1234567890-1234567890-2000"),
			},
			expectedKey:     "sid",
			expectedContext: "user",
			checkValueFunc: func(t *testing.T, val interface{}) {
				assert.Equal(t, "S-1-5-21-1234567890-1234567890-1234567890-2000", val)
			},
		},
		{
			name: "uid",
			rawValues: map[string]attr.Value{
				"uid": types.Int64Value(2001),
			},
			expectedKey:     "uid",
			expectedContext: "user",
			checkValueFunc: func(t *testing.T, val interface{}) {
				assert.Equal(t, int64(2001), val)
			},
		},
		{
			name: "username",
			rawValues: map[string]attr.Value{
				"username": types.StringValue("test-user"),
			},
			expectedKey:     "username",
			expectedContext: "user",
			checkValueFunc: func(t *testing.T, val interface{}) {
				assert.Equal(t, "test-user", val)
			},
		},
		{
			name: "gid_checked_last",
			rawValues: map[string]attr.Value{
				"gid":     types.Int64Value(1001),
				"context": types.StringValue("ad"),
			},
			expectedKey:     "gid",
			expectedContext: "group",
			checkValueFunc: func(t *testing.T, val interface{}) {
				assert.Equal(t, int64(1001), val)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := rschema.Schema{
				Attributes: map[string]rschema.Attribute{
					"group_sid":      rschema.StringAttribute{Optional: true},
					"groupname":      rschema.StringAttribute{Optional: true},
					"gid":            rschema.Int64Attribute{Optional: true},
					"uid":            rschema.Int64Attribute{Optional: true},
					"sid":            rschema.StringAttribute{Optional: true},
					"username":       rschema.StringAttribute{Optional: true},
					"s3_policy_id":   rschema.Int64Attribute{Optional: true},
					"s3_policy_guid": rschema.StringAttribute{Optional: true, Computed: true},
					"ignore_present": rschema.BoolAttribute{Optional: true, Computed: true},
					"context":        rschema.StringAttribute{Optional: true},
					"tenant_id":      rschema.Int64Attribute{Optional: true},
				},
			}

			// Ensure all schema fields have values (null if not provided)
			fullRawValues := map[string]attr.Value{
				"group_sid":      types.StringNull(),
				"groupname":      types.StringNull(),
				"gid":            types.Int64Null(),
				"uid":            types.Int64Null(),
				"sid":            types.StringNull(),
				"username":       types.StringNull(),
				"s3_policy_id":   types.Int64Null(),
				"s3_policy_guid": types.StringNull(),
				"ignore_present": types.BoolNull(),
				"context":        types.StringNull(),
				"tenant_id":      types.Int64Null(),
			}

			// Override with provided values
			for k, v := range tt.rawValues {
				fullRawValues[k] = v
			}

			attachment := &S3PolicyAttachment{
				tfstate: is.NewTFStateMust(fullRawValues, schema, nil),
			}

			searchParams, attachContext := attachment.getSearchParamsFromState(attachment.tfstate)

			assert.Equal(t, tt.expectedContext, attachContext, "Attach context should match")

			// Check that the expected key is present in searchParams
			val, ok := searchParams[tt.expectedKey]
			require.True(t, ok, "Expected key %s should be present in search params", tt.expectedKey)

			// Use the checkValueFunc to verify the value
			if tt.checkValueFunc != nil {
				tt.checkValueFunc(t, val)
			}
		})
	}
}
