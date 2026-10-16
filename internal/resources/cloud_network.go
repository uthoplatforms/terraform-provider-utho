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

// ══════════════════════════════════════════════════════════════
// PUBLIC IP RESOURCE — utho_cloud_public_ip
// ══════════════════════════════════════════════════════════════

type CloudPublicIPResource struct {
	client *client.Client
}

type CloudPublicIPModel struct {
	ID      types.String `tfsdk:"id"`
	CloudID types.String `tfsdk:"cloud_id"`
	IP      types.String `tfsdk:"ip"`
}

func NewCloudPublicIPResource() resource.Resource {
	return &CloudPublicIPResource{}
}

func (r *CloudPublicIPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_public_ip"
}

func (r *CloudPublicIPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Assign or release an additional public IP on a Utho Cloud instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for this IP assignment.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloud_id": schema.StringAttribute{
				Required:    true,
				Description: "Cloud instance ID to assign the IP to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ip": schema.StringAttribute{
				Computed:    true,
				Description: "Public IP address assigned by Utho.",
			},
		},
	}
}

func (r *CloudPublicIPResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudPublicIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudPublicIPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ip, err := r.client.AssignPublicIP(plan.CloudID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error assigning public IP", fmt.Sprintf("%s", err))
		return
	}

	plan.IP = types.StringValue(ip)
	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", plan.CloudID.ValueString(), ip))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *CloudPublicIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *CloudPublicIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *CloudPublicIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CloudPublicIPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePublicIP(
		state.CloudID.ValueString(),
		state.IP.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error releasing public IP", fmt.Sprintf("%s", err))
		return
	}
}

// ══════════════════════════════════════════════════════════════
// VPC ATTACHMENT RESOURCE — utho_cloud_vpc
// ══════════════════════════════════════════════════════════════

type CloudVPCResource struct {
	client *client.Client
}

type CloudVPCModel struct {
	ID        types.String `tfsdk:"id"`
	CloudID   types.String `tfsdk:"cloud_id"`
	SubnetID  types.String `tfsdk:"subnet_id"`
	PrivateIP types.String `tfsdk:"private_ip"`
}

func NewCloudVPCResource() resource.Resource {
	return &CloudVPCResource{}
}

func (r *CloudVPCResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_vpc"
}

func (r *CloudVPCResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Attach or detach a VPC subnet from a Utho Cloud instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for this VPC attachment.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloud_id": schema.StringAttribute{
				Required:    true,
				Description: "Cloud instance ID to attach the VPC to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "Subnet UUID to attach (e.g. b0825dae-fd86-4d67-8238-c438774da894).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"private_ip": schema.StringAttribute{
				Computed:    true,
				Description: "Private IP assigned from the subnet after attachment.",
			},
		},
	}
}

func (r *CloudVPCResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudVPCResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudVPCModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	privateIP, err := r.client.AttachVPC(
		plan.CloudID.ValueString(),
		plan.SubnetID.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error attaching VPC", fmt.Sprintf("%s", err))
		return
	}

	plan.PrivateIP = types.StringValue(privateIP)
	plan.ID = types.StringValue(
		fmt.Sprintf("%s:%s", plan.CloudID.ValueString(), plan.SubnetID.ValueString()),
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *CloudVPCResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *CloudVPCResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *CloudVPCResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CloudVPCModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DetachVPC(
		state.CloudID.ValueString(),
		state.SubnetID.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error detaching VPC", fmt.Sprintf("%s", err))
		return
	}
}
