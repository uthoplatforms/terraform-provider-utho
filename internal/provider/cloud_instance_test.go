package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudInstanceResource(t *testing.T) {
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
}		  
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("utho_cloud_instance.example", "name", "example-name"),
					resource.TestCheckResourceAttr("utho_cloud_instance.example", "dclocation.dc", "inbangalore"),
					resource.TestCheckResourceAttr("utho_cloud_instance.example", "enablebackup", "false"),
					resource.TestCheckResourceAttr("utho_cloud_instance.example", "billingcycle", "hourly"),
					resource.TestCheckResourceAttr("utho_cloud_instance.example", "image", "rocky-8.8-x86_64"),

					resource.TestCheckResourceAttrSet("utho_cloud_instance.example", "id"),
					resource.TestCheckResourceAttrSet("utho_cloud_instance.example", "ip"),
					resource.TestCheckResourceAttrSet("utho_cloud_instance.example", "billingcycle"),
					resource.TestCheckResourceAttrSet("utho_cloud_instance.example", "disksize"),
					resource.TestCheckResourceAttrSet("utho_cloud_instance.example", "public_network.0.ip_address"),
					resource.TestCheckResourceAttrSet("utho_cloud_instance.example", "storages.0.id"),
					resource.TestCheckResourceAttrSet("utho_cloud_instance.example", "vmcost"),
				),
			},
			{
				ResourceName:      "utho_cloud_instance.example",
				ImportState:       true,
				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"powerstatus", "root_password", "ha", "planid"},
			},
		},
	})
}
