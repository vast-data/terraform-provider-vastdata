resource "vastdata_view_policy" "broker" {
  name                 = local.view_policy_name
  tenant_id            = vastdata_tenant.dataengine.id
  flavor               = "S3_NATIVE"
  auth_source          = var.s3_view_policy_auth_source
  s3_visibility_groups = [local.app_group_name]
}

resource "vastdata_view" "broker" {
  name         = local.broker_bucket
  path         = var.broker_view_path
  create_dir   = true
  delete_dir   = var.delete_dir
  policy_id    = vastdata_view_policy.broker.id
  tenant_id    = vastdata_tenant.dataengine.id
  bucket       = local.broker_bucket
  bucket_owner = local.broker_owner
  protocols    = ["S3", "DATABASE", "KAFKA"]

  kafka_vip_pools                   = [vastdata_vip_pool.dataengine.id]
  is_kafka_unencrypted_conn_allowed = true
  kafka_unencrypted_auth_mechanism  = var.kafka_unencrypted_auth_mechanism
  kafka_is_authorization_required   = var.kafka_is_authorization_required
}
