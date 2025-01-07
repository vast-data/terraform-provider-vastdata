# Create a view with NFS & NFSv4 protocols
resource "vastdata_view_policy" "example" {
  name   = "example"
  flavor = "NFS"
}


resource "vastdata_view" "example-view" {
  path       = "/example"
  policy_id  = vastdata_view_policy.example.id
  create_dir = "true"
  protocols  = ["NFS", "NFS4"]
}

#Creating a KAFKA type view#
#When creating a view from the type of KAKFA protocols should also include DATABASE & S3 and kafka_vip_pools should be defined with at least one vippool id.
data "vastdata_view_policy" "default_s3_policy" {
  name = "s3_default_policy"
}

resource "vastdata_user" "example-user" {
  name = "user1"
  uid  = 9000
  allow_create_bucket = true
  allow_delete_bucket = true
  s3_superuser = true
}

resource "vastdata_vip_pool" "pool1" {
  name        = "pool1"
  role        = "PROTOCOLS"
  subnet_cidr = "24"
  ip_ranges {
    end_ip   = "11.0.0.40"
    start_ip = "11.0.0.20"
  }

  ip_ranges {
    start_ip = "11.0.0.5"
    end_ip   = "11.0.0.10"
  }
}

resource "vastdata_view" "kafka-view" {
  path       = "/kafka-view"
  policy_id  = data.vastdata_view_policy.default_s3_policy.id
  create_dir = "true"
  tenant_id = 1
  bucket = "kafkabucket"
  bucket_owner = vastdata_user.example-user.name
  kafka_vip_pools = [vastdata_vip_pool.pool1.id]
  protocols  = ["KAFKA","S3","DATABASE"]
