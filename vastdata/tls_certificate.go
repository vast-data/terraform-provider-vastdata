// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vast-data/go-vast-client/core"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
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
	tfstate *is.TFState
}

func (m *TlsCertificate) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &TlsCertificate{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			Importable: &notImportable,
			SchemaRef:  TlsCertificateSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				// File-upload fields not present in the GET response schema.
				"ca_certificate_file": rschema.StringAttribute{
					Required:    true,
					Sensitive:   true,
					Description: "CA certificate content (PEM/DER). Uploaded as a file to POST /tlscertificates/.",
				},
				"revocation_file": rschema.StringAttribute{
					Optional:    true,
					Sensitive:   true,
					Description: "Certificate Revocation List (CRL) content. When cleared, the CRL is deleted via DELETE /tlscertificates/{id}/crl/.",
				},
				"tenant_id": rschema.Int64Attribute{
					Optional:    true,
					Computed:    true,
					Description: "Tenant ID to associate with the TLS certificate.",
				},
			},
			SensitiveFields:         []string{"ca_certificate_file", "revocation_file"},
			PreserveUserValueFields: []string{"ca_certificate_file", "revocation_file"},
		},
	)}
}

func (m *TlsCertificate) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &TlsCertificate{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: TlsCertificateSchemaRef,
		},
	)}
}

func (m *TlsCertificate) TfState() *is.TFState {
	return m.tfstate
}

func (m *TlsCertificate) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.TlsCertificates
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
				tflog.Warn(ctx, fmt.Sprintf("TlsCertificate is_operation_healthy warning: %v", w))
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

// AfterUpdateResource deletes the CRL when revocation_file is cleared.
func (m *TlsCertificate) AfterUpdateResource(ctx context.Context, plan AfterUpdateResource, rest *VMSRest, record Record) error {
	planTs := plan.(*TlsCertificate).TfState()

	// If revocation_file was previously set and is now cleared, delete the CRL.
	if m.tfstate.IsKnownAndNotNull("revocation_file") && !planTs.IsKnownAndNotNull("revocation_file") {
		id := record.RecordID()
		return rest.TlsCertificates.TlsCertificateCrlWithContext_DELETE(ctx, id)
	}
	return nil
}
