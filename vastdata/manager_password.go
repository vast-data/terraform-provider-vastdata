// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/schema_generation"
)

var ManagerPasswordSchemaRef = is.NewSchemaReference(
	http.MethodPatch,
	"managers/password",
	"",
	"",
)

type ManagerPassword struct {
	tfstate *is.TFState
}

func (m *ManagerPassword) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &ManagerPassword{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SensitiveFields: []string{"password"},
			SchemaRef:       ManagerPasswordSchemaRef,
			CommonModifiersMapping: map[string]string{
				"password": schema_generation.ModifierForceNew,
			},
		},
	)}
}
func (m *ManagerPassword) TfState() *is.TFState {
	return m.tfstate
}

func (m *ManagerPassword) API(_ *VMSRest) VastResourceAPIWithContext {
	return nil
}

func (m *ManagerPassword) ReadResource(_ context.Context, _ *VMSRest) (DisplayableRecord, error) {
	return nil, nil
}

func (m *ManagerPassword) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	password := m.tfstate.String("password")
	err := rest.Managers.ManagerPasswordWithContext_PATCH(ctx, password)
	return nil, err
}

func (m *ManagerPassword) UpdateResource(_ context.Context, plan UpdateResource, _ *VMSRest) (DisplayableRecord, error) {
	// With force_new modifiers, Terraform will handle replacements automatically
	// This method should not be called for updates since all fields have RequiresReplace()
	// But we'll keep it as a safety net in case it's called
	return nil, fmt.Errorf("manager password operations should be replaced, not updated")
}

func (m *ManagerPassword) DeleteResource(_ context.Context, _ *VMSRest) error {
	// Password cannot be "deleted" - this is a no-op
	return nil
}
