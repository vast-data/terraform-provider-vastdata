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
			"gid":            rschema.Int64Attribute{Optional: true},
			"uid":            rschema.Int64Attribute{Optional: true},
			"s3_policy_id":   rschema.Int64Attribute{Optional: true},
			"s3_policy_guid": rschema.StringAttribute{Optional: true, Computed: true},
			"ignore_present": rschema.BoolAttribute{Optional: true, Computed: true},
			"context":        rschema.StringAttribute{Optional: true},
			"tenant_id":      rschema.Int64Attribute{Optional: true},
		},
	}

	// Ensure all schema fields have values (null if not provided)
	fullRawValues := map[string]attr.Value{
		"gid":            types.Int64Null(),
		"uid":            types.Int64Null(),
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
			name: "valid_with_uid_and_s3_policy_guid",
			rawValues: map[string]attr.Value{
				"uid":            types.Int64Value(2001),
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
			name: "invalid_both_user_and_group",
			rawValues: map[string]attr.Value{
				"gid":          types.Int64Value(1001),
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
