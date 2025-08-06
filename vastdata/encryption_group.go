// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var EncryptionGroupSchemaRef = is.NewSchemaReference(
	"",
	"",
	http.MethodGet,
	"encryptiongroups",
)

type EncryptionGroup struct {
	tfstate *is.TFState
}

func (m *EncryptionGroup) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &EncryptionGroup{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: EncryptionGroupSchemaRef,
		},
	)}
}

func (m *EncryptionGroup) TfState() *is.TFState {
	return m.tfstate
}

func (m *EncryptionGroup) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.EncryptionGroups
}
