resource "vastdata_webhook" "vastdb_webhook" {
  name    = "vastdb-webhook"
  url     = "https://webhook.site/example-endpoint"
  method  = "POST"
  data    = jsonencode({ event = "storage.alert", severity = "high" })
  enabled = true
}
