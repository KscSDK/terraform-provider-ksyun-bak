package ksyun

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccKsyunKrdsRR_basic(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_krds_rr.rds-rr-1",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKrdsDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccKrdsRRConfig,

				Check: resource.ComposeTestCheckFunc(
					testCheckKrdsRRExists("ksyun_krds_rr.rds-rr-1", &val),
				),
			},
		},
	})
}

func testCheckKrdsRRExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
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

func testAccCheckKrdsRRDestroy(s *terraform.State) error {
	for _, res := range s.RootModule().Resources {
		if res.Type != "ksyun_krds_rr" {
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

const testAccKrdsRRConfig = `
resource "ksyun_krds_rr" "rds-rr-1"{
  output_file = "output_file"
  db_instance_identifier= "ca79f6e5-cf73-42d8-9fd2-8c40fd5a22a3"
  db_instance_class= "db.ram.2|db.disk.50"
  db_instance_name = "houbin_terraform_888_rr_1"
  bill_type = "DAY"
  security_group_id = "62185"

  parameters {
    name = "auto_increment_increment"
    value = "7"
  }

  parameters {
    name = "binlog_format"
    value = "ROW"
  }

  parameters {
    name = "delayed_insert_limit"
    value = "107"
  }
}

`
