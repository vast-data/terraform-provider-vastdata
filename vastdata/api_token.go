// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var ApiTokenSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"apitokens",
	http.MethodGet,
	"apitokens",
)

type ApiToken struct {
	tfstate *is.TFState
}

func (m *ApiToken) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &ApiToken{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:      ApiTokenSchemaRef,
			ReadOnlyFields: []string{"archived"},
		},
	)}
}

func (m *ApiToken) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ApiToken{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: ApiTokenSchemaRef,
		},
	)}
}

func (m *ApiToken) TfState() *is.TFState {
	return m.tfstate
}

func (m *ApiToken) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.ApiTokens
}

func (m *ApiToken) DeleteResource(ctx context.Context, rest *VMSRest) error {
	ts := m.tfstate
	_, err := rest.ApiTokens.ApiTokenRevokeWithContext_PATCH(ctx, ts.String("id"), nil)
	return ignoreStatusCodes(err, http.StatusNotFound)
}
