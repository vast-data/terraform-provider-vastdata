// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"
	"time"

	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

const vipPoolAllocateTimeout = 5 * time.Minute

var VipPoolSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"vippools",
	http.MethodGet,
	"vippools",
)

type VipPool struct {
	tfstate *is.TFState
}

func (m *VipPool) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &VipPool{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:               VipPoolSchemaRef,
			NotRequiredSchemaFields: []string{"subnet_cidr", "ip_ranges"},
			ReadOnlyFields:          []string{"serves_tenant"},
			PreserveOrderFields:     []string{"ip_ranges", "client_monitoring_ips"},
			PreserveUserValueFields: []string{"ips_count"},
			AdditionalSchemaAttributes: map[string]any{
				"ips_count": rschema.Int64Attribute{
					Optional: true,
					Description: "Number of IPs to allocate automatically. When set, the provider uses " +
						"POST /vippools/allocate/ for creation and PATCH /vippools/{id}/reallocate/ for " +
						"updates instead of the standard endpoints. Mutually exclusive with ip_ranges.",
				},
			},
		},
	)}
}

func (m *VipPool) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &VipPool{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:           VipPoolSchemaRef,
			PreserveOrderFields: []string{"ip_ranges", "client_monitoring_ips"},
		}),
	}
}

func (m *VipPool) TfState() *is.TFState {
	return m.tfstate
}

func (m *VipPool) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.VipPools
}

// CreateResource handles both the standard /vippools/ POST and the /vippools/allocate/ POST.
// When ips_count is set in the configuration, the allocate endpoint is used and ip_ranges
// must not be specified. Otherwise the standard endpoint is used with ip_ranges.
func (m *VipPool) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	createParams := ts.GetCreateParams()

	ipsCount := createParams["ips_count"]
	if ipsCount != nil {
		delete(createParams, "ip_ranges")
		return rest.VipPools.VipPoolAllocateWithContext_POST(ctx, createParams)
	}

	tflog.Debug(ctx, "VipPool.CreateResource: standard mode")
	delete(createParams, "ips_count")
	return rest.VipPools.CreateWithContext(ctx, createParams)
}

// UpdateResource handles both the standard /vippools/{id}/ PATCH and the
// /vippools/{id}/reallocate/ PATCH endpoint.
// When ips_count is set in the plan, the reallocate endpoint is used:
// the currently allocated ip_ranges are removed and ips_count new IPs are allocated.
func (m *VipPool) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	stateTs := m.tfstate
	planTs := plan.(*VipPool).TfState()

	planParams := planTs.GetCreateParams()
	ipsCount := planParams["ips_count"]

	id := stateTs.Int64("id")

	if ipsCount != nil {
		tflog.Debug(ctx, fmt.Sprintf("VipPool.UpdateResource: reallocate mode, id=%d, ips_count=%v", id, ipsCount))

		body := params{"ips_count_to_add": ipsCount}
		if ipRangesVal, ok := stateTs.Raw["ip_ranges"]; ok && !ipRangesVal.IsNull() && !ipRangesVal.IsUnknown() {
			body["ip_ranges_to_remove"] = is.ConvertAttrValueToRaw(ipRangesVal, stateTs.Type("ip_ranges"))
		}

		asyncResult, err := rest.VipPools.VipPoolReallocateWithContext_PATCH(ctx, id, body, vipPoolAllocateTimeout)
		if err != nil {
			return nil, fmt.Errorf("vippool reallocate failed: %w", err)
		}
		if asyncResult != nil && asyncResult.IsFailed() {
			return nil, fmt.Errorf("vippool reallocate task failed: %v", asyncResult.Err)
		}

		return rest.VipPools.GetByIdWithContext(ctx, id)
	}

	// Standard PATCH update.
	updateParams := planTs.GetUpdateParams(stateTs)
	delete(updateParams, "ips_count")
	delete(updateParams, "id")

	if len(updateParams) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("VipPool.UpdateResource: no changes for id=%d, refreshing state", id))
		return rest.VipPools.GetByIdWithContext(ctx, id)
	}

	tflog.Debug(ctx, fmt.Sprintf("VipPool.UpdateResource: standard update, id=%d", id))
	return rest.VipPools.UpdateWithContext(ctx, id, updateParams)
}
