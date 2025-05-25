#Create a user for which keys will be created
#!!! Each user can have up to two keys. 
resource "vastdata_user" "example-user" {
  name = "example"
  uid  = 9000
}

#You can create the key and provide a PGP public key so that the secret will be encrypted using this public key.
#The PGP public key must be provided in the ASCII armor format. The encrypted secret key returned by the cluster will be set to the `encrypted_secret_key` field.
#This example creates a key and sets it to be disabled:
resource "vastdata_user_key" "key1" {
  user_id        = vastdata_user.example-user.id
  enabled        = false
  pgp_public_key = <<EOT
  -----BEGIN PGP PUBLIC KEY BLOCK-----
  .
  .  <public pgp key content>
  .
  -----END PGP PUBLIC KEY BLOCK-----
  EOT
}

#If you do not set the PGP public key during key creation, the returned secret key will be set to the `secret_key` field. It is highly recommended to avoid using this option. Otherwise, ensure that your Terraform backend is well secured.
resource "vastdata_user_key" "key2" {
  user_id = vastdata_user.example-user.id
}

