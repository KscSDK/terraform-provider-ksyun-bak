---
layout: "ksyun"
page_title: "Ksyun: ksyun_lb_acl_entry"
sidebar_current: "docs-ksyun-resource-lb-acl"
description: |-
  Provides a Load Balancer acl entry resource to add content forwarding policies for Load Balancer backend resource.
---

# ksyun_lb_acl_entry

Provides a Load Balancer acl entry resource to add content forwarding policies for Load Balancer backend resource.

## Example Usage

```hcl
resource "ksyun_lb_acl_entry" "default" {
  load_balancer_acl_id = "8e6d0871-da8a-481e-8bee-b3343e2a6166"
  cidr_block = "192.168.11.2/32"
  rule_number = 10
  rule_action = "allow"
  protocol = "ip"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_acl_id` - (Required) The ID of a load balancer acl.