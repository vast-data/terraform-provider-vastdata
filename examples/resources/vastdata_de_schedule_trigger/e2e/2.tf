# ignore:e2e
# Daily cron; broker accessibility checks disabled.

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
