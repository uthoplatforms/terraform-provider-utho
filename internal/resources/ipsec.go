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
// IPSEC TUNNEL — utho_ipsec
// ══════════════════════════════════════════════════════════════

type IPSecResource struct{ client *client.Client }

type IPSecModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	DCSlug       types.String `tfsdk:"dcslug"`
	VPC          types.String `tfsdk:"vpc"`
	BillingCycle types.String `tfsdk:"billingcycle"`
	PSK          types.String `tfsdk:"psk"`
	Status       types.String `tfsdk:"status"`
	CreatedAt    types.String `tfsdk:"created_at"`
	CloudID      types.String `tfsdk:"cloudid"`
}

func NewIPSecResource() resource.Resource { return &IPSecResource{} }

func (r *IPSecResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipsec"
}

func (r *IPSecResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage Utho IPSec site-to-site VPN tunnels.",
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"psk":          schema.StringAttribute{Computed: true, Sensitive: true, Description: "Auto-generated pre-shared key for the IPSec tunnel.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"status":       schema.StringAttribute{Computed: true, Description: "Tunnel status (active, pending)."},
			"created_at":   schema.StringAttribute{Computed: true, Description: "Creation timestamp.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"cloudid":      schema.StringAttribute{Computed: true, Description: "Cloud instance ID of the IPSec gateway.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":         schema.StringAttribute{Required: true, Description: "Tunnel name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"dcslug":       schema.StringAttribute{Required: true, Description: "Data center slug (e.g. inmumbaizone2).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"vpc":          schema.StringAttribute{Required: true, Description: "VPC ID or subnet ID. If the VPC has subnets, provide the subnet ID.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"billingcycle": schema.StringAttribute{Required: true, Description: "Billing cycle: monthly, 3month, 6month, 12month.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *IPSecResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IPSecResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan IPSecModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateIPSec(&client.IPSecCreateRequest{
		Name:         plan.Name.ValueString(),
		DCSlug:       plan.DCSlug.ValueString(),
		VPC:          plan.VPC.ValueString(),
		CPUModel:     "amd",
		BillingCycle: plan.BillingCycle.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating IPSec tunnel", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	plan.PSK = types.StringValue("")
	plan.Status = types.StringValue("pending")
	plan.CreatedAt = types.StringValue("")
	plan.CloudID = types.StringValue("")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Poll until active (up to 10 minutes)
	if waitErr := r.client.WaitForIPSecReady(id); waitErr != nil {
		if waitErr.Error() == "STILL_PROVISIONING" {
			resp.Diagnostics.AddWarning(
				"IPSec tunnel still provisioning",
				fmt.Sprintf("Tunnel %q (ID: %s) is still deploying. Run terraform apply again once active.", plan.Name.ValueString(), id),
			)
		}
	}

	// Refresh from API
	if ipsec, err2 := r.client.GetIPSec(id); err2 == nil && ipsec != nil {
		plan.PSK = types.StringValue(ipsec.PSK)
		plan.Status = types.StringValue(ipsec.Status)
		plan.CreatedAt = types.StringValue(ipsec.CreatedAt)
		plan.CloudID = types.StringValue(ipsec.CloudID)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *IPSecResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IPSecModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ipsec, err := r.client.GetIPSec(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading IPSec tunnel", fmt.Sprintf("%s", err))
		return
	}
	if ipsec == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(ipsec.Name)
	state.DCSlug = types.StringValue(ipsec.DCSlug)
	state.PSK = types.StringValue(ipsec.PSK)
	state.Status = types.StringValue(ipsec.Status)
	state.BillingCycle = types.StringValue(ipsec.BillingCycle)
	state.CreatedAt = types.StringValue(ipsec.CreatedAt)
	state.CloudID = types.StringValue(ipsec.CloudID)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IPSecResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "IPSec tunnels cannot be updated in place. All fields require replacement.")
}

func (r *IPSecResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state IPSecModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteIPSec(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting IPSec tunnel", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// IPSEC CONNECTION — utho_ipsec_connection
// Pairs two Utho IPSec tunnels together
// ══════════════════════════════════════════════════════════════

type IPSecConnectionResource struct{ client *client.Client }

type IPSecConnectionModel struct {
	ID          types.String `tfsdk:"id"`
	IPSecID     types.String `tfsdk:"ipsec_id"`
	PeerIPSecID types.String `tfsdk:"peer_ipsec_id"`
	LocalSubnet types.String `tfsdk:"local_subnet"`
	PeerSubnet  types.String `tfsdk:"peer_subnet"`
	Name        types.String `tfsdk:"name"`
	Status      types.String `tfsdk:"status"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

func NewIPSecConnectionResource() resource.Resource { return &IPSecConnectionResource{} }

func (r *IPSecConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipsec_connection"
}

func (r *IPSecConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Pair two Utho IPSec tunnels together to create a site-to-site VPN connection.",
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"status":        schema.StringAttribute{Computed: true, Description: "Connection status."},
			"created_at":    schema.StringAttribute{Computed: true, Description: "Creation timestamp.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"ipsec_id":      schema.StringAttribute{Required: true, Description: "Primary IPSec tunnel ID.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"peer_ipsec_id": schema.StringAttribute{Required: true, Description: "Peer IPSec tunnel ID to connect with.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"local_subnet":  schema.StringAttribute{Required: true, Description: "Local network CIDR (e.g. 192.168.50.0/24).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"peer_subnet":   schema.StringAttribute{Required: true, Description: "Peer network CIDR (e.g. 192.168.60.0/24).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"name":          schema.StringAttribute{Required: true, Description: "Connection name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *IPSecConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IPSecConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan IPSecConnectionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateIPSecPair(&client.IPSecPairRequest{
		IPSecID:     plan.IPSecID.ValueString(),
		PeerIPSecID: plan.PeerIPSecID.ValueString(),
		LocalSubnet: plan.LocalSubnet.ValueString(),
		PeerSubnet:  plan.PeerSubnet.ValueString(),
		Name:        plan.Name.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating IPSec connection", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	plan.Status = types.StringValue("active")
	plan.CreatedAt = types.StringValue("")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *IPSecConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IPSecConnectionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connections, err := r.client.ListIPSecConnections(state.IPSecID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading IPSec connections", fmt.Sprintf("%s", err))
		return
	}

	for _, conn := range connections {
		if conn.ID == state.ID.ValueString() {
			state.Status = types.StringValue(conn.Status)
			state.CreatedAt = types.StringValue(conn.CreatedAt)
			resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
			return
		}
	}
	resp.State.RemoveResource(ctx)
}

func (r *IPSecConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "IPSec connections cannot be updated. Destroy and recreate.")
}

func (r *IPSecConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state IPSecConnectionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteIPSecConnection(state.IPSecID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting IPSec connection", fmt.Sprintf("%s", err))
	}
}
