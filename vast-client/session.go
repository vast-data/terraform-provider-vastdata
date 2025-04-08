package vast_client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type VMSSession struct {
	config         *RestClientConfig
	client         *http.Client
	mu             sync.Mutex
	auth           Authenticator
	clusterVersion *string
}

func NewSession(ctx context.Context, config *RestClientConfig) *VMSSession {
	//Create a new session object
	return &VMSSession{
		config: config,
		client: nil,
		auth:   CreateAuthenticator(ctx, config),
	}
}

func (s *VMSSession) Start() error {
	config := s.GetConfig()
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: !config.SslVerify}
	customTransport.MaxConnsPerHost = 10
	customTransport.IdleConnTimeout = time.Duration(time.Second * 30)
	client := &http.Client{Transport: customTransport}
	s.client = client
	return nil
}

func setupHeaders(ctx context.Context, s *VMSSession, r *http.Request, headers map[string]string) error {
	if err := s.auth.SetAuthHeader(ctx, s, &r.Header); err != nil {
		return err
	}
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-type", "application/json")
	for k, v := range headers {
		r.Header.Add(k, v)
	}
	r.Header.Set("User-Agent", GetUserAgent())
	return nil

}

func buildUrl(s *VMSSession, path, query string) url.URL {
	_url := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s:%v", s.config.Host, s.config.Port),
		Path:   path,
	}
	if query != "" {
		_url.RawQuery = query
	}
	return _url
}

/*
This function validate that the response is OK by validating that the error is nill and that the
Exist code of the response is part of the allwed list
*/
func validateResponse(response *http.Response, err error, allowed ...int) (*http.Response, error) {
	if err != nil {
		return response, err
	}
	if response == nil {
		return response, errors.New("Nil response was provided")
	}
	for _, i := range allowed {
		if response.StatusCode == i {
			return response, err
		}
	}
	return response, errors.New(fmt.Sprintf("Response Status code is %d , which is not allowed", response.StatusCode))
}

/*Define basic HTTP methods to be used with the session*/
func (s *VMSSession) Get(ctx context.Context, path string, query string, headers map[string]string) (response *http.Response, err error) {

	_url := buildUrl(s, path, query)
	url := _url.String()
	req, err := http.NewRequest("GET", url, nil)
	e := setupHeaders(ctx, s, req, headers)
	if e != nil {
		return nil, e
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling Get with URL %s", url))
	response, response_error := s.client.Do(req)
	return validateResponse(response, response_error, 200, 201, 204)
}

func (s *VMSSession) Post(ctx context.Context, path string, body io.Reader, headers map[string]string) (response *http.Response, err error) {

	_url := buildUrl(s, path, "")
	url := _url.String()
	req, err := http.NewRequest("POST", url, body)

	e := setupHeaders(ctx, s, req, headers)
	if e != nil {
		return nil, e
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling Post with URL %s", url))
	response, response_error := s.client.Do(req)
	return validateResponse(response, response_error, 200, 201, 204)
}

func (s *VMSSession) Put(ctx context.Context, path string, body io.Reader, headers map[string]string) (response *http.Response, err error) {
	_url := buildUrl(s, path, "")
	url := _url.String()
	req, err := http.NewRequest("PUT", url, body)
	e := setupHeaders(ctx, s, req, headers)
	if e != nil {
		return nil, e
	}
	response, response_error := s.client.Do(req)
	tflog.Debug(ctx, fmt.Sprintf("Calling Put with URL %s", url))
	return validateResponse(response, response_error, 200, 201, 204)
}

func (s *VMSSession) Patch(ctx context.Context, path string, body io.Reader, headers map[string]string) (response *http.Response, err error) {
	_url := buildUrl(s, path, "")
	url := fmt.Sprintf("%s/", _url.String())
	req, err := http.NewRequest("PATCH", url, body)
	e := setupHeaders(ctx, s, req, headers)
	if e != nil {
		return nil, e
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling Patch with URL %s", url))
	response, response_error := s.client.Do(req)
	return validateResponse(response, response_error, 200, 201, 204)
}

func (s *VMSSession) Delete(ctx context.Context, path, query string, body io.Reader, headers map[string]string) (response *http.Response, err error) {
	_url := buildUrl(s, path, query)
	url := fmt.Sprintf("%s/", _url.String())
	req, err := http.NewRequest("DELETE", url, body)
	e := setupHeaders(ctx, s, req, headers)
	if e != nil {
		return nil, e
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling Delete with URL %s", url))
	response, response_error := s.client.Do(req)
	return validateResponse(response, response_error, 200, 201, 204)
}

func (s *VMSSession) ClusterVersion(ctx context.Context) (version string, response *http.Response, err error) {
	//Cache the cluster version//
	type ClusterVersion struct {
		ClusterVersion string `json:"sw_version"`
	}

	if s.clusterVersion != nil {
		return *s.clusterVersion, nil, nil
	}
	var b []byte
	response, response_error := s.Get(ctx, "/api/clusters/", "", map[string]string{})
	response, response_error = validateResponse(response, response_error, 200, 201, 204)
	if response_error != nil {
		return "", response, response_error
	}
	clustersVersions := []ClusterVersion{}
	b, response_error = io.ReadAll(response.Body)
	if response_error != nil {
		return "", response, errors.New("Failed to read http response body")
	}
	response_error = json.Unmarshal(b, &clustersVersions)
	if response_error != nil {
		return "", response, errors.New("Falied to extract list of clusters version from server response")
	}
	if len(clustersVersions) <= 0 {
		return "", response, errors.New("Could not found clusters to obtain version from")
	}
	//For now we as assume that there is only one cluster so we always grab the first cluster in the list.
	clusterVersion := clustersVersions[0]
	if clusterVersion.ClusterVersion == "" {
		return "", response, errors.New("Empty Cluster version returned")
	}
	s.clusterVersion = &clusterVersion.ClusterVersion
	return clusterVersion.ClusterVersion, response, nil
}

func (s *VMSSession) GetConfig() *RestClientConfig {
	return s.config
}
func (s *VMSSession) Lock() { s.mu.Lock() }

func (s *VMSSession) Unlock() { s.mu.Unlock() }
