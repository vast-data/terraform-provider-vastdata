// Copyright (c) HashiCorp, Inc.

//go:build ignore
// +build ignore

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	vastdata "github.com/vast-data/terraform-provider-vastdata/vastdata"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/schema_generation"
)

var (
	mode       string
	filter     string
	automation bool
)

func main() {
	flag.StringVar(&mode, "type", "", "Specify type: 'datasource' or 'resource'")
	flag.StringVar(&filter, "filter", "", "Optional filter for resource path substring")
	flag.BoolVar(&automation, "automation", false, "Use automation marker (🔸 instead of 🔹)")
	flag.Parse()

	switch mode {
	case "datasource", "d":
		printDataSourceSchemas(filter, automation)
	case "resource", "r":
		printResourceSchemas(filter, automation)
	default:
		fmt.Fprintln(os.Stderr, "Usage: go run main.go -type=datasource|resource [-filter=...] [-automation]")
		os.Exit(1)
	}
}

func printDataSourceSchemas(filter string, automation bool) {
	for _, factory := range vastdata.GetDatasourceFactories() {
		instance := factory()
		d, ok := instance.(*vastdata.Datasource)
		if !ok {
			panic(fmt.Sprintf("unexpected type: %T", instance))
		}

		var ctx context.Context
		manager := d.EmptyManager()
		tfState := manager.TfState()
		hints := tfState.Hints
		resourceName := fmt.Sprintf("vastdata_%s", is.SnakeCaseName(manager))

		if filter != "" && !containsIgnoreCase(resourceName, filter) {
			continue
		}
		if automation {
			ctx = context.Background()
		}

		schema, err := schema_generation.GetDatasourceSchema(ctx, hints)
		if err != nil {
			panic(err)
		}

		// Handle custom datasources that don't have SchemaRef
		var pathInfo string
		if hints.TFStateHintsForCustom != nil {
			pathInfo = "custom"
		} else if hints.SchemaRef != nil && hints.SchemaRef.Read != nil {
			pathInfo = hints.SchemaRef.Read.Path
		} else {
			pathInfo = "unknown"
		}

		fmt.Printf("%s: # Datasource (%s)\n", resourceName, pathInfo)
		visualization := is.BuildDataSourceAttributesString(schema.Attributes, automation, 2)
		fmt.Println(visualization)
		fmt.Println()
	}
}

func printResourceSchemas(filter string, automation bool) {
	for _, factory := range vastdata.GetResourceFactories() {
		instance := factory()
		r, ok := instance.(*vastdata.Resource)
		if !ok {
			panic(fmt.Sprintf("unexpected type: %T", instance))
		}

		var ctx context.Context
		manager := r.EmptyManager()
		tfState := manager.TfState()
		hints := tfState.Hints
		resourceName := fmt.Sprintf("vastdata_%s", is.SnakeCaseName(manager))

		if filter != "" && !containsIgnoreCase(resourceName, filter) {
			continue
		}

		if automation {
			ctx = context.Background()
		}

		schema, err := schema_generation.GetResourceSchema(ctx, hints)
		if err != nil {
			panic(err)
		}

		// Handle custom resources that don't have SchemaRef
		var pathInfo string
		if hints.TFStateHintsForCustom != nil {
			pathInfo = "custom"
		} else if hints.SchemaRef != nil && hints.SchemaRef.Create != nil {
			pathInfo = hints.SchemaRef.Create.Path
		} else {
			pathInfo = "unknown"
		}

		fmt.Printf("%s: # Resource (%s)\n", is.SnakeCaseName(manager), pathInfo)
		visualization := is.BuildResourceAttributesString(schema.Attributes, automation, 2)
		fmt.Println(visualization)
		fmt.Println()
	}
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
