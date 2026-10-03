# ignore:e2e
# Element trigger with name filters and tag filters.
# Tag filters require ObjectTagging:Put / ObjectTagging:Delete among events.

resource "vastdata_de_element_trigger" "with_filters" {
  name         = "de-element-filtered"
  tenant_name  = "de-lab"
  topic_name   = "engine-main"
  description  = "Element trigger with name and tag filters"
  source_types = ["S3"]
  tags         = ["example", "filtered"]

  broker = {
    name = "de-lab-broker"
    type = "Internal"
  }

  config = {
    source_bucket_name = "my-bucket"
    events = [
      "ObjectCreated:*",
      "ObjectTagging:Put",
      "ObjectTagging:Delete",
    ]
    name_filters = {
      prefixes = ["incoming/"]
      suffixes = [".json"]
    }
    tag_filters = {
      prefixes = ["env"]
      suffixes = ["prod"]
    }
  }
}
