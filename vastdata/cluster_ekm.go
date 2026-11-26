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
	tfstate is.ResourceState
}

func (m *ClusterEkm) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &ClusterEkm{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			Importable: &notImportable,
			SchemaRef:  ClusterEkmSchemaRef,
			WriteOnlyFields: []string{
				"ekm_private_key", "ekm_certificate", "ekm_ca_certificate",
			},
			WriteOnlyMD5Fields: []string{
				"ekm_private_key", "ekm_certificate", "ekm_ca_certificate",
			},
			SensitiveFields: []string{
				"ekm_private_key", "ekm_certificate", "ekm_ca_certificate",
			},
			AdditionalSchemaAttributes: map[string]any{
				"cluster_id": rschema.Int64Attribute{
					Required:    true,
					Description: "Cluster ID to which the EKM should be added",
				},
				"ekm_private_key_md5":    ContentMD5SchemaAttribute("ekm_private_key", false),
				"ekm_certificate_md5":    ContentMD5SchemaAttribute("ekm_certificate", false),
				"ekm_ca_certificate_md5": ContentMD5SchemaAttribute("ekm_ca_certificate", false),
			},
		},
	)}
}

func (m *ClusterEkm) TfState() is.ResourceState {
	return m.tfstate
}

func (m *ClusterEkm) TransformResponseRecord(record Record) Record {
	DeleteRecordFields(record, "ekm_private_key", "ekm_certificate", "ekm_ca_certificate")
	return record
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
	// POST /clusters/{id}/set_certificates/ expects the full EKM configuration
	// rather than a partial diff: omitting encryption_type causes VMS to default
	// to INTERNAL and reject certificate fields with HTTP 400 (TERF-322).
	params := planTs.GetCreateParams()
	params.Without("cluster_id") // Remove cluster_id from params, as it's part of the URL

	return setClusterEkmCertificates(ctx, clusterId, params, rest)
}

func (m *ClusterEkm) DeleteResource(_ context.Context, _ *VMSRest) error {
	// No-op: ClusterEkm cannot be deleted - it's a one-time operation
	return nil
}

// addClusterEkm adds EKM configuration to a specific cluster via POST /clusters/{id}/add_ekm/.
func addClusterEkm(ctx context.Context, clusterId int64, params map[string]any, rest *VMSRest) (DisplayableRecord, error) {
	if clusterId == 0 {
		return nil, fmt.Errorf("failed to get cluster ID: cluster ID is empty")
	}

	err := rest.Clusters.ClusterAddEkmWithContext_POST(ctx, clusterId, params)
	return nil, err
}

// setClusterEkmCertificates updates EKM certificates and configuration on a cluster
// via POST /clusters/{id}/set_certificates/.
func setClusterEkmCertificates(ctx context.Context, clusterId int64, params map[string]any, rest *VMSRest) (DisplayableRecord, error) {
	if clusterId == 0 {
		return nil, fmt.Errorf("failed to get cluster ID: cluster ID is empty")
	}

	err := rest.Clusters.ClusterSetCertificatesWithContext_POST(ctx, clusterId, params)
	return nil, err
}
