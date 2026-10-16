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

type VPCPeeringResource struct{ client *client.Client }
type VPCPeeringModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	DCSlug         types.String `tfsdk:"dcslug"`
	RequesterVPCID types.String `tfsdk:"requester_vpc_id"`
	AccepterVPCID  types.String `tfsdk:"accepter_vpc_id"`
	Status         types.String `tfsdk:"status"`
}

func NewVPCPeeringResource() resource.Resource { return &VPCPeeringResource{} }

func (r *VPCPeeringResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_peering"
}

func (r *VPCPeeringResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage VPC peering connections between two Utho VPCs.",
		Attributes: map[string]schema.Attribute{
			"id":               schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":             schema.StringAttribute{Required: true, Description: "Peering connection name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"dcslug":           schema.StringAttribute{Required: true, Description: "Data center slug.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"requester_vpc_id": schema.StringAttribute{Required: true, Description: "ID of the requester VPC.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"accepter_vpc_id":  schema.StringAttribute{Required: true, Description: "ID of the accepter VPC.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"status":           schema.StringAttribute{Computed: true, Description: "Peering connection status."},
		},
	}
}

func (r *VPCPeeringResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VPCPeeringResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VPCPeeringModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateVPCPeering(&client.VPCPeeringCreateRequest{
		Name: plan.Name.ValueString(), DCSlug: plan.DCSlug.ValueString(),
		RequesterVPCID: plan.RequesterVPCID.ValueString(), AccepterVPCID: plan.AccepterVPCID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating VPC peering", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	plan.Status = types.StringValue("Active")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *VPCPeeringResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VPCPeeringModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	peerings, err := r.client.ListVPCPeering()
	if err != nil {
		resp.Diagnostics.AddError("Error reading VPC peering", fmt.Sprintf("%s", err))
		return
	}
	for _, p := range peerings {
		if p.ID == state.ID.ValueString() {
			state.Status = types.StringValue(p.Status)
			state.Name = types.StringValue(p.Name)
			resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
			return
		}
	}
	resp.State.RemoveResource(ctx)
}

func (r *VPCPeeringResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "VPC peering connections cannot be updated. Destroy and recreate.")
}

func (r *VPCPeeringResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VPCPeeringModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteVPCPeering(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting VPC peering", fmt.Sprintf("%s", err))
	}
}
