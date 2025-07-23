// Copyright (c) HashiCorp, Inc.

package client

import (
	vast_client "github.com/vast-data/go-vast-client"
	"time"
)

func NewRest(
	host string,
	port int64,
	username, password, apiToken string,
	sslVerify bool,
	pluginVer string,
	timeout time.Duration,
) (*vast_client.VMSRest, error) {
	vmsConfig := &vast_client.VMSConfig{
		Host:      host,
		Port:      uint64(port),
		Username:  username,
		Password:  password,
		ApiToken:  apiToken,
		SslVerify: sslVerify,
		UserAgent: getUserAgent(pluginVer),
		Timeout:   &timeout,

		BeforeRequestFn: BeforeRequestFnCallback,
		AfterRequestFn:  AfterRequestFnCallback,
	}

	return vast_client.NewVMSRest(vmsConfig)
}
