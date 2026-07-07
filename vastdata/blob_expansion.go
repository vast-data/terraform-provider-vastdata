// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/vast-data/go-vast-client/core"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/schema_generation"
)

var BlobExpansionSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"blobexpansions",
	http.MethodGet,
	"blobexpansions/show",
)

type BlobExpansion struct {
	tfstate *is.TFState
}

func blobExpansionHints() *is.TFStateHints {
	return &is.TFStateHints{
		SchemaRef: BlobExpansionSchemaRef,
		CommonModifiersMapping: map[string]string{
			"database_name":       schema_generation.ModifierForceNew,
			"table_name":          schema_generation.ModifierForceNew,
			"source_column_name":  schema_generation.ModifierForceNew,
			"target_table_name":   schema_generation.ModifierForceNew,
			"expansion_format":    schema_generation.ModifierForceNew,
			"target_table_schema": schema_generation.ModifierForceNew,
			"flatten_path":        schema_generation.ModifierForceNew,
			"flatten_delimiter":   schema_generation.ModifierForceNew,
		},
		ExcludedSchemaFields:    []string{"columns"},
		PreserveOrderFields:     []string{"arrow_schema"},
		PreserveUserValueFields: []string{"arrow_schema"},
		ReadOnlyFields:          []string{"tenant_id"},
	}
}

func (m *BlobExpansion) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &BlobExpansion{tfstate: is.NewTFStateMust(raw, schema, blobExpansionHints())}
}

func (m *BlobExpansion) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &BlobExpansion{tfstate: is.NewTFStateMust(raw, schema, blobExpansionHints())}
}

func (m *BlobExpansion) TfState() *is.TFState {
	return m.tfstate
}

func (m *BlobExpansion) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.BlobExpansions
}

func (m *BlobExpansion) blobExpansionSearchParams() params {
	searchParams, _ := m.tfstate.SetIfAvailable("database_name", "table_name", "source_column_name", "tenant_id")
	return searchParams
}

// normalizeBlobExpansionRecord converts the API's fully-qualified
// target_table_name path (e.g. "/default/db/schema/table") to just the bare
// table name so that the TF state stays consistent with what the user
// configured.
func normalizeBlobExpansionRecord(record DisplayableRecord) DisplayableRecord {
	if record == nil {
		return nil
	}
	rec, ok := record.(core.Record)
	if !ok {
		return record
	}
	if full, ok := rec["target_table_name"].(string); ok && strings.Contains(full, "/") {
		parts := strings.Split(full, "/")
		rec["target_table_name"] = parts[len(parts)-1]
	}
	return rec
}

func (m *BlobExpansion) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	searchParams := m.blobExpansionSearchParams()
	record, err := rest.BlobExpansions.BlobExpansionShowWithContext_GET(ctx, searchParams)
	if isApiError(err) {
		if strings.Contains(err.(*ApiError).Body, "Invalid blob expansion configuration") {
			return nil, nil
		}
	}
	if expectStatusCodes(err, http.StatusNotFound) {
		return nil, nil
	}
	return normalizeBlobExpansionRecord(record), err
}

func (m *BlobExpansion) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.ReadDatasource(ctx, rest)
}

func (m *BlobExpansion) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	record, err := m.ReadDatasource(ctx, rest)
	if err != nil {
		return nil, err
	}
	if record != nil {
		return record, nil
	}
	createParams := m.tfstate.GetCreateParams()
	if _, err := rest.BlobExpansions.CreateWithContext(ctx, createParams); err != nil {
		return nil, err
	}
	return m.ReadDatasource(ctx, rest)
}

func (m *BlobExpansion) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	stateTs := m.tfstate
	planTs := plan.(*BlobExpansion).tfstate

	hierarchy, _ := stateTs.SetIfAvailable("database_name", "table_name", "source_column_name", "tenant_id")

	stateSchema := arrowSchemaFromState(stateTs)
	planSchema := arrowSchemaFromState(planTs)
	added, removed := diffArrowSchema(stateSchema, planSchema)

	if len(added) > 0 {
		body := copyParams(hierarchy)
		body["arrow_schema"] = added
		applyBlobExpansionAddFlags(body, stateTs, planTs)
		if err := rest.BlobExpansions.BlobExpansionAddColumnsWithContext_PATCH(ctx, body); err != nil {
			return nil, err
		}
	}
	if len(removed) > 0 {
		body := copyParams(hierarchy)
		body["arrow_schema"] = removed
		applyBlobExpansionDropFlags(body, stateTs, planTs)
		if err := rest.BlobExpansions.BlobExpansionDropColumnsWithContext_PATCH(ctx, body); err != nil {
			return nil, err
		}
	}
	if len(added) == 0 && len(removed) == 0 {
		if err := patchBlobExpansionFlags(ctx, rest, hierarchy, stateTs, planTs); err != nil {
			return nil, err
		}
	}

	return m.ReadDatasource(ctx, rest)
}

func (m *BlobExpansion) DeleteResource(ctx context.Context, rest *VMSRest) error {
	deleteParams := m.blobExpansionSearchParams()
	err := rest.BlobExpansions.BlobExpansionDeleteWithContext_DELETE(ctx, deleteParams)
	return ignoreStatusCodes(err, http.StatusNotFound)
}

func arrowSchemaFromState(ts *is.TFState) []any {
	raw, ok := ts.GetAllValues()["arrow_schema"].([]any)
	if !ok {
		return nil
	}
	return raw
}

func diffArrowSchema(state, plan []any) (added, removed []any) {
	stateByName := indexArrowSchemaByName(state)
	planByName := indexArrowSchemaByName(plan)

	for name, col := range planByName {
		if _, ok := stateByName[name]; !ok {
			added = append(added, col)
		}
	}
	for name, col := range stateByName {
		if _, ok := planByName[name]; !ok {
			removed = append(removed, col)
		}
	}
	return added, removed
}

func indexArrowSchemaByName(cols []any) map[string]any {
	out := make(map[string]any, len(cols))
	for _, col := range cols {
		m, ok := col.(map[string]any)
		if !ok {
			continue
		}
		name, _ := m["name"].(string)
		if name != "" {
			out[name] = col
		}
	}
	return out
}

func copyParams(in params) params {
	out := make(params, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func applyBlobExpansionAddFlags(body params, stateTs, planTs *is.TFState) {
	if boolChanged(stateTs, planTs, "copy_source_column") && boolFromState(planTs, "copy_source_column") {
		body["add_copy_source_column"] = true
	}
	if boolChanged(stateTs, planTs, "add_missing_values_output") && boolFromState(planTs, "add_missing_values_output") {
		body["add_missing_values_output"] = true
	}
	if boolChanged(stateTs, planTs, "add_excessive_values_output") && boolFromState(planTs, "add_excessive_values_output") {
		body["add_excessive_values_output"] = true
	}
}

func applyBlobExpansionDropFlags(body params, stateTs, planTs *is.TFState) {
	if boolChanged(stateTs, planTs, "copy_source_column") && !boolFromState(planTs, "copy_source_column") {
		body["remove_copy_source_column"] = true
	}
	if boolChanged(stateTs, planTs, "add_missing_values_output") && !boolFromState(planTs, "add_missing_values_output") {
		body["remove_missing_values_output"] = true
	}
	if boolChanged(stateTs, planTs, "add_excessive_values_output") && !boolFromState(planTs, "add_excessive_values_output") {
		body["remove_excessive_values_output"] = true
	}
}

func patchBlobExpansionFlags(ctx context.Context, rest *VMSRest, hierarchy params, stateTs, planTs *is.TFState) error {
	needsAdd := (boolChanged(stateTs, planTs, "copy_source_column") && boolFromState(planTs, "copy_source_column")) ||
		(boolChanged(stateTs, planTs, "add_missing_values_output") && boolFromState(planTs, "add_missing_values_output")) ||
		(boolChanged(stateTs, planTs, "add_excessive_values_output") && boolFromState(planTs, "add_excessive_values_output"))
	needsDrop := (boolChanged(stateTs, planTs, "copy_source_column") && !boolFromState(planTs, "copy_source_column")) ||
		(boolChanged(stateTs, planTs, "add_missing_values_output") && !boolFromState(planTs, "add_missing_values_output")) ||
		(boolChanged(stateTs, planTs, "add_excessive_values_output") && !boolFromState(planTs, "add_excessive_values_output"))

	if needsAdd {
		body := copyParams(hierarchy)
		applyBlobExpansionAddFlags(body, stateTs, planTs)
		if err := rest.BlobExpansions.BlobExpansionAddColumnsWithContext_PATCH(ctx, body); err != nil {
			return err
		}
	}
	if needsDrop {
		body := copyParams(hierarchy)
		applyBlobExpansionDropFlags(body, stateTs, planTs)
		if err := rest.BlobExpansions.BlobExpansionDropColumnsWithContext_PATCH(ctx, body); err != nil {
			return err
		}
	}
	return nil
}

func boolFromState(ts *is.TFState, boolField string) bool {
	v, ok := ts.GetAllValues()[boolField].(bool)
	return ok && v
}

func boolChanged(stateTs, planTs *is.TFState, boolField string) bool {
	return boolFromState(stateTs, boolField) != boolFromState(planTs, boolField)
}
