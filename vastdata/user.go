// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var UserSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"users",
	http.MethodGet,
	"users",
)

type User struct {
	tfstate *is.TFState
}

func (m *User) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &User{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:       UserSchemaRef,
			SensitiveFields: []string{"password"},
		},
	)}
}

func (m *User) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &User{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: UserSchemaRef,
		}),
	}
}

func (m *User) TfState() *is.TFState {
	return m.tfstate
}

func (m *User) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Users
}
