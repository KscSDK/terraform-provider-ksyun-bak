package ksyun

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccKsyunSubnetAllocatedIpAddressesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSubnetAllocatedIpAddressesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_subnet_allocated_ip_addresses.foo"),
				),
			},
		},
	})
}

const testAccDataSubnetAllocatedIpAddressesConfig = `
data "ksyun_subnet_allocated_ip_addresses" "foo" {
  output_file="output_result"
  ids=["d8f6f5dd-b0ee-4106-bf33-52042b70032d"]
  subnet_id=["d8f6f5dd-b0ee-4106-bf33-52042b70032d"]
}
`
