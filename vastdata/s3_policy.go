// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
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
			SchemaRef:            S3PolicySchemaRef,
			EditOnlyFields:       []string{"enabled"},
			OptionalSchemaFields: []string{"enabled"},
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
