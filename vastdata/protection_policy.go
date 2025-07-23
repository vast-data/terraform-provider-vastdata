// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var ProtectionPolicySchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"protectionpolicies",
	http.MethodGet,
	"protectionpolicies",
)

type ProtectionPolicy struct {
	tfstate *is.TFState
}

func (m *ProtectionPolicy) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &ProtectionPolicy{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: ProtectionPolicySchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				// NOTE: original fields from OpenAPI spec with "-" is not acceptable in Terraform schema.
				// We replace "frames" property 'in-place' here.
				"frames": rschema.ListNestedAttribute{
					Computed:    true,
					Optional:    true,
					Description: "Defines the schedule for snapshot creation and the local and remote retention policies. Example: every 90m start-at 2025-07-27 20:10:35 keep-local 10h keep-remote 30d",
					NestedObject: rschema.NestedAttributeObject{
						Attributes: map[string]rschema.Attribute{
							"keep_remote": rschema.StringAttribute{
								Computed:    true,
								Optional:    true,
								Description: "Remote retention period (e.g., '30d').",
							},
							"start_at": rschema.StringAttribute{
								Computed:    true,
								Optional:    true,
								Description: "Start time for the snapshot schedule (e.g., '2025-07-27 20:10:35', in UTC).",
							},
							"every": rschema.StringAttribute{
								Computed:    true,
								Optional:    true,
								Description: "Snapshot frequency (e.g., '1d' or '12h').",
							},
							"keep_local": rschema.StringAttribute{
								Computed:    true,
								Optional:    true,
								Description: "Local retention period (e.g., '7d').",
							},
						},
					},
				},
			},
		},
	)}
}

func (m *ProtectionPolicy) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ProtectionPolicy{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: ProtectionPolicySchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				// NOTE: original fields from OpenAPI spec with "-" is not acceptable in Terraform schema.
				// We replace "frames" property 'in-place' here.
				"frames": dschema.ListNestedAttribute{
					Computed:    true,
					Description: "Defines the schedule for snapshot creation and the local and remote retention policies. Example: every 90m start-at 2025-07-27 20:10:35 keep-local 10h keep-remote 30d",
					NestedObject: dschema.NestedAttributeObject{
						Attributes: map[string]dschema.Attribute{
							"keep_remote": dschema.StringAttribute{
								Computed:    true,
								Description: "Remote retention period (e.g., '30d').",
							},
							"start_at": dschema.StringAttribute{
								Computed:    true,
								Description: "Start time for the snapshot schedule (e.g., '2025-07-27 20:10:35', in UTC).",
							},
							"every": dschema.StringAttribute{
								Computed:    true,
								Description: "Snapshot frequency (e.g., '1d' or '12h').",
							},
							"keep_local": dschema.StringAttribute{
								Computed:    true,
								Description: "Local retention period (e.g., '7d').",
							},
						},
					},
				},
			},
		},
	)}
}

func (m *ProtectionPolicy) TfState() *is.TFState {
	return m.tfstate
}

func (m *ProtectionPolicy) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.ProtectionPolicies
}

// TransformRequestBody applies a recursive key transformation to the "frames" field
// in the outgoing request body, converting all map keys from snake_case to dash-case.
// This ensures compatibility with the backend API, which expects keys like "start-at"
// instead of "start_at".
func (m *ProtectionPolicy) TransformRequestBody(body params) params {
	if frames, ok := body["frames"]; ok {
		body["frames"] = convertMapKeysRecursive(frames, underscoreToDash)
	}
	return body
}

// TransformResponseRecord applies a recursive key transformation to the "frames" field
// in the backend response, converting all map keys from dash-case to snake_case.
// This ensures the data conforms to Terraform schema expectations (snake_case).
func (m *ProtectionPolicy) TransformResponseRecord(record Record) Record {
	if frames, ok := record["frames"]; ok {
		record["frames"] = convertMapKeysRecursive(frames, dashToUnderscore)
	}
	return record
}
