provider "ksyun"{
  region = "cn-shanghai-3"
  access_key = ""
  secret_key = ""
}

resource "ksyun_krds_security_group" "krds_sec_group_12" {
  output_file = "output_file"
  security_group_name = "terraform_security_group_12"
  security_group_description = "terraform-security-group-12"
  security_group_rule{
    security_group_rule_protocol = "0.0.0.0/0"
    security_group_rule_name = "all-shit"
  }
  security_group_rule{
    security_group_rule_protocol = "182.56.0.0/16"
    security_group_rule_name = "all-shit"
  }
  security_group_rule{
    security_group_rule_protocol = "182.133.0.0/16"
    security_group_rule_name = "wtf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.134.0.0/16"
    security_group_rule_name = "wtf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.135.0.0/16"
    security_group_rule_name = "wtf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.136.0.0/16"
    security_group_rule_name = "wtf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.137.0.0/16"
    security_group_rule_name = "wtf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.138.0.0/16"
    security_group_rule_name = "wtf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.139.0.0/16"
    security_group_rule_name = "wtf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.140.0.0/16"
    security_group_rule_name = "wtf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.141.0.0/16"
    security_group_rule_name = "wtf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.142.0.0/16"
    security_group_rule_name = "wtf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.143.0.0/16"
    security_group_rule_name = "wtf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.144.0.0/16"
    security_group_rule_name = "wtf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.145.0.0/16"
    security_group_rule_name = "wtf"
  }
}