// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var ProtectedPathSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"protectedpaths",
	http.MethodGet,
	"protectedpaths",
)

type ProtectedPath struct {
	tfstate *is.TFState
}

func (m *ProtectedPath) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &ProtectedPath{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: ProtectedPathSchemaRef,
		},
	)}
}

func (m *ProtectedPath) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ProtectedPath{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: ProtectedPathSchemaRef,
		},
	)}
}

func (m *ProtectedPath) TfState() *is.TFState {
	return m.tfstate
}

func (m *ProtectedPath) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.ProtectedPaths
}
