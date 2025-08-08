resource "vastdata_nis" "vastdb_nis" {
  name        = "vastdb_nis"
  domain_name = "my.nis.domain.example.com"
  ips         = ["1.1.1.1", "2.2.2.2"]
}
