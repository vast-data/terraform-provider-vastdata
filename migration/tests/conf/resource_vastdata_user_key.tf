# Copyright (c) HashiCorp, Inc.

# To refer to a specific system, need to add `provider = vastdata.system<IDX>` where
# IDX is the system's index (in the order of appearance in Comet's commandline)
variable user_name {
    type = string
}

variable user_uid {
    type = number
}

resource vastdata_user user2 {
    name = var.user_name
    uid = var.user_uid
}

#Create Key and provide pgp public key so that the secret will be encrypted using this public key
#The pgp public key should be provided at the ascii armor format, the encrypted secret_key retuend
#will be set to the encrypted_secret_key field
#This key will be created and set to be disabled.
resource vastdata_user_key user_key1 {
    user_id        = vastdata_user.user2.id
    enabled        = false
    pgp_public_key = <<EOT
    -----BEGIN PGP PUBLIC KEY BLOCK-----
Version: Keybase OpenPGP v1.0.0
Comment: https://keybase.io/crypto

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
qyjqu/s3g8Fn6+sNO0X5aRM=
=Wtzq
-----END PGP PUBLIC KEY BLOCK-----
    EOT
}

#This key is provided without setting the pgp public key this means that after key creation
#The secret key returned will be stored set to the secret_key field, it is highly recomanded
#not to use this option and if so please make sure that your terraform backend is secured.
resource vastdata_user_key user_key2 {
    user_id = vastdata_user.user2.id
}

output tf_user {
  value = vastdata_user.user2
}

output tf_key1 {
  value = vastdata_user_key.user_key1
  sensitive = true
}

output tf_key2 {
  value = vastdata_user_key.user_key2
  sensitive = true
}