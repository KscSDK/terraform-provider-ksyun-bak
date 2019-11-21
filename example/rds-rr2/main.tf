provider "ksyun"{
  region = "cn-beijing-6"
  access_key = "-"
  secret_key = ""
}


resource "ksyun_krds_rr" "houbin_terraform_3"{
  output_file = "output_file"
  db_instance_identifier= "ca79f6e5-cf73-42d8-9fd2-8c40fd5a22a3"
  db_instance_class= "db.ram.2|db.disk.66"
  db_instance_name = "houbin_terraform_888_rr_2"
  bill_type = "DAY"
  security_group_id = "62185"

  parameters {
    name = "auto_increment_increment"
    value = "3"
  }

  parameters {
    name = "binlog_format"
    value = "ROW"
  }

  parameters {
    name = "delayed_insert_limit"
    value = "103"
  }
}




