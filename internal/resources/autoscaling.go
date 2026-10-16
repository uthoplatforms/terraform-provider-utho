package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
)

// ══════════════════════════════════════════════════════════════
// AUTOSCALING GROUP — utho_autoscaling
// ══════════════════════════════════════════════════════════════

type AutoScalingResource struct{ client *client.Client }

type ASGPolicyModel struct {
	Name     types.String `tfsdk:"name"`
	Type     types.String `tfsdk:"type"`
	Compare  types.String `tfsdk:"compare"`
	Value    types.String `tfsdk:"value"`
	Adjust   types.Int64  `tfsdk:"adjust"`
	Period   types.String `tfsdk:"period"`
	Cooldown types.String `tfsdk:"cooldown"`
}

type ASGScheduleModel struct {
	Name               types.String `tfsdk:"name"`
	DesiredSize        types.String `tfsdk:"desiredsize"`
	Timezone           types.String `tfsdk:"timezone"`
	Recurrence         types.String `tfsdk:"recurrence"`
	RecurrenceDuration types.String `tfsdk:"recurrence_duration"`
	RecurrenceWeek     types.String `tfsdk:"recurrence_week"`
	SelectedTime       types.String `tfsdk:"selected_time"`
	SelectedDate       types.String `tfsdk:"selected_date"`
	StartDate          types.String `tfsdk:"start_date"`
}

type AutoScalingModel struct {
	ID              types.String       `tfsdk:"id"`
	Name            types.String       `tfsdk:"name"`
	DCSlug          types.String       `tfsdk:"dcslug"`
	MinSize         types.String       `tfsdk:"minsize"`
	MaxSize         types.String       `tfsdk:"maxsize"`
	DesiredSize     types.String       `tfsdk:"desiredsize"`
	PlanID          types.String       `tfsdk:"planid"`
	PlanName        types.String       `tfsdk:"planname"`
	OSDiskSize      types.Int64        `tfsdk:"os_disk_size"`
	PublicIPEnabled types.Int64        `tfsdk:"public_ip_enabled"`
	ImageName       types.String       `tfsdk:"image_name"`
	SnapshotID      types.String       `tfsdk:"snapshotid"`
	Stack           types.String       `tfsdk:"stack"`
	StackID         types.String       `tfsdk:"stackid"`
	StackImage      types.String       `tfsdk:"stackimage"`
	VPC             types.String       `tfsdk:"vpc"`
	LoadBalancers   types.String       `tfsdk:"load_balancers"`
	SecurityGroups  types.String       `tfsdk:"security_groups"`
	TargetGroups    types.String       `tfsdk:"target_groups"`
	BackupID        types.String       `tfsdk:"backupid"`
	CPUModel        types.String       `tfsdk:"cpumodel"`
	Status          types.String       `tfsdk:"status"`
	CreatedAt       types.String       `tfsdk:"created_at"`
	Policies        []ASGPolicyModel   `tfsdk:"policies"`
	Schedules       []ASGScheduleModel `tfsdk:"schedules"`
}

func NewAutoScalingResource() resource.Resource { return &AutoScalingResource{} }

func (r *AutoScalingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	asg, err := r.client.GetAutoScaling(req.ID)
	if err != nil || asg == nil {
		resp.Diagnostics.AddError("Error importing auto scaling group", fmt.Sprintf("Could not find ASG with ID %s", req.ID))
		return
	}
	// Helper to get string from interface{}
	strVal := func(v interface{}) string {
		if v == nil {
			return ""
		}
		if s, ok := v.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	// Helper to get int64 from interface{}
	int64Val := func(v interface{}) int64 {
		if v == nil {
			return 0
		}
		switch n := v.(type) {
		case float64:
			return int64(n)
		case string:
			if n == "" {
				return 0
			}
			var i int64
			fmt.Sscanf(n, "%d", &i)
			return i
		}
		return 0
	}
	state := AutoScalingModel{
		ID:              types.StringValue(asg.ID),
		Name:            types.StringValue(asg.Name),
		DCSlug:          types.StringValue(asg.DCSlug),
		MinSize:         types.StringValue(asg.MinSize),
		MaxSize:         types.StringValue(asg.MaxSize),
		DesiredSize:     types.StringValue(asg.DesiredSize),
		PlanID:          types.StringValue(asg.PlanID),
		PlanName:        types.StringValue(asg.PlanName),
		Status:          types.StringValue(asg.Status),
		CreatedAt:       types.StringValue(asg.CreatedAt),
		OSDiskSize:      types.Int64Value(int64Val(asg.OSDiskSize)),
		PublicIPEnabled: types.Int64Value(int64Val(asg.PublicIPEnabled)),
		ImageName:       types.StringValue(asg.ImageName),
		SnapshotID:      types.StringValue(asg.SnapshotID),
		Stack:           types.StringValue(asg.Stack),
		StackID:         types.StringValue(asg.StackID),
		StackImage:      types.StringValue(asg.StackImage),
		VPC:             types.StringValue(strVal(asg.VPC)),
		LoadBalancers:   types.StringValue(""),
		SecurityGroups:  types.StringValue(""),
		TargetGroups:    types.StringValue(""),
		BackupID:        types.StringValue(asg.BackupID),
		CPUModel:        types.StringValue(asg.CPUModel),
		Policies:        []ASGPolicyModel{},
		Schedules:       []ASGScheduleModel{},
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AutoScalingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_autoscaling"
}

func (r *AutoScalingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage Utho Auto Scaling groups. Automatically scales cloud instances based on CPU/RAM policies or schedules.",
		Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"status":     schema.StringAttribute{Computed: true, Description: "Auto scaling group status."},
			"created_at": schema.StringAttribute{Computed: true, Description: "Creation timestamp.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},

			// Required
			"name":              schema.StringAttribute{Required: true, Description: "Unique auto scaling group name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"dcslug":            schema.StringAttribute{Required: true, Description: "Data center slug.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"minsize":           schema.StringAttribute{Required: true, Description: "Minimum number of instances."},
			"maxsize":           schema.StringAttribute{Required: true, Description: "Maximum number of instances."},
			"desiredsize":       schema.StringAttribute{Required: true, Description: "Desired number of instances."},
			"planid":            schema.StringAttribute{Required: true, Description: "Plan ID for the instance size (e.g. 10404). Use the utho_cloud_plans data source to find available plan IDs for your dcslug.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"planname":          schema.StringAttribute{Required: true, Description: "Plan name or label corresponding to the planid (e.g. basic, standard). This is a display label sent to the API alongside planid — it does not need to match the planid exactly as the API does not cross-validate them. Use the name field from utho_cloud_plans data source.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"os_disk_size":      schema.Int64Attribute{Required: true, Description: "OS disk size in GB.", PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()}},
			"public_ip_enabled": schema.Int64Attribute{Required: true, Description: "Enable public IP: 1 or 0."},

			// Deployment source — at least one required (validated at Create)
			"image_name": schema.StringAttribute{Optional: true, Description: "Boot image name. For snapshot deployment: set to the snapshot name. For stack deployment: set to the stack image slug (same as stackimage). Required when using snapshot or stack deployment.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"snapshotid": schema.StringAttribute{Optional: true, Description: "Snapshot ID for snapshot-based deployment. Use this OR stackid+stackimage — not both. When set, also set image_name to the snapshot name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"stack":      schema.StringAttribute{Optional: true, Description: "Stack slug (same as stackid). Use this with stackid and stackimage for stack-based deployment.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"stackid":    schema.StringAttribute{Optional: true, Description: "Stack ID for stack-based deployment. Use this OR snapshotid — not both. When set, also set stackimage and image_name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"stackimage": schema.StringAttribute{Optional: true, Description: "Stack image slug. Required when using stack-based deployment (stackid). Also set image_name to the same value.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},

			// Optional
			"vpc":             schema.StringAttribute{Optional: true, Description: "VPC subnet ID.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"load_balancers":  schema.StringAttribute{Optional: true, Description: "Load balancer ID to attach. Can be updated in place."},
			"security_groups": schema.StringAttribute{Optional: true, Description: "Security group ID to attach. Can be updated in place."},
			"target_groups":   schema.StringAttribute{Optional: true, Description: "Target group IDs. Can be updated in place."},
			"backupid":        schema.StringAttribute{Optional: true, Description: "Backup ID.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"cpumodel":        schema.StringAttribute{Optional: true, Description: "CPU model: amd or intel.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},

			// Policies — at least 1 required
			"policies": schema.ListNestedAttribute{
				Required:    true,
				Description: "Scaling policies. At least one policy is required.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":     schema.StringAttribute{Required: true, Description: "Policy name."},
						"type":     schema.StringAttribute{Required: true, Description: "Metric: cpu or ram."},
						"compare":  schema.StringAttribute{Required: true, Description: "Comparison: above or below."},
						"value":    schema.StringAttribute{Required: true, Description: "Threshold value (e.g. 80 for 80%)."},
						"adjust":   schema.Int64Attribute{Required: true, Description: "Instances to add/remove."},
						"period":   schema.StringAttribute{Required: true, Description: "Evaluation period (e.g. 5m)."},
						"cooldown": schema.StringAttribute{Required: true, Description: "Cooldown in seconds (e.g. 300)."},
					},
				},
			},

			// Schedules — optional, 0 or more
			"schedules": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Scheduled scaling. Optional, 0 or more.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":                schema.StringAttribute{Required: true, Description: "Schedule name."},
						"desiredsize":         schema.StringAttribute{Required: true, Description: "Desired instance count."},
						"timezone":            schema.StringAttribute{Required: true, Description: "Timezone (e.g. Asia/Kolkata)."},
						"recurrence":          schema.StringAttribute{Required: true, Description: "Recurrence expression (e.g. Every day  12:00)."},
						"recurrence_duration": schema.StringAttribute{Required: true, Description: "Recurrence duration (e.g. Every day)."},
						"recurrence_week":     schema.StringAttribute{Optional: true, Description: "Recurrence week day."},
						"selected_time":       schema.StringAttribute{Required: true, Description: "Time (HH:MM)."},
						"selected_date":       schema.StringAttribute{Required: true, Description: "Date (YYYY-MM-DD)."},
						"start_date":          schema.StringAttribute{Required: true, Description: "Start datetime ISO 8601."},
					},
				},
			},
		},
	}
}

func (r *AutoScalingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AutoScalingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AutoScalingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate at least 1 policy
	if len(plan.Policies) == 0 {
		resp.Diagnostics.AddError("Validation Error", "At least one scaling policy is required in the policies block.")
		return
	}

	// Validate deployment source
	hasSnapshot := !plan.SnapshotID.IsNull() && plan.SnapshotID.ValueString() != ""
	hasStack := !plan.StackID.IsNull() && plan.StackID.ValueString() != ""
	hasImage := false // image_id not used for autoscaling
	if !hasSnapshot && !hasStack && !hasImage {
		resp.Diagnostics.AddError("Validation Error", "At least one deployment source is required: set snapshotid, stackid+stackimage, or image_name.")
		return
	}

	var policies []client.AutoScalingPolicy
	for _, p := range plan.Policies {
		policies = append(policies, client.AutoScalingPolicy{
			Name: p.Name.ValueString(), Type: p.Type.ValueString(),
			Compare: p.Compare.ValueString(), Value: p.Value.ValueString(),
			Adjust: int(p.Adjust.ValueInt64()), Period: p.Period.ValueString(),
			Cooldown: p.Cooldown.ValueString(),
		})
	}

	var schedules []client.AutoScalingSchedule
	for _, s := range plan.Schedules {
		schedules = append(schedules, client.AutoScalingSchedule{
			Name: s.Name.ValueString(), DesiredSize: s.DesiredSize.ValueString(),
			Timezone: s.Timezone.ValueString(), Recurrence: s.Recurrence.ValueString(),
			RecurrenceDuration: s.RecurrenceDuration.ValueString(),
			RecurrenceWeek:     s.RecurrenceWeek.ValueString(),
			SelectedTime:       s.SelectedTime.ValueString(),
			SelectedDate:       s.SelectedDate.ValueString(),
			StartDate:          s.StartDate.ValueString(),
		})
	}

	id, err := r.client.CreateAutoScaling(&client.AutoScalingCreateRequest{
		Name: plan.Name.ValueString(), OSDiskSize: int(plan.OSDiskSize.ValueInt64()),
		DCSlug: plan.DCSlug.ValueString(), MinSize: plan.MinSize.ValueString(),
		MaxSize: plan.MaxSize.ValueString(), DesiredSize: plan.DesiredSize.ValueString(),
		PlanID: plan.PlanID.ValueString(), PlanName: plan.PlanName.ValueString(),
		InstanceTemplateID: "none",
		ImageID:            "",
		ImageName:          plan.ImageName.ValueString(),
		PublicIPEnabled:    int(plan.PublicIPEnabled.ValueInt64()),
		VPC:                plan.VPC.ValueString(),
		LoadBalancers:      plan.LoadBalancers.ValueString(),
		SecurityGroups:     plan.SecurityGroups.ValueString(),
		TargetGroups:       plan.TargetGroups.ValueString(),
		SnapshotID:         plan.SnapshotID.ValueString(),
		Stack:              plan.Stack.ValueString(),
		StackID:            plan.StackID.ValueString(),
		StackImage:         plan.StackImage.ValueString(),
		BackupID:           plan.BackupID.ValueString(),
		CPUModel:           plan.CPUModel.ValueString(),
		Policies:           policies,
		Schedules:          schedules,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating auto scaling group", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	plan.Status = types.StringValue("Deploying")
	plan.CreatedAt = types.StringValue("")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Poll until Active (up to 6 minutes)
	if waitErr := r.client.WaitForASGReady(id); waitErr != nil {
		if waitErr.Error() == "STILL_PROVISIONING" {
			resp.Diagnostics.AddWarning(
				"Auto scaling group still provisioning",
				fmt.Sprintf("Group %q (ID: %s) was created but is still deploying instances. "+
					"Run terraform apply again once provisioning completes in the Utho Console.", plan.Name.ValueString(), id),
			)
		}
	}

	// Refresh status
	if asg, err := r.client.GetAutoScaling(id); err == nil && asg != nil {
		plan.Status = types.StringValue(asg.Status)
		plan.CreatedAt = types.StringValue(asg.CreatedAt)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *AutoScalingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AutoScalingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	asg, err := r.client.GetAutoScaling(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading auto scaling group", fmt.Sprintf("%s", err))
		return
	}
	if asg == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Status = types.StringValue(asg.Status)
	state.CreatedAt = types.StringValue(asg.CreatedAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AutoScalingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state AutoScalingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	// Block updates when still provisioning
	if strings.ToLower(state.Status.ValueString()) == "deploying" {
		resp.Diagnostics.AddError(
			"Auto scaling group is still provisioning",
			"Cannot update while the group is deploying. Wait for provisioning to complete and run terraform apply again.",
		)
		return
	}

	// Update scaling config if changed
	if plan.MinSize != state.MinSize || plan.MaxSize != state.MaxSize ||
		plan.DesiredSize != state.DesiredSize || plan.PublicIPEnabled != state.PublicIPEnabled ||
		plan.SnapshotID != state.SnapshotID || plan.StackID != state.StackID || plan.StackImage != state.StackImage {
		updateReq := &client.AutoScalingUpdateRequest{
			Name:        plan.Name.ValueString(),
			MinSize:     plan.MinSize.ValueString(),
			MaxSize:     plan.MaxSize.ValueString(),
			DesiredSize: plan.DesiredSize.ValueString(),
		}
		// Send only one deployment source
		if plan.SnapshotID.ValueString() != "" {
			updateReq.SnapshotID = plan.SnapshotID.ValueString()
		} else if plan.StackID.ValueString() != "" {
			updateReq.Stack = plan.Stack.ValueString()
			updateReq.StackID = plan.StackID.ValueString()
			updateReq.StackImage = plan.StackImage.ValueString()
		}
		if err := r.client.UpdateAutoScaling(id, updateReq); err != nil {
			resp.Diagnostics.AddError("Error updating auto scaling group", fmt.Sprintf("%s", err))
			return
		}
	}

	// Handle security group change
	oldSG := state.SecurityGroups.ValueString()
	newSG := plan.SecurityGroups.ValueString()
	if oldSG != newSG {
		if oldSG != "" {
			r.client.DetachASGSecurityGroup(id, oldSG)
		}
		if newSG != "" {
			if err := r.client.AttachASGSecurityGroup(id, newSG); err != nil {
				resp.Diagnostics.AddError("Error updating security group", fmt.Sprintf("%s", err))
				return
			}
		}
	}

	// Handle load balancer change
	oldLB := state.LoadBalancers.ValueString()
	newLB := plan.LoadBalancers.ValueString()
	if oldLB != newLB {
		if oldLB != "" {
			r.client.DetachASGLoadBalancer(id, oldLB)
		}
		if newLB != "" {
			if err := r.client.AttachASGLoadBalancer(id, newLB); err != nil {
				resp.Diagnostics.AddError("Error updating load balancer", fmt.Sprintf("%s", err))
				return
			}
		}
	}

	// Handle target group change
	oldTG := state.TargetGroups.ValueString()
	newTG := plan.TargetGroups.ValueString()
	if oldTG != newTG {
		if oldTG != "" {
			r.client.DetachASGTargetGroup(id, oldTG)
		}
		if newTG != "" {
			if err := r.client.AttachASGTargetGroup(id, newTG); err != nil {
				resp.Diagnostics.AddError("Error updating target group", fmt.Sprintf("%s", err))
				return
			}
		}
	}

	plan.ID = state.ID
	plan.Status = state.Status
	plan.CreatedAt = state.CreatedAt
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *AutoScalingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AutoScalingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteAutoScaling(state.ID.ValueString(), state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting auto scaling group", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// SCALING POLICY — utho_autoscaling_policy
// ══════════════════════════════════════════════════════════════

type ScalingPolicyResource struct{ client *client.Client }

type ScalingPolicyModel struct {
	ID          types.String `tfsdk:"id"`
	ASGId       types.String `tfsdk:"asg_id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Compare     types.String `tfsdk:"compare"`
	Value       types.String `tfsdk:"value"`
	Adjust      types.Int64  `tfsdk:"adjust"`
	Period      types.String `tfsdk:"period"`
	Cooldown    types.String `tfsdk:"cooldown"`
	ScalingType types.String `tfsdk:"scaling_type"`
}

func NewScalingPolicyResource() resource.Resource { return &ScalingPolicyResource{} }

func (r *ScalingPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_autoscaling_policy"
}

func (r *ScalingPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage scaling policies for a Utho Auto Scaling group.",
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"asg_id":       schema.StringAttribute{Required: true, Description: "Auto scaling group ID.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"name":         schema.StringAttribute{Required: true, Description: "Policy name."},
			"type":         schema.StringAttribute{Required: true, Description: "Metric type: cpu or ram.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"compare":      schema.StringAttribute{Required: true, Description: "Comparison: above or below."},
			"value":        schema.StringAttribute{Required: true, Description: "Threshold value (e.g. 80 for 80%)."},
			"adjust":       schema.Int64Attribute{Required: true, Description: "Number of instances to add or remove."},
			"period":       schema.StringAttribute{Required: true, Description: "Evaluation period (e.g. 1m, 5m)."},
			"cooldown":     schema.StringAttribute{Required: true, Description: "Cooldown period in seconds."},
			"scaling_type": schema.StringAttribute{Required: true, Description: "Scaling type: horizontal.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *ScalingPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ScalingPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ScalingPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateScalingPolicy(&client.ScalingPolicyCreateRequest{
		Name: plan.Name.ValueString(), Type: plan.Type.ValueString(),
		Compare: plan.Compare.ValueString(), Value: plan.Value.ValueString(),
		Adjust: fmt.Sprintf("%d", plan.Adjust.ValueInt64()),
		Period: plan.Period.ValueString(), Cooldown: plan.Cooldown.ValueString(),
		Product: "autoscaling", ProductID: plan.ASGId.ValueString(),
		ScalingType: plan.ScalingType.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating scaling policy", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ScalingPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Policies are embedded in ASG — no standalone get endpoint
	// Just keep state as-is
	var state ScalingPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ScalingPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ScalingPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateScalingPolicy(plan.ID.ValueString(), &client.ScalingPolicyUpdateRequest{
		Name:     plan.Name.ValueString(),
		Type:     plan.Type.ValueString(),
		Compare:  plan.Compare.ValueString(),
		Value:    plan.Value.ValueString(),
		Adjust:   int(plan.Adjust.ValueInt64()),
		Period:   plan.Period.ValueString(),
		Cooldown: plan.Cooldown.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating scaling policy", fmt.Sprintf("%s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ScalingPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ScalingPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteScalingPolicy(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting scaling policy", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// SCALING SCHEDULE — utho_autoscaling_schedule
// ══════════════════════════════════════════════════════════════

type ScalingScheduleResource struct{ client *client.Client }

type ScalingScheduleModel struct {
	ID          types.String `tfsdk:"id"`
	ASGId       types.String `tfsdk:"asg_id"`
	Name        types.String `tfsdk:"name"`
	DesiredSize types.String `tfsdk:"desiredsize"`
	Timezone    types.String `tfsdk:"timezone"`
	Recurrence  types.String `tfsdk:"recurrence"`
	StartDate   types.String `tfsdk:"start_date"`
	Status      types.Int64  `tfsdk:"status"`
}

func NewScalingScheduleResource() resource.Resource { return &ScalingScheduleResource{} }

func (r *ScalingScheduleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_autoscaling_schedule"
}

func (r *ScalingScheduleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage scheduled scaling policies for a Utho Auto Scaling group.",
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"asg_id":      schema.StringAttribute{Required: true, Description: "Auto scaling group ID.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"name":        schema.StringAttribute{Required: true, Description: "Schedule name."},
			"desiredsize": schema.StringAttribute{Required: true, Description: "Desired instance count at schedule time."},
			"timezone":    schema.StringAttribute{Required: true, Description: "Timezone (e.g. Asia/Kolkata)."},
			"recurrence":  schema.StringAttribute{Required: true, Description: "Recurrence expression."},
			"start_date":  schema.StringAttribute{Required: true, Description: "Start datetime (e.g. 2026-09-15 17:04:00)."},
			"status":      schema.Int64Attribute{Optional: true, Description: "Schedule status: 1 (active) or 0 (inactive). Default: 1."},
		},
	}
}

func (r *ScalingScheduleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ScalingScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ScalingScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Check ASG is Active before adding schedule
	asg, err := r.client.GetAutoScaling(plan.ASGId.ValueString())
	if err == nil && asg != nil && strings.ToLower(asg.Status) == "deploying" {
		resp.Diagnostics.AddError(
			"Auto scaling group is still provisioning",
			"Cannot add a schedule while the group is deploying. Wait for provisioning to complete.",
		)
		return
	}
	// Schedules are created via the ASG update — use a composite ID
	plan.ID = types.StringValue(fmt.Sprintf("%s:schedule:%s", plan.ASGId.ValueString(), plan.Name.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ScalingScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Schedules are embedded in ASG — no standalone get endpoint
	// Just keep state as-is
	var state ScalingScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ScalingScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ScalingScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateScalingSchedule(plan.ASGId.ValueString(), plan.ID.ValueString(), &client.ScalingScheduleUpdateRequest{
		GroupID:     plan.ASGId.ValueString(),
		ID:          plan.ID.ValueString(),
		Name:        plan.Name.ValueString(),
		DesiredSize: plan.DesiredSize.ValueString(),
		Timezone:    plan.Timezone.ValueString(),
		Recurrence:  plan.Recurrence.ValueString(),
		StartDate:   plan.StartDate.ValueString(),
		Status:      int(plan.Status.ValueInt64()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating scaling schedule", fmt.Sprintf("%s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ScalingScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ScalingScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteScalingSchedule(state.ASGId.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting scaling schedule", fmt.Sprintf("%s", err))
	}
}
