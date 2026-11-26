# Manages Data Engine on the tenant (POST/GET/DELETE /serverless/data-engine —
# not deprecated /setup-provisioning).
# Works with tenant-admin (provider tenant) or cluster admin + tenant_name query.
resource "vastdata_de_engine" "this" {
  count = var.enable_data_engine ? 1 : 0

  tenant_name            = vastdata_tenant.dataengine.name
  default_topic_name     = vastdata_topic.main.name
  dead_letter_topic_name = vastdata_topic.deadletter.name

  kafka_broker = {
    name = vastdata_view.broker.bucket
    type = "Internal"
  }

  depends_on = [
    vastdata_topic.main,
    vastdata_topic.deadletter,
    vastdata_vip_pool.dataengine,
  ]
}
