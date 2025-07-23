// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var SnapshotSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"snapshots",
	http.MethodGet,
	"snapshots",
)

type Snapshot struct {
	tfstate *is.TFState
}

func (m *Snapshot) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Snapshot{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:      SnapshotSchemaRef,
			ReadOnlyFields: []string{"volume_id"},
			CommonValidatorsMapping: map[string]string{
				"path":            ValidatorPathStartsEndsWithSlash,
				"expiration_time": ValidatorRFC3339Format,
			},
		},
	)}
}

func (m *Snapshot) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Snapshot{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: SnapshotSchemaRef,
		}),
	}
}

func (m *Snapshot) TfState() *is.TFState {
	return m.tfstate
}

func (m *Snapshot) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Snapshots
}
