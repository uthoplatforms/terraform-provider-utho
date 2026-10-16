package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ══════════════════════════════════════════════════════════════
// DNS ZONE — utho_dns_zone
// ══════════════════════════════════════════════════════════════

type DNSZoneResource struct{ client *client.Client }

type DNSZoneModel struct {
	ID          types.String `tfsdk:"id"`
	Domain      types.String `tfsdk:"domain"`
	NSPoint     types.String `tfsdk:"nspoint"`
	RecordCount types.String `tfsdk:"record_count"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

func NewDNSZoneResource() resource.Resource { return &DNSZoneResource{} }

func (r *DNSZoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_zone"
}

func (r *DNSZoneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage Utho DNS zones.",
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"domain":       schema.StringAttribute{Required: true, Description: "Domain name (e.g. example.com).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"nspoint":      schema.StringAttribute{Computed: true, Description: "Whether nameservers are pointed to Utho."},
			"record_count": schema.StringAttribute{Computed: true, Description: "Number of DNS records in this zone."},
			"created_at":   schema.StringAttribute{Computed: true, Description: "Timestamp when the zone was created."},
		},
	}
}

func (r *DNSZoneResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DNSZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DNSZoneModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.CreateDNSZone(plan.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error creating DNS zone", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(plan.Domain.ValueString())
	plan.NSPoint = types.StringValue("NO")
	plan.RecordCount = types.StringValue("0")
	plan.CreatedAt = types.StringValue("")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *DNSZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DNSZoneModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zone, err := r.client.GetDNSZone(state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading DNS zone", fmt.Sprintf("%s", err))
		return
	}
	if zone == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.NSPoint = types.StringValue(zone.NSPoint)
	state.RecordCount = types.StringValue(zone.RecordCount)
	state.CreatedAt = types.StringValue(zone.CreatedAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DNSZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "DNS zones cannot be updated. Destroy and recreate.")
}

func (r *DNSZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DNSZoneModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteDNSZone(state.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting DNS zone", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// DNS RECORD — utho_dns_record
// ══════════════════════════════════════════════════════════════

type DNSRecordResource struct{ client *client.Client }

type DNSRecordModel struct {
	ID       types.String `tfsdk:"id"`
	Domain   types.String `tfsdk:"domain"`
	Type     types.String `tfsdk:"type"`
	Hostname types.String `tfsdk:"hostname"`
	Value    types.String `tfsdk:"value"`
	TTL      types.String `tfsdk:"ttl"`
}

func NewDNSRecordResource() resource.Resource { return &DNSRecordResource{} }

func (r *DNSRecordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

func (r *DNSRecordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage DNS records in a Utho DNS zone.",
		Attributes: map[string]schema.Attribute{
			"id":       schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"domain":   schema.StringAttribute{Required: true, Description: "Domain name the record belongs to.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"type":     schema.StringAttribute{Required: true, Description: "Record type: A, AAAA, CNAME, MX, TXT, SRV, NS."},
			"hostname": schema.StringAttribute{Required: true, Description: "Hostname for the record. Use @ for the root domain."},
			"value":    schema.StringAttribute{Required: true, Description: "Record value (IP address, target hostname, etc.)."},
			"ttl":      schema.StringAttribute{Required: true, Description: "Time to live in seconds (e.g. 3600)."},
		},
	}
}

func (r *DNSRecordResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DNSRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DNSRecordModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateDNSRecord(plan.Domain.ValueString(), &client.DNSRecordRequest{
		Type: plan.Type.ValueString(), Hostname: plan.Hostname.ValueString(),
		Value: plan.Value.ValueString(), TTL: plan.TTL.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating DNS record", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *DNSRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DNSRecordModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zone, err := r.client.GetDNSZone(state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading DNS zone", fmt.Sprintf("%s", err))
		return
	}
	if zone == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Find record by ID
	domain := state.Domain.ValueString()
	for _, record := range zone.Records {
		if record.ID == state.ID.ValueString() {
			state.Type = types.StringValue(record.Type)
			// API returns hostname with domain appended — strip it
			hostname := record.Hostname
			if hostname == domain {
				hostname = "@"
			} else {
				hostname = strings.TrimSuffix(hostname, "."+domain)
				hostname = strings.TrimSuffix(hostname, domain)
			}
			state.Hostname = types.StringValue(hostname)
			state.Value = types.StringValue(record.Value)
			state.TTL = types.StringValue(record.TTL)
			resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
			return
		}
	}

	// Record not found
	resp.State.RemoveResource(ctx)
}

func (r *DNSRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DNSRecordModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateDNSRecord(plan.Domain.ValueString(), plan.ID.ValueString(), &client.DNSRecordRequest{
		Type: plan.Type.ValueString(), Hostname: plan.Hostname.ValueString(),
		Value: plan.Value.ValueString(), TTL: plan.TTL.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating DNS record", fmt.Sprintf("%s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *DNSRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DNSRecordModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteDNSRecord(state.Domain.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting DNS record", fmt.Sprintf("%s", err))
	}
}
