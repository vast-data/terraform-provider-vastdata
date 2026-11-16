// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var ManagerPasswordSchemaRef = is.NewSchemaReference(
	http.MethodPatch,
	"managers/password",
	"",
	"",
)

type ManagerPassword struct {
	tfstate is.ResourceState
}

func (m *ManagerPassword) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &ManagerPassword{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:          ManagerPasswordSchemaRef,
			WriteOnlyFields:    []string{"password"},
			WriteOnlyMD5Fields: []string{"password"},
			SensitiveFields:    []string{"password"},
			AdditionalSchemaAttributes: map[string]any{
				"password_md5": ContentMD5SchemaAttribute("password", true),
			},
		},
	)}
}
func (m *ManagerPassword) TfState() is.ResourceState {
	return m.tfstate
}

func (m *ManagerPassword) API(_ *VMSRest) VastResourceAPIWithContext {
	return nil
}

func (m *ManagerPassword) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) error {
	SetPlanContentMD5(ctx, req, resp, "password", "password_md5")
	return nil
}

func (m *ManagerPassword) TransformResponseRecord(record Record) Record {
	DeleteRecordFields(record, "password")
	return record
}

func (m *ManagerPassword) ReadResource(_ context.Context, _ *VMSRest) (DisplayableRecord, error) {
	return nil, nil
}

func (m *ManagerPassword) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	password := m.tfstate.String("password")
	err := rest.Managers.ManagerPasswordWithContext_PATCH(ctx, params{"password": password})
	return nil, err
}

func (m *ManagerPassword) UpdateResource(_ context.Context, plan UpdateResource, _ *VMSRest) (DisplayableRecord, error) {
	// password_md5 has RequiresReplace; Terraform should replace, not update.
	return nil, fmt.Errorf("manager password operations should be replaced, not updated")
}

func (m *ManagerPassword) DeleteResource(_ context.Context, _ *VMSRest) error {
	// Password cannot be "deleted" - this is a no-op
	return nil
}
