# Creating a non-local user key without encryption.
# This key will be created and set to be enabled.
# Each non-local user (identified by UID) can have up to 2 keys.
resource "vastdata_non_local_user_key" "ExternalUserKey" {
    uid                 = 1097416930
    tenant_id           = 1
    enabled             = true
}

# Creating a non-local user key with encryption.
# A PGP public key is provided in ASCII-armored format.
# The corresponding secret_key will be encrypted with this key,
# and the result will be available in the encrypted_secret_key field.
# This key will be created and set to be enabled.
resource "vastdata_non_local_user_key" "EncryptedExternalUserKey" {
    uid                 = 1097416930
    tenant_id           = 1
    enabled             = true
    pgp_public_key = <<EOT
    -----BEGIN PGP PUBLIC KEY BLOCK-----
    .
    .  <public pgp key content>
    .
    -----END PGP PUBLIC KEY BLOCK-----
    EOT
}
