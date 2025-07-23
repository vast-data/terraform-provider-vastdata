resource "vastdata_kafka_broker" "vastdb_kafka_broker" {
  name = "vastdb_kafka_broker"
  addresses = [
    {
      host = "10.131.21.121"
      port = 31485
    },
    {
      host = "10.131.21.121"
      port = 31486
    }
  ]
}
