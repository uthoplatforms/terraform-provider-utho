package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccountsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "utho_account" "example" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.utho_account.example", "user.verify", "1"),

					resource.TestCheckResourceAttrSet("data.utho_account.example", "user.id"),
					resource.TestCheckResourceAttrSet("data.utho_account.fullname", "user.fullname"),
					resource.TestCheckResourceAttrSet("data.utho_account.email", "user.email"),
				),
			},
		},
	})
}
