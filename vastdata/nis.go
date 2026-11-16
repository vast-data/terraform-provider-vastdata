// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var NisSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"nis",
	http.MethodGet,
	"nis",
)

type Nis struct {
	tfstate is.ResourceState
}

func (m *Nis) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Nis{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: NisSchemaRef,
		},
	)}
}

func (m *Nis) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Nis{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: NisSchemaRef,
		},
	)}
}

func (m *Nis) TfState() is.ResourceState {
	return m.tfstate
}

func (m *Nis) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Nis
}
