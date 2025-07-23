// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var DnsSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"dns",
	http.MethodGet,
	"dns",
)

type Dns struct {
	tfstate *is.TFState
}

func (m *Dns) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Dns{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: DnsSchemaRef,
		},
	)}
}

func (m *Dns) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Dns{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: DnsSchemaRef,
		}),
	}
}

func (m *Dns) TfState() *is.TFState {
	return m.tfstate
}

func (m *Dns) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Dns
}
