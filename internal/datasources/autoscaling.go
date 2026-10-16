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

// ── data.utho_autoscalings ────────────────────────────────────────────────

type AutoScalingsDataSource struct{ client *client.Client }

type AutoScalingsModel struct {
	AutoScalings []AutoScalingItem `tfsdk:"autoscalings"`
}

type AutoScalingItem struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	DCSlug      types.String `tfsdk:"dcslug"`
	MinSize     types.String `tfsdk:"minsize"`
	MaxSize     types.String `tfsdk:"maxsize"`
	DesiredSize types.String `tfsdk:"desiredsize"`
	Status      types.String `tfsdk:"status"`
}

func NewAutoScalingsDataSource() datasource.DataSource {
	return &AutoScalingsDataSource{}
}

func (d *AutoScalingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_autoscalings"
}

func (d *AutoScalingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all auto scaling groups in your Utho account.",
		Attributes: map[string]schema.Attribute{
			"autoscalings": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.StringAttribute{Computed: true, Description: "Auto scaling group ID."},
						"name":        schema.StringAttribute{Computed: true, Description: "Auto scaling group name."},
						"dcslug":      schema.StringAttribute{Computed: true, Description: "Data center slug."},
						"minsize":     schema.StringAttribute{Computed: true, Description: "Minimum instances."},
						"maxsize":     schema.StringAttribute{Computed: true, Description: "Maximum instances."},
						"desiredsize": schema.StringAttribute{Computed: true, Description: "Desired instances."},
						"status":      schema.StringAttribute{Computed: true, Description: "Status."},
					},
				},
			},
		},
	}
}

func (d *AutoScalingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AutoScalingsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	respBytes, err := d.client.Get("/autoscaling")
	if err != nil {
		resp.Diagnostics.AddError("Error fetching auto scaling groups", fmt.Sprintf("%s", err))
		return
	}

	var result struct {
		Autoscaling []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			DCSlug      string `json:"dcslug"`
			MinSize     string `json:"minsize"`
			MaxSize     string `json:"maxsize"`
			DesiredSize string `json:"desiredsize"`
			Status      string `json:"status"`
		} `json:"autoscaling"`
	}

	if err := json.Unmarshal(respBytes, &result); err != nil {
		resp.Diagnostics.AddError("Error parsing auto scaling groups", fmt.Sprintf("%s", err))
		return
	}

	var items []AutoScalingItem
	for _, a := range result.Autoscaling {
		items = append(items, AutoScalingItem{
			ID:          types.StringValue(a.ID),
			Name:        types.StringValue(a.Name),
			DCSlug:      types.StringValue(a.DCSlug),
			MinSize:     types.StringValue(a.MinSize),
			MaxSize:     types.StringValue(a.MaxSize),
			DesiredSize: types.StringValue(a.DesiredSize),
			Status:      types.StringValue(a.Status),
		})
	}
	if items == nil {
		items = []AutoScalingItem{}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &AutoScalingsModel{AutoScalings: items})...)
}

// ── data.utho_target_groups ───────────────────────────────────────────────

type TargetGroupsDataSource struct{ client *client.Client }

type TargetGroupsModel struct {
	TargetGroups []TargetGroupItem `tfsdk:"target_groups"`
}

type TargetGroupItem struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Protocol types.String `tfsdk:"protocol"`
	Port     types.String `tfsdk:"port"`
}

func NewTargetGroupsDataSource() datasource.DataSource {
	return &TargetGroupsDataSource{}
}

func (d *TargetGroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_target_groups"
}

func (d *TargetGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all target groups in your Utho account.",
		Attributes: map[string]schema.Attribute{
			"target_groups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":       schema.StringAttribute{Computed: true, Description: "Target group ID."},
						"name":     schema.StringAttribute{Computed: true, Description: "Target group name."},
						"protocol": schema.StringAttribute{Computed: true, Description: "Protocol."},
						"port":     schema.StringAttribute{Computed: true, Description: "Port."},
					},
				},
			},
		},
	}
}

func (d *TargetGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TargetGroupsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	respBytes, err := d.client.Get("/targetgroup")
	if err != nil {
		resp.Diagnostics.AddError("Error fetching target groups", fmt.Sprintf("%s", err))
		return
	}

	var result struct {
		TargetGroups []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Protocol string `json:"protocol"`
			Port     string `json:"port"`
		} `json:"targetgroups"`
	}

	if err := json.Unmarshal(respBytes, &result); err != nil {
		resp.Diagnostics.AddError("Error parsing target groups", fmt.Sprintf("%s", err))
		return
	}

	var items []TargetGroupItem
	for _, t := range result.TargetGroups {
		items = append(items, TargetGroupItem{
			ID:       types.StringValue(t.ID),
			Name:     types.StringValue(t.Name),
			Protocol: types.StringValue(t.Protocol),
			Port:     types.StringValue(t.Port),
		})
	}
	if items == nil {
		items = []TargetGroupItem{}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &TargetGroupsModel{TargetGroups: items})...)
}
