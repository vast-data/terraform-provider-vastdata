// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var RackSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"racks",
	http.MethodGet,
	"racks",
)

type Rack struct {
	tfstate *is.TFState
}

func (m *Rack) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Rack{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			ComputedSchemaFields: []string{"ip_range"},
			SchemaRef:            RackSchemaRef,
		},
	)}
}

func (m *Rack) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Rack{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: RackSchemaRef,
		},
	)}
}

func (m *Rack) TfState() *is.TFState {
	return m.tfstate
}

func (m *Rack) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Racks
}
