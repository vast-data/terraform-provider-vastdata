resource "vastdata_active_directory2" "ad1" {
  machine_account_name = "vast-cluster01"
  organizational_unit  = "OU=VASTs,OU=VastENG,DC=VastENG,DC=lab"
  use_auto_discovery   = false
  binddn               = "cn=admin,dc=qa,dc=vastdata,dc=com"
  searchbase           = "dc=qa,dc=vastdata,dc=com"
  bindpw               = "<password>"
  use_ldaps            = "false"
  domain_name          = "VastEng.lab"
  method               = "simple"
  query_groups_mode    = "COMPATIBLE"
  use_tls              = "false"
  urls                 = ["ldap://198.51.100.3"]

}
