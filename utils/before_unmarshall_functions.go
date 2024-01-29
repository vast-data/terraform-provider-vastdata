package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type ReadBeforeUnmarshallFunc func(context.Context, *http.Response) ([]byte, error)

func ReadResultField(ctx context.Context, response *http.Response) ([]byte, error) {
	var data []map[string]interface{}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("Http response returned 0 records of schemas")
	}
	m, merr := json.Marshal(data[0])
	if merr != nil {
		return nil, errors.New(fmt.Sprintf("Error occured converting results list to json : %s ", merr))
	}
	tflog.Debug(ctx, fmt.Sprintf("Returning first result field %v", m))
	return m, nil
}

func GenerateTableReadResponse(ctx context.Context, response *http.Response) ([]byte, error) {
	var m []map[string]interface{} = []map[string]interface{}{}
	var out map[string]interface{} = map[string]interface{}{}
	body, err := io.ReadAll(response.Body)
	tflog.Debug(ctx, fmt.Sprintf("Table Body Returned %v", string(body)))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}

	if len(m) < 1 {
		return nil, errors.New("No Fields returned")
	}
	out["database_name"] = m[0]["database_name"]
	out["schema_name"] = m[0]["schema_name"]
	out["schema_identifier"] = fmt.Sprintf("%v/%v", m[0]["database_name"], m[0]["schema_name"])

	out["name"] = m[0]["table_name"]
	fields := []map[string]interface{}{}
	for _, e := range m {
		fields = append(fields, map[string]interface{}{"name": e["name"], "field": e["field"]})
	}
	f, e := json.Marshal(fields)
	if e != nil {
		return nil, e
	}
	out["fields"] = string(f)
	return json.Marshal(out)
}
