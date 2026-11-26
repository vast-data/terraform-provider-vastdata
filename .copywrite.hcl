schema_version = 1

project {
  copyright_year = 2025

  header_ignore = [
    # examples used within documentation (prose)
    "examples/**",

    # DataEngine modules are VAST-owned; do not stamp HashiCorp headers
    "modules/**",

    # GitHub Actions workflow-specific configurations
    ".github/labeler-*.yml",

    # golangci-lint tooling configuration
    ".golangci.yml",

    # Release Engineering tooling configuration
    ".release/*.hcl",

    # Test files and virtual environment
    "tests/**"
  ]
}
