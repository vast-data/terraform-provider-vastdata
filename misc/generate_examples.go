// Copyright (c) HashiCorp, Inc.

//go:build ignore
// +build ignore

package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	baseDirs := []struct {
		path   string
		output string
	}{
		{"examples/resources", "resource.tf"},
		{"examples/data-sources", "data-source.tf"},
	}

	for _, base := range baseDirs {
		err := filepath.WalkDir(base.path, func(path string, d fs.DirEntry, err error) error {
			if err != nil || !d.IsDir() || path == base.path {
				return nil
			}

			basicPath := filepath.Join(path, "basic.tf")
			e2ePath := filepath.Join(path, "e2e")

			if _, err := os.Stat(basicPath); err != nil {
				return nil // skip if no basic.tf
			}

			var builder strings.Builder

			// Append cleaned basic.tf
			hasContent, err := appendCleaned(&builder, basicPath)
			if err != nil {
				return fmt.Errorf("error reading basic.tf in %s: %w", path, err)
			}
			if !hasContent {
				return nil // skip if basic.tf is empty or only contains ignored lines
			}

			// Add E2E section if there are any *.tf files
			e2eFiles, err := filepath.Glob(filepath.Join(e2ePath, "*.tf"))
			if err == nil && len(e2eFiles) > 0 {
				sort.Strings(e2eFiles)

				// Filter out files where all non-empty lines are commented out (start with '#')
				eligible := make([]string, 0, len(e2eFiles))
				for _, file := range e2eFiles {
					hasContent, err := hasUncommentedNonEmptyLine(file)
					if err != nil {
						return fmt.Errorf("error scanning %s: %w", file, err)
					}
					if hasContent {
						eligible = append(eligible, file)
					}
				}

				if len(eligible) > 0 {
					builder.WriteString("\n# ---------------------\n")
					builder.WriteString("# Complete examples\n")
					builder.WriteString("# ---------------------\n\n")

					for _, file := range eligible {
						if _, err := appendCleaned(&builder, file); err != nil {
							return fmt.Errorf("error reading %s: %w", file, err)
						}
						builder.WriteString("\n# --------------------\n\n")
					}
				}
			}

			outPath := filepath.Join(path, base.output)
			if err := os.WriteFile(outPath, []byte(builder.String()), 0644); err != nil {
				return fmt.Errorf("error writing %s: %w", outPath, err)
			}

			fmt.Printf("✅ Generated: %s\n", outPath)
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		}
	}
}

func appendCleaned(builder *strings.Builder, path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var hasContent bool
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Copyright (c) HashiCorp") {
			continue
		}
		builder.WriteString(line + "\n")
		if strings.TrimSpace(line) != "" {
			hasContent = true
		}
	}
	return hasContent, scanner.Err()
}

// hasUncommentedNonEmptyLine returns true if the file contains at least one
// non-empty line that does not start with '#'. Lines with only whitespace are ignored.
func hasUncommentedNonEmptyLine(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		trimmed := strings.TrimSpace(scanner.Text())
		if trimmed == "" {
			continue
		}
		if !strings.HasPrefix(trimmed, "#") {
			return true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return false, nil
}
