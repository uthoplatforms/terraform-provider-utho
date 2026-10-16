package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudSnapshotResource struct {
	client *client.Client
}

type CloudSnapshotModel struct {
	ID         types.String `tfsdk:"id"`
	CloudID    types.String `tfsdk:"cloud_id"`
	Name       types.String `tfsdk:"name"`
	SnapshotID types.String `tfsdk:"snapshot_id"`
	Restore    types.Bool   `tfsdk:"restore"`
}

func NewCloudSnapshotResource() resource.Resource {
	return &CloudSnapshotResource{}
}

func (r *CloudSnapshotResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_snapshot"
}

func (r *CloudSnapshotResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create, restore, and delete snapshots of a Utho Cloud instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for this snapshot resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloud_id": schema.StringAttribute{
				Required:    true,
				Description: "Cloud instance ID to snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"snapshot_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Snapshot ID. Set this after creation to enable restore/delete operations.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"restore": schema.BoolAttribute{
				Optional:    true,
				Description: "Set to true to restore the snapshot. Requires snapshot_id to be set.",
			},
		},
	}
}

func (r *CloudSnapshotResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// CREATE — takes a snapshot
func (r *CloudSnapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudSnapshotModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateSnapshot(
		plan.CloudID.ValueString(),
		plan.Name.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error creating snapshot", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(
		fmt.Sprintf("%s:%s", plan.CloudID.ValueString(), plan.Name.ValueString()),
	)

	// Get the cloud instance dcslug for targeted snapshot lookup
	dcslug := ""
	if cloud, err2 := r.client.GetCloud(plan.CloudID.ValueString()); err2 == nil && cloud != nil {
		dcslug = cloud.DCSlug
	}

	// Poll up to 3 minutes for snapshot to appear in deploy data
	snapshotID := ""
	for attempt := 0; attempt < 18; attempt++ {
		if attempt > 0 {
			time.Sleep(10 * time.Second)
		}
		deployData, err2 := r.client.GetCloudDeployData(dcslug)
		if err2 != nil {
			continue
		}
		for _, s := range deployData.Snapshots {
			if s.Name == plan.Name.ValueString() && s.CloudID == plan.CloudID.ValueString() {
				snapshotID = s.ID
				break
			}
		}
		if snapshotID != "" {
			break
		}
	}

	if snapshotID == "" {
		resp.Diagnostics.AddWarning(
			"Snapshot ID not yet available",
			"The snapshot was created but its ID could not be retrieved automatically. Run terraform refresh or set snapshot_id manually from the Utho console.",
		)
	}
	plan.SnapshotID = types.StringValue(snapshotID)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// READ — look up snapshot_id from deploy data if not already set
func (r *CloudSnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CloudSnapshotModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If snapshot_id already set, keep it
	if state.SnapshotID.ValueString() != "" {
		resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
		return
	}

	// Try to find snapshot_id from deploy data
	dcslug := ""
	if cloud, err := r.client.GetCloud(state.CloudID.ValueString()); err == nil && cloud != nil {
		dcslug = cloud.DCSlug
	}
	deployData, err := r.client.GetCloudDeployData(dcslug)
	if err == nil {
		for _, s := range deployData.Snapshots {
			if s.Name == state.Name.ValueString() && s.CloudID == state.CloudID.ValueString() {
				state.SnapshotID = types.StringValue(s.ID)
				break
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// UPDATE — restore snapshot if restore = true
func (r *CloudSnapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CloudSnapshotModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If restore = true, restore the snapshot
	if !plan.Restore.IsNull() && plan.Restore.ValueBool() {
		if plan.SnapshotID.ValueString() == "" {
			resp.Diagnostics.AddError(
				"Missing snapshot_id",
				"snapshot_id is required to restore. Get it from Utho dashboard.",
			)
			return
		}

		err := r.client.RestoreSnapshot(
			plan.CloudID.ValueString(),
			plan.SnapshotID.ValueString(),
		)
		if err != nil {
			resp.Diagnostics.AddError("Error restoring snapshot", fmt.Sprintf("%s", err))
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// DELETE — deletes the snapshot
func (r *CloudSnapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CloudSnapshotModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.SnapshotID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing snapshot_id",
			"snapshot_id is required to delete. Get it from Utho dashboard.",
		)
		return
	}

	err := r.client.DeleteSnapshot(
		state.CloudID.ValueString(),
		state.SnapshotID.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting snapshot", fmt.Sprintf("%s", err))
		return
	}
}
