// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"net/http"

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
	tfstate *is.TFState
}

func (m *EventDefinitionConfig) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &EventDefinitionConfig{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: EventDefinitionConfigSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"id": rschema.Int64Attribute{
					Required:    true,
					Description: "Id of the event definition configuration",
				},
			},
		},
	)}
}

func (m *EventDefinitionConfig) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &EventDefinitionConfig{tfstate: is.NewTFStateMust(
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

func (m *EventDefinitionConfig) TfState() *is.TFState {
	return m.tfstate
}

func (m *EventDefinitionConfig) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.EventDefinitionConfigs
}

func (m *EventDefinitionConfig) DeleteResource(ctx context.Context, rest *VMSRest) error {
	// No-op for event definitions, as they cannot be deleted through this API.
	// Event definitions are managed by the system and can only be updated/patched.
	return nil
}
