---
layout: "ksyun"
page_title: "Ksyun: ksyun_lb"
sidebar_current: "docs-ksyun-resource-lb"
description: |-
  Provides a Load Balancer resource.
---

# ksyun_lb

Provides a Load Balancer resource.

## Example Usage

```hcl
resource "ksyun_lb" "default" {
  vpc_id = "74d0a45b-472d-49fc-84ad-221e21ee23aa"
  load_balancer_name = "tf-xun1"
  type = "public"
  subnet_id = "609d1736-d8d7-492d-abd3-1183bb60329e"
  load_balancer_state = "stop"
  private_ip_address = "10.0.77.11"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_name` - (Optional) The name of the load balancer. 
* `vpc_id` - (Optional) The ID of the VPC linked to the Load Balancers, This argumnet is not required if default VPC.
* `subnet_id` - (Optional) The ID of subnet that intrant load balancer belongs to. This argumnet is not required if default subnet.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The time of creation for load balancer, formatted in RFC3339 time string.
* `private_ip` - The IP address of intranet IP. It is `""` if `internal` is `false`.

## Import

LB can be imported using the `id`, e.g.

```
$ terraform import ksyun_lb.example ulb-abc123456
```