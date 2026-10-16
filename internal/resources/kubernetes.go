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
// KUBERNETES CLUSTER — utho_kubernetes
// ══════════════════════════════════════════════════════════════

type KubernetesResource struct{ client *client.Client }

type K8sNodePoolModel struct {
	Label     types.String `tfsdk:"label"`
	Size      types.String `tfsdk:"size"`
	NodeCount types.String `tfsdk:"node_count"`
	MinNodes  types.String `tfsdk:"min_nodes"`
	MaxNodes  types.String `tfsdk:"max_nodes"`
	DiskSize  types.String `tfsdk:"disk_size"`
	DiskType  types.String `tfsdk:"disk_type"`
}

type KubernetesModel struct {
	ID             types.String       `tfsdk:"id"`
	DCSlug         types.String       `tfsdk:"dcslug"`
	ClusterLabel   types.String       `tfsdk:"cluster_label"`
	ClusterVersion types.String       `tfsdk:"cluster_version"`
	NetworkType    types.String       `tfsdk:"network_type"`
	VPC            types.String       `tfsdk:"vpc"`
	CPUModel       types.String       `tfsdk:"cpumodel"`
	Status         types.String       `tfsdk:"status"`
	IP             types.String       `tfsdk:"ip"`
	DNS            types.String       `tfsdk:"dns"`
	CreatedAt      types.String       `tfsdk:"created_at"`
	NodePools      []K8sNodePoolModel `tfsdk:"nodepools"`
}

func NewKubernetesResource() resource.Resource { return &KubernetesResource{} }

func (r *KubernetesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes"
}

func (r *KubernetesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage Utho Kubernetes clusters.",
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"status":          schema.StringAttribute{Computed: true, Description: "Cluster status."},
			"ip":              schema.StringAttribute{Computed: true, Description: "Cluster control plane IP."},
			"dns":             schema.StringAttribute{Computed: true, Description: "Cluster DNS endpoint."},
			"created_at":      schema.StringAttribute{Computed: true, Description: "Timestamp when the cluster was created."},
			"dcslug":          schema.StringAttribute{Required: true, Description: "Data center slug.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"cluster_label":   schema.StringAttribute{Required: true, Description: "Cluster name/label.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"cluster_version": schema.StringAttribute{Required: true, Description: "Kubernetes version (e.g. 1.30.0-utho).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"network_type":    schema.StringAttribute{Required: true, Description: "Network type: public, private, or publicprivate.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"vpc":             schema.StringAttribute{Optional: true, Description: "VPC subnet ID.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"cpumodel":        schema.StringAttribute{Optional: true, Description: "CPU model: amd or intel.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"nodepools": schema.ListNestedAttribute{
				Required:    true,
				Description: "Initial node pools for the cluster.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"label":      schema.StringAttribute{Required: true, Description: "Node pool label."},
						"size":       schema.StringAttribute{Required: true, Description: "Plan ID for worker node size."},
						"node_count": schema.StringAttribute{Required: true, Description: "Number of worker nodes."},
						"min_nodes":  schema.StringAttribute{Required: true, Description: "Minimum nodes for autoscaling."},
						"max_nodes":  schema.StringAttribute{Required: true, Description: "Maximum nodes for autoscaling."},
						"disk_size":  schema.StringAttribute{Optional: true, Description: "Additional EBS disk size in GB."},
						"disk_type":  schema.StringAttribute{Optional: true, Description: "EBS disk type: nvme or ssd."},
					},
				},
			},
		},
	}
}

func (r *KubernetesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *KubernetesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan KubernetesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var nodePools []client.KubernetesNodePool
	for _, np := range plan.NodePools {
		pool := client.KubernetesNodePool{
			Label:    np.Label.ValueString(),
			Size:     np.Size.ValueString(),
			Count:    np.NodeCount.ValueString(),
			MinNodes: np.MinNodes.ValueString(),
			MaxNodes: np.MaxNodes.ValueString(),
		}
		if !np.DiskSize.IsNull() && !np.DiskType.IsNull() &&
			np.DiskSize.ValueString() != "" && np.DiskType.ValueString() != "" {
			pool.EBS = []client.K8sEBSDisk{{
				Disk: np.DiskSize.ValueString(),
				Type: np.DiskType.ValueString(),
			}}
		}
		nodePools = append(nodePools, pool)
	}

	id, err := r.client.CreateKubernetesCluster(&client.KubernetesDeployRequest{
		DCSlug:         plan.DCSlug.ValueString(),
		ClusterLabel:   plan.ClusterLabel.ValueString(),
		ClusterVersion: plan.ClusterVersion.ValueString(),
		NodePools:      nodePools,
		VPC:            plan.VPC.ValueString(),
		NetworkType:    plan.NetworkType.ValueString(),
		CPUModel:       plan.CPUModel.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating Kubernetes cluster", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	plan.Status = types.StringValue("Pending")
	plan.IP = types.StringValue(r.client.GetK8sIP(id))
	plan.DNS = types.StringValue("")
	plan.CreatedAt = types.StringValue("")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *KubernetesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state KubernetesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cluster, err := r.client.GetKubernetesCluster(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading Kubernetes cluster", fmt.Sprintf("%s", err))
		return
	}
	if cluster == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Status = types.StringValue(cluster.Status)
	state.DNS = types.StringValue(cluster.DNS)
	state.CreatedAt = types.StringValue(cluster.CreatedAt)
	if cluster.IP != "" {
		state.IP = types.StringValue(cluster.IP)
	} else {
		state.IP = types.StringValue(r.client.GetK8sIP(state.ID.ValueString()))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *KubernetesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Kubernetes clusters cannot be updated in place. Destroy and recreate.")
}

func (r *KubernetesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state KubernetesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteKubernetesCluster(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting Kubernetes cluster", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// NODE POOL — utho_kubernetes_node_pool
// ══════════════════════════════════════════════════════════════

type KubernetesNodePoolResource struct{ client *client.Client }

type KubernetesNodePoolModel struct {
	ID        types.String `tfsdk:"id"`
	ClusterID types.String `tfsdk:"cluster_id"`
	PoolID    types.String `tfsdk:"pool_id"`
	Label     types.String `tfsdk:"label"`
	Size      types.String `tfsdk:"size"`
	NodeCount types.String `tfsdk:"node_count"`
	MinNodes  types.String `tfsdk:"min_nodes"`
	MaxNodes  types.String `tfsdk:"max_nodes"`
	DiskSize  types.String `tfsdk:"disk_size"`
	DiskType  types.String `tfsdk:"disk_type"`
}

func NewKubernetesNodePoolResource() resource.Resource { return &KubernetesNodePoolResource{} }

func (r *KubernetesNodePoolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_node_pool"
}

func (r *KubernetesNodePoolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Add and manage node pools in a Utho Kubernetes cluster.",
		Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"cluster_id": schema.StringAttribute{Required: true, Description: "Kubernetes cluster ID.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"pool_id":    schema.StringAttribute{Computed: true, Description: "Node pool ID assigned by Utho.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"label":      schema.StringAttribute{Required: true, Description: "Node pool label.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"size":       schema.StringAttribute{Required: true, Description: "Plan ID for worker node size.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"node_count": schema.StringAttribute{Required: true, Description: "Desired number of worker nodes. Updatable."},
			"min_nodes":  schema.StringAttribute{Required: true, Description: "Minimum nodes for autoscaling. Updatable."},
			"max_nodes":  schema.StringAttribute{Required: true, Description: "Maximum nodes for autoscaling. Updatable."},
			"disk_size":  schema.StringAttribute{Optional: true, Description: "EBS disk size in GB.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"disk_type":  schema.StringAttribute{Optional: true, Description: "EBS disk type: nvme or ssd.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *KubernetesNodePoolResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *KubernetesNodePoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan KubernetesNodePoolModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pool := client.KubernetesNodePool{
		Label:    plan.Label.ValueString(),
		Size:     plan.Size.ValueString(),
		Count:    plan.NodeCount.ValueString(),
		MinNodes: plan.MinNodes.ValueString(),
		MaxNodes: plan.MaxNodes.ValueString(),
	}
	if !plan.DiskSize.IsNull() && !plan.DiskType.IsNull() &&
		plan.DiskSize.ValueString() != "" && plan.DiskType.ValueString() != "" {
		pool.EBS = []client.K8sEBSDisk{{
			Disk: plan.DiskSize.ValueString(),
			Type: plan.DiskType.ValueString(),
		}}
	}

	err := r.client.AddNodePool(plan.ClusterID.ValueString(), []client.KubernetesNodePool{pool})
	if err != nil {
		resp.Diagnostics.AddError("Error adding node pool", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", plan.ClusterID.ValueString(), plan.Label.ValueString()))
	plan.PoolID = types.StringValue(plan.Label.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *KubernetesNodePoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state KubernetesNodePoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Node pool state is maintained locally — no GET API for individual pool
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *KubernetesNodePoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan KubernetesNodePoolModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateNodePool(plan.ClusterID.ValueString(), plan.PoolID.ValueString(), &client.NodePoolUpdateRequest{
		Count:    plan.NodeCount.ValueString(),
		MinNodes: plan.MinNodes.ValueString(),
		MaxNodes: plan.MaxNodes.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating node pool", fmt.Sprintf("%s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *KubernetesNodePoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state KubernetesNodePoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteNodePool(state.ClusterID.ValueString(), state.PoolID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting node pool", fmt.Sprintf("%s", err))
	}
}
