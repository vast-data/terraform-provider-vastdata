# Register a Data Engine function that points at an image already in a
# container registry (build/push is outside Terraform).
#
# Prerequisites: vastdata_de_engine, vastdata_de_container_registry.

resource "vastdata_de_function" "hello" {
  name                   = "de-test-function"
  tenant_name            = "de-lab" # required for cluster-admin auth
  container_registry_vrn = vastdata_de_container_registry.example.vrn
  artifact_source        = "de-lab-test/hello-world"
  image_tag              = "latest"
  revision_alias         = "latest"
  description            = "Hello World function"
}

# ---------------------
# Complete examples
# ---------------------

resource "vastdata_de_function" "tagged" {
  name                   = "de-function-tagged"
  tenant_name            = "de-lab"
  container_registry_vrn = vastdata_de_container_registry.example.vrn
  artifact_source        = "de-lab-test/hello-world"
  image_tag              = "latest"
  skip_image_validation  = true
  revision_alias         = "lab"
  revision_description   = "published lab revision"
  description            = "Tagged published function"
  is_published           = true
  tags                   = ["example", "tagged"]
}

# --------------------

resource "vastdata_de_function" "plain" {
  name                   = "de-function-plain"
  tenant_name            = "de-lab"
  container_registry_vrn = vastdata_de_container_registry.example.vrn
  artifact_source        = "de-lab-test/hello-world"
  image_tag              = "latest"
  skip_image_validation  = true
  description            = "Plain unpublished function"
  is_published           = false
}
