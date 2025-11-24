resource "vastdata_active_directory" "vastdb_active_directory" {
  machine_account_name = "machine_acc"
  organizational_unit  = "OU=VASTs,OU=VastENG,DC=VastENG,DC=lab"
  use_auto_discovery   = "false"
  binddn               = "cn=admin,dc=qa,dc=vastdata,dc=com"
  searchbase           = "dc=qa,dc=vastdata,dc=com"
  bindpw               = "vastdata"
  use_ldaps            = "false"
  domain_name          = "VastEng.lab"
  method               = "simple"
  query_groups_mode    = "COMPATIBLE"
  use_tls              = "false"
  urls                 = ["ldap://10.27.252.30"]
}

# ---------------------
# Complete examples
# ---------------------


resource "vastdata_ldap" "vastdb_ldap" {
  domain_name        = "VastEng.lab"
  urls               = ["ldap://10.27.252.30"]
  binddn             = "cn=admin,dc=qa,dc=vastdata,dc=com"
  searchbase         = "dc=qa,dc=vastdata,dc=com"
  bindpw             = "vastdata"
  use_auto_discovery = "false"
  use_ldaps          = "false"
  port               = "389"
  method             = "simple"
  query_groups_mode  = "COMPATIBLE"
  use_tls            = "false"
}

resource "vastdata_active_directory" "ad1" {
  machine_account_name = "sales-devvm-tal"
  organizational_unit  = "OU=VASTs,OU=VastENG,DC=VastENG,DC=lab"
  ldap_id              = vastdata_ldap.vastdb_ldap.id
}

# --------------------


resource "vastdata_ldap" "vastdb_ldap" {
  domain_name        = "VastEng.lab"
  urls               = ["ldap://10.27.252.30"]
  binddn             = "cn=admin,dc=qa,dc=vastdata,dc=com"
  searchbase         = "dc=qa,dc=vastdata,dc=com"
  bindpw             = "vastdata"
  use_auto_discovery = "false"
  use_ldaps          = "false"
  port               = "389"
  method             = "simple"
  query_groups_mode  = "COMPATIBLE"
  use_tls            = "false"
}

resource "vastdata_active_directory" "ad1" {
  machine_account_name = "sales-devvm-tal"
  organizational_unit  = "OU=VASTs,OU=VastENG,DC=VastENG,DC=lab"
  ldap_id              = vastdata_ldap.vastdb_ldap.id
}

# --------------------


resource "vastdata_active_directory" "active_dir3" {
  machine_account_name            = "machine_acc"
  organizational_unit             = "OU=VASTs,OU=VastENG,DC=VastENG,DC=lab"
  use_auto_discovery              = true
  binddn                          = "CN=vastdata,OU=SmbOU,OU=VastENG,DC=VastENG,DC=lab"
  bindpw                          = "vastdata"
  use_ldaps                       = false
  domain_name                     = "VastENG.lab"
  method                          = "simple"
  query_groups_mode               = "COMPATIBLE"
  use_tls                         = false
  is_vms_auth_provider            = true
  smb_allowed                     = true
  ntlm_enabled                    = true
  port                            = 389
  posix_attributes_source         = "JOINED_DOMAIN"
  reverse_lookup                  = false
  gid_number                      = "gidNumber"
  uid                             = "sAMAccountName"
  uid_number                      = "uidNumber"
  use_multi_forest                = true
  match_user                      = "sAMAccountName"
  uid_member_value_property_name  = "sAMAccountName"
  uid_member                      = "member"
  posix_account                   = "user"
  posix_group                     = "group"
  username_property_name          = "name"
  user_login_name                 = "sAMAccountName"
  group_login_name                = "sAMAccountName"
  mail_property_name              = "mail"
  monitor_action                  = "PING"
  abac_read_only_value_name       = "ro"
  abac_read_write_value_name      = "rw"
  scheduled_ma_pwd_change_enabled = false
}

# --------------------

