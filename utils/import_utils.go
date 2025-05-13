package utils

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ImportInterface interface {
	GetFunc() GetFuncType
	GetDoc() []string
}

/*
An importer that will create an import by obtaining a list of http fields in specifc order
To preserve backward compatibility we allow also to provide GUID as an optional feature
*/
type HttpFieldTuple struct {
	DisplayName string
	FieldName   string
}

type ImportByHttpFields struct {
	http_fields       []HttpFieldTuple
	DisableGuidImport bool
	fields_String     []string
}

func NewImportByHttpFields(disable_guid_import bool, fields []HttpFieldTuple) *ImportByHttpFields {
	i := ImportByHttpFields{DisableGuidImport: disable_guid_import, http_fields: []HttpFieldTuple{}}
	if !disable_guid_import {
		i.http_fields = append(i.http_fields, HttpFieldTuple{DisplayName: "guid", FieldName: "guid"})
	}
	i.setFields(fields)
	i.fields_String = i.genNeededFieldsString()
	return &i
}

func (i *ImportByHttpFields) GetFieldsString() []string {
	return i.fields_String
}

func (i *ImportByHttpFields) setFields(http_fields []HttpFieldTuple) {

	if i.DisableGuidImport && len(http_fields) == 0 {
		panic("When Definning ImportByHttpFields with DisableGuidImport=True http_fields must be provided ")
	}
	for _, t := range http_fields {
		if t.DisplayName == "guid" && i.DisableGuidImport == false {
			panic("\"guid\" is a reserved word when definning fields")
		}
	}
	i.http_fields = append(i.http_fields, http_fields...)

}

func (i *ImportByHttpFields) genNeededFieldsString() []string {
	out := []string{}
	if !i.DisableGuidImport {
		out = append(out, "<guid>")
	}
	import_string := ""
	for j, r := range i.http_fields {
		if r.DisplayName == "guid" {
			continue
		}
		import_string += "<" + r.DisplayName + ">"
		if j != len(i.http_fields)-1 {
			import_string += "|"
		}
	}
	if import_string != "" {
		out = append(out, import_string)
	}
	return out
}

func (i *ImportByHttpFields) genQuery(s string) (string, error) {
	values := url.Values{}
	availableFieldsCount := len(i.http_fields)
	if !i.DisableGuidImport {
		// GUID import enabled
		_, err := guid.FromString(s)
		if err == nil {
			// The given string is a GUID, query string will return guid=<GUID>
			values.Add("guid", s)
			return values.Encode(), nil
		}
		// Remove GUID from the count
		availableFieldsCount--
	}
	// GUID import is disabled
	passedValues := strings.Split(s, "|")
	if len(passedValues) != availableFieldsCount {
		return "", fmt.Errorf(
			"expected %d field values (%s), but got %d (%v)",
			availableFieldsCount,
			i.genNeededFieldsString(),
			len(passedValues),
			passedValues,
		)
	}
	var fieldNamesWithoutGUID []string
	for _, field := range i.http_fields {
		if field.DisplayName != "guid" {
			fieldNamesWithoutGUID = append(fieldNamesWithoutGUID, field.FieldName)
		}
	}
	for n, f := range passedValues {
		values.Add(fieldNamesWithoutGUID[n], f)
	}
	return values.Encode(), nil
}

func (i *ImportByHttpFields) GetDoc() []string {
	return i.fields_String
}

func (i *ImportByHttpFields) getFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	s := fmt.Sprintf("%v", d.Id())
	query, err := i.genQuery(s)
	if err != nil {
		return nil, err
	}
	attr["query"] = query
	return DefaultGetFunc(ctx, _client, attr, d, headers)
}

func (i *ImportByHttpFields) GetFunc() GetFuncType {
	return i.getFunc
}
func GetDefaultImporter() ImportInterface {
	return NewImportByHttpFields(false, []HttpFieldTuple{})

}
