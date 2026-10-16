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

type ElasticIPResource struct{ client *client.Client }

type ElasticIPModel struct {
	ID           types.String `tfsdk:"id"`
	IP           types.String `tfsdk:"ip"`
	DCSlug       types.String `tfsdk:"dcslug"`
	BillingCycle types.String `tfsdk:"billingcycle"`
	CloudID      types.String `tfsdk:"cloud_id"`
	AssignedAt   types.String `tfsdk:"assigned_at"`
}

func NewElasticIPResource() resource.Resource { return &ElasticIPResource{} }

func (r *ElasticIPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_elastic_ip"
}

func (r *ElasticIPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Allocate and manage Elastic IPs on Utho Cloud. Elastic IPs are static public IPs that can be attached to and detached from cloud instances.",
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"ip":           schema.StringAttribute{Computed: true, Description: "Allocated elastic IP address.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"dcslug":       schema.StringAttribute{Required: true, Description: "Data center slug.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"billingcycle": schema.StringAttribute{Required: true, Description: "Billing cycle: monthly or hourly.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"cloud_id":     schema.StringAttribute{Optional: true, Computed: true, Description: "Cloud instance ID to attach this IP to. Leave empty to keep unattached."},
			"assigned_at":  schema.StringAttribute{Computed: true, Description: "Timestamp when the IP was assigned.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		},
	}
}

func (r *ElasticIPResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ElasticIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ElasticIPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	eipAlloc, err := r.client.AllocateElasticIP(&client.ElasticIPAllocateRequest{
		DCSlug:       plan.DCSlug.ValueString(),
		BillingCycle: plan.BillingCycle.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error allocating Elastic IP", fmt.Sprintf("%s", err))
		return
	}

	plan.IP = types.StringValue(eipAlloc.IP)
	plan.ID = types.StringValue(eipAlloc.IP)
	plan.AssignedAt = types.StringValue("")

	// Attach to cloud instance if specified
	if !plan.CloudID.IsNull() && !plan.CloudID.IsUnknown() && plan.CloudID.ValueString() != "" {
		if err := r.client.AttachElasticIP(plan.CloudID.ValueString(), eipAlloc.IP); err != nil {
			resp.Diagnostics.AddError("Error attaching Elastic IP", fmt.Sprintf("%s", err))
			return
		}
	} else {
		plan.CloudID = types.StringValue("")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ElasticIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ElasticIPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	eip, err := r.client.GetElasticIP(state.IP.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading Elastic IP", fmt.Sprintf("%s", err))
		return
	}
	if eip == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.CloudID = types.StringValue(eip.CloudID)
	state.AssignedAt = types.StringValue(eip.AssignedAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ElasticIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ElasticIPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle cloud_id change — detach from old, attach to new
	oldCloudID := state.CloudID.ValueString()
	newCloudID := plan.CloudID.ValueString()

	if oldCloudID != newCloudID {
		// Detach from old instance if attached
		if oldCloudID != "" && oldCloudID != "0" {
			if err := r.client.DetachElasticIP(oldCloudID, state.IP.ValueString()); err != nil {
				resp.Diagnostics.AddError("Error detaching Elastic IP", fmt.Sprintf("%s", err))
				return
			}
		}
		// Attach to new instance if specified
		if newCloudID != "" {
			if err := r.client.AttachElasticIP(newCloudID, state.IP.ValueString()); err != nil {
				resp.Diagnostics.AddError("Error attaching Elastic IP", fmt.Sprintf("%s", err))
				return
			}
		}
	}

	plan.IP = state.IP
	plan.ID = state.ID
	plan.AssignedAt = state.AssignedAt
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ElasticIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ElasticIPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Detach from cloud instance first if attached
	if state.CloudID.ValueString() != "" && state.CloudID.ValueString() != "0" {
		if err := r.client.DetachElasticIP(state.CloudID.ValueString(), state.IP.ValueString()); err != nil {
			resp.Diagnostics.AddError("Error detaching Elastic IP before release", fmt.Sprintf("%s", err))
			return
		}
	}

	if err := r.client.ReleaseElasticIP(state.IP.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error releasing Elastic IP", fmt.Sprintf("%s", err))
	}
}
