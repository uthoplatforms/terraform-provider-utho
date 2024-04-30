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
	"github.com/uthoplatforms/terraform-provider-utho/api"
)

// implement resource interfaces.
var (
	_ resource.Resource                = &DnsRecordResource{}
	_ resource.ResourceWithConfigure   = &DnsRecordResource{}
	_ resource.ResourceWithImportState = &DnsRecordResource{}
)

// NewDnsRecordResource is a helper function to simplify the provider implementation.
func NewDnsRecordResource() resource.Resource {
	return &DnsRecordResource{}
}

// DnsRecordResource is the resource implementation.
type DnsRecordResource struct {
	client *api.Client
}

type DnsRecordResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Domain   types.String `tfsdk:"domain"`
	Type     types.String `tfsdk:"type"`
	Hostname types.String `tfsdk:"hostname"`
	Value    types.String `tfsdk:"value"`
	TTL      types.String `tfsdk:"ttl"`
	Porttype types.String `tfsdk:"porttype"`
	Port     types.String `tfsdk:"port"`
	Priority types.String `tfsdk:"priority"`
	Wight    types.String `tfsdk:"wight"`
}

// Metadata returns the resource type name.
func (s *DnsRecordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

// Configure adds the provider configured client to the data source.
func (d *DnsRecordResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DnsRecord Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

// Schema defines the schema for the resource.
func (s *DnsRecordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":       schema.StringAttribute{Computed: true, Description: "id"},
			"domain":   schema.StringAttribute{Required: true, Description: "Name of the domain", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"type":     schema.StringAttribute{Required: true, Description: "The Record Type (A, AAAA, CAA, CNAME, MX, TXT, SRV, NS)", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"hostname": schema.StringAttribute{Required: true, Description: "Name (Hostname) The host name, alias, or service being defined by the record.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"value":    schema.StringAttribute{Required: true, Description: "Variable data depending on record type. For example, the value for an A record would be the IPv4 address to which the domain will be mapped. For a CAA record, it would contain the domain name of the CA being granted permission to issue certificates.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"ttl":      schema.StringAttribute{Required: true, Description: "The priority of the host (for SRV and MX records. null otherwise).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"porttype": schema.StringAttribute{Required: true, Description: "This value is the time to live for the record, in seconds. This defines the time frame that ", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"port":     schema.StringAttribute{Optional: true, Description: "The port that the service is accessible on (for SRV records only. null otherwise).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"priority": schema.StringAttribute{Optional: true, Description: "priority", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"wight":    schema.StringAttribute{Optional: true, Description: "The weight of records with the same priority (for SRV records only. null otherwise). ", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

// Import using dns record as the attribute
func (s *DnsRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource.
func (s *DnsRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "create dns record")
	// Retrieve values from plan
	var plan DnsRecordResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	dnsRecordRequest := api.DnsRecordRequest{
		Domain:   plan.Domain.ValueString(),
		Type:     plan.Type.ValueString(),
		Hostname: plan.Hostname.ValueString(),
		Value:    plan.Value.ValueString(),
		TTL:      plan.TTL.ValueString(),
		Porttype: plan.Porttype.ValueString(),
		Port:     plan.Port.ValueString(),
		Priority: plan.Priority.ValueString(),
		Wight:    plan.Wight.ValueString(),
	}
	tflog.Debug(ctx, "send create dns record request")
	dnsRecord, err := s.client.CreateDnsRecord(ctx, dnsRecordRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating dns record",
			"Could not create dns record, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = DnsRecordResourceModel{
		ID:       types.StringValue(dnsRecord.ID),
		Domain:   types.StringValue(plan.Domain.ValueString()),
		Type:     types.StringValue(plan.Type.ValueString()),
		Hostname: types.StringValue(plan.Hostname.ValueString()),
		Value:    types.StringValue(plan.Value.ValueString()),
		TTL:      types.StringValue(plan.TTL.ValueString()),
		Porttype: types.StringValue(plan.Porttype.ValueString()),
		Port:     types.StringValue(plan.Port.ValueString()),
		Priority: types.StringValue(plan.Priority.ValueString()),
		Wight:    types.StringValue(plan.Wight.ValueString()),
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish create dns record")
}

// Read resource information.
func (s *DnsRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "read dns record")

	// Get current state
	var state DnsRecordResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "send get dns record request")
	// Get refreshed dns record value from utho
	domain, err := s.client.GetDomain(ctx, state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading utho dns record",
			"Could not read utho dns record "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// get record id
	recordValues := api.Record{}
	for _, record := range domain.Records {
		if state.ID.ValueString() == record.ID {
			recordValues = record
		}
	}
	if recordValues.ID == "" {
		resp.Diagnostics.AddError(
			"Error Reading utho dns record",
			"Could not read utho dns record "+state.ID.ValueString()+": ",
		)
		return
	}

	// Overwrite items with refreshed state
	state = DnsRecordResourceModel{
		ID:       types.StringValue(recordValues.ID),
		Domain:   types.StringValue(state.Domain.ValueString()),
		Type:     types.StringValue(recordValues.Type),
		Hostname: types.StringValue(recordValues.Hostname),
		Value:    types.StringValue(recordValues.Value),
		TTL:      types.StringValue(recordValues.TTL),
		Porttype: types.StringValue(state.Porttype.ValueString()),
		Port:     types.StringValue(state.Port.ValueString()),
		Priority: types.StringValue(recordValues.Priority),
		Wight:    types.StringValue(state.Wight.ValueString()),
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish get dns record request")
}

func (s *DnsRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// updating resource is not supported
}

// Delete deletes the resource and removes the Terraform state on success.
func (s *DnsRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "delete dns record")
	// Get current state
	var state DnsRecordResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "send delete dns record request")
	// delete dns record
	err := s.client.DeleteDnsRecord(ctx, state.Domain.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing utho dns record",
			"Could not delete utho dns record "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}
