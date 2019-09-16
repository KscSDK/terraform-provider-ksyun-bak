---
layout: "ksyun"
page_title: "Ksyun: ksyun_security_group_entry"
sidebar_current: "docs-ksyun-resource-security-group"
description: |-
  Provides a Security Group resource.
---

# ksyun_security_group_entry

Provides a Security Group resource.

## Example Usage

```hcl
resource "ksyun_security_group_entry" "default" {
  description = ""
  security_group_id="7385c8ea-79f7-4e9c-b99f-517fc3726256"
  cidr_block="10.0.0.1/32"
  direction="in"
  protocol="ip"
  icmp_type=0
  icmp_code=0
  port_range_from=0
  port_range_to=0
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the security group which contains 1-63 characters and only support Chinese, English, numbers, '-', '_' and '.'. If not specified, terraform will autogenerate a name beginning with `tf-security-group`.
* `port_range_from` - (Optional) The start of port numbers.
* `port_range_to` - (Optional) The end of port numbers.
* `cidr_block` - (Optional) The cidr block of source.
* `protocol` - (Optional) The protocol. Possible values are: `tcp`, `udp`, `icmp`, `ip`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The time of creation of security group, formatted in RFC3339 time string.

## Import

Security Group can be imported using the `id`, e.g.

```
$ terraform import ksyun_security_group_entry.example firewall-abc123456
```