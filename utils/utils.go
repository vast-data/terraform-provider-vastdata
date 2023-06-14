package utils

import (
	"regexp"
	"strings"
)

func GetTFformatName(s string) string {
	re := regexp.MustCompile("([A-Z])")
	t := re.ReplaceAllString(s, "_${1}")
	return strings.ToLower(strings.TrimLeft(t, "_"))
}
