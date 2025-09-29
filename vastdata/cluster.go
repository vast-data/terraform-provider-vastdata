// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var ClusterSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"clusters",
	http.MethodGet,
	"clusters",
)

type Cluster struct {
	tfstate *is.TFState
}

func (m *Cluster) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Cluster{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: ClusterSchemaRef,
		},
	)}
}

func (m *Cluster) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Cluster{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: ClusterSchemaRef,
		},
	)}
}

func (m *Cluster) TfState() *is.TFState {
	return m.tfstate
}

func (m *Cluster) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Clusters
}
