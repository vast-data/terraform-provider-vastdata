/*
   This code will implemant a Generic HTTP interfaces to be used to perform CRUD operations
   aginst a vastdata cluster.
*/

package vast_client

import (
	"context"
	"io"
	"net/http"
	"sync"
)

type RestClient interface {
	GetConfig() *RestClientConfig
	Post(ctx context.Context, path string, query string, body io.Reader, headers map[string]string) (response *http.Response, err error)
	Get(ctx context.Context, path, query string, headers map[string]string) (response *http.Response, err error)
	Put(ctx context.Context, path string, body io.Reader, headers map[string]string) (response *http.Response, err error)
	Patch(ctx context.Context, path string, query string, body io.Reader, headers map[string]string) (response *http.Response, err error)
	Delete(ctx context.Context, path, query string, body io.Reader, headers map[string]string) (response *http.Response, err error)
}

type VersionGetter interface {
	ClusterVersion(ctx context.Context) (string, *http.Response, error)
}

type Lockable interface {
	sync.Locker
}

type Session interface {
	RestClient
	VersionGetter
	Lockable
	Start() error
}

type RestClientConfig struct {
	Host      string
	Port      uint64
	Username  string
	Password  string
	ApiToken  string
	SslVerify bool
}
