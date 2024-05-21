# Terraform Provider Vastdata


The VastData Terrafrom provider is a provider to manage VastData clusters [resources](./resources).

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
## Submitting Bugs/Feature Requests

While it is common to submit Bugs/Feature Requests using github issues.
We would rather if you open a Bug/Feature Request to Vast Data support at [VastData Support Portal](https://support.vastdata.com/s/login/),
Or by sending mail to support@vastdata.com

