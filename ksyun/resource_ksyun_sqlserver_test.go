package ksyun

import (
	"github.com/hashicorp/terraform/helper/resource"
	"os"
	"testing"
)

func TestAccKsyunSqlserver_basic(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	os.Setenv("KSYUN_ACCESS_KEY", "")
	os.Setenv("KSYUN_SECRET_KEY", "")
	os.Setenv("KSYUN_REGION", "cn-beijing-6")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_sqlserver.houbin-2",
		Providers:     testAccProviders,
		//CheckDestroy:  testAccCheckSqlserverDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfig,

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-2", "dbinstanceclass", "db.ram.2|db.disk.100"),
					resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-2", "dbinstancename", "ksyun_sqlserver_2"),
					resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-2", "dbinstancetype", "HRDS_SS"),
					resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-2", "engine", "SQLServer"),
					resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-2", "engineversion", "2008r2"),
					resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-2", "masterusername", "admin"),
					resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-2", "masteruserpassword", "123qweASD"),
					//resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-2", "version", "2019-04-25"),
					resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-2", "vpcid", "40e0c2e0-3607-4f17-abb5-1a6efe3951c8"),
					resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-2", "subnetid", "bc159134-4c94-4a6b-bec0-d97c75d83774"),
					resource.TestCheckResourceAttr("ksyun_sqlserver.houbin-2", "billtype", "DAY"),

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
resource "ksyun_sqlserver" "houbin-2"{

  dbinstanceclass= "db.ram.2|db.disk.100"
  dbinstancename = "ksyun_sqlserver_2"
  dbinstancetype = "HRDS_SS"
  engine = "SQLServer"
  engineversion = "2008r2"
  masterusername = "admin"
  masteruserpassword = "123qweASD"
  
  vpcid =	"40e0c2e0-3607-4f17-abb5-1a6efe3951c8"
  subnetid = "bc159134-4c94-4a6b-bec0-d97c75d83774"
  billtype = "DAY"
}
`

//version = "2008r2"

//const testAccSqlserverConfigUpdate = `
//resource "ksyun_vpc" "foo" {
//	vpc_name        = "tf-acc-vpc-1"
//    cidr_block      = "192.168.0.0/16"
//}
//`
