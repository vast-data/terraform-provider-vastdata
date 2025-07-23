// Copyright (c) HashiCorp, Inc.

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"
	"strconv"
)

func getenvOr(val types.String, envKey string) string {
	if !val.IsNull() {
		return val.ValueString()
	}
	if v := os.Getenv(envKey); v != "" {
		return v
	}
	return ""
}

func int64Or(val types.Int64, envKey string, def int64) int64 {
	if !val.IsNull() {
		return val.ValueInt64()
	}
	if v := os.Getenv(envKey); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			return parsed
		}
	}
	return def
}

func boolOr(val types.Bool, envKey string, def bool) bool {
	if !val.IsNull() {
		return val.ValueBool()
	}
	if v := os.Getenv(envKey); v != "" {
		return v == "1" || v == "true" || v == "TRUE"
	}
	return def
}
