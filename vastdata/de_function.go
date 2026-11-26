// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var DeFunctionSchemaRef = is.NewDataEngineSchemaReference(
	http.MethodPost,
	"functions",
	http.MethodGet,
	"functions",
)

type DeFunction struct {
	tfstate is.ResourceState
}

func (m *DeFunction) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &DeFunction{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:   DeFunctionSchemaRef,
			PathIDField: "guid",
			ExcludedSchemaFields: []string{
				"cursor", "limit", "order_by",
				"owner_name", "created_after", "created_before", "updated_after", "updated_before",
			},
			// OpenAPI marks container_registry_vrn required; TF allows name OR vrn.
			NotRequiredSchemaFields: []string{deContainerRegistryVrnAttr},
			// OpenAPI leaves image_tag optional; require it for a valid create body.
			RequiredSchemaFields: []string{"image_tag"},
			// Revision/image fields are not returned on GET /functions/{guid}.
			CreateOnlyFields: []string{deTenantNameAttr},
			PreserveUserValueFields: []string{
				"artifact_source",
				"artifact_type",
				"image_tag",
				"revision_alias",
				"revision_description",
				"force_new_revision",
				"skip_image_validation",
				"is_published",
				"architecture",
				"runtime",
				deContainerRegistryVrnAttr,
				deContainerRegistryNameAttr,
				deTenantNameAttr,
			},
			CommonModifiersMapping: map[string]string{
				deTenantNameAttr: ModifierForceNew,
			},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr:            deTenantNameSchemaAttribute(),
				deContainerRegistryVrnAttr:  deContainerRegistryVrnSchemaAttribute(),
				deContainerRegistryNameAttr: deContainerRegistryNameSchemaAttribute(),
			},
		},
	)}
}

func (m *DeFunction) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &DeFunction{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:   DeFunctionSchemaRef,
			PathIDField: "guid",
			ExcludedSchemaFields: []string{
				"cursor", "limit", "order_by",
				"owner_name", "created_after", "created_before", "updated_after", "updated_before",
			},
			SearchableFields: []string{"guid", "vrn", "name"},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameSchemaAttribute(),
			},
		},
	)}
}

func (m *DeFunction) TfState() is.ResourceState {
	return m.tfstate
}

func (m *DeFunction) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.DataEngine.Functions
}

func (m *DeFunction) PrepareCreateResource(ctx context.Context, rest *VMSRest) error {
	return resolveContainerRegistryNameOrVrnInState(ctx, m.tfstate, rest)
}

func (m *DeFunction) PrepareUpdateResource(ctx context.Context, plan PrepareUpdateResource, rest *VMSRest) error {
	return resolveContainerRegistryNameOrVrnInState(ctx, plan.(*DeFunction).tfstate, rest)
}

// UpdateResource uses PUT /functions/{guid} with a full create-shaped body.
// Diff-only updates omit unchanged required fields (e.g. name) and the API returns 422.
func (m *DeFunction) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	planTS := plan.(*DeFunction).tfstate
	guid, err := pathGUIDOrPlan(m.tfstate, planTS, "de_function")
	if err != nil {
		return nil, err
	}
	body := plan.(*DeFunction).TransformRequestBody(planTS.GetCreateParams())
	if len(body) == 0 {
		return m.ReadResource(ctx, rest)
	}

	api := m.API(rest)
	query := plan.(*DeFunction).TransformQueryParams(nil, rest)
	rec, err := api.UpdateWithContext(ctx, guid, query, body)
	if err != nil {
		return nil, err
	}
	m.tfstate = planTS
	return rec, nil
}

func (m *DeFunction) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	guid, err := requirePathGUID(m.tfstate, "de_function")
	if err != nil {
		return nil, err
	}
	api := m.API(rest)
	query := m.TransformQueryParams(nil, rest)
	return api.GetByIdWithContext(ctx, guid, query)
}

func (m *DeFunction) TransformQueryParams(query params, rest *VMSRest) params {
	return deTransformQueryParams(query, m.tfstate, rest)
}

func (m *DeFunction) TransformRequestBody(body params) params {
	body = deTransformRequestBody(body)
	if body != nil {
		body.Without(
			"guid", "id", "vrn", "created_at", "updated_at", "owner", "tenant_guid",
			"last_published_revision_number", "last_revision_number", "default_revision_number",
			deContainerRegistryNameAttr,
		)
		// 5.5 rejects the field entirely; omit unless the user explicitly set true (5.6+).
		if v, ok := body["force_new_revision"]; !ok || v == nil || v == false {
			body.Without("force_new_revision")
		}
	}
	return body
}

func (m *DeFunction) TransformResponseRecord(response Record) Record {
	response = ensureTenantNameInRecord(response, m.tfstate, nil)
	return preserveContainerRegistryNameInRecord(response, m.tfstate)
}
