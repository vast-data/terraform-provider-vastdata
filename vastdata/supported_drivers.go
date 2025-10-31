// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var SupportedDriversSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"supporteddrives",
	http.MethodGet,
	"supporteddrives",
)

type SupportedDrivers struct {
	tfstate *is.TFState
}

func (m *SupportedDrivers) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &SupportedDrivers{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: SupportedDriversSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"model_name": rschema.StringAttribute{
					Required:    true,
					Description: "Model name for the supported driver",
				},
			},
		},
	)}
}

func (m *SupportedDrivers) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &SupportedDrivers{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: SupportedDriversSchemaRef,
		},
	)}
}

func (m *SupportedDrivers) TfState() *is.TFState {
	return m.tfstate
}

func (m *SupportedDrivers) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.SupportedDrivers
}
