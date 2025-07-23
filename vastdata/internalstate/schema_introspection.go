// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/client"
)

func extractTypesFromSchema(schema any, kind SchemaContext) (map[string]attr.Type, error) {
	var attrTypes map[string]attr.Type
	switch kind {
	case SchemaForDataSource:
		schema := schema.(dsschema.Schema)
		attrTypes = make(map[string]attr.Type, len(schema.Attributes))
		for k, v := range schema.Attributes {
			attrTypes[k] = v.GetType()
		}
	case SchemaForResource:
		schema := schema.(rschema.Schema)
		attrTypes = make(map[string]attr.Type, len(schema.Attributes))
		for k, v := range schema.Attributes {
			attrTypes[k] = v.GetType()
		}
	default:
		return nil, fmt.Errorf("unsupported schema type: %T", schema)
	}
	return attrTypes, nil
}

func extractMetaFromSchema(schema any, kind SchemaContext, hints *TFStateHints) (map[string]attrMeta, error) {
	meta := make(map[string]attrMeta)

	switch kind {
	case SchemaForDataSource:
		s, ok := schema.(dsschema.Schema)
		if !ok {
			return nil, fmt.Errorf("schema is not a datasource schema")
		}
		attrs := make(map[string]any, len(s.Attributes))
		for k, v := range s.Attributes {
			attrs[k] = v
		}
		extractAttrsRecursive(meta, "", attrs, hints)

	case SchemaForResource:
		s, ok := schema.(rschema.Schema)
		if !ok {
			return nil, fmt.Errorf("schema is not a resource schema")
		}
		attrs := make(map[string]any, len(s.Attributes))
		for k, v := range s.Attributes {
			attrs[k] = v
		}
		extractAttrsRecursive(meta, "", attrs, hints)

	default:
		return nil, fmt.Errorf("unknown schema kind: %d", kind)
	}

	return meta, nil
}

func extractAttrsRecursive(meta map[string]attrMeta, prefix string, attrs map[string]any, hints *TFStateHints) {
	for name, att := range attrs {
		fullKey := name
		if prefix != "" {
			fullKey = prefix + "." + name
		}

		switch att.(type) {
		case dsschema.StringAttribute, dsschema.Int64Attribute, dsschema.BoolAttribute,
			dsschema.Float64Attribute, dsschema.NumberAttribute, dsschema.SetAttribute,
			dsschema.ListAttribute, dsschema.MapAttribute,
			rschema.StringAttribute, rschema.Int64Attribute, rschema.BoolAttribute,
			rschema.Float64Attribute, rschema.NumberAttribute, rschema.SetAttribute,
			rschema.ListAttribute, rschema.MapAttribute:
			var searchableFields []string
			if hints != nil {
				searchableFields = hints.SearchableFields

				if hints.SchemaRef != nil && hints.SchemaRef.Read != nil && prefix == "" {
					// Top-level attributes.
					searchableQueryParams, err := client.SearchableQueryParams(hints.SchemaRef.Read.Path)
					if err == nil && len(searchableQueryParams) > 0 {
						searchableFields = append(searchableFields, searchableQueryParams...)
					}
				}
			}

			metaObj := extractAttrMeta(att)
			if hints != nil {
				if contains(searchableFields, fullKey) {
					metaObj.Searchable = true
				}
				if contains(hints.ReadOnlyFields, fullKey) {
					metaObj.ReadOnly = true
				}
				if contains(hints.WriteOnlyFields, fullKey) {
					metaObj.WriteOnly = true
				}
				if contains(hints.EditOnlyFields, fullKey) {
					metaObj.EditOnly = true
				}
			}
			meta[fullKey] = metaObj

		case dsschema.SingleNestedAttribute, rschema.SingleNestedAttribute,
			dsschema.ListNestedAttribute, rschema.ListNestedAttribute,
			dsschema.MapNestedAttribute, rschema.MapNestedAttribute,
			dsschema.SetNestedAttribute, rschema.SetNestedAttribute:
			metaObj := extractAttrMeta(att)
			if hints != nil {
				if contains(hints.SearchableFields, fullKey) {
					metaObj.Searchable = true
				}
				if contains(hints.ReadOnlyFields, fullKey) {
					metaObj.ReadOnly = true
				}
				if contains(hints.WriteOnlyFields, fullKey) {
					metaObj.WriteOnly = true
				}
				if contains(hints.EditOnlyFields, fullKey) {
					metaObj.EditOnly = true
				}
			}
			meta[fullKey] = metaObj
			if nestedAttrs := getNestedAttributes(att); nestedAttrs != nil {
				extractAttrsRecursive(meta, fullKey, nestedAttrs, hints)
			}

		default:
			panic(fmt.Sprintf("unsupported attribute type %T at %q", att, fullKey))
		}
	}
}

func extractAttrMeta(attr any) attrMeta {
	if m, ok := attr.(interface {
		IsRequired() bool
		IsOptional() bool
		IsComputed() bool
		IsWriteOnly() bool
	}); ok {
		return attrMeta{
			Required:  m.IsRequired(),
			Optional:  m.IsOptional(),
			Computed:  m.IsComputed(),
			WriteOnly: m.IsWriteOnly(),
		}
	}
	panic(fmt.Sprintf("attribute %T does not implement metadata accessors", attr))
}
