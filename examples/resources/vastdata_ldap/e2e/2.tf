resource "vastdata_ldap" "vastdb_ldap" {
  use_auto_discovery             = false
  binddn                         = "cn=admin,dc=qa,dc=vastdata,dc=com"
  searchbase                     = "dc=qa,dc=vastdata,dc=com"
  bindpw                         = "vastdata"
  use_ldaps                      = false
  method                         = "simple"
  query_groups_mode              = "COMPATIBLE"
  use_tls                        = false
  urls                           = ["ldap://10.27.252.30"]
  is_vms_auth_provider           = false
  port                           = 389
  posix_attributes_source        = "JOINED_DOMAIN"
  reverse_lookup                 = false
  gid_number                     = "gidNumber"
  uid                            = "uid"
  uid_number                     = "uid_number"
  match_user                     = "uid"
  uid_member_value_property_name = "uid"
  uid_member                     = "memberUID"
  posix_account                  = "posixAccount"
  posix_group                    = "posixGroup"
  username_property_name         = "cn"
  user_login_name                = "uid"
  group_login_name               = "cn"
  mail_property_name             = "mail"
  monitor_action                 = "PING"
}
