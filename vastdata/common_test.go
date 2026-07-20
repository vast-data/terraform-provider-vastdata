// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

func TestPopulateIDFieldsFromNestedObjects(t *testing.T) {
	tests := []struct {
		name           string
		tfstateKeys    map[string]attr.Value
		record         Record
		expectedRecord Record
		description    string
	}{
		{
			name: "extract local_provider_id from local_provider",
			tfstateKeys: map[string]attr.Value{
				"id":                types.Int64Value(123),
				"name":              types.StringValue("test-group"),
				"local_provider_id": types.Int64Null(),
			},
			record: Record{
				"id":   int64(123),
				"name": "test-group",
				"local_provider": map[string]any{
					"id":   int64(5),
					"name": "provider-name",
				},
			},
			expectedRecord: Record{
				"id":   int64(123),
				"name": "test-group",
				"local_provider": map[string]any{
					"id":   int64(5),
					"name": "provider-name",
				},
				"local_provider_id": int64(5),
			},
			description: "Should extract local_provider_id from local_provider.id",
		},
		{
			name: "extract tenant_id from tenant",
			tfstateKeys: map[string]attr.Value{
				"id":        types.Int64Value(456),
				"name":      types.StringValue("test-role"),
				"tenant_id": types.Int64Null(),
			},
			record: Record{
				"id":   int64(456),
				"name": "test-role",
				"tenant": map[string]any{
					"id":   int64(28),
					"name": "tenant-name",
					"guid": "abc-123",
				},
			},
			expectedRecord: Record{
				"id":   int64(456),
				"name": "test-role",
				"tenant": map[string]any{
					"id":   int64(28),
					"name": "tenant-name",
					"guid": "abc-123",
				},
				"tenant_id": int64(28),
			},
			description: "Should extract tenant_id from tenant.id",
		},
		{
			name: "multiple _id fields",
			tfstateKeys: map[string]attr.Value{
				"id":                 types.Int64Value(789),
				"name":               types.StringValue("test-snapshot"),
				"loanee_tenant_id":   types.Int64Null(),
				"remote_target_id":   types.Int64Null(),
				"loanee_snapshot_id": types.Int64Null(),
			},
			record: Record{
				"id":   int64(789),
				"name": "test-snapshot",
				"loanee_tenant": map[string]any{
					"id":   int64(1),
					"name": "default",
					"guid": "xyz-456",
				},
				"remote_target": map[string]any{
					"id":   int64(10),
					"name": "remote-cluster",
				},
				"loanee_snapshot": map[string]any{
					"id":   int64(42),
					"name": "snap-1",
				},
			},
			expectedRecord: Record{
				"id":   int64(789),
				"name": "test-snapshot",
				"loanee_tenant": map[string]any{
					"id":   int64(1),
					"name": "default",
					"guid": "xyz-456",
				},
				"loanee_tenant_id": int64(1),
				"remote_target": map[string]any{
					"id":   int64(10),
					"name": "remote-cluster",
				},
				"remote_target_id": int64(10),
				"loanee_snapshot": map[string]any{
					"id":   int64(42),
					"name": "snap-1",
				},
				"loanee_snapshot_id": int64(42),
			},
			description: "Should extract multiple _id fields from nested objects",
		},
		{
			name: "id field already exists in record",
			tfstateKeys: map[string]attr.Value{
				"id":                types.Int64Value(111),
				"name":              types.StringValue("test-existing"),
				"local_provider_id": types.Int64Value(99),
			},
			record: Record{
				"id":                int64(111),
				"name":              "test-existing",
				"local_provider_id": int64(99), // Already exists
				"local_provider": map[string]any{
					"id":   int64(5),
					"name": "provider-name",
				},
			},
			expectedRecord: Record{
				"id":                int64(111),
				"name":              "test-existing",
				"local_provider_id": int64(99), // Should remain unchanged
				"local_provider": map[string]any{
					"id":   int64(5),
					"name": "provider-name",
				},
			},
			description: "Should not overwrite existing _id field",
		},
		{
			name: "parent object does not exist",
			tfstateKeys: map[string]attr.Value{
				"id":                types.Int64Value(222),
				"name":              types.StringValue("test-no-parent"),
				"local_provider_id": types.Int64Null(),
			},
			record: Record{
				"id":   int64(222),
				"name": "test-no-parent",
				// No local_provider object
			},
			expectedRecord: Record{
				"id":   int64(222),
				"name": "test-no-parent",
				// local_provider_id should not be added
			},
			description: "Should not add _id field if parent object doesn't exist",
		},
		{
			name: "parent object is null",
			tfstateKeys: map[string]attr.Value{
				"id":                types.Int64Value(333),
				"name":              types.StringValue("test-null-parent"),
				"local_provider_id": types.Int64Null(),
			},
			record: Record{
				"id":             int64(333),
				"name":           "test-null-parent",
				"local_provider": nil,
			},
			expectedRecord: Record{
				"id":             int64(333),
				"name":           "test-null-parent",
				"local_provider": nil,
			},
			description: "Should not add _id field if parent object is nil",
		},
		{
			name: "parent object is not a map",
			tfstateKeys: map[string]attr.Value{
				"id":                types.Int64Value(444),
				"name":              types.StringValue("test-invalid-parent"),
				"local_provider_id": types.Int64Null(),
			},
			record: Record{
				"id":             int64(444),
				"name":           "test-invalid-parent",
				"local_provider": "not-a-map",
			},
			expectedRecord: Record{
				"id":             int64(444),
				"name":           "test-invalid-parent",
				"local_provider": "not-a-map",
			},
			description: "Should not add _id field if parent is not a map",
		},
		{
			name: "parent object has no id field",
			tfstateKeys: map[string]attr.Value{
				"id":                types.Int64Value(555),
				"name":              types.StringValue("test-no-id"),
				"local_provider_id": types.Int64Null(),
			},
			record: Record{
				"id":   int64(555),
				"name": "test-no-id",
				"local_provider": map[string]any{
					"name": "provider-name",
					// No "id" field
				},
			},
			expectedRecord: Record{
				"id":   int64(555),
				"name": "test-no-id",
				"local_provider": map[string]any{
					"name": "provider-name",
				},
			},
			description: "Should not add _id field if parent has no id field",
		},
		{
			name: "parent id field is nil",
			tfstateKeys: map[string]attr.Value{
				"id":                types.Int64Value(666),
				"name":              types.StringValue("test-nil-id"),
				"local_provider_id": types.Int64Null(),
			},
			record: Record{
				"id":   int64(666),
				"name": "test-nil-id",
				"local_provider": map[string]any{
					"id":   nil,
					"name": "provider-name",
				},
			},
			expectedRecord: Record{
				"id":   int64(666),
				"name": "test-nil-id",
				"local_provider": map[string]any{
					"id":   nil,
					"name": "provider-name",
				},
			},
			description: "Should not add _id field if parent id is nil",
		},
		{
			name: "no _id fields in tfstate",
			tfstateKeys: map[string]attr.Value{
				"id":   types.Int64Value(777),
				"name": types.StringValue("test-no-id-fields"),
			},
			record: Record{
				"id":   int64(777),
				"name": "test-no-id-fields",
				"local_provider": map[string]any{
					"id":   int64(5),
					"name": "provider-name",
				},
			},
			expectedRecord: Record{
				"id":   int64(777),
				"name": "test-no-id-fields",
				"local_provider": map[string]any{
					"id":   int64(5),
					"name": "provider-name",
				},
			},
			description: "Should not modify record if no _id fields in tfstate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a minimal tfstate with the test keys
			tfstate := &is.TFState{
				Raw: tt.tfstateKeys,
			}

			// Make a copy of the record for testing
			recordCopy := make(Record)
			for k, v := range tt.record {
				recordCopy[k] = v
			}

			// Call the function
			ctx := context.Background()
			PopulateIDFieldsFromNestedObjects(ctx, tfstate, recordCopy)

			// Verify the result
			for expectedKey, expectedValue := range tt.expectedRecord {
				actualValue, exists := recordCopy[expectedKey]
				if !exists {
					if expectedValue != nil {
						t.Errorf("%s: expected key %q to exist in record, but it doesn't", tt.description, expectedKey)
					}
					continue
				}

				// Use deep equality for comparison (handles maps, slices, etc.)
				if !reflect.DeepEqual(expectedValue, actualValue) {
					t.Errorf("%s: for key %q, expected value %v, got %v", tt.description, expectedKey, expectedValue, actualValue)
				}
			}

			// Verify no unexpected keys were added
			for actualKey := range recordCopy {
				if _, expected := tt.expectedRecord[actualKey]; !expected {
					t.Errorf("%s: unexpected key %q was added to record", tt.description, actualKey)
				}
			}
		})
	}
}

func TestPopulateIDFieldsFromNestedObjects_NilInputs(t *testing.T) {
	ctx := context.Background()

	t.Run("nil record", func(t *testing.T) {
		tfstate := &is.TFState{
			Raw: map[string]attr.Value{
				"local_provider_id": types.Int64Null(),
			},
		}
		// Should not panic
		PopulateIDFieldsFromNestedObjects(ctx, tfstate, nil)
	})

	t.Run("nil tfstate", func(t *testing.T) {
		record := Record{
			"id": int64(123),
		}
		// Should not panic
		PopulateIDFieldsFromNestedObjects(ctx, nil, record)
	})

	t.Run("both nil", func(t *testing.T) {
		// Should not panic
		PopulateIDFieldsFromNestedObjects(ctx, nil, nil)
	})
}

func TestIsTransientAsyncPollError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"canceled", context.Canceled, false},
		{"deadline", context.DeadlineExceeded, false},
		{"verification failed", fmt.Errorf("WaitAPICondition verification failed: boom"), false},
		{"timeout", fmt.Errorf("WaitAPICondition timeout after 10m0s"), false},
		{"cancelled", fmt.Errorf("WaitAPICondition cancelled: %w", context.Canceled), false},
		{"generic", errors.New("task failed"), false},
		{"api 400", &ApiError{StatusCode: 400, Body: "bad"}, false},
		{"api 503", &ApiError{StatusCode: 503, Body: "unavailable"}, true},
		{"api unreachable", &ApiError{StatusCode: 0, Body: "unreachable"}, true},
		{
			"wrapped connection reset",
			fmt.Errorf("WaitAPICondition API call failed: %w",
				fmt.Errorf(`failed to perform GET request to https://v162:443/api/latest/vtasks/20/, error Get "https://v162:443/api/latest/vtasks/20/": read tcp 10.241.12.3:62794->10.141.200.162:443: connection reset by peer`)),
			true,
		},
		{"eof", io.EOF, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTransientAsyncPollError(tt.err); got != tt.want {
				t.Errorf("isTransientAsyncPollError() = %v, want %v", got, tt.want)
			}
		})
	}
}
