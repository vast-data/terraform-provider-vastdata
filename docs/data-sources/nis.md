---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_nis Data Source - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_nis (Data Source)



## Example Usage

```terraform
data "vastdata_nis" "nis1" {
  domain_name = "my.nis.domain.example.com"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain_name` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The nis server domain name

### Read-Only

- `guid` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) A uniq guid given to the nis server configuration
- `hosts` (List of String) (Valid for versions: 5.0.0,5.1.0,5.2.0) List of ip addresses/hostnames of nis servers
- `id` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0) A uniq id given to the nis server configuration
