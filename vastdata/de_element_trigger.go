// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

// Element triggers are created/updated at /triggers/element but read/deleted
// at /triggers/{guid}. Schema comes from POST /triggers/element.
var DeElementTriggerSchemaRef = is.NewDataEngineSchemaReference(
	http.MethodPost,
	"triggers/element",
	http.MethodGet,
	"triggers",
)

type DeElementTrigger struct {
	tfstate is.ResourceState
}

func (m *DeElementTrigger) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &DeElementTrigger{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:            DeElementTriggerSchemaRef,
			PathIDField:          "guid",
			ExcludedSchemaFields: []string{"cursor", "limit", "order_by", "custom_extensions"},
			CreateOnlyFields:     []string{deTenantNameAttr},
			PreserveUserValueFields: []string{
				deTenantNameAttr,
				"config",
				"broker",
				"topic_name",
				"source_types",
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

func (m *DeElementTrigger) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &DeElementTrigger{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:            DeElementTriggerSchemaRef,
			PathIDField:          "guid",
			ExcludedSchemaFields: []string{"cursor", "limit", "order_by", "custom_extensions"},
			SearchableFields:     []string{"guid", "vrn", "name"},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameSchemaAttribute(),
			},
		},
	)}
}

func (m *DeElementTrigger) TfState() is.ResourceState {
	return m.tfstate
}

// API returns ElementTriggers. Client overrides route get/delete to /triggers/{guid}
// while create/update stay on /triggers/element.
func (m *DeElementTrigger) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.DataEngine.ElementTriggers
}

func (m *DeElementTrigger) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	api := m.API(rest)
	body := m.TransformRequestBody(m.tfstate.GetCreateParams())
	query := m.TransformQueryParams(nil, rest)
	rec, err := api.CreateWithContext(ctx, query, body)
	if err != nil {
		return nil, err
	}
	return normalizeTriggerRecord(rec), nil
}

func (m *DeElementTrigger) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	guid, err := requirePathGUID(m.tfstate, "de_element_trigger")
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

func (m *DeElementTrigger) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	planMgr := plan.(*DeElementTrigger)
	planTS := planMgr.tfstate
	guid, err := pathGUIDOrPlan(m.tfstate, planTS, "de_element_trigger")
	if err != nil {
		return nil, err
	}
	// PUT /triggers/element/{guid} expects a full create-shaped body (config required).
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

func (m *DeElementTrigger) DeleteResource(ctx context.Context, rest *VMSRest) error {
	guid, err := requirePathGUID(m.tfstate, "de_element_trigger")
	if err != nil {
		return err
	}
	api := m.API(rest)
	query := m.TransformQueryParams(nil, rest)
	_, err = api.DeleteByIdWithContext(ctx, guid, query, nil)
	return ignoreResourceGone(err)
}

func (m *DeElementTrigger) TransformQueryParams(query params, rest *VMSRest) params {
	return deTransformQueryParams(query, m.tfstate, rest)
}

func (m *DeElementTrigger) TransformRequestBody(body params) params {
	body = deTransformRequestBody(body)
	if body != nil {
		body.Without(
			"guid", "id", "vrn", "created_at", "updated_at", "owner", "tenant_guid",
			"status", "type", "topic", "events", "source_bucket_name",
		)
	}
	return body
}

func (m *DeElementTrigger) TransformResponseRecord(response Record) Record {
	return ensureTenantNameInRecord(normalizeTriggerRecord(response), m.tfstate, nil)
}
