// Copyright (c) HashiCorp, Inc.

package schema_generation

import (
	"context"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

// BenchmarkGetResourceSchema measures the performance of resource schema generation
func BenchmarkGetResourceSchema(b *testing.B) {
	hints := &internalstate.TFStateHints{
		SchemaRef: &internalstate.SchemaReference{
			Create: &internalstate.OpenAPIEndpointRef{
				Method: http.MethodPost,
				Path:   "users",
			},
			Read: &internalstate.OpenAPIEndpointRef{
				Method: http.MethodGet,
				Path:   "users",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GetResourceSchema(context.Background(), hints)
		if err != nil {
			b.Fatalf("Schema generation failed: %v", err)
		}
	}
}

// BenchmarkGetDatasourceSchema measures the performance of datasource schema generation
func BenchmarkGetDatasourceSchema(b *testing.B) {
	hints := &internalstate.TFStateHints{
		SchemaRef: &internalstate.SchemaReference{
			Read: &internalstate.OpenAPIEndpointRef{
				Method: http.MethodGet,
				Path:   "users",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GetDatasourceSchema(context.Background(), hints)
		if err != nil {
			b.Fatalf("Schema generation failed: %v", err)
		}
	}
}

// BenchmarkBuildResourceAttribute measures attribute building performance
func BenchmarkBuildResourceAttribute(b *testing.B) {
	schema := &openapi3.Schema{
		Type:        &openapi3.Types{openapi3.TypeString},
		Description: "Test string field",
	}

	entry := &SchemaEntry{
		Prop:        schema,
		Required:    true,
		Description: "Test field",
	}

	hints := &internalstate.TFStateHints{}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		attr := buildResourceAttribute(ctx, "test_field", schema, entry, hints)
		if attr == nil {
			b.Fatal("Expected non-nil attribute")
		}
	}
}

// BenchmarkBuildResourceAttributeComplex measures complex attribute building performance
func BenchmarkBuildResourceAttributeComplex(b *testing.B) {
	// Create a complex nested object schema
	schema := &openapi3.Schema{
		Type: &openapi3.Types{openapi3.TypeObject},
		Properties: map[string]*openapi3.SchemaRef{
			"name": {Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}}},
			"age":  {Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeInteger}}},
			"tags": {Value: &openapi3.Schema{
				Type: &openapi3.Types{openapi3.TypeArray},
				Items: &openapi3.SchemaRef{
					Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}},
				},
			}},
			"metadata": {Value: &openapi3.Schema{
				Type: &openapi3.Types{openapi3.TypeObject},
				Properties: map[string]*openapi3.SchemaRef{
					"created": {Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}}},
					"updated": {Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}}},
				},
			}},
		},
	}

	entry := &SchemaEntry{
		Prop:        schema,
		Optional:    true,
		Description: "Complex nested object",
	}

	hints := &internalstate.TFStateHints{}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		attr := buildResourceAttribute(ctx, "complex_field", schema, entry, hints)
		if attr == nil {
			b.Fatal("Expected non-nil attribute")
		}
	}
}

// BenchmarkAddSchemaEntries measures the performance of adding schema entries
func BenchmarkAddSchemaEntries(b *testing.B) {
	props := make(map[string]*openapi3.SchemaRef)
	for i := 0; i < 100; i++ {
		fieldName := "field_" + string(rune('a'+i%26)) + string(rune('0'+i%10))
		props[fieldName] = &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:        &openapi3.Types{openapi3.TypeString},
				Description: "Test field " + fieldName,
			},
		}
	}

	hints := &internalstate.TFStateHints{}
	requiredFields := []string{"field_a0", "field_b1", "field_c2"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		target := make(map[string]*SchemaEntry)
		addSchemaEntries(props, requiredFields, hints, target, false, true, false, false, false, false)
		if len(target) == 0 {
			b.Fatal("Expected entries to be added")
		}
	}
}

// BenchmarkBuildAttrTypeFromSchema measures type building performance
func BenchmarkBuildAttrTypeFromSchema(b *testing.B) {
	tests := []struct {
		name   string
		schema *openapi3.Schema
	}{
		{
			name: "string",
			schema: &openapi3.Schema{
				Type: &openapi3.Types{openapi3.TypeString},
			},
		},
		{
			name: "integer",
			schema: &openapi3.Schema{
				Type: &openapi3.Types{openapi3.TypeInteger},
			},
		},
		{
			name: "array_of_strings",
			schema: &openapi3.Schema{
				Type: &openapi3.Types{openapi3.TypeArray},
				Items: &openapi3.SchemaRef{
					Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}},
				},
			},
		},
		{
			name: "object",
			schema: &openapi3.Schema{
				Type: &openapi3.Types{openapi3.TypeObject},
				Properties: map[string]*openapi3.SchemaRef{
					"name": {Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}}},
					"age":  {Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeInteger}}},
				},
			},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				attrType := buildAttrTypeFromSchema(tt.schema)
				if attrType == nil {
					b.Fatal("Expected non-nil attribute type")
				}
			}
		})
	}
}

// BenchmarkResolveComposedSchema measures schema resolution performance
func BenchmarkResolveComposedSchema(b *testing.B) {
	// Create a schema with allOf composition
	schema := &openapi3.Schema{
		AllOf: []*openapi3.SchemaRef{
			{Value: &openapi3.Schema{
				Type: &openapi3.Types{openapi3.TypeObject},
				Properties: map[string]*openapi3.SchemaRef{
					"name": {Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}}},
				},
			}},
			{Value: &openapi3.Schema{
				Type: &openapi3.Types{openapi3.TypeObject},
				Properties: map[string]*openapi3.SchemaRef{
					"age": {Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeInteger}}},
				},
			}},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resolved := resolveComposedSchema(schema)
		if resolved == nil {
			b.Fatal("Expected non-nil resolved schema")
		}
	}
}

// BenchmarkBuildResourceAttributesFromMap measures building multiple attributes
func BenchmarkBuildResourceAttributesFromMap(b *testing.B) {
	entries := make(map[string]*SchemaEntry)

	// Create 50 different field types to simulate a realistic resource
	fieldTypes := []openapi3.Types{
		{openapi3.TypeString},
		{openapi3.TypeInteger},
		{openapi3.TypeBoolean},
		{openapi3.TypeNumber},
	}

	for i := 0; i < 50; i++ {
		fieldName := "field_" + string(rune('a'+i%26)) + string(rune('0'+i%10))
		fieldType := fieldTypes[i%len(fieldTypes)]

		entries[fieldName] = &SchemaEntry{
			Prop: &openapi3.Schema{
				Type:        &fieldType,
				Description: "Test field " + fieldName,
			},
			Required:    i%3 == 0,
			Optional:    i%3 == 1,
			Computed:    i%3 == 2,
			Description: "Test field " + fieldName,
		}
	}

	hints := &internalstate.TFStateHints{}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		attrs := buildResourceAttributesFromMap(ctx, entries, hints)
		if len(attrs) == 0 {
			b.Fatal("Expected attributes to be built")
		}
	}
}

// BenchmarkIsEmptySchema measures empty schema detection performance
func BenchmarkIsEmptySchema(b *testing.B) {
	schemas := []*openapi3.SchemaRef{
		nil,
		{Value: nil},
		{Value: &openapi3.Schema{}},
		{Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}}},
		{Value: &openapi3.Schema{
			Properties: map[string]*openapi3.SchemaRef{
				"test": {Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}}},
			},
		}},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, schema := range schemas {
			_ = IsEmptySchema(schema)
		}
	}
}

// BenchmarkValidatorApplication measures validator application performance
func BenchmarkValidatorApplication(b *testing.B) {
	schema := &openapi3.Schema{
		Type: &openapi3.Types{openapi3.TypeString},
		Enum: []interface{}{"value1", "value2", "value3", "value4", "value5"},
	}

	entry := &SchemaEntry{
		Prop:        schema,
		Required:    true,
		Description: "Test field with enum validation",
	}

	hints := &internalstate.TFStateHints{
		CommonValidatorsMapping: map[string]string{
			"test_field": ValidatorPathStartsWithSlash,
		},
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		attr := buildResourceAttribute(ctx, "test_field", schema, entry, hints)
		if attr == nil {
			b.Fatal("Expected non-nil attribute")
		}
	}
}
