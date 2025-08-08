resource "vastdata_nis" "vastdb_nis" {
  name        = "vastdb_nis"
  domain_name = "my.nis.domain.example.com"
  ips         = ["1.1.1.1", "2.2.2.2"]
}

# ---------------------
# Complete examples
# ---------------------

resource "vastdata_nis" "vastdb_nis" {
  name        = "vastdb_nis"
  domain_name = "my.nis.domain.example.com"
  ips         = ["1.1.1.1", "2.2.2.2"]
}

# --------------------

resource "vastdata_nis" "vastdb_nis" {
  name        = "vastdb_nis"
  domain_name = "my.nis.domain.example.com"
  hosts       = ["my.domain.example.com", "my.domain2.example.com"]
  servers     = ["server1", "server2"]
  ips         = ["1.1.1.1", "2.2.2.2"]
}

# --------------------

