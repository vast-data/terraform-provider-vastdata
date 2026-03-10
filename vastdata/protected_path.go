// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var ProtectedPathSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"protectedpaths",
	http.MethodGet,
	"protectedpaths",
)

type ProtectedPath struct {
	tfstate *is.TFState
}

func (m *ProtectedPath) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &ProtectedPath{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:            ProtectedPathSchemaRef,
			ExcludedSchemaFields: []string{"members_info"},
			PreserveUserValueFields: []string{
				"role",
				"target_exported_dir",
			},
			AdditionalSchemaAttributes: map[string]any{
				// TERF-186
				"estimated_read_only_time": rschema.StringAttribute{
					Computed:    true,
					Description: "",
				},
			},
		},
	)}
}

func (m *ProtectedPath) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ProtectedPath{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: ProtectedPathSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"estimated_read_only_time": dschema.StringAttribute{
					Computed:    true,
					Description: "",
				},
			},
		},
	)}
}

func (m *ProtectedPath) TfState() *is.TFState {
	return m.tfstate
}

func (m *ProtectedPath) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.ProtectedPaths
}
