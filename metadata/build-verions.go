package metadata

import (
	version "github.com/hashicorp/go-version"
)

var build_version, _ = version.NewVersion("5.3.0")
var min_version, _ = version.NewVersion("5.0.0")

const (
	CLUSTER_VERSION_EQUALS int = 0
	CLUSTER_VERSION_LOWER  int = 1
	CLUSTER_VERSION_GRATER int = -1
)

func GetBuildVersion() version.Version {
	return *build_version
}

func ClusterVersionCompare() int {
	return build_version.Compare(cluster_version)

}

func GetMinVersion() version.Version {
	return *min_version
}

func BuildVersionString() string {
	bv := GetBuildVersion()
	return bv.String()

}

func IsLowerThanMinVersion() bool {
	return cluster_version.LessThan(min_version)

}
