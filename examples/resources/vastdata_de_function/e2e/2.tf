# ignore:e2e
# Minimal unpublished function (no tags / revision alias).

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
