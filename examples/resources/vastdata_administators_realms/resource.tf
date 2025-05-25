#Realms are avaliable in version 5.2.0 and later

resource "vastdata_administators_realms" "realm01" {
  name         = "realm01"
  object_types = ["dnode", "nic", "viewpolicy"]
}

