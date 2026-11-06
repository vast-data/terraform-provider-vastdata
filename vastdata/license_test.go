// Copyright (c) HashiCorp, Inc.
package provider

import (
	"testing"

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
