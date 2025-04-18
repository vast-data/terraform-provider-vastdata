package metadata

import (
	"github.com/hashicorp/go-version"
)

func extractVersion(v *version.Version, e error) version.Version {
	if e != nil {
		panic("Broken version")
	}
	return *v
}

type VastVersionStruct struct {
	Ver         version.Version
	VastVersion string
}

func (v *VastVersionStruct) GetVersion() *version.Version {
	return &v.Ver
}

func (v *VastVersionStruct) GetVastVersion() string {
	return v.VastVersion
}

var apiVersions = []VastVersionStruct{
	{Ver: extractVersion(version.NewVersion("4.6.0")), VastVersion: "v2"},
	{Ver: extractVersion(version.NewVersion("4.7.0")), VastVersion: "v3"},
	{Ver: extractVersion(version.NewVersion("5.0.0")), VastVersion: "v4"},
	{Ver: extractVersion(version.NewVersion("5.1.0")), VastVersion: "v5"},
	{Ver: extractVersion(version.NewVersion("5.2.0")), VastVersion: "v5"},
}

func MaxVastVersion() string {
	/*
	   return the latest available supported, vast version
	   This is for situations where we have a cluster version grater than our
	   build version.
	*/
	last := apiVersions[len(apiVersions)-1]
	return last.GetVastVersion()
}

func MinVastVersion() string {
	/*
	   This will return minimal version , however unlike max where we know that at least
	   the latest the build_version is supported , we can never know the version, so we simply return
	   the latest version which in terms of the cluster is the latest capble version it supports.
	*/
	return "latest"
}

func FindVastVersion(ver string) string {
	newVersion, err := version.NewVersion(ver)
	_clusterVersion := extractVersion(newVersion, err)
	for i := range apiVersions {
		c := _clusterVersion.Compare(apiVersions[i].GetVersion())
		if c == 0 {
			return apiVersions[i].GetVastVersion()
		} else if c == -1 {
			if i == 0 {
				//Version is smaller than the minimal version
				return MinVastVersion()
			}
			/*
			   If current version is smaller than this version and the index is not 0
			   than the maxversion is the previous version
			*/
			return apiVersions[i-1].GetVastVersion()
		}
	}
	/*
	   Reaching to this stage means that the cluster version is bigger than out build (last) version
	*/
	return MaxVastVersion()
}
