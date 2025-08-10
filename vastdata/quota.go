// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var QuotaSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"quotas",
	http.MethodGet,
	"quotas",
)

type Quota struct {
	tfstate *is.TFState
}

func (m *Quota) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Quota{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: QuotaSchemaRef,
		},
	)}
}

func (m *Quota) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Quota{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: QuotaSchemaRef,
		}),
	}
}

func (m *Quota) TfState() *is.TFState {
	return m.tfstate
}

func (m *Quota) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Quotas
}
