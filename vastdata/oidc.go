// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var OidcSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"oidcs",
	http.MethodGet,
	"oidcs",
)

type Oidc struct {
	tfstate is.ResourceState
}

func (m *Oidc) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Oidc{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: OidcSchemaRef,
		},
	)}
}

func (m *Oidc) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Oidc{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: OidcSchemaRef,
		},
	)}
}

func (m *Oidc) TfState() is.ResourceState {
	return m.tfstate
}

func (m *Oidc) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Oidcs
}
