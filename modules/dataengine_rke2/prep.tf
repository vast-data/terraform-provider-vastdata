resource "null_resource" "prep" {
  for_each = var.prepare_nodes ? local.all_nodes_map : {}

  triggers = {
    fingerprint = local.install_fingerprint
    upgrade     = tostring(var.apt_upgrade)
    ping_vip    = var.ping_vip != null ? var.ping_vip : ""
  }

  connection {
    type        = "ssh"
    host        = each.value.management_ip
    user        = var.ssh_user
    port        = var.ssh_port
    private_key = local.ssh_private_key
    timeout     = "5m"
  }

  provisioner "remote-exec" {
    inline = concat(
      [
        "set -eu",
        "sudo apt-get update -y",
      ],
      var.apt_upgrade ? ["sudo DEBIAN_FRONTEND=noninteractive apt-get upgrade -y"] : [],
      [
        "sudo DEBIAN_FRONTEND=noninteractive apt-get install -y nfs-common curl",
      ],
      var.ping_vip != null ? [
        "ping -c 4 ${var.ping_vip} | grep 'packet loss'"
      ] : []
    )
  }
}
