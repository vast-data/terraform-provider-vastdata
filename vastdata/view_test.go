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

var corsRuleObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"allowed_methods": types.ListType{ElemType: types.StringType},
		"allowed_origins": types.ListType{ElemType: types.StringType},
		"allowed_headers": types.ListType{ElemType: types.StringType},
		"expose_headers":  types.ListType{ElemType: types.StringType},
		"max_age_seconds": types.Int64Type,
	},
}

func testStringList(t *testing.T, vals ...string) types.List {
	t.Helper()
	elems := make([]attr.Value, len(vals))
	for i, v := range vals {
		elems[i] = types.StringValue(v)
	}
	return types.ListValueMust(types.StringType, elems)
}

func testCorsRule(t *testing.T, overrides map[string]attr.Value) types.Object {
	t.Helper()
	attrs := map[string]attr.Value{
		"allowed_methods": testStringList(t, "GET"),
		"allowed_origins": testStringList(t, "https://example.com"),
		"allowed_headers": types.ListNull(types.StringType),
		"expose_headers":  types.ListNull(types.StringType),
		"max_age_seconds": types.Int64Null(),
	}
	for k, v := range overrides {
		attrs[k] = v
	}
	return types.ObjectValueMust(corsRuleObjectType.AttrTypes, attrs)
}

func testViewWithCorsRules(t *testing.T, rules ...types.Object) *View {
	t.Helper()
	ruleValues := make([]attr.Value, len(rules))
	for i, rule := range rules {
		ruleValues[i] = rule
	}
	corsRules := types.ListValueMust(corsRuleObjectType, ruleValues)
	corsConfig := types.ObjectValueMust(
		map[string]attr.Type{"cors_rules": types.ListType{ElemType: corsRuleObjectType}},
		map[string]attr.Value{"cors_rules": corsRules},
	)
	return &View{
		tfstate: is.NewTFStateMust(
			map[string]attr.Value{"s3cors_configuration": corsConfig},
			nil,
			nil,
		),
	}
}

func testCorsBodyRule(t *testing.T, body params, index int) map[string]any {
	t.Helper()
	rules, ok := body["cors_rules"].([]any)
	require.True(t, ok, "cors_rules should be a slice")
	require.Greater(t, len(rules), index)
	rule, ok := rules[index].(map[string]any)
	require.True(t, ok, "cors rule should be a map")
	return rule
}

func TestView_corsBody_maxAgeSecondsZero(t *testing.T) {
	view := testViewWithCorsRules(t, testCorsRule(t, map[string]attr.Value{
		"max_age_seconds": types.Int64Value(0),
	}))

	body, ok := view.corsBody()
	require.True(t, ok)
	require.NotNil(t, body)

	rule := testCorsBodyRule(t, body, 0)
	value, present := rule["max_age_seconds"]
	require.True(t, present, "max_age_seconds must be present for explicit 0")
	assert.Equal(t, int64(0), value)
}

func TestView_corsBody_maxAgeSecondsOmittedWhenNull(t *testing.T) {
	view := testViewWithCorsRules(t, testCorsRule(t, nil))

	body, ok := view.corsBody()
	require.True(t, ok)
	require.NotNil(t, body)

	rule := testCorsBodyRule(t, body, 0)
	_, present := rule["max_age_seconds"]
	assert.False(t, present, "max_age_seconds should be omitted when unset")
}

func TestView_corsBody_maxAgeSecondsNonZero(t *testing.T) {
	view := testViewWithCorsRules(t, testCorsRule(t, map[string]attr.Value{
		"max_age_seconds": types.Int64Value(3600),
	}))

	body, ok := view.corsBody()
	require.True(t, ok)

	rule := testCorsBodyRule(t, body, 0)
	assert.Equal(t, int64(3600), rule["max_age_seconds"])
}

func TestView_corsBody_maxAgeSecondsZero_roundTrip(t *testing.T) {
	view := testViewWithCorsRules(t, testCorsRule(t, map[string]attr.Value{
		"max_age_seconds": types.Int64Value(0),
	}))

	body, ok := view.corsBody()
	require.True(t, ok)

	// GetSubResources passes cors_rules through verbatim from the API response.
	apiResponse := map[string]any{"cors_rules": body["cors_rules"]}
	rules, ok := apiResponse["cors_rules"].([]any)
	require.True(t, ok)
	require.Len(t, rules, 1)

	rule, ok := rules[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, int64(0), rule["max_age_seconds"])
}

func viewCorsTestSchema() rschema.Schema {
	return rschema.Schema{Attributes: map[string]rschema.Attribute{
		"get_s3cors_configuration": rschema.BoolAttribute{Optional: true},
		"s3cors_configuration": rschema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]rschema.Attribute{
				"cors_rules": rschema.ListNestedAttribute{Optional: true},
			},
		},
	}}
}

func TestIsS3View(t *testing.T) {
	tests := []struct {
		name     string
		record   Record
		expected bool
	}{
		{name: "nil record", record: nil, expected: false},
		{name: "missing protocols", record: Record{"id": 1}, expected: false},
		{name: "nil protocols", record: Record{"protocols": nil}, expected: false},
		{name: "NFS only", record: Record{"protocols": []any{"NFS"}}, expected: false},
		{name: "NFS and NFS4", record: Record{"protocols": []any{"NFS", "NFS4"}}, expected: false},
		{name: "S3 only", record: Record{"protocols": []any{"S3"}}, expected: true},
		{name: "S3 and NFS", record: Record{"protocols": []any{"NFS", "S3"}}, expected: true},
		{name: "string slice S3", record: Record{"protocols": []string{"S3"}}, expected: true},
		{name: "string slice NFS", record: Record{"protocols": []string{"NFS"}}, expected: false},
		{name: "BLOCK only", record: Record{"protocols": []any{"BLOCK"}}, expected: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isS3View(tt.record))
		})
	}
}

func TestView_hasCorsConfigured(t *testing.T) {
	t.Parallel()
	schema := viewCorsTestSchema()

	t.Run("disabled", func(t *testing.T) {
		t.Parallel()
		view := &View{
			tfstate: is.NewTFStateMust(map[string]attr.Value{}, nil, nil),
		}
		assert.False(t, view.hasCorsConfigured())
	})

	t.Run("null cors", func(t *testing.T) {
		t.Parallel()
		view := &View{
			tfstate: is.NewTFStateMust(
				map[string]attr.Value{
					"s3cors_configuration": types.ObjectNull(map[string]attr.Type{
						"cors_rules": types.ListType{ElemType: corsRuleObjectType},
					}),
				},
				schema,
				nil,
			),
		}
		assert.False(t, view.hasCorsConfigured())
	})

	t.Run("with cors rules", func(t *testing.T) {
		t.Parallel()
		view := testViewWithCorsRules(t, testCorsRule(t, nil))
		view.tfstate = is.NewTFStateMust(view.tfstate.Raw, schema, nil)
		assert.True(t, view.hasCorsConfigured())
	})
}

func TestView_AfterUpdateResource_SkipsCorsDeleteWhenNeverConfigured(t *testing.T) {
	t.Parallel()

	schema := viewCorsTestSchema()
	nullCors := types.ObjectNull(map[string]attr.Type{
		"cors_rules": types.ListType{ElemType: corsRuleObjectType},
	})
	raw := map[string]attr.Value{
		"get_s3cors_configuration": types.BoolNull(),
		"s3cors_configuration":     nullCors,
	}

	// NFS-style prior/plan: no CORS configured. Must not call DELETE (rest=nil would panic).
	prior := &View{tfstate: is.NewTFStateMust(raw, schema, nil)}
	plan := &View{tfstate: is.NewTFStateMust(raw, schema, nil)}

	err := prior.AfterUpdateResource(t.Context(), plan, nil, Record{
		"id":        int64(41),
		"protocols": []any{"NFS"},
	})
	require.NoError(t, err)
}

func TestView_AfterUpdateResource_RejectsCorsOnNonS3View(t *testing.T) {
	t.Parallel()

	schema := viewCorsTestSchema()
	planView := testViewWithCorsRules(t, testCorsRule(t, nil))
	raw := planView.tfstate.Raw
	raw["get_s3cors_configuration"] = types.BoolNull()
	planView.tfstate = is.NewTFStateMust(raw, schema, nil)
	prior := &View{tfstate: is.NewTFStateMust(map[string]attr.Value{
		"get_s3cors_configuration": types.BoolNull(),
		"s3cors_configuration": types.ObjectNull(map[string]attr.Type{
			"cors_rules": types.ListType{ElemType: corsRuleObjectType},
		}),
	}, schema, nil)}

	err := prior.AfterUpdateResource(t.Context(), planView, nil, Record{
		"id":        int64(661),
		"protocols": []any{"NFS"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), `protocols to include "S3"`)
}
