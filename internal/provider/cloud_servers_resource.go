package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/uthoterraform/terraform-provider-utho/api"
)

// implement resource interfaces.
var (
	_ resource.Resource                = &CloudServerResource{}
	_ resource.ResourceWithConfigure   = &CloudServerResource{}
	_ resource.ResourceWithImportState = &CloudServerResource{}
)

// NewCloudServerResource is a helper function to simplify the provider implementation.
func NewCloudServerResource() resource.Resource {
	return &CloudServerResource{}
}

// CloudServerResource is the resource implementation.
type CloudServerResource struct {
	client *api.Client
}

type CloudServerResourceModel struct {
	Name         types.String `tfsdk:"name"`
	Dcslug       types.String `tfsdk:"dcslug"`
	Image        types.String `tfsdk:"image"`
	Planid       types.String `tfsdk:"planid"`
	RootPassword types.String `tfsdk:"root_password"`
	Firewall     types.String `tfsdk:"firewall"`
	Enablebackup types.Bool   `tfsdk:"enablebackup"`
	Backupid     types.String `tfsdk:"backupid"`
	Snapshotid   types.String `tfsdk:"snapshotid"`
	Sshkeys      types.String `tfsdk:"sshkeys"`
	Billingcycle types.String `tfsdk:"billingcycle"`
	////////////////////////
	Cloudid           types.String  `tfsdk:"cloudid"`
	IP                types.String  `tfsdk:"ip"`
	CPU               types.String  `tfsdk:"cpu"`
	RAM               types.String  `tfsdk:"ram"`
	ManagedOs         types.String  `tfsdk:"managed_os"`
	ManagedFull       types.String  `tfsdk:"managed_full"`
	ManagedOnetime    types.String  `tfsdk:"managed_onetime"`
	PlanDisksize      types.Int64   `tfsdk:"plan_disksize"`
	Disksize          types.Int64   `tfsdk:"disksize"`
	Ha                types.String  `tfsdk:"ha"`
	Status            types.String  `tfsdk:"status"`
	Iso               types.String  `tfsdk:"iso"`
	Cost              types.Float64 `tfsdk:"cost"`
	Vmcost            types.Float64 `tfsdk:"vmcost"`
	Imagecost         types.Int64   `tfsdk:"imagecost"`
	Backupcost        types.Int64   `tfsdk:"backupcost"`
	Hourlycost        types.Float64 `tfsdk:"hourlycost"`
	Cloudhourlycost   types.Float64 `tfsdk:"cloudhourlycost"`
	Imagehourlycost   types.Int64   `tfsdk:"imagehourlycost"`
	Backuphourlycost  types.Int64   `tfsdk:"backuphourlycost"`
	Creditrequired    types.Float64 `tfsdk:"creditrequired"`
	Creditreserved    types.Int64   `tfsdk:"creditreserved"`
	Nextinvoiceamount types.Float64 `tfsdk:"nextinvoiceamount"`
	Nextinvoicehours  types.String  `tfsdk:"nextinvoicehours"`
	Consolepassword   types.String  `tfsdk:"consolepassword"`
	Powerstatus       types.String  `tfsdk:"powerstatus"`
	CreatedAt         types.String  `tfsdk:"created_at"`
	UpdatedAt         types.String  `tfsdk:"updated_at"`
	Nextduedate       types.String  `tfsdk:"nextduedate"`
	Bandwidth         types.String  `tfsdk:"bandwidth"`
	BandwidthUsed     types.Int64   `tfsdk:"bandwidth_used"`
	BandwidthFree     types.Int64   `tfsdk:"bandwidth_free"`
	GpuAvailable      types.String  `tfsdk:"gpu_available"`
	/////////////////////////
	Features       types.Object `tfsdk:"features"`
	Dclocation     types.Object `tfsdk:"dclocation"`
	PublicNetwork  types.List   `tfsdk:"public_network"`
	PrivateNetwork types.List   `tfsdk:"private_network"`
	Storages       types.List   `tfsdk:"storages"`
	Snapshots      types.List   `tfsdk:"snapshots"`
	Firewalls      types.List   `tfsdk:"firewalls"`
}
type PublicNetworkResourceModel struct {
	IPAddress types.String `tfsdk:"ip_address"`
	Netmask   types.String `tfsdk:"netmask"`
	Gateway   types.String `tfsdk:"gateway"`
	Type      types.String `tfsdk:"type"`
	Nat       types.Bool   `tfsdk:"nat"`
	Primary   types.String `tfsdk:"primary"`
}
type PrivateNetworkResourceModel struct {
	Noip      types.Int64  `tfsdk:"noip"`
	IPAddress types.String `tfsdk:"ip_address"`
	VpcName   types.String `tfsdk:"vpc_name"`
	Network   types.String `tfsdk:"network"`
	VpcID     types.String `tfsdk:"vpc_id"`
	Netmask   types.String `tfsdk:"netmask"`
	Gateway   types.String `tfsdk:"gateway"`
	Type      types.String `tfsdk:"type"`
	Primary   types.String `tfsdk:"primary"`
}
type FeaturesResourceModel struct {
	Backups types.String `tfsdk:"backups"`
}
type DclocationResourceModel struct {
	Location types.String `tfsdk:"location"`
	Country  types.String `tfsdk:"country"`
	Dc       types.String `tfsdk:"dc"`
	Dccc     types.String `tfsdk:"dccc"`
}
type StoragesResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Size      types.Int64  `tfsdk:"size"`
	DiskUsed  types.String `tfsdk:"disk_used"`
	DiskFree  types.String `tfsdk:"disk_free"`
	DiskUsedp types.String `tfsdk:"disk_usedp"`
	CreatedAt types.String `tfsdk:"created_at"`
	Bus       types.String `tfsdk:"bus"`
	Type      types.String `tfsdk:"type"`
}
type SnapshotsResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Size      types.String `tfsdk:"size"`
	CreatedAt types.String `tfsdk:"created_at"`
	Note      types.String `tfsdk:"note"`
	Name      types.String `tfsdk:"name"`
}
type FirewallsResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	CreatedAt types.String `tfsdk:"created_at"`
}

// Metadata returns the resource type name.
func (s *CloudServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_server"
}

// Configure adds the provider configured client to the data source.
func (d *CloudServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected CloudServer Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

// Schema defines the schema for the resource.
func (s *CloudServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"name":         schema.StringAttribute{Required: true, Description: "", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"dcslug":       schema.StringAttribute{Required: true, Description: "", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"image":        schema.StringAttribute{Required: true, Description: "", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"planid":       schema.StringAttribute{Optional: true, Description: "", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"firewall":     schema.StringAttribute{Optional: true, Description: "", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"enablebackup": schema.BoolAttribute{Optional: true, Description: "", PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()}},
		"billingcycle": schema.StringAttribute{Optional: true, Description: "", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"backupid":     schema.StringAttribute{Optional: true, Description: "", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"snapshotid":   schema.StringAttribute{Optional: true, Description: "", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"sshkeys":      schema.StringAttribute{Optional: true, Description: "", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},

		"root_password":     schema.StringAttribute{Computed: true, Description: ""},
		"cloudid":           schema.StringAttribute{Computed: true, Description: "Cloudid"},
		"ip":                schema.StringAttribute{Computed: true, Description: "Ip"},
		"cpu":               schema.StringAttribute{Computed: true, Description: "Cpu"},
		"ram":               schema.StringAttribute{Computed: true, Description: "Ram"},
		"managed_os":        schema.StringAttribute{Computed: true, Description: "Managed Os"},
		"managed_full":      schema.StringAttribute{Computed: true, Description: "Managed Full"},
		"managed_onetime":   schema.StringAttribute{Computed: true, Description: "Managed Onetime"},
		"plan_disksize":     schema.Int64Attribute{Computed: true, Description: "Plan Disksize"},
		"disksize":          schema.Int64Attribute{Computed: true, Description: "Disksize"},
		"ha":                schema.StringAttribute{Computed: true, Description: "Ha"},
		"status":            schema.StringAttribute{Computed: true, Description: "Status"},
		"iso":               schema.StringAttribute{Computed: true, Description: "Iso"},
		"cost":              schema.Float64Attribute{Computed: true, Description: "Cost"},
		"vmcost":            schema.Float64Attribute{Computed: true, Description: "Vmcost"},
		"imagecost":         schema.Int64Attribute{Computed: true, Description: "Imagecost"},
		"backupcost":        schema.Int64Attribute{Computed: true, Description: "Backupcost"},
		"hourlycost":        schema.Float64Attribute{Computed: true, Description: "Hourlycost"},
		"cloudhourlycost":   schema.Float64Attribute{Computed: true, Description: "Cloudhourlycost"},
		"imagehourlycost":   schema.Int64Attribute{Computed: true, Description: "Imagehourlycost"},
		"backuphourlycost":  schema.Int64Attribute{Computed: true, Description: "Backuphourlycost"},
		"creditrequired":    schema.Float64Attribute{Computed: true, Description: "Creditrequired"},
		"creditreserved":    schema.Int64Attribute{Computed: true, Description: "Creditreserved"},
		"nextinvoiceamount": schema.Float64Attribute{Computed: true, Description: "Nextinvoiceamount"},
		"nextinvoicehours":  schema.StringAttribute{Computed: true, Description: "Nextinvoicehours"},
		"consolepassword":   schema.StringAttribute{Computed: true, Description: "Consolepassword"},
		"powerstatus":       schema.StringAttribute{Computed: true, Description: "Powerstatus"},
		"created_at":        schema.StringAttribute{Computed: true, Description: "Created At"},
		"updated_at":        schema.StringAttribute{Computed: true, Description: "Updated At"},
		"nextduedate":       schema.StringAttribute{Computed: true, Description: "Nextduedate"},
		"bandwidth":         schema.StringAttribute{Computed: true, Description: "Bandwidth"},
		"bandwidth_used":    schema.Int64Attribute{Computed: true, Description: "Bandwidth Used"},
		"bandwidth_free":    schema.Int64Attribute{Computed: true, Description: "Bandwidth Free"},
		"gpu_available":     schema.StringAttribute{Computed: true, Description: "Gpu Available"},
		"features": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Features",
			Attributes: map[string]schema.Attribute{
				"backups": schema.StringAttribute{
					Computed:    true,
					Description: "backups",
				},
			},
		},
		"dclocation": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "dclocation",
			Attributes: map[string]schema.Attribute{
				"location": schema.StringAttribute{Computed: true, Description: ""},
				"country":  schema.StringAttribute{Computed: true, Description: ""},
				"dc":       schema.StringAttribute{Computed: true, Description: ""},
				"dccc":     schema.StringAttribute{Computed: true, Description: ""},
			},
		},
		"public_network": schema.ListNestedAttribute{
			Computed:    true,
			Description: "",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"ip_address": schema.StringAttribute{Computed: true, Description: ""},
					"netmask":    schema.StringAttribute{Computed: true, Description: ""},
					"gateway":    schema.StringAttribute{Computed: true, Description: ""},
					"type":       schema.StringAttribute{Computed: true, Description: ""},
					"nat":        schema.BoolAttribute{Computed: true, Description: ""},
					"primary":    schema.StringAttribute{Computed: true, Description: ""},
				},
			},
		},
		"private_network": schema.ListNestedAttribute{
			Computed:    true,
			Description: "",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"noip":       schema.Int64Attribute{Computed: true, Description: ""},
					"ip_address": schema.StringAttribute{Computed: true},
					"vpc_name":   schema.StringAttribute{Computed: true},
					"network":    schema.StringAttribute{Computed: true},
					"vpc_id":     schema.StringAttribute{Computed: true},
					"netmask":    schema.StringAttribute{Computed: true},
					"gateway":    schema.StringAttribute{Computed: true},
					"type":       schema.StringAttribute{Computed: true},
					"primary":    schema.StringAttribute{Computed: true},
				},
			},
		},
		"storages": schema.ListNestedAttribute{
			Computed:    true,
			Description: "",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id":         schema.StringAttribute{Computed: true, Description: ""},
					"size":       schema.Int64Attribute{Computed: true, Description: ""},
					"disk_used":  schema.StringAttribute{Computed: true, Description: ""},
					"disk_free":  schema.StringAttribute{Computed: true, Description: ""},
					"disk_usedp": schema.StringAttribute{Computed: true, Description: ""},
					"created_at": schema.StringAttribute{Computed: true, Description: ""},
					"bus":        schema.StringAttribute{Computed: true, Description: ""},
					"type":       schema.StringAttribute{Computed: true, Description: ""},
				},
			},
		},
		"snapshots": schema.ListNestedAttribute{
			Computed:    true,
			Description: "",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id":         schema.StringAttribute{Computed: true},
					"size":       schema.StringAttribute{Computed: true},
					"created_at": schema.StringAttribute{Computed: true},
					"note":       schema.StringAttribute{Computed: true},
					"name":       schema.StringAttribute{Computed: true},
				},
			},
		},
		"firewalls": schema.ListNestedAttribute{
			Computed:    true,
			Description: "",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id":         schema.StringAttribute{Computed: true},
					"name":       schema.StringAttribute{Computed: true},
					"created_at": schema.StringAttribute{Computed: true},
				},
			},
		},
	},
	}
}

// Import using cloud server as the attribute
func (s *CloudServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("cloudid"), req, resp)
}

// Create a new resource.
func (s *CloudServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "create cloud server")
	// Retrieve values from plan
	var plan CloudServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// map bool to string
	enableBackupMapStrBool := map[bool]string{
		false: "false",
		true:  "true",
	}
	hostName := []api.CloudHostname{}
	hostName = append(hostName, api.CloudHostname{Hostname: plan.Name.ValueString()})
	// Generate API request body from plan
	firewallRequest := api.CreateCloudServerArgs{
		Dcslug:       plan.Dcslug.ValueString(),
		Image:        plan.Image.ValueString(),
		Planid:       plan.Planid.ValueString(),
		RootPassword: plan.RootPassword.ValueString(),
		Firewall:     plan.Firewall.ValueString(),
		Enablebackup: enableBackupMapStrBool[plan.Enablebackup.ValueBool()],
		Billingcycle: plan.Billingcycle.ValueString(),
		Backupid:     plan.Backupid.ValueString(),
		Snapshotid:   plan.Snapshotid.ValueString(),
		Sshkeys:      plan.Sshkeys.ValueString(),
		Cloud:        hostName,
	}

	tflog.Debug(ctx, "send create cloud server request")

	cloudServer, err := s.client.CreateCloudServer(ctx, firewallRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cloud server",
			"Could not create cloud server, unexpected error: "+err.Error(),
		)
		return
	}

	getCloudServer, err := s.client.GetCloudServer(ctx, plan.Cloudid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading utho cloud server",
			"Could not read utho cloud server "+plan.Cloudid.ValueString()+": "+err.Error(),
		)
		return
	}

	// map response value to more readable value
	enableBackupMap := map[string]bool{
		"0": false,
		"1": true,
	}

	// Map response body to schema and populate Computed attribute values
	plan.RootPassword = types.StringValue(cloudServer.Password)
	plan.Cloudid = types.StringValue(cloudServer.Cloudid)
	plan.IP = types.StringValue(cloudServer.Ipv4)
	plan.CPU = types.StringValue(getCloudServer.CPU)
	plan.RAM = types.StringValue(getCloudServer.RAM)
	plan.ManagedOs = types.StringValue(getCloudServer.ManagedOs)
	plan.ManagedFull = types.StringValue(getCloudServer.ManagedFull)
	plan.ManagedOnetime = types.StringValue(getCloudServer.ManagedOnetime)
	plan.PlanDisksize = types.Int64Value(int64(getCloudServer.PlanDisksize))
	plan.Disksize = types.Int64Value(int64(getCloudServer.Disksize))
	plan.Ha = types.StringValue(getCloudServer.Ha)
	plan.Status = types.StringValue(getCloudServer.Status)
	plan.Iso = types.StringValue(getCloudServer.Iso)
	plan.Cost = types.Float64Value(getCloudServer.Cost)
	plan.Vmcost = types.Float64Value(getCloudServer.Vmcost)
	plan.Imagecost = types.Int64Value(int64(getCloudServer.Imagecost))
	plan.Backupcost = types.Int64Value(int64(getCloudServer.Backupcost))
	plan.Hourlycost = types.Float64Value(getCloudServer.Hourlycost)
	plan.Cloudhourlycost = types.Float64Value(getCloudServer.Cloudhourlycost)
	plan.Imagehourlycost = types.Int64Value(int64(getCloudServer.Imagehourlycost))
	plan.Backuphourlycost = types.Int64Value(int64(getCloudServer.Backuphourlycost))
	plan.Creditrequired = types.Float64Value(getCloudServer.Creditrequired)
	plan.Creditreserved = types.Int64Value(int64(getCloudServer.Creditreserved))
	plan.Nextinvoiceamount = types.Float64Value(getCloudServer.Nextinvoiceamount)
	plan.Nextinvoicehours = types.StringValue(getCloudServer.Nextinvoicehours)
	plan.Consolepassword = types.StringValue(getCloudServer.Consolepassword)
	plan.Powerstatus = types.StringValue(getCloudServer.Powerstatus)
	plan.CreatedAt = types.StringValue(getCloudServer.CreatedAt)
	plan.UpdatedAt = types.StringValue(getCloudServer.UpdatedAt)
	plan.Nextduedate = types.StringValue(getCloudServer.Nextduedate)
	plan.Bandwidth = types.StringValue(getCloudServer.Bandwidth)
	plan.BandwidthUsed = types.Int64Value(int64(getCloudServer.BandwidthUsed))
	plan.BandwidthFree = types.Int64Value(int64(getCloudServer.BandwidthFree))
	plan.Enablebackup = types.BoolValue(enableBackupMap[getCloudServer.Features.Backups])
	plan.GpuAvailable = types.StringValue(getCloudServer.GpuAvailable)

	// set state for primary types
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// set state fro compex types
	featuresResourceModel := FeaturesResourceModel{Backups: types.StringValue(getCloudServer.Features.Backups)}
	diags = resp.State.SetAttribute(ctx, path.Root("features"), featuresResourceModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// dclocatio
	dclocationResourceModel := DclocationResourceModel{
		Location: types.StringValue(getCloudServer.Dclocation.Location),
		Country:  types.StringValue(getCloudServer.Dclocation.Country),
		Dc:       types.StringValue(getCloudServer.Dclocation.Dc),
		Dccc:     types.StringValue(getCloudServer.Dclocation.Dccc),
	}
	diags = resp.State.SetAttribute(ctx, path.Root("dclocation"), dclocationResourceModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var privateNetworkObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"noip":       types.Int64Type,
		"ip_address": types.StringType,
		"vpc_name":   types.StringType,
		"network":    types.StringType,
		"vpc_id":     types.StringType,
		"netmask":    types.StringType,
		"gateway":    types.StringType,
		"type":       types.StringType,
		"primary":    types.StringType,
	}}
	privateNetworkModel := make([]PrivateNetworkResourceModel, len(getCloudServer.Networks.Private.V4))
	for i, v := range getCloudServer.Networks.Private.V4 {
		privateNetworkModel[i] = PrivateNetworkResourceModel{
			Noip:      types.Int64Value(int64(v.Noip)),
			IPAddress: types.StringValue(v.IPAddress),
			VpcName:   types.StringValue(v.VpcName),
			Network:   types.StringValue(v.Network),
			VpcID:     types.StringValue(v.VpcID),
			Netmask:   types.StringValue(v.Netmask),
			Gateway:   types.StringValue(v.Gateway),
			Type:      types.StringValue(v.Type),
			Primary:   types.StringValue(v.Primary),
		}
	}
	privateNetworkList, diags := types.ListValueFrom(ctx, privateNetworkObjType, privateNetworkModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.SetAttribute(ctx, path.Root("private_network"), privateNetworkList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var publicNetworkObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"ip_address": types.StringType,
		"netmask":    types.StringType,
		"gateway":    types.StringType,
		"type":       types.StringType,
		"nat":        types.BoolType,
		"primary":    types.StringType,
	}}
	PublicNetworkModel := make([]PublicNetworkResourceModel, len(getCloudServer.Networks.Public.V4))
	for i, v := range getCloudServer.Networks.Public.V4 {
		PublicNetworkModel[i] = PublicNetworkResourceModel{
			IPAddress: types.StringValue(v.IPAddress),
			Netmask:   types.StringValue(v.Netmask),
			Gateway:   types.StringValue(v.Gateway),
			Type:      types.StringValue(v.Type),
			Nat:       types.BoolValue(v.Nat),
			Primary:   types.StringValue(v.Primary),
		}
	}
	publicNetworkList, diags := types.ListValueFrom(ctx, publicNetworkObjType, PublicNetworkModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.SetAttribute(ctx, path.Root("public_network"), publicNetworkList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var storageObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":         types.StringType,
		"size":       types.Int64Type,
		"disk_used":  types.StringType,
		"disk_free":  types.StringType,
		"disk_usedp": types.StringType,
		"created_at": types.StringType,
		"bus":        types.StringType,
		"type":       types.StringType,
	}}
	storageModel := make([]StoragesResourceModel, len(getCloudServer.Storages))
	for i, v := range getCloudServer.Storages {
		storageModel[i] = StoragesResourceModel{
			ID:        types.StringValue(v.ID),
			Size:      types.Int64Value(int64(v.Size)),
			DiskUsed:  types.StringValue(v.DiskUsed),
			DiskFree:  types.StringValue(v.DiskFree),
			DiskUsedp: types.StringValue(v.DiskUsedp),
			CreatedAt: types.StringValue(v.CreatedAt),
			Bus:       types.StringValue(v.Bus),
			Type:      types.StringValue(v.Type),
		}
	}
	storageList, diags := types.ListValueFrom(ctx, storageObjType, storageModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.SetAttribute(ctx, path.Root("storages"), storageList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var snapshotObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":         types.StringType,
		"size":       types.StringType,
		"created_at": types.StringType,
		"note":       types.StringType,
		"name":       types.StringType,
	}}
	snapshotModel := make([]SnapshotsResourceModel, len(getCloudServer.Snapshots))
	for i, v := range getCloudServer.Snapshots {
		snapshotModel[i] = SnapshotsResourceModel{
			ID:        types.StringValue(v.ID),
			Size:      types.StringValue(v.Size),
			CreatedAt: types.StringValue(v.CreatedAt),
			Note:      types.StringValue(v.Note),
			Name:      types.StringValue(v.Name),
		}
	}
	snapshotList, diags := types.ListValueFrom(ctx, snapshotObjType, snapshotModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.SetAttribute(ctx, path.Root("snapshots"), snapshotList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var firewallObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":         types.StringType,
		"created_at": types.StringType,
		"name":       types.StringType,
	}}
	firewallModel := make([]FirewallsResourceModel, len(getCloudServer.Firewalls))
	for i, v := range getCloudServer.Firewalls {
		firewallModel[i] = FirewallsResourceModel{
			ID:        types.StringValue(v.ID),
			CreatedAt: types.StringValue(v.CreatedAt),
			Name:      types.StringValue(v.Name),
		}
	}
	firewallList, diags := types.ListValueFrom(ctx, firewallObjType, firewallModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.SetAttribute(ctx, path.Root("firewalls"), firewallList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "finish create cloud server")
}

// Read resource information.
func (s *CloudServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "read cloud server")

	// Get current state
	var state CloudServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "send get cloud server request")
	// Get refreshed cloud server value from utho
	cloudServer, err := s.client.GetCloudServer(ctx, state.Cloudid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading utho cloud server",
			"Could not read utho cloud server "+state.Cloudid.ValueString()+": "+err.Error(),
		)
		return
	}

	// map response value to more readable value
	enableBackupMap := map[string]bool{
		"0": false,
		"1": true,
	}

	// Overwrite items with refreshed state
	state.Name = types.StringValue(cloudServer.Hostname)
	state.Dcslug = types.StringValue(cloudServer.Dclocation.Dc)
	state.Image = types.StringValue(cloudServer.Image.Image)
	state.Firewall = types.StringValue(cloudServer.Firewalls[0].ID)
	state.Enablebackup = types.BoolValue(enableBackupMap[cloudServer.Features.Backups])
	state.Billingcycle = types.StringValue(cloudServer.Billingcycle)
	state.Cloudid = types.StringValue(cloudServer.Cloudid)
	state.IP = types.StringValue(cloudServer.IP)
	state.CPU = types.StringValue(cloudServer.CPU)
	state.RAM = types.StringValue(cloudServer.RAM)
	state.ManagedOs = types.StringValue(cloudServer.ManagedOs)
	state.ManagedFull = types.StringValue(cloudServer.ManagedFull)
	state.ManagedOnetime = types.StringValue(cloudServer.ManagedOnetime)
	state.PlanDisksize = types.Int64Value(int64(cloudServer.PlanDisksize))
	state.Disksize = types.Int64Value(int64(cloudServer.Disksize))
	state.Ha = types.StringValue(cloudServer.Ha)
	state.Status = types.StringValue(cloudServer.Status)
	state.Iso = types.StringValue(cloudServer.Iso)
	state.Cost = types.Float64Value(cloudServer.Cost)
	state.Vmcost = types.Float64Value(cloudServer.Vmcost)
	state.Imagecost = types.Int64Value(int64(cloudServer.Imagecost))
	state.Backupcost = types.Int64Value(int64(cloudServer.Backupcost))
	state.Hourlycost = types.Float64Value(cloudServer.Hourlycost)
	state.Cloudhourlycost = types.Float64Value(cloudServer.Cloudhourlycost)
	state.Imagehourlycost = types.Int64Value(int64(cloudServer.Imagehourlycost))
	state.Backuphourlycost = types.Int64Value(int64(cloudServer.Backuphourlycost))
	state.Creditrequired = types.Float64Value(cloudServer.Creditrequired)
	state.Creditreserved = types.Int64Value(int64(cloudServer.Creditreserved))
	state.Nextinvoiceamount = types.Float64Value(cloudServer.Nextinvoiceamount)
	state.Nextinvoicehours = types.StringValue(cloudServer.Nextinvoicehours)
	state.Consolepassword = types.StringValue(cloudServer.Consolepassword)
	state.Powerstatus = types.StringValue(cloudServer.Powerstatus)
	state.CreatedAt = types.StringValue(cloudServer.CreatedAt)
	state.UpdatedAt = types.StringValue(cloudServer.UpdatedAt)
	state.Nextduedate = types.StringValue(cloudServer.Nextduedate)
	state.Bandwidth = types.StringValue(cloudServer.Bandwidth)
	state.BandwidthUsed = types.Int64Value(int64(cloudServer.BandwidthUsed))
	state.BandwidthFree = types.Int64Value(int64(cloudServer.BandwidthFree))
	state.GpuAvailable = types.StringValue(cloudServer.GpuAvailable)

	if !state.Snapshotid.IsNull() {
		state.Snapshotid = types.StringValue(state.Snapshotid.ValueString())
	}
	if !state.Planid.IsNull() {
		state.Planid = types.StringValue(state.Planid.ValueString())
	}
	if !state.RootPassword.IsNull() {
		state.RootPassword = types.StringValue(state.RootPassword.ValueString())
	}
	if !state.Firewall.IsNull() {
		state.Firewall = types.StringValue(state.Firewall.ValueString())
	}
	if !state.Enablebackup.IsNull() {
		state.Enablebackup = types.BoolValue(state.Enablebackup.ValueBool())
	}
	if !state.Billingcycle.IsNull() {
		state.Billingcycle = types.StringValue(state.Billingcycle.ValueString())
	}
	if !state.Backupid.IsNull() {
		state.Backupid = types.StringValue(state.Backupid.ValueString())
	}
	if !state.Sshkeys.IsNull() {
		state.Sshkeys = types.StringValue(state.Sshkeys.ValueString())
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// set state fro compex types
	featuresResourceModel := FeaturesResourceModel{Backups: types.StringValue(cloudServer.Features.Backups)}
	diags = resp.State.SetAttribute(ctx, path.Root("features"), featuresResourceModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// dclocatio
	dclocationResourceModel := DclocationResourceModel{
		Location: types.StringValue(cloudServer.Dclocation.Location),
		Country:  types.StringValue(cloudServer.Dclocation.Country),
		Dc:       types.StringValue(cloudServer.Dclocation.Dc),
		Dccc:     types.StringValue(cloudServer.Dclocation.Dccc),
	}
	diags = resp.State.SetAttribute(ctx, path.Root("dclocation"), dclocationResourceModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var privateNetworkObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"noip":       types.Int64Type,
		"ip_address": types.StringType,
		"vpc_name":   types.StringType,
		"network":    types.StringType,
		"vpc_id":     types.StringType,
		"netmask":    types.StringType,
		"gateway":    types.StringType,
		"type":       types.StringType,
		"primary":    types.StringType,
	}}
	privateNetworkModel := make([]PrivateNetworkResourceModel, len(cloudServer.Networks.Private.V4))
	for i, v := range cloudServer.Networks.Private.V4 {
		privateNetworkModel[i] = PrivateNetworkResourceModel{
			Noip:      types.Int64Value(int64(v.Noip)),
			IPAddress: types.StringValue(v.IPAddress),
			VpcName:   types.StringValue(v.VpcName),
			Network:   types.StringValue(v.Network),
			VpcID:     types.StringValue(v.VpcID),
			Netmask:   types.StringValue(v.Netmask),
			Gateway:   types.StringValue(v.Gateway),
			Type:      types.StringValue(v.Type),
			Primary:   types.StringValue(v.Primary),
		}
	}
	privateNetworkList, diags := types.ListValueFrom(ctx, privateNetworkObjType, privateNetworkModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.SetAttribute(ctx, path.Root("private_network"), privateNetworkList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var publicNetworkObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"ip_address": types.StringType,
		"netmask":    types.StringType,
		"gateway":    types.StringType,
		"type":       types.StringType,
		"nat":        types.BoolType,
		"primary":    types.StringType,
	}}
	PublicNetworkModel := make([]PublicNetworkResourceModel, len(cloudServer.Networks.Public.V4))
	for i, v := range cloudServer.Networks.Public.V4 {
		PublicNetworkModel[i] = PublicNetworkResourceModel{
			IPAddress: types.StringValue(v.IPAddress),
			Netmask:   types.StringValue(v.Netmask),
			Gateway:   types.StringValue(v.Gateway),
			Type:      types.StringValue(v.Type),
			Nat:       types.BoolValue(v.Nat),
			Primary:   types.StringValue(v.Primary),
		}
	}
	publicNetworkList, diags := types.ListValueFrom(ctx, publicNetworkObjType, PublicNetworkModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.SetAttribute(ctx, path.Root("public_network"), publicNetworkList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var storageObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":         types.StringType,
		"size":       types.Int64Type,
		"disk_used":  types.StringType,
		"disk_free":  types.StringType,
		"disk_usedp": types.StringType,
		"created_at": types.StringType,
		"bus":        types.StringType,
		"type":       types.StringType,
	}}
	storageModel := make([]StoragesResourceModel, len(cloudServer.Storages))
	for i, v := range cloudServer.Storages {
		storageModel[i] = StoragesResourceModel{
			ID:        types.StringValue(v.ID),
			Size:      types.Int64Value(int64(v.Size)),
			DiskUsed:  types.StringValue(v.DiskUsed),
			DiskFree:  types.StringValue(v.DiskFree),
			DiskUsedp: types.StringValue(v.DiskUsedp),
			CreatedAt: types.StringValue(v.CreatedAt),
			Bus:       types.StringValue(v.Bus),
			Type:      types.StringValue(v.Type),
		}
	}
	storageList, diags := types.ListValueFrom(ctx, storageObjType, storageModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.SetAttribute(ctx, path.Root("storages"), storageList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var snapshotObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":         types.StringType,
		"size":       types.StringType,
		"created_at": types.StringType,
		"note":       types.StringType,
		"name":       types.StringType,
	}}
	snapshotModel := make([]SnapshotsResourceModel, len(cloudServer.Snapshots))
	for i, v := range cloudServer.Snapshots {
		snapshotModel[i] = SnapshotsResourceModel{
			ID:        types.StringValue(v.ID),
			Size:      types.StringValue(v.Size),
			CreatedAt: types.StringValue(v.CreatedAt),
			Note:      types.StringValue(v.Note),
			Name:      types.StringValue(v.Name),
		}
	}
	snapshotList, diags := types.ListValueFrom(ctx, snapshotObjType, snapshotModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.SetAttribute(ctx, path.Root("snapshots"), snapshotList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var firewallObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":         types.StringType,
		"created_at": types.StringType,
		"name":       types.StringType,
	}}
	firewallModel := make([]FirewallsResourceModel, len(cloudServer.Firewalls))
	for i, v := range cloudServer.Firewalls {
		firewallModel[i] = FirewallsResourceModel{
			ID:        types.StringValue(v.ID),
			CreatedAt: types.StringValue(v.CreatedAt),
			Name:      types.StringValue(v.Name),
		}
	}
	firewallList, diags := types.ListValueFrom(ctx, firewallObjType, firewallModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.SetAttribute(ctx, path.Root("firewalls"), firewallList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "finish get cloud server request")
}

func (s *CloudServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// updating resource is not supported
}

// Delete deletes the resource and removes the Terraform state on success.
func (s *CloudServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "delete cloud server")
	// Get current state
	var state CloudServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "send delete cloud server request")
	// delete cloud server
	err := s.client.DeleteCloudServer(ctx, state.Cloudid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing utho cloud server",
			"Could not delete utho cloud server "+state.Cloudid.ValueString()+": "+err.Error(),
		)
		return
	}
}
