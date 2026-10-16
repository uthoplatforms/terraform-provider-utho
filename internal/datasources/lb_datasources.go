package datasources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ══════════════════════════════════════════════════════════════
// data.utho_loadbalancers
// ══════════════════════════════════════════════════════════════

type LoadBalancersDataSource struct{ client *client.Client }

type LoadBalancersModel struct {
	LoadBalancers []LoadBalancerItem `tfsdk:"loadbalancers"`
}

type LoadBalancerItem struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	IP             types.String `tfsdk:"ip"`
	DNS            types.String `tfsdk:"dns"`
	Type           types.String `tfsdk:"type"`
	DCSlug         types.String `tfsdk:"dcslug"`
	Status         types.String `tfsdk:"status"`
	EnablePublicIP types.String `tfsdk:"enable_publicip"`
	CreatedAt      types.String `tfsdk:"created_at"`
}

func NewLoadBalancersDataSource() datasource.DataSource { return &LoadBalancersDataSource{} }

func (d *LoadBalancersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancers"
}

func (d *LoadBalancersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all load balancers in your Utho account.",
		Attributes: map[string]schema.Attribute{
			"loadbalancers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":              schema.StringAttribute{Computed: true, Description: "Load balancer ID."},
						"name":            schema.StringAttribute{Computed: true, Description: "Load balancer name."},
						"ip":              schema.StringAttribute{Computed: true, Description: "Public IP address."},
						"dns":             schema.StringAttribute{Computed: true, Description: "DNS hostname."},
						"type":            schema.StringAttribute{Computed: true, Description: "Load balancer type: network or application."},
						"dcslug":          schema.StringAttribute{Computed: true, Description: "Data center slug."},
						"status":          schema.StringAttribute{Computed: true, Description: "Current status."},
						"enable_publicip": schema.StringAttribute{Computed: true, Description: "Whether public IP is enabled."},
						"created_at":      schema.StringAttribute{Computed: true, Description: "Creation timestamp."},
					},
				},
			},
		},
	}
}

func (d *LoadBalancersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", "Expected *client.Client")
		return
	}
	d.client = c
}

func (d *LoadBalancersDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	respBytes, err := d.client.Get("/loadbalancer")
	if err != nil {
		resp.Diagnostics.AddError("Error fetching load balancers", fmt.Sprintf("%s", err))
		return
	}

	var result struct {
		LoadBalancers []struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			IP             string `json:"ip"`
			DNS            string `json:"dns"`
			Type           string `json:"type"`
			DCSlug         string `json:"dcslug"`
			Status         string `json:"status"`
			EnablePublicIP string `json:"enable_publicip"`
			CreatedAt      string `json:"created_at"`
		} `json:"loadbalancers"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		resp.Diagnostics.AddError("Error parsing load balancers", fmt.Sprintf("%s", err))
		return
	}

	var lbs []LoadBalancerItem
	for _, lb := range result.LoadBalancers {
		lbs = append(lbs, LoadBalancerItem{
			ID:             types.StringValue(lb.ID),
			Name:           types.StringValue(lb.Name),
			IP:             types.StringValue(lb.IP),
			DNS:            types.StringValue(lb.DNS),
			Type:           types.StringValue(lb.Type),
			DCSlug:         types.StringValue(lb.DCSlug),
			Status:         types.StringValue(lb.Status),
			EnablePublicIP: types.StringValue(lb.EnablePublicIP),
			CreatedAt:      types.StringValue(lb.CreatedAt),
		})
	}
	if lbs == nil {
		lbs = []LoadBalancerItem{}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &LoadBalancersModel{LoadBalancers: lbs})...)
}

// ══════════════════════════════════════════════════════════════
// data.utho_ssl_certificates
// ══════════════════════════════════════════════════════════════

type SSLCertificatesDataSource struct{ client *client.Client }

type SSLCertificatesModel struct {
	Certificates []SSLCertificateItem `tfsdk:"certificates"`
}

type SSLCertificateItem struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Type               types.String `tfsdk:"type"`
	State              types.String `tfsdk:"state"`
	Status             types.String `tfsdk:"status"`
	PrimaryDomain      types.String `tfsdk:"primary_domain"`
	Issuer             types.String `tfsdk:"issuer"`
	IssuedAt           types.String `tfsdk:"issued_at"`
	ExpireAt           types.String `tfsdk:"expire_at"`
	RemainingDays      types.Int64  `tfsdk:"remaining_days"`
	KeyAlgorithm       types.String `tfsdk:"key_algorithm"`
	SignatureAlgorithm types.String `tfsdk:"signature_algorithm"`
	CreatedAt          types.String `tfsdk:"created_at"`
}

func NewSSLCertificatesDataSource() datasource.DataSource { return &SSLCertificatesDataSource{} }

func (d *SSLCertificatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssl_certificates"
}

func (d *SSLCertificatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all SSL certificates in your Utho account. Use certificate IDs when creating HTTPS frontends on load balancers.",
		Attributes: map[string]schema.Attribute{
			"certificates": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                  schema.StringAttribute{Computed: true, Description: "Certificate ID. Use this as certificate_id in utho_loadbalancer_frontend."},
						"name":                schema.StringAttribute{Computed: true, Description: "Certificate name."},
						"type":                schema.StringAttribute{Computed: true, Description: "Certificate type: custom or lets_encrypt."},
						"state":               schema.StringAttribute{Computed: true, Description: "Certificate state: verified, pending, failed."},
						"status":              schema.StringAttribute{Computed: true, Description: "Certificate status."},
						"primary_domain":      schema.StringAttribute{Computed: true, Description: "Primary domain the certificate is issued for."},
						"issuer":              schema.StringAttribute{Computed: true, Description: "Certificate issuer details."},
						"issued_at":           schema.StringAttribute{Computed: true, Description: "Issue timestamp."},
						"expire_at":           schema.StringAttribute{Computed: true, Description: "Expiry timestamp."},
						"remaining_days":      schema.Int64Attribute{Computed: true, Description: "Days remaining before expiry."},
						"key_algorithm":       schema.StringAttribute{Computed: true, Description: "Key algorithm: RSA or EC."},
						"signature_algorithm": schema.StringAttribute{Computed: true, Description: "Signature algorithm."},
						"created_at":          schema.StringAttribute{Computed: true, Description: "Creation timestamp."},
					},
				},
			},
		},
	}
}

func (d *SSLCertificatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", "Expected *client.Client")
		return
	}
	d.client = c
}

func (d *SSLCertificatesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	respBytes, err := d.client.Get("/certificates")
	if err != nil {
		resp.Diagnostics.AddError("Error fetching SSL certificates", fmt.Sprintf("%s", err))
		return
	}

	var result struct {
		Certificates []struct {
			ID                 string `json:"id"`
			Name               string `json:"name"`
			Type               string `json:"type"`
			State              string `json:"state"`
			Status             string `json:"status"`
			PrimaryDomain      string `json:"primary_domain"`
			Issuer             string `json:"issuer"`
			IssuedAt           string `json:"issued_at"`
			ExpireAt           string `json:"expire_at"`
			RemainingDays      int64  `json:"remaining_days"`
			KeyAlgorithm       string `json:"key_algorithm"`
			SignatureAlgorithm string `json:"signature_algorithm"`
			CreatedAt          string `json:"created_at"`
		} `json:"certificates"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		resp.Diagnostics.AddError("Error parsing SSL certificates", fmt.Sprintf("%s", err))
		return
	}

	var certs []SSLCertificateItem
	for _, c := range result.Certificates {
		certs = append(certs, SSLCertificateItem{
			ID:                 types.StringValue(c.ID),
			Name:               types.StringValue(c.Name),
			Type:               types.StringValue(c.Type),
			State:              types.StringValue(c.State),
			Status:             types.StringValue(c.Status),
			PrimaryDomain:      types.StringValue(c.PrimaryDomain),
			Issuer:             types.StringValue(c.Issuer),
			IssuedAt:           types.StringValue(c.IssuedAt),
			ExpireAt:           types.StringValue(c.ExpireAt),
			RemainingDays:      types.Int64Value(c.RemainingDays),
			KeyAlgorithm:       types.StringValue(c.KeyAlgorithm),
			SignatureAlgorithm: types.StringValue(c.SignatureAlgorithm),
			CreatedAt:          types.StringValue(c.CreatedAt),
		})
	}
	if certs == nil {
		certs = []SSLCertificateItem{}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &SSLCertificatesModel{Certificates: certs})...)
}
