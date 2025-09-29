// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var ManagerAuthorizedStatusSchemaRef = is.NewSchemaReference(
	"",
	"",
	http.MethodGet,
	"managers/authorized_status",
)

type ManagerAuthorizedStatus struct {
	tfstate *is.TFState
}

func (m *ManagerAuthorizedStatus) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ManagerAuthorizedStatus{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: ManagerAuthorizedStatusSchemaRef,
		},
	)}
}

func (m *ManagerAuthorizedStatus) TfState() *is.TFState {
	return m.tfstate
}

func (m *ManagerAuthorizedStatus) API(rest *VMSRest) VastResourceAPIWithContext {
	return nil
}

func (m *ManagerAuthorizedStatus) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	params := getSearchParams(ctx, m.tfstate, nil)
	return rest.Managers.GetAuthorizedStatusWithContext(ctx, params)

}
