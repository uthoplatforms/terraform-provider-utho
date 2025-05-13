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
	"github.com/uthoplatforms/utho-go/utho"
)

// implement resource interfaces.
var (
	_ resource.Resource                = &CloudInstanceResource{}
	_ resource.ResourceWithConfigure   = &CloudInstanceResource{}
	_ resource.ResourceWithImportState = &CloudInstanceResource{}
)

// NewCloudInstanceResource is a helper function to simplify the provider implementation.
func NewCloudInstanceResource() resource.Resource {
	return &CloudInstanceResource{}
}

// CloudInstanceResource is the resource implementation.
type CloudInstanceResource struct {
	client utho.Client
}

type CloudInstanceResourceModel struct {
	Name         types.String `tfsdk:"name"`
	Dcslug       types.String `tfsdk:"dcslug"`
	Image        types.String `tfsdk:"image"`
	Planid       types.String `tfsdk:"planid"`
	Vpcid        types.String `tfsdk:"vpc_id"`
	RootPassword types.String `tfsdk:"root_password"`
	Firewall     types.String `tfsdk:"firewall"`
	Enablebackup types.Bool   `tfsdk:"enablebackup"`
	Backupid     types.String `tfsdk:"backupid"`
	Snapshotid   types.String `tfsdk:"snapshotid"`
	Sshkeys      types.String `tfsdk:"sshkeys"`
	Billingcycle types.String `tfsdk:"billingcycle"`
	////////////////////////
	ID                types.String  `tfsdk:"id"`
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
func (s *CloudInstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance"
}

// Configure adds the provider configured client to the data source.
func (d *CloudInstanceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(utho.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected CloudInstance Data Source Configure Type",
			fmt.Sprintf("Expected utho.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

// Schema defines the schema for the resource.
func (s *CloudInstanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":                schema.StringAttribute{Computed: true, Description: "Cloud id"},
		"name":              schema.StringAttribute{Required: true, Description: "Give a name to your cloud server eg: myweb1.server.com", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"dcslug":            schema.StringAttribute{Required: true, MarkdownDescription: "Provide Zone dcslug eg: innoida. You can find a list of available dcslug on [Utho API documentation](https://utho.com/api-docs/#api-Cloud-Servers-AVAILABLEDCZONES).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"image":             schema.StringAttribute{Required: true, Description: "Image name eg: centos-7.4-x86_64", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"planid":            schema.StringAttribute{Optional: true, MarkdownDescription: "The unique ID that identifies the type of Instance plane. You can find a list of available IDs on [Utho API documentation](https://utho.com/api-docs/#api-Cloud-Servers-GETPLANS).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"vpc_id":            schema.StringAttribute{Optional: true, MarkdownDescription: "The unique ID that identifies the VPC. You can list all VPCs id on [Utho API documentation](https://utho.com/api-docs/#api-VPC-VPCList).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"firewall":          schema.StringAttribute{Optional: true, Description: "Firewall Id", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"enablebackup":      schema.BoolAttribute{Optional: true, Description: "Please pass value on to enable weekly backups*", PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()}},
		"billingcycle":      schema.StringAttribute{Optional: true, Description: "If you required billing cycle other then hourly billing you can pass value as eg: monthly, 3month, 6month, 12month. by default its selected as hourly", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"backupid":          schema.StringAttribute{Optional: true, Description: "Provide a backupid if you have a backup in same datacenter location.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"snapshotid":        schema.StringAttribute{Optional: true, Description: "Provide a snapshot id if you have a snapshot in same datacenter location.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"sshkeys":           schema.StringAttribute{Optional: true, Description: "Provide SSH Key ids or pass multiple SSH Key ids with commans (eg: 432,331).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"root_password":     schema.StringAttribute{Computed: true, Description: "Root Password"},
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

// Import using cloud instance as the attribute
func (s *CloudInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource.
func (s *CloudInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "create cloud instance")
	// Retrieve values from plan
	var plan CloudInstanceResourceModel
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
	hostName := []utho.CloudHostname{}
	hostName = append(hostName, utho.CloudHostname{Hostname: plan.Name.ValueString()})
	// Generate API request body from plan
	cloudinstanceRequest := utho.CreateCloudInstanceParams{
		Dcslug:       plan.Dcslug.ValueString(),
		Image:        plan.Image.ValueString(),
		Planid:       plan.Planid.ValueString(),
		Vpcid:        plan.Vpcid.ValueString(),
		RootPassword: plan.RootPassword.ValueString(),
		Firewall:     plan.Firewall.ValueString(),
		Enablebackup: enableBackupMapStrBool[plan.Enablebackup.ValueBool()],
		Billingcycle: plan.Billingcycle.ValueString(),
		Backupid:     plan.Backupid.ValueString(),
		Snapshotid:   plan.Snapshotid.ValueString(),
		Sshkeys:      plan.Sshkeys.ValueString(),
		Cloud:        hostName,
	}

	tflog.Debug(ctx, "send create cloud instance request")

	cloudinstance, err := s.client.CloudInstances().Create(cloudinstanceRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cloud instance",
			"Could not create cloud instance, unexpected error: "+err.Error(),
		)
		return
	}

	getCloudInstance, err := s.client.CloudInstances().Read(cloudinstance.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading utho cloud instance",
			"Could not read utho cloud instance in create func"+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// map response value to more readable value
	enableBackupMap := map[string]bool{
		"0": false,
		"1": true,
	}

	// Map response body to schema and populate Computed attribute values
	plan.RootPassword = types.StringValue(cloudinstance.Password)
	plan.ID = types.StringValue(cloudinstance.ID)
	plan.IP = types.StringValue(cloudinstance.Ipv4)
	plan.CPU = types.StringValue(getCloudInstance.CPU)
	plan.RAM = types.StringValue(getCloudInstance.RAM)
	plan.ManagedOs = types.StringValue(getCloudInstance.ManagedOs)
	plan.ManagedFull = types.StringValue(getCloudInstance.ManagedFull)
	plan.ManagedOnetime = types.StringValue(getCloudInstance.ManagedOnetime)
	plan.PlanDisksize = types.Int64Value(int64(getCloudInstance.PlanDisksize))
	plan.Disksize = types.Int64Value(int64(getCloudInstance.Disksize))
	plan.Ha = types.StringValue(getCloudInstance.Ha)
	plan.Status = types.StringValue(getCloudInstance.Status)
	plan.Iso = types.StringValue(getCloudInstance.Iso)
	plan.Cost = types.Float64Value(getCloudInstance.Cost)
	plan.Vmcost = types.Float64Value(getCloudInstance.Vmcost)
	plan.Imagecost = types.Int64Value(int64(getCloudInstance.Imagecost))
	plan.Backupcost = types.Int64Value(int64(getCloudInstance.Backupcost))
	plan.Hourlycost = types.Float64Value(getCloudInstance.Hourlycost)
	plan.Cloudhourlycost = types.Float64Value(getCloudInstance.Cloudhourlycost)
	plan.Imagehourlycost = types.Int64Value(int64(getCloudInstance.Imagehourlycost))
	plan.Backuphourlycost = types.Int64Value(int64(getCloudInstance.Backuphourlycost))
	plan.Creditrequired = types.Float64Value(getCloudInstance.Creditrequired)
	plan.Creditreserved = types.Int64Value(int64(getCloudInstance.Creditreserved))
	plan.Nextinvoiceamount = types.Float64Value(getCloudInstance.Nextinvoiceamount)
	plan.Nextinvoicehours = types.StringValue(getCloudInstance.Nextinvoicehours)
	plan.Consolepassword = types.StringValue(getCloudInstance.Consolepassword)
	plan.Powerstatus = types.StringValue(getCloudInstance.Powerstatus)
	plan.CreatedAt = types.StringValue(getCloudInstance.CreatedAt)
	plan.UpdatedAt = types.StringValue(getCloudInstance.UpdatedAt)
	plan.Nextduedate = types.StringValue(getCloudInstance.Nextduedate)
	plan.Bandwidth = types.StringValue(getCloudInstance.Bandwidth)
	plan.BandwidthUsed = types.Int64Value(int64(getCloudInstance.BandwidthUsed))
	plan.BandwidthFree = types.Int64Value(int64(getCloudInstance.BandwidthFree))
	plan.Enablebackup = types.BoolValue(enableBackupMap[getCloudInstance.Features.Backups])
	plan.GpuAvailable = types.StringValue(getCloudInstance.GpuAvailable)

	// set state for primary types
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// set state fro complex types
	// dclocatio
	dclocationResourceModel := DclocationResourceModel{
		Location: types.StringValue(getCloudInstance.Dclocation.Location),
		Country:  types.StringValue(getCloudInstance.Dclocation.Country),
		Dc:       types.StringValue(getCloudInstance.Dclocation.Dc),
		Dccc:     types.StringValue(getCloudInstance.Dclocation.Dccc),
	}
	diags = resp.State.SetAttribute(ctx, path.Root("dclocation"), dclocationResourceModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// PrivateNetwork
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
	privateNetworkModel := make([]PrivateNetworkResourceModel, len(getCloudInstance.Networks.Private.V4))
	for i, v := range getCloudInstance.Networks.Private.V4 {
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

	// PublicNetwork
	var publicNetworkObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"ip_address": types.StringType,
		"netmask":    types.StringType,
		"gateway":    types.StringType,
		"type":       types.StringType,
		"nat":        types.BoolType,
		"primary":    types.StringType,
	}}
	PublicNetworkModel := make([]PublicNetworkResourceModel, len(getCloudInstance.Networks.Public.V4))
	for i, v := range getCloudInstance.Networks.Public.V4 {
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

	// Storages
	var storageObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":         types.StringType,
		"size":       types.Int64Type,
		"disk_used":  types.StringType,
		"disk_free":  types.StringType,
		"disk_usedp": types.StringType,
		"bus":        types.StringType,
		"type":       types.StringType,
	}}
	storageModel := make([]StoragesResourceModel, len(getCloudInstance.Storages))
	for i, v := range getCloudInstance.Storages {
		storageModel[i] = StoragesResourceModel{
			ID:        types.StringValue(v.ID),
			Size:      types.Int64Value(int64(v.Size)),
			DiskUsed:  types.StringValue(v.DiskUsed),
			DiskFree:  types.StringValue(v.DiskFree),
			DiskUsedp: types.StringValue(v.DiskUsedp),
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

	// Snapshots
	var snapshotObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":         types.StringType,
		"size":       types.StringType,
		"created_at": types.StringType,
		"note":       types.StringType,
		"name":       types.StringType,
	}}
	snapshotModel := make([]SnapshotsResourceModel, len(getCloudInstance.Snapshots))
	for i, v := range getCloudInstance.Snapshots {
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

	// Firewalls
	var firewallObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":         types.StringType,
		"created_at": types.StringType,
		"name":       types.StringType,
	}}
	firewallModel := make([]FirewallsResourceModel, len(getCloudInstance.Firewalls))
	for i, v := range getCloudInstance.Firewalls {
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

	tflog.Debug(ctx, "finish create cloud instance")
}

// Read resource information.
func (s *CloudInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "read cloud instance")

	// Get current state
	var state CloudInstanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "send get cloud instance request")
	// Get refreshed cloud instance value from utho
	cloudinstance, err := s.client.CloudInstances().Read(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading utho cloud instance",
			"Could not read utho cloud instance "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// map response value to more readable value
	enableBackupMap := map[string]bool{
		"0": false,
		"1": true,
	}

	// Overwrite items with refreshed state
	state.Name = types.StringValue(cloudinstance.Hostname)
	state.Dcslug = types.StringValue(cloudinstance.Dclocation.Dc)
	state.Image = types.StringValue(cloudinstance.Image.Image)
	state.Enablebackup = types.BoolValue(enableBackupMap[cloudinstance.Features.Backups])
	state.Billingcycle = types.StringValue(cloudinstance.Billingcycle)
	state.ID = types.StringValue(cloudinstance.ID)
	state.IP = types.StringValue(cloudinstance.IP)
	state.CPU = types.StringValue(cloudinstance.CPU)
	state.RAM = types.StringValue(cloudinstance.RAM)
	state.ManagedOs = types.StringValue(cloudinstance.ManagedOs)
	state.ManagedFull = types.StringValue(cloudinstance.ManagedFull)
	state.ManagedOnetime = types.StringValue(cloudinstance.ManagedOnetime)
	state.PlanDisksize = types.Int64Value(int64(cloudinstance.PlanDisksize))
	state.Disksize = types.Int64Value(int64(cloudinstance.Disksize))
	state.Ha = types.StringValue(cloudinstance.Ha)
	state.Status = types.StringValue(cloudinstance.Status)
	state.Iso = types.StringValue(cloudinstance.Iso)
	state.Cost = types.Float64Value(cloudinstance.Cost)
	state.Vmcost = types.Float64Value(cloudinstance.Vmcost)
	state.Imagecost = types.Int64Value(int64(cloudinstance.Imagecost))
	state.Backupcost = types.Int64Value(int64(cloudinstance.Backupcost))
	state.Hourlycost = types.Float64Value(cloudinstance.Hourlycost)
	state.Cloudhourlycost = types.Float64Value(cloudinstance.Cloudhourlycost)
	state.Imagehourlycost = types.Int64Value(int64(cloudinstance.Imagehourlycost))
	state.Backuphourlycost = types.Int64Value(int64(cloudinstance.Backuphourlycost))
	state.Creditrequired = types.Float64Value(cloudinstance.Creditrequired)
	state.Creditreserved = types.Int64Value(int64(cloudinstance.Creditreserved))
	state.Nextinvoiceamount = types.Float64Value(cloudinstance.Nextinvoiceamount)
	state.Nextinvoicehours = types.StringValue(cloudinstance.Nextinvoicehours)
	state.Consolepassword = types.StringValue(cloudinstance.Consolepassword)
	state.Powerstatus = types.StringValue(cloudinstance.Powerstatus)
	state.CreatedAt = types.StringValue(cloudinstance.CreatedAt)
	state.UpdatedAt = types.StringValue(cloudinstance.UpdatedAt)
	state.Nextduedate = types.StringValue(cloudinstance.Nextduedate)
	state.Bandwidth = types.StringValue(cloudinstance.Bandwidth)
	state.BandwidthUsed = types.Int64Value(int64(cloudinstance.BandwidthUsed))
	state.BandwidthFree = types.Int64Value(int64(cloudinstance.BandwidthFree))
	state.GpuAvailable = types.StringValue(cloudinstance.GpuAvailable)

	if !state.Snapshotid.IsNull() {
		state.Snapshotid = types.StringValue(state.Snapshotid.ValueString())
	}
	if !state.Planid.IsNull() {
		state.Planid = types.StringValue(state.Planid.ValueString())
	}
	if !state.Vpcid.IsNull() {
		state.Vpcid = types.StringValue(state.Vpcid.ValueString())
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
	// dclocatio
	dclocationResourceModel := DclocationResourceModel{
		Location: types.StringValue(cloudinstance.Dclocation.Location),
		Country:  types.StringValue(cloudinstance.Dclocation.Country),
		Dc:       types.StringValue(cloudinstance.Dclocation.Dc),
		Dccc:     types.StringValue(cloudinstance.Dclocation.Dccc),
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
	privateNetworkModel := make([]PrivateNetworkResourceModel, len(cloudinstance.Networks.Private.V4))
	for i, v := range cloudinstance.Networks.Private.V4 {
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
	PublicNetworkModel := make([]PublicNetworkResourceModel, len(cloudinstance.Networks.Public.V4))
	for i, v := range cloudinstance.Networks.Public.V4 {
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
		"bus":        types.StringType,
		"type":       types.StringType,
	}}
	storageModel := make([]StoragesResourceModel, len(cloudinstance.Storages))
	for i, v := range cloudinstance.Storages {
		storageModel[i] = StoragesResourceModel{
			ID:        types.StringValue(v.ID),
			Size:      types.Int64Value(int64(v.Size)),
			DiskUsed:  types.StringValue(v.DiskUsed),
			DiskFree:  types.StringValue(v.DiskFree),
			DiskUsedp: types.StringValue(v.DiskUsedp),
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
	snapshotModel := make([]SnapshotsResourceModel, len(cloudinstance.Snapshots))
	for i, v := range cloudinstance.Snapshots {
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
	firewallModel := make([]FirewallsResourceModel, len(cloudinstance.Firewalls))
	for i, v := range cloudinstance.Firewalls {
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

	tflog.Debug(ctx, "finish get cloud instance request")
}

func (s *CloudInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// updating resource is not supported
}

// Delete deletes the resource and removes the Terraform state on success.
func (s *CloudInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "delete cloud instance")
	// Get current state
	var state CloudInstanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "send delete cloud instance request")
	// delete cloud instance
	deleteCloudInstanceParams := utho.DeleteCloudInstanceParams{Confirm: "I am aware this action will delete data and server permanently"}
	_, err := s.client.CloudInstances().Delete(state.ID.ValueString(), deleteCloudInstanceParams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing utho cloud instance",
			"Could not delete utho cloud instance "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}
