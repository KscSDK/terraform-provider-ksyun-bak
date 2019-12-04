package ksyun

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

func TestAccKsyunKrds_basic(t *testing.T) {
	var val map[string]interface{}

	os.Setenv("KSYUN_ACCESS_KEY", "AKLTZ2u0gyI-TNKNpKe7668O8g")
	os.Setenv("KSYUN_SECRET_KEY", "OHZHhGDHXFgnFSh8ImHc0OwyxzRFsnqhKlqqWEr5X2DA0fU0PjTHNPQNsKOH/PD37w==")
	os.Setenv("KSYUN_REGION", "cn-shanghai-3")
	os.Setenv("TF_ACC", "1")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_krds.rds_terraform_3",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKrdsDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccKrdsConfig,

				Check: resource.ComposeTestCheckFunc(
					testCheckKrdsExists("ksyun_krds.rds_terraform_3", &val),
				),
			},
		},
	})
}

func testCheckKrdsExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
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

func testAccCheckKrdsDestroy(s *terraform.State) error {
	for _, res := range s.RootModule().Resources {
		if res.Type != "ksyun_krds" {
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

const testAccKrdsConfig = `
resource "ksyun_krds" "rds_terraform_4"{
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

`
