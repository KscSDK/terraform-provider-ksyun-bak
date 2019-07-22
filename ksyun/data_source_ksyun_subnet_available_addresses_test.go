package ksyun

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccKsyunSubnetAvailableAddressesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSubnetAvailableAddressesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_subnet_available_addresses.foo"),
				),
			},
		},
	})
}

const testAccDataSubnetAvailableAddressesConfig = `
data "ksyun_subnet_available_addresses" "foo" {
  output_file="output_result"
  ids=[]
  subnet_id=["d8f6f5dd-b0ee-4106-bf33-52042b70032d"]
}
`
