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

type EBSDCZonesDataSource struct{ client *client.Client }

type EBSDCZoneModel struct {
	ID      types.String `tfsdk:"id"`
	Slug    types.String `tfsdk:"slug"`
	City    types.String `tfsdk:"city"`
	Country types.String `tfsdk:"country"`
	CC      types.String `tfsdk:"cc"`
	Status  types.String `tfsdk:"status"`
}

type EBSDCZonesModel struct {
	ID      types.String     `tfsdk:"id"`
	DCZones []EBSDCZoneModel `tfsdk:"dczones"`
}

func NewEBSDCZonesDataSource() datasource.DataSource { return &EBSDCZonesDataSource{} }

func (d *EBSDCZonesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ebs_dczones"
}

func (d *EBSDCZonesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve available data center zones for Utho Elastic Block Storage (EBS) deployment. EBS is currently available in Delhi (Noida), Mumbai, and Bangalore.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"dczones": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of data center zones where EBS is available.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":      schema.StringAttribute{Computed: true, Description: "Data center ID."},
						"slug":    schema.StringAttribute{Computed: true, Description: "Data center slug to use in utho_ebs dcslug field."},
						"city":    schema.StringAttribute{Computed: true, Description: "City name."},
						"country": schema.StringAttribute{Computed: true, Description: "Country name."},
						"cc":      schema.StringAttribute{Computed: true, Description: "Country code."},
						"status":  schema.StringAttribute{Computed: true, Description: "Data center status."},
					},
				},
			},
		},
	}
}

func (d *EBSDCZonesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EBSDCZonesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	respBytes, err := d.client.Get("/cloud/getdeploy")
	if err != nil {
		resp.Diagnostics.AddError("Error fetching EBS DC zones", fmt.Sprintf("%s", err))
		return
	}

	var result struct {
		DCZones []struct {
			ID           string `json:"id"`
			Slug         string `json:"slug"`
			City         string `json:"city"`
			Country      string `json:"country"`
			CC           string `json:"cc"`
			Status       string `json:"status"`
			EBSAvailable bool   `json:"ebs_available"`
		} `json:"dczones"`
	}

	if err := json.Unmarshal(respBytes, &result); err != nil {
		resp.Diagnostics.AddError("Error parsing EBS DC zones", fmt.Sprintf("%s", err))
		return
	}

	state := EBSDCZonesModel{ID: types.StringValue("ebs_dczones")}
	for _, dc := range result.DCZones {
		if dc.EBSAvailable && dc.Status == "active" {
			state.DCZones = append(state.DCZones, EBSDCZoneModel{
				ID:      types.StringValue(dc.ID),
				Slug:    types.StringValue(dc.Slug),
				City:    types.StringValue(dc.City),
				Country: types.StringValue(dc.Country),
				CC:      types.StringValue(dc.CC),
				Status:  types.StringValue(dc.Status),
			})
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
