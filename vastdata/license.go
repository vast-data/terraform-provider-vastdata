// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
			// No PATCH on /licenses/{id}/ — changing the key requires replace.
			CommonModifiersMapping: map[string]string{
				"license_key": ModifierForceNew,
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

// TransformResponseRecord drops license_key from GET responses so the secret
// is never written into Terraform state (write-only on create only).
func (m *License) TransformResponseRecord(record Record) Record {
	delete(record, "license_key")
	return record
}
