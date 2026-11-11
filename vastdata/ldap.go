// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var LdapSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"ldaps",
	http.MethodGet,
	"ldaps",
)

type Ldap struct {
	tfstate is.ResourceState
}

func (m *Ldap) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Ldap{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:               LdapSchemaRef,
			NotComputedSchemaFields: []string{"bindpw"},
			SearchableFields:        []string{"domain_name"},
			WriteOnlyFields:         []string{"bindpw"},
			WriteOnlyMD5Fields:      []string{"bindpw"},
			SensitiveFields:         []string{"bindpw"},
			AdditionalSchemaAttributes: map[string]any{
				"bindpw_md5": ContentMD5SchemaAttribute("bindpw", false),
			},
		},
	)}
}

func (m *Ldap) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Ldap{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: LdapSchemaRef,
		},
	)}
}

func (m *Ldap) TfState() is.ResourceState {
	return m.tfstate
}

func (m *Ldap) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Ldaps
}

func (m *Ldap) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) error {
	SetPlanContentMD5(ctx, req, resp, "bindpw", "bindpw_md5")
	return nil
}

func (m *Ldap) TransformResponseRecord(record Record) Record {
	DeleteRecordFields(record, "bindpw")
	return record
}
