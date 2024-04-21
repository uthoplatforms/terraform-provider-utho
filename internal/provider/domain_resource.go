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
	"github.com/uthoterraform/terraform-provider-utho/api"
)

// implement resource interfaces.
var (
	_ resource.Resource                = &DomainResource{}
	_ resource.ResourceWithConfigure   = &DomainResource{}
	_ resource.ResourceWithImportState = &DomainResource{}
)

// NewDomainResource is a helper function to simplify the provider implementation.
func NewDomainResource() resource.Resource {
	return &DomainResource{}
}

// DomainResource is the resource implementation.
type (
	DomainResource struct {
		client *api.Client
	}

	// DomainResource is the model implementation.
	DomainResourceModel struct {
		Domain types.String `tfsdk:"domain"`
	}
)

// Metadata returns the resource type name.
func (s *DomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

// Configure adds the provider configured client to the data source.
func (d *DomainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Domain Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

// Schema defines the schema for the resource.
func (s *DomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Required: true,
				// Requires Replace if the value change
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Must be unique",
			},
		},
	}
}

// Import using domain as the attribute
func (s *DomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("domain"), req, resp)
}

// Create a new resource.
func (s *DomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "create domain")
	// Retrieve values from plan
	var plan DomainResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	domainRequest := api.DomainRequest{
		Domain: plan.Domain.ValueString(),
	}
	tflog.Debug(ctx, "send create domain request")
	domain, err := s.client.CreateDomain(ctx, domainRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating domain",
			"Could not create domain, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = DomainResourceModel{
		Domain: types.StringValue(domain.Domain),
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish create domain")
}

// Read resource information.
func (s *DomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "read domain")

	// Get current state
	var state DomainResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "send get domain request")
	// Get refreshed domain value from utho
	domain, err := s.client.GetDomain(ctx, state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading utho domain",
			"Could not read utho domain domain "+state.Domain.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state = DomainResourceModel{
		Domain: types.StringValue(domain.Domain),
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish get domain request")
}

func (s *DomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// updating resource is not supported
}

// Delete deletes the resource and removes the Terraform state on success.
func (s *DomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "delete domain")
	// Get current state
	var state DomainResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "send delete domain request")
	// delete domain
	err := s.client.DeleteDomain(ctx, state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing utho domain",
			"Could not delete utho domain "+state.Domain.ValueString()+": "+err.Error(),
		)
		return
	}
}
