// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/vast-data/go-vast-client/core"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/logging"
)

var DePipelineSchemaRef = is.NewDataEngineSchemaReference(
	http.MethodPost,
	"pipelines",
	http.MethodGet,
	"pipelines",
)

var dePipelineDeployTimeout = 10 * time.Minute

// Bound for re-POST /deploy while status=Error with a known-transient Knative race.
var dePipelineTransientRetryTimeout = 2 * time.Minute

// Sleep between waitReady / waitGone / stack-delete retry cycles.
var dePipelineWaitSleep = 5 * time.Second

// Substrings in pipeline status.reason that indicate Knative/operator races that
// normally clear once brokers/services become Addressable. Retry deploy for a
// bounded window instead of failing Terraform on the first Error.
var transientPipelineDeployReasonSubstrings = []string{
	"address not set",
	"failed to resolve spec.delivery.deadlettersink",
	"failed to resolve subscriber",
	"probestatus",
	"status: notready",
}

// VMS sometimes reports status=Error with these reasons while the stack is
// already applying. Keep polling until Ready; do not fail or re-POST /deploy.
var inProgressPipelineDeployReasonSubstrings = []string{
	"deployed",
}

type DePipeline struct {
	tfstate is.ResourceState
}

func (m *DePipeline) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &DePipeline{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:   DePipelineSchemaRef,
			PathIDField: "guid",
			ExcludedSchemaFields: []string{
				"cursor", "limit", "order_by",
				"owner_name", "created_after", "created_before", "updated_after", "updated_before",
				// OpenAPI uses camelCase property names that Terraform rejects.
				"binaryData",
				"conditional_trigger_labels",
			},
			// OpenAPI PipelineCreateInfo only marks name required; keep TF UX of requiring the graph fields.
			// kubernetes_cluster_vrn is required by OpenAPI; TF allows name OR vrn.
			RequiredSchemaFields: []string{
				"name", "manifest", "namespace",
			},
			NotRequiredSchemaFields: []string{deKubernetesClusterVrnAttr},
			// Match prior custom schema: ordered lists (not sets) for graph nodes/edges and tags.
			PreserveOrderFields: []string{
				"tags",
				"triggers",
				"function_deployments",
				"links",
				"source",
				"destination",
			},
			AsyncTaskTimeout: &dePipelineDeployTimeout,
			CreateOnlyFields: []string{deTenantNameAttr},
			PreserveUserValueFields: []string{
				"manifest",
				deTenantNameAttr,
				"deploy",
				deKubernetesClusterVrnAttr,
				deKubernetesClusterNameAttr,
			},
			CommonModifiersMapping: map[string]string{
				deTenantNameAttr: ModifierForceNew,
			},
			// Failed create+deploy leaves the pipeline in state (tainted). Next apply replaces it;
			// VMS stack teardown is async and returns 409 "being deleted" until finished.
			RetryOn: &is.RetryPolicy{
				Create: &is.RetryExpression{
					StatusCodes:  []int{http.StatusConflict},
					BodyContains: []string{"being deleted"},
					Times:        30,
					SleepSeconds: 5,
				},
			},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr:            deTenantNameSchemaAttribute(),
				deKubernetesClusterVrnAttr:  deKubernetesClusterVrnSchemaAttribute(),
				deKubernetesClusterNameAttr: deKubernetesClusterNameSchemaAttribute(),
				"deploy": rschema.BoolAttribute{
					Optional: true,
					Computed: true,
					Default:  booldefault.StaticBool(true),
					Description: "When true (default), POST /pipeline/{guid}/deploy after create/update " +
						"and wait until status is Ready.",
				},
			},
		},
	)}
}

func (m *DePipeline) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &DePipeline{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:   DePipelineSchemaRef,
			PathIDField: "guid",
			ExcludedSchemaFields: []string{
				"cursor", "limit", "order_by",
				"owner_name", "created_after", "created_before", "updated_after", "updated_before",
				"binaryData",
				"conditional_trigger_labels",
			},
			SearchableFields: []string{"guid", "vrn", "name"},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameSchemaAttribute(),
			},
		},
	)}
}

func (m *DePipeline) PrepareCreateResource(ctx context.Context, rest *VMSRest) error {
	return resolveKubernetesClusterNameOrVrnInState(ctx, m.tfstate, rest)
}

func (m *DePipeline) PrepareUpdateResource(ctx context.Context, plan PrepareUpdateResource, rest *VMSRest) error {
	return resolveKubernetesClusterNameOrVrnInState(ctx, plan.(*DePipeline).tfstate, rest)
}

func (m *DePipeline) TransformQueryParams(query params, rest *VMSRest) params {
	return deTransformQueryParams(query, m.tfstate, rest)
}

func (m *DePipeline) TransformResponseRecord(response Record) Record {
	response = ensureTenantNameInRecord(response, m.tfstate, nil)
	return preserveKubernetesClusterNameInRecord(response, m.tfstate)
}

func (m *DePipeline) TfState() is.ResourceState {
	return m.tfstate
}

func (m *DePipeline) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.DataEngine.Pipelines
}

func (m *DePipeline) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	api := m.API(rest)
	body := m.TransformRequestBody(m.tfstate.GetCreateParams())
	query := m.TransformQueryParams(nil, rest)

	out, err := m.createOrAdopt(ctx, rest, api, query, body)
	if err != nil {
		return nil, err
	}

	guid := nestedString(out, "guid")
	if guid == "" {
		return nil, fmt.Errorf("de_pipeline: create response missing guid")
	}
	if m.shouldDeploy(m.tfstate) {
		if err := m.deployAndWait(ctx, rest, api, guid); err != nil {
			if refreshed, getErr := m.getByGUID(ctx, rest, api, guid); getErr == nil {
				out = refreshed
			}
			return m.stateRecord(out), err
		}
		out, err = m.getByGUID(ctx, rest, api, guid)
		if err != nil {
			return nil, err
		}
	}
	return m.stateRecord(out), nil
}

// createOrAdopt POSTs a new pipeline, or adopts a same-name orphan after 409.
// When the conflict/orphan is mid-teardown ("stack is being deleted"), wait until
// it is gone and retry create — replace after a tainted create races with async delete.
func (m *DePipeline) createOrAdopt(
	ctx context.Context,
	rest *VMSRest,
	api VastResourceAPIWithContext,
	query core.Params,
	body core.Params,
) (Record, error) {
	var out Record
	var lastErr error
	err := waitWithTimeout(ctx, waitWithTimeoutConfig{
		Timeout: dePipelineDeployTimeout,
		Sleep:   dePipelineWaitSleep,
		OnTimeout: func() error {
			if lastErr != nil {
				return lastErr
			}
			return fmt.Errorf("de_pipeline: timed out waiting to create/adopt after stack conflict")
		},
	}, func(ctx context.Context) (bool, error) {
		created, err := api.CreateWithContext(ctx, query, body)
		if err == nil {
			out = created
			return true, nil
		}
		lastErr = err
		if !expectStatusCodes(err, http.StatusConflict) {
			return false, err
		}

		existing, lookupErr := m.findByName(ctx, rest, api)
		if lookupErr != nil {
			return false, fmt.Errorf("de_pipeline create conflict, adopt failed: %w (original: %v)", lookupErr, err)
		}
		if existing == nil {
			// Name conflict without a readable object — usually mid-delete; wait and retry.
			logging.Warn(ctx, "DePipeline Create: 409 with no adoptable pipeline — waiting for teardown")
			return false, nil
		}

		guid := nestedString(existing, "guid")
		if guid == "" {
			return false, fmt.Errorf("de_pipeline: adopted pipeline missing guid")
		}
		logging.Warn(ctx, "DePipeline Create: adopting existing pipeline after 409 conflict")
		if len(body) > 0 {
			if _, updErr := api.UpdateWithContext(ctx, guid, query, body); updErr != nil {
				if isStackBeingDeleted(updErr) {
					logging.Warn(ctx, fmt.Sprintf(
						"DePipeline Create: stack being deleted for %s — waiting then recreating", guid,
					))
					if waitErr := m.waitGone(ctx, rest, api, guid); waitErr != nil {
						return false, waitErr
					}
					return false, nil
				}
				out = existing
				return false, updErr
			}
		}
		out = existing
		return true, nil
	})
	return out, err
}

func (m *DePipeline) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	guid, err := requirePathGUID(m.tfstate, "de_pipeline")
	if err != nil {
		return nil, err
	}
	api := m.API(rest)
	rec, err := m.getByGUID(ctx, rest, api, guid)
	if err != nil {
		return nil, err
	}
	return m.stateRecord(rec), nil
}

func (m *DePipeline) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	planMgr := plan.(*DePipeline)
	planTS := planMgr.tfstate
	guid, err := pathGUIDOrPlan(m.tfstate, planTS, "de_pipeline")
	if err != nil {
		return nil, err
	}
	api := m.API(rest)
	body := planMgr.TransformRequestBody(planTS.GetCreateParams())
	query := planMgr.TransformQueryParams(nil, rest)

	if len(body) > 0 {
		if err := m.updateWithStackDeleteRetry(ctx, rest, api, guid, query, body); err != nil {
			return nil, err
		}
	}

	m.tfstate = planTS
	if m.shouldDeploy(planTS) {
		if err := m.deployAndWait(ctx, rest, api, guid); err != nil {
			return nil, err
		}
	}
	out, err := m.getByGUID(ctx, rest, api, guid)
	if err != nil {
		return nil, err
	}
	return m.stateRecord(out), nil
}

func (m *DePipeline) DeleteResource(ctx context.Context, rest *VMSRest) error {
	guid, err := requirePathGUID(m.tfstate, "de_pipeline")
	if err != nil {
		return err
	}
	api := m.API(rest)
	query := core.Params(m.TransformQueryParams(nil, rest))
	_, err = api.DeleteByIdWithContext(ctx, guid, query, nil)
	if err = ignoreResourceGone(err); err != nil {
		return err
	}
	// Stack teardown is async; wait so a follow-up create does not race with "being deleted".
	return m.waitGone(ctx, rest, api, guid)
}

// updateWithStackDeleteRetry PATCHes the pipeline, retrying while VMS reports the stack
// is still being deleted (transient 409 after a prior replace/delete).
func (m *DePipeline) updateWithStackDeleteRetry(
	ctx context.Context,
	rest *VMSRest,
	api VastResourceAPIWithContext,
	guid string,
	query core.Params,
	body core.Params,
) error {
	var lastErr error
	return waitWithTimeout(ctx, waitWithTimeoutConfig{
		Timeout: dePipelineDeployTimeout,
		Sleep:   dePipelineWaitSleep,
		OnTimeout: func() error {
			if lastErr != nil {
				return lastErr
			}
			return fmt.Errorf("de_pipeline: timed out retrying update while stack is being deleted")
		},
	}, func(ctx context.Context) (bool, error) {
		_, err := api.UpdateWithContext(ctx, guid, query, body)
		if err == nil {
			return true, nil
		}
		lastErr = err
		if !isStackBeingDeleted(err) {
			return false, err
		}
		logging.Warn(ctx, fmt.Sprintf(
			"DePipeline Update: stack being deleted for %s — retrying", guid,
		))
		return false, nil
	})
}

func (m *DePipeline) TransformRequestBody(body params) params {
	body = deTransformRequestBody(body)
	if body != nil {
		body.Without(
			"deploy",
			deKubernetesClusterNameAttr,
			"guid", "id", "status", "reason", "vrn",
			"created_at", "updated_at", "owner", "tenant_guid",
			"last_revision_number", "stack_guid", "stack_revision_number",
		)
	}
	return body
}

func (m *DePipeline) shouldDeploy(ts is.ResourceState) bool {
	if ts == nil || !ts.HasAttribute("deploy") {
		return true
	}
	if !ts.IsKnownAndNotNull("deploy") {
		return true
	}
	return ts.Bool("deploy")
}

func (m *DePipeline) deployAndWait(ctx context.Context, rest *VMSRest, api VastResourceAPIWithContext, guid string) error {
	if err := m.deploy(ctx, rest, guid); err != nil {
		return err
	}
	return m.waitReady(ctx, rest, api, guid)
}

func (m *DePipeline) deploy(ctx context.Context, rest *VMSRest, guid string) error {
	query := m.TransformQueryParams(nil, rest)
	logging.Info(ctx, fmt.Sprintf("DePipeline: deploying %s", guid))
	if _, err := rest.DataEngine.Pipelines.DeployWithContext(ctx, guid, query); err != nil {
		return fmt.Errorf("de_pipeline deploy: %w", err)
	}
	return nil
}

func (m *DePipeline) waitReady(ctx context.Context, rest *VMSRest, api VastResourceAPIWithContext, guid string) error {
	var lastStatus, lastReason string
	var transientDeadline time.Time
	return waitWithTimeout(ctx, waitWithTimeoutConfig{
		Timeout: dePipelineDeployTimeout,
		Sleep:   dePipelineWaitSleep,
		OnTimeout: func() error {
			return fmt.Errorf(
				"timed out waiting for de_pipeline Ready (last status=%q reason=%q)",
				lastStatus, lastReason,
			)
		},
	}, func(ctx context.Context) (bool, error) {
		rec, err := m.getByGUID(ctx, rest, api, guid)
		if err != nil {
			return false, err
		}
		status := nestedString(rec, "status")
		reason := nestedString(rec, "reason")
		lastStatus, lastReason = status, reason
		logging.Debug(ctx, fmt.Sprintf("DePipeline wait: status=%q reason=%q", status, reason))
		switch strings.ToLower(strings.TrimSpace(status)) {
		case "ready":
			return true, nil
		case "failed", "error":
			if isInProgressPipelineDeployReason(reason) {
				transientDeadline = time.Time{}
				return false, nil
			}
			if !isTransientPipelineDeployReason(reason) {
				return false, fmt.Errorf("de_pipeline deploy failed: status=%q reason=%q", status, reason)
			}
			if transientDeadline.IsZero() {
				transientDeadline = time.Now().Add(dePipelineTransientRetryTimeout)
			}
			if time.Now().After(transientDeadline) {
				return false, fmt.Errorf(
					"de_pipeline deploy failed after %v of transient retries: status=%q reason=%q",
					dePipelineTransientRetryTimeout, status, reason,
				)
			}
			if err := m.deploy(ctx, rest, guid); err != nil {
				logging.Warn(ctx, fmt.Sprintf(
					"DePipeline wait: redeploy after transient error failed (will retry): %v", err,
				))
			}
			return false, nil
		default:
			// Left Error for an in-progress status; allow a fresh transient window if it returns.
			transientDeadline = time.Time{}
			return false, nil
		}
	})
}

// isTransientPipelineDeployReason reports Knative/operator races that clear once
// Addressable resources (brokers, services) finish provisioning.
func isTransientPipelineDeployReason(reason string) bool {
	return reasonContainsAny(reason, transientPipelineDeployReasonSubstrings)
}

func isInProgressPipelineDeployReason(reason string) bool {
	return reasonContainsAny(reason, inProgressPipelineDeployReasonSubstrings)
}

func reasonContainsAny(reason string, subs []string) bool {
	r := strings.ToLower(strings.TrimSpace(reason))
	if r == "" {
		return false
	}
	for _, sub := range subs {
		if strings.Contains(r, sub) {
			return true
		}
	}
	return false
}

// waitGone blocks until GET /pipelines/{guid} returns 404 (stack teardown finished).
func (m *DePipeline) waitGone(ctx context.Context, rest *VMSRest, api VastResourceAPIWithContext, guid string) error {
	return waitWithTimeout(ctx, waitWithTimeoutConfig{
		Timeout:        dePipelineDeployTimeout,
		Sleep:          dePipelineWaitSleep,
		TimeoutMessage: fmt.Sprintf("timed out waiting for de_pipeline %s to be deleted", guid),
	}, func(ctx context.Context) (bool, error) {
		_, err := m.getByGUID(ctx, rest, api, guid)
		if err != nil {
			if isNotFoundErr(err) || expectStatusCodes(err, http.StatusNotFound) {
				return true, nil
			}
			return false, err
		}
		logging.Debug(ctx, fmt.Sprintf("DePipeline waitGone: %s still present", guid))
		return false, nil
	})
}

func (m *DePipeline) getByGUID(ctx context.Context, rest *VMSRest, api VastResourceAPIWithContext, guid string) (Record, error) {
	query := m.TransformQueryParams(nil, rest)
	return api.GetByIdWithContext(ctx, guid, query)
}

func (m *DePipeline) findByName(ctx context.Context, rest *VMSRest, api VastResourceAPIWithContext) (Record, error) {
	if m.tfstate == nil || !m.tfstate.IsKnownAndNotNull("name") {
		return nil, fmt.Errorf("de_pipeline: name is required to adopt")
	}
	name := strings.TrimSpace(m.tfstate.String("name"))
	query := m.TransformQueryParams(nil, rest)
	query["name"] = name
	rec, err := api.GetWithContext(ctx, query)
	if err != nil {
		if isNotFoundErr(err) || expectStatusCodes(err, http.StatusNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if nestedString(rec, "guid") == "" && nestedString(rec, "name") == "" {
		return nil, nil
	}
	return rec, nil
}

func (m *DePipeline) stateRecord(apiRec Record) Record {
	out := Record{}
	for k, v := range apiRec {
		out[k] = v
	}
	for _, key := range []string{"manifest", "deploy", deTenantNameAttr, deKubernetesClusterNameAttr} {
		if m.tfstate == nil || !m.tfstate.HasAttribute(key) || !m.tfstate.IsKnownAndNotNull(key) {
			continue
		}
		out[key] = convertAttrToRaw(m.tfstate, key)
	}
	return out
}

// isStackBeingDeleted reports a 409 conflict while a pipeline stack is still tearing down.
func isStackBeingDeleted(err error) bool {
	if err == nil || !expectStatusCodes(err, http.StatusConflict) {
		return false
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "being deleted") {
		return true
	}
	if isApiError(err) {
		return strings.Contains(strings.ToLower(err.(*ApiError).Body), "being deleted")
	}
	return false
}
