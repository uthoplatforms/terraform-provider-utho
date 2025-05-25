package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/uthoplatforms/utho-go/utho"
)

// implement resource interfaces.
var (
	_ resource.Resource                = &AutoScalingResource{}
	_ resource.ResourceWithConfigure   = &AutoScalingResource{}
	_ resource.ResourceWithImportState = &AutoScalingResource{}
)

// NewAutoScalingResource is a helper function to simplify the provider implementation.
func NewAutoScalingResource() resource.Resource {
	return &AutoScalingResource{}
}

// AutoScalingResource is the resource implementation.
type AutoScalingResource struct {
	client utho.Client
}

// AutoScalingResource is the model implementation.
type AutoScalingResourceModel struct {
	ID                 types.String  `tfsdk:"id"`
	Userid             types.String  `tfsdk:"userid"`
	Name               types.String  `tfsdk:"name"`
	Dcslug             types.String  `tfsdk:"dcslug"`
	Minsize            types.String  `tfsdk:"minsize"`
	Maxsize            types.String  `tfsdk:"maxsize"`
	Desiredsize        types.String  `tfsdk:"desiredsize"`
	Planid             types.String  `tfsdk:"planid"`
	Planname           types.String  `tfsdk:"planname"`
	InstanceTemplateid types.String  `tfsdk:"instance_templateid"`
	Image              types.String  `tfsdk:"image"`
	ImageName          types.String  `tfsdk:"image_name"`
	Snapshotid         types.String  `tfsdk:"snapshotid"`
	Status             types.String  `tfsdk:"status"`
	CreatedAt          types.String  `tfsdk:"created_at"`
	SuspendedAt        types.String  `tfsdk:"suspended_at"`
	StoppedAt          types.String  `tfsdk:"stopped_at"`
	StartedAt          types.String  `tfsdk:"started_at"`
	DeletedAt          types.String  `tfsdk:"deleted_at"`
	PublicIPEnabled    types.Bool    `tfsdk:"public_ip_enabled"`
	CooldownTill       types.String  `tfsdk:"cooldown_till"`
	Backupid           types.String  `tfsdk:"backupid"`
	Stackid            types.String  `tfsdk:"stackid"`
	Stackimage         types.String  `tfsdk:"stackimage"`
	VpcID              types.String  `tfsdk:"vpc_id"`
	LoadbalancersID    types.String  `tfsdk:"loadbalancers_id"`
	SecurityGroupID    types.String  `tfsdk:"security_group_id"`
	TargetGroupsID     types.String  `tfsdk:"target_groups_id"`
	OsDiskSize         types.Int64   `tfsdk:"os_disk_size"`
	Policies           []PolicyModel `tfsdk:"policies"`
	Schedules          types.List    `tfsdk:"schedules"`
	Vpc                types.List    `tfsdk:"vpc"`
	Loadbalancers      types.List    `tfsdk:"load_balancers"`
	TargetGroups       types.List    `tfsdk:"target_groups"`
	SecurityGroups     types.List    `tfsdk:"security_groups"`
	Instances          types.List    `tfsdk:"instances"`
	Dclocation         types.Object  `tfsdk:"dclocation"`
	Plan               types.Object  `tfsdk:"plan"`
}
type AutoScalingVpcModel struct {
	Total     types.Int64  `tfsdk:"total"`
	Available types.Int64  `tfsdk:"available"`
	Network   types.String `tfsdk:"network"`
	Name      types.String `tfsdk:"name"`
	Size      types.String `tfsdk:"size"`
	Dcslug    types.String `tfsdk:"dcslug"`
}
type AutoScalingLoadbalancersModel struct {
	ID   types.String `tfsdk:"lbid"`
	Name types.String `tfsdk:"name"`
	IP   types.String `tfsdk:"ip"`
}
type AutoScalingTargetGroupModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Protocol types.String `tfsdk:"protocol"`
	Port     types.String `tfsdk:"port"`
}
type SecurityGroupModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}
type InstancesModel struct {
	Cloudid   types.String `tfsdk:"cloudid"`
	Hostname  types.String `tfsdk:"hostname"`
	CreatedAt types.String `tfsdk:"created_at"`
	IP        types.String `tfsdk:"ip"`
	Status    types.String `tfsdk:"status"`
}
type PolicyModel struct {
	ID                 types.String `tfsdk:"id"`
	Userid             types.String `tfsdk:"userid"`
	Product            types.String `tfsdk:"product"`
	Productid          types.String `tfsdk:"productid"`
	Groupid            types.String `tfsdk:"groupid"`
	Name               types.String `tfsdk:"name"`
	Type               types.String `tfsdk:"type"`
	Adjust             types.String `tfsdk:"adjust"`
	Period             types.String `tfsdk:"period"`
	Cooldown           types.String `tfsdk:"cooldown"`
	CooldownTill       types.String `tfsdk:"cooldown_till"`
	Compare            types.String `tfsdk:"compare"`
	Value              types.String `tfsdk:"value"`
	AlertID            types.String `tfsdk:"alert_id"`
	Status             types.String `tfsdk:"status"`
	KubernetesID       types.String `tfsdk:"kubernetes_id"`
	KubernetesNodepool types.String `tfsdk:"kubernetes_nodepool"`
	Cloudid            types.String `tfsdk:"cloudid"`
	Maxsize            types.String `tfsdk:"maxsize"`
	Minsize            types.String `tfsdk:"minsize"`
}
type ScheduleModel struct {
	ID          types.String `tfsdk:"id"`
	Groupid     types.String `tfsdk:"groupid"`
	Name        types.String `tfsdk:"name"`
	Desiredsize types.String `tfsdk:"desiredsize"`
	Recurrence  types.String `tfsdk:"recurrence"`
	StartDate   types.String `tfsdk:"start_date"`
	Status      types.String `tfsdk:"status"`
	Timezone    types.String `tfsdk:"timezone"`
}
type DclocationModel struct {
	Location types.String `tfsdk:"location"`
	Country  types.String `tfsdk:"country"`
	Dc       types.String `tfsdk:"dc"`
	Dccc     types.String `tfsdk:"dccc"`
}
type AutoScalingPlanModel struct {
	Planid         types.String `tfsdk:"planid"`
	RAM            types.String `tfsdk:"ram"`
	CPU            types.String `tfsdk:"cpu"`
	Disk           types.String `tfsdk:"disk"`
	Bandwidth      types.String `tfsdk:"bandwidth"`
	DedicatedVcore types.String `tfsdk:"dedicated_vcore"`
}

// Metadata returns the resource type name.
func (s *AutoScalingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auto_scaling"
}

// Configure adds the provider configured client to the data source.
func (d *AutoScalingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(utho.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected AutoScaling Data Source Configure Type",
			fmt.Sprintf("Expected utho.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

// Schema defines the schema for the resource.
func (s *AutoScalingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true},
			"name":        schema.StringAttribute{Required: true, Description: "Provide AUTOSCALING name eg: autoscaling-ywqo2pmc"},
			"minsize":     schema.StringAttribute{Required: true},
			"maxsize":     schema.StringAttribute{Required: true},
			"desiredsize": schema.StringAttribute{Required: true},
			"dcslug": schema.StringAttribute{Required: true, Description: "Provide dcslug eg: innoida",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"planid": schema.StringAttribute{Required: true, Description: "Provide the planid eg: 10045",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"os_disk_size": schema.Int64Attribute{Required: true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"planname": schema.StringAttribute{Required: true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"instance_templateid": schema.StringAttribute{Required: true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"public_ip_enabled": schema.BoolAttribute{Required: true,
				PlanModifiers: []planmodifier.Bool{},
			},
			"stackid": schema.StringAttribute{Required: true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"stackimage": schema.StringAttribute{Required: true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"vpc_id": schema.StringAttribute{Required: true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"loadbalancers_id": schema.StringAttribute{Optional: true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"security_group_id": schema.StringAttribute{Optional: true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"target_groups_id": schema.StringAttribute{Optional: true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"policies": schema.ListNestedAttribute{
				Optional:    true,
				Description: "policies",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Computed: true, Description: "id"},
						"name": schema.StringAttribute{Required: true, Description: "name",
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"type": schema.StringAttribute{Required: true, Description: "type",
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"compare": schema.StringAttribute{Required: true, Description: "compare",
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"value": schema.StringAttribute{Required: true, Description: "value",
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"adjust": schema.StringAttribute{Required: true, Description: "adjust",
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"period": schema.StringAttribute{Required: true, Description: "period",
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"cooldown": schema.StringAttribute{Required: true, Description: "cooldown",
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"userid":              schema.StringAttribute{Computed: true, Description: "userid"},
						"product":             schema.StringAttribute{Computed: true, Description: "product"},
						"productid":           schema.StringAttribute{Computed: true, Description: "productid"},
						"groupid":             schema.StringAttribute{Computed: true, Description: "groupid"},
						"cooldown_till":       schema.StringAttribute{Computed: true, Description: "cooldown_till"},
						"alert_id":            schema.StringAttribute{Computed: true, Description: "alert_id"},
						"status":              schema.StringAttribute{Computed: true, Description: "status"},
						"kubernetes_id":       schema.StringAttribute{Computed: true, Description: "kubernetes_id"},
						"kubernetes_nodepool": schema.StringAttribute{Computed: true, Description: "kubernetes_nodepool"},
						"cloudid":             schema.StringAttribute{Computed: true, Description: "cloudid"},
						"maxsize":             schema.StringAttribute{Computed: true, Description: "maxsize"},
						"minsize":             schema.StringAttribute{Computed: true, Description: "minsize"},
					},
				},
			},
			////////////////////
			"userid":        schema.StringAttribute{Computed: true},
			"image":         schema.StringAttribute{Computed: true},
			"image_name":    schema.StringAttribute{Computed: true},
			"snapshotid":    schema.StringAttribute{Computed: true},
			"status":        schema.StringAttribute{Computed: true},
			"created_at":    schema.StringAttribute{Computed: true},
			"suspended_at":  schema.StringAttribute{Computed: true},
			"stopped_at":    schema.StringAttribute{Computed: true},
			"started_at":    schema.StringAttribute{Computed: true},
			"deleted_at":    schema.StringAttribute{Computed: true},
			"cooldown_till": schema.StringAttribute{Computed: true},
			"backupid":      schema.StringAttribute{Computed: true},
			"vpc": schema.ListNestedAttribute{
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"total":     schema.Int64Attribute{Computed: true},
						"available": schema.Int64Attribute{Computed: true},
						"network":   schema.StringAttribute{Computed: true},
						"name":      schema.StringAttribute{Computed: true},
						"size":      schema.StringAttribute{Computed: true},
						"dcslug":    schema.StringAttribute{Computed: true},
					},
				},
			},
			"load_balancers": schema.ListNestedAttribute{
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"lbid": schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
						"ip":   schema.StringAttribute{Computed: true},
					},
				},
			},
			"target_groups": schema.ListNestedAttribute{
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":       schema.StringAttribute{Computed: true},
						"name":     schema.StringAttribute{Computed: true},
						"protocol": schema.StringAttribute{Computed: true},
						"port":     schema.StringAttribute{Computed: true},
					},
				},
			},
			"security_groups": schema.ListNestedAttribute{
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":   schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
					},
				},
			},
			"instances": schema.ListNestedAttribute{
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cloudid":    schema.StringAttribute{Computed: true},
						"hostname":   schema.StringAttribute{Computed: true},
						"created_at": schema.StringAttribute{Computed: true},
						"ip":         schema.StringAttribute{Computed: true},
						"status":     schema.StringAttribute{Computed: true},
					},
				},
			},
			"schedules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.StringAttribute{Computed: true},
						"groupid":     schema.StringAttribute{Computed: true},
						"name":        schema.StringAttribute{Computed: true},
						"desiredsize": schema.StringAttribute{Computed: true},
						"recurrence":  schema.StringAttribute{Computed: true},
						"start_date":  schema.StringAttribute{Computed: true},
						"status":      schema.StringAttribute{Computed: true},
						"timezone":    schema.StringAttribute{Computed: true},
					},
				},
			},
			"dclocation": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "",
				Attributes: map[string]schema.Attribute{
					"location": schema.StringAttribute{Computed: true, Description: ""},
					"country":  schema.StringAttribute{Computed: true, Description: ""},
					"dc":       schema.StringAttribute{Computed: true, Description: ""},
					"dccc":     schema.StringAttribute{Computed: true, Description: ""},
				},
			},
			"plan": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "",
				Attributes: map[string]schema.Attribute{
					"planid":          schema.StringAttribute{Computed: true},
					"ram":             schema.StringAttribute{Computed: true},
					"cpu":             schema.StringAttribute{Computed: true},
					"disk":            schema.StringAttribute{Computed: true},
					"bandwidth":       schema.StringAttribute{Computed: true},
					"dedicated_vcore": schema.StringAttribute{Computed: true},
				},
			},
		},
	}
}

// Import using autoscaling as the attribute
func (s *AutoScalingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource.
func (s *AutoScalingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "create autoscaling")
	// Retrieve values from plan
	var plan AutoScalingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	policies := []utho.CreatePoliciesParams{}
	for _, v := range plan.Policies {
		policy := utho.CreatePoliciesParams{
			Name:     v.Name.ValueString(),
			Type:     v.Type.ValueString(),
			Compare:  v.Compare.ValueString(),
			Value:    v.Value.ValueString(),
			Adjust:   v.Adjust.ValueString(),
			Period:   v.Period.ValueString(),
			Cooldown: v.Cooldown.ValueString(),
		}
		policies = append(policies, policy)
	}

	schedules := []utho.CreateSchedulesParams{}
	if len(plan.Schedules.Elements()) > 0 {
		_ = plan.Schedules.ElementsAs(ctx, &schedules, false)
		for _, v := range schedules {
			schedule := utho.CreateSchedulesParams{
				Name:         v.Name,
				Desiredsize:  v.Desiredsize,
				StartDate:    time.Now(),
				SelectedTime: v.SelectedTime,
				SelectedDate: v.SelectedDate,
			}
			schedules = append(schedules, schedule)
		}

	}

	autoscalingRequest := utho.CreateAutoScalingParams{
		Name:               plan.Name.ValueString(),
		Dcslug:             plan.Dcslug.ValueString(),
		Planid:             plan.Planid.ValueString(),
		OsDiskSize:         int(plan.OsDiskSize.ValueInt64()),
		Minsize:            plan.Minsize.ValueString(),
		Maxsize:            plan.Maxsize.ValueString(),
		Desiredsize:        plan.Desiredsize.ValueString(),
		Planname:           plan.Planname.ValueString(),
		InstanceTemplateid: plan.InstanceTemplateid.ValueString(),
		PublicIPEnabled:    plan.PublicIPEnabled.ValueBool(),
		Vpc:                plan.VpcID.ValueString(),
		LoadBalancers:      plan.LoadbalancersID.ValueString(),
		Stack:              plan.Stackid.ValueString(),
		Stackid:            plan.Stackid.ValueString(),
		Stackimage:         plan.Stackimage.ValueString(),
		Policies:           policies,
		SecurityGroups:     plan.SecurityGroupID.ValueString(),
		TargetGroups:       plan.TargetGroupsID.ValueString(),
		Schedules:          schedules,
	}
	tflog.Debug(ctx, "send create autoscaling request")
	autoscaling, err := s.client.AutoScaling().Create(autoscalingRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating autoscaling",
			"Could not create autoscaling, unexpected error: "+err.Error(),
		)
		return
	}

	// get autoscaling data
	getAutoScaling, err := s.client.AutoScaling().Read(strconv.Itoa(autoscaling.ID))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading utho autoscaling",
			"Could not read utho autoscaling "+strconv.Itoa(autoscaling.ID)+": "+err.Error(),
		)
		return
	}

	osDiskSize, _ := strconv.Atoi(getAutoScaling.Plan.Disk)

	policiesModel := []PolicyModel{}
	for _, v := range getAutoScaling.Policies {
		policy := PolicyModel{
			ID:                 types.StringValue(v.ID),
			Userid:             types.StringValue(v.Userid),
			Product:            types.StringValue(v.Product),
			Productid:          types.StringValue(v.Productid),
			Groupid:            types.StringValue(v.Groupid),
			Name:               types.StringValue(v.Name),
			Type:               types.StringValue(v.Type),
			Adjust:             types.StringValue(v.Adjust),
			Period:             types.StringValue(v.Period),
			Cooldown:           types.StringValue(v.Cooldown),
			CooldownTill:       types.StringValue(v.CooldownTill),
			Compare:            types.StringValue(v.Compare),
			Value:              types.StringValue(v.Value),
			AlertID:            types.StringValue(v.AlertID),
			Status:             types.StringValue(v.Status),
			KubernetesID:       types.StringValue(v.KubernetesID),
			KubernetesNodepool: types.StringValue(v.KubernetesNodepool),
			Cloudid:            types.StringValue(v.Cloudid),
			Maxsize:            types.StringValue(v.Maxsize),
			Minsize:            types.StringValue(v.Minsize),
		}
		policiesModel = append(policiesModel, policy)
	}

	if !plan.LoadbalancersID.IsNull() {
		plan.LoadbalancersID = types.StringValue(plan.LoadbalancersID.ValueString())
	}
	if !plan.SecurityGroupID.IsNull() {
		plan.SecurityGroupID = types.StringValue(plan.SecurityGroupID.ValueString())
	}
	if !plan.TargetGroupsID.IsNull() {
		plan.TargetGroupsID = types.StringValue(plan.TargetGroupsID.ValueString())
	}

	publicIPEnabled, err := strconv.ParseBool(getAutoScaling.PublicIPEnabled)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing PublicIPEnabled",
			fmt.Sprintf("Could not parse PublicIPEnabled value '%s' for autoscaling %s: %s", getAutoScaling.PublicIPEnabled, getAutoScaling.ID, err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(autoscaling.ID))
	plan.Dcslug = types.StringValue(plan.Dcslug.ValueString())
	plan.Userid = types.StringValue(getAutoScaling.Userid)
	plan.Name = types.StringValue(getAutoScaling.Name)
	plan.Minsize = types.StringValue(getAutoScaling.Minsize)
	plan.Maxsize = types.StringValue(getAutoScaling.Maxsize)
	plan.Desiredsize = types.StringValue(getAutoScaling.Desiredsize)
	plan.Planid = types.StringValue(getAutoScaling.Planid)
	plan.Planname = types.StringValue(getAutoScaling.Planname)
	plan.InstanceTemplateid = types.StringValue(getAutoScaling.InstanceTemplateid)
	plan.Image = types.StringValue(getAutoScaling.Image)
	plan.ImageName = types.StringValue(getAutoScaling.ImageName)
	plan.Snapshotid = types.StringValue(getAutoScaling.Snapshotid)
	plan.Status = types.StringValue(getAutoScaling.Status)
	plan.CreatedAt = types.StringValue(getAutoScaling.CreatedAt)
	plan.SuspendedAt = types.StringValue(getAutoScaling.SuspendedAt)
	plan.StoppedAt = types.StringValue(getAutoScaling.StoppedAt)
	plan.StartedAt = types.StringValue(getAutoScaling.StartedAt)
	plan.DeletedAt = types.StringValue(getAutoScaling.DeletedAt)
	plan.PublicIPEnabled = types.BoolValue(publicIPEnabled)
	plan.CooldownTill = types.StringValue(getAutoScaling.CooldownTill)
	plan.Backupid = types.StringValue(getAutoScaling.Backupid)
	// TODO
	plan.Stackid = types.StringValue(plan.Stackid.ValueString())
	plan.Stackimage = types.StringValue(getAutoScaling.Image)
	plan.OsDiskSize = types.Int64Value(int64(osDiskSize * 10))
	plan.Policies = policiesModel
	plan.VpcID = types.StringValue(getAutoScaling.Vpc[0].ID)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// dclocatio
	dclocationModel := DclocationModel{
		Location: types.StringValue(getAutoScaling.Dclocation.Location),
		Country:  types.StringValue(getAutoScaling.Dclocation.Country),
		Dc:       types.StringValue(getAutoScaling.Dclocation.DC),
		Dccc:     types.StringValue(getAutoScaling.Dclocation.Dccc),
	}
	diags = resp.State.SetAttribute(ctx, path.Root("dclocation"), dclocationModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Auto Scaling Plan
	autoScalingPlanModel := AutoScalingPlanModel{
		Planid:         types.StringValue(getAutoScaling.Plan.Planid),
		RAM:            types.StringValue(getAutoScaling.Plan.RAM),
		CPU:            types.StringValue(getAutoScaling.Plan.CPU),
		Disk:           types.StringValue(getAutoScaling.Plan.Disk),
		Bandwidth:      types.StringValue(getAutoScaling.Plan.Bandwidth),
		DedicatedVcore: types.StringValue(getAutoScaling.Plan.DedicatedVcore),
	}
	diags = resp.State.SetAttribute(ctx, path.Root("plan"), autoScalingPlanModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// vpc
	vpcModel := []AutoScalingVpcModel{}
	for _, v := range getAutoScaling.Vpc {
		vpc := AutoScalingVpcModel{
			Total:     types.Int64Value(int64(v.Total)),
			Available: types.Int64Value(int64(v.Available)),
			Network:   types.StringValue(v.Network),
			Name:      types.StringValue(v.Name),
			Size:      types.StringValue(v.Size),
			Dcslug:    types.StringValue(v.Dcslug),
		}
		vpcModel = append(vpcModel, vpc)
	}
	diags = resp.State.SetAttribute(ctx, path.Root("vpc"), vpcModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Load balancers
	lbModel := []AutoScalingLoadbalancersModel{}
	for _, v := range getAutoScaling.Loadbalancers {
		lb := AutoScalingLoadbalancersModel{
			ID:   types.StringValue(v.ID),
			Name: types.StringValue(v.Name),
			IP:   types.StringValue(v.IP),
		}
		lbModel = append(lbModel, lb)
	}
	diags = resp.State.SetAttribute(ctx, path.Root("load_balancers"), lbModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Target Groups
	targetGroupsModek := []AutoScalingTargetGroupModel{}
	for _, v := range getAutoScaling.TargetGroups {
		targetGroup := AutoScalingTargetGroupModel{
			ID:       types.StringValue(v.ID),
			Name:     types.StringValue(v.Name),
			Protocol: types.StringValue(v.Protocol),
			Port:     types.StringValue(v.Port),
		}
		targetGroupsModek = append(targetGroupsModek, targetGroup)
	}
	diags = resp.State.SetAttribute(ctx, path.Root("target_groups"), targetGroupsModek)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Security Group
	securityGroupModel := []SecurityGroupModel{}
	for _, v := range getAutoScaling.SecurityGroups {
		securityGroup := SecurityGroupModel{
			ID:   types.StringValue(v.ID),
			Name: types.StringValue(v.Name),
		}
		securityGroupModel = append(securityGroupModel, securityGroup)
	}
	diags = resp.State.SetAttribute(ctx, path.Root("security_groups"), securityGroupModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instancesModel := []InstancesModel{}
	for _, v := range getAutoScaling.Instances {
		Instances := InstancesModel{
			Cloudid:   types.StringValue(v.ID),
			Hostname:  types.StringValue(v.Hostname),
			CreatedAt: types.StringValue(v.CreatedAt),
			IP:        types.StringValue(v.IP),
			Status:    types.StringValue(v.Status),
		}
		instancesModel = append(instancesModel, Instances)
	}
	diags = resp.State.SetAttribute(ctx, path.Root("instances"), instancesModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Schedule
	scheduleModel := []ScheduleModel{}
	for _, v := range getAutoScaling.Schedules {
		Schedule := ScheduleModel{
			ID:          types.StringValue(v.ID),
			Groupid:     types.StringValue(v.Groupid),
			Name:        types.StringValue(v.Name),
			Desiredsize: types.StringValue(v.Desiredsize),
			Recurrence:  types.StringValue(v.Recurrence),
			StartDate:   types.StringValue(v.StartDate),
			Status:      types.StringValue(v.Status),
			Timezone:    types.StringValue(v.Timezone),
		}
		scheduleModel = append(scheduleModel, Schedule)
	}
	diags = resp.State.SetAttribute(ctx, path.Root("schedules"), scheduleModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "finish create autoscaling")
}

// Read resource information.
func (s *AutoScalingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "read autoscaling")

	// Get current state
	var state AutoScalingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "send get autoscaling request")
	// Get refreshed autoscaling value from utho
	getAutoScaling, err := s.client.AutoScaling().Read(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading utho autoscaling",
			"Could not read utho autoscaling "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	osDiskSize, _ := strconv.Atoi(getAutoScaling.Plan.Disk)

	policiesModel := []PolicyModel{}
	for _, v := range getAutoScaling.Policies {
		policy := PolicyModel{
			ID:                 types.StringValue(v.ID),
			Userid:             types.StringValue(v.Userid),
			Product:            types.StringValue(v.Product),
			Productid:          types.StringValue(v.Productid),
			Groupid:            types.StringValue(v.Groupid),
			Name:               types.StringValue(v.Name),
			Type:               types.StringValue(v.Type),
			Adjust:             types.StringValue(v.Adjust),
			Period:             types.StringValue(v.Period),
			Cooldown:           types.StringValue(v.Cooldown),
			CooldownTill:       types.StringValue(v.CooldownTill),
			Compare:            types.StringValue(v.Compare),
			Value:              types.StringValue(v.Value),
			AlertID:            types.StringValue(v.AlertID),
			Status:             types.StringValue(v.Status),
			KubernetesID:       types.StringValue(v.KubernetesID),
			KubernetesNodepool: types.StringValue(v.KubernetesNodepool),
			Cloudid:            types.StringValue(v.Cloudid),
			Maxsize:            types.StringValue(v.Maxsize),
			Minsize:            types.StringValue(v.Minsize),
		}
		policiesModel = append(policiesModel, policy)
	}

	if !state.LoadbalancersID.IsNull() {
		state.LoadbalancersID = types.StringValue(state.LoadbalancersID.ValueString())
	}
	if !state.SecurityGroupID.IsNull() {
		state.SecurityGroupID = types.StringValue(state.SecurityGroupID.ValueString())
	}
	if !state.TargetGroupsID.IsNull() {
		state.TargetGroupsID = types.StringValue(state.TargetGroupsID.ValueString())
	}

	publicIPEnabled, err := strconv.ParseBool(getAutoScaling.PublicIPEnabled)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing PublicIPEnabled",
			fmt.Sprintf("Could not parse PublicIPEnabled value '%s' for autoscaling %s: %s", getAutoScaling.PublicIPEnabled, getAutoScaling.ID, err.Error()),
		)
		return
	}
	state.ID = types.StringValue(getAutoScaling.ID)
	state.Dcslug = types.StringValue(getAutoScaling.Dcslug)
	state.Userid = types.StringValue(getAutoScaling.Userid)
	state.Name = types.StringValue(getAutoScaling.Name)
	state.Minsize = types.StringValue(getAutoScaling.Minsize)
	state.Maxsize = types.StringValue(getAutoScaling.Maxsize)
	state.Desiredsize = types.StringValue(getAutoScaling.Desiredsize)
	state.Planid = types.StringValue(getAutoScaling.Planid)
	state.Planname = types.StringValue(getAutoScaling.Planname)
	state.InstanceTemplateid = types.StringValue(getAutoScaling.InstanceTemplateid)
	state.Image = types.StringValue(getAutoScaling.Image)
	state.ImageName = types.StringValue(getAutoScaling.ImageName)
	state.Snapshotid = types.StringValue(getAutoScaling.Snapshotid)
	state.Status = types.StringValue(getAutoScaling.Status)
	state.CreatedAt = types.StringValue(getAutoScaling.CreatedAt)
	state.SuspendedAt = types.StringValue(getAutoScaling.SuspendedAt)
	state.StoppedAt = types.StringValue(getAutoScaling.StoppedAt)
	state.StartedAt = types.StringValue(getAutoScaling.StartedAt)
	state.DeletedAt = types.StringValue(getAutoScaling.DeletedAt)
	state.PublicIPEnabled = types.BoolValue(publicIPEnabled)
	state.CooldownTill = types.StringValue(getAutoScaling.CooldownTill)
	state.Backupid = types.StringValue(getAutoScaling.Backupid)
	state.Stackid = types.StringValue(getAutoScaling.Stack)
	state.Stackimage = types.StringValue(getAutoScaling.Image)
	state.OsDiskSize = types.Int64Value(int64(osDiskSize * 10))
	state.Policies = policiesModel
	state.VpcID = types.StringValue(getAutoScaling.Vpc[0].ID)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// dclocatio
	dclocationModel := DclocationModel{
		Location: types.StringValue(getAutoScaling.Dclocation.Location),
		Country:  types.StringValue(getAutoScaling.Dclocation.Country),
		Dc:       types.StringValue(getAutoScaling.Dclocation.DC),
		Dccc:     types.StringValue(getAutoScaling.Dclocation.Dccc),
	}
	diags = resp.State.SetAttribute(ctx, path.Root("dclocation"), dclocationModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Auto Scaling Plan
	autoScalingPlanModel := AutoScalingPlanModel{
		Planid:         types.StringValue(getAutoScaling.Plan.Planid),
		RAM:            types.StringValue(getAutoScaling.Plan.RAM),
		CPU:            types.StringValue(getAutoScaling.Plan.CPU),
		Disk:           types.StringValue(getAutoScaling.Plan.Disk),
		Bandwidth:      types.StringValue(getAutoScaling.Plan.Bandwidth),
		DedicatedVcore: types.StringValue(getAutoScaling.Plan.DedicatedVcore),
	}
	diags = resp.State.SetAttribute(ctx, path.Root("plan"), autoScalingPlanModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// vpc
	vpcModel := []AutoScalingVpcModel{}
	for _, v := range getAutoScaling.Vpc {
		vpc := AutoScalingVpcModel{
			Total:     types.Int64Value(int64(v.Total)),
			Available: types.Int64Value(int64(v.Available)),
			Network:   types.StringValue(v.Network),
			Name:      types.StringValue(v.Name),
			Size:      types.StringValue(v.Size),
			Dcslug:    types.StringValue(v.Dcslug),
		}
		vpcModel = append(vpcModel, vpc)
	}
	diags = resp.State.SetAttribute(ctx, path.Root("vpc"), vpcModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Load balancers
	lbModel := []AutoScalingLoadbalancersModel{}
	for _, v := range getAutoScaling.Loadbalancers {
		lb := AutoScalingLoadbalancersModel{
			ID:   types.StringValue(v.ID),
			Name: types.StringValue(v.Name),
			IP:   types.StringValue(v.IP),
		}
		lbModel = append(lbModel, lb)
	}
	diags = resp.State.SetAttribute(ctx, path.Root("load_balancers"), lbModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Target Groups
	targetGroupsModek := []AutoScalingTargetGroupModel{}
	for _, v := range getAutoScaling.TargetGroups {
		targetGroup := AutoScalingTargetGroupModel{
			ID:       types.StringValue(v.ID),
			Name:     types.StringValue(v.Name),
			Protocol: types.StringValue(v.Protocol),
			Port:     types.StringValue(v.Port),
		}
		targetGroupsModek = append(targetGroupsModek, targetGroup)
	}
	diags = resp.State.SetAttribute(ctx, path.Root("target_groups"), targetGroupsModek)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Security Group
	securityGroupModel := []SecurityGroupModel{}
	for _, v := range getAutoScaling.SecurityGroups {
		securityGroup := SecurityGroupModel{
			ID:   types.StringValue(v.ID),
			Name: types.StringValue(v.Name),
		}
		securityGroupModel = append(securityGroupModel, securityGroup)
	}
	diags = resp.State.SetAttribute(ctx, path.Root("security_groups"), securityGroupModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Instances
	instancesModel := []InstancesModel{}
	for _, v := range getAutoScaling.Instances {
		Instances := InstancesModel{
			Cloudid:   types.StringValue(v.ID),
			Hostname:  types.StringValue(v.Hostname),
			CreatedAt: types.StringValue(v.CreatedAt),
			IP:        types.StringValue(v.IP),
			Status:    types.StringValue(v.Status),
		}
		instancesModel = append(instancesModel, Instances)
	}
	diags = resp.State.SetAttribute(ctx, path.Root("instances"), instancesModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Schedule
	scheduleModel := []ScheduleModel{}
	for _, v := range getAutoScaling.Schedules {
		Schedule := ScheduleModel{
			ID:          types.StringValue(v.ID),
			Groupid:     types.StringValue(v.Groupid),
			Name:        types.StringValue(v.Name),
			Desiredsize: types.StringValue(v.Desiredsize),
			Recurrence:  types.StringValue(v.Recurrence),
			StartDate:   types.StringValue(v.StartDate),
			Status:      types.StringValue(v.Status),
			Timezone:    types.StringValue(v.Timezone),
		}
		scheduleModel = append(scheduleModel, Schedule)
	}
	diags = resp.State.SetAttribute(ctx, path.Root("schedules"), scheduleModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "finish get autoscaling request")
}

func (s *AutoScalingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AutoScalingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state AutoScalingResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := utho.UpdateAutoScalingParams{
		AutoScalingId: state.ID.ValueString(),
		Name:          plan.Name.ValueString(),
		Minsize:       plan.Minsize.ValueString(),
		Maxsize:       plan.Maxsize.ValueString(),
		Desiredsize:   plan.Desiredsize.ValueString(),
	}

	tflog.Debug(ctx, "send update autoscaling request")
	// Get refreshed autoscaling value from utho
	_, err := s.client.AutoScaling().Update(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading utho autoscaling",
			"Could not update utho autoscaling "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "send get autoscaling request")
	// Get refreshed autoscaling value from utho
	getAutoScaling, err := s.client.AutoScaling().Read(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading utho autoscaling",
			"Could not update utho autoscaling "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
	state.Name = types.StringValue(getAutoScaling.Name)
	state.Minsize = types.StringValue(getAutoScaling.Minsize)
	state.Maxsize = types.StringValue(getAutoScaling.Maxsize)
	state.Desiredsize = types.StringValue(getAutoScaling.Desiredsize)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish get autoscaling request")
}

// Delete deletes the resource and removes the Terraform state on success.
func (s *AutoScalingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "delete autoscaling")
	// Get current state
	var state AutoScalingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "send delete autoscaling request")
	// delete autoscaling
	_, err := s.client.AutoScaling().Delete(state.ID.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing utho autoscaling",
			"Could not delete utho autoscaling "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}
