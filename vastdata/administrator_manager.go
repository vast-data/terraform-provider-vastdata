// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
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
			PreserveUserValueFieldsWhenApiReturnsNull: []string{"password"},
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
