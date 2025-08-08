data "vastdata_nonlocal_user" "vastdb_nonlocal_user_by_id" {
  id = 1
}

data "vastdata_nonlocal_user" "vastdb_nonlocal_user_by_uid" {
  uid = 15959802
}

data "vastdata_nonlocal_user" "vastdb_nonlocal_user_by_username" {
  username  = "ndb.user"
  tenant_id = 1
}


data "vastdata_nonlocal_user" "vastdb_nonlocal_user_by_sid" {
  sid     = "S-1-111-2410027339-2891882760-1490074538-1286424047-4"
  context = "ldap"
}
