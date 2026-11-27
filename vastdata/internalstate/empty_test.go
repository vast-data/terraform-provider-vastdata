// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func secretState(t *testing.T, secret attr.Value) *TFState {
	t.Helper()
	sch := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"name":   rschema.StringAttribute{Required: true},
		"secret": rschema.StringAttribute{Optional: true, WriteOnly: true},
	}}
	st, err := NewTFStateFromSchema(map[string]attr.Value{"name": types.StringValue("n"), "secret": secret}, sch,
		&TFStateHints{WriteOnlyFields: []string{"secret"}})
	require.NoError(t, err)
	return st
}

func TestMarkEmpty(t *testing.T) {
	st := secretState(t, types.StringNull())
	assert.False(t, st.IsEmpty())

	assert.Same(t, st, MarkEmpty(st))
	assert.True(t, st.IsEmpty())
	assert.Equal(t, "n", st.String("name"), "an empty state still reads like any other")
	assert.True(t, st.Copy().IsEmpty(), "a copy stays empty")

	var none *TFState
	assert.True(t, none.IsEmpty())
}

func TestEmptyState(t *testing.T) {
	st := EmptyState()
	assert.True(t, st.IsEmpty())
	assert.False(t, st.IsEnabled())
	assert.Empty(t, st.RawValues())
	assert.Nil(t, st.StateHints())
}

func TestOverlayWriteOnly(t *testing.T) {
	dst := secretState(t, types.StringNull())
	OverlayWriteOnly(dst, secretState(t, types.StringValue("s3cr3t")))
	assert.Equal(t, "s3cr3t", dst.String("secret"))

	OverlayWriteOnly(dst, secretState(t, types.StringNull()))
	assert.Equal(t, "s3cr3t", dst.String("secret"), "a value config leaves unset is kept")

	assert.NotPanics(t, func() { OverlayWriteOnly(dst, nil) })
}
