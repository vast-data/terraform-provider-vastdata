// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var TenantNfs4DelegationSchemaRef = is.NewSchemaReference(
	http.MethodGet,
	"tenants/{id}/nfs4_delegs",
	http.MethodGet,
	"tenants/{id}/nfs4_delegs",
)

type TenantNfs4Delegation struct {
	tfstate *is.TFState
}

func (m *TenantNfs4Delegation) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &TenantNfs4Delegation{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: TenantNfs4DelegationSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"tenant_id": rschema.Int64Attribute{
					Required:    true,
					Description: "ID of the tenant",
				},
				"file_path": rschema.StringAttribute{
					Required:    true,
					Description: "File path for the NFS4 delegation",
				},
			},
		},
	)}
}

func (m *TenantNfs4Delegation) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &TenantNfs4Delegation{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: TenantNfs4DelegationSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"tenant_id": dschema.Int64Attribute{
					Required:    true,
					Description: "ID of the tenant",
				},
				"file_path": dschema.StringAttribute{
					Required:    true,
					Description: "File path for the NFS4 delegation",
				},
			},
		},
	)}
}

func (m *TenantNfs4Delegation) TfState() *is.TFState {
	return m.tfstate
}

func (m *TenantNfs4Delegation) API(rest *VMSRest) VastResourceAPIWithContext {
	return nil
}

func (m *TenantNfs4Delegation) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	tenantId := m.tfstate.Int64("tenant_id")
	params := getSearchParams(ctx, m.tfstate, nil)
	params.Without("tenant_id")
	return rest.Tenants.TenantNfs4DelegsWithContext_GET(ctx, tenantId, params)
}

func (m *TenantNfs4Delegation) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.ReadDatasource(ctx, rest)
}

func (m *TenantNfs4Delegation) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.ReadDatasource(ctx, rest)
}

func (m *TenantNfs4Delegation) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	return m.ReadDatasource(ctx, rest)
}

func (m *TenantNfs4Delegation) DeleteResource(ctx context.Context, rest *VMSRest) error {
	tenantId := m.tfstate.Int64("tenant_id")
	params := getSearchParams(ctx, m.tfstate, nil)
	return rest.Tenants.TenantNfs4DelegWithContext_DELETE(ctx, tenantId, params)
}
