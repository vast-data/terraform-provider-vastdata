// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

func TestPopulateIDFieldsFromNestedObjects(t *testing.T) {
	tests := []struct {
		name           string
		tfstateKeys    map[string]attr.Value
		record         Record
		expectedRecord Record
		description    string
	}{
		{
			name: "extract local_provider_id from local_provider",
			tfstateKeys: map[string]attr.Value{
				"id":                types.Int64Value(123),
				"name":              types.StringValue("test-group"),
				"local_provider_id": types.Int64Null(),
			},
			record: Record{
				"id":   int64(123),
				"name": "test-group",
				"local_provider": map[string]any{
					"id":   int64(5),
					"name": "provider-name",
				},
			},
			expectedRecord: Record{
				"id":   int64(123),
				"name": "test-group",
				"local_provider": map[string]any{
					"id":   int64(5),
					"name": "provider-name",
				},
				"local_provider_id": int64(5),
			},
			description: "Should extract local_provider_id from local_provider.id",
		},
		{
			name: "extract tenant_id from tenant",
			tfstateKeys: map[string]attr.Value{
				"id":        types.Int64Value(456),
				"name":      types.StringValue("test-role"),
				"tenant_id": types.Int64Null(),
			},
			record: Record{
				"id":   int64(456),
				"name": "test-role",
				"tenant": map[string]any{
					"id":   int64(28),
					"name": "tenant-name",
					"guid": "abc-123",
				},
			},
			expectedRecord: Record{
				"id":   int64(456),
				"name": "test-role",
				"tenant": map[string]any{
					"id":   int64(28),
					"name": "tenant-name",
					"guid": "abc-123",
				},
				"tenant_id": int64(28),
			},
			description: "Should extract tenant_id from tenant.id",
		},
		{
			name: "multiple _id fields",
			tfstateKeys: map[string]attr.Value{
				"id":                 types.Int64Value(789),
				"name":               types.StringValue("test-snapshot"),
				"loanee_tenant_id":   types.Int64Null(),
				"remote_target_id":   types.Int64Null(),
				"loanee_snapshot_id": types.Int64Null(),
			},
			record: Record{
				"id":   int64(789),
				"name": "test-snapshot",
				"loanee_tenant": map[string]any{
					"id":   int64(1),
					"name": "default",
					"guid": "xyz-456",
				},
				"remote_target": map[string]any{
					"id":   int64(10),
					"name": "remote-cluster",
				},
				"loanee_snapshot": map[string]any{
					"id":   int64(42),
					"name": "snap-1",
				},
			},
			expectedRecord: Record{
				"id":   int64(789),
				"name": "test-snapshot",
				"loanee_tenant": map[string]any{
					"id":   int64(1),
					"name": "default",
					"guid": "xyz-456",
				},
				"loanee_tenant_id": int64(1),
				"remote_target": map[string]any{
					"id":   int64(10),
					"name": "remote-cluster",
				},
				"remote_target_id": int64(10),
				"loanee_snapshot": map[string]any{
					"id":   int64(42),
					"name": "snap-1",
				},
				"loanee_snapshot_id": int64(42),
			},
			description: "Should extract multiple _id fields from nested objects",
		},
		{
			name: "id field already exists in record",
			tfstateKeys: map[string]attr.Value{
				"id":                types.Int64Value(111),
				"name":              types.StringValue("test-existing"),
				"local_provider_id": types.Int64Value(99),
			},
			record: Record{
				"id":                int64(111),
				"name":              "test-existing",
				"local_provider_id": int64(99), // Already exists
				"local_provider": map[string]any{
					"id":   int64(5),
					"name": "provider-name",
				},
			},
			expectedRecord: Record{
				"id":                int64(111),
				"name":              "test-existing",
				"local_provider_id": int64(99), // Should remain unchanged
				"local_provider": map[string]any{
					"id":   int64(5),
					"name": "provider-name",
				},
			},
			description: "Should not overwrite existing _id field",
		},
		{
			name: "parent object does not exist",
			tfstateKeys: map[string]attr.Value{
				"id":                types.Int64Value(222),
				"name":              types.StringValue("test-no-parent"),
				"local_provider_id": types.Int64Null(),
			},
			record: Record{
				"id":   int64(222),
				"name": "test-no-parent",
				// No local_provider object
			},
			expectedRecord: Record{
				"id":   int64(222),
				"name": "test-no-parent",
				// local_provider_id should not be added
			},
			description: "Should not add _id field if parent object doesn't exist",
		},
		{
			name: "parent object is null",
			tfstateKeys: map[string]attr.Value{
				"id":                types.Int64Value(333),
				"name":              types.StringValue("test-null-parent"),
				"local_provider_id": types.Int64Null(),
			},
			record: Record{
				"id":             int64(333),
				"name":           "test-null-parent",
				"local_provider": nil,
			},
			expectedRecord: Record{
				"id":             int64(333),
				"name":           "test-null-parent",
				"local_provider": nil,
			},
			description: "Should not add _id field if parent object is nil",
		},
		{
			name: "parent object is not a map",
			tfstateKeys: map[string]attr.Value{
				"id":                types.Int64Value(444),
				"name":              types.StringValue("test-invalid-parent"),
				"local_provider_id": types.Int64Null(),
			},
			record: Record{
				"id":             int64(444),
				"name":           "test-invalid-parent",
				"local_provider": "not-a-map",
			},
			expectedRecord: Record{
				"id":             int64(444),
				"name":           "test-invalid-parent",
				"local_provider": "not-a-map",
			},
			description: "Should not add _id field if parent is not a map",
		},
		{
			name: "parent object has no id field",
			tfstateKeys: map[string]attr.Value{
				"id":                types.Int64Value(555),
				"name":              types.StringValue("test-no-id"),
				"local_provider_id": types.Int64Null(),
			},
			record: Record{
				"id":   int64(555),
				"name": "test-no-id",
				"local_provider": map[string]any{
					"name": "provider-name",
					// No "id" field
				},
			},
			expectedRecord: Record{
				"id":   int64(555),
				"name": "test-no-id",
				"local_provider": map[string]any{
					"name": "provider-name",
				},
			},
			description: "Should not add _id field if parent has no id field",
		},
		{
			name: "parent id field is nil",
			tfstateKeys: map[string]attr.Value{
				"id":                types.Int64Value(666),
				"name":              types.StringValue("test-nil-id"),
				"local_provider_id": types.Int64Null(),
			},
			record: Record{
				"id":   int64(666),
				"name": "test-nil-id",
				"local_provider": map[string]any{
					"id":   nil,
					"name": "provider-name",
				},
			},
			expectedRecord: Record{
				"id":   int64(666),
				"name": "test-nil-id",
				"local_provider": map[string]any{
					"id":   nil,
					"name": "provider-name",
				},
			},
			description: "Should not add _id field if parent id is nil",
		},
		{
			name: "no _id fields in tfstate",
			tfstateKeys: map[string]attr.Value{
				"id":   types.Int64Value(777),
				"name": types.StringValue("test-no-id-fields"),
			},
			record: Record{
				"id":   int64(777),
				"name": "test-no-id-fields",
				"local_provider": map[string]any{
					"id":   int64(5),
					"name": "provider-name",
				},
			},
			expectedRecord: Record{
				"id":   int64(777),
				"name": "test-no-id-fields",
				"local_provider": map[string]any{
					"id":   int64(5),
					"name": "provider-name",
				},
			},
			description: "Should not modify record if no _id fields in tfstate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a minimal tfstate with the test keys
			tfstate := &is.TFState{
				Raw: tt.tfstateKeys,
			}

			// Make a copy of the record for testing
			recordCopy := make(Record)
			for k, v := range tt.record {
				recordCopy[k] = v
			}

			// Call the function
			ctx := context.Background()
			PopulateIDFieldsFromNestedObjects(ctx, tfstate, recordCopy)

			// Verify the result
			for expectedKey, expectedValue := range tt.expectedRecord {
				actualValue, exists := recordCopy[expectedKey]
				if !exists {
					if expectedValue != nil {
						t.Errorf("%s: expected key %q to exist in record, but it doesn't", tt.description, expectedKey)
					}
					continue
				}

				// Use deep equality for comparison (handles maps, slices, etc.)
				if !reflect.DeepEqual(expectedValue, actualValue) {
					t.Errorf("%s: for key %q, expected value %v, got %v", tt.description, expectedKey, expectedValue, actualValue)
				}
			}

			// Verify no unexpected keys were added
			for actualKey := range recordCopy {
				if _, expected := tt.expectedRecord[actualKey]; !expected {
					t.Errorf("%s: unexpected key %q was added to record", tt.description, actualKey)
				}
			}
		})
	}
}

func TestPopulateIDFieldsFromNestedObjects_NilInputs(t *testing.T) {
	ctx := context.Background()

	t.Run("nil record", func(t *testing.T) {
		tfstate := &is.TFState{
			Raw: map[string]attr.Value{
				"local_provider_id": types.Int64Null(),
			},
		}
		// Should not panic
		PopulateIDFieldsFromNestedObjects(ctx, tfstate, nil)
	})

	t.Run("nil tfstate", func(t *testing.T) {
		record := Record{
			"id": int64(123),
		}
		// Should not panic
		PopulateIDFieldsFromNestedObjects(ctx, nil, record)
	})

	t.Run("both nil", func(t *testing.T) {
		// Should not panic
		PopulateIDFieldsFromNestedObjects(ctx, nil, nil)
	})
}

func TestContentMD5Hex(t *testing.T) {
	a := ContentMD5Hex("key-a")
	b := ContentMD5Hex("key-b")
	if a == "" || len(a) != 32 {
		t.Fatalf("ContentMD5Hex(key-a) = %q, want 32-char hex", a)
	}
	if a == b {
		t.Fatalf("expected different hashes for different inputs")
	}
	if a != ContentMD5Hex("key-a") {
		t.Fatalf("ContentMD5Hex not stable")
	}
}

func TestBytesMD5Hex(t *testing.T) {
	if ContentMD5Hex("abc") != BytesMD5Hex([]byte("abc")) {
		t.Fatalf("ContentMD5Hex and BytesMD5Hex disagree")
	}
}

func TestSetRecordContentMD5(t *testing.T) {
	rec := Record{"license_key": "secret-value", "state": "ACTIVE"}
	SetRecordContentMD5(rec, "license_key", "license_key_md5")
	if rec["license_key_md5"] != ContentMD5Hex("secret-value") {
		t.Fatalf("license_key_md5 = %v, want %s", rec["license_key_md5"], ContentMD5Hex("secret-value"))
	}
	if rec["state"] != "ACTIVE" {
		t.Fatalf("state mutated: %v", rec["state"])
	}

	SetRecordContentMD5(nil, "license_key", "license_key_md5") // no panic
	empty := Record{"license_key": ""}
	SetRecordContentMD5(empty, "license_key", "license_key_md5")
	if _, has := empty["license_key_md5"]; has {
		t.Fatalf("empty source should not set md5")
	}
}

// ----------------------------------
// Runtime-agnostic registry
// ----------------------------------

var cacheProbeCalls int

func cacheProbe() int {
	cacheProbeCalls++
	return cacheProbeCalls
}

func TestCached_RunsOnceAndNamesTheFunction(t *testing.T) {
	cacheProbeCalls = 0
	fn := cached(cacheProbe)
	if got := fn(); got != 1 {
		t.Fatalf("first call = %d", got)
	}
	if got := fn(); got != 1 {
		t.Fatalf("second call = %d, want the cached value", got)
	}
	if cacheProbeCalls != 1 {
		t.Fatalf("wrapped function ran %d times", cacheProbeCalls)
	}
	if got := callableName(cacheProbe); got != "cacheProbe" {
		t.Fatalf("callableName = %q", got)
	}
}

func TestResourceManagerPrototypes_ContainsKnownResources(t *testing.T) {
	protos := ResourceManagerPrototypes()
	if reflect.ValueOf(protos).Pointer() != reflect.ValueOf(ResourceManagerPrototypes()).Pointer() {
		t.Fatal("ResourceManagerPrototypes rebuilt the map")
	}
	for _, name := range []string{"user", "vip_pool", "license", "tenant", "view"} {
		if _, ok := protos[name]; !ok {
			t.Errorf("expected resource %q in registry", name)
		}
	}

	// The registry must expose exactly what Terraform registers as resources.
	if got, want := len(protos), len(GetResourceFactories()); got != want {
		t.Errorf("registry has %d resources, Terraform registers %d", got, want)
	}
}

func TestNewResourceManagerByName(t *testing.T) {
	ctx := context.Background()

	m, err := NewResourceManagerByName(ctx, "user")
	if err != nil {
		t.Fatalf("NewResourceManagerByName(user): %v", err)
	}
	if !m.TfState().IsEnabled() || !m.TfState().HasAttribute("name") {
		t.Fatalf("expected a schema-backed state with a 'name' attribute")
	}
	if _, ok := m.TfState().(*is.TFState); !ok {
		t.Fatalf("default state must be *TFState, got %T", m.TfState())
	}

	if _, err := NewResourceManagerByName(ctx, "does_not_exist"); err == nil {
		t.Fatalf("expected an error for an unknown resource")
	}
}

// otherRuntimeState stands in for CrossPlaneState: it lives outside the internalstate
// package, embeds *TFState and overrides only what differs.
type otherRuntimeState struct{ *is.TFState }

func (o *otherRuntimeState) IsUnknown(string) bool { return false }
func (o *otherRuntimeState) Copy() is.ResourceState {
	return &otherRuntimeState{TFState: o.TFState.CopyTFState()}
}

func TestStateFactory_ReplacesStateForEveryManager(t *testing.T) {
	is.SetStateFactory(func(raw map[string]attr.Value, schema any, hints *is.TFStateHints) (is.ResourceState, error) {
		base, err := is.NewTFState(raw, schema, schemaKindOf(schema), hints)
		if err != nil {
			return nil, err
		}
		return &otherRuntimeState{TFState: base}, nil
	})
	t.Cleanup(func() { is.SetStateFactory(nil) })

	// No manager constructor was touched: they all pick the new implementation up.
	for _, name := range []string{"user", "vip_pool", "tenant"} {
		m, err := NewResourceManagerByName(context.Background(), name)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if _, ok := m.TfState().(*otherRuntimeState); !ok {
			t.Errorf("%s: expected *otherRuntimeState, got %T", name, m.TfState())
		}
	}
}

// ----------------------------------
// Framework-agnostic lifecycle
// ----------------------------------

// hookManager is a manager that implements every CRUD hook and records the call order.
type hookManager struct {
	state is.ResourceState
	log   *[]string

	createErr error
}

func (h *hookManager) record(s string) { *h.log = append(*h.log, s) }

func (h *hookManager) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return h
}
func (h *hookManager) TfState() is.ResourceState                 { return h.state }
func (h *hookManager) API(_ *VMSRest) VastResourceAPIWithContext { return nil }

func (h *hookManager) PrepareCreateResource(context.Context, *VMSRest) error {
	h.record("PrepareCreate")
	return nil
}
func (h *hookManager) CreateResource(context.Context, *VMSRest) (DisplayableRecord, error) {
	h.record("Create")
	if h.createErr != nil {
		return nil, h.createErr
	}
	return Record{"id": int64(42), "name": h.state.String("name")}, nil
}
func (h *hookManager) AfterCreateResource(context.Context, *VMSRest, Record) error {
	h.record("AfterCreate")
	return nil
}
func (h *hookManager) ReadResource(context.Context, *VMSRest) (DisplayableRecord, error) {
	h.record("Read")
	if h.state.String("name") == "gone" {
		return nil, ForceCleanState{}
	}
	return Record{"id": int64(42), "name": "observed"}, nil
}
func (h *hookManager) UpdateResource(_ context.Context, desired UpdateResource, _ *VMSRest) (DisplayableRecord, error) {
	h.record("Update->" + desired.(*hookManager).state.String("name"))
	return Record{"id": int64(42), "name": desired.(*hookManager).state.String("name")}, nil
}
func (h *hookManager) DeleteResource(context.Context, *VMSRest) error {
	h.record("Delete")
	return nil
}
func (h *hookManager) AfterDeleteResource(context.Context, *VMSRest) error {
	h.record("AfterDelete")
	return nil
}

func newHookManager(t *testing.T, log *[]string, name string) *hookManager {
	t.Helper()
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"id":   rschema.Int64Attribute{Computed: true},
		"name": rschema.StringAttribute{Required: true},
	}}
	raw := map[string]attr.Value{"id": types.Int64Null(), "name": types.StringValue(name)}
	hints := &is.TFStateHints{TFStateHintsForCustom: &is.TFStateHintsForCustom{
		Description: "test", MarkdownDescription: "test",
		SchemaAttributes: map[string]any{"id": schema.Attributes["id"], "name": schema.Attributes["name"]},
	}}
	return &hookManager{state: is.NewResourceState(raw, schema, hints), log: log}
}

func TestResourceLifecycle_Create_RunsHooksInOrderAndFillsState(t *testing.T) {
	var log []string
	m := newHookManager(t, &log, "alpha")

	if err := NewResourceLifecycle(nil, "hook").Create(context.Background(), m); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got, want := log, []string{"PrepareCreate", "Create", "AfterCreate"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("hook order = %v, want %v", got, want)
	}
	if got := m.state.Int64("id"); got != 42 {
		t.Fatalf("state id = %d, want 42 (filled from the API record)", got)
	}
}

func TestResourceLifecycle_Create_ErrorIsWrappedAndDoesNotRollbackCustomCreate(t *testing.T) {
	var log []string
	m := newHookManager(t, &log, "alpha")
	boom := errors.New("boom")
	m.createErr = boom

	err := NewResourceLifecycle(nil, "hook").Create(context.Background(), m)
	if !errors.Is(err, boom) {
		t.Fatalf("expected wrapped boom, got %v", err)
	}
	if !strings.Contains(err.Error(), "create[hook]") {
		t.Fatalf("error should name the stage and resource, got %q", err)
	}
	for _, entry := range log {
		if entry == "Delete" || entry == "AfterDelete" {
			t.Fatalf("custom CreateResource failures must not trigger rollback, log=%v", log)
		}
	}
}

func TestFindHook_SetterOverridesAndKeepsInnerHook(t *testing.T) {
	var log []string
	inner := newHookManager(t, &log, "alpha")
	wrapped := WithReadResource(inner, func(context.Context, *VMSRest) (DisplayableRecord, error) {
		return Record{"id": int64(7), "name": "from-setter"}, nil
	})
	wrapped = WithAfterDeleteResource(wrapped, func(context.Context, *VMSRest) error {
		inner.record("setter-after-delete")
		return nil
	})
	if _, ok := WithCreateResource(wrapped, nil).(createResourceHook); ok {
		t.Fatal("a nil function must leave the manager unchanged")
	}

	reader, ok := findHook[ReadResource](wrapped)
	if !ok {
		t.Fatal("ReadResource setter was not found")
	}
	rec, err := reader.ReadResource(context.Background(), nil)
	if err != nil || rec.(Record)["name"] != "from-setter" {
		t.Fatalf("setter read = (%v, %v)", rec, err)
	}
	if _, ok := findHook[CreateResource](wrapped); !ok {
		t.Fatal("wrappers hid the inner CreateResource")
	}
	after, ok := findHook[AfterDeleteResource](wrapped)
	if !ok {
		t.Fatal("AfterDeleteResource setter was not found")
	}
	if err := after.AfterDeleteResource(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if len(log) != 1 || log[0] != "setter-after-delete" {
		t.Fatalf("log = %v", log)
	}
}

func TestHasOwnHook_IgnoresWrappers(t *testing.T) {
	var log []string
	inner := newHookManager(t, &log, "alpha")
	wrapped := WithAfterDeleteResource(inner, func(context.Context, *VMSRest) error { return nil })

	if !hasOwnHook[AfterDeleteResource](wrapped) {
		t.Fatal("the original manager implements AfterDeleteResource")
	}
	// A manager that does not implement the hook must not look like it does just
	// because a wrapper added one from outside.
	bare := &customDeleteManager{}
	wrappedBare := WithAfterCreateResource(bare, func(context.Context, *VMSRest, Record) error { return nil })
	if hasOwnHook[AfterCreateResource](wrappedBare) {
		t.Fatal("a setter-provided hook must not count as the resource's own hook")
	}
	if _, ok := findHook[AfterCreateResource](wrappedBare); !ok {
		t.Fatal("the setter hook must still be found")
	}
}

func TestForEachHook_RunsSetterAndOriginalAndJoinsErrors(t *testing.T) {
	var log []string
	inner := newHookManager(t, &log, "alpha")
	wrapped := WithReadResource(inner, func(context.Context, *VMSRest) (DisplayableRecord, error) {
		return Record{"id": int64(7), "name": "from-setter"}, nil
	})
	wrapped = WithAfterDeleteResource(wrapped, func(context.Context, *VMSRest) error {
		inner.record("setter-after-delete")
		return errors.New("setter-boom")
	})

	err := forEachHook(wrapped, func(imp AfterDeleteResource) error {
		return imp.AfterDeleteResource(context.Background(), nil)
	})
	if got, want := log, []string{"setter-after-delete", "AfterDelete"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("hook order = %v, want %v", got, want)
	}
	if err == nil || !strings.Contains(err.Error(), "setter-boom") {
		t.Fatalf("expected joined setter error, got %v", err)
	}

	// Unrelated wrappers must not make forEachHook re-run promoted CRUD hooks.
	var creates int
	err = forEachHook(wrapped, func(imp CreateResource) error {
		creates++
		_, createErr := imp.CreateResource(context.Background(), nil)
		return createErr
	})
	if err != nil {
		t.Fatalf("CreateResource forEachHook: %v", err)
	}
	if creates != 1 {
		t.Fatalf("CreateResource invocations = %d, want 1 (original only; no promotion duplicates)", creates)
	}
}

func TestResourceLifecycle_Read(t *testing.T) {
	var log []string
	lc := NewResourceLifecycle(nil, "hook")

	m := newHookManager(t, &log, "alpha")
	exists, err := lc.Read(context.Background(), m)
	if err != nil || !exists {
		t.Fatalf("Read = (%v, %v), want (true, nil)", exists, err)
	}
	if got := m.state.String("name"); got != "observed" {
		t.Fatalf("state name = %q, want the read value", got)
	}

	gone := newHookManager(t, &log, "gone")
	exists, err = lc.Read(context.Background(), gone)
	if err != nil || exists {
		t.Fatalf("Read(gone) = (%v, %v), want (false, nil)", exists, err)
	}
}

func TestResourceLifecycle_UpdateAndDelete(t *testing.T) {
	var log []string
	lc := NewResourceLifecycle(nil, "hook")

	current := newHookManager(t, &log, "old")
	desired := newHookManager(t, &log, "new")

	if err := lc.Update(context.Background(), current, desired); err != nil {
		t.Fatalf("Update: %v", err)
	}
	if got := current.state.String("name"); got != "new" {
		t.Fatalf("current name = %q, want desired value after update", got)
	}

	if err := lc.Delete(context.Background(), current); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if got, want := log, []string{"Update->new", "Delete", "AfterDelete"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("hook order = %v, want %v", got, want)
	}
}

type customDeleteManager struct {
	st     is.ResourceState
	called bool
}

func (m *customDeleteManager) NewResourceManager(map[string]attr.Value, any) ResourceManager {
	return m
}
func (m *customDeleteManager) TfState() is.ResourceState               { return m.st }
func (m *customDeleteManager) API(*VMSRest) VastResourceAPIWithContext { return nil }
func (m *customDeleteManager) DeleteResource(context.Context, *VMSRest) error {
	m.called = true
	return nil
}

func TestDelete_CustomMethodRuns(t *testing.T) {
	st, err := is.NewTFState(map[string]attr.Value{}, rschema.Schema{}, is.SchemaForResource, nil)
	if err != nil {
		t.Fatalf("state: %v", err)
	}
	m := &customDeleteManager{st: st}
	if err := NewResourceLifecycle(nil, "api_token").Delete(context.Background(), m); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if !m.called {
		t.Fatal("Delete must run a custom DeleteResource")
	}
}

func TestDelete_EmptyDeleteResourceSkipsAPI(t *testing.T) {
	ctx := context.Background()
	m, err := NewResourceManagerByName(ctx, "callhome_config")
	if err != nil {
		t.Fatalf("manager: %v", err)
	}
	// A nil REST client panics if Delete reaches the API. The empty method returns first.
	if err := NewResourceLifecycle(nil, "callhome_config").Delete(ctx, m); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestRecordByGUID(t *testing.T) {
	records := RecordSet{
		{"id": int64(6), "guid": "aaa", "name": "other"},
		{"id": int64(7), "guid": "bbb", "name": "wanted"},
		nil,
	}

	got, err := recordByGUID(records, "bbb")
	if err != nil {
		t.Fatalf("matching guid: %v", err)
	}
	if got["name"] != "wanted" {
		t.Fatalf("record = %v, want the one whose guid matches", got)
	}

	_, err = recordByGUID(records, "missing")
	if !expectStatusCodes(err, http.StatusNotFound) {
		t.Fatalf("missing guid error = %v, want a 404 ApiError", err)
	}
}

func schemaKindOf(schema any) is.SchemaContext {
	if _, ok := schema.(rschema.Schema); ok {
		return is.SchemaForResource
	}
	return is.SchemaForDataSource
}
