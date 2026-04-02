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

var CnodeBgpConfigSchemaRef = is.NewSchemaReference(
	http.MethodPatch,
	"cnodes/{id}/bgpconfig",
	http.MethodGet,
	"cnodes/{id}/bgpconfig",
)

type CnodeBgpConfig struct {
	tfstate *is.TFState
}

func (m *CnodeBgpConfig) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &CnodeBgpConfig{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: CnodeBgpConfigSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"cnode_id": rschema.Int64Attribute{
					Required:    true,
					Description: "ID of the Cnode",
				},
			},
		},
	)}
}

func (m *CnodeBgpConfig) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &CnodeBgpConfig{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: CnodeBgpConfigSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"cnode_id": dschema.Int64Attribute{
					Required:    true,
					Description: "ID of the Cnode",
				},
			},
		},
	)}
}

func (m *CnodeBgpConfig) TfState() *is.TFState {
	return m.tfstate
}

func (m *CnodeBgpConfig) API(_ *VMSRest) VastResourceAPIWithContext {
	return nil
}

func (m *CnodeBgpConfig) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	cnodeId := m.tfstate.Int64("cnode_id")
	return rest.Cnodes.CnodeBgpconfigWithContext_GET(ctx, cnodeId)
}

func (m *CnodeBgpConfig) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	record, err := m.ReadDatasource(ctx, rest)
	if err != nil {
		return nil, fmt.Errorf("error reading cnode bgp config: %w", err)
	}
	rec := record.(Record)
	m.tfstate.Set("enabled", rec["enabled"])
	m.tfstate.Set("subnet_bits", rec["subnet_bits"])
	m.tfstate.Set("self_asn", rec["self_asn"])
	m.tfstate.Set("port1_self_address", rec["port1_self_address"])
	m.tfstate.Set("port2_self_address", rec["port2_self_address"])
	m.tfstate.Set("port1_peer_address", rec["port1_peer_address"])
	m.tfstate.Set("port2_peer_address", rec["port2_peer_address"])
	return record, nil
}

func (m *CnodeBgpConfig) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	cnodeId := m.tfstate.Int64("cnode_id")
	params := m.tfstate.GetCreateParams()
	params.Without("cnode_id") // Remove cnode_id from params, as it's part of the URL

	return updateCnodeBgpConfig(ctx, cnodeId, params, rest)
}

func (m *CnodeBgpConfig) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	cnodeId := m.tfstate.Int64("cnode_id")
	planTs := plan.(*CnodeBgpConfig).TfState()
	params := planTs.GetChangedParams(m.tfstate)
	params.Without("cnode_id") // Remove cnode_id from params, as it's part of the URL

	return updateCnodeBgpConfig(ctx, cnodeId, params, rest)
}

func (m *CnodeBgpConfig) DeleteResource(ctx context.Context, rest *VMSRest) error {
	// No-op: CnodeBgpConfig cannot be deleted - it's a configuration operation
	return nil
}

// updateCnodeBgpConfig updates the BGP configuration for a specific cnode with the given parameters.
// This is used by both CreateResource and UpdateResource for CnodeBgpConfig.
func updateCnodeBgpConfig(ctx context.Context, cnodeId int64, params map[string]any, rest *VMSRest) (DisplayableRecord, error) {
	if cnodeId == 0 {
		return nil, fmt.Errorf("failed to get cnode ID: cnode ID is empty")
	}
	deleteZeroValues(params, []string{
		"self_asn",
		"port1_peer_address",
		"port2_peer_address",
		"port1_self_address",
		"port2_self_address",
	})
	if len(params) > 0 {
		err := rest.Cnodes.CnodeBgpconfigWithContext_PATCH(ctx, cnodeId, params)
		if err != nil {
			return nil, err
		}
	}
	// Return the updated config
	return rest.Cnodes.CnodeBgpconfigWithContext_GET(ctx, cnodeId)
}
