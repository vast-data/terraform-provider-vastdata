// Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
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
			SchemaRef: AdministratorRoleSchemaRef,
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
