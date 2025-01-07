#### Mapping blockhosts to a volume #####
data "vastdata_view_policy" "default" {
  name = "default"
}

resource "vastdata_view" "blockA-view" {
  path                 = "/blockA"
  name                 = "blockA"
  is_default_subsystem = false
  policy_id            = data.vastdata_view_policy.default.id
  create_dir           = "true"
  protocols            = ["BLOCK"]
}

resource "vastdata_blockhost" "hostA" {
  name      = "hostA"
  tenant_id = 1
  nqn       = "nqn.2014-08.org.nvmexpress:ABCDEFGHIJKL"
}


resource "vastdata_volume" "volume01" {
  name           = "/volume1"
  size           = 150000000000
  view_id        = vastdata_view.blockA-view.id
  volume_tags    = ["key1:value1", "key2:value2"]
  block_host_ids = [vastdata_blockhost.hostA.id]
}

#### Mapping multiple blockhosts to a volume generate byt terraform count (will also work for for_each) #####
#this example mapps 50 blockhosts to the same volume

data "vastdata_view_policy" "default" {
  name = "default"
}

resource "vastdata_view" "blockA-view" {
  path                 = "/blockA"
  name                 = "blockA"
  is_default_subsystem = false
  policy_id            = data.vastdata_view_policy.default.id
  create_dir           = "true"
  protocols            = ["BLOCK"]
}

resource "vastdata_blockhost" "blockhost" {
  count = 50
  name  = "blockhost-${count.index}"
  nqn   = "nqn.2014-08.org.nvmexpress:BLCOKHOST_${count.index}"
}

resource "vastdata_volume" "volume01" {
  name           = "/volume1"
  size           = 150000000000
  view_id        = vastdata_view.blockA-view.id
  volume_tags    = ["key1:value1", "key2:value2"]
  block_host_ids = vastdata_blockhost.blockhost[*].id
}
