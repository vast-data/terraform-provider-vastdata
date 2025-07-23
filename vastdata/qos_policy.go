// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

const QosPolicyResourcePath = "qospolicies"

var QosPolicySchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"qospolicies",
	http.MethodGet,
	"qospolicies",
)

type QosPolicy struct {
	tfstate *is.TFState
}

func (m *QosPolicy) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &QosPolicy{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: QosPolicySchemaRef,
		},
	)}
}

func (m *QosPolicy) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &QosPolicy{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: QosPolicySchemaRef,
		}),
	}
}

func (m *QosPolicy) TfState() *is.TFState {
	return m.tfstate
}

func (m *QosPolicy) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.QosPolicies
}
