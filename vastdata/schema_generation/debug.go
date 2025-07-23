// Copyright (c) HashiCorp, Inc.

package schema_generation

import (
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
)

func printSchema(label string, s *openapi3.Schema) {
	if s == nil {
		fmt.Printf("%s: <nil>\n", label)
		return
	}
	j, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Printf("%s: error marshaling schema: %v\n", label, err)
		return
	}
	fmt.Printf("=== %s ===\n%s\n", label, string(j))
}

func printSchemaRef(label string, ref *openapi3.SchemaRef) {
	fmt.Printf("\n--- %s ---\n", label)

	if ref == nil {
		fmt.Println("nil SchemaRef")
		return
	}

	fmt.Printf("Ref: %q\n", ref.Ref)

	if ref.Value != nil {
		tmp := struct {
			Title       string                 `json:"title,omitempty"`
			Description string                 `json:"description,omitempty"`
			Type        []string               `json:"type,omitempty"`
			Properties  map[string]interface{} `json:"properties,omitempty"`
		}{
			Title:       ref.Value.Title,
			Description: ref.Value.Description,
		}

		if ref.Value.Type != nil {
			tmp.Type = *ref.Value.Type
		}

		// Print property names only to avoid recursive depth
		tmp.Properties = map[string]interface{}{}
		for k := range ref.Value.Properties {
			tmp.Properties[k] = "...omitted..."
		}

		pretty, _ := json.MarshalIndent(tmp, "", "  ")
		fmt.Println(string(pretty))
	} else {
		fmt.Println("SchemaRef.Value is nil")
	}
}

func schemaToJSONString(s *openapi3.Schema) string {
	if s == nil {
		return "<nil>"
	}
	j, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("error marshaling schema: %v", err)
	}
	return string(j)
}
