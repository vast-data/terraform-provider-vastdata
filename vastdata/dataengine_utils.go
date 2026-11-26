// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	planmodifiers "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vast-data/go-vast-client/core"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

const (
	deTenantNameAttr            = "tenant_name"
	deKubernetesClusterVrnAttr  = "kubernetes_cluster_vrn"
	deKubernetesClusterNameAttr = "kubernetes_cluster_name"
	deContainerRegistryVrnAttr  = "container_registry_vrn"
	deContainerRegistryNameAttr = "container_registry_name"
	deCustomExtensionsAttr      = "custom_extensions"
)

// deCustomExtensionsSchemaAttribute is a flat string map. OpenAPI only declares
// type:object (no properties/additionalProperties), so codegen skips it.
func deCustomExtensionsSchemaAttribute() rschema.MapAttribute {
	return rschema.MapAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Description: "Custom key/value pairs passed as event output extensions. " +
			"Flat map only (no nesting). Keys must be lowercase letters and digits, " +
			"start with a letter, and be 1–20 characters.",
	}
}

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

func isDataEngineVrn(s string) bool {
	return strings.HasPrefix(strings.TrimSpace(s), "vast:dataengine:")
}

// deKubernetesClusterVrnSchemaAttribute is optional+computed: filled from
// kubernetes_cluster_name when that alternative is used.
func deKubernetesClusterVrnSchemaAttribute() rschema.StringAttribute {
	return rschema.StringAttribute{
		Optional: true,
		Computed: true,
		Description: "Kubernetes cluster VRN (vast:dataengine:kubernetes-clusters:...). " +
			"Exactly one of kubernetes_cluster_vrn or kubernetes_cluster_name must be set.",
		PlanModifiers: []planmodifiers.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

// deKubernetesClusterNameSchemaAttribute is a provider-only alternative to
// kubernetes_cluster_vrn; resolved before create/update and not sent to the API.
func deKubernetesClusterNameSchemaAttribute() rschema.StringAttribute {
	return rschema.StringAttribute{
		Optional: true,
		Description: "Kubernetes cluster name. Resolved to a VRN before create/update. " +
			"Exactly one of kubernetes_cluster_vrn or kubernetes_cluster_name must be set.",
	}
}

// resolveKubernetesClusterVRN resolves a cluster name to its VRN.
// If nameOrVrn is already a VRN, it is returned unchanged.
func resolveKubernetesClusterVRN(
	ctx context.Context,
	api VastResourceAPIWithContext,
	listQuery params,
	nameOrVrn string,
) (string, error) {
	nameOrVrn = strings.TrimSpace(nameOrVrn)
	if nameOrVrn == "" {
		return "", fmt.Errorf("kubernetes cluster name or VRN is required")
	}
	if isDataEngineVrn(nameOrVrn) {
		return nameOrVrn, nil
	}
	records, err := api.ListWithContext(ctx, listQuery)
	if err != nil {
		return "", fmt.Errorf("failed to list Kubernetes clusters: %w", err)
	}
	var matches RecordSet
	for _, rec := range records {
		if strings.TrimSpace(fmt.Sprint(rec["name"])) == nameOrVrn {
			matches = append(matches, rec)
		}
	}
	switch len(matches) {
	case 0:
		return "", fmt.Errorf("kubernetes cluster with name %q not found", nameOrVrn)
	case 1:
		vrn := strings.TrimSpace(fmt.Sprint(matches[0]["vrn"]))
		if vrn == "" {
			return "", fmt.Errorf("kubernetes cluster %q has no VRN", nameOrVrn)
		}
		return vrn, nil
	default:
		return "", fmt.Errorf("multiple Kubernetes clusters found with name %q", nameOrVrn)
	}
}

// validateKubernetesClusterNameOrVrnAttrs ensures exactly one of
// kubernetes_cluster_vrn / kubernetes_cluster_name is set. Skips when either
// value is still unknown (plan-time references).
func validateKubernetesClusterNameOrVrnAttrs(label string, attrs map[string]attr.Value) error {
	if attrs == nil {
		return nil
	}
	vrn, hasVrn := attrs[deKubernetesClusterVrnAttr]
	name, hasName := attrs[deKubernetesClusterNameAttr]
	if (hasVrn && vrn.IsUnknown()) || (hasName && name.IsUnknown()) {
		return nil
	}
	vrnSet := hasVrn && !vrn.IsNull() && !vrn.IsUnknown()
	nameSet := hasName && !name.IsNull() && !name.IsUnknown()
	switch {
	case !vrnSet && !nameSet:
		return fmt.Errorf("%s: one of %q or %q must be set", label, deKubernetesClusterVrnAttr, deKubernetesClusterNameAttr)
	case vrnSet && nameSet:
		return fmt.Errorf("%s: only one of %q or %q can be set", label, deKubernetesClusterVrnAttr, deKubernetesClusterNameAttr)
	default:
		return nil
	}
}

// resolveKubernetesClusterNameOrVrnInMap ensures kubernetes_cluster_vrn is set
// on a nest/body map from either vrn or name (exactly one).
func resolveKubernetesClusterNameOrVrnInMap(
	ctx context.Context,
	api VastResourceAPIWithContext,
	listQuery params,
	m map[string]any,
	label string,
) error {
	if m == nil {
		return fmt.Errorf("%s: cluster reference is required", label)
	}
	vrn, _ := m[deKubernetesClusterVrnAttr].(string)
	name, _ := m[deKubernetesClusterNameAttr].(string)
	vrn = strings.TrimSpace(vrn)
	name = strings.TrimSpace(name)
	switch {
	case vrn != "" && name != "":
		return fmt.Errorf("%s: only one of %q or %q can be set", label, deKubernetesClusterVrnAttr, deKubernetesClusterNameAttr)
	case vrn != "":
		return nil
	case name != "":
		resolved, err := resolveKubernetesClusterVRN(ctx, api, listQuery, name)
		if err != nil {
			return fmt.Errorf("%s: %w", label, err)
		}
		m[deKubernetesClusterVrnAttr] = resolved
		return nil
	default:
		return fmt.Errorf("%s: one of %q or %q must be set", label, deKubernetesClusterVrnAttr, deKubernetesClusterNameAttr)
	}
}

// validateKubernetesClusterNameOrVrnInState ensures exactly one of
// kubernetes_cluster_vrn / kubernetes_cluster_name. Skips while either is unknown.
func validateKubernetesClusterNameOrVrnInState(ts is.ResourceState) error {
	if ts == nil {
		return nil
	}
	if v, ok := ts.RawValues()[deKubernetesClusterVrnAttr]; ok && v.IsUnknown() {
		return nil
	}
	if v, ok := ts.RawValues()[deKubernetesClusterNameAttr]; ok && v.IsUnknown() {
		return nil
	}
	return validateOneOf(ts, deKubernetesClusterVrnAttr, deKubernetesClusterNameAttr)
}

// resolveKubernetesClusterNameOrVrnInState resolves a top-level
// kubernetes_cluster_name into kubernetes_cluster_vrn on tfstate (pipeline, secret, …).
func resolveKubernetesClusterNameOrVrnInState(
	ctx context.Context,
	ts is.ResourceState,
	rest *VMSRest,
) error {
	if ts == nil || rest == nil {
		return fmt.Errorf("kubernetes cluster resolution: empty state or client")
	}
	m := map[string]any{}
	if v, ok := ts.RawValues()[deKubernetesClusterVrnAttr]; ok {
		if v.IsUnknown() {
			return nil
		}
		if !v.IsNull() {
			m[deKubernetesClusterVrnAttr] = ts.String(deKubernetesClusterVrnAttr)
		}
	}
	if v, ok := ts.RawValues()[deKubernetesClusterNameAttr]; ok {
		if v.IsUnknown() {
			return nil
		}
		if !v.IsNull() {
			m[deKubernetesClusterNameAttr] = ts.String(deKubernetesClusterNameAttr)
		}
	}
	if len(m) == 0 && ts.RawValues()[deKubernetesClusterVrnAttr] == nil && ts.RawValues()[deKubernetesClusterNameAttr] == nil {
		return nil
	}
	if err := resolveKubernetesClusterNameOrVrnInMap(
		ctx, rest.DataEngine.KubernetesClusters, resolveTenantQuery(ts, rest), m, "resource",
	); err != nil {
		return err
	}
	if vrn, ok := m[deKubernetesClusterVrnAttr].(string); ok && strings.TrimSpace(vrn) != "" {
		ts.Set(deKubernetesClusterVrnAttr, vrn)
	}
	return nil
}

// preserveKubernetesClusterNameInRecord copies kubernetes_cluster_name from tfstate
// into the record (API returns VRN only).
func preserveKubernetesClusterNameInRecord(rec Record, ts is.ResourceState) Record {
	if rec == nil || ts == nil {
		return rec
	}
	if v, ok := ts.RawValues()[deKubernetesClusterNameAttr]; !ok || v.IsNull() || v.IsUnknown() {
		return rec
	}
	name := strings.TrimSpace(ts.String(deKubernetesClusterNameAttr))
	if name == "" {
		return rec
	}
	out := Record{}
	for k, v := range rec {
		out[k] = v
	}
	out[deKubernetesClusterNameAttr] = name
	return out
}

func deContainerRegistryVrnSchemaAttribute() rschema.StringAttribute {
	return rschema.StringAttribute{
		Optional: true,
		Computed: true,
		Description: "Container registry VRN (vast:dataengine:container-registries:...). " +
			"Exactly one of container_registry_vrn or container_registry_name must be set.",
		PlanModifiers: []planmodifiers.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func deContainerRegistryNameSchemaAttribute() rschema.StringAttribute {
	return rschema.StringAttribute{
		Optional: true,
		Description: "Container registry name. Resolved to a VRN before create/update. " +
			"Exactly one of container_registry_vrn or container_registry_name must be set.",
	}
}

// resolveContainerRegistryVRN resolves a registry name to its VRN.
// If nameOrVrn is already a VRN, it is returned unchanged.
// List ignores name filters — match locally after List.
func resolveContainerRegistryVRN(
	ctx context.Context,
	api VastResourceAPIWithContext,
	listQuery params,
	nameOrVrn string,
) (string, error) {
	nameOrVrn = strings.TrimSpace(nameOrVrn)
	if nameOrVrn == "" {
		return "", fmt.Errorf("container registry name or VRN is required")
	}
	if isDataEngineVrn(nameOrVrn) {
		return nameOrVrn, nil
	}
	records, err := api.ListWithContext(ctx, listQuery)
	if err != nil {
		return "", fmt.Errorf("failed to list container registries: %w", err)
	}
	var matches RecordSet
	for _, rec := range records {
		if strings.TrimSpace(fmt.Sprint(rec["name"])) == nameOrVrn {
			matches = append(matches, rec)
		}
	}
	switch len(matches) {
	case 0:
		return "", fmt.Errorf("container registry with name %q not found", nameOrVrn)
	case 1:
		vrn := strings.TrimSpace(fmt.Sprint(matches[0]["vrn"]))
		if vrn == "" {
			return "", fmt.Errorf("container registry %q has no VRN", nameOrVrn)
		}
		return vrn, nil
	default:
		return "", fmt.Errorf("multiple container registries found with name %q", nameOrVrn)
	}
}

// validateContainerRegistryNameOrVrnInState ensures exactly one of
// container_registry_vrn / container_registry_name. Skips while either is unknown.
func validateContainerRegistryNameOrVrnInState(ts is.ResourceState) error {
	if ts == nil {
		return nil
	}
	if v, ok := ts.RawValues()[deContainerRegistryVrnAttr]; ok && v.IsUnknown() {
		return nil
	}
	if v, ok := ts.RawValues()[deContainerRegistryNameAttr]; ok && v.IsUnknown() {
		return nil
	}
	return validateOneOf(ts, deContainerRegistryVrnAttr, deContainerRegistryNameAttr)
}

// resolveContainerRegistryNameOrVrnInState fills container_registry_vrn from name.
func resolveContainerRegistryNameOrVrnInState(
	ctx context.Context,
	ts is.ResourceState,
	rest *VMSRest,
) error {
	if ts == nil || rest == nil {
		return fmt.Errorf("container registry resolution: empty state or client")
	}
	vrnSet := false
	nameSet := false
	var name string
	if v, ok := ts.RawValues()[deContainerRegistryVrnAttr]; ok {
		if v.IsUnknown() {
			return nil
		}
		if !v.IsNull() && strings.TrimSpace(ts.String(deContainerRegistryVrnAttr)) != "" {
			vrnSet = true
		}
	}
	if v, ok := ts.RawValues()[deContainerRegistryNameAttr]; ok {
		if v.IsUnknown() {
			return nil
		}
		if !v.IsNull() {
			name = strings.TrimSpace(ts.String(deContainerRegistryNameAttr))
			nameSet = name != ""
		}
	}
	switch {
	case vrnSet && nameSet:
		return fmt.Errorf("only one of %q or %q can be set", deContainerRegistryVrnAttr, deContainerRegistryNameAttr)
	case vrnSet:
		return nil
	case nameSet:
		resolved, err := resolveContainerRegistryVRN(
			ctx, rest.DataEngine.ContainerRegistries, resolveTenantQuery(ts, rest), name,
		)
		if err != nil {
			return err
		}
		ts.Set(deContainerRegistryVrnAttr, resolved)
		return nil
	default:
		return fmt.Errorf("one of %q or %q must be set", deContainerRegistryVrnAttr, deContainerRegistryNameAttr)
	}
}

// preserveContainerRegistryNameInRecord copies container_registry_name from tfstate
// into the API record (API returns VRN only).
func preserveContainerRegistryNameInRecord(rec Record, ts is.ResourceState) Record {
	if rec == nil || ts == nil {
		return rec
	}
	if v, ok := ts.RawValues()[deContainerRegistryNameAttr]; !ok || v.IsNull() || v.IsUnknown() {
		return rec
	}
	name := strings.TrimSpace(ts.String(deContainerRegistryNameAttr))
	if name == "" {
		return rec
	}
	out := Record{}
	for k, v := range rec {
		out[k] = v
	}
	out[deContainerRegistryNameAttr] = name
	return out
}

// ---------------------------------------------------------------------------
// Telemetries-query data source helpers (de_traces, de_logs, de_metrics, …)
// ---------------------------------------------------------------------------

func deTenantNameDatasourceAttribute() dschema.StringAttribute {
	return dschema.StringAttribute{
		Optional: true,
		Description: "Tenant name for Data Engine API queries. Required for cluster-admin " +
			"auth; optional when the provider already authenticates as that tenant.",
	}
}

// deTelemetryQueryParams collects all user-supplied Required and Optional
// attributes from the Swagger-generated datasource state and merges tenant_name.
func deTelemetryQueryParams(tfState is.ResourceState, rest *VMSRest) params {
	query := params(tfState.GetFilteredValues(
		is.FilterOr,
		nil,
		is.SearchRequired,
		is.SearchOptional,
	))
	return deTransformQueryParams(query, tfState, rest)
}
