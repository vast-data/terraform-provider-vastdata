# Copyright (c) HashiCorp, Inc.

variable protection_policy_name {
    type = string
}

variable protection_policy_every {
    type = string
}

variable protection_policy_prefix {
    type = string
}

resource vastdata_protection_policy ppolicy1 {
        name = var.protection_policy_name
        indestructible = "false"
        prefix = var.protection_policy_prefix
        clone_type = "LOCAL"
        frames {
                every = var.protection_policy_every
                keep_local = "14D"
                start_at = "2023-06-04 09:00:00"
        }
        frames {
                every = var.protection_policy_every
                keep_local = "8D"
                start_at = "2023-07-04 09:00:00"
        }
}

output tf_protection_policy {
  value = vastdata_protection_policy.ppolicy1
}