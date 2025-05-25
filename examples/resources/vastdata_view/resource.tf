# Create a view with NFSv3 and NFSv4 protocols
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
