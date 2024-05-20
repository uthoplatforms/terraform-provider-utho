package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAutoScalingResource(t *testing.T) {
	resourceName := "utho_auto_scaling.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "utho_auto_scaling" "example" {
	name                = "example-name"
	os_disk_size        = 800
	dcslug              = "inmumbaizone2"
	minsize             = "1"
	maxsize             = "2"
	desiredsize         = "1"
	planid              = "10045"
	planname            = "basic"
	instance_templateid = "none"
	public_ip_enabled   = true
	loadbalancers_id    = ""
	stackid    = "6669341"
	stackimage = "ubuntu-22.04-x86_64"
	vpc_id            = ""
	security_group_id = ""
	target_groups_id  = ""
	policies = [
		{
		name     = "Policy-16H2jh"
		type     = "cpu"
		compare  = "above"
		value    = "80"
		adjust   = "1"
		period   = "5m"
		cooldown = "300"
		}
	]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "example-name"),
					resource.TestCheckResourceAttr(resourceName, "minsize", "1"),
					resource.TestCheckResourceAttr(resourceName, "maxsize", "2"),
					resource.TestCheckResourceAttr(resourceName, "desiredsize", "1"),
					resource.TestCheckResourceAttr(resourceName, "image", "ubuntu-22.04-x86_64"),

					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip_enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "userid"),
					resource.TestCheckResourceAttrSet(resourceName, "dclocation.location"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"loadbalancers_id", "security_group_id", "target_groups_id", "vpc_id"},
			},
		},
	})
}
