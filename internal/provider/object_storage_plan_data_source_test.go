package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccObjectStoragePlanDataSource(t *testing.T) {
	resourcename := "data.utho_object_storage_plan.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "utho_object_storage_plan" "example" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourcename, "pricing.0.disk"),
					resource.TestCheckResourceAttrSet(resourcename, "pricing.0.id"),
					resource.TestCheckResourceAttrSet(resourcename, "pricing.0.monthly"),
					resource.TestCheckResourceAttrSet(resourcename, "pricing.0.name"),
					resource.TestCheckResourceAttrSet(resourcename, "pricing.0.price"),
					resource.TestCheckResourceAttrSet(resourcename, "pricing.0.type"),
				),
			},
		},
	})
}
