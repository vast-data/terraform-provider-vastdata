# Copyright (c) HashiCorp, Inc.

variable idp_name {
    type = string
}

variable vms_id {
    type = number
}

data "vastdata_saml" "saml_ds1" {
    idp_name = var.idp_name
    vms_id = var.vms_id
}

output tf_saml_ds {
    value = data.vastdata_saml.saml_ds1
    sensitive = true
}
