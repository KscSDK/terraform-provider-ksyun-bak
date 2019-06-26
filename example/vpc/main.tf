provider "ksyun" {
  region = "cn-beijing-6"
}

resource "ksyun_vpc" "test" {
  vpc_name   = "${var.vpc_name}"
  cidr_block = "${var.vpc_cidr}"
}
