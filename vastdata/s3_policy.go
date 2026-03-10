// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var S3PolicySchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"s3policies",
	http.MethodGet,
	"s3policies",
)

type S3Policy struct {
	tfstate *is.TFState
}

func (m *S3Policy) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &S3Policy{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			PreserveUserValueFields: []string{"policy"},
			SchemaRef:               S3PolicySchemaRef,
			EditOnlyFields:          []string{"enabled"},
			OptionalSchemaFields:    []string{"enabled"},
		},
	)}
}

func (m *S3Policy) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &S3Policy{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: S3PolicySchemaRef,
		}),
	}
}

func (m *S3Policy) TfState() *is.TFState {
	return m.tfstate
}

func (m *S3Policy) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.S3Policies
}

// TransformResponseRecord deduplicates the "users" and "groups" list fields
// in the backend response.  The API may return duplicate entries (e.g. multiple
// empty strings) which cause Terraform to reject the state with
// "Duplicate Set Element" errors.
func (m *S3Policy) TransformResponseRecord(record Record) Record {
	for _, field := range []string{"users", "groups"} {
		if val, ok := record[field]; ok {
			if list, ok := val.([]any); ok {
				record[field] = deleteDuplicates(list)
			}
		}
	}
	return record
}
