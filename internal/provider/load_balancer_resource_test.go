package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLoadBalancerResource(t *testing.T) {
	resourceName := "utho_loadbalancer.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "utho_loadbalancer" "example" {
  dcslug          = "inmumbaizone2"
  name            = "example-utho"
  type            = "application"
  vpc_id          = "4de5f07a-f51c-4323-b39a-ef66130e1bd9"
  firewall        = "23432614"
  cpu_model       = "amd"
  enable_publicip = "true"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "dcslug", "inmumbaizone2"),
					resource.TestCheckResourceAttr(resourceName, "name", "example-utho"),
					resource.TestCheckResourceAttr(resourceName, "type", "application"),

					resource.TestCheckResourceAttrSet(resourceName, "dcslug"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
				),
			},
			// {
			// 	ResourceName:      resourceName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}
