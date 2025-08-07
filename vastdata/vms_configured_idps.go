// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

type VmsConfiguredIdps struct {
	tfstate *is.TFState
}

func (m *VmsConfiguredIdps) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &VmsConfiguredIdps{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			TFStateHintsForCustom: &is.TFStateHintsForCustom{
				Description: "Configured SAML Identify Providers (IdPs).",
				SchemaAttributes: map[string]any{
					"vms_id": dschema.Int64Attribute{
						Required:    true,
						Description: "Unique ID of the VMS.",
					},
					"idps": dschema.ListAttribute{
						ElementType: types.StringType,
						Optional:    false,
						Computed:    true,
						Description: "List of configured IdPs.",
					},
				},
			},
		},
	)}
}

func (m *VmsConfiguredIdps) TfState() *is.TFState {
	return m.tfstate
}

func (m *VmsConfiguredIdps) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Vms
}

func (m *VmsConfiguredIdps) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	vmsId := m.tfstate.Int64("vms_id")
	idps, err := rest.Vms.GetConfiguredIdPsWithContext(ctx, vmsId)
	if err == nil {
		// Convert []string to []any for Terraform compatibility
		idpsAny := make([]any, len(idps))
		for i, idp := range idps {
			idpsAny[i] = idp
		}
		m.tfstate.Set("idps", idpsAny)
	}
	return nil, err
}
