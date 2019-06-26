package ksyun

import (
	//"aws-sdk-go/aws"
	"fmt"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAwsEip_Instance(t *testing.T) {
	dataSourceName := "data.ksyun_eip.eip1234"
	resourceName := "ksyun_eip.eip1234"

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsEipConfigInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "next_token", resourceName, "next_token"),
					resource.TestCheckResourceAttrPair(dataSourceName, "max_results", resourceName, "max_results"),
					resource.TestCheckResourceAttrPair(dataSourceName, "line_id", resourceName, "line_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "band_width", resourceName, "band_width"),
					resource.TestCheckResourceAttrPair(dataSourceName, "charge_type", resourceName, "charge_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "purchase_time", resourceName, "purchase_time"),
				),
			},
		},
	})
}

const testAccDataSourceAwsEipConfigInstance = `
resource "ksyun_eip" "eip1234" {
  next_token="1"
  max_results="5"
  line_id="a2403858-2550-4612-850c-ea840fa343f9"
  band_width="10"
  charge_type="PostPaidByDay"
  purchase_time="1"
}
`

func TestAccAWSEIP_basic(t *testing.T) {
	var conf map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		IDRefreshName: "ksyun_eip.bar",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEIPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKSCEIPExists("aws_eip.bar", &conf),
				),
			},
		},
	})
}
func testAccCheckKSCEIPExists(n string, res *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		log.Println("rssssssss:", rs)

		conn := testAccProvider.Meta().(*KsyunClient).eipconn

		var req map[string]interface{}
		describe, err := conn.DescribeAddresses(&req)
		if err != nil {
			return err
		}

		*res = *describe

		return nil
	}
}

// Regression test for https://github.com/hashicorp/terraform/issues/3429 (now
// https://github.com/terraform-providers/terraform-provider-aws/issues/42)
const testAccAWSEIPConfig = `
resource "ksyun_eip" "bar" {
}
`
