---
layout: "ksyun"
page_title: "Ksyun: ksyun_instances"
sidebar_current: "docs-ksyun-datasource-instances"
description: |-
  Provides a list of UHost instance resources in the current region.
---

# ksyun_instances

This data source providers a list of UHost instance resources according to their availability zone, instance ID and tag.

## Example Usage

```hcl
# Get  instances
data "ksyun_instances" "default" {
  output_file = "output_result"
  ids = []
  project_id = []
  network_interface {
    network_interface_id = []
    subnet_id = []
    group_id = []
  }
  instance_state {
    name =  []
  }
  availability_zone {
    name =  []
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) Availability zone where instances are located. Such as: "cn-bj2-02". You may refer to [list of availability zone](https://docs.ksyun.cn/api/summary/regionlist)
* `ids` - (Optional) A list of instance IDs, all the instances belongs to the defined region will be retrieved if this argument is "".
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - It is a nested type. instances documented below.
* `total_count` - Total number of instances that satisfy the condition.

The attribute (`instances`) support the following:

* `availability_zone` - Availability zone where instances are located.
* `id` - The ID of instance.
* `instance_state` - The state of instance.
* `vpc_id` - The ID of VPC linked to the instance.
* `subnet_id` - The ID of subnet linked to the instance.
