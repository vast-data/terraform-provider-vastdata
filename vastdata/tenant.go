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

// tenantViewsCountSubResource declares the read-only sub-resource for
// GET /tenants/{id}/views_count/ (VAST >= 5.5.0).
// Set get_views_count = false to explicitly opt out of fetching.
var tenantViewsCountSubResource = is.SubResourceHint{
	MinVastVersion: VastVersion550,
	FieldTrigger:   "get_views_count",
	SchemaKey:      "view_count",
	SchemaAttributes: map[string]any{
		"current_views_count": rschema.Int64Attribute{
			Computed:    true,
			Description: "The number of views currently present in this tenant.",
		},
		"max_views": rschema.Int64Attribute{
			Computed:    true,
			Description: "Maximum number of views allowed in this tenant.",
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
			SchemaRef:               TenantSchemaRef,
			NotComputedSchemaFields: []string{"local_provider_id"},
			DeleteOnlyParamFields: map[string]string{
				"force_delete":          "force",
				"delete_local_provider": "", // custom: handled via LocalProviders API (destroy + create rollback)
			},
			PreserveOrderFields: []string{"client_ip_ranges"},
			SubResources:        []is.SubResourceHint{tenantViewsCountSubResource},
			AdditionalSchemaAttributes: map[string]any{
				"force_delete": rschema.BoolAttribute{
					Optional: true,
					Description: "If set to true, forces deletion of the tenant even if it has empty subdirectories" +
						" or other removable remnants. Use with caution, as this will bypass standard cleanup checks.",
				},
				"delete_local_provider": rschema.BoolAttribute{
					Optional: true,
					Description: "If set to true, deletes the VMS auto-created local provider for this tenant " +
						"on destroy and on create rollback. Cannot be used together with local_provider_id.",
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

// ValidateResourceConfig rejects delete_local_provider together with a user-supplied
// local_provider_id. Only VMS auto-created providers (provider-<tenant_name>) may be
// deleted; a configured local_provider_id is out of this resource's control.
func (m *Tenant) ValidateResourceConfig(context.Context) error {
	if m.tfstate == nil {
		return nil
	}
	deleteLP := m.tfstate.IsKnownAndNotNull("delete_local_provider") && m.tfstate.Bool("delete_local_provider")
	hasLPID := m.tfstate.IsKnownAndNotNull("local_provider_id")
	if deleteLP && hasLPID {
		return fmt.Errorf(
			"delete_local_provider cannot be set together with local_provider_id: " +
				"only VMS auto-created local providers (provider-<tenant_name>) may be deleted; " +
				"a user-supplied local_provider_id is out of this resource's control",
		)
	}
	return nil
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

	return Record{
		"view_count": map[string]any{
			"current_views_count": rec["current_views_count"],
			"max_views":           rec["max_views"],
		},
	}, nil
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

func dedicatedLocalProviderName(tenantName string) string {
	return "provider-" + tenantName
}

func (m *Tenant) findDedicatedLocalProvider(ctx context.Context, rest *VMSRest) (Record, error) {
	tenantName := m.tfstate.String("name")
	if tenantName == "" {
		return nil, nil
	}
	rec, err := rest.LocalProviders.GetWithContext(ctx, params{"name": dedicatedLocalProviderName(tenantName)})
	if err != nil {
		if isNotFoundErr(err) {
			return nil, nil
		}
		return nil, err
	}
	if rec != nil && !rec.Empty() {
		return rec, nil
	}
	return nil, nil
}

func (m *Tenant) AfterDeleteResource(ctx context.Context, rest *VMSRest) error {
	if m.tfstate == nil || !m.tfstate.IsKnownAndNotNull("delete_local_provider") || !m.tfstate.Bool("delete_local_provider") {
		return nil
	}

	rec, err := m.findDedicatedLocalProvider(ctx, rest)
	if err != nil {
		return err
	}
	if rec == nil || rec.Empty() {
		return nil
	}
	_, err = rest.LocalProviders.DeleteByIdWithContext(ctx, rec.RecordID(), nil, nil)
	return ignoreResourceGone(err)
}
