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
