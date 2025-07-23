---
name: Feature request
about: Suggest a new resource, data source, or feature for the VastData Terraform Provider
title: '[FEATURE] '
labels: ['enhancement']
assignees: ''

---

## Is your feature request related to a problem?

A clear and concise description of what the problem is. Ex. I'm always frustrated when [...]

## Describe the solution you'd like

A clear and concise description of what you want to happen.

## Feature Type

What type of feature are you requesting?

- [ ] New resource (e.g. `vastdata_new_resource`)
- [ ] New data source (e.g. `data "vastdata_new_data_source"`)
- [ ] Enhancement to existing resource/data source
- [ ] Provider configuration improvement
- [ ] Documentation improvement
- [ ] Other: ____________

## Proposed Terraform Configuration

If you have ideas about how the Terraform configuration should look, please provide examples:

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

# Example of how you envision using this feature
resource "vastdata_new_resource" "example" {
  name        = "example-resource"
  description = "Example description"
  
  # Your proposed configuration here
  setting1 = "value1"
  setting2 = 42
  
  tags = {
    environment = "production"
    team        = "devops"
  }
}

# Or for data sources
data "vastdata_new_data_source" "example" {
  name = "existing-resource"
}

output "resource_info" {
  value = data.vastdata_new_data_source.example
}
```

## Use Case

Describe the specific use case for this feature:

- **What are you trying to accomplish?**
- **How would this feature help?**
- **Who would benefit from this feature?**
- **What Terraform workflows would this enable?**

## VAST API Reference

If this feature relates to a specific VAST API endpoint, please provide:

- **API endpoint**: [e.g. `/api/views/`, `/api/users/`]
- **VAST documentation link**: [if available]
- **API version**: [e.g. v5.1, v5.2]
- **HTTP methods**: [e.g. GET, POST, PUT, DELETE]

## Expected Resource/Data Source Schema

If requesting a new resource or data source, describe the expected attributes:

```hcl
# Expected schema structure
resource "vastdata_new_resource" "example" {
  # Required attributes
  name = string
  
  # Optional attributes
  description = string
  enabled     = bool
  
  # Complex nested attributes
  configuration {
    setting1 = string
    setting2 = number
    
    nested_block {
      option = string
    }
  }
  
  # Computed attributes (read-only)
  id           = string (computed)
  created_time = string (computed)
  status       = string (computed)
}
```

## Describe alternatives you've considered

A clear and concise description of any alternative solutions or features you've considered.

## Additional Context

Add any other context, screenshots, VAST GUI references, or examples about the feature request here.

## Implementation Notes

If you have technical insights about implementation:

- Related VAST concepts or terminology
- Dependencies on other resources
- State management considerations
- Import requirements
- Validation needs

## Community Impact

How would this feature benefit the broader Terraform and VAST community? 