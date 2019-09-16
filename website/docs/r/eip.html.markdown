---
layout: "ksyun"
page_title: "Ksyun: ksyun_eip"
sidebar_current: "docs-ksyun-resource-eip"
description: |-
  Provides an Elastic IP resource.
---

# ksyun_eip

Provides an Elastic IP resource.

## Example Usage

```hcl
data "ksyun_lines" "default" {
  output_file="output_result1"
  line_name="BGP"
}
resource "ksyun_eip" "default1" {
  line_id ="${data.ksyun_lines.default.lines.0.line_id}"
  band_width =1
  charge_type = "PostPaidByDay"
  purchase_time =1
  project_id=0
}
```

## Argument Reference

The following arguments are supported:

* `band_width` - (Optional) Maximum bandwidth to the elastic public network, measured in Mbps (Mega bit per second).
* `charge_type` - (Optional) Elastic IP charge type.

