## Description

Brief description of what this PR does.

## Type of Change

Please delete options that are not relevant.

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New resource or data source
- [ ] Enhancement to existing resource/data source
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Code refactoring (no functional changes)
- [ ] Performance improvement
- [ ] Test coverage improvement
- [ ] Provider configuration improvement

## Changes Made

- 
- 
- 

## Resources/Data Sources Added or Modified

- [ ] vastdata_user
- [ ] vastdata_group
- [ ] vastdata_view
- [ ] vastdata_quota
- [ ] vastdata_tenant
- [ ] vastdata_snapshot
- [ ] vastdata_protection_policy
- [ ] vastdata_qos_policy
- [ ] vastdata_s3_policy
- [ ] vastdata_replication_peer
- [ ] Other: ____________

## Testing

### Unit Tests
- [ ] I have added unit tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally (`make test`)
- [ ] Unit tests pass with race detection (`make test-unit`)
- [ ] Code passes go vet (`make vet`)
- [ ] Code passes formatting check (`make fmt`)

### Terraform Testing
- [ ] I have tested the Terraform configuration examples in `examples/`
- [ ] `terraform init` works with the provider changes
- [ ] `terraform validate` passes for affected examples
- [ ] `terraform plan` works as expected (where applicable)

### Acceptance Tests
- [ ] I have added/updated acceptance tests (`make testacc`)
- [ ] Acceptance tests pass against a real VAST cluster (if applicable)
- [ ] Tests cover both success and failure scenarios

### Manual Testing
- [ ] I have manually tested the changes with a VAST cluster
- [ ] Import functionality works (if applicable)
- [ ] Resource updates work correctly (if applicable)
- [ ] Resource deletion works correctly (if applicable)

## Documentation

- [ ] I have updated the resource/data source documentation
- [ ] I have updated/added examples in `examples/`
- [ ] I have updated the CHANGELOG.md (if applicable)
- [ ] Documentation generation works (`make generate-docs`)
- [ ] My changes generate no new warnings in the documentation build

## Code Quality

- [ ] My code follows the Go style guidelines
- [ ] My code follows Terraform provider best practices
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] My changes generate no new lint warnings
- [ ] I have run the full test suite locally (`make test-all`)

## VAST API Integration

- [ ] I have verified the VAST API endpoints used
- [ ] Error handling covers VAST API error responses
- [ ] Resource attributes map correctly to VAST API fields
- [ ] API version compatibility is maintained

## Configuration Example

If this PR adds or modifies resources/data sources, provide a working example:

```hcl
terraform {
  required_providers {
    vastdata = {
      source = "vast-data/vastdata"
      version = "~> 1.0"
    }
  }
}

provider "vastdata" {
  host     = "vast.example.com"
  username = "admin"
  password = "password"
}

# Example configuration for the new/modified resource
resource "vastdata_example" "test" {
  name        = "test-resource"
  description = "Test resource for PR validation"
  
  # Include relevant attributes
}
```

## Breaking Changes

If this PR contains breaking changes, please describe them here and how users should adapt their code:

## Performance Impact

- [ ] No performance impact
- [ ] Positive performance impact
- [ ] Potential performance impact (please describe below)

Performance notes:

## Security Considerations

- [ ] No security impact
- [ ] Security improvement
- [ ] Potential security impact (please describe below)

Security notes:

## Additional Notes

Any additional information that reviewers should know about this PR.

## Related Issues

Fixes #(issue number)
Closes #(issue number)
Relates to #(issue number)

## Checklist for Maintainers

- [ ] PR title follows conventional commit format
- [ ] All CI checks pass
- [ ] Documentation is complete and accurate
- [ ] Breaking changes are properly documented
- [ ] Version bump is appropriate (if needed) 