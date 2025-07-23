// Copyright (c) HashiCorp, Inc.

package schema_generation

import (
	"context"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/stretchr/testify/require"
)

func Test_buildDatasourceAttribute_Primitives(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"name": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeString),
			},
			Required: true,
		},
		"age": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeInteger),
			},
			Optional: true,
		},
		"score": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeNumber),
			},
			Computed: true,
		},
		"is_active": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeBoolean),
			},
			Optional: true,
		},
	}

	attrs := buildDatasourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})

	require.IsType(t, dschema.StringAttribute{}, attrs["name"])
	require.IsType(t, dschema.Int64Attribute{}, attrs["age"])
	require.IsType(t, dschema.Float64Attribute{}, attrs["score"])
	require.IsType(t, dschema.BoolAttribute{}, attrs["is_active"])
}

func Test_buildDatasourceAttribute_ListOfStrings(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"tags": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeArray),
				Items: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: toTypes(openapi3.TypeString),
					},
				},
			},
			Computed: true,
		},
	}

	attrs := buildDatasourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.IsType(t, dschema.SetAttribute{}, attrs["tags"])
}

func Test_buildDatasourceAttribute_SkipNonComputedComplex(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"complex_array": {
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
			// Not computed, so should be skipped
		},
	}

	attrs := buildDatasourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	_, found := attrs["complex_array"]
	require.False(t, found, "non-computed complex array should be skipped")
}

func Test_buildDatasourceAttribute_AllOf(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"allof_field": {
			Prop: &openapi3.Schema{
				AllOf: []*openapi3.SchemaRef{
					{Value: &openapi3.Schema{
						Type:        toTypes(openapi3.TypeString),
						Description: "from allOf",
					}},
				},
			},
			Computed: true,
		},
	}

	attrs := buildDatasourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.Contains(t, attrs, "allof_field")
	require.IsType(t, dschema.StringAttribute{}, attrs["allof_field"])
}

func Test_buildDatasourceAttribute_AnyOf(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"anyof_field": {
			Prop: &openapi3.Schema{
				AnyOf: []*openapi3.SchemaRef{
					{Value: &openapi3.Schema{
						Type: toTypes(openapi3.TypeInteger),
					}},
				},
			},
			Computed: true,
		},
	}

	attrs := buildDatasourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.Contains(t, attrs, "anyof_field")
	require.IsType(t, dschema.Int64Attribute{}, attrs["anyof_field"])
}

func Test_buildDatasourceAttribute_OneOf(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"oneof_field": {
			Prop: &openapi3.Schema{
				OneOf: []*openapi3.SchemaRef{
					{Value: &openapi3.Schema{
						Type: toTypes(openapi3.TypeBoolean),
					}},
				},
			},
			Computed: true,
		},
	}

	attrs := buildDatasourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.Contains(t, attrs, "oneof_field")
	require.IsType(t, dschema.BoolAttribute{}, attrs["oneof_field"])
}

func Test_buildDatasourceAttribute_ArrayOfObjects(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"items": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeArray),
				Items: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: toTypes(openapi3.TypeObject),
						Properties: map[string]*openapi3.SchemaRef{
							"id": {Value: &openapi3.Schema{Type: toTypes(openapi3.TypeString)}},
						},
					},
				},
			},
			Computed: true,
		},
	}

	attrs := buildDatasourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.Contains(t, attrs, "items")
	require.IsType(t, dschema.SetNestedAttribute{}, attrs["items"])
}

func Test_buildDatasourceAttribute_ObjectWithMapOfObject(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"configs": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeObject),
				AdditionalProperties: openapi3.AdditionalProperties{
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: toTypes(openapi3.TypeObject),
							Properties: map[string]*openapi3.SchemaRef{
								"enabled": {Value: &openapi3.Schema{Type: toTypes(openapi3.TypeBoolean)}},
							},
						},
					},
				},
			},
			Computed: true,
		},
	}

	attrs := buildDatasourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.Contains(t, attrs, "configs")
	require.IsType(t, dschema.MapNestedAttribute{}, attrs["configs"])
}

func Test_buildDatasourceAttribute_ArrayOfObjectsWithArrayOfObjects(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"matrix": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeArray),
				Items: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: toTypes(openapi3.TypeObject),
						Properties: map[string]*openapi3.SchemaRef{
							"rows": {
								Value: &openapi3.Schema{
									Type: toTypes(openapi3.TypeArray),
									Items: &openapi3.SchemaRef{
										Value: &openapi3.Schema{
											Type: toTypes(openapi3.TypeObject),
											Properties: map[string]*openapi3.SchemaRef{
												"value": {Value: &openapi3.Schema{Type: toTypes(openapi3.TypeNumber)}},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			Computed: true,
		},
	}

	attrs := buildDatasourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.Contains(t, attrs, "matrix")
	require.IsType(t, dschema.SetNestedAttribute{}, attrs["matrix"])
}

func Test_buildDatasourceAttribute_ObjectWithNestedArrayObject(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"report": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeObject),
				Properties: map[string]*openapi3.SchemaRef{
					"entries": {
						Value: &openapi3.Schema{
							Type: toTypes(openapi3.TypeArray),
							Items: &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: toTypes(openapi3.TypeObject),
									Properties: map[string]*openapi3.SchemaRef{
										"message": {Value: &openapi3.Schema{Type: toTypes(openapi3.TypeString)}},
									},
								},
							},
						},
					},
				},
			},
			Computed: true,
		},
	}

	attrs := buildDatasourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.Contains(t, attrs, "report")
	require.IsType(t, dschema.SingleNestedAttribute{}, attrs["report"])
}

func Test_buildDatasourceAttribute_AnyOfWithinArray(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"tags": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeArray),
				Items: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						AnyOf: []*openapi3.SchemaRef{
							{Value: &openapi3.Schema{Type: toTypes(openapi3.TypeString)}},
						},
					},
				},
			},
			Computed: true,
		},
	}

	attrs := buildDatasourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.Contains(t, attrs, "tags")
	require.IsType(t, dschema.SetAttribute{}, attrs["tags"])
}

func Test_buildDatasourceAttribute_AllOfWithObjectProps(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"user": {
			Prop: &openapi3.Schema{
				AllOf: []*openapi3.SchemaRef{
					{Value: &openapi3.Schema{
						Type: toTypes(openapi3.TypeObject),
						Properties: map[string]*openapi3.SchemaRef{
							"id":   {Value: &openapi3.Schema{Type: toTypes(openapi3.TypeString)}},
							"name": {Value: &openapi3.Schema{Type: toTypes(openapi3.TypeString)}},
						},
					}},
					{Value: &openapi3.Schema{
						Type: toTypes(openapi3.TypeObject),
						Properties: map[string]*openapi3.SchemaRef{
							"email": {Value: &openapi3.Schema{Type: toTypes(openapi3.TypeString)}},
						},
					}},
				},
			},
			Computed: true,
		},
	}

	attrs := buildDatasourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.Contains(t, attrs, "user")
	require.IsType(t, dschema.SingleNestedAttribute{}, attrs["user"])
}

func Test_buildDatasourceAttribute_DeepMixedComposition(t *testing.T) {
	entries := map[string]*SchemaEntry{
		"deep": {
			Prop: &openapi3.Schema{
				Type: toTypes(openapi3.TypeArray),
				Items: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: toTypes(openapi3.TypeObject),
						Properties: map[string]*openapi3.SchemaRef{
							"level1": {
								Value: &openapi3.Schema{
									Type: toTypes(openapi3.TypeArray),
									Items: &openapi3.SchemaRef{
										Value: &openapi3.Schema{
											Type: toTypes(openapi3.TypeObject),
											Properties: map[string]*openapi3.SchemaRef{
												"level2": {
													Value: &openapi3.Schema{
														AnyOf: []*openapi3.SchemaRef{
															{Value: &openapi3.Schema{
																Type: toTypes(openapi3.TypeObject),
																Properties: map[string]*openapi3.SchemaRef{
																	"level3": {
																		Value: &openapi3.Schema{
																			AllOf: []*openapi3.SchemaRef{
																				{Value: &openapi3.Schema{
																					Type: toTypes(openapi3.TypeObject),
																					Properties: map[string]*openapi3.SchemaRef{
																						"level4": {
																							Value: &openapi3.Schema{
																								Type: toTypes(openapi3.TypeString),
																							},
																						},
																					},
																				}},
																			},
																		},
																	},
																},
															}},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			Computed: true,
		},
	}

	attrs := buildDatasourceAttributesFromMap(context.TODO(), entries, &TFStateHints{})
	require.Contains(t, attrs, "deep")
	require.IsType(t, dschema.SetNestedAttribute{}, attrs["deep"])
}

// Helper to convert a single type string into an *openapi3.Types
func toTypes(t string) *openapi3.Types {
	return (*openapi3.Types)(&[]string{t})
}
