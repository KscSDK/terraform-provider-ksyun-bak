package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strings"
	"testing"
)

func TestAccKsyunMongodbShardInstance_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongodbShardInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbShardInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbShardInstanceExists("ksyun_mongodb_shard_instance.default"),
				),
			},
		},
	})
}

func testAccCheckMongodbShardInstanceExists(n string) resource.TestCheckFunc {
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

func testAccCheckMongodbShardInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*KsyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ksyun_mongodb_shard_instance" {
			instanceCheck := make(map[string]interface{})
			instanceCheck["InstanceId"] = rs.Primary.ID
			_, err := client.mongodbconn.DescribeMongoDBInstance(&instanceCheck)

			if err != nil {
				if strings.Contains(err.Error(), "InstanceNotFound") {
					return nil
				} else {
					return fmt.Errorf("mongodb instance delete failure")
				}
			}
		}
	}

	return nil
}

const testAccMongodbShardInstanceConfig = `
resource "ksyun_mongodb_shard_instance" "default" {
  name = "InstanceName"
  instance_account = "root"
  instance_password = "admin"
  mongos_class = "1C2G"
  mongos_num = 2
  shard_class = "1C2G"
  shard_num = 2
  storage = 5
  vpc_id = "VpcId"
  vnet_id = "VnetId"
  db_version = "3.6"
  pay_type = "hourlyInstantSettlement"
  iam_project_id = "0"
  availability_zone = "cn-shanghai-3b"
}
`
