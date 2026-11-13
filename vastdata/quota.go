// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

const quotaSchemaVersion int64 = 1

var QuotaSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"quotas",
	http.MethodGet,
	"quotas",
)

type Quota struct {
	tfstate is.ResourceState
}

func (m *Quota) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Quota{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			PreserveUserValueFields: []string{
				"default_group_quota",
				"default_user_quota",
				"grace_period",
				"group_quotas",
				"user_quotas",
			},
			// The VAST API stores 0 (no limit) as null.
			PreserveUserValueFieldsWhenApiReturnsNull: []string{
				"hard_limit",
				"soft_limit",
				"hard_limit_inodes",
				"soft_limit_inodes",
			},
			// Create-only: API rejects PATCH (including null). Omitting after create must not plan a clear (TERF-224).
			CreateOnlyFields: []string{"is_physical_quota"},
			SchemaRef:        QuotaSchemaRef,
		},
	)}
}

func (m *Quota) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Quota{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: QuotaSchemaRef,
		}),
	}
}

func (m *Quota) TfState() is.ResourceState {
	return m.tfstate
}

func (m *Quota) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Quotas
}

func (m *Quota) PrepareUpdateResource(_ context.Context, plan PrepareUpdateResource, _ *VMSRest) error {
	planTs := plan.(*Quota).tfstate
	if !planTs.IsKnownAndNotNull("is_physical_quota") {
		return nil
	}
	return ensureNotChanged(m.tfstate, planTs, "is_physical_quota")
}

// TransformRequestBody removes is_group from embedded user_quotas / group_quotas
// if present. Sending it in PATCH/POST /quotas/ breaks VMS
// (permissions_userquota has no such column).
func (m *Quota) TransformRequestBody(body params) params {
	stripIsGroupFromQuotaRuleLists(body)
	return body
}

func (m *Quota) SchemaVersion() int64 { return quotaSchemaVersion }

// UpgradeResourceSchema drops is_group from nested quota rules left in prior state.
// Shared UpgradeState already unmarshals with IgnoreUndefinedAttributes; this
// makes the cleanup explicit for nested set elements.
func (m *Quota) UpgradeResourceSchema(_ context.Context, _ int64, ts *is.TFState) error {
	if ts == nil || !ts.IsEnabled() {
		return nil
	}
	for _, key := range []string{"user_quotas", "group_quotas"} {
		v, ok := ts.RawValues()[key]
		if !ok || v == nil || v.IsNull() || v.IsUnknown() {
			continue
		}
		typ := ts.Type(key)
		raw := is.ConvertAttrValueToRaw(v, typ)
		stripIsGroupFromQuotaRuleLists(map[string]any{key: raw})
		newVal, _, err := is.BuildAttrValueFromAny(typ, raw)
		if err != nil {
			return err
		}
		ts.RawValues()[key] = newVal
	}
	return nil
}

func stripIsGroupFromQuotaRuleLists(body map[string]any) {
	for _, key := range []string{"user_quotas", "group_quotas"} {
		raw, ok := body[key]
		if !ok || raw == nil {
			continue
		}
		items, ok := raw.([]any)
		if !ok {
			continue
		}
		for _, item := range items {
			if rule, ok := item.(map[string]any); ok {
				delete(rule, "is_group")
			}
		}
	}
}
