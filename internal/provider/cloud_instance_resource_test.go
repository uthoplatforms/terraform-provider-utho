package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudInstanceResource(t *testing.T) {
	resourceName := "utho_cloud_instance.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "utho_cloud_instance" "example" {
	name = "example-name"
	# country slug
	dcslug        = "inmumbaizone2"
	image         = "ubuntu-22.04-x86_64"
	planid        = "10045"
	enablebackup  = "false"
	billingcycle  = "hourly"
	firewall      = "23432614"
	vpc_id		  = "4de5f07a-f51c-4323-b39a-ef66130e1bd9"
	root_password = "2uDsQ1$Ioqa@uFj"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "example-name"),
					resource.TestCheckResourceAttr(resourceName, "dclocation.dc", "inmumbaizone2"),
					resource.TestCheckResourceAttr(resourceName, "enablebackup", "false"),
					resource.TestCheckResourceAttr(resourceName, "billingcycle", "hourly"),
					resource.TestCheckResourceAttr(resourceName, "image", "ubuntu-22.04-x86_64"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", "4de5f07a-f51c-4323-b39a-ef66130e1bd9"),

					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "billingcycle"),
					resource.TestCheckResourceAttrSet(resourceName, "disksize"),
					resource.TestCheckResourceAttrSet(resourceName, "storages.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "vmcost"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"powerstatus",
					"root_password",
					"ha",
					"planid",
					"firewall",
					"ip",
					"public_network",
					"vmcost",
					"disksize",
					"vpc_id",
				},
			},
		},
	})
}
