// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	planmodifiers "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
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
			Importable: &notImportable,
			TFStateHintsForCustom: &is.TFStateHintsForCustom{
				Description: "Make a Folder Read-Only",
				SchemaAttributes: map[string]any{
					"path": rschema.StringAttribute{
						Required:    true,
						Description: "Path of the folder to be read-only.",
						PlanModifiers: []planmodifiers.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^/([^/].*/)?$`),
								"must start and end with '/'",
							),
						},
					},
					"tenant_id": rschema.Int64Attribute{
						Required:    true,
						Description: "ID of the tenant to which the folder belongs.",
						PlanModifiers: []planmodifiers.Int64{
							int64planmodifier.RequiresReplace(),
						},
					},
				},
			},
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

func (m *FolderReadOnly) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	searchParams := getSearchParams(ctx, m.tfstate, nil)
	return rest.Folders.FolderReadOnlyWithContext_GET(ctx, searchParams)
}

func (m *FolderReadOnly) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	record, err := m.ReadDatasource(ctx, rest)
	if err != nil {
		if err = ignoreStatusCodes(err, http.StatusBadRequest, http.StatusNotFound); err == nil {
			return nil, ForceCleanState{}
		}
	}
	return record, err
}

func (m *FolderReadOnly) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	searchParams := getSearchParams(ctx, m.tfstate, nil)
	record, err := rest.Folders.FolderReadOnlyWithContext_POST(ctx, searchParams)
	return record, err
}

func (m *FolderReadOnly) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	searchParams := getSearchParams(ctx, m.tfstate, nil)
	record, err := rest.Folders.FolderReadOnlyWithContext_POST(ctx, searchParams)
	return record, err
}

func (m *FolderReadOnly) DeleteResource(ctx context.Context, rest *VMSRest) error {
	searchParams := getSearchParams(ctx, m.tfstate, nil)
	err := rest.Folders.FolderReadOnlyWithContext_DELETE(ctx, searchParams)
	return err
}
