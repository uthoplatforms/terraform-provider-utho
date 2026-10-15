package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
)

// Ensure CloudResource implements ResourceWithImportState
var _ resource.ResourceWithImportState = &CloudResource{}

type CloudResource struct {
	client *client.Client
}

type EBSVolumeModel struct {
	ID   types.String `tfsdk:"id"`
	Disk types.Int64  `tfsdk:"disk"`
	Type types.String `tfsdk:"type"`
}

type CloudResourceModel struct {
	// Computed
	ID          types.String `tfsdk:"id"`
	IP          types.String `tfsdk:"ip"`
	Status      types.String `tfsdk:"status"`
	PowerStatus types.String `tfsdk:"power_status"`
	CreatedAt   types.String `tfsdk:"created_at"`

	// Required
	Hostname     types.String `tfsdk:"hostname"`
	DCSlug       types.String `tfsdk:"dcslug"`
	PlanID       types.String `tfsdk:"planid"`
	BillingCycle types.String `tfsdk:"billingcycle"`
	Auth         types.String `tfsdk:"auth"`

	// Optional
	EnablePublicIP types.String     `tfsdk:"enable_publicip"`
	SubnetRequired types.String     `tfsdk:"subnet_required"`
	CPUModel       types.String     `tfsdk:"cpumodel"`
	EnableBackup   types.String     `tfsdk:"enablebackup"`
	Support        types.String     `tfsdk:"support"`
	Firewall       types.String     `tfsdk:"firewall"`
	RootPassword   types.String     `tfsdk:"root_password"`
	SSHKeys        types.String     `tfsdk:"sshkeys"`
	Image          types.String     `tfsdk:"image"`
	VPC            types.String     `tfsdk:"vpc"`
	SnapshotID     types.String     `tfsdk:"snapshotid"`
	BackupID       types.String     `tfsdk:"backupid"`
	ISO            types.String     `tfsdk:"iso"`
	Stack          types.String     `tfsdk:"stack"`
	DeleteEBS      types.Bool       `tfsdk:"delete_ebs"`
	EBS            []EBSVolumeModel `tfsdk:"ebs"`
}

func NewCloudResource() resource.Resource {
	return &CloudResource{}
}

func (r *CloudResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud"
}

func (r *CloudResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage Utho Cloud instances.",
		Attributes: map[string]schema.Attribute{

			// Computed
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique cloud instance ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ip": schema.StringAttribute{
				Computed:    true,
				Description: "Public IP address.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Instance status (e.g. Active).",
			},
			"power_status": schema.StringAttribute{
				Computed:    true,
				Description: "Power status (e.g. Running).",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Creation timestamp.",
			},

			// Required
			"hostname": schema.StringAttribute{
				Required:    true,
				Description: "Hostname for the instance. Must be unique.",
			},
			"dcslug": schema.StringAttribute{
				Required:    true,
				Description: "Data center slug (e.g. innoida, inmumbaizone2). Use the utho_cloud_dczones data source to list available zones and their supported CPU models.",
			},
			"planid": schema.StringAttribute{
				Required:    true,
				Description: "Plan ID for the instance size. Use the utho_cloud_plans data source to list plans available in your chosen dcslug. Not all plans are available in every zone.",
			},
			"billingcycle": schema.StringAttribute{
				Required:    true,
				Description: "Billing cycle: hourly, monthly, 12month.",
			},
			"auth": schema.StringAttribute{
				Required:    true,
				Description: "Auth method: option1 (password) or option2 (SSH key).",
			},

			// Optional
			"enable_publicip": schema.StringAttribute{
				Optional:    true,
				Description: "Enable public IP: true or false.",
			},
			"subnet_required": schema.StringAttribute{
				Optional:    true,
				Description: "Whether subnet is required: true or false.",
			},
			"cpumodel": schema.StringAttribute{
				Required:    true,
				Description: "CPU model: amd or intel. Availability varies by data center — check the utho_cloud_dczones data source or the Utho console before setting this value.",
			},
			"enablebackup": schema.StringAttribute{
				Optional:    true,
				Description: "Enable backup: true or false.",
			},
			"support": schema.StringAttribute{
				Optional:    true,
				Description: "Support type: unmanaged or managed.",
			},
			"firewall": schema.StringAttribute{
				Optional:    true,
				Description: "Firewall/security group ID to attach.",
			},
			"root_password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Root password. Required when auth = option1.",
			},
			"sshkeys": schema.StringAttribute{
				Optional:    true,
				Description: "SSH key ID. Required when auth = option2.",
			},
			"image": schema.StringAttribute{
				Optional:    true,
				Description: "OS image slug (e.g. ubuntu-22.04-x86_64).",
			},
			"vpc": schema.StringAttribute{
				Optional:    true,
				Description: "VPC ID for deploying inside a VPC.",
			},
			"snapshotid": schema.StringAttribute{
				Optional:    true,
				Description: "Snapshot ID to deploy from snapshot.",
			},
			"backupid": schema.StringAttribute{
				Optional:    true,
				Description: "Backup ID to restore from backup.",
			},
			"iso": schema.StringAttribute{
				Optional:    true,
				Description: "ISO ID to deploy from custom ISO.",
			},
			"stack": schema.StringAttribute{
				Optional:    true,
				Description: "Stack ID to deploy from marketplace stack.",
			},
			"delete_ebs": schema.BoolAttribute{
				Optional:    true,
				Description: "Delete EBS volumes on destroy. Default: false.",
			},

			// EBS volumes
			"ebs": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Additional EBS volumes to attach.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:    true,
							Description: "EBS volume index ID (e.g. 1, 2).",
						},
						"disk": schema.Int64Attribute{
							Required:    true,
							Description: "Disk size in GB.",
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "Disk type: nvme or ssd.",
						},
					},
				},
			},
		},
	}
}

func (r *CloudResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build EBS volumes
	var ebsVolumes []client.EBSVolume
	for _, e := range plan.EBS {
		ebsVolumes = append(ebsVolumes, client.EBSVolume{
			ID:   e.ID.ValueString(),
			Disk: int(e.Disk.ValueInt64()),
			Type: e.Type.ValueString(),
		})
	}

	deployReq := &client.CloudDeployRequest{
		DCSlug:         plan.DCSlug.ValueString(),
		PlanID:         plan.PlanID.ValueString(),
		BillingCycle:   plan.BillingCycle.ValueString(),
		Auth:           plan.Auth.ValueString(),
		EnablePublicIP: plan.EnablePublicIP.ValueString(),
		SubnetRequired: plan.SubnetRequired.ValueString(),
		CPUModel:       plan.CPUModel.ValueString(),
		EnableBackup:   plan.EnableBackup.ValueString(),
		Support:        plan.Support.ValueString(),
		Firewall:       plan.Firewall.ValueString(),
		RootPassword:   plan.RootPassword.ValueString(),
		SSHKeys:        plan.SSHKeys.ValueString(),
		Image:          plan.Image.ValueString(),
		VPC:            plan.VPC.ValueString(),
		SnapshotID:     plan.SnapshotID.ValueString(),
		BackupID:       plan.BackupID.ValueString(),
		ISO:            plan.ISO.ValueString(),
		Stack:          plan.Stack.ValueString(),
		EBS:            ebsVolumes,
		Cloud: []client.CloudServer{
			{Hostname: plan.Hostname.ValueString()},
		},
	}

	instance, err := r.client.CreateCloud(deployReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Utho Cloud instance",
			fmt.Sprintf("%s", err),
		)
		return
	}

	plan.ID = types.StringValue(instance.CloudID)
	plan.IP = types.StringValue(instance.IP)
	plan.Status = types.StringValue(instance.Status)
	plan.PowerStatus = types.StringValue("Running")
	plan.CreatedAt = types.StringValue(instance.CreatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *CloudResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CloudResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instance, err := r.client.GetCloud(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading Utho Cloud instance", fmt.Sprintf("%s", err))
		return
	}

	if instance == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.IP = types.StringValue(instance.IP)
	state.Status = types.StringValue(instance.Status)
	state.PowerStatus = types.StringValue(instance.PowerStatus)
	state.Hostname = types.StringValue(instance.Hostname)
	state.DCSlug = types.StringValue(instance.DCSlug)
	state.BillingCycle = types.StringValue(instance.BillingCycle)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CloudResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"Utho Cloud instances cannot be updated in place. Destroy and recreate.",
	)
}

func (r *CloudResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CloudResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteEBS := false
	if !state.DeleteEBS.IsNull() {
		deleteEBS = state.DeleteEBS.ValueBool()
	}

	// Pass hostname — required by the Utho delete API
	err := r.client.DeleteCloud(state.ID.ValueString(), state.Hostname.ValueString(), deleteEBS)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Utho Cloud instance", fmt.Sprintf("%s", err))
		return
	}
}

func (r *CloudResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	instance, err := r.client.GetCloud(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing Utho Cloud instance",
			fmt.Sprintf("%s", err),
		)
		return
	}

	if instance == nil {
		resp.Diagnostics.AddError(
			"Cloud instance not found",
			fmt.Sprintf("No instance found with ID: %s", req.ID),
		)
		return
	}

	state := CloudResourceModel{
		ID:             types.StringValue(instance.CloudID),
		IP:             types.StringValue(instance.IP),
		Status:         types.StringValue(instance.Status),
		PowerStatus:    types.StringValue(instance.PowerStatus),
		Hostname:       types.StringValue(instance.Hostname),
		DCSlug:         types.StringValue(instance.DCSlug),
		BillingCycle:   types.StringValue(instance.BillingCycle),
		CreatedAt:      types.StringValue(instance.CreatedAt),
		PlanID:         types.StringValue(""),
		Auth:           types.StringValue(""),
		EnablePublicIP: types.StringNull(),
		SubnetRequired: types.StringNull(),
		CPUModel:       types.StringNull(),
		EnableBackup:   types.StringNull(),
		Support:        types.StringNull(),
		Firewall:       types.StringNull(),
		RootPassword:   types.StringNull(),
		SSHKeys:        types.StringNull(),
		Image:          types.StringNull(),
		VPC:            types.StringNull(),
		SnapshotID:     types.StringNull(),
		BackupID:       types.StringNull(),
		ISO:            types.StringNull(),
		Stack:          types.StringNull(),
		DeleteEBS:      types.BoolNull(),
		EBS:            []EBSVolumeModel{},
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
