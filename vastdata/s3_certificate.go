// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var s3CertificateHints = endpointCertificateHints(
	"cluster_id",
	"ID of the cluster that hosts the S3 service (typically 1).",
	"Installs the S3 HTTPS server certificate on a cluster via PATCH /clusters/{id}/. "+
		"The PEM certificate and private key are write-only and are never stored in Terraform state.",
	"Manages the **S3** HTTPS endpoint certificate.\n\n"+
		"Uses `PATCH /clusters/{id}/` (`s3_certificate`, `s3_private_key`). Available on VAST 5.4+.\n\n"+
		"Destroy is a no-op: removing this resource does not clear the S3 certificate.\n\n"+
		"This resource is not importable.",
)

type S3Certificate struct {
	tfstate *is.TFState
}

func (m *S3Certificate) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &S3Certificate{tfstate: newResourceTFState(raw, schema, s3CertificateHints)}
}

func (m *S3Certificate) TfState() *is.TFState {
	return m.tfstate
}

func (m *S3Certificate) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Clusters
}

func (m *S3Certificate) ops() *endpointCertificateOps {
	return &endpointCertificateOps{
		parentAttr:  "cluster_id",
		certAPI:     "s3_certificate",
		certInfoAPI: "s3_certificate_info",
		apply: func(ctx context.Context, rest *VMSRest, parentID int64, cert, key string) (Record, error) {
			return rest.Clusters.UpdateWithContext(ctx, parentID, params{
				"s3_certificate": cert,
				"s3_private_key": key,
			})
		},
		get: func(ctx context.Context, rest *VMSRest, parentID int64) (Record, error) {
			return rest.Clusters.GetByIdWithContext(ctx, parentID)
		},
	}
}

func (m *S3Certificate) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.ops().Read(ctx, rest, m.tfstate)
}

func (m *S3Certificate) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.ops().Apply(ctx, rest, m.tfstate)
}

func (m *S3Certificate) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	return m.ops().Apply(ctx, rest, plan.(*S3Certificate).tfstate)
}

func (m *S3Certificate) DeleteResource(_ context.Context, _ *VMSRest) error {
	return nil
}

func (m *S3Certificate) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) error {
	return m.ops().ModifyPlan(ctx, req, resp)
}
