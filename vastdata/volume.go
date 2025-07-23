// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var VolumeSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"volumes",
	http.MethodGet,
	"volumes",
)

type Volume struct {
	tfstate *is.TFState
}

func (m *Volume) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Volume{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: VolumeSchemaRef,
		},
	)}
}

func (m *Volume) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Volume{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: VolumeSchemaRef,
		}),
	}
}

func (m *Volume) TfState() *is.TFState {
	return m.tfstate
}

func (m *Volume) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Volumes
}
