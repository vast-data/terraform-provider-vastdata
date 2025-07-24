// Copyright (c) HashiCorp, Inc.

package schema_generation

import (
	"context"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

func TestGetResourceSchema_EdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		hints         *internalstate.TFStateHints
		expectError   bool
		errorContains string
	}{
		{
			name: "empty_resource_path",
			hints: &internalstate.TFStateHints{
				SchemaRef: &internalstate.SchemaReference{
					Create: &internalstate.OpenAPIEndpointRef{
						Method: http.MethodPost,
						Path:   "", // Empty path
					},
				},
			},
			expectError:   true,
			errorContains: "resource path is required but was empty",
		},
		{
			name: "empty_method",
			hints: &internalstate.TFStateHints{
				SchemaRef: &internalstate.SchemaReference{
					Create: &internalstate.OpenAPIEndpointRef{
						Method: "", // Empty method
						Path:   "users",
					},
				},
			},
			expectError:   true,
			errorContains: "resource method is required but was empty",
		},
		{
			name: "unsupported_method",
			hints: &internalstate.TFStateHints{
				SchemaRef: &internalstate.SchemaReference{
					Create: &internalstate.OpenAPIEndpointRef{
						Method: "DELETE", // Unsupported for create
						Path:   "users",
					},
				},
			},
			expectError:   true,
			errorContains: "unsupported method",
		},
		{
			name: "nil_schema_ref",
			hints: &internalstate.TFStateHints{
				SchemaRef: nil,
			},
			expectError: true,
		},
		{
			name: "custom_schema_empty_attributes",
			hints: &internalstate.TFStateHints{
				TFStateHintsForCustom: &internalstate.TFStateHintsForCustom{
					Description:      "Test custom resource",
					SchemaAttributes: map[string]any{}, // Empty attributes
				},
			},
			expectError:   true,
			errorContains: "custom datasource schema attributes are required but were empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetResourceSchema(context.Background(), tt.hints)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetDatasourceSchema_EdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		hints         *internalstate.TFStateHints
		expectError   bool
		errorContains string
	}{
		{
			name: "empty_resource_path",
			hints: &internalstate.TFStateHints{
				SchemaRef: &internalstate.SchemaReference{
					Read: &internalstate.OpenAPIEndpointRef{
						Method: http.MethodGet,
						Path:   "", // Empty path
					},
				},
			},
			expectError:   true,
			errorContains: "resource path is required but was empty",
		},
		{
			name: "unsupported_method_for_datasource",
			hints: &internalstate.TFStateHints{
				SchemaRef: &internalstate.SchemaReference{
					Read: &internalstate.OpenAPIEndpointRef{
						Method: http.MethodPost, // POST not supported for datasource
						Path:   "users",
					},
				},
			},
			expectError:   true,
			errorContains: "not supported resource method",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetDatasourceSchema(context.Background(), tt.hints)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBuildResourceAttribute_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		schema    *openapi3.Schema
		entry     *SchemaEntry
		hints     *internalstate.TFStateHints
		expectNil bool
	}{
		{
			name: "nil_schema_type",
			schema: &openapi3.Schema{
				Type: nil, // This should cause a panic
			},
			entry: &SchemaEntry{
				Required: true,
			},
			hints:     &internalstate.TFStateHints{},
			expectNil: true, // Function should handle this gracefully or panic
		},
		{
			name: "empty_schema_type",
			schema: &openapi3.Schema{
				Type: &openapi3.Types{}, // Empty type array
			},
			entry: &SchemaEntry{
				Required: true,
			},
			hints:     &internalstate.TFStateHints{},
			expectNil: true,
		},
		{
			name: "object_with_no_properties",
			schema: &openapi3.Schema{
				Type:       &openapi3.Types{openapi3.TypeObject},
				Properties: map[string]*openapi3.SchemaRef{}, // No properties
			},
			entry: &SchemaEntry{
				Optional: true,
			},
			hints:     &internalstate.TFStateHints{},
			expectNil: true, // Should return nil for empty objects
		},
		{
			name: "array_with_empty_object_items",
			schema: &openapi3.Schema{
				Type: &openapi3.Types{openapi3.TypeArray},
				Items: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:       &openapi3.Types{openapi3.TypeObject},
						Properties: map[string]*openapi3.SchemaRef{}, // Empty object
					},
				},
			},
			entry: &SchemaEntry{
				Optional: true,
			},
			hints:     &internalstate.TFStateHints{},
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.expectNil {
					t.Errorf("Unexpected panic: %v", r)
				}
			}()

			result := buildResourceAttribute(context.Background(), "test_field", tt.schema, tt.entry, tt.hints)

			if tt.expectNil {
				require.Nil(t, result, "Expected nil result for edge case")
			} else {
				require.NotNil(t, result, "Expected non-nil result")
			}
		})
	}
}

func TestSchemaEntry_String(t *testing.T) {
	entry := &SchemaEntry{
		Required:    true,
		Optional:    false,
		Computed:    true,
		WriteOnly:   false,
		Sensitive:   true,
		Ordered:     false,
		Description: "Test field",
	}

	result := entry.String()
	require.Contains(t, result, "required=true")
	require.Contains(t, result, "optional=false")
	require.Contains(t, result, "computed=true")
	require.Contains(t, result, "sensitive=true")
}

func TestAddSchemaEntries_WithExclusions(t *testing.T) {
	props := map[string]*openapi3.SchemaRef{
		"included_field": {
			Value: &openapi3.Schema{
				Type:        &openapi3.Types{openapi3.TypeString},
				Description: "This should be included",
			},
		},
		"excluded_field": {
			Value: &openapi3.Schema{
				Type:        &openapi3.Types{openapi3.TypeString},
				Description: "This should be excluded",
			},
		},
	}

	hints := &internalstate.TFStateHints{
		ExcludedSchemaFields: []string{"excluded_field"},
	}

	target := make(map[string]*SchemaEntry)

	addSchemaEntries(props, []string{}, hints, target, false, true, false, false, false, false)

	require.Contains(t, target, "included_field")
	require.NotContains(t, target, "excluded_field")
}

func TestIsEmptySchema(t *testing.T) {
	tests := []struct {
		name     string
		schema   *openapi3.SchemaRef
		expected bool
	}{
		{
			name:     "nil_schema_ref",
			schema:   nil,
			expected: true,
		},
		{
			name: "nil_schema_value",
			schema: &openapi3.SchemaRef{
				Value: nil,
			},
			expected: true,
		},
		{
			name: "empty_schema",
			schema: &openapi3.SchemaRef{
				Value: &openapi3.Schema{},
			},
			expected: true,
		},
		{
			name: "schema_with_type",
			schema: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{openapi3.TypeString},
				},
			},
			expected: false,
		},
		{
			name: "schema_with_properties",
			schema: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Properties: map[string]*openapi3.SchemaRef{
						"test": {Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}}},
					},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmptySchema(tt.schema)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestFlagsFromHintsForResource_OverrideLogic(t *testing.T) {
	hints := &internalstate.TFStateHints{
		RequiredSchemaFields:    []string{"field1"},
		OptionalSchemaFields:    []string{"field2"},
		ComputedSchemaFields:    []string{"field3"},
		NotRequiredSchemaFields: []string{"field4"},
		NotOptionalSchemaFields: []string{"field5"},
		WriteOnlyFields:         []string{"field6"},
		SensitiveFields:         []string{"field7"},
		PreserveOrderFields:     []string{"field8"},
	}

	tests := []struct {
		name              string
		fieldName         string
		initialRequired   bool
		initialOptional   bool
		initialComputed   bool
		expectedRequired  bool
		expectedOptional  bool
		expectedComputed  bool
		expectedSensitive bool
		expectedOrdered   bool
		expectedWriteOnly bool
	}{
		{
			name:              "required_field_override",
			fieldName:         "field1",
			initialRequired:   false,
			initialOptional:   true,
			initialComputed:   false,
			expectedRequired:  true,
			expectedOptional:  false, // Should be false when required is true
			expectedComputed:  false, // Should be false when required is true
			expectedSensitive: false,
			expectedOrdered:   false,
			expectedWriteOnly: false,
		},
		{
			name:              "optional_field_override",
			fieldName:         "field2",
			initialRequired:   true,
			initialOptional:   false,
			initialComputed:   false,
			expectedRequired:  false, // Should be overridden
			expectedOptional:  true,
			expectedComputed:  false,
			expectedSensitive: false,
			expectedOrdered:   false,
			expectedWriteOnly: false,
		},
		{
			name:              "computed_field_override",
			fieldName:         "field3",
			initialRequired:   false,
			initialOptional:   false,
			initialComputed:   false,
			expectedRequired:  false,
			expectedOptional:  false,
			expectedComputed:  true,
			expectedSensitive: false,
			expectedOrdered:   false,
			expectedWriteOnly: false,
		},
		{
			name:              "write_only_field",
			fieldName:         "field6",
			initialRequired:   false,
			initialOptional:   false,
			initialComputed:   true,
			expectedRequired:  false,
			expectedOptional:  true,  // Write-only fields should be optional
			expectedComputed:  false, // Write-only fields cannot be computed
			expectedSensitive: false,
			expectedOrdered:   false,
			expectedWriteOnly: true,
		},
		{
			name:              "sensitive_field",
			fieldName:         "field7",
			initialRequired:   false,
			initialOptional:   true,
			initialComputed:   false,
			expectedRequired:  false,
			expectedOptional:  true,
			expectedComputed:  false,
			expectedSensitive: true,
			expectedOrdered:   false,
			expectedWriteOnly: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			required, optional, computed, writeOnly, sensitive, ordered := flagsFromHintsForResource(
				tt.fieldName, hints,
				tt.initialRequired, tt.initialOptional, tt.initialComputed,
				tt.expectedSensitive, tt.expectedOrdered, tt.expectedWriteOnly,
			)

			require.Equal(t, tt.expectedRequired, required, "Required flag mismatch")
			require.Equal(t, tt.expectedOptional, optional, "Optional flag mismatch")
			require.Equal(t, tt.expectedComputed, computed, "Computed flag mismatch")
			require.Equal(t, tt.expectedWriteOnly, writeOnly, "WriteOnly flag mismatch")
			require.Equal(t, tt.expectedSensitive, sensitive, "Sensitive flag mismatch")
			require.Equal(t, tt.expectedOrdered, ordered, "Ordered flag mismatch")
		})
	}
}

func TestBuildAttrTypeFromSchema_ComplexTypes(t *testing.T) {
	tests := []struct {
		name        string
		schema      *openapi3.Schema
		expectPanic bool
	}{
		{
			name: "valid_string_type",
			schema: &openapi3.Schema{
				Type: &openapi3.Types{openapi3.TypeString},
			},
			expectPanic: false,
		},
		{
			name: "valid_array_type",
			schema: &openapi3.Schema{
				Type: &openapi3.Types{openapi3.TypeArray},
				Items: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeString},
					},
				},
			},
			expectPanic: false,
		},
		{
			name: "array_missing_items",
			schema: &openapi3.Schema{
				Type:  &openapi3.Types{openapi3.TypeArray},
				Items: nil, // Missing items should cause panic
			},
			expectPanic: true,
		},
		{
			name:        "nil_schema",
			schema:      nil,
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.expectPanic && r == nil {
					t.Error("Expected panic but didn't get one")
				} else if !tt.expectPanic && r != nil {
					t.Errorf("Unexpected panic: %v", r)
				}
			}()

			result := buildAttrTypeFromSchema(tt.schema)
			if !tt.expectPanic {
				require.NotNil(t, result)
			}
		})
	}
}
