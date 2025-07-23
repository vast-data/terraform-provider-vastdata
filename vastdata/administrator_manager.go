// Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var AdministratorManagerSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"managers",
	http.MethodGet,
	"managers",
)

type AdministratorManager struct {
	tfstate *is.TFState
}

func (m *AdministratorManager) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &AdministratorManager{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:       AdministratorManagerSchemaRef,
			SensitiveFields: []string{"password"},
		},
	)}
}

func (m *AdministratorManager) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &AdministratorManager{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:       AdministratorManagerSchemaRef,
			SensitiveFields: []string{"password"},
		},
	)}
}

func (m *AdministratorManager) TfState() *is.TFState {
	return m.tfstate
}

func (m *AdministratorManager) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Managers
}
