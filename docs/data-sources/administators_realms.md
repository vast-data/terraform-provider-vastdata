---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_administators_realms Data Source - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_administators_realms (Data Source)



## Example Usage

```terraform
data "vastdata_administators_realms" "realm01" {
  name = "realm01"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) (Valid for versions: 5.2.0) The unique name of the realm.

### Read-Only

- `guid` (String) (Valid for versions: 5.2.0) The unique GUID of the realm.
- `id` (String) The ID of this resource.
- `object_types` (List of String) (Valid for versions: 5.2.0) A list of permissions granted. Allowed Values are [cnodegroup managedapplication managedapplicationset alarm event eventdefinition eventdefinitionconfig cbox cnode carrier cluster dbox dnode dtray ebox fan host nic nvram psu port rack ssd subnetmanager switch dns globalsnapstream kafkabroker nativereplicationremotetarget protectedpath protectionpolicy qospolicy quota quotaentityinfo replicationrestorepoint replicationstream replicationtarget s3lifecyclerule snapshot userquota vip vippool view viewpolicy monitor activedirectory encryptedpath encryptiongroup group indestructibility ldap manager nis permission realm role s3policy tenant user vms callhomeconfig challengetoken env license module supportbundle systemsettingsdiff]
