// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var GlobalSnapshotSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"globalsnapstreams",
	http.MethodGet,
	"globalsnapstreams",
)

type GlobalSnapshot struct {
	tfstate *is.TFState
}

func (m *GlobalSnapshot) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &GlobalSnapshot{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: GlobalSnapshotSchemaRef,
		},
	)}
}

func (m *GlobalSnapshot) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &GlobalSnapshot{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: GlobalSnapshotSchemaRef,
		}),
	}
}

func (m *GlobalSnapshot) TfState() *is.TFState {
	return m.tfstate
}

func (m *GlobalSnapshot) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.GlobalSnapshotStreams
}
