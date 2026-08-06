// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"time"

	version "github.com/hashicorp/go-version"
)

// TFStateHints defines metadata and overrides used during schema generation for
// Terraform resources and data sources. These hints allow customizing required,
// optional, excluded, and searchable fields beyond what is defined in the OpenAPI schema.
type TFStateHints struct {
	// Importable controls whether the resource supports import operations.
	// If nil or true, the resource is considered importable by default.
	// Set to BoolPtr(false) to explicitly disable import for a resource.
	Importable *bool
	// TFStateHintsForCustom - presence indicates whether the resource or data source is custom-defined
	// rather than generated from OpenAPI definitions. Schema generation from OpenAPI
	// for such resources is skipped, and the hints are used to define the schema.
	// Use `AdditionalSchemaAttributes` to add custom attributes.
	TFStateHintsForCustom *TFStateHintsForCustom

	// SchemaRef defines where to get request (create) and response (read) schemas from OpenAPI
	SchemaRef *SchemaReference

	// ImportFields defines ordered field names used to support composite import IDs.
	// When an import ID is provided without key=value pairs, it will be split by
	// a supported delimiter and mapped to these fields in order. When key=value
	// pairs are provided, any subset and order is accepted; keys must exist in the schema.
	ImportFields []string

	// SearchableFields lists field names that should be treated as searchable
	// when constructing lookup parameters (e.g., for API GET calls).
	SearchableFields []string

	// AllowEmptySearchParams skips the datasource "at least one search field" check.
	// Use for reads that are valid with no selector (e.g. cluster-wide dashboard).
	AllowEmptySearchParams bool

	// SearchFilterFields lists field names that are always appended to search
	// query params as additional AND-filters when they are set (non-null).
	// Unlike SearchableFields (which picks ONE field from a set of alternatives),
	// every non-null field listed here is always included alongside the primary
	// search key.
	SearchFilterFields []string

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

	// ReadOnlyFields indicates fields only for search IOW only read operations.
	ReadOnlyFields []string

	// WriteOnlyFields indicates fields whose values Terraform will not store
	// in the plan or state artifacts. If a field is write-only, it must be either
	// optional or required. Write-only fields cannot be computed.
	WriteOnlyFields []string

	// CreateOnlyFields lists fields that may be set on create (POST) but must never
	// be sent on update (PATCH). Omitting them from configuration after create must
	// not plan a null clear. These fields are also treated as Computed so
	// UseStateForUnknown keeps the prior state value when the attribute is omitted.
	CreateOnlyFields []string

	// EditOnlyFields lists fields that can be updated only during PATCH request.
	// For instance some resources have field "enabled" that cannot be set to false along with create (POST) request.
	EditOnlyFields []string

	// PreserveUserValueFields lists fields whose user-declared values should be preserved
	// during refresh/read operations, rather than being overwritten with API values.
	// By default, all fields are updated from API during refresh to sync external changes.
	// Only fields listed here will preserve their user-configured values.
	// Typical use case: optional fields that should not be overwritten even if the API returns different values.
	PreserveUserValueFields []string

	// PreserveUserValueFieldsWhenApiReturnsNull lists fields whose user-declared values should be preserved
	// It is similar to PreserveUserValueFields:
	// PreserveUserValueFields - means once set value in tf state it will never be changed!
	// PreserveUserValueFieldsWhenApiReturnsNull - once set value can be changed
	// but only if API returns value which is logically not nil value
	PreserveUserValueFieldsWhenApiReturnsNull []string

	// DeleteOnlyBodyFields maps Terraform attribute names to API body field names
	// for fields that are only allowed to be sent in DELETE request bodies.
	// Key: Terraform schema field name; Value: API body field name.
	DeleteOnlyBodyFields map[string]string

	// DeleteOnlyParamFields maps Terraform attribute names to API query parameter names
	// for fields that are only allowed to be sent as DELETE request query params.
	// Key: Terraform schema field name; Value: API query parameter name.
	DeleteOnlyParamFields map[string]string

	// PreserveOrderFields defines fields where the order matters (e.g., for lists instead of sets).
	PreserveOrderFields []string

	// IPSetEquivalenceFields lists fields that contain lists/sets of IP addresses or CIDR
	// ranges.  During read/refresh the API may return a semantically equivalent but
	// differently-formatted value (e.g. individual IPs collapsed into a CIDR block).
	// For these fields the provider expands both the user's current state value and the
	// API response to their full constituent IP sets and compares them.  If the sets are
	// identical the user's original form is preserved in state (no spurious drift).
	// If the sets genuinely differ (e.g. someone added or removed an IP via the UI)
	// the API value is written to state so Terraform correctly detects the change.
	IPSetEquivalenceFields []string

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

	// Skip API calls during refresh and use current tfstate instead.
	// Useful for offline or performance-critical scenarios. Default is false.
	SkipRefreshAPICall bool

	// RetryOn configures retry behaviour for resource create and delete API calls only.
	// It is not used by the datasource read path.
	RetryOn *RetryPolicy

	// AsyncTaskTimeout overrides the default timeout used when waiting
	// for asynchronous task completion after create/update operations.
	// When nil, the provider uses the default timeout (10 minutes).
	AsyncTaskTimeout *time.Duration

	// SubResources declares nested API endpoints whose responses are fetched
	// and merged (flattened) into the parent resource's Terraform state.
	// Schema attributes from each SubResourceHint are automatically injected
	// into the parent schema alongside AdditionalSchemaAttributes.
	SubResources []SubResourceHint
}

// SubResourceHint declares a nested sub-endpoint that is fetched after the
// parent resource is read/created/updated and whose response fields are
// flattened into the parent resource's Terraform state.
type SubResourceHint struct {
	// FieldTrigger is the name of a bool attribute on the parent resource.
	// The sub-resource is fetched only when this field evaluates to true.
	// When empty (and MinVastVersion is also unset) the sub-resource is always fetched.
	FieldTrigger string

	// MinVastVersion is the minimum VAST cluster version from which this sub-resource
	// was introduced. When set, the sub-resource is fetched only if the connected
	// cluster's version is greater than or equal to MinVastVersion. This check is
	// evaluated in addition to FieldTrigger (if both are set, both must pass).
	// When MinVastVersion alone is set (FieldTrigger is empty), the sub-resource is
	// fetched automatically whenever the cluster version meets the minimum.
	MinVastVersion *version.Version

	// SchemaKey is both the URL segment appended after the parent resource's
	// base path and ID (e.g. "s3_true_ip_config" → GET /clusters/{id}/s3_true_ip_config/)
	// and the key used when embedding the sub-record into the parent record.
	// When empty the sub-record is flattened directly into the parent.
	SchemaKey string

	// SchemaAttributes defines the Terraform schema attributes contributed
	// by this sub-resource (including the optional trigger field), flattened
	// into the parent resource's schema.
	SchemaAttributes map[string]any

	// Writable indicates that the sub-resource supports write operations
	// (e.g. POST / DELETE) managed via hooks (AfterCreateResource / AfterUpdateResource).
	// When true and SchemaKey is non-empty, the generated nested SingleNestedAttribute
	// is Optional+Computed instead of Computed-only, allowing users to configure it.
	Writable bool
}

// RetryPolicy groups retry rules for different resource lifecycle operations.
type RetryPolicy struct {
	// Create configures retries for resource creation (custom CreateResource and the
	// default create path).
	Create *RetryExpression

	// Delete configures retries for resource deletion (custom DeleteResource and the
	// default delete path).
	Delete *RetryExpression
}

// RetryExpression defines the conditions and parameters for retrying a failed API request.
// Retries are triggered when the API returns one of the configured StatusCodes and, if
// BodyContains is non-empty, at least one of the listed substrings is found in the response body.
type RetryExpression struct {
	// StatusCodes lists the HTTP response status codes that should trigger a retry.
	// At least one code must match for a retry to occur.
	StatusCodes []int

	// BodyContains is an optional list of substrings to search for in the response body.
	// When non-empty, at least one substring must be present in the response body for
	// a retry to be triggered (in addition to the status code check).
	BodyContains []string

	// Times is the maximum number of attempts (including the first).
	// Defaults to 5 when zero or negative.
	Times int

	// SleepSeconds is the number of seconds to sleep between retry attempts.
	// Defaults to 10 when zero or negative.
	SleepSeconds int
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
