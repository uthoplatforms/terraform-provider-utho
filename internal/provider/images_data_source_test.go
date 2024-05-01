package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImagesDataSource(t *testing.T) {
	resourceName := "data.utho_images.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "utho_images" "example" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "images.0.distribution"),
					resource.TestCheckResourceAttrSet(resourceName, "images.0.image"),
					resource.TestCheckResourceAttrSet(resourceName, "images.0.version"),
				),
			},
		},
	})
}
