# Copyright (c) HashiCorp, Inc.

variable role_name {
    type = string
}

variable permissions_list {
    type = list(string)
}

variable ldap_groups {
    type = list(string)
}

resource "vastdata_administators_roles" "role1" {
  name = var.role_name
  permissions_list = var.permissions_list
  ldap_groups = var.ldap_groups
}

output tf_role {
  value = vastdata_administators_roles.role1
}
