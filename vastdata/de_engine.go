// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vast-data/go-vast-client/core"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

// DeEngine manages tenant Data Engine via POST/GET/DELETE /serverless/data-engine.
// Not a normal id-based collection — default CRUD cannot be used as-is.
type DeEngine struct {
	tfstate is.ResourceState
}

var DeEngineSchemaRef = is.NewDataEngineSchemaReference(
	http.MethodPost,
	"data-engine",
	http.MethodGet,
	"data-engine",
)

var deEngineAsyncTimeout = 30 * time.Minute

func (m *DeEngine) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &DeEngine{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:        DeEngineSchemaRef,
			Importable:       &notImportable,
			AsyncTaskTimeout: &deEngineAsyncTimeout,
			SensitiveFields:  []string{"kafka_ca_certificate"},
			CreateOnlyFields: []string{
				"default_topic_name",
				"dead_letter_topic_name",
				"kafka_broker",
				"kafka_ca_certificate",
				"vip_pools",
				deTenantNameAttr,
			},
			// Keep user create params across refresh — API returns expanded kafka_broker
			// (url, is_accessible_to_data_engine) that must not force perpetual replace.
			PreserveUserValueFields: []string{
				"default_topic_name",
				"dead_letter_topic_name",
				"kafka_broker",
				"kafka_ca_certificate",
				"vip_pools",
				deTenantNameAttr,
			},
			CommonModifiersMapping: map[string]string{
				"default_topic_name":     ModifierForceNew,
				"dead_letter_topic_name": ModifierForceNew,
				"kafka_ca_certificate":   ModifierForceNew,
				deTenantNameAttr:         ModifierForceNew,
			},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameSchemaAttribute(),
			},
		},
	)}
}

func (m *DeEngine) TransformQueryParams(query params, rest *VMSRest) params {
	return deTransformQueryParams(query, m.tfstate, rest)
}

func (m *DeEngine) TransformRequestBody(body params) params {
	return deTransformRequestBody(body)
}

func (m *DeEngine) TransformResponseRecord(response Record) Record {
	return ensureTenantNameInRecord(response, m.tfstate, nil)
}

func (m *DeEngine) TfState() is.ResourceState { return m.tfstate }

func (m *DeEngine) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.DataEngine.DataEngine
}

func (m *DeEngine) LookupForCreate(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.ReadResource(ctx, rest)
}

func (m *DeEngine) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	api := m.API(rest)

	// Adopt if already enabled (custom CreateResource skips the default LookupForCreate path).
	if rec, ok := m.tryAdopt(ctx, rest, api); ok {
		tflog.Info(ctx, "DeEngine Create: already enabled — adopting existing Data Engine")
		return rec, nil
	}

	body := m.requestBody()
	query := m.TransformQueryParams(nil, rest)

	// DELETE can report setup=not_started before VMS fully clears DE-enabled; retry 409 briefly.
	var out DisplayableRecord
	var lastErr error
	var attempt int
	err := waitWithTimeout(ctx, waitWithTimeoutConfig{
		Cycles: 10,
		Sleep:  5 * time.Second,
		OnTimeout: func() error {
			if lastErr != nil {
				return lastErr
			}
			return fmt.Errorf("data-engine: create conflict retries exhausted")
		},
	}, func(ctx context.Context) (bool, error) {
		attempt++
		if attempt > 1 {
			if rec, ok := m.tryAdopt(ctx, rest, api); ok {
				tflog.Info(ctx, "DeEngine Create: adopted after conflict/retry")
				out = rec
				return true, nil
			}
		}
		_, err := api.CreateWithContext(ctx, query, body)
		if err == nil {
			status, err := m.waitStatus(ctx, rest, api, "completed")
			if err != nil {
				return false, err
			}
			out = m.stateRecord(status)
			return true, nil
		}
		lastErr = err
		if !expectStatusCodes(err, http.StatusConflict) {
			return false, err
		}
		tflog.Warn(ctx, fmt.Sprintf("DeEngine Create: 409 on attempt %d — waiting to adopt/retry", attempt))
		return false, nil
	})
	if err != nil {
		if rec, ok := m.tryAdopt(ctx, rest, api); ok {
			return rec, nil
		}
		return nil, err
	}
	return out, nil
}

func (m *DeEngine) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	api := m.API(rest)
	status, err := m.getStatus(ctx, rest, api)
	if err != nil {
		return nil, err
	}
	setupStatus := nestedString(status, "setup", "status")
	if !isDataEnginePresent(setupStatus) {
		return nil, notFoundErr("data-engine", "setup.status="+setupStatus)
	}
	return m.stateRecord(status), nil
}

func (m *DeEngine) DeleteResource(ctx context.Context, rest *VMSRest) error {
	api := m.API(rest)
	// Singleton collection DELETE /data-engine (no guid); client override omits body so VMS defaults apply.
	_, err := api.DeleteWithContext(ctx, nil, m.TransformQueryParams(nil, rest), nil)
	if ignoreStatusCodes(err, http.StatusNotFound) != nil {
		return err
	}
	_, err = m.waitStatus(ctx, rest, api, "idle")
	return err
}

func (m *DeEngine) tryAdopt(ctx context.Context, rest *VMSRest, api VastResourceAPIWithContext) (DisplayableRecord, bool) {
	status, err := m.getStatus(ctx, rest, api)
	if err != nil {
		return nil, false
	}
	if !isDataEnginePresent(nestedString(status, "setup", "status")) {
		return nil, false
	}
	return m.stateRecord(status), true
}

func (m *DeEngine) requestBody() core.Params {
	body := core.Params{
		"default_topic_name":     m.tfstate.String("default_topic_name"),
		"dead_letter_topic_name": m.tfstate.String("dead_letter_topic_name"),
		"kafka_broker":           m.kafkaBrokerForRequest(),
	}
	if m.tfstate.HasAttribute("kafka_ca_certificate") && m.tfstate.IsKnownAndNotNull("kafka_ca_certificate") {
		if v := m.tfstate.String("kafka_ca_certificate"); v != "" {
			body["kafka_ca_certificate"] = v
		}
	}
	if m.tfstate.HasAttribute("vip_pools") && m.tfstate.IsKnownAndNotNull("vip_pools") {
		if pools := m.tfstate.ToSlice("vip_pools"); len(pools) > 0 {
			body["vip_pools"] = pools
		}
	}
	return body
}

// kafkaBrokerForRequest sends only create-relevant broker fields (name/type[+optional flag]).
func (m *DeEngine) kafkaBrokerForRequest() map[string]any {
	raw := m.tfstate.ToMap("kafka_broker")
	out := map[string]any{}
	if v, ok := raw["name"]; ok && v != nil {
		out["name"] = v
	}
	if v, ok := raw["type"]; ok && v != nil {
		out["type"] = v
	}
	if v, ok := raw["is_accessible_to_data_engine"]; ok && v != nil {
		out["is_accessible_to_data_engine"] = v
	}
	return out
}

func (m *DeEngine) getStatus(ctx context.Context, rest *VMSRest, api VastResourceAPIWithContext) (Record, error) {
	return api.GetWithContext(ctx, m.TransformQueryParams(nil, rest))
}

func (m *DeEngine) waitStatus(ctx context.Context, rest *VMSRest, api VastResourceAPIWithContext, want string) (Record, error) {
	timeout := deEngineAsyncTimeout

	return core.WaitAPICondition(
		ctx,
		api,
		m.TransformQueryParams(nil, rest),
		&core.WaitAPIConditionConfig{
			Timeout:  timeout,
			Interval: 5 * time.Second,
		},
		func(status Record) (bool, error) {
			cur := nestedString(status, "setup", "status")
			msg := nestedString(status, "setup", "message")
			del := nestedString(status, "deletion", "status")
			tflog.Debug(ctx, fmt.Sprintf("DeEngine wait: setup.status=%q want=%q message=%q deletion.status=%q",
				cur, want, msg, del))

			switch {
			case strings.EqualFold(cur, want):
				// For delete, also require deletion not in progress.
				if want == "idle" && strings.EqualFold(del, "in_progress") {
					return false, nil
				}
				return true, nil
			case strings.EqualFold(cur, "failed"), strings.EqualFold(cur, "error"):
				return false, fmt.Errorf("data-engine setup %s: %s", cur, msg)
			case want == "idle" && !isDataEnginePresent(cur) && !strings.EqualFold(cur, "in_progress") && !strings.EqualFold(cur, "running") && !strings.EqualFold(del, "in_progress"):
				// not_started / idle after delete — but VMS may still 409 briefly; Create retries.
				return true, nil
			}
			return false, nil
		},
	)
}

func (m *DeEngine) stateRecord(status Record) Record {
	out := Record{}
	for _, key := range []string{"setup", "deletion", "deprecated_api_used"} {
		if v, ok := status[key]; ok && v != nil {
			out[key] = v
		}
	}

	// Prefer user create params from TF state so API-expanded kafka_broker does not drift.
	for _, key := range []string{"default_topic_name", "dead_letter_topic_name", "kafka_broker", "kafka_ca_certificate", "vip_pools", deTenantNameAttr} {
		if m.tfstate == nil || !m.tfstate.HasAttribute(key) || !m.tfstate.IsKnownAndNotNull(key) {
			continue
		}
		switch key {
		case "kafka_broker":
			out[key] = m.kafkaBrokerForRequest()
		case "vip_pools":
			out[key] = m.tfstate.ToSlice(key)
		default:
			out[key] = convertAttrToRaw(m.tfstate, key)
		}
	}

	// Fill any still-missing create fields from API setup.parameters (normalized).
	params := nestedMap(status, "setup", "parameters")
	if len(params) > 0 {
		for _, key := range []string{"default_topic_name", "dead_letter_topic_name", "kafka_broker", "kafka_ca_certificate", "vip_pools"} {
			if _, ok := out[key]; ok {
				continue
			}
			if v, ok := params[key]; ok && v != nil {
				if key == "kafka_broker" {
					out[key] = normalizeKafkaBroker(v)
				} else {
					out[key] = v
				}
			}
		}
	}
	return out
}

func convertAttrToRaw(ts is.ResourceState, key string) any {
	v, ok := ts.RawValues()[key]
	if !ok || v.IsNull() || v.IsUnknown() {
		return nil
	}
	return is.ConvertAttrValueToRaw(v, ts.Type(key))
}

func normalizeKafkaBroker(v any) map[string]any {
	m, ok := v.(map[string]any)
	if !ok {
		if r, ok := v.(Record); ok {
			m = map[string]any(r)
		} else {
			return nil
		}
	}
	out := map[string]any{}
	if name, ok := m["name"]; ok && name != nil {
		out["name"] = name
	}
	if typ, ok := m["type"]; ok && typ != nil {
		out["type"] = typ
	}
	if acc, ok := m["is_accessible_to_data_engine"]; ok && acc != nil {
		out["is_accessible_to_data_engine"] = acc
	}
	return out
}
