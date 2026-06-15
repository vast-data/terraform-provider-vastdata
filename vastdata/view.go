// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"
	"strings"

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
var viewS3CorsSubResource = is.SubResourceHint{
	SchemaKey: "s3cors_configuration",
	Writable:  true,
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
			EditOnlyFields: []string{"enable_nfs4_triggers"},
			SubResources:   []is.SubResourceHint{viewS3CorsSubResource},
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
				"enable_nfs4_triggers": rschema.BoolAttribute{
					Optional:    true,
					Computed:    true,
					Description: "When true, enables NFSv4 triggers for this view (/views/{id}/nfs4_triggers/).",
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
			SchemaRef: ViewSchemaRef,
		}),
	}
}

func (m *View) TfState() *is.TFState {
	return m.tfstate
}

func (m *View) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Views
}

func (m *View) GetSubResources(ctx context.Context, rest *VMSRest, record Record) (Record, error) {
	id, ok := record["id"]
	if !ok || id == nil || id == "" {
		return nil, nil
	}

	result := Record{}

	// --- nfs4_triggers ---
	// Only applicable to NFS4 views;
	if isNFS4View(record) {
		nfs4Rec, err := rest.Views.ViewNfs4TriggersWithContext_GET(ctx, id)
		if err = ignoreStatusCodes(err, http.StatusNotFound, http.StatusForbidden); err != nil {
			return nil, err
		}
		if nfs4Rec != nil {
			result["enable_nfs4_triggers"] = nfs4Rec["enabled"]
		} else {
			result["enable_nfs4_triggers"] = nil
		}
	} else {
		result["enable_nfs4_triggers"] = nil
	}

	// --- s3cors_configuration ---
	corsRec, err := rest.Views.ViewS3corsConfigurationWithContext_GET(ctx, id, params{})
	if err != nil && !isNotFoundErr(err) {
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

// AfterCreateResource enables nfs4_triggers and/or creates s3cors_configuration
// when the user specifies those fields in the Terraform configuration.
func (m *View) AfterCreateResource(ctx context.Context, rest *VMSRest, record Record) error {
	id := record["id"]

	if m.tfstate.IsKnownAndNotNull("enable_nfs4_triggers") && m.tfstate.Bool("enable_nfs4_triggers") {
		if err := rest.Views.ViewNfs4TriggersWithContext_POST(ctx, id, params{"enabled": true}); err != nil {
			return err
		}
	}

	if body, ok := m.corsBody(); ok && body != nil {
		if _, err := rest.Views.ViewS3corsConfigurationWithContext_POST(ctx, id, body); err != nil {
			return err
		}
	}
	return nil
}

// AfterUpdateResource syncs nfs4_triggers and s3cors_configuration with the plan state.
// If s3cors_configuration is removed from config, all CORS rules are deleted.
// plan is the desired state (what the user configured); m is the prior state.
func (m *View) AfterUpdateResource(ctx context.Context, plan AfterUpdateResource, rest *VMSRest, record Record) error {
	planView := plan.(*View)
	id := record["id"]

	// Sync enable_nfs4_triggers
	if planView.tfstate.IsKnownAndNotNull("enable_nfs4_triggers") {
		enabled := planView.tfstate.Bool("enable_nfs4_triggers")
		if err := rest.Views.ViewNfs4TriggersWithContext_POST(ctx, id, params{"enabled": enabled}); err != nil {
			return err
		}
	}

	// Sync s3cors_configuration
	rawCors, hasCors := planView.tfstate.Raw["s3cors_configuration"]
	if !hasCors || rawCors.IsNull() || rawCors.IsUnknown() {
		// User removed s3cors_configuration — delete all CORS rules.
		if err := rest.Views.ViewS3corsConfigurationWithContext_DELETE(ctx, id); err != nil && !isNotFoundErr(err) {
			return err
		}
	} else {
		if body, ok := planView.corsBody(); ok && body != nil {
			if _, err := rest.Views.ViewS3corsConfigurationWithContext_POST(ctx, id, body); err != nil {
				return err
			}
		}
	}
	return nil
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
			if i, ok := v.(types.Int64); ok && i.ValueInt64() != 0 {
				rule["max_age_seconds"] = i.ValueInt64()
			}
		}
		rules = append(rules, rule)
	}
	return params{"cors_rules": rules}, true
}

// isNFS4View returns true when the view record includes "NFS4" in its protocols list.
func isNFS4View(record Record) bool {
	raw, ok := record["protocols"]
	if !ok || raw == nil {
		return false
	}
	list, ok := raw.([]any)
	if !ok {
		return false
	}
	for _, p := range list {
		if s, ok := p.(string); ok && s == "NFS4" {
			return true
		}
	}
	return false
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
