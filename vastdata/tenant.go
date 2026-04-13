// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"fmt"
	"net/http"

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
		}),
	}
}

func (m *Tenant) TfState() *is.TFState {
	return m.tfstate
}

func (m *Tenant) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Tenants
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
