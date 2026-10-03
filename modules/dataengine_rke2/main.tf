resource "null_resource" "first_server" {
  depends_on = [null_resource.prep]

  triggers = {
    fingerprint = local.install_fingerprint
    role        = "first"
    host        = local.first_cp.management_ip
  }

  connection {
    type        = "ssh"
    host        = local.first_cp.management_ip
    user        = var.ssh_user
    port        = var.ssh_port
    private_key = local.ssh_private_key
    timeout     = "10m"
  }

  provisioner "file" {
    source      = "${path.module}/scripts/bootstrap_rke2.sh"
    destination = "/tmp/bootstrap_rke2.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "set -eu",
      "chmod 755 /tmp/bootstrap_rke2.sh",
      "sudo env RKE2_CHANNEL='${var.rke2_channel}' RKE2_TOKEN='${var.cluster_token}' RKE2_WRITE_KUBECONFIG='${var.write_kubeconfig_mode}' RKE2_EXTRA_TLS_SANS='${join(" ", var.extra_tls_sans)}' /tmp/bootstrap_rke2.sh first",
      # Wait until the local API answers before joining peers.
      "for i in $(seq 1 90); do sudo /var/lib/rancher/rke2/bin/kubectl --kubeconfig /etc/rancher/rke2/rke2.yaml get nodes >/dev/null 2>&1 && exit 0; sleep 10; done",
      "echo 'Timed out waiting for first RKE2 server API' >&2",
      "exit 1",
    ]
  }
}

# Additional control planes must join one at a time (etcd).
resource "null_resource" "additional_servers" {
  count = length(local.additional_control_planes) > 0 ? 1 : 0

  depends_on = [null_resource.first_server]

  triggers = {
    fingerprint = local.install_fingerprint
    role        = "next"
    hosts       = join(",", [for n in local.additional_control_planes : n.management_ip])
  }

  provisioner "local-exec" {
    command     = abspath("${path.module}/scripts/join_additional_servers.sh")
    interpreter = ["bash"]
    environment = merge(local.ssh_env, {
      JOIN_ADDRESS     = local.join_address
      FIRST_HOST       = local.first_cp.management_ip
      ADDITIONAL_HOSTS = join(" ", [for n in local.additional_control_planes : n.management_ip])
      EXPECTED_CPS     = tostring(length(var.control_plane_nodes))
    })
  }
}

# One resource so workers join one at a time. for_each would download the
# RKE2 tarball from GitHub on every worker in parallel.
resource "null_resource" "agents" {
  count = length(var.worker_nodes) > 0 ? 1 : 0

  depends_on = [
    null_resource.first_server,
    null_resource.additional_servers,
  ]

  triggers = {
    fingerprint = local.install_fingerprint
    role        = "agent"
    hosts       = join(",", [for n in var.worker_nodes : n.management_ip])
  }

  provisioner "local-exec" {
    command     = abspath("${path.module}/scripts/join_agent.sh")
    interpreter = ["bash"]
    environment = merge(local.ssh_env, {
      AGENT_HOSTS  = join(" ", [for n in var.worker_nodes : n.management_ip])
      JOIN_ADDRESS = local.join_address
    })
  }
}
