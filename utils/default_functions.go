package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
)

func getAttributeOrDefault(name string, attribute_default *string, attrs map[string]interface{}) *string {
	_attr, attr_exists := attrs[name]
	var value *string = new(string)
	if !attr_exists {
		value = attribute_default
		return value
	}
	attr, is_string := _attr.(string)
	if !is_string {
		value = attribute_default
		return value
	}
	*value = attr
	return value

}

func getArrtibuteAsString(name string, attrs map[string]interface{}) (string, error) {
	_attr, attr_exists := attrs[name]
	if !attr_exists {
		return "", errors.New(fmt.Sprintf("Attribute with the name \"%v\" does not exists", name))
	}
	attr, is_string := _attr.(string)
	if !is_string {
		return "", errors.New(fmt.Sprintf("Attribute with the name \"%v\" is not a string %v", name, _attr))
	}
	return attr, nil
}

func getAttributesAsString(names []string, attrs map[string]interface{}) (*map[string]string, error) {
	var m map[string]string = make(map[string]string)
	for _, name := range names {
		attr, err := getArrtibuteAsString(name, attrs)
		if err != nil {
			return nil, err
		}
		m[name] = attr
	}
	return &m, nil
}

type CreateFuncType func(context.Context, interface{}, map[string]interface{}, map[string]interface{}, map[string]string) (*http.Response, error)

func DefaultCreateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	client := _client.(vast_client.JwtSession)
	attributes, err := getAttributesAsString([]string{"path"}, attr)
	if err != nil {
		return nil, err
	}
	b, marshal_error := json.Marshal(data)
	if marshal_error != nil {
		return nil, marshal_error
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling POST to path \"%v\"", (*attributes)["path"]))
	return client.Post(ctx, (*attributes)["path"], bytes.NewReader(b), map[string]string{})
}

type UpdateFuncType func(context.Context, interface{}, map[string]interface{}, map[string]interface{}, map[string]string) (*http.Response, error)

func DefaultUpdateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	client := _client.(vast_client.JwtSession)
	attributes, err := getAttributesAsString([]string{"path", "id"}, attr)
	if err != nil {
		return nil, err
	}
	update_path := fmt.Sprintf("%v/%v", (*attributes)["path"], (*attributes)["id"])
	b, marshal_error := json.Marshal(data)
	if marshal_error != nil {
		return nil, marshal_error
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling PATCH to path \"%v\"", update_path))
	return client.Patch(ctx, update_path, "application/json", bytes.NewReader(b), map[string]string{})
}

type DeleteFuncType func(context.Context, interface{}, map[string]interface{}, map[string]interface{}, map[string]string) (*http.Response, error)

func DefaultDeleteFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	client := _client.(vast_client.JwtSession)
	attributes, err := getAttributesAsString([]string{"path", "id"}, attr)
	if err != nil {
		return nil, err
	}
	delete_path := fmt.Sprintf("%v/%v", (*attributes)["path"], (*attributes)["id"])
	query := ""
	_query := getAttributeOrDefault("query", nil, attr)
	if _query != nil {
		query = *_query
	}
	var r *bytes.Reader = nil
	if !reflect.DeepEqual(data, map[string]interface{}{}) {

		b, marshal_error := json.Marshal(data)
		if marshal_error != nil {
			return nil, marshal_error
		}
		r = bytes.NewReader(b)
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling PATCH to path \"%v\"", delete_path))
	return client.Delete(ctx, delete_path, query, r, map[string]string{})
}

type GetFuncType func(context.Context, interface{}, map[string]interface{}, map[string]string) (*http.Response, error)

func DefaultGetFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, headers map[string]string) (*http.Response, error) {
	client := _client.(vast_client.JwtSession)
	attributes, err := getAttributesAsString([]string{"path"}, attr)
	if err != nil {
		return nil, err
	}
	path := (*attributes)["path"]
	id := getAttributeOrDefault("id", nil, attr)
	if id != nil {
		path = fmt.Sprintf("%v/%v", path, *id)
	}
	query := ""
	_query := getAttributeOrDefault("query", nil, attr)
	if _query != nil {
		query = *_query
	}

	tflog.Debug(ctx, fmt.Sprintf("Calling GET to path \"%v\" , with Query %v", path, query))
	return client.Get(ctx, path, query, headers)
}
