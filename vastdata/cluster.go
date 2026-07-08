// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var ClusterSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"clusters",
	http.MethodGet,
	"clusters",
)

// clusterS3TrueIPSubResource declares the read-only sub-resource for the
// /clusters/{id}/s3_true_ip_config/ endpoint (VAST >= 5.5.0, GET only).
// Set get_s3_true_ip_config = false to explicitly opt out of fetching.
var clusterS3TrueIPSubResource = is.SubResourceHint{
	MinVastVersion: VastVersion550,
	FieldTrigger:   "get_s3_true_ip_config",
	SchemaKey:      "s3_true_ip_config",
	SchemaAttributes: map[string]any{
		"s3_true_client_ip_header": rschema.StringAttribute{
			Computed:    true,
			Description: "True-client-IP header name for S3 requests.",
		},
		"s3_included_addresses": rschema.ListNestedAttribute{
			Computed:    true,
			Description: "IP address ranges included in the S3 True-IP configuration.",
			NestedObject: rschema.NestedAttributeObject{
				Attributes: map[string]rschema.Attribute{
					"start_ip": rschema.StringAttribute{
						Computed:    true,
						Description: "Starting IP address of the range.",
					},
					"range": rschema.Int64Attribute{
						Computed:    true,
						Description: "Number of addresses in the range.",
					},
				},
			},
		},
	},
}

type Cluster struct {
	tfstate *is.TFState
}

func (m *Cluster) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Cluster{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:    ClusterSchemaRef,
			SubResources: []is.SubResourceHint{clusterS3TrueIPSubResource},
		},
	)}
}

func (m *Cluster) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Cluster{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:    ClusterSchemaRef,
			SubResources: []is.SubResourceHint{clusterS3TrueIPSubResource},
		},
	)}
}

func (m *Cluster) TfState() *is.TFState {
	return m.tfstate
}

func (m *Cluster) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Clusters
}

// GetSubResources fetches /clusters/{id}/s3_true_ip_config/ on VAST clusters
// running version >= 5.5.0. Set get_s3_true_ip_config = false to opt out.
func (m *Cluster) GetSubResources(ctx context.Context, rest *VMSRest, record Record, clusterVersion *version.Version) (Record, error) {
	if clusterVersion == nil || clusterVersion.LessThan(VastVersion550) {
		return nil, nil
	}
	// User explicitly opted out.
	if m.tfstate.IsKnownAndNotNull("get_s3_true_ip_config") && !m.tfstate.Bool("get_s3_true_ip_config") {
		return nil, nil
	}

	id := record.RecordID()

	rec, err := rest.Clusters.ClusterS3TrueIpConfigWithContext_GET(ctx, id)
	if err != nil {
		return nil, err
	}

	config := map[string]any{
		"s3_true_client_ip_header": rec["true_client_ip_header"],
		"s3_included_addresses":    parseIncludedAddresses(rec["included_addresses"]),
	}
	return Record{clusterS3TrueIPSubResource.SchemaKey: config}, nil
}

func parseIncludedAddresses(raw any) []map[string]any {
	items, ok := raw.([]any)
	if !ok || len(items) == 0 {
		return nil // preserve null, not []
	}
	addrs := make([]map[string]any, 0, len(items))
	for _, item := range items {
		entry, ok := item.(map[string]any)
		if !ok {
			continue
		}
		addrs = append(addrs, map[string]any{"start_ip": entry["start_ip"], "range": entry["range"]})
	}
	return addrs
}
