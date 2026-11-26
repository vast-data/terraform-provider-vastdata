# kubernetes/helm stat config_path while Terraform configures providers, which
# is before the module's finalize step writes the kubeconfig. The module
# returns null until that file exists, so the providers stay unconfigured and
# the RKE2 bootstrap can apply. VAST CSI resources need a live API server, so
# enable them only after out/kubeconfig exists and run apply again.

provider "kubernetes" {
  config_path = module.rke2.kubernetes_config_path
}

provider "helm" {
  kubernetes = {
    config_path = module.rke2.kubernetes_config_path
  }
}
