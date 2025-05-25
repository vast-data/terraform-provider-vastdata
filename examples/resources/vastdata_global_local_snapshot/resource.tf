#Create a local global snapshot from a snapshot on the same tenant
resource "vastdata_tenant" "tenant" {
  name = "tenant1"
  client_ip_ranges {
    start_ip = "192.168.0.100"
    end_ip   = "192.168.0.200"
  }
}

resource "vastdata_view" "view" {
  path       = "/view1"
  policy_id  = vastdata_view_policy.view-policy.id
  tenant_id  = vastdata_tenant.tenant.id
  create_dir = "true"
}

resource "vastdata_snapshot" "snapshot" {
  name            = "snapshot1"
  path            = vastdata_view.view.path
  tenant_id       = vastdata_tenant.tenant.id
  indestructible  = false
  expiration_time = "2023-11-20T12:22:32Z"
  lifecycle {
    ignore_changes = [path]
  }

}

resource "vastdata_global_local_snapshot" "local_snapshot" {
  name               = "local_snapshot1"
  loanee_root_path   = "/local_snapshot1"
  loanee_snapshot_id = vastdata_snapshot.snapshot.id
  loanee_tenant_id   = vastdata_tenant.tenant.id
  owner_tenant {
    name = vastdata_tenant.tenant.name
    guid = vastdata_tenant.tenant.guid
  }

}
