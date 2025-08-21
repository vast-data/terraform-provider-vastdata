# Terraform Provider Vastdata

[![CI](https://github.com/vast-data/terraform-provider-vastdata/workflows/CI/badge.svg)](https://github.com/vast-data/terraform-provider-vastdata/actions/workflows/ci.yml)
[![Release](https://github.com/vast-data/terraform-provider-vastdata/workflows/Release/badge.svg)](https://github.com/vast-data/terraform-provider-vastdata/actions/workflows/release.yml)
[![License: Apache2](https://img.shields.io/badge/License-Apache2-yellow.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/vast-data/terraform-provider-vastdata)](https://goreportcard.com/report/github.com/vast-data/terraform-provider-vastdata)
[![Coverage Status](https://codecov.io/gh/vast-data/terraform-provider-vastdata/branch/main/graph/badge.svg)](https://codecov.io/gh/vast-data/terraform-provider-vastdata)
[![Go Reference](https://pkg.go.dev/badge/github.com/vast-data/terraform-provider-vastdata.svg)](https://pkg.go.dev/github.com/vast-data/terraform-provider-vastdata)
[![Terraform Registry](https://img.shields.io/badge/terraform-registry-blue.svg)](https://registry.terraform.io/providers/vast-data/vastdata/latest)
[![Latest Release](https://img.shields.io/github/v/release/vast-data/terraform-provider-vastdata)](https://github.com/vast-data/terraform-provider-vastdata/releases/latest)

The VastData Terraform provider is a provider to manage VastData clusters [resources](./resources).

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [Terraform Provider Vastdata](#terraform-provider-vastdata)
    - [Configuring Provider to be downloaded from Terraform Registry](#configuring-provider-to-be-downloaded-from-terraform-registryhttpsregistryterraformioprovidersvast-datavastdatalatest)
    - [Building The Provider Directly From Github](#building-the-provider-directly-from-github)
    - [Building The Provider Locally](#building-the-provider-locally)
        - [Cloning The Source Code](#cloning-the-source-code)
        - [Building The Code Locally](#building-the-code-locally)
        - [Building Into A Differant directory](#building-into-a-differant-directory)
    - [Using Local Copy Using dev_overrides](#using-local-copy-using-dev_overrides)
        - [Edit ~/.terraformrc file](#edit-terraformrc-file)
    - [Importing Existing Resources](#importing-existing-resources)
        - [Import Formats](#import-formats)
        - [Import Field Types](#import-field-types)
        - [Examples](#examples)
        - [Troubleshooting Import](#troubleshooting-import)
        - [Finding Import Fields](#finding-import-fields)
    - [Submitting Bugs/Feature Requests](#submitting-bugsfeature-requests)

<!-- markdown-toc end -->


## Configuring Provider to be downloaded from [Terraform Registry](https://registry.terraform.io/providers/vast-data/vastdata/latest)
In order to configure the provider to be used directly from Terraform registry, use the following provider defenition. 
```hcl
terraform {
  required_providers {
    vastdata = {
      source = "vast-data/vastdata"
    }
  }
}
```
Now when running `terraform init` it will download the VastData provider from the [Terraform Registry](https://registry.terraform.io/providers/vast-data/vastdata/latest)

## Building The Provider Directly From Github

In order to build the provider you can simpy use go install 

```bash
$ go install github.com/vast-data/terraform-provider-vastdata
```

To install a specfic tag/branch version use the following syntax 
```bash
$ go install github.com/vast-data/terraform-provider-vastdata@<branch/tag>
```

Check you go GOBIN path for the compiled file named `terraform-provider-vastdata`

## Building The Provider Locally 

In order to build the provider locally you will need to first clone the repo.

### Cloning The Source Code

```bash
$ git clone https://github.com/vast-data/terraform-provider-vastdata.git
$ cd terraform-provider-vastdata
```

If you wish to checkout a specific tag check out the tag.

```bash
$ git checkout <tag name>
```

### Building The Code Locally 

***In order to build the code locally you need to have [GNU Make](https://www.gnu.org/software/make/) installed***

To build the code ru the following command

```bash
make build
```

This will build the provider binary with the name *terraform-provider-vastdata* directly into the build directory.

### Building Into A Differant directory

In order to build to a differant directory other than build , specify BUILD_DEST=<*build directory*>

```bash
	make build BUILD_DEST=some/other/directory
```

## Using Local Copy Using dev_overrides

dev_overrides is a terraform configuration that will allow to overrides any other method of obtaining the terrafrom plugin binary and forces terrafrom to obtain the provider binary from a specific path locally.

### Edit ~/.terraformrc file 

*This file is being scanned every time you run terrafrom if you dont want to create/edit it you can specify the environment vasriable TF_CLI_CONFIG_FILE=path/to/configuration/file*

add the following configurations to the configuration file.

```hcl
provider_installation {
  dev_overrides {
    "vastdata/vastdata" = "/some/directory/where/the/binary/is/stored"
  }
  direct {}
} 
```

When creating a terrafrom configration specify the following at the file 

```hcl
terraform {
  required_providers {
    vastdata = {
      source  = "vastdata/vastdata"
    }
  }
}
```

Now you can define providers.

```hcl
provider vastdata {
username = "<username>"
port = <port>
password = "<password>"
host = "<address>"
skip_ssl_verify = true
version_validation_mode = "warn"
}
```

## Importing Existing Resources

The VastData provider supports importing existing resources using various ID formats, including composite keys for resources that require multiple identifiers.

### Import Formats

#### 1. Simple ID Import
For resources that use a single identifier:
```bash
terraform import vastdata_example.my_resource "12345"
```

#### 2. Key-Value Pairs Import
For resources requiring multiple fields, use key=value format with comma or semicolon separators:
```bash
# Using comma separator
terraform import vastdata_example.my_resource "gid=1001,tenant_id=22,context=ad"

# Using semicolon separator  
terraform import vastdata_example.my_resource "gid=1001;tenant_id=22;context=ad"
```

#### 3. Ordered Values Import (Pipe-separated)
For resources with predefined import field order, use pipe-separated values:
```bash
terraform import vastdata_example.my_resource "1001|22|ad"
```

### Import Field Types

The provider automatically handles type conversion for imported values:

- **String fields**: Values are imported as-is
- **Integer fields**: Numeric strings are converted to integers
- **Boolean fields**: Accepts `true`, `false`, `1`, or `0`

### Examples

#### Import a User with Multiple Identifiers
```bash
# Key-value format
terraform import vastdata_user.admin "username=admin,tenant_id=1,domain=local"

# Ordered format (if resource supports it)
terraform import vastdata_user.admin "admin|1|local"
```

#### Import a Quota with Composite Key
```bash
terraform import vastdata_quota.project_quota "name=project1,path=/data/project1,tenant_id=5"
```

#### Import a Network Interface
```bash
terraform import vastdata_network_interface.eth0 "name=eth0,node_id=1"
```

### Troubleshooting Import

**Error: "field 'x' is not present in the resource schema"**
- Ensure the field name matches exactly what's defined in the resource schema
- Check the resource documentation for correct field names

**Error: "expected X values for fields [...], got Y"**
- When using pipe-separated format, ensure the number of values matches the expected import fields
- Use key=value format instead if you need to specify only some fields

**Error: "invalid int64 for field 'x'"**
- Ensure numeric fields contain valid integer values
- Check for extra spaces or non-numeric characters

### Finding Import Fields

To determine which fields are required for importing a specific resource:

1. Check the resource documentation
2. Look at the resource's required attributes
3. Use `terraform plan` after creating a minimal resource configuration to see required fields
4. Refer to the VastData API documentation for the underlying resource identifiers

# Submitting Bugs/Feature Requests

While it is common to submit Bugs/Feature Requests using github issues,
we would rather if you open a Bug/Feature Request to Vast Data support at customer.support@vastdata.com

