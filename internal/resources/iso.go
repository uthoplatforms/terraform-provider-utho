package resources

import (
	"context"
	"fmt"

	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ISOResource struct{ client *client.Client }

type ISOModel struct {
	ID       types.String  `tfsdk:"id"`
	Name     types.String  `tfsdk:"name"`
	DCSlug   types.String  `tfsdk:"dcslug"`
	URL      types.String  `tfsdk:"url"`
	File     types.String  `tfsdk:"file"`
	Size     types.Float64 `tfsdk:"size"`
	AddedAt  types.String  `tfsdk:"added_at"`
	Location types.String  `tfsdk:"location"`
}

func NewISOResource() resource.Resource { return &ISOResource{} }

func (r *ISOResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iso"
}

func (r *ISOResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Upload and manage an ISO image on Utho by providing a public download URL. The ISO is downloaded asynchronously. Use the file attribute as input to utho_cloud_iso to mount it on a cloud instance.",
		Attributes: map[string]schema.Attribute{
			"id":       schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"file":     schema.StringAttribute{Computed: true, Description: "ISO filename assigned by Utho (e.g. PzFZp-194454.iso). Use this as input to utho_cloud_iso.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"added_at": schema.StringAttribute{Computed: true, Description: "Timestamp when the ISO was added.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"size":     schema.Float64Attribute{Computed: true, Description: "ISO size in GB."},
			"location": schema.StringAttribute{Computed: true, Description: "Data center location name."},
			"name":     schema.StringAttribute{Required: true, Description: "ISO name. Changing forces new resource.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"dcslug":   schema.StringAttribute{Required: true, Description: "Data center slug where the ISO will be stored (e.g. innoida, inmumbaizone2). Changing forces new resource.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"url":      schema.StringAttribute{Required: true, Description: "Public download URL of the ISO file. Changing forces new resource.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *ISOResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ISOResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ISOModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.AddISO(&client.ISOAddRequest{
		DCSlug: plan.DCSlug.ValueString(),
		Name:   plan.Name.ValueString(),
		URL:    plan.URL.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error adding ISO", fmt.Sprintf("%s", err))
		return
	}

	plan.ID = types.StringValue(id)
	plan.File = types.StringValue("")
	plan.AddedAt = types.StringValue("")
	plan.Size = types.Float64Value(0)
	plan.Location = types.StringValue("")

	// Try to find the ISO in the list by name
	isos, err2 := r.client.ListISOs()
	if err2 == nil {
		for _, iso := range isos {
			if iso.Name == plan.Name.ValueString() && iso.DC == plan.DCSlug.ValueString() {
				plan.File = types.StringValue(iso.File)
				plan.AddedAt = types.StringValue(iso.AddedAt)
				plan.Size = types.Float64Value(iso.Size)
				plan.Location = types.StringValue(iso.DCLocation.Location)
				break
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ISOResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ISOModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If file is known use it to look up, else find by name
	if state.File.ValueString() != "" {
		iso, err := r.client.GetISOByFile(state.File.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error reading ISO", fmt.Sprintf("%s", err))
			return
		}
		if iso == nil {
			resp.State.RemoveResource(ctx)
			return
		}
		state.Size = types.Float64Value(iso.Size)
		state.Location = types.StringValue(iso.DCLocation.Location)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ISOResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All fields are RequiresReplace — no in-place update needed
}

func (r *ISOResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ISOModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	file := state.File.ValueString()
	if file == "" {
		// If file not captured, find by name
		isos, err := r.client.ListISOs()
		if err == nil {
			for _, iso := range isos {
				if iso.Name == state.Name.ValueString() {
					file = iso.File
					break
				}
			}
		}
	}

	if file == "" {
		// Already gone
		return
	}

	if err := r.client.DeleteISO(file); err != nil {
		resp.Diagnostics.AddError("Error deleting ISO", fmt.Sprintf("%s", err))
	}
}
