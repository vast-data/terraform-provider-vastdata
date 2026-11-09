// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var CallhomeConfigSchemaRef = is.NewSchemaReference(
	http.MethodPatch,
	"callhomeconfigs/{id}",
	http.MethodGet,
	"callhomeconfigs",
)

var callhomeConfigRegisterTimeout = 10 * time.Minute

type CallhomeConfig struct {
	tfstate is.ResourceState
}

func (m *CallhomeConfig) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &CallhomeConfig{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:  CallhomeConfigSchemaRef,
			Importable: &notImportable, // write-only secrets (*_md5) cannot be restored from the API
			// uplink_email is stored (not write-only) so changing it alone plans an
			// update and re-runs register-cluster.
			WriteOnlyFields: []string{
				"proxy_password",
				"aws_s3_ak",
				"aws_s3_sk",
				"cloud_api_key",
				"uplink_password",
				"uplink_token",
			},
			WriteOnlyMD5Fields: []string{
				"proxy_password",
				"aws_s3_ak",
				"aws_s3_sk",
				"cloud_api_key",
				"uplink_password",
				"uplink_token",
			},
			SensitiveFields: []string{
				"proxy_password",
				"aws_s3_ak",
				"aws_s3_sk",
				"cloud_api_key",
				"uplink_password",
				"uplink_token",
			},
			AdditionalSchemaAttributes: map[string]any{
				"id": rschema.Int64Attribute{
					Required:    true,
					Description: "ID of the Call Home configuration (typically 1).",
				},
				// GET-only (not in PATCH body) — keep out of Optional without changing
				// global PATCH schema generation.
				"cloud_registered": rschema.BoolAttribute{
					Computed:    true,
					Description: "Whether the cluster is registered with VAST Cloud Services / Uplink.",
				},
				"callhome_upload_destination": rschema.StringAttribute{
					Computed:    true,
					Description: "Callhome upload destination.",
				},
				"ssl_certificate": rschema.StringAttribute{
					Computed:    true,
					Description: "SSL Certificate file content, including the BEGIN CERTIFICATE and END CERTIFICATE lines.",
				},
				"uplink_email": rschema.StringAttribute{
					Optional:    true,
					Description: "VAST Uplink super user login email.",
				},
				"uplink_password": rschema.StringAttribute{
					Optional:    true,
					Sensitive:   true,
					WriteOnly:   true,
					Description: "VAST Uplink super user password.",
				},
				"uplink_token": rschema.StringAttribute{
					Optional:    true,
					Sensitive:   true,
					WriteOnly:   true,
					Description: "VAST Uplink cluster registration token.",
				},
				"proxy_password_md5":  ContentMD5SchemaAttribute("proxy_password", false),
				"aws_s3_ak_md5":       ContentMD5SchemaAttribute("aws_s3_ak", false),
				"aws_s3_sk_md5":       ContentMD5SchemaAttribute("aws_s3_sk", false),
				"cloud_api_key_md5":   ContentMD5SchemaAttribute("cloud_api_key", false),
				"uplink_password_md5": ContentMD5SchemaAttribute("uplink_password", false),
				"uplink_token_md5":    ContentMD5SchemaAttribute("uplink_token", false),
			},
		},
	)}
}

func (m *CallhomeConfig) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &CallhomeConfig{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: CallhomeConfigSchemaRef,
			AdditionalSchemaAttributes: map[string]any{
				"id": dschema.Int64Attribute{
					Required:    true,
					Description: "ID of the Call Home configuration (typically 1).",
				},
			},
		},
	)}
}

func (m *CallhomeConfig) TfState() is.ResourceState {
	return m.tfstate
}

func (m *CallhomeConfig) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.CallhomeConfigs
}

// ModifyPlan fingerprints write-only secrets so rotations plan an update
// (TERF-317: proxy_password and other secrets were invisible to the plan).
func (m *CallhomeConfig) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) error {
	SetPlanContentMD5Fields(ctx, req, resp,
		"proxy_password", "proxy_password_md5",
		"aws_s3_ak", "aws_s3_ak_md5",
		"aws_s3_sk", "aws_s3_sk_md5",
		"cloud_api_key", "cloud_api_key_md5",
		"uplink_password", "uplink_password_md5",
		"uplink_token", "uplink_token_md5",
	)
	return nil
}

func (m *CallhomeConfig) TransformResponseRecord(record Record) Record {
	DeleteRecordFields(record,
		"proxy_password",
		"aws_s3_ak",
		"aws_s3_sk",
		"cloud_api_key",
		"uplink_password",
		"uplink_token",
	)
	return record
}

func (m *CallhomeConfig) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	id := m.tfstate.Int64("id")
	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}
	return rest.CallhomeConfigs.GetByIdWithContext(ctx, id)
}

func (m *CallhomeConfig) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.ReadResource(ctx, rest)
}

// CreateResource adopts the existing Call Home config (API has no POST), PATCHes
// configured fields, then optionally registers the cluster with Uplink.
func (m *CallhomeConfig) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return ensureCallhomeConfigUpdatedWith(ctx, m.tfstate, m.tfstate, rest)
}

func (m *CallhomeConfig) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	planTs := plan.(*CallhomeConfig).TfState()
	if err := ensureNotChanged(m.tfstate, planTs, "id"); err != nil {
		return nil, err
	}
	return ensureCallhomeConfigUpdatedWith(ctx, m.tfstate, planTs, rest)
}

func (m *CallhomeConfig) DeleteResource(_ context.Context, _ *VMSRest) error {
	// No-op. Call Home configs cannot be deleted; destroy only clears Terraform state.
	return nil
}

func ensureCallhomeConfigUpdatedWith(ctx context.Context, stateTs, fieldsTs is.ResourceState, rest *VMSRest) (DisplayableRecord, error) {
	id := stateTs.Int64("id")
	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}

	if _, err := rest.CallhomeConfigs.GetByIdWithContext(ctx, id); err != nil {
		return nil, fmt.Errorf("callhome config %d not found (API has no create): %w", id, err)
	}

	body := fieldsTs.GetCreateParams()
	if stateTs != fieldsTs {
		// GetUpdateParams respects WriteOnlyMD5Fields so unchanged secrets are
		// not re-PATCHed after overlayWriteOnlyFromConfig.
		body = fieldsTs.GetUpdateParams(stateTs)
	}
	body.Without(
		"id",
		"uplink_email",
		"uplink_password",
		"uplink_token",
		"proxy_password_md5",
		"aws_s3_ak_md5",
		"aws_s3_sk_md5",
		"cloud_api_key_md5",
		"uplink_password_md5",
		"uplink_token_md5",
	)
	// Never PATCH null secrets — VMS rejects e.g. cloud_api_key=null.
	for key, value := range body {
		if value == nil {
			delete(body, key)
		}
	}
	if len(body) > 0 {
		if _, err := rest.CallhomeConfigs.UpdateWithContext(ctx, id, body); err != nil {
			return nil, fmt.Errorf("failed to update callhome config %d: %w", id, err)
		}
	}

	if err := registerCallhomeClusterIfNeeded(ctx, id, fieldsTs, rest); err != nil {
		return nil, err
	}

	return rest.CallhomeConfigs.GetByIdWithContext(ctx, id)
}

func registerCallhomeClusterIfNeeded(ctx context.Context, id int64, ts is.ResourceState, rest *VMSRest) error {
	email, password, token := "", "", ""
	if ts.IsKnownAndNotNull("uplink_email") {
		email = ts.String("uplink_email")
	}
	if ts.IsKnownAndNotNull("uplink_password") {
		password = ts.String("uplink_password")
	}
	if ts.IsKnownAndNotNull("uplink_token") {
		token = ts.String("uplink_token")
	}
	if email == "" && password == "" && token == "" {
		return nil
	}

	body := params{}
	if email != "" {
		body["email"] = email
	}
	if password != "" {
		body["password"] = password
	}
	if token != "" {
		body["token"] = token
	}

	asyncResult, err := rest.CallhomeConfigs.CallhomeConfigsRegisterClusterWithContext_PATCH(
		ctx, id, body, callhomeConfigRegisterTimeout,
	)
	if err != nil {
		return fmt.Errorf("failed to register cluster with Uplink: %w", err)
	}
	if asyncResult != nil && asyncResult.IsFailed() {
		return fmt.Errorf("Uplink register-cluster task failed: %v", asyncResult.Err)
	}
	return nil
}
