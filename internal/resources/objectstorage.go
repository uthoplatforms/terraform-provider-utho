package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
)

// ══════════════════════════════════════════════════════════════
// OBJECT STORAGE BUCKET — utho_object_storage
// ══════════════════════════════════════════════════════════════

type ObjectStorageResource struct{ client *client.Client }

type ObjectStorageModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	DCSlug         types.String `tfsdk:"dcslug"`
	Access         types.String `tfsdk:"access"`
	VersionEnabled types.Bool   `tfsdk:"version_enabled"`
	Status         types.String `tfsdk:"status"`
	AccessKey      types.String `tfsdk:"access_key"`
	SecretKey      types.String `tfsdk:"secret_key"`
	CreatedAt      types.String `tfsdk:"created_at"`
	PlanGB         types.Int64  `tfsdk:"plan_gb"`
	UsedGB         types.String `tfsdk:"used_gb"`
	BillingCycle   types.String `tfsdk:"billingcycle"`
}

func NewObjectStorageResource() resource.Resource { return &ObjectStorageResource{} }

func (r *ObjectStorageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object_storage"
}

func (r *ObjectStorageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage Utho Object Storage buckets (S3-compatible).",
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":            schema.StringAttribute{Required: true, Description: "Bucket name. Must be globally unique, lowercase only, and may contain letters, numbers, and hyphens. Uppercase letters and underscores are not allowed.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"dcslug":          schema.StringAttribute{Required: true, Description: "Data center slug. Currently only innoida is supported.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"access":          schema.StringAttribute{Optional: true, Computed: true, Description: "Bucket access policy: private, public, or upload. Default: private."},
			"version_enabled": schema.BoolAttribute{Optional: true, Computed: true, Description: "Enable object versioning. Default: false.", PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()}},
			"billingcycle":    schema.StringAttribute{Optional: true, Description: "Billing cycle. Default: monthly.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"status":          schema.StringAttribute{Computed: true, Description: "Bucket status."},
			"access_key":      schema.StringAttribute{Computed: true, Sensitive: true, Description: "S3-compatible access key for this bucket."},
			"secret_key":      schema.StringAttribute{Computed: true, Sensitive: true, Description: "S3-compatible secret key for this bucket."},
			"created_at":      schema.StringAttribute{Computed: true, Description: "Timestamp when the bucket was created."},
			"plan_gb":         schema.Int64Attribute{Computed: true, Description: "Plan storage in GB (100 GB minimum)."},
			"used_gb":         schema.StringAttribute{Computed: true, Description: "Current storage usage in GB."},
		},
	}
}

func (r *ObjectStorageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ObjectStorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ObjectStorageModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	billingcycle := "monthly"
	if !plan.BillingCycle.IsNull() && plan.BillingCycle.ValueString() != "" {
		billingcycle = plan.BillingCycle.ValueString()
	}

	_, err := r.client.CreateBucket(&client.BucketCreateRequest{
		DCSlug: plan.DCSlug.ValueString(), Name: plan.Name.ValueString(),
		Size: "obs-flat-100", PlanID: "obs-flat-100",
		Cycle: billingcycle, Billing: billingcycle,
		BillingCycle: billingcycle, Price: "500",
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating bucket", fmt.Sprintf("%s", err))
		return
	}

	// Set access policy if specified
	if !plan.Access.IsNull() && plan.Access.ValueString() != "" && plan.Access.ValueString() != "private" {
		if err := r.client.UpdateBucketPolicy(plan.DCSlug.ValueString(), plan.Name.ValueString(), plan.Access.ValueString()); err != nil {
			resp.Diagnostics.AddError("Error setting bucket policy", fmt.Sprintf("%s", err))
			return
		}
	}

	// Enable versioning if requested
	if !plan.VersionEnabled.IsNull() && plan.VersionEnabled.ValueBool() {
		if err := r.client.EnableBucketVersioning(plan.DCSlug.ValueString(), plan.Name.ValueString()); err != nil {
			resp.Diagnostics.AddError("Error enabling versioning", fmt.Sprintf("%s", err))
			return
		}
	}

	// Fetch bucket details
	bucket, err := r.client.GetBucket(plan.DCSlug.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading bucket after create", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(plan.Name.ValueString())
	plan.Status = types.StringValue(bucket.Status)
	plan.AccessKey = types.StringValue(bucket.AccessKey)
	plan.SecretKey = types.StringValue(bucket.SecretKey)
	plan.CreatedAt = types.StringValue(bucket.CreatedAt)
	plan.PlanGB = types.Int64Value(int64(bucket.PlanGB))
	plan.UsedGB = types.StringValue(fmt.Sprintf("%v", bucket.UsedGB))
	plan.Access = types.StringValue(bucket.Access)
	// List API doesn't return version_enabled reliably
	// Keep value from plan if set, default to false
	if plan.VersionEnabled.IsNull() || plan.VersionEnabled.IsUnknown() {
		if bucket.VersionEnabled {
			plan.VersionEnabled = types.BoolValue(true)
		} else {
			plan.VersionEnabled = types.BoolValue(false)
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ObjectStorageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ObjectStorageModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucket, err := r.client.GetBucket(state.DCSlug.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading bucket", fmt.Sprintf("%s", err))
		return
	}
	if bucket == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Status = types.StringValue(bucket.Status)
	state.Access = types.StringValue(bucket.Access)
	// List API doesn't reliably return version_enabled — keep state value if API returns false
	if bucket.VersionEnabled {
		state.VersionEnabled = types.BoolValue(true)
	} else if state.VersionEnabled.IsNull() || state.VersionEnabled.IsUnknown() {
		state.VersionEnabled = types.BoolValue(false)
	}
	if bucket.AccessKey != "" {
		state.AccessKey = types.StringValue(bucket.AccessKey)
	}
	if bucket.SecretKey != "" {
		state.SecretKey = types.StringValue(bucket.SecretKey)
	}
	state.PlanGB = types.Int64Value(int64(bucket.PlanGB))
	state.UsedGB = types.StringValue(fmt.Sprintf("%v", bucket.UsedGB))
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ObjectStorageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ObjectStorageModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update access policy if changed
	if plan.Access.ValueString() != state.Access.ValueString() {
		if err := r.client.UpdateBucketPolicy(plan.DCSlug.ValueString(), plan.Name.ValueString(), plan.Access.ValueString()); err != nil {
			resp.Diagnostics.AddError("Error updating bucket policy", fmt.Sprintf("%s", err))
			return
		}
	}

	// Update versioning if changed
	if plan.VersionEnabled.ValueBool() != state.VersionEnabled.ValueBool() {
		if plan.VersionEnabled.ValueBool() {
			if err := r.client.EnableBucketVersioning(plan.DCSlug.ValueString(), plan.Name.ValueString()); err != nil {
				resp.Diagnostics.AddError("Error enabling versioning", fmt.Sprintf("%s", err))
				return
			}
		} else {
			if err := r.client.DisableBucketVersioning(plan.DCSlug.ValueString(), plan.Name.ValueString()); err != nil {
				resp.Diagnostics.AddError("Error disabling versioning", fmt.Sprintf("%s", err))
				return
			}
		}
	}

	// Copy known values from state into plan before saving
	plan.ID = state.ID
	plan.CreatedAt = state.CreatedAt
	plan.Status = state.Status
	plan.AccessKey = state.AccessKey
	plan.SecretKey = state.SecretKey
	plan.PlanGB = state.PlanGB
	plan.UsedGB = state.UsedGB
	if plan.VersionEnabled.IsUnknown() || plan.VersionEnabled.IsNull() {
		plan.VersionEnabled = state.VersionEnabled
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ObjectStorageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ObjectStorageModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteBucket(state.DCSlug.ValueString(), state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting bucket", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// BUCKET PERMISSION — utho_object_storage_permission
// ══════════════════════════════════════════════════════════════

type ObjectStoragePermissionResource struct{ client *client.Client }

type ObjectStoragePermissionModel struct {
	ID         types.String `tfsdk:"id"`
	DCSlug     types.String `tfsdk:"dcslug"`
	BucketName types.String `tfsdk:"bucket_name"`
	Permission types.String `tfsdk:"permission"`
	AccessKey  types.String `tfsdk:"access_key"`
}

func NewObjectStoragePermissionResource() resource.Resource {
	return &ObjectStoragePermissionResource{}
}

func (r *ObjectStoragePermissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object_storage_permission"
}

func (r *ObjectStoragePermissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Grant permissions on a Utho Object Storage bucket to an access key.",
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"dcslug":      schema.StringAttribute{Required: true, Description: "Data center slug.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"bucket_name": schema.StringAttribute{Required: true, Description: "Bucket name to grant permission on.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"permission":  schema.StringAttribute{Required: true, Description: "Permission level: read, write, full, or none."},
			"access_key":  schema.StringAttribute{Required: true, Description: "Access key to grant permission to.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *ObjectStoragePermissionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ObjectStoragePermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ObjectStoragePermissionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.GrantBucketPermission(plan.DCSlug.ValueString(), plan.BucketName.ValueString(), plan.Permission.ValueString(), plan.AccessKey.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error granting bucket permission", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s:%s:%s", plan.BucketName.ValueString(), plan.AccessKey.ValueString(), plan.Permission.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ObjectStoragePermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *ObjectStoragePermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ObjectStoragePermissionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.GrantBucketPermission(plan.DCSlug.ValueString(), plan.BucketName.ValueString(), plan.Permission.ValueString(), plan.AccessKey.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error updating bucket permission", fmt.Sprintf("%s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ObjectStoragePermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ObjectStoragePermissionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Revoke by setting to none
	if err := r.client.GrantBucketPermission(state.DCSlug.ValueString(), state.BucketName.ValueString(), "none", state.AccessKey.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error revoking bucket permission", fmt.Sprintf("%s", err))
	}
}

// ══════════════════════════════════════════════════════════════
// ACCESS KEY — utho_object_storage_key
// ══════════════════════════════════════════════════════════════

type ObjectStorageKeyResource struct{ client *client.Client }

type ObjectStorageKeyModel struct {
	ID        types.String `tfsdk:"id"`
	DCSlug    types.String `tfsdk:"dcslug"`
	Name      types.String `tfsdk:"name"`
	AccessKey types.String `tfsdk:"access_key"`
	SecretKey types.String `tfsdk:"secret_key"`
	Status    types.String `tfsdk:"status"`
}

func NewObjectStorageKeyResource() resource.Resource { return &ObjectStorageKeyResource{} }

func (r *ObjectStorageKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object_storage_key"
}

func (r *ObjectStorageKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage access keys for Utho Object Storage.",
		Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"dcslug":     schema.StringAttribute{Required: true, Description: "Data center slug.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"name":       schema.StringAttribute{Required: true, Description: "Access key name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"access_key": schema.StringAttribute{Computed: true, Sensitive: true, Description: "Generated access key."},
			"secret_key": schema.StringAttribute{Computed: true, Sensitive: true, Description: "Generated secret key."},
			"status":     schema.StringAttribute{Optional: true, Computed: true, Description: "Access key status: enable or disable."},
		},
	}
}

func (r *ObjectStorageKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ObjectStorageKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ObjectStorageKeyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	key, err := r.client.CreateAccessKey(plan.DCSlug.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating access key", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(key.AccessKey)
	plan.AccessKey = types.StringValue(key.AccessKey)
	plan.SecretKey = types.StringValue(key.SecretKey)

	// Keys are created as enabled by default
	// If user requested disabled, apply it now
	requestedStatus := plan.Status.ValueString()
	if requestedStatus == "disable" {
		if err := r.client.UpdateAccessKeyStatus(plan.DCSlug.ValueString(), key.AccessKey, "disable"); err != nil {
			resp.Diagnostics.AddError("Error disabling access key", fmt.Sprintf("%s", err))
			return
		}
		plan.Status = types.StringValue("disable")
	} else {
		plan.Status = types.StringValue("enable")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ObjectStorageKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *ObjectStorageKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ObjectStorageKeyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use access key from state — plan has it as unknown (sensitive computed)
	accessKey := state.AccessKey.ValueString()
	if accessKey == "" {
		accessKey = plan.AccessKey.ValueString()
	}

	if !plan.Status.IsNull() && plan.Status.ValueString() != state.Status.ValueString() {
		if err := r.client.UpdateAccessKeyStatus(plan.DCSlug.ValueString(), accessKey, plan.Status.ValueString()); err != nil {
			resp.Diagnostics.AddError("Error updating access key status", fmt.Sprintf("%s", err))
			return
		}
	}

	// Copy known values from state
	plan.AccessKey = state.AccessKey
	plan.SecretKey = state.SecretKey
	if plan.Status.IsNull() || plan.Status.IsUnknown() {
		plan.Status = state.Status
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ObjectStorageKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ObjectStorageKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteAccessKey(state.DCSlug.ValueString(), state.AccessKey.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting access key", fmt.Sprintf("%s", err))
	}
}
