// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
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
			SchemaRef:            VolumeSchemaRef,
			OptionalSchemaFields: []string{"is_monitored"},
			EditOnlyFields:       []string{"is_monitored"},
			ReadOnlyFields:       []string{"tenant_id", "snapshot_id", "mapped_snapshot_id"},
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
