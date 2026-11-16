// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var LocalS3KeySchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"locals3keys",
	http.MethodGet,
	"locals3keys",
)

type LocalS3Key struct {
	tfstate is.ResourceState
}

func (m *LocalS3Key) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &LocalS3Key{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: LocalS3KeySchemaRef,
		},
	)}
}

func (m *LocalS3Key) TfState() is.ResourceState {
	return m.tfstate
}

func (m *LocalS3Key) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.LocalS3Keys
}
