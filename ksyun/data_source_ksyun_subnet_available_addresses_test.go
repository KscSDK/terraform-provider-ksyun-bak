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
  ids=["f5e3b70e-493d-4473-8072-9f50640f4ae3"]
  #subnet_id=["f5e3b70e-493d-4473-8072-9f50640f4ae3"]
}
`
