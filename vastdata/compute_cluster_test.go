// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

func TestRetryOnExpression_Retries503(t *testing.T) {
	attempts := 0
	expr := &is.RetryExpression{
		StatusCodes:  []int{http.StatusServiceUnavailable},
		Times:        3,
		SleepSeconds: 0,
	}

	result, err := retryOnExpression(context.Background(), expr, "GetSubResources", "compute_cluster/nodes", func() (string, error) {
		attempts++
		if attempts < 3 {
			return "", &ApiError{StatusCode: http.StatusServiceUnavailable, Body: "service_unavailable"}
		}
		return "ok", nil
	})

	require.NoError(t, err)
	require.Equal(t, "ok", result)
	require.Equal(t, 3, attempts)
}

func TestRetryOnExpression_DoesNotRetry400(t *testing.T) {
	attempts := 0
	expr := &is.RetryExpression{
		StatusCodes:  []int{http.StatusServiceUnavailable},
		Times:        3,
		SleepSeconds: 0,
	}

	_, err := retryOnExpression(context.Background(), expr, "GetSubResources", "compute_cluster/nodes", func() (string, error) {
		attempts++
		return "", &ApiError{StatusCode: http.StatusBadRequest, Body: "bad request"}
	})

	require.Error(t, err)
	require.Equal(t, 1, attempts)
}

func TestWithComputeClusterSubResourceRetry_NilExpression(t *testing.T) {
	calls := 0
	got, err := withComputeClusterSubResourceRetry(context.Background(), nil, "nodes", func() (int, error) {
		calls++
		return 42, nil
	})
	require.NoError(t, err)
	require.Equal(t, 42, got)
	require.Equal(t, 1, calls)
}

func TestRetryOnExpression_ExhaustsAttempts(t *testing.T) {
	expr := &is.RetryExpression{
		StatusCodes:  []int{http.StatusServiceUnavailable},
		Times:        2,
		SleepSeconds: 0,
	}

	_, err := retryOnExpression(context.Background(), expr, "GetSubResources", "compute_cluster/nodes", func() (struct{}, error) {
		return struct{}{}, &ApiError{StatusCode: http.StatusServiceUnavailable, Body: "service_unavailable"}
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "all 2 attempts failed")
}

func TestConvertRecordSetToList_RecordSet(t *testing.T) {
	records := RecordSet{
		{"name": "node-1", "status": "Ready"},
		{"name": "node-2", "status": "Ready"},
	}

	got := convertRecordSetToList(records, func(r map[string]any) map[string]any {
		return map[string]any{"name": r["name"], "status": r["status"]}
	})

	require.Len(t, got, 2)
	require.Equal(t, "node-1", got[0].(map[string]any)["name"])
	require.Equal(t, "node-2", got[1].(map[string]any)["name"])
}

func TestConvertRecordSetToList_PaginatedRecord(t *testing.T) {
	records := Record{
		"results": []any{
			map[string]any{"name": "node-1"},
		},
	}

	got := convertRecordSetToList(records, func(r map[string]any) map[string]any {
		return map[string]any{"name": r["name"]}
	})

	require.Len(t, got, 1)
	require.Equal(t, "node-1", got[0].(map[string]any)["name"])
}

func TestConvertRecordSetToList_Nil(t *testing.T) {
	got := convertRecordSetToList(nil, func(r map[string]any) map[string]any { return r })
	require.Equal(t, []any{}, got)
}

func TestMergeComputeClusterRecords(t *testing.T) {
	got := mergeComputeClusterRecords(
		Record{"compute_cluster_id": int64(2), "compute_cluster_name": "tf-test-cluster"},
		Record{"nodes": []any{map[string]any{"name": "node-1"}}},
	)
	require.Equal(t, int64(2), got["compute_cluster_id"])
	require.Equal(t, "tf-test-cluster", got["compute_cluster_name"])
	require.Len(t, got["nodes"], 1)
}

func TestIsComputeClusterAsyncTaskRecord(t *testing.T) {
	require.False(t, isComputeClusterAsyncTaskRecord(Record{"id": int64(2), "name": "tf-test-cluster"}))
	require.True(t, isComputeClusterAsyncTaskRecord(Record{
		"id":             int64(29),
		"@resourceType":  "VTask",
		"name":           "update_compute_cluster",
		"info": map[string]any{
			"kwargs": map[string]any{"compute_cluster_id": float64(2)},
		},
	}))
	require.Equal(t, int64(2), computeClusterIDFromAsyncRecord(Record{
		"id":            int64(31),
		"@resourceType": "VTask",
		"info": map[string]any{
			"kwargs": map[string]any{"compute_cluster_id": float64(2)},
		},
	}))
}

func TestNormalizeRecordForCreateAdopt_Cnodes(t *testing.T) {
	m := &ComputeCluster{}
	got := m.NormalizeRecordForCreateAdopt(Record{
		"cnodes": []any{
			map[string]any{"id": int64(1), "resource_preset": "BALANCED", "ip": "172.21.87.50"},
		},
	})
	cnodes := got["cnodes"].([]any)
	require.Equal(t, map[string]any{"id": int64(1), "resource_preset": "BALANCED"}, cnodes[0])
}
