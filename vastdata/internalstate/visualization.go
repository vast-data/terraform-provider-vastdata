// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"sort"
	"strings"
)

var bold = color.New(color.Bold).SprintFunc()

// BuildDataSourceAttributesString recursively formats the schema attributes of a Terraform
// data source into a readable string. It includes attribute name, type, and modifier flags
// such as required, optional, computed, and sensitive.
//
// Parameters:
//   - attrs: map of attribute names to dschema.Attribute instances.
//   - indent: indentation level for nested attributes.
//
// Returns:
//   - A formatted string of all attributes and their characteristics for the data source schema.
func BuildDataSourceAttributesString(attrs map[string]dschema.Attribute, automation bool, indent int) string {
	var b strings.Builder
	indentStr := strings.Repeat("  ", indent)
	marker := "🔹"
	if automation {
		marker = ""
	}

	for name, attr := range attrs {
		label := bold(name)
		if automation {
			label = fmt.Sprintf("%s:", label)
		}

		switch a := attr.(type) {
		case dschema.StringAttribute:
			fmt.Fprintf(&b, "%s%s %s (string) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, a.Required, a.Optional, a.Computed, a.Sensitive, a.IsWriteOnly())

		case dschema.Int64Attribute:
			fmt.Fprintf(&b, "%s%s %s (int64) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, a.Required, a.Optional, a.Computed, a.Sensitive, a.IsWriteOnly())

		case dschema.Float64Attribute:
			fmt.Fprintf(&b, "%s%s %s (float64) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, a.Required, a.Optional, a.Computed, a.Sensitive, a.IsWriteOnly())

		case dschema.NumberAttribute:
			fmt.Fprintf(&b, "%s%s %s (number) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, a.Required, a.Optional, a.Computed, a.Sensitive, a.IsWriteOnly())

		case dschema.BoolAttribute:
			fmt.Fprintf(&b, "%s%s %s (bool) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, a.Required, a.Optional, a.Computed, a.Sensitive, a.IsWriteOnly())

		case dschema.ListAttribute:
			fmt.Fprintf(&b, "%s%s %s (list<%s>) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, simplifyType(a.ElementType), a.Required, a.Optional, a.Computed, a.Sensitive, a.IsWriteOnly())

		case dschema.SetAttribute:
			fmt.Fprintf(&b, "%s%s %s (set<%s>) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, simplifyType(a.ElementType), a.Required, a.Optional, a.Computed, a.Sensitive, a.IsWriteOnly())

		case dschema.MapAttribute:
			fmt.Fprintf(&b, "%s%s %s (map<%s>) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, simplifyType(a.ElementType), a.Required, a.Optional, a.Computed, a.Sensitive, a.IsWriteOnly())

		case dschema.SingleNestedAttribute:
			fmt.Fprintf(&b, "%s%s %s (object) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, a.Required, a.Optional, a.Computed, a.Sensitive, a.IsWriteOnly())
			b.WriteString(BuildDataSourceAttributesString(a.Attributes, automation, indent+2))

		case dschema.ListNestedAttribute:
			fmt.Fprintf(&b, "%s%s %s (list<object>) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, a.Required, a.Optional, a.Computed, a.Sensitive, a.IsWriteOnly())
			b.WriteString(BuildDataSourceAttributesString(a.NestedObject.Attributes, automation, indent+2))

		case dschema.SetNestedAttribute:
			fmt.Fprintf(&b, "%s%s %s (set<object>) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, a.Required, a.Optional, a.Computed, a.Sensitive, a.IsWriteOnly())
			b.WriteString(BuildDataSourceAttributesString(a.NestedObject.Attributes, automation, indent+2))

		case dschema.MapNestedAttribute:
			fmt.Fprintf(&b, "%s%s %s (map<object>) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, a.Required, a.Optional, a.Computed, a.Sensitive, a.IsWriteOnly())
			b.WriteString(BuildDataSourceAttributesString(a.NestedObject.Attributes, automation, indent+2))

		default:
			fmt.Fprintf(&b, "%s%s %s (unknown type %T)\n", indentStr, marker, label, a)
		}
	}

	return b.String()
}

// BuildResourceAttributesString recursively formats the schema attributes of a Terraform
// resource into a readable string. It includes attribute name, type, and modifier flags
// such as required, optional, computed, and sensitive.
//
// Parameters:
//   - attrs: map of attribute names to rschema.Attribute instances.
//   - indent: indentation level for nested attributes.
//
// Returns:
//   - A formatted string of all attributes and their characteristics for the resource schema.
func BuildResourceAttributesString(attrs map[string]rschema.Attribute, automation bool, indent int) string {
	var b strings.Builder
	indentStr := strings.Repeat("  ", indent)
	marker := "🔸"
	if automation {
		marker = ""
	}

	for name, attr := range attrs {
		label := bold(name)
		if automation {
			label = fmt.Sprintf("%s:", label)
		}

		switch a := attr.(type) {
		case rschema.StringAttribute:
			fmt.Fprintf(&b, "%s%s %s (string) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, a.Required, a.Optional, a.Computed, a.Sensitive, a.WriteOnly)

		case rschema.Int64Attribute:
			fmt.Fprintf(&b, "%s%s %s (int64) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, a.Required, a.Optional, a.Computed, a.Sensitive, a.WriteOnly)

		case rschema.Float64Attribute:
			fmt.Fprintf(&b, "%s%s %s (float64) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, a.Required, a.Optional, a.Computed, a.Sensitive, a.WriteOnly)

		case rschema.NumberAttribute:
			fmt.Fprintf(&b, "%s%s %s (number) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, a.Required, a.Optional, a.Computed, a.Sensitive, a.WriteOnly)

		case rschema.BoolAttribute:
			fmt.Fprintf(&b, "%s%s %s (bool) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, a.Required, a.Optional, a.Computed, a.Sensitive, a.WriteOnly)

		case rschema.ListAttribute:
			fmt.Fprintf(&b, "%s%s %s (list<%s>) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, simplifyType(a.ElementType), a.Required, a.Optional, a.Computed, a.Sensitive, a.WriteOnly)

		case rschema.SetAttribute:
			fmt.Fprintf(&b, "%s%s %s (set<%s>) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, simplifyType(a.ElementType), a.Required, a.Optional, a.Computed, a.Sensitive, a.IsWriteOnly())

		case rschema.MapAttribute:
			fmt.Fprintf(&b, "%s%s %s (map<%s>) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, simplifyType(a.ElementType), a.Required, a.Optional, a.Computed, a.Sensitive, a.WriteOnly)

		case rschema.SingleNestedAttribute:
			fmt.Fprintf(&b, "%s%s %s (object) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, a.Required, a.Optional, a.Computed, a.Sensitive, a.WriteOnly)
			b.WriteString(BuildResourceAttributesString(a.Attributes, automation, indent+2))

		case rschema.ListNestedAttribute:
			fmt.Fprintf(&b, "%s%s %s (list<object>) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, a.Required, a.Optional, a.Computed, a.Sensitive, a.WriteOnly)
			b.WriteString(BuildResourceAttributesString(a.NestedObject.Attributes, automation, indent+2))

		case rschema.SetNestedAttribute:
			fmt.Fprintf(&b, "%s%s %s (set<object>) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, a.Required, a.Optional, a.Computed, a.Sensitive, a.IsWriteOnly())
			b.WriteString(BuildResourceAttributesString(a.NestedObject.Attributes, automation, indent+2))

		case rschema.MapNestedAttribute:
			fmt.Fprintf(&b, "%s%s %s (map<object>) [req=%v opt=%v comp=%v sens=%v wo=%v]\n",
				indentStr, marker, label, a.Required, a.Optional, a.Computed, a.Sensitive, a.WriteOnly)
			b.WriteString(BuildResourceAttributesString(a.NestedObject.Attributes, automation, indent+2))

		default:
			fmt.Fprintf(&b, "%s%s %s (unknown type %T)\n", indentStr, marker, label, a)
		}
	}

	return b.String()
}

func simplifyType(t attr.Type) string {
	switch t.String() {
	case types.StringType.String():
		return "string"
	case types.Int64Type.String():
		return "int64"
	case types.Float64Type.String():
		return "float64"
	case types.BoolType.String():
		return "bool"
	case types.NumberType.String():
		return "number"
	case types.DynamicType.String():
		return "dynamic"
	}

	switch tt := t.(type) {
	case types.ListType:
		return fmt.Sprintf("list<%s>", simplifyType(tt.ElemType))
	case types.SetType:
		return fmt.Sprintf("set<%s>", simplifyType(tt.ElemType))
	case types.MapType:
		return fmt.Sprintf("map<%s>", simplifyType(tt.ElemType))
	case types.ObjectType:
		return "object"
	default:
		return t.String() // fallback
	}
}

func schemaVisualization(schema any, kind SchemaContext, automation bool) string {
	switch kind {
	case SchemaForDataSource:
		return BuildDataSourceAttributesString(schema.(dschema.Schema).Attributes, automation, 2)
	case SchemaForResource:
		return BuildResourceAttributesString(schema.(rschema.Schema).Attributes, automation, 2)
	default:
		panic(fmt.Sprintf("unknown schema visualization kind: %v", kind))
	}
}

// prettyWithMeta formats a nested map of Terraform attribute values with associated metadata
// into a human-readable string. It recursively walks through complex attribute types (objects,
// lists, sets), and indents each nested level.
//
// The meta map provides attribute modifiers like `required`, `optional`, or `computed`,
// keyed by the full dotted path to each attribute.
//
// Parameters:
//   - data: a map of attribute names to Terraform attr.Value instances.
//   - meta: a map of attribute metadata (required/optional/computed) by dotted path.
//
// Returns:
//   - A formatted multiline string representing the attribute tree with type, modifiers,
//     and value for each field.
func prettyWithMeta(data map[string]attr.Value, meta map[string]attrMeta) string {
	type entry struct {
		Level    int
		Path     string
		Name     string
		Type     string
		Mod      string
		ValueStr string
	}

	var entries []entry

	shortenType := func(t string) string {
		t = strings.TrimPrefix(t, "basetypes.")
		t = strings.TrimPrefix(t, "types.")
		if idx := strings.Index(t, "["); idx > 0 {
			return t[:idx] + "[…]"
		}
		return t
	}

	var walk func(key string, val attr.Value, level int, path string)
	walk = func(key string, val attr.Value, level int, path string) {
		fullPath := key
		if path != "" {
			fullPath = path + "." + key
		}

		// --- Modifiers
		mod := ""
		if m, ok := meta[fullPath]; ok {
			var mods []string
			if m.Required {
				mods = append(mods, "required")
			}
			if m.Optional {
				mods = append(mods, "optional")
			}
			if m.Computed {
				mods = append(mods, "computed")
			}
			if m.Searchable {
				mods = append(mods, "searchable")
			}
			if len(mods) > 0 {
				mod = fmt.Sprintf("(%s)", strings.Join(mods, ","))
			}
		}

		var valStr string
		skipEntry := false

		switch v := val.(type) {
		case types.String:
			if val.IsNull() || val.IsUnknown() {
				valStr = "<null>"
			} else {
				valStr = fmt.Sprintf("%q", v.ValueString())
			}

		case types.Int64:
			if val.IsNull() || val.IsUnknown() {
				valStr = "<null>"
			} else {
				valStr = fmt.Sprintf("%d", v.ValueInt64())
			}

		case types.Float64:
			if val.IsNull() || val.IsUnknown() {
				valStr = "<null>"
			} else {
				valStr = fmt.Sprintf("%f", v.ValueFloat64())
			}

		case types.Bool:
			if val.IsNull() || val.IsUnknown() {
				valStr = "<null>"
			} else {
				valStr = fmt.Sprintf("%v", v.ValueBool())
			}

		case types.Object:
			if val.IsNull() || val.IsUnknown() {
				valStr = "<null>"
			} else {
				skipEntry = true
				keys := make([]string, 0, len(v.Attributes()))
				for k := range v.Attributes() {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					walk(k, v.Attributes()[k], level+1, fullPath)
				}
			}

		case types.List:
			if val.IsNull() || val.IsUnknown() {
				valStr = "<null>"
			} else {
				elems := v.Elements()
				values := make([]string, 0, len(elems))
				for _, e := range elems {
					if e.IsNull() || e.IsUnknown() {
						values = append(values, "<null>")
					} else {
						values = append(values, fmt.Sprintf("%v", ConvertAttrValueToRaw(e, v.ElementType(context.Background()))))
					}
				}
				valStr = "[" + strings.Join(values, ", ") + "]"
			}

		case types.Set:
			if val.IsNull() || val.IsUnknown() {
				valStr = "<null>"
			} else {
				elems := v.Elements()
				values := make([]string, 0, len(elems))
				for _, e := range elems {
					if e.IsNull() || e.IsUnknown() {
						values = append(values, "<null>")
					} else {
						values = append(values, fmt.Sprintf("%v", ConvertAttrValueToRaw(e, v.ElementType(context.Background()))))
					}
				}
				valStr = "[" + strings.Join(values, ", ") + "]"
			}

		default:
			valStr = "<unsupported>"
		}

		if !skipEntry {
			typeStr := shortenType(val.Type(context.Background()).String())

			entries = append(entries, entry{
				Level:    level,
				Path:     fullPath,
				Name:     key,
				Type:     typeStr,
				Mod:      mod,
				ValueStr: valStr,
			})
		}
	}

	// Walk all top-level keys
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		walk(k, data[k], 0, "")
	}

	// Align and print
	maxType, maxName, maxMod := 0, 0, 0
	for _, e := range entries {
		if len(e.Type) > maxType {
			maxType = len(e.Type)
		}
		if len(e.Name)+2*e.Level > maxName {
			maxName = len(e.Name) + 2*e.Level
		}
		if len(e.Mod) > maxMod {
			maxMod = len(e.Mod)
		}
	}

	var b strings.Builder
	for _, e := range entries {
		indent := strings.Repeat("  ", e.Level)
		fmt.Fprintf(&b, "[%-*s] %-*s %-*s : %s\n",
			maxType, e.Type,
			maxName, indent+e.Name,
			maxMod, e.Mod,
			e.ValueStr,
		)
	}

	return b.String()
}
