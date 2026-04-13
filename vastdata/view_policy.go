// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var ViewPolicySchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"viewpolicies",
	http.MethodGet,
	"viewpolicies",
)

type ViewPolicy struct {
	tfstate *is.TFState
}

func (m *ViewPolicy) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &ViewPolicy{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: ViewPolicySchemaRef,
			// The VAST backend runs collapse_addresses (Python ipaddress module) on NFS
			// IP allow-list fields when a policy is saved, consolidating adjacent host IPs
			// into the minimal CIDR set.  PreserveUserValueFields keeps the user's declared
			// list intact in state (only when the field is non-null, i.e. the user set it),
			// so the API's rewritten CIDRs never cause spurious drift.  Fields the user
			// leaves unset remain null in plan and are still populated from the API normally.
			PreserveUserValueFields: []string{
				"protocols_audit",
				"nfs_read_write",
				"nfs_read_only",
				"nfs_no_squash",
				"nfs_root_squash",
				"nfs_all_squash",
			},
			ReadOnlyFields: []string{"serves_tenant"},
			ImportFields:   []string{"name", "tenant_name"},
		},
	)}
}

func (m *ViewPolicy) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ViewPolicy{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: ViewPolicySchemaRef,
		}),
	}
}

func (m *ViewPolicy) TfState() *is.TFState {
	return m.tfstate
}

func (m *ViewPolicy) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.ViewPolicies
}
