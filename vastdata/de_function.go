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
				"container_registry_vrn",
				deTenantNameAttr,
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
	// Image/registry changes should always mint a new revision so pipelines pull the new artifact.
	if m.tfstate != nil && planTS != nil {
		changed := planTS.GetUpdateParams(m.tfstate)
		for _, k := range []string{"artifact_source", "artifact_type", "image_tag", "container_registry_vrn"} {
			if _, ok := changed[k]; ok {
				body["force_new_revision"] = true
				break
			}
		}
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
		)
	}
	return body
}

func (m *DeFunction) TransformResponseRecord(response Record) Record {
	return ensureTenantNameInRecord(response, m.tfstate, nil)
}
