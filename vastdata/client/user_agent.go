// Copyright (c) HashiCorp, Inc.

package client

import (
	"fmt"
	"runtime"
)

func getUserAgent(pluginVer string) string {
	return fmt.Sprintf(
		"Terraform Provider VASTData ,OS:%s, Arch:%s , Version:%s",
		runtime.GOOS,
		runtime.GOARCH,
		pluginVer,
	)
}
