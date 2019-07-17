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
					testAccCheckIDExists("data.ksyun_SubnetAllocatedIpAddresses.foo"),
				),
			},
		},
	})
}

const testAccDataSubnetAllocatedIpAddressesConfig = `
data "ksyun_SubnetAllocatedIpAddresses" "foo" {
  output_file="output_result"
  ids=[]
  state=""
  vpc_id=[]
}
`
