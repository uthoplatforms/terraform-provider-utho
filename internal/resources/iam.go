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

type IAMUserResource struct{ client *client.Client }

type IAMUserModel struct {
	ID          types.String `tfsdk:"id"`
	FullName    types.String `tfsdk:"fullname"`
	Email       types.String `tfsdk:"email"`
	MobileCC    types.String `tfsdk:"mobilecc"`
	Mobile      types.String `tfsdk:"mobile"`
	Permissions types.String `tfsdk:"permissions"`
	Status      types.String `tfsdk:"status"`
	DateAdded   types.String `tfsdk:"date_added"`
}

func NewIAMUserResource() resource.Resource { return &IAMUserResource{} }

func (r *IAMUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_user"
}

func (r *IAMUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage IAM sub-users in your Utho account.",
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"status":      schema.StringAttribute{Computed: true, Description: "User status (Pending, Active)."},
			"date_added":  schema.StringAttribute{Computed: true, Description: "Date user was added.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"fullname":    schema.StringAttribute{Required: true, Description: "Full name of the IAM user.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"email":       schema.StringAttribute{Required: true, Description: "Email address. Must be unique.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"mobilecc":    schema.StringAttribute{Required: true, Description: "Mobile country code (e.g. 91).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"mobile":      schema.StringAttribute{Required: true, Description: "Mobile number.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"permissions": schema.StringAttribute{Required: true, Description: "Comma-separated permissions string."},
		},
	}
}

func (r *IAMUserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IAMUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan IAMUserModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateIAMUser(&client.IAMUserCreateRequest{
		FullName:    plan.FullName.ValueString(),
		Email:       plan.Email.ValueString(),
		MobileCC:    plan.MobileCC.ValueString(),
		Mobile:      plan.Mobile.ValueString(),
		Permissions: plan.Permissions.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating IAM user", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	plan.Status = types.StringValue("Pending")
	plan.DateAdded = types.StringValue("")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *IAMUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IAMUserModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.GetIAMUser(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading IAM user", fmt.Sprintf("%s", err))
		return
	}
	if user == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Status = types.StringValue(user.Status)
	state.DateAdded = types.StringValue(user.DateAdded)
	state.Permissions = types.StringValue(user.Permissions)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IAMUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state IAMUserModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateIAMUser(state.ID.ValueString(), &client.IAMUserUpdateRequest{
		Permissions: plan.Permissions.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating IAM user", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = state.ID
	plan.Status = state.Status
	plan.DateAdded = state.DateAdded
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *IAMUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state IAMUserModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteIAMUser(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting IAM user", fmt.Sprintf("%s", err))
	}
}
