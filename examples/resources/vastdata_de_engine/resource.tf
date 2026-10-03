# Manage Data Engine for a tenant via POST/GET/DELETE /serverless/data-engine/
# (not the deprecated /setup-provisioning endpoint).
#
# Tenant-admin: set provider.tenant (or omit tenant_name when already scoped).
# Cluster admin: pass tenant_name on the resource (API query).
#
# Prerequisites (e.g. modules/dataengine_vms_setup): tenant, VIP pool,
# S3-native broker view/bucket, and Kafka topics.

provider "vastdata" {
  host     = "vast-cluster"
  username = "admin"
  password = "******"
  # tenant = "de-lab"  # optional for tenant-admin auth
}

resource "vastdata_de_engine" "this" {
  tenant_name            = "de-lab" # required for cluster-admin auth
  default_topic_name     = "engine-main"
  dead_letter_topic_name = "engine-deadletter"

  kafka_broker = {
    name = "de-lab-broker" # Internal: broker view/bucket name
    type = "Internal"
  }
}
