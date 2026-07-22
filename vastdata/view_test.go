// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	t.Run("null cors", func(t *testing.T) {
		view := &View{
			tfstate: is.NewTFStateMust(
				map[string]attr.Value{
					"s3cors_configuration": types.ObjectNull(map[string]attr.Type{
						"cors_rules": types.ListType{ElemType: corsRuleObjectType},
					}),
				},
				nil,
				nil,
			),
		}
		assert.False(t, view.hasCorsConfigured())
	})

	t.Run("with cors rules", func(t *testing.T) {
		view := testViewWithCorsRules(t, testCorsRule(t, nil))
		assert.True(t, view.hasCorsConfigured())
	})

	t.Run("missing key", func(t *testing.T) {
		view := &View{
			tfstate: is.NewTFStateMust(map[string]attr.Value{}, nil, nil),
		}
		assert.False(t, view.hasCorsConfigured())
	})
}
