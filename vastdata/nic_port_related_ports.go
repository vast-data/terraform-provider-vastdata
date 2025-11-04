// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

type NicPortRelatedPorts struct {
	tfstate *is.TFState
}

func (m *NicPortRelatedPorts) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &NicPortRelatedPorts{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			TFStateHintsForCustom: &is.TFStateHintsForCustom{
				Description: "Datasource for retrieving related NIC port IDs",
				SchemaAttributes: map[string]any{
					"nic_port_id": dschema.Int64Attribute{
						Required:    true,
						Description: "ID of the NIC port to get related ports for",
					},
					"related_nic_port_ids": dschema.ListAttribute{
						ElementType: types.Int64Type,
						Computed:    true,
						Description: "List of IDs of NIC ports that are logically the same physical port",
					},
				},
			},
		},
	)}
}

func (m *NicPortRelatedPorts) TfState() *is.TFState {
	return m.tfstate
}

func (m *NicPortRelatedPorts) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.NicPorts
}

func (m *NicPortRelatedPorts) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	nicPortId := m.tfstate.Int64("nic_port_id")
	if nicPortId == 0 {
		return nil, fmt.Errorf("nic_port_id is required")
	}

	relatedIds, err := rest.NicPorts.NicPortRelatedNicportsWithContext_GET(ctx, nicPortId)
	if err != nil {
		return nil, fmt.Errorf("failed to get related NIC ports: %w", err)
	}

	// Convert []int64 to []any for the record
	relatedIdsAny := make([]any, len(relatedIds))
	for i, id := range relatedIds {
		relatedIdsAny[i] = id
	}

	return Record{
		"nic_port_id":          nicPortId,
		"related_nic_port_ids": relatedIdsAny,
	}, nil
}
