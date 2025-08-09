// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
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
			SchemaRef:      ViewPolicySchemaRef,
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
	return rest.ViewPolies
}
