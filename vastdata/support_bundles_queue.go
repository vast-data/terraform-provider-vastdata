// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
	tfstate is.ResourceState
}

func (m *SupportBundlesQueue) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &SupportBundlesQueue{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:          SupportBundlesQueueSchemaRef,
			WriteOnlyFields:    []string{"access_key", "secret_key"},
			WriteOnlyMD5Fields: []string{"access_key", "secret_key"},
			SensitiveFields:    []string{"access_key", "secret_key"},
			AdditionalSchemaAttributes: map[string]any{
				"position": rschema.Int64Attribute{
					Optional:    true,
					Computed:    true,
					Description: "Position of the support bundle in the queue (1 = first). Use this to reorder via PATCH /supportbundlesqueue/{id}/move/.",
				},
				"access_key_md5": ContentMD5SchemaAttribute("access_key", false),
				"secret_key_md5": ContentMD5SchemaAttribute("secret_key", false),
			},
			PreserveUserValueFields: []string{"position"},
		},
	)}
}

func (m *SupportBundlesQueue) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &SupportBundlesQueue{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: SupportBundlesQueueSchemaRef,
		},
	)}
}

func (m *SupportBundlesQueue) TfState() is.ResourceState {
	return m.tfstate
}

func (m *SupportBundlesQueue) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.SupportBundlesQueue
}

func (m *SupportBundlesQueue) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) error {
	SetPlanContentMD5Fields(ctx, req, resp,
		"access_key", "access_key_md5",
		"secret_key", "secret_key_md5",
	)
	return nil
}

func (m *SupportBundlesQueue) TransformResponseRecord(record Record) Record {
	DeleteRecordFields(record, "access_key", "secret_key")
	return record
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

func (m *SupportBundlesQueue) moveToPositionIfNeeded(ctx context.Context, rest *VMSRest, record Record, position int64) (DisplayableRecord, error) {
	if position == 0 {
		return record, nil
	}
	// Create-only skip: Update passes Record{"id": id} (no position_in_queue), so it always PATCHes — which is fine because the framework only calls Update when position changed.
	if v, ok := record["position_in_queue"]; ok {
		if current, err := is.ToInt(v); err == nil && current == position {
			return record, nil
		}
	}
	id := record.RecordID()
	return rest.SupportBundlesQueue.SupportBundlesQueueMoveWithContext_PATCH(ctx, id, params{"position": position})
}

// DeleteResource is a no-op — queue entries are managed by the support bundle
// lifecycle and cannot be deleted directly from the queue.
func (m *SupportBundlesQueue) DeleteResource(_ context.Context, _ *VMSRest) error {
	return nil
}
