// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

func TestFinalizeUserKeyRecord_GeneratedKeepsSecret(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"secret_key":           rschema.StringAttribute{Optional: true, Computed: true, Sensitive: true},
		"encrypted_secret_key": rschema.StringAttribute{Computed: true},
		"access_key":           rschema.StringAttribute{Optional: true, Computed: true},
		"pgp_public_key":       rschema.StringAttribute{Optional: true, Sensitive: true},
	}}
	ts := is.NewTFStateMust(map[string]attr.Value{
		"secret_key":           types.StringNull(),
		"encrypted_secret_key": types.StringNull(),
		"access_key":           types.StringNull(),
		"pgp_public_key":       types.StringNull(),
	}, schema, &is.TFStateHints{})

	out, err := finalizeUserKeyRecord(Record{
		"access_key": "AKIA_GENERATED",
		"secret_key": "cluster-generated-secret",
	}, ts)
	require.NoError(t, err)
	assert.Equal(t, "cluster-generated-secret", out["secret_key"])
	assert.NotContains(t, ts.StateHints().WriteOnlyFields, "secret_key")
}

func TestFinalizeUserKeyRecord_UserSuppliedKeepsSecret(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"secret_key":           rschema.StringAttribute{Optional: true, Computed: true, Sensitive: true},
		"encrypted_secret_key": rschema.StringAttribute{Computed: true},
		"access_key":           rschema.StringAttribute{Optional: true, Computed: true},
		"pgp_public_key":       rschema.StringAttribute{Optional: true, Sensitive: true},
	}}
	ts := is.NewTFStateMust(map[string]attr.Value{
		"secret_key":           types.StringValue("user-supplied-secret"),
		"encrypted_secret_key": types.StringNull(),
		"access_key":           types.StringValue("AKIA_USER"),
		"pgp_public_key":       types.StringNull(),
	}, schema, &is.TFStateHints{})

	out, err := finalizeUserKeyRecord(Record{
		"access_key": "AKIA_USER",
		"secret_key": "ignored-response",
	}, ts)
	require.NoError(t, err)
	assert.Equal(t, "user-supplied-secret", out["secret_key"])
	assert.NotContains(t, out, "secret_key_md5")
	assert.NotContains(t, ts.StateHints().WriteOnlyFields, "secret_key")
}

func TestValidateSecretKeyFormat(t *testing.T) {
	assert.NoError(t, validateSecretKeyFormat("TfProbeSecret3cccccccccccccccccccccccccc"))
	assert.Error(t, validateSecretKeyFormat("TfProbeSecret3ccccccccccccccccccccccccccc")) // 41
	assert.Error(t, validateSecretKeyFormat("short"))
	assert.Error(t, validateSecretKeyFormat("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA!!")) // 40 but invalid chars
}

func TestUserKeyHints_NotSchemaWriteOnly(t *testing.T) {
	m := (&UserKey{}).NewResourceManager(map[string]attr.Value{}, nil)
	hints := m.TfState().StateHints()
	require.NotNil(t, hints)
	assert.NotContains(t, hints.WriteOnlyFields, "secret_key")
	assert.NotContains(t, hints.WriteOnlyMD5Fields, "secret_key")
	assert.NotContains(t, hints.AdditionalSchemaAttributes, "secret_key_md5")
	assert.Contains(t, hints.PreserveUserValueFields, "secret_key")
	assert.Contains(t, hints.ComputedSchemaFields, "secret_key")
}
