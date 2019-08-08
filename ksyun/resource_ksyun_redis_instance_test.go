package ksyun

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
	"github.com/hashicorp/terraform/terraform"
	"fmt"
	"os"
	"strings"
	"github.com/pkg/errors"
)

func TestAccKcs_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKcsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKcsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKcsInstanceExists("ksyun_redis_instance.default"),
				),
			},
		},
	})
}

func testAccCheckKcsInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find resource or data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("kcs instance is create failure")
		}
		return nil
	}
}

func testAccCheckKcsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*KsyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ksyun_redis_instance" {
			instanceCheck := make(map[string]interface{})
			instanceCheck["CacheId"] = rs.Primary.ID
			ptr, err := client.kcsv1conn.DescribeCacheCluster(&instanceCheck)
			// Verify the error is what we want
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "cannot be found") {
					return nil
				}
				return err
			}
			if ptr != nil {
				if (*ptr)["Data"] != nil {
					return errors.New("delete instance failure")
				}
			}
		}
	}

	return nil
}

const testAccKcsConfig = `
variable "available_zone" {
  default = "cn-shanghai-3a"
}

variable "protocol" {
  default = "4.0"
}

resource "ksyun_redis_instance" "default" {
  available_zone        = "${var.available_zone}"
  name                  = "MyRedisInstance1101"
  mode                  = 2
  capacity              = 1
  net_type              = 2
  vnet_id               = "bb94ca2d-5865-44b4-a0b4-3c7117185459"
  vpc_id                = "585c8a05-47f1-4934-a556-1aae96841692"
  bill_type             = 5
  duration              = ""
  duration_unit         = ""
  pass_word             = "Shiwo1101"
  iam_project_id        = "0"
  protocol              = "${var.protocol}"
  reset_all_parameters  = false
  parameters = {
    "appendonly"                  = "no",
    "appendfsync"                 = "everysec",
    "maxmemory-policy"            = "volatile-lru",
    "hash-max-ziplist-entries"    = "513",
    "zset-max-ziplist-entries"    = "129",
    "list-max-ziplist-size"       = "-2",
    "hash-max-ziplist-value"      = "64",
    "notify-keyspace-events"      = "",
    "zset-max-ziplist-value"      = "64",
    "maxmemory-samples"           = "5",
    "set-max-intset-entries"      = "512",
    "timeout"                     = "600",
  }
}
`