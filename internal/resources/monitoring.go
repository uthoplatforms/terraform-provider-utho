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

// ══════════════════════════════════════════════════════════════
// ALERT CONTACT — utho_alert_contact
// ══════════════════════════════════════════════════════════════

type AlertContactResource struct{ client *client.Client }

type AlertContactModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Email        types.String `tfsdk:"email"`
	MobileNumber types.String `tfsdk:"mobilenumber"`
	Status       types.String `tfsdk:"status"`
}

func NewAlertContactResource() resource.Resource { return &AlertContactResource{} }

func (r *AlertContactResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_contact"
}

func (r *AlertContactResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage alert contacts for Utho monitoring alerts.",
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":         schema.StringAttribute{Required: true, Description: "Contact name."},
			"email":        schema.StringAttribute{Required: true, Description: "Contact email address."},
			"mobilenumber": schema.StringAttribute{Required: true, Description: "Contact mobile number."},
			"status":       schema.StringAttribute{Optional: true, Description: "Contact status: 1 (active) or 0 (inactive). Default: 1."},
		},
	}
}

func (r *AlertContactResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AlertContactResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AlertContactModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	status := "1"
	if !plan.Status.IsNull() && plan.Status.ValueString() != "" {
		status = plan.Status.ValueString()
	}

	id, err := r.client.CreateAlertContact(&client.AlertContactCreateRequest{
		Name: plan.Name.ValueString(), Email: plan.Email.ValueString(),
		MobileNumber: plan.MobileNumber.ValueString(), Status: status,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating alert contact", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	plan.Status = types.StringValue(status)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *AlertContactResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AlertContactModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	contact, err := r.client.GetAlertContact(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading alert contact", fmt.Sprintf("%s", err))
		return
	}
	if contact == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(contact.Name)
	state.Email = types.StringValue(contact.Email)
	state.MobileNumber = types.StringValue(contact.MobileNumber)
	state.Status = types.StringValue(contact.Status)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AlertContactResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AlertContactModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateAlertContact(plan.ID.ValueString(), &client.AlertContactCreateRequest{
		Name: plan.Name.ValueString(), Email: plan.Email.ValueString(),
		MobileNumber: plan.MobileNumber.ValueString(), Status: plan.Status.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating alert contact", fmt.Sprintf("%s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *AlertContactResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AlertContactModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteAlertContact(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting alert contact", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// ALERT RULE — utho_alert
// ══════════════════════════════════════════════════════════════

type AlertResource struct{ client *client.Client }

type AlertModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	RefType  types.String `tfsdk:"ref_type"`
	Type     types.String `tfsdk:"type"`
	Compare  types.String `tfsdk:"compare"`
	Value    types.String `tfsdk:"value"`
	For      types.String `tfsdk:"for"`
	Contacts types.String `tfsdk:"contacts"`
	Status   types.String `tfsdk:"status"`
	RefIDs   types.String `tfsdk:"ref_ids"`
}

func NewAlertResource() resource.Resource { return &AlertResource{} }

func (r *AlertResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert"
}

func (r *AlertResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage Utho monitoring alert rules.",
		Attributes: map[string]schema.Attribute{
			"id":       schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":     schema.StringAttribute{Required: true, Description: "Alert rule name."},
			"ref_type": schema.StringAttribute{Required: true, Description: "Resource type to monitor: cloud.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"type":     schema.StringAttribute{Required: true, Description: "Metric type: cpu, ram, disk, bandwidth.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"compare":  schema.StringAttribute{Required: true, Description: "Comparison operator: above or below."},
			"value":    schema.StringAttribute{Required: true, Description: "Threshold value as percentage (e.g. 80)."},
			"for":      schema.StringAttribute{Required: true, Description: "Duration: 5m, 10m, 15m, 30m, 1h."},
			"contacts": schema.StringAttribute{Required: true, Description: "Comma-separated contact IDs to notify (e.g. 426,508)."},
			"status":   schema.StringAttribute{Required: true, Description: "Alert status: active or inactive."},
			"ref_ids":  schema.StringAttribute{Required: true, Description: "Comma-separated cloud instance IDs to monitor (e.g. 1671990,990001511).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *AlertResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AlertResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AlertModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateAlert(&client.AlertCreateRequest{
		Name: plan.Name.ValueString(), RefType: plan.RefType.ValueString(),
		Type: plan.Type.ValueString(), Compare: plan.Compare.ValueString(),
		Value: plan.Value.ValueString(), For: plan.For.ValueString(),
		Contacts: plan.Contacts.ValueString(), Status: plan.Status.ValueString(),
		RefIDs: plan.RefIDs.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating alert", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *AlertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AlertModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	alert, err := r.client.GetAlert(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading alert", fmt.Sprintf("%s", err))
		return
	}
	if alert == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(alert.Name)
	state.Type = types.StringValue(alert.Type)
	state.Compare = types.StringValue(alert.Compare)
	state.Value = types.StringValue(alert.Value)
	state.For = types.StringValue(alert.For)
	state.Status = types.StringValue(alert.Status)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AlertResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AlertModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateAlert(plan.ID.ValueString(), &client.AlertCreateRequest{
		Name: plan.Name.ValueString(), RefType: plan.RefType.ValueString(),
		Type: plan.Type.ValueString(), Compare: plan.Compare.ValueString(),
		Value: plan.Value.ValueString(), For: plan.For.ValueString(),
		Contacts: plan.Contacts.ValueString(), Status: plan.Status.ValueString(),
		RefIDs: plan.RefIDs.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating alert", fmt.Sprintf("%s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *AlertResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AlertModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteAlert(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting alert", fmt.Sprintf("%s", err))
	}
}
