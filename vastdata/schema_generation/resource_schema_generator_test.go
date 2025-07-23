// Copyright (c) HashiCorp, Inc.

package schema_generation

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/stretchr/testify/require"
)

func containsUseStateForUnknown(pm any) bool {
	useStateStr := fmt.Sprintf("%#v", stringplanmodifier.UseStateForUnknown())

	switch mods := pm.(type) {
	case []planmodifier.String:
		for _, m := range mods {
			if fmt.Sprintf("%#v", m) == useStateStr {
				return true
			}
		}
	case []planmodifier.Int64:
		useStateStr = fmt.Sprintf("%#v", int64planmodifier.UseStateForUnknown())
		for _, m := range mods {
			if fmt.Sprintf("%#v", m) == useStateStr {
				return true
			}
		}
	case []planmodifier.Float64:
		useStateStr = fmt.Sprintf("%#v", float64planmodifier.UseStateForUnknown())
		for _, m := range mods {
			if fmt.Sprintf("%#v", m) == useStateStr {
				return true
			}
		}
	case []planmodifier.Bool:
		useStateStr = fmt.Sprintf("%#v", boolplanmodifier.UseStateForUnknown())
		for _, m := range mods {
			if fmt.Sprintf("%#v", m) == useStateStr {
				return true
			}
		}
	case []planmodifier.List:
		useStateStr = fmt.Sprintf("%#v", listplanmodifier.UseStateForUnknown())
		for _, m := range mods {
			if fmt.Sprintf("%#v", m) == useStateStr {
				return true
			}
		}
	case []planmodifier.Set:
		useStateStr = fmt.Sprintf("%#v", setplanmodifier.UseStateForUnknown())
		for _, m := range mods {
			if fmt.Sprintf("%#v", m) == useStateStr {
				return true
			}
		}
	case []planmodifier.Map:
		useStateStr = fmt.Sprintf("%#v", mapplanmodifier.UseStateForUnknown())
		for _, m := range mods {
			if fmt.Sprintf("%#v", m) == useStateStr {
				return true
			}
		}
	}
	return false
}

func Test_buildResourceAttribute_Primitives(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"name": {
			Prop:     &openapi3.Schema{Type: toTypes(openapi3.TypeString)},
			Required: true,
		},
		"count": {
			Prop:     &openapi3.Schema{Type: toTypes(openapi3.TypeInteger)},
			Optional: true,
		},
		"score": {
			Prop:     &openapi3.Schema{Type: toTypes(openapi3.TypeNumber)},
			Computed: true,
		},
		"is_active": {
			Prop:     &openapi3.Schema{Type: toTypes(openapi3.TypeBoolean)},
			Optional: true,
		},
	}

	attrs := buildResourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})

	require.IsType(t, rschema.StringAttribute{}, attrs["name"])
	require.IsType(t, rschema.Int64Attribute{}, attrs["count"])
	require.IsType(t, rschema.Float64Attribute{}, attrs["score"])
	require.IsType(t, rschema.BoolAttribute{}, attrs["is_active"])
}

func Test_buildResourceAttribute_ArrayOfStrings(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"tags": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeArray),
				Items: &openapi3.SchemaRef{
					Value: &openapi3.Schema{Type: toTypes(openapi3.TypeString)},
				},
			},
			Optional: true,
		},
	}

	attrs := buildResourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.IsType(t, rschema.SetAttribute{}, attrs["tags"])
}

func Test_buildResourceAttribute_MapOfInt(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"metrics": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeObject),
				AdditionalProperties: openapi3.AdditionalProperties{
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{Type: toTypes(openapi3.TypeInteger)},
					},
				},
			},
			Computed: true,
		},
	}

	attrs := buildResourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.IsType(t, rschema.MapAttribute{}, attrs["metrics"])
}

func Test_buildResourceAttribute_EmptyObjectsAreSkipped(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"empty_object": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeObject),
			},
			Computed: true,
		},
		"map_with_empty_object_values": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeObject),
				AdditionalProperties: openapi3.AdditionalProperties{
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type:       toTypes(openapi3.TypeObject),
							Properties: map[string]*openapi3.SchemaRef{},
						},
					},
				},
			},
			Computed: true,
		},
	}

	attrs := buildResourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})

	// Assert that both fields are skipped (not present in final schema)
	_, exists1 := attrs["empty_object"]
	_, exists2 := attrs["map_with_empty_object_values"]

	require.False(t, exists1, "empty_object should be skipped")
	require.False(t, exists2, "map_with_empty_object_values should be skipped")
}

func Test_buildResourceAttribute_OrderedListOfStrings(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"ordered_tags": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeArray),
				Items: &openapi3.SchemaRef{
					Value: &openapi3.Schema{Type: toTypes(openapi3.TypeString)},
				},
			},
			Optional: true,
			Ordered:  true,
		},
	}

	attrs := buildResourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.Contains(t, attrs, "ordered_tags")
	require.IsType(t, rschema.ListAttribute{}, attrs["ordered_tags"])
}

func Test_buildResourceAttribute_OrderedListOfObjects(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"ordered_items": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeArray),
				Items: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: toTypes(openapi3.TypeObject),
						Properties: map[string]*openapi3.SchemaRef{
							"id":   {Value: &openapi3.Schema{Type: toTypes(openapi3.TypeString)}},
							"rank": {Value: &openapi3.Schema{Type: toTypes(openapi3.TypeInteger)}},
						},
					},
				},
			},
			Computed: true,
			Ordered:  true,
		},
	}

	attrs := buildResourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.Contains(t, attrs, "ordered_items")
	require.IsType(t, rschema.ListNestedAttribute{}, attrs["ordered_items"])
}

func Test_buildResourceAttribute_UnorderedSetOfObjects(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"items": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeArray),
				Items: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: toTypes(openapi3.TypeObject),
						Properties: map[string]*openapi3.SchemaRef{
							"name": {Value: &openapi3.Schema{Type: toTypes(openapi3.TypeString)}},
						},
					},
				},
			},
			Optional: true,
			Ordered:  false,
		},
	}

	attrs := buildResourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.Contains(t, attrs, "items")
	require.IsType(t, rschema.SetNestedAttribute{}, attrs["items"])
}

func Test_buildResourceAttribute_ListOfLists(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"matrix": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeArray),
				Items: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: toTypes(openapi3.TypeArray),
						Items: &openapi3.SchemaRef{
							Value: &openapi3.Schema{Type: toTypes(openapi3.TypeInteger)},
						},
					},
				},
			},
			Optional: true,
			Ordered:  true,
		},
	}

	attrs := buildResourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.Contains(t, attrs, "matrix")
	require.IsType(t, rschema.ListAttribute{}, attrs["matrix"])
}

func Test_buildResourceAttribute_MapOfNestedObjects(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"labels": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeObject),
				AdditionalProperties: openapi3.AdditionalProperties{
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: toTypes(openapi3.TypeObject),
							Properties: map[string]*openapi3.SchemaRef{
								"value": {Value: &openapi3.Schema{Type: toTypes(openapi3.TypeString)}},
							},
						},
					},
				},
			},
			Computed: true,
		},
	}

	attrs := buildResourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.Contains(t, attrs, "labels")
	require.IsType(t, rschema.MapNestedAttribute{}, attrs["labels"])
}

func Test_buildResourceAttribute_SingleNestedObject(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"config": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeObject),
				Properties: map[string]*openapi3.SchemaRef{
					"enabled": {Value: &openapi3.Schema{Type: toTypes(openapi3.TypeBoolean)}},
				},
			},
			Optional: true,
		},
	}

	attrs := buildResourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.Contains(t, attrs, "config")
	require.IsType(t, rschema.SingleNestedAttribute{}, attrs["config"])
}

func Test_injectModifiers_ComputedOnlyGetsUseState(t *testing.T) {
	att := rschema.StringAttribute{
		Computed: true,
	}

	modified := injectModifiers(att, "name", &TFStateHints{})
	mod, ok := modified.(rschema.StringAttribute)
	require.True(t, ok)
	require.True(t, containsUseStateForUnknown(mod.PlanModifiers), "Expected UseStateForUnknown for computed-only attribute")
}

func Test_injectModifiers_OptionalDoesNotGetUseState(t *testing.T) {
	att := rschema.StringAttribute{
		Optional: true,
		Computed: true,
	}

	modified := injectModifiers(att, "name", &TFStateHints{})
	mod, ok := modified.(rschema.StringAttribute)
	require.True(t, ok)
	require.False(t, containsUseStateForUnknown(mod.PlanModifiers), "Should not apply UseStateForUnknown to optional+computed field")
}

func Test_injectModifiers_FromHints(t *testing.T) {
	att := rschema.StringAttribute{
		Optional: true,
	}

	hints := &TFStateHints{
		CommonModifiersMapping: map[string]string{
			"force_field": ModifierForceNew,
		},
	}

	modified := injectModifiers(att, "force_field", hints)
	mod, ok := modified.(rschema.StringAttribute)
	require.True(t, ok)
	require.Len(t, mod.PlanModifiers, 1)
}
