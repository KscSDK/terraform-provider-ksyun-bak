package ksyun

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccKsyunSqlServer_basic(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_sqlserver.ks-ss-1",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSqlServerDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccSqlServerConfig,

				Check: resource.ComposeTestCheckFunc(
					testCheckSqlServerExists("ksyun_sqlserver.ks-ss-1", &val),
				),
			},
		},
	})
}

func testCheckSqlServerExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found : %s", n)
		}
		if res.Primary.ID == "" {
			return fmt.Errorf("instance is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		req := map[string]interface{}{
			"DBInstanceIdentifier": res.Primary.ID,
		}
		resp, err := client.krdsconn.DescribeDBInstances(&req)
		if err != nil {
			return err
		}
		if resp != nil {
			bodyData, dataOk := (*resp)["Data"].(map[string]interface{})
			if !dataOk {
				return fmt.Errorf("error on reading Instance(krds)  %+v", (*resp)["Error"])
			}
			instances := bodyData["Instances"].([]interface{})
			if len(instances) == 0 {
				return fmt.Errorf("no instance find, instance number is 0")
			}
		}
		*val = *resp
		return nil
	}
}

func testAccCheckSqlServerDestroy(s *terraform.State) error {
	for _, res := range s.RootModule().Resources {
		if res.Type != "ksyun_sqlserver" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		req := map[string]interface{}{
			"DBInstanceIdentifier": res.Primary.ID,
		}
		resp, err := client.krdsconn.DescribeDBInstances(&req)
		if err != nil {
			if err.(awserr.Error).Code() == "NOT_FOUND" {
				return nil
			}
			return err
		}
		if resp != nil {
			bodyData, dataOk := (*resp)["Data"].(map[string]interface{})
			if !dataOk {
				return fmt.Errorf("error on reading Instance(krds)  %+v", (*resp)["Error"])
			}
			instances := bodyData["Instances"].([]interface{})
			if len(instances) != 0 {
				return fmt.Errorf("no instance find, instance number is 0")
			}
		}
	}

	return nil
}

const testAccSqlServerConfig = `
resource "ksyun_sqlserver" "ks-ss-1"{
  output_file = "output_file"
  db_instance_class= "db.ram.1|db.disk.15"
  db_instance_name = "ksyun_sqlserver_1"
  db_instance_type = "HRDS_SS"
  engine = "SQLServer"
  engine_version = "2008r2"
  master_user_name = "admin"
  master_user_password = "123qweASD"
  vpc_id = "9d4d768e-36ff-49b0-b72d-dc8b760c51cb"
  subnet_id = "d3f42a3b-5b6b-4a24-9e6c-785941540d23"
  bill_type = "DAY"
  security_group_id = "27813"
  port = "54321"
}
`
