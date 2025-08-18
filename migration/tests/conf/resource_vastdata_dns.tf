# Copyright (c) HashiCorp, Inc.

# To refer to a specific system, need to add `provider = vastdata.system<IDX>` where
# IDX is the system's index (in the order of appearance in Comet's commandline)

variable dns_name {
    type = string
}

variable dns_net_type {
    type = string
}

variable dns_invalid_name_response {
    type = string
}

variable dns_invalid_type_response {
    type = string
}

resource vastdata_dns dns1 {
  name = var.dns_name
  vip = "11.0.0.1"
  domain_suffix = "my.example.com"
  net_type = var.dns_net_type
  invalid_name_response = var.dns_invalid_name_response
  invalid_type_response = var.dns_invalid_type_response
  ttl = 1900
}

output tf_dns {
  value = vastdata_dns.dns1
}
