# Copyright (c) HashiCorp, Inc.

variable machine_name {
    type = string
}

resource vastdata_ldap ldap1 {
  domain_name = "VastEng.lab"
  urls = ["ldap://10.27.252.30"]
  binddn = "cn=admin,dc=qa,dc=vastdata,dc=com"
  searchbase = "dc=qa,dc=vastdata,dc=com"
  bindpw = "vastdata"
  use_auto_discovery = "false"
  use_ldaps = "false"
  port = "389"
  method = "simple"
  query_groups_mode = "COMPATIBLE"
  use_tls = "false"
}

resource vastdata_active_directory active_dir1 {
  ldap_id = vastdata_ldap.ldap1.id
  machine_account_name = var.machine_name
  organizational_unit = "OU=VASTs,OU=VastENG,DC=VastENG,DC=lab"
}

output tf_active_directory {
  value = vastdata_active_directory.active_dir1
  sensitive = true
}

output tf_active_directory_ldap {
  value = vastdata_ldap.ldap1
  sensitive = true
}