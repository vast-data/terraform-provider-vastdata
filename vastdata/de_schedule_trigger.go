// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/vast-data/go-vast-client/core"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

// Schedule triggers are created/updated at /triggers/schedule but read/deleted
// at /triggers/{guid}. Schema comes from POST /triggers/schedule.
var DeScheduleTriggerSchemaRef = is.NewDataEngineSchemaReference(
	http.MethodPost,
	"triggers/schedule",
	http.MethodGet,
	"triggers",
)

type DeScheduleTrigger struct {
	tfstate is.ResourceState
}

func (m *DeScheduleTrigger) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &DeScheduleTrigger{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:            DeScheduleTriggerSchemaRef,
			PathIDField:          "guid",
			ExcludedSchemaFields: []string{"cursor", "limit", "order_by", "custom_extensions"},
			CreateOnlyFields:     []string{deTenantNameAttr},
			PreserveUserValueFields: []string{
				deTenantNameAttr,
				"config",
				"broker",
				"topic_name",
			},
			CommonModifiersMapping: map[string]string{
				deTenantNameAttr: ModifierForceNew,
			},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameSchemaAttribute(),
			},
		},
	)}
}

func (m *DeScheduleTrigger) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &DeScheduleTrigger{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:            DeScheduleTriggerSchemaRef,
			PathIDField:          "guid",
			ExcludedSchemaFields: []string{"cursor", "limit", "order_by", "custom_extensions"},
			SearchableFields:     []string{"guid", "vrn", "name"},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameSchemaAttribute(),
			},
		},
	)}
}

func (m *DeScheduleTrigger) TfState() is.ResourceState {
	return m.tfstate
}

// API returns ScheduleTriggers. Client overrides route get/delete to /triggers/{guid}
// while create/update stay on /triggers/schedule.
func (m *DeScheduleTrigger) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.DataEngine.ScheduleTriggers
}

func (m *DeScheduleTrigger) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	api := m.API(rest)
	body := m.TransformRequestBody(m.tfstate.GetCreateParams())
	query := core.Params(m.TransformQueryParams(nil, rest))
	rec, err := api.CreateWithContext(ctx, query, core.Params(body))
	if err != nil {
		return nil, err
	}
	return normalizeTriggerRecord(rec), nil
}

func (m *DeScheduleTrigger) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	guid, err := requirePathGUID(m.tfstate, "de_schedule_trigger")
	if err != nil {
		return nil, err
	}
	api := m.API(rest)
	query := m.TransformQueryParams(nil, rest)
	rec, err := api.GetByIdWithContext(ctx, guid, query)
	if err != nil {
		return nil, err
	}
	return normalizeTriggerRecord(rec), nil
}

func (m *DeScheduleTrigger) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	planMgr := plan.(*DeScheduleTrigger)
	planTS := planMgr.tfstate
	guid, err := pathGUIDOrPlan(m.tfstate, planTS, "de_schedule_trigger")
	if err != nil {
		return nil, err
	}
	// PUT /triggers/schedule/{guid} expects a full create-shaped body (config required).
	body := planMgr.TransformRequestBody(planTS.GetCreateParams())
	if len(body) == 0 {
		return m.ReadResource(ctx, rest)
	}
	api := m.API(rest)
	query := planMgr.TransformQueryParams(nil, rest)
	rec, err := api.UpdateWithContext(ctx, guid, query, body)
	if err != nil {
		return nil, err
	}
	m.tfstate = planTS
	return normalizeTriggerRecord(rec), nil
}

func (m *DeScheduleTrigger) DeleteResource(ctx context.Context, rest *VMSRest) error {
	guid, err := requirePathGUID(m.tfstate, "de_schedule_trigger")
	if err != nil {
		return err
	}
	api := m.API(rest)
	query := m.TransformQueryParams(nil, rest)
	_, err = api.DeleteByIdWithContext(ctx, guid, query, nil)
	return ignoreResourceGone(err)
}

func (m *DeScheduleTrigger) TransformQueryParams(query params, rest *VMSRest) params {
	return deTransformQueryParams(query, m.tfstate, rest)
}

func (m *DeScheduleTrigger) TransformRequestBody(body params) params {
	body = deTransformRequestBody(body)
	if body != nil {
		body.Without(
			"guid", "id", "vrn", "created_at", "updated_at", "owner", "tenant_guid",
			"status", "type", "topic", "events", "source_bucket_name",
		)
	}
	return body
}

func (m *DeScheduleTrigger) TransformResponseRecord(response Record) Record {
	return ensureTenantNameInRecord(normalizeTriggerRecord(response), m.tfstate, nil)
}

// normalizeTriggerRecord coerces GET /triggers/{guid} config (often a JSON string)
// into a map so it matches the nested ScheduleTriggerCreate schema.
func normalizeTriggerRecord(rec Record) Record {
	if rec == nil {
		return nil
	}
	out := Record{}
	for k, v := range rec {
		out[k] = v
	}
	if cfg, ok := out["config"]; ok {
		out["config"] = normalizeTriggerConfig(cfg)
	}
	if broker, ok := out["broker"]; ok {
		out["broker"] = normalizeKafkaBroker(broker)
	}
	return out
}

func normalizeTriggerConfig(v any) any {
	switch c := v.(type) {
	case map[string]any:
		return c
	case Record:
		return map[string]any(c)
	case string:
		s := strings.TrimSpace(c)
		if s == "" {
			return nil
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(s), &m); err != nil {
			// Treat bare cron expression as cron_schedule.
			return map[string]any{"cron_schedule": s}
		}
		return m
	default:
		return v
	}
}
