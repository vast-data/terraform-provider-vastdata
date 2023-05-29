#Create a tenant with the name tenent01 with client_ip_ranges
resource vastdata_tenant tenant1 {
 name = "tenant01"
 client_ip_ranges {
         start_ip = "192.168.0.100"
         end_ip = "192.168.0.200"
    }
}
