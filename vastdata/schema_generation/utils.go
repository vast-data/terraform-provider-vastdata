// Copyright (c) HashiCorp, Inc.

package schema_generation

import (
	"context"
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vast-data/go-vast-client/openapi_schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

type TFStateHints = is.TFStateHints

// Aliases to shared OpenAPI schema utilities from go-vast-client/openapi_schema
var (
	resolveComposedSchema = openapi_schema.ResolveComposedSchema
	resolveAllRefs        = openapi_schema.ResolveAllRefs
	isObject              = openapi_schema.IsObject
	isAmbiguousObject     = openapi_schema.IsAmbiguousObject
	isPrimitive           = openapi_schema.IsPrimitive
	isStringOrInteger     = openapi_schema.IsStringOrInteger
	IsEmptySchema         = openapi_schema.IsEmptySchema
	compareSchemaValues   = openapi_schema.CompareSchemaValues
	getSchemaType         = openapi_schema.GetSchemaType
)

// getKnownDescription returns known descriptions for properties that lost their descriptions
func getKnownDescription(propertyName string, schema *openapi3.Schema) string {
	// Known properties that have descriptions in source Swagger but lost during conversion
	knownDescriptions := map[string]string{
		"permission_per_vip_pool": "VIP pools permissions map - {vippool_id: permission}. Example - {1: 'RW'}.",
	}

	if desc, ok := knownDescriptions[propertyName]; ok {
		return desc
	}

	return ""
}

type SchemaEntry struct {
	Prop        *openapi3.Schema
	Required    bool
	Optional    bool
	Computed    bool
	WriteOnly   bool
	Sensitive   bool
	Ordered     bool
	Description string
	Children    map[string]*SchemaEntry
}

func (s *SchemaEntry) String() string {
	return fmt.Sprintf(
		"<entry required=%t, optional=%t, computed=%t, writeonly=%t, sensitive=%t, ordered=%t>",
		s.Required,
		s.Optional,
		s.Computed,
		s.WriteOnly,
		s.Sensitive,
		s.Ordered,
	)
}

func addSchemaEntries(
	props map[string]*openapi3.SchemaRef,
	requiredFields []string,
	hints *TFStateHints,
	target map[string]*SchemaEntry,
	required, optional, computed, writeOnly, sensitive, ordered bool,
) {
	for name, ref := range props {
		schema := ref.Value
		if schema == nil || isExcluded(name, hints) {
			continue
		}
		if isAmbiguousObject(schema) {
			// Skip ambiguous objects with no properties
			warnWithContext(context.Background(), fmt.Sprintf("Skipping ambiguous object schema for field '%s' with no properties", name))
			continue
		}
		// Requiredness is determined strictly by this object's required array
		fieldRequired := contains(requiredFields, name)
		fieldOptional := optional
		fieldComputed := computed
		fieldWriteOnly := writeOnly
		fieldSensitive := sensitive
		fieldOrdered := ordered

		fieldRequired, fieldOptional, fieldComputed, fieldWriteOnly, fieldSensitive, fieldOrdered = flagsFromHintsForResource(name, hints, fieldRequired, fieldOptional, fieldComputed, fieldSensitive, fieldOrdered, fieldWriteOnly)

		desc := schema.Description
		if desc == "" {
			desc = schema.Title
		}
		if desc == "" {
			desc = getKnownDescription(name, schema)
		}

		// Ensure at least one flag is set; default to Optional when none provided
		if !fieldRequired && !fieldOptional && !fieldComputed {
			fieldOptional = true
		}

		entry := &SchemaEntry{
			Prop:        schema,
			Required:    fieldRequired,
			Optional:    fieldOptional,
			Computed:    fieldComputed,
			WriteOnly:   fieldWriteOnly,
			Sensitive:   fieldSensitive,
			Ordered:     fieldOrdered,
			Description: desc,
		}

		if isObject(schema) && schema.Properties != nil {
			entry.Children = make(map[string]*SchemaEntry)
			// Do NOT propagate parent required to children; children requiredness
			// must be determined solely by the child object's own required array.
			addSchemaEntries(
				schema.Properties,
				schema.Required,
				hints,
				entry.Children,
				false, // required: do not inherit from parent
				fieldOptional,
				fieldComputed,
				fieldWriteOnly,
				fieldSensitive,
				fieldOrdered,
			)
		}

		target[name] = entry
	}
}

func buildTmpSchemaRefFromParam(p *openapi3.Parameter) *openapi3.SchemaRef {
	if p == nil || p.Schema == nil || p.Schema.Value == nil {
		return nil
	}

	// Shallow copy
	schemaCopy := *p.Schema.Value

	// Inject description
	if p.Description != "" {
		schemaCopy.Description = p.Description
	}

	// Optional: derive "title" as first sentence of description
	if schemaCopy.Title == "" && p.Description != "" {
		if idx := strings.Index(p.Description, "."); idx > 0 {
			schemaCopy.Title = strings.TrimSpace(p.Description[:idx])
		}
	}

	return &openapi3.SchemaRef{
		Value: &schemaCopy,
	}
}

func flagsFromHintsForResource(name string, hints *TFStateHints, required, optional, computed, sensitive, ordered, writeOnly bool) (bool, bool, bool, bool, bool, bool) {
	if hints != nil {
		if contains(hints.RequiredSchemaFields, name) {
			required, optional, computed = true, false, false
		}
		if contains(hints.OptionalSchemaFields, name) {
			required, optional = false, true
		}
		if contains(hints.ComputedSchemaFields, name) {
			computed = true
		}
		if contains(hints.WriteOnlyFields, name) {
			writeOnly = true
			computed = false
			optional = true
		}
		if contains(hints.CreateOnlyFields, name) {
			// Keep prior state when omitted from config (via UseStateForUnknown).
			computed = true
		}
		if contains(hints.SensitiveFields, name) {
			sensitive = true
		}
		if contains(hints.NotComputedSchemaFields, name) {
			computed = false
		}
		if contains(hints.NotOptionalSchemaFields, name) {
			optional = false
		}
		if contains(hints.NotRequiredSchemaFields, name) {
			required = false
		}
		if contains(hints.PreserveOrderFields, name) {
			ordered = true
		}
	}
	if required {
		optional = false
		computed = false
	}
	return required, optional, computed, writeOnly, sensitive, ordered
}

func isExcluded(name string, hints *TFStateHints) bool {
	return hints != nil && hints.ExcludedSchemaFields != nil && contains(hints.ExcludedSchemaFields, name)
}

func contains[T comparable](list []T, key T) bool {
	if list == nil {
		return false
	}
	for _, item := range list {
		if item == key {
			return true
		}
	}
	return false
}

func buildAttrTypeFromSchema(schema *openapi3.Schema) attr.Type {
	schema = resolveComposedSchema(schema)
	if schema == nil || schema.Type == nil || len(*schema.Type) == 0 {
		panic("invalid schema type")
	}

	switch (*schema.Type)[0] {
	case openapi3.TypeString:
		return types.StringType
	case openapi3.TypeInteger:
		return types.Int64Type
	case openapi3.TypeNumber:
		return types.Float64Type
	case openapi3.TypeBoolean:
		return types.BoolType
	case openapi3.TypeArray:
		if schema.Items == nil || schema.Items.Value == nil {
			panic("array schema missing items")
		}
		return types.SetType{
			ElemType: buildAttrTypeFromSchema(resolveComposedSchema(resolveAllRefs(schema.Items))),
		}
	case openapi3.TypeObject:
		attrTypes := make(map[string]attr.Type)
		for name, prop := range schema.Properties {
			attrTypes[name] = buildAttrTypeFromSchema(resolveComposedSchema(resolveAllRefs(prop)))
		}
		return types.ObjectType{AttrTypes: attrTypes}
	default:
		panic(fmt.Sprintf("unsupported schema type: %q", (*schema.Type)[0]))
	}
}

func warnWithContext(ctx context.Context, message string) {
	if ctx != nil {
		tflog.Debug(ctx, fmt.Sprintf("⚠  %s", message))
	} else {
		fmt.Printf("#===> ⚠️  %s\n", message)
	}
}

func infoWithContext(ctx context.Context, message string) {
	if ctx != nil {
		tflog.Debug(ctx, fmt.Sprintf("◉  %s", message))
	} else {
		fmt.Printf("#===> 🟢  %s\n", message)
	}
}

// injectModifiers applies plan modifiers from hints and automatically adds UseStateForUnknown()
// to computed attributes (both computed-only and optional+computed).
// This prevents "known after apply" noise for unchanged computed fields.
func injectModifiers(attr schema.Attribute, name string, hints *TFStateHints) schema.Attribute {
	switch a := attr.(type) {

	case schema.StringAttribute:
		if hints != nil {
			if key := hints.CommonModifiersMapping[name]; key != "" {
				if mods, ok := commonStringModifiers[key]; ok {
					a.PlanModifiers = append(a.PlanModifiers, mods...)
				}
			}
		}
		// Add UseStateForUnknown for all computed attributes (optional+computed or computed-only)
		// This prevents "known after apply" noise when the field hasn't changed
		if a.Computed && !a.Required && !a.Sensitive {
			a.PlanModifiers = append(a.PlanModifiers, stringplanmodifier.UseStateForUnknown())
		}
		return a

	case schema.Int64Attribute:
		if hints != nil {
			if key := hints.CommonModifiersMapping[name]; key != "" {
				if mods, ok := commonIntModifiers[key]; ok {
					a.PlanModifiers = append(a.PlanModifiers, mods...)
				}
			}
		}
		if a.Computed && !a.Required && !a.Sensitive {
			a.PlanModifiers = append(a.PlanModifiers, int64planmodifier.UseStateForUnknown())
		}
		return a

	case schema.Float64Attribute:
		if hints != nil {
			if key := hints.CommonModifiersMapping[name]; key != "" {
				if mods, ok := commonFloatModifiers[key]; ok {
					a.PlanModifiers = append(a.PlanModifiers, mods...)
				}
			}
		}
		if a.Computed && !a.Required && !a.Sensitive {
			a.PlanModifiers = append(a.PlanModifiers, float64planmodifier.UseStateForUnknown())
		}
		return a

	case schema.BoolAttribute:
		if a.Computed && !a.Required && !a.Sensitive {
			a.PlanModifiers = append(a.PlanModifiers, boolplanmodifier.UseStateForUnknown())
		}
		return a

	case schema.ListAttribute:
		if a.Computed && !a.Required && !a.Sensitive {
			a.PlanModifiers = append(a.PlanModifiers, listplanmodifier.UseStateForUnknown())
		}
		return a

	case schema.SetAttribute:
		if a.Computed && !a.Required && !a.Sensitive {
			a.PlanModifiers = append(a.PlanModifiers, setplanmodifier.UseStateForUnknown())
		}
		return a

	case schema.MapAttribute:
		if a.Computed && !a.Required && !a.Sensitive {
			a.PlanModifiers = append(a.PlanModifiers, mapplanmodifier.UseStateForUnknown())
		}
		return a

	case schema.SingleNestedAttribute:
		if a.Computed && !a.Required && !a.Sensitive {
			a.PlanModifiers = append(a.PlanModifiers, objectplanmodifier.UseStateForUnknown())
		}
		return a

	case schema.ObjectAttribute:
		if a.Computed && !a.Required && !a.Sensitive {
			a.PlanModifiers = append(a.PlanModifiers, objectplanmodifier.UseStateForUnknown())
		}
		return a

	case schema.ListNestedAttribute:
		if a.Computed && !a.Required && !a.Sensitive {
			a.PlanModifiers = append(a.PlanModifiers, listplanmodifier.UseStateForUnknown())
		}
		return a

	case schema.SetNestedAttribute:
		if a.Computed && !a.Required && !a.Sensitive {
			a.PlanModifiers = append(a.PlanModifiers, setplanmodifier.UseStateForUnknown())
		}
		return a

	// For other unsupported attribute types, return as-is
	default:
		return attr
	}
}
