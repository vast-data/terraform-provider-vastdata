// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
	"strings"
)

var ViewSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"views",
	http.MethodGet,
	"views",
)

type View struct {
	tfstate *is.TFState
}

func (m *View) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &View{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:            ViewSchemaRef,
			DeleteOnlyBodyFields: map[string]string{"delete_dir": ""},
			ImportFields:         []string{"path", "tenant_name"},
			CommonValidatorsMapping: map[string]string{
				"path":                     ValidatorPathStartsWithSlash,
				"alias":                    ValidatorPathStartsWithSlash,
				"max_retention_period":     ValidatorRetentionFormat,
				"min_retention_period":     ValidatorRetentionFormat,
				"default_retention_period": ValidatorRetentionFormat,
				"auto_commit":              ValidatorRetentionFormat,
			},
			AdditionalSchemaAttributes: map[string]any{
				"delete_dir": rschema.BoolAttribute{
					Optional: true,
					Description: "If set to true during view deletion, the underlying directory will also be deleted. " +
						"This behavior is only effective during delete operations. " +
						"For it to work properly, the Trash API must be enabled on the VAST cluster.",
				},
			},
		},
	)}
}

func (m *View) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &View{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: ViewSchemaRef,
		}),
	}
}

func (m *View) TfState() *is.TFState {
	return m.tfstate
}

func (m *View) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Views
}

func (m *View) PrepareDeleteResource(ctx context.Context, rest *VMSRest) error {
	tfstate := m.tfstate
	var err error
	if tfstate.IsKnownAndNotNull("delete_dir") && tfstate.Bool("delete_dir") {
		// If delete_dir is true, we delete the directory.
		deleteParams, _ := tfstate.SetIfAvailable("path", "tenant_id")
		if _, err = rest.Folders.DeleteFolderWithContext(ctx, deleteParams); isApiError(err) {
			body := err.(*ApiError).Body
			if strings.Contains(body, "no such directory") {
				return nil
			}
		}
	}
	return err

}
