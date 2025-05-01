resource "vastdata_non_local_group" "ExternalGroup" {
    gid                 = 10000
    tenant_id           = 1
    s3_policies_ids     = [
        1
    ]
}
