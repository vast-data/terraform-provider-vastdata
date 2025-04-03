resource "vastdata_non_local_user" "ExternalUser" {
    uid                 = 1097416930
    tenant_id           = 1
    allow_create_bucket = true
    allow_delete_bucket = false
    s3_policies_ids     = [
        1
    ]
}