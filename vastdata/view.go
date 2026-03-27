// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
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
			SchemaRef:             ViewSchemaRef,
			DeleteOnlyBodyFields:  map[string]string{"delete_dir": ""},
			DeleteOnlyParamFields: map[string]string{"force": "force"},
			ImportFields:          []string{"path", "tenant_name"},
			// qos_policy: removing Computed allows Terraform
			// to plan a change when the user sets them to null. (TERF-225)
			NotComputedSchemaFields: []string{"qos_policy"},
			CommonValidatorsMapping: map[string]string{
				"path":                     ValidatorPathStartsWithSlash,
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
				"force": rschema.BoolAttribute{
					Optional:    true,
					Description: "Force View removal.",
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
		path := tfstate.String("path")
		var tenantId int64
		if tfstate.IsKnownAndNotNull("tenant_id") {
			tenantId = tfstate.Int64("tenant_id")
		}
		if _, err = rest.Folders.FolderDeleteFolderWithContext_DELETE(ctx, params{"path": path, "tenant_id": tenantId}); isApiError(err) {
			body := err.(*ApiError).Body
			if strings.Contains(body, "no such directory") {
				return nil
			}
		}
	}
	return err

}
