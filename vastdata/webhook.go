// Copyright (c) HashiCorp, Inc.
package provider

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var WebhookSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"webhooks",
	http.MethodGet,
	"webhooks",
)

type Webhook struct {
	tfstate *is.TFState
}

func (m *Webhook) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Webhook{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: WebhookSchemaRef,
		},
	)}
}

func (m *Webhook) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Webhook{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: WebhookSchemaRef,
		},
	)}
}

func (m *Webhook) TfState() *is.TFState {
	return m.tfstate
}

func (m *Webhook) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.WebHooks
}
