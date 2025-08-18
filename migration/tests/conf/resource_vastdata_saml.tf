# Copyright (c) HashiCorp, Inc.

variable idp_entityid {
    type = string
}

variable idp_name {
    type = string
}

variable idp_metadata_url {
    type = string
}

variable vms_id {
    type = number
}

variable encrypt_assertion {
    type = bool
}

variable want_assertions_or_response_signed {
    type = bool
}

resource vastdata_saml saml1 {
  idp_entityid = var.idp_entityid
  idp_name = var.idp_name
  vms_id = var.vms_id
  idp_metadata_url = var.idp_metadata_url
  encrypt_assertion = var.encrypt_assertion
  want_assertions_or_response_signed = var.want_assertions_or_response_signed
}

output tf_saml {
  value = vastdata_saml.saml1
  sensitive = true
}