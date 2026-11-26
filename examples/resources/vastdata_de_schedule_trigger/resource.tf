# Create a Data Engine schedule trigger (quartz cron).
# Runs every minute in this example: "0 * * ? * * *"

resource "vastdata_de_schedule_trigger" "hello" {
  name        = "de-test-trigger"
  tenant_name = "de-lab" # required for cluster-admin auth
  topic_name  = "engine-main"

  broker = {
    name = "de-lab-broker"
    type = "Internal"
  }

  config = {
    cron_schedule = "0 * * ? * * *"
  }
}

# ---------------------
# Complete examples
# ---------------------

resource "vastdata_de_schedule_trigger" "hourly" {
  name        = "de-sched-hourly"
  tenant_name = "de-lab"
  topic_name  = "engine-main"
  description = "Hourly schedule trigger"
  tags        = ["example", "hourly"]

  broker = {
    name                         = "de-lab-broker"
    type                         = "Internal"
    is_accessible_to_data_engine = true
  }

  config = {
    cron_schedule = "0 0 * ? * * *"
  }

  custom_extensions = {
    env     = "lab"
    purpose = "hourly"
  }
}

# --------------------

resource "vastdata_de_schedule_trigger" "daily" {
  name        = "de-sched-daily"
  tenant_name = "de-lab"
  topic_name  = "engine-main"
  description = "Daily schedule trigger"
  tags        = ["example", "daily"]

  broker = {
    name                         = "de-lab-broker"
    type                         = "Internal"
    is_accessible_to_data_engine = false
  }

  config = {
    cron_schedule = "0 0 0 ? * * *"
  }
}
