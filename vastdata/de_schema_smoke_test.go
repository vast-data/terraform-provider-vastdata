// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestAllResourceAndDatasourceSchemasBuild(t *testing.T) {
	ctx := context.Background()
	for _, factory := range GetResourceFactories() {
		r := factory()
		resp := &resource.SchemaResponse{}
		r.Schema(ctx, resource.SchemaRequest{}, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("resource %T schema: %v", r, resp.Diagnostics)
		}
	}
	for _, factory := range GetDatasourceFactories() {
		d := factory()
		resp := &datasource.SchemaResponse{}
		d.Schema(ctx, datasource.SchemaRequest{}, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("datasource %T schema: %v", d, resp.Diagnostics)
		}
	}
}
