// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResourceManagerFromState(t *testing.T) {
	fresh, err := NewResourceManagerByName(context.Background(), "view")
	require.NoError(t, err)
	st := fresh.TfState()
	st.SetOrAdd("path", "/a/")

	m, err := NewResourceManagerFromState("view", st)
	require.NoError(t, err)
	assert.IsType(t, fresh, m, "the manager of the named resource, with its hooks")
	assert.Equal(t, "/a/", m.TfState().String("path"))

	m.TfState().SetOrAdd("path", "/b/")
	assert.Equal(t, "/a/", st.String("path"), "the manager holds a copy")

	_, err = NewResourceManagerFromState("no_such_resource", st)
	assert.Error(t, err)
	_, err = NewResourceManagerFromState("view", nil)
	assert.Error(t, err)
}
