// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var S3LifeCycleRuleSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"s3lifecyclerules",
	http.MethodGet,
	"s3lifecyclerules",
)

type S3LifeCycleRule struct {
	tfstate *is.TFState
}

func (m *S3LifeCycleRule) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &S3LifeCycleRule{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: S3LifeCycleRuleSchemaRef,
		},
	)}
}

func (m *S3LifeCycleRule) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &S3LifeCycleRule{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: S3LifeCycleRuleSchemaRef,
		}),
	}
}

func (m *S3LifeCycleRule) TfState() *is.TFState {
	return m.tfstate
}

func (m *S3LifeCycleRule) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.S3LifeCycleRules
}
