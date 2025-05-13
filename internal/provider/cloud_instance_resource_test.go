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
	dcslug       = "inbangalore"
	image        = "rocky-8.8-x86_64"
	planid       = "10045"
	enablebackup = "false"
	billingcycle = "hourly"
	firewall     = "23432614"
	vpc_id		 = "f1qq22aa-11aa-11dd-8b94-f69f312c0245"
}		  
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "example-name"),
					resource.TestCheckResourceAttr(resourceName, "dclocation.dc", "inbangalore"),
					resource.TestCheckResourceAttr(resourceName, "enablebackup", "false"),
					resource.TestCheckResourceAttr(resourceName, "billingcycle", "hourly"),
					resource.TestCheckResourceAttr(resourceName, "image", "rocky-8.8-x86_64"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", "f1qq22aa-11aa-11dd-8b94-f69f312c0245"),

					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip"),
					resource.TestCheckResourceAttrSet(resourceName, "billingcycle"),
					resource.TestCheckResourceAttrSet(resourceName, "disksize"),
					resource.TestCheckResourceAttrSet(resourceName, "public_network.0.ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "storages.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "vmcost"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"powerstatus", "root_password", "ha", "planid", "firewall"},
			},
		},
	})
}
