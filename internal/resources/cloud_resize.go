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

type CloudResizeResource struct {
	client *client.Client
}

type CloudResizeModel struct {
	ID         types.String `tfsdk:"id"`
	CloudID    types.String `tfsdk:"cloud_id"`
	PlanID     types.String `tfsdk:"plan_id"`
	ResizeType types.String `tfsdk:"resize_type"`
}

func NewCloudResizeResource() resource.Resource {
	return &CloudResizeResource{}
}

func (r *CloudResizeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_resize"
}

func (r *CloudResizeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resize a Utho Cloud instance — CPU/RAM only or full resize including disk.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for this resize operation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloud_id": schema.StringAttribute{
				Required:    true,
				Description: "Cloud instance ID to resize.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"plan_id": schema.StringAttribute{
				Required:    true,
				Description: "Plan ID to resize to (e.g. 10315).",
			},
			"resize_type": schema.StringAttribute{
				Required:    true,
				Description: "Resize type: ramcpu (CPU/RAM only) or full (CPU/RAM + disk).",
			},
		},
	}
}

func (r *CloudResizeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// CREATE — performs the resize
func (r *CloudResizeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudResizeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate resize type
	if plan.ResizeType.ValueString() != "ramcpu" && plan.ResizeType.ValueString() != "full" {
		resp.Diagnostics.AddError(
			"Invalid resize_type",
			"resize_type must be 'ramcpu' (CPU/RAM only) or 'full' (CPU/RAM + disk).",
		)
		return
	}

	err := r.client.ResizeCloud(
		plan.CloudID.ValueString(),
		plan.PlanID.ValueString(),
		plan.ResizeType.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error resizing cloud instance", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(
		fmt.Sprintf("%s:%s", plan.CloudID.ValueString(), plan.PlanID.ValueString()),
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// READ — no dedicated GET endpoint for resize status
func (r *CloudResizeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// UPDATE — re-apply resize if plan or type changes
func (r *CloudResizeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CloudResizeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ResizeType.ValueString() != "ramcpu" && plan.ResizeType.ValueString() != "full" {
		resp.Diagnostics.AddError(
			"Invalid resize_type",
			"resize_type must be 'ramcpu' or 'full'.",
		)
		return
	}

	err := r.client.ResizeCloud(
		plan.CloudID.ValueString(),
		plan.PlanID.ValueString(),
		plan.ResizeType.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error resizing cloud instance", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(
		fmt.Sprintf("%s:%s", plan.CloudID.ValueString(), plan.PlanID.ValueString()),
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// DELETE — resize is not reversible via API, just remove from state
func (r *CloudResizeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
