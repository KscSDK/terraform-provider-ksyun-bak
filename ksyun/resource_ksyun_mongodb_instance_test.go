package ksyun

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strings"
	"testing"
)

func TestAccKsyunMongodbInstance_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongodbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("ksyun_mongodb_instance.default"),
				),
			},
		},
	})
}

func testAccCheckMongodbInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find resource or data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("mongodb instance create failure")
		}
		return nil
	}
}

func testAccCheckMongodbInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*KsyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ksyun_mongodb_instance" {
			instanceCheck := make(map[string]interface{})
			instanceCheck["InstanceId"] = rs.Primary.ID
			resp, err := client.mongodbconn.DescribeMongoDBInstance(&instanceCheck)

			if err != nil {
				if strings.Contains(err.Error(), "InstanceNotFound") {
					return nil
				}
				return err
			}
			if resp != nil {
				if (*resp)["MongoDBInstancesResult"] != nil {
					return errors.New("delete instance failure")
				}
			}
		}
	}

	return nil
}

const testAccMongodbInstanceConfig = `
resource "ksyun_mongodb_instance" "default" {
  name = "InstanceName"
  instance_account = "root"
  instance_password = "admin"
  instance_class = "1C2G"
  storage = 5
  node_num = 3
  vpc_id = "VpcId"
  vnet_id = "VnetId"
  db_version = "3.6"
  pay_type = "byDay"
  iam_project_id = "0"
  availability_zone = "cn-shanghai-3b"
}
`
