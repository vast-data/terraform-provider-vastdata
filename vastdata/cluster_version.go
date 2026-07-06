// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"sync"

	version "github.com/hashicorp/go-version"
)

// VastVersion550 is the minimum cluster version that introduced sub-resources
// such as s3_true_ip_config (cluster) and s3cors_configuration (view).
var VastVersion550 = version.Must(version.NewVersion("5.5.0"))

// clusterVersionStore is a process-wide cache of VAST cluster versions, keyed by
// the cluster's host address. It is populated lazily on the first sub-resource
// fetch per cluster and reused for the lifetime of the Terraform process.
var clusterVersionStore struct {
	mu    sync.Mutex
	cache map[string]*version.Version
}

func init() {
	clusterVersionStore.cache = make(map[string]*version.Version)
}

// GetCachedClusterVersion returns the VAST cluster version for the given REST
// client, retrieving it from the API exactly once per cluster host address and
// caching the result for all subsequent calls within the same process.
func GetCachedClusterVersion(ctx context.Context, rest *VMSRest) (*version.Version, error) {
	host := rest.Session.GetConfig().Host

	clusterVersionStore.mu.Lock()
	if v, ok := clusterVersionStore.cache[host]; ok {
		clusterVersionStore.mu.Unlock()
		return v, nil
	}
	clusterVersionStore.mu.Unlock()

	v, err := rest.Versions.GetVersionWithContext(ctx)
	if err != nil {
		return nil, err
	}

	clusterVersionStore.mu.Lock()
	clusterVersionStore.cache[host] = v
	clusterVersionStore.mu.Unlock()
	return v, nil
}
