// Copyright (c) HashiCorp, Inc.

package client

import (
	"time"

	vast_client "github.com/vast-data/go-vast-client"
)

func NewRest(
	host string,
	port int64,
	username, password, apiToken, tenant string,
	sslVerify bool,
	pluginVer string,
	timeout time.Duration,
) (*vast_client.VMSRest, error) {
	vmsConfig := &vast_client.VMSConfig{
		ApiVersion:   "latest",
		Host:         host,
		Port:         uint64(port),
		Username:     username,
		Password:     password,
		ApiToken:     apiToken,
		Tenant:       tenant,
		SslVerify:    sslVerify,
		UserAgent:    getUserAgent(pluginVer),
		Timeout:      &timeout,
		UseBasicAuth: false,
		RespectProxy: true,

		BeforeRequestFn: BeforeRequestFnCallback,
		AfterRequestFn:  AfterRequestFnCallback,
	}

	return vast_client.NewVMSRest(vmsConfig)
}
