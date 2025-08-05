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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/client"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

type TFStateHints = is.TFStateHints

var resolveComposedSchema = client.ResolveComposedSchema
var resolveAllRefs = client.ResolveAllRefs

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
		fieldRequired := required || contains(requiredFields, name)
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
			addSchemaEntries(schema.Properties, schema.Required, hints, entry.Children, fieldRequired, fieldOptional, fieldComputed, fieldWriteOnly, fieldSensitive, fieldOrdered)
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

func isObject(prop *openapi3.Schema) bool {
	return prop.Type != nil && len(*prop.Type) > 0 && (*prop.Type)[0] == openapi3.TypeObject
}

func isAmbiguousObject(prop *openapi3.Schema) bool {
	return isObject(prop) && len(prop.Properties) == 0
}

func isExcluded(name string, hints *TFStateHints) bool {
	return hints != nil && hints.ExcludedSchemaFields != nil && contains(hints.ExcludedSchemaFields, name)
}

// isPrimitive returns true if the given OpenAPI schema represents a primitive type
// supported by Terraform input parameters (string, integer, number, or boolean).
//
// In data source schema generation, this is used to restrict search (input) parameters
// to primitive types, while non-primitives are only allowed if they are computed.
func isPrimitive(prop *openapi3.Schema) bool {
	if prop == nil || prop.Type == nil || len(*prop.Type) == 0 {
		return false
	}
	switch (*prop.Type)[0] {
	case openapi3.TypeString,
		openapi3.TypeInteger,
		openapi3.TypeNumber,
		openapi3.TypeBoolean:
		return true
	default:
		return false
	}
}

// isStringOrInteger returns true if the given OpenAPI schema represents string or integer
func isStringOrInteger(prop *openapi3.Schema) bool {
	if prop == nil || prop.Type == nil || len(*prop.Type) == 0 {
		return false
	}
	switch (*prop.Type)[0] {
	case openapi3.TypeString, openapi3.TypeInteger:
		return true
	default:
		return false
	}
}

func IsEmptySchema(ref *openapi3.SchemaRef) bool {
	if ref == nil || ref.Value == nil {
		return true
	}
	schema := ref.Value
	return (schema.Type == nil || len(*schema.Type) == 0) &&
		len(schema.Properties) == 0 &&
		schema.Items == nil &&
		len(schema.AllOf) == 0 &&
		len(schema.OneOf) == 0 &&
		len(schema.AnyOf) == 0 &&
		len(schema.Required) == 0
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

func compareSchemaValues(a, b *openapi3.Schema) (string, bool) {
	if a == nil || b == nil {
		if a == b {
			return "", true
		}
		return "One schema is nil while the other is not", false
	}

	typeA := getSchemaType(a)
	typeB := getSchemaType(b)
	if typeA != typeB {
		return fmt.Sprintf("Type mismatch: %q vs %q", typeA, typeB), false
	}

	// Compare array items
	if typeA == "array" {
		if a.Items == nil || b.Items == nil {
			if a.Items == b.Items {
				return "", true
			}
			return "Array item schema is nil in one but not the other", false
		}
		msg, ok := compareSchemaValues(a.Items.Value, b.Items.Value)
		if !ok {
			return fmt.Sprintf("Array item mismatch: %s", msg), false
		}
		return "", true
	}

	// Compare object properties
	if typeA == "object" {
		if len(a.Properties) != len(b.Properties) {
			return fmt.Sprintf("Object property count mismatch: %d vs %d", len(a.Properties), len(b.Properties)), false
		}
		for key, valA := range a.Properties {
			valB, ok := b.Properties[key]
			if !ok {
				return fmt.Sprintf("Property %q missing in one schema", key), false
			}
			msg, ok := compareSchemaValues(valA.Value, valB.Value)
			if !ok {
				return fmt.Sprintf("Property %q mismatch: %s", key, msg), false
			}
		}
	}

	// Optionally: compare format, enum, etc.
	return "", true
}

func getSchemaType(s *openapi3.Schema) string {
	if s == nil || s.Type == nil || len(*s.Type) == 0 {
		return ""
	}
	return (*s.Type)[0]
}

// injectModifiers applies plan modifiers from hints and automatically adds UseStateForUnknown()
// to computed-only attributes (i.e., not required, optional, or sensitive).
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
		if a.Computed && !a.Required && !a.Optional && !a.Sensitive {
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
		if a.Computed && !a.Required && !a.Optional && !a.Sensitive {
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
		if a.Computed && !a.Required && !a.Optional && !a.Sensitive {
			a.PlanModifiers = append(a.PlanModifiers, float64planmodifier.UseStateForUnknown())
		}
		return a

	case schema.BoolAttribute:
		if a.Computed && !a.Required && !a.Optional && !a.Sensitive {
			a.PlanModifiers = append(a.PlanModifiers, boolplanmodifier.UseStateForUnknown())
		}
		return a

	case schema.ListAttribute:
		if a.Computed && !a.Required && !a.Optional && !a.Sensitive {
			a.PlanModifiers = append(a.PlanModifiers, listplanmodifier.UseStateForUnknown())
		}
		return a

	case schema.SetAttribute:
		if a.Computed && !a.Required && !a.Optional && !a.Sensitive {
			a.PlanModifiers = append(a.PlanModifiers, setplanmodifier.UseStateForUnknown())
		}
		return a

	case schema.MapAttribute:
		if a.Computed && !a.Required && !a.Optional && !a.Sensitive {
			a.PlanModifiers = append(a.PlanModifiers, mapplanmodifier.UseStateForUnknown())
		}
		return a

	// Not implemented for nested attributes
	default:
		return attr
	}
}
