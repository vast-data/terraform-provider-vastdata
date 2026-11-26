resource "vastdata_api_token" "vastde" {
  count = var.create_api_token ? 1 : 0

  name        = local.api_token_name
  expiry_date = var.api_token_expiry
  # Leave null to create for the authenticated user. Some VMS builds reject
  # POST body field "owner" with {"owner":["Unexpected field."]}.
  owner = var.api_token_owner
}
