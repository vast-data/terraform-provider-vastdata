// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

func TestResolveKubernetesClusterVRN(t *testing.T) {
	api := &deDeleteProbeAPI{listed: RecordSet{
		{"name": "alpha", "vrn": "vast:dataengine:kubernetes-clusters:a"},
		{"name": "beta", "vrn": "vast:dataengine:kubernetes-clusters:b"},
	}}

	t.Run("passthrough_vrn", func(t *testing.T) {
		vrn, err := resolveKubernetesClusterVRN(
			context.Background(), api, nil,
			"vast:dataengine:kubernetes-clusters:already",
		)
		require.NoError(t, err)
		require.Equal(t, "vast:dataengine:kubernetes-clusters:already", vrn)
	})

	t.Run("resolve_name", func(t *testing.T) {
		vrn, err := resolveKubernetesClusterVRN(context.Background(), api, nil, "beta")
		require.NoError(t, err)
		require.Equal(t, "vast:dataengine:kubernetes-clusters:b", vrn)
	})

	t.Run("not_found", func(t *testing.T) {
		_, err := resolveKubernetesClusterVRN(context.Background(), api, nil, "missing")
		require.Error(t, err)
	})
}

func TestResolveContainerRegistryVRN(t *testing.T) {
	api := &deDeleteProbeAPI{listed: RecordSet{
		{"name": "reg-a", "vrn": "vast:dataengine:container-registries:a"},
		{"name": "reg-b", "vrn": "vast:dataengine:container-registries:b"},
	}}

	t.Run("passthrough_vrn", func(t *testing.T) {
		vrn, err := resolveContainerRegistryVRN(
			context.Background(), api, nil,
			"vast:dataengine:container-registries:already",
		)
		require.NoError(t, err)
		require.Equal(t, "vast:dataengine:container-registries:already", vrn)
	})

	t.Run("resolve_name", func(t *testing.T) {
		vrn, err := resolveContainerRegistryVRN(context.Background(), api, nil, "reg-b")
		require.NoError(t, err)
		require.Equal(t, "vast:dataengine:container-registries:b", vrn)
	})

	t.Run("not_found", func(t *testing.T) {
		_, err := resolveContainerRegistryVRN(context.Background(), api, nil, "missing")
		require.Error(t, err)
	})
}

func TestValidateContainerRegistryNameOrVrnInState(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		deContainerRegistryVrnAttr:  deContainerRegistryVrnSchemaAttribute(),
		deContainerRegistryNameAttr: deContainerRegistryNameSchemaAttribute(),
	}}
	newState := func(vrn, name attr.Value) is.ResourceState {
		return is.NewResourceState(map[string]attr.Value{
			deContainerRegistryVrnAttr:  vrn,
			deContainerRegistryNameAttr: name,
		}, schema, &is.TFStateHints{})
	}

	require.NoError(t, validateContainerRegistryNameOrVrnInState(newState(
		types.StringValue("vast:dataengine:container-registries:x"), types.StringNull(),
	)))
	require.NoError(t, validateContainerRegistryNameOrVrnInState(newState(
		types.StringNull(), types.StringValue("my-reg"),
	)))
	require.Error(t, validateContainerRegistryNameOrVrnInState(newState(
		types.StringValue("vast:dataengine:container-registries:x"), types.StringValue("my-reg"),
	)))
	require.Error(t, validateContainerRegistryNameOrVrnInState(newState(
		types.StringNull(), types.StringNull(),
	)))
}

func TestValidateKubernetesClusterNameOrVrnInState(t *testing.T) {
	schema := rschema.Schema{Attributes: map[string]rschema.Attribute{
		deKubernetesClusterVrnAttr:  deKubernetesClusterVrnSchemaAttribute(),
		deKubernetesClusterNameAttr: deKubernetesClusterNameSchemaAttribute(),
	}}
	newState := func(vrn, name attr.Value) is.ResourceState {
		return is.NewResourceState(map[string]attr.Value{
			deKubernetesClusterVrnAttr:  vrn,
			deKubernetesClusterNameAttr: name,
		}, schema, &is.TFStateHints{})
	}

	require.NoError(t, validateKubernetesClusterNameOrVrnInState(newState(
		types.StringValue("vast:dataengine:kubernetes-clusters:x"), types.StringNull(),
	)))
	require.NoError(t, validateKubernetesClusterNameOrVrnInState(newState(
		types.StringNull(), types.StringValue("my-cluster"),
	)))
	require.Error(t, validateKubernetesClusterNameOrVrnInState(newState(
		types.StringValue("vast:dataengine:kubernetes-clusters:x"), types.StringValue("my-cluster"),
	)))
	require.Error(t, validateKubernetesClusterNameOrVrnInState(newState(
		types.StringNull(), types.StringNull(),
	)))
}

func TestValidateKubernetesClusterNameOrVrnAttrs(t *testing.T) {
	t.Run("vrn_ok", func(t *testing.T) {
		require.NoError(t, validateKubernetesClusterNameOrVrnAttrs("ref", map[string]attr.Value{
			deKubernetesClusterVrnAttr:  types.StringValue("vast:dataengine:kubernetes-clusters:x"),
			deKubernetesClusterNameAttr: types.StringNull(),
		}))
	})
	t.Run("name_ok", func(t *testing.T) {
		require.NoError(t, validateKubernetesClusterNameOrVrnAttrs("ref", map[string]attr.Value{
			deKubernetesClusterVrnAttr:  types.StringNull(),
			deKubernetesClusterNameAttr: types.StringValue("my-cluster"),
		}))
	})
	t.Run("both_rejected", func(t *testing.T) {
		require.Error(t, validateKubernetesClusterNameOrVrnAttrs("ref", map[string]attr.Value{
			deKubernetesClusterVrnAttr:  types.StringValue("vast:dataengine:kubernetes-clusters:x"),
			deKubernetesClusterNameAttr: types.StringValue("my-cluster"),
		}))
	})
	t.Run("neither_rejected", func(t *testing.T) {
		require.Error(t, validateKubernetesClusterNameOrVrnAttrs("ref", map[string]attr.Value{
			deKubernetesClusterVrnAttr:  types.StringNull(),
			deKubernetesClusterNameAttr: types.StringNull(),
		}))
	})
	t.Run("unknown_skipped", func(t *testing.T) {
		require.NoError(t, validateKubernetesClusterNameOrVrnAttrs("ref", map[string]attr.Value{
			deKubernetesClusterVrnAttr:  types.StringUnknown(),
			deKubernetesClusterNameAttr: types.StringNull(),
		}))
	})
}
