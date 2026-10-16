package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SSHKeyResource struct{ client *client.Client }

type SSHKeyModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	SSHKey    types.String `tfsdk:"sshkey"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func NewSSHKeyResource() resource.Resource { return &SSHKeyResource{} }

func (r *SSHKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key"
}

func (r *SSHKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Import and manage SSH keys on Utho Cloud.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique SSH key ID assigned by Utho.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the SSH key.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sshkey": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Public SSH key content (e.g. ssh-ed25519 AAAA... user@host).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the key was created.",
			},
		},
	}
}

func (r *SSHKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SSHKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SSHKeyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateSSHKey(&client.SSHKeyCreateRequest{
		Name:   plan.Name.ValueString(),
		SSHKey: plan.SSHKey.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating SSH key", fmt.Sprintf("%s", err))
		return
	}

	// API does not return key ID on create — fetch list and find by name
	respBytes, err := r.client.Get("/key")
	if err != nil {
		resp.Diagnostics.AddError("Error fetching SSH keys after create", fmt.Sprintf("%s", err))
		return
	}

	var listResp client.SSHKeyListResponse
	if err := json.Unmarshal(respBytes, &listResp); err != nil {
		resp.Diagnostics.AddError("Error parsing SSH key list", fmt.Sprintf("%s", err))
		return
	}

	// Find by name — take the last match (most recently created)
	foundID := ""
	foundAt := ""
	for _, key := range listResp.Key {
		if key.Name == plan.Name.ValueString() {
			foundID = key.ID
			foundAt = key.CreatedAt
		}
	}

	if foundID == "" {
		resp.Diagnostics.AddError(
			"SSH key not found after create",
			fmt.Sprintf("Key with name '%s' not found after creation.", plan.Name.ValueString()),
		)
		return
	}

	plan.ID = types.StringValue(foundID)
	plan.CreatedAt = types.StringValue(foundAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *SSHKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SSHKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	key, err := r.client.GetSSHKey(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading SSH key", fmt.Sprintf("%s", err))
		return
	}

	if key == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(key.Name)
	state.CreatedAt = types.StringValue(key.CreatedAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SSHKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "SSH keys cannot be updated. Destroy and recreate.")
}

func (r *SSHKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SSHKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSSHKey(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting SSH key", fmt.Sprintf("%s", err))
	}
}
