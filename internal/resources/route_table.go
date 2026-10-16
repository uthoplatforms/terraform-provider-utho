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

type RouteTableResource struct{ client *client.Client }
type RouteTableModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	VPCID  types.String `tfsdk:"vpc_id"`
	DCSlug types.String `tfsdk:"dcslug"`
}

func NewRouteTableResource() resource.Resource { return &RouteTableResource{} }

func (r *RouteTableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route_table"
}

func (r *RouteTableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage route tables for Utho VPCs.",
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":   schema.StringAttribute{Required: true, Description: "Route table name."},
			"vpc_id": schema.StringAttribute{Required: true, Description: "VPC ID to associate the route table with.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"dcslug": schema.StringAttribute{Required: true, Description: "Data center slug.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *RouteTableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RouteTableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RouteTableModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateRouteTable(&client.RouteTableCreateRequest{
		Name: plan.Name.ValueString(), VPC: plan.VPCID.ValueString(), DCSlug: plan.DCSlug.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating route table", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *RouteTableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RouteTableModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tables, err := r.client.GetRouteTable(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading route table", fmt.Sprintf("%s", err))
		return
	}
	if len(tables) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *RouteTableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Route tables cannot be updated in place.")
}

func (r *RouteTableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RouteTableModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Detach from VPC first then delete
	if err := r.client.DetachRouteTable(state.VPCID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error detaching route table", fmt.Sprintf("%s", err))
		return
	}
	if err := r.client.DeleteRouteTable(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting route table", fmt.Sprintf("%s", err))
	}
}
