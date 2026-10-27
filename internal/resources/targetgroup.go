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
// utho_target_group
// ══════════════════════════════════════════════════════════════

type TargetGroupResource struct{ client *client.Client }

type TargetGroupModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Protocol            types.String `tfsdk:"protocol"`
	Port                types.String `tfsdk:"port"`
	HealthCheckPath     types.String `tfsdk:"health_check_path"`
	HealthCheckProtocol types.String `tfsdk:"health_check_protocol"`
	HealthCheckInterval types.String `tfsdk:"health_check_interval"`
	HealthCheckTimeout  types.String `tfsdk:"health_check_timeout"`
	HealthyThreshold    types.String `tfsdk:"healthy_threshold"`
	UnhealthyThreshold  types.String `tfsdk:"unhealthy_threshold"`
	CreatedAt           types.String `tfsdk:"created_at"`
}

func NewTargetGroupResource() resource.Resource { return &TargetGroupResource{} }

func (r *TargetGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_target_group"
}

func (r *TargetGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage Utho Target Groups. A target group defines a pool of backend servers and health check configuration. Attach it to a Load Balancer frontend or Auto Scaling group to route traffic.",
		Attributes: map[string]schema.Attribute{
			"id":                    schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"created_at":            schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":                  schema.StringAttribute{Required: true, Description: "Target group name. Must be unique in your account."},
			"protocol":              schema.StringAttribute{Required: true, Description: "Traffic protocol: HTTP, HTTPS, TCP, UDP."},
			"port":                  schema.StringAttribute{Required: true, Description: "Port on which registered targets receive traffic (e.g. 80, 443)."},
			"health_check_path":     schema.StringAttribute{Required: true, Description: "Health check endpoint path (e.g. /health, /status)."},
			"health_check_protocol": schema.StringAttribute{Required: true, Description: "Health check protocol: HTTP, HTTPS, TCP, UDP."},
			"health_check_interval": schema.StringAttribute{Required: true, Description: "Health check interval in seconds (e.g. 300)."},
			"health_check_timeout":  schema.StringAttribute{Required: true, Description: "Health check timeout in seconds (e.g. 300)."},
			"healthy_threshold":     schema.StringAttribute{Required: true, Description: "Consecutive successful checks to mark a target healthy (e.g. 3)."},
			"unhealthy_threshold":   schema.StringAttribute{Required: true, Description: "Consecutive failed checks to mark a target unhealthy (e.g. 3)."},
		},
	}
}

func (r *TargetGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TargetGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TargetGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateTargetGroup(&client.TargetGroupCreateRequest{
		Name:                plan.Name.ValueString(),
		Protocol:            plan.Protocol.ValueString(),
		Port:                plan.Port.ValueString(),
		HealthCheckPath:     plan.HealthCheckPath.ValueString(),
		HealthCheckProtocol: plan.HealthCheckProtocol.ValueString(),
		HealthCheckInterval: plan.HealthCheckInterval.ValueString(),
		HealthCheckTimeout:  plan.HealthCheckTimeout.ValueString(),
		HealthyThreshold:    plan.HealthyThreshold.ValueString(),
		UnhealthyThreshold:  plan.UnhealthyThreshold.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating target group", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	plan.CreatedAt = types.StringValue("")

	if tg, err2 := r.client.GetTargetGroup(id); err2 == nil && tg != nil {
		plan.CreatedAt = types.StringValue(tg.CreatedAt)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *TargetGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TargetGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tg, err := r.client.GetTargetGroup(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading target group", fmt.Sprintf("%s", err))
		return
	}
	if tg == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(tg.Name)
	state.Protocol = types.StringValue(tg.Protocol)
	state.Port = types.StringValue(tg.Port)
	state.HealthCheckPath = types.StringValue(tg.HealthCheckPath)
	state.HealthCheckProtocol = types.StringValue(tg.HealthCheckProtocol)
	state.HealthCheckInterval = types.StringValue(tg.HealthCheckInterval)
	state.HealthCheckTimeout = types.StringValue(tg.HealthCheckTimeout)
	state.HealthyThreshold = types.StringValue(tg.HealthyThreshold)
	state.UnhealthyThreshold = types.StringValue(tg.UnhealthyThreshold)
	state.CreatedAt = types.StringValue(tg.CreatedAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *TargetGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state TargetGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateTargetGroup(state.ID.ValueString(), &client.TargetGroupUpdateRequest{
		Name:                plan.Name.ValueString(),
		Protocol:            plan.Protocol.ValueString(),
		Port:                plan.Port.ValueString(),
		HealthCheckPath:     plan.HealthCheckPath.ValueString(),
		HealthCheckProtocol: plan.HealthCheckProtocol.ValueString(),
		HealthCheckInterval: plan.HealthCheckInterval.ValueString(),
		HealthCheckTimeout:  plan.HealthCheckTimeout.ValueString(),
		HealthyThreshold:    plan.HealthyThreshold.ValueString(),
		UnhealthyThreshold:  plan.UnhealthyThreshold.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error updating target group", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = state.ID
	plan.CreatedAt = state.CreatedAt
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *TargetGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TargetGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteTargetGroup(state.ID.ValueString(), state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting target group", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// utho_target_group_target
// ══════════════════════════════════════════════════════════════

type TargetGroupTargetResource struct{ client *client.Client }

type TargetGroupTargetModel struct {
	ID                  types.String `tfsdk:"id"`
	TargetGroupID       types.String `tfsdk:"target_group_id"`
	Type                types.String `tfsdk:"type"`
	BackendPort         types.String `tfsdk:"backend_port"`
	CloudID             types.String `tfsdk:"cloud_id"`
	IP                  types.String `tfsdk:"ip"`
	KubernetesClusterID types.String `tfsdk:"kubernetes_cluster_id"`
	PoolName            types.String `tfsdk:"pool_name"`
}

func NewTargetGroupTargetResource() resource.Resource { return &TargetGroupTargetResource{} }

func (r *TargetGroupTargetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_target_group_target"
}

func (r *TargetGroupTargetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Register a backend server as a target in a Utho Target Group. Supports cloud instances, Kubernetes node pools, and custom IPs.",
		Attributes: map[string]schema.Attribute{
			"id":                    schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"target_group_id":       schema.StringAttribute{Required: true, Description: "Target group ID. Changing forces new resource.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"type":                  schema.StringAttribute{Required: true, Description: "Target type: cloud, kubernetes, custom. Changing forces new resource.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"backend_port":          schema.StringAttribute{Required: true, Description: "Port on which the target receives traffic. Changing forces new resource.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"cloud_id":              schema.StringAttribute{Optional: true, Description: "Cloud instance ID. Required when type is cloud. Changing forces new resource.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"ip":                    schema.StringAttribute{Optional: true, Computed: true, Description: "IP address. Required when type is custom. For cloud targets, provide the instance public IP. Changing forces new resource.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"kubernetes_cluster_id": schema.StringAttribute{Optional: true, Description: "Kubernetes cluster ID. Required when type is kubernetes. Changing forces new resource.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"pool_name":             schema.StringAttribute{Optional: true, Description: "Kubernetes node pool name. Required when type is kubernetes. Changing forces new resource.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *TargetGroupTargetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TargetGroupTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TargetGroupTargetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &client.TargetGroupAddTargetRequest{
		BackendPort:         plan.BackendPort.ValueString(),
		Type:                plan.Type.ValueString(),
		CloudID:             plan.CloudID.ValueString(),
		IP:                  plan.IP.ValueString(),
		KubernetesClusterID: plan.KubernetesClusterID.ValueString(),
		PoolName:            plan.PoolName.ValueString(),
	}

	if err := r.client.AddTargetGroupTarget(plan.TargetGroupID.ValueString(), addReq); err != nil {
		resp.Diagnostics.AddError("Error adding target", fmt.Sprintf("%s", err))
		return
	}

	// Find the target ID from the TG targets list
	targetID := ""
	if tg, err2 := r.client.GetTargetGroup(plan.TargetGroupID.ValueString()); err2 == nil && tg != nil {
		for _, t := range tg.Targets {
			if t.BackendPort == plan.BackendPort.ValueString() {
				if plan.Type.ValueString() == "cloud" && t.CloudID == plan.CloudID.ValueString() {
					targetID = t.ID
					if plan.IP.IsNull() || plan.IP.IsUnknown() {
						plan.IP = types.StringValue(t.IP)
					}
					break
				}
				if plan.Type.ValueString() == "custom" && t.IP == plan.IP.ValueString() {
					targetID = t.ID
					break
				}
				if plan.Type.ValueString() == "kubernetes" && t.KubernetesClusterID == plan.KubernetesClusterID.ValueString() {
					targetID = t.ID
					break
				}
			}
		}
	}

	plan.ID = types.StringValue(targetID)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *TargetGroupTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TargetGroupTargetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tg, err := r.client.GetTargetGroup(state.TargetGroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading target group", fmt.Sprintf("%s", err))
		return
	}
	if tg == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	found := false
	for _, t := range tg.Targets {
		if t.ID == state.ID.ValueString() {
			found = true
			break
		}
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *TargetGroupTargetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Target group targets cannot be updated. Destroy and recreate.")
}

func (r *TargetGroupTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TargetGroupTargetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteTargetGroupTarget(state.TargetGroupID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting target", fmt.Sprintf("%s", err))
	}
}
