package resources

import (
	"context"
	"fmt"

	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ══════════════════════════════════════════════════════════════
// FIREWALL RESOURCE — utho_firewall
// ══════════════════════════════════════════════════════════════

type FirewallResource struct{ client *client.Client }

type FirewallModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	CreatedAt    types.String `tfsdk:"created_at"`
	RuleCount    types.String `tfsdk:"rule_count"`
	ServersCount types.String `tfsdk:"servers_count"`
}

func NewFirewallResource() resource.Resource { return &FirewallResource{} }

func (r *FirewallResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall"
}

func (r *FirewallResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage Utho Firewall (security groups).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Firewall name.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the firewall was created.",
			},
			"rule_count": schema.StringAttribute{
				Computed:    true,
				Description: "Number of rules in this firewall.",
			},
			"servers_count": schema.StringAttribute{
				Computed:    true,
				Description: "Number of servers attached to this firewall.",
			},
		},
	}
}

func (r *FirewallResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", "Expected *client.Client")
		return
	}
	r.client = c
}

func (r *FirewallResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FirewallModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateFirewall(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating firewall", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	plan.CreatedAt = types.StringValue("")
	plan.RuleCount = types.StringValue("0")
	plan.ServersCount = types.StringValue("0")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *FirewallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FirewallModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fw, err := r.client.GetFirewall(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading firewall", fmt.Sprintf("%s", err))
		return
	}
	if fw == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fw.Name)
	state.CreatedAt = types.StringValue(fw.CreatedAt)
	state.RuleCount = types.StringValue(fw.RuleCount)
	state.ServersCount = types.StringValue(fw.ServersCount)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *FirewallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Firewalls cannot be updated in place.")
}

func (r *FirewallResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FirewallModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteFirewall(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting firewall", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// FIREWALL RULE RESOURCE — utho_firewall_rule
// ══════════════════════════════════════════════════════════════

type FirewallRuleResource struct{ client *client.Client }

type FirewallRuleModel struct {
	ID          types.String `tfsdk:"id"`
	FirewallID  types.String `tfsdk:"firewall_id"`
	Type        types.String `tfsdk:"type"`
	Service     types.String `tfsdk:"service"`
	Protocol    types.String `tfsdk:"protocol"`
	Port        types.String `tfsdk:"port"`
	PortRange   types.String `tfsdk:"port_range"`
	Addresses   types.String `tfsdk:"addresses"`
	SourceRange types.String `tfsdk:"source_range"`
}

func NewFirewallRuleResource() resource.Resource { return &FirewallRuleResource{} }

func (r *FirewallRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_rule"
}

func (r *FirewallRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Add and manage rules in a Utho Firewall.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"firewall_id": schema.StringAttribute{
				Required:    true,
				Description: "Firewall ID to add this rule to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Rule direction: incoming or outgoing.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"service": schema.StringAttribute{
				Required:    true,
				Description: "Service name (e.g. SSH, HTTP, CUSTOM).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "Protocol: tcp, udp, or icmp.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"port": schema.StringAttribute{
				Required:    true,
				Description: "Port number or ALL.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"port_range": schema.StringAttribute{
				Required:    true,
				Description: "Port range (e.g. 8000-9000) or same as port for single port.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"addresses": schema.StringAttribute{
				Required:    true,
				Description: "CIDR address range (e.g. 0.0.0.0/0).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_range": schema.StringAttribute{
				Required:    true,
				Description: "Source CIDR range (e.g. 0.0.0.0/0).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *FirewallRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", "Expected *client.Client")
		return
	}
	r.client = c
}

func (r *FirewallRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FirewallRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID, err := r.client.AddFirewallRule(
		plan.FirewallID.ValueString(),
		&client.FirewallRuleRequest{
			Type:        plan.Type.ValueString(),
			Service:     plan.Service.ValueString(),
			Protocol:    plan.Protocol.ValueString(),
			Port:        plan.Port.ValueString(),
			PortRange:   plan.PortRange.ValueString(),
			Addresses:   plan.Addresses.ValueString(),
			SourceRange: plan.SourceRange.ValueString(),
		},
	)
	if err != nil {
		resp.Diagnostics.AddError("Error adding firewall rule", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(ruleID)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *FirewallRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *FirewallRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Firewall rules cannot be updated. Destroy and recreate.")
}

func (r *FirewallRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FirewallRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteFirewallRule(state.FirewallID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting firewall rule", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// FIREWALL SERVER RESOURCE — utho_firewall_server
// ══════════════════════════════════════════════════════════════

type FirewallServerResource struct{ client *client.Client }

type FirewallServerModel struct {
	ID         types.String `tfsdk:"id"`
	FirewallID types.String `tfsdk:"firewall_id"`
	CloudID    types.String `tfsdk:"cloud_id"`
}

func NewFirewallServerResource() resource.Resource { return &FirewallServerResource{} }

func (r *FirewallServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_server"
}

func (r *FirewallServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Attach or detach a cloud instance from a Utho Firewall.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"firewall_id": schema.StringAttribute{
				Required:    true,
				Description: "Firewall ID to attach the server to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cloud_id": schema.StringAttribute{
				Required:    true,
				Description: "Cloud instance ID to attach.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *FirewallServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", "Expected *client.Client")
		return
	}
	r.client = c
}

func (r *FirewallServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FirewallServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.AttachFirewallServer(plan.FirewallID.ValueString(), plan.CloudID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error attaching server to firewall", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", plan.FirewallID.ValueString(), plan.CloudID.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *FirewallServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *FirewallServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Firewall server attachments cannot be updated.")
}

func (r *FirewallServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FirewallServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DetachFirewallServer(state.FirewallID.ValueString(), state.CloudID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error detaching server from firewall", fmt.Sprintf("%s", err))
	}
}
