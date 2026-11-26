# ignore:e2e
# Function with tags, published revision alias, and revision description.

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
