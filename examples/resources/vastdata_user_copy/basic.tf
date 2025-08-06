# Example: Copy users by tenant ID
resource "vastdata_user_copy" "copy_by_tenant" {
  destination_provider_id = 123
  tenant_id               = 456
}

# Example: Copy specific users by their IDs
resource "vastdata_user_copy" "copy_specific_users" {
  destination_provider_id = 123
  user_ids                = [1, 2, 3, 4, 5]
}