---
layout: "ksyun"
page_title: "Ksyun: ksyun_instance"
sidebar_current: "docs-ksyun-resource-instance"
description: |-
  Provides an Host Instance resource.
---

# ksyun_instance

Provides an Host Instance resource.

~> **Note** The instance will reboot automatically to make the change take effect when update `instance_type`, `root_password`, `boot_disk_size`, `data_disk_size`.

## Example Usage

```hcl
data "ksyun_availability_zones" "default" {
  output_file=""
  ids=[]
}
data "ksyun_lines" "default" {
  output_file=""
  line_name="BGP"
}

resource "ksyun_vpc" "default" {
  vpc_name   = "${var.vpc_name}"
  cidr_block = "${var.vpc_cidr}"
}
resource "ksyun_subnet" "default" {
  subnet_name      = "${var.subnet_name}"
  cidr_block = "10.1.0.0/21"
  subnet_type = "Normal"
  dhcp_ip_from = "10.1.0.2"
  dhcp_ip_to = "10.1.0.253"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.1.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}

resource "ksyun_security_group" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  security_group_name="${var.security_group_name}"
}
resource "ksyun_security_group_entry" "default" {
  description = "test1"
  security_group_id="${ksyun_security_group.default.id}"
  cidr_block="10.0.1.1/32"
  direction="in"
  protocol="ip"
  icmp_type=0
  icmp_code=0
  port_range_from=0
  port_range_to=0
}

resource "ksyun_ssh_key" "default" {
  key_name="ssh_key_tf"
  public_key=""
}
resource "ksyun_instance" "default" {
  image_id="${data.ksyun_images.centos-7_5.images.0.image_id}"
  instance_type="N3.2B"
  system_disk{
    disk_type="SSD3.0"
    disk_size=30
  }
  data_disk_gb=0
  data_disk =[
    {
      type="SSD3.0"
      size=20
      delete_with_instance=true
    }
  ]
  subnet_id="${ksyun_subnet.default.id}"
  instance_password="Xuan663222"
  keep_image_login=false
  charge_type="Daily"
  purchase_time=1
  security_group_id=["${ksyun_security_group.default.id}"]
  private_ip_address=""
  instance_name="xuan-tf-combine"
  instance_name_suffix=""
  sriov_net_support=false
  project_id=0
  data_guard_id=""
  key_id=["${ksyun_ssh_key.default.id}"]
  force_delete=true
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required) Availability zone where instance is located.
* `image_id` - (Required) The ID for the image to use for the instance.
* `data_disk_type` - (Optional) The type of local data disk. Possible values are: `local_normal` and `local_ssd` for local data disk. (Default: `local_normal`). The `local_ssd` is not supported in all regions as data disk type, please proceed to Ksyun console for more details.
* `data_disk_size` - (Optional) The size of local data disk, measured in GB (GigaByte), range: 0-8000 (Default: `20`), 0-8000 for cloud disk, 0-2000 for local sata disk and 100-1000 for local ssd disk (all the GPU type instances are included). The volume adjustment must be a multiple of 10 GB. In addition, any reduction of data disk size is not supported. 
* `name` - (Optional) The name of instance, which contains 1-63 characters and only support Chinese, English, numbers, '-', '_', '.'. If not specified, terraform will autogenerate a name beginning with `tf-instance`.
* `security_group` - (Optional) The ID of the associated security group.
* `subnet_id` - (Optional) The ID of subnet. If defined `vpc_id`, the `subnet_id` is Required. If not defined `vpc_id` and `subnet_id`, the instance will use the default subnet in the current region.
* `vpc_id` - (Optional) The ID of VPC linked to the instance. If not defined `vpc_id`, the instance will use the default VPC in the current region.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The time of creation for instance, formatted in RFC3339 time string.
* `status` - Instance current status. Possible values are `Initializing`, `Starting`, `Running`, `Stopping`, `Stopped`, `Install Fail`, `ResizeFail` and `Rebooting`.
* `private_ip` - The private IP address assigned to the instance.
* `disk_set` - It is a nested type which documented below.

The attribute (`disk_set`) supports the following:

* `id` - The ID of disk.
* `size` - The size of disk, measured in GB (Gigabyte).
* `type` - The type of disk.


## Import

Instance can be imported using the `id`, e.g.

```
$ terraform import ksyun_instance.example uhost-abcdefg
```