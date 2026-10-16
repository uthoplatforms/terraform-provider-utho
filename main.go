package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/uthoplatforms/terraform-provider-utho/internal/provider"
)

func main() {
	err := providerserver.Serve(
		context.Background(),
		provider.New(),
		providerserver.ServeOpts{
			Address: "registry.terraform.io/uthoplatforms/utho",
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
