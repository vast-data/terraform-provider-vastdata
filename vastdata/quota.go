// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var QuotaSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"quotas",
	http.MethodGet,
	"quotas",
)

type Quota struct {
	tfstate *is.TFState
}

func (m *Quota) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Quota{tfstate: is.NewTFStateMust(
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
	return &Quota{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: QuotaSchemaRef,
		}),
	}
}

func (m *Quota) TfState() *is.TFState {
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
