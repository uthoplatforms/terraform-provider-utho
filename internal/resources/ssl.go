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

type SSLCertificateResource struct{ client *client.Client }

type SSLCertificateModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	CertificateKey   types.String `tfsdk:"certificate_key"`
	PrivateKey       types.String `tfsdk:"private_key"`
	CertificateChain types.String `tfsdk:"certificate_chain"`
	Type             types.String `tfsdk:"type"`
	Status           types.String `tfsdk:"status"`
	PrimaryDomain    types.String `tfsdk:"primary_domain"`
	IsWildcard       types.String `tfsdk:"is_wildcard"`
	AutoRenew        types.String `tfsdk:"auto_renew"`
	CreatedAt        types.String `tfsdk:"created_at"`
	ExpireAt         types.String `tfsdk:"expire_at"`
	Issuer           types.String `tfsdk:"issuer"`
}

func NewSSLCertificateResource() resource.Resource { return &SSLCertificateResource{} }

func (r *SSLCertificateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssl_certificate"
}

func (r *SSLCertificateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Upload and manage SSL certificates on Utho Cloud.",
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name": schema.StringAttribute{Required: true, Description: "Certificate name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"type": schema.StringAttribute{Required: true, Description: "Certificate type: Custom or Managed.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"certificate_key": schema.StringAttribute{
				Required:      true,
				Sensitive:     true,
				Description:   "PEM-encoded certificate content.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"private_key": schema.StringAttribute{
				Required:      true,
				Sensitive:     true,
				Description:   "PEM-encoded private key.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"certificate_chain": schema.StringAttribute{
				Optional:      true,
				Sensitive:     true,
				Description:   "PEM-encoded certificate chain (optional).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"status":         schema.StringAttribute{Computed: true, Description: "Certificate status."},
			"primary_domain": schema.StringAttribute{Computed: true, Description: "Primary domain on the certificate."},
			"is_wildcard":    schema.StringAttribute{Computed: true, Description: "Whether this is a wildcard certificate."},
			"auto_renew":     schema.StringAttribute{Computed: true, Description: "Whether auto-renewal is enabled."},
			"created_at":     schema.StringAttribute{Computed: true, Description: "Timestamp when the certificate was created."},
			"expire_at":      schema.StringAttribute{Computed: true, Description: "Certificate expiry timestamp."},
			"issuer":         schema.StringAttribute{Computed: true, Description: "Certificate issuer."},
		},
	}
}

func (r *SSLCertificateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SSLCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SSLCertificateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateSSLCertificate(&client.SSLCertificateCreateRequest{
		Name:             plan.Name.ValueString(),
		CertificateKey:   plan.CertificateKey.ValueString(),
		PrivateKey:       plan.PrivateKey.ValueString(),
		CertificateChain: plan.CertificateChain.ValueString(),
		Type:             plan.Type.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating SSL certificate", fmt.Sprintf("%s", err))
		return
	}

	// If API doesn't return ID find by name in list
	if id == "" {
		certs, err := r.client.ListSSLCertificates()
		if err != nil {
			resp.Diagnostics.AddError("Error listing SSL certificates", fmt.Sprintf("%s", err))
			return
		}
		for _, cert := range certs {
			if cert.Name == plan.Name.ValueString() {
				id = cert.ID
			}
		}
	}

	plan.ID = types.StringValue(id)
	plan.Status = types.StringValue("")
	plan.PrimaryDomain = types.StringValue("")
	plan.IsWildcard = types.StringValue("0")
	plan.AutoRenew = types.StringValue("0")
	plan.CreatedAt = types.StringValue("")
	plan.ExpireAt = types.StringValue("")
	plan.Issuer = types.StringValue("")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *SSLCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SSLCertificateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cert, err := r.client.GetSSLCertificate(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading SSL certificate", fmt.Sprintf("%s", err))
		return
	}
	if cert == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(cert.Name)
	state.Status = types.StringValue(cert.Status)
	state.PrimaryDomain = types.StringValue(cert.PrimaryDomain)
	state.IsWildcard = types.StringValue(cert.IsWildcard)
	state.AutoRenew = types.StringValue(cert.AutoRenew)
	state.CreatedAt = types.StringValue(cert.CreatedAt)
	state.ExpireAt = types.StringValue(cert.ExpireAt)
	state.Issuer = types.StringValue(cert.Issuer)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SSLCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "SSL certificates cannot be updated. Destroy and recreate.")
}

func (r *SSLCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SSLCertificateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteSSLCertificate(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting SSL certificate", fmt.Sprintf("%s", err))
	}
}
