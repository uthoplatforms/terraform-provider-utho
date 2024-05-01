package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccountsDataSource(t *testing.T) {
	resourceName := "data.utho_account.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "utho_account" "example" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user.verify", "1"),

					resource.TestCheckResourceAttrSet(resourceName, "user.id"),
					resource.TestCheckResourceAttrSet(resourceName, "user.fullname"),
					resource.TestCheckResourceAttrSet(resourceName, "user.email"),
				),
			},
		},
	})
}
