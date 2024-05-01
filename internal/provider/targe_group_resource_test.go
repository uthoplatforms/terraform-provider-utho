package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTargetGroupResource(t *testing.T) {
	resourceName := "utho_target_group.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "utho_target_group" "example" {
	name                  = "example-utho"
	protocol              = "HTTP"
	port                  = "12"
	health_check_path     = "1"
	health_check_protocol = "HTTP"
	health_check_timeout  = "1"
	unhealthy_threshold   = "1"
	health_check_interval = "1"
	healthy_threshold     = "1"
	targets = [
		{
			ip               = "103.146.242.55"
			backend_port     = "12"
			backend_protocol = "HTTP"
		},
		{
			ip               = "103.146.200.55"
			backend_port     = "15"
			backend_protocol = "HTTPS"
		}
	]
	}		  
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "example-utho"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "port", "12"),
					resource.TestCheckResourceAttr(resourceName, "health_check_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "healthy_threshold", "1"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.ip", "103.146.242.55"),
					resource.TestCheckResourceAttr(resourceName, "targets.1.ip", "103.146.200.55"),

					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.1.id"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.0.cloudid"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.1.cloudid"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
