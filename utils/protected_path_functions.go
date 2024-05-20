package utils

import (
	"context"
	"errors"
	"net/http"

	"github.com/hashicorp/go-version"
	metadata "github.com/vast-data/terraform-provider-vastdata/metadata"
)

func ProtectedPathCreateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	//If cluster version is 5.1 and above we force the capabilities attributes to be set otherwise we get differant results than previous versions
	cluster_version := metadata.GetClusterVersion()
	min_cluster_version, _ := version.NewVersion("5.1.0")
	if cluster_version.GreaterThanOrEqual(min_cluster_version) {
		_, capabilities_exists := data["capabilities"]
		if !capabilities_exists {
			return nil, errors.New("When cluster version is 5.1.0 and above the \"capabilities\" attribute ,ust be given")
		}
	}
	return DefaultCreateFunc(ctx, _client, attr, data, headers)
}
