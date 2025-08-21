# Copyright (c) HashiCorp, Inc.

variable gss_name {
    type = string
}

resource vastdata_snapshot snapshot {
    name = "snap_1"
    path = "/"
    tenant_id = 1
    indestructible = false
    expiration_time = "2025-06-02T12:22:32Z"
}

resource vastdata_global_snapshot gsnap1 {
	name             = var.gss_name
	enabled          = true
	loanee_tenant_id = 1
	loanee_root_path = "/${var.gss_name}"

	#remote_target_guid = gss_tenant1.guid
	#remote_target_id   = vastdata_replication_peers.clusterA-clusterB-peer.id
	owner_root_snapshot {
		name = vastdata_snapshot.snapshot.name

	}
	owner_tenant {
		guid = 1
		name = "default"
	}
	lifecycle {
		ignore_changes = [remote_target_guid, remote_target_path, loanee_root_path]
	}
}

output tf_gss_snapshot {
  value = vastdata_snapshot.snapshot
}