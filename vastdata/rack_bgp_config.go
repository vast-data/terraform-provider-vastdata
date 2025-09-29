// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var RackBgpConfigSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"racks/{id}/bgpconfig",
	"",
	"",
)

type RackBgpConfig struct {
	tfstate *is.TFState
}

func (m *RackBgpConfig) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &RackBgpConfig{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: RackBgpConfigSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"rack_id": rschema.Int64Attribute{
					Required:    true,
					Description: "ID of the Rack",
				},
			},
		},
	)}
}

func (m *RackBgpConfig) TfState() *is.TFState {
	return m.tfstate
}

func (m *RackBgpConfig) API(_ *VMSRest) VastResourceAPIWithContext {
	return nil
}

func (m *RackBgpConfig) ReadResource(_ context.Context, _ *VMSRest) (DisplayableRecord, error) {
	return nil, nil
}

func (m *RackBgpConfig) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	rackId := ts.Int64("rack_id")
	params := ts.GetCreateParams()
	params.Without("rack_id") // Remove rack_id from params, as it's part of the URL
	return rest.Racks.UpdateBgpConfigWithContext(ctx, rackId, params)
}

func (m *RackBgpConfig) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	planTs := plan.(*RackBgpConfig).TfState()
	rackId := ts.Int64("rack_id")
	params := ts.GetCreateParams()
	// Merge plan changes into params
	params.Update(planTs.GetCreateParams(), true)
	params.Without("rack_id") // Remove rack_id from params, as it's part of the URL
	return rest.Racks.UpdateBgpConfigWithContext(ctx, rackId, params)
}

func (m *RackBgpConfig) DeleteResource(_ context.Context, _ *VMSRest) error {
	// No-op: KerberosKeytab cannot be deleted - it's a one-time operation
	return nil
}
