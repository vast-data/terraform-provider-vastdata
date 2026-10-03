terraform {
  required_version = ">= 1.5.0"

  required_providers {
    vastdata = {
      source  = "vast-data/vastdata"
      version = ">= 2.0.0"
    }
  }
}
