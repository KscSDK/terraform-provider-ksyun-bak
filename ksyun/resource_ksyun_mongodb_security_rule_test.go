package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccKsyunMongodbSecurityRule_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongodbSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbSecurityRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbSecurityRuleExists("ksyun_mongodb_security_rule.default"),
				),
			},
		},
	})
}

func testAccCheckMongodbSecurityRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find resource or data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("mongodb instance is not exist")
		}

		client := testAccProvider.Meta().(*KsyunClient)
		securityRuleCheck := make(map[string]interface{})
		securityRuleCheck["InstanceId"] = rs.Primary.ID
		resp, err := client.mongodbconn.ListSecurityGroupRules(&securityRuleCheck)
		if err != nil {
			return fmt.Errorf("error on reading mongodb instance security rule %q, %s", rs.Primary.ID, err)
		}
		rules := (*resp)["MongoDBSecurityGroupRule"].([]interface{})
		if len(rules) == 0 {
			return fmt.Errorf("mongodb instance security rule is not exist")
		}

		return nil
	}
}

func testAccCheckMongodbSecurityRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*KsyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ksyun_mongodb_security_rule" {
			securityRuleCheck := make(map[string]interface{})
			securityRuleCheck["InstanceId"] = rs.Primary.ID
			resp, err := client.mongodbconn.ListSecurityGroupRules(&securityRuleCheck)

			if err != nil {
				return fmt.Errorf("error on reading mongodb instance security rule %q, %s", rs.Primary.ID, err)
			}
			rules := (*resp)["MongoDBSecurityGroupRule"].([]interface{})
			if len(rules) > 0 {
				return fmt.Errorf("delete mongodb security rule failure")
			}

			return nil
		}
	}

	return nil
}

const testAccMongodbSecurityRuleConfig = `
resource "ksyun_mongodb_security_rule" "default" {
	instance_id = "InstanceId"
    cidrs = "192.16.10.1/32"
}
`
