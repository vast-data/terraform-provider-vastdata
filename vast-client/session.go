package vast_client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Time duration is set to 10 min after this we refresh the token
var TOKEN_REFRESH_TIME_IN_SECONDS time.Duration = time.Duration(time.Minute * 10)

type jwt_token struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
	Created time.Time
}

type JwtSession struct {
	token                      *jwt_token
	valid                      bool
	server, username, password string
	port                       uint64
	no_ssl_verify              bool
	client                     *http.Client
	mu                         sync.Mutex
	clusterVersion             *string
}

func (s *JwtSession) IsValid() bool {
	return s.valid
}

func parseToken(rsp *http.Response) (*jwt_token, error) {

	var tokens jwt_token
	out, e := io.ReadAll(rsp.Body)
	if e != nil {
		return nil, e
	}
	e = json.Unmarshal(out, &tokens)
	if e != nil {
		return nil, e
	}
	tokens.Created = time.Now()
	return &tokens, nil
}

func (s *JwtSession) getJwtAccessToken() (string, error) {
	/*This function will return the access token
	/if the time duration between the token creation time & current time greated or equal to
	TOKEN_REFRESH_TIME_IN_SECONDS we try refresh it
	*/
	s.mu.Lock()
	if !s.IsValid() {
		return "", errors.New("Session has not been in initialized")
	}
	defer s.mu.Unlock()
	if time.Now().Sub(s.token.Created) >= TOKEN_REFRESH_TIME_IN_SECONDS {
		//Refresh the token
		path := url.URL{
			Scheme: "https",
			Host:   s.server,
			Path:   "api/token/refresh/",
		}
		b, e := json.Marshal(map[string]string{"refresh": s.token.Refresh})
		if e != nil {
			return "", nil
		}
		reader := bytes.NewReader(b)
		rsp, e := s.client.Post(path.String(), "application/json", reader)
		if e != nil {
			return "", e
		}
		t, e := parseToken(rsp)
		if e != nil {
			return "", e
		}
		s.token = t

	}
	return s.token.Access, nil
}

func NewJwtSession(server, username, password string, port uint64, no_ssl_verify bool) JwtSession {
	//Create a new session object
	return JwtSession{token: nil,
		valid:         false,
		server:        server,
		username:      username,
		password:      password,
		port:          port,
		no_ssl_verify: no_ssl_verify,
		client:        nil,
	}
}

func (s *JwtSession) Start() error {
	//Start a new session , try to obtain new access & refresh tokens
	user_pass := map[string]string{"username": s.username, "password": s.password}
	server := s.server + ":" + strconv.FormatUint(s.port, 10)

	body, _ := json.Marshal(user_pass)
	request := bytes.NewBuffer(body)

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: s.no_ssl_verify}
	customTransport.MaxConnsPerHost = 10
	customTransport.IdleConnTimeout = time.Duration(time.Second * 30)
	client := &http.Client{Transport: customTransport}
	s.client = client

	/*Generate URL to obtain token keys*/
	path := url.URL{
		Scheme: "https",
		Host:   server,
		Path:   "api/token/",
	}

	/* Obtain access & session tokens from VMS */
	resp, err := s.client.Post(path.String(), "application/json", request)
	if err != nil {
		return err
	}
	t, e := parseToken(resp)
	if e != nil {
		return e
	}
	s.token = t
	s.valid = true
	return nil

}

func setupHeaders(s *JwtSession, r *http.Request, headers map[string]string) error {
	t, e := s.getJwtAccessToken()
	if e != nil {
		return e
	}
	r.Header.Add("authorization", "Bearer "+t)
	r.Header.Add("accept", "application/json")
	r.Header.Add("content-type", "application/json")
	for k, v := range headers {
		r.Header.Add(k, v)
	}
	return nil

}

func buildUrl(s *JwtSession, path, query string) url.URL {
	_url := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s:%v", s.server, s.port),
		Path:   path,
	}
	if query != "" {
		_url.RawQuery = query
	}
	return _url
}

/*
This dunction validate that the response is OK by validating that the error is nill and that the
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
func (s *JwtSession) Get(ctx context.Context, path string, query string, headers map[string]string) (response *http.Response, err error) {

	_url := buildUrl(s, path, query)
	url := _url.String()
	req, err := http.NewRequest("GET", url, nil)
	e := setupHeaders(s, req, headers)
	if e != nil {
		return nil, e
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling Get with URL %s", url))
	response, response_error := s.client.Do(req)
	return validateResponse(response, response_error, 200, 201, 204)
}

func (s *JwtSession) Post(ctx context.Context, path string, body io.Reader, headers map[string]string) (response *http.Response, err error) {

	_url := buildUrl(s, path, "")
	url := _url.String()
	req, err := http.NewRequest("POST", url, body)

	e := setupHeaders(s, req, headers)
	if e != nil {
		return nil, e
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling Post with URL %s", url))
	response, response_error := s.client.Do(req)
	return validateResponse(response, response_error, 200, 201, 204)
}

func (s *JwtSession) Put(ctx context.Context, path string, body io.Reader, headers map[string]string) (response *http.Response, err error) {
	_url := buildUrl(s, path, "")
	url := _url.String()
	req, err := http.NewRequest("PUT", url, body)
	e := setupHeaders(s, req, headers)
	if e != nil {
		return nil, e
	}
	response, response_error := s.client.Do(req)
	tflog.Debug(ctx, fmt.Sprintf("Calling Put with URL %s", url))
	return validateResponse(response, response_error, 200, 201, 204)
}

func (s *JwtSession) Patch(ctx context.Context, path, contentType string, body io.Reader, headers map[string]string) (response *http.Response, err error) {
	_url := buildUrl(s, path, "")
	url := fmt.Sprintf("%s/", _url.String())
	req, err := http.NewRequest("PATCH", url, body)
	e := setupHeaders(s, req, headers)
	if e != nil {
		return nil, e
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling Patch with URL %s", url))
	response, response_error := s.client.Do(req)
	return validateResponse(response, response_error, 200, 201, 204)
}

func (s *JwtSession) Delete(ctx context.Context, path, query string, body io.Reader, headers map[string]string) (response *http.Response, err error) {
	_url := buildUrl(s, path, query)
	url := fmt.Sprintf("%s/", _url.String())
	req, err := http.NewRequest("DELETE", url, body)
	e := setupHeaders(s, req, headers)
	if e != nil {
		return nil, e
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling Delete with URL %s", url))
	response, response_error := s.client.Do(req)
	return validateResponse(response, response_error, 200, 201, 204)
}

func (s *JwtSession) ClusterVersion(ctx context.Context) (version string, response *http.Response, err error) {
	//Cache the cluster version//
	type ClusterVersion struct {
		ClusterVersion string `json:"sw_version"`
	}

	if s.clusterVersion != nil {
		return *s.clusterVersion, nil, nil
	}
	var b []byte
	response, response_error := s.Get(ctx, "/api/clusters/1/", "", map[string]string{})
	response, response_error = validateResponse(response, response_error, 200, 201, 204)
	if response_error != nil {
		return "", response, response_error
	}
	clusterVersion := ClusterVersion{}
	b, response_error = io.ReadAll(response.Body)
	if response_error != nil {
		return "", response, errors.New("Failed to read http response body")
	}
	response_error = json.Unmarshal(b, &clusterVersion)
	if response_error != nil {
		return "", response, errors.New("Falied to extract sw_version from server response")
	}
	s.clusterVersion = &clusterVersion.ClusterVersion
	return clusterVersion.ClusterVersion, response, nil

}
