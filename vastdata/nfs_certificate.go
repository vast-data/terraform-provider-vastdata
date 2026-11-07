// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var nfsCertificateHints = endpointCertificateHints(
	"cluster_id",
	"ID of the cluster that hosts NFSv4 (typically 1).",
	"Installs the NFSv4 RPC-with-TLS certificate on a cluster via PATCH /clusters/{id}/. "+
		"The PEM certificate and private key are write-only and are never stored in Terraform state.",
	"Manages the **NFS** (NFSv4 RPC-with-TLS) server certificate.\n\n"+
		"Uses `PATCH /clusters/{id}/` (`nfs4_certificate`, `nfs4_private_key`). Available on VAST 5.4+.\n\n"+
		"Destroy is a no-op: removing this resource does not clear the NFS certificate.\n\n"+
		"This resource is not importable.",
)

type NfsCertificate struct {
	tfstate is.ResourceState
}

func (m *NfsCertificate) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &NfsCertificate{tfstate: newResourceTFState(raw, schema, nfsCertificateHints)}
}

func (m *NfsCertificate) TfState() is.ResourceState {
	return m.tfstate
}

func (m *NfsCertificate) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Clusters
}

func (m *NfsCertificate) ops() *endpointCertificateOps {
	return &endpointCertificateOps{
		parentAttr:  "cluster_id",
		certAPI:     "nfs4_certificate",
		certInfoAPI: "nfs4_certificate_info",
		apply: func(ctx context.Context, rest *VMSRest, parentID int64, cert, key string) (Record, error) {
			return rest.Clusters.UpdateWithContext(ctx, parentID, params{
				"nfs4_certificate": cert,
				"nfs4_private_key": key,
			})
		},
		get: func(ctx context.Context, rest *VMSRest, parentID int64) (Record, error) {
			return rest.Clusters.GetByIdWithContext(ctx, parentID)
		},
	}
}

func (m *NfsCertificate) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.ops().Read(ctx, rest, m.tfstate)
}

func (m *NfsCertificate) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.ops().Apply(ctx, rest, m.tfstate)
}

func (m *NfsCertificate) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	return m.ops().Apply(ctx, rest, plan.(*NfsCertificate).tfstate)
}

func (m *NfsCertificate) DeleteResource(_ context.Context, _ *VMSRest) error {
	return nil
}

func (m *NfsCertificate) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) error {
	return m.ops().ModifyPlan(ctx, req, resp)
}
