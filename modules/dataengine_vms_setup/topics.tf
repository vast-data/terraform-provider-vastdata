resource "vastdata_topic" "main" {
  database_name    = vastdata_view.broker.bucket
  name             = var.main_topic_name
  topic_partitions = var.topic_partitions
  retention_ms     = var.topic_retention_ms
  tenant_id        = vastdata_tenant.dataengine.id
}

resource "vastdata_topic" "deadletter" {
  database_name    = vastdata_view.broker.bucket
  name             = var.deadletter_topic_name
  topic_partitions = var.topic_partitions
  retention_ms     = var.topic_retention_ms
  tenant_id        = vastdata_tenant.dataengine.id
}
