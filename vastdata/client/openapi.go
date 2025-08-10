// Copyright (c) HashiCorp, Inc.

package client

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
)

var (
	//go:embed api/**/*
	FS             embed.FS
	openApiDocOnce sync.Once
	openApiDoc     *openapi3.T
	openApiDocErr  error
)

// loadOpenAPIDocOnce loads and parses the OpenAPI v3 document from a .tar.gz archive exactly once.
// It looks for a file named "openapi-v3.json" inside the archive located at "api/openapi-v3.tar.gz".
// The document is parsed using the kin-openapi loader and cached for future calls.
//
// Returns:
//   - *openapi3.T: the parsed OpenAPI document.
//   - error: if the archive cannot be read, the JSON file is not found, or the document fails to parse.
//
// Notes:
//   - This function is thread-safe and memoized via sync.Once to ensure the document is only loaded once.
//   - Errors encountered during the initial load are also cached and returned on subsequent calls.
func loadOpenAPIDocOnce() (*openapi3.T, error) {
	openApiDocOnce.Do(func() {
		data, err := FS.ReadFile("api/5.3.0/api.tar.gz")
		if err != nil {
			openApiDocErr = fmt.Errorf("read embedded tar.gz: %w", err)
			return
		}

		gzr, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			openApiDocErr = fmt.Errorf("gzip reader: %w", err)
			return
		}
		defer gzr.Close()

		tr := tar.NewReader(gzr)

		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				openApiDocErr = fmt.Errorf("api.json not found in embedded archive")
				return
			}
			if err != nil {
				openApiDocErr = fmt.Errorf("tar read error: %w", err)
				return
			}

			if strings.HasSuffix(hdr.Name, "api.json") {
				var buf bytes.Buffer
				if _, err := io.Copy(&buf, tr); err != nil {
					openApiDocErr = fmt.Errorf("copy api.json from tar: %w", err)
					return
				}

				loader := openapi3.NewLoader()
				openApiDoc, openApiDocErr = loader.LoadFromData(buf.Bytes())
				return
			}
		}
	})

	return openApiDoc, openApiDocErr
}

func GetOpenApiResource(resourcePath string) (*openapi3.PathItem, error) {

	// Normalize path to ensure format like /users/
	resourcePath = "/" + strings.Trim(resourcePath, "/") + "/"

	doc, err := loadOpenAPIDocOnce()
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI document: %w", err)
	}

	resource := doc.Paths.Map()[resourcePath]
	if resource == nil {
		// Collect all available paths
		var available []string
		for path := range doc.Paths.Map() {
			available = append(available, path)
		}

		return nil, fmt.Errorf(
			"path %q not found in OpenAPI schema. Available paths:\n  - %s",
			resourcePath,
			strings.Join(available, "\n  - "),
		)
	}

	return resource, nil
}

func GetOpenApiComponents() (*openapi3.Components, error) {
	doc, err := loadOpenAPIDocOnce()

	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI document: %w", err)
	}

	if doc.Components == nil {
		return nil, fmt.Errorf("OpenAPI document has no components defined")
	}

	return doc.Components, nil
}

func GetOpenApiComponentSchema(ref string) (*openapi3.SchemaRef, error) {
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		ref = parts[len(parts)-1]
	} else {
		panic("invalid schema reference: " + ref)
	}
	components, err := GetOpenApiComponents()
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAPI components: %w", err)
	}
	schemaRef := components.Schemas[ref]
	return schemaRef, nil
}

// GetSchema_POST_RequestBody extracts the request body schema from a POST operation.
// It expects the request body to be defined with content type "application/json".
// Returns the schema reference for the POST body payload.
// Returns an error if the POST operation or its schema is not properly defined.
func GetSchema_POST_RequestBody(resourcePath string) (*openapi3.SchemaRef, error) {
	resource, err := GetOpenApiResource(resourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAPI resource %q: %w", resourcePath, err)
	}

	if resource == nil || resource.Post == nil || resource.Post.RequestBody == nil ||
		resource.Post.RequestBody.Value == nil {
		return &openapi3.SchemaRef{Value: &openapi3.Schema{}}, nil
	}

	// Try application/json, then fallback to */*
	content := resource.Post.RequestBody.Value.Content["application/json"]
	if content == nil {
		content = resource.Post.RequestBody.Value.Content["*/*"]
	}
	if content == nil || content.Schema == nil {
		return &openapi3.SchemaRef{Value: &openapi3.Schema{}}, nil
	}

	// Resolve and compose if necessary
	final := ResolveComposedSchema(ResolveAllRefs(content.Schema))
	return &openapi3.SchemaRef{Value: final}, nil
}

// GetSchema_PATCH_RequestBody extracts the request body schema from a PATCH operation.
// Returns an empty schema if PATCH or application/json content is missing.
func GetSchema_PATCH_RequestBody(resourcePath string) (*openapi3.SchemaRef, error) {
	resource, err := GetOpenApiResource(resourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAPI resource %q: %w", resourcePath, err)
	}

	if resource == nil || resource.Patch == nil || resource.Patch.RequestBody == nil ||
		resource.Patch.RequestBody.Value == nil {
		return &openapi3.SchemaRef{Value: &openapi3.Schema{}}, nil
	}

	content := resource.Patch.RequestBody.Value.Content["application/json"]
	if content == nil {
		content = resource.Patch.RequestBody.Value.Content["*/*"]
	}
	if content == nil || content.Schema == nil {
		return &openapi3.SchemaRef{Value: &openapi3.Schema{}}, nil
	}

	final := ResolveComposedSchema(ResolveAllRefs(content.Schema))
	return &openapi3.SchemaRef{Value: final}, nil
}

// GetSchema_POST_StatusOk extracts the schema from a POST operation's response,
// checking status codes 200, 201, 202 (in that order of preference).
// It returns the schema if available under "application/json" content type.
func GetSchema_POST_StatusOk(resourcePath string) (*openapi3.SchemaRef, error) {
	resource, err := GetOpenApiResource(resourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAPI resource %q: %w", resourcePath, err)
	}

	if resource == nil || resource.Post == nil {
		return &openapi3.SchemaRef{Value: &openapi3.Schema{}}, nil
	}

	for _, code := range []int{200, 201, 202} {
		resp := resource.Post.Responses.Status(code)
		schemaRef := extractSchemaFromResponse(resp)
		if schemaRef != nil {
			final := ResolveComposedSchema(ResolveAllRefs(schemaRef))
			return &openapi3.SchemaRef{Value: final}, nil
		}
	}

	return nil, fmt.Errorf(
		"no valid schema found in POST response (200/201/202) for resource %s", resourcePath,
	)
}

// GetSchema_GET_StatusOk attempts to extract the item schema from a GET 200 response.
// It supports paginated (results[]), flat list ([]), and single-object responses.
func GetSchema_GET_StatusOk(resourcePath string) (*openapi3.SchemaRef, error) {
	resource, err := GetOpenApiResource(resourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAPI resource %q: %w", resourcePath, err)
	}

	if resource == nil || resource.Get == nil {
		return &openapi3.SchemaRef{Value: &openapi3.Schema{}}, nil
	}

	resp := resource.Get.Responses.Status(200)
	if resp == nil || resp.Value == nil {
		return nil, fmt.Errorf("GET missing 200 response for resource %s", resourcePath)
	}

	content := resp.Value.Content["application/json"]
	if content == nil || content.Schema == nil {
		return nil, fmt.Errorf("GET response missing or malformed schema")
	}

	rootSchema := ResolveComposedSchema(ResolveAllRefs(content.Schema))

	// Case 1: paginated response { "results": [...] }
	if results, ok := rootSchema.Properties["results"]; ok && results != nil {
		resolvedResults := ResolveComposedSchema(ResolveAllRefs(results))
		if resolvedResults.Type != nil && len(*resolvedResults.Type) > 0 && (*resolvedResults.Type)[0] == "array" {
			if resolvedResults.Items != nil {
				item := ResolveComposedSchema(ResolveAllRefs(resolvedResults.Items))
				return &openapi3.SchemaRef{Value: item}, nil
			}
			return nil, fmt.Errorf("GET response 'results' array has no items schema")
		}
		return nil, fmt.Errorf("GET response 'results' is not an array")
	}

	// Case 2: root is array itself
	if rootSchema.Type != nil && len(*rootSchema.Type) > 0 && (*rootSchema.Type)[0] == "array" {
		if rootSchema.Items != nil {
			item := ResolveComposedSchema(ResolveAllRefs(rootSchema.Items))
			return &openapi3.SchemaRef{Value: item}, nil
		}
		return nil, fmt.Errorf("GET root array has no items schema")
	}

	// Case 3: single object
	return &openapi3.SchemaRef{Value: rootSchema}, nil
}

// QueryParametersGET extracts all query parameters from the GET operation of a given OpenAPI path item.
// It returns a slice of openapi3.Parameter objects whose "in" field is "query".
// These typically represent optional or required query string inputs accepted by the endpoint.
//
// Parameters:
//   - resource: an *openapi3.PathItem representing a specific OpenAPI path (e.g., "/users/").
//
// Returns:
//   - []*openapi3.Parameter: a slice of query parameter definitions.
//   - error: if the GET operation is not defined.
func QueryParametersGET(resourcePath string) ([]*openapi3.Parameter, error) {
	resource, err := GetOpenApiResource(resourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAPI resource %q: %w", resourcePath, err)
	}

	if resource == nil || resource.Get == nil {
		// No GET operation — treat as no query parameters
		return []*openapi3.Parameter{}, nil
	}

	var queryParams []*openapi3.Parameter
	for _, paramRef := range resource.Get.Parameters {
		if paramRef == nil || paramRef.Value == nil {
			continue
		}
		if strings.EqualFold(paramRef.Value.In, "query") {
			queryParams = append(queryParams, paramRef.Value)
		}
	}

	return queryParams, nil
}

// extractSchemaFromResponse attempts to extract an application/json schema from a response.
func extractSchemaFromResponse(resp *openapi3.ResponseRef) *openapi3.SchemaRef {
	if resp == nil || resp.Value == nil {
		return nil
	}
	content := resp.Value.Content["application/json"]
	if content == nil || content.Schema == nil {
		return nil
	}
	return content.Schema
}

// SearchableQueryParams returns only query parameters that are primitive types
// (string, integer) from the GET operation of the given resource path.
func SearchableQueryParams(resourcePath string) ([]string, error) {
	params, err := QueryParametersGET(resourcePath)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, p := range params {
		if p == nil || p.Schema == nil || p.Schema.Value == nil {
			continue
		}
		schema := p.Schema.Value

		// Skip primitive or read-only fields
		if !isStringOrInteger(schema) || schema.ReadOnly {
			continue
		}

		result = append(result, p.Name)
	}

	return result, nil
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

func ResolveComposedSchema(schema *openapi3.Schema) *openapi3.Schema {
	if schema == nil {
		return nil
	}
	// Resolve allOf first, regardless of whether Type is set on the current schema.
	if len(schema.AllOf) > 0 {
		merged := &openapi3.Schema{
			Properties:   map[string]*openapi3.SchemaRef{},
			Required:     []string{},
			Title:        schema.Title,
			Description:  schema.Description,
			ExternalDocs: schema.ExternalDocs,
		}

		// First, copy properties from the original schema itself
		for name, prop := range schema.Properties {
			merged.Properties[name] = prop
		}
		merged.Required = append(merged.Required, schema.Required...)
		if schema.Type != nil && len(*schema.Type) > 0 {
			merged.Type = schema.Type
		}

		// Then, merge properties from allOf sub-schemas
		for _, subRef := range schema.AllOf {
			// Resolve refs and also compose nested allOf/anyOf/oneOf
			sub := ResolveComposedSchema(ResolveAllRefs(subRef))
			if sub == nil {
				continue
			}
			for name, prop := range sub.Properties {
				merged.Properties[name] = prop
			}
			merged.Required = append(merged.Required, sub.Required...)
			if sub.Type != nil && len(*sub.Type) > 0 {
				merged.Type = sub.Type
			}
		}
		return merged
	}

	// If there is no composition to resolve, return as-is.
	if schema.Type != nil && len(*schema.Type) > 0 {
		return schema
	}

	// Resolve oneOf or anyOf by picking the first resolvable schema with a type
	for _, refList := range [][]*openapi3.SchemaRef{schema.OneOf, schema.AnyOf} {
		for _, subRef := range refList {
			sub := ResolveAllRefs(subRef)
			if sub != nil && sub.Type != nil && len(*sub.Type) > 0 {
				return sub
			}
		}
	}
	return schema
}

func ResolveAllRefs(ref *openapi3.SchemaRef) *openapi3.Schema {
	seen := map[string]bool{}
	for ref != nil && ref.Ref != "" && !seen[ref.Ref] {
		seen[ref.Ref] = true
		ref, _ = GetOpenApiComponentSchema(ref.Ref)
	}
	if ref == nil || ref.Value == nil {
		panic(fmt.Sprintf("cannot resolve final schema from ref: %+v", ref))
	}
	return ref.Value
}
