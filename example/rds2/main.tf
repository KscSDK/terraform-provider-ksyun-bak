provider "ksyun"{
  region = "cn-beijing-6"
  access_key = ""
  secret_key = ""
}


//resource "ksyun_krds" "houbin_terraform_2"{
//  output_file = "output_file"
//  db_instance_class= "db.ram.4|db.disk.101"
//  db_instance_name = "houbin_terraform_666"
//  db_instance_type = "HRDS"
//  engine = "mysql"
//  engine_version = "5.6"
//  master_user_name = "admin"
//  master_user_password = "123qweASD123"
//  vpc_id = "7f8ea0b7-a624-4279-9ca7-2284a1380878"
//  subnet_id = "c0a3de22-5f5b-47eb-bdeb-89b7414bd5f2"
//  bill_type = "DAY"
//  //db_parameter_group_id = "0a0708ac-f7da-4b2f-ba2b-599bc9f77759"
//  security_group_id = "62185"
//
//  parameters {
//    name = "auto_increment_increment"
//    value = "3"
//  }
//
//  parameters {
//    name = "binlog_format"
//    value = "ROW"
//  }
//
//  parameters {
//    name = "delayed_insert_limit"
//    value = "101"
//  }
//}
