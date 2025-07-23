// Copyright (c) HashiCorp, Inc.

package internalstate

// TFStateHints defines metadata and overrides used during schema generation for
// Terraform resources and data sources. These hints allow customizing required,
// optional, excluded, and searchable fields beyond what is defined in the OpenAPI schema.
type TFStateHints struct {
	// TFStateHintsForCustom - presence indicates whether the resource or data source is custom-defined
	// rather than generated from OpenAPI definitions. Schema generation from OpenAPI
	// for such resources is skipped, and the hints are used to define the schema.
	// Use `AdditionalSchemaAttributes` to add custom attributes.
	TFStateHintsForCustom *TFStateHintsForCustom

	// SchemaRef defines where to get request (create) and response (read) schemas from OpenAPI
	SchemaRef *SchemaReference

	// SearchableFields lists field names that should be treated as searchable
	// when constructing lookup parameters (e.g., for API GET calls).
	SearchableFields []string

	// RequiredSchemaFields explicitly marks these fields as required in the
	// Terraform schema, regardless of whether they are marked optional or read-only
	// in the OpenAPI definition.
	RequiredSchemaFields []string

	// NotRequiredSchemaFields forces the specified fields to not be required,
	// even if marked required in the OpenAPI definition or RequiredSchemaFields list.
	NotRequiredSchemaFields []string

	// OptionalSchemaFields explicitly marks these fields as optional in the
	// Terraform schema, even if they are marked required or read-only in OpenAPI.
	OptionalSchemaFields []string

	// NotOptionalSchemaFields disables the optional behavior for the specified fields,
	// even if listed in OptionalSchemaFields or inferred from OpenAPI.
	NotOptionalSchemaFields []string

	// ExcludedSchemaFields lists fields that should be completely excluded from
	// the Terraform schema, regardless of their presence in the OpenAPI definition.
	ExcludedSchemaFields []string

	// ComputedSchemaFields forces the Computed flag for the specified fields,
	ComputedSchemaFields []string

	// NotComputedSchemaFields disables the Computed flag for the specified fields,
	// typically used in data source schemas to mark fields that are not returned
	// by the backend and thus should not be treated as computed.
	NotComputedSchemaFields []string

	// WriteOnlyFields indicates fields only for search IOW only read operations.
	ReadOnlyFields []string

	// WriteOnlyFields indicates fields whose values Terraform will not store
	// in the plan or state artifacts. If a field is write-only, it must be either
	// optional or required. Write-only fields cannot be computed.
	WriteOnlyFields []string

	// EditOnlyFields lists fields that can be updated only during PATCH request.
	// For instance some resources have field "enabled" that cannot be set to false along with create (POST) request.
	EditOnlyFields []string

	// DeleteOnlyFields lists fields that can be provided only during DELETE request.
	DeleteOnlyFields []string

	// PreserveOrderFields defines fields where the order matters (e.g., for lists instead of sets).
	PreserveOrderFields []string

	// SensitiveFields marks fields as sensitive, so their values are redacted
	// from logs and plan output.
	SensitiveFields []string

	// AdditionalSchemaAttributes defines extra schema attributes to inject into
	// the Terraform schema even if they are not present in the OpenAPI schema.
	// The key is the attribute name, and the value is the schema definition.
	AdditionalSchemaAttributes map[string]any

	// CommonValidatorsMapping defines a mapping between resource field names and common validator identifiers.
	//
	// Each entry maps a specific resource field (as a string) to a common validator name (also a string or validator definition).
	// This allows reuse of predefined validator logic across multiple fields or resources.
	//
	// Example:
	//   CommonValidatorsMapping: map[string]string{
	//       "bucket_name": "s3_name",
	//       "fqdn":        "hostname",
	//   }
	// NOTE: All common validators in is in: vastdata/schema_generation/common_validators.go
	CommonValidatorsMapping map[string]string

	// CommonModifiersMapping defines a mapping between resource field names and common modifier names.
	//
	// Each key represents a resource field, and the corresponding value is the name of a predefined modifier
	// function or transformation to apply to that field (e.g., normalization, trimming, lowercasing).
	//
	// This allows centralized reuse of common field modification logic across multiple resources.
	//
	// Example:
	//   CommonModifiersMapping: map[string]string{
	//       "username": "trim_space",
	//       "email":    "to_lower",
	//   }
	// NOTE: All common validators in is in: vastdata/schema_generation/common_modifiers.go
	CommonModifiersMapping map[string]string
}

// SchemaReference encapsulates both create and read endpoints for a resource.
// Used to extract the POST request schema (for resources) and the GET response schema (for resources or data sources).
type SchemaReference struct {
	// Create specifies the OpenAPI endpoint to use for extracting the creation schema (e.g., POST /volumes).
	Create *OpenAPIEndpointRef

	// Read specifies the OpenAPI endpoint to use for extracting the read schema (e.g., GET /volumes/{id}).
	Read *OpenAPIEndpointRef
}

func NewSchemaReference(
	createMethod, createPath string,
	readMethod, readPath string,
) *SchemaReference {
	var createRef, readRef *OpenAPIEndpointRef

	if createMethod != "" && createPath != "" {
		createRef = &OpenAPIEndpointRef{
			Method: createMethod,
			Path:   createPath,
		}
	}

	if readMethod != "" && readPath != "" {
		readRef = &OpenAPIEndpointRef{
			Method: readMethod,
			Path:   readPath,
		}
	}

	return &SchemaReference{
		Create: createRef,
		Read:   readRef,
	}
}

// OpenAPIEndpointRef defines a reference to a specific HTTP method + path
// in the OpenAPI schema, used for schema extraction.
type OpenAPIEndpointRef struct {
	// HTTP method (e.g., "get", "post", "patch")
	Method string
	// Path in OpenAPI (e.g., "/volumes", "/volumes/{id}")
	Path string
}

type TFStateHintsForCustom struct {
	// Description provides a detailed explanation of the resource or data source.
	Description string
	// MarkdownDescription provides a markdown-formatted description for the resource or data source.
	MarkdownDescription string
	// SchemaAttributes defines schema attributes to inject into
	SchemaAttributes map[string]any
}
