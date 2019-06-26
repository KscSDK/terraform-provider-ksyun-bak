package ksyun

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccKsyunVPCsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataVPCsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_vpcs.foo"),
				),
			},
		},
	})
}

const testAccDataVPCsConfig = `

data "ksyun_vpcs" "foo" {
    ids = ["b6533caa-2f78-4d07-806e-e5dfdd03f331"]
	output_file = "output_result"
}
`
