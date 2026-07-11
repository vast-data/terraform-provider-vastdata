resource "vastdata_tenant_metric_label_values_bulk" "vastdb_tenant_metric_values_bulk" {
  # Tenant whose metric label values are managed in bulk.
  tenant_id = vastdata_tenant.example.id

  # Label keys must match vastdata_tenant_metric_labels.key values.
  values = {
    environment = "production"
    region      = "us-east-1"
  }
}

# ---------------------
# Complete examples
# ---------------------

resource "vastdata_tenant_metric_labels" "environment" {
  key           = "vastdb_environment"
  default_value = "staging"
  description   = "Deployment environment tag for tenant metrics."
}

resource "vastdata_tenant_metric_labels" "region" {
  key           = "vastdb_region"
  default_value = "unknown"
  description   = "Geographic region tag for tenant metrics."
}

resource "vastdata_tenant" "vastdb_tenant" {
  name         = "vastdb-tenant-mlvbulk"
  force_delete = true
}

resource "vastdata_tenant_metric_label_values_bulk" "vastdb_metric_label_values_bulk" {
  tenant_id = vastdata_tenant.vastdb_tenant.id
  values = {
    vastdb_environment = "production"
    vastdb_region      = "us-east-1"
  }
}

# --------------------

