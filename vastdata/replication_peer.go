// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var ReplicationPeersSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"nativereplicationremotetargets",
	http.MethodGet,
	"nativereplicationremotetargets",
)

type ReplicationPeer struct {
	tfstate *is.TFState
}

func (m *ReplicationPeer) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &ReplicationPeer{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: ReplicationPeersSchemaRef,
			// 503 SERVICE_UNAVAILABLE / HANDSHAKE_IN_PROGRESS is transient and occurs when
			// two peers are created concurrently. Retry until the handshake completes.
			RetryOn: &is.RetryExpression{
				StatusCodes: []int{503},
			},
		},
	)}
}

func (m *ReplicationPeer) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ReplicationPeer{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: ReplicationPeersSchemaRef,
		}),
	}
}

func (m *ReplicationPeer) TfState() *is.TFState {
	return m.tfstate
}

func (m *ReplicationPeer) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.ReplicationPeers
}
