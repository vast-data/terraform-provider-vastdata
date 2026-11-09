// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var EventDefinitionConfigSchemaRef = is.NewSchemaReference(
	http.MethodPatch,
	"eventdefinitionconfigs/{id}",
	http.MethodGet,
	"eventdefinitionconfigs",
)

type EventDefinitionConfig struct {
	tfstate is.ResourceState
}

func (m *EventDefinitionConfig) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &EventDefinitionConfig{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:          EventDefinitionConfigSchemaRef,
			WriteOnlyFields:    []string{"smtp_password"},
			WriteOnlyMD5Fields: []string{"smtp_password"},
			SensitiveFields:    []string{"smtp_password"},
			AdditionalSchemaAttributes: map[string]any{
				"id": rschema.Int64Attribute{
					Required:    true,
					Description: "Id of the event definition configuration",
				},
				"smtp_password_md5": ContentMD5SchemaAttribute("smtp_password", false),
			},
		},
	)}
}

func (m *EventDefinitionConfig) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &EventDefinitionConfig{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: EventDefinitionConfigSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"id": dschema.Int64Attribute{
					Required:    true,
					Description: "Id of the event definition configuration",
				},
			},
		},
	)}
}

func (m *EventDefinitionConfig) TfState() is.ResourceState {
	return m.tfstate
}

func (m *EventDefinitionConfig) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.EventDefinitionConfigs
}

func (m *EventDefinitionConfig) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) error {
	SetPlanContentMD5(ctx, req, resp, "smtp_password", "smtp_password_md5")
	return nil
}

func (m *EventDefinitionConfig) TransformResponseRecord(record Record) Record {
	DeleteRecordFields(record, "smtp_password")
	return record
}

func (m *EventDefinitionConfig) DeleteResource(ctx context.Context, rest *VMSRest) error {
	// No-op for event definitions, as they cannot be deleted through this API.
	// Event definitions are managed by the system and can only be updated/patched.
	return nil
}
