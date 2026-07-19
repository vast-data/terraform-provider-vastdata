// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var ComputeClusterSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"computeclusters",
	http.MethodGet,
	"computeclusters",
)

var computeClusterNodesSubResource = is.SubResourceHint{
	FieldTrigger: "get_nodes",
	SchemaKey:    "",
	SchemaAttributes: map[string]any{
		"nodes": rschema.ListNestedAttribute{
			Computed:    true,
			Description: "List of Kubernetes nodes in the compute cluster. Populated when get_nodes is true.",
			NestedObject: rschema.NestedAttributeObject{
				Attributes: map[string]rschema.Attribute{
					"name":            rschema.StringAttribute{Computed: true, Description: "Node name."},
					"status":          rschema.StringAttribute{Computed: true, Description: "Node status (Ready, Unknown)."},
					"architecture":    rschema.StringAttribute{Computed: true, Description: "Node architecture (amd64, arm64, etc.)."},
					"capacity_cpu":    rschema.StringAttribute{Computed: true, Description: "CPU capacity of the node."},
					"capacity_memory": rschema.StringAttribute{Computed: true, Description: "Memory capacity of the node."},
					"os_image":        rschema.StringAttribute{Computed: true, Description: "Operating system image."},
					"version":         rschema.StringAttribute{Computed: true, Description: "Kubernetes version running on the node."},
					"message":         rschema.StringAttribute{Computed: true, Description: "Unhealthy condition messages from the node."},
					"created_at":      rschema.StringAttribute{Computed: true, Description: "Node creation timestamp."},
				},
			},
		},
	},
}

var computeClusterPodsSubResource = is.SubResourceHint{
	FieldTrigger: "get_pods",
	SchemaKey:    "",
	SchemaAttributes: map[string]any{
		"pods": rschema.ListNestedAttribute{
			Computed:    true,
			Description: "List of Kubernetes pods in the compute cluster. Populated when get_pods is true.",
			NestedObject: rschema.NestedAttributeObject{
				Attributes: map[string]rschema.Attribute{
					"name":       rschema.StringAttribute{Computed: true, Description: "Pod name."},
					"namespace":  rschema.StringAttribute{Computed: true, Description: "Pod namespace."},
					"ip":         rschema.StringAttribute{Computed: true, Description: "Pod IP address."},
					"node":       rschema.StringAttribute{Computed: true, Description: "Node name where pod is scheduled."},
					"status":     rschema.StringAttribute{Computed: true, Description: "Pod status (Running, Pending, Failed, etc.)."},
					"created_at": rschema.StringAttribute{Computed: true, Description: "Pod creation timestamp."},
					"conditions": rschema.ListNestedAttribute{
						Computed:    true,
						Description: "Pod status conditions reported by Kubernetes.",
						NestedObject: rschema.NestedAttributeObject{
							Attributes: map[string]rschema.Attribute{
								"type":                 rschema.StringAttribute{Computed: true, Description: "Condition type (e.g. Ready, Initialized, ContainersReady, PodScheduled)."},
								"status":               rschema.StringAttribute{Computed: true, Description: "Condition status (True, False, Unknown)."},
								"reason":               rschema.StringAttribute{Computed: true, Description: "Machine-readable reason for the condition's last transition."},
								"message":              rschema.StringAttribute{Computed: true, Description: "Human-readable message about the condition's last transition."},
								"last_transition_time": rschema.StringAttribute{Computed: true, Description: "Timestamp of the condition's last transition."},
							},
						},
					},
				},
			},
		},
	},
}

var computeClusterNamespacesSubResource = is.SubResourceHint{
	FieldTrigger: "get_namespaces",
	SchemaKey:    "",
	SchemaAttributes: map[string]any{
		"namespaces": rschema.ListNestedAttribute{
			Computed:    true,
			Description: "List of Kubernetes namespaces in the compute cluster. Populated when get_namespaces is true.",
			NestedObject: rschema.NestedAttributeObject{
				Attributes: map[string]rschema.Attribute{
					"name":       rschema.StringAttribute{Computed: true, Description: "Namespace name."},
					"status":     rschema.StringAttribute{Computed: true, Description: "Namespace status (Active, Terminating, etc.)."},
					"created_at": rschema.StringAttribute{Computed: true, Description: "Namespace creation timestamp."},
				},
			},
		},
	},
}

var computeClusterServicesSubResource = is.SubResourceHint{
	FieldTrigger: "get_services",
	SchemaKey:    "",
	SchemaAttributes: map[string]any{
		"services": rschema.ListNestedAttribute{
			Computed:    true,
			Description: "List of Kubernetes services in the compute cluster. Populated when get_services is true.",
			NestedObject: rschema.NestedAttributeObject{
				Attributes: map[string]rschema.Attribute{
					"name":         rschema.StringAttribute{Computed: true, Description: "Service name."},
					"namespace":    rschema.StringAttribute{Computed: true, Description: "Namespace where the service is located."},
					"cluster_ip":   rschema.StringAttribute{Computed: true, Description: "Cluster IP address."},
					"external_ip":  rschema.StringAttribute{Computed: true, Description: "External IP address."},
					"service_type": rschema.StringAttribute{Computed: true, Description: "Service type (ClusterIP, NodePort, LoadBalancer, etc.)."},
					"created_at":   rschema.StringAttribute{Computed: true, Description: "Service creation timestamp."},
				},
			},
		},
	},
}

var computeClusterDeploymentsSubResource = is.SubResourceHint{
	FieldTrigger: "get_deployments",
	SchemaKey:    "",
	SchemaAttributes: map[string]any{
		"deployments": rschema.ListNestedAttribute{
			Computed:    true,
			Description: "List of Kubernetes deployments in the compute cluster. Populated when get_deployments is true.",
			NestedObject: rschema.NestedAttributeObject{
				Attributes: map[string]rschema.Attribute{
					"name":       rschema.StringAttribute{Computed: true, Description: "Deployment name."},
					"namespace":  rschema.StringAttribute{Computed: true, Description: "Namespace where the deployment is located."},
					"ready":      rschema.StringAttribute{Computed: true, Description: "Ready replicas in \"X/Y\" format (like kubectl)."},
					"available":  rschema.Int64Attribute{Computed: true, Description: "Number of available replicas."},
					"up_to_date": rschema.Int64Attribute{Computed: true, Description: "Number of up-to-date replicas."},
					"created_at": rschema.StringAttribute{Computed: true, Description: "Deployment creation timestamp."},
				},
			},
		},
	},
}

var computeClusterTenantsSubResource = is.SubResourceHint{
	FieldTrigger: "get_tenants",
	SchemaKey:    "",
	SchemaAttributes: map[string]any{
		"cluster_tenants": rschema.ListNestedAttribute{
			Computed:    true,
			Description: "List of tenants associated with the compute cluster. Populated when get_tenants is true.",
			NestedObject: rschema.NestedAttributeObject{
				Attributes: map[string]rschema.Attribute{
					"tenant_id":         rschema.Int64Attribute{Computed: true, Description: "Tenant ID."},
					"compute_cluster":   rschema.Int64Attribute{Computed: true, Description: "Compute Cluster ID."},
					"status":            rschema.StringAttribute{Computed: true, Description: "Status of the tenant association."},
					"de_compute_guid":   rschema.StringAttribute{Computed: true, Description: "Data Engine compute GUID."},
					"de_mtls_cert_guid": rschema.StringAttribute{Computed: true, Description: "Data Engine mTLS certificate GUID."},
				},
			},
		},
	},
}

var computeClusterDashboardSubResource = is.SubResourceHint{
	FieldTrigger: "get_dashboard",
	SchemaKey:    "",
	SchemaAttributes: map[string]any{
		"dashboard": rschema.SingleNestedAttribute{
			Computed:    true,
			Description: "Dashboard statistics for the compute cluster. Populated when get_dashboard is true.",
			Attributes: map[string]rschema.Attribute{
				"namespace_counts": rschema.Int64Attribute{Computed: true, Description: "Number of namespaces in the cluster."},
				"service_counts":   rschema.Int64Attribute{Computed: true, Description: "Number of services in the cluster."},
				"tenant_counts":    rschema.Int64Attribute{Computed: true, Description: "Number of tenants attached to the compute cluster."},
				"cnodes":           rschema.StringAttribute{Computed: true, Description: "JSON map of cnode status → count."},
				"compute_clusters": rschema.StringAttribute{Computed: true, Description: "JSON map of compute cluster status → count."},
				"deployments":      rschema.StringAttribute{Computed: true, Description: "JSON map of deployment status → count."},
				"pods":             rschema.StringAttribute{Computed: true, Description: "JSON map of pod status → count."},
			},
		},
	},
}

var allComputeClusterSubResources = []is.SubResourceHint{
	computeClusterNodesSubResource,
	computeClusterPodsSubResource,
	computeClusterNamespacesSubResource,
	computeClusterServicesSubResource,
	computeClusterDeploymentsSubResource,
	computeClusterTenantsSubResource,
	computeClusterDashboardSubResource,
}

var computeClusterAsyncTaskTimeout = 40 * time.Minute

// 503 SERVICE_UNAVAILABLE is transient when VMS cannot reach the internal Kubernetes API.
var computeClusterSubResourceRetryOn = &is.RetryExpression{
	StatusCodes: []int{http.StatusServiceUnavailable},
	Times:       5,
}

type ComputeCluster struct {
	tfstate *is.TFState
}

func (m *ComputeCluster) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &ComputeCluster{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:               ComputeClusterSchemaRef,
			SubResources:            allComputeClusterSubResources,
			AsyncTaskTimeout:        &computeClusterAsyncTaskTimeout,
			RetryOn:                 computeClusterSubResourceRetryOn,
			NotRequiredSchemaFields: []string{"cnodes"},
			PreserveOrderFields:     []string{"static_ip_ranges"},
		},
	)}
}

func (m *ComputeCluster) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ComputeCluster{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:           ComputeClusterSchemaRef,
			SubResources:        allComputeClusterSubResources,
			RetryOn:             computeClusterSubResourceRetryOn,
			PreserveOrderFields: []string{"static_ip_ranges"},
		},
	)}
}

func (m *ComputeCluster) TfState() *is.TFState {
	return m.tfstate
}

func (m *ComputeCluster) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.ComputeClusters
}

func (m *ComputeCluster) NormalizeRecordForCreateAdopt(record Record) Record {
	if record == nil {
		return record
	}
	if cnodes, ok := record["cnodes"]; ok {
		record["cnodes"] = normalizeComputeClusterCnodes(cnodes)
	}
	return record
}

func (m *ComputeCluster) ResolveRecordAfterAsyncTask(ctx context.Context, rest *VMSRest, record Record) (Record, error) {
	if !isComputeClusterAsyncTaskRecord(record) {
		return record, nil
	}
	clusterID := computeClusterIDFromAsyncRecord(record)
	if clusterID == 0 && m.tfstate.IsKnownAndNotNull("id") {
		clusterID = m.tfstate.Int64("id")
	}
	if clusterID == 0 {
		return record, fmt.Errorf("compute cluster async task completed but cluster id is unknown")
	}
	return rest.ComputeClusters.GetByIdWithContext(ctx, clusterID)
}

func normalizeComputeClusterCnodes(raw any) any {
	items, ok := raw.([]any)
	if !ok {
		return raw
	}
	out := make([]any, 0, len(items))
	for _, item := range items {
		cnode, ok := item.(map[string]any)
		if !ok {
			continue
		}
		normalized := map[string]any{"id": cnode["id"]}
		if preset, ok := cnode["resource_preset"]; ok {
			normalized["resource_preset"] = preset
		}
		out = append(out, normalized)
	}
	return out
}

func computeClusterIDFromAsyncRecord(record Record) int64 {
	info, ok := record["info"].(map[string]any)
	if !ok {
		if asyncTask, ok := record["async_task"].(map[string]any); ok {
			info, ok = asyncTask["info"].(map[string]any)
		}
		if !ok {
			return 0
		}
	}
	if id := int64FromAny(info["compute_cluster_id"]); id != 0 {
		return id
	}
	if kwargs, ok := info["kwargs"].(map[string]any); ok {
		return int64FromAny(kwargs["compute_cluster_id"])
	}
	return 0
}

func isComputeClusterAsyncTaskRecord(record Record) bool {
	if record == nil || record.Empty() {
		return false
	}
	if rt, ok := record["@resourceType"]; ok && fmt.Sprintf("%v", rt) == "VTask" {
		return true
	}
	if _, ok := record["async_task"]; ok {
		return true
	}
	return computeClusterIDFromAsyncRecord(record) != 0
}

func int64FromAny(v any) int64 {
	switch n := v.(type) {
	case int:
		return int64(n)
	case int32:
		return int64(n)
	case int64:
		return n
	case float64:
		return int64(n)
	case json.Number:
		i, err := n.Int64()
		if err == nil {
			return i
		}
	}
	return 0
}

func (m *ComputeCluster) GetSubResources(ctx context.Context, rest *VMSRest, record Record, _ *version.Version) (Record, error) {
	if isComputeClusterK8sSubResourceFetchUnavailable(clusterStateFromRecord(record)) {
		tflog.Debug(ctx, fmt.Sprintf(
			"GetSubResources[compute_cluster]: cluster state %q — reusing sub-resources from tfstate",
			record["state"],
		))
		return computeClusterSubResourcesFromTFState(m.tfstate), nil
	}

	id := record.RecordID()
	result := Record{}

	if m.tfstate.Bool("get_nodes") {
		nodes, err := fetchComputeClusterNodes(ctx, rest, id)
		if err != nil {
			return nil, err
		}
		result["nodes"] = nodes
	}

	if m.tfstate.Bool("get_pods") {
		pods, err := fetchComputeClusterPods(ctx, rest, id, nil)
		if err != nil {
			return nil, err
		}
		result["pods"] = pods
	}

	if m.tfstate.Bool("get_namespaces") {
		namespaces, err := fetchComputeClusterNamespaces(ctx, rest, id)
		if err != nil {
			return nil, err
		}
		result["namespaces"] = namespaces
	}

	if m.tfstate.Bool("get_services") {
		services, err := fetchComputeClusterServices(ctx, rest, id, nil)
		if err != nil {
			return nil, err
		}
		result["services"] = services
	}

	if m.tfstate.Bool("get_deployments") {
		deployments, err := fetchComputeClusterDeployments(ctx, rest, id, nil)
		if err != nil {
			return nil, err
		}
		result["deployments"] = deployments
	}

	if m.tfstate.Bool("get_tenants") {
		tenants, err := fetchComputeClusterTenants(ctx, rest, id)
		if err != nil {
			return nil, err
		}
		result["cluster_tenants"] = tenants
	}

	if m.tfstate.Bool("get_dashboard") {
		dashboard, err := fetchComputeClusterDashboard(ctx, rest)
		if err != nil {
			return nil, err
		}
		if dashboard != nil {
			result["dashboard"] = dashboard
		}
	}

	return result, nil
}

func clusterStateFromRecord(record Record) string {
	if record == nil {
		return ""
	}
	state, _ := record["state"].(string)
	return state
}

func isComputeClusterK8sSubResourceFetchUnavailable(state string) bool {
	switch strings.ToUpper(strings.TrimSpace(state)) {
	case "STOPPED", "STOPPING":
		return true
	default:
		return false
	}
}

func computeClusterSubResourcesFromTFState(tfstate *is.TFState) Record {
	if tfstate == nil || tfstate.Hints == nil {
		return Record{}
	}
	all := tfstate.GetAllValues()
	result := Record{}
	for _, sr := range tfstate.Hints.SubResources {
		if sr.SchemaKey != "" {
			if v, ok := all[sr.SchemaKey]; ok {
				result[sr.SchemaKey] = v
			}
			continue
		}
		for k := range sr.SchemaAttributes {
			if v, ok := all[k]; ok {
				result[k] = v
			}
		}
	}
	return result
}

func withComputeClusterSubResourceRetry[T any](
	ctx context.Context,
	expr *is.RetryExpression,
	name string,
	fn func() (T, error),
) (T, error) {
	if expr == nil {
		return fn()
	}
	return retryOnExpression(ctx, expr, "GetSubResources", "compute_cluster/"+name, fn)
}

// convertRecordSetToList converts an API list response into []any for FillFromRecord.
func convertRecordSetToList(records any, mapFn func(map[string]any) map[string]any) []any {
	if records == nil {
		return []any{}
	}

	switch slice := records.(type) {
	case RecordSet:
		return mapRecordSliceToAny(slice, mapFn)
	case []map[string]any:
		return mapRecordSliceToAny(slice, mapFn)
	case []any:
		out := make([]any, 0, len(slice))
		for _, item := range slice {
			if m, ok := item.(map[string]any); ok {
				out = append(out, mapFn(m))
			}
		}
		return out
	case Record:
		if raw, ok := slice["results"].([]any); ok {
			return convertRecordSetToList(raw, mapFn)
		}
	case map[string]any:
		if raw, ok := slice["results"].([]any); ok {
			return convertRecordSetToList(raw, mapFn)
		}
	}

	return []any{}
}

func mapRecordSliceToAny[T ~map[string]any](slice []T, mapFn func(map[string]any) map[string]any) []any {
	out := make([]any, 0, len(slice))
	for _, r := range slice {
		out = append(out, mapFn(r))
	}
	return out
}

func computeClusterLookupSchemaAttributes() map[string]any {
	return map[string]any{
		"compute_cluster_id": dschema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Compute cluster ID. Provide this or compute_cluster_name.",
		},
		"compute_cluster_name": dschema.StringAttribute{
			Optional:    true,
			Description: "Compute cluster name. Used to look up the cluster ID when compute_cluster_id is not set.",
		},
	}
}

func resolveComputeClusterContext(ctx context.Context, rest *VMSRest, tfstate *is.TFState) (int64, Record, error) {
	if tfstate.IsKnownAndNotNull("compute_cluster_id") {
		id := tfstate.Int64("compute_cluster_id")
		if id != 0 {
			record := Record{"compute_cluster_id": id}
			if tfstate.IsKnownAndNotNull("compute_cluster_name") {
				record["compute_cluster_name"] = tfstate.String("compute_cluster_name")
			}
			return id, record, nil
		}
	}

	name := tfstate.String("compute_cluster_name")
	if name == "" {
		return 0, nil, fmt.Errorf("either compute_cluster_id or compute_cluster_name must be set")
	}

	cluster, err := rest.ComputeClusters.GetWithContext(ctx, params{"name": name})
	if err != nil {
		return 0, nil, fmt.Errorf("lookup compute cluster by name %q: %w", name, err)
	}

	id := cluster.RecordID()
	return id, Record{
		"compute_cluster_id":   id,
		"compute_cluster_name": name,
	}, nil
}

func mergeComputeClusterRecords(base Record, extra Record) Record {
	out := Record{}
	for k, v := range base {
		out[k] = v
	}
	for k, v := range extra {
		out[k] = v
	}
	return out
}

func mapComputeClusterNode(r map[string]any) map[string]any {
	return map[string]any{
		"name":            r["name"],
		"status":          r["status"],
		"architecture":    r["architecture"],
		"capacity_cpu":    r["capacity_cpu"],
		"capacity_memory": r["capacity_memory"],
		"os_image":        r["os_image"],
		"version":         r["version"],
		"message":         r["message"],
		"created_at":      r["created_at"],
	}
}

func mapComputeClusterPod(r map[string]any) map[string]any {
	out := map[string]any{
		"name":       r["name"],
		"namespace":  r["namespace"],
		"ip":         r["ip"],
		"node":       r["node"],
		"status":     r["status"],
		"created_at": r["created_at"],
	}
	if conds, ok := r["conditions"]; ok && conds != nil {
		if items, ok := conds.([]any); ok {
			condList := make([]any, 0, len(items))
			for _, c := range items {
				if cm, ok := c.(map[string]any); ok {
					condList = append(condList, map[string]any{
						"type":                 cm["type"],
						"status":               cm["status"],
						"reason":               cm["reason"],
						"message":              cm["message"],
						"last_transition_time": cm["last_transition_time"],
					})
				}
			}
			out["conditions"] = condList
		}
	}
	if out["conditions"] == nil {
		out["conditions"] = []any{}
	}
	return out
}

func mapComputeClusterNamespace(r map[string]any) map[string]any {
	return map[string]any{
		"name":       r["name"],
		"status":     r["status"],
		"created_at": r["created_at"],
	}
}

func mapComputeClusterService(r map[string]any) map[string]any {
	return map[string]any{
		"name":         r["name"],
		"namespace":    r["namespace"],
		"cluster_ip":   r["cluster_ip"],
		"external_ip":  r["external_ip"],
		"service_type": r["service_type"],
		"created_at":   r["created_at"],
	}
}

func mapComputeClusterDeployment(r map[string]any) map[string]any {
	return map[string]any{
		"name":       r["name"],
		"namespace":  r["namespace"],
		"ready":      r["ready"],
		"available":  r["available"],
		"up_to_date": r["up_to_date"],
		"created_at": r["created_at"],
	}
}

func mapComputeClusterTenant(r map[string]any) map[string]any {
	return map[string]any{
		"tenant_id":         r["tenant_id"],
		"compute_cluster":   r["compute_cluster"],
		"status":            r["status"],
		"de_compute_guid":   r["de_compute_guid"],
		"de_mtls_cert_guid": r["de_mtls_cert_guid"],
	}
}

func fetchComputeClusterNodes(ctx context.Context, rest *VMSRest, id int64) ([]any, error) {
	records, err := withComputeClusterSubResourceRetry(ctx, computeClusterSubResourceRetryOn, "nodes", func() (RecordSet, error) {
		return rest.ComputeClusters.ComputeClusterNodesWithContext_GET(ctx, id)
	})
	if err != nil {
		return nil, err
	}
	return convertRecordSetToList(records, mapComputeClusterNode), nil
}

func fetchComputeClusterPods(ctx context.Context, rest *VMSRest, id int64, query params) ([]any, error) {
	records, err := withComputeClusterSubResourceRetry(ctx, computeClusterSubResourceRetryOn, "pods", func() (RecordSet, error) {
		return rest.ComputeClusters.ComputeClusterPodsWithContext_GET(ctx, id, query)
	})
	if err != nil {
		return nil, err
	}
	return convertRecordSetToList(records, mapComputeClusterPod), nil
}

func fetchComputeClusterNamespaces(ctx context.Context, rest *VMSRest, id int64) ([]any, error) {
	records, err := withComputeClusterSubResourceRetry(ctx, computeClusterSubResourceRetryOn, "namespaces", func() (RecordSet, error) {
		return rest.ComputeClusters.ComputeClusterNamespacesWithContext_GET(ctx, id)
	})
	if err != nil {
		return nil, err
	}
	return convertRecordSetToList(records, mapComputeClusterNamespace), nil
}

func fetchComputeClusterServices(ctx context.Context, rest *VMSRest, id int64, query params) ([]any, error) {
	records, err := withComputeClusterSubResourceRetry(ctx, computeClusterSubResourceRetryOn, "services", func() (RecordSet, error) {
		return rest.ComputeClusters.ComputeClusterServicesWithContext_GET(ctx, id, query)
	})
	if err != nil {
		return nil, err
	}
	return convertRecordSetToList(records, mapComputeClusterService), nil
}

func fetchComputeClusterDeployments(ctx context.Context, rest *VMSRest, id int64, query params) ([]any, error) {
	records, err := withComputeClusterSubResourceRetry(ctx, computeClusterSubResourceRetryOn, "deployments", func() (RecordSet, error) {
		return rest.ComputeClusters.ComputeClusterDeploymentsWithContext_GET(ctx, id, query)
	})
	if err != nil {
		return nil, err
	}
	return convertRecordSetToList(records, mapComputeClusterDeployment), nil
}

func fetchComputeClusterTenants(ctx context.Context, rest *VMSRest, id int64) ([]any, error) {
	records, err := withComputeClusterSubResourceRetry(ctx, computeClusterSubResourceRetryOn, "tenants", func() (RecordSet, error) {
		return rest.ComputeClusters.ComputeClusterTenantsWithContext_GET(ctx, id)
	})
	if err != nil {
		return nil, err
	}
	return convertRecordSetToList(records, mapComputeClusterTenant), nil
}

func fetchComputeClusterReplicaSets(ctx context.Context, rest *VMSRest, id int64, query params) (string, error) {
	records, err := withComputeClusterSubResourceRetry(ctx, computeClusterSubResourceRetryOn, "replica_sets", func() (RecordSet, error) {
		return rest.ComputeClusters.ComputeClusterReplicaSetsWithContext_GET(ctx, id, query)
	})
	if err != nil {
		return "", err
	}
	raw, _ := json.Marshal(records)
	return string(raw), nil
}

func fetchComputeClusterEvents(ctx context.Context, rest *VMSRest, id int64, query params) (string, error) {
	records, err := withComputeClusterSubResourceRetry(ctx, computeClusterSubResourceRetryOn, "events", func() (RecordSet, error) {
		return rest.ComputeClusters.ComputeClusterEventsWithContext_GET(ctx, id, query)
	})
	if err != nil {
		return "", err
	}
	raw, _ := json.Marshal(records)
	return string(raw), nil
}

func fetchComputeClusterDashboard(ctx context.Context, rest *VMSRest) (map[string]any, error) {
	dash, err := withComputeClusterSubResourceRetry(ctx, computeClusterSubResourceRetryOn, "dashboard", func() (Record, error) {
		return rest.ComputeClusters.ComputeClusterDashboardWithContext_GET(ctx)
	})
	if err != nil {
		return nil, err
	}
	if dash == nil {
		return nil, nil
	}
	dashRecord := map[string]any{
		"namespace_counts": dash["namespace_counts"],
		"service_counts":   dash["service_counts"],
		"tenant_counts":    dash["tenant_counts"],
	}
	for _, mapField := range []string{"cnodes", "compute_clusters", "deployments", "pods"} {
		if v, ok := dash[mapField]; ok {
			raw, _ := json.Marshal(v)
			dashRecord[mapField] = string(raw)
		} else {
			dashRecord[mapField] = types.StringNull().ValueString()
		}
	}
	return dashRecord, nil
}

func computeClusterPodConditionsSchema() dschema.ListNestedAttribute {
	return dschema.ListNestedAttribute{
		Computed:    true,
		Description: "Pod status conditions reported by Kubernetes.",
		NestedObject: dschema.NestedAttributeObject{
			Attributes: map[string]dschema.Attribute{
				"type":                 dschema.StringAttribute{Computed: true, Description: "Condition type (e.g. Ready, Initialized, ContainersReady, PodScheduled)."},
				"status":               dschema.StringAttribute{Computed: true, Description: "Condition status (True, False, Unknown)."},
				"reason":               dschema.StringAttribute{Computed: true, Description: "Machine-readable reason for the condition's last transition."},
				"message":              dschema.StringAttribute{Computed: true, Description: "Human-readable message about the condition's last transition."},
				"last_transition_time": dschema.StringAttribute{Computed: true, Description: "Timestamp of the condition's last transition."},
			},
		},
	}
}

func computeClusterNodesListSchema() dschema.ListNestedAttribute {
	return dschema.ListNestedAttribute{
		Computed:    true,
		Description: "Kubernetes nodes in the compute cluster.",
		NestedObject: dschema.NestedAttributeObject{
			Attributes: map[string]dschema.Attribute{
				"name":            dschema.StringAttribute{Computed: true, Description: "Node name."},
				"status":          dschema.StringAttribute{Computed: true, Description: "Node status (Ready, Unknown)."},
				"architecture":    dschema.StringAttribute{Computed: true, Description: "Node architecture (amd64, arm64, etc.)."},
				"capacity_cpu":    dschema.StringAttribute{Computed: true, Description: "CPU capacity of the node."},
				"capacity_memory": dschema.StringAttribute{Computed: true, Description: "Memory capacity of the node."},
				"os_image":        dschema.StringAttribute{Computed: true, Description: "Operating system image."},
				"version":         dschema.StringAttribute{Computed: true, Description: "Kubernetes version running on the node."},
				"message":         dschema.StringAttribute{Computed: true, Description: "Unhealthy condition messages from the node."},
				"created_at":      dschema.StringAttribute{Computed: true, Description: "Node creation timestamp."},
			},
		},
	}
}

func computeClusterPodsListSchema() dschema.ListNestedAttribute {
	return dschema.ListNestedAttribute{
		Computed:    true,
		Description: "Kubernetes pods in the compute cluster.",
		NestedObject: dschema.NestedAttributeObject{
			Attributes: map[string]dschema.Attribute{
				"name":       dschema.StringAttribute{Computed: true, Description: "Pod name."},
				"namespace":  dschema.StringAttribute{Computed: true, Description: "Pod namespace."},
				"ip":         dschema.StringAttribute{Computed: true, Description: "Pod IP address."},
				"node":       dschema.StringAttribute{Computed: true, Description: "Node name where pod is scheduled."},
				"status":     dschema.StringAttribute{Computed: true, Description: "Pod status (Running, Pending, Failed, etc.)."},
				"created_at": dschema.StringAttribute{Computed: true, Description: "Pod creation timestamp."},
				"conditions": computeClusterPodConditionsSchema(),
			},
		},
	}
}

func computeClusterNamespacesListSchema() dschema.ListNestedAttribute {
	return dschema.ListNestedAttribute{
		Computed:    true,
		Description: "Kubernetes namespaces in the compute cluster.",
		NestedObject: dschema.NestedAttributeObject{
			Attributes: map[string]dschema.Attribute{
				"name":       dschema.StringAttribute{Computed: true, Description: "Namespace name."},
				"status":     dschema.StringAttribute{Computed: true, Description: "Namespace status (Active, Terminating, etc.)."},
				"created_at": dschema.StringAttribute{Computed: true, Description: "Namespace creation timestamp."},
			},
		},
	}
}

func computeClusterServicesListSchema() dschema.ListNestedAttribute {
	return dschema.ListNestedAttribute{
		Computed:    true,
		Description: "Kubernetes services in the compute cluster.",
		NestedObject: dschema.NestedAttributeObject{
			Attributes: map[string]dschema.Attribute{
				"name":         dschema.StringAttribute{Computed: true, Description: "Service name."},
				"namespace":    dschema.StringAttribute{Computed: true, Description: "Namespace where the service is located."},
				"cluster_ip":   dschema.StringAttribute{Computed: true, Description: "Cluster IP address."},
				"external_ip":  dschema.StringAttribute{Computed: true, Description: "External IP address."},
				"service_type": dschema.StringAttribute{Computed: true, Description: "Service type (ClusterIP, NodePort, LoadBalancer, etc.)."},
				"created_at":   dschema.StringAttribute{Computed: true, Description: "Service creation timestamp."},
			},
		},
	}
}

func computeClusterDeploymentsListSchema() dschema.ListNestedAttribute {
	return dschema.ListNestedAttribute{
		Computed:    true,
		Description: "Kubernetes deployments in the compute cluster.",
		NestedObject: dschema.NestedAttributeObject{
			Attributes: map[string]dschema.Attribute{
				"name":       dschema.StringAttribute{Computed: true, Description: "Deployment name."},
				"namespace":  dschema.StringAttribute{Computed: true, Description: "Namespace where the deployment is located."},
				"ready":      dschema.StringAttribute{Computed: true, Description: "Ready replicas in \"X/Y\" format (like kubectl)."},
				"available":  dschema.Int64Attribute{Computed: true, Description: "Number of available replicas."},
				"up_to_date": dschema.Int64Attribute{Computed: true, Description: "Number of up-to-date replicas."},
				"created_at": dschema.StringAttribute{Computed: true, Description: "Deployment creation timestamp."},
			},
		},
	}
}

func computeClusterTenantsListSchema() dschema.ListNestedAttribute {
	return dschema.ListNestedAttribute{
		Computed:    true,
		Description: "Tenants associated with the compute cluster.",
		NestedObject: dschema.NestedAttributeObject{
			Attributes: map[string]dschema.Attribute{
				"tenant_id":         dschema.Int64Attribute{Computed: true, Description: "Tenant ID."},
				"compute_cluster":   dschema.Int64Attribute{Computed: true, Description: "Compute Cluster ID."},
				"status":            dschema.StringAttribute{Computed: true, Description: "Status of the tenant association."},
				"de_compute_guid":   dschema.StringAttribute{Computed: true, Description: "Data Engine compute GUID."},
				"de_mtls_cert_guid": dschema.StringAttribute{Computed: true, Description: "Data Engine mTLS certificate GUID."},
			},
		},
	}
}

func computeClusterDashboardObjectSchema() dschema.SingleNestedAttribute {
	return dschema.SingleNestedAttribute{
		Computed:    true,
		Description: "Dashboard statistics for the compute cluster.",
		Attributes: map[string]dschema.Attribute{
			"namespace_counts": dschema.Int64Attribute{Computed: true, Description: "Number of namespaces in the cluster."},
			"service_counts":   dschema.Int64Attribute{Computed: true, Description: "Number of services in the cluster."},
			"tenant_counts":    dschema.Int64Attribute{Computed: true, Description: "Number of tenants attached to the compute cluster."},
			"cnodes":           dschema.StringAttribute{Computed: true, Description: "JSON map of cnode status → count."},
			"compute_clusters": dschema.StringAttribute{Computed: true, Description: "JSON map of compute cluster status → count."},
			"deployments":      dschema.StringAttribute{Computed: true, Description: "JSON map of deployment status → count."},
			"pods":             dschema.StringAttribute{Computed: true, Description: "JSON map of pod status → count."},
		},
	}
}

func newComputeClusterSubresourceTFState(
	raw map[string]attr.Value,
	schema any,
	description string,
	attrs map[string]any,
) *is.TFState {
	merged := computeClusterLookupSchemaAttributes()
	for k, v := range attrs {
		merged[k] = v
	}
	return is.NewTFStateMust(raw, schema, &is.TFStateHints{
		TFStateHintsForCustom: &is.TFStateHintsForCustom{
			Description:      description,
			SchemaAttributes: merged,
		},
	})
}

type computeClusterSubresourceDatasource struct {
	tfstate *is.TFState
}

func (m *computeClusterSubresourceDatasource) TfState() *is.TFState { return m.tfstate }

func (m *computeClusterSubresourceDatasource) API(_ *VMSRest) VastResourceAPIWithContext { return nil }

func (m *computeClusterSubresourceDatasource) lookupRecord(ctx context.Context, rest *VMSRest) (Record, error) {
	_, base, err := resolveComputeClusterContext(ctx, rest, m.tfstate)
	if err != nil {
		return nil, fmt.Errorf("resolve compute cluster: %w", err)
	}
	return base, nil
}

type ComputeClusterNodes struct {
	computeClusterSubresourceDatasource
}

func (m *ComputeClusterNodes) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ComputeClusterNodes{computeClusterSubresourceDatasource{
		tfstate: newComputeClusterSubresourceTFState(raw, schema,
			"Kubernetes nodes in a compute cluster.",
			map[string]any{"nodes": computeClusterNodesListSchema()}),
	}}
}

func (m *ComputeClusterNodes) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	id, base, err := resolveComputeClusterContext(ctx, rest, m.tfstate)
	if err != nil {
		return nil, err
	}
	nodes, err := fetchComputeClusterNodes(ctx, rest, id)
	if err != nil {
		return nil, err
	}
	return mergeComputeClusterRecords(base, Record{"nodes": nodes}), nil
}

type ComputeClusterPods struct {
	computeClusterSubresourceDatasource
}

func (m *ComputeClusterPods) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ComputeClusterPods{computeClusterSubresourceDatasource{
		tfstate: newComputeClusterSubresourceTFState(raw, schema,
			"Kubernetes pods in a compute cluster.",
			map[string]any{
				"namespace":  dschema.StringAttribute{Optional: true, Description: "Filter pods by namespace."},
				"deployment": dschema.StringAttribute{Optional: true, Description: "Filter pods by deployment name. Requires namespace."},
				"pods":       computeClusterPodsListSchema(),
			}),
	}}
}

func (m *ComputeClusterPods) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	id, base, err := resolveComputeClusterContext(ctx, rest, m.tfstate)
	if err != nil {
		return nil, err
	}
	query := params{}
	if m.tfstate.IsKnownAndNotNull("namespace") {
		query["namespace"] = m.tfstate.String("namespace")
	}
	if m.tfstate.IsKnownAndNotNull("deployment") {
		query["deployment"] = m.tfstate.String("deployment")
	}
	pods, err := fetchComputeClusterPods(ctx, rest, id, query)
	if err != nil {
		return nil, err
	}
	return mergeComputeClusterRecords(base, Record{"pods": pods}), nil
}

type ComputeClusterNamespaces struct {
	computeClusterSubresourceDatasource
}

func (m *ComputeClusterNamespaces) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ComputeClusterNamespaces{computeClusterSubresourceDatasource{
		tfstate: newComputeClusterSubresourceTFState(raw, schema,
			"Kubernetes namespaces in a compute cluster.",
			map[string]any{"namespaces": computeClusterNamespacesListSchema()}),
	}}
}

func (m *ComputeClusterNamespaces) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	id, base, err := resolveComputeClusterContext(ctx, rest, m.tfstate)
	if err != nil {
		return nil, err
	}
	namespaces, err := fetchComputeClusterNamespaces(ctx, rest, id)
	if err != nil {
		return nil, err
	}
	return mergeComputeClusterRecords(base, Record{"namespaces": namespaces}), nil
}

type ComputeClusterServices struct {
	computeClusterSubresourceDatasource
}

func (m *ComputeClusterServices) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ComputeClusterServices{computeClusterSubresourceDatasource{
		tfstate: newComputeClusterSubresourceTFState(raw, schema,
			"Kubernetes services in a compute cluster.",
			map[string]any{"services": computeClusterServicesListSchema()}),
	}}
}

func (m *ComputeClusterServices) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	id, base, err := resolveComputeClusterContext(ctx, rest, m.tfstate)
	if err != nil {
		return nil, err
	}
	services, err := fetchComputeClusterServices(ctx, rest, id, nil)
	if err != nil {
		return nil, err
	}
	return mergeComputeClusterRecords(base, Record{"services": services}), nil
}

type ComputeClusterDeployments struct {
	computeClusterSubresourceDatasource
}

func (m *ComputeClusterDeployments) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ComputeClusterDeployments{computeClusterSubresourceDatasource{
		tfstate: newComputeClusterSubresourceTFState(raw, schema,
			"Kubernetes deployments in a compute cluster.",
			map[string]any{"deployments": computeClusterDeploymentsListSchema()}),
	}}
}

func (m *ComputeClusterDeployments) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	id, base, err := resolveComputeClusterContext(ctx, rest, m.tfstate)
	if err != nil {
		return nil, err
	}
	deployments, err := fetchComputeClusterDeployments(ctx, rest, id, nil)
	if err != nil {
		return nil, err
	}
	return mergeComputeClusterRecords(base, Record{"deployments": deployments}), nil
}

type ComputeClusterTenants struct {
	computeClusterSubresourceDatasource
}

func (m *ComputeClusterTenants) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ComputeClusterTenants{computeClusterSubresourceDatasource{
		tfstate: newComputeClusterSubresourceTFState(raw, schema,
			"Tenants associated with a compute cluster.",
			map[string]any{"cluster_tenants": computeClusterTenantsListSchema()}),
	}}
}

func (m *ComputeClusterTenants) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	id, base, err := resolveComputeClusterContext(ctx, rest, m.tfstate)
	if err != nil {
		return nil, err
	}
	tenants, err := fetchComputeClusterTenants(ctx, rest, id)
	if err != nil {
		return nil, err
	}
	return mergeComputeClusterRecords(base, Record{"cluster_tenants": tenants}), nil
}

type ComputeClusterDashboard struct {
	computeClusterSubresourceDatasource
}

func (m *ComputeClusterDashboard) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ComputeClusterDashboard{computeClusterSubresourceDatasource{
		tfstate: newComputeClusterSubresourceTFState(raw, schema,
			"Dashboard statistics for a compute cluster.",
			map[string]any{"dashboard": computeClusterDashboardObjectSchema()}),
	}}
}

func (m *ComputeClusterDashboard) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	base, err := m.lookupRecord(ctx, rest)
	if err != nil {
		return nil, err
	}
	dashboard, err := fetchComputeClusterDashboard(ctx, rest)
	if err != nil {
		return nil, err
	}
	return mergeComputeClusterRecords(base, Record{"dashboard": dashboard}), nil
}

type ComputeClusterReplicaSets struct {
	computeClusterSubresourceDatasource
}

func (m *ComputeClusterReplicaSets) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ComputeClusterReplicaSets{computeClusterSubresourceDatasource{
		tfstate: newComputeClusterSubresourceTFState(raw, schema,
			"ReplicaSets for a deployment in a compute cluster.",
			map[string]any{
				"resource_name":      dschema.StringAttribute{Required: true, Description: "Deployment name to list ReplicaSets for."},
				"resource_namespace": dschema.StringAttribute{Required: true, Description: "Namespace of the deployment."},
				"replica_sets":       dschema.StringAttribute{Computed: true, Description: "JSON-encoded list of ReplicaSets for the deployment."},
			}),
	}}
}

func (m *ComputeClusterReplicaSets) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	id, base, err := resolveComputeClusterContext(ctx, rest, m.tfstate)
	if err != nil {
		return nil, err
	}
	query := params{
		"resource_name":      m.tfstate.String("resource_name"),
		"resource_namespace": m.tfstate.String("resource_namespace"),
	}
	replicaSets, err := fetchComputeClusterReplicaSets(ctx, rest, id, query)
	if err != nil {
		return nil, err
	}
	return mergeComputeClusterRecords(base, Record{"replica_sets": replicaSets}), nil
}

type ComputeClusterEvents struct {
	computeClusterSubresourceDatasource
}

func (m *ComputeClusterEvents) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ComputeClusterEvents{computeClusterSubresourceDatasource{
		tfstate: newComputeClusterSubresourceTFState(raw, schema,
			"Events for a Kubernetes resource in a compute cluster.",
			map[string]any{
				"resource_name":      dschema.StringAttribute{Required: true, Description: "Name of the Kubernetes resource."},
				"resource_namespace": dschema.StringAttribute{Required: true, Description: "Namespace of the Kubernetes resource."},
				"resource_kind":      dschema.StringAttribute{Required: true, Description: "Kind of the resource (Pod, Deployment, Service, ReplicaSet)."},
				"events":             dschema.StringAttribute{Computed: true, Description: "JSON-encoded list of events for the resource."},
			}),
	}}
}

func (m *ComputeClusterEvents) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	id, base, err := resolveComputeClusterContext(ctx, rest, m.tfstate)
	if err != nil {
		return nil, err
	}
	query := params{
		"resource_name":      m.tfstate.String("resource_name"),
		"resource_namespace": m.tfstate.String("resource_namespace"),
		"resource_kind":      m.tfstate.String("resource_kind"),
	}
	events, err := fetchComputeClusterEvents(ctx, rest, id, query)
	if err != nil {
		return nil, err
	}
	return mergeComputeClusterRecords(base, Record{"events": events}), nil
}
