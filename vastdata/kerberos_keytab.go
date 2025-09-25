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

var KerberosKeytabSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"kerberos/{id}/keytab",
	"",
	"",
)

type KerberosKeytab struct {
	tfstate *is.TFState
}

func (m *KerberosKeytab) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &KerberosKeytab{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			Importable: &notImportable,
			SchemaRef:  KerberosKeytabSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"kerberos_id": rschema.Int64Attribute{
					Required:    true,
					Description: "ID of the Kerberos configuration",
				},
				"keytab_file": rschema.StringAttribute{
					Optional:    true,
					Description: "Keytab file content",
				},
				"filename": rschema.StringAttribute{
					Optional:    true,
					Description: "Custom filename for the keytab file",
				},
			},
			CommonModifiersMapping: map[string]string{
				"kerberos_id":    schema_generation.ModifierForceNew,
				"admin_username": schema_generation.ModifierForceNew,
				"admin_password": schema_generation.ModifierForceNew,
				"keytab_file":    schema_generation.ModifierForceNew,
				"filename":       schema_generation.ModifierForceNew,
			},
		},
	)}
}

func (m *KerberosKeytab) TfState() *is.TFState {
	return m.tfstate
}

func (m *KerberosKeytab) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Kerberos
}

func (m *KerberosKeytab) ReadResource(_ context.Context, _ *VMSRest) (DisplayableRecord, error) {
	return nil, nil
}

func (m *KerberosKeytab) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate

	kerberosId := ts.Int64("kerberos_id")
	adminUsername := ts.String("admin_username")
	adminPassword := ts.String("admin_password")

	// Prepare parameters for keytab generation
	params := params{
		"admin_username": adminUsername,
		"admin_password": adminPassword,
	}

	// Generate the keytab
	_, err := rest.Kerberos.GenerateKeytabWithContext(ctx, kerberosId, params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate keytab: %w", err)
	}

	// Check if keytab_file is provided for upload
	if ts.IsKnownAndNotNull("keytab_file") {
		keytabFileData := ts.String("keytab_file")
		filename := "keytab"

		// Use custom filename if provided
		if ts.IsKnownAndNotNull("filename") {
			filename = ts.String("filename")
		}

		// Upload the keytab file
		uploadRecord, err := rest.Kerberos.UploadKeytabWithContext(ctx, kerberosId, []byte(keytabFileData), filename)
		if err != nil {
			return nil, fmt.Errorf("failed to upload keytab: %w", err)
		}

		// Return the upload record if available, otherwise the generation record
		if uploadRecord != nil {
			return uploadRecord, nil
		}
	}

	return nil, nil
}

func (m *KerberosKeytab) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	// With force_new modifiers, Terraform will handle replacements automatically
	// This method should not be called for updates since all fields have RequiresReplace()
	// But we'll keep it as a safety net in case it's called
	return nil, fmt.Errorf("kerberos keytab operations should be replaced, not updated")
}

func (m *KerberosKeytab) DeleteResource(ctx context.Context, rest *VMSRest) error {
	// No-op: KerberosKeytab cannot be deleted - it's a one-time operation
	return nil
}
