// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

// ---------- normalizeNumber ----------

func TestNormalizeNumber(t *testing.T) {
	assert.Equal(t, int64(10), normalizeNumber(float64(10)))
	assert.Equal(t, float64(10.5), normalizeNumber(float64(10.5)))
	assert.Equal(t, "abc", normalizeNumber("abc"))

	// Slices
	inputSlice := []any{float64(1), float64(2.5), "str"}
	expected := []any{int64(1), float64(2.5), "str"}
	assert.Equal(t, expected, normalizeNumber(inputSlice))

	// Maps
	inputMap := map[string]any{
		"x": float64(7),
		"y": []any{float64(3), "b"},
	}
	expectedMap := map[string]any{
		"x": int64(7),
		"y": []any{int64(3), "b"},
	}
	assert.Equal(t, expectedMap, normalizeNumber(inputMap))
}

// ---------- convertMapKeysRecursive ----------

func TestConvertMapKeysRecursive_UnderscoreToDash(t *testing.T) {
	input := map[string]any{
		"start_at": "val",
		"nested": map[string]any{
			"keep_local": "val2",
		},
	}
	expected := map[string]any{
		"start-at": "val",
		"nested": map[string]any{
			"keep-local": "val2",
		},
	}
	out := convertMapKeysRecursive(input, underscoreToDash)
	assert.Equal(t, expected, out)
}

func TestConvertMapKeysRecursive_DashToUnderscore(t *testing.T) {
	input := map[string]any{
		"start-at": "val",
		"nested": map[string]any{
			"keep-local": "val2",
		},
	}
	expected := map[string]any{
		"start_at": "val",
		"nested": map[string]any{
			"keep_local": "val2",
		},
	}
	out := convertMapKeysRecursive(input, dashToUnderscore)
	assert.Equal(t, expected, out)
}

// ---------- underscoreToDash / dashToUnderscore ----------

func TestKeyTransformHelpers(t *testing.T) {
	assert.Equal(t, "start-at", underscoreToDash("start_at"))
	assert.Equal(t, "keep_local", dashToUnderscore("keep-local"))
}

// ---------- validateOneOf / validateAllOf / validateNoneOf ----------

func TestValidateOneOf(t *testing.T) {
	tf := mustTFState(map[string]attr.Value{
		"a": types.StringValue("x"),
		"b": types.StringNull(), // defined, but null initially
		"c": types.StringNull(),
	})

	err := validateOneOf(tf, "a", "b", "c")
	assert.NoError(t, err)

	// Create new TFState with both a and b set
	tf = mustTFState(map[string]attr.Value{
		"a": types.StringValue("x"),
		"b": types.StringValue("y"),
		"c": types.StringNull(),
	})
	err = validateOneOf(tf, "a", "b", "c")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "only one of")
}

func TestValidateAllOf(t *testing.T) {
	tf := mustTFState(map[string]attr.Value{
		"a": types.StringValue("x"),
		"b": types.StringValue("y"),
		"c": types.StringNull(), // present, but intentionally null
	})
	err := validateAllOf(tf, "a", "b", "c")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be set")

	// Now make all present
	tf = mustTFState(map[string]attr.Value{
		"a": types.StringValue("x"),
		"b": types.StringValue("y"),
		"c": types.StringValue("z"),
	})
	err = validateAllOf(tf, "a", "b", "c")
	assert.NoError(t, err)
}

func TestValidateNoneOf(t *testing.T) {
	tf := mustTFState(map[string]attr.Value{
		"a": types.StringValue("x"),
		"b": types.StringValue("y"),
	})

	err := validateNoneOf(tf, "a", "b")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `none of ["a" "b"] should be set`)
	assert.Contains(t, err.Error(), `[a b]`)

	tf = mustTFState(map[string]attr.Value{
		"a": types.StringValue("x"),
		"b": types.StringNull(), // explicitly null
	})
	err = validateNoneOf(tf, "a", "b")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `none of ["a" "b"] should be set`)
	assert.Contains(t, err.Error(), `[a]`)

	tf = mustTFState(map[string]attr.Value{
		"a": types.StringNull(),
		"b": types.StringNull(),
	})
	err = validateNoneOf(tf, "a", "b")
	assert.NoError(t, err)
}

// ---------- TFState: Getters ----------

func TestTFState_String(t *testing.T) {
	tf := mustTFState(map[string]attr.Value{
		"foo": types.StringValue("bar"),
	})
	assert.Equal(t, "bar", tf.String("foo"))
}

func TestTFState_IsKnownAndNotNull(t *testing.T) {
	tf := mustTFState(map[string]attr.Value{
		"x": types.StringValue("test"),
	})
	assert.True(t, tf.IsKnownAndNotNull("x"))
}

// ---------- TFState: ToMap / ToSlice ----------

func TestTFState_ToMap(t *testing.T) {
	tf := mustTFState(map[string]attr.Value{
		"config": mustBuildAttr(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"enabled": types.BoolType,
			},
		}, map[string]any{"enabled": true}),
	})
	result := tf.ToMap("config")
	assert.Equal(t, map[string]any{"enabled": true}, result)
}

func TestTFState_ToSlice(t *testing.T) {
	tf := mustTFState(map[string]attr.Value{
		"items": mustBuildAttr(types.ListType{ElemType: types.StringType}, []any{"a", "b"}),
	})
	assert.Equal(t, []any{"a", "b"}, tf.ToSlice("items"))
}

// ---------- Helpers ----------

func mustTFState(raw map[string]attr.Value) *internalstate.TFState {
	// Provide a minimal dummy schema to enable TFState
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"a":     rschema.StringAttribute{Optional: true},
			"b":     rschema.StringAttribute{Optional: true},
			"c":     rschema.StringAttribute{Optional: true},
			"foo":   rschema.StringAttribute{Optional: true},
			"bar":   rschema.StringAttribute{Optional: true},
			"items": rschema.ListAttribute{ElementType: types.StringType, Optional: true},
			"config": rschema.SingleNestedAttribute{
				Attributes: map[string]rschema.Attribute{
					"enabled": rschema.BoolAttribute{Optional: true},
				},
			},
		},
	}
	return internalstate.NewTFStateMust(raw, schema, nil)
}

func mustBuildAttr(t attr.Type, val any) attr.Value {
	v, err := internalstate.BuildAttrValueFromAny(t, val)
	if err != nil {
		panic(fmt.Sprintf("build attr failed: %v", err))
	}
	return v
}
