package datasources

import (
	"context"
	"fmt"

	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ── Data source struct ────────────────────────────────────────────────────

type CloudsDataSource struct {
	client *client.Client
}

// ── Model ─────────────────────────────────────────────────────────────────

type CloudsDataSourceModel struct {
	Clouds []CloudDataModel `tfsdk:"clouds"`
}

type CloudDataModel struct {
	CloudID      types.String `tfsdk:"cloud_id"`
	Hostname     types.String `tfsdk:"hostname"`
	Status       types.String `tfsdk:"status"`
	PowerStatus  types.String `tfsdk:"power_status"`
	IP           types.String `tfsdk:"ip"`
	DCSlug       types.String `tfsdk:"dcslug"`
	BillingCycle types.String `tfsdk:"billingcycle"`
	CreatedAt    types.String `tfsdk:"created_at"`
}

// ── Constructor ───────────────────────────────────────────────────────────

func NewCloudsDataSource() datasource.DataSource {
	return &CloudsDataSource{}
}

// ── Metadata ──────────────────────────────────────────────────────────────

func (d *CloudsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clouds"
}

// ── Schema ────────────────────────────────────────────────────────────────

func (d *CloudsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all Utho Cloud instances in your account.",
		Attributes: map[string]schema.Attribute{
			"clouds": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cloud_id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique cloud instance ID.",
						},
						"hostname": schema.StringAttribute{
							Computed:    true,
							Description: "Hostname of the instance.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Current status.",
						},
						"power_status": schema.StringAttribute{
							Computed:    true,
							Description: "Power status.",
						},
						"ip": schema.StringAttribute{
							Computed:    true,
							Description: "Public IP address.",
						},
						"dcslug": schema.StringAttribute{
							Computed:    true,
							Description: "Data center slug.",
						},
						"billingcycle": schema.StringAttribute{
							Computed:    true,
							Description: "Billing cycle.",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "Creation timestamp.",
						},
					},
				},
			},
		},
	}
}

// ── Configure ─────────────────────────────────────────────────────────────

func (d *CloudsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	uthoClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", "Expected *client.Client")
		return
	}
	d.client = uthoClient
}

// ── READ ──────────────────────────────────────────────────────────────────

func (d *CloudsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	instances, err := d.client.ListClouds()
	if err != nil {
		resp.Diagnostics.AddError("Error listing cloud instances", fmt.Sprintf("%s", err))
		return
	}

	var state CloudsDataSourceModel
	for _, instance := range instances {
		state.Clouds = append(state.Clouds, CloudDataModel{
			CloudID:      types.StringValue(instance.CloudID),
			Hostname:     types.StringValue(instance.Hostname),
			Status:       types.StringValue(instance.Status),
			PowerStatus:  types.StringValue(instance.PowerStatus),
			IP:           types.StringValue(instance.IP),
			DCSlug:       types.StringValue(instance.DCSlug),
			BillingCycle: types.StringValue(instance.BillingCycle),
			CreatedAt:    types.StringValue(instance.CreatedAt),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
