// Copyright (c) HashiCorp, Inc.

// This file implements dynamic data-source schema generation from OpenAPI schemas.
// It builds Terraform schema definitions by analyzing request/response models
// from the VAST API (via kin-openapi), supporting fallback logic and validators.

package schema_generation

import (
	"context"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/client"
)

var excludeSearchParams = []string{"page", "page_size", "sync", "created", "sync_time"}

func GetDatasourceSchema(ctx context.Context, hints *TFStateHints) (*dschema.Schema, error) {

	if hints.TFStateHintsForCustom != nil {
		// Build schema for custom datasource
		return getDatasourceSchemaForCustom(hints)
	}

	if hints.SchemaRef == nil {
		return nil, fmt.Errorf("schema reference is required but was nil")
	}
	if hints.SchemaRef.Read == nil {
		return nil, fmt.Errorf("read schema reference is required but was nil")
	}

	resourcePath := hints.SchemaRef.Read.Path
	resourceMethod := hints.SchemaRef.Read.Method
	if resourcePath == "" {
		return nil, fmt.Errorf("resource path is required but was empty")
	}
	if resourceMethod == "" {
		return nil, fmt.Errorf("resource method is required but was empty")
	}

	var (
		readSchemaRef *openapi3.SchemaRef // Will be computed fields (Response schema definition)
		err           error
	)

	switch resourceMethod {
	case http.MethodGet:
		if readSchemaRef, err = client.GetSchema_GET_StatusOk(resourcePath); err != nil {
			return nil, fmt.Errorf("failed to get schema for GET %q: %w", resourcePath, err)
		}
	default:
		return nil, fmt.Errorf(
			"not supported resource method %q for datasource schema generation", resourceMethod,
		)
	}

	// Will be write-only fields (Query parameters) unless present in Response schema
	params, err := client.QueryParametersGET(resourcePath)
	if err != nil {
		return nil, err
	}

	// Extract main schema attributes
	allProps := map[string]*SchemaEntry{}

	// First, add query parameters with correct required status
	requiredParams := []string{}
	paramSchemas := map[string]*openapi3.SchemaRef{}

	for _, p := range params {
		if !isPrimitive(p.Schema.Value) {
			// We search only for primitive types
			continue
		}
		name := p.Name
		if contains(excludeSearchParams, name) {
			continue
		}
		if p.Schema == nil || p.Schema.Value == nil || p.Schema.Value.Type == nil || len(*p.Schema.Value.Type) == 0 {
			continue
		}

		paramSchemas[name] = buildTmpSchemaRefFromParam(p)

		// Check if parameter is required
		if p.Required {
			requiredParams = append(requiredParams, name)
		}
	}

	// Add query parameters with proper required status
	addSchemaEntries(paramSchemas, requiredParams, hints, allProps, false, true, false, false, false, false)

	// Then, add response schema properties as computed fields (ignore response required fields)
	if readSchemaRef.Value != nil {
		addSchemaEntries(readSchemaRef.Value.Properties, []string{}, hints, allProps, false, true, true, false, false, false)
	}

	attrs := buildDatasourceAttributesFromMap(ctx, allProps, hints)
	if hints.AdditionalSchemaAttributes != nil {
		for k, v := range hints.AdditionalSchemaAttributes {
			if att, ok := v.(dschema.Attribute); ok {
				attrs[k] = att
			} else {
				return nil, fmt.Errorf("additional schema attribute %q is not a valid dschema.Attribute (got %T)", k, v)
			}
		}
	}

	// Description fallback
	var description, summary string
	if readSchemaRef.Value != nil {
		summary = readSchemaRef.Value.Title
		description = readSchemaRef.Value.Description
	}
	if summary == "" {
		summary = description
	}

	return &dschema.Schema{
		Description:         summary,
		MarkdownDescription: summary,
		Attributes:          attrs,
	}, nil
}

func getDatasourceSchemaForCustom(hints *TFStateHints) (*dschema.Schema, error) {
	customHints := hints.TFStateHintsForCustom
	if customHints.SchemaAttributes == nil || len(customHints.SchemaAttributes) == 0 {
		return nil, fmt.Errorf("custom datasource schema attributes are required but were empty")
	}

	attrs := make(map[string]dschema.Attribute)
	for k, v := range customHints.SchemaAttributes {
		if att, ok := v.(dschema.Attribute); ok {
			attrs[k] = att
		} else {
			return nil, fmt.Errorf("additional schema attribute %q is not a valid dschema.Attribute (got %T)", k, v)
		}
	}

	var description, markdownDescription string
	description = customHints.Description
	if customHints.MarkdownDescription != "" {
		markdownDescription = customHints.MarkdownDescription
	} else {
		markdownDescription = description
	}

	return &dschema.Schema{
		Description:         description,
		MarkdownDescription: markdownDescription,
		Attributes:          attrs,
	}, nil

}

func buildDatasourceAttributesFromMap(
	ctx context.Context,
	entries map[string]*SchemaEntry,
	hints *TFStateHints,
) map[string]dschema.Attribute {
	result := make(map[string]dschema.Attribute, len(entries))
	for name, entry := range entries {
		schema := resolveComposedSchema(resolveAllRefs(&openapi3.SchemaRef{Value: entry.Prop}))

		if !isPrimitive(schema) {
			if !entry.Computed {
				// Skip non-primitive fields if they are not computed.
				//These fields are not part of result object schema and used only for request.
				continue
			}
			// We're not going to create fields in datasource so prohibit searching by non-primitive fields
			entry.Optional = false
		}

		attr := buildDatasourceAttribute(ctx, name, schema, entry, hints)
		if attr == nil {
			continue
		}
		result[name] = attr
	}
	return result
}

func buildDatasourceAttribute(
	ctx context.Context,
	name string,
	schema *openapi3.Schema,
	entry *SchemaEntry,
	hints *TFStateHints,
) dschema.Attribute {
	schema = resolveComposedSchema(schema)
	if schema == nil || schema.Type == nil || len(*schema.Type) == 0 {
		panic(fmt.Sprintf("missing or invalid schema type for attribute %q", name))
	}

	desc := entry.Description

	switch (*schema.Type)[0] {
	case openapi3.TypeString:
		return dschema.StringAttribute{
			Required:            entry.Required,
			Optional:            entry.Optional,
			Computed:            entry.Computed,
			Sensitive:           entry.Sensitive,
			Description:         desc,
			MarkdownDescription: desc,
		}

	case openapi3.TypeInteger:
		return dschema.Int64Attribute{
			Required:            entry.Required,
			Optional:            entry.Optional,
			Computed:            entry.Computed,
			Sensitive:           entry.Sensitive,
			Description:         desc,
			MarkdownDescription: desc,
		}

	case openapi3.TypeNumber:
		return dschema.Float64Attribute{
			Required:            entry.Required,
			Optional:            entry.Optional,
			Computed:            entry.Computed,
			Sensitive:           entry.Sensitive,
			Description:         desc,
			MarkdownDescription: desc,
		}

	case openapi3.TypeBoolean:
		return dschema.BoolAttribute{
			Required:            entry.Required,
			Optional:            entry.Optional,
			Computed:            entry.Computed,
			Sensitive:           entry.Sensitive,
			Description:         desc,
			MarkdownDescription: desc,
		}

	case openapi3.TypeArray:
		itemSchema := resolveComposedSchema(resolveAllRefs(schema.Items))
		if itemSchema == nil || itemSchema.Type == nil || len(*itemSchema.Type) == 0 {
			panic("invalid item schema for array")
		}

		if (*itemSchema.Type)[0] == openapi3.TypeObject && len(itemSchema.Properties) == 0 {
			return nil // skip empty nested object
		}

		isOrdered := entry.Ordered

		switch (*itemSchema.Type)[0] {
		case openapi3.TypeObject:
			nested := make(map[string]*SchemaEntry)
			addSchemaEntries(itemSchema.Properties, itemSchema.Required, hints, nested, false, false, true, false, false, entry.Ordered)
			attrs := buildDatasourceAttributesFromMap(ctx, nested, hints)

			if isOrdered {
				return dschema.ListNestedAttribute{
					NestedObject: dschema.NestedAttributeObject{
						Attributes: attrs,
					},
					Computed:            true,
					Description:         entry.Description,
					MarkdownDescription: entry.Description,
				}
			}
			return dschema.SetNestedAttribute{
				NestedObject: dschema.NestedAttributeObject{
					Attributes: attrs,
				},
				Computed:            true,
				Description:         entry.Description,
				MarkdownDescription: entry.Description,
			}

		case openapi3.TypeArray:
			inner := resolveComposedSchema(resolveAllRefs(itemSchema.Items))
			innerType := buildAttrTypeFromSchema(inner)

			if isOrdered {
				return dschema.ListAttribute{
					ElementType:         types.ListType{ElemType: innerType},
					Computed:            true,
					Description:         entry.Description,
					MarkdownDescription: entry.Description,
				}
			}
			return dschema.SetAttribute{
				ElementType:         types.SetType{ElemType: innerType},
				Computed:            true,
				Description:         entry.Description,
				MarkdownDescription: entry.Description,
			}

		default:
			elemType := buildAttrTypeFromSchema(itemSchema)
			if isOrdered {
				return dschema.ListAttribute{
					ElementType:         elemType,
					Computed:            true,
					Description:         entry.Description,
					MarkdownDescription: entry.Description,
				}
			}
			return dschema.SetAttribute{
				ElementType:         elemType,
				Computed:            true,
				Description:         entry.Description,
				MarkdownDescription: entry.Description,
			}
		}

	case openapi3.TypeObject:
		if schema.AdditionalProperties.Schema != nil {
			valSchema := resolveComposedSchema(resolveAllRefs(schema.AdditionalProperties.Schema))
			if valSchema.Type != nil && (*valSchema.Type)[0] == openapi3.TypeObject && len(valSchema.Properties) == 0 {
				warnWithContext(ctx, fmt.Sprintf("Skipping. Map value schema for %q is empty object", name))
				return nil
			}

			valType := buildAttrTypeFromSchema(valSchema)

			if valSchema.Type != nil && (*valSchema.Type)[0] == openapi3.TypeObject && len(valSchema.Properties) > 0 {
				nested := make(map[string]*SchemaEntry)
				addSchemaEntries(valSchema.Properties, valSchema.Required, hints, nested, entry.Required, entry.Optional, entry.Computed, entry.WriteOnly, entry.Sensitive, entry.Ordered)

				return dschema.MapNestedAttribute{
					NestedObject: dschema.NestedAttributeObject{
						Attributes: buildDatasourceAttributesFromMap(ctx, nested, hints),
					},
					Required:            entry.Required,
					Optional:            entry.Optional,
					Computed:            entry.Computed,
					Sensitive:           entry.Sensitive,
					Description:         desc,
					MarkdownDescription: desc,
				}
			}

			return dschema.MapAttribute{
				ElementType:         valType,
				Required:            entry.Required,
				Optional:            entry.Optional,
				Computed:            entry.Computed,
				Sensitive:           entry.Sensitive,
				Description:         desc,
				MarkdownDescription: desc,
			}
		}

		if len(schema.Properties) > 0 {
			nested := make(map[string]*SchemaEntry)
			addSchemaEntries(schema.Properties, schema.Required, hints, nested, entry.Required, entry.Optional, entry.Computed, entry.WriteOnly, entry.Sensitive, entry.Ordered)

			return dschema.SingleNestedAttribute{
				Attributes:          buildDatasourceAttributesFromMap(ctx, nested, hints),
				Required:            entry.Required,
				Optional:            entry.Optional,
				Computed:            entry.Computed,
				Sensitive:           entry.Sensitive,
				Description:         desc,
				MarkdownDescription: desc,
			}
		}
		warnWithContext(ctx, fmt.Sprintf("Skipping. Object attribute %q has no properties or additionalProperties", name))
		return nil

	default:
		panic(fmt.Sprintf("unsupported schema type %q for attribute %q", (*schema.Type)[0], name))
	}
}
