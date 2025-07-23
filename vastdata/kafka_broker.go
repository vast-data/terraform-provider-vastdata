// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var KafkaBrokerSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"kafkabrokers",
	http.MethodGet,
	"kafkabrokers",
)

type KafkaBroker struct {
	tfstate *is.TFState
}

func (m *KafkaBroker) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &KafkaBroker{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: KafkaBrokerSchemaRef,
		},
	)}
}

func (m *KafkaBroker) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &KafkaBroker{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: KafkaBrokerSchemaRef,
		},
	)}
}

func (m *KafkaBroker) TfState() *is.TFState {
	return m.tfstate
}

func (m *KafkaBroker) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.KafkaBrokers
}
