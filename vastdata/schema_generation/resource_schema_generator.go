// Copyright (c) HashiCorp, Inc.

// This file implements dynamic resource schema generation from OpenAPI schemas.
// It builds Terraform schema definitions by analyzing request/response models
// from the VAST API (via kin-openapi), supporting fallback logic and validators.

package schema_generation

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/client"
)

func GetResourceSchema(ctx context.Context, hints *TFStateHints) (*rschema.Schema, error) {

	if hints.TFStateHintsForCustom != nil {
		// Build schema for custom resource
		return getResourceSchemaForCustom(hints)
	}

	if hints.SchemaRef == nil {
		return nil, fmt.Errorf("schema reference is required but was nil")
	}
	if hints.SchemaRef.Create == nil {
		return nil, fmt.Errorf("create schema reference is required but was nil")
	}

	resourcePath := hints.SchemaRef.Create.Path
	resourceMethod := hints.SchemaRef.Create.Method
	if resourcePath == "" {
		return nil, fmt.Errorf("resource path is required but was empty")
	}
	if resourceMethod == "" {
		return nil, fmt.Errorf("resource method is required but was empty")
	}

	var (
		createSchemaRef *openapi3.SchemaRef // Schema for request body (Create json body)
		modelSchemaRef  *openapi3.SchemaRef // Schema for response body (Model json body)
		err             error
	)

	// Resolve createSchemaRef
	switch resourceMethod {
	case http.MethodPost:
		if createSchemaRef, err = client.GetSchema_POST_RequestBody(resourcePath); err != nil {
			return nil, fmt.Errorf("failed to get POST schema for resource %q: %w", resourcePath, err)
		}
	case http.MethodGet:
		if createSchemaRef, err = client.GetSchema_GET_StatusOk(resourcePath); err != nil {
			return nil, fmt.Errorf("failed to get GET schema for resource %q: %w", resourcePath, err)
		}
	case http.MethodPatch:
		if createSchemaRef, err = client.GetSchema_PATCH_RequestBody(resourcePath); err != nil {
			return nil, fmt.Errorf("failed to get PATCH schema for resource %q: %w", resourcePath, err)
		}
	default:
		return nil, fmt.Errorf(
			"unsupported method %q for resource %q (CreateSchema)", resourceMethod, resourcePath,
		)
	}

	// Resolve modelSchemaRef
	switch resourceMethod {
	case http.MethodPost:
		if modelSchemaRef, err = client.GetSchema_POST_StatusOk(resourcePath); err != nil {
			return nil, fmt.Errorf("failed to get POST model schema for resource %q: %w", resourcePath, err)
		}
	case http.MethodPatch:
		if modelSchemaRef, err = client.GetSchema_PATCH_RequestBody(resourcePath); err != nil {
			return nil, fmt.Errorf("failed to get patch model schema for resource %q: %w", resourcePath, err)
		}

	default:
		return nil, fmt.Errorf(
			"unsupported method %q for resource %q (ModelSchema)", resourceMethod, resourcePath,
		)
	}

	if IsEmptySchema(createSchemaRef) {
		return nil, fmt.Errorf("create schema is empty for resource %q", resourcePath)
	}
	if modelSchemaRef == nil {
		return nil, fmt.Errorf("model schema reference is nil for resource path %q", resourcePath)
	}

	//createSchemaRef, err := client.GetSchema_POST_RequestBody(resourcePath)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get POST request schema for resource %q: %w", resourcePath, err)
	//}
	//
	//if IsEmptySchema(createSchemaRef) {
	//	if createSchemaRef, err = client.GetSchema_GET_StatusOk(resourcePath); err != nil {
	//		return nil, fmt.Errorf("failed to get GET status schema for resource %q: %w", resourcePath, err)
	//	}
	//	infoWithContext(
	//		ctx,
	//		fmt.Sprintf("POST schema is not present for resource %q, using GET status schema instead", resourcePath),
	//	)
	//}

	//modelSchemaRef, err := client.GetSchema_POST_StatusOk(resourcePath)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get POST status schema: %w", err)
	//}
	//if modelSchemaRef == nil {
	//	return nil, fmt.Errorf("model schema reference is nil for resource path %q", resourcePath)
	//}
	//if IsEmptySchema(modelSchemaRef) {
	//	if modelSchemaRef, err = client.GetSchema_PATCH_RequestBody(resourcePath); err != nil {
	//		return nil, fmt.Errorf("failed to get PATCH request schema for resource %q: %w", resourcePath, err)
	//	}
	//	infoWithContext(
	//		ctx,
	//		fmt.Sprintf("Model schema is not present for resource %q, using PATCH request schema instead", resourcePath),
	//	)
	//}

	var description, title string
	if modelSchemaRef.Value != nil {
		title = modelSchemaRef.Value.Title
		description = modelSchemaRef.Value.Description
	}
	if description == "" {
		description = title
	}

	allProps := make(map[string]*SchemaEntry)
	// Required fields only from the POST request schema.
	// Required fields from response (modelSchemaRef) is not considered to be required for request.
	requiredFields := createSchemaRef.Value.Required

	if createSchemaRef.Value != nil {
		// Add fields from the POST request schema.
		// fields are optional=true - means fields can be passed via POST request.
		// fields are computed=false - at this stage we don't know if fields are returned by API or not. Might be overridden later.
		addSchemaEntries(createSchemaRef.Value.Properties, requiredFields, hints, allProps, false, true, false, false, false, false)
	}

	if _, ok := allProps["id"]; ok {
		warnWithContext(ctx, "Field 'id' is present in the schema but it is not expected to be used in POST request. It is usually a read-only field.")
	}
	if _, ok := allProps["guid"]; ok {
		warnWithContext(ctx, "Field 'guid' is present in the schema but it is not expected to be used in POST request. It is usually a read-only field.")
	}

	if modelSchemaRef.Value != nil {
		// Add fields from the POST request schema (typically Optional)
		for name, ref := range modelSchemaRef.Value.Properties {
			if existing, ok := allProps[name]; ok {
				if isAmbiguousObject(ref.Value) {
					existing.Computed = false
					existingSchemaJson := schemaToJSONString(existing.Prop)
					warnWithContext(
						ctx,
						fmt.Sprintf(
							"POST schema for '%s' is ambiguous. Using schema %q from 'GET' in 'read-only' mode ",
							name,
							existingSchemaJson,
						),
					)
					continue
				}

				if diffReason, ok := compareSchemaValues(existing.Prop, ref.Value); ok {
					if !existing.Required {
						// Field is present in both POST and GET with identical schema.
						// Normally mark as computed, unless explicitly overridden by hints.NotComputedSchemaFields.
						if hints != nil && contains(hints.NotComputedSchemaFields, name) {
							existing.Computed = false
						} else {
							existing.Computed = true
						}
					}

				} else {
					// Schema for request body and response body is different. We cannot expect this field to be computed.
					// Force re-assign.
					warnWithContext(
						ctx,
						fmt.Sprintf("Schema for field %q is different in POST request and GET response. Reason = %s", name, diffReason),
					)
					existing.Computed = false
				}
				continue
			}
			// optional false - field is part of response schema only so cannot be passed via POST request.
			// computed true - field is returned by API so it is computed.
			addSchemaEntries(map[string]*openapi3.SchemaRef{name: ref}, requiredFields, hints, allProps, false, false, true, false, false, false)
		}
	}

	// Will be optional only fields (Query parameters).
	params, err := client.QueryParametersGET(resourcePath)
	if err != nil {
		return nil, err
	}

	for _, p := range params {
		if !isStringOrInteger(p.Schema.Value) {
			// We search only for primitive types
			continue
		}
		name := p.Name
		if contains(excludeSearchParams, name) {
			continue
		}
		if p.Required {
			continue
		}
		if strings.Contains(name, "__") {
			continue
		}
		if _, exists := allProps[name]; exists {
			continue
		}
		if !contains(excludeSearchParams, name) && !contains(hints.ReadOnlyFields, name) {
			return nil, fmt.Errorf("field %q was detected as read-only for resource %q"+
				" but is not listed in ReadOnlyFields or ExcludedSchemaFields."+
				" Such fields cannot be used in context of create/update but only as search query params", name, resourcePath)
		}
		addSchemaEntries(map[string]*openapi3.SchemaRef{name: buildTmpSchemaRefFromParam(p)}, nil, hints, allProps, false, true, false, false, false, false)
	}

	// ---- NOTE: Just for printing information about searchable fields ----
	// Search intersection of QueryParams (GET method) and paramers from POST request schema.
	// Here is it just for printing information about searchable fields.
	// Due to complexity of this logic it is important to understand it model has common searchable fields (bettween GET and POST).
	if len(hints.ReadOnlyFields) > 0 {
		infoWithContext(ctx, fmt.Sprintf("ReadOnly fields: %v", hints.ReadOnlyFields))
	}

	if searchableQueryParams, err := client.SearchableQueryParams(resourcePath); err == nil {
		var filteredSearchableParams []string
		for _, name := range searchableQueryParams {
			if _, ok := allProps[name]; ok {
				filteredSearchableParams = append(filteredSearchableParams, name)
			}
		}
		if len(filteredSearchableParams) > 0 {
			infoWithContext(ctx, fmt.Sprintf("Searchable fields: %v", searchableQueryParams))
		} else {
			warnWithContext(ctx, fmt.Sprintf("No searchable fields (intersection between GET and POST models) found: %q", resourcePath))
		}
	}

	attrs := buildResourceAttributesFromMap(ctx, allProps, hints)
	if hints.AdditionalSchemaAttributes != nil {
		for k, v := range hints.AdditionalSchemaAttributes {
			att, ok := v.(rschema.Attribute)
			if !ok {
				return nil, fmt.Errorf("additional schema attribute %q is not a valid rschema.Attribute (got %T)", k, v)
			}
			attrs[k] = att
		}
	}

	return &rschema.Schema{
		Description:         description,
		MarkdownDescription: description,
		Attributes:          attrs,
	}, nil
}

func getResourceSchemaForCustom(hints *TFStateHints) (*rschema.Schema, error) {
	customHints := hints.TFStateHintsForCustom
	if customHints.SchemaAttributes == nil || len(customHints.SchemaAttributes) == 0 {
		return nil, fmt.Errorf("custom datasource schema attributes are required but were empty")
	}

	attrs := make(map[string]rschema.Attribute)
	for k, v := range customHints.SchemaAttributes {
		if att, ok := v.(rschema.Attribute); ok {
			attrs[k] = att
		} else {
			return nil, fmt.Errorf("additional schema attribute %q is not a valid schema.Attribute (got %T)", k, v)
		}
	}

	var description, markdownDescription string
	description = customHints.Description
	if customHints.MarkdownDescription != "" {
		markdownDescription = customHints.MarkdownDescription
	} else {
		markdownDescription = description
	}

	return &rschema.Schema{
		Description:         description,
		MarkdownDescription: markdownDescription,
		Attributes:          attrs,
	}, nil

}

func buildResourceAttributesFromMap(ctx context.Context, entries map[string]*SchemaEntry, hints *TFStateHints) map[string]rschema.Attribute {
	result := make(map[string]rschema.Attribute, len(entries))
	for name, entry := range entries {
		schema := resolveComposedSchema(resolveAllRefs(&openapi3.SchemaRef{Value: entry.Prop}))
		attr := buildResourceAttribute(ctx, name, schema, entry, hints)
		if attr == nil {
			continue
		}
		result[name] = attr
	}
	return result
}

func buildResourceAttribute(
	ctx context.Context,
	name string,
	schema *openapi3.Schema,
	entry *SchemaEntry,
	hints *TFStateHints,
) rschema.Attribute {
	schema = resolveComposedSchema(schema)
	if schema == nil || schema.Type == nil || len(*schema.Type) == 0 {
		panic("missing or invalid schema type")
	}

	desc := entry.Description

	switch (*schema.Type)[0] {
	case openapi3.TypeString:
		att := rschema.StringAttribute{
			Required:            entry.Required,
			Optional:            entry.Optional,
			Computed:            entry.Computed,
			WriteOnly:           entry.WriteOnly,
			Sensitive:           entry.Sensitive,
			Description:         desc,
			MarkdownDescription: desc,
		}

		if validatorName := hints.CommonValidatorsMapping[name]; validatorName != "" {
			if validators, ok := commonStringValidators[validatorName]; ok {
				att.Validators = append(att.Validators, validators...)
			}
		}
		// OpenAPI enum
		if len(schema.Enum) > 0 {
			var enumValues []string
			for _, val := range schema.Enum {
				if str, ok := val.(string); ok {
					enumValues = append(enumValues, str)
				}
			}
			if len(enumValues) > 0 {
				att.Validators = append(att.Validators, stringvalidator.OneOf(enumValues...))
			}
		}

		return injectModifiers(att, name, hints)

	case openapi3.TypeInteger:
		att := rschema.Int64Attribute{
			Required:            entry.Required,
			Optional:            entry.Optional,
			Computed:            entry.Computed,
			WriteOnly:           entry.WriteOnly,
			Sensitive:           entry.Sensitive,
			Description:         desc,
			MarkdownDescription: desc,
		}

		if validatorName := hints.CommonValidatorsMapping[name]; validatorName != "" {
			if validators, ok := commonIntValidators[validatorName]; ok {
				att.Validators = append(att.Validators, validators...)
			}
		}

		// OpenAPI enum
		if len(schema.Enum) > 0 {
			var enumValues []int64
			for _, val := range schema.Enum {
				switch v := val.(type) {
				case int:
					enumValues = append(enumValues, int64(v))
				case int64:
					enumValues = append(enumValues, v)
				case float64:
					enumValues = append(enumValues, int64(v))
				}
			}
			if len(enumValues) > 0 {
				att.Validators = append(att.Validators, int64validator.OneOf(enumValues...))
			}
		}
		// Min/max validation
		if schema.Min != nil && schema.Max != nil {
			att.Validators = append(att.Validators, int64validator.Between(int64(*schema.Min), int64(*schema.Max)))
		} else if schema.Min != nil {
			att.Validators = append(att.Validators, int64validator.Between(int64(*schema.Min), int64(^uint64(0)>>1))) // MaxInt64
		} else if schema.Max != nil {
			att.Validators = append(att.Validators, int64validator.Between(int64(-1<<63), int64(*schema.Max))) // MinInt64
		}

		return injectModifiers(att, name, hints)

	case openapi3.TypeNumber:
		att := rschema.Float64Attribute{
			Required:            entry.Required,
			Optional:            entry.Optional,
			Computed:            entry.Computed,
			WriteOnly:           entry.WriteOnly,
			Sensitive:           entry.Sensitive,
			Description:         desc,
			MarkdownDescription: desc,
		}

		if validatorName := hints.CommonValidatorsMapping[name]; validatorName != "" {
			if validators, ok := commonFloatValidators[validatorName]; ok {
				att.Validators = append(att.Validators, validators...)
			}
		}

		// OpenAPI enum
		if len(schema.Enum) > 0 {
			var floatValues []float64
			for _, val := range schema.Enum {
				switch v := val.(type) {
				case float64:
					floatValues = append(floatValues, v)
				case int: // fallback handling
					floatValues = append(floatValues, float64(v))
				}
			}
			if len(floatValues) > 0 {
				att.Validators = append(att.Validators, float64validator.OneOf(floatValues...))
			}
		}

		return injectModifiers(att, name, hints)

	case openapi3.TypeBoolean:
		att := rschema.BoolAttribute{
			Required:            entry.Required,
			Optional:            entry.Optional,
			Computed:            entry.Computed,
			WriteOnly:           entry.WriteOnly,
			Sensitive:           entry.Sensitive,
			Description:         desc,
			MarkdownDescription: desc,
		}

		return injectModifiers(att, name, hints)

	case openapi3.TypeArray:
		itemSchema := resolveComposedSchema(resolveAllRefs(schema.Items))
		if itemSchema == nil || itemSchema.Type == nil || len(*itemSchema.Type) == 0 {
			panic("invalid item schema for array")
		}

		if (*itemSchema.Type)[0] == openapi3.TypeObject && len(itemSchema.Properties) == 0 {
			warnWithContext(ctx, fmt.Sprintf("Skipping. Nested object for field %q: no properties defined", name))
			return nil
		}

		isOrdered := entry.Ordered

		switch (*itemSchema.Type)[0] {
		case openapi3.TypeObject:
			nested := make(map[string]*SchemaEntry)
			addSchemaEntries(itemSchema.Properties, itemSchema.Required, hints, nested, entry.Required, entry.Optional, entry.Computed, entry.WriteOnly, entry.Sensitive, entry.Ordered)
			attributes := buildResourceAttributesFromMap(ctx, nested, hints)

			if isOrdered {
				return rschema.ListNestedAttribute{
					NestedObject: rschema.NestedAttributeObject{
						Attributes: attributes,
					},
					Required:            entry.Required,
					Optional:            entry.Optional,
					Computed:            entry.Computed,
					Sensitive:           entry.Sensitive,
					Description:         desc,
					MarkdownDescription: desc,
				}
			}

			return rschema.SetNestedAttribute{
				NestedObject: rschema.NestedAttributeObject{
					Attributes: attributes,
				},
				Required:            entry.Required,
				Optional:            entry.Optional,
				Computed:            entry.Computed,
				Sensitive:           entry.Sensitive,
				Description:         desc,
				MarkdownDescription: desc,
			}

		case openapi3.TypeArray:
			inner := resolveComposedSchema(resolveAllRefs(itemSchema.Items))
			innerType := buildAttrTypeFromSchema(inner)

			if isOrdered {
				att := rschema.ListAttribute{
					ElementType:         types.ListType{ElemType: innerType},
					Required:            entry.Required,
					Optional:            entry.Optional,
					Computed:            entry.Computed,
					Sensitive:           entry.Sensitive,
					Description:         desc,
					MarkdownDescription: desc,
				}
				return injectModifiers(att, name, hints)
			}

			att := rschema.SetAttribute{
				ElementType:         types.SetType{ElemType: innerType},
				Required:            entry.Required,
				Optional:            entry.Optional,
				Computed:            entry.Computed,
				Sensitive:           entry.Sensitive,
				Description:         desc,
				MarkdownDescription: desc,
			}
			return injectModifiers(att, name, hints)

		default:
			elemType := buildAttrTypeFromSchema(itemSchema)

			if isOrdered {
				att := rschema.ListAttribute{
					ElementType:         elemType,
					Required:            entry.Required,
					Optional:            entry.Optional,
					Computed:            entry.Computed,
					Sensitive:           entry.Sensitive,
					Description:         desc,
					MarkdownDescription: desc,
				}
				return injectModifiers(att, name, hints)
			}

			att := rschema.SetAttribute{
				ElementType:         elemType,
				Required:            entry.Required,
				Optional:            entry.Optional,
				Computed:            entry.Computed,
				Sensitive:           entry.Sensitive,
				Description:         desc,
				MarkdownDescription: desc,
			}
			return injectModifiers(att, name, hints)
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
				return rschema.MapNestedAttribute{
					NestedObject: rschema.NestedAttributeObject{
						Attributes: buildResourceAttributesFromMap(ctx, nested, hints),
					},
					Required:            entry.Required,
					Optional:            entry.Optional,
					Computed:            entry.Computed,
					Sensitive:           entry.Sensitive,
					Description:         desc,
					MarkdownDescription: desc,
				}
			}

			att := rschema.MapAttribute{
				ElementType:         valType,
				Required:            entry.Required,
				Optional:            entry.Optional,
				Computed:            entry.Computed,
				Sensitive:           entry.Sensitive,
				Description:         desc,
				MarkdownDescription: desc,
			}
			return injectModifiers(att, name, hints)
		}

		if len(schema.Properties) > 0 {
			nested := make(map[string]*SchemaEntry)
			addSchemaEntries(schema.Properties, schema.Required, hints, nested, entry.Required, entry.Optional, entry.Computed, entry.WriteOnly, entry.Sensitive, entry.Ordered)
			return rschema.SingleNestedAttribute{
				Attributes:          buildResourceAttributesFromMap(ctx, nested, hints),
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
