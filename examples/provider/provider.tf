#A Simple provider with the version validation level of `warn`
provider vastdata {
  username = "<username>"
  port = 443
  password = "<password>"
  host = "<address>"
  skip_ssl_verify = true
  version_validation_mode = "warn"
}

#Define two providers for two different clusters named clusterA and clusterB, one with port 443 and the other with port 9443

provider vastdata {
  api_token = "<api_token>"
  port = 443
  host = "<address>"
  skip_ssl_verify = true
  alias = clusterA
}

provider vastdata {
  username = "<username>"
  port = 9443
  password = "<password>"
  host = "<address>"
  skip_ssl_verify = true
  alias = clusterB
}


