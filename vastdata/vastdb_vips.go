// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

type VastDbVips struct {
	tfstate *is.TFState
}

func (m *VastDbVips) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &VastDbVips{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			TFStateHintsForCustom: &is.TFStateHintsForCustom{
				Description: "Datasource for retrieving VastDB VIPs",
				SchemaAttributes: map[string]any{
					"tenant_id": dschema.Int64Attribute{
						Description: "Tenant id",
						Required:    true,
					},
					"vips": dschema.ListNestedAttribute{
						NestedObject: dschema.NestedAttributeObject{
							Attributes: map[string]dschema.Attribute{
								"tenant_id": dschema.Int64Attribute{
									Computed:    true,
									Description: "ID of the tenant",
								},
								"port": dschema.Int64Attribute{
									Computed:    true,
									Description: "Port number",
								},
								"vip": dschema.StringAttribute{
									Computed:    true,
									Description: "VIP address",
								},
							},
						},
						Computed:    true,
						Description: "List of VastDB VIPs",
					},
				},
			},
		},
	)}
}

func (m *VastDbVips) TfState() *is.TFState {
	return m.tfstate
}

func (m *VastDbVips) API(_ *VMSRest) VastResourceAPIWithContext {
	return nil
}

func (m *VastDbVips) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	tenantId := m.tfstate.Int64("tenant_id")
	recordSet, err := rest.Vips.VipVipsWithContext_GET(ctx, params{"tenant_id": tenantId})
	if err != nil {
		return nil, err
	}

	// Convert RecordSet to []any for the list field
	// Each item needs to be map[string]any, not Record
	vipsList := make([]any, len(recordSet))
	for i, record := range recordSet {
		// Convert Record to map[string]any
		recordMap := make(map[string]any)
		for key, value := range record {
			recordMap[key] = value
		}
		vipsList[i] = recordMap
	}

	return Record{"vips": vipsList}, nil
}
