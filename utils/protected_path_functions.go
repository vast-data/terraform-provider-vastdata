package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	metadata "github.com/vast-data/terraform-provider-vastdata/metadata"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
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

func ProtectedPathDeleteFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	//First we call delete than we wait up to 10 min till the protected path is deleted
	response, err := DefaultDeleteFunc(ctx, _client, attr, data, headers)
	if err != nil {
		return response, err
	}
	//Now we wait for the protected path deletion
	client := _client.(vast_client.JwtSession)
	attributes, err := getAttributesAsString([]string{"path", "id"}, attr)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("%v%v", (*attributes)["path"], (*attributes)["id"])
	deleted := false

	for i := time.Minute * 10; i > 0; i = (i - 10*time.Second) {
		tflog.Debug(ctx, fmt.Sprintf("Wating for protectedpath: %v to be deleted", path))
		time.Sleep(time.Second * 10)
		r, _ := client.Get(ctx, path, "", map[string]string{})
		switch r.StatusCode {
		case 404:
			tflog.Debug(ctx, fmt.Sprintf("Protected Path %v Deleted", path))
			deleted = true
			break
		default:
			tflog.Debug(ctx, fmt.Sprintf("Protected Path %v is still being deleted", path))

		}
		if deleted {
			break
		}

	}
	return response, nil

}
