// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var KerberosSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"kerberos",
	http.MethodGet,
	"kerberos",
)

type Kerberos struct {
	tfstate *is.TFState
}

func (m *Kerberos) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Kerberos{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:            KerberosSchemaRef,
			EditOnlyFields:       []string{"enabled"},
			OptionalSchemaFields: []string{"enabled"},
		},
	)}
}

func (m *Kerberos) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Kerberos{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: KerberosSchemaRef,
		}),
	}
}

func (m *Kerberos) TfState() *is.TFState {
	return m.tfstate
}

func (m *Kerberos) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Kerberos
}
