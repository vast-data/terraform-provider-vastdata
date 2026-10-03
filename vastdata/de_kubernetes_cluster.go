// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var DeKubernetesClusterSchemaRef = is.NewDataEngineSchemaReference(
	http.MethodPost,
	"kubernetes-clusters",
	http.MethodGet,
	"kubernetes-clusters",
)

type DeKubernetesCluster struct {
	tfstate is.ResourceState
}

func (m *DeKubernetesCluster) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &DeKubernetesCluster{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:            DeKubernetesClusterSchemaRef,
			PathIDField:          "guid",
			ExcludedSchemaFields: []string{"cursor", "limit", "order_by"},
			// List ignores name/vrn/guid; lookup is List + local filter.
			SearchableFields:        []string{"guid", "vrn", "name"},
			CreateOnlyFields:        []string{deTenantNameAttr},
			PreserveUserValueFields: []string{deTenantNameAttr},
			CommonModifiersMapping: map[string]string{
				deTenantNameAttr: ModifierForceNew,
			},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameSchemaAttribute(),
			},
		},
	)}
}

func (m *DeKubernetesCluster) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &DeKubernetesCluster{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:            DeKubernetesClusterSchemaRef,
			PathIDField:          "guid",
			ExcludedSchemaFields: []string{"cursor", "limit", "order_by"},
			SearchableFields:     []string{"guid", "vrn", "name"},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameSchemaAttribute(),
			},
		},
	)}
}

func (m *DeKubernetesCluster) TransformQueryParams(query params, rest *VMSRest) params {
	return deTransformQueryParams(query, m.tfstate, rest)
}

func (m *DeKubernetesCluster) TransformRequestBody(body params) params {
	return deTransformRequestBody(body)
}

func (m *DeKubernetesCluster) TransformResponseRecord(response Record) Record {
	return ensureTenantNameInRecord(response, m.tfstate, nil)
}

func (m *DeKubernetesCluster) TfState() is.ResourceState {
	return m.tfstate
}

func (m *DeKubernetesCluster) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.DataEngine.KubernetesClusters
}

func (m *DeKubernetesCluster) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return findRecordByLocalFilter(ctx, m.API(rest), m.tfstate, m.TransformQueryParams(nil, rest), "kubernetes-clusters")
}

func (m *DeKubernetesCluster) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return findRecordByLocalFilter(ctx, m.API(rest), m.tfstate, m.TransformQueryParams(nil, rest), "kubernetes-clusters")
}

func (m *DeKubernetesCluster) LookupForCreate(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return findRecordByLocalFilter(ctx, m.API(rest), m.tfstate, m.TransformQueryParams(nil, rest), "kubernetes-clusters")
}

func (m *DeKubernetesCluster) DeleteResource(ctx context.Context, rest *VMSRest) error {
	return deleteRecordByLocalFilter(ctx, m.API(rest), m.tfstate, m.TransformQueryParams(nil, rest), "kubernetes-clusters")
}
