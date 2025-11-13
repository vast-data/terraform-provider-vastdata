// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var ClusterEkmSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"clusters/{id}/add_ekm",
	"",
	"",
)

type ClusterEkm struct {
	tfstate *is.TFState
}

func (m *ClusterEkm) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &ClusterEkm{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			Importable: &notImportable,
			SchemaRef:  ClusterEkmSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"cluster_id": rschema.Int64Attribute{
					Required:    true,
					Description: "Cluster ID to which the EKM should be added",
				},
			},
		},
	)}
}

func (m *ClusterEkm) TfState() *is.TFState {
	return m.tfstate
}

func (m *ClusterEkm) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Clusters
}

func (m *ClusterEkm) ReadResource(_ context.Context, _ *VMSRest) (DisplayableRecord, error) {
	return nil, nil
}

func (m *ClusterEkm) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	clusterId := m.tfstate.Int64("cluster_id")
	params := m.tfstate.GetCreateParams()
	params.Without("cluster_id") // Remove cluster_id from params, as it's part of the URL

	return addClusterEkm(ctx, clusterId, params, rest)
}

func (m *ClusterEkm) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	clusterId := m.tfstate.Int64("cluster_id")
	planTs := plan.(*ClusterEkm).TfState()
	params := planTs.GetChangedParams(m.tfstate)
	params.Without("cluster_id") // Remove cluster_id from params, as it's part of the URL

	return addClusterEkm(ctx, clusterId, params, rest)
}

func (m *ClusterEkm) DeleteResource(_ context.Context, _ *VMSRest) error {
	// No-op: ClusterEkm cannot be deleted - it's a one-time operation
	return nil
}

// addClusterEkm adds EKM configuration to a specific cluster with the given parameters.
// This is used by both CreateResource and UpdateResource for ClusterEkm.
func addClusterEkm(ctx context.Context, clusterId int64, params map[string]any, rest *VMSRest) (DisplayableRecord, error) {
	if clusterId == 0 {
		return nil, fmt.Errorf("failed to get cluster ID: cluster ID is empty")
	}

	err := rest.Clusters.ClusterAddEkmWithContext_POST(ctx, clusterId, params)
	return nil, err
}
