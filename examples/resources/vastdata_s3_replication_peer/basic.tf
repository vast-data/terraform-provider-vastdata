#Create an S3 replication peer for an AWS bucket
resource "vastdata_s3_replication_peer" "vastdb_s3peer-aws" {
  name          = "vastdb_s3peer-aws"
  bucket_name   = "my-aws-s3-bucket"
  http_protocol = "https"
  type          = "AWS_S3"
  aws_region    = "eu-west-1"
  access_key    = "W21E6X5ZQEOODYB6J0UY"
  secret_key    = "fcESVNih9Ykb/bDSmKipQdinnHObrRyv9nre+nR1"
}