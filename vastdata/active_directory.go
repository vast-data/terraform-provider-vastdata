// Copyright (c) HashiCorp, Inc.

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var ActiveDirectorySchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"activedirectory",
	http.MethodGet,
	"activedirectory",
)

type ActiveDirectory struct {
	tfstate *is.TFState
}

func (m *ActiveDirectory) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &ActiveDirectory{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:        ActiveDirectorySchemaRef,
			SearchableFields: []string{"ldap_id", "domain_name", "machine_account_name"},
		},
	)}
}

func (m *ActiveDirectory) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ActiveDirectory{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: ActiveDirectorySchemaRef,
		},
	)}
}

func (m *ActiveDirectory) TfState() *is.TFState {
	return m.tfstate
}

func (m *ActiveDirectory) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.ActiveDirectories
}
