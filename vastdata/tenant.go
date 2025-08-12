// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
			SchemaRef:             TenantSchemaRef,
			DeleteOnlyParamFields: []string{"force"},
			PreserveOrderFields:   []string{"client_ip_ranges"},
			AdditionalSchemaAttributes: map[string]any{
				"force": rschema.BoolAttribute{
					Optional: true,
					Description: "If set to true, forces deletion of the tenant even if it has empty subdirectories" +
						" or other removable remnants. Use with caution, as this will bypass standard cleanup checks.",
				},
			},
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
