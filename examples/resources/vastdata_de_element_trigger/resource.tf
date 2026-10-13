# Create a Data Engine element trigger (bucket/object events).

resource "vastdata_de_element_trigger" "example" {
  name        = "de-element-trigger"
  tenant_name = "de-lab" # required for cluster-admin auth
  topic_name  = "engine-main"

  broker = {
    name = "de-lab-broker"
    type = "Internal"
  }

  source_types = ["S3"]

  config = {
    source_bucket_name = "my-bucket"
    events             = ["ObjectCreated:*", "ObjectRemoved:*"]
    name_filters = {
      prefixes = ["incoming/"]
      suffixes = [".json"]
    }
  }
}

# ---------------------
# Complete examples
# ---------------------

# Tag filters require ObjectTagging:Put / ObjectTagging:Delete among events.
resource "vastdata_de_element_trigger" "with_filters" {
  name         = "de-element-filtered"
  tenant_name  = "de-lab"
  topic_name   = "engine-main"
  description  = "Element trigger with name and tag filters"
  source_types = ["S3"]
  tags         = ["example", "filtered"]

  custom_extensions = {
    env     = "lab"
    purpose = "filtered"
  }

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

# --------------------

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
