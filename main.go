// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/vast-data/terraform-provider-vastdata/vastdata"
)

var version string = "dev"

func main() {
	fmt.Fprintln(os.Stderr, ">>> starting terraform-provider-vastdata")

	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/vastdata/vastdata",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), vastdata.New(version), opts)
	if err != nil {
		log.Fatalf("provider failed to start: %s", err)
	}
}
