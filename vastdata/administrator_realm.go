// Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var AdministratorRealmSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"realms",
	http.MethodGet,
	"realms",
)

type AdministratorRealm struct {
	tfstate *is.TFState
}

func (m *AdministratorRealm) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &AdministratorRealm{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: AdministratorRealmSchemaRef,
		},
	)}
}

func (m *AdministratorRealm) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &AdministratorRealm{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: AdministratorRealmSchemaRef,
		},
	)}
}

func (m *AdministratorRealm) TfState() *is.TFState {
	return m.tfstate
}

func (m *AdministratorRealm) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Realms
}
