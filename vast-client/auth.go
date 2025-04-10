package vast_client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// TOKEN_REFRESH_TIME_IN_SECONDS Time duration is set to 10 min after this we refresh the token
const TOKEN_REFRESH_TIME_IN_SECONDS time.Duration = time.Duration(time.Minute * 10)

type Authenticator interface {
	Authorize(ctx context.Context, ses Session) error
	SetAuthHeader(ctx context.Context, ses Session, headers *http.Header) error
}

func CreateAuthenticator(ctx context.Context, config *RestClientConfig) Authenticator {
	// Check if username and password are provided
	if config.Username != "" && config.Password != "" {
		// Return a new JWTAuthenticator
		tflog.Debug(ctx, "VMS session is using username/password authentication (Bearer token will be acquired).")
		return &JWTAuthenticator{
			Username: config.Username,
			Password: config.Password,
			Token:    nil, // Initially no token
		}
	}
	// If apiToken is provided, return a new ApiRTokenAuthenticator
	if config.ApiToken != "" {
		tflog.Debug(ctx, "VMS session is using API token authentication.")
		return &ApiRTokenAuthenticator{
			Token: config.ApiToken,
		}
	}
	// If neither are provided, panic with an error message
	panic("CreateAuthenticator: neither username/password nor apiToken are provided")
}

type jwtToken struct {
	Access    string `json:"access"`
	Refresh   string `json:"refresh"`
	CreatedAt time.Time
}

type JWTAuthenticator struct {
	Username    string
	Password    string
	Token       *jwtToken
	initialized bool
}

func parseToken(rsp *http.Response) (*jwtToken, error) {
	var tokens jwtToken
	out, e := io.ReadAll(rsp.Body)
	if e != nil {
		return nil, e
	}
	e = json.Unmarshal(out, &tokens)
	if e != nil {
		return nil, e
	}
	tokens.CreatedAt = time.Now()
	return &tokens, nil
}

func (auth *JWTAuthenticator) refreshToken(ctx context.Context, client *http.Client, config RestClientConfig) (*http.Response, error) {
	tflog.Debug(ctx, "Refreshing Bearer token using refresh token.")
	var resp *http.Response
	path := url.URL{
		Scheme: "https",
		Host:   config.Host,
		Path:   "api/token/refresh/",
	}
	body, err := json.Marshal(map[string]string{"refresh": auth.Token.Refresh})
	if err != nil {
		return nil, err
	}
	resp, err = client.Post(path.String(), "application/json", bytes.NewBuffer(body))
	if err != nil {
		tflog.Debug(ctx, "Failed to refresh token: "+err.Error())
		return nil, err
	}
	return resp, nil
}

func (auth *JWTAuthenticator) acquireToken(ctx context.Context, client *http.Client, config RestClientConfig) (*http.Response, error) {
	// obtain new access & refresh tokens
	tflog.Debug(ctx, "Acquiring new Bearer token using username and password.")
	var resp *http.Response
	userPass := map[string]string{"username": config.Username, "password": config.Password}
	server := config.Host + ":" + strconv.FormatUint(config.Port, 10)
	body, err := json.Marshal(userPass)
	if err != nil {
		return nil, err
	}
	// Generate URL to obtain token keys
	path := url.URL{
		Scheme: "https",
		Host:   server,
		Path:   "api/token/",
	}
	resp, err = client.Post(path.String(), "application/json", bytes.NewBuffer(body))
	if err != nil {
		tflog.Debug(ctx, "Failed to acquire token: "+err.Error())
		return nil, err
	}
	return resp, nil
}

func (auth *JWTAuthenticator) Authorize(ctx context.Context, ses Session) error {
	ses.Lock()
	defer ses.Unlock()
	var (
		resp *http.Response
		err  error
	)
	config := ses.GetConfig()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.SslVerify},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	if auth.initialized {
		tokenExpired := time.Now().Sub(auth.Token.CreatedAt) >= TOKEN_REFRESH_TIME_IN_SECONDS
		if !tokenExpired {
			return nil
		}
		resp, err = auth.refreshToken(ctx, client, *config)
	} else {
		resp, err = auth.acquireToken(ctx, client, *config)
		auth.initialized = true
	}
	if _, err = validateResponse(resp, err, 200, 201); err != nil {
		return err
	}
	// Read response
	token, err := parseToken(resp)
	if err != nil {
		return err
	}
	auth.Token = token
	return nil
}

func (auth *JWTAuthenticator) SetAuthHeader(ctx context.Context, ses Session, headers *http.Header) error {
	if err := auth.Authorize(ctx, ses); err != nil {
		return err
	}
	headers.Add("Authorization", "Bearer "+auth.Token.Access)
	return nil
}

type ApiRTokenAuthenticator struct {
	Token string
}

func (auth *ApiRTokenAuthenticator) Authorize(ctx context.Context, ses Session) error {
	if auth.Token == "" {
		auth.Token = ses.GetConfig().ApiToken
	}
	return nil
}

func (auth *ApiRTokenAuthenticator) SetAuthHeader(ctx context.Context, ses Session, headers *http.Header) error {
	if err := auth.Authorize(ctx, ses); err != nil {
		return err
	}
	headers.Add("Authorization", "Api-Token "+auth.Token)
	return nil
}
