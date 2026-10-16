package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ══════════════════════════════════════════════════════════════
// VPC RESOURCE — utho_vpc
// ══════════════════════════════════════════════════════════════

type VPCResource struct{ client *client.Client }
type VPCModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Network   types.String `tfsdk:"network"`
	Size      types.String `tfsdk:"size"`
	DCSlug    types.String `tfsdk:"dcslug"`
	IsDefault types.String `tfsdk:"is_default"`
	Total     types.Int64  `tfsdk:"total"`
	Available types.Int64  `tfsdk:"available"`
}

func NewVPCResource() resource.Resource { return &VPCResource{} }

func (r *VPCResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc"
}

func (r *VPCResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage Utho VPC networks.",
		Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":       schema.StringAttribute{Required: true, Description: "VPC name."},
			"network":    schema.StringAttribute{Required: true, Description: "Network address (e.g. 10.0.3.0).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"size":       schema.StringAttribute{Required: true, Description: "CIDR prefix size (e.g. 24).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"dcslug":     schema.StringAttribute{Required: true, Description: "Data center slug.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"is_default": schema.StringAttribute{Computed: true, Description: "Whether this is the default VPC."},
			"total":      schema.Int64Attribute{Computed: true, Description: "Total IPs in the VPC."},
			"available":  schema.Int64Attribute{Computed: true, Description: "Available IPs in the VPC."},
		},
	}
}

func (r *VPCResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VPCResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VPCModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateVPC(&client.VPCCreateRequest{
		Name: plan.Name.ValueString(), Network: plan.Network.ValueString(),
		Size: plan.Size.ValueString(), DCSlug: plan.DCSlug.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating VPC", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	plan.IsDefault = types.StringValue("0")
	plan.Total = types.Int64Value(0)
	plan.Available = types.Int64Value(0)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *VPCResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VPCModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vpc, err := r.client.GetVPC(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading VPC", fmt.Sprintf("%s", err))
		return
	}
	if vpc == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(vpc.Name)
	state.Network = types.StringValue(vpc.Network)
	state.Size = types.StringValue(vpc.Size)
	state.DCSlug = types.StringValue(vpc.DCSlug)
	state.IsDefault = types.StringValue(vpc.IsDefault)
	state.Total = types.Int64Value(int64(vpc.Total))
	state.Available = types.Int64Value(int64(vpc.Available))
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *VPCResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "VPCs cannot be updated in place. Destroy and recreate.")
}

func (r *VPCResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VPCModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteVPC(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting VPC", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// SUBNET RESOURCE — utho_subnet
// ══════════════════════════════════════════════════════════════

type SubnetResource struct{ client *client.Client }
type SubnetModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	VPCID          types.String `tfsdk:"vpc_id"`
	Network        types.String `tfsdk:"network"`
	Size           types.Int64  `tfsdk:"size"`
	Type           types.String `tfsdk:"type"`
	AssignPublicIP types.Int64  `tfsdk:"assign_publicip"`
	IsDefault      types.Int64  `tfsdk:"is_default"`
	Gateway        types.String `tfsdk:"gateway"`
	DCSlug         types.String `tfsdk:"dcslug"`
	Status         types.String `tfsdk:"status"`
}

func NewSubnetResource() resource.Resource { return &SubnetResource{} }

func (r *SubnetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnet"
}

func (r *SubnetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage subnets inside a Utho VPC.",
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":            schema.StringAttribute{Required: true, Description: "Subnet name."},
			"vpc_id":          schema.StringAttribute{Required: true, Description: "VPC ID to create the subnet in.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"network":         schema.StringAttribute{Required: true, Description: "Subnet network address (e.g. 10.135.4.0).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"size":            schema.Int64Attribute{Required: true, Description: "CIDR prefix size (e.g. 22)."},
			"type":            schema.StringAttribute{Required: true, Description: "Subnet type: public or private."},
			"assign_publicip": schema.Int64Attribute{Optional: true, Description: "Assign public IP to instances: 1 or 0."},
			"is_default":      schema.Int64Attribute{Optional: true, Description: "Set as default subnet: 1 or 0."},
			"gateway":         schema.StringAttribute{Computed: true, Description: "Gateway IP assigned to the subnet."},
			"dcslug":          schema.StringAttribute{Computed: true, Description: "Data center slug."},
			"status":          schema.StringAttribute{Computed: true, Description: "Subnet status."},
		},
	}
}

func (r *SubnetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SubnetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SubnetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateSubnet(&client.SubnetCreateRequest{
		Name: plan.Name.ValueString(), VPCID: plan.VPCID.ValueString(),
		Network: plan.Network.ValueString(), Size: int(plan.Size.ValueInt64()),
		Type: plan.Type.ValueString(), AssignPublicIP: int(plan.AssignPublicIP.ValueInt64()),
		IsDefault: int(plan.IsDefault.ValueInt64()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating subnet", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	plan.Gateway = types.StringValue("")
	plan.Status = types.StringValue("Active")
	plan.DCSlug = types.StringValue("")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *SubnetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SubnetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subnet, err := r.client.GetSubnet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading subnet", fmt.Sprintf("%s", err))
		return
	}
	if subnet == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(subnet.Name)
	state.Network = types.StringValue(subnet.Network)
	state.Gateway = types.StringValue(subnet.Gateway)
	state.DCSlug = types.StringValue(subnet.DCSlug)
	state.Status = types.StringValue(subnet.Status)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SubnetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Subnets cannot be updated in place.")
}

func (r *SubnetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SubnetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteSubnet(state.ID.ValueString()); err != nil {
		errMsg := err.Error()
		// If subnet is still attached to resources, remove from state and warn user.
		// The subnet must be manually detached from all resources (LB, cloud instances,
		// NAT gateways, IPSec etc.) before it can be deleted.
		if strings.Contains(errMsg, "attached with resources") || strings.Contains(errMsg, "Please release subnet") {
			resp.Diagnostics.AddWarning(
				"Subnet not deleted — still attached to resources",
				fmt.Sprintf(
					"Subnet %s could not be deleted because it is still attached to one or more resources (load balancers, cloud instances, NAT gateways, etc.). "+
						"Please detach the subnet from all resources in the Utho Console, then run `terraform destroy` again or manually delete the subnet and run `terraform state rm` to remove it from state. "+
						"Original error: %s", state.ID.ValueString(), errMsg,
				),
			)
			return
		}
		resp.Diagnostics.AddError("Error deleting subnet", errMsg)
	}
}
