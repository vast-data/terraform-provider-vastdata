# Copyright (c) HashiCorp, Inc.

resource "vastdata_view" "artifactory-prod" {
    path       = "/service/artifactory-prod"
    create_dir = "true"
    policy_id  = vastdata_view_policy.services.id
    protocols  = ["NFS"]
}

resource "vastdata_view" "artifactory-backup-prod" {
    path       = "/service/artifactory-backup-prod"
    create_dir = "true"
    policy_id  = vastdata_view_policy.services.id
    protocols  = ["NFS"]
}

resource "vastdata_view" "tempo_cs_us_dev" {
  path         = "/service/grafana-tempo-us-cs-dev"
  create_dir   = "true"
  policy_id    = vastdata_view_policy.services.id
  protocols    = ["S3"]
  bucket       = "grafana-tempo-us-cs-dev"
  bucket_owner = "svcTempoCSUSDev@mavensecurities.com"
}

resource "vastdata_view" "tempo_cs_us_prod" {
  path         = "/service/grafana-tempo-us-cs-prod"
  create_dir   = "true"
  policy_id    = vastdata_view_policy.services.id
  protocols    = ["S3"]
  bucket       = "grafana-tempo-us-cs-prod"
  bucket_owner = "svcTempoCSUSProd@mavensecurities.com"
}

resource "vastdata_quota" "artifactory-prod" {
  name          = "artifactory-prod"
  default_email = "core.platform@mavensecurities.com"
  path          = vastdata_view.artifactory-prod.path
  // limits are in bytes 9Tb and 10Tb
  soft_limit    = 9000000000000
  hard_limit    = 10000000000000
  is_user_quota = false
}

resource "vastdata_quota" "artifactory-backup-prod" {
  name          = "artifactory-backup-prod"
  default_email = "core.platform@mavensecurities.com"
  path          = vastdata_view.artifactory-backup-prod.path
  // limits are in bytes 9Tb and 10Tb
  soft_limit    = 9000000000000
  hard_limit    = 10000000000000
  is_user_quota = false
}

resource "vastdata_quota" "tempo_cs_us_dev" {
  name          = "tempo-us-cs-dev"
  default_email = "core.platform@mavensecurities.com"
  path          = vastdata_view.tempo_cs_us_dev.path
  // 60GB soft limit, 100GB hard limit
  soft_limit = 60000000000
  hard_limit = 100000000000
}

resource "vastdata_quota" "tempo_cs_us_prod" {
  name          = "tempo-us-cs-prod"
  default_email = "core.platform@mavensecurities.com"
  path          = vastdata_view.tempo_cs_us_prod.path
  // 400GB soft limit, 500GB hard limit
  soft_limit = 400000000000
  hard_limit = 500000000000
}
