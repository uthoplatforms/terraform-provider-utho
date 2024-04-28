package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/uthoplatforms/terraform-provider-utho/api"
)

var (
	_ datasource.DataSource              = &ObjectStoragePlanDataSource{}
	_ datasource.DataSourceWithConfigure = &ObjectStoragePlanDataSource{}
)

type ObjectStoragePlanDataSource struct {
	client *api.Client
}
type ObjectStoragePlanDataSourceModel struct {
	Pricing []PricingDataSourceModel `tfsdk:"pricing"`
}
type PricingDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	UUID           types.String `tfsdk:"uuid"`
	Type           types.String `tfsdk:"type"`
	Slug           types.String `tfsdk:"slug"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Disk           types.String `tfsdk:"disk"`
	RAM            types.String `tfsdk:"ram"`
	CPU            types.String `tfsdk:"cpu"`
	Bandwidth      types.String `tfsdk:"bandwidth"`
	IsFeatured     types.String `tfsdk:"is_featured"`
	DedicatedVcore types.String `tfsdk:"dedicated_vcore"`
	Price          types.String `tfsdk:"price"`
	Monthly        types.String `tfsdk:"monthly"`
}

func NewObjectStoragePlanDataSource() datasource.DataSource {
	return &ObjectStoragePlanDataSource{}
}

func (*ObjectStoragePlanDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object_storage_plan"
}

// Schema defines the schema for the data source.
func (d *ObjectStoragePlanDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"pricing": schema.ListNestedAttribute{
				Computed:    true,
				Description: "object storage plan",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":              schema.StringAttribute{Computed: true, Description: "id"},
						"uuid":            schema.StringAttribute{Computed: true, Description: "uuid"},
						"type":            schema.StringAttribute{Computed: true, Description: "type"},
						"slug":            schema.StringAttribute{Computed: true, Description: "slug"},
						"name":            schema.StringAttribute{Computed: true, Description: "name"},
						"description":     schema.StringAttribute{Computed: true, Description: "description"},
						"disk":            schema.StringAttribute{Computed: true, Description: "disk"},
						"ram":             schema.StringAttribute{Computed: true, Description: "ram"},
						"cpu":             schema.StringAttribute{Computed: true, Description: "cpu"},
						"bandwidth":       schema.StringAttribute{Computed: true, Description: "bandwidth"},
						"is_featured":     schema.StringAttribute{Computed: true, Description: "is_featured"},
						"dedicated_vcore": schema.StringAttribute{Computed: true, Description: "dedicated_vcore"},
						"price":           schema.StringAttribute{Computed: true, Description: "price"},
						"monthly":         schema.StringAttribute{Computed: true, Description: "monthly"},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ObjectStoragePlanDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ObjectStoragePlan Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

// Read refreshes the Terraform state with the latest data
func (d *ObjectStoragePlanDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read `item` data source")
	// get object_storage_plan
	objectStoragePlan, err := d.client.GetObjectStoragePlan(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list `Object Storage Plan`",
			err.Error(),
		)
		return
	}
	// Map response body to model
	state := ObjectStoragePlanDataSourceModel{}
	for _, price := range objectStoragePlan.Pricing {
		resourceState := PricingDataSourceModel{
			ID:             types.StringValue(price.ID),
			UUID:           types.StringValue(price.UUID),
			Type:           types.StringValue(price.Type),
			Slug:           types.StringValue(price.Slug),
			Name:           types.StringValue(price.Name),
			Description:    types.StringValue(price.Description),
			Disk:           types.StringValue(price.Disk),
			RAM:            types.StringValue(price.RAM),
			CPU:            types.StringValue(price.CPU),
			Bandwidth:      types.StringValue(price.Bandwidth),
			IsFeatured:     types.StringValue(price.IsFeatured),
			DedicatedVcore: types.StringValue(price.DedicatedVcore),
			Price:          types.StringValue(price.Price),
			Monthly:        types.StringValue(price.Monthly),
		}
		state.Pricing = append(state.Pricing, resourceState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Finished reading `objectStoragePlan` data source", map[string]any{"success": true})
}
