# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    vastdata = {
      source = "vast-data/vastdata"
      version = "1.7.0"
    }
  }
}

provider "vault" {
  address = "https://vault.mavensecurities.com"
}

data "vault_generic_secret" "vast_terraform" {
  path = "kv/teams/linfra/vast-terraform"
}

# us cluster
provider vastdata {
  username = "admin"
  port = 443
  password = data.vault_generic_secret.vast_terraform.data["vast-admin-password"]
  host = "vast-us-1-admin.mavensecurities.com"
  skip_ssl_verify = true
}

# Backend configuration moved to environment variables for security
# DB that we use postgres-vast-terraform-state-prod.mavensecurities.com
# This is pulled from the vault key teams/linfra/vast-terraform
terraform {
  backend "pg" {
    # conn_str will be provided via TF_VAR_conn_str environment variable
    # or through -backend-config flag during terraform init
  }
}

resource "vastdata_ldap" "ldap1" {
  domain_name        = "mavensecurities.com"
  urls               = ["ldap://10.66.4.19","ldap://10.66.4.20"]
  binddn             = "cn=svcVast,ou=MavenServiceAcc,dc=mavensecurities,dc=com"
  searchbase         = "dc=mavensecurities,dc=com"
  bindpw             = data.vault_generic_secret.vast_terraform.data["ad-password"]
  use_auto_discovery = "false"
  is_vms_auth_provider = "true"
  use_ldaps          = "false"
  port               = "389"
  method             = "simple"
  query_groups_mode  = "COMPATIBLE"
  use_tls            = "false"
  uid                = "sAMAccountName"
  group_login_name   = "sAMAccountName"
  user_login_name    = "sAMAccountName"
  match_user         = "sAMAccountName"
  posix_account      = "user"
  posix_group        = "group"
  username_property_name  = "name"
  uid_member_value_property_name  = "sAMAccountName"
  uid_member         = "member"
}

# the role read_only is a built-in role. In order to create a user assigned to this role, we need to import it
resource "vastdata_administators_roles" "read_only" {
  name             = "read_only"
  ldap_groups      = ["sgCorePlatform","ops","sgLinfra","sgdTechnology-MarketMaking-Data","sgDataMarketMaking"]
}

# the superuser role is a built-in role. In order to create a user assigned to this role, or manage the role further, we need to import it
resource "vastdata_administators_roles" "superuser" {
  name             = "superuser"
  ldap_groups      = ["sgVastAdmin"]

}

resource "vastdata_administators_roles" "csi" {
  name             = "csi"
}

resource "vastdata_administators_managers" "prometheus_reader" {
  username         = "prometheus_reader"
  password         =  data.vault_generic_secret.vast_terraform.data["prometheus_reader-password"]
  roles            = [vastdata_administators_roles.read_only.id]
  permissions_list = ["create_monitoring"]
}

resource "vastdata_administators_managers" "csi" {
  username         = "csi"
  password         = data.vault_generic_secret.vast_terraform.data["csi-password"]
  roles            = [vastdata_administators_roles.csi.id]
}
