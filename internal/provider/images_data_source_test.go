package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImagesDataSource(t *testing.T) {
	resourcename := "data.utho_images.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "utho_images" "example" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourcename, "images.0.distribution"),
					resource.TestCheckResourceAttrSet(resourcename, "images.0.image"),
					resource.TestCheckResourceAttrSet(resourcename, "images.0.version"),
				),
			},
		},
	})
}
