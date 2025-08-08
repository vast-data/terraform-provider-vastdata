data "vastdata_nonlocal_group" "nonlocal_group_by_name" {
  groupname = "QA-GRP-0077"
  context   = "ldap"
  tenant_id = 2
}

data "vastdata_nonlocal_group" "nonlocal_group_by_gid" {
  gid = 4
}

data "vastdata_nonlocal_group" "nonlocal_group_by_sid" {
  sid = "S-1-111-2410027339-2891882760-1490074538-1286424047-118"
}
