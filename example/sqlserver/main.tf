provider "ksyun"{
  region = "cn-shanghai-3"
  access_key = ""
  secret_key = ""
}

//resource "ksyun_sqlserver" "houbin-10"{
//  output_file = "output_file"
//  dbinstanceclass= "db.ram.2|db.disk.20"
//  dbinstancename = "ksyun_sqlserver_9"
//  dbinstancetype = "HRDS_SS"
//  engine = "SQLServer"
//  engineversion = "2008r2"
//  masterusername = "admin"
//  masteruserpassword = "123qweASD"
//  vpcid = "cbfdbc08-912a-4ec1-972d-e80bc6fe8aae"
//  subnetid = "87df0198-71e5-4c65-8b7f-a860fbdbeb47"
//  billtype = "DAY"
//}

resource "ksyun_sqlserver" "houbin-1"{
  output_file = "output_file"
  dbinstanceclass= "db.ram.2|db.disk.22"
  dbinstancename = "ksyun_sqlserver_1"
  dbinstancetype = "HRDS_SS"
  engine = "SQLServer"
  engineversion = "2008r2"
  masterusername = "admin"
  masteruserpassword = "123qweASD"
  vpcid = "cbfdbc08-912a-4ec1-972d-e80bc6fe8aae"
  subnetid = "87df0198-71e5-4c65-8b7f-a860fbdbeb47"
  billtype = "DAY"
}

//resource "ksyun_sqlserver" "houbin-11"{
//  output_file = "output_file"
//  dbinstanceclass= "db.ram.2|db.disk.20"
//  dbinstancename = "ksyun_sqlserver_11"
//  dbinstancetype = "HRDS_SS"
//  engine = "SQLServer"
//  engineversion = "2008r2"
//  masterusername = "admin"
//  masteruserpassword = "123qweASD"
//  vpcid = "cbfdbc08-912a-4ec1-972d-e80bc6fe8aae"
//  subnetid = "87df0198-71e5-4c65-8b7f-a860fbdbeb47"
//  billtype = "DAY"
//}

//data "ksyun_sqlservers" "hou_desc" {
//  dbinstancestatus="active"
//  dbinstanceidentifier="0b0adac8-73c4-4d05-9b8b-982ca09dd313"
//}