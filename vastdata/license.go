// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var LicenseSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"licenses",
	http.MethodGet,
	"licenses",
)

type License struct {
	tfstate *is.TFState
}

func (m *License) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &License{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:            LicenseSchemaRef,
			ImportFields:         []string{"id"},
			RequiredSchemaFields: []string{"license_key"},
			WriteOnlyFields:      []string{"license_key"},
			SensitiveFields:      []string{"license_key"},
			// license_key is write-only (not in state).
			//Track a computed MD5 of the key and
			// force replace when that changes.
			AdditionalSchemaAttributes: map[string]any{
				"license_key_md5": ContentMD5SchemaAttribute("license_key", true),
			},
		},
	)}
}

func (m *License) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &License{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:        LicenseSchemaRef,
			SearchableFields: []string{"guid", "name"},
			// Response metadata is read-only; only guid/name (plus id) are search params.
			NotOptionalSchemaFields: []string{
				"accelerated_cpu_cores",
				"capacity",
				"cpu_cores",
				"expiration_date",
				"license_key",
				"period",
				"start_date",
				"state",
			},
		},
	)}
}

func (m *License) TfState() *is.TFState {
	return m.tfstate
}

func (m *License) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Licenses
}

func (m *License) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return rest.Licenses.CreateWithContext(ctx, m.tfstate.GetCreateParams())
}

// TransformResponseRecord derives license_key_md5 for drift detection, then drops
// license_key so the secret is never written into Terraform state.
func (m *License) TransformResponseRecord(record Record) Record {
	SetRecordContentMD5(record, "license_key", "license_key_md5")
	delete(record, "license_key")
	return record
}

// ModifyPlan sets license_key_md5 from the write-only license_key in config so
func (m *License) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) error {
	SetPlanContentMD5(ctx, req, resp, "license_key", "license_key_md5")
	return nil
}
