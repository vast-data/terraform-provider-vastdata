#Create view policy with the name vpolicy1 
resource "vastdata_vip_pool" "pool3" {
  name        = "pool3"
  role        = "PROTOCOLS"
  subnet_cidr = "24"
  ip_ranges {
    end_ip   = "11.0.0.80"
    start_ip = "11.0.0.70"
  }

}


resource "vastresource" "vastdata_tenant" "tenant1" {
  name = "tenant01"
  client_ip_ranges {
    start_ip = "192.168.0.100"
    end_ip   = "192.168.0.200"
  }
}


resource "vastdata_view_policy" "vpolicy1" {
  name          = "vpolicy1"
  vip_pools     = [vastdata_vip_pool.pool3.id]
  tenant_id     = vastdata_tenant.tenant1.id
  flavor        = "NFS"
  nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
}

#Definning vippools#
#the attribute vip_pools will work on vast version 5.0.0 and greater.
#however when working with version greater than 5.1.0 is is advisable wither to migrate using vippool_permissions if previously vip_pools was used on cluster from version 5.0.0.
#Important things to notice is that when using vip_pools on vast versions starting from version 5.1.0 the vippool permissions will always be set to RW.

#Example one using vip_pools versions 5.0.0 and above#

resource "vastdata_vip_pool" "pool4" {
  name        = "pool4"
  role        = "PROTOCOLS"
  subnet_cidr = "24"
  ip_ranges {
    end_ip   = "11.0.0.80"
    start_ip = "11.0.0.70"
  }

}

resource "vastdata_vip_pool" "pool5" {
  name        = "pool5"
  role        = "PROTOCOLS"
  subnet_cidr = "24"
  ip_ranges {
    end_ip   = "11.0.0.110"
    start_ip = "11.0.0.90"
  }

}


resource "vastresource" "vastdata_tenant" "tenant2" {
  name = "tenant02"
  client_ip_ranges {
    start_ip = "192.168.0.100"
    end_ip   = "192.168.0.200"
  }
  vippool_ids = [vastdata_vip_pool.pool4.id,vastdata_vip_pool.pool5.id]

}


resource "vastdata_view_policy" "vpolicy1" {
  name          = "vpolicy1"
  vip_pools     = [vastdata_vip_pool.pool3.id]
  tenant_id     = vastdata_tenant.tenant1.id
  flavor        = "NFS"
  nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
  vip_pools =  [vastdata_vip_pool.pool4.id,vastdata_vip_pool.pool5.id]
}

#Example one using vippool_permissions for vast versions 5.1.0 and above, while setting view policy vippool permissions setting one vippool with read/write permissions and the other one with read only#

resource "vastdata_vip_pool" "pool4" {
  name        = "pool4"
  role        = "PROTOCOLS"
  subnet_cidr = "24"
  ip_ranges {
    end_ip   = "11.0.0.80"
    start_ip = "11.0.0.70"
  }

}

resource "vastdata_vip_pool" "pool5" {
  name        = "pool5"
  role        = "PROTOCOLS"
  subnet_cidr = "24"
  ip_ranges {
    end_ip   = "11.0.0.110"
    start_ip = "11.0.0.90"
  }

}


resource "vastresource" "vastdata_tenant" "tenant2"{
  name = "tenant02"
  client_ip_ranges {
    start_ip = "192.168.0.100"
    end_ip   = "192.168.0.200"
  }
  vippool_ids = [vastdata_vip_pool.pool4.id,vastdata_vip_pool.pool5.id]

}


resource "vastdata_view_policy" "vpolicy1" {
  name          = "vpolicy1"
  vip_pools     = [vastdata_vip_pool.pool3.id]
  tenant_id     = vastdata_tenant.tenant1.id
  flavor        = "NFS"
  nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
  vippool_permissions {
    vippool_id = vastdata_vip_pool.pool4.id
    vippool_permissions  = "RW"
  }
  vippool_permissions {
    vippool_id = vastdata_vip_pool.pool5.id
    vippool_permissions  = "RO"
  }  
}


