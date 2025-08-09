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

var NonlocalUserSchemaRef = is.NewSchemaReference(
	http.MethodPatch,
	"users/query",
	http.MethodGet,
	"users/query",
)

type NonlocalUser struct {
	tfstate *is.TFState
}

func (m *NonlocalUser) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &NonlocalUser{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			ReadOnlyFields: []string{"context", "vid"},
			SchemaRef:      NonlocalUserSchemaRef,
			ImportFields:   []string{"username", "context", "tenant_id"},
			AdditionalSchemaAttributes: map[string]any{
				"access_keys": rschema.SetNestedAttribute{
					Computed:    true,
					Description: "A set of access keys with creation time, key, remote, and status.",
					NestedObject: rschema.NestedAttributeObject{
						Attributes: map[string]rschema.Attribute{
							"creation_time": rschema.StringAttribute{Computed: true},
							"key":           rschema.StringAttribute{Computed: true},
							"remote":        rschema.StringAttribute{Computed: true},
							"status":        rschema.StringAttribute{Computed: true},
						},
					},
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

func (m *NonlocalUser) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &NonlocalUser{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:      NonlocalUserSchemaRef,
			ReadOnlyFields: []string{"context", "vid"},
			AdditionalSchemaAttributes: map[string]any{
				"access_keys": dschema.SetNestedAttribute{
					Computed:    true,
					Optional:    false,
					Description: "A set of access keys with creation time, key, remote, and status.",
					NestedObject: dschema.NestedAttributeObject{
						Attributes: map[string]dschema.Attribute{
							"creation_time": dschema.StringAttribute{Computed: true},
							"key":           dschema.StringAttribute{Computed: true},
							"remote":        dschema.StringAttribute{Computed: true},
							"status":        dschema.StringAttribute{Computed: true},
						},
					},
				},
			},
		}),
	}
}

func (m *NonlocalUser) TfState() *is.TFState {
	return m.tfstate
}

func (m *NonlocalUser) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.NonLocalUsers
}

func (m *NonlocalUser) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	return ensureNonlocalUserUpdatedWith(ctx, ts, ts, rest)
}

func (m *NonlocalUser) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	stateTs := m.tfstate
	planTs := plan.(*NonlocalUser).TfState()
	return ensureNonlocalUserUpdatedWith(ctx, stateTs, planTs, rest)
}

func (m *NonlocalUser) DeleteResource(ctx context.Context, rest *VMSRest) error {
	// No-op for non-local users, as they cannot be deleted.
	return nil
}

// ensureNonlocalUserUpdatedWith verifies if the given NonLocalUser (looked up by state)
// needs to be updated with new fields (like tenant_id or s3_policies_ids) and performs the update if necessary.
//
// This is used in both CreateResource and UpdateResource for NonLocalUser.
func ensureNonlocalUserUpdatedWith(ctx context.Context, stateTs, fieldsTs *is.TFState, rest *VMSRest) (DisplayableRecord, error) {
	searchParams := getSearchParams(ctx, stateTs, fieldsTs)
	record, err := rest.NonLocalUsers.GetWithContext(ctx, searchParams)
	if err != nil {
		return nil, err
	}

	searchParams.Without("context", "vid")
	if ok := fieldsTs.SetToMapIfAvailable(
		searchParams,
		"tenant_id",
		"s3_policies_ids",
		"allow_create_bucket",
		"allow_delete_bucket",
		"username",
		"s3_superuser",
	); ok {
		if _, err = rest.NonLocalUsers.UpdateNonLocalUserWithContext(ctx, searchParams); err != nil {
			return nil, err
		}
	}
	return record, nil

}
