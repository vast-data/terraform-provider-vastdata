// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"maps"
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
			"database_name":             schema_generation.ModifierForceNew,
			"table_name":                schema_generation.ModifierForceNew,
			"source_column_name":        schema_generation.ModifierForceNew,
			"target_table_name":         schema_generation.ModifierForceNew,
			"expansion_format":          schema_generation.ModifierForceNew,
			"target_table_schema":       schema_generation.ModifierForceNew,
			"flatten_path":              schema_generation.ModifierForceNew,
			"flatten_delimiter":         schema_generation.ModifierForceNew,
			"add_missing_values_output":   schema_generation.ModifierForceNew,
			"add_excessive_values_output": schema_generation.ModifierForceNew,
		},
		ExcludedSchemaFields:    []string{"columns"},
		PreserveOrderFields:     []string{"arrow_schema"},
		PreserveUserValueFields: []string{"arrow_schema"},
		ReadOnlyFields:          []string{"tenant_id"},
	}
}

func blobExpansionDatasourceHints() *is.TFStateHints {
	return &is.TFStateHints{
		SchemaRef:            BlobExpansionSchemaRef,
		ExcludedSchemaFields: []string{"columns"},
	}
}

func (m *BlobExpansion) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &BlobExpansion{tfstate: is.NewTFStateMust(raw, schema, blobExpansionHints())}
}

func (m *BlobExpansion) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &BlobExpansion{tfstate: is.NewTFStateMust(raw, schema, blobExpansionDatasourceHints())}
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
	if _, err := core.Request[core.Record](ctx, rest.BlobExpansions, http.MethodPost, "/blobexpansions/", nil, createParams); err != nil {
		return nil, err
	}
	return m.ReadDatasource(ctx, rest)
}

func (m *BlobExpansion) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	stateTs := m.tfstate
	planTs := plan.(*BlobExpansion).tfstate

	hierarchy, _ := stateTs.SetIfAvailable("database_name", "table_name", "source_column_name", "tenant_id")

	stateSchema := stateTs.ToSlice("arrow_schema")
	planSchema := planTs.ToSlice("arrow_schema")
	added, removed := diffArrowSchema(stateSchema, planSchema)

	addBody := maps.Clone(hierarchy)
	dropBody := maps.Clone(hierarchy)
	hasAddFlags, hasDropFlags := false, false
	for _, spec := range blobExpansionFlags {
		if !boolChanged(stateTs, planTs, spec.field) {
			continue
		}
		if boolFromState(planTs, spec.field) {
			addBody[spec.addFlag] = true
			hasAddFlags = true
		} else {
			dropBody[spec.dropFlag] = true
			hasDropFlags = true
		}
	}
	if len(added) > 0 {
		addBody["arrow_schema"] = added
	} else if hasAddFlags {
		addBody["arrow_schema"] = []any{}
	}
	if len(removed) > 0 {
		dropBody["arrow_schema"] = removed
	} else if hasDropFlags {
		dropBody["arrow_schema"] = []any{}
	}
	if len(added) > 0 || hasAddFlags {
		if err := rest.BlobExpansions.BlobExpansionAddColumnsWithContext_PATCH(ctx, addBody); err != nil {
			return nil, err
		}
	}
	if len(removed) > 0 || hasDropFlags {
		if err := rest.BlobExpansions.BlobExpansionDropColumnsWithContext_PATCH(ctx, dropBody); err != nil {
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

var blobExpansionFlags = []struct{ field, addFlag, dropFlag string }{
	{"copy_source_column", "add_copy_source_column", "remove_copy_source_column"},
}

func boolFromState(ts *is.TFState, boolField string) bool {
	v, ok := ts.GetAllValues()[boolField].(bool)
	return ok && v
}

func boolChanged(stateTs, planTs *is.TFState, boolField string) bool {
	return boolFromState(stateTs, boolField) != boolFromState(planTs, boolField)
}
