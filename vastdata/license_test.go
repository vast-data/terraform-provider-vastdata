// Copyright (c) HashiCorp, Inc.
package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/stretchr/testify/assert"
)

func TestLicense_TransformResponseRecord_SetsMD5AndDropsKey(t *testing.T) {
	m := &License{}
	key := "gAAAAABq7b3_test-license-key"
	out := m.TransformResponseRecord(Record{
		"id":          int64(12),
		"license_key": key,
		"state":       "ACTIVE",
	})
	_, hasKey := out["license_key"]
	assert.False(t, hasKey, "license_key must not remain in the record")
	assert.Equal(t, ContentMD5Hex(key), out["license_key_md5"])
	assert.Equal(t, "ACTIVE", out["state"])
}

func TestLicense_ResourceHints_WriteOnlyMD5(t *testing.T) {
	m := (&License{}).NewResourceManager(map[string]attr.Value{}, nil)
	hints := m.TfState().StateHints()
	assert.Contains(t, hints.WriteOnlyFields, "license_key")
	assert.Contains(t, hints.WriteOnlyMD5Fields, "license_key")
	assert.Contains(t, hints.AdditionalSchemaAttributes, "license_key_md5")
}
