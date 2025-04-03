package metadata

import (
	version "github.com/hashicorp/go-version"
	"strconv"
	"strings"
)

var clusterVersion, _ = version.NewVersion("0.0.0")

// SanitizeVersion truncates segments of Cluster Version so that each segment can fit within int64.
// This is needed, because hashicorp's go-version package parses each segment into int64.
func SanitizeVersion(version string) (string, bool) {
	segments := strings.Split(version, ".")
	truncated := false
	for i, segment := range segments {
		for {
			if _, err := strconv.ParseInt(segment, 10, 64); err == nil {
				break
			}
			segment = segment[1:]
			truncated = true
		}
		segments[i] = segment
	}
	return strings.Join(segments, "."), truncated
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
