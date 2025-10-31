resource "vastdata_supported_drivers" "example" {
  model_name = "ExampleModel"
  version    = "1.0.0"
  replace    = false

  drives = [
    {
      name        = "SSD_MODEL_X"
      type        = "SSD"
      model       = "ModelX"
      capacity_tb = 15
      hw_platform = ["platform-a", "platform-b"]

      fw = [
        {
          name   = "FW_1.0.0"
          latest = true
        }
      ]
    }
  ]

  platforms = [
    {
      name                   = "platform-a"
      ssd_capacity_tb        = 15
      default_section_layout = "default"
    }
  ]
}

