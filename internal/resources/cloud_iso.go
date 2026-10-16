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

type CloudISOResource struct {
	client *client.Client
}

type CloudISOModel struct {
	ID      types.String `tfsdk:"id"`
	CloudID types.String `tfsdk:"cloud_id"`
	ISO     types.String `tfsdk:"iso"`
}

func NewCloudISOResource() resource.Resource {
	return &CloudISOResource{}
}

func (r *CloudISOResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_iso"
}

func (r *CloudISOResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Mount an ISO on a Utho Cloud instance and boot from it.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for this ISO mount.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloud_id": schema.StringAttribute{
				Required:    true,
				Description: "Cloud instance ID to mount the ISO on.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"iso": schema.StringAttribute{
				Required:    true,
				Description: "ISO filename to mount (e.g. ATUjo-194454.iso).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *CloudISOResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// CREATE — mounts ISO on cloud instance
func (r *CloudISOResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudISOModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.MountISO(
		plan.CloudID.ValueString(),
		plan.ISO.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error mounting ISO", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(
		fmt.Sprintf("%s:%s", plan.CloudID.ValueString(), plan.ISO.ValueString()),
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// READ — no dedicated GET endpoint
func (r *CloudISOResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// UPDATE — ISO and cloud_id both have RequiresReplace
// any change triggers destroy + recreate
func (r *CloudISOResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// DELETE — ISO mount is a one-time boot action, nothing to undo via API
func (r *CloudISOResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
