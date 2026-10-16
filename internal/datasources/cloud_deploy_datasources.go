package datasources

import (
	"context"
	"fmt"

	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ══════════════════════════════════════════════════════════════
// data.utho_cloud_dczones
// ══════════════════════════════════════════════════════════════

type CloudDCZonesDataSource struct{ client *client.Client }

type CloudDCZonesModel struct {
	Zones []CloudDCZoneItem `tfsdk:"zones"`
}

type CloudDCZoneItem struct {
	ID                types.String `tfsdk:"id"`
	Slug              types.String `tfsdk:"slug"`
	City              types.String `tfsdk:"city"`
	Country           types.String `tfsdk:"country"`
	CC                types.String `tfsdk:"cc"`
	Status            types.String `tfsdk:"status"`
	DefaultCPU        types.String `tfsdk:"default_cpu"`
	EBSAvailable      types.Bool   `tfsdk:"ebs_available"`
	PublicIPAvailable types.Bool   `tfsdk:"public_ip_available"`
}

func NewCloudDCZonesDataSource() datasource.DataSource { return &CloudDCZonesDataSource{} }

func (d *CloudDCZonesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_dczones"
}

func (d *CloudDCZonesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all available data center zones for cloud instances.",
		Attributes: map[string]schema.Attribute{
			"zones": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                  schema.StringAttribute{Computed: true, Description: "DC zone ID."},
						"slug":                schema.StringAttribute{Computed: true, Description: "DC slug to use in resources (e.g. inmumbaizone2)."},
						"city":                schema.StringAttribute{Computed: true, Description: "City name."},
						"country":             schema.StringAttribute{Computed: true, Description: "Country name."},
						"cc":                  schema.StringAttribute{Computed: true, Description: "Country code (e.g. in, de, us)."},
						"status":              schema.StringAttribute{Computed: true, Description: "Zone status: active or inactive."},
						"default_cpu":         schema.StringAttribute{Computed: true, Description: "Default CPU model: amd or intel."},
						"ebs_available":       schema.BoolAttribute{Computed: true, Description: "Whether EBS volumes are available in this zone."},
						"public_ip_available": schema.BoolAttribute{Computed: true, Description: "Whether public IPs are available."},
					},
				},
			},
		},
	}
}

func (d *CloudDCZonesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudDCZonesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	data, err := d.client.GetCloudDeployData("")
	if err != nil {
		resp.Diagnostics.AddError("Error fetching DC zones", fmt.Sprintf("%s", err))
		return
	}

	var zones []CloudDCZoneItem
	for _, z := range data.DCZones {
		zones = append(zones, CloudDCZoneItem{
			ID:                types.StringValue(z.ID),
			Slug:              types.StringValue(z.Slug),
			City:              types.StringValue(z.City),
			Country:           types.StringValue(z.Country),
			CC:                types.StringValue(z.CC),
			Status:            types.StringValue(z.Status),
			DefaultCPU:        types.StringValue(z.DefaultCPU),
			EBSAvailable:      types.BoolValue(z.EBSAvailable),
			PublicIPAvailable: types.BoolValue(z.PublicIPAvailable),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &CloudDCZonesModel{Zones: zones})...)
}

// ══════════════════════════════════════════════════════════════
// data.utho_cloud_plans
// ══════════════════════════════════════════════════════════════

type CloudPlansDataSource struct{ client *client.Client }

type CloudPlansModel struct {
	DCSlug types.String    `tfsdk:"dcslug"`
	Plans  []CloudPlanItem `tfsdk:"plans"`
}

type CloudPlanItem struct {
	ID             types.String  `tfsdk:"id"`
	Name           types.String  `tfsdk:"name"`
	Slug           types.String  `tfsdk:"slug"`
	CPU            types.String  `tfsdk:"cpu"`
	RAM            types.String  `tfsdk:"ram"`
	Disk           types.String  `tfsdk:"disk"`
	Bandwidth      types.String  `tfsdk:"bandwidth"`
	DedicatedVCore types.String  `tfsdk:"dedicated_vcore"`
	Price          types.Float64 `tfsdk:"price"`
	IsAvailable    types.String  `tfsdk:"is_available"`
}

func NewCloudPlansDataSource() datasource.DataSource { return &CloudPlansDataSource{} }

func (d *CloudPlansDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_plans"
}

func (d *CloudPlansDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List available plans for cloud instances, optionally filtered by data center.",
		Attributes: map[string]schema.Attribute{
			"dcslug": schema.StringAttribute{
				Optional:    true,
				Description: "Filter plans by data center slug (e.g. inmumbaizone2). If omitted, returns plans for all DCs.",
			},
			"plans": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":              schema.StringAttribute{Computed: true, Description: "Plan ID. Use this as planid in utho_cloud."},
						"name":            schema.StringAttribute{Computed: true, Description: "Plan name."},
						"slug":            schema.StringAttribute{Computed: true, Description: "Plan type: basic, dedicated-cpu, dedicated-memory, gpu."},
						"cpu":             schema.StringAttribute{Computed: true, Description: "Number of vCPUs."},
						"ram":             schema.StringAttribute{Computed: true, Description: "RAM in MB."},
						"disk":            schema.StringAttribute{Computed: true, Description: "Disk size in GB. 0 means EBS-only plan."},
						"bandwidth":       schema.StringAttribute{Computed: true, Description: "Bandwidth in GB."},
						"dedicated_vcore": schema.StringAttribute{Computed: true, Description: "1 if dedicated vCPU, 0 if shared."},
						"price":           schema.Float64Attribute{Computed: true, Description: "Monthly price in INR."},
						"is_available":    schema.StringAttribute{Computed: true, Description: "YES or NO."},
					},
				},
			},
		},
	}
}

func (d *CloudPlansDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudPlansDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CloudPlansModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetCloudDeployData(config.DCSlug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error fetching cloud plans", fmt.Sprintf("%s", err))
		return
	}

	var plans []CloudPlanItem
	for _, p := range data.Plans {
		if p.IsAvailable != "YES" {
			continue
		}
		plans = append(plans, CloudPlanItem{
			ID:             types.StringValue(p.ID),
			Name:           types.StringValue(p.Name),
			Slug:           types.StringValue(p.Slug),
			CPU:            types.StringValue(p.CPU),
			RAM:            types.StringValue(p.RAM),
			Disk:           types.StringValue(p.Disk),
			Bandwidth:      types.StringValue(p.Bandwidth),
			DedicatedVCore: types.StringValue(p.DedicatedVCore),
			Price:          types.Float64Value(p.Price),
			IsAvailable:    types.StringValue(p.IsAvailable),
		})
	}
	if plans == nil {
		plans = []CloudPlanItem{}
	}

	config.Plans = plans
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ══════════════════════════════════════════════════════════════
// data.utho_cloud_images
// ══════════════════════════════════════════════════════════════

type CloudImagesDataSource struct{ client *client.Client }

type CloudImagesModel struct {
	Distro types.String     `tfsdk:"distro"`
	Images []CloudImageItem `tfsdk:"images"`
}

type CloudImageItem struct {
	ID           types.String `tfsdk:"id"`
	Image        types.String `tfsdk:"image"`
	Version      types.String `tfsdk:"version"`
	Distribution types.String `tfsdk:"distribution"`
	Distro       types.String `tfsdk:"distro"`
	Category     types.String `tfsdk:"category"`
}

func NewCloudImagesDataSource() datasource.DataSource { return &CloudImagesDataSource{} }

func (d *CloudImagesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_images"
}

func (d *CloudImagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List available OS images for cloud instances, optionally filtered by distro.",
		Attributes: map[string]schema.Attribute{
			"distro": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by distro: ubuntu, centos, debian, almalinux, fedora, rockylinux, windows. If omitted, returns all.",
			},
			"images": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":           schema.StringAttribute{Computed: true, Description: "Image ID."},
						"image":        schema.StringAttribute{Computed: true, Description: "Image slug. Use this as image in utho_cloud."},
						"version":      schema.StringAttribute{Computed: true, Description: "OS version."},
						"distribution": schema.StringAttribute{Computed: true, Description: "Full distribution name."},
						"distro":       schema.StringAttribute{Computed: true, Description: "Distro identifier."},
						"category":     schema.StringAttribute{Computed: true, Description: "Image category: distro or app."},
					},
				},
			},
		},
	}
}

func (d *CloudImagesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudImagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CloudImagesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetCloudDeployData("")
	if err != nil {
		resp.Diagnostics.AddError("Error fetching cloud images", fmt.Sprintf("%s", err))
		return
	}

	var images []CloudImageItem
	for _, distro := range data.Distros {
		if !config.Distro.IsNull() && config.Distro.ValueString() != "" && distro.Distro != config.Distro.ValueString() {
			continue
		}
		for _, img := range distro.Images {
			images = append(images, CloudImageItem{
				ID:           types.StringValue(img.ID),
				Image:        types.StringValue(img.Image),
				Version:      types.StringValue(img.Version),
				Distribution: types.StringValue(img.Distribution),
				Distro:       types.StringValue(img.Distro),
				Category:     types.StringValue(img.Category),
			})
		}
	}
	if images == nil {
		images = []CloudImageItem{}
	}

	config.Images = images
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ══════════════════════════════════════════════════════════════
// data.utho_cloud_snapshots
// ══════════════════════════════════════════════════════════════

type CloudSnapshotsDataSource struct{ client *client.Client }

type CloudSnapshotsModel struct {
	Snapshots []CloudSnapshotItem `tfsdk:"snapshots"`
}

type CloudSnapshotItem struct {
	ID         types.String `tfsdk:"id"`
	CloudID    types.String `tfsdk:"cloud_id"`
	Name       types.String `tfsdk:"name"`
	Image      types.String `tfsdk:"image"`
	Size       types.String `tfsdk:"size"`
	Status     types.String `tfsdk:"status"`
	Storage    types.String `tfsdk:"storage"`
	CreateDate types.String `tfsdk:"create_date"`
}

func NewCloudSnapshotsDataSource() datasource.DataSource { return &CloudSnapshotsDataSource{} }

func (d *CloudSnapshotsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_snapshots"
}

func (d *CloudSnapshotsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all snapshots in your account.",
		Attributes: map[string]schema.Attribute{
			"snapshots": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.StringAttribute{Computed: true, Description: "Snapshot ID. Use as snapshotid in utho_cloud."},
						"cloud_id":    schema.StringAttribute{Computed: true, Description: "Source cloud instance ID."},
						"name":        schema.StringAttribute{Computed: true, Description: "Snapshot name."},
						"image":       schema.StringAttribute{Computed: true, Description: "Base image the snapshot was created from."},
						"size":        schema.StringAttribute{Computed: true, Description: "Snapshot size in GB."},
						"status":      schema.StringAttribute{Computed: true, Description: "Snapshot status."},
						"storage":     schema.StringAttribute{Computed: true, Description: "Storage type: general, ebs, other."},
						"create_date": schema.StringAttribute{Computed: true, Description: "Creation timestamp."},
					},
				},
			},
		},
	}
}

func (d *CloudSnapshotsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudSnapshotsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	data, err := d.client.GetCloudDeployData("")
	if err != nil {
		resp.Diagnostics.AddError("Error fetching snapshots", fmt.Sprintf("%s", err))
		return
	}

	var snapshots []CloudSnapshotItem
	for _, s := range data.Snapshots {
		snapshots = append(snapshots, CloudSnapshotItem{
			ID:         types.StringValue(s.ID),
			CloudID:    types.StringValue(s.CloudID),
			Name:       types.StringValue(s.Name),
			Image:      types.StringValue(s.Image),
			Size:       types.StringValue(s.Size),
			Status:     types.StringValue(s.Status),
			Storage:    types.StringValue(s.Storage),
			CreateDate: types.StringValue(s.CreateDate),
		})
	}
	if snapshots == nil {
		snapshots = []CloudSnapshotItem{}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &CloudSnapshotsModel{Snapshots: snapshots})...)
}

// ══════════════════════════════════════════════════════════════
// data.utho_cloud_isos
// ══════════════════════════════════════════════════════════════

type CloudISOsDataSource struct{ client *client.Client }

type CloudISOsModel struct {
	DCSlug types.String   `tfsdk:"dcslug"`
	ISOs   []CloudISOItem `tfsdk:"isos"`
}

type CloudISOItem struct {
	Name    types.String `tfsdk:"name"`
	File    types.String `tfsdk:"file"`
	DC      types.String `tfsdk:"dc"`
	AddedAt types.String `tfsdk:"added_at"`
}

func NewCloudISOsDataSource() datasource.DataSource { return &CloudISOsDataSource{} }

func (d *CloudISOsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_isos"
}

func (d *CloudISOsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List available ISOs in your account.",
		Attributes: map[string]schema.Attribute{
			"dcslug": schema.StringAttribute{
				Optional:    true,
				Description: "Filter ISOs by data center slug.",
			},
			"isos": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":     schema.StringAttribute{Computed: true, Description: "ISO name. Use this as iso in utho_cloud."},
						"file":     schema.StringAttribute{Computed: true, Description: "ISO filename."},
						"dc":       schema.StringAttribute{Computed: true, Description: "Data center where this ISO is available."},
						"added_at": schema.StringAttribute{Computed: true, Description: "Upload timestamp."},
					},
				},
			},
		},
	}
}

func (d *CloudISOsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudISOsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CloudISOsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetCloudDeployData(config.DCSlug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error fetching ISOs", fmt.Sprintf("%s", err))
		return
	}

	var isos []CloudISOItem
	for _, iso := range data.ISOs {
		if !config.DCSlug.IsNull() && config.DCSlug.ValueString() != "" && iso.DC != config.DCSlug.ValueString() {
			continue
		}
		isos = append(isos, CloudISOItem{
			Name:    types.StringValue(iso.Name),
			File:    types.StringValue(iso.File),
			DC:      types.StringValue(iso.DC),
			AddedAt: types.StringValue(iso.AddedAt),
		})
	}
	if isos == nil {
		isos = []CloudISOItem{}
	}

	config.ISOs = isos
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ══════════════════════════════════════════════════════════════
// data.utho_vpcs
// ══════════════════════════════════════════════════════════════

type VPCsDataSource struct{ client *client.Client }

type VPCsModel struct {
	DCSlug types.String  `tfsdk:"dcslug"`
	VPCs   []VPCDataItem `tfsdk:"vpcs"`
}

type VPCSubnetDataItem struct {
	ID             types.String `tfsdk:"id"`
	UUID           types.String `tfsdk:"uuid"`
	Name           types.String `tfsdk:"name"`
	Network        types.String `tfsdk:"network"`
	Size           types.String `tfsdk:"size"`
	SubnetType     types.String `tfsdk:"subnet_type"`
	AssignPublicIP types.String `tfsdk:"assign_publicip"`
	Status         types.String `tfsdk:"status"`
	Gateway        types.String `tfsdk:"gateway"`
	DCSlug         types.String `tfsdk:"dcslug"`
}

type VPCDataItem struct {
	ID        types.String        `tfsdk:"id"`
	Name      types.String        `tfsdk:"name"`
	Network   types.String        `tfsdk:"network"`
	Size      types.String        `tfsdk:"size"`
	DCSlug    types.String        `tfsdk:"dcslug"`
	Total     types.Int64         `tfsdk:"total"`
	Available types.Int64         `tfsdk:"available"`
	Subnets   []VPCSubnetDataItem `tfsdk:"subnets"`
}

func NewVPCsDataSource() datasource.DataSource { return &VPCsDataSource{} }

func (d *VPCsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpcs"
}

func (d *VPCsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all VPCs in your account with their subnets, optionally filtered by data center.",
		Attributes: map[string]schema.Attribute{
			"dcslug": schema.StringAttribute{
				Optional:    true,
				Description: "Filter VPCs by data center slug (e.g. inmumbaizone2).",
			},
			"vpcs": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":        schema.StringAttribute{Computed: true, Description: "VPC UUID."},
						"name":      schema.StringAttribute{Computed: true, Description: "VPC name."},
						"network":   schema.StringAttribute{Computed: true, Description: "VPC network CIDR base."},
						"size":      schema.StringAttribute{Computed: true, Description: "Network prefix length."},
						"dcslug":    schema.StringAttribute{Computed: true, Description: "Data center slug."},
						"total":     schema.Int64Attribute{Computed: true, Description: "Total IPs in VPC."},
						"available": schema.Int64Attribute{Computed: true, Description: "Available IPs remaining."},
						"subnets": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id":              schema.StringAttribute{Computed: true, Description: "Subnet numeric ID. Use this as vpc in utho_cloud."},
									"uuid":            schema.StringAttribute{Computed: true, Description: "Subnet UUID."},
									"name":            schema.StringAttribute{Computed: true, Description: "Subnet name."},
									"network":         schema.StringAttribute{Computed: true, Description: "Subnet network CIDR base."},
									"size":            schema.StringAttribute{Computed: true, Description: "Subnet prefix length."},
									"subnet_type":     schema.StringAttribute{Computed: true, Description: "Subnet type: public or private."},
									"assign_publicip": schema.StringAttribute{Computed: true, Description: "1 if public IPs assigned automatically, 0 if not."},
									"status":          schema.StringAttribute{Computed: true, Description: "Subnet status."},
									"gateway":         schema.StringAttribute{Computed: true, Description: "Subnet gateway IP."},
									"dcslug":          schema.StringAttribute{Computed: true, Description: "Data center slug."},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *VPCsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VPCsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config VPCsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vpcs, err := d.client.ListVPCsFull()
	if err != nil {
		resp.Diagnostics.AddError("Error fetching VPCs", fmt.Sprintf("%s", err))
		return
	}

	var result []VPCDataItem
	for _, vpc := range vpcs {
		if !config.DCSlug.IsNull() && config.DCSlug.ValueString() != "" && vpc.DCSlug != config.DCSlug.ValueString() {
			continue
		}

		var subnets []VPCSubnetDataItem
		for _, s := range vpc.Subnets {
			subnets = append(subnets, VPCSubnetDataItem{
				ID:             types.StringValue(s.ID),
				UUID:           types.StringValue(s.UUID),
				Name:           types.StringValue(s.Name),
				Network:        types.StringValue(s.Network),
				Size:           types.StringValue(s.Size),
				SubnetType:     types.StringValue(s.SubnetType),
				AssignPublicIP: types.StringValue(s.AssignPublicIP),
				Status:         types.StringValue(s.Status),
				Gateway:        types.StringValue(s.Gateway),
				DCSlug:         types.StringValue(s.DCSlug),
			})
		}
		if subnets == nil {
			subnets = []VPCSubnetDataItem{}
		}

		result = append(result, VPCDataItem{
			ID:        types.StringValue(vpc.ID),
			Name:      types.StringValue(vpc.Name),
			Network:   types.StringValue(vpc.Network),
			Size:      types.StringValue(vpc.Size),
			DCSlug:    types.StringValue(vpc.DCSlug),
			Total:     types.Int64Value(int64(vpc.Total)),
			Available: types.Int64Value(int64(vpc.Available)),
			Subnets:   subnets,
		})
	}
	if result == nil {
		result = []VPCDataItem{}
	}

	config.VPCs = result
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
