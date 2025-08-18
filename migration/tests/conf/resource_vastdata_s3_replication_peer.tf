# Copyright (c) HashiCorp, Inc.

variable aws_access_key {
        type = string
}

variable aws_secret_key {
        type = string
}

variable aws_region {
        type = string
}

variable aws_bucket_name {
        type = string
}

variable s3_peer_name {
        type = string
}

resource vastdata_s3_replication_peers s3peer1 {
        name = var.s3_peer_name
        bucket_name = var.aws_bucket_name
        http_protocol = "https"
        type_ = "AWS_S3"
        aws_region = "eu-west-1"
        access_key = var.aws_access_key
        secret_key = var.aws_secret_key
}

output tf_s3_peer {
        sensitive = true
        value = vastdata_s3_replication_peers.s3peer1
}