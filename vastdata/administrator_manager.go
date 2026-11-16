// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var AdministratorManagerSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"managers",
	http.MethodGet,
	"managers",
)

type AdministratorManager struct {
	tfstate is.ResourceState
}

func (m *AdministratorManager) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &AdministratorManager{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:          AdministratorManagerSchemaRef,
			WriteOnlyFields:    []string{"password", "password_retype"},
			WriteOnlyMD5Fields: []string{"password"},
			SensitiveFields:    []string{"password", "password_retype"},
			// password is write-only (not in state). Track a computed MD5 so
			// password rotations plan an update without storing the secret.
			AdditionalSchemaAttributes: map[string]any{
				"password_md5": ContentMD5SchemaAttribute("password", false),
			},
		},
	)}
}

func (m *AdministratorManager) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &AdministratorManager{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:       AdministratorManagerSchemaRef,
			SensitiveFields: []string{"password"},
		},
	)}
}

func (m *AdministratorManager) TfState() is.ResourceState {
	return m.tfstate
}

func (m *AdministratorManager) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Managers
}

// ModifyPlan sets password_md5 from the write-only password in config so
// Terraform can detect password changes without storing the secret in state.
func (m *AdministratorManager) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) error {
	SetPlanContentMD5(ctx, req, resp, "password", "password_md5")
	return nil
}

// TransformResponseRecord drops password material if the API ever returns it
// so the secret cannot land in Terraform state.
func (m *AdministratorManager) TransformResponseRecord(record Record) Record {
	delete(record, "password")
	delete(record, "password_retype")
	return record
}

func (m *AdministratorManager) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	// Transform roles from list of objects to list of integers
	// API returns: [{"id": 5, "name": "csi"}, {"id": 1, "name": "administrators"}]
	// Terraform expects: [5, 1]
	searchParams := getSearchParams(ctx, m.tfstate, nil)
	record, err := rest.Managers.GetWithContext(ctx, searchParams)
	if err != nil {
		return nil, err
	}

	if roles, ok := record["roles"]; ok && roles != nil {
		if rolesList, ok := roles.([]any); ok {
			roleIds := make([]any, 0, len(rolesList))
			for _, role := range rolesList {
				if roleMap, ok := role.(map[string]any); ok {
					if id, exists := roleMap["id"]; exists {
						roleIds = append(roleIds, id)
					}
				}
			}
			record["roles"] = roleIds
		}

	}
	return record, nil

}

func (m *AdministratorManager) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.ReadDatasource(ctx, rest)
}
