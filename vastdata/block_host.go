// Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var BlockHostSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"blockhosts",
	http.MethodGet,
	"blockhosts",
)

type BlockHost struct {
	tfstate is.ResourceState
}

func (m *BlockHost) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &BlockHost{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: BlockHostSchemaRef,
		},
	)}
}

func (m *BlockHost) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &BlockHost{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: BlockHostSchemaRef,
		},
	)}
}

func (m *BlockHost) TfState() is.ResourceState {
	return m.tfstate
}

func (m *BlockHost) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.BlockHosts
}
