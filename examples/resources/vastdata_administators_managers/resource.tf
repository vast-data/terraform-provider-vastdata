resource "vastdata_administators_roles" "roleC" {
  name             = "rolec"
  permissions_list = ["create_support", "create_settings", "create_security", "create_logical", "create_hardware"]
}

resource "vastdata_administators_managers" "managerC" {
  username         = "managerc"
  password         = "<some password>"
  roles            = [vastdata_administators_roles.roleB.id]
  permissions_list = ["create_monitoring"]
}

#Starting with version 5.2.0, a role can reference one or more realms. You can specify 4 types of actions referring to realms: create, delete, view, edit. To configure the action, specify <action>_<realm name>.
#For example, to configure edit for a realm named `realm01`, add the following to the `permissions_list` attribute: `edit_realm01`

resource "vastdata_administators_realms" "realmc" {
  name         = "realmc"
  object_types = ["nic", "viewpolicy"]
}

resource "vastdata_administators_roles" "rolec" {
  name             = "rolec"
  permissions_list = ["create_support", "create_settings", "create_security", "create_logical", "create_hardware"]
}

resource "vastdata_administators_managers" "managerc" {
  username = "managerc"
  password = "<some password>"
  roles    = [vastdata_administators_roles.rolec.id]
  permissions_list = [
    "create_monitoring",
    "edit_${vastdata_administators_realms.realmc.name}",
    "view_${vastdata_administators_realms.realmc.name}"
  ]
}

