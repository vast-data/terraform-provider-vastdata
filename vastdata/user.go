// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var UserSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"users",
	http.MethodGet,
	"users",
)

type User struct {
	tfstate is.ResourceState
}

func (m *User) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &User{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:          UserSchemaRef,
			WriteOnlyFields:    []string{"password"},
			WriteOnlyMD5Fields: []string{"password"},
			SensitiveFields:    []string{"password"},
			AdditionalSchemaAttributes: map[string]any{
				"password_md5": ContentMD5SchemaAttribute("password", false),
				"s3_policies_ids": rschema.SetAttribute{
					ElementType: types.Int64Type,
					Optional:    true,
					Computed:    true,
					Description: "S3 policies IDs, denoting which S3 identity policies are associated with the user. The user is granted and denied S3 permissions according to the associated S3 identity policies.",
				},
			},
		},
	)}
}

func (m *User) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &User{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: UserSchemaRef,
		}),
	}
}

func (m *User) TfState() is.ResourceState {
	return m.tfstate
}

func (m *User) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Users
}

func (m *User) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) error {
	SetPlanContentMD5(ctx, req, resp, "password", "password_md5")
	return nil
}

func (m *User) TransformResponseRecord(record Record) Record {
	DeleteRecordFields(record, "password")
	return record
}
