// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var EventDefinitionSchemaRef = is.NewSchemaReference(
	http.MethodPatch,
	"eventdefinitions/{id}",
	http.MethodGet,
	"eventdefinitions",
)

type EventDefinition struct {
	tfstate *is.TFState
}

func (m *EventDefinition) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &EventDefinition{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: EventDefinitionSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"id": rschema.Int64Attribute{
					Optional:    true,
					Computed:    true,
					Description: "Unique identifier for the event definition",
				},
				"name": rschema.StringAttribute{
					Required:    true,
					Description: "Name of the event definition",
				},
			},
		},
	)}
}

func (m *EventDefinition) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &EventDefinition{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: EventDefinitionSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"name": dschema.StringAttribute{
					Required:    true,
					Description: "Name of the event definition",
				},
			},
		},
	)}
}

func (m *EventDefinition) TfState() *is.TFState {
	return m.tfstate
}

func (m *EventDefinition) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.EventDefinitions
}

func (m *EventDefinition) DeleteResource(ctx context.Context, rest *VMSRest) error {
	// No-op for event definitions, as they cannot be deleted through this API.
	// Event definitions are managed by the system and can only be updated/patched.
	return nil
}
