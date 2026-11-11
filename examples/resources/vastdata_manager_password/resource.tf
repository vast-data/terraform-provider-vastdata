resource "vastdata_manager_password" "example" {
  password = "new_secure_password_123!"
}

# ---------------------
# Complete examples
# ---------------------

# e2e: rotate current manager password (write-only + password_md5 RequiresReplace).
# Password must match VASTDATA_PASSWORD used by the e2e harness or later applies fail.
resource "vastdata_manager_password" "vastdb_manager_password" {
  password = "Www##12345678"
}

# --------------------

