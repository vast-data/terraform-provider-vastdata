// Copyright (c) HashiCorp, Inc.

package client

import (
	"context"
	"io"
	"net/http"
	"time"

	vast_client "github.com/vast-data/go-vast-client"
)

// Hooks are the request interceptors installed on the REST client.
type Hooks struct {
	// BeforeRequest runs before every request is sent.
	BeforeRequest func(ctx context.Context, r *http.Request, verb, url string, body io.Reader) error
	// AfterRequest runs after every HTTP attempt; response is nil when it failed.
	AfterRequest func(ctx context.Context, response vast_client.Renderable, statusCode int) (vast_client.Renderable, error)
}

// DefaultHooks are the interceptors used by NewRest: they log through the provider's
// logging package, which defaults to tflog.
func DefaultHooks() Hooks {
	return Hooks{BeforeRequest: BeforeRequestFnCallback, AfterRequest: AfterRequestFnCallback}
}

// NewRest creates a REST client with the default (Terraform) interceptors.
func NewRest(
	host string,
	port int64,
	username, password, apiToken, tenant string,
	sslVerify bool,
	pluginVer string,
	timeout time.Duration,
) (*vast_client.VMSRest, error) {
	return NewRestWithHooks(host, port, username, password, apiToken, tenant, sslVerify, pluginVer, timeout, DefaultHooks())
}

// NewRestWithHooks is NewRest with caller-supplied request interceptors.
func NewRestWithHooks(
	host string,
	port int64,
	username, password, apiToken, tenant string,
	sslVerify bool,
	pluginVer string,
	timeout time.Duration,
	hooks Hooks,
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

		BeforeRequestFn: hooks.BeforeRequest,
		AfterRequestFn:  hooks.AfterRequest,
	}

	return vast_client.NewVMSRest(vmsConfig)
}
