resource "null_resource" "vast_csi_snapshotter" {
  count = var.enable_vast_csi && var.vast_csi_install_snapshotter ? 1 : 0

  depends_on = [
    null_resource.finalize,
    null_resource.local_path_storage,
  ]

  triggers = {
    kubeconfig = local.kubeconfig_abs
    version    = var.vast_csi_snapshotter_version
  }

  provisioner "local-exec" {
    command     = abspath("${path.module}/scripts/install_snapshotter.sh")
    interpreter = ["bash"]
    environment = {
      KUBECONFIG          = local.kubeconfig_abs
      SNAPSHOTTER_VERSION = var.vast_csi_snapshotter_version
    }
  }
}

resource "kubernetes_namespace_v1" "vast_csi" {
  count = var.enable_vast_csi ? 1 : 0

  depends_on = [
    null_resource.finalize,
    null_resource.vast_csi_snapshotter,
  ]

  metadata {
    name = var.vast_csi_namespace
  }
}

resource "kubernetes_secret_v1" "vast_csi" {
  count = var.enable_vast_csi ? 1 : 0

  depends_on = [kubernetes_namespace_v1.vast_csi]

  metadata {
    name      = var.vast_csi_secret_name
    namespace = var.vast_csi_namespace
  }

  type = "Opaque"

  data = {
    endpoint = var.vast_csi_endpoint
    username = var.vast_csi_username
    password = var.vast_csi_password
    tenant   = var.vast_csi_tenant
  }
}

resource "helm_release" "vast_csi" {
  count = var.enable_vast_csi ? 1 : 0

  depends_on = [
    kubernetes_secret_v1.vast_csi,
    null_resource.vast_csi_snapshotter,
    null_resource.local_path_storage,
  ]

  name       = var.vast_csi_release_name
  repository = "https://vast-data.github.io/vast-csi"
  chart      = "vastcsi"
  version    = var.vast_csi_chart_version
  namespace  = var.vast_csi_namespace
  wait       = true
  timeout    = 600

  values = [
    yamlencode({
      truncateVolumeName = 128
      storageClasses = {
        (var.vast_csi_storage_class_name) = {
          vipPool                   = var.vast_csi_vip_pool
          storagePath               = var.vast_csi_storage_path
          viewPolicy                = var.vast_csi_view_policy
          volumeNameFormat          = "csi:{namespace}:{name}:{id}"
          ephemeralVolumeNameFormat = "csi:{namespace}:{name}:{id}"
          mountOptions = [
            "vers=3",
            "proto=tcp",
            "nconnect=16",
          ]
          secretName             = var.vast_csi_secret_name
          secretNamespace        = var.vast_csi_namespace
          setDefaultStorageClass = local.vast_csi_is_default
        }
      }
      controller = {
        runOnMaster       = false
        runOnControlPlane = false
        nodeSelector = {
          "node-role.kubernetes.io/control-plane" = "true"
        }
      }
    })
  ]
}
