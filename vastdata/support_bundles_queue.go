// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var SupportBundlesQueueSchemaRef = is.NewSchemaReference(
	http.MethodGet,
	"supportbundlesqueue",
	http.MethodGet,
	"supportbundlesqueue",
)

type SupportBundlesQueue struct {
	tfstate *is.TFState
}

func (m *SupportBundlesQueue) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &SupportBundlesQueue{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: SupportBundlesQueueSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"position": rschema.Int64Attribute{
					Optional:    true,
					Computed:    true,
					Description: "Position of the support bundle in the queue (1 = first). Use this to reorder via PATCH /supportbundlesqueue/{id}/move/.",
				},
			},
			PreserveUserValueFields: []string{"position"},
		},
	)}
}

func (m *SupportBundlesQueue) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &SupportBundlesQueue{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: SupportBundlesQueueSchemaRef,
		},
	)}
}

func (m *SupportBundlesQueue) TfState() *is.TFState {
	return m.tfstate
}

func (m *SupportBundlesQueue) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.SupportBundlesQueue
}

// CreateResource adopts an existing queue entry — bundles enter the queue when
// created via POST /supportbundles/. There is no POST /supportbundlesqueue/.
func (m *SupportBundlesQueue) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	record, err := rest.SupportBundlesQueue.GetWithContext(ctx, m.tfstate.GetGenericSearchParams(ctx))
	if err != nil {
		return nil, fmt.Errorf("could not find support bundle queue entry: %w", err)
	}
	return m.moveToPositionIfNeeded(ctx, rest, record, m.tfstate.Int64("position"))
}

// UpdateResource calls PATCH /supportbundlesqueue/{id}/move/ when position changes.
func (m *SupportBundlesQueue) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	id := m.tfstate.Int64("id")
	if id == 0 {
		return nil, fmt.Errorf("support bundle queue entry id is not set")
	}
	planTs := plan.(*SupportBundlesQueue).TfState()
	position := planTs.Int64("position")
	if position == 0 {
		return nil, fmt.Errorf("position must be >= 1")
	}
	record := Record{"id": id}
	return m.moveToPositionIfNeeded(ctx, rest, record, position)
}

func (m *SupportBundlesQueue) moveToPositionIfNeeded(ctx context.Context, rest *VMSRest, record DisplayableRecord, position int64) (DisplayableRecord, error) {
	if position == 0 {
		return record, nil
	}
	if current, ok := positionInQueueFromRecord(record); ok && current == position {
		return record, nil
	}
	id := record.(Record).RecordID()
	if id == 0 {
		return nil, fmt.Errorf("support bundle queue entry id is not set")
	}
	return rest.SupportBundlesQueue.SupportBundlesQueueMoveWithContext_PATCH(ctx, id, params{"position": position})
}

func positionInQueueFromRecord(record DisplayableRecord) (int64, bool) {
	v, ok := record.(Record)["position_in_queue"]
	if !ok {
		return 0, false
	}
	pos, err := is.ToInt(v)
	if err != nil {
		return 0, false
	}
	return pos, true
}

// DeleteResource is a no-op — queue entries are managed by the support bundle
// lifecycle and cannot be deleted directly from the queue.
func (m *SupportBundlesQueue) DeleteResource(_ context.Context, _ *VMSRest) error {
	return nil
}
