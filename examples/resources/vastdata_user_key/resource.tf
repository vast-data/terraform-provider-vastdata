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

