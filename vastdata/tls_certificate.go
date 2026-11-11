// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/vast-data/go-vast-client/core"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/logging"
)

// POST /tlscertificates/ uses multipart/form-data which the schema generator
// cannot parse. Use GET schema for both create and read; file-upload fields
// are injected via AdditionalSchemaAttributes.
var TlsCertificateSchemaRef = is.NewSchemaReference(
	http.MethodGet,
	"tlscertificates",
	http.MethodGet,
	"tlscertificates",
)

type TlsCertificate struct {
	tfstate is.ResourceState
}

func (m *TlsCertificate) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &TlsCertificate{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			Importable:         &notImportable,
			SchemaRef:          TlsCertificateSchemaRef,
			WriteOnlyFields:    []string{"ca_certificate_file", "revocation_file"},
			WriteOnlyMD5Fields: []string{"ca_certificate_file", "revocation_file"},
			SensitiveFields:    []string{"ca_certificate_file", "revocation_file"},
			AdditionalSchemaAttributes: map[string]any{
				// File-upload fields not present in the GET response schema.
				"ca_certificate_file": rschema.StringAttribute{
					Required:    true,
					Sensitive:   true,
					WriteOnly:   true,
					Description: "CA certificate content (PEM/DER). Uploaded as a file to POST /tlscertificates/.",
				},
				"revocation_file": rschema.StringAttribute{
					Optional:    true,
					Sensitive:   true,
					WriteOnly:   true,
					Description: "Certificate Revocation List (CRL) content.",
				},
				"ca_certificate_file_md5": ContentMD5SchemaAttribute("ca_certificate_file", false),
				"revocation_file_md5":     ContentMD5SchemaAttribute("revocation_file", false),
				"tenant_id": rschema.Int64Attribute{
					Optional:    true,
					Computed:    true,
					Description: "Tenant ID to associate with the TLS certificate.",
				},
			},
		},
	)}
}

func (m *TlsCertificate) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &TlsCertificate{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: TlsCertificateSchemaRef,
		},
	)}
}

func (m *TlsCertificate) TfState() is.ResourceState {
	return m.tfstate
}

func (m *TlsCertificate) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.TlsCertificates
}

func (m *TlsCertificate) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) error {
	SetPlanContentMD5Fields(ctx, req, resp,
		"ca_certificate_file", "ca_certificate_file_md5",
		"revocation_file", "revocation_file_md5",
	)
	return nil
}

func (m *TlsCertificate) TransformResponseRecord(record Record) Record {
	DeleteRecordFields(record, "ca_certificate_file", "revocation_file")
	return record
}

// PrepareCreateResource calls POST /tlscertificates/is_operation_healthy/ to
// verify the system can perform the TLS certificate operation.
// The endpoint is called with an empty body (system-level health check only —
// the API requires multipart file upload for cert content which is not
// representable as a plain JSON param).
func (m *TlsCertificate) PrepareCreateResource(ctx context.Context, rest *VMSRest) error {
	result, err := rest.TlsCertificates.TlsCertificateIsOperationHealthyWithContext_POST(ctx, params{})
	if err != nil {
		return err
	}
	if warnings, ok := result["warnings"]; ok {
		if ws, ok := warnings.([]any); ok && len(ws) > 0 {
			for _, w := range ws {
				logging.Warn(ctx, fmt.Sprintf("TlsCertificate is_operation_healthy warning: %v", w))
			}
		}
	}
	return nil
}

func (m *TlsCertificate) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	body := ts.GetCreateParams()
	delete(body, "ca_certificate_file")
	delete(body, "revocation_file")

	// Wrap file fields as core.FileData so the session serialises them as
	// multipart parts with the correct content-type.
	if ts.IsKnownAndNotNull("ca_certificate_file") {
		body["ca_certificate_file"] = core.FileData{
			Filename:    ts.String("ca_certificate_name") + ".pem",
			Content:     []byte(ts.String("ca_certificate_file")),
			ContentType: "application/x-pem-file",
		}
	}
	if ts.IsKnownAndNotNull("revocation_file") {
		body["revocation_file"] = core.FileData{
			Filename:    ts.String("ca_certificate_name") + "-crl.pem",
			Content:     []byte(ts.String("revocation_file")),
			ContentType: "application/x-pem-file",
		}
	}

	multipartHeader := http.Header{}
	multipartHeader.Set("Content-Type", core.ContentTypeMultipartForm)

	result, err := core.RequestWithHeaders[core.Record](
		ctx, rest.TlsCertificates,
		http.MethodPost, "/tlscertificates/",
		nil, body,
		[]http.Header{multipartHeader},
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AfterUpdateResource deletes the CRL when revocation_file_md5 is cleared.
func (m *TlsCertificate) AfterUpdateResource(ctx context.Context, plan AfterUpdateResource, rest *VMSRest, record Record) error {
	planTs := plan.(*TlsCertificate).TfState()
	stateHad := m.tfstate.IsKnownAndNotNull("revocation_file_md5") && m.tfstate.String("revocation_file_md5") != ""
	planHas := planTs.IsKnownAndNotNull("revocation_file_md5") && planTs.String("revocation_file_md5") != ""
	if !stateHad || planHas {
		return nil
	}
	if exist, ok := record["revocations_exist"].(bool); ok && !exist {
		return nil
	}
	id := record.RecordID()
	err := rest.TlsCertificates.TlsCertificateCrlWithContext_DELETE(ctx, id)
	// Already removed (e.g. prior apply deleted CRL but failed to clear md5).
	return ignoreStatusCodes(err, http.StatusNotFound, http.StatusBadRequest)
}
