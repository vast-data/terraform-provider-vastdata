#Create DNS with ip 11.0.0.1 for domain mu.example.com
resource "vastdata_dns" "dns1" {
  name          = "dns1"
  vip           = "11.0.0.1"
  domain_suffix = "my.example.com"
}
