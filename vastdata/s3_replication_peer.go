// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var S3ReplicationPeerSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"replicationtargets",
	http.MethodGet,
	"replicationtargets",
)

type S3ReplicationPeer struct {
	tfstate is.ResourceState
}

func (m *S3ReplicationPeer) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &S3ReplicationPeer{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:          S3ReplicationPeerSchemaRef,
			WriteOnlyFields:    []string{"access_key", "secret_key"},
			WriteOnlyMD5Fields: []string{"access_key", "secret_key"},
			SensitiveFields:    []string{"access_key", "secret_key"},
			AdditionalSchemaAttributes: map[string]any{
				"access_key_md5": ContentMD5SchemaAttribute("access_key", false),
				"secret_key_md5": ContentMD5SchemaAttribute("secret_key", false),
			},
		},
	)}
}

func (m *S3ReplicationPeer) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &S3ReplicationPeer{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: S3ReplicationPeerSchemaRef,
		},
	)}
}

func (m *S3ReplicationPeer) TfState() is.ResourceState {
	return m.tfstate
}

func (m *S3ReplicationPeer) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.S3replicationPeers
}

func (m *S3ReplicationPeer) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) error {
	SetPlanContentMD5Fields(ctx, req, resp,
		"access_key", "access_key_md5",
		"secret_key", "secret_key_md5",
	)
	return nil
}

func (m *S3ReplicationPeer) TransformResponseRecord(record Record) Record {
	DeleteRecordFields(record, "access_key", "secret_key")
	return record
}
