// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var SamlConfigSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"vms/{id}/saml_config",
	http.MethodGet,
	"vms/{id}/saml_config",
)

type SamlConfig struct {
	tfstate *is.TFState
}

func (m *SamlConfig) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &SamlConfig{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: SamlConfigSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"vms_id": rschema.Int64Attribute{
					Required:    true,
					Description: "Unique ID of the VMS.",
				},
				"idp_name": rschema.StringAttribute{
					Required:    true,
					Description: "SAML IDP name.",
				},
			},
		},
	)}
}

func (m *SamlConfig) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &SamlConfig{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: SamlConfigSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"vms_id": dschema.Int64Attribute{
					Required:    true,
					Description: "Unique ID of the VMS.",
				},
				"idp_name": dschema.StringAttribute{
					Required:    true,
					Description: "SAML IDP name.",
				},
			},
		},
	)}
}

func (m *SamlConfig) TfState() *is.TFState {
	return m.tfstate
}

func (m *SamlConfig) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.SamlConfigs
}

func (m *SamlConfig) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	vmsId := m.tfstate.Int64("vms_id")
	idpName := m.tfstate.String("idp_name")
	record, err := rest.SamlConfigs.GetConfigWithContext(ctx, vmsId, idpName)
	if err != nil {
		return nil, err
	}

	// Backend response flattens SAML config across nested keys. Rebuild saml_settings map for TF state.
	samlSettings := map[string]any{}

	// Extract all fields from sp_settings
	if sp, ok := record["sp_settings"].(map[string]any); ok {
		for k, v := range sp {
			samlSettings[k] = v
		}
	}

	// Extract idp_entityid from idp map keyed by entityid
	if idp, ok := record["idp"].(map[string]any); ok {
		for entity := range idp {
			samlSettings["idp_entityid"] = entity
			break
		}
	}

	// Extract idp_metadata_url from metadata.remote[0].url
	if md, ok := record["metadata"].(map[string]any); ok {
		if remote, ok := md["remote"].([]any); ok && len(remote) > 0 {
			if first, ok := remote[0].(map[string]any); ok {
				if url, ok := first["url"]; ok {
					samlSettings["idp_metadata_url"] = url
				}
			}
		}
	}

	if len(samlSettings) > 0 {
		record["saml_settings"] = samlSettings
	}

	return record, nil
}

func (m *SamlConfig) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.ReadDatasource(ctx, rest)
}

func (m *SamlConfig) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	return ensureSamlConfigUpdatedWith(ctx, ts, ts, rest)
}

func (m *SamlConfig) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	stateTs := m.tfstate
	planTs := plan.(*SamlConfig).TfState()
	if err := ensureNotChanged(stateTs, planTs, "idp_name", "vms_id"); err != nil {
		return nil, err
	}
	return ensureSamlConfigUpdatedWith(ctx, stateTs, planTs, rest)
}

func (m *SamlConfig) DeleteResource(ctx context.Context, rest *VMSRest) error {
	vmsId := m.tfstate.Int64("vms_id")
	idpName := m.tfstate.String("idp_name")
	_, err := rest.SamlConfigs.DeleteConfigWithContext(ctx, vmsId, idpName)
	return ignoreStatusCodes(err, http.StatusNotFound)
}

// ensureSamlConfigUpdatedWith verifies if the given SamlConfig (looked up by state)
// needs to be updated with new fields and performs the update if necessary.
//
// This is used in both CreateResource and UpdateResource for SamlConfig.
func ensureSamlConfigUpdatedWith(ctx context.Context, stateTs, fieldsTs *is.TFState, rest *VMSRest) (DisplayableRecord, error) {
	vmsId := stateTs.Int64("vms_id")
	if vmsId == 0 {
		return nil, fmt.Errorf("failed to get VMS ID: VMS ID cannot be 0")
	}

	idpName := stateTs.String("idp_name")
	if idpName == "" {
		return nil, fmt.Errorf("failed to get IDP name: IDP name is empty")
	}

	// Create params with the SAML configuration fields
	if data, ok := fieldsTs.SetIfAvailable(
		"saml_settings",
	); ok {
		// Use the custom API method to update SAML config
		if _, err := rest.SamlConfigs.UpdateConfigWithContext(ctx, vmsId, idpName, data); err != nil {
			return nil, err
		}
		return nil, nil
	}

	// If no fields to update, just get the current SAML config
	return rest.SamlConfigs.GetConfigWithContext(ctx, vmsId, idpName)
}
