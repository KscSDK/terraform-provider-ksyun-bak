---
layout: "ksyun"
page_title: "Ksyun: ksyun_images"
sidebar_current: "docs-ksyun-datasource-images"
description: |-
  Provides a list of available image resources in the current region.
---

# ksyun_images

This data source providers a list of available image resources according to their availability zone, image ID and other fields.

## Example Usage

```hcl
# Get  ksyun_images
data "ksyun_images" "default" {
  output_file="output_result"
  ids=[]
  name_regex="centos-7.0-20180927115228"
  is_public=true
  image_source="system"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) The ID of image.
* `name_regex` - (Optional) A regex string to filter resulting images by name. (Such as: `^CentOS 7.[1-2] 64` means CentOS 7.1 of 64-bit operating system or CentOS 7.2 of 64-bit operating system, "^Ubuntu 16.04 64" means Ubuntu 16.04 of 64-bit operating system).
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `images` - It is a nested type which documented below.
* `total_count` - Total number of images that satisfy the condition.


