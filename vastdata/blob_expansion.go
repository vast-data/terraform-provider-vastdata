// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	planmodifiers "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

type BlobExpansion struct {
	tfstate *is.TFState
}

func (m *BlobExpansion) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &BlobExpansion{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			PreserveUserValueFields: []string{"arrow_schema", "target_table_name"},
			TFStateHintsForCustom: &is.TFStateHintsForCustom{
				Description: "Manages a Blob Expansion resource on a VAST cluster. A Blob Expansion parses binary blobs (e.g., Kafka topic values) into structured columns in a target VastDB table.",
				SchemaAttributes: map[string]any{
					"database_name": rschema.StringAttribute{
						Required:    true,
						Description: "Kafka bucket name containing the source topic.",
						PlanModifiers: []planmodifiers.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"table_name": rschema.StringAttribute{
						Required:    true,
						Description: "Source Kafka topic name.",
						PlanModifiers: []planmodifiers.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"expansion_format": rschema.StringAttribute{
						Required:    true,
						Description: "Expansion format (e.g., 'protobuf', 'avro', 'json').",
						PlanModifiers: []planmodifiers.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"target_table_name": rschema.StringAttribute{
						Required:    true,
						Description: "Target table name for expanded data.",
						PlanModifiers: []planmodifiers.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"arrow_schema": rschema.ListNestedAttribute{
						Required:    true,
						Description: "Schema of the columns to parse from the blob. Changes are applied incrementally via add_columns/drop_columns.",
						NestedObject: rschema.NestedAttributeObject{
							Attributes: map[string]rschema.Attribute{
								"name": rschema.StringAttribute{
									Required:    true,
									Description: "Name of the column.",
								},
								"field": rschema.SingleNestedAttribute{
									Required:    true,
									Description: "Column type definition.",
									Attributes: map[string]rschema.Attribute{
										"column_type": rschema.StringAttribute{
											Optional:    true,
											Description: "Column type (e.g., 'string', 'int32', 'bool', 'map').",
										},
										"key_type": rschema.SingleNestedAttribute{
											Optional:    true,
											Description: "Key type for map columns.",
											Attributes: map[string]rschema.Attribute{
												"column_type": rschema.StringAttribute{
													Optional:    true,
													Description: "Column type for the map key (e.g., 'string').",
												},
											},
										},
										"value_type": rschema.SingleNestedAttribute{
											Optional:    true,
											Description: "Value type for map columns.",
											Attributes: map[string]rschema.Attribute{
												"column_type": rschema.StringAttribute{
													Optional:    true,
													Description: "Column type for the map value (e.g., 'string').",
												},
											},
										},
									},
								},
							},
						},
					},
					"source_column_name": rschema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Source column name (defaults to 'value' for Kafka topics).",
						PlanModifiers: []planmodifiers.String{
							stringplanmodifier.UseStateForUnknown(),
							stringplanmodifier.RequiresReplace(),
						},
					},
					"tenant_id": rschema.Int64Attribute{
						Optional:    true,
						Description: "Tenant ID (uses cluster default when unset).",
						PlanModifiers: []planmodifiers.Int64{
							int64planmodifier.RequiresReplace(),
						},
					},
					"copy_source_column": rschema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Whether to copy the raw source column to the target table.",
						PlanModifiers: []planmodifiers.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"add_excessive_values_output": rschema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Add an output column tracking source-blob fields not declared in the schema.",
						PlanModifiers: []planmodifiers.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"add_missing_values_output": rschema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Add an output column tracking schema fields missing from the source blob.",
						PlanModifiers: []planmodifiers.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"flatten_path": rschema.BoolAttribute{
						Optional:    true,
						Description: "Whether to flatten nested struct columns into separate columns.",
						PlanModifiers: []planmodifiers.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"flatten_delimiter": rschema.StringAttribute{
						Optional:    true,
						Description: "Delimiter for flattened path column names (default '__').",
						PlanModifiers: []planmodifiers.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"target_table_schema": rschema.StringAttribute{
						Optional:    true,
						Description: "Target table schema (defaults to the source table's schema).",
						PlanModifiers: []planmodifiers.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"columns": rschema.ListAttribute{
						Computed:    true,
						Description: "List of expanded column names (returned by the API).",
						ElementType: types.StringType,
					},
				},
			},
		},
	)}
}

func (m *BlobExpansion) TfState() *is.TFState {
	return m.tfstate
}

func (m *BlobExpansion) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.BlobExpansions
}

// identityParams returns the minimal set of params that uniquely identify this blob expansion.
func (m *BlobExpansion) identityParams() params {
	p := params{}
	if m.tfstate.IsKnownAndNotNull("database_name") {
		p["database_name"] = m.tfstate.String("database_name")
	}
	if m.tfstate.IsKnownAndNotNull("table_name") {
		p["table_name"] = m.tfstate.String("table_name")
	}
	if m.tfstate.IsKnownAndNotNull("source_column_name") {
		p["source_column_name"] = m.tfstate.String("source_column_name")
	}
	if m.tfstate.IsKnownAndNotNull("tenant_id") {
		p["tenant_id"] = m.tfstate.Int64("tenant_id")
	}
	return p
}

func (m *BlobExpansion) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	record, err := rest.BlobExpansions.BlobExpansionShowWithContext_GET(ctx, m.identityParams())
	if expectStatusCodes(err, http.StatusNotFound) {
		return nil, nil
	}
	// The API returns 500 with "Invalid blob expansion configuration" when no expansion exists.
	if isApiError(err) && strings.Contains(err.(*ApiError).Body, "Invalid blob expansion configuration") {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if len(record) == 0 {
		return nil, nil
	}

	// The API returns `columns` as a list of column-definition objects.
	// Extract just the column names for TF state (matches the ListAttribute{ElementType: StringType} schema).
	if rawCols, ok := record["columns"]; ok {
		if cols, ok := rawCols.([]any); ok {
			names := make([]any, 0, len(cols))
			for _, c := range cols {
				if cm, ok := c.(map[string]any); ok {
					if name, ok := cm["name"].(string); ok {
						names = append(names, name)
					}
				}
			}
			record["columns"] = names
		}
	}

	return record, nil
}

func (m *BlobExpansion) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	existing, err := m.ReadResource(ctx, rest)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	createParams := m.tfstate.GetCreateParams()
	if _, err := rest.BlobExpansions.CreateWithContext(ctx, createParams); err != nil {
		return nil, err
	}

	record, err := m.ReadResource(ctx, rest)
	if err != nil || record == nil {
		return record, err
	}

	// The API enables add_excessive_values_output and add_missing_values_output by
	// default regardless of what is sent in the create body. Drop them if the config
	// requests false to avoid drift on the next plan.
	dropBody := m.identityParams()
	dropBody["arrow_schema"] = []any{}
	needsDrop := false

	apiRecord, _ := record.(Record)
	if !m.tfstate.Bool("add_excessive_values_output") {
		if v, ok := apiRecord["add_excessive_values_output"].(bool); ok && v {
			dropBody["remove_excessive_values_output"] = true
			needsDrop = true
		}
	}
	if !m.tfstate.Bool("add_missing_values_output") {
		if v, ok := apiRecord["add_missing_values_output"].(bool); ok && v {
			dropBody["remove_missing_values_output"] = true
			needsDrop = true
		}
	}
	if !m.tfstate.Bool("copy_source_column") {
		if v, ok := apiRecord["copy_source_column"].(bool); ok && v {
			dropBody["remove_copy_source_column"] = true
			needsDrop = true
		}
	}

	if needsDrop {
		if err := rest.BlobExpansions.BlobExpansionDropColumnsWithContext_PATCH(ctx, dropBody); err != nil {
			return nil, err
		}
		return m.ReadResource(ctx, rest)
	}

	return record, nil
}

func (m *BlobExpansion) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	planTS := plan.(*BlobExpansion).TfState()
	identity := m.identityParams()

	currentSchemaItems := m.tfstate.ToSlice("arrow_schema")
	planSchemaItems := planTS.ToSlice("arrow_schema")

	currentByName := schemaItemsByName(currentSchemaItems)
	planByName := schemaItemsByName(planSchemaItems)

	var toAdd []any
	for name, item := range planByName {
		if _, exists := currentByName[name]; !exists {
			toAdd = append(toAdd, item)
		}
	}

	var toDrop []any
	for name, item := range currentByName {
		if _, exists := planByName[name]; !exists {
			toDrop = append(toDrop, item)
		}
	}

	if len(toAdd) > 0 {
		addBody := maps_clone(identity)
		addBody["arrow_schema"] = toAdd
		if err := rest.BlobExpansions.BlobExpansionAddColumnsWithContext_PATCH(ctx, addBody); err != nil {
			return nil, err
		}
	}

	if len(toDrop) > 0 {
		dropBody := maps_clone(identity)
		dropBody["arrow_schema"] = toDrop
		if err := rest.BlobExpansions.BlobExpansionDropColumnsWithContext_PATCH(ctx, dropBody); err != nil {
			return nil, err
		}
	}

	copySourcePlan := planTS.Bool("copy_source_column")
	copySourceCurrent := m.tfstate.Bool("copy_source_column")

	excessivePlan := planTS.Bool("add_excessive_values_output")
	excessiveCurrent := m.tfstate.Bool("add_excessive_values_output")

	missingPlan := planTS.Bool("add_missing_values_output")
	missingCurrent := m.tfstate.Bool("add_missing_values_output")

	needsAdd := (!copySourceCurrent && copySourcePlan) ||
		(!excessiveCurrent && excessivePlan) ||
		(!missingCurrent && missingPlan)

	needsDrop := (copySourceCurrent && !copySourcePlan) ||
		(excessiveCurrent && !excessivePlan) ||
		(missingCurrent && !missingPlan)

	if needsAdd {
		addBody := maps_clone(identity)
		// arrow_schema is required by the API even when only toggling metadata flags.
		if len(toAdd) > 0 {
			addBody["arrow_schema"] = toAdd
		} else {
			addBody["arrow_schema"] = []any{}
		}
		if !copySourceCurrent && copySourcePlan {
			addBody["add_copy_source_column"] = true
		}
		if !excessiveCurrent && excessivePlan {
			addBody["add_excessive_values_output"] = true
		}
		if !missingCurrent && missingPlan {
			addBody["add_missing_values_output"] = true
		}
		if err := rest.BlobExpansions.BlobExpansionAddColumnsWithContext_PATCH(ctx, addBody); err != nil {
			return nil, err
		}
	}

	if needsDrop {
		dropBody := maps_clone(identity)
		if len(toDrop) > 0 {
			dropBody["arrow_schema"] = toDrop
		} else {
			dropBody["arrow_schema"] = []any{}
		}
		if copySourceCurrent && !copySourcePlan {
			dropBody["remove_copy_source_column"] = true
		}
		if excessiveCurrent && !excessivePlan {
			dropBody["remove_excessive_values_output"] = true
		}
		if missingCurrent && !missingPlan {
			dropBody["remove_missing_values_output"] = true
		}
		if err := rest.BlobExpansions.BlobExpansionDropColumnsWithContext_PATCH(ctx, dropBody); err != nil {
			return nil, err
		}
	}

	return m.ReadResource(ctx, rest)
}

func (m *BlobExpansion) DeleteResource(ctx context.Context, rest *VMSRest) error {
	err := rest.BlobExpansions.BlobExpansionDeleteWithContext_DELETE(ctx, m.identityParams())
	return ignoreStatusCodes(err, http.StatusNotFound)
}

// schemaItemsByName indexes arrow_schema items by their "name" field.
func schemaItemsByName(items []any) map[string]any {
	result := make(map[string]any, len(items))
	for _, item := range items {
		if m, ok := item.(map[string]any); ok {
			if name, ok := m["name"].(string); ok {
				result[name] = item
			}
		}
	}
	return result
}

// maps_clone returns a shallow copy of a params map.
func maps_clone(src params) params {
	dst := make(params, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
