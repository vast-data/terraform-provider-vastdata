// Copyright (c) HashiCorp, Inc.

//go:build ignore
// +build ignore

package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
)

func main() {
	// Get paths from environment variables or use defaults
	inputPath := os.Getenv("INPUT_PATH")
	if inputPath == "" {
		inputPath = "/tmp/apiconv/swagger.json"
	}

	outputDir := os.Getenv("OUTPUT_DIR")
	if outputDir == "" {
		outputDir = "/tmp/apiconv"
	}

	jsonOut := filepath.Join(outputDir, "api.json")
	tarOut := filepath.Join(outputDir, "api.tar.gz")

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("❌ Failed to create output directory: %v", err)
	}

	// Read Swagger v2 (OpenAPI 2.0) JSON
	data, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("❌ Failed to read input: %v", err)
	}

	var docV2 openapi2.T
	if err := json.Unmarshal(data, &docV2); err != nil {
		log.Fatalf("❌ Failed to parse Swagger v2: %v", err)
	}

	// Convert to OpenAPI v3
	docV3, err := openapi2conv.ToV3(&docV2)
	if err != nil {
		log.Fatalf("❌ Failed to convert to OpenAPI v3: %v", err)
	}

	// Save as compact JSON (no spaces)
	f, err := os.Create(jsonOut)
	if err != nil {
		log.Fatalf("❌ Failed to create api.json: %v", err)
	}
	if err := json.NewEncoder(f).Encode(docV3); err != nil {
		log.Fatalf("❌ Failed to write api.json: %v", err)
	}
	f.Close()
	log.Println("✅ Saved OpenAPI v3 to", jsonOut)

	// Create tar.gz archive
	tarFile, err := os.Create(tarOut)
	if err != nil {
		log.Fatalf("❌ Failed to create tar.gz: %v", err)
	}
	defer tarFile.Close()

	gzw := gzip.NewWriter(tarFile)
	defer gzw.Close()
	tw := tar.NewWriter(gzw)
	defer tw.Close()

	fileInfo, err := os.Stat(jsonOut)
	if err != nil {
		log.Fatalf("❌ Cannot stat output file: %v", err)
	}

	header, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		log.Fatalf("❌ Failed to create tar header: %v", err)
	}
	header.Name = "api.json"

	if err := tw.WriteHeader(header); err != nil {
		log.Fatalf("❌ Failed to write tar header: %v", err)
	}

	fileData, err := os.ReadFile(jsonOut)
	if err != nil {
		log.Fatalf("❌ Failed to read api.json: %v", err)
	}
	if _, err := tw.Write(fileData); err != nil {
		log.Fatalf("❌ Failed to write file to tar: %v", err)
	}

	log.Println("📦 Created archive", tarOut)
}
