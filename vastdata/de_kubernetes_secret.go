// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/vast-data/go-vast-client/core"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

// Kubernetes secrets have no GET/{guid}. Identity is
// (kubernetes_cluster_vrn, namespace, secret_name). Create is upsert (POST);
// delete sends the identity as a request body on the collection.
var DeKubernetesSecretSchemaRef = is.NewDataEngineSchemaReference(
	http.MethodPost,
	"kubernetes-secrets",
	"",
	"",
)

type DeKubernetesSecret struct {
	tfstate is.ResourceState
}

func (m *DeKubernetesSecret) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &DeKubernetesSecret{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			Importable:           &notImportable,
			SchemaRef:            DeKubernetesSecretSchemaRef,
			SkipRefreshAPICall:   true,
			ExcludedSchemaFields: []string{"cursor", "limit", "order_by"},
			SensitiveFields:      []string{"entries"},
			CreateOnlyFields: []string{
				"kubernetes_cluster_vrn",
				"namespace",
				"secret_name",
				deTenantNameAttr,
			},
			PreserveUserValueFields: []string{
				"kubernetes_cluster_vrn",
				"namespace",
				"secret_name",
				"entries",
				deTenantNameAttr,
			},
			CommonModifiersMapping: map[string]string{
				"kubernetes_cluster_vrn": ModifierForceNew,
				"namespace":              ModifierForceNew,
				"secret_name":            ModifierForceNew,
				deTenantNameAttr:         ModifierForceNew,
			},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameSchemaAttribute(),
			},
		},
	)}
}

func (m *DeKubernetesSecret) TransformQueryParams(query params, rest *VMSRest) params {
	return deTransformQueryParams(query, m.tfstate, rest)
}

func (m *DeKubernetesSecret) TransformRequestBody(body params) params {
	return deTransformRequestBody(body)
}

func (m *DeKubernetesSecret) TransformResponseRecord(response Record) Record {
	return ensureTenantNameInRecord(response, m.tfstate, nil)
}

func (m *DeKubernetesSecret) TfState() is.ResourceState {
	return m.tfstate
}

func (m *DeKubernetesSecret) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.DataEngine.KubernetesSecrets
}

func (m *DeKubernetesSecret) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	api := m.API(rest)
	body := m.TransformRequestBody(m.tfstate.GetCreateParams())
	query := m.TransformQueryParams(nil, rest)
	// POST is upsert; response is often empty — echo the request as state.
	if _, err := api.CreateWithContext(ctx, query, body); err != nil {
		return nil, err
	}
	return Record(body), nil
}

func (m *DeKubernetesSecret) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	planMgr := plan.(*DeKubernetesSecret)
	api := m.API(rest)
	body := planMgr.TransformRequestBody(planMgr.tfstate.GetCreateParams())
	query := planMgr.TransformQueryParams(nil, rest)
	// Same POST upsert as create (no PATCH/{guid} API).
	if _, err := api.CreateWithContext(ctx, query, body); err != nil {
		return nil, err
	}
	m.tfstate = planMgr.tfstate
	return Record(body), nil
}

func (m *DeKubernetesSecret) DeleteResource(ctx context.Context, rest *VMSRest) error {
	ts := m.tfstate
	if ts == nil {
		return fmt.Errorf("de_kubernetes_secret: empty state")
	}
	body := core.Params{}
	for _, k := range []string{"kubernetes_cluster_vrn", "namespace", "secret_name"} {
		if !ts.IsKnownAndNotNull(k) {
			return fmt.Errorf("de_kubernetes_secret delete: missing %q", k)
		}
		body[k] = ts.String(k)
	}
	api := m.API(rest)
	query := m.TransformQueryParams(nil, rest)
	_, err := api.DeleteWithContext(ctx, nil, query, body)
	return ignoreResourceGone(err)
}
