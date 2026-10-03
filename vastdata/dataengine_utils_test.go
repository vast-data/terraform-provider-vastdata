// Copyright (c) HashiCorp, Inc.
package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

func TestEnsureTenantNameInRecord(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"name":        rschema.StringAttribute{Optional: true},
		"tenant_name": rschema.StringAttribute{Optional: true},
	}}
	ts := is.NewResourceState(
		map[string]attr.Value{
			"name":        types.StringValue("fn"),
			"tenant_name": types.StringValue("de-lab"),
		},
		schema,
		&is.TFStateHints{},
	)

	got := ensureTenantNameInRecord(Record{"name": "fn"}, ts, nil)
	if got["tenant_name"] != "de-lab" {
		t.Fatalf("inject: %#v", got)
	}

	kept := ensureTenantNameInRecord(Record{"name": "fn", "tenant_name": "other"}, ts, nil)
	if kept["tenant_name"] != "other" {
		t.Fatalf("keep existing: %#v", kept)
	}

	emptyState := is.NewResourceState(
		map[string]attr.Value{
			"name":        types.StringValue("fn"),
			"tenant_name": types.StringNull(),
		},
		schema,
		&is.TFStateHints{},
	)
	unchanged := ensureTenantNameInRecord(Record{"name": "fn"}, emptyState, nil)
	if _, ok := unchanged["tenant_name"]; ok {
		t.Fatalf("no tenant in state: %#v", unchanged)
	}
}

func TestResolveTenantQuery(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"name":        rschema.StringAttribute{Optional: true},
		"tenant_name": rschema.StringAttribute{Optional: true},
	}}
	ts := is.NewResourceState(
		map[string]attr.Value{
			"name":        types.StringValue("de-mtls"),
			"tenant_name": types.StringValue("de-lab"),
		},
		schema,
		&is.TFStateHints{},
	)
	q := resolveTenantQuery(ts, nil)
	if q["tenant_name"] != "de-lab" {
		t.Fatalf("from state: got %#v", q)
	}

	tsEmpty := is.NewResourceState(
		map[string]attr.Value{
			"name":        types.StringValue("de-mtls"),
			"tenant_name": types.StringNull(),
		},
		schema,
		&is.TFStateHints{},
	)
	q = resolveTenantQuery(tsEmpty, nil)
	if len(q) != 0 {
		t.Fatalf("expected empty query without session tenant, got %#v", q)
	}
}

func TestParamsUpdateAndWithout(t *testing.T) {
	merged := params{"name": "x"}
	merged.Update(params{"tenant_name": "de-lab"}, true)
	if merged["name"] != "x" || merged["tenant_name"] != "de-lab" {
		t.Fatalf("merge: %#v", merged)
	}
	merged = params{"tenant_name": "keep"}
	merged.Update(params{"tenant_name": "other"}, true)
	if merged["tenant_name"] != "keep" {
		t.Fatalf("overwrite: %#v", merged)
	}
	body := params{"name": "x", "tenant_name": "de-lab"}
	body.Without("tenant_name")
	if _, ok := body["tenant_name"]; ok {
		t.Fatalf("strip failed: %#v", body)
	}
}

func TestDeTransformQueryParams(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"tenant_name": rschema.StringAttribute{Optional: true},
	}}
	ts := is.NewResourceState(
		map[string]attr.Value{"tenant_name": types.StringValue("de-lab")},
		schema,
		&is.TFStateHints{},
	)

	got := deTransformQueryParams(params{"name": "fn"}, ts, nil)
	if got["name"] != "fn" || got["tenant_name"] != "de-lab" {
		t.Fatalf("merge: %#v", got)
	}

	kept := deTransformQueryParams(params{"tenant_name": "keep"}, ts, nil)
	if kept["tenant_name"] != "keep" {
		t.Fatalf("must not overwrite existing tenant_name: %#v", kept)
	}
}

func TestDeTransformRequestBody(t *testing.T) {
	body := params{"name": "fn", "tenant_name": "de-lab"}
	got := deTransformRequestBody(body)
	if _, ok := got["tenant_name"]; ok {
		t.Fatalf("tenant_name should be stripped: %#v", got)
	}
	if got["name"] != "fn" {
		t.Fatalf("other fields kept: %#v", got)
	}
}
