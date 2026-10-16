package provider

import (
	"context"
	"os"

	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
	"github.com/uthoplatforms/terraform-provider-utho/internal/datasources"
	"github.com/uthoplatforms/terraform-provider-utho/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type uthoProvider struct{}

type uthoProviderModel struct {
	APIKey  types.String `tfsdk:"api_key"`
	BaseURL types.String `tfsdk:"base_url"`
}

func New() func() provider.Provider {
	return func() provider.Provider {
		return &uthoProvider{}
	}
}

func (p *uthoProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "utho"
	resp.Version = "1.0.0"
}

func (p *uthoProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Utho Cloud resources.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Utho API key. Can also be set via UTHO_API_KEY environment variable.",
			},
			"base_url": schema.StringAttribute{
				Optional:    true,
				Description: "Override the Utho API base URL. Defaults to https://api.utho.com/v2",
			},
		},
	}
}

func (p *uthoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config uthoProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := os.Getenv("UTHO_API_KEY")
	if !config.APIKey.IsNull() && config.APIKey.ValueString() != "" {
		apiKey = config.APIKey.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key",
			"Set api_key in provider block or export UTHO_API_KEY=your-key",
		)
		return
	}

	baseURL := config.BaseURL.ValueString()
	uthoClient := client.NewClientWithURL(apiKey, baseURL)
	resp.DataSourceData = uthoClient
	resp.ResourceData = uthoClient
}

func (p *uthoProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewCloudResource,                   // utho_cloud
		resources.NewCloudPowerResource,              // utho_cloud_power
		resources.NewCloudFirewallResource,           //utho_cloud_firewall
		resources.NewCloudStorageResource,            // utho_cloud_storage
		resources.NewCloudEBSResource,                // utho_cloud_ebs
		resources.NewCloudSnapshotResource,           // utho_cloud_snapshot
		resources.NewCloudISOResource,                // utho_cloud_iso
		resources.NewCloudResizeResource,             // utho_cloud_resize
		resources.NewCloudPublicIPResource,           // utho_cloud_public_ip
		resources.NewCloudVPCResource,                // utho_cloud_vpc
		resources.NewVPCResource,                     // utho_vpc
		resources.NewSubnetResource,                  // utho_subnet
		resources.NewNATGatewayResource,              // utho_nat_gateway
		resources.NewRouteTableResource,              // utho_route_table
		resources.NewRouteResource,                   // utho_route
		resources.NewElasticIPResource,               // utho_elastic_ip
		resources.NewVPCPeeringResource,              // utho_vpc_peering
		resources.NewSSHKeyResource,                  // utho_ssh_key
		resources.NewFirewallResource,                // utho_firewall
		resources.NewFirewallRuleResource,            // utho_firewall_rule
		resources.NewFirewallServerResource,          // utho_firewall_server
		resources.NewLoadBalancerResource,            // utho_loadbalancer
		resources.NewLBFrontendResource,              // utho_loadbalancer_frontend
		resources.NewLBBackendResource,               // utho_loadbalancer_backend
		resources.NewLBACLResource,                   // utho_loadbalancer_acl
		resources.NewLBSettingsResource,              // utho_loadbalancer_settings
		resources.NewObjectStorageResource,           // utho_object_storage
		resources.NewObjectStoragePermissionResource, // utho_object_storage_permission
		resources.NewObjectStorageKeyResource,        // utho_object_storage_key
		resources.NewDNSZoneResource,                 // utho_dns_zone
		resources.NewDNSRecordResource,               // utho_dns_record
		resources.NewSSLCertificateResource,          // utho_ssl_certificate
		resources.NewKubernetesResource,              // utho_kubernetes
		resources.NewKubernetesNodePoolResource,      // utho_kubernetes_node_pool
		resources.NewDatabaseResource,                // utho_database
		resources.NewDatabaseDBResource,              // utho_database_db
		resources.NewDatabaseUserResource,            // utho_database_user
		resources.NewDatabasePoolResource,            // utho_database_pool
		resources.NewAutoScalingResource,             // utho_autoscaling
		resources.NewScalingPolicyResource,           // utho_autoscaling_policy
		resources.NewScalingScheduleResource,         // utho_autoscaling_schedule
		resources.NewAPITokenResource,                // utho_api_token
		resources.NewAlertContactResource,            // utho_alert_contact
		resources.NewAlertResource,                   // utho_alert
		resources.NewIPSecResource,                   // utho_ipsec
		resources.NewIPSecConnectionResource,         // utho_ipsec_connection
		resources.NewIAMUserResource,                 // utho_iam_user
		resources.NewProjectResource,                 // utho_project
		resources.NewProjectMemberResource,           // utho_project_member
		resources.NewContainerRegistryResource,       // utho_container_registry
		resources.NewRegistryRobotResource,           // utho_container_registry_robot
		resources.NewRegistryWebhookResource,         // utho_container_registry_webhook
		resources.NewRegistryImmutableRuleResource,   // utho_container_registry_immutable_rule
		resources.NewEBSResource,                     // utho_ebs
		resources.NewEBSAttachmentResource,           // utho_ebs_attachment

	}
}

func (p *uthoProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewCloudsDataSource,     // data.utho_clouds
		datasources.NewKubeconfigDataSource, // data.utho_kubernetes_config
		datasources.NewCloudDCZonesDataSource,
		datasources.NewCloudPlansDataSource,
		datasources.NewCloudImagesDataSource,
		datasources.NewCloudSnapshotsDataSource,
		datasources.NewCloudISOsDataSource,
		datasources.NewVPCsDataSource,
		datasources.NewFirewallsDataSource,
		datasources.NewSSHKeysDataSource,
		datasources.NewVPCSubnetsDataSource,
		datasources.NewBillingCyclesDataSource,
		datasources.NewDatabasePlansDataSource,        // data.utho_database_plans
		datasources.NewAutoScalingsDataSource,         // data.utho_autoscalings
		datasources.NewTargetGroupsDataSource,         // data.utho_target_groups
		datasources.NewBillingUsageDataSource,         // data.utho_billing_usage
		datasources.NewBillingInvoicesDataSource,      // data.utho_billing_invoices
		datasources.NewBillingCostByProjectDataSource, // data.utho_billing_cost_by_project
		datasources.NewEBSDCZonesDataSource,           // data.utho_ebs_dczones
	}
}
