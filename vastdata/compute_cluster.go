// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

var computeClusterReplicaSetsSubResource = is.SubResourceHint{
	FieldTrigger: "get_replica_sets",
	SchemaKey:    "",
	SchemaAttributes: map[string]any{
		"replica_sets": rschema.StringAttribute{
			Computed:    true,
			Description: "JSON-encoded list of ReplicaSets in the compute cluster. Populated when get_replica_sets is true.",
		},
	},
}

var computeClusterEventsSubResource = is.SubResourceHint{
	FieldTrigger: "get_events",
	SchemaKey:    "",
	SchemaAttributes: map[string]any{
		"events": rschema.StringAttribute{
			Computed:    true,
			Description: "JSON-encoded list of events in the compute cluster. Populated when get_events is true.",
		},
	},
}

var computeClusterDashboardSubResource = is.SubResourceHint{
	FieldTrigger: "get_dashboard",
	SchemaKey:    "",
	SchemaAttributes: map[string]any{
		"dashboard": rschema.SingleNestedAttribute{
			Computed:    true,
			Description: "Aggregated dashboard statistics across all compute clusters. Populated when get_dashboard is true.",
			Attributes: map[string]rschema.Attribute{
				"namespace_counts": rschema.Int64Attribute{Computed: true, Description: "Total number of namespaces across all clusters."},
				"service_counts":   rschema.Int64Attribute{Computed: true, Description: "Total number of services across all clusters."},
				"tenant_counts":    rschema.Int64Attribute{Computed: true, Description: "Total number of tenants attached to compute clusters."},
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
	computeClusterReplicaSetsSubResource,
	computeClusterEventsSubResource,
	computeClusterDashboardSubResource,
}

var computeClusterAsyncTaskTimeout = 30 * time.Minute

type ComputeCluster struct {
	tfstate *is.TFState
}

func (m *ComputeCluster) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &ComputeCluster{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:        ComputeClusterSchemaRef,
			SubResources:     allComputeClusterSubResources,
			AsyncTaskTimeout: &computeClusterAsyncTaskTimeout,
		},
	)}
}

func (m *ComputeCluster) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ComputeCluster{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:    ComputeClusterSchemaRef,
			SubResources: allComputeClusterSubResources,
		},
	)}
}

func (m *ComputeCluster) TfState() *is.TFState {
	return m.tfstate
}

func (m *ComputeCluster) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.ComputeClusters
}

func (m *ComputeCluster) GetSubResources(ctx context.Context, rest *VMSRest, record Record, _ *version.Version) (Record, error) {
	id := record.RecordID()
	result := Record{}

	// nodes
	if m.tfstate.Bool("get_nodes") {
		records, err := rest.ComputeClusters.ComputeClusterNodesWithContext_GET(ctx, id)
		if err = ignoreStatusCodes(err, http.StatusNotFound); err != nil {
			return nil, err
		}
		result["nodes"] = convertRecordSetToList(records, func(r map[string]any) map[string]any {
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
		})
	}

	// pods
	if m.tfstate.Bool("get_pods") {
		records, err := rest.ComputeClusters.ComputeClusterPodsWithContext_GET(ctx, id, nil)
		if err = ignoreStatusCodes(err, http.StatusNotFound); err != nil {
			return nil, err
		}
		result["pods"] = convertRecordSetToList(records, func(r map[string]any) map[string]any {
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
					condList := make([]map[string]any, 0, len(items))
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
				out["conditions"] = []map[string]any{}
			}
			return out
		})
	}

	// namespaces
	if m.tfstate.Bool("get_namespaces") {
		records, err := rest.ComputeClusters.ComputeClusterNamespacesWithContext_GET(ctx, id)
		if err = ignoreStatusCodes(err, http.StatusNotFound); err != nil {
			return nil, err
		}
		result["namespaces"] = convertRecordSetToList(records, func(r map[string]any) map[string]any {
			return map[string]any{
				"name":       r["name"],
				"status":     r["status"],
				"created_at": r["created_at"],
			}
		})
	}

	// services
	if m.tfstate.Bool("get_services") {
		records, err := rest.ComputeClusters.ComputeClusterServicesWithContext_GET(ctx, id, nil)
		if err = ignoreStatusCodes(err, http.StatusNotFound); err != nil {
			return nil, err
		}
		result["services"] = convertRecordSetToList(records, func(r map[string]any) map[string]any {
			return map[string]any{
				"name":         r["name"],
				"namespace":    r["namespace"],
				"cluster_ip":   r["cluster_ip"],
				"external_ip":  r["external_ip"],
				"service_type": r["service_type"],
				"created_at":   r["created_at"],
			}
		})
	}

	// deployments
	if m.tfstate.Bool("get_deployments") {
		records, err := rest.ComputeClusters.ComputeClusterDeploymentsWithContext_GET(ctx, id, nil)
		if err = ignoreStatusCodes(err, http.StatusNotFound); err != nil {
			return nil, err
		}
		result["deployments"] = convertRecordSetToList(records, func(r map[string]any) map[string]any {
			return map[string]any{
				"name":       r["name"],
				"namespace":  r["namespace"],
				"ready":      r["ready"],
				"available":  r["available"],
				"up_to_date": r["up_to_date"],
				"created_at": r["created_at"],
			}
		})
	}

	// tenants
	if m.tfstate.Bool("get_tenants") {
		records, err := rest.ComputeClusters.ComputeClusterTenantsWithContext_GET(ctx, id)
		if err = ignoreStatusCodes(err, http.StatusNotFound); err != nil {
			return nil, err
		}
		result["cluster_tenants"] = convertRecordSetToList(records, func(r map[string]any) map[string]any {
			return map[string]any{
				"tenant_id":         r["tenant_id"],
				"compute_cluster":   r["compute_cluster"],
				"status":            r["status"],
				"de_compute_guid":   r["de_compute_guid"],
				"de_mtls_cert_guid": r["de_mtls_cert_guid"],
			}
		})
	}

	// replica_sets (ambiguous schema — JSON string)
	if m.tfstate.Bool("get_replica_sets") {
		records, err := rest.ComputeClusters.ComputeClusterReplicaSetsWithContext_GET(ctx, id, nil)
		if err = ignoreStatusCodes(err, http.StatusNotFound); err != nil {
			return nil, err
		}
		raw, _ := json.Marshal(records)
		result["replica_sets"] = string(raw)
	}

	// events (ambiguous schema — JSON string)
	if m.tfstate.Bool("get_events") {
		records, err := rest.ComputeClusters.ComputeClusterEventsWithContext_GET(ctx, id, nil)
		if err = ignoreStatusCodes(err, http.StatusNotFound); err != nil {
			return nil, err
		}
		raw, _ := json.Marshal(records)
		result["events"] = string(raw)
	}

	// dashboard (global — no per-cluster id)
	if m.tfstate.Bool("get_dashboard") {
		dash, err := rest.ComputeClusters.ComputeClusterDashboardWithContext_GET(ctx)
		if err = ignoreStatusCodes(err, http.StatusNotFound); err != nil {
			return nil, err
		}
		if dash != nil {
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
			result["dashboard"] = dashRecord
		}
	}

	return result, nil
}

// convertRecordSetToList converts an untyped RecordSet ([]map[string]any) into
// a []map[string]any with only the fields selected by mapFn.
func convertRecordSetToList(records any, mapFn func(map[string]any) map[string]any) []map[string]any {
	if records == nil {
		return []map[string]any{}
	}
	slice, ok := records.([]map[string]any)
	if !ok {
		// RecordSet is []core.Record which is []map[string]any; try []any fallback
		if anySlice, ok2 := records.([]any); ok2 {
			out := make([]map[string]any, 0, len(anySlice))
			for _, item := range anySlice {
				if m, ok := item.(map[string]any); ok {
					out = append(out, mapFn(m))
				}
			}
			return out
		}
		return []map[string]any{}
	}
	out := make([]map[string]any, 0, len(slice))
	for _, r := range slice {
		out = append(out, mapFn(r))
	}
	return out
}
