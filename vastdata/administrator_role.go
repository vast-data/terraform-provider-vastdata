// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var AdministratorRoleSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"roles",
	http.MethodGet,
	"roles",
)

type AdministratorRole struct {
	tfstate is.ResourceState
}

func (m *AdministratorRole) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &AdministratorRole{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:            AdministratorRoleSchemaRef,
			ComputedSchemaFields: []string{"permissions_list"},
			// Request "permissions" are legacy types (view/create/edit/delete).
			// Response "permissions" are permission codenames (mapped to permissions_list).
			// Same JSON key — do not let GET overwrite the configured types with [].
			PreserveUserValueFields: []string{"permissions"},
			NotComputedSchemaFields: []string{"permissions"},
		},
	)}
}

func (m *AdministratorRole) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &AdministratorRole{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: AdministratorRoleSchemaRef,
		},
	)}
}

func (m *AdministratorRole) TfState() is.ResourceState {
	return m.tfstate
}

func (m *AdministratorRole) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Roles
}

func (m *AdministratorRole) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	searchParams := getSearchParams(ctx, m.tfstate, nil)
	record, err := rest.Roles.GetWithContext(ctx, searchParams)
	if err != nil {
		return nil, err
	}

	if permissions, ok := record["permissions"]; ok && permissions != nil {
		// Response permissions are codenames; surface them on permissions_list.
		record["permissions_list"] = permissions
	}
	// Drop response permissions so they cannot clobber TF "permissions" (legacy types).
	delete(record, "permissions")
	return record, nil

}

func (m *AdministratorRole) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.ReadDatasource(ctx, rest)
}
