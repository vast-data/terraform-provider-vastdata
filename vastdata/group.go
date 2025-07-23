// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var GroupSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"groups",
	http.MethodGet,
	"groups",
)

type Group struct {
	tfstate *is.TFState
}

func (m *Group) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Group{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: GroupSchemaRef,
		},
	)}
}

func (m *Group) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Group{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: GroupSchemaRef,
		},
	)}
}

func (m *Group) TfState() *is.TFState {
	return m.tfstate
}

func (m *Group) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Groups
}
