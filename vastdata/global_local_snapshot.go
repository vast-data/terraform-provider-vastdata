// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var GlobalLocalSnapshotSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"snapshots/{id}/clone",
	http.MethodGet,
	"snapshots/{id}/clone",
)

type GlobalLocalSnapshot struct {
	tfstate *is.TFState
}

func (m *GlobalLocalSnapshot) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &GlobalLocalSnapshot{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:            GlobalLocalSnapshotSchemaRef,
			RequiredSchemaFields: []string{"name", "loanee_root_path", "loanee_tenant_id", "loanee_snapshot_id"},
		},
	)}
}

func (m *GlobalLocalSnapshot) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &GlobalLocalSnapshot{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: GlobalLocalSnapshotSchemaRef,
		}),
	}
}

func (m *GlobalLocalSnapshot) TfState() *is.TFState {
	return m.tfstate
}

func (m *GlobalLocalSnapshot) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.GlobalSnapshotStreams
}

func (m *GlobalLocalSnapshot) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	snapId := ts.Int64("loanee_snapshot_id")
	createParams := make(params)
	ts.SetToMapIfAvailable(createParams, "name", "loanee_root_path", "loanee_tenant_id", "enabled")
	record, err := rest.GlobalSnapshotStreams.CloneSnapshotWithContext(
		ctx,
		snapId,
		createParams,
	)
	return record, err
}

func (m *GlobalLocalSnapshot) DeleteResource(ctx context.Context, rest *VMSRest) error {
	ts := m.tfstate
	name := ts.String("name")
	_, err := rest.GlobalSnapshotStreams.EnsureCloneSnapshotDeletedWithContext(ctx, params{"name": name})
	return err
}
