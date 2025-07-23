// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"testing"

	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestExtractTypesFromSchema_Resource(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"name":  rschema.StringAttribute{Optional: true},
			"count": rschema.Int64Attribute{Required: true},
		},
	}

	typesMap, err := extractTypesFromSchema(schema, SchemaForResource)
	require.NoError(t, err)
	require.Equal(t, types.StringType, typesMap["name"])
	require.Equal(t, types.Int64Type, typesMap["count"])
}

func TestExtractTypesFromSchema_Datasource(t *testing.T) {
	schema := dsschema.Schema{
		Attributes: map[string]dsschema.Attribute{
			"enabled": dsschema.BoolAttribute{Required: true},
		},
	}

	typesMap, err := extractTypesFromSchema(schema, SchemaForDataSource)
	require.NoError(t, err)
	require.Equal(t, types.BoolType, typesMap["enabled"])
}

func TestExtractMetaFromSchema_Resource(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"name": rschema.StringAttribute{Required: true},
			"tags": rschema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}

	metaMap, err := extractMetaFromSchema(schema, SchemaForResource, nil)
	require.NoError(t, err)

	require.True(t, metaMap["name"].Required)
	require.False(t, metaMap["name"].Optional)
	require.False(t, metaMap["name"].Computed)

	require.True(t, metaMap["tags"].Optional)
	require.False(t, metaMap["tags"].Required)
	require.False(t, metaMap["tags"].Computed)
}

func TestExtractMetaFromSchema_NestedAttributes(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"config": rschema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]rschema.Attribute{
					"enabled": rschema.BoolAttribute{Required: true},
				},
			},
		},
	}

	metaMap, err := extractMetaFromSchema(schema, SchemaForResource, nil)
	require.NoError(t, err)

	require.True(t, metaMap["config"].Optional)
	require.True(t, metaMap["config.enabled"].Required)
}

func TestExtractMetaWithHints(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"name": rschema.StringAttribute{Optional: true},
		},
	}
	hints := &TFStateHints{
		SearchableFields: []string{"name"},
	}
	metaMap, err := extractMetaFromSchema(schema, SchemaForResource, hints)
	require.NoError(t, err)
	require.True(t, metaMap["name"].Searchable)
}
