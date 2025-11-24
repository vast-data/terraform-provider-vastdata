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
	tfstate *is.TFState
}

func (m *AdministratorRole) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &AdministratorRole{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:            AdministratorRoleSchemaRef,
			ComputedSchemaFields: []string{"permissions_list"},
		},
	)}
}

func (m *AdministratorRole) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &AdministratorRole{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: AdministratorRoleSchemaRef,
		},
	)}
}

func (m *AdministratorRole) TfState() *is.TFState {
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
		record["permissions_list"] = permissions
	}
	return record, nil

}

func (m *AdministratorRole) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.ReadDatasource(ctx, rest)
}
