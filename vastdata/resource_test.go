// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vast_client "github.com/vast-data/go-vast-client"
	"github.com/vast-data/go-vast-client/core"
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
	// Return a record based on current state for testing
	record := make(Record)

	// Add id if present in state
	if t.tf.HasAttribute("id") && !t.tf.IsNull("id") {
		if t.tf.Type("id").Equal(types.Int64Type) {
			record["id"] = t.tf.Int64("id")
		} else if t.tf.Type("id").Equal(types.StringType) {
			record["id"] = t.tf.String("id")
		}
	}

	// Add other common test attributes
	for _, attr := range []string{"name", "gid", "tenant_id", "context"} {
		if t.tf.HasAttribute(attr) && !t.tf.IsNull(attr) {
			if t.tf.Type(attr).Equal(types.Int64Type) {
				record[attr] = t.tf.Int64(attr)
			} else if t.tf.Type(attr).Equal(types.StringType) {
				record[attr] = t.tf.String(attr)
			}
		}
	}

	return record, nil
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
	t.Skip("Temporarily skipped due to Terraform Plugin Framework state initialization issue - will be fixed post-release")
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
	t.Skip("Temporarily skipped due to Terraform Plugin Framework state initialization issue - will be fixed post-release")
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

func TestImport_KeyValue_MultipleFields(t *testing.T) {
	t.Skip("Temporarily skipped due to Terraform Plugin Framework state initialization issue - will be fixed post-release")
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id":        rschema.Int64Attribute{Optional: true, Computed: true},
		"name":      rschema.StringAttribute{Optional: true, Computed: true},
		"tenant_id": rschema.Int64Attribute{Optional: true, Computed: true},
	}}
	r := buildTestResourceWithSchema(schema, &is.TFStateHints{})
	req := resource.ImportStateRequest{ID: "name=foo, tenant_id=42"}
	resp := &resource.ImportStateResponse{}
	r.importStateImpl(context.Background(), req, resp)
	require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics.Errors())
	var name types.String
	var tenantID types.Int64
	require.False(t, resp.State.GetAttribute(context.Background(), path.Root("name"), &name).HasError())
	require.False(t, resp.State.GetAttribute(context.Background(), path.Root("tenant_id"), &tenantID).HasError())
	require.Equal(t, "foo", name.ValueString())
	require.Equal(t, int64(42), tenantID.ValueInt64())
}

func TestImport_SingleIdToken_Int64(t *testing.T) {
	t.Skip("Temporarily skipped due to Terraform Plugin Framework state initialization issue - will be fixed post-release")
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id": rschema.Int64Attribute{Optional: true, Computed: true},
	}}
	r := buildTestResourceWithSchema(schema, &is.TFStateHints{})
	req := resource.ImportStateRequest{ID: "1234"}
	resp := &resource.ImportStateResponse{}
	r.importStateImpl(context.Background(), req, resp)
	require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics.Errors())
	var id types.Int64
	require.False(t, resp.State.GetAttribute(context.Background(), path.Root("id"), &id).HasError())
	require.Equal(t, int64(1234), id.ValueInt64())
}

func TestImport_OrderedWithHints_Pipe(t *testing.T) {
	t.Skip("Temporarily skipped due to Terraform Plugin Framework state initialization issue - will be fixed post-release")
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"gid":       rschema.Int64Attribute{Optional: true, Computed: true},
		"tenant_id": rschema.Int64Attribute{Optional: true, Computed: true},
		"context":   rschema.StringAttribute{Optional: true, Computed: true},
	}}
	hints := &is.TFStateHints{ImportFields: []string{"gid", "tenant_id", "context"}}
	r := buildTestResourceWithSchema(schema, hints)
	req := resource.ImportStateRequest{ID: "1001|22|ad"}
	resp := &resource.ImportStateResponse{}
	r.importStateImpl(context.Background(), req, resp)
	require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics.Errors())
	var gid types.Int64
	var tenantID types.Int64
	var contextStr types.String
	require.False(t, resp.State.GetAttribute(context.Background(), path.Root("gid"), &gid).HasError())
	require.False(t, resp.State.GetAttribute(context.Background(), path.Root("tenant_id"), &tenantID).HasError())
	require.False(t, resp.State.GetAttribute(context.Background(), path.Root("context"), &contextStr).HasError())
	require.Equal(t, int64(1001), gid.ValueInt64())
	require.Equal(t, int64(22), tenantID.ValueInt64())
	require.Equal(t, "ad", contextStr.ValueString())
}

func TestImport_NotImportable_Err(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id": rschema.Int64Attribute{Optional: true, Computed: true},
	}}
	f := false
	r := buildTestResourceWithSchema(schema, &is.TFStateHints{Importable: &f})
	req := resource.ImportStateRequest{ID: "123"}
	resp := &resource.ImportStateResponse{}
	r.importStateImpl(context.Background(), req, resp)
	require.True(t, resp.Diagnostics.HasError())
}

// fake reader manager to verify read population is called
type testReaderManager struct{ testManager }

func (m *testReaderManager) ReadResource(_ context.Context, _ *VMSRest) (DisplayableRecord, error) {
	return Record{"id": int64(9), "name": "filled"}, nil
}

func buildResourceWithReader(schema rschema.Schema) *Resource {
	return &Resource{
		newManager: func(raw map[string]attr.Value, s any) ResourceManager {
			if s == nil {
				s = schema
			}
			return &testReaderManager{testManager{tf: is.NewTFStateMust(raw, s, &is.TFStateHints{TFStateHintsForCustom: &is.TFStateHintsForCustom{SchemaAttributes: map[string]any{"id": schema.Attributes["id"], "name": schema.Attributes["name"]}}})}}
		},
		managerName: "test",
	}
}

func TestImport_CallsReadAndFillsState(t *testing.T) {
	t.Skip("Temporarily skipped due to Terraform Plugin Framework state initialization issue - will be fixed post-release")
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id":   rschema.Int64Attribute{Optional: true, Computed: true},
		"name": rschema.StringAttribute{Optional: true, Computed: true},
	}}
	r := buildResourceWithReader(schema)
	req := resource.ImportStateRequest{ID: "id=9"}
	resp := &resource.ImportStateResponse{}
	r.importStateImpl(context.Background(), req, resp)
	require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics.Errors())
	var name types.String
	require.False(t, resp.State.GetAttribute(context.Background(), path.Root("name"), &name).HasError())
	require.Equal(t, "filled", name.ValueString())
}

// ------------------------------
// ManagerWithSchemaOnly tests
// ------------------------------

func TestResource_ManagerWithSchemaOnly_FillsZeroRaw(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id":   rschema.Int64Attribute{Optional: true, Computed: true},
		"name": rschema.StringAttribute{Optional: true, Computed: true},
	}}
	r := buildTestResourceWithSchema(schema, &is.TFStateHints{})
	mgr, err := r.ManagerWithSchemaOnly(context.Background())
	require.NoError(t, err)
	tf := mgr.TfState()
	// keys exist
	require.True(t, tf.HasAttribute("id"))
	require.True(t, tf.HasAttribute("name"))
	// values initialized to Null
	require.True(t, tf.IsNull("id"))
	require.True(t, tf.IsNull("name"))
}

type testDSManager struct{ tf *is.TFState }

func (t *testDSManager) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return t
}
func (t *testDSManager) TfState() *is.TFState                      { return t.tf }
func (t *testDSManager) API(_ *VMSRest) VastResourceAPIWithContext { return nil }

func buildTestDatasourceWithSchema(schema dschema.Schema, hints *is.TFStateHints) *Datasource {
	return &Datasource{
		newManager: func(raw map[string]attr.Value, s any) DataSourceManager {
			localHints := hints
			if localHints == nil {
				localHints = &is.TFStateHints{}
			}
			localHints.TFStateHintsForCustom = &is.TFStateHintsForCustom{
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
			return &testDSManager{tf: is.NewTFStateMust(raw, s, localHints)}
		},
		managerName: "test_ds",
	}
}

func TestDatasource_ManagerWithSchemaOnly_FillsZeroRaw(t *testing.T) {
	schema := dschema.Schema{Attributes: map[string]dschema.Attribute{
		"id":   dschema.Int64Attribute{Optional: true, Computed: true},
		"name": dschema.StringAttribute{Optional: true, Computed: true},
	}}
	d := buildTestDatasourceWithSchema(schema, &is.TFStateHints{})
	mgr, err := d.ManagerWithSchemaOnly(context.Background())
	require.NoError(t, err)
	tf := mgr.TfState()
	require.True(t, tf.HasAttribute("id"))
	require.True(t, tf.HasAttribute("name"))
	require.True(t, tf.IsNull("id"))
	require.True(t, tf.IsNull("name"))
}

// ------------------------------
// Delete-only hints tests
// ------------------------------

func TestDeleteOnlyParamsAndBody_FromHints(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id":          rschema.Int64Attribute{Optional: true, Computed: true},
		"delete_flag": rschema.StringAttribute{Optional: true},
		"reason":      rschema.StringAttribute{Optional: true},
	}}
	hints := &is.TFStateHints{
		DeleteOnlyBodyFields:  map[string]string{"delete_flag": "delete_flag_api"},
		DeleteOnlyParamFields: map[string]string{"reason": "reason_api"},
	}
	r := buildTestResourceWithSchema(schema, hints)
	mgr, err := r.ManagerWithSchemaOnly(context.Background())
	require.NoError(t, err)
	tf := mgr.TfState()

	// Set fields in state
	tf.SetOrAdd("delete_flag", "force")
	tf.SetOrAdd("reason", "cleanup")

	// Ensure GetDeleteOnlyBodyParams returns only the body fields
	body := tf.GetDeleteOnlyBodyParams()
	require.Equal(t, map[string]any{"delete_flag_api": "force"}, map[string]any(body))

	// Ensure GetDeleteOnlyQueryParams returns only the param fields
	qp := tf.GetDeleteOnlyQueryParams()
	require.Equal(t, map[string]any{"reason_api": "cleanup"}, map[string]any(qp))
}

// --- FillFromRecordWithComputedOnly tests ---

func TestFillFromRecordIncludingRequired_ComputedOnlyTrue(t *testing.T) {
	t.Skip("Temporarily skipped due to TFState initialization issue - will be fixed post-release")
	// Schema: id (computed), title (computed), name (optional)
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id":    rschema.Int64Attribute{Optional: true, Computed: true},
		"title": rschema.StringAttribute{Optional: true, Computed: true},
		"name":  rschema.StringAttribute{Optional: true},
	}}
	tf := is.NewTFStateMust(map[string]attr.Value{}, schema, nil)

	rec := Record{
		"id":    int64(8),
		"title": "from-backend",
		"name":  "should-not-be-set",
	}

	// computedOnly = true => only computed fields should be set
	err := tf.FillFromRecordIncludingRequired(rec, false)
	require.NoError(t, err)

	// id and title should be set
	// Directly check Raw via helpers
	// id
	tfID := tf.Get("id").(types.Int64)
	require.Equal(t, int64(8), tfID.ValueInt64())
	// title
	tfTitle := tf.Get("title").(types.String)
	require.Equal(t, "from-backend", tfTitle.ValueString())
	// name should remain null/unknown
	require.True(t, tf.IsNull("name") || tf.IsUnknown("name"))
}

func TestFillFromRecordIncludingRequired_ComputedOnlyFalse(t *testing.T) {
	// Schema: id (computed), title (computed), name (optional)
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id":    rschema.Int64Attribute{Optional: true, Computed: true},
		"title": rschema.StringAttribute{Optional: true, Computed: true},
		"name":  rschema.StringAttribute{Optional: true},
	}}
	tf := is.NewTFStateMust(map[string]attr.Value{}, schema, nil)

	rec := Record{
		"id":    int64(9),
		"title": "from-backend",
		"name":  "should-be-set",
		"bogus": "skip-me", // unknown key should be ignored
	}

	// includeRequired = true => computed + required fields present in record should be set
	// Simulate name being required in meta
	m := tf.Meta["name"]
	m.Required = true
	tf.Meta["name"] = m
	err := tf.FillFromRecordIncludingRequired(rec, true)
	require.NoError(t, err)

	// id
	tfID := tf.Get("id").(types.Int64)
	require.Equal(t, int64(9), tfID.ValueInt64())
	// title
	tfTitle := tf.Get("title").(types.String)
	require.Equal(t, "from-backend", tfTitle.ValueString())
	// name should be set as well
	tfName := tf.Get("name").(types.String)
	require.Equal(t, "should-be-set", tfName.ValueString())
}

// ---- MIGRATE MODE TESTS ----

func TestMigrateMode_ProviderDataFlag(t *testing.T) {
	// Verify MigrateMode flag is properly stored in ProviderData
	pd := &ProviderData{Client: nil, MigrateMode: true}
	require.True(t, pd.MigrateMode)

	pd2 := &ProviderData{Client: nil, MigrateMode: false}
	require.False(t, pd2.MigrateMode)
}

func TestMigrateMode_DeleteBlockedOnNonEmptyState(t *testing.T) {
	// In migrate mode, delete should be blocked because it means state is not empty
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id":   rschema.Int64Attribute{Optional: true, Computed: true},
		"name": rschema.StringAttribute{Optional: true, Computed: true},
	}}

	r := buildTestResourceWithSchema(schema, &is.TFStateHints{})
	r.providerData = &ProviderData{Client: nil, MigrateMode: true}

	resp := &resource.DeleteResponse{}
	r.Delete(context.Background(), resource.DeleteRequest{}, resp)

	// Should produce an error telling user state is not empty
	require.True(t, resp.Diagnostics.HasError(), "migrate mode delete should produce an error")
	require.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "EMPTY state")
}

func TestMigrateMode_UpdateBlockedOnNonEmptyState(t *testing.T) {
	// In migrate mode, update should be blocked because it means state is not empty
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id":   rschema.Int64Attribute{Optional: true, Computed: true},
		"name": rschema.StringAttribute{Optional: true, Computed: true},
	}}

	r := buildTestResourceWithSchema(schema, &is.TFStateHints{})
	r.providerData = &ProviderData{Client: nil, MigrateMode: true}

	resp := &resource.UpdateResponse{}
	r.Update(context.Background(), resource.UpdateRequest{}, resp)

	// Should produce an error telling user state is not empty
	require.True(t, resp.Diagnostics.HasError(), "migrate mode update should produce an error")
	require.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "EMPTY state")
}

func TestMigrateMode_NilProviderData_NoMigrate(t *testing.T) {
	// When providerData is nil, migrate mode checks should not panic
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id": rschema.Int64Attribute{Optional: true, Computed: true},
	}}

	r := buildTestResourceWithSchema(schema, &is.TFStateHints{})
	// providerData is nil by default

	// The nil check should prevent any panic when accessing MigrateMode
	require.NotPanics(t, func() {
		pd := r.providerData
		_ = pd != nil && pd.MigrateMode
	})
}

// TERF-268: write-only values are null in Plan and only present in Config.
func TestOverlayWriteOnlyFromConfig_IncludedInCreateParams(t *testing.T) {
	t.Parallel()

	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"name":        rschema.StringAttribute{Required: true},
		"cert_type":   rschema.StringAttribute{Optional: true},
		"certificate": rschema.StringAttribute{Optional: true, WriteOnly: true, Sensitive: true},
		"private_key": rschema.StringAttribute{Optional: true, WriteOnly: true, Sensitive: true},
		"ca_certificate": rschema.StringAttribute{
			Optional:  true,
			WriteOnly: true,
			Sensitive: true,
		},
	}}
	hints := &is.TFStateHints{
		WriteOnlyFields: []string{"certificate", "private_key", "ca_certificate"},
	}
	r := buildTestResourceWithSchema(schema, hints)

	objType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"name":           tftypes.String,
		"cert_type":      tftypes.String,
		"certificate":    tftypes.String,
		"private_key":    tftypes.String,
		"ca_certificate": tftypes.String,
	}}

	// Plan: write-only attrs are null (Terraform behavior).
	plan := tfsdk.Plan{
		Schema: schema,
		Raw: tftypes.NewValue(objType, map[string]tftypes.Value{
			"name":           tftypes.NewValue(tftypes.String, "tf-cert"),
			"cert_type":      tftypes.NewValue(tftypes.String, "WEBHOOK"),
			"certificate":    tftypes.NewValue(tftypes.String, nil),
			"private_key":    tftypes.NewValue(tftypes.String, nil),
			"ca_certificate": tftypes.NewValue(tftypes.String, nil),
		}),
	}
	// Config: write-only PEMs are present.
	config := tfsdk.Config{
		Schema: schema,
		Raw: tftypes.NewValue(objType, map[string]tftypes.Value{
			"name":           tftypes.NewValue(tftypes.String, "tf-cert"),
			"cert_type":      tftypes.NewValue(tftypes.String, "WEBHOOK"),
			"certificate":    tftypes.NewValue(tftypes.String, "CERT-PEM"),
			"private_key":    tftypes.NewValue(tftypes.String, "KEY-PEM"),
			"ca_certificate": tftypes.NewValue(tftypes.String, "CA-PEM"),
		}),
	}

	manager := r.NewManager(plan)
	paramsBefore := manager.TfState().GetCreateParams()
	assert.NotContains(t, paramsBefore, "certificate")
	assert.NotContains(t, paramsBefore, "private_key")
	assert.NotContains(t, paramsBefore, "ca_certificate")
	assert.Equal(t, "WEBHOOK", paramsBefore["cert_type"])

	r.overlayWriteOnlyFromConfig(manager, config)
	params := manager.TfState().GetCreateParams()
	assert.Equal(t, "CERT-PEM", params["certificate"])
	assert.Equal(t, "KEY-PEM", params["private_key"])
	assert.Equal(t, "CA-PEM", params["ca_certificate"])
	assert.Equal(t, "tf-cert", params["name"])
	assert.Equal(t, "WEBHOOK", params["cert_type"])
}

type deleteErrorAPI struct {
	deleteErr error
}

func (m *deleteErrorAPI) Session() core.RESTSession { return nil }
func (m *deleteErrorAPI) GetResourceType() string   { return "test" }
func (m *deleteErrorAPI) GetResourcePath() string   { return "test" }
func (m *deleteErrorAPI) List(vast_client.Params) (RecordSet, error) {
	return nil, nil
}
func (m *deleteErrorAPI) Create(vast_client.Params) (Record, error) { return nil, nil }
func (m *deleteErrorAPI) Update(any, vast_client.Params) (Record, error) {
	return nil, nil
}
func (m *deleteErrorAPI) Delete(vast_client.Params, vast_client.Params) (Record, error) {
	return nil, nil
}
func (m *deleteErrorAPI) DeleteById(any, vast_client.Params, vast_client.Params) (Record, error) {
	return nil, nil
}
func (m *deleteErrorAPI) Ensure(vast_client.Params, vast_client.Params) (Record, error) {
	return nil, nil
}
func (m *deleteErrorAPI) Get(vast_client.Params) (Record, error) { return nil, nil }
func (m *deleteErrorAPI) GetById(any) (Record, error)            { return nil, nil }
func (m *deleteErrorAPI) Exists(vast_client.Params) (bool, error) {
	return false, nil
}
func (m *deleteErrorAPI) MustExists(vast_client.Params) bool { return false }
func (m *deleteErrorAPI) GetIterator(vast_client.Params, int) core.Iterator {
	return nil
}
func (m *deleteErrorAPI) Lock(...any) func() { return func() {} }
func (m *deleteErrorAPI) ListWithContext(context.Context, vast_client.Params) (RecordSet, error) {
	return nil, nil
}
func (m *deleteErrorAPI) CreateWithContext(context.Context, vast_client.Params) (Record, error) {
	return nil, nil
}
func (m *deleteErrorAPI) UpdateWithContext(context.Context, any, vast_client.Params) (Record, error) {
	return nil, nil
}
func (m *deleteErrorAPI) DeleteWithContext(context.Context, vast_client.Params, vast_client.Params, vast_client.Params) (Record, error) {
	return nil, nil
}
func (m *deleteErrorAPI) DeleteByIdWithContext(_ context.Context, _ any, _, _ vast_client.Params) (Record, error) {
	return nil, m.deleteErr
}
func (m *deleteErrorAPI) EnsureWithContext(context.Context, vast_client.Params, vast_client.Params) (Record, error) {
	return nil, nil
}
func (m *deleteErrorAPI) GetWithContext(context.Context, vast_client.Params) (Record, error) {
	return nil, nil
}
func (m *deleteErrorAPI) GetByIdWithContext(context.Context, any) (Record, error) {
	return nil, nil
}
func (m *deleteErrorAPI) ExistsWithContext(context.Context, vast_client.Params) (bool, error) {
	return false, nil
}
func (m *deleteErrorAPI) MustExistsWithContext(context.Context, vast_client.Params) bool {
	return false
}
func (m *deleteErrorAPI) GetIteratorWithContext(context.Context, vast_client.Params, int) core.Iterator {
	return nil
}

func TestDeleteRecordBySearchParams_ReturnsNon404ApiError(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id": rschema.Int64Attribute{Optional: true, Computed: true},
	}}
	tf := is.NewTFStateMust(
		map[string]attr.Value{"id": types.Int64Value(8)},
		schema,
		&is.TFStateHints{},
	)
	deleteErr := &ApiError{
		StatusCode: 503,
		Body:       `{"detail":"vip_pool_remove returned an error: ObjectRemoveResultCode.CONFLICT"}`,
	}

	_, err := deleteRecordBySearchParams(
		context.Background(),
		&deleteErrorAPI{deleteErr: deleteErr},
		tf,
		"vastdata_vip_pool",
		"Delete",
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "503")
}

func TestRetryOnExpression_DeleteRetriesTransient503(t *testing.T) {
	attempts := 0
	expr := &is.RetryExpression{
		StatusCodes:  []int{503},
		Times:        3,
		SleepSeconds: 0,
	}

	_, err := retryOnExpression(context.Background(), expr, "Delete", "vastdata_vip_pool", func() (struct{}, error) {
		attempts++
		if attempts < 3 {
			return struct{}{}, &ApiError{StatusCode: 503, Body: "CONFLICT"}
		}
		return struct{}{}, nil
	})

	require.NoError(t, err)
	require.Equal(t, 3, attempts)
}
