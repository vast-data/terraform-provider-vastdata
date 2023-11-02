#Create a quota with user & group quota
resource "vastdata_view_policy" "example" {
  name   = "example"
  flavor = "NFS"
}

resource "vastdata_view" "quotas2-view" {
  path       = "/quota2"
  policy_id  = vastdata_view_policy.example.id
  create_dir = "true"
  protocols  = ["NFS", "NFS4"]
}


resource "vastdata_quota" "quota2" {
  name          = "quota2"
  default_email = "user@example.com"
  path          = vastdata_view.quotas2-view.path
  soft_limit    = 100000
  hard_limit    = 100000
  is_user_quota = true
  default_user_quota {
    grace_period      = "02:00:00"
    hard_limit        = 2000
    soft_limit        = 1000
    hard_limit_inodes = 20000000
  }
  user_quotas {
    grace_period = "01:00:00"
    hard_limit   = 15000
    soft_limit   = 15000
    entity {
      name            = "user1"
      email           = "user1@example.com"
      identifier      = "user1"
      identifier_type = "username"
      is_group        = "false"
    }
  }
  group_quotas {
    grace_period = "01:00:00"
    hard_limit   = 15000
    soft_limit   = 15000
    entity {
      name            = "group1"
      identifier      = "group1"
      identifier_type = "group"
      is_group        = "false"
    }
  }

}
