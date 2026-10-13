// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var vmsCertificateHints = endpointCertificateHints(
	"vms_id",
	"ID of the VMS object (typically 1).",
	"Installs the VMS web/API TLS certificate via PATCH /vms/{id}/set_certificate/. "+
		"The PEM certificate and private key are write-only and are never stored in Terraform state.",
	"Manages the **VMS** (web UI / API) TLS certificate.\n\n"+
		"Uses `PATCH /vms/{id}/set_certificate/` (`ssl_certificate`, `ssl_keyfile`). "+
		"Destroy is a no-op: removing this resource does not reset the VMS certificate.\n\n"+
		"This resource is not importable.",
)

type VmsCertificate struct {
	tfstate *is.TFState
}

func (m *VmsCertificate) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &VmsCertificate{tfstate: newResourceTFState(raw, schema, vmsCertificateHints)}
}

func (m *VmsCertificate) TfState() *is.TFState {
	return m.tfstate
}

func (m *VmsCertificate) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Vms
}

func (m *VmsCertificate) ops() *endpointCertificateOps {
	return &endpointCertificateOps{
		parentAttr:  "vms_id",
		certAPI:     "ssl_certificate",
		certInfoAPI: "ssl_certificate_info",
		apply: func(ctx context.Context, rest *VMSRest, parentID int64, cert, key string) (Record, error) {
			return rest.Vms.VmsSetCertificateWithContext_PATCH(ctx, parentID, params{
				"ssl_certificate": cert,
				"ssl_keyfile":     key,
			})
		},
		get: func(ctx context.Context, rest *VMSRest, parentID int64) (Record, error) {
			return rest.Vms.GetByIdWithContext(ctx, parentID)
		},
	}
}

func (m *VmsCertificate) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.ops().Read(ctx, rest, m.tfstate)
}

func (m *VmsCertificate) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.ops().Apply(ctx, rest, m.tfstate)
}

func (m *VmsCertificate) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	return m.ops().Apply(ctx, rest, plan.(*VmsCertificate).tfstate)
}

func (m *VmsCertificate) DeleteResource(_ context.Context, _ *VMSRest) error {
	return nil
}

func (m *VmsCertificate) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) error {
	return m.ops().ModifyPlan(ctx, req, resp)
}
