#Define a nis with domain my.nis.domain.example.com , with nis servers 1.1.1.1 , 2.2.2.2
resource vastdata_nis nis1{
  domain_name = "my.nis.domain.example.com"
  hosts = ["1.1.1.1","2.2.2.2"]
}
