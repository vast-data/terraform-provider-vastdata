// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var BgpConfigSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"bgpconfigs",
	http.MethodGet,
	"bgpconfigs",
)

type BgpConfig struct {
	tfstate is.ResourceState
}

func (m *BgpConfig) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &BgpConfig{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:          BgpConfigSchemaRef,
			WriteOnlyFields:    []string{"md5_password"},
			WriteOnlyMD5Fields: []string{"md5_password"},
			SensitiveFields:    []string{"md5_password"},
			AdditionalSchemaAttributes: map[string]any{
				"md5_password_md5": ContentMD5SchemaAttribute("md5_password", false),
			},
		},
	)}
}

func (m *BgpConfig) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &BgpConfig{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: BgpConfigSchemaRef,
		},
	)}
}

func (m *BgpConfig) TfState() is.ResourceState {
	return m.tfstate
}

func (m *BgpConfig) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.BGPConfigs
}

func (m *BgpConfig) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) error {
	SetPlanContentMD5(ctx, req, resp, "md5_password", "md5_password_md5")
	return nil
}

func (m *BgpConfig) TransformResponseRecord(record Record) Record {
	DeleteRecordFields(record, "md5_password")
	return record
}
