/*
   This code will implemant a Generic HTTP client to be used to perform CRUD operations
   aginst a vastdata cluster.

*/

package vast_client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// A generic session interface to be used by the vast client to suppor various types of authentication
type Session interface {
	Start() error
	Post(ctx context.Context, path string, body io.Reader, headers map[string]string) (response *http.Response, err error)
	Get(ctx context.Context, path, query string, headers map[string]string) (response *http.Response, err error)
	Put(ctx context.Context, path string, body io.Reader, headers map[string]string) (response *http.Response, err error)
	Patch(ctx context.Context, path string, body io.Reader, headers map[string]string) (response *http.Response, err error)
	Delete(ctx context.Context, path, query string, headers map[string]string) (response *http.Response, err error)
	ClusterVersion(ctx context.Context) (version string, response *http.Response, err error)
}

// A struct of the HTTP client that will be performing CRUD operations agint the Vast Data cluster.
type VastClient struct {
	//	session Session
	session Session
}

// Perform POST method requests to a VastData Cluster
func (v *VastClient) Post(ctx context.Context, path string, headers map[string]string, in interface{}) (response *http.Response, err error) {
	//Marshal the input in preperation to send as post data
	b, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	rsp, err := v.session.Post(ctx, path, bytes.NewReader(b), headers)
	return rsp, err
}

// Perform PATCH method requests to a VastData Cluster
func (v *VastClient) Patch(ctx context.Context, path string, headers map[string]string, in interface{}) (response *http.Response, err error) {
	//Covert data to json
	b, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	rsp, err := v.session.Patch(ctx, path, bytes.NewReader(b), headers)
	return rsp, err
}

// Perform PUT method requests to a VastData Cluster
func (v *VastClient) Put(ctx context.Context, path string, headers map[string]string, in interface{}) (response *http.Response, err error) {
	//Marshal the input in preperation to send as put data
	b, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	rsp, err := v.session.Put(ctx, path, bytes.NewReader(b), headers)
	return rsp, err
}

// Perform GET method requests to a VastData Cluster
func (v *VastClient) Get(ctx context.Context, path, query string, headers map[string]string) (response *http.Response, err error) {
	//Simply call get
	rsp, err := v.session.Get(ctx, path, query, headers)
	return rsp, err
}

// Perform GET method requests to a VastData Cluster building a query from list of attributes related to a struct (in)
func (v *VastClient) GetByAttributesListFromInterface(ctx context.Context, path string, in interface{}, attributes []string, headers map[string]string) (response *http.Response, err error) {
	//Build Query from attributes of interface in
	query := Query_builder(in, attributes)
	rsp, err := v.session.Get(ctx, path, query, headers)
	return rsp, err
}

// Perform GET method requests to a VastData Cluster building a query from a map[string]string of values
func (v *VastClient) GetByAttributesMap(ctx context.Context, path string, attributes map[string]string, headers map[string]string) (response *http.Response, err error) {
	//Build Query from a map
	u := url.Values{}
	for k, v := range attributes {
		u.Add(k, v)
	}
	rsp, err := v.session.Get(ctx, path, u.Encode(), headers)
	return rsp, err
}

// Perform DELETE method requests to a VastData Cluster
func (v *VastClient) Delete(ctx context.Context, path, query string, headers map[string]string) (response *http.Response, err error) {

	rsp, err := v.session.Delete(ctx, path, query, headers)
	return rsp, err
}

func (v *VastClient) ClusterVersion(ctx context.Context) (version string, response *http.Response, err error) {
	version, rsp, err := v.session.ClusterVersion(ctx)
	return version, rsp, err

}
