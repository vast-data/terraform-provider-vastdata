package metadata

import (
	"fmt"
	"strings"

	version "github.com/hashicorp/go-version"
)

var cluster_version, _ = version.NewVersion("0.0.0")

func VastVersionToSemVer(v string) (string, error) {
	/*
	   Ever since vast version 5.3.0 , the returned version from the doe not match the SemVer specs and it is now
	   at the format of Major.Minor.Patch.<something>.<something>....
	   To cope with this we take only the first 3 numbers to obtain the Major.Minor.Patch
	*/

	_v := strings.SplitN(v, ".", 4)
	if len(_v) < 3 {
		return v, fmt.Errorf("Cluster reported wrong version: %v ,which does not have Major.Minor.Patch version format", v)
	}
	ver := strings.Join(_v[0:3], ".")
	return ver, nil
}

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
