// Copyright (c) HashiCorp, Inc.

package planner

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	vsd "github.com/vast-data/terraform-provider-vastdata/vastdata"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

// ---- nested sets ----------------------------------------------------------------------

// A set element the server enriched with computed values still matches the configured
// element it came from, in any order.
func TestPlan_SetNestedCorrelatesServerFilledElements(t *testing.T) {
	ace := func(grantee, perm string) map[string]any { return map[string]any{"grantee": grantee, "perm": perm} }
	served := func(grantee, perm, sid string, uid int64) map[string]any {
		return merge(ace(grantee, perm), map[string]any{"fqdn": "lab.local", "name": grantee, "sid_str": sid, "uid_or_gid": uid})
	}
	base := map[string]any{"path": "/a/", "policy_id": int64(1)}
	prior := merge(base, map[string]any{"id": int64(7), "tenant_id": int64(1), "share_acl": map[string]any{
		"enabled": true,
		"acl":     []any{served("users", "FULL", "S-1", 5), served("groups", "READ", "S-2", 6)},
	}})
	acl := func(entries ...any) map[string]any {
		return merge(base, map[string]any{"share_acl": map[string]any{"enabled": true, "acl": entries}})
	}

	res := plan(t, "view", prior, acl(ace("users", "FULL"), ace("groups", "READ")))
	assert.Equal(t, NoOp, res.Action, "changed=%v", res.Changed)

	res = plan(t, "view", prior, acl(ace("groups", "READ"), ace("users", "FULL")))
	assert.Equal(t, NoOp, res.Action, "set order does not matter (changed=%v)", res.Changed)

	res = plan(t, "view", prior, acl(ace("users", "READ"), ace("groups", "READ")))
	assert.Equal(t, Update, res.Action)
	assert.Equal(t, []string{"share_acl"}, res.Changed)

	res = plan(t, "view", prior, acl(ace("users", "FULL"), ace("groups", "READ"), ace("groups", "FULL")))
	assert.Equal(t, Update, res.Action, "a new element is a change")

	res = plan(t, "view", prior, acl(ace("users", "FULL")))
	assert.Equal(t, Update, res.Action, "a removed element is a change")
}

func TestPlan_SetNestedQuotaEntries(t *testing.T) {
	entry := func(limit int64) map[string]any {
		return map[string]any{"identifier": "alice", "identifier_type": "username", "hard_limit": limit}
	}
	base := map[string]any{"name": "q", "path": "/q/"}
	prior := merge(base, map[string]any{"id": int64(3), "user_quotas": []any{
		merge(entry(10), map[string]any{"used_capacity": int64(4), "state": "OK", "vast_id": int64(9), "name": "alice"}),
	}})

	res := plan(t, "quota", prior, merge(base, map[string]any{"user_quotas": []any{entry(10)}}))
	assert.Equal(t, NoOp, res.Action, "usage the server reports is not drift (changed=%v)", res.Changed)

	res = plan(t, "quota", prior, merge(base, map[string]any{"user_quotas": []any{entry(20)}}))
	assert.Equal(t, Update, res.Action)
	assert.Equal(t, []string{"user_quotas"}, res.Changed)
}

func TestPlan_SetNestedWithoutComputedValues(t *testing.T) {
	addr := func(host string, port int64) map[string]any { return map[string]any{"host": host, "port": port} }
	config := func(addrs ...any) map[string]any { return map[string]any{"name": "k", "addresses": addrs} }
	prior := config(addr("a", 9092), addr("b", 9092))

	assert.Equal(t, NoOp, plan(t, "kafka_broker", prior, config(addr("b", 9092), addr("a", 9092))).Action)

	res := plan(t, "kafka_broker", prior, config(addr("a", 9092), addr("b", 9093)))
	assert.Equal(t, Update, res.Action)
	assert.Equal(t, []string{"addresses"}, res.Changed)
}

// ---- nested lists and maps -------------------------------------------------------------

func TestPlan_ListNestedCorrelatesByIndex(t *testing.T) {
	frame := func(every string) map[string]any {
		return map[string]any{"every": every, "keep_local": "1d", "start_at": "2025-01-01 00:00:00"}
	}
	base := map[string]any{"name": "p", "prefix": "p", "clone_type": "LOCAL"}
	served := func(every string) map[string]any { return merge(frame(every), map[string]any{"keep_remote": "2d"}) }
	prior := merge(base, map[string]any{"id": int64(1), "frames": []any{served("1h"), served("1d")}})

	res := plan(t, "protection_policy", prior, merge(base, map[string]any{"frames": []any{frame("1h"), frame("1d")}}))
	assert.Equal(t, NoOp, res.Action, "computed values in list elements are carried by index (changed=%v)", res.Changed)

	res = plan(t, "protection_policy", prior, merge(base, map[string]any{"frames": []any{frame("1d"), frame("1h")}}))
	assert.Equal(t, Update, res.Action, "list order matters")
	assert.Equal(t, []string{"frames"}, res.Changed)
}

func TestPlan_MapNestedCorrelatesByKey(t *testing.T) {
	pool := func(ip string) map[string]any { return map[string]any{"cnode_ip_pool": []any{ip}} }
	base := map[string]any{"name": "rack-1"}
	prior := merge(base, map[string]any{"id": int64(1), "rack_pools": map[string]any{"r1": pool("10.0.0.1")}})

	assert.Equal(t, NoOp, plan(t, "cluster", prior, merge(base, map[string]any{"rack_pools": map[string]any{"r1": pool("10.0.0.1")}})).Action)

	res := plan(t, "cluster", prior, merge(base, map[string]any{"rack_pools": map[string]any{"r1": pool("10.0.0.1"), "r2": pool("10.0.0.2")}}))
	assert.Equal(t, Update, res.Action)
	assert.Equal(t, []string{"rack_pools"}, res.Changed)

	res = plan(t, "cluster", prior, merge(base, map[string]any{"rack_pools": map[string]any{"r1": pool("10.0.0.9")}}))
	assert.Equal(t, Update, res.Action)
}

// ---- unknown values -------------------------------------------------------------------

// Computed values the plan cannot know are reported as Unknown and reach Planned as null, the
// way Resource.NewManager(req.Plan) reads them. They are also Changed when they replace a
// prior value (`~ x = "a" -> (known after apply)`), but not when there was none.
func TestPlan_UnknownIsSentAsNull(t *testing.T) {
	f := newFixture(t, "view")
	config := f.minimalConfig()
	prior := f.mustApply(config)
	// A computed value the server never returned.
	for _, k := range sortedAttributeNames(f.sch.Attributes) {
		a := f.sch.Attributes[k]
		if a.IsComputed() && !a.IsOptional() {
			prior = with(t, prior, k, tftypes.NewValue(f.typ.AttributeTypes[k], nil))
		}
	}
	changed := copyConfig(config)
	changed["policy_id"] = tftypes.NewValue(tftypes.Number, big.NewFloat(42))
	res := f.mustPlan(prior, changed)

	require.Equal(t, Update, res.Action)
	assert.Equal(t, []string{"policy_id"}, knownChanges(res))
	require.NotEmpty(t, res.Unknown, "a change makes computed values without a prior value unknown")
	priorSt := f.state(prior)
	withoutPrior := 0
	for _, k := range res.Unknown {
		assert.True(t, res.RawValues()[k].IsNull(), "%s: unknown is sent as null", k)
		if priorSt.RawValues()[k].IsNull() {
			withoutPrior++
			assert.NotContains(t, res.Changed, k)
		} else {
			assert.Contains(t, res.Changed, k)
		}
	}
	assert.Positive(t, withoutPrior)
}

// ---- requires replace -----------------------------------------------------------------

// Terraform Core keeps a path the provider asked to replace for only when its planned value
// differs from prior: a later plan modifier can restore the prior value.
func TestRequiresReplace_KeepsOnlyChangedPaths(t *testing.T) {
	ty := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"a": tftypes.String, "b": tftypes.String, "c": tftypes.String}}
	obj := func(a, b, c any) tftypes.Value {
		return tftypes.NewValue(ty, map[string]tftypes.Value{
			"a": tftypes.NewValue(tftypes.String, a), "b": tftypes.NewValue(tftypes.String, b), "c": tftypes.NewValue(tftypes.String, c),
		})
	}
	paths := []*tftypes.AttributePath{
		tftypes.NewAttributePath().WithAttributeName("a"),
		tftypes.NewAttributePath().WithAttributeName("b"),
		tftypes.NewAttributePath().WithAttributeName("c"),
	}
	prior := obj("x", "y", "z")
	planned := obj("x", "changed", tftypes.UnknownValue)

	assert.Equal(t, []string{"b", "c"}, topLevelNames(requiresReplace(paths, prior, planned)),
		"an equal value is dropped, a changed or unknown one kept")
	assert.Empty(t, requiresReplace(paths, tftypes.NewValue(ty, nil), planned), "a creation replaces nothing")
}

func TestPlan_ForceNewOnlyWhenTheValueChanges(t *testing.T) {
	f := newFixture(t, "blob_expansion")
	config := f.minimalConfig()
	prior := f.mustApply(config)
	changed := copyConfig(config)
	changed["tenant_id"] = tftypes.NewValue(tftypes.Number, big.NewFloat(5))

	res := f.mustPlan(prior, changed)
	assert.Equal(t, Update, res.Action, "source_column_name keeps its value, so it does not force a replacement")
	assert.Empty(t, res.RequiresReplace)

	changed["database_name"] = f.value("database_name", 1)
	res = f.mustPlan(prior, changed)
	assert.Equal(t, Replace, res.Action)
	assert.Equal(t, []string{"database_name"}, res.RequiresReplace)
}

// ---- configuration and provider errors -------------------------------------------------

func TestPlan_NestedConfigShapeIsChecked(t *testing.T) {
	ctx := context.Background()
	config := newState(t, "view", map[string]any{
		"path": "/a/", "policy_id": int64(1),
		"create_dir_acl": []any{map[string]any{"grantee": "users"}},
	})
	_, err := Plan(ctx, "view", nil, config)
	require.Error(t, err)
	assert.Contains(t, err.Error(), `AttributeName("create_dir_acl").AttributeName("perm")`)
}

func TestPlan_ProviderValidationIsReported(t *testing.T) {
	config := newState(t, "tenant", map[string]any{"name": "t", "delete_local_provider": true, "local_provider_id": int64(3)})
	_, err := Plan(context.Background(), "tenant", nil, config)
	var de *DiagnosticsError
	require.True(t, errors.As(err, &de), "got %v", err)
	assert.Equal(t, "ValidateResourceConfig", de.Stage)
	assert.Contains(t, err.Error(), "delete_local_provider cannot be set together with local_provider_id")
}

func TestPlan_SchemaValidatorIsReported(t *testing.T) {
	_, err := Plan(context.Background(), "view", nil, newState(t, "view", map[string]any{"path": "no-slash", "policy_id": int64(1)}))
	var de *DiagnosticsError
	require.True(t, errors.As(err, &de), "got %v", err)
	assert.Contains(t, err.Error(), `AttributeName("path")`)
}

func TestPlan_PrepareDesiredStateErrorIsReported(t *testing.T) {
	config := newState(t, "nfs_certificate", map[string]any{"certificate": "not a pem", "private_key": "k", "cluster_id": int64(1)})
	_, err := Plan(context.Background(), "nfs_certificate", nil, config)
	var de *DiagnosticsError
	require.True(t, errors.As(err, &de), "got %v", err)
	assert.Equal(t, "PlanResourceChange", de.Stage)
}

func TestPlan_RequiresConfig(t *testing.T) {
	_, err := Plan(context.Background(), "view", nil, nil)
	require.Error(t, err)
}

func TestPlan_UnknownResourceType(t *testing.T) {
	_, err := Plan(context.Background(), "no_such_resource", nil, newState(t, "view", map[string]any{"path": "/a/", "policy_id": int64(1)}))
	require.Error(t, err)
}

// ---- create ---------------------------------------------------------------------------

// The plan of a creation holds the fingerprint, and Create takes the secret from config.
func TestPlan_CreateCarriesFingerprint(t *testing.T) {
	const pem = "pem-bytes"
	config := newState(t, "tls_certificate", map[string]any{"ca_certificate_file": pem})
	res, err := Plan(context.Background(), "tls_certificate", nil, config)
	require.NoError(t, err)
	require.Equal(t, Create, res.Action)
	assert.True(t, res.IsNull("ca_certificate_file"), "the plan holds no write-only value")
	assert.Equal(t, vsd.ContentMD5Hex(pem), res.String("ca_certificate_file_md5"))
	assert.Contains(t, res.Unknown, "id", "the server assigns the id")

	desired := res.Copy()
	is.OverlayWriteOnly(desired, config)
	assert.Equal(t, pem, desired.String("ca_certificate_file"))
}

// An empty prior is a resource that does not exist: its plan is a creation.
func TestPlan_EmptyPriorPlansCreation(t *testing.T) {
	config := newState(t, "view", map[string]any{"path": "/a/", "policy_id": int64(1)})
	prior := is.MarkEmpty(newState(t, "view", map[string]any{"id": int64(7), "path": "/old/"}))
	res, err := Plan(context.Background(), "view", prior, config)
	require.NoError(t, err)
	assert.Equal(t, Create, res.Action)
	assert.False(t, res.IsEmpty())
	assert.Equal(t, "/a/", res.String("path"))
}

func TestPlan_EmptyConfigIsRefused(t *testing.T) {
	_, err := Plan(context.Background(), "view", nil, is.MarkEmpty(newState(t, "view", nil)))
	require.Error(t, err)
}

// ---- concurrency ----------------------------------------------------------------------

// One planner serves every controller of a process at once.
func TestPlan_ConcurrentUse(t *testing.T) {
	names := sortedResourceNames()
	fixtures := make([]*fixture, len(names))
	for i, name := range names {
		fixtures[i] = newFixture(t, name)
	}
	var wg sync.WaitGroup
	errs := make(chan error, len(names)*4)
	for round := 0; round < 4; round++ {
		for _, f := range fixtures {
			wg.Add(1)
			go func(f *fixture) {
				defer wg.Done()
				config := f.minimalConfig()
				res, err := f.plan(f.null(), config)
				if err == nil && res.Action != Create {
					err = fmt.Errorf("%s: action %s", f.name, res.Action)
				}
				if err != nil {
					errs <- fmt.Errorf("%s: %w", f.name, err)
				}
			}(f)
		}
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
}

func TestNew_IndependentPlanner(t *testing.T) {
	p, err := New("test")
	require.NoError(t, err)
	res, err := p.Plan(context.Background(), "view", nil, newState(t, "view", map[string]any{"path": "/a/", "policy_id": int64(1)}))
	require.NoError(t, err)
	assert.Equal(t, Create, res.Action)
}
