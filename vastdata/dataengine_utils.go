// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	planmodifiers "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/vast-data/go-vast-client/core"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

const deTenantNameAttr = "tenant_name"

// deTenantNameSchemaAttribute is the optional ForceNew tenant scoping field for
// Data Engine APIs that require ?tenant_name= when using cluster-admin auth.
func deTenantNameSchemaAttribute() rschema.StringAttribute {
	return rschema.StringAttribute{
		Optional: true,
		Description: "Tenant name for Data Engine API queries. Required for cluster-admin " +
			"auth; optional when the provider already authenticates as that tenant.",
		PlanModifiers: []planmodifiers.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
}

// resolveTenantQuery returns {"tenant_name": ...} from resource state, then the
// VMS session config tenant (provider block / VASTDATA_TENANT). Empty when unset.
func resolveTenantQuery(tfState is.ResourceState, rest *VMSRest) params {
	if tfState != nil && tfState.HasAttribute(deTenantNameAttr) && tfState.IsKnownAndNotNull(deTenantNameAttr) {
		if name := strings.TrimSpace(tfState.String(deTenantNameAttr)); name != "" {
			return params{deTenantNameAttr: name}
		}
	}
	if rest != nil && rest.GetSession() != nil {
		if cfg := rest.GetSession().GetConfig(); cfg != nil {
			if name := strings.TrimSpace(cfg.Tenant); name != "" {
				return params{deTenantNameAttr: name}
			}
		}
	}
	return params{}
}

// deTransformQueryParams implements TransformQueryParams for Data Engine
// resources: merge tenant_name from resolveTenantQuery without overwriting a
// caller-supplied tenant_name.
func deTransformQueryParams(query params, tfState is.ResourceState, rest *VMSRest) params {
	if query == nil {
		query = params{}
	}
	query.Update(resolveTenantQuery(tfState, rest), true)
	return query
}

// deTransformRequestBody strips provider-only tenant_name from API request bodies.
func deTransformRequestBody(body params) params {
	if body != nil {
		body.Without(deTenantNameAttr)
	}
	return body
}

// ensureTenantNameInRecord injects tenant_name from resolveTenantQuery when the
// Data Engine API response omitted it. Uses resource state first, then the VMS
// session tenant (provider block / VASTDATA_TENANT). No-op when the record
// already has a non-empty tenant_name, or when neither source provides one.
func ensureTenantNameInRecord(rec Record, tfState is.ResourceState, rest *VMSRest) Record {
	if rec == nil {
		return nil
	}
	if v, ok := rec[deTenantNameAttr]; ok {
		if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
			return rec
		}
		if v != nil {
			if _, ok := v.(string); !ok {
				return rec
			}
		}
	}
	q := resolveTenantQuery(tfState, rest)
	name, _ := q[deTenantNameAttr].(string)
	if strings.TrimSpace(name) == "" {
		return rec
	}
	out := Record{}
	for k, v := range rec {
		out[k] = v
	}
	out[deTenantNameAttr] = name
	return out
}

func isDataEnginePresent(setupStatus string) bool {
	switch strings.ToLower(strings.TrimSpace(setupStatus)) {
	case "completed", "complete", "ready", "in_progress", "running", "pending":
		return true
	default:
		return false
	}
}

// findRecordByLocalFilter lists a collection with listQuery and matches
// guid/vrn/name in-memory. Use only where List ignores those query params
// (kubernetes-clusters, container-registries, mtls-authentication-credentials).
// Callers should pass TransformQueryParams(..., rest) as listQuery.
func findRecordByLocalFilter(
	ctx context.Context,
	api VastResourceAPIWithContext,
	tfState is.ResourceState,
	listQuery params,
	resourcePath string,
) (DisplayableRecord, error) {
	records, err := api.ListWithContext(ctx, listQuery)
	if err != nil {
		return nil, err
	}
	want := params{}
	for _, key := range []string{"guid", "vrn", "name"} {
		if tfState != nil && tfState.HasAttribute(key) && tfState.IsKnownAndNotNull(key) {
			want[key] = is.ConvertAttrValueToRaw(tfState.Get(key), tfState.Type(key))
		}
	}
	filtered := filterRecordSetByAttrs(records, want)
	switch len(filtered) {
	case 0:
		return nil, notFoundErr(resourcePath, want)
	case 1:
		if filtered[0].Empty() {
			return nil, notFoundErr(resourcePath, want)
		}
		return filtered[0], nil
	default:
		return nil, &core.TooManyRecordsError{ResourcePath: resourcePath, Params: want}
	}
}

// deleteRecordByLocalFilter finds the row with findRecordByLocalFilter, then
// DELETE /{guid}/. Ignores not-found. query should include tenant_name when required.
func deleteRecordByLocalFilter(
	ctx context.Context,
	api VastResourceAPIWithContext,
	tfState is.ResourceState,
	query params,
	resourcePath string,
) error {
	rec, err := findRecordByLocalFilter(ctx, api, tfState, query, resourcePath)
	if err != nil {
		_, err = ignoreNotFound(nil, err)
		return err
	}
	guid, ok := rec.(Record)["guid"]
	if !ok || guid == nil || fmt.Sprint(guid) == "" {
		return fmt.Errorf("%s: record has no guid; refusing delete by numeric id", resourcePath)
	}
	_, err = api.DeleteByIdWithContext(ctx, guid, query, nil)
	return ignoreStatusCodes(err, http.StatusNotFound)
}
