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
// utho_ebs
// ══════════════════════════════════════════════════════════════

type EBSResource struct{ client *client.Client }

type EBSModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	DCSlug     types.String `tfsdk:"dcslug"`
	DiskType   types.String `tfsdk:"disk_type"`
	Disk       types.String `tfsdk:"disk"`
	IOPS       types.String `tfsdk:"iops"`
	Throughput types.String `tfsdk:"throughput"`
	Status     types.String `tfsdk:"status"`
	CloudID    types.String `tfsdk:"cloudid"`
	CreatedAt  types.String `tfsdk:"created_at"`
}

func NewEBSResource() resource.Resource { return &EBSResource{} }

func (r *EBSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ebs"
}

func (r *EBSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage Utho Elastic Block Storage (EBS) volumes. Volumes can be attached to cloud instances using the utho_ebs_attachment resource.",
		Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"status":     schema.StringAttribute{Computed: true, Description: "Volume status (Active, Inactive)."},
			"cloudid":    schema.StringAttribute{Computed: true, Description: "Cloud instance ID the volume is attached to."},
			"created_at": schema.StringAttribute{Computed: true, Description: "Creation timestamp.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":       schema.StringAttribute{Required: true, Description: "Volume name. Can be updated in place."},
			"dcslug":     schema.StringAttribute{Required: true, Description: "Data center slug (e.g. innoida). Changing forces new resource.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"disk_type":  schema.StringAttribute{Required: true, Description: "Disk type: SSD or HDD. Changing forces new resource.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"disk":       schema.StringAttribute{Required: true, Description: "Volume size in GB (10–10240). Can only be increased, never decreased."},
			"iops":       schema.StringAttribute{Required: true, Description: "IOPS (1 to 5×disk size). Can be updated in place."},
			"throughput": schema.StringAttribute{Required: true, Description: "Throughput in MB/s (max 500). Can be updated in place."},
		},
	}
}

func (r *EBSResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EBSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EBSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateEBS(&client.EBSCreateRequest{
		DCSlug:     plan.DCSlug.ValueString(),
		DiskType:   plan.DiskType.ValueString(),
		Name:       plan.Name.ValueString(),
		Disk:       plan.Disk.ValueString(),
		IOPS:       plan.IOPS.ValueString(),
		Throughput: plan.Throughput.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating EBS volume", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	plan.Status = types.StringValue("Active")
	plan.CloudID = types.StringValue("")
	plan.CreatedAt = types.StringValue("")

	// Refresh from API
	if ebs, err2 := r.client.GetEBS(id); err2 == nil && ebs != nil {
		plan.Status = types.StringValue(ebs.Status)
		plan.CloudID = types.StringValue(ebs.CloudID)
		plan.CreatedAt = types.StringValue(ebs.CreatedAt)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *EBSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EBSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ebs, err := r.client.GetEBS(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading EBS volume", fmt.Sprintf("%s", err))
		return
	}
	if ebs == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Keep name from state - API may return stale name after update
	state.Status = types.StringValue(ebs.Status)
	state.CloudID = types.StringValue(ebs.CloudID)
	state.CreatedAt = types.StringValue(ebs.CreatedAt)
	state.IOPS = types.StringValue(ebs.IOPS)
	state.Throughput = types.StringValue(ebs.Throughput)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *EBSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state EBSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update name if changed
	if plan.Name.ValueString() != state.Name.ValueString() {
		if err := r.client.UpdateEBSName(state.ID.ValueString(), &client.EBSUpdateNameRequest{
			Name: plan.Name.ValueString(),
		}); err != nil {
			resp.Diagnostics.AddError("Error updating EBS name", fmt.Sprintf("%s", err))
			return
		}
	}

	// Resize if disk/iops/throughput changed
	if plan.Disk.ValueString() != state.Disk.ValueString() ||
		plan.IOPS.ValueString() != state.IOPS.ValueString() ||
		plan.Throughput.ValueString() != state.Throughput.ValueString() {
		if err := r.client.ResizeEBS(state.ID.ValueString(), &client.EBSResizeRequest{
			Disk:       plan.Disk.ValueString(),
			IOPS:       plan.IOPS.ValueString(),
			Throughput: plan.Throughput.ValueString(),
		}); err != nil {
			resp.Diagnostics.AddError("Error resizing EBS volume", fmt.Sprintf("%s", err))
			return
		}
	}

	plan.ID = state.ID
	plan.Status = state.Status
	plan.CloudID = state.CloudID
	plan.CreatedAt = state.CreatedAt
	// Refresh from API to get latest values
	if ebs, err2 := r.client.GetEBS(state.ID.ValueString()); err2 == nil && ebs != nil {
		plan.Status = types.StringValue(ebs.Status)
		plan.CloudID = types.StringValue(ebs.CloudID)
		// Keep name from plan - API may return stale data
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *EBSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EBSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteEBS(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting EBS volume", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// utho_ebs_attachment
// ══════════════════════════════════════════════════════════════

type EBSAttachmentResource struct{ client *client.Client }

type EBSAttachmentModel struct {
	ID      types.String `tfsdk:"id"`
	EBSID   types.String `tfsdk:"ebs_id"`
	CloudID types.String `tfsdk:"cloud_id"`
	Device  types.String `tfsdk:"device"`
}

func NewEBSAttachmentResource() resource.Resource { return &EBSAttachmentResource{} }

func (r *EBSAttachmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ebs_attachment"
}

func (r *EBSAttachmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Attach a Utho EBS volume to a cloud instance.",
		Attributes: map[string]schema.Attribute{
			"id":       schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"device":   schema.StringAttribute{Computed: true, Description: "Device name assigned (e.g. vdb, vdc)."},
			"ebs_id":   schema.StringAttribute{Required: true, Description: "EBS volume ID. Changing forces new resource.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"cloud_id": schema.StringAttribute{Required: true, Description: "Cloud instance ID to attach the volume to. Changing forces new resource.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *EBSAttachmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EBSAttachmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EBSAttachmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.AttachEBSVolume(plan.EBSID.ValueString(), &client.EBSAttachRequest{
		ResourceID: plan.CloudID.ValueString(),
		Type:       "cloud",
	}); err != nil {
		resp.Diagnostics.AddError("Error attaching EBS volume", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", plan.EBSID.ValueString(), plan.CloudID.ValueString()))
	plan.Device = types.StringValue("")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *EBSAttachmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EBSAttachmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ebs, err := r.client.GetEBS(state.EBSID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading EBS attachment", fmt.Sprintf("%s", err))
		return
	}
	if ebs == nil || ebs.CloudID != state.CloudID.ValueString() {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *EBSAttachmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "EBS attachments cannot be updated. Destroy and recreate.")
}

func (r *EBSAttachmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EBSAttachmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DetachEBSVolume(state.EBSID.ValueString(), &client.EBSAttachRequest{
		ResourceID: state.CloudID.ValueString(),
		Type:       "cloud",
	}); err != nil {
		resp.Diagnostics.AddError("Error detaching EBS volume", fmt.Sprintf("%s", err))
	}
}
