# Copyright (c) HashiCorp, Inc.

variable machine_name {
    type = string
}

resource vastdata_active_directory2 active_dir3 {
    machine_account_name = var.machine_name
    organizational_unit = "OU=VASTs,OU=VastENG,DC=VastENG,DC=lab"
    use_auto_discovery = "false"
    binddn = "cn=admin,dc=qa,dc=vastdata,dc=com"
    searchbase = "dc=qa,dc=vastdata,dc=com"
    bindpw = "vastdata"
    use_ldaps = "false"
    domain_name = "VastEng.lab"
    method = "simple"
    query_groups_mode = "COMPATIBLE"
    use_tls = "false"
    urls = ["ldap://10.27.252.30"]
    is_vms_auth_provider = "true"
}

output tf_active_directory {
  value = vastdata_active_directory2.active_dir3
  sensitive = true
}
