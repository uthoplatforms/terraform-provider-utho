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
	"github.com/uthoterraform/terraform-provider-utho/api"
)

// implement resource interfaces.
var (
	_ resource.Resource                = &VpcResource{}
	_ resource.ResourceWithConfigure   = &VpcResource{}
	_ resource.ResourceWithImportState = &VpcResource{}
	_ resource.ResourceWithModifyPlan  = &VpcResource{}
)

// NewVpcResource is a helper function to simplify the provider implementation.
func NewVpcResource() resource.Resource {
	return &VpcResource{}
}

// VpcResource is the resource implementation.
type (
	VpcResource struct {
		client *api.Client
	}

	// VpcResource is the model implementation.
	VpcResourceModel struct {
		Id        types.String `tfsdk:"id"`
		Name      types.String `tfsdk:"name"`
		Dcslug    types.String `tfsdk:"dcslug"`
		Planid    types.String `tfsdk:"planid"`
		Network   types.String `tfsdk:"network"`
		Size      types.String `tfsdk:"size"`
		Total     types.Int64  `tfsdk:"total"`
		Available types.Int64  `tfsdk:"available"`
	}
)

// Metadata returns the resource type name.
func (s *VpcResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc"
}

// Configure adds the provider configured client to the data source.
func (d *VpcResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Vpc Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

// Schema defines the schema for the resource.
func (s *VpcResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required:      true,
				Description:   "Provide VPC name eg: vpc1",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"dcslug": schema.StringAttribute{
				Required:      true,
				Description:   "Provide Zone dcslug eg: innoida",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"planid": schema.StringAttribute{
				Required:      true,
				Description:   "Provide network eg: 10.210.100.0",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"network": schema.StringAttribute{
				Required:    true,
				Description: "Provide the planid eg: 1008",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"size": schema.StringAttribute{
				Required:      true,
				Description:   "Provide subnet size eg: 24",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"total":     schema.Int64Attribute{Computed: true, Description: "total"},
			"available": schema.Int64Attribute{Computed: true, Description: "k8s available"},
		},
	}
}

// ModifyPlan tailor the plan to match the expected end state.
func (s *VpcResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Check if the resource is being created.
	if req.State.Raw.IsNull() {
		tflog.Debug(ctx, "start ModifyPlan")

		var state VpcResourceModel
		diags := req.Plan.Get(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		if state.Planid.ValueString() == "" {
			state.Planid = types.StringValue("1008")
		}
	}
}

// Import using vpc as the attribute
func (s *VpcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource.
func (s *VpcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "create vpc")
	// Retrieve values from plan
	var plan VpcResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	vpcRequest := api.VpcRequest{
		Dcslug:  plan.Dcslug.ValueString(),
		Name:    plan.Name.ValueString(),
		Planid:  plan.Planid.ValueString(),
		Network: plan.Network.ValueString(),
		Size:    plan.Size.ValueString(),
	}
	tflog.Debug(ctx, "send create vpc request")
	vpc, err := s.client.CreateVpc(ctx, vpcRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vpc",
			"Could not create vpc, unexpected error: "+err.Error(),
		)
		return
	}

	// get vpc data
	getVpc, err := s.client.GetVpc(ctx, vpc.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading utho vpc",
			"Could not read utho vpc "+vpc.Id+": "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = VpcResourceModel{
		Id:        types.StringValue(vpc.Id),
		Dcslug:    types.StringValue(plan.Dcslug.ValueString()),
		Name:      types.StringValue(plan.Name.ValueString()),
		Planid:    types.StringValue(plan.Planid.ValueString()),
		Network:   types.StringValue(plan.Network.ValueString()),
		Size:      types.StringValue(plan.Size.ValueString()),
		Total:     types.Int64Value(int64(getVpc.Total)),
		Available: types.Int64Value(int64(getVpc.Available)),
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish create vpc")
}

// Read resource information.
func (s *VpcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "read vpc")

	// Get current state
	var state VpcResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "send get vpc request")
	// Get refreshed vpc value from utho
	vpc, err := s.client.GetVpc(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading utho vpc",
			"Could not read utho vpc vpc "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state = VpcResourceModel{
		Id:        types.StringValue(vpc.Id),
		Dcslug:    types.StringValue(vpc.Dcslug),
		Name:      types.StringValue(vpc.Name),
		Planid:    types.StringValue(state.Planid.ValueString()),
		Network:   types.StringValue(vpc.Network),
		Size:      types.StringValue(vpc.Size),
		Total:     types.Int64Value(int64(vpc.Total)),
		Available: types.Int64Value(int64(vpc.Available)),
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish get vpc request")
}

func (s *VpcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// updating resource is not supported
}

// Delete deletes the resource and removes the Terraform state on success.
func (s *VpcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError(
		"Error deleteing utho vpc",
		"Could not delete utho vpc ",
	)
	return
	// tflog.Debug(ctx, "delete vpc")
	// // Get current state
	// var state VpcResourceModel
	// diags := req.State.Get(ctx, &state)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
	// tflog.Debug(ctx, "send delete vpc request")
	// // delete vpc
	// err := s.client.DeleteVpc(ctx, state.Vpc.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error deleteing utho vpc",
	// 		"Could not delete utho vpc "+state.Vpc.ValueString()+": "+err.Error(),
	// 	)
	// 	return
	// }
}
