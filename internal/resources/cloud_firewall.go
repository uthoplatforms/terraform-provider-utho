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

// ── Resource struct ───────────────────────────────────────────────────────

type CloudFirewallResource struct {
	client *client.Client
}

// ── Model ─────────────────────────────────────────────────────────────────

type CloudFirewallModel struct {
	ID         types.String `tfsdk:"id"`
	CloudID    types.String `tfsdk:"cloud_id"`
	FirewallID types.String `tfsdk:"firewall_id"`
}

// ── Constructor ───────────────────────────────────────────────────────────

func NewCloudFirewallResource() resource.Resource {
	return &CloudFirewallResource{}
}

// ── Metadata ──────────────────────────────────────────────────────────────

func (r *CloudFirewallResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_firewall"
}

// ── Schema ────────────────────────────────────────────────────────────────

func (r *CloudFirewallResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Attach or detach a security group from a Utho Cloud instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for this firewall attachment.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloud_id": schema.StringAttribute{
				Required:    true,
				Description: "The cloud instance ID to attach the security group to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"firewall_id": schema.StringAttribute{
				Required:    true,
				Description: "The security group ID to attach.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// ── Configure ─────────────────────────────────────────────────────────────

func (r *CloudFirewallResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	uthoClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", "Expected *client.Client")
		return
	}
	r.client = uthoClient
}

// ── CREATE — attaches firewall to cloud instance ───────────────────────────

func (r *CloudFirewallResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudFirewallModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.AttachFirewall(
		plan.FirewallID.ValueString(),
		plan.CloudID.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error attaching security group", fmt.Sprintf("%s", err))
		return
	}

	// ID = firewall_id:cloud_id so it is unique and reversible
	plan.ID = types.StringValue(
		fmt.Sprintf("%s:%s", plan.FirewallID.ValueString(), plan.CloudID.ValueString()),
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// ── READ — nothing to read, attachment has no readable state ──────────────

func (r *CloudFirewallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Firewall attachments have no dedicated GET endpoint
	// State is maintained by Terraform
}

// ── UPDATE — firewall_id and cloud_id both have RequiresReplace ───────────
// so any change triggers destroy + recreate automatically

func (r *CloudFirewallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// RequiresReplace on both fields means Update is never called
	// Terraform will destroy and recreate instead
}

// ── DELETE — detaches firewall from cloud instance ────────────────────────

func (r *CloudFirewallResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CloudFirewallModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DetachFirewall(
		state.FirewallID.ValueString(),
		state.CloudID.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error detaching security group", fmt.Sprintf("%s", err))
		return
	}
}
