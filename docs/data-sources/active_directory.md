---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_active_directory Data Source - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_active_directory (Data Source)



## Example Usage

```terraform
data "vastdata_active_directory" "ad1" {
  machine_account_name = "machine1"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `machine_account_name` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0)

### Read-Only

- `guid` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `id` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `ldap_id` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `organizational_unit` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0)
