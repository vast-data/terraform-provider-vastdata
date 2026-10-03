// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var DeContainerRegistrySchemaRef = is.NewDataEngineSchemaReference(
	http.MethodPost,
	"container-registries",
	http.MethodGet,
	"container-registries",
)

type DeContainerRegistry struct {
	tfstate is.ResourceState
}

func (m *DeContainerRegistry) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &DeContainerRegistry{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:            DeContainerRegistrySchemaRef,
			PathIDField:          "guid",
			ExcludedSchemaFields: []string{"cursor", "limit", "order_by"},
			SearchableFields:     []string{"guid", "vrn", "name"},
			SensitiveFields: []string{
				"password",
				"secret",
			},
			// API rejects PATCH of name ("Can't update container registry name").
			CreateOnlyFields:        []string{"name", deTenantNameAttr},
			PreserveUserValueFields: []string{deTenantNameAttr},
			CommonModifiersMapping: map[string]string{
				"name":           ModifierForceNew,
				deTenantNameAttr: ModifierForceNew,
			},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameSchemaAttribute(),
			},
		},
	)}
}

func (m *DeContainerRegistry) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &DeContainerRegistry{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:            DeContainerRegistrySchemaRef,
			PathIDField:          "guid",
			ExcludedSchemaFields: []string{"cursor", "limit", "order_by"},
			SensitiveFields:      []string{"secret"},
			SearchableFields:     []string{"guid", "vrn", "name"},
			NotOptionalSchemaFields: []string{
				"created_at",
				"updated_at",
				"url",
				"is_accessible_to_data_engine",
				"secret",
				"description",
				"auth_type",
				"id",
				"tenant_guid",
			},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameSchemaAttribute(),
			},
		},
	)}
}

func (m *DeContainerRegistry) TransformQueryParams(query params, rest *VMSRest) params {
	return deTransformQueryParams(query, m.tfstate, rest)
}

func (m *DeContainerRegistry) TransformRequestBody(body params) params {
	return deTransformRequestBody(body)
}

func (m *DeContainerRegistry) TransformResponseRecord(response Record) Record {
	return ensureTenantNameInRecord(response, m.tfstate, nil)
}

func (m *DeContainerRegistry) TfState() is.ResourceState {
	return m.tfstate
}

func (m *DeContainerRegistry) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.DataEngine.ContainerRegistries
}

func (m *DeContainerRegistry) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return findRecordByLocalFilter(ctx, m.API(rest), m.tfstate, m.TransformQueryParams(nil, rest), "container-registries")
}

func (m *DeContainerRegistry) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return findRecordByLocalFilter(ctx, m.API(rest), m.tfstate, m.TransformQueryParams(nil, rest), "container-registries")
}

func (m *DeContainerRegistry) LookupForCreate(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return findRecordByLocalFilter(ctx, m.API(rest), m.tfstate, m.TransformQueryParams(nil, rest), "container-registries")
}

func (m *DeContainerRegistry) DeleteResource(ctx context.Context, rest *VMSRest) error {
	return deleteRecordByLocalFilter(ctx, m.API(rest), m.tfstate, m.TransformQueryParams(nil, rest), "container-registries")
}
