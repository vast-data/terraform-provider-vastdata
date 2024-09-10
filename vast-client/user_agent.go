package vast_client

import (
	"fmt"
	"runtime"
)

var user_agent_version string

func init() {
	user_agent_version = "UNKOWN"
}

func GetUserAgent() string {
	return fmt.Sprintf("Terraform Provider VastData ,OS:%s, Arch:%s , Version:%s", runtime.GOOS, runtime.GOARCH, user_agent_version)
}
