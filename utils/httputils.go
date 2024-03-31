package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func GetResponseBodyAsStr(r *http.Response) string {
	var b bytes.Buffer
	if r == nil {
		return ""
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	//Let's try to make it a pretty json if not we will just dump the body
	err = json.Indent(&b, body, "", "  ")
	if err == nil {
		return string(b.Bytes())
	}
	return string(body)
}

func UnmarshelBodyToMap(r *http.Response, i *map[string]interface{}) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, i)
	if err != nil {
		return err
	}
	return nil
}

func FakeHttpResponse(orig *http.Response, m map[string]interface{}) (*http.Response, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	r := http.Response{
		Request:    orig.Request,
		Status:     orig.Status,
		StatusCode: orig.StatusCode,
		Body:       io.NopCloser(bytes.NewBuffer(b)),
	}

	return &r, nil
}
