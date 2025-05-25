#Create an S3 replication peer for a custom bucket
resource "vastdata_s3_replication_peers" "s3peers" {
  name              = "s3peer"
  bucket_name       = "s3bucket"
  http_protocol     = "https"
  type_             = "CUSTOM_S3"
  custom_bucket_url = "customs3.bucket.com"
  access_key        = "W21E6X5ZQEOODYB6J0UY"
  secret_key        = "fcESVNih9Ykb/bDSmKipQdinnHObrRyv9nre+nR1"

}

#Create an S3 replication peer for an AWS bucket
resource "vastdata_s3_replication_peers" "s3peer-aws" {
  name          = "s3peer-aws"
  bucket_name   = "my-aws-s3-bucket"
  http_protocol = "https"
  type_         = "AWS_S3"
  aws_region    = "eu-west-1"
  access_key    = "W21E6X5ZQEOODYB6J0UY"
  secret_key    = "fcESVNih9Ykb/bDSmKipQdinnHObrRyv9nre+nR1"

}
