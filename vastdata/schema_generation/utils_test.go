// Copyright (c) HashiCorp, Inc.

package schema_generation

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
)

func TestAddSchemaEntries_PrimitiveTypes(t *testing.T) {
	props := map[string]*openapi3.SchemaRef{
		"foo": {Value: &openapi3.Schema{Type: (*openapi3.Types)(&[]string{openapi3.TypeString}), Description: "desc"}},
		"bar": {Value: &openapi3.Schema{Type: (*openapi3.Types)(&[]string{openapi3.TypeInteger})}},
	}
	target := map[string]*SchemaEntry{}
	addSchemaEntries(props, []string{"foo"}, nil, target, false, true, false, false, false, false)
	require.Contains(t, target, "foo")
	require.True(t, target["foo"].Required)
	require.Equal(t, "desc", target["foo"].Description)
	require.Contains(t, target, "bar")
	require.True(t, target["bar"].Optional)
}

func TestGetSchemaType(t *testing.T) {
	t.Run("nil schema", func(t *testing.T) {
		require.Equal(t, "", getSchemaType(nil))
	})

	t.Run("empty type", func(t *testing.T) {
		s := &openapi3.Schema{}
		require.Equal(t, "", getSchemaType(s))
	})

	t.Run("string type", func(t *testing.T) {
		s := &openapi3.Schema{
			Type: (*openapi3.Types)(&[]string{openapi3.TypeString}),
		}
		require.Equal(t, "string", getSchemaType(s))
	})

	t.Run("array type", func(t *testing.T) {
		s := &openapi3.Schema{
			Type: (*openapi3.Types)(&[]string{openapi3.TypeArray}),
		}
		require.Equal(t, "array", getSchemaType(s))
	})
}

func TestCompareSchemaValues(t *testing.T) {
	t.Run("equal simple strings", func(t *testing.T) {
		a := &openapi3.Schema{Type: (*openapi3.Types)(&[]string{"string"})}
		b := &openapi3.Schema{Type: (*openapi3.Types)(&[]string{"string"})}
		reason, ok := compareSchemaValues(a, b)
		require.True(t, ok)
		require.Equal(t, "", reason)
	})

	t.Run("type mismatch", func(t *testing.T) {
		a := &openapi3.Schema{Type: (*openapi3.Types)(&[]string{"string"})}
		b := &openapi3.Schema{Type: (*openapi3.Types)(&[]string{"array"})}
		reason, ok := compareSchemaValues(a, b)
		require.False(t, ok)
		require.Contains(t, reason, "Type mismatch")
	})

	t.Run("array item mismatch", func(t *testing.T) {
		a := &openapi3.Schema{
			Type:  (*openapi3.Types)(&[]string{"array"}),
			Items: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: (*openapi3.Types)(&[]string{"string"})}},
		}
		b := &openapi3.Schema{
			Type:  (*openapi3.Types)(&[]string{"array"}),
			Items: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: (*openapi3.Types)(&[]string{"integer"})}},
		}
		reason, ok := compareSchemaValues(a, b)
		require.False(t, ok)
		require.Contains(t, reason, "Array item mismatch")
	})

	t.Run("object with different property count", func(t *testing.T) {
		a := &openapi3.Schema{
			Type:       (*openapi3.Types)(&[]string{"object"}),
			Properties: map[string]*openapi3.SchemaRef{"foo": {Value: &openapi3.Schema{Type: (*openapi3.Types)(&[]string{"string"})}}},
		}
		b := &openapi3.Schema{
			Type:       (*openapi3.Types)(&[]string{"object"}),
			Properties: map[string]*openapi3.SchemaRef{},
		}
		reason, ok := compareSchemaValues(a, b)
		require.False(t, ok)
		require.Contains(t, reason, "Object property count mismatch")
	})

	t.Run("object with different nested property types", func(t *testing.T) {
		a := &openapi3.Schema{
			Type: (*openapi3.Types)(&[]string{"object"}),
			Properties: map[string]*openapi3.SchemaRef{
				"id": {Value: &openapi3.Schema{Type: (*openapi3.Types)(&[]string{"string"})}},
			},
		}
		b := &openapi3.Schema{
			Type: (*openapi3.Types)(&[]string{"object"}),
			Properties: map[string]*openapi3.SchemaRef{
				"id": {Value: &openapi3.Schema{Type: (*openapi3.Types)(&[]string{"integer"})}},
			},
		}
		reason, ok := compareSchemaValues(a, b)
		require.False(t, ok)
		require.Contains(t, reason, `Property "id" mismatch`)
	})
}
