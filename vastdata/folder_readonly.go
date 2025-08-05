// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var FolderReadOnlySchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"folders/read_only",
	http.MethodGet,
	"folders/read_only",
)

type FolderReadOnly struct {
	tfstate *is.TFState
}

func (m *FolderReadOnly) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &FolderReadOnly{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: FolderReadOnlySchemaRef,
		},
	)}
}

func (m *FolderReadOnly) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &FolderReadOnly{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: FolderReadOnlySchemaRef,
		},
	)}
}

func (m *FolderReadOnly) TfState() *is.TFState {
	return m.tfstate
}

func (m *FolderReadOnly) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Folders
}

// TODO: implemente read/create/update/delete methods for FolderReadOnly
func (m *FolderReadOnly) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return nil, nil
}

func (m *FolderReadOnly) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return nil, nil
}

func (m *FolderReadOnly) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	return nil, nil

}

func (m *FolderReadOnly) DeleteResource(ctx context.Context, rest *VMSRest) error {
	return nil
}
