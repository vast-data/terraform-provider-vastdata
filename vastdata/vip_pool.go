// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var VipPoolSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"vippools",
	http.MethodGet,
	"vippools",
)

type VipPool struct {
	tfstate *is.TFState
}

func (m *VipPool) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &VipPool{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:           VipPoolSchemaRef,
			ReadOnlyFields:      []string{"serves_tenant"},
			PreserveOrderFields: []string{"ip_ranges"},
		},
	)}
}

func (m *VipPool) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &VipPool{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:           VipPoolSchemaRef,
			PreserveOrderFields: []string{"ip_ranges"},
		}),
	}
}

func (m *VipPool) TfState() *is.TFState {
	return m.tfstate
}

func (m *VipPool) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.VipPools
}
