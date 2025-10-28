// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var IamRoleSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"iamroles",
	http.MethodGet,
	"iamroles",
)

type IamRole struct {
	tfstate *is.TFState
}

func (m *IamRole) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &IamRole{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:      IamRoleSchemaRef,
			EditOnlyFields: []string{"revoke_access_keys"},
			AdditionalSchemaAttributes: map[string]any{
				"revoke_access_keys": rschema.BoolAttribute{
					Optional:    true,
					Description: "Indicates whether access keys should be revoked for this iam role.",
				},
			},
		},
	)}
}

func (m *IamRole) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &IamRole{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: IamRoleSchemaRef,
		},
	)}
}

func (m *IamRole) TfState() *is.TFState {
	return m.tfstate
}

func (m *IamRole) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.IamRoles
}

// AfterCreateResource is a custom post-creation handler for IamRole resources.
// This function handles the revoke_access_keys EditOnlyField by calling a separate
// API endpoint to revoke access keys for the newly created IAM role.
//
// Behavior:
// - Only executes when revoke_access_keys is explicitly set to true in the Terraform configuration
// - Calls the IamRoleRevokeAccessKeysWithContext_PATCH API endpoint using the newly created role's ID
// - This is a separate API call from the main IAM role creation endpoint
// - Used for immediate access key revocation after role creation
func (m *IamRole) AfterCreateResource(ctx context.Context, rest *VMSRest, record Record) error {
	ts := m.tfstate
	if ts.IsKnownAndNotNull("revoke_access_keys") && ts.Bool("revoke_access_keys") {
		iamRoleId := record.RecordID()
		if err := rest.IamRoles.IamRoleRevokeAccessKeysWithContext_PATCH(ctx, iamRoleId, nil); err != nil {
			return err
		}
	}
	return nil
}

// AfterUpdateResource is a custom post-update handler for IamRole resources.
// This function handles the revoke_access_keys EditOnlyField by calling a separate
// API endpoint to revoke access keys for the updated IAM role.
//
// Behavior:
// - Only executes when revoke_access_keys is explicitly set to true in the Terraform configuration
// - Calls the IamRoleRevokeAccessKeysWithContext_PATCH API endpoint using the updated role's ID
// - This is a separate API call from the main IAM role update endpoint
// - Used for immediate access key revocation after role update
func (m *IamRole) AfterUpdateResource(ctx context.Context, plan AfterUpdateResource, rest *VMSRest, record Record) error {
	ts := plan.(*IamRole).tfstate
	if ts.IsKnownAndNotNull("revoke_access_keys") && ts.Bool("revoke_access_keys") {
		iamRoleId := record.RecordID()
		if err := rest.IamRoles.IamRoleRevokeAccessKeysWithContext_PATCH(ctx, iamRoleId, nil); err != nil {
			return err
		}
	}
	return nil
}
