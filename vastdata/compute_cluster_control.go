// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	planmodifiers "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

// ComputeClusterControl is a custom action resource that triggers lifecycle
// operations on a Compute Cluster:
//
//	start, stop, rotate_leaf_certificates, rotate_service_key,
//	reconcile_create, rotate_base_certificates, metric_viewer_certificates
//
// All fields are RequiresReplace so any change to action or id will destroy
// the old "execution" and re-trigger the new one.
type ComputeClusterControl struct {
	tfstate *is.TFState
}

func (m *ComputeClusterControl) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &ComputeClusterControl{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			TFStateHintsForCustom: &is.TFStateHintsForCustom{
				Description: "Trigger lifecycle actions on a Compute Cluster (start, stop, certificate rotation, etc.). " +
					"Because Terraform is declarative, each distinct combination of compute_cluster_id + action represents " +
					"a single execution intent; changing either field destroys and re-creates the resource, re-triggering " +
					"the action.",
				SchemaAttributes: map[string]any{
					"compute_cluster_id": rschema.Int64Attribute{
						Required:    true,
						Description: "ID of the Compute Cluster to act on.",
						PlanModifiers: []planmodifiers.Int64{
							int64planmodifier.RequiresReplace(),
						},
					},
					"action": rschema.StringAttribute{
						Required: true,
						Description: "Action to perform. Valid values: start, stop, rotate_leaf_certificates, " +
							"rotate_service_key, reconcile_create, rotate_base_certificates, metric_viewer_certificates.",
						PlanModifiers: []planmodifiers.String{
							stringplanmodifier.RequiresReplace(),
						},
					},

					// ── rotate_base_certificates optional parameters ──

					"keep_root": rschema.BoolAttribute{
						Optional: true,
						Description: "For rotate_base_certificates: if true, keep the existing root certificate " +
							"and only regenerate intermediate and leaf certificates.",
						PlanModifiers: []planmodifiers.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"intermediate_certificate": rschema.StringAttribute{
						Optional:  true,
						Sensitive: true,
						Description: "For rotate_base_certificates: RKE2 intermediate certificate PEM. " +
							"Must be provided together with root_certificate and intermediate_key.",
						PlanModifiers: []planmodifiers.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"intermediate_key": rschema.StringAttribute{
						Optional:  true,
						Sensitive: true,
						Description: "For rotate_base_certificates: RKE2 intermediate private key PEM. " +
							"Must be provided together with root_certificate and intermediate_certificate.",
						PlanModifiers: []planmodifiers.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"root_certificate": rschema.StringAttribute{
						Optional:  true,
						Sensitive: true,
						Description: "For rotate_base_certificates: RKE2 root certificate PEM. " +
							"Must be provided together with intermediate_certificate and intermediate_key.",
						PlanModifiers: []planmodifiers.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
		},
	)}
}

func (m *ComputeClusterControl) TfState() *is.TFState {
	return m.tfstate
}

func (m *ComputeClusterControl) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.ComputeClusters
}

func (m *ComputeClusterControl) ValidateResourceConfig(context.Context) error {
	return ValidateFieldIsOneOf(
		m.tfstate, "action",
		"start", "stop",
		"rotate_leaf_certificates", "rotate_service_key",
		"reconcile_create", "rotate_base_certificates",
		"metric_viewer_certificates",
	)
}

func (m *ComputeClusterControl) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.performAction(ctx, rest)
}

func (m *ComputeClusterControl) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	return m.performAction(ctx, rest)
}

func (m *ComputeClusterControl) DeleteResource(_ context.Context, _ *VMSRest) error {
	// Lifecycle actions are not reversible via Terraform delete.
	return nil
}

// performAction executes the configured action on the compute cluster.
func (m *ComputeClusterControl) performAction(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	id := m.tfstate.Int64("compute_cluster_id")
	action := m.tfstate.String("action")

	if id == 0 {
		return nil, fmt.Errorf("compute_cluster_id must not be 0")
	}
	if action == "" {
		return nil, fmt.Errorf("action must not be empty")
	}

	const noWait = time.Duration(0)

	switch action {
	case "start":
		_, err := rest.ComputeClusters.ComputeClusterStartWithContext_POST(ctx, id, nil, noWait)
		return nil, err

	case "stop":
		_, err := rest.ComputeClusters.ComputeClusterStopWithContext_POST(ctx, id, nil, noWait)
		return nil, err

	case "rotate_leaf_certificates":
		_, err := rest.ComputeClusters.ComputeClusterRotateLeafCertificatesWithContext_POST(ctx, id, nil, noWait)
		return nil, err

	case "rotate_service_key":
		_, err := rest.ComputeClusters.ComputeClusterRotateServiceKeyWithContext_POST(ctx, id, nil, noWait)
		return nil, err

	case "reconcile_create":
		_, err := rest.ComputeClusters.ComputeClusterReconcileCreateWithContext_POST(ctx, id, nil, noWait)
		return nil, err

	case "rotate_base_certificates":
		queryParams := params{}
		if m.tfstate.IsKnownAndNotNull("keep_root") {
			queryParams["keep_root"] = m.tfstate.Bool("keep_root")
		}
		body := params{}
		if m.tfstate.IsKnownAndNotNull("intermediate_certificate") {
			body["intermediate_certificate"] = m.tfstate.String("intermediate_certificate")
		}
		if m.tfstate.IsKnownAndNotNull("intermediate_key") {
			body["intermediate_key"] = m.tfstate.String("intermediate_key")
		}
		if m.tfstate.IsKnownAndNotNull("root_certificate") {
			body["root_certificate"] = m.tfstate.String("root_certificate")
		}
		_, err := rest.ComputeClusters.ComputeClusterRotateBaseCertificatesWithContext_POST(ctx, id, queryParams, body, noWait)
		return nil, err

	case "metric_viewer_certificates":
		_, err := rest.ComputeClusters.ComputeClusterMetricViewerCertificatesWithContext_POST(ctx, id, nil)
		return nil, err

	default:
		return nil, fmt.Errorf(
			"unknown action %q. Valid: start, stop, rotate_leaf_certificates, rotate_service_key, "+
				"reconcile_create, rotate_base_certificates, metric_viewer_certificates", action,
		)
	}
}
