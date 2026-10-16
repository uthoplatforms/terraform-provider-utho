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

type DatabasePlansDataSource struct{ client *client.Client }

type DatabasePlansModel struct {
	Plans []DatabasePlanItem `tfsdk:"plans"`
}

type DatabasePlanItem struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	CPU  types.String `tfsdk:"cpu"`
	RAM  types.String `tfsdk:"ram"`
	Disk types.String `tfsdk:"disk"`
	Cost types.String `tfsdk:"cost"`
}

func NewDatabasePlansDataSource() datasource.DataSource {
	return &DatabasePlansDataSource{}
}

func (d *DatabasePlansDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_plans"
}

func (d *DatabasePlansDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List available database plans.",
		Attributes: map[string]schema.Attribute{
			"plans": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":   schema.StringAttribute{Computed: true, Description: "Plan ID — use as size in utho_database."},
						"name": schema.StringAttribute{Computed: true, Description: "Plan name."},
						"cpu":  schema.StringAttribute{Computed: true, Description: "vCPU count."},
						"ram":  schema.StringAttribute{Computed: true, Description: "RAM in MB."},
						"disk": schema.StringAttribute{Computed: true, Description: "Disk in GB."},
						"cost": schema.StringAttribute{Computed: true, Description: "Monthly cost."},
					},
				},
			},
		},
	}
}

func (d *DatabasePlansDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DatabasePlansDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	respBytes, err := d.client.Get("/plans?type=db_rdbms&currency=INR")
	if err != nil {
		resp.Diagnostics.AddError("Error fetching database plans", fmt.Sprintf("%s", err))
		return
	}

	var result struct {
		Plans []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			CPU  string `json:"cpu"`
			RAM  string `json:"ram"`
			Disk string `json:"disk"`
			Cost string `json:"cost"`
		} `json:"plans"`
	}

	if err := json.Unmarshal(respBytes, &result); err != nil {
		resp.Diagnostics.AddError("Error parsing database plans", fmt.Sprintf("%s", err))
		return
	}

	var plans []DatabasePlanItem
	for _, p := range result.Plans {
		plans = append(plans, DatabasePlanItem{
			ID:   types.StringValue(p.ID),
			Name: types.StringValue(p.Name),
			CPU:  types.StringValue(p.CPU),
			RAM:  types.StringValue(p.RAM),
			Disk: types.StringValue(p.Disk),
			Cost: types.StringValue(p.Cost),
		})
	}
	if plans == nil {
		plans = []DatabasePlanItem{}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &DatabasePlansModel{Plans: plans})...)
}
