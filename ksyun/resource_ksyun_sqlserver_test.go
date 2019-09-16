package ksyun

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccKsyunSqlserver_basic(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_sqlserver.houbin-1",
		Providers:     testAccProviders,
		//CheckDestroy:  testAccCheckSqlserverDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfig,

				Check: resource.ComposeTestCheckFunc(
				//resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-1", "dbinstanceclass", "db.ram.2%7Cdb.disk.50"),
				//resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-1", "dbinstancename", "ksyun_sqlserver_1"),
				//resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-1", "dbinstancetype", "HRDS_SS"),
				//resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-1", "engine", "SQLServer"),
				//resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-1", "engineversion", "2008r2"),
				//resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-1", "masterusername", "admin"),
				//resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-1", "masteruserpassword", "123qweASD"),
				//resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-1", "version", "2008r2"),
				//resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-1", "vpcid", "3c12ccdf-9b8f-4d9b-8aa6-a523897e97a1"),
				//resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-1", "subnetid", "293c16a5-c757-405c-a693-3b2a3adead50"),
				//resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-1", "billtype", "DAY"),

				),
			},
		},
	})
}

//func TestAccKsyunSqlserver_update(t *testing.T) {
//	var val map[string]interface{}
//
//	resource.ParallelTest(t, resource.TestCase{
//		PreCheck: func() {
//			testAccPreCheck(t)
//		},
//
//		IDRefreshName: "ksyun_vpc.foo",
//		Providers:     testAccProviders,
//		CheckDestroy:  testAccCheckSqlserverDestroy,
//
//		Steps: []resource.TestStep{
//			{
//				Config: testAccVPCConfig,
//
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckSqlserverExists("ksyun_vpc.foo", &val),
//					testAccCheckSqlserverAttributes(&val),
//					//resource.TestCheckResourceAttr("ksyun_vpc.foo", "vpc_name", "tf-acc-vpc"),
//					//resource.TestCheckResourceAttr("ksyun_vpc.foo", "cidr_block", "192.168.0.0/16"),
//				),
//			},
//			{
//				Config: testAccVPCConfigUpdate,
//
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckSqlserverExists("ksyun_vpc.foo", &val),
//					testAccCheckSqlserverAttributes(&val),
//					//resource.TestCheckResourceAttr("ksyun_vpc.foo", "vpc_name", "tf-acc-vpc-1"),
//				),
//			},
//		},
//	})
//}
//
//func testAccCheckSqlserverExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		rs, ok := s.RootModule().Resources[n]
//
//		if !ok {
//			return fmt.Errorf("not found: %s", n)
//		}
//
//		if rs.Primary.ID == "" {
//			return fmt.Errorf("vpc id is empty")
//		}
//
//		client := testAccProvider.Meta().(*KsyunClient)
//		vpc := make(map[string]interface{})
//		vpc["VpcId.1"] = rs.Primary.ID
//		ptr, err := client.vpcconn.DescribeVpcs(&vpc)
//
//		if err != nil {
//			return err
//		}
//		if ptr != nil {
//			l := (*ptr)["VpcSet"].([]interface{})
//			if len(l) == 0 {
//				return err
//			}
//		}
//
//		*val = *ptr
//		return nil
//	}
//}
//
//func testAccCheckSqlserverAttributes(val *map[string]interface{}) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		if val != nil {
//			l := (*val)["VpcSet"].([]interface{})
//			if len(l) == 0 {
//				return fmt.Errorf("vpc id is empty")
//			}
//		}
//		return nil
//	}
//}
//
//func testAccCheckSqlserverDestroy(s *terraform.State) error {
//	for _, rs := range s.RootModule().Resources {
//		if rs.Type != "ksyun_vpc" {
//			continue
//		}
//
//		client := testAccProvider.Meta().(*KsyunClient)
//		vpc := make(map[string]interface{})
//		vpc["VpcId.1"] = rs.Primary.ID
//		ptr, err := client.vpcconn.DescribeVpcs(&vpc)
//
//		// Verify the error is what we want
//		if err != nil {
//			return err
//		}
//		if ptr != nil {
//			l := (*ptr)["VpcSet"].([]interface{})
//			if len(l) == 0 {
//				continue
//			} else {
//				return fmt.Errorf("VPC still exist")
//			}
//		}
//	}
//
//	return nil
//}

const testAccSqlserverConfig = `
resource "ksyun_sqlserver" "houbin-1"{

  dbinstanceclass= "db.ram.2%7Cdb.disk.50"
  dbinstancename = "ksyun_sqlserver_1"
  dbinstancetype = "HRDS_SS"
  engine = "SQLServer"
  engineversion = "2008r2"
  masterusername = "admin"
  masteruserpassword = "123qweASD"
  version = "2008r2"
  vpcid =	"3c12ccdf-9b8f-4d9b-8aa6-a523897e97a1"
  subnetid = "293c16a5-c757-405c-a693-3b2a3adead50"
  billtype = "DAY"
}
`

//const testAccSqlserverConfigUpdate = `
//resource "ksyun_vpc" "foo" {
//	vpc_name        = "tf-acc-vpc-1"
//    cidr_block      = "192.168.0.0/16"
//}
//`
