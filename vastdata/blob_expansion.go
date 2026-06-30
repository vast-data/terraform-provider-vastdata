// Copyright (c) HashiCorp, Inc.
package provider

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var BlobExpansionSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"blobexpansions",
	http.MethodGet,
	"blobexpansions/show",
)

type BlobExpansion struct {
	tfstate *is.TFState
}

func (m *BlobExpansion) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &BlobExpansion{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: BlobExpansionSchemaRef,
		},
	)}
}

func (m *BlobExpansion) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &BlobExpansion{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: BlobExpansionSchemaRef,
		},
	)}
}

func (m *BlobExpansion) TfState() *is.TFState {
	return m.tfstate
}

func (m *BlobExpansion) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.BlobExpansions
}
