package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccObjectStoragePlanDataSource(t *testing.T) {
	resourceName := "data.utho_object_storage_plan.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "utho_object_storage_plan" "example" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "pricing.0.disk"),
					resource.TestCheckResourceAttrSet(resourceName, "pricing.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "pricing.0.monthly"),
					resource.TestCheckResourceAttrSet(resourceName, "pricing.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "pricing.0.price"),
					resource.TestCheckResourceAttrSet(resourceName, "pricing.0.type"),
				),
			},
		},
	})
}
