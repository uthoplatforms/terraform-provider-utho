package resources

import (
	"context"
	"fmt"

	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ══════════════════════════════════════════════════════════════
// GENERAL STORAGE RESOURCE — utho_cloud_storage
// ══════════════════════════════════════════════════════════════

type CloudStorageResource struct {
	client *client.Client
}

type CloudStorageModel struct {
	ID       types.String `tfsdk:"id"`
	CloudID  types.String `tfsdk:"cloud_id"`
	DiskID   types.String `tfsdk:"disk_id"`
	SizeGB   types.Int64  `tfsdk:"size_gb"`
	Bus      types.String `tfsdk:"bus"`
	DiskType types.String `tfsdk:"disk_type"`
}

func NewCloudStorageResource() resource.Resource {
	return &CloudStorageResource{}
}

func (r *CloudStorageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage"
}

func (r *CloudStorageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Add and manage general storage disks on a Utho Cloud instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for this storage resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloud_id": schema.StringAttribute{
				Required:    true,
				Description: "Cloud instance ID to attach the storage to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"disk_id": schema.StringAttribute{
				Computed:    true,
				Description: "Disk ID assigned by Utho after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"size_gb": schema.Int64Attribute{
				Required:    true,
				Description: "Disk size in GB.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"bus": schema.StringAttribute{
				Optional:    true,
				Description: "Bus type: virtio or ide. Default: virtio.",
			},
			"disk_type": schema.StringAttribute{
				Optional:    true,
				Description: "Disk type: Primary or Additional. Only one Primary allowed per instance.",
			},
		},
	}
}

func (r *CloudStorageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudStorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudStorageModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.AddGeneralStorage(
		plan.CloudID.ValueString(),
		int(plan.SizeGB.ValueInt64()),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error adding storage", fmt.Sprintf("%s", err))
		return
	}

	// Look up disk ID from cloud instance after adding storage
	diskID := ""
	cloud, err2 := r.client.GetCloud(plan.CloudID.ValueString())
	if err2 == nil && cloud != nil {
		// Find the newly added disk — look for one matching our size
		targetSize := float64(plan.SizeGB.ValueInt64())
		for _, s := range cloud.Storages {
			if s.Size == targetSize {
				diskID = s.ID
				break
			}
		}
		// If not found by size, take the last one
		if diskID == "" && len(cloud.Storages) > 0 {
			diskID = cloud.Storages[len(cloud.Storages)-1].ID
		}
	}

	plan.DiskID = types.StringValue(diskID)
	plan.ID = types.StringValue(
		fmt.Sprintf("%s:%d", plan.CloudID.ValueString(), plan.SizeGB.ValueInt64()),
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *CloudStorageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// No dedicated GET endpoint for individual disk
	// State maintained by Terraform
}

func (r *CloudStorageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CloudStorageModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.DiskID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing disk_id",
			"disk_id is required to update storage. Get it from Utho dashboard.",
		)
		return
	}

	err := r.client.UpdateGeneralStorage(
		plan.CloudID.ValueString(),
		plan.DiskID.ValueString(),
		plan.Bus.ValueString(),
		plan.DiskType.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error updating storage", fmt.Sprintf("%s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *CloudStorageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CloudStorageModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.DiskID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing disk_id",
			"disk_id is required to delete storage. Get it from Utho dashboard.",
		)
		return
	}

	err := r.client.DeleteGeneralStorage(
		state.CloudID.ValueString(),
		state.DiskID.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting storage", fmt.Sprintf("%s", err))
		return
	}
}

// ══════════════════════════════════════════════════════════════
// EBS RESOURCE — utho_cloud_ebs
// ══════════════════════════════════════════════════════════════

type CloudEBSResource struct {
	client *client.Client
}

type CloudEBSModel struct {
	ID       types.String `tfsdk:"id"`
	CloudID  types.String `tfsdk:"cloud_id"`
	EBSID    types.String `tfsdk:"ebs_id"`
	Bus      types.String `tfsdk:"bus"`
	DiskType types.String `tfsdk:"disk_type"`
}

func NewCloudEBSResource() resource.Resource {
	return &CloudEBSResource{}
}

func (r *CloudEBSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_ebs"
}

func (r *CloudEBSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Attach and manage EBS volumes on a Utho Cloud instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for this EBS attachment.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloud_id": schema.StringAttribute{
				Required:    true,
				Description: "Cloud instance ID to attach the EBS volume to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ebs_id": schema.StringAttribute{
				Required:    true,
				Description: "EBS volume ID to attach.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"bus": schema.StringAttribute{
				Optional:    true,
				Description: "Bus type: virtio or ide.",
			},
			"disk_type": schema.StringAttribute{
				Optional:    true,
				Description: "Disk type: Primary or Additional.",
			},
		},
	}
}

func (r *CloudEBSResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudEBSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudEBSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.AttachEBS(
		plan.EBSID.ValueString(),
		plan.CloudID.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error attaching EBS volume", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(
		fmt.Sprintf("%s:%s", plan.EBSID.ValueString(), plan.CloudID.ValueString()),
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *CloudEBSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// No dedicated GET endpoint for EBS attachment
	// State maintained by Terraform
}

func (r *CloudEBSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CloudEBSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateEBS(
		plan.CloudID.ValueString(),
		plan.EBSID.ValueString(),
		plan.Bus.ValueString(),
		plan.DiskType.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error updating EBS volume", fmt.Sprintf("%s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *CloudEBSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CloudEBSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DetachEBS(
		state.EBSID.ValueString(),
		state.CloudID.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error detaching EBS volume", fmt.Sprintf("%s", err))
		return
	}
}
