// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

func TestQuotaSchema_DropsIsGroupAndRegistersUpgrade(t *testing.T) {
	t.Parallel()
	r := &Resource{managerName: "Quota", newManager: (&Quota{}).NewResourceManager}
	mgr, err := r.ManagerWithSchemaOnly(context.Background())
	require.NoError(t, err)
	sch := mgr.TfState().SchemaDefinition().(rschema.Schema)
	assert.Equal(t, quotaSchemaVersion, sch.Version)

	uq, ok := sch.Attributes["user_quotas"].(rschema.SetNestedAttribute)
	require.True(t, ok)
	_, has := uq.NestedObject.Attributes["is_group"]
	assert.False(t, has, "is_group must not remain in schema")

	upgraders := r.UpgradeState(context.Background())
	require.Contains(t, upgraders, int64(0))
}

func TestQuota_TransformRequestBody_RemovesIsGroup(t *testing.T) {
	t.Parallel()
	m := &Quota{}
	body := params{
		"name": "q",
		"user_quotas": []any{
			map[string]any{"name": "u1", "is_group": false, "soft_limit": int64(1)},
		},
		"group_quotas": []any{
			map[string]any{"name": "g1", "is_group": true, "hard_limit": int64(2)},
		},
	}
	out := m.TransformRequestBody(body)
	assert.NotContains(t, out["user_quotas"].([]any)[0].(map[string]any), "is_group")
	assert.NotContains(t, out["group_quotas"].([]any)[0].(map[string]any), "is_group")
	assert.Equal(t, "u1", out["user_quotas"].([]any)[0].(map[string]any)["name"])
	assert.Equal(t, "q", out["name"])
}

func TestQuota_UpgradeResourceSchema_DropsIsGroupFromPriorState(t *testing.T) {
	t.Parallel()
	r := &Resource{managerName: "Quota", newManager: (&Quota{}).NewResourceManager}
	mgr, err := r.ManagerWithSchemaOnly(context.Background())
	require.NoError(t, err)
	sch := mgr.TfState().SchemaDefinition().(rschema.Schema)

	// Prior state still has is_group; new schema does not. Upgrade must succeed
	// and leave no is_group on nested rules.
	priorJSON := []byte(`{
		"user_quotas": [{
			"name": "u1",
			"identifier": "u1",
			"identifier_type": "username",
			"is_group": false,
			"soft_limit": 1
		}]
	}`)

	upgraders := r.UpgradeState(context.Background())
	require.Contains(t, upgraders, int64(0))

	req := resource.UpgradeStateRequest{
		RawState: &tfprotov6.RawState{JSON: priorJSON},
	}
	resp := &resource.UpgradeStateResponse{State: tfsdk.State{Schema: sch}}
	upgraders[0].StateUpgrader(context.Background(), req, resp)
	require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics.Errors())

	var userQuotas types.Set
	require.False(t, resp.State.GetAttribute(context.Background(), path.Root("user_quotas"), &userQuotas).HasError())
	require.False(t, userQuotas.IsNull())
	elems := userQuotas.Elements()
	require.Len(t, elems, 1)
	obj := elems[0].(types.Object)
	_, has := obj.Attributes()["is_group"]
	assert.False(t, has, "is_group must be gone after UpgradeResourceSchema")
}

func TestResourceLifecycle_UpgradeResourceSchema_CallbacksAfterHook(t *testing.T) {
	t.Parallel()
	r := &Resource{managerName: "Quota", newManager: (&Quota{}).NewResourceManager}
	mgr, err := r.ManagerWithSchemaOnly(context.Background())
	require.NoError(t, err)
	sch := mgr.TfState().SchemaDefinition().(rschema.Schema)

	zeroRaw := make(map[string]attr.Value, len(sch.Attributes))
	for k, a := range sch.Attributes {
		zeroRaw[k], _, _ = is.BuildAttrValueFromAny(a.GetType(), nil)
	}
	ts := is.NewTFStateMust(zeroRaw, sch, mgr.TfState().StateHints())

	var order []string
	lc := NewResourceLifecycle(nil, "Quota")
	lc.AddUpgradeResourceSchemaCallback(func(_ context.Context, from int64, state *is.TFState) error {
		assert.Equal(t, int64(0), from)
		require.NotNil(t, state)
		order = append(order, "callback")
		return nil
	})

	// Wrap manager so we can observe the own hook running before callbacks.
	wrapped := upgradeOrderManager{ResourceManager: mgr, order: &order}
	require.NoError(t, lc.UpgradeResourceSchema(context.Background(), 0, wrapped, ts))
	assert.Equal(t, []string{"hook", "callback"}, order)
}

// upgradeOrderManager proxies Quota and records UpgradeResourceSchema calls.
type upgradeOrderManager struct {
	ResourceManager
	order *[]string
}

func (m upgradeOrderManager) SchemaVersion() int64 { return quotaSchemaVersion }

func (m upgradeOrderManager) UpgradeResourceSchema(ctx context.Context, from int64, ts *is.TFState) error {
	*m.order = append(*m.order, "hook")
	if u, ok := m.ResourceManager.(UpgradeResourceSchema); ok {
		return u.UpgradeResourceSchema(ctx, from, ts)
	}
	return nil
}
