package utils

import (
	"fmt"

	metadata "github.com/vast-data/terraform-provider-vastdata/metadata"
)

func GenPath(resource_path string) string {
	ver, exists := metadata.GetClusterConfig("vast_version")
	if !exists {
		panic("Could have not find vast_version")
	}
	return fmt.Sprintf("api/%v/%v/", ver, resource_path)

}
