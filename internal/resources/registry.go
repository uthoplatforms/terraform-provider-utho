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
// CONTAINER REGISTRY — utho_container_registry
// ══════════════════════════════════════════════════════════════

type ContainerRegistryResource struct{ client *client.Client }

type ContainerRegistryModel struct {
	ID           types.String `tfsdk:"id"`
	ProjectName  types.String `tfsdk:"project_name"`
	DCSlug       types.String `tfsdk:"dcslug"`
	PlanID       types.Int64  `tfsdk:"planid"`
	BillingCycle types.String `tfsdk:"billingcycle"`
	Public       types.Bool   `tfsdk:"public"`
	CreatedAt    types.String `tfsdk:"created_at"`
}

func NewContainerRegistryResource() resource.Resource { return &ContainerRegistryResource{} }

func (r *ContainerRegistryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_container_registry"
}

func (r *ContainerRegistryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage Utho Container Registries.",
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"created_at":   schema.StringAttribute{Computed: true, Description: "Creation timestamp."},
			"project_name": schema.StringAttribute{Required: true, Description: "Registry project name (e.g. my-registry).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"dcslug":       schema.StringAttribute{Required: true, Description: "Data center slug.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"planid":       schema.Int64Attribute{Required: true, Description: "Registry plan ID (e.g. 10276 for Starter).", PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()}},
			"billingcycle": schema.StringAttribute{Required: true, Description: "Billing cycle: monthly.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"public":       schema.BoolAttribute{Required: true, Description: "Set true for public registry, false for private."},
		},
	}
}

func (r *ContainerRegistryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ContainerRegistryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ContainerRegistryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.CreateContainerRegistry(&client.ContainerRegistryCreateRequest{
		DCSlug:       plan.DCSlug.ValueString(),
		PlanID:       int(plan.PlanID.ValueInt64()),
		BillingCycle: plan.BillingCycle.ValueString(),
		Public:       fmt.Sprintf("%v", plan.Public.ValueBool()),
		ProjectName:  plan.ProjectName.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating container registry", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(plan.ProjectName.ValueString())
	plan.CreatedAt = types.StringValue("")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ContainerRegistryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ContainerRegistryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GetContainerRegistry may fail due to permission timing — keep state on error
	reg, err := r.client.GetContainerRegistry(state.ProjectName.ValueString())
	if err != nil || reg == nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
		return
	}
	state.CreatedAt = types.StringValue(reg.CreatedAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ContainerRegistryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ContainerRegistryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateContainerRegistry(plan.ProjectName.ValueString(), fmt.Sprintf("%v", plan.Public.ValueBool())); err != nil {
		resp.Diagnostics.AddError("Error updating container registry", fmt.Sprintf("%s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ContainerRegistryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ContainerRegistryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteContainerRegistry(state.ProjectName.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting container registry", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// REGISTRY ROBOT — utho_container_registry_robot
// ══════════════════════════════════════════════════════════════

type RegistryRobotResource struct{ client *client.Client }

type RegistryRobotAccessModel struct {
	Resource types.String `tfsdk:"resource"`
	Action   types.String `tfsdk:"action"`
}

type RegistryRobotModel struct {
	ID          types.String               `tfsdk:"id"`
	ProjectName types.String               `tfsdk:"project_name"`
	Name        types.String               `tfsdk:"name"`
	Description types.String               `tfsdk:"description"`
	Duration    types.Int64                `tfsdk:"duration"`
	Secret      types.String               `tfsdk:"secret"`
	ExpiresAt   types.Int64                `tfsdk:"expires_at"`
	Access      []RegistryRobotAccessModel `tfsdk:"access"`
}

func NewRegistryRobotResource() resource.Resource { return &RegistryRobotResource{} }

func (r *RegistryRobotResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_container_registry_robot"
}

func (r *RegistryRobotResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage robot accounts for a Utho Container Registry.",
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"secret":       schema.StringAttribute{Computed: true, Sensitive: true, Description: "Robot account secret. Shown only at creation time."},
			"expires_at":   schema.Int64Attribute{Computed: true, Description: "Expiry timestamp."},
			"project_name": schema.StringAttribute{Required: true, Description: "Registry project name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"name":         schema.StringAttribute{Required: true, Description: "Robot account name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"description":  schema.StringAttribute{Optional: true, Description: "Robot account description.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"duration":     schema.Int64Attribute{Required: true, Description: "Token validity in days (e.g. 30).", PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()}},
			"access": schema.ListNestedAttribute{
				Required:    true,
				Description: "Access permissions for the robot account.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"resource": schema.StringAttribute{Required: true, Description: "Resource type: repository."},
						"action":   schema.StringAttribute{Required: true, Description: "Action: pull or push."},
					},
				},
			},
		},
	}
}

func (r *RegistryRobotResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RegistryRobotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RegistryRobotModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var access []client.RegistryRobotAccess
	for _, a := range plan.Access {
		access = append(access, client.RegistryRobotAccess{
			Resource: a.Resource.ValueString(),
			Action:   a.Action.ValueString(),
		})
	}

	robot, err := r.client.CreateRegistryRobot(plan.ProjectName.ValueString(), &client.RegistryRobotCreateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Duration:    int(plan.Duration.ValueInt64()),
		Access:      access,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating registry robot", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%d", robot.ID))
	plan.Secret = types.StringValue(robot.Secret)
	plan.ExpiresAt = types.Int64Value(robot.ExpiresAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *RegistryRobotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}
func (r *RegistryRobotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Registry robots cannot be updated. Destroy and recreate.")
}
func (r *RegistryRobotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RegistryRobotModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	robotID := 0
	fmt.Sscanf(state.ID.ValueString(), "%d", &robotID)
	if err := r.client.DeleteRegistryRobot(state.ProjectName.ValueString(), robotID); err != nil {
		resp.Diagnostics.AddError("Error deleting registry robot", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// REGISTRY WEBHOOK — utho_container_registry_webhook
// ══════════════════════════════════════════════════════════════

type RegistryWebhookResource struct{ client *client.Client }

type RegistryWebhookModel struct {
	ID             types.String `tfsdk:"id"`
	ProjectName    types.String `tfsdk:"project_name"`
	Name           types.String `tfsdk:"name"`
	Address        types.String `tfsdk:"address"`
	AuthHeader     types.String `tfsdk:"auth_header"`
	SkipCertVerify types.Bool   `tfsdk:"skip_cert_verify"`
	PayloadFormat  types.String `tfsdk:"payload_format"`
	EventTypes     types.List   `tfsdk:"event_types"`
}

func NewRegistryWebhookResource() resource.Resource { return &RegistryWebhookResource{} }

func (r *RegistryWebhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_container_registry_webhook"
}

func (r *RegistryWebhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage webhooks for a Utho Container Registry.",
		Attributes: map[string]schema.Attribute{
			"id":               schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"project_name":     schema.StringAttribute{Required: true, Description: "Registry project name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"name":             schema.StringAttribute{Required: true, Description: "Webhook name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"address":          schema.StringAttribute{Required: true, Description: "Webhook endpoint URL.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"payload_format":   schema.StringAttribute{Required: true, Description: "Payload format: Default or CloudEvents.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"skip_cert_verify": schema.BoolAttribute{Optional: true, Description: "Skip TLS certificate verification.", PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()}},
			"auth_header":      schema.StringAttribute{Optional: true, Sensitive: true, Description: "Authorization header value.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"event_types": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "Event types to trigger the webhook.",
			},
		},
	}
}

func (r *RegistryWebhookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RegistryWebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RegistryWebhookModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var eventTypes []string
	resp.Diagnostics.Append(plan.EventTypes.ElementsAs(ctx, &eventTypes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateRegistryWebhook(plan.ProjectName.ValueString(), &client.RegistryWebhookCreateRequest{
		Name:           plan.Name.ValueString(),
		Address:        plan.Address.ValueString(),
		EventTypes:     eventTypes,
		AuthHeader:     plan.AuthHeader.ValueString(),
		SkipCertVerify: plan.SkipCertVerify.ValueBool(),
		PayloadFormat:  plan.PayloadFormat.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating registry webhook", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *RegistryWebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}
func (r *RegistryWebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Registry webhooks cannot be updated. Destroy and recreate.")
}
func (r *RegistryWebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RegistryWebhookModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	webhookID := 0
	fmt.Sscanf(state.ID.ValueString(), "%d", &webhookID)
	if err := r.client.DeleteRegistryWebhook(state.ProjectName.ValueString(), webhookID); err != nil {
		resp.Diagnostics.AddError("Error deleting registry webhook", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// IMMUTABLE TAG RULE — utho_container_registry_immutable_rule
// ══════════════════════════════════════════════════════════════

type RegistryImmutableRuleResource struct{ client *client.Client }

type RegistryImmutableRuleModel struct {
	ID          types.String `tfsdk:"id"`
	ProjectName types.String `tfsdk:"project_name"`
	TagPattern  types.String `tfsdk:"tag_pattern"`
	RepoPattern types.String `tfsdk:"repo_pattern"`
	RuleID      types.Int64  `tfsdk:"rule_id"`
}

func NewRegistryImmutableRuleResource() resource.Resource { return &RegistryImmutableRuleResource{} }

func (r *RegistryImmutableRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_container_registry_immutable_rule"
}

func (r *RegistryImmutableRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create immutable tag rules for a Utho Container Registry.",
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"rule_id":      schema.Int64Attribute{Computed: true, Description: "Rule ID assigned by the registry.", PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}},
			"project_name": schema.StringAttribute{Required: true, Description: "Registry project name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"tag_pattern":  schema.StringAttribute{Required: true, Description: "Tag pattern to make immutable (e.g. v*).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"repo_pattern": schema.StringAttribute{Required: true, Description: "Repository pattern (e.g. **).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *RegistryImmutableRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RegistryImmutableRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RegistryImmutableRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateImmutableRule(plan.ProjectName.ValueString(), &client.RegistryImmutableCreateRequest{
		TagPattern:  plan.TagPattern.ValueString(),
		RepoPattern: plan.RepoPattern.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating immutable rule", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", plan.ProjectName.ValueString(), plan.TagPattern.ValueString()))
	plan.RuleID = types.Int64Value(0)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *RegistryImmutableRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}
func (r *RegistryImmutableRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Immutable rules cannot be updated. Destroy and recreate.")
}
func (r *RegistryImmutableRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RegistryImmutableRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteImmutableRule(state.ProjectName.ValueString(), int(state.RuleID.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Error deleting immutable rule", fmt.Sprintf("%s", err))
	}
}
