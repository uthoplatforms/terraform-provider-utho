package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsRecordResource(t *testing.T) {
	resourceName := "utho_dns_record.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "utho_domain" "example" {
	domain = "example-test-utho.com"
}
resource "utho_dns_record" "example" {
	domain   = utho_domain.example.domain
	type     = "A"
	hostname = "subdomain"
	value    = "1.1.1.1"
	ttl      = "65444"
	porttype = "TCP"
	port     = "5060"
	priority = "10"
	weight    = "100"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "domain", "example-test-utho.com"),
					resource.TestCheckResourceAttr(resourceName, "hostname", "subdomain"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "65444"),
					resource.TestCheckResourceAttr(resourceName, "type", "A"),
					resource.TestCheckResourceAttr(resourceName, "value", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "65444"),
					resource.TestCheckResourceAttr(resourceName, "porttype", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "port", "5060"),
					resource.TestCheckResourceAttr(resourceName, "priority", "10"),
					resource.TestCheckResourceAttr(resourceName, "weight", "100"),

					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}
