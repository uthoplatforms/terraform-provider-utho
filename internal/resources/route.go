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

type RouteResource struct{ client *client.Client }
type RouteModel struct {
	ID              types.String `tfsdk:"id"`
	RouteTableID    types.String `tfsdk:"route_table_id"`
	DestinationCIDR types.String `tfsdk:"destination_cidr_block"`
	RouteType       types.String `tfsdk:"route_type"`
	Target          types.String `tfsdk:"target"`
}

func NewRouteResource() resource.Resource { return &RouteResource{} }

func (r *RouteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route"
}

func (r *RouteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage individual routes inside a Utho route table.",
		Attributes: map[string]schema.Attribute{
			"id":                     schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"route_table_id":         schema.StringAttribute{Required: true, Description: "Route table ID to add this route to.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"destination_cidr_block": schema.StringAttribute{Required: true, Description: "Destination CIDR block (e.g. 10.0.0.0/16)."},
			"route_type":             schema.StringAttribute{Required: true, Description: "Route type (e.g. igw)."},
			"target":                 schema.StringAttribute{Required: true, Description: "Route target (e.g. igw)."},
		},
	}
}

func (r *RouteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RouteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RouteModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	routeID, err := r.client.CreateRoute(&client.RouteCreateRequest{
		DestinationCIDR: plan.DestinationCIDR.ValueString(),
		RouteType:       plan.RouteType.ValueString(),
		Target:          plan.Target.ValueString(),
		RouteTableID:    plan.RouteTableID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating route", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(routeID)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *RouteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RouteModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *RouteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RouteModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateRoute(&client.RouteUpdateRequest{
		DestinationCIDR: plan.DestinationCIDR.ValueString(),
		RouteType:       plan.RouteType.ValueString(),
		Target:          plan.Target.ValueString(),
		RouteTableID:    plan.RouteTableID.ValueString(),
		RouteID:         plan.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating route", fmt.Sprintf("%s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *RouteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RouteModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteRoute(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting route", fmt.Sprintf("%s", err))
	}
}
