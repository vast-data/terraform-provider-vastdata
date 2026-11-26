# ignore:e2e
# Schedule trigger with tags and broker accessibility flag.

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
}
