package ksyun

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccKsyunSubnetsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSubnetsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_subnets.foo"),
				),
			},
		},
	})
}

const testAccDataSubnetsConfig = `

data "ksyun_subnets" "foo" {
    subnet_types = ["Normal"]
    vpc_ids = ["db00bf66-f904-4749-9ddf-94f170f62d05"]
	output_file = "output_result"
}
`
