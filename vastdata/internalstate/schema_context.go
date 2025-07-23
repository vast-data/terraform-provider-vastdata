// Copyright (c) HashiCorp, Inc.

// Shared types and utilities to support dynamic schema generation
// for Terraform providers built using the Terraform Plugin Framework.
//
// It defines schema contexts for distinguishing between data sources and managed resources,
// as well as tagging logic for reflecting appropriate field-level annotations depending
// on the context.

package internalstate

import (
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// SchemaContext represents whether the schema is being generated for
// a Terraform data source or a resource block.
type SchemaContext int

const (
	SchemaForDataSource SchemaContext = iota
	SchemaForResource
	SchemaTypeUnknown
)

func (sc SchemaContext) toSchema() any {
	switch sc {
	case SchemaForDataSource:
		return datasource_schema.Schema{}
	case SchemaForResource:
		return resource_schema.Schema{}
	default:
		panic("unknown schema context")
	}
}
