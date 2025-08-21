# Copyright (c) HashiCorp, Inc.

resource "vastdata_view" "instrument-reference-data-cache-chi-dev" {
    path       = "/service/instrument-reference-data-cache-chi-dev"
    bucket     = "instrument-reference-data-cache-chi-dev"
    share      = "instrument-reference-data-cache-chi-dev"
    create_dir = "true"
    policy_id  = vastdata_view_policy.services.id
    protocols  = ["NFS","S3","SMB"]
    bucket_owner = "svcInsRefDataCacheD@mavensecurities.com"
}

resource "vastdata_view" "instrument-reference-data-cache-chi" {
    path       = "/service/instrument-reference-data-cache-chi"
    bucket     = "instrument-reference-data-cache-chi"
    share      = "instrument-reference-data-cache-chi"
    create_dir = "true"
    policy_id  = vastdata_view_policy.services.id
    protocols  = ["NFS","S3","SMB"]
    bucket_owner = "svcInsRefDataCache@mavensecurities.com"
}

resource "vastdata_quota" "instrument-reference-data-cache-chi-dev" {
  name          = "instrument-reference-data-cache-chi-dev"
  default_email = "coredata@mavensecurities.com"
  path          = vastdata_view.instrument-reference-data-cache-chi-dev.path
  // limits are in bytes 150Gb and 200Gb
  soft_limit    = 150000000000
  hard_limit    = 200000000000
  is_user_quota = false
}

resource "vastdata_quota" "instrument-reference-data-cache-chi" {
  name          = "instrument-reference-data-cache-chi"
  default_email = "coredata@mavensecurities.com"
  path          = vastdata_view.instrument-reference-data-cache-chi.path
  // limits are in bytes 150Gb and 200Gb
  soft_limit    = 150000000000
  hard_limit    = 200000000000
  is_user_quota = false
}
