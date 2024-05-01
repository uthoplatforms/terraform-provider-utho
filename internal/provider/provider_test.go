package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// set utho token
	providerConfig = `
provider "utho" {
	token = "UTHO_TOKEN"
}
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"utho": providerserver.NewProtocol6WithError(New("test")()),
	}
)
