// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"

	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var TenantSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"tenants",
	http.MethodGet,
	"tenants",
)

// tenantViewsCountSubResource adds views_count as a flattened, trigger-gated sub-resource
// on the tenant (VAST >= 5.5.0). Set get_views_count = false to opt out of fetching
// GET /tenants/{id}/views_count/.
var tenantViewsCountSubResource = is.SubResourceHint{
	MinVastVersion: VastVersion550,
	FieldTrigger:   "get_views_count",
	SchemaAttributes: map[string]any{
		"get_views_count": rschema.BoolAttribute{
			Optional: true,
			Description: "Controls fetching of tenant views count (requires VAST >= 5.5.0). " +
				"When unset or true the count is fetched automatically on supported clusters. " +
				"Set to false to explicitly opt out.",
		},
		"views_count": rschema.Int64Attribute{
			Computed:    true,
			Description: "The number of views currently present in this tenant. Populated automatically on VAST >= 5.5.0 unless get_views_count is false.",
		},
	},
}

type Tenant struct {
	tfstate *is.TFState
}

func (m *Tenant) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Tenant{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:             TenantSchemaRef,
			DeleteOnlyParamFields: map[string]string{"force_delete": "force"},
			PreserveOrderFields:   []string{"client_ip_ranges"},
			SubResources:          []is.SubResourceHint{tenantViewsCountSubResource},
			AdditionalSchemaAttributes: map[string]any{
				"force_delete": rschema.BoolAttribute{
					Optional: true,
					Description: "If set to true, forces deletion of the tenant even if it has empty subdirectories" +
						" or other removable remnants. Use with caution, as this will bypass standard cleanup checks.",
				},
			},
		},
	)}
}

func (m *Tenant) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Tenant{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:           TenantSchemaRef,
			PreserveOrderFields: []string{"client_ip_ranges"},
			SubResources:        []is.SubResourceHint{tenantViewsCountSubResource},
		}),
	}
}

func (m *Tenant) TfState() *is.TFState {
	return m.tfstate
}

func (m *Tenant) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Tenants
}

// GetSubResources fetches /tenants/{id}/views_count/ on VAST clusters running
// version >= 5.5.0. Set get_views_count = false to opt out.
func (m *Tenant) GetSubResources(ctx context.Context, rest *VMSRest, record Record, clusterVersion *version.Version) (Record, error) {
	if clusterVersion == nil || clusterVersion.LessThan(VastVersion550) {
		return nil, nil
	}
	if m.tfstate.IsKnownAndNotNull("get_views_count") && !m.tfstate.Bool("get_views_count") {
		return nil, nil
	}

	id := record.RecordID()
	rec, err := rest.Tenants.TenantViewsCountWithContext_GET(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch views_count for tenant %v: %w", id, err)
	}

	return Record{"views_count": rec["current_views_count"]}, nil
}

// TransformResponseRecord normalizes the "vippools" field in the tenant API response.
// When a VIP pool has tenant_id=null (shared across all tenants), the VAST API returns
// each entry as a tuple ["name", id] instead of an object {"name": ..., "id": ...}.
// This converts any such tuples to the proper map form before deserialization.
func (m *Tenant) TransformResponseRecord(record Record) Record {
	raw, ok := record["vippools"]
	if !ok {
		return record
	}
	list, ok := raw.([]any)
	if !ok {
		return record
	}
	normalized := make([]any, 0, len(list))
	for _, item := range list {
		tuple, ok := item.([]any)
		if !ok || len(tuple) != 2 {
			normalized = append(normalized, item)
			continue
		}
		normalized = append(normalized, map[string]any{
			"name": fmt.Sprintf("%v", tuple[0]),
			"id":   tuple[1],
		})
	}
	record["vippools"] = normalized
	return record
}
