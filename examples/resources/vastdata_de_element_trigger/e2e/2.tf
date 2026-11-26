# ignore:e2e
# Minimal ObjectRemoved-only element trigger (no name/tag filters).

resource "vastdata_de_element_trigger" "removed_only" {
  name         = "de-element-removed"
  tenant_name  = "de-lab"
  topic_name   = "engine-main"
  description  = "ObjectRemoved-only element trigger"
  source_types = ["S3"]
  tags         = ["example", "removed"]

  broker = {
    name = "de-lab-broker"
    type = "Internal"
  }

  config = {
    source_bucket_name = "my-bucket"
    events             = ["ObjectRemoved:*"]
  }
}
