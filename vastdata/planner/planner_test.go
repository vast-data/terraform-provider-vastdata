// Copyright (c) HashiCorp, Inc.

package planner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	vsd "github.com/vast-data/terraform-provider-vastdata/vastdata"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

func newState(t *testing.T, name string, vals map[string]any) is.ResourceState {
	t.Helper()
	m, err := vsd.NewResourceManagerByName(context.Background(), name)
	require.NoError(t, err)
	for k, v := range vals {
		require.True(t, m.TfState().HasAttribute(k), "%s has no attribute %q", name, k)
		m.TfState().SetOrAdd(k, v)
	}
	return m.TfState()
}

func plan(t *testing.T, name string, prior, config map[string]any) *Result {
	t.Helper()
	var priorSt is.ResourceState
	if prior != nil {
		priorSt = newState(t, name, prior)
	}
	res, err := Plan(context.Background(), name, priorSt, newState(t, name, config))
	require.NoError(t, err)
	return res
}

func merge(base map[string]any, over map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range base {
		out[k] = v
	}
	for k, v := range over {
		out[k] = v
	}
	return out
}

// Computed values the config leaves unset (server defaults) are carried from prior state,
// so an unchanged config plans no change.
func TestPlan_UnchangedConfigIsNoOp(t *testing.T) {
	config := map[string]any{"path": "/a", "policy_id": int64(1)}
	prior := merge(config, map[string]any{"id": int64(7), "tenant_id": int64(1), "name": "a", "protocols": []any{"NFS"}})

	res := plan(t, "view", prior, config)
	assert.Equal(t, NoOp, res.Action, "changed=%v unknown=%v", res.Changed, res.Unknown)
	assert.Empty(t, res.Changed)
}

func TestPlan_ConfigChangeIsUpdate(t *testing.T) {
	prior := map[string]any{"id": int64(7), "path": "/a", "policy_id": int64(1), "tenant_id": int64(1)}
	res := plan(t, "view", prior, map[string]any{"path": "/a", "policy_id": int64(2)})
	assert.Equal(t, Update, res.Action)
	assert.Contains(t, res.Changed, "policy_id")
}

// EditOnly fields never reach GetUpdateParams, but a real plan still schedules Update for
// them (ActiveDirectory join).
func TestPlan_EditOnlyChangeIsUpdate(t *testing.T) {
	base := map[string]any{"machine_account_name": "m", "domain_name": "d.lab", "admin_username": "u", "admin_passwd": "p"}
	prior := merge(base, map[string]any{"id": int64(1), "enabled": false})

	res := plan(t, "active_directory", prior, merge(base, map[string]any{"enabled": true}))
	assert.Equal(t, Update, res.Action)
	assert.Equal(t, []string{"enabled"}, res.Changed)
	assert.NotContains(t, res.GetUpdateParams(newState(t, "active_directory", prior)), "enabled",
		"the main PATCH still leaves EditOnly fields out")
}

// Clearing a write-only secret nulls its *_md5 companion in the plan (TERF-320), so Update
// runs and AfterUpdateResource deletes the CRL.
func TestPlan_ClearedWriteOnlySecretIsUpdate(t *testing.T) {
	const pem, crl = "pem-bytes", "crl-bytes"
	base := map[string]any{"ca_certificate_name": "c", "tenant_id": int64(1), "protocols": []any{"NFS"}}
	prior := merge(base, map[string]any{
		"id": int64(1), "ca_certificate_file_md5": vsd.ContentMD5Hex(pem),
		"revocation_file_md5": vsd.ContentMD5Hex(crl), "revocations_name": "r",
	})

	res := plan(t, "tls_certificate", prior, merge(base, map[string]any{"ca_certificate_file": pem}))
	assert.Equal(t, Update, res.Action)
	assert.Equal(t, []string{"revocation_file_md5"}, res.Changed)
	assert.True(t, res.IsNull("ca_certificate_file"), "the plan holds no write-only value")
}

func TestPlan_WriteOnlySecrets(t *testing.T) {
	const pem, crl = "pem-bytes", "crl-bytes"
	base := map[string]any{"ca_certificate_name": "c", "tenant_id": int64(1), "protocols": []any{"NFS"}}
	prior := merge(base, map[string]any{
		"id": int64(1), "ca_certificate_file_md5": vsd.ContentMD5Hex(pem),
		"revocation_file_md5": vsd.ContentMD5Hex(crl), "revocations_name": "r",
	})

	same := plan(t, "tls_certificate", prior, merge(base, map[string]any{"ca_certificate_file": pem, "revocation_file": crl}))
	assert.Equal(t, NoOp, same.Action, "changed=%v", same.Changed)

	rotated := plan(t, "tls_certificate", prior, merge(base, map[string]any{"ca_certificate_file": pem, "revocation_file": "new-crl"}))
	assert.Equal(t, Update, rotated.Action)
	assert.Equal(t, []string{"revocation_file_md5"}, rotated.Changed)
}

// Sub-resource attributes are excluded from the main PATCH but are part of the plan.
func TestPlan_SubResourceChangeIsUpdate(t *testing.T) {
	rule := func(origin string) map[string]any {
		return map[string]any{"allowed_methods": []any{"GET"}, "allowed_origins": []any{origin}}
	}
	base := map[string]any{"path": "/a", "policy_id": int64(1), "get_s3cors_configuration": true}
	priorRule := merge(rule("https://a"), map[string]any{"allowed_headers": []any{"*"}})
	prior := merge(base, map[string]any{
		"id": int64(7), "tenant_id": int64(1),
		"s3cors_configuration": map[string]any{"cors_rules": []any{priorRule}},
	})

	unchanged := plan(t, "view", prior, merge(base, map[string]any{
		"s3cors_configuration": map[string]any{"cors_rules": []any{rule("https://a")}},
	}))
	assert.Equal(t, NoOp, unchanged.Action,
		"a computed nested value the config leaves unset is carried from prior (changed=%v)", unchanged.Changed)

	changed := plan(t, "view", prior, merge(base, map[string]any{
		"s3cors_configuration": map[string]any{"cors_rules": []any{rule("https://b")}},
	}))
	assert.Equal(t, Update, changed.Action)
	assert.Equal(t, []string{"s3cors_configuration"}, changed.Changed)
}

func TestPlan_RequiresReplace(t *testing.T) {
	prior := map[string]any{"sid": "S-1", "s3_policy_id": int64(1)}
	res := plan(t, "s3_policy_attachment", prior, map[string]any{"sid": "S-2", "s3_policy_id": int64(1)})
	assert.Equal(t, Replace, res.Action)
	assert.Equal(t, []string{"sid"}, res.RequiresReplace)
}

func TestPlan_NoPriorIsCreate(t *testing.T) {
	res := plan(t, "view", nil, map[string]any{"path": "/a", "policy_id": int64(1)})
	assert.Equal(t, Create, res.Action)
}

func TestPlan_ConfigShapeIsChecked(t *testing.T) {
	ctx := context.Background()

	_, err := Plan(ctx, "view", nil, newState(t, "view", map[string]any{"policy_id": int64(1)}))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing required argument: AttributeName(\"path\")")

	_, err = Plan(ctx, "view", nil, newState(t, "view", map[string]any{"path": "/a", "policy_id": int64(1), "id": int64(5)}))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read-only attribute")
}
