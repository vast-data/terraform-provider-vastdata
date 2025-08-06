// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var TenantConfiguredIdpSchemaRef = is.NewSchemaReference(
	"",
	"",
	http.MethodGet,
	"tenants/configured_idp",
)

type TenantConfiguredIdp struct {
	tfstate *is.TFState
}

func (m *TenantConfiguredIdp) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &TenantConfiguredIdp{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: TenantConfiguredIdpSchemaRef,
		},
	)}
}

func (m *TenantConfiguredIdp) TfState() *is.TFState {
	return m.tfstate
}

func (m *TenantConfiguredIdp) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Tenants
}

func (m *TenantConfiguredIdp) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	tenantName := m.tfstate.String("name")
	return rest.Tenants.GetConfiguredIdPWithContext(ctx, tenantName)
}
