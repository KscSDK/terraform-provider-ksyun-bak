# Specify the provider and access details
provider "ksyun" {
  access_key = "ak"
  secret_key = "sk"
  region = "cn-shanghai-3"
}

resource "ksyun_mongodb_instance" "default" {
  name = "InstanceName"
  instance_account = "root"
  instance_password = "admin"
  instance_class = "1C2G"
  storage = 5
  node_num = 3
  vpc_id = "VpcId"
  vnet_id = "VnetId"
  db_version = "3.6"
  pay_type = "byDay"
  iam_project_id = "0"
  availability_zone = "cn-shanghai-3b"
}

resource "ksyun_mongodb_shard_instance" "default" {
  name = "InstanceName"
  instance_account = "root"
  instance_password = "admin"
  mongos_class = "1C2G"
  mongos_num = 2
  shard_class = "1C2G"
  shard_num = 2
  storage = 5
  vpc_id = "VpcId"
  vnet_id = "VnetId"
  db_version = "3.6"
  pay_type = "hourlyInstantSettlement"
  iam_project_id = "0"
  availability_zone = "cn-shanghai-3b"
}

resource "ksyun_mongodb_security_rule" "repset" {
  instance_id = "${ksyun_mongodb_instance.default.id}"
  cidrs = "192.168.10.1/32"
}

resource "ksyun_mongodb_security_rule" "cluster" {
  instance_id = "${ksyun_mongodb_shard_instance.default.id}"
  cidrs = "192.168.10.1/32,192.168.20.1/32"
}

