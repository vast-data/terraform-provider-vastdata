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
	tfstate *is.TFState
}

func (m *Oidc) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Oidc{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: OidcSchemaRef,
		},
	)}
}

func (m *Oidc) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Oidc{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: OidcSchemaRef,
		},
	)}
}

func (m *Oidc) TfState() *is.TFState {
	return m.tfstate
}

func (m *Oidc) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Oidcs
}
