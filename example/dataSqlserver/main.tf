provider "ksyun"{
  region = "cn-shanghai-3"
  access_key = "AKLTpmhJ3QlBQEmB401iSDl0dA"
  secret_key = "OAeYHxiil7rl7nbVcpsUPnbvzJEkY6zQM4ExOR4aOYUx4SZhwLqrKnlaCETVyVv7gw=="
}



data "ksyun_sqlservers" "hou_desc" {
  output_file = "output_file"
  dbinstancestatus="active"
//  dbinstanceidentifier="0b0adac8-73c4-4d05-9b8b-982ca09dd313"
}