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
	_ resource.Resource                = &SqsResource{}
	_ resource.ResourceWithConfigure   = &SqsResource{}
	_ resource.ResourceWithImportState = &SqsResource{}
)

// NewSqsResource is a helper function to simplify the provider implementation.
func NewSqsResource() resource.Resource {
	return &SqsResource{}
}

// SqsResource is the resource implementation.
type SqsResource struct {
	client utho.Client
}

// SqsResource is the model implementation.
type SqsResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Dcslug    types.String `tfsdk:"dcslug"`
	Planid    types.String `tfsdk:"planid"`
	Userid    types.String `tfsdk:"userid"`
	Cloudid   types.String `tfsdk:"cloudid"`
	Status    types.String `tfsdk:"status"`
	CreatedAt types.String `tfsdk:"created_at"`
	IP        types.String `tfsdk:"ip"`
	Count     types.String `tfsdk:"count"`
}

// Metadata returns the resource type name.
func (s *SqsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sqs"
}

// Configure adds the provider configured client to the data source.
func (d *SqsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(utho.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Sqs Data Source Configure Type",
			fmt.Sprintf("Expected utho.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

// Schema defines the schema for the resource.
func (s *SqsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "id"},
			"name": schema.StringAttribute{Required: true, Description: "Provide SQS name eg: sqs-ywqo2pmc",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"dcslug": schema.StringAttribute{Required: true, Description: "Provide dcslug eg: innoida",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"planid": schema.StringAttribute{Required: true, Description: "Provide the planid eg: 10045",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"userid":     schema.StringAttribute{Computed: true, Description: "userid"},
			"cloudid":    schema.StringAttribute{Computed: true, Description: "cloudid"},
			"status":     schema.StringAttribute{Computed: true, Description: "status"},
			"created_at": schema.StringAttribute{Computed: true, Description: "created_at"},
			"ip":         schema.StringAttribute{Computed: true, Description: "ip"},
			"count":      schema.StringAttribute{Computed: true, Description: "count"},
		},
	}
}

// Import using sqs as the attribute
func (s *SqsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource.
func (s *SqsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "create sqs")
	// Retrieve values from plan
	var plan SqsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	sqsRequest := utho.CreateSqsParams{
		Name:   plan.Name.ValueString(),
		Dcslug: plan.Dcslug.ValueString(),
		Planid: plan.Planid.ValueString(),
	}
	tflog.Debug(ctx, "send create sqs request")
	sqs, err := s.client.Sqs().Create(sqsRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating sqs",
			"Could not create sqs, unexpected error: "+err.Error(),
		)
		return
	}

	// get sqs data
	getSqs, err := s.client.Sqs().Read(sqs.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading utho sqs",
			"Could not read utho sqs "+sqs.ID+": "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = SqsResourceModel{
		ID:        types.StringValue(sqs.ID),
		Dcslug:    types.StringValue(plan.Dcslug.ValueString()),
		Name:      types.StringValue(plan.Name.ValueString()),
		Planid:    types.StringValue(plan.Planid.ValueString()),
		Userid:    types.StringValue(getSqs.Userid),
		Cloudid:   types.StringValue(getSqs.Cloudid),
		Status:    types.StringValue(getSqs.Status),
		CreatedAt: types.StringValue(getSqs.CreatedAt),
		IP:        types.StringValue(getSqs.IP),
		Count:     types.StringValue(getSqs.Count),
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish create sqs")
}

// Read resource information.
func (s *SqsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "read sqs")

	// Get current state
	var state SqsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "send get sqs request")
	// Get refreshed sqs value from utho
	sqs, err := s.client.Sqs().Read(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading utho sqs",
			"Could not read utho sqs sqs "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state = SqsResourceModel{
		ID: types.StringValue(sqs.ID),
		// Dcslug:    types.StringValue(sqs.Dcslug),
		Name: types.StringValue(sqs.Name),
		// Planid:    types.StringValue(sqs.Planid),
		Userid:    types.StringValue(sqs.Userid),
		Cloudid:   types.StringValue(sqs.Cloudid),
		Status:    types.StringValue(sqs.Status),
		CreatedAt: types.StringValue(sqs.CreatedAt),
		IP:        types.StringValue(sqs.IP),
		Count:     types.StringValue(sqs.Count),
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish get sqs request")
}

func (s *SqsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// updating resource is
}

// Delete deletes the resource and removes the Terraform state on success.
func (s *SqsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError(
		"Error deleteing utho sqs",
		"Could not delete utho sqs ",
	)
	tflog.Debug(ctx, "delete sqs")
	// Get current state
	var state SqsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "send delete sqs request")
	// delete sqs
	_, err := s.client.Sqs().Delete(state.ID.ValueString(), state.Name.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing utho sqs",
			"Could not delete utho sqs "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}
