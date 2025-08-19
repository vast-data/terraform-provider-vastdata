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
	"regexp"
	"sort"
	"strings"
)

// sensitiveDataPatterns defines regex patterns for sensitive data that should be sanitized
var sensitiveDataPatterns = []struct {
	name    string
	pattern *regexp.Regexp
	replace string
}{
	{
		name:    "PGP Public Key",
		pattern: regexp.MustCompile(`(?s)-----BEGIN PGP PUBLIC KEY BLOCK-----\s*\n.*?\n\s*-----END PGP PUBLIC KEY BLOCK-----`),
		replace: "-----BEGIN PGP PUBLIC KEY BLOCK-----\n    .\n    .  <content>\n    .\n-----END PGP PUBLIC KEY BLOCK-----",
	},
	{
		name:    "PGP Private Key",
		pattern: regexp.MustCompile(`(?s)-----BEGIN PGP PRIVATE KEY BLOCK-----\s*\n.*?\n\s*-----END PGP PRIVATE KEY BLOCK-----`),
		replace: "-----BEGIN PGP PRIVATE KEY BLOCK-----\n    .\n    .  <content>\n    .\n-----END PGP PRIVATE KEY BLOCK-----",
	},
	{
		name:    "X.509 Certificate",
		pattern: regexp.MustCompile(`(?s)-----BEGIN CERTIFICATE-----\s*\n.*?\n\s*-----END CERTIFICATE-----`),
		replace: "-----BEGIN CERTIFICATE-----\n    .\n    .  <content>\n    .\n-----END CERTIFICATE-----",
	},
	{
		name:    "X.509 Private Key",
		pattern: regexp.MustCompile(`(?s)-----BEGIN PRIVATE KEY-----\s*\n.*?\n\s*-----END PRIVATE KEY-----`),
		replace: "-----BEGIN PRIVATE KEY-----\n    .\n    .  <content>\n    .\n-----END PRIVATE KEY-----",
	},
	{
		name:    "RSA Private Key",
		pattern: regexp.MustCompile(`(?s)-----BEGIN RSA PRIVATE KEY-----\s*\n.*?\n\s*-----END RSA PRIVATE KEY-----`),
		replace: "-----BEGIN RSA PRIVATE KEY-----\n    .\n    .  <content>\n    .\n-----END RSA PRIVATE KEY-----",
	},
	{
		name:    "DSA Private Key",
		pattern: regexp.MustCompile(`(?s)-----BEGIN DSA PRIVATE KEY-----\s*\n.*?\n\s*-----END DSA PRIVATE KEY-----`),
		replace: "-----BEGIN DSA PRIVATE KEY-----\n    .\n    .  <content>\n    .\n-----END DSA PRIVATE KEY-----",
	},
	{
		name:    "EC Private Key",
		pattern: regexp.MustCompile(`(?s)-----BEGIN EC PRIVATE KEY-----\s*\n.*?\n\s*-----END EC PRIVATE KEY-----`),
		replace: "-----BEGIN EC PRIVATE KEY-----\n    .\n    .  <content>\n    .\n-----END EC PRIVATE KEY-----",
	},
	{
		name:    "OpenSSH Private Key",
		pattern: regexp.MustCompile(`(?s)-----BEGIN OPENSSH PRIVATE KEY-----\s*\n.*?\n\s*-----END OPENSSH PRIVATE KEY-----`),
		replace: "-----BEGIN OPENSSH PRIVATE KEY-----\n    .\n    .  <content>\n    .\n-----END OPENSSH PRIVATE KEY-----",
	},
	{
		name:    "SSH2 Private Key",
		pattern: regexp.MustCompile(`(?s)---- BEGIN SSH2 PRIVATE KEY ----\s*\n.*?\n\s*---- END SSH2 PRIVATE KEY ----`),
		replace: "---- BEGIN SSH2 PRIVATE KEY ----\n    .\n    .  <content>\n    .\n---- END SSH2 PRIVATE KEY ----",
	},
}

func main() {
	baseDirs := []struct {
		path   string
		output string
	}{
		{"../examples/resources", "resource.tf"},
		{"../examples/data-sources", "data-source.tf"},
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

				// Filter out files explicitly ignored or where all non-empty lines are commented out (start with '#')
				eligible := make([]string, 0, len(e2eFiles))
				for _, file := range e2eFiles {
					// Skip examples explicitly marked to ignore via first non-empty line
					if ignore, err := hasIgnoreExampleDirective(file); err != nil {
						return fmt.Errorf("error scanning %s: %w", file, err)
					} else if ignore {
						continue
					}

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
	var content strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Copyright (c) HashiCorp") {
			continue
		}
		content.WriteString(line + "\n")
		if strings.TrimSpace(line) != "" {
			hasContent = true
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	// Sanitize sensitive data
	sanitizedContent := sanitizeSensitiveData(content.String())
	builder.WriteString(sanitizedContent)

	return hasContent, nil
}

// sanitizeSensitiveData replaces sensitive data with placeholder content
func sanitizeSensitiveData(content string) string {
	sanitized := content

	for _, pattern := range sensitiveDataPatterns {
		if pattern.pattern.MatchString(sanitized) {
			fmt.Printf("Sanitized %s in content\n", pattern.name)
			sanitized = pattern.pattern.ReplaceAllString(sanitized, pattern.replace)
		}
	}

	return sanitized
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

// hasIgnoreExampleDirective returns true if the first non-empty line equals
// "# ignore:example" (exact match after trimming whitespace).
func hasIgnoreExampleDirective(path string) (bool, error) {
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
		return trimmed == "# ignore:example", nil
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return false, nil
}
