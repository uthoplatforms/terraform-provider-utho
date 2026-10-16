package resources

import (
	"context"
	"fmt"

	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ══════════════════════════════════════════════════════════════
// PROJECT — utho_project
// ══════════════════════════════════════════════════════════════

type ProjectResource struct{ client *client.Client }

type ProjectModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Environment   types.String `tfsdk:"environment"`
	Status        types.String `tfsdk:"status"`
	IsDefault     types.Bool   `tfsdk:"is_default"`
	MemberCount   types.Int64  `tfsdk:"member_count"`
	ResourceCount types.Int64  `tfsdk:"resource_count"`
	CreatedAt     types.String `tfsdk:"created_at"`
}

func NewProjectResource() resource.Resource { return &ProjectResource{} }

func (r *ProjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage Utho Projects to organize resources and team members.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Project status.",
			},
			"is_default": schema.BoolAttribute{
				Computed:      true,
				Description:   "Whether this is the default project.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"member_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of members in the project.",
			},
			"resource_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of resources in the project.",
			},
			"created_at": schema.StringAttribute{
				Computed:      true,
				Description:   "Creation timestamp.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Project name.",
			},
			"environment": schema.StringAttribute{
				Required:    true,
				Description: "Environment: development, staging, production, or testing.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Project description.",
			},
		},
	}
}

func (r *ProjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ProjectModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.client.CreateProject(&client.ProjectCreateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Environment: plan.Environment.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating project", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%d", project.ID))
	plan.Status = types.StringValue(project.Status)
	plan.IsDefault = types.BoolValue(project.IsDefault)
	plan.MemberCount = types.Int64Value(int64(project.MemberCount))
	plan.ResourceCount = types.Int64Value(int64(project.ResourceCount))
	plan.CreatedAt = types.StringValue(project.CreatedAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ProjectModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.client.GetProject(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading project", fmt.Sprintf("%s", err))
		return
	}
	if project == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(project.Name)
	state.Description = types.StringValue(project.Description)
	state.Environment = types.StringValue(project.Environment)
	state.Status = types.StringValue(project.Status)
	state.IsDefault = types.BoolValue(project.IsDefault)
	state.MemberCount = types.Int64Value(int64(project.MemberCount))
	state.ResourceCount = types.Int64Value(int64(project.ResourceCount))
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ProjectModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.client.UpdateProject(plan.ID.ValueString(), &client.ProjectCreateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Environment: plan.Environment.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating project", fmt.Sprintf("%s", err))
		return
	}

	plan.Status = types.StringValue(project.Status)
	plan.MemberCount = types.Int64Value(int64(project.MemberCount))
	plan.ResourceCount = types.Int64Value(int64(project.ResourceCount))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ProjectModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteProject(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting project", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// PROJECT MEMBER — utho_project_member
// ══════════════════════════════════════════════════════════════

type ProjectMemberResource struct{ client *client.Client }

type ProjectMemberModel struct {
	ID        types.String `tfsdk:"id"`
	ProjectID types.String `tfsdk:"project_id"`
	UserID    types.Int64  `tfsdk:"user_id"`
}

func NewProjectMemberResource() resource.Resource { return &ProjectMemberResource{} }

func (r *ProjectMemberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_member"
}

func (r *ProjectMemberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Add and manage members in a Utho Project.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"project_id": schema.StringAttribute{
				Required:      true,
				Description:   "Project ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"user_id": schema.Int64Attribute{
				Required:      true,
				Description:   "IAM user ID to add as member.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *ProjectMemberResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ProjectMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ProjectMemberModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.AddProjectMember(plan.ProjectID.ValueString(), int(plan.UserID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error adding project member", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s:%d", plan.ProjectID.ValueString(), plan.UserID.ValueInt64()))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ProjectMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *ProjectMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ProjectMemberModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove then re-add with new role
	if err := r.client.RemoveProjectMember(plan.ProjectID.ValueString(), fmt.Sprintf("%d", plan.UserID.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Error removing project member for update", fmt.Sprintf("%s", err))
		return
	}
	if err := r.client.AddProjectMember(plan.ProjectID.ValueString(), int(plan.UserID.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Error re-adding project member with new role", fmt.Sprintf("%s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ProjectMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ProjectMemberModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.RemoveProjectMember(state.ProjectID.ValueString(), fmt.Sprintf("%d", state.UserID.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Error removing project member", fmt.Sprintf("%s", err))
	}
}
