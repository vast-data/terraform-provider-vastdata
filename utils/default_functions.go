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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
)

func getAttributeOrDefault(name string, attribute_default *string, attrs map[string]interface{}) *string {
	_attr, attributeExists := attrs[name]
	var value *string = new(string)
	if !attributeExists {
		value = attribute_default
		return value
	}
	attr, isString := _attr.(string)
	if !isString {
		value = attribute_default
		return value
	}
	*value = attr
	return value

}

func getAttributeAsString(name string, attrs map[string]interface{}) (string, error) {
	_attr, attributeExists := attrs[name]
	if !attributeExists {
		return "", errors.New(fmt.Sprintf("Attribute with the name \"%v\" does not exists", name))
	}
	attr, isString := _attr.(string)
	if !isString {
		return "", errors.New(fmt.Sprintf("Attribute with the name \"%v\" is not a string %v", name, _attr))
	}
	return attr, nil
}

func getAttributesAsString(names []string, attrs map[string]interface{}) (*map[string]string, error) {
	var m map[string]string = make(map[string]string)
	for _, name := range names {
		attr, err := getAttributeAsString(name, attrs)
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
	b, marshallingError := json.Marshal(data)
	if marshallingError != nil {
		return nil, marshallingError
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling POST to path \"%v\"", (*attributes)["path"]))
	return client.Post(ctx, (*attributes)["path"], bytes.NewReader(b), map[string]string{})
}

type UpdateFuncType func(context.Context, interface{}, map[string]interface{}, map[string]interface{}, *schema.ResourceData, map[string]string) (*http.Response, error)

func DefaultUpdateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	client := _client.(vast_client.JwtSession)
	attributes, err := getAttributesAsString([]string{"path", "id"}, attr)
	if err != nil {
		return nil, err
	}
	updatePath := fmt.Sprintf("%v/%v", (*attributes)["path"], (*attributes)["id"])
	b, marshallingError := json.Marshal(data)
	if marshallingError != nil {
		return nil, marshallingError
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling PATCH to path \"%v\"", updatePath))
	tflog.Debug(ctx, fmt.Sprintf("Calling PATCH with payload: %v", string(b)))
	return client.Patch(ctx, updatePath, "application/json", bytes.NewReader(b), map[string]string{})
}

type DeleteFuncType func(context.Context, interface{}, map[string]interface{}, map[string]interface{}, map[string]string) (*http.Response, error)

func DefaultDeleteFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	client := _client.(vast_client.JwtSession)
	attributes, err := getAttributesAsString([]string{"path", "id"}, attr)
	if err != nil {
		return nil, err
	}
	deletePath := fmt.Sprintf("%v%v", (*attributes)["path"], (*attributes)["id"])
	query := ""
	_query := getAttributeOrDefault("query", nil, attr)
	if _query != nil {
		query = *_query
	}
	var r *bytes.Reader = nil
	if !reflect.DeepEqual(data, map[string]interface{}{}) {

		b, marshallingError := json.Marshal(data)
		if marshallingError != nil {
			return nil, marshallingError
		}
		r = bytes.NewReader(b)
	}
	tflog.Debug(ctx, fmt.Sprintf("Calling Delete to path \"%v\"", deletePath))
	return client.Delete(ctx, deletePath, query, r, map[string]string{})
}

type GetFuncType func(context.Context, interface{}, map[string]interface{}, *schema.ResourceData, map[string]string) (*http.Response, error)

func DefaultGetFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
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

type IdFuncType func(context.Context, interface{}, interface{}, *schema.ResourceData) error

func DefaultIdFunc(ctx context.Context, _client interface{}, _id interface{}, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%v", _id))
	return nil
}

type ImportFunc func(context.Context, interface{}, map[string]interface{}, *schema.ResourceData, GetFuncType) (*http.Response, error)

func DefaultImportFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, g GetFuncType) (*http.Response, error) {
	tflog.Debug(ctx, fmt.Sprintf("Calling Import Func %v", g))
	return g(ctx, _client, attr, d, map[string]string{})
}
