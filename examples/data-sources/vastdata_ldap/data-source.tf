data "vastdata_ldap" "vastdb_ldap_by_id" {
  id = 1
}

data "vastdata_ldap" "vastdb_ldap_by_guid" {
  guid = "00000000-0000-0000-0000-000000000001"
}

data "vastdata_ldap" "vastdb_ldap_by_domain_name" {
  domain_name = "VastEng.lab"
}

data "vastdata_ldap" "vastdb_ldap_by_name" {
  name = "ldap1"
}
