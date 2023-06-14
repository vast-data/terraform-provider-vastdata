#When there is only one quota with that name at the entire cluster
data vastdata_quota quota1 {
  name = "quota1"
}

#When there is more than one quota with that name at the cluster with differant tenant id

resource vastdata_tenant tenant1 {
 name = "tenant01"
 client_ip_ranges {
         start_ip = "192.168.0.100"
         end_ip = "192.168.0.200"
    }
}


data vastdata_quota quota3 {
        name = "quota3"
        tenant_id = vastdata_tenant.tenant1.id

}
