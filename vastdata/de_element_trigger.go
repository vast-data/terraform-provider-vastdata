// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
			ExcludedSchemaFields: []string{"cursor", "limit", "order_by"},
			CreateOnlyFields:     []string{deTenantNameAttr},
			PreserveUserValueFields: []string{
				deTenantNameAttr,
				"config",
				"broker",
				"topic_name",
				"source_types",
				deCustomExtensionsAttr,
			},
			CommonModifiersMapping: map[string]string{
				deTenantNameAttr: ModifierForceNew,
			},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr:       deTenantNameSchemaAttribute(),
				deCustomExtensionsAttr: deCustomExtensionsSchemaAttribute(),
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
			ExcludedSchemaFields: []string{"cursor", "limit", "order_by"},
			SearchableFields:     []string{"guid", "vrn", "name"},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr:       deTenantNameSchemaAttribute(),
				deCustomExtensionsAttr: deCustomExtensionsSchemaAttribute(),
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

// PrepareCreateResource enforces create checks that OpenAPI leaves optional:
// non-empty config.events + config.source_bucket_name, and topic_name ↔ broker together.
func (m *DeElementTrigger) PrepareCreateResource(context.Context, *VMSRest) error {
	return validateDeElementTriggerConfig(m.tfstate)
}

func (m *DeElementTrigger) PrepareUpdateResource(_ context.Context, plan PrepareUpdateResource, _ *VMSRest) error {
	return validateDeElementTriggerConfig(plan.(*DeElementTrigger).tfstate)
}

func validateDeElementTriggerConfig(ts is.ResourceState) error {
	if ts == nil {
		return nil
	}

	// Skip while plan refs are still unknown.
	if v, ok := ts.RawValues()["config"]; ok && v.IsUnknown() {
		return nil
	}
	if v, ok := ts.RawValues()["topic_name"]; ok && v.IsUnknown() {
		return nil
	}
	if v, ok := ts.RawValues()["broker"]; ok && v.IsUnknown() {
		return nil
	}

	if err := validateDeElementTriggerRequiredConfig(ts); err != nil {
		return err
	}
	return validateDeElementTriggerBrokerTopic(ts)
}

func validateDeElementTriggerRequiredConfig(ts is.ResourceState) error {
	cfgVal, ok := ts.RawValues()["config"]
	if !ok || cfgVal.IsNull() {
		return fmt.Errorf("config is required")
	}
	cfg, ok := cfgVal.(types.Object)
	if !ok {
		return fmt.Errorf("config: expected object")
	}
	attrs := cfg.Attributes()

	eventsVal, hasEvents := attrs["events"]
	if !hasEvents || eventsVal.IsNull() {
		return fmt.Errorf("config.events is required")
	}
	if eventsVal.IsUnknown() {
		return nil
	}
	switch e := eventsVal.(type) {
	case types.Set:
		if len(e.Elements()) == 0 {
			return fmt.Errorf("config.events must not be empty")
		}
	case types.List:
		if len(e.Elements()) == 0 {
			return fmt.Errorf("config.events must not be empty")
		}
	}

	bucketVal, hasBucket := attrs["source_bucket_name"]
	if !hasBucket || bucketVal.IsNull() {
		return fmt.Errorf("config.source_bucket_name is required")
	}
	if bucketVal.IsUnknown() {
		return nil
	}
	if s, ok := bucketVal.(types.String); ok && strings.TrimSpace(s.ValueString()) == "" {
		return fmt.Errorf("config.source_bucket_name is required")
	}
	return nil
}

// validateDeElementTriggerBrokerTopic enforces broker/topic pairing:
// topic_name and broker are both set or both omitted. When either is present,
// validateAllOf requires the other.
func validateDeElementTriggerBrokerTopic(ts is.ResourceState) error {
	topicSet := false
	if v, ok := ts.RawValues()["topic_name"]; ok && !v.IsNull() && !v.IsUnknown() {
		if s, ok := v.(types.String); ok && strings.TrimSpace(s.ValueString()) != "" {
			topicSet = true
		}
	}
	brokerSet := false
	if v, ok := ts.RawValues()["broker"]; ok && !v.IsNull() && !v.IsUnknown() {
		if obj, ok := v.(types.Object); ok {
			attrs := obj.Attributes()
			name, hasName := attrs["name"]
			typ, hasType := attrs["type"]
			nameSet := hasName && !name.IsNull() && !name.IsUnknown() &&
				strings.TrimSpace(name.(types.String).ValueString()) != ""
			typeSet := hasType && !typ.IsNull() && !typ.IsUnknown() &&
				strings.TrimSpace(typ.(types.String).ValueString()) != ""
			brokerSet = nameSet || typeSet
		}
	}
	if !topicSet && !brokerSet {
		return nil
	}
	return validateAllOf(ts, "topic_name", "broker")
}
