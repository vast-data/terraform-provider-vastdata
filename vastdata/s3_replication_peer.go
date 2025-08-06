// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var S3ReplicationPeerSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"replicationtargets",
	http.MethodGet,
	"replicationtargets",
)

type S3ReplicationPeer struct {
	tfstate *is.TFState
}

func (m *S3ReplicationPeer) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &S3ReplicationPeer{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: S3ReplicationPeerSchemaRef,
		},
	)}
}

func (m *S3ReplicationPeer) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &S3ReplicationPeer{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: S3ReplicationPeerSchemaRef,
		},
	)}
}

func (m *S3ReplicationPeer) TfState() *is.TFState {
	return m.tfstate
}

func (m *S3ReplicationPeer) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.S3replicationPeers
}
