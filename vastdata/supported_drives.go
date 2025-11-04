// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var SupportedDriversSchemaRef = is.NewSchemaReference(
	http.MethodGet,
	"supporteddrives",
	http.MethodGet,
	"supporteddrives",
)

type SupportedDrives struct {
	tfstate *is.TFState
}

func (m *SupportedDrives) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &SupportedDrives{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			RequiredSchemaFields:    []string{"model_name", "drive_type"},
			NotOptionalSchemaFields: []string{"model", "type"},
			SchemaRef:               SupportedDriversSchemaRef,
		},
	)}
}

func (m *SupportedDrives) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &SupportedDrives{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SearchableFields:        []string{"model_name", "drive_type"},
			NotOptionalSchemaFields: []string{"model", "type"},
			SchemaRef:               SupportedDriversSchemaRef,
		},
	)}
}

func (m *SupportedDrives) TfState() *is.TFState {
	return m.tfstate
}

func (m *SupportedDrives) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.SupportedDrivers
}

// transformToBulkFormat converts the state data to the bulk POST format required by the API.
// User configures: model_name, drive_type (required), raw_data (but model and type are not optional in raw_data)
// POST format: {replace: false, drives: [{...raw_data, type: drive_type, model: model_name}]}
func (m *SupportedDrives) transformToBulkFormat(stateParams map[string]any) map[string]any {
	bulkParams := make(map[string]any)
	bulkParams["replace"] = false

	var drive map[string]any
	if rawData, ok := stateParams["raw_data"].(map[string]any); ok {
		drive = make(map[string]any)
		for k, v := range rawData {
			drive[k] = v
		}
	} else {
		drive = make(map[string]any)
	}

	// Map top-level fields to raw_data fields
	// drive_type -> raw_data.type
	if driveType, ok := stateParams["drive_type"]; ok {
		drive["type"] = driveType
	}
	// model_name -> raw_data.model
	if modelName, ok := stateParams["model_name"]; ok {
		drive["model"] = modelName
	}

	// Wrap drive in drives array
	if len(drive) > 0 {
		bulkParams["drives"] = []any{drive}
	} else {
		bulkParams["drives"] = []any{}
	}

	return bulkParams
}

func (m *SupportedDrives) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	createParams := m.tfstate.GetCreateParams()
	bulkParams := m.transformToBulkFormat(createParams)
	return rest.SupportedDrivers.CreateWithContext(ctx, bulkParams)
}

func (m *SupportedDrives) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	planState := plan.(*SupportedDrives).TfState()
	updateParams := planState.GetAllValues()
	bulkParams := m.transformToBulkFormat(updateParams)
	return rest.SupportedDrivers.CreateWithContext(ctx, bulkParams)
}
