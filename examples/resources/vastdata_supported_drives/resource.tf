resource "vastdata_supported_drives" "vastdb_drive" {
  model_name = "ModelX"
  drive_type = "ssd"

  raw_data = {
    name        = "SSD_MODEL_X"
    capacity_tb = 15
    hw_platform = ["platform-a", "platform-b"]

    fw = [
      {
        name   = "FW_2.0.0"
        latest = true
      }
    ]
  }
}
