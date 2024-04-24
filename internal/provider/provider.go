package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/uthoterraform/terraform-provider-utho/api"
)

var _ provider.Provider = &uthoProvider{}

// uthoProvider is the provider implementation.
type (
	uthoProvider struct {
		version string
	}

	uthoProviderModel struct {
		Token types.String `tfsdk:"token"`
	}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &uthoProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *uthoProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "utho"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *uthoProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Utho token",
			},
		},
	}
}

func (p *uthoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Token client")
	// Retrieve provider data from configuration
	var config uthoProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Token API Token",
			"The provider cannot create the Token API client as there is an unknown configuration value for the Token API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the UTHO_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	token := os.Getenv("UTHO_TOKEN")

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing utho API Token",
			"The provider cannot create the utho API client as there is a missing or empty value for the utho API token. "+
				"Set the token value in the configuration or use the UTHO_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "utho_token", token)

	tflog.Debug(ctx, "Creating Token client")

	// Create a new Token client using the configuration values
	client := api.NewClient(token)

	// Make the Token client available during DataSource and Resource

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured utho client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *uthoProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAccountDataSource,
		NewImagesDataSource,
		NewObjectStoragePlanDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *uthoProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDomainResource,
		NewVpcResource,
	}
}
