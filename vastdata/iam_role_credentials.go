// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var IamRoleCredentialsSchemaRef = is.NewSchemaReference(
	http.MethodGet,
	"iamroles/{id}/credentials",
	http.MethodGet,
	"iamroles/{id}/credentials",
)

type IamRoleCredentials struct {
	tfstate *is.TFState
}

func (m *IamRoleCredentials) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &IamRoleCredentials{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: IamRoleCredentialsSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"iam_role_id": dschema.Int64Attribute{
					Required:    true,
					Description: "ID of the IAM role",
				},
			},
		},
	)}
}

func (m *IamRoleCredentials) TfState() *is.TFState {
	return m.tfstate
}

func (m *IamRoleCredentials) API(_ *VMSRest) VastResourceAPIWithContext {
	return nil
}

func (m *IamRoleCredentials) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	iamRoleId := ts.Int64("iam_role_id")
	accessKey := ts.String("access_key")
	return rest.IamRoles.IamRoleCredentialsWithContext_GET(ctx, iamRoleId, params{"access_key": accessKey})
}
