resource "vastdata_callhome_config" "vastdb_callhome_config" {
  id         = 1
  obfuscated = true
}

# ---------------------
# Complete examples
# ---------------------

# e2e: adopt Call Home config id=1 and set a non-secret field.
resource "vastdata_callhome_config" "vastdb_callhome_config" {
  id         = 1
  obfuscated = true
}

# --------------------

# e2e: proxy settings including write-only proxy_password (MD5 change detection).
resource "vastdata_callhome_config" "vastdb_callhome_config" {
  id             = 1
  proxy_host     = "proxy.e2e.local"
  proxy_port     = "3128"
  proxy_username = "e2e-proxy-user"
  proxy_password = "e2e-proxy-pw-A"
}

# --------------------

