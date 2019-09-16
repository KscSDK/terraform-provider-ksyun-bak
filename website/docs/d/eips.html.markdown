---
layout: "ksyun"
page_title: "Ksyun: ksyun_eips"
sidebar_current: "docs-ksyun-datasource-eips"
description: |-
  Provides a list of EIP resources in the current region.
---

# ksyun_eips

This data source provides a list of EIP resources (Elastic IP address) according to their EIP ID.

## Example Usage

```hcl
# Get  eips
data "ksyun_eips" "default" {
  output_file="output_result"

  ids=[]
  project_id=["0"]
  instance_type=["Ipfwd"]
  network_interface_id=[]
  internet_gateway_id=[]
  band_width_share_id=[]
  line_id=[]
  public_ip=[]
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional)  A list of Elastic IP IDs, all the EIPs belong to this region will be retrieved if the ID is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `eips` - It is a nested type which documented below.
* `total_count` - Total number of Elastic IPs that satisfy the condition.

