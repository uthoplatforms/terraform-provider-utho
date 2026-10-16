package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
)

// ══════════════════════════════════════════════════════════════
// DATABASE CLUSTER — utho_database
// ══════════════════════════════════════════════════════════════

type DatabaseResource struct{ client *client.Client }

type DatabaseModel struct {
	ID            types.String `tfsdk:"id"`
	CloudID       types.String `tfsdk:"cloud_id"`
	ClusterName   types.String `tfsdk:"cluster_name"`
	DCSlug        types.String `tfsdk:"dcslug"`
	Engine        types.String `tfsdk:"engine"`
	Version       types.String `tfsdk:"version"`
	Size          types.String `tfsdk:"size"`
	NetworkType   types.String `tfsdk:"network_type"`
	VPC           types.String `tfsdk:"vpc"`
	Firewall      types.String `tfsdk:"firewall"`
	Billing       types.String `tfsdk:"billing"`
	PITREnabled   types.String `tfsdk:"pitr_enabled"`
	ReplicaCount  types.String `tfsdk:"replica_count"`
	Status        types.String `tfsdk:"status"`
	DefaultUser   types.String `tfsdk:"default_user"`
	DefaultPass   types.String `tfsdk:"default_pass"`
	DefaultDBName types.String `tfsdk:"default_dbname"`
	Port          types.String `tfsdk:"port"`
	CreatedAt     types.String `tfsdk:"created_at"`
	Host          types.String `tfsdk:"host"`
	HostPrivate   types.String `tfsdk:"host_private"`
	URI           types.String `tfsdk:"uri"`
	URIPrivate    types.String `tfsdk:"uri_private"`
}

func NewDatabaseResource() resource.Resource { return &DatabaseResource{} }

func (r *DatabaseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (r *DatabaseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage Utho Managed PostgreSQL database clusters.",
		Attributes: map[string]schema.Attribute{
			"id":             schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"cloud_id":       schema.StringAttribute{Computed: true, Description: "Primary node cloud ID.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"status":         schema.StringAttribute{Computed: true, Description: "Cluster status."},
			"default_user":   schema.StringAttribute{Computed: true, Description: "Default admin username."},
			"default_pass":   schema.StringAttribute{Computed: true, Sensitive: true, Description: "Default admin password."},
			"default_dbname": schema.StringAttribute{Computed: true, Description: "Default database name."},
			"port":           schema.StringAttribute{Computed: true, Description: "Database port."},
			"created_at":     schema.StringAttribute{Computed: true, Description: "Cluster creation timestamp."},
			"host":           schema.StringAttribute{Computed: true, Description: "Public connection hostname.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"host_private":   schema.StringAttribute{Computed: true, Description: "Private connection hostname.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"uri":            schema.StringAttribute{Computed: true, Sensitive: true, Description: "Public connection URI.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"uri_private":    schema.StringAttribute{Computed: true, Sensitive: true, Description: "Private connection URI.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"cluster_name":   schema.StringAttribute{Required: true, Description: "Cluster label/name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"dcslug":         schema.StringAttribute{Required: true, Description: "Data center slug.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"engine":         schema.StringAttribute{Required: true, Description: "Database engine. Use pg for PostgreSQL.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"version":        schema.StringAttribute{Required: true, Description: "Database version (e.g. 17).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"size":           schema.StringAttribute{Required: true, Description: "Plan ID for the database node size."},
			"network_type":   schema.StringAttribute{Required: true, Description: "Network type: public or private.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"billing":        schema.StringAttribute{Required: true, Description: "Billing cycle for the database cluster. Required. Accepted values: hourly, monthly. Use hourly for development/testing and monthly for production workloads."},
			"pitr_enabled":   schema.StringAttribute{Optional: true, Description: "Enable point-in-time recovery: 1 or 0."},
			"replica_count":  schema.StringAttribute{Optional: true, Description: "Number of replica nodes to deploy."},
			"vpc":            schema.StringAttribute{Optional: true, Description: "VPC subnet ID.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"firewall":       schema.StringAttribute{Optional: true, Description: "Security group ID to attach."},
		},
	}
}

func (r *DatabaseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DatabaseModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateDatabase(&client.DatabaseCreateRequest{
		DCSlug:         plan.DCSlug.ValueString(),
		NetworkType:    plan.NetworkType.ValueString(),
		VPC:            plan.VPC.ValueString(),
		Firewall:       plan.Firewall.ValueString(),
		ClusterLabel:   plan.ClusterName.ValueString(),
		ClusterEngine:  plan.Engine.ValueString(),
		ClusterVersion: plan.Version.ValueString(),
		Size:           plan.Size.ValueString(),
		Billing:        plan.Billing.ValueString(),
		PITREnabled:    plan.PITREnabled.ValueString(),
		ReplicaCount:   plan.ReplicaCount.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating database cluster", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	plan.Status = types.StringValue("Pending")
	plan.CreatedAt = types.StringValue("")

	// Wait for cluster to be Active (takes 5-10 minutes)
	if waitErr := r.client.WaitForDatabaseReady(id); waitErr != nil {
		if waitErr.Error() == "STILL_PROVISIONING" {
			msg := fmt.Sprintf("Cluster %q (ID: %s) was created successfully but has not become Active within 10 minutes. This is normal for larger plans or busy data centers. The cluster ID has been saved to state. Once provisioning completes in the Utho Console, run: terraform apply — Terraform will detect the cluster is now Active and continue.", plan.ClusterName.ValueString(), id)
			resp.Diagnostics.AddWarning("Database cluster still provisioning", msg)
		} else {
			resp.Diagnostics.AddWarning("Database cluster warning", fmt.Sprintf("%s", waitErr))
		}
	}

	// Fetch cloud_id, credentials and connection strings from list API
	db, err := r.client.GetDatabase(id, "")
	if err == nil && db != nil {
		plan.CloudID = types.StringValue(db.GetCloudID())
		plan.DefaultUser = types.StringValue(db.DefaultUser)
		plan.DefaultPass = types.StringValue(db.DefaultPass)
		plan.DefaultDBName = types.StringValue(db.DefaultDBName)
		plan.Port = types.StringValue(db.Port)
		plan.Status = types.StringValue(db.Status)
		plan.CreatedAt = types.StringValue(db.CreatedAt)
		plan.Host = types.StringValue(db.GetHost())
		plan.HostPrivate = types.StringValue(db.GetHostPrivate())
		plan.URI = types.StringValue(db.GetURI())
		plan.URIPrivate = types.StringValue(db.GetURIPrivate())
	} else {
		plan.CloudID = types.StringValue("")
		plan.DefaultUser = types.StringValue("")
		plan.DefaultPass = types.StringValue("")
		plan.DefaultDBName = types.StringValue("")
		plan.Port = types.StringValue("")
		plan.Host = types.StringValue("")
		plan.HostPrivate = types.StringValue("")
		plan.URI = types.StringValue("")
		plan.URIPrivate = types.StringValue("")
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *DatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DatabaseModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	db, err := r.client.GetDatabase(state.ID.ValueString(), state.CloudID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading database", fmt.Sprintf("%s", err))
		return
	}
	if db == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Status = types.StringValue(db.Status)
	state.DefaultUser = types.StringValue(db.DefaultUser)
	state.DefaultPass = types.StringValue(db.DefaultPass)
	state.DefaultDBName = types.StringValue(db.DefaultDBName)
	state.Port = types.StringValue(db.Port)
	state.CreatedAt = types.StringValue(db.CreatedAt)
	// Update cloud_id and connection strings from nodes if available
	if cloudID := db.GetCloudID(); cloudID != "" {
		state.CloudID = types.StringValue(cloudID)
	}
	if h := db.GetHost(); h != "" {
		state.Host = types.StringValue(h)
	}
	if h := db.GetHostPrivate(); h != "" {
		state.HostPrivate = types.StringValue(h)
	}
	if u := db.GetURI(); u != "" {
		state.URI = types.StringValue(u)
	}
	if u := db.GetURIPrivate(); u != "" {
		state.URIPrivate = types.StringValue(u)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DatabaseModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Resize if size changed
	if !plan.CloudID.IsNull() && plan.CloudID.ValueString() != "" {
		err := r.client.ResizeDatabase(plan.ID.ValueString(), &client.DatabaseResizeRequest{
			CloudID: plan.CloudID.ValueString(),
			Plan:    plan.Size.ValueString(),
		})
		if err != nil {
			resp.Diagnostics.AddError("Error resizing database", fmt.Sprintf("%s", err))
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *DatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DatabaseModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteDatabase(state.ID.ValueString(), state.CloudID.ValueString(), state.ClusterName.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting database", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// DATABASE — utho_database_db
// ══════════════════════════════════════════════════════════════

type DatabaseDBResource struct{ client *client.Client }

type DatabaseDBModel struct {
	ID        types.String `tfsdk:"id"`
	ClusterID types.String `tfsdk:"cluster_id"`
	Name      types.String `tfsdk:"name"`
}

func NewDatabaseDBResource() resource.Resource { return &DatabaseDBResource{} }

func (r *DatabaseDBResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_db"
}

func (r *DatabaseDBResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage databases inside a Utho database cluster.",
		Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"cluster_id": schema.StringAttribute{Required: true, Description: "Database cluster ID.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"name":       schema.StringAttribute{Required: true, Description: "Database name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *DatabaseDBResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabaseDBResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DatabaseDBModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Wait for cluster to be ready before adding database
	if err := r.client.WaitForDatabaseReady(plan.ClusterID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error waiting for database cluster", fmt.Sprintf("%s", err))
		return
	}
	if err := r.client.CreateDBDatabase(plan.ClusterID.ValueString(), plan.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error creating database", fmt.Sprintf("%s", err))
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", plan.ClusterID.ValueString(), plan.Name.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *DatabaseDBResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}
func (r *DatabaseDBResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Databases cannot be updated. Destroy and recreate.")
}
func (r *DatabaseDBResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DatabaseDBModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteDBDatabase(state.ClusterID.ValueString(), state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting database", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// DATABASE USER — utho_database_user
// ══════════════════════════════════════════════════════════════

type DatabaseUserResource struct{ client *client.Client }

type DatabaseUserModel struct {
	ID                types.String `tfsdk:"id"`
	ClusterID         types.String `tfsdk:"cluster_id"`
	Name              types.String `tfsdk:"name"`
	Password          types.String `tfsdk:"password"`
	GeneratedPassword types.String `tfsdk:"generated_password"`
}

func NewDatabaseUserResource() resource.Resource { return &DatabaseUserResource{} }

func (r *DatabaseUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_user"
}

func (r *DatabaseUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage users in a Utho database cluster.",
		Attributes: map[string]schema.Attribute{
			"id":                 schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"cluster_id":         schema.StringAttribute{Required: true, Description: "Database cluster ID.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"name":               schema.StringAttribute{Required: true, Description: "Username.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"password":           schema.StringAttribute{Optional: true, Sensitive: true, Description: "Password. Leave empty to auto-generate.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"generated_password": schema.StringAttribute{Computed: true, Sensitive: true, Description: "Auto-generated password. Save this — it will not be shown again."},
		},
	}
}

func (r *DatabaseUserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabaseUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DatabaseUserModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateDBUser(plan.ClusterID.ValueString(), &client.DatabaseUserRequest{
		Name:     plan.Name.ValueString(),
		Password: plan.Password.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", plan.ClusterID.ValueString(), plan.Name.ValueString()))
	plan.GeneratedPassword = types.StringValue(result.GeneratedPassword)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *DatabaseUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}
func (r *DatabaseUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Users cannot be updated. Destroy and recreate.")
}
func (r *DatabaseUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DatabaseUserModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteDBUser(state.ClusterID.ValueString(), state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting user", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// CONNECTION POOL — utho_database_pool
// ══════════════════════════════════════════════════════════════

type DatabasePoolResource struct{ client *client.Client }

type DatabasePoolModel struct {
	ID        types.String `tfsdk:"id"`
	ClusterID types.String `tfsdk:"cluster_id"`
	CloudID   types.String `tfsdk:"cloud_id"`
	Name      types.String `tfsdk:"name"`
	DB        types.String `tfsdk:"db"`
	User      types.String `tfsdk:"user"`
	Mode      types.String `tfsdk:"mode"`
	Size      types.Int64  `tfsdk:"size"`
}

func NewDatabasePoolResource() resource.Resource { return &DatabasePoolResource{} }

func (r *DatabasePoolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_pool"
}

func (r *DatabasePoolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage connection pools in a Utho database cluster.",
		Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"cluster_id": schema.StringAttribute{Required: true, Description: "Database cluster ID.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"cloud_id":   schema.StringAttribute{Required: true, Description: "Primary node cloud ID.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"name":       schema.StringAttribute{Required: true, Description: "Connection pool name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"db":         schema.StringAttribute{Required: true, Description: "Database name to pool connections for.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"user":       schema.StringAttribute{Required: true, Description: "Database user for the pool.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"mode":       schema.StringAttribute{Required: true, Description: "Pool mode: transaction, session, or statement.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"size":       schema.Int64Attribute{Required: true, Description: "Pool size (max connections)."},
		},
	}
}

func (r *DatabasePoolResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabasePoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DatabasePoolModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateConnectionPool(plan.ClusterID.ValueString(), plan.CloudID.ValueString(), &client.DatabasePoolRequest{
		Name: plan.Name.ValueString(), DB: plan.DB.ValueString(),
		User: plan.User.ValueString(), Mode: plan.Mode.ValueString(),
		Size: int(plan.Size.ValueInt64()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating connection pool", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *DatabasePoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}
func (r *DatabasePoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Connection pools cannot be updated. Destroy and recreate.")
}
func (r *DatabasePoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DatabasePoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteConnectionPool(state.ClusterID.ValueString(), state.CloudID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting connection pool", fmt.Sprintf("%s", err))
	}
}
