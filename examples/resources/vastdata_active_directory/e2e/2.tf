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