package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
)

type WithURL struct {
	Url string `json:"url,omitempty"`
}

type ResponseProcessingFunc func(context.Context, *http.Response, *schema.ResourceData) ([]byte, error)

// The default processing func for http responses for read
func DefaultProcessingFunc(ctx context.Context, response *http.Response, d *schema.ResourceData) ([]byte, error) {
	body, err := io.ReadAll(response.Body)
	tflog.Debug(ctx, fmt.Sprintf("HTTP Response body %s", string(body)))
	return body, err
}

/*
A dedicate function to process reponses from /api/s3lifecyclerules/ as response is structured
In a differant way and we need to extract the results list returned
Ex: of response

	{
	    "count": 1,
	    "next": null,
	    "previous": null,
	    "results": [
	        {
	            "id": 1,
	            "guid": "fcfb4523-f10b-5a69-bc37-759cc7e76c74",
	            "name": "bla",
	            "url": "https://172.31.59.73/api/s3lifecyclerules/1/",
	            "title": "bla",
	            "enabled": true,
	            "prefix": "/bla/*",
	            "min_size": null,
	            "max_size": null,
	            "expiration_days": 7,
	            "view_path": "/bla",
	            "view_id": 1,
	            "expiration_date": null,
	            "expired_obj_delete_marker": null,
	            "noncurrent_days": 30,
	            "newer_noncurrent_versions": 2,
	            "abort_mpu_days_after_initiation": 365
	        }
	    ]
	}
*/
func ProcessingResultsListResponse(ctx context.Context, response *http.Response, d *schema.ResourceData) ([]byte, error) {
	m := new(map[string]interface{})
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	err = json.Unmarshal(body, m)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Resonse From Cluster %v", string(body)))
		return []byte{}, err
	}
	results, exists := (*m)["results"]
	if !exists {
		return []byte{}, nil
	}
	return json.Marshal(results)

}

/*
This function will get a []byte response representing a json at the following format.

	[
	   {...
	    url: <some url>
	    ...},
	]
	We than re-query only the url to rebuild the resonse the following way (we assume response to each query is a json.
	[
	  {json response for url1},
	  {json response for url2},
	  ...
	]
*/
func ResponseGetByURL(ctx context.Context, body []byte, client *vast_client.VMSSession) ([]byte, error) {
	var marsheled []byte
	responses := []map[string]interface{}{}
	urls := []WithURL{}
	err := json.Unmarshal(body, &urls)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Error occured Unmarsheling  %s", err.Error()))
		return []byte{}, err
	}
	tflog.Debug(ctx, fmt.Sprintf("URLs when Unmarsheling  %v", urls))
	for _, u := range urls {
		url, err := url.Parse(u.Url)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error occured parsing url %s", err.Error()))
			return []byte{}, err
		}
		tflog.Debug(ctx, fmt.Sprintf("Trying to read data from URL: %v", url.Path))
		response, err := client.Get(ctx, url.Path, "", map[string]string{})
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error occured reading data from URL: %v", url.Path))
			return []byte{}, err
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error occured reading response body from URL: %v", url.Path))
			return []byte{}, err
		}
		m := new(map[string]interface{})
		err = json.Unmarshal(body, m)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error occured Unmarsheling  %s", err.Error()))
			return []byte{}, err
		}
		responses = append(responses, *m)
		tflog.Debug(ctx, fmt.Sprintf("Response recived from URL :%v \n %v", url.Path, string(body)))

	}

	marsheled, err = json.Marshal(responses)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to marshel results  %s", err.Error()))
		tflog.Debug(ctx, fmt.Sprintf("Failed to marshel %v into a json format with the error %s", responses, err.Error()))
		return []byte{}, err
	}
	tflog.Debug(ctx, fmt.Sprintf("Marsheled responses when collecting data from url %s", marsheled))
	return marsheled, nil
}
