// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var DnsSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"dns",
	http.MethodGet,
	"dns",
)

type Dns struct {
	tfstate *is.TFState
}

func (m *Dns) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Dns{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: DnsSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				// vip is required in POST /dns/ but auto-assigned by POST /dns/allocate/.
				"vip": rschema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Description: "The VIP assigned to the DNS server. Auto-assigned when allocate is true.",
				},
				"allocate": rschema.BoolAttribute{
					Optional:    true,
					Description: "When true, use POST /dns/allocate/ to create the DNS entry with an auto-assigned VIP instead of requiring an explicit vip.",
				},
			},
			PreserveUserValueFields: []string{"allocate"},
		},
	)}
}

func (m *Dns) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Dns{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: DnsSchemaRef,
		}),
	}
}

func (m *Dns) TfState() *is.TFState {
	return m.tfstate
}

func (m *Dns) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Dns
}

// CreateResource handles both DNS creation paths:
//   - allocate=true  → POST /dns/allocate/ (VIP auto-assigned, async task)
//   - allocate=false → POST /dns/ (standard creation, vip required)
func (m *Dns) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate

	if ts.Bool("allocate") {
		if ts.IsKnownAndNotNull("vip") {
			return nil, fmt.Errorf("'vip' must not be set when allocate is true — the VIP is auto-assigned by the cluster")
		}
		body := ts.GetCreateParams()
		delete(body, "allocate")

		asyncRecord, err := rest.Dns.DnsAllocateWithContext_POST(ctx, body)
		if err != nil {
			return nil, err
		}
		if err := handleMaybeAsyncTask(ctx, rest, asyncRecord, 10*time.Minute); err != nil {
			return nil, fmt.Errorf("DNS allocate task failed: %w", err)
		}

		return rest.Dns.GetWithContext(ctx, ts.GetGenericSearchParams(ctx))
	}

	// Standard creation: POST /dns/ — vip is required.
	if !ts.IsKnownAndNotNull("vip") {
		return nil, fmt.Errorf("'vip' is required when allocate is false")
	}
	body := ts.GetCreateParams()
	delete(body, "allocate")
	return rest.Dns.CreateWithContext(ctx, body)
}
