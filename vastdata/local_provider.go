// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var LocalProviderSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"localproviders",
	http.MethodGet,
	"localproviders",
)

type LocalProvider struct {
	tfstate is.ResourceState
}

func (m *LocalProvider) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &LocalProvider{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: LocalProviderSchemaRef,
		},
	)}
}

func (m *LocalProvider) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &LocalProvider{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: LocalProviderSchemaRef,
		},
	)}
}

func (m *LocalProvider) TfState() is.ResourceState {
	return m.tfstate
}

func (m *LocalProvider) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.LocalProviders
}
