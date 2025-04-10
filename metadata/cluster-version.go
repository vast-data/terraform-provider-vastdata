package metadata

import (
	version "github.com/hashicorp/go-version"
	"strings"
)

var clusterVersion, _ = version.NewVersion("0.0.0")

// SanitizeVersion truncates all segments of Cluster Version above core (x.y.z)
func SanitizeVersion(version string) (string, bool) {
	segments := strings.Split(version, ".")
	truncated := len(segments) > 3
	return strings.Join(segments[:3], "."), truncated
}

func UpdateClusterVersion(v string) error {
	newVersion, err := version.NewVersion(v)
	if err != nil {
		return err
	}
	//We only work with core version//
	clusterVersion = newVersion.Core()
	return nil
}

func GetClusterVersion() version.Version {
	return *clusterVersion
}

func ClusterVersionString() string {
	cv := GetClusterVersion()
	return cv.String()

}
