provider "ksyun"{
  region = "cn-shanghai-3"
  access_key = ""
  secret_key = ""
}


resource "ksyun_krds" "houbin_terraform_4"{
  output_file = "output_file"
  db_instance_class= "db.ram.2|db.disk.21"
  db_instance_name = "houbin_terraform_1-n"
  db_instance_type = "HRDS"
  engine = "mysql"
  engine_version = "5.5"
  master_user_name = "admin"
  master_user_password = "123qweASD123"
  vpc_id = "19b422fa-74b2-45ac-8b03-fe4d955f27cc"
  subnet_id = "02744161-e5c0-4a70-9bf4-975a3f8ae0be"
  bill_type = "DAY"
  security_group_id = "27936" //"27936"
  preferred_backup_time = "01:00-02:00"
  db_instance_identifier= "666666"
  parameters {
    name = "auto_increment_increment"
    value = "8"
  }

  parameters {
    name = "binlog_format"
    value = "ROW"
  }

  parameters {
    name = "delayed_insert_limit"
    value = "108"
  }
  parameters {
    name = "auto_increment_offset"
    value= "2"
  }
}




