---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "upcloud_loadbalancer_frontend_tls_config Resource - terraform-provider-upcloud"
subcategory: ""
description: |-
  This resource represents frontend TLS config
---

# upcloud_loadbalancer_frontend_tls_config (Resource)

This resource represents frontend TLS config

## Example Usage

```terraform
resource "upcloud_loadbalancer_frontend_tls_config" "lb_fe_1_tls1" {
  frontend           = resource.upcloud_loadbalancer_frontend.lb_fe_1.id
  name               = "lb-fe-1-tls1-test"
  certificate_bundle = resource.upcloud_loadbalancer_manual_certificate_bundle.lb-cb-m1.id
}

variable "lb_zone" {
  type    = string
  default = "fi-hel2"
}

resource "upcloud_network" "lb_network" {
  name = "lb-test-net"
  zone = var.lb_zone
  ip_network {
    address = "10.0.0.0/24"
    dhcp    = true
    family  = "IPv4"
  }
}

resource "upcloud_loadbalancer_manual_certificate_bundle" "lb-cb-m1" {
  name        = "lb-cb-m1-test"
  certificate = "LS0tLS1CRUdJTiBDRVJ..."
  private_key = "LS0tLS1CRUdJTiBQUkl..."
}

resource "upcloud_loadbalancer_frontend" "lb_fe_1" {
  loadbalancer         = resource.upcloud_loadbalancer.lb.id
  name                 = "lb-fe-1-test"
  mode                 = "http"
  port                 = 8080
  default_backend_name = resource.upcloud_loadbalancer_backend.lb_be_1.name
}

resource "upcloud_loadbalancer" "lb" {
  configured_status = "started"
  name              = "lb-test"
  plan              = "development"
  zone              = var.lb_zone
  network           = resource.upcloud_network.lb_network.id
}

resource "upcloud_loadbalancer_backend" "lb_be_1" {
  loadbalancer = resource.upcloud_loadbalancer.lb.id
  name         = "lb-be-1-test"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **certificate_bundle** (String) Reference to certificate bundle ID.
- **frontend** (String) ID of the load balancer frontend to which the TLS config is connected.
- **name** (String) The name of the TLS config must be unique within service frontend.

### Optional

- **id** (String) The ID of this resource.

