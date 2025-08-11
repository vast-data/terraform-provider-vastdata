// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	planmodifiers "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

type BlockHostMapping struct {
	tfstate *is.TFState
}

func (m *BlockHostMapping) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &BlockHostMapping{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			TFStateHintsForCustom: &is.TFStateHintsForCustom{
				Description: "One-to-one mapping between a block host and a volume. This resource attaches a host to a volume via explicit identifiers.",
				SchemaAttributes: map[string]any{
					"id": rschema.Int64Attribute{
						Computed:    true,
						Description: "Unique ID of the block host mapping.",
					},
					"host_id": rschema.Int64Attribute{
						Required:    true,
						Description: "ID of the host to be mapped.",
						PlanModifiers: []planmodifiers.Int64{
							int64planmodifier.RequiresReplace(),
						},
					},
					"volume_id": rschema.Int64Attribute{
						Required:    true,
						Description: "ID of the volume to be mapped.",
						PlanModifiers: []planmodifiers.Int64{
							int64planmodifier.RequiresReplace(),
						},
					},
				},
			},
		},
	)}
}

func (m *BlockHostMapping) TfState() *is.TFState {
	return m.tfstate
}

func (m *BlockHostMapping) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.BlockHostMappings
}

func (m *BlockHostMapping) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	if !m.tfstate.IsKnownAndNotNull("volume_id") || !m.tfstate.IsKnownAndNotNull("host_id") {
		return nil, errors.New("volume_id and host_id must be known and not null for BlockHostMapping")
	}
	volumeId := m.tfstate.Int64("volume_id")
	hostId := m.tfstate.Int64("host_id")
	return rest.BlockHostMappings.GetWithContext(ctx, params{"volume__id": volumeId, "block_host__id": hostId})
}

func (m *BlockHostMapping) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	volumeId := m.tfstate.Int64("volume_id")
	hostId := m.tfstate.Int64("host_id")
	return rest.BlockHostMappings.EnsureMapWithContext(ctx, hostId, volumeId)
}

func (m *BlockHostMapping) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	planTfstate := plan.(*BlockHostMapping).tfstate
	volumeId := planTfstate.Int64("volume_id")
	hostId := planTfstate.Int64("host_id")
	return rest.BlockHostMappings.EnsureMapWithContext(ctx, hostId, volumeId)
}

func (m *BlockHostMapping) DeleteResource(ctx context.Context, rest *VMSRest) error {
	volumeId := m.tfstate.Int64("volume_id")
	hostId := m.tfstate.Int64("host_id")
	_, err := rest.BlockHostMappings.UnMapWithContext(ctx, hostId, volumeId)
	return err
}
