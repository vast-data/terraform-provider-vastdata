resource "vastdata_non_local_group" "ExternalGroup" {
    gid                 = 10000
    s3_policies_ids     = [
        1
    ]
}
