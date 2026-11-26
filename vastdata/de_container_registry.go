// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var DeContainerRegistrySchemaRef = is.NewDataEngineSchemaReference(
	http.MethodPost,
	"container-registries",
	http.MethodGet,
	"container-registries",
)

const (
	deCRPrimaryKubernetesCluster     = "primary_kubernetes_cluster"
	deCRAdditionalKubernetesClusters = "additional_kubernetes_clusters"
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
			NotRequiredSchemaFields: []string{deKubernetesClusterVrnAttr},
			PreserveUserValueFields: []string{deTenantNameAttr},
			CommonModifiersMapping: map[string]string{
				"name":           ModifierForceNew,
				deTenantNameAttr: ModifierForceNew,
			},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr:                 deTenantNameSchemaAttribute(),
				deCRPrimaryKubernetesCluster:     deCRPrimaryKubernetesClusterSchemaAttribute(),
				deCRAdditionalKubernetesClusters: deCRAdditionalKubernetesClustersSchemaAttribute(),
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
	body = deTransformRequestBody(body)
	return deCRStripKubernetesClusterNamesFromBody(body)
}

func (m *DeContainerRegistry) TransformResponseRecord(response Record) Record {
	response = ensureTenantNameInRecord(response, m.tfstate, nil)
	return deCRPreserveKubernetesClusterNamesInRecord(response, m.tfstate)
}

func (m *DeContainerRegistry) TfState() is.ResourceState {
	return m.tfstate
}

func (m *DeContainerRegistry) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.DataEngine.ContainerRegistries
}

func (m *DeContainerRegistry) PrepareCreateResource(ctx context.Context, rest *VMSRest) error {
	if err := deCRValidateKubernetesClusterRefsConfig(m.tfstate); err != nil {
		return err
	}
	return deCRResolveKubernetesClusterRefsInState(ctx, m.tfstate, rest)
}

func (m *DeContainerRegistry) PrepareUpdateResource(ctx context.Context, plan PrepareUpdateResource, rest *VMSRest) error {
	ts := plan.(*DeContainerRegistry).tfstate
	if err := deCRValidateKubernetesClusterRefsConfig(ts); err != nil {
		return err
	}
	return deCRResolveKubernetesClusterRefsInState(ctx, ts, rest)
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

func deCRKubernetesClusterRefAttributes() map[string]rschema.Attribute {
	return map[string]rschema.Attribute{
		deKubernetesClusterVrnAttr:  deKubernetesClusterVrnSchemaAttribute(),
		deKubernetesClusterNameAttr: deKubernetesClusterNameSchemaAttribute(),
		"namespace": rschema.StringAttribute{
			Optional:    true,
			Description: "Kubernetes namespace on the cluster.",
		},
	}
}

func deCRPrimaryKubernetesClusterSchemaAttribute() rschema.SingleNestedAttribute {
	return rschema.SingleNestedAttribute{
		Required:    true,
		Attributes:  deCRKubernetesClusterRefAttributes(),
		Description: "Primary Kubernetes cluster and namespace for the container registry.",
	}
}

func deCRAdditionalKubernetesClustersSchemaAttribute() rschema.SetNestedAttribute {
	return rschema.SetNestedAttribute{
		Optional:    true,
		Description: "Additional Kubernetes clusters and namespaces where the container registry is replicated.",
		NestedObject: rschema.NestedAttributeObject{
			Attributes: deCRKubernetesClusterRefAttributes(),
		},
	}
}

func deCRResolveKubernetesClusterRefsInState(
	ctx context.Context,
	ts is.ResourceState,
	rest *VMSRest,
) error {
	if ts == nil || rest == nil {
		return fmt.Errorf("kubernetes cluster resolution: empty state or client")
	}
	listQuery := resolveTenantQuery(ts, rest)
	api := rest.DataEngine.KubernetesClusters

	if ts.IsKnownAndNotNull(deCRPrimaryKubernetesCluster) {
		primary := ts.ToMap(deCRPrimaryKubernetesCluster)
		if err := resolveKubernetesClusterNameOrVrnInMap(ctx, api, listQuery, primary, deCRPrimaryKubernetesCluster); err != nil {
			return err
		}
		ts.Set(deCRPrimaryKubernetesCluster, primary)
	}

	if ts.HasAttribute(deCRAdditionalKubernetesClusters) &&
		ts.IsKnownAndNotNull(deCRAdditionalKubernetesClusters) {
		items := ts.ToSlice(deCRAdditionalKubernetesClusters)
		for i, item := range items {
			m, ok := item.(map[string]any)
			if !ok {
				return fmt.Errorf("%s[%d]: expected object", deCRAdditionalKubernetesClusters, i)
			}
			if err := resolveKubernetesClusterNameOrVrnInMap(
				ctx, api, listQuery, m,
				fmt.Sprintf("%s[%d]", deCRAdditionalKubernetesClusters, i),
			); err != nil {
				return err
			}
			items[i] = m
		}
		ts.Set(deCRAdditionalKubernetesClusters, items)
	}
	return nil
}

func deCRStripKubernetesClusterNamesFromBody(body params) params {
	if body == nil {
		return body
	}
	if primary, ok := body[deCRPrimaryKubernetesCluster].(map[string]any); ok {
		delete(primary, deKubernetesClusterNameAttr)
	}
	switch add := body[deCRAdditionalKubernetesClusters].(type) {
	case []any:
		for _, item := range add {
			if m, ok := item.(map[string]any); ok {
				delete(m, deKubernetesClusterNameAttr)
			}
		}
	case []map[string]any:
		for _, m := range add {
			delete(m, deKubernetesClusterNameAttr)
		}
	}
	return body
}

func deCRPreserveKubernetesClusterNamesInRecord(rec Record, ts is.ResourceState) Record {
	if rec == nil || ts == nil {
		return rec
	}
	out := Record{}
	for k, v := range rec {
		out[k] = v
	}

	if ts.IsKnownAndNotNull(deCRPrimaryKubernetesCluster) {
		statePrimary := ts.ToMap(deCRPrimaryKubernetesCluster)
		if name, _ := statePrimary[deKubernetesClusterNameAttr].(string); strings.TrimSpace(name) != "" {
			apiPrimary, _ := out[deCRPrimaryKubernetesCluster].(map[string]any)
			if apiPrimary == nil {
				apiPrimary = map[string]any{}
			} else {
				copied := map[string]any{}
				for k, v := range apiPrimary {
					copied[k] = v
				}
				apiPrimary = copied
			}
			apiPrimary[deKubernetesClusterNameAttr] = strings.TrimSpace(name)
			out[deCRPrimaryKubernetesCluster] = apiPrimary
		}
	}

	if !ts.HasAttribute(deCRAdditionalKubernetesClusters) ||
		!ts.IsKnownAndNotNull(deCRAdditionalKubernetesClusters) {
		return out
	}
	nameByVrn := map[string]string{}
	for _, item := range ts.ToSlice(deCRAdditionalKubernetesClusters) {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		name, _ := m[deKubernetesClusterNameAttr].(string)
		vrn, _ := m[deKubernetesClusterVrnAttr].(string)
		name = strings.TrimSpace(name)
		vrn = strings.TrimSpace(vrn)
		if name != "" && vrn != "" {
			nameByVrn[vrn] = name
		}
	}
	if len(nameByVrn) == 0 {
		return out
	}

	var apiAdd []any
	switch v := out[deCRAdditionalKubernetesClusters].(type) {
	case []any:
		apiAdd = v
	case RecordSet:
		for _, r := range v {
			apiAdd = append(apiAdd, map[string]any(r))
		}
	default:
		return out
	}
	for i, item := range apiAdd {
		m, ok := item.(map[string]any)
		if !ok {
			if r, ok := item.(Record); ok {
				m = map[string]any(r)
			} else {
				continue
			}
		}
		copied := map[string]any{}
		for k, v := range m {
			copied[k] = v
		}
		vrn := strings.TrimSpace(fmt.Sprint(copied[deKubernetesClusterVrnAttr]))
		if name, ok := nameByVrn[vrn]; ok {
			copied[deKubernetesClusterNameAttr] = name
		}
		apiAdd[i] = copied
	}
	out[deCRAdditionalKubernetesClusters] = apiAdd
	return out
}

func deCRValidateKubernetesClusterRefsConfig(ts is.ResourceState) error {
	if ts == nil {
		return nil
	}
	if v, ok := ts.RawValues()[deCRPrimaryKubernetesCluster]; ok {
		if v.IsUnknown() {
			return nil
		}
		if !v.IsNull() {
			obj, ok := v.(types.Object)
			if !ok {
				return fmt.Errorf("%s: expected object", deCRPrimaryKubernetesCluster)
			}
			if err := validateKubernetesClusterNameOrVrnAttrs(deCRPrimaryKubernetesCluster, obj.Attributes()); err != nil {
				return err
			}
		}
	}
	if v, ok := ts.RawValues()[deCRAdditionalKubernetesClusters]; ok {
		if v.IsUnknown() || v.IsNull() {
			return nil
		}
		set, ok := v.(types.Set)
		if !ok {
			return fmt.Errorf("%s: expected set", deCRAdditionalKubernetesClusters)
		}
		for i, elem := range set.Elements() {
			if elem.IsUnknown() || elem.IsNull() {
				continue
			}
			obj, ok := elem.(types.Object)
			if !ok {
				return fmt.Errorf("%s[%d]: expected object", deCRAdditionalKubernetesClusters, i)
			}
			if err := validateKubernetesClusterNameOrVrnAttrs(
				fmt.Sprintf("%s[%d]", deCRAdditionalKubernetesClusters, i),
				obj.Attributes(),
			); err != nil {
				return err
			}
		}
	}
	return nil
}
