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
	kerberosId := m.tfstate.Int64("kerberos_id")
	// For create, all params are considered "changed"
	allParams := m.tfstate.GetCreateParams()

	return processKerberosKeytab(ctx, kerberosId, m.tfstate, allParams, rest)
}

func (m *KerberosKeytab) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	kerberosId := m.tfstate.Int64("kerberos_id")
	planTs := plan.(*KerberosKeytab).TfState()
	changedParams := planTs.GetChangedParams(m.tfstate)

	return processKerberosKeytab(ctx, kerberosId, planTs, changedParams, rest)
}

func (m *KerberosKeytab) DeleteResource(ctx context.Context, rest *VMSRest) error {
	// No-op: KerberosKeytab cannot be deleted - it's a one-time operation
	return nil
}

// processKerberosKeytab handles the complex keytab generation and optional upload logic.
// This is used by both CreateResource and UpdateResource for KerberosKeytab.
func processKerberosKeytab(ctx context.Context, kerberosId int64, tfstate *is.TFState, changedParams map[string]any, rest *VMSRest) (DisplayableRecord, error) {
	if kerberosId == 0 {
		return nil, fmt.Errorf("failed to get kerberos ID: kerberos ID is empty")
	}

	// Check if we need to generate keytab (admin credentials changed or this is a create operation)
	needsGeneration := false
	if _, hasUsername := changedParams["admin_username"]; hasUsername {
		needsGeneration = true
	}
	if _, hasPassword := changedParams["admin_password"]; hasPassword {
		needsGeneration = true
	}
	// For create operations, changedParams will contain all params, so we always generate
	if len(changedParams) > 0 && (changedParams["admin_username"] != nil || changedParams["admin_password"] != nil) {
		needsGeneration = true
	}

	if needsGeneration {
		// Get both admin credentials from the current tfstate
		adminUsername := tfstate.String("admin_username")
		adminPassword := tfstate.String("admin_password")

		if adminUsername == "" || adminPassword == "" {
			return nil, fmt.Errorf("admin_username and admin_password are required for keytab generation")
		}

		_, err := rest.Kerberos.KerberosKeytabWithContext_POST(ctx, kerberosId, params{
			"admin_password": adminPassword,
			"admin_username": adminUsername,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to generate keytab: %w", err)
		}
	}

	// Check if keytab_file is provided for upload
	if tfstate.IsKnownAndNotNull("keytab_file") {
		keytabFileData := tfstate.String("keytab_file")
		filename := "keytab"

		// Use custom filename if provided
		if tfstate.IsKnownAndNotNull("filename") {
			filename = tfstate.String("filename")
		}

		// Upload the keytab file
		uploadRecord, err := rest.Kerberos.KerberosKeytabWithContext_PUT(ctx, kerberosId, []byte(keytabFileData), filename)
		if err != nil {
			return nil, fmt.Errorf("failed to upload keytab: %w", err)
		}

		// Return the upload record if available
		if uploadRecord != nil {
			return uploadRecord, nil
		}
	}

	// Return nil for successful generation without upload
	return nil, nil
}
