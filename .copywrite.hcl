schema_version = 1

project {
  copyright_year = 2025

  header_ignore = [
    # examples used within documentation (prose)
    "examples/**",

    # GitHub Actions workflow-specific configurations
    ".github/labeler-*.yml",

    # golangci-lint tooling configuration
    ".golangci.yml",

    # Release Engineering tooling configuration
    ".release/*.hcl",
  ]
}
