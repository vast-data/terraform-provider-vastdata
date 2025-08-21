# Copyright (c) HashiCorp, Inc.

# To refer to a specific system, need to add `provider = vastdata.system<IDX>` where
# IDX is the system's index (in the order of appearance in Comet's commandline)

#Define a nis with domain my.nis.domain.example.com , with nis servers 1.1.1.1 , 2.2.2.2

variable nis_domain_name {
    type = string
}

resource vastdata_nis nis1 {
  domain_name = var.nis_domain_name
  hosts = ["10.27.252.101", "10.27.252.102"]
}

output tf_nis {
  value = vastdata_nis.nis1
}