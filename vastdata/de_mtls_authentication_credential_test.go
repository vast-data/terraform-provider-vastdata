// Copyright (c) HashiCorp, Inc.
package provider

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

func TestDeMtlsValidateResourceConfig(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"name":                      rschema.StringAttribute{Optional: true},
		"kubeconfig_data":           rschema.StringAttribute{Optional: true, WriteOnly: true},
		"context_name":              rschema.StringAttribute{Optional: true},
		"certificate_authority_b64": rschema.StringAttribute{Optional: true},
		"client_certificate_b64":    rschema.StringAttribute{Optional: true},
		"client_key_b64":            rschema.StringAttribute{Optional: true},
	}}

	newMgr := func(vals map[string]attr.Value) *DeMtlsAuthenticationCredential {
		return &DeMtlsAuthenticationCredential{tfstate: is.NewResourceState(vals, schema, &is.TFStateHints{})}
	}

	t.Run("empty_ok", func(t *testing.T) {
		m := newMgr(map[string]attr.Value{
			"kubeconfig_data":           types.StringNull(),
			"certificate_authority_b64": types.StringNull(),
			"client_certificate_b64":    types.StringNull(),
			"client_key_b64":            types.StringNull(),
		})
		if err := m.ValidateResourceConfig(t.Context()); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("kubeconfig_ok", func(t *testing.T) {
		m := newMgr(map[string]attr.Value{
			"kubeconfig_data":           types.StringValue("apiVersion: v1\n"),
			"certificate_authority_b64": types.StringNull(),
			"client_certificate_b64":    types.StringNull(),
			"client_key_b64":            types.StringNull(),
		})
		if err := m.ValidateResourceConfig(t.Context()); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("all_pems_ok", func(t *testing.T) {
		m := newMgr(map[string]attr.Value{
			"kubeconfig_data":           types.StringNull(),
			"certificate_authority_b64": types.StringValue("ca"),
			"client_certificate_b64":    types.StringValue("cert"),
			"client_key_b64":            types.StringValue("key"),
		})
		if err := m.ValidateResourceConfig(t.Context()); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("mix_rejected_via_validateOneOf", func(t *testing.T) {
		m := newMgr(map[string]attr.Value{
			"kubeconfig_data":           types.StringValue("apiVersion: v1\n"),
			"certificate_authority_b64": types.StringValue("ca"),
			"client_certificate_b64":    types.StringNull(),
			"client_key_b64":            types.StringNull(),
		})
		if err := m.ValidateResourceConfig(t.Context()); err == nil {
			t.Fatal("expected validateOneOf conflict")
		}
	})

	t.Run("partial_pems_rejected_via_validateAllOf", func(t *testing.T) {
		m := newMgr(map[string]attr.Value{
			"kubeconfig_data":           types.StringNull(),
			"certificate_authority_b64": types.StringValue("ca"),
			"client_certificate_b64":    types.StringNull(),
			"client_key_b64":            types.StringNull(),
		})
		if err := m.ValidateResourceConfig(t.Context()); err == nil {
			t.Fatal("expected validateAllOf missing fields")
		}
	})
}

func TestMtlsMaterialFromKubeconfig_DataFields(t *testing.T) {
	ca := base64.StdEncoding.EncodeToString([]byte("-----BEGIN CERTIFICATE-----\nCA\n-----END CERTIFICATE-----"))
	cert := base64.StdEncoding.EncodeToString([]byte("-----BEGIN CERTIFICATE-----\nCLIENT\n-----END CERTIFICATE-----"))
	key := base64.StdEncoding.EncodeToString([]byte("-----BEGIN PRIVATE KEY-----\nKEY\n-----END PRIVATE KEY-----"))

	kc := `
apiVersion: v1
kind: Config
current-context: de-ctx
contexts:
- name: de-ctx
  context:
    cluster: de-cluster
    user: de-user
- name: other
  context:
    cluster: other-cluster
    user: other-user
clusters:
- name: de-cluster
  cluster:
    certificate-authority-data: ` + ca + `
- name: other-cluster
  cluster:
    certificate-authority-data: ` + base64.StdEncoding.EncodeToString([]byte("OTHER")) + `
users:
- name: de-user
  user:
    client-certificate-data: ` + cert + `
    client-key-data: ` + key + `
- name: other-user
  user:
    client-certificate-data: ` + base64.StdEncoding.EncodeToString([]byte("OC")) + `
    client-key-data: ` + base64.StdEncoding.EncodeToString([]byte("OK")) + `
`
	gotCA, gotCert, gotKey, err := mtlsMaterialFromKubeconfig(kc, "")
	if err != nil {
		t.Fatal(err)
	}
	if gotCA != ca || gotCert != cert || gotKey != key {
		t.Fatalf("mismatch: ca/cert/key from current-context")
	}

	gotCA, gotCert, gotKey, err = mtlsMaterialFromKubeconfig(kc, "other")
	if err != nil {
		t.Fatal(err)
	}
	if gotCA == ca {
		t.Fatal("expected other context CA")
	}
	_ = gotCert
	_ = gotKey
}

func TestMtlsMaterialFromKubeconfig_FilePaths(t *testing.T) {
	dir := t.TempDir()
	caPath := filepath.Join(dir, "ca.pem")
	certPath := filepath.Join(dir, "cert.pem")
	keyPath := filepath.Join(dir, "key.pem")
	for path, content := range map[string]string{
		caPath:   "-----BEGIN CERTIFICATE-----\nCA\n-----END CERTIFICATE-----",
		certPath: "-----BEGIN CERTIFICATE-----\nCLIENT\n-----END CERTIFICATE-----",
		keyPath:  "-----BEGIN PRIVATE KEY-----\nKEY\n-----END PRIVATE KEY-----",
	} {
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	kc := `
current-context: c
contexts:
- name: c
  context:
    cluster: cl
    user: u
clusters:
- name: cl
  cluster:
    certificate-authority: ` + caPath + `
users:
- name: u
  user:
    client-certificate: ` + certPath + `
    client-key: ` + keyPath + `
`
	ca, cert, key, err := mtlsMaterialFromKubeconfig(kc, "c")
	if err != nil {
		t.Fatal(err)
	}
	wantCA := base64.StdEncoding.EncodeToString([]byte("-----BEGIN CERTIFICATE-----\nCA\n-----END CERTIFICATE-----"))
	if ca != wantCA {
		t.Fatalf("ca: got %q want %q", ca, wantCA)
	}
	if cert == "" || key == "" {
		t.Fatal("empty cert/key")
	}
}

func TestMtlsMaterialFromKubeconfig_Errors(t *testing.T) {
	if _, _, _, err := mtlsMaterialFromKubeconfig("not: [yaml", ""); err == nil {
		t.Fatal("expected parse error")
	}
	if _, _, _, err := mtlsMaterialFromKubeconfig("current-context: missing\n", ""); err == nil {
		t.Fatal("expected missing context")
	}
}

func TestDeMtlsTransformResponseRecord_DropsRedactedKey(t *testing.T) {
	m := &DeMtlsAuthenticationCredential{tfstate: is.NewResourceState(
		map[string]attr.Value{"tenant_name": types.StringValue("de-shared")},
		rschema.Schema{Attributes: map[string]rschema.Attribute{
			"tenant_name":               rschema.StringAttribute{Optional: true},
			"client_key_b64":            rschema.StringAttribute{Optional: true},
			"client_certificate_b64":    rschema.StringAttribute{Optional: true},
			"certificate_authority_b64": rschema.StringAttribute{Optional: true},
		}},
		&is.TFStateHints{},
	)}

	out := m.TransformResponseRecord(Record{
		"name":                      "tf-mtls-pem-repro",
		"client_key_b64":            "******",
		"client_certificate_b64":    "LS0tLS1CRUdJTi...",
		"certificate_authority_b64": "LS0tLS1CRUdJTi...",
		"tenant_guid":               "d707214e-3c70-4cae-9c89-03be8b75abf8",
	})

	if _, ok := out["client_key_b64"]; ok {
		t.Fatal("redacted client_key_b64 must be removed from the response record")
	}
	if out["client_certificate_b64"] != "LS0tLS1CRUdJTi..." {
		t.Fatalf("real client_certificate_b64 must be kept: %v", out["client_certificate_b64"])
	}
	if out["tenant_name"] != "de-shared" {
		t.Fatalf("tenant_name: got %v", out["tenant_name"])
	}
}

func TestIsRedacted(t *testing.T) {
	if !isRedacted("******") {
		t.Fatal("expected ****** to be redacted")
	}
	if isRedacted("LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQo=") {
		t.Fatal("real PEM must not be treated as redacted")
	}
}
