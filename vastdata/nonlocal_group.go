// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var NonlocalGroupSchemaRef = is.NewSchemaReference(
	http.MethodPatch,
	"groups/query",
	http.MethodGet,
	"groups/query",
)

type NonlocalGroup struct {
	tfstate *is.TFState
}

func (m *NonlocalGroup) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &NonlocalGroup{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:      NonlocalGroupSchemaRef,
			ReadOnlyFields: []string{"context"},
			AdditionalSchemaAttributes: map[string]any{
				"name": rschema.StringAttribute{
					Optional:    false,
					Computed:    true,
					Description: "The name of the non-local group.",
				},
				"s3_policies_ids": rschema.SetAttribute{
					ElementType: types.Int64Type,
					Optional:    true,
					Computed:    true,
					Description: "A set of IDs of S3 policies associated with the non-local group.",
				},
			},
		},
	)}
}

func (m *NonlocalGroup) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &NonlocalGroup{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:      NonlocalGroupSchemaRef,
			ReadOnlyFields: []string{"context"},
			AdditionalSchemaAttributes: map[string]any{
				"name": dschema.StringAttribute{
					Optional:    false,
					Computed:    true,
					Description: "The name of the non-local group.",
				},
				"gid": dschema.Int64Attribute{
					Optional:    true,
					Computed:    true,
					Description: "The gid of the non-local group.",
				},
				"sid": dschema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Description: "The sid of the non-local group.",
				},
				"s3_policies": dschema.SetAttribute{
					ElementType: types.StringType,
					Optional:    false,
					Computed:    true,
					Description: "A set of S3 policies associated with the non-local group.",
				},
				"s3_policies_ids": dschema.SetAttribute{
					ElementType: types.Int64Type,
					Optional:    false,
					Computed:    true,
					Description: "A set of IDs of S3 policies associated with the non-local group.",
				},
			},
		},
	)}
}

func (m *NonlocalGroup) TfState() *is.TFState {
	return m.tfstate
}

func (m *NonlocalGroup) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.NonLocalGroups
}

func (m *NonlocalGroup) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	return ensureNonlocalGroupUpdatedWith(ctx, ts, ts, rest)
}

func (m *NonlocalGroup) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	stateTs := m.tfstate
	planTs := plan.(*NonlocalGroup).TfState()
	return ensureNonlocalGroupUpdatedWith(ctx, stateTs, planTs, rest)
}

func (m *NonlocalGroup) DeleteResource(ctx context.Context, rest *VMSRest) error {
	// No-op for non-local groups, as they cannot be deleted.
	return nil
}

// ensureNonlocalGroupUpdatedWith verifies if the given NonLocalGroup (looked up by state)
// needs to be updated with new fields (like tenant_id or s3_policies_ids) and performs the update if necessary.
//
// This is used in both CreateResource and UpdateResource for NonLocalGroup.
func ensureNonlocalGroupUpdatedWith(ctx context.Context, stateTs, fieldsTs *is.TFState, rest *VMSRest) (DisplayableRecord, error) {
	searchParams := getSearchParams(ctx, stateTs, fieldsTs)
	record, err := rest.NonLocalGroups.GetWithContext(ctx, searchParams)
	if err != nil {
		return nil, err
	}
	searchParams.Without("context")
	if ok := fieldsTs.SetToMapIfAvailable(searchParams, "tenant_id", "s3_policies_ids"); ok {
		if _, err = rest.NonLocalGroups.UpdateNonLocalGroupWithContext(ctx, searchParams); err != nil {
			return nil, err
		}
	}
	return record, nil

}
