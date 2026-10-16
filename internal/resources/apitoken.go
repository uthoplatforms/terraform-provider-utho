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

type APITokenResource struct{ client *client.Client }

type APITokenModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Write     types.String `tfsdk:"write"`
	Token     types.String `tfsdk:"token"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func NewAPITokenResource() resource.Resource { return &APITokenResource{} }

func (r *APITokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_token"
}

func (r *APITokenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage Utho API tokens.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique API token ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Token name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"write": schema.StringAttribute{
				Required:    true,
				Description: "Write access: on (read/write) or off (read only).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"token": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Generated API token value. Sensitive — shown only at creation time.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the token was created.",
			},
		},
	}
}

func (r *APITokenResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *APITokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan APITokenModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apikey, _, err := r.client.CreateAPIToken(&client.APITokenCreateRequest{
		Name:  plan.Name.ValueString(),
		Write: plan.Write.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating API token", fmt.Sprintf("%s", err))
		return
	}

	// API doesn't return ID on create — fetch list and find by name
	tokens, err := r.client.ListAPITokens()
	if err != nil {
		resp.Diagnostics.AddError("Error listing tokens after create", fmt.Sprintf("%s", err))
		return
	}

	foundID := ""
	foundAt := ""
	for _, t := range tokens {
		if t.Name == plan.Name.ValueString() {
			foundID = t.ID
			foundAt = t.CreatedAt
		}
	}

	if foundID == "" {
		resp.Diagnostics.AddError("Token not found after create",
			fmt.Sprintf("Token with name '%s' not found after creation.", plan.Name.ValueString()))
		return
	}

	plan.ID = types.StringValue(foundID)
	plan.Token = types.StringValue(apikey)
	plan.CreatedAt = types.StringValue(foundAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *APITokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state APITokenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	token, err := r.client.GetAPIToken(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading API token", fmt.Sprintf("%s", err))
		return
	}
	if token == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(token.Name)
	state.Write = types.StringValue(token.Write)
	state.CreatedAt = types.StringValue(token.CreatedAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *APITokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "API tokens cannot be updated. Destroy and recreate.")
}

func (r *APITokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state APITokenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteAPIToken(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting API token", fmt.Sprintf("%s", err))
	}
}
