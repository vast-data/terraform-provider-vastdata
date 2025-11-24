resource "vastdata_rack" "vastdb_rack" {
  name = "rack-1"
}

# ---------------------
# Complete examples
# ---------------------

resource "vastdata_rack" "vastdb_rack" {
  name        = "rack-2"
  ip_range    = ["1.1.1.1", "1.1.1.10"]
  description = "This is test rack"
}

# --------------------

