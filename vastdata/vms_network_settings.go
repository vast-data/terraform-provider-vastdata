// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/logging"
)

var VmsNetworkSettingsSchemaRef = is.NewSchemaReference(
	http.MethodPatch,
	"vms/{id}/network_settings",
	http.MethodPatch,
	"vms/{id}/network_settings",
)

var vmsNetworkSettingsAsyncTaskTimeout = 40 * time.Minute

// Fields accepted by PATCH /vms/{id}/network_settings/. VMS currently rejects
// partial bodies, so TF fills any omitted
// keys from a GET before sending the request.
var vmsNetworkSettingsPatchFields = []string{
	"auto_ports_ext_iface",
	"b2b_ipmi",
	"dns",
	"eth_mtu",
	"ext_netmask",
	"external_gateways",
	"hosts",
	"ib_mtu",
	"ipmi_gateway",
	"ipmi_netmask",
	"management_vips",
	"ntp",
}

type VmsNetworkSettings struct {
	tfstate is.ResourceState
}

func (m *VmsNetworkSettings) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &VmsNetworkSettings{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:    VmsNetworkSettingsSchemaRef,
			ImportFields: []string{"vms_id"},
			// OpenAPI PATCH response is only {async_task}, so request-body fields
			// would stay optional-only and FillFromRecordForImport would skip them
			// (TERF-312). Mark them computed so import/read populate state; schema
			// generation also attaches UseStateForUnknown for partial configs.
			ComputedSchemaFields: append([]string{}, vmsNetworkSettingsPatchFields...),
			ExcludedSchemaFields: []string{"async_task"},
			AsyncTaskTimeout:     &vmsNetworkSettingsAsyncTaskTimeout,
			AdditionalSchemaAttributes: map[string]any{
				"vms_id": rschema.Int64Attribute{
					Required:    true,
					Description: "ID of the VMS object.",
				},
			},
		},
	)}
}

func (m *VmsNetworkSettings) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &VmsNetworkSettings{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: VmsNetworkSettingsSchemaRef,
			NotOptionalSchemaFields: []string{
				"auto_ports_ext_iface",
				"b2b_ipmi",
				"dns",
				"eth_mtu",
				"ext_netmask",
				"external_gateways",
				"hosts",
				"ib_mtu",
				"ipmi_gateway",
				"ipmi_netmask",
				"management_vips",
				"ntp",
			},
			AdditionalSchemaAttributes: map[string]any{
				"vms_id": dschema.Int64Attribute{
					Required:    true,
					Description: "ID of the VMS object.",
				},
			},
		},
	)}
}

func (m *VmsNetworkSettings) TfState() is.ResourceState {
	return m.tfstate
}

func (m *VmsNetworkSettings) API(_ *VMSRest) VastResourceAPIWithContext {
	return nil
}

func (m *VmsNetworkSettings) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	vmsId := m.tfstate.Int64("vms_id")
	if vmsId == 0 {
		return nil, fmt.Errorf("vms_id is required")
	}
	return rest.Vms.VmsNetworkSettingsWithContext_GET(ctx, vmsId)
}

func (m *VmsNetworkSettings) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.ReadDatasource(ctx, rest)
}

func (m *VmsNetworkSettings) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	vmsId := m.tfstate.Int64("vms_id")
	params := m.tfstate.GetCreateParams()
	params.Without("vms_id")
	return updateVmsNetworkSettings(ctx, vmsId, params, rest, m)
}

func (m *VmsNetworkSettings) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	vmsId := m.tfstate.Int64("vms_id")
	planTs := plan.(*VmsNetworkSettings).TfState()
	params := planTs.GetChangedParams(m.tfstate)
	params.Without("vms_id")
	return updateVmsNetworkSettings(ctx, vmsId, params, rest, m)
}

func (m *VmsNetworkSettings) DeleteResource(_ context.Context, _ *VMSRest) error {
	// No-op. VMS network settings cannot be deleted; destroy only clears Terraform state.
	return nil
}

// TransformResponseRecord unwraps the optional API "data" envelope and maps
// GET /vms/{id}/network_settings/ (boxes[].hosts) onto the writable hosts
// attribute used by PATCH / Terraform state.
func (m *VmsNetworkSettings) TransformResponseRecord(record Record) Record {
	if nested, ok := record["data"].(map[string]any); ok {
		record = Record(nested)
	}
	if hosts := flattenNetworkSettingsHosts(record["boxes"]); hosts != nil {
		record["hosts"] = hosts
	}
	delete(record, "boxes")
	if m.tfstate != nil && m.tfstate.IsKnownAndNotNull("vms_id") {
		record["vms_id"] = m.tfstate.Int64("vms_id")
	}
	return record
}

func updateVmsNetworkSettings(ctx context.Context, vmsId int64, overlay params, rest *VMSRest, m *VmsNetworkSettings) (DisplayableRecord, error) {
	if vmsId == 0 {
		return nil, fmt.Errorf("vms_id is required")
	}

	current, err := rest.Vms.VmsNetworkSettingsWithContext_GET(ctx, vmsId)
	if err != nil {
		return nil, fmt.Errorf("failed to read VMS network settings before update: %w", err)
	}
	transformed := m.TransformResponseRecord(current)

	if len(overlay) == 0 || vmsNetworkSettingsOverlayMatches(transformed, overlay) {
		if len(overlay) > 0 {
			logging.Info(ctx, "VMS network settings already match config; skipping PATCH")
		}
		return transformed, nil
	}

	body := mergeVmsNetworkSettingsPatchBody(transformed, overlay)
	if _, err := rest.Vms.VmsNetworkSettingsWithContext_PATCH(ctx, vmsId, body, vmsNetworkSettingsAsyncTaskTimeout); err != nil {
		return nil, fmt.Errorf("failed to update VMS network settings: %w", err)
	}
	return rest.Vms.VmsNetworkSettingsWithContext_GET(ctx, vmsId)
}

// vmsNetworkSettingsOverlayMatches reports whether every overlay key already
// equals the corresponding value on the transformed GET record.
func vmsNetworkSettingsOverlayMatches(current Record, overlay params) bool {
	for key, want := range overlay {
		if !vmsNetworkSettingsValuesEqual(current[key], want) {
			return false
		}
	}
	return true
}

func vmsNetworkSettingsValuesEqual(a, b any) bool {
	return reflect.DeepEqual(normalizeVmsNetworkSettingsValue(a), normalizeVmsNetworkSettingsValue(b))
}

func normalizeVmsNetworkSettingsValue(v any) any {
	switch x := v.(type) {
	case nil:
		return nil
	case float64:
		if x == float64(int64(x)) {
			return int64(x)
		}
		return x
	case float32:
		return normalizeVmsNetworkSettingsValue(float64(x))
	case int:
		return int64(x)
	case int32:
		return int64(x)
	case []string:
		out := append([]string(nil), x...)
		sort.Strings(out)
		return out
	case []any:
		out := make([]any, len(x))
		for i, e := range x {
			out[i] = normalizeVmsNetworkSettingsValue(e)
		}
		// Order-insensitive for primitive lists (ntp, dns, …).
		if allPrimitiveScalars(out) {
			sort.Slice(out, func(i, j int) bool {
				return fmt.Sprint(out[i]) < fmt.Sprint(out[j])
			})
		}
		return out
	case map[string]any:
		out := make(map[string]any, len(x))
		for k, e := range x {
			out[k] = normalizeVmsNetworkSettingsValue(e)
		}
		return out
	default:
		return v
	}
}

// mergeVmsNetworkSettingsPatchBody builds a full PATCH body: current GET values
// for all patchable fields, then overlay (Terraform-configured changes) wins.
func mergeVmsNetworkSettingsPatchBody(current Record, overlay params) params {
	body := make(params, len(vmsNetworkSettingsPatchFields))
	for _, key := range vmsNetworkSettingsPatchFields {
		if v, ok := current[key]; ok {
			body[key] = v
		}
	}
	for key, value := range overlay {
		body[key] = value
	}
	body.Without("vms_id", "boxes", "data")
	return body
}

func flattenNetworkSettingsHosts(raw any) []any {
	boxes, ok := raw.([]any)
	if !ok || len(boxes) == 0 {
		return nil
	}
	hosts := make([]any, 0)
	for _, box := range boxes {
		boxMap, ok := box.(map[string]any)
		if !ok {
			continue
		}
		boxHosts, ok := boxMap["hosts"].([]any)
		if !ok {
			continue
		}
		hosts = append(hosts, boxHosts...)
	}
	if len(hosts) == 0 {
		return nil
	}
	return hosts
}
