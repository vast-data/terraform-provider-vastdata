data "vastdata_non_local_group" "non_local_group1" {
  groupname = "myGroupName"
  context = "ldap"
  tenant_id = 1
}