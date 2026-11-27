// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/client"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/schema_generation"
)

func newTestClusterEkm(t *testing.T, overrides map[string]attr.Value) *ClusterEkm {
	t.Helper()
	proto := (&ClusterEkm{}).NewResourceManager(nil, nil)
	schema, err := schema_generation.GetResourceSchema(context.Background(), proto.TfState().StateHints())
	require.NoError(t, err)

	raw := make(map[string]attr.Value, len(schema.Attributes))
	for name, a := range schema.Attributes {
		switch a.(type) {
		case rschema.Int64Attribute:
			raw[name] = types.Int64Null()
		case rschema.BoolAttribute:
			raw[name] = types.BoolNull()
		default:
			raw[name] = types.StringNull()
		}
	}
	for k, v := range overrides {
		raw[k] = v
	}
	return (&ClusterEkm{}).NewResourceManager(raw, *schema).(*ClusterEkm)
}

func newMockClusterEkmRest(t *testing.T, handler http.HandlerFunc) *VMSRest {
	t.Helper()
	srv := httptest.NewTLSServer(handler)
	t.Cleanup(srv.Close)

	u, err := url.Parse(srv.URL)
	require.NoError(t, err)
	port, err := strconv.ParseInt(u.Port(), 10, 64)
	require.NoError(t, err)

	rest, err := client.NewRest(u.Hostname(), port, "", "", "token", "", false, "test", 5*time.Second)
	require.NoError(t, err)
	return rest
}

func TestClusterEkm_CreateCallsAddEkm(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any

	rest := newMockClusterEkmRest(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusNoContent)
	})

	mgr := newTestClusterEkm(t, map[string]attr.Value{
		"cluster_id":         types.Int64Value(1),
		"ekm_servers":        types.StringValue("ekm.example:5696"),
		"encryption_type":    types.StringValue("CIPHER_TRUST_KMIP"),
		"ekm_auth_domain":    types.StringValue("terraform_test_domain"),
		"ekm_certificate":    types.StringValue("CERT-V1"),
		"ekm_private_key":    types.StringValue("KEY-V1"),
		"ekm_ca_certificate": types.StringValue("CA-V1"),
	})

	_, err := mgr.CreateResource(context.Background(), rest)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/api/latest/clusters/1/add_ekm/", gotPath)
	assert.Equal(t, map[string]any{
		"ekm_servers":        "ekm.example:5696",
		"encryption_type":    "CIPHER_TRUST_KMIP",
		"ekm_auth_domain":    "terraform_test_domain",
		"ekm_certificate":    "CERT-V1",
		"ekm_private_key":    "KEY-V1",
		"ekm_ca_certificate": "CA-V1",
	}, gotBody)
}

func TestClusterEkm_UpdateCallsSetCertificatesWithFullPayload(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any

	rest := newMockClusterEkmRest(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusNoContent)
	})

	// Prior state: write-only certs are cleared, MD5s stored, non-secret fields intact.
	state := newTestClusterEkm(t, map[string]attr.Value{
		"cluster_id":             types.Int64Value(1),
		"ekm_servers":            types.StringValue("ekm.example:5696"),
		"encryption_type":        types.StringValue("CIPHER_TRUST_KMIP"),
		"ekm_auth_domain":        types.StringValue("terraform_test_domain"),
		"ekm_certificate_md5":    types.StringValue(ContentMD5Hex("CERT-V1")),
		"ekm_private_key_md5":    types.StringValue(ContentMD5Hex("KEY-V1")),
		"ekm_ca_certificate_md5": types.StringValue(ContentMD5Hex("CA-V1")),
	})

	// Planned state (after overlayWriteOnlyFromConfig): only ekm_certificate changed.
	plan := newTestClusterEkm(t, map[string]attr.Value{
		"cluster_id":             types.Int64Value(1),
		"ekm_servers":            types.StringValue("ekm.example:5696"),
		"encryption_type":        types.StringValue("CIPHER_TRUST_KMIP"),
		"ekm_auth_domain":        types.StringValue("terraform_test_domain"),
		"ekm_certificate":        types.StringValue("CERT-V2"),
		"ekm_private_key":        types.StringValue("KEY-V1"),
		"ekm_ca_certificate":     types.StringValue("CA-V1"),
		"ekm_certificate_md5":    types.StringValue(ContentMD5Hex("CERT-V2")),
		"ekm_private_key_md5":    types.StringValue(ContentMD5Hex("KEY-V1")),
		"ekm_ca_certificate_md5": types.StringValue(ContentMD5Hex("CA-V1")),
	})

	_, err := state.UpdateResource(context.Background(), plan, rest)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/api/latest/clusters/1/set_certificates/", gotPath)
	assert.Equal(t, map[string]any{
		"ekm_servers":        "ekm.example:5696",
		"encryption_type":    "CIPHER_TRUST_KMIP",
		"ekm_auth_domain":    "terraform_test_domain",
		"ekm_certificate":    "CERT-V2",
		"ekm_private_key":    "KEY-V1",
		"ekm_ca_certificate": "CA-V1",
	}, gotBody)
	assert.NotContains(t, gotBody, "cluster_id")
	assert.NotContains(t, gotBody, "ekm_certificate_md5")
	assert.NotContains(t, gotBody, "ekm_private_key_md5")
	assert.NotContains(t, gotBody, "ekm_ca_certificate_md5")
}

func TestClusterEkm_UpdateEmptyClusterIDReturnsError(t *testing.T) {
	called := false
	rest := newMockClusterEkmRest(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})

	state := newTestClusterEkm(t, map[string]attr.Value{
		"cluster_id": types.Int64Value(0),
	})
	plan := newTestClusterEkm(t, map[string]attr.Value{
		"cluster_id":      types.Int64Value(0),
		"ekm_certificate": types.StringValue("CERT-V2"),
	})

	_, err := state.UpdateResource(context.Background(), plan, rest)
	require.Error(t, err)
	assert.False(t, called)
}
