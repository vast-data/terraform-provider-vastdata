package metadata

import (
	version "github.com/hashicorp/go-version"
)

var cluster_version, _ = version.NewVersion("0.0.0")

func UpdateClusterVersion(v string) error {
	_cluster_version, err := version.NewVersion(v)
	if err != nil {
		return err
	}
	//We only work with core version//
	cluster_version = _cluster_version.Core()
	return nil
}

func GetClusterVersion() version.Version {
	return *cluster_version
}

func ClusterVersionString() string {
	cv := GetClusterVersion()
	return cv.String()

}
