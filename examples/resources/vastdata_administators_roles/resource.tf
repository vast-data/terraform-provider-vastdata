#Basic Role creation without referring to realms (realms are avaliable with version 5.2.0 and later)
resource "vastdata_administators_roles" "role01" {
  name             = "role01"
  permissions_list = ["create_support", "create_settings", "create_security", "create_monitoring", "create_logical", "create_hardware"]
}

#Starting with version 5.2.0, a role can reference one or more realms. You can specify 4 types of actions referring to realms: create, delete, view, edit. To configure the action, specify <action>_<realm name>.
#For example, to configure edit for a realm named `realm01`, add the following to the `permissions_list` attribute: `edit_realm01`

resource "vastdata_administators_realms" "realmc" {
  name         = "realmc"
  object_types = ["nic", "viewpolicy"]
}


resource "vastdata_administators_roles" "rolec" {
  name = "rolec"
  permissions_list = ["create_support",
    "create_settings",
    "edit_${vastdata_administators_realms.realmc.name}",
    "view_${vastdata_administators_realms.realmc.name}"
  ]
}


