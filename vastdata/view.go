// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var ViewSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"views",
	http.MethodGet,
	"views",
)

// viewS3CorsSubResource declares the nested writable sub-resource for /views/{id}/s3cors_configuration/.
// Fetched automatically on VAST clusters running >= 5.5.0.
// Set get_s3cors_configuration = false to explicitly opt out.
var viewS3CorsSubResource = is.SubResourceHint{
	MinVastVersion: VastVersion550,
	FieldTrigger:   "get_s3cors_configuration",
	SchemaKey:      "s3cors_configuration",
	Writable:       true,
	SchemaAttributes: map[string]any{
		"cors_rules": rschema.ListNestedAttribute{
			Optional:    true,
			Computed:    true,
			Description: "S3 CORS rules for this view.",
			NestedObject: rschema.NestedAttributeObject{
				Attributes: map[string]rschema.Attribute{
					"allowed_methods": rschema.ListAttribute{
						Required:    true,
						ElementType: types.StringType,
						Description: "CORS allowed HTTP methods (e.g. GET, POST).",
					},
					"allowed_origins": rschema.ListAttribute{
						Required:    true,
						ElementType: types.StringType,
						Description: "CORS allowed origins.",
					},
					"allowed_headers": rschema.ListAttribute{
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
						Description: "CORS allowed request headers.",
					},
					"expose_headers": rschema.ListAttribute{
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
						Description: "Headers the browser may expose to the client-side script.",
					},
					"max_age_seconds": rschema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: "Time in seconds to cache the CORS preflight response.",
					},
				},
			},
		},
	},
}

type View struct {
	tfstate *is.TFState
}

func (m *View) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &View{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:             ViewSchemaRef,
			DeleteOnlyBodyFields:  map[string]string{"delete_dir": ""},
			DeleteOnlyParamFields: map[string]string{"force": "force"},
			ImportFields:          []string{"path", "tenant_name"},
			// qos_policy: removing Computed allows Terraform
			// to plan a change when the user sets them to null. (TERF-225)
			NotComputedSchemaFields: []string{"qos_policy", "qos_policy_id"},
			// s3cors_configuration is auto-excluded via SubResources hints.
			SubResources: []is.SubResourceHint{viewS3CorsSubResource},
			CommonValidatorsMapping: map[string]string{
				"path":                     ValidatorPathStartsWithSlash,
				"max_retention_period":     ValidatorRetentionFormat,
				"min_retention_period":     ValidatorRetentionFormat,
				"default_retention_period": ValidatorRetentionFormat,
				"auto_commit":              ValidatorRetentionFormat,
			},
			AdditionalSchemaAttributes: map[string]any{
				"delete_dir": rschema.BoolAttribute{
					Optional: true,
					Description: "If set to true during view deletion, the underlying directory will also be deleted. " +
						"This behavior is only effective during delete operations. " +
						"For it to work properly, the Trash API must be enabled on the VAST cluster.",
				},
				"force": rschema.BoolAttribute{
					Optional:    true,
					Description: "Force View removal.",
				},
			},
		},
	)}
}

func (m *View) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &View{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:    ViewSchemaRef,
			SubResources: []is.SubResourceHint{viewS3CorsSubResource},
		}),
	}
}

func (m *View) TfState() *is.TFState {
	return m.tfstate
}

func (m *View) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Views
}

// GetSubResources fetches /views/{id}/s3cors_configuration/ on VAST clusters
// running version >= 5.5.0. Set get_s3cors_configuration = false to opt out.
// Only applicable to S3 bucket views; non-S3 views are skipped.
func (m *View) GetSubResources(ctx context.Context, rest *VMSRest, record Record, clusterVersion *version.Version) (Record, error) {
	if clusterVersion == nil || clusterVersion.LessThan(VastVersion550) {
		return nil, nil
	}
	// User explicitly opted out.
	if m.tfstate.IsKnownAndNotNull("get_s3cors_configuration") && !m.tfstate.Bool("get_s3cors_configuration") {
		return nil, nil
	}
	// S3 CORS is only valid for S3 bucket views.
	if !isS3View(record) {
		return nil, nil
	}

	id := record.RecordID()
	result := Record{}

	// --- s3cors_configuration ---
	corsRec, err := rest.Views.ViewS3corsConfigurationWithContext_GET(ctx, id, nil)
	// 404: no CORS config; 400: view is not an S3 bucket.
	if err = ignoreStatusCodes(err, http.StatusNotFound, http.StatusBadRequest); err != nil {
		return nil, err
	}
	if corsRec != nil {
		if rules, ok := corsRec["cors_rules"]; ok {
			if items, isSlice := rules.([]any); isSlice && len(items) > 0 {
				result["s3cors_configuration"] = map[string]any{"cors_rules": rules}
			} else {
				result["s3cors_configuration"] = nil
			}
		} else {
			result["s3cors_configuration"] = nil
		}
	} else {
		result["s3cors_configuration"] = nil
	}

	return result, nil
}

// corsVersionCheck returns an error when the connected cluster does not support
// s3cors_configuration (requires VAST >= 5.5.0).
func corsVersionCheck(ctx context.Context, rest *VMSRest) error {
	clusterVer, err := GetCachedClusterVersion(ctx, rest)
	if err != nil {
		return fmt.Errorf("s3cors_configuration: failed to retrieve cluster version: %w", err)
	}
	if clusterVer.LessThan(VastVersion550) {
		return fmt.Errorf(
			"s3cors_configuration requires VAST cluster version >= %s, but cluster is running %s",
			VastVersion550, clusterVer,
		)
	}
	return nil
}

// AfterCreateResource creates s3cors_configuration when the user specifies it.
// Returns an error if the cluster version does not support this sub-resource.
// Skipped entirely when get_s3cors_configuration is explicitly set to false.
func (m *View) AfterCreateResource(ctx context.Context, rest *VMSRest, record Record) error {
	if m.tfstate.IsKnownAndNotNull("get_s3cors_configuration") && !m.tfstate.Bool("get_s3cors_configuration") {
		return nil
	}
	body, hasCors := m.corsBody()
	if !hasCors || body == nil {
		return nil
	}
	if !isS3View(record) {
		return fmt.Errorf(`s3cors_configuration requires protocols to include "S3"`)
	}
	if err := corsVersionCheck(ctx, rest); err != nil {
		return err
	}
	id := record.RecordID()
	_, err := rest.Views.ViewS3corsConfigurationWithContext_POST(ctx, id, body)
	return err
}

// AfterUpdateResource syncs s3cors_configuration with the plan state.
// Returns an error if the user provides s3cors_configuration on an unsupported cluster.
// If s3cors_configuration is removed from config, CORS rules are deleted (version permitting).
// Skipped entirely when get_s3cors_configuration is explicitly set to false in the plan.
// plan is the desired state (what the user configured); m is the prior state.
func (m *View) AfterUpdateResource(ctx context.Context, plan AfterUpdateResource, rest *VMSRest, record Record) error {
	planView := plan.(*View)
	id := record.RecordID()

	// Respect explicit opt-out on the plan side.
	if planView.tfstate.IsKnownAndNotNull("get_s3cors_configuration") && !planView.tfstate.Bool("get_s3cors_configuration") {
		return nil
	}

	rawCors, hasCors := planView.tfstate.Raw["s3cors_configuration"]
	if !hasCors || rawCors.IsNull() || rawCors.IsUnknown() {
		// Only delete when prior state actually had CORS configured.
		// Otherwise every update of a view without CORS would hit DELETE
		// (and non-S3 views return 400 "not an S3 bucket").
		if !m.hasCorsConfigured() || !isS3View(record) {
			return nil
		}
		clusterVer, err := GetCachedClusterVersion(ctx, rest)
		if err != nil || clusterVer.LessThan(VastVersion550) {
			return nil // endpoint not available on this cluster, nothing to delete
		}
		err = rest.Views.ViewS3corsConfigurationWithContext_DELETE(ctx, id)
		return ignoreStatusCodes(err, http.StatusNotFound, http.StatusBadRequest)
	}

	// User has s3cors_configuration block — only valid on S3 bucket views.
	if !isS3View(record) {
		return fmt.Errorf(`s3cors_configuration requires protocols to include "S3"`)
	}
	if err := corsVersionCheck(ctx, rest); err != nil {
		return err
	}
	if body, ok := planView.corsBody(); ok && body != nil {
		if _, err := rest.Views.ViewS3corsConfigurationWithContext_POST(ctx, id, body); err != nil {
			return err
		}
	}
	return nil
}

// isS3View returns true when the view record includes "S3" in its protocols list.
func isS3View(record Record) bool {
	raw, ok := record["protocols"]
	if !ok || raw == nil {
		return false
	}
	switch list := raw.(type) {
	case []any:
		for _, p := range list {
			if s, ok := p.(string); ok && s == "S3" {
				return true
			}
		}
	case []string:
		for _, s := range list {
			if s == "S3" {
				return true
			}
		}
	}
	return false
}

// hasCorsConfigured reports whether s3cors_configuration is present and non-null in tfstate.
func (m *View) hasCorsConfigured() bool {
	return m.tfstate != nil && m.tfstate.Enabled && m.tfstate.IsKnownAndNotNull("s3cors_configuration")
}

// corsBody converts the s3cors_configuration from tfstate into a params map
// suitable for POST /views/{id}/s3cors_configuration/.
// Returns (nil, false) when s3cors_configuration is absent/null,
// (nil, true) when present but has no rules, and (body, true) otherwise.
func (m *View) corsBody() (params, bool) {
	rawCors, hasCors := m.tfstate.Raw["s3cors_configuration"]
	if !hasCors || rawCors.IsNull() || rawCors.IsUnknown() {
		return nil, false
	}
	corsObj, ok := rawCors.(types.Object)
	if !ok {
		return nil, false
	}
	corsAttrs := corsObj.Attributes()
	rawRules, hasRules := corsAttrs["cors_rules"]
	if !hasRules || rawRules.IsNull() || rawRules.IsUnknown() {
		return nil, true
	}
	rulesList, ok := rawRules.(types.List)
	if !ok || len(rulesList.Elements()) == 0 {
		return nil, true
	}

	extractStringList := func(obj types.Object, key string) []string {
		v, ok := obj.Attributes()[key]
		if !ok || v.IsNull() || v.IsUnknown() {
			return nil
		}
		list, ok := v.(types.List)
		if !ok {
			return nil
		}
		result := make([]string, 0, len(list.Elements()))
		for _, e := range list.Elements() {
			if s, ok := e.(types.String); ok {
				result = append(result, s.ValueString())
			}
		}
		return result
	}

	rules := make([]any, 0, len(rulesList.Elements()))
	for _, elem := range rulesList.Elements() {
		obj, ok := elem.(types.Object)
		if !ok {
			continue
		}
		rule := make(map[string]any)
		if v := extractStringList(obj, "allowed_methods"); v != nil {
			rule["allowed_methods"] = v
		}
		if v := extractStringList(obj, "allowed_origins"); v != nil {
			rule["allowed_origins"] = v
		}
		if v := extractStringList(obj, "allowed_headers"); v != nil {
			rule["allowed_headers"] = v
		}
		if v := extractStringList(obj, "expose_headers"); v != nil {
			rule["expose_headers"] = v
		}
		if v, ok := obj.Attributes()["max_age_seconds"]; ok && !v.IsNull() && !v.IsUnknown() {
			if i, ok := v.(types.Int64); ok {
				rule["max_age_seconds"] = i.ValueInt64()
			}
		}
		rules = append(rules, rule)
	}
	return params{"cors_rules": rules}, true
}

func (m *View) PrepareDeleteResource(ctx context.Context, rest *VMSRest) error {
	tfstate := m.tfstate
	var err error
	if tfstate.IsKnownAndNotNull("delete_dir") && tfstate.Bool("delete_dir") {
		// If delete_dir is true, we delete the directory.
		path := tfstate.String("path")
		var tenantId int64
		if tfstate.IsKnownAndNotNull("tenant_id") {
			tenantId = tfstate.Int64("tenant_id")
		}
		if _, err = rest.Folders.FolderDeleteFolderWithContext_DELETE(ctx, params{"path": path, "tenant_id": tenantId}); isApiError(err) {
			body := err.(*ApiError).Body
			if strings.Contains(body, "no such directory") {
				return nil
			}
		}
	}
	return err

}
