---
name: Bug report
about: Create a report to help us improve the VastData Terraform Provider
title: '[BUG] '
labels: ['bug']
assignees: ''

---

## Bug Description

A clear and concise description of what the bug is.

## Terraform Configuration

Please provide the relevant Terraform configuration that reproduces the issue:

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

# Your resource configuration that causes the issue
resource "vastdata_user" "example" {
  name = "test-user"
  # ... other configuration
}
```

## Steps to Reproduce

Steps to reproduce the behavior:

1. Create the Terraform configuration above
2. Run `terraform init`
3. Run `terraform plan` / `terraform apply`
4. See error

## Expected Behavior

A clear and concise description of what you expected to happen.

## Actual Behavior

A clear and concise description of what actually happened.

## Error Message

```
Paste the full error message here, including any stack traces
```

## Environment

- **Terraform version**: [e.g. v1.5.7] (run `terraform version`)
- **VastData provider version**: [e.g. v1.0.0] (check your `terraform.lock.hcl`)
- **VAST OS version**: [e.g. 5.1.0] (from your VAST cluster)
- **Operating System**: [e.g. Ubuntu 22.04, macOS 13.0, Windows 11]
- **Go version**: [e.g. go1.22.0] (if building from source)

## Debug Information

Please provide debug output if available:

```bash
# Run with debug logging enabled
TF_LOG=DEBUG terraform apply
```

<details>
<summary>Debug Output (click to expand)</summary>

```
Paste debug output here
```

</details>

## Resource/Data Source Affected

- [ ] vastdata_user
- [ ] vastdata_group
- [ ] vastdata_view
- [ ] vastdata_quota
- [ ] vastdata_tenant
- [ ] vastdata_snapshot
- [ ] vastdata_protection_policy
- [ ] vastdata_qos_policy
- [ ] Other: ____________

## Additional Context

Add any other context about the problem here, such as:
- Does this happen consistently or intermittently?
- Are there any workarounds?
- Related VAST API endpoints or documentation
- Screenshots of VAST GUI (if relevant)

## Possible Solution

If you have ideas about what might be causing the issue or how to fix it, please share them here. 