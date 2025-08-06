data "vastdata_vms_configured_idps" "vastdb_vms_configured_idps" {
  vms_id = 1

}

resource "vastdata_saml_config" "vastdb_saml_config" {
  vms_id   = data.vastdata_vms_configured_idps.vastdb_vms_configured_idps.vms_id
  idp_name = data.vastdata_vms_configured_idps.vastdb_vms_configured_idps.idps[0]
  saml_settings = {
    encrypt_assertion                  = true
    force_authn                        = true
    idp_entityid                       = "https://my-idp.com/entry"
    idp_metadata_url                   = "https://my-idp.com/metadata"
    want_assertions_or_response_signed = true
  }
}
