---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_user_key Resource - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_user_key (Resource)



## Example Usage

```terraform
#Creating a user to create keys for.
#!!! it is important to note that each user can have up to 2 keys.
resource "vastdata_user" "example-user" {
  name = "example"
  uid  = 9000
}

#Create Key and provide pgp public key so that the secret will be encrypted using this public key
#The pgp public key should be provided at the ascii armor format, the encrypted secret_key retuend
#will be set to the encrypted_secret_key field
#This key will be created and set to be disabled.
resource "vastdata_user_key" "key1" {
  user_id        = vastdata_user.example-user.id
  enabled        = false
  pgp_public_key = <<-EOT
  -----BEGIN PGP PUBLIC KEY BLOCK-----
  .
  .  <public pgp key content>
  .
  -----END PGP PUBLIC KEY BLOCK-----
  EOT
}

#This key is provided without setting the pgp public key this means that after key creation
#The secret key returned will be stored set to the secret_key field, it is highly recomanded
#not to use this option and if so please make sure that your terraform backend is secured.
resource "vastdata_user_key" "key2" {
  user_id = vastdata_user.example-user.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `user_id` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0) The user id to create the Key for

### Optional

- `enabled` (Boolean) (Valid for versions: 5.0.0,5.1.0,5.2.0) Should the key be enabled or disabled
- `pgp_public_key` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The PGP public key at ascii armor format to encrypt the secret id returned from vast cluster, if this option is set than the encrypted_secret_key will be returned and secret_key will be empty, changing it after apply will have no affect

### Read-Only

- `access_key` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The access id of the user key
- `encrypted_secret_key` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The secret id returned from the vast cluster encrypted with the public key provided at pgp_public_key
- `id` (String) The ID of this resource.
- `secret_key` (String, Sensitive) (Valid for versions: 5.0.0,5.1.0,5.2.0) The secret id of the user key, please note that that the secret id is not encrypted and should be kept in an highly secure backend ,this field will only be returned if pgp_public_key is not provided
