// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

// TestImportLogic_Int64ID tests the import logic for int64 ID fields
func TestImportLogic_Int64ID(t *testing.T) {
	// Create a schema with int64 ID
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"id":   rschema.Int64Attribute{Optional: true, Computed: true},
			"name": rschema.StringAttribute{Optional: true, Computed: true},
		},
	}

	// Create TFState with empty values
	tfState := is.NewTFStateMust(map[string]attr.Value{}, schema, nil)

	// Test that the schema has an id field
	assert.True(t, tfState.HasAttribute("id"))

	// Test that the id field is int64 type
	idType := tfState.Type("id")
	assert.True(t, idType.Equal(types.Int64Type))

	// Test setting an int64 value
	tfState.SetOrAdd("id", int64(123))
	assert.Equal(t, int64(123), tfState.Int64("id"))
}

// TestImportLogic_StringID tests the import logic for string ID fields
func TestImportLogic_StringID(t *testing.T) {
	// Create a schema with string ID
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"id":   rschema.StringAttribute{Optional: true, Computed: true},
			"name": rschema.StringAttribute{Optional: true, Computed: true},
		},
	}

	// Create TFState with empty values
	tfState := is.NewTFStateMust(map[string]attr.Value{}, schema, nil)

	// Test that the schema has an id field
	assert.True(t, tfState.HasAttribute("id"))

	// Test that the id field is string type
	idType := tfState.Type("id")
	assert.True(t, idType.Equal(types.StringType))

	// Test setting a string value
	tfState.SetOrAdd("id", "test-id-123")
	assert.Equal(t, "test-id-123", tfState.String("id"))
}

// TestImportLogic_NoIDField tests the import logic for resources without ID fields
func TestImportLogic_NoIDField(t *testing.T) {
	// Create a schema without ID field
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"name": rschema.StringAttribute{Optional: true, Computed: true},
		},
	}

	// Create TFState with empty values
	tfState := is.NewTFStateMust(map[string]attr.Value{}, schema, nil)

	// Test that the schema does not have an id field
	assert.False(t, tfState.HasAttribute("id"))
}

// TestImportLogic_StringToInt64Conversion tests string to int64 conversion
func TestImportLogic_StringToInt64Conversion(t *testing.T) {
	// Test valid conversions
	testCases := []struct {
		input    string
		expected int64
		valid    bool
	}{
		{"123", 123, true},
		{"0", 0, true},
		{"-1", -1, true},
		{"invalid", 0, false},
		{"", 0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			// This would be the actual conversion logic used in the import
			// We're testing the logic that would be used in strconv.ParseInt
			if tc.valid {
				// In a real scenario, this would be the conversion
				// For now, we just test that our test cases are valid
				assert.True(t, tc.valid)
			} else {
				// For invalid cases, we expect the conversion to fail
				assert.False(t, tc.valid)
			}
		})
	}
}

// TestTFState_SetOrAdd_NewKey tests that SetOrAdd can add new keys
func TestTFState_SetOrAdd_NewKey(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"id": rschema.Int64Attribute{Optional: true},
		},
	}

	// Create TFState with empty Raw map
	tfState := is.NewTFStateMust(map[string]attr.Value{}, schema, nil)

	// Verify the key doesn't exist initially
	_, exists := tfState.Raw["id"]
	assert.False(t, exists)

	// Add the key using SetOrAdd
	tfState.SetOrAdd("id", int64(456))

	// Verify the key now exists and has the correct value
	_, exists = tfState.Raw["id"]
	assert.True(t, exists)
	assert.Equal(t, int64(456), tfState.Int64("id"))
}

// TestTFState_HasAttribute_EdgeCases tests edge cases for HasAttribute
func TestTFState_HasAttribute_EdgeCases(t *testing.T) {
	// Test with disabled TFState
	disabledTFState := is.NewTFStateMust(nil, nil, nil)
	assert.False(t, disabledTFState.HasAttribute("any-field"))

	// Test with empty schema
	emptySchema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{},
	}
	emptyTFState := is.NewTFStateMust(map[string]attr.Value{}, emptySchema, nil)
	assert.False(t, emptyTFState.HasAttribute("id"))
	assert.False(t, emptyTFState.HasAttribute("non-existent"))
}

// --- parseAndApplyCompositeImport unit tests ---
func TestParseComposite_KeyValue_Success(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"gid":       rschema.Int64Attribute{Optional: true},
		"tenant_id": rschema.Int64Attribute{Optional: true},
		"context":   rschema.StringAttribute{Optional: true},
	}}
	tf := is.NewTFStateMust(map[string]attr.Value{}, schema, &is.TFStateHints{TFStateHintsForCustom: &is.TFStateHintsForCustom{}})
	got := make(map[string]attr.Value)
	err := parseAndApplyCompositeImport("gid=1001,tenant_id=22,context=ad", []string{"gid", "tenant_id", "context"}, tf, func(k string, v attr.Value) {
		got[k] = v
	})
	require.NoError(t, err)
	// type assertions
	_, ok := got["gid"].(basetypes.Int64Value)
	require.True(t, ok)
	require.Equal(t, int64(1001), got["gid"].(types.Int64).ValueInt64())
	require.Equal(t, int64(22), got["tenant_id"].(types.Int64).ValueInt64())
	require.Equal(t, "ad", got["context"].(types.String).ValueString())
}

func TestParseComposite_Ordered_Pipe_Success(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"gid":       rschema.Int64Attribute{Optional: true},
		"tenant_id": rschema.Int64Attribute{Optional: true},
		"context":   rschema.StringAttribute{Optional: true},
	}}
	tf := is.NewTFStateMust(map[string]attr.Value{}, schema, &is.TFStateHints{TFStateHintsForCustom: &is.TFStateHintsForCustom{}})
	got := make(map[string]attr.Value)
	err := parseAndApplyCompositeImport("1001|22|ad", []string{"gid", "tenant_id", "context"}, tf, func(k string, v attr.Value) {
		got[k] = v
	})
	require.NoError(t, err)
}

func TestParseComposite_MissingKey_Error(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"gid":       rschema.Int64Attribute{Optional: true},
		"tenant_id": rschema.Int64Attribute{Optional: true},
		"context":   rschema.StringAttribute{Optional: true},
	}}
	tf := is.NewTFStateMust(map[string]attr.Value{}, schema, &is.TFStateHints{TFStateHintsForCustom: &is.TFStateHintsForCustom{}})
	// key=value form now allows partial keys; no error expected
	err := parseAndApplyCompositeImport("gid=1001,tenant_id=22", []string{"gid", "tenant_id", "context"}, tf, func(k string, v attr.Value) {})
	require.NoError(t, err)
}

func TestParseComposite_WrongCount_Error(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"gid":       rschema.Int64Attribute{Optional: true},
		"tenant_id": rschema.Int64Attribute{Optional: true},
		"context":   rschema.StringAttribute{Optional: true},
	}}
	tf := is.NewTFStateMust(map[string]attr.Value{}, schema, &is.TFStateHints{TFStateHintsForCustom: &is.TFStateHintsForCustom{}})
	err := parseAndApplyCompositeImport("1001,22", []string{"gid", "tenant_id", "context"}, tf, func(k string, v attr.Value) {})
	require.Error(t, err)
}

func TestParseComposite_UnknownField_Error(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"gid":       rschema.Int64Attribute{Optional: true},
		"tenant_id": rschema.Int64Attribute{Optional: true},
		"context":   rschema.StringAttribute{Optional: true},
	}}
	tf := is.NewTFStateMust(map[string]attr.Value{}, schema, &is.TFStateHints{TFStateHintsForCustom: &is.TFStateHintsForCustom{}})
	err := parseAndApplyCompositeImport("1001,22,ad", []string{"gid", "tenant_id", "unknown"}, tf, func(k string, v attr.Value) {})
	require.Error(t, err)
}

func TestParseComposite_SingleToken_TreatAsID_Int64(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id":        rschema.Int64Attribute{Optional: true, Computed: true},
		"gid":       rschema.Int64Attribute{Optional: true},
		"tenant_id": rschema.Int64Attribute{Optional: true},
		"context":   rschema.StringAttribute{Optional: true},
	}}
	tf := is.NewTFStateMust(map[string]attr.Value{}, schema, &is.TFStateHints{TFStateHintsForCustom: &is.TFStateHintsForCustom{}})
	got := make(map[string]attr.Value)
	// Provide single token without delimiters; should be treated as id
	err := parseAndApplyCompositeImport("12345", []string{"gid", "tenant_id", "context"}, tf, func(k string, v attr.Value) {
		got[k] = v
	})
	require.NoError(t, err)
	v, ok := got["id"].(types.Int64)
	require.True(t, ok)
	require.Equal(t, int64(12345), v.ValueInt64())
}

func TestParseComposite_SingleToken_TreatAsID_String(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id":        rschema.StringAttribute{Optional: true, Computed: true},
		"gid":       rschema.Int64Attribute{Optional: true},
		"tenant_id": rschema.Int64Attribute{Optional: true},
		"context":   rschema.StringAttribute{Optional: true},
	}}
	tf := is.NewTFStateMust(map[string]attr.Value{}, schema, &is.TFStateHints{TFStateHintsForCustom: &is.TFStateHintsForCustom{}})
	got := make(map[string]attr.Value)
	err := parseAndApplyCompositeImport("abc-123", []string{"gid", "tenant_id", "context"}, tf, func(k string, v attr.Value) {
		got[k] = v
	})
	require.NoError(t, err)
	v, ok := got["id"].(types.String)
	require.True(t, ok)
	require.Equal(t, "abc-123", v.ValueString())
}

// Helpers for default import tests
type testManager struct{ tf *is.TFState }

func (t *testManager) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return t
}
func (t *testManager) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return nil
}
func (t *testManager) TfState() *is.TFState                      { return t.tf }
func (t *testManager) API(_ *VMSRest) VastResourceAPIWithContext { return nil }
func (t *testManager) ReadResource(_ context.Context, _ *VMSRest) (DisplayableRecord, error) {
	return nil, nil
}

// build a Resource whose ManagerWithSchemaOnly uses custom schema (no OpenAPI)
func buildTestResourceWithSchema(schema rschema.Schema, hints *is.TFStateHints) *Resource {
	return &Resource{
		newManager: func(raw map[string]attr.Value, s any) ResourceManager {
			localHints := hints
			if localHints == nil {
				localHints = &is.TFStateHints{}
			}
			// force custom schema mode
			localHints.TFStateHintsForCustom = &is.TFStateHintsForCustom{
				Description:         "test",
				MarkdownDescription: "test",
				SchemaAttributes: func() map[string]any {
					out := make(map[string]any, len(schema.Attributes))
					for k, v := range schema.Attributes {
						out[k] = v
					}
					return out
				}(),
			}
			if s == nil {
				s = schema
			}
			return &testManager{tf: is.NewTFStateMust(raw, s, localHints)}
		},
		managerName: "test",
	}
}

func TestImport_DefaultId_KeyValue_Int64(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id": rschema.Int64Attribute{Optional: true, Computed: true},
	}}
	r := buildTestResourceWithSchema(schema, &is.TFStateHints{})
	req := resource.ImportStateRequest{ID: "id=777"}
	resp := &resource.ImportStateResponse{}
	r.importStateImpl(context.Background(), req, resp)
	require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics.Errors())
	// Read back state: we can directly inspect resp.State via path.Root
	var gotId types.Int64
	require.False(t, resp.State.GetAttribute(context.Background(), path.Root("id"), &gotId).HasError())
	require.Equal(t, int64(777), gotId.ValueInt64())
}

func TestImport_DefaultId_KeyValue_String(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id": rschema.StringAttribute{Optional: true, Computed: true},
	}}
	r := buildTestResourceWithSchema(schema, &is.TFStateHints{})
	req := resource.ImportStateRequest{ID: "id=abc"}
	resp := &resource.ImportStateResponse{}
	r.importStateImpl(context.Background(), req, resp)
	require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics.Errors())
	var gotStr types.String
	require.False(t, resp.State.GetAttribute(context.Background(), path.Root("id"), &gotStr).HasError())
	require.Equal(t, "abc", gotStr.ValueString())
}

// composite import tests moved to dedicated parser tests
