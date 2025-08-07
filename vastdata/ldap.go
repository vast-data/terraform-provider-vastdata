// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var LdapSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"ldaps",
	http.MethodGet,
	"ldaps",
)

type Ldap struct {
	tfstate *is.TFState
}

func (m *Ldap) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Ldap{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:        LdapSchemaRef,
			SearchableFields: []string{"domain_name"},
		},
	)}
}

func (m *Ldap) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Ldap{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: LdapSchemaRef,
		},
	)}
}

func (m *Ldap) TfState() *is.TFState {
	return m.tfstate
}

func (m *Ldap) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Ldaps
}
