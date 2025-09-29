// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var CnodeSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"cnodes",
	http.MethodGet,
	"cnodes",
)

type Cnode struct {
	tfstate *is.TFState
}

func (m *Cnode) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Cnode{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			ReadOnlyFields: []string{"cluster_name", "vippool_id"},
			SchemaRef:      CnodeSchemaRef,
		},
	)}
}

func (m *Cnode) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Cnode{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: CnodeSchemaRef,
		},
	)}
}

func (m *Cnode) TfState() *is.TFState {
	return m.tfstate
}

func (m *Cnode) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Cnodes
}
