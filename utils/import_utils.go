package utils

import (
	"context"
	"errors"
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
	if !i.DisableGuidImport {
		// we check if GUID was provided
		_, err := guid.FromString(s)
		if err == nil {
			//The given string is a GUID, query string will return guid=<GUID>
			values.Add("guid", s)
			return values.Encode(), nil
		} else if len(i.http_fields) == 1 && err != nil {
			return "", err
		}
	}
	//If we got here wither importing by GUID is disabled or that the given string is not a valid GUID
	l := len(i.http_fields) - 1
	q := strings.Split(s, "|")
	if len(q) != l {
		return "", errors.New(fmt.Sprintf("No Enough Fields provider, needed fields are %s", i.genNeededFieldsString()))
	}
	for n, f := range q {
		values.Add(i.http_fields[n+1].FieldName, f)
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
