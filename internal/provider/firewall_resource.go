package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/uthoplatforms/terraform-provider-utho/api"
)

// implement resource interfaces.
var (
	_ resource.Resource                = &FirewallResource{}
	_ resource.ResourceWithConfigure   = &FirewallResource{}
	_ resource.ResourceWithImportState = &FirewallResource{}
)

// NewFirewallResource is a helper function to simplify the provider implementation.
func NewFirewallResource() resource.Resource {
	return &FirewallResource{}
}

// FirewallResource is the resource implementation.
type FirewallResource struct {
	client *api.Client
}

type FirewallResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	CreatedAt    types.String `tfsdk:"created_at"`
	Rulecount    types.String `tfsdk:"rulecount"`
	Serverscount types.String `tfsdk:"serverscount"`
}

// Metadata returns the resource type name.
func (s *FirewallResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall"
}

// Configure adds the provider configured client to the data source.
func (d *FirewallResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Firewall Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

// Schema defines the schema for the resource.
func (s *FirewallResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "id"},
			"name": schema.StringAttribute{
				Required: true,
				// Requires Replace if the value change
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Name of the firewall",
			},
			"created_at":   schema.StringAttribute{Computed: true, Description: "created_at"},
			"rulecount":    schema.StringAttribute{Computed: true, Description: "rulecount"},
			"serverscount": schema.StringAttribute{Computed: true, Description: "serverscount"},
		},
	}
}

// Import using firewall as the attribute
func (s *FirewallResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("firewall"), req, resp)
}

// Create a new resource.
func (s *FirewallResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "create firewall")
	// Retrieve values from plan
	var plan FirewallResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	firewallRequest := api.FirewallRequest{
		Name: plan.Name.ValueString(),
	}
	tflog.Debug(ctx, "send create firewall request")
	firewall, err := s.client.CreateFirewall(ctx, firewallRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating firewall",
			"Could not create firewall, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = FirewallResourceModel{
		ID:           types.StringValue(firewall.ID),
		Name:         types.StringValue(plan.Name.ValueString()),
		CreatedAt:    types.StringValue(plan.CreatedAt.ValueString()),
		Rulecount:    types.StringValue(plan.Rulecount.ValueString()),
		Serverscount: types.StringValue(plan.Serverscount.ValueString()),
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish create firewall")
}

// Read resource information.
func (s *FirewallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "read firewall")

	// Get current state
	var state FirewallResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "send get firewall request")
	// Get refreshed firewall value from utho
	firewall, err := s.client.GetFirewall(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading utho firewall",
			"Could not read utho firewall "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state = FirewallResourceModel{
		ID:           types.StringValue(firewall.ID),
		Name:         types.StringValue(firewall.Name),
		CreatedAt:    types.StringValue(firewall.CreatedAt),
		Rulecount:    types.StringValue(firewall.Rulecount),
		Serverscount: types.StringValue(firewall.Serverscount),
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish get firewall request")
}

func (s *FirewallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// updating resource is not supported
}

// Delete deletes the resource and removes the Terraform state on success.
func (s *FirewallResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "delete firewall")
	// Get current state
	var state FirewallResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "send delete firewall request")
	// delete firewall
	err := s.client.DeleteFirewall(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing utho firewall",
			"Could not delete utho firewall "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}
