# Copyright (c) HashiCorp, Inc.

variable realm_name {
    type = string
}

variable object_types {
    type = list(string)
}

resource "vastdata_administators_realms" "realm1" {
  name = var.realm_name
  object_types = var.object_types
}

output tf_realm {
  value = vastdata_administators_realms.realm1
}
