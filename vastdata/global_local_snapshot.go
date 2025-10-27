// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
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
	record, err := rest.Snapshots.SnapshotCloneWithContext_POST(ctx, snapId, createParams)
	return record, err
}

func (m *GlobalLocalSnapshot) DeleteResource(ctx context.Context, rest *VMSRest) error {
	ts := m.tfstate
	name := ts.String("name")

	// Get the global snapshot stream by name
	response, err := rest.GlobalSnapshotStreams.GetWithContext(ctx, params{"name": name})
	if err != nil {
		// If not found, nothing to delete
		return ignoreStatusCodes(err, http.StatusNotFound)
	}

	// Extract ID and state from response
	type GssContainer struct {
		Id    int64  `json:"id"`
		State string `json:"state,omitempty"`
	}
	gssContainer := GssContainer{}
	if err = response.Fill(&gssContainer); err != nil {
		return err
	}

	// If not finished, stop the clone snapshot operation first
	if gssContainer.State != "finished" {
		_, err := rest.GlobalSnapshotStreams.GlobalSnapshotStreamStopWithContext_PATCH(ctx, gssContainer.Id, nil, 3*time.Minute)
		if err != nil {
			return err
		}
	}

	// Delete the global snapshot stream with remove_dir parameter
	_, err = rest.GlobalSnapshotStreams.DeleteByIdWithContext(ctx, response.RecordID(), nil, params{"remove_dir": true})

	// Ignore 404 errors (already deleted)
	return ignoreStatusCodes(err, http.StatusNotFound)
}
