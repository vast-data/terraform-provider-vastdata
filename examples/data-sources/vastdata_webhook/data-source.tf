
data "vastdata_webhook" "vastdb_webhook_by_url" {
  url = "https://webhook.site/example-endpoint"
}

data "vastdata_webhook" "vastdb_webhook_by_name" {
  name = "vastdb-webhook"
}
