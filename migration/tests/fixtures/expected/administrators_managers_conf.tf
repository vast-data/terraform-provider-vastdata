# Copyright (c) HashiCorp, Inc.

variable role_name {
    type = string
}

variable man_name {
    type = string
}

variable man_password {
    type = string
}

variable role_permissions_list {
    type = list(string)
}

variable permissions_list {
    type = list(string)
}

resource "vastdata_administrator_role" "man_role1" {
  name = var.role_name
  permissions_list = var.role_permissions_list
}

resource "vastdata_administrator_manager" "manager1" {
  username = var.man_name
  password = var.man_password
  roles = [vastdata_administrator_role.man_role1.id]
  permissions_list = var.permissions_list
}

output tf_role {
  value = vastdata_administrator_role.man_role1
}

output tf_man {
  value = vastdata_administrator_manager.manager1
  sensitive = true
}
