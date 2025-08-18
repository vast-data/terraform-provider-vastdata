# Copyright (c) HashiCorp, Inc.

variable broker_name {
    type = string
}

variable broker_host {
    type = string
}

variable broker_port {
    type = number
}

variable tenant_name {
    type = string
}

variable tenant_client_start_ip {
    type = string
}

variable tenant_client_end_ip {
    type = string
}

resource vastdata_tenant broker_tenant1 {
    name = var.tenant_name
    client_ip_ranges {
        start_ip = var.tenant_client_start_ip
        end_ip = var.tenant_client_end_ip
    }
}

resource vastdata_kafka_brokers broker1 {
    name = var.broker_name
    addresses {
        host = var.broker_host
        port = var.broker_port
    }
    tenant_id = vastdata_tenant.broker_tenant1.id
}

output tf_broker {
  value = vastdata_kafka_brokers.broker1
}

output tf_broker_tenant {
  value = vastdata_tenant.broker_tenant1
}