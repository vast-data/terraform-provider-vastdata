// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
			AdditionalSchemaAttributes: map[string]any{
				"s3_policies_ids": rschema.SetAttribute{
					ElementType: types.Int64Type,
					Optional:    true,
					Computed:    true,
					Description: "S3 policies IDs, denoting which S3 identity policies are associated with the user. The user is granted and denied S3 permissions according to the associated S3 identity policies.",
				},
			},
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
