// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var TenantSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"tenants",
	http.MethodGet,
	"tenants",
)

type Tenant struct {
	tfstate *is.TFState
}

func (m *Tenant) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Tenant{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:           TenantSchemaRef,
			PreserveOrderFields: []string{"client_ip_ranges"},
		},
	)}
}

func (m *Tenant) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Tenant{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:           TenantSchemaRef,
			PreserveOrderFields: []string{"client_ip_ranges"},
		}),
	}
}

func (m *Tenant) TfState() *is.TFState {
	return m.tfstate
}

func (m *Tenant) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Tenants
}
