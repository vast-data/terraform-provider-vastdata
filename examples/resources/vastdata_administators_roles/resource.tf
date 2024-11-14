#Basic Role creation without referring to realms (realms are avaliable from version 5.2.0).
resource "vastdata_administators_roles" "role01" {
  name             = "role01"
  permissions_list = ["create_support", "create_settings", "create_security", "create_monitoring", "create_logical", "create_hardware"]
}

#Since version 5.2.0 reamls are supported and can be referanced by a role
#you can specify 4 types of actions refering to realms create,delete,view,edit
#in order to cofigure the action specify <action>_<realm name>.
#Ex: if the realm name is realm01 in order to configure edit for this realm add the following to the permissions_list attribute edit_realm01

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


