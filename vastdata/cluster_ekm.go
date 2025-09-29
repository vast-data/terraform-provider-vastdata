// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/schema_generation"
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
			CommonModifiersMapping: map[string]string{
				"cluster_id":            schema_generation.ModifierForceNew,
				"ekm_auth_domain":       schema_generation.ModifierForceNew,
				"ekm_bypass_validation": schema_generation.ModifierForceNew,
				"ekm_ca_certificate":    schema_generation.ModifierForceNew,
				"ekm_certificate":       schema_generation.ModifierForceNew,
				"ekm_domain":            schema_generation.ModifierForceNew,
				"ekm_private_key":       schema_generation.ModifierForceNew,
				"ekm_proxy_address":     schema_generation.ModifierForceNew,
				"ekm_servers":           schema_generation.ModifierForceNew,
				"encryption_type":       schema_generation.ModifierForceNew,
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
	ts := m.tfstate
	clusterId := ts.Int64("cluster_id")
	createParams := ts.GetCreateParams()
	createParams.Without("cluster_id") // Remove cluster_id from the body since it's in the URL path
	_, err := rest.Clusters.AddEkmWithContext(ctx, clusterId, createParams)
	return nil, err
}

func (m *ClusterEkm) UpdateResource(_ context.Context, _ UpdateResource, _ *VMSRest) (DisplayableRecord, error) {
	// With force_new modifiers, Terraform will handle replacements automatically
	// This method should not be called for updates since all fields have RequiresReplace()
	// But we'll keep it as a safety net in case it's called
	return nil, fmt.Errorf("cluster EKM operations should be replaced, not updated")
}

func (m *ClusterEkm) DeleteResource(_ context.Context, _ *VMSRest) error {
	// No-op: KerberosKeytab cannot be deleted - it's a one-time operation
	return nil
}
