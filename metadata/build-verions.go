package metadata

import (
	version "github.com/hashicorp/go-version"
)

var buildVersion, _ = version.NewVersion("5.2.0")
var minVersion, _ = version.NewVersion("5.0.0")

const (
	CLUSTER_VERSION_EQUALS int = 0
	CLUSTER_VERSION_LOWER  int = 1
	CLUSTER_VERSION_GRATER int = -1
)

func GetBuildVersion() version.Version {
	return *buildVersion
}

func ClusterVersionCompare() int {
	return buildVersion.Compare(clusterVersion)

}

func GetMinVersion() version.Version {
	return *minVersion
}

func BuildVersionString() string {
	bv := GetBuildVersion()
	return bv.String()

}

func IsLowerThanMinVersion() bool {
	return clusterVersion.LessThan(minVersion)

}
