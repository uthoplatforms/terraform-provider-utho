package datasources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ══════════════════════════════════════════════════════════════
// data.utho_firewalls
// ══════════════════════════════════════════════════════════════

type FirewallsDataSource struct{ client *client.Client }

type FirewallsModel struct {
	Firewalls []FirewallItem `tfsdk:"firewalls"`
}

type FirewallItem struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	CreatedAt    types.String `tfsdk:"created_at"`
	RuleCount    types.String `tfsdk:"rule_count"`
	ServersCount types.String `tfsdk:"servers_count"`
}

func NewFirewallsDataSource() datasource.DataSource { return &FirewallsDataSource{} }

func (d *FirewallsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewalls"
}

func (d *FirewallsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all security groups (firewalls) in your account.",
		Attributes: map[string]schema.Attribute{
			"firewalls": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":            schema.StringAttribute{Computed: true, Description: "Security group ID. Use this as firewall in utho_cloud."},
						"name":          schema.StringAttribute{Computed: true, Description: "Security group name."},
						"created_at":    schema.StringAttribute{Computed: true, Description: "Creation timestamp."},
						"rule_count":    schema.StringAttribute{Computed: true, Description: "Number of rules in this security group."},
						"servers_count": schema.StringAttribute{Computed: true, Description: "Number of instances attached to this security group."},
					},
				},
			},
		},
	}
}

func (d *FirewallsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", "Expected *client.Client")
		return
	}
	d.client = c
}

func (d *FirewallsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	respBytes, err := d.client.Get("/firewall")
	if err != nil {
		resp.Diagnostics.AddError("Error fetching firewalls", fmt.Sprintf("%s", err))
		return
	}

	var result struct {
		Firewalls []struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			CreatedAt    string `json:"created_at"`
			RuleCount    string `json:"rulecount"`
			ServersCount string `json:"serverscount"`
		} `json:"firewalls"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		resp.Diagnostics.AddError("Error parsing firewalls", fmt.Sprintf("%s", err))
		return
	}

	var firewalls []FirewallItem
	for _, f := range result.Firewalls {
		firewalls = append(firewalls, FirewallItem{
			ID:           types.StringValue(f.ID),
			Name:         types.StringValue(f.Name),
			CreatedAt:    types.StringValue(f.CreatedAt),
			RuleCount:    types.StringValue(f.RuleCount),
			ServersCount: types.StringValue(f.ServersCount),
		})
	}
	if firewalls == nil {
		firewalls = []FirewallItem{}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &FirewallsModel{Firewalls: firewalls})...)
}

// ══════════════════════════════════════════════════════════════
// data.utho_ssh_keys
// ══════════════════════════════════════════════════════════════

type SSHKeysDataSource struct{ client *client.Client }

type SSHKeysModel struct {
	Keys []SSHKeyItem `tfsdk:"keys"`
}

type SSHKeyItem struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func NewSSHKeysDataSource() datasource.DataSource { return &SSHKeysDataSource{} }

func (d *SSHKeysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_keys"
}

func (d *SSHKeysDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all SSH keys in your account.",
		Attributes: map[string]schema.Attribute{
			"keys": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":         schema.StringAttribute{Computed: true, Description: "SSH key ID. Use this as sshkeys in utho_cloud."},
						"name":       schema.StringAttribute{Computed: true, Description: "SSH key name."},
						"created_at": schema.StringAttribute{Computed: true, Description: "Creation timestamp."},
					},
				},
			},
		},
	}
}

func (d *SSHKeysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", "Expected *client.Client")
		return
	}
	d.client = c
}

func (d *SSHKeysDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	respBytes, err := d.client.Get("/key")
	if err != nil {
		resp.Diagnostics.AddError("Error fetching SSH keys", fmt.Sprintf("%s", err))
		return
	}

	var result struct {
		Keys []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			CreatedAt string `json:"created_at"`
		} `json:"key"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		resp.Diagnostics.AddError("Error parsing SSH keys", fmt.Sprintf("%s", err))
		return
	}

	var keys []SSHKeyItem
	for _, k := range result.Keys {
		keys = append(keys, SSHKeyItem{
			ID:        types.StringValue(k.ID),
			Name:      types.StringValue(k.Name),
			CreatedAt: types.StringValue(k.CreatedAt),
		})
	}
	if keys == nil {
		keys = []SSHKeyItem{}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &SSHKeysModel{Keys: keys})...)
}

// ══════════════════════════════════════════════════════════════
// data.utho_vpc_subnets
// ══════════════════════════════════════════════════════════════

type VPCSubnetsDataSource struct{ client *client.Client }

type VPCSubnetsModel struct {
	VPCID   types.String        `tfsdk:"vpc_id"`
	Subnets []VPCSubnetFullItem `tfsdk:"subnets"`
}

type VPCSubnetFullItem struct {
	ID             types.String `tfsdk:"id"`
	UUID           types.String `tfsdk:"uuid"`
	Name           types.String `tfsdk:"name"`
	Network        types.String `tfsdk:"network"`
	Size           types.String `tfsdk:"size"`
	SubnetType     types.String `tfsdk:"subnet_type"`
	AssignPublicIP types.String `tfsdk:"assign_publicip"`
	Status         types.String `tfsdk:"status"`
	Gateway        types.String `tfsdk:"gateway"`
	DCSlug         types.String `tfsdk:"dcslug"`
}

func NewVPCSubnetsDataSource() datasource.DataSource { return &VPCSubnetsDataSource{} }

func (d *VPCSubnetsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_subnets"
}

func (d *VPCSubnetsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all subnets within a specific VPC.",
		Attributes: map[string]schema.Attribute{
			"vpc_id": schema.StringAttribute{Required: true, Description: "VPC UUID. Get this from data.utho_vpcs."},
			"subnets": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":              schema.StringAttribute{Computed: true, Description: "Subnet numeric ID. Use this as vpc in utho_cloud."},
						"uuid":            schema.StringAttribute{Computed: true, Description: "Subnet UUID."},
						"name":            schema.StringAttribute{Computed: true, Description: "Subnet name."},
						"network":         schema.StringAttribute{Computed: true, Description: "Subnet base network address."},
						"size":            schema.StringAttribute{Computed: true, Description: "Subnet prefix length."},
						"subnet_type":     schema.StringAttribute{Computed: true, Description: "public or private."},
						"assign_publicip": schema.StringAttribute{Computed: true, Description: "1 if public IPs are auto-assigned."},
						"status":          schema.StringAttribute{Computed: true, Description: "Subnet status."},
						"gateway":         schema.StringAttribute{Computed: true, Description: "Subnet gateway IP."},
						"dcslug":          schema.StringAttribute{Computed: true, Description: "Data center slug."},
					},
				},
			},
		},
	}
}

func (d *VPCSubnetsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", "Expected *client.Client")
		return
	}
	d.client = c
}

func (d *VPCSubnetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config VPCSubnetsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	respBytes, err := d.client.Get(fmt.Sprintf("/vpc/%s", config.VPCID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error fetching VPC subnets", fmt.Sprintf("%s", err))
		return
	}

	var result struct {
		Subnets []struct {
			ID             string `json:"id"`
			UUID           string `json:"uuid"`
			Name           string `json:"name"`
			Network        string `json:"network"`
			Size           string `json:"size"`
			SubnetType     string `json:"subnet_type"`
			AssignPublicIP string `json:"assign_publicip"`
			Status         string `json:"status"`
			Gateway        string `json:"gateway"`
			DCSlug         string `json:"dcslug"`
		} `json:"subnets"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		resp.Diagnostics.AddError("Error parsing VPC subnets", fmt.Sprintf("%s", err))
		return
	}

	var subnets []VPCSubnetFullItem
	for _, s := range result.Subnets {
		subnets = append(subnets, VPCSubnetFullItem{
			ID:             types.StringValue(s.ID),
			UUID:           types.StringValue(s.UUID),
			Name:           types.StringValue(s.Name),
			Network:        types.StringValue(s.Network),
			Size:           types.StringValue(s.Size),
			SubnetType:     types.StringValue(s.SubnetType),
			AssignPublicIP: types.StringValue(s.AssignPublicIP),
			Status:         types.StringValue(s.Status),
			Gateway:        types.StringValue(s.Gateway),
			DCSlug:         types.StringValue(s.DCSlug),
		})
	}
	if subnets == nil {
		subnets = []VPCSubnetFullItem{}
	}

	config.Subnets = subnets
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ══════════════════════════════════════════════════════════════
// data.utho_billing_cycles
// ══════════════════════════════════════════════════════════════

type BillingCyclesDataSource struct{ client *client.Client }

type BillingCyclesModel struct {
	Product types.String   `tfsdk:"product"`
	Cycles  []types.String `tfsdk:"cycles"`
}

func NewBillingCyclesDataSource() datasource.DataSource { return &BillingCyclesDataSource{} }

func (d *BillingCyclesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_billing_cycles"
}

func (d *BillingCyclesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List available billing cycles for a Utho product.",
		Attributes: map[string]schema.Attribute{
			"product": schema.StringAttribute{
				Required:    true,
				Description: "Product name: cloud, kubernetes, database, objectstorage, etc.",
			},
			"cycles": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of valid billing cycle values for this product.",
			},
		},
	}
}

func (d *BillingCyclesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", "Expected *client.Client")
		return
	}
	d.client = c
}

func (d *BillingCyclesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config BillingCyclesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	respBytes, err := d.client.Get(fmt.Sprintf("/billingcycles/get?product=%s", config.Product.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error fetching billing cycles", fmt.Sprintf("%s", err))
		return
	}

	var result struct {
		Cycles []string `json:"cycles"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		resp.Diagnostics.AddError("Error parsing billing cycles", fmt.Sprintf("%s", err))
		return
	}

	var cycles []types.String
	for _, c := range result.Cycles {
		cycles = append(cycles, types.StringValue(c))
	}
	if cycles == nil {
		cycles = []types.String{}
	}

	config.Cycles = cycles
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
