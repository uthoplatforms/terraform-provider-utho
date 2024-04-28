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
	_ datasource.DataSource              = &ImagesDataSource{}
	_ datasource.DataSourceWithConfigure = &ImagesDataSource{}
)

type ImagesDataSource struct {
	client *api.Client
}
type ImagesDataSourceModel struct {
	Images []ImageDataSourceModel `tfsdk:"images"`
}
type ImageDataSourceModel struct {
	Distro       types.String `tfsdk:"distro"`
	Distribution types.String `tfsdk:"distribution"`
	Version      types.String `tfsdk:"version"`
	Image        types.String `tfsdk:"image"`
	Cost         types.Int64  `tfsdk:"cost"`
}

func NewImagesDataSource() datasource.DataSource {
	return &ImagesDataSource{}
}

func (*ImagesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_images"
}

// Schema defines the schema for the data source.
func (d *ImagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"images": schema.ListNestedAttribute{
				Computed:    true,
				Description: "OS images",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"distro":       schema.StringAttribute{Computed: true, Description: "distro"},
						"distribution": schema.StringAttribute{Computed: true, Description: "distribution"},
						"version":      schema.StringAttribute{Computed: true, Description: "version"},
						"image":        schema.StringAttribute{Computed: true, Description: "image"},
						"cost":         schema.Int64Attribute{Computed: true, Description: "cost"},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ImagesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Images Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

// Read refreshes the Terraform state with the latest data
func (d *ImagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read `item` data source")
	// get images
	images, err := d.client.GetImages(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list `images`",
			err.Error(),
		)
		return
	}
	// Map response body to model
	state := ImagesDataSourceModel{}
	for _, image := range images.Images {
		resourceState := ImageDataSourceModel{
			Distro:       types.StringValue(image.Distro),
			Distribution: types.StringValue(image.Distribution),
			Version:      types.StringValue(image.Version),
			Image:        types.StringValue(image.Image),
			Cost:         types.Int64Value(int64(image.Cost)),
		}
		state.Images = append(state.Images, resourceState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Finished reading `images` data source", map[string]any{"success": true})
}
