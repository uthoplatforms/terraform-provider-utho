package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ══════════════════════════════════════════════════════════════
// LOAD BALANCER — utho_loadbalancer
// ══════════════════════════════════════════════════════════════

type LoadBalancerResource struct{ client *client.Client }

type LoadBalancerModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
	DCSlug         types.String `tfsdk:"dcslug"`
	VPC            types.String `tfsdk:"vpc"`
	EnablePublicIP types.String `tfsdk:"enable_publicip"`
	Firewall       types.String `tfsdk:"firewall"`
	IP             types.String `tfsdk:"ip"`
	DNS            types.String `tfsdk:"dns"`
	Status         types.String `tfsdk:"status"`
	CreatedAt      types.String `tfsdk:"created_at"`
}

func NewLoadBalancerResource() resource.Resource { return &LoadBalancerResource{} }

func (r *LoadBalancerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer"
}

func (r *LoadBalancerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage Utho Load Balancers.",
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"ip":              schema.StringAttribute{Computed: true, Description: "Public IP address of the load balancer."},
			"dns":             schema.StringAttribute{Computed: true, Description: "DNS hostname of the load balancer."},
			"status":          schema.StringAttribute{Computed: true, Description: "Current status of the load balancer."},
			"created_at":      schema.StringAttribute{Computed: true, Description: "Timestamp when the load balancer was created."},
			"name":            schema.StringAttribute{Required: true, Description: "Load balancer name."},
			"type":            schema.StringAttribute{Required: true, Description: "Load balancer type: network or application.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"dcslug":          schema.StringAttribute{Required: true, Description: "Data center slug.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"enable_publicip": schema.StringAttribute{Required: true, Description: "Enable public IP: true or false."},
			"vpc":             schema.StringAttribute{Optional: true, Description: "Subnet ID (not VPC ID) to deploy the load balancer in. Despite the field name, this accepts a subnet UUID. The Utho API may return null for this field after creation — Terraform state retains the configured value.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"firewall":        schema.StringAttribute{Optional: true, Description: "Security group ID to attach.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *LoadBalancerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *LoadBalancerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan LoadBalancerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateLoadBalancer(&client.LoadBalancerCreateRequest{
		Type: plan.Type.ValueString(), Name: plan.Name.ValueString(),
		DCSlug: plan.DCSlug.ValueString(), VPC: plan.VPC.ValueString(),
		EnablePublicIP: plan.EnablePublicIP.ValueString(), Firewall: plan.Firewall.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating load balancer", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	plan.IP = types.StringValue("")
	plan.DNS = types.StringValue("")
	plan.Status = types.StringValue("Active")
	plan.CreatedAt = types.StringValue("")

	// Poll for IP/DNS to be assigned (usually takes 10-20s)
	for i := 0; i < 12; i++ {
		time.Sleep(5 * time.Second)
		lb, err2 := r.client.GetLoadBalancer(id)
		if err2 == nil && lb != nil && lb.IP != "" {
			plan.IP = types.StringValue(lb.IP)
			plan.DNS = types.StringValue(lb.DNS)
			plan.Status = types.StringValue(lb.Status)
			plan.CreatedAt = types.StringValue(lb.CreatedAt)
			break
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *LoadBalancerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state LoadBalancerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lb, err := r.client.GetLoadBalancer(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading load balancer", fmt.Sprintf("%s", err))
		return
	}
	if lb == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(lb.Name)
	state.Type = types.StringValue(lb.Type)
	state.IP = types.StringValue(lb.IP)
	state.DNS = types.StringValue(lb.DNS)
	state.Status = types.StringValue(lb.Status)
	state.CreatedAt = types.StringValue(lb.CreatedAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *LoadBalancerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Load balancers cannot be updated in place. Destroy and recreate.")
}

func (r *LoadBalancerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state LoadBalancerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteLoadBalancer(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting load balancer", fmt.Sprintf("%s", err))
		return
	}
	// Wait for Utho backend to release VPC subnet attachment
	time.Sleep(15 * time.Second)
}

// ══════════════════════════════════════════════════════════════
// FRONTEND — utho_loadbalancer_frontend
// ══════════════════════════════════════════════════════════════

type LBFrontendResource struct{ client *client.Client }

type LBFrontendModel struct {
	ID             types.String `tfsdk:"id"`
	LoadBalancerID types.String `tfsdk:"loadbalancer_id"`
	Name           types.String `tfsdk:"name"`
	Algorithm      types.String `tfsdk:"algorithm"`
	Proto          types.String `tfsdk:"proto"`
	Port           types.String `tfsdk:"port"`
	Cookie         types.String `tfsdk:"cookie"`
	RedirectHTTPS  types.String `tfsdk:"redirecthttps"`
	CertificateID  types.String `tfsdk:"certificate_id"`
	CookieName     types.String `tfsdk:"cookiename"`
}

func NewLBFrontendResource() resource.Resource { return &LBFrontendResource{} }

func (r *LBFrontendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_frontend"
}

func (r *LBFrontendResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Add and manage frontends on a Utho Load Balancer.",
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"loadbalancer_id": schema.StringAttribute{Required: true, Description: "Load balancer ID.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"name":            schema.StringAttribute{Required: true, Description: "Frontend name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"algorithm":       schema.StringAttribute{Required: true, Description: "Algorithm: roundrobin, leastconn, first.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"proto":           schema.StringAttribute{Required: true, Description: "Protocol: http, https, tcp, udp.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"port":            schema.StringAttribute{Required: true, Description: "Port number.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"cookie":          schema.StringAttribute{Required: true, Description: "Enable cookie: 1 or 0.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"redirecthttps":   schema.StringAttribute{Required: true, Description: "Redirect HTTP to HTTPS: 1 or 0.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"certificate_id":  schema.StringAttribute{Required: true, Description: "SSL certificate ID. Use 0 if not applicable.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"cookiename":      schema.StringAttribute{Optional: true, Description: "Cookie name when cookie is enabled.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *LBFrontendResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *LBFrontendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan LBFrontendModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.AddLBFrontend(plan.LoadBalancerID.ValueString(), &client.LoadBalancerFrontendRequest{
		Name: plan.Name.ValueString(), Algorithm: plan.Algorithm.ValueString(),
		Proto: plan.Proto.ValueString(), Port: plan.Port.ValueString(),
		Cookie: plan.Cookie.ValueString(), RedirectHTTPS: plan.RedirectHTTPS.ValueString(),
		CertificateID: plan.CertificateID.ValueString(), CookieName: plan.CookieName.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error adding frontend", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *LBFrontendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}
func (r *LBFrontendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Frontends cannot be updated. Destroy and recreate.")
}
func (r *LBFrontendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state LBFrontendModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteLBFrontend(state.LoadBalancerID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting frontend", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// BACKEND — utho_loadbalancer_backend
// ══════════════════════════════════════════════════════════════

type LBBackendResource struct{ client *client.Client }

type LBBackendModel struct {
	ID             types.String `tfsdk:"id"`
	LoadBalancerID types.String `tfsdk:"loadbalancer_id"`
	FrontendID     types.String `tfsdk:"frontend_id"`
	BackendPort    types.String `tfsdk:"backend_port"`
	Weight         types.String `tfsdk:"weight"`
	Type           types.String `tfsdk:"type"`
	CloudID        types.String `tfsdk:"cloudid"`
	IP             types.String `tfsdk:"ip"`
}

func NewLBBackendResource() resource.Resource { return &LBBackendResource{} }

func (r *LBBackendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_backend"
}

func (r *LBBackendResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Add and manage backends on a Utho Load Balancer frontend.",
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"loadbalancer_id": schema.StringAttribute{Required: true, Description: "Load balancer ID.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"frontend_id":     schema.StringAttribute{Required: true, Description: "Frontend ID to attach this backend to.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"backend_port":    schema.StringAttribute{Required: true, Description: "Port the backend listens on.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"weight":          schema.StringAttribute{Required: true, Description: "Backend weight for load distribution (e.g. 1).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"type":            schema.StringAttribute{Required: true, Description: "Backend type: cloud or custom.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"cloudid":         schema.StringAttribute{Optional: true, Description: "Cloud instance ID. Required when type = cloud.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"ip":              schema.StringAttribute{Optional: true, Description: "Custom IP address. Required when type = custom.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *LBBackendResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *LBBackendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan LBBackendModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.AddLBBackend(plan.LoadBalancerID.ValueString(), &client.LoadBalancerBackendRequest{
		FrontendID: plan.FrontendID.ValueString(), BackendPort: plan.BackendPort.ValueString(),
		Weight: plan.Weight.ValueString(), Type: plan.Type.ValueString(),
		CloudID: plan.CloudID.ValueString(), IP: plan.IP.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error adding backend", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *LBBackendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}
func (r *LBBackendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Backends cannot be updated. Destroy and recreate.")
}
func (r *LBBackendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state LBBackendModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteLBBackend(state.LoadBalancerID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting backend", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// ACL RULE — utho_loadbalancer_acl
// ══════════════════════════════════════════════════════════════

type LBACLResource struct{ client *client.Client }

type LBACLModel struct {
	ID             types.String `tfsdk:"id"`
	LoadBalancerID types.String `tfsdk:"loadbalancer_id"`
	FrontendID     types.String `tfsdk:"frontend_id"`
	Name           types.String `tfsdk:"name"`
	ConditionType  types.String `tfsdk:"condition_type"`
	Value          types.String `tfsdk:"value"`
}

func NewLBACLResource() resource.Resource { return &LBACLResource{} }

func (r *LBACLResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_acl"
}

func (r *LBACLResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Add and manage ACL rules on a Utho Load Balancer frontend.",
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"loadbalancer_id": schema.StringAttribute{Required: true, Description: "Load balancer ID.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"frontend_id":     schema.StringAttribute{Required: true, Description: "Frontend ID to attach this ACL rule to."},
			"name":            schema.StringAttribute{Required: true, Description: "ACL rule name."},
			"condition_type":  schema.StringAttribute{Required: true, Description: "ACL condition type. Accepted values: host, url_path, url_path_regex, http_user_agent, http_referer. Note: http_method is not currently supported by the Utho API."},
			"value":           schema.StringAttribute{Required: true, Description: "ACL rule value as a JSON string. Format depends on condition_type: host: {\"host\": \"example.com\"}, url_path: {\"path\": \"/api\"}, url_path_regex: {\"regex\": \"/api/.*\"}, http_user_agent: {\"user_agent\": \"Mozilla\"}, http_referer: {\"referer\": \"example.com\"}."},
		},
	}
}

func (r *LBACLResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *LBACLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan LBACLModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.AddLBACL(plan.LoadBalancerID.ValueString(), &client.LoadBalancerACLRequest{
		FrontendID: plan.FrontendID.ValueString(), Name: plan.Name.ValueString(),
		ConditionType: plan.ConditionType.ValueString(), Value: plan.Value.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error adding ACL rule", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *LBACLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *LBACLResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan LBACLModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateLBACL(plan.LoadBalancerID.ValueString(), plan.ID.ValueString(), &client.LoadBalancerACLRequest{
		FrontendID: plan.FrontendID.ValueString(), Name: plan.Name.ValueString(),
		ConditionType: plan.ConditionType.ValueString(), Value: plan.Value.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating ACL rule", fmt.Sprintf("%s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *LBACLResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state LBACLModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteLBACL(state.LoadBalancerID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting ACL rule", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// SETTINGS — utho_loadbalancer_settings
// ══════════════════════════════════════════════════════════════

type LBSettingsResource struct{ client *client.Client }

type LBSettingsModel struct {
	ID                   types.String `tfsdk:"id"`
	LoadBalancerID       types.String `tfsdk:"loadbalancer_id"`
	TimeoutConnect       types.String `tfsdk:"timeout_connect"`
	TimeoutClient        types.String `tfsdk:"timeout_client"`
	TimeoutServer        types.String `tfsdk:"timeout_server"`
	TimeoutHTTPRequest   types.String `tfsdk:"timeout_http_request"`
	TimeoutHTTPKeepalive types.String `tfsdk:"timeout_http_keepalive"`
	TimeoutTunnel        types.String `tfsdk:"timeout_tunnel"`
	MaxConnections       types.String `tfsdk:"max_connections"`
	HTTP2                types.String `tfsdk:"http2"`
	Compression          types.String `tfsdk:"compression"`
	HSTS                 types.String `tfsdk:"hsts"`
}

func NewLBSettingsResource() resource.Resource { return &LBSettingsResource{} }

func (r *LBSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancer_settings"
}

func (r *LBSettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configure advanced settings for a Utho Load Balancer.",
		Attributes: map[string]schema.Attribute{
			"id":                     schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"loadbalancer_id":        schema.StringAttribute{Required: true, Description: "Load balancer ID.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"timeout_connect":        schema.StringAttribute{Optional: true, Description: "Connection timeout in seconds. Default: 5."},
			"timeout_client":         schema.StringAttribute{Optional: true, Description: "Client timeout in seconds. Default: 60."},
			"timeout_server":         schema.StringAttribute{Optional: true, Description: "Server timeout in seconds. Default: 60."},
			"timeout_http_request":   schema.StringAttribute{Optional: true, Description: "HTTP request timeout in seconds. Default: 10."},
			"timeout_http_keepalive": schema.StringAttribute{Optional: true, Description: "HTTP keepalive timeout in seconds. Default: 15."},
			"timeout_tunnel":         schema.StringAttribute{Optional: true, Description: "Tunnel timeout in seconds. Default: 3600."},
			"max_connections":        schema.StringAttribute{Optional: true, Description: "Maximum concurrent connections. Default: 10000."},
			"http2":                  schema.StringAttribute{Optional: true, Description: "Enable HTTP/2: 1 or 0."},
			"compression":            schema.StringAttribute{Optional: true, Description: "Enable compression: 1 or 0."},
			"hsts":                   schema.StringAttribute{Optional: true, Description: "Enable HSTS: 1 or 0."},
		},
	}
}

func (r *LBSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *LBSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan LBSettingsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateLBSettings(plan.LoadBalancerID.ValueString(), &client.LoadBalancerSettingsRequest{
		TimeoutConnect: plan.TimeoutConnect.ValueString(), TimeoutClient: plan.TimeoutClient.ValueString(),
		TimeoutServer: plan.TimeoutServer.ValueString(), TimeoutHTTPRequest: plan.TimeoutHTTPRequest.ValueString(),
		TimeoutHTTPKeepalive: plan.TimeoutHTTPKeepalive.ValueString(), TimeoutTunnel: plan.TimeoutTunnel.ValueString(),
		MaxConnections: plan.MaxConnections.ValueString(), HTTP2: plan.HTTP2.ValueString(),
		Compression: plan.Compression.ValueString(), HSTS: plan.HSTS.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating load balancer settings", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(plan.LoadBalancerID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *LBSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *LBSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan LBSettingsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateLBSettings(plan.LoadBalancerID.ValueString(), &client.LoadBalancerSettingsRequest{
		TimeoutConnect: plan.TimeoutConnect.ValueString(), TimeoutClient: plan.TimeoutClient.ValueString(),
		TimeoutServer: plan.TimeoutServer.ValueString(), TimeoutHTTPRequest: plan.TimeoutHTTPRequest.ValueString(),
		TimeoutHTTPKeepalive: plan.TimeoutHTTPKeepalive.ValueString(), TimeoutTunnel: plan.TimeoutTunnel.ValueString(),
		MaxConnections: plan.MaxConnections.ValueString(), HTTP2: plan.HTTP2.ValueString(),
		Compression: plan.Compression.ValueString(), HSTS: plan.HSTS.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating load balancer settings", fmt.Sprintf("%s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *LBSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Settings are part of the LB — nothing to delete separately
}
