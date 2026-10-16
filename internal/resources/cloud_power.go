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

type CloudPowerResource struct {
	client *client.Client
}

// ── Model ─────────────────────────────────────────────────────────────────

type CloudPowerModel struct {
	ID      types.String `tfsdk:"id"`
	CloudID types.String `tfsdk:"cloud_id"`
	Action  types.String `tfsdk:"action"`
}

// ── Constructor ───────────────────────────────────────────────────────────

func NewCloudPowerResource() resource.Resource {
	return &CloudPowerResource{}
}

// ── Metadata ──────────────────────────────────────────────────────────────

func (r *CloudPowerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_power"
}

// ── Schema ────────────────────────────────────────────────────────────────

func (r *CloudPowerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage power state of a Utho Cloud instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloud_id": schema.StringAttribute{
				Required:    true,
				Description: "The cloud instance ID to manage power for.",
			},
			"action": schema.StringAttribute{
				Required:    true,
				Description: "Power action: poweron, poweroff, hardreboot, powercycle.",
			},
		},
	}
}

// ── Configure ─────────────────────────────────────────────────────────────

func (r *CloudPowerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// ── CREATE ────────────────────────────────────────────────────────────────

func (r *CloudPowerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudPowerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate action
	validActions := map[string]bool{
		"poweron":    true,
		"poweroff":   true,
		"hardreboot": true,
		"powercycle": true,
	}
	if !validActions[plan.Action.ValueString()] {
		resp.Diagnostics.AddError(
			"Invalid power action",
			fmt.Sprintf("Action must be one of: poweron, poweroff, hardreboot, powercycle. Got: %s", plan.Action.ValueString()),
		)
		return
	}

	err := r.client.PowerAction(plan.CloudID.ValueString(), plan.Action.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error performing power action", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(plan.CloudID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// ── READ ──────────────────────────────────────────────────────────────────

func (r *CloudPowerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Power actions are one-time operations — nothing to read
}

// ── UPDATE ────────────────────────────────────────────────────────────────

func (r *CloudPowerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CloudPowerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.PowerAction(plan.CloudID.ValueString(), plan.Action.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error performing power action", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(plan.CloudID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// ── DELETE ────────────────────────────────────────────────────────────────

func (r *CloudPowerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Nothing to delete — power action is already done
}
