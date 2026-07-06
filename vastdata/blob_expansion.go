// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/vast-data/go-vast-client/core"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
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

func (m *BlobExpansion) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &BlobExpansion{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: BlobExpansionSchemaRef,
		},
	)}
}

func (m *BlobExpansion) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &BlobExpansion{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: BlobExpansionSchemaRef,
		},
	)}
}

func (m *BlobExpansion) TfState() *is.TFState {
	return m.tfstate
}

func (m *BlobExpansion) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.BlobExpansions
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
	searchParams, _ := m.tfstate.SetIfAvailable("database_name", "table_name", "source_column_name", "tenant_id")
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
	// Idempotency: if the expansion already exists, return it as-is.
	existing, err := m.ReadDatasource(ctx, rest)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	createParams := m.tfstate.GetCreateParams()
	if _, err := core.Request[core.Record](ctx, rest.BlobExpansions, http.MethodPost, "/blobexpansions/", nil, createParams); err != nil {
		return nil, err
	}
	return m.ReadDatasource(ctx, rest)
}

func (m *BlobExpansion) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	planState := plan.(*BlobExpansion).tfstate

	// Build name-keyed sets from current state and plan so we can diff them.
	colKey := func(col map[string]any) string {
		if name, ok := col["name"].(string); ok {
			return name
		}
		return ""
	}

	currentCols := m.tfstate.ToSlice("arrow_schema")
	planCols := planState.ToSlice("arrow_schema")

	currentByName := make(map[string]map[string]any, len(currentCols))
	for _, c := range currentCols {
		if col, ok := c.(map[string]any); ok {
			currentByName[colKey(col)] = col
		}
	}
	planByName := make(map[string]map[string]any, len(planCols))
	for _, c := range planCols {
		if col, ok := c.(map[string]any); ok {
			planByName[colKey(col)] = col
		}
	}

	// Columns present in plan but missing from current state -> add.
	var toAdd []any
	for name, col := range planByName {
		if _, exists := currentByName[name]; !exists {
			toAdd = append(toAdd, col)
		}
	}

	// Columns present in current state but missing from plan -> drop.
	var toDrop []any
	for name, col := range currentByName {
		if _, exists := planByName[name]; !exists {
			toDrop = append(toDrop, col)
		}
	}

	// Build the hierarchy params shared by both operations.
	hierarchy, _ := m.tfstate.SetIfAvailable("database_name", "table_name", "source_column_name", "tenant_id")

	if len(toAdd) > 0 {
		body := core.Params{}
		for k, v := range hierarchy {
			body[k] = v
		}
		body["arrow_schema"] = toAdd
		if err := rest.BlobExpansions.BlobExpansionAddColumnsWithContext_PATCH(ctx, body); err != nil {
			return nil, err
		}
	}

	if len(toDrop) > 0 {
		body := core.Params{}
		for k, v := range hierarchy {
			body[k] = v
		}
		body["arrow_schema"] = toDrop
		if err := rest.BlobExpansions.BlobExpansionDropColumnsWithContext_PATCH(ctx, body); err != nil {
			return nil, err
		}
	}

	return m.ReadDatasource(ctx, rest)
}

func (m *BlobExpansion) DeleteResource(ctx context.Context, rest *VMSRest) error {
	deleteParams, _ := m.tfstate.SetIfAvailable("database_name", "table_name", "source_column_name", "tenant_id")
	err := rest.BlobExpansions.BlobExpansionDeleteWithContext_DELETE(ctx, deleteParams)
	return ignoreStatusCodes(err, http.StatusNotFound)
}
