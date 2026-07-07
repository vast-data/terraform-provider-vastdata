// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var CertificateSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"certificates",
	http.MethodGet,
	"certificates",
)

type Certificate struct {
	tfstate *is.TFState
}

func (m *Certificate) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Certificate{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:       CertificateSchemaRef,
			ReadOnlyFields:  []string{"id", "created", "guid", "state"},
			SensitiveFields: []string{"private_key"},
			// The API does not return PEM fields in GET responses.
			SkipRefreshAPICall: true,
			AdditionalSchemaAttributes: map[string]any{
				"pgp_public_key": rschema.StringAttribute{
					Optional:    true,
					Sensitive:   true,
					Description: "Optional PGP public key used to encrypt the private_key. When provided, private_key will be cleared from state and the encrypted value stored in encrypted_private_key.",
				},
				"encrypted_private_key": rschema.StringAttribute{
					Computed:    true,
					Description: "The PGP-encrypted private key, populated when pgp_public_key is provided.",
				},
				"validate": rschema.BoolAttribute{
					Optional:    true,
					Description: "When true, calls POST /certificates/validate_compute_cluster_certificates/ with ca_certificate, certificate, and private_key before creating the certificate. Defaults to false.",
				},
			},
		},
	)}
}

func (m *Certificate) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Certificate{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: CertificateSchemaRef,
		},
	)}
}

func (m *Certificate) TfState() *is.TFState {
	return m.tfstate
}

func (m *Certificate) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Certificates
}

// PrepareCreateResource validates the certificate material against the cluster
// before creation when validate=true.
func (m *Certificate) PrepareCreateResource(ctx context.Context, rest *VMSRest) error {
	ts := m.tfstate
	if !ts.IsKnownAndNotNull("validate") || !ts.Bool("validate") {
		return nil
	}
	body := map[string]any{
		"root_certificate":         ts.String("ca_certificate"),
		"intermediate_certificate": ts.String("certificate"),
		"intermediate_key":         ts.String("private_key"),
	}
	_, err := rest.Certificates.CertificateValidateComputeClusterCertificatesWithContext_POST(ctx, body)
	return err
}

// AfterCreateResource encrypts private_key with pgp_public_key (if provided)
// and stores the result in encrypted_private_key. private_key is always kept
// in state.
func (m *Certificate) AfterCreateResource(_ context.Context, _ *VMSRest, record Record) error {
	ts := m.tfstate
	if !ts.IsKnownAndNotNull("pgp_public_key") {
		record["encrypted_private_key"] = types.StringNull()
		return nil
	}

	privateKey := ts.String("private_key")
	if privateKey == "" {
		record["encrypted_private_key"] = types.StringNull()
		return nil
	}

	encrypted, err := helper.EncryptMessageArmored(ts.String("pgp_public_key"), privateKey)
	if err != nil {
		return err
	}
	record["encrypted_private_key"] = encrypted
	return nil
}
