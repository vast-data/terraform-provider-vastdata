// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/schema_generation"
)

// TERF-312: import must fill optional network-settings fields from the API.
// They are optional in config but must be Optional+Computed so
// FillFromRecordForImport populates state (OpenAPI PATCH response is only async_task).
func TestVmsNetworkSettings_ImportFillsComputedFields(t *testing.T) {
	mgr := (&VmsNetworkSettings{}).NewResourceManager(nil, nil).(*VmsNetworkSettings)
	hints := mgr.TfState().Hints
	require.NotNil(t, hints)

	schema, err := schema_generation.GetResourceSchema(context.Background(), hints)
	require.NoError(t, err)

	for _, name := range vmsNetworkSettingsPatchFields {
		attrSchema, ok := schema.Attributes[name]
		require.True(t, ok, "missing schema attribute %q", name)
		assert.True(t, attrSchema.IsOptional(), "%s should remain optional for partial configs", name)
		assert.True(t, attrSchema.IsComputed(), "%s must be computed so import fills state (TERF-312)", name)
	}

	state := is.NewTFStateMust(nil, *schema, hints)
	for name, typ := range state.TypeMap {
		nullVal, ok, err := is.BuildAttrValueFromAny(typ, nil)
		require.True(t, ok)
		require.NoError(t, err)
		state.Raw[name] = nullVal
	}
	state.SetOrAdd("vms_id", int64(1))

	record := Record{
		"vms_id":  int64(1),
		"ntp":     []any{"10.130.0.17"},
		"eth_mtu": int64(9000),
		"ib_mtu":  int64(65520),
		"dns":     []any{"10.130.3.248"},
	}
	require.NoError(t, state.FillFromRecordForImport(record))

	assert.Equal(t, int64(1), state.Int64("vms_id"))
	assert.False(t, state.IsNull("ntp"), "ntp should be filled on import")
	assert.False(t, state.IsNull("eth_mtu"), "eth_mtu should be filled on import")
	assert.False(t, state.IsNull("ib_mtu"), "ib_mtu should be filled on import")
	assert.Equal(t, int64(9000), state.Int64("eth_mtu"))
	assert.Equal(t, int64(65520), state.Int64("ib_mtu"))
}

func TestVmsNetworkSettingsOverlayMatches(t *testing.T) {
	current := Record{
		"ntp":     []any{"10.140.3.248", "10.140.0.12"},
		"eth_mtu": float64(9000),
		"ib_mtu":  int64(65520),
		"dns":     []any{"10.140.3.248"},
	}
	assert.True(t, vmsNetworkSettingsOverlayMatches(current, params{
		"ntp":     []any{"10.140.0.12", "10.140.3.248"},
		"eth_mtu": int64(9000),
		"ib_mtu":  int64(65520),
	}))
	assert.False(t, vmsNetworkSettingsOverlayMatches(current, params{
		"ntp": []any{"10.140.0.12"},
	}))
}

func TestVmsNetworkSettings_TransformResponseRecord(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"vms_id": rschema.Int64Attribute{Required: true},
			"hosts":  rschema.SetAttribute{ElementType: types.DynamicType, Optional: true, Computed: true},
			"ntp":    rschema.SetAttribute{ElementType: types.StringType, Optional: true, Computed: true},
		},
	}
	m := &VmsNetworkSettings{tfstate: is.NewTFStateMust(map[string]attr.Value{
		"vms_id": types.Int64Value(1),
		"hosts":  types.SetNull(types.DynamicType),
		"ntp":    types.SetNull(types.StringType),
	}, schema, nil)}

	out := m.TransformResponseRecord(Record{
		"data": map[string]any{
			"ntp": []any{"10.130.0.17"},
			"boxes": []any{
				map[string]any{
					"box_name": "dbox-1",
					"hosts": []any{
						map[string]any{"id": float64(6), "hostname": "v2144dn1"},
					},
				},
			},
		},
	})
	assert.Equal(t, int64(1), out["vms_id"])
	assert.Equal(t, []any{"10.130.0.17"}, out["ntp"])
	_, hasBoxes := out["boxes"]
	assert.False(t, hasBoxes)
	hosts, ok := out["hosts"].([]any)
	require.True(t, ok)
	require.Len(t, hosts, 1)
}
