// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var GroupSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"groups",
	http.MethodGet,
	"groups",
)

type Group struct {
	tfstate is.ResourceState
}

func (m *Group) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Group{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			ComputedSchemaFields: []string{"local_provider_id"},
			SchemaRef:            GroupSchemaRef,
		},
	)}
}

func (m *Group) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Group{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			ComputedSchemaFields: []string{"local_provider_id"},
			SchemaRef:            GroupSchemaRef,
		},
	)}
}

func (m *Group) TfState() is.ResourceState {
	return m.tfstate
}

func (m *Group) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Groups
}
