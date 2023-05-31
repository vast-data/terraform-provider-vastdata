#When there is only one snapshot with that name at the entire cluster
data vastdata_snapshot snapshot1 {
  name = "snapshot01"
}

#When there is more than one snapshot with that name at the cluster with differant tenant id

resource vastdata_tenant tenant1 {
 name = "tenant01"
 client_ip_ranges {
         start_ip = "192.168.0.100"
         end_ip = "192.168.0.200"
    }
}


data vastdata_snapshot snapshot1 {
        name = "snapshot01"
        tenant_id = vastdata_tenant.tenant1.id

}
