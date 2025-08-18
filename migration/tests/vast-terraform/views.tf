# Copyright (c) HashiCorp, Inc.


## holds definitions of shared views that are likely to serve multiple purposes across multiple teams and serives

resource "vastdata_view" "data" {
  path       = "/data"
  policy_id  = vastdata_view_policy.data-worldread.id
  create_dir = "true"
  protocols  = ["NFS"]
}


resource "vastdata_view" "logs" {
  path       = "/logs"
  policy_id  = vastdata_view_policy.logs.id
  create_dir = "true"
  protocols  = ["NFS"]
}

resource "vastdata_view" "logs-system" {
  path       = "/logs/system"
  policy_id  = vastdata_view_policy.logs.id
  create_dir = "true"
  protocols  = ["NFS"]
}

resource "vastdata_view" "logs-maven-apps" {
  path       = "/logs/maven-apps"
  policy_id  = vastdata_view_policy.logs.id
  create_dir = "true"
  protocols  = ["NFS"]
}

resource "vastdata_view" "logs-maven-k8s" {
  path       = "/logs/k8s"
  policy_id  = vastdata_view_policy.logs.id
  create_dir = "true"
  protocols  = ["NFS"]
}

resource "vastdata_view" "logs-mm-app" {
  path       = "/logs/maven-apps/marketmaking"
  policy_id  = vastdata_view_policy.logs.id
  create_dir = "true"
  protocols  = ["NFS"]
}

resource "vastdata_quota" "logs" {
  name          = "logs"
  default_email = "michael.moyles@mavensecurities.com"
  path          = vastdata_view.logs.path
    // limits are in bytes
  soft_limit    = 400000000000
  hard_limit    = 500000000000
  is_user_quota = false
}




