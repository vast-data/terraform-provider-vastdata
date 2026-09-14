#A Simple provider with warn version validation level
provider "vastdata" {
  username                = "<username>"
  port                    = 443
  password                = "<password>"
  host                    = "<address>"
  skip_ssl_verify         = true
  version_validation_mode = "warn"
}

# Two clusters via Terraform provider aliases.
# alias is a Terraform meta-argument, not a provider schema attribute.
# Resources must set provider = vastdata.<alias> to select a cluster.

provider "vastdata" {
  alias           = "clusterA"
  api_token       = "<api_token>"
  port            = 443
  host            = "<address>"
  skip_ssl_verify = true
}

# Trigger Terraform to ask for password instead of hardcoding it
variable "password" {
  sensitive = true
}

provider "vastdata" {
  alias           = "clusterB"
  username        = "<username>"
  port            = 9443
  password        = var.password
  host            = "<address>"
  skip_ssl_verify = true
}

# resource "vastdata_tenant" "example" {
#   provider = vastdata.clusterA
#   name     = "example"
# }
