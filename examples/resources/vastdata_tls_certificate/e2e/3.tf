# e2e: TLS certificate — different CA, KAFKA protocol, no CRL
resource "vastdata_tls_certificate" "vastdb_tls_cert" {
  ca_certificate_name = "vastdb-tls-cert-3"
  protocols           = ["KAFKA"]
  tenant_id           = 1

  ca_certificate_file = <<-EOT
    -----BEGIN CERTIFICATE-----
    MIIDRzCCAi+gAwIBAgIUes3Uoln7gjggDaAlY82R16n0I9IwDQYJKoZIhvcNAQEL
    BQAwMzEQMA4GA1UEAwwHVGVzdENBMjESMBAGA1UECgwJVmFzdFRlc3QyMQswCQYD
    VQQGEwJVUzAeFw0yNjA2MTUwMDMxMzBaFw0zNjA2MTIwMDMxMzBaMDMxEDAOBgNV
    BAMMB1Rlc3RDQTIxEjAQBgNVBAoMCVZhc3RUZXN0MjELMAkGA1UEBhMCVVMwggEi
    MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDCLHuKqa/XBIql19bfx6QaYP0j
    j553O8Asp9TVq62xu5Ky1Ddk11yS+Bzy/BwMK1msvlfm9NRXgkj7cn0uWjCmjQJO
    0iKRP5fE7enVEbl8zqTT9FuMdOPSot5OV8MqJcZIrOc7AYuZPTUFJm85j11FdLkE
    7r7wLPuMqIFhFqy2CA5ym7HAZPn+UORNwD2A2A5nEomHBaOhQFNXFnssZnFSR8S9
    2KrbgF8aUhtKjTa0+C67O6JJh71kEdSrtWjKpZHXjz7XOzmRtksJLnG7uLiLq0A1
    8NP042AhCukSBQKkpM9hREJtbI3t+57B7AJgj1ezmbx5M/o+d2DUp9RePsHRAgMB
    AAGjUzBRMB0GA1UdDgQWBBSrO19E4Y+iHou0sy7eDV8hVSoZEjAfBgNVHSMEGDAW
    gBSrO19E4Y+iHou0sy7eDV8hVSoZEjAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3
    DQEBCwUAA4IBAQAl6V+ZziAG6Vo2/7pKf32tWzDRAMx1lVsGdIQSa9UGihgV05Hg
    afiJapMffvYNDGCl/WbXfc3hfDy+BcHgWcrHBiNp/bSG51pWYofCXlKbtsBuy9yR
    oJdqFflOY55bgrofaT6LTalZHe7r9M3raXiQAbYDJSETQoUQFdAzf4hhWIFx29tc
    IQ4nMAyM5Sq/HjRnUDqXnqK2Fqnpa2bCCYBzqSoWo4fht9UOsPeYXb5iM8QXjKf3
    KoXUXCgok6Uee0yJbEuDspflfNioSgsiqRVSWBggvjO8SmU40xIFpTUNvstE3rIE
    0JluiqFiNQZyrDnDmsNqz9rk6mMrvxwmKEed
    -----END CERTIFICATE-----
  EOT
}
