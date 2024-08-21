package provider

import (
	"context"
	"fmt"

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
	_ resource.Resource                = &LoadbalancerResource{}
	_ resource.ResourceWithConfigure   = &LoadbalancerResource{}
	_ resource.ResourceWithImportState = &LoadbalancerResource{}
)

// NewLoadbalancerResource is a helper function to simplify the provider implementation.
func NewLoadbalancerResource() resource.Resource {
	return &LoadbalancerResource{}
}

// LoadbalancerResource is the resource implementation.
type (
	LoadbalancerResource struct {
		client utho.Client
	}

	// LoadbalancerResourceModel is the model implementation.
	LoadbalancerResourceModel struct {
		ID            types.String `tfsdk:"id"`
		Type          types.String `tfsdk:"type"`
		Dcslug        types.String `tfsdk:"dcslug"`
		Userid        types.String `tfsdk:"userid"`
		IP            types.String `tfsdk:"ip"`
		Name          types.String `tfsdk:"name"`
		Algorithm     types.String `tfsdk:"algorithm"`
		Cookie        types.String `tfsdk:"cookie"`
		Cookiename    types.String `tfsdk:"cookiename"`
		Redirecthttps types.String `tfsdk:"redirecthttps"`
		Country       types.String `tfsdk:"country"`
		Cc            types.String `tfsdk:"cc"`
		City          types.String `tfsdk:"city"`
		Backendcount  types.String `tfsdk:"backendcount"`
		CreatedAt     types.String `tfsdk:"created_at"`
		Status        types.String `tfsdk:"status"`
	}

	RuleResourceModel struct {
		ID          types.String `tfsdk:"id"`
		Lb          types.String `tfsdk:"lb"`
		SrcProto    types.String `tfsdk:"src_proto"`
		SrcPort     types.String `tfsdk:"src_port"`
		DstProto    types.String `tfsdk:"dst_proto"`
		DstPort     types.String `tfsdk:"dst_port"`
		Timeadded   types.String `tfsdk:"timeadded"`
		Timeupdated types.String `tfsdk:"timeupdated"`
	}
)

// Metadata returns the resource type name.
func (s *LoadbalancerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer"
}

// Configure adds the provider configured client to the data source.
func (d *LoadbalancerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(utho.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Loadbalancer Data Source Configure Type",
			fmt.Sprintf("Expected utho.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

// Schema defines the schema for the resource.
func (s *LoadbalancerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"dcslug": schema.StringAttribute{
				Required:    true,
				Description: "Provide Zone dcslug eg: innoida",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Load Balancer name eg: webapplb",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Load-Balancer type must be either application or network. The default value is application",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id":            schema.StringAttribute{Computed: true, Description: "Id"},
			"userid":        schema.StringAttribute{Computed: true, Description: "User id"},
			"ip":            schema.StringAttribute{Computed: true, Description: "Ip"},
			"algorithm":     schema.StringAttribute{Computed: true, Description: "Algorithm"},
			"cookie":        schema.StringAttribute{Computed: true, Description: "Cookie"},
			"cookiename":    schema.StringAttribute{Computed: true, Description: "Cookie name"},
			"redirecthttps": schema.StringAttribute{Computed: true, Description: "Redirect https"},
			"country":       schema.StringAttribute{Computed: true, Description: "Country"},
			"cc":            schema.StringAttribute{Computed: true, Description: "Cc"},
			"city":          schema.StringAttribute{Computed: true, Description: "City"},
			"backendcount":  schema.StringAttribute{Computed: true, Description: "Backend count"},
			"created_at":    schema.StringAttribute{Computed: true, Description: "Created At"},
			"status":        schema.StringAttribute{Computed: true, Description: "Status"},
		},
	}
}

// Import using id as the attribute
func (s *LoadbalancerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource.
func (s *LoadbalancerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "create load balncer")
	// Retrieve values from plan
	var plan LoadbalancerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	loadbalancerRequest := utho.CreateLoadbalancerParams{
		Dcslug: plan.Dcslug.ValueString(),
		Type:   plan.Type.ValueString(),
		Name:   plan.Name.ValueString(),
	}
	tflog.Debug(ctx, "send create loadbalancer request")

	createloadbalancer, err := s.client.Loadbalancers().Create(loadbalancerRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating loadbalancer",
			"Could not create loadbalancer, unexpected error: "+err.Error(),
		)
		return
	}

	loadbalancer, err := s.client.Loadbalancers().Read(createloadbalancer.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating loadbalancer",
			"Could not create loadbalancer, unexpected error: "+err.Error(),
		)
		return
	}

	plan = LoadbalancerResourceModel{
		ID:            types.StringValue(loadbalancer.ID),
		Type:          types.StringValue(loadbalancer.Type),
		Dcslug:        types.StringValue(plan.Dcslug.ValueString()),
		Userid:        types.StringValue(loadbalancer.Userid),
		IP:            types.StringValue(loadbalancer.IP),
		Name:          types.StringValue(loadbalancer.Name),
		Algorithm:     types.StringValue(loadbalancer.Algorithm),
		Cookie:        types.StringValue(loadbalancer.Cookie),
		Cookiename:    types.StringValue(loadbalancer.Cookiename),
		Redirecthttps: types.StringValue(loadbalancer.Redirecthttps),
		Country:       types.StringValue(loadbalancer.Country),
		Cc:            types.StringValue(loadbalancer.Cc),
		City:          types.StringValue(loadbalancer.City),
		Backendcount:  types.StringValue(loadbalancer.Backendcount),
		CreatedAt:     types.StringValue(loadbalancer.CreatedAt),
		Status:        types.StringValue(loadbalancer.Status),
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish create loadbalancer")
}

// Read resource information.
func (s *LoadbalancerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "read loadbalancer")

	// Get current state
	var state LoadbalancerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "send get loadbalancer request")
	// Get refreshed loadbalancer value from utho
	loadbalancer, err := s.client.Loadbalancers().Read(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading utho loadbalancer",
			"Could not read utho loadbalancer loadbalancer "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state = LoadbalancerResourceModel{
		ID:            types.StringValue(loadbalancer.ID),
		Type:          types.StringValue(loadbalancer.Type),
		Dcslug:        types.StringValue(state.Dcslug.ValueString()),
		Userid:        types.StringValue(loadbalancer.Userid),
		IP:            types.StringValue(loadbalancer.IP),
		Name:          types.StringValue(loadbalancer.Name),
		Algorithm:     types.StringValue(loadbalancer.Algorithm),
		Cookie:        types.StringValue(loadbalancer.Cookie),
		Cookiename:    types.StringValue(loadbalancer.Cookiename),
		Redirecthttps: types.StringValue(loadbalancer.Redirecthttps),
		Country:       types.StringValue(loadbalancer.Country),
		Cc:            types.StringValue(loadbalancer.Cc),
		City:          types.StringValue(loadbalancer.City),
		Backendcount:  types.StringValue(loadbalancer.Backendcount),
		CreatedAt:     types.StringValue(loadbalancer.CreatedAt),
		Status:        types.StringValue(loadbalancer.Status),
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish get loadbalancer request")
}

func (s *LoadbalancerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// updating resource is not supported
}

// Delete deletes the resource and removes the Terraform state on success.
func (s *LoadbalancerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "delete loadbalancer")
	// Get current state
	var state LoadbalancerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "send delete loadbalancer request")
	// delete loadbalancer
	_, err := s.client.Loadbalancers().Delete(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing utho loadbalancer",
			"Could not delete utho loadbalancer "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}
