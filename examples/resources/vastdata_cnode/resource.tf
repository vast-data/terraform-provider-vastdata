resource "vastdata_cnode" "vastdb_cnode" {
  ip         = "10.0.1.100"
  cores      = 16
  cluster_id = 1
  force      = false
}
