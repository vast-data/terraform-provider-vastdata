#Realms are avaliable from version 5.2.0

resource "vastdata_administators_realms" "realm01" {
  name         = "realm01"
  object_types = ["dnode", "nic", "viewpolicy"]
}

