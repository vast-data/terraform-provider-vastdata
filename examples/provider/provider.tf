#A Simple provider with warn version validation level
provider vastdata {
  username                = "<username>"
  port                    = 443
  password                = "<password>"
  host                    = "<address>"
  skip_ssl_verify         = true
  version_validation_mode = "warn"
}

#Define 2 providers for 2 different cluster with alias one with port 443 and one with port 9443

provider vastdata {
  api_token       = "<api_token>"
  port            = 443
  host            = "<address>"
  skip_ssl_verify = true
  alias           = clusterA
}

# Trigger Terraform to ask for password instead of hardcoding it
variable "password" {
  sensitive = true
}

provider vastdata {
  username        = "<username>"
  port            = 9443
  password        = var.password
  host            = "<address>"
  skip_ssl_verify = true
  alias           = clusterB
}
