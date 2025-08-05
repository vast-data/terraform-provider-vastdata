// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var BgpConfigSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"bgpconfigs",
	http.MethodGet,
	"bgpconfigs",
)

type BgpConfig struct {
	tfstate *is.TFState
}

func (m *BgpConfig) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &BgpConfig{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: BgpConfigSchemaRef,
		},
	)}
}

func (m *BgpConfig) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &BgpConfig{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: BgpConfigSchemaRef,
		},
	)}
}

func (m *BgpConfig) TfState() *is.TFState {
	return m.tfstate
}

func (m *BgpConfig) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.BGPConfigs
}
