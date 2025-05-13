package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/uthoplatforms/utho-go/utho"
)

// implement resource interfaces.
var (
	_ resource.Resource                = &TargetGroupResource{}
	_ resource.ResourceWithConfigure   = &TargetGroupResource{}
	_ resource.ResourceWithImportState = &TargetGroupResource{}
)

// NewTargetGroupResource is a helper function to simplify the provider implementation.
func NewTargetGroupResource() resource.Resource {
	return &TargetGroupResource{}
}

// TargetGroupResource is the resource implementation.
type TargetGroupResource struct {
	client utho.Client
}

type TargetGroupResourceModel struct {
	ID                  types.String          `tfsdk:"id"`
	Name                types.String          `tfsdk:"name"`
	Port                types.String          `tfsdk:"port"`
	Protocol            types.String          `tfsdk:"protocol"`
	HealthCheckPath     types.String          `tfsdk:"health_check_path"`
	HealthCheckInterval types.String          `tfsdk:"health_check_interval"`
	HealthCheckProtocol types.String          `tfsdk:"health_check_protocol"`
	HealthCheckTimeout  types.String          `tfsdk:"health_check_timeout"`
	HealthyThreshold    types.String          `tfsdk:"healthy_threshold"`
	UnhealthyThreshold  types.String          `tfsdk:"unhealthy_threshold"`
	CreatedAt           types.String          `tfsdk:"created_at"`
	UpdatedAt           types.String          `tfsdk:"updated_at"`
	Targets             []TargetResourceModel `tfsdk:"targets"`
}
type TargetResourceModel struct {
	Lbid                types.String `tfsdk:"lbid"`
	IP                  types.String `tfsdk:"ip"`
	Cloudid             types.String `tfsdk:"cloudid"`
	Status              types.String `tfsdk:"status"`
	ScalingGroupid      types.String `tfsdk:"scaling_groupid"`
	KubernetesClusterid types.String `tfsdk:"kubernetes_clusterid"`
	BackendPort         types.String `tfsdk:"backend_port"`
	BackendProtocol     types.String `tfsdk:"backend_protocol"`
	TargetgroupID       types.String `tfsdk:"targetgroup_id"`
	FrontendID          types.String `tfsdk:"frontend_id"`
	ID                  types.String `tfsdk:"id"`
}

// Metadata returns the resource type name.
func (s *TargetGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_target_group"
}

// Configure adds the provider configured client to the data source.
func (d *TargetGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(utho.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected TargetGroup Data Source Configure Type",
			fmt.Sprintf("Expected utho.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

// Schema defines the schema for the resource.
func (s *TargetGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                    schema.StringAttribute{Computed: true, Description: "id"},
			"name":                  schema.StringAttribute{Required: true, Description: "Provide Target Group name eg: my_group", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"protocol":              schema.StringAttribute{Required: true, Description: "Provide protocol eg: HTTP, HTTPS, TCP, UDP", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"port":                  schema.StringAttribute{Required: true, Description: "Provide the port according to protocol eg: 80", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"health_check_path":     schema.StringAttribute{Required: true, Description: "Provide health check path for the target group", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"health_check_protocol": schema.StringAttribute{Required: true, Description: "Provide health check protocol for the target group", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"health_check_timeout":  schema.StringAttribute{Required: true, Description: "Provide health check timeout for the target group", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"unhealthy_threshold":   schema.StringAttribute{Required: true, Description: "Provide unhealthy threshold for the target group", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"health_check_interval": schema.StringAttribute{Required: true, Description: "Provide health check interval for the target group", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"healthy_threshold":     schema.StringAttribute{Required: true, Description: "Provide healthy threshold for the target group", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"created_at":            schema.StringAttribute{Computed: true, Description: "created at"},
			"updated_at":            schema.StringAttribute{Computed: true, Description: "updated at"},
			"targets": schema.ListNestedAttribute{
				Optional:    true,
				Description: "targets",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                   schema.StringAttribute{Computed: true, Description: "Id"},
						"ip":                   schema.StringAttribute{Required: true, Description: "Target Ip", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
						"backend_port":         schema.StringAttribute{Required: true, Description: "Backend Port", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
						"backend_protocol":     schema.StringAttribute{Required: true, Description: "Backend Protocol", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
						"lbid":                 schema.StringAttribute{Computed: true, Description: "Lbid"},
						"cloudid":              schema.StringAttribute{Computed: true, Description: "Cloudid"},
						"status":               schema.StringAttribute{Computed: true, Description: "Status"},
						"scaling_groupid":      schema.StringAttribute{Computed: true, Description: "Scaling Groupid"},
						"kubernetes_clusterid": schema.StringAttribute{Computed: true, Description: "Kubernetes Clusterid"},
						"targetgroup_id":       schema.StringAttribute{Computed: true, Description: "Targetgroup Id"},
						"frontend_id":          schema.StringAttribute{Computed: true, Description: "Frontend Id"},
					},
				},
			},
		},
	}
}

// Import using target group as the attribute
func (s *TargetGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource.
func (s *TargetGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "create target group")
	// Retrieve values from plan
	var plan TargetGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	targetGroupRequest := utho.CreateTargetGroupParams{
		Name:                plan.Name.ValueString(),
		Protocol:            plan.Protocol.ValueString(),
		Port:                plan.Port.ValueString(),
		HealthCheckPath:     plan.HealthCheckPath.ValueString(),
		HealthCheckProtocol: plan.HealthCheckProtocol.ValueString(),
		HealthCheckInterval: plan.HealthCheckInterval.ValueString(),
		HealthCheckTimeout:  plan.HealthCheckTimeout.ValueString(),
		HealthyThreshold:    plan.HealthyThreshold.ValueString(),
		UnhealthyThreshold:  plan.UnhealthyThreshold.ValueString(),
	}
	tflog.Debug(ctx, "send create target group request")
	createTargetGroupResponse, err := s.client.TargetGroup().Create(targetGroupRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating target group",
			"Could not create target group, unexpected error: "+err.Error(),
		)
		return
	}
	targetGroupId := strconv.Itoa(createTargetGroupResponse.ID)

	var targets []string
	for _, target := range plan.Targets {
		createTargetGroupTargetParams := utho.CreateTargetGroupTargetParams{
			TargetGroupId:   targetGroupId,
			IP:              target.IP.ValueString(),
			BackendPort:     target.BackendPort.ValueString(),
			BackendProtocol: target.BackendProtocol.ValueString(),
			Cloudid:         target.Cloudid.ValueString(),
		}
		createTargetsRes, err := s.client.TargetGroup().CreateTarget(createTargetGroupTargetParams)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating target group",
				"Could not create target group target, unexpected error: "+err.Error(),
			)
			return
		}
		targets = append(targets, createTargetsRes.ID)
	}

	targetGroup, err := s.client.TargetGroup().Read(targetGroupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading utho target group",
			"Could not read utho target group "+targetGroupId+": "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	var targetsResourceModel []TargetResourceModel
	for _, target := range targetGroup.Targets {
		targetResourceModel := TargetResourceModel{
			Lbid:                types.StringValue(target.Lbid),
			IP:                  types.StringValue(target.IP),
			Cloudid:             types.StringValue(target.Cloudid),
			Status:              types.StringValue(target.Status),
			ScalingGroupid:      types.StringValue(target.ScalingGroupid),
			KubernetesClusterid: types.StringValue(target.KubernetesClusterid),
			BackendPort:         types.StringValue(target.BackendPort),
			BackendProtocol:     types.StringValue(target.BackendProtocol),
			TargetgroupID:       types.StringValue(target.TargetgroupID),
			FrontendID:          types.StringValue(target.FrontendID),
			ID:                  types.StringValue(target.ID),
		}
		targetsResourceModel = append(targetsResourceModel, targetResourceModel)
	}

	plan = TargetGroupResourceModel{
		ID:                  types.StringValue(targetGroupId),
		Name:                types.StringValue(targetGroup.Name),
		Port:                types.StringValue(targetGroup.Port),
		Protocol:            types.StringValue(targetGroup.Protocol),
		HealthCheckPath:     types.StringValue(targetGroup.HealthCheckPath),
		HealthCheckInterval: types.StringValue(targetGroup.HealthCheckInterval),
		HealthCheckProtocol: types.StringValue(targetGroup.HealthCheckProtocol),
		HealthCheckTimeout:  types.StringValue(targetGroup.HealthCheckTimeout),
		HealthyThreshold:    types.StringValue(targetGroup.HealthyThreshold),
		UnhealthyThreshold:  types.StringValue(targetGroup.UnhealthyThreshold),
		CreatedAt:           types.StringValue(targetGroup.CreatedAt),
		UpdatedAt:           types.StringValue(targetGroup.UpdatedAt),
		Targets:             targetsResourceModel,
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish create target group")
}

// Read resource information.
func (s *TargetGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "read target group")

	// Get current state
	var state TargetGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "send get target group request")
	// Get refreshed target group value from utho
	targetGroup, err := s.client.TargetGroup().Read(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading utho target group",
			"Could not read utho target group "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	var targetsResourceModel []TargetResourceModel
	for _, target := range targetGroup.Targets {
		targetResourceModel := TargetResourceModel{
			Lbid:                types.StringValue(target.Lbid),
			IP:                  types.StringValue(target.IP),
			Cloudid:             types.StringValue(target.Cloudid),
			Status:              types.StringValue(target.Status),
			ScalingGroupid:      types.StringValue(target.ScalingGroupid),
			KubernetesClusterid: types.StringValue(target.KubernetesClusterid),
			BackendPort:         types.StringValue(target.BackendPort),
			BackendProtocol:     types.StringValue(target.BackendProtocol),
			TargetgroupID:       types.StringValue(target.TargetgroupID),
			FrontendID:          types.StringValue(target.FrontendID),
			ID:                  types.StringValue(target.ID),
		}
		targetsResourceModel = append(targetsResourceModel, targetResourceModel)
	}

	state = TargetGroupResourceModel{
		ID:                  types.StringValue(targetGroup.ID),
		Name:                types.StringValue(targetGroup.Name),
		Port:                types.StringValue(targetGroup.Port),
		Protocol:            types.StringValue(targetGroup.Protocol),
		HealthCheckPath:     types.StringValue(targetGroup.HealthCheckPath),
		HealthCheckInterval: types.StringValue(targetGroup.HealthCheckInterval),
		HealthCheckProtocol: types.StringValue(targetGroup.HealthCheckProtocol),
		HealthCheckTimeout:  types.StringValue(targetGroup.HealthCheckTimeout),
		HealthyThreshold:    types.StringValue(targetGroup.HealthyThreshold),
		UnhealthyThreshold:  types.StringValue(targetGroup.UnhealthyThreshold),
		CreatedAt:           types.StringValue(targetGroup.CreatedAt),
		UpdatedAt:           types.StringValue(targetGroup.UpdatedAt),
		Targets:             targetsResourceModel,
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish get target group request")
}

func (s *TargetGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// updating resource is not supported
}

// Delete deletes the resource and removes the Terraform state on success.
func (s *TargetGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "delete target group")
	// Get current state
	var state TargetGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "send delete target group request")
	// delete target group
	_, err := s.client.TargetGroup().Delete(state.ID.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing utho target group",
			"Could not delete utho target group "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}
