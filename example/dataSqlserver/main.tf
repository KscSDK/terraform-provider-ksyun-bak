provider "ksyun"{
  region = "cn-shanghai-3"
  access_key = "aaaaaaaa"
  secret_key = "aaaaaaaa"
}



data "ksyun_sqlservers" "hou_desc" {
  output_file = "output_file"
  dbinstancestatus="active"
//  dbinstanceidentifier="0b0adac8-73c4-4d05-9b8b-982ca09dd313"
}