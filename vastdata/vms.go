// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var VmsSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"vms",
	http.MethodGet,
	"vms",
)

type Vms struct {
	tfstate *is.TFState
}

func (m *Vms) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Vms{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			TFStateHintsForCustom: &is.TFStateHintsForCustom{
				Description: "Vast Management Service (VMS) settings.",
				SchemaAttributes: map[string]any{
					"id": rschema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: "Unique ID of the VMS.",
					},
					"name": rschema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Name of the VMS.",
					},
					"max_api_tokens_per_user": rschema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: "Maximum number of API tokens per user.",
					},
				},
			},
		},
	)}
}

func (m *Vms) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Vms{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: VmsSchemaRef,
		}),
	}
}

func (m *Vms) TfState() *is.TFState {
	return m.tfstate
}

func (m *Vms) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Vms
}

func (m *Vms) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	if !ts.IsKnownAndNotNull("id") {
		if ts.IsKnownAndNotNull("name") {
			// If name is known, we can fetch the VMS by name.
			record, err := rest.Vms.GetWithContext(ctx, params{"name": ts.String("name")})
			if err != nil {
				return nil, fmt.Errorf("failed to get VMS by name: %w", err)
			}
			ts.Set("id", record.RecordID())
			return record, nil
		} else {
			return nil, fmt.Errorf("VMS ID or name must be provided to read the resource")
		}
	}
	return rest.Vms.GetByIdWithContext(ctx, m.tfstate.Int64("id"))
}

func (m *Vms) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	if _, err := m.ReadResource(ctx, rest); err != nil {
		return nil, err
	}
	vmsId := ts.Int64("id")
	record, err := rest.Vms.GetByIdWithContext(ctx, vmsId)
	if err != nil {
		return nil, err
	}
	if ts.IsKnownAndNotNull("max_api_tokens_per_user") {
		maxTokens := ts.Int64("max_api_tokens_per_user")
		if _, err = rest.Vms.SetMaxApiTokensPerUser(vmsId, maxTokens); err != nil {
			return nil, fmt.Errorf("failed to set max_api_tokens_per_user: %w", err)
		}
		record["max_api_tokens_per_user"] = maxTokens
	}
	return record, err
}

func (m *Vms) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	var (
		ts          = m.tfstate
		vmsId       = ts.Int64("id")
		planManager = plan.(*Vms)
		planTs      = planManager.tfstate
		err         error
	)

	if planTs.IsKnownAndNotNull("max_api_tokens_per_user") {
		maxTokens := planTs.Int64("max_api_tokens_per_user")
		if _, err = rest.Vms.SetMaxApiTokensPerUser(vmsId, maxTokens); err != nil {
			return nil, fmt.Errorf("failed to set max_api_tokens_per_user: %w", err)
		}
		ts.Set("max_api_tokens_per_user", maxTokens)
	}
	return nil, nil

}

func (m *Vms) DeleteResource(ctx context.Context, rest *VMSRest) error {
	// No-op. Vms cannot be deleted.
	return nil
}
