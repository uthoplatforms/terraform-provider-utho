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
// data.utho_billing_usage
// ══════════════════════════════════════════════════════════════

type BillingUsageDataSource struct{ client *client.Client }

type BillingCategoryModel struct {
	Key           types.String  `tfsdk:"key"`
	Label         types.String  `tfsdk:"label"`
	Subtotal      types.Float64 `tfsdk:"subtotal"`
	ResourceCount types.Int64   `tfsdk:"resource_count"`
	SharePct      types.Float64 `tfsdk:"share_pct"`
}

type BillingTopResourceModel struct {
	Name          types.String  `tfsdk:"name"`
	Category      types.String  `tfsdk:"category"`
	CategoryLabel types.String  `tfsdk:"category_label"`
	Amount        types.Float64 `tfsdk:"amount"`
	SharePct      types.Float64 `tfsdk:"share_pct"`
}

type BillingUsageModel struct {
	ID                types.String              `tfsdk:"id"`
	CycleStart        types.String              `tfsdk:"cycle_start"`
	CycleEnd          types.String              `tfsdk:"cycle_end"`
	DaysIn            types.Int64               `tfsdk:"days_in"`
	DaysLeft          types.Int64               `tfsdk:"days_left"`
	UsageExclTax      types.Float64             `tfsdk:"usage_excl_tax"`
	Tax               types.Float64             `tfsdk:"tax"`
	TotalBeforeWallet types.Float64             `tfsdk:"total_before_wallet"`
	DueAmount         types.Float64             `tfsdk:"due_amount"`
	WalletBalance     types.Float64             `tfsdk:"wallet_balance"`
	AvailableCredit   types.Float64             `tfsdk:"available_credit"`
	ByCategory        []BillingCategoryModel    `tfsdk:"by_category"`
	TopResources      []BillingTopResourceModel `tfsdk:"top_resources"`
}

func NewBillingUsageDataSource() datasource.DataSource { return &BillingUsageDataSource{} }

func (d *BillingUsageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_billing_usage"
}

func (d *BillingUsageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve current billing cycle usage breakdown by resource category.",
		Attributes: map[string]schema.Attribute{
			"id":                  schema.StringAttribute{Computed: true},
			"cycle_start":         schema.StringAttribute{Computed: true, Description: "Billing cycle start date."},
			"cycle_end":           schema.StringAttribute{Computed: true, Description: "Billing cycle end date."},
			"days_in":             schema.Int64Attribute{Computed: true, Description: "Days elapsed in current cycle."},
			"days_left":           schema.Int64Attribute{Computed: true, Description: "Days remaining in current cycle."},
			"usage_excl_tax":      schema.Float64Attribute{Computed: true, Description: "Total usage excluding tax (INR)."},
			"tax":                 schema.Float64Attribute{Computed: true, Description: "Tax amount (INR)."},
			"total_before_wallet": schema.Float64Attribute{Computed: true, Description: "Total before wallet deduction (INR)."},
			"due_amount":          schema.Float64Attribute{Computed: true, Description: "Amount due (INR)."},
			"wallet_balance":      schema.Float64Attribute{Computed: true, Description: "Wallet balance (INR)."},
			"available_credit":    schema.Float64Attribute{Computed: true, Description: "Available free credits (INR)."},
			"by_category": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Usage breakdown by resource category.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key":            schema.StringAttribute{Computed: true},
						"label":          schema.StringAttribute{Computed: true},
						"subtotal":       schema.Float64Attribute{Computed: true},
						"resource_count": schema.Int64Attribute{Computed: true},
						"share_pct":      schema.Float64Attribute{Computed: true},
					},
				},
			},
			"top_resources": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Top resources by spend this cycle.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":           schema.StringAttribute{Computed: true},
						"category":       schema.StringAttribute{Computed: true},
						"category_label": schema.StringAttribute{Computed: true},
						"amount":         schema.Float64Attribute{Computed: true},
						"share_pct":      schema.Float64Attribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *BillingUsageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BillingUsageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	usage, err := d.client.GetBillingUsage()
	if err != nil {
		resp.Diagnostics.AddError("Error reading billing usage", fmt.Sprintf("%s", err))
		return
	}

	state := BillingUsageModel{
		ID:                types.StringValue("billing_usage"),
		CycleStart:        types.StringValue(usage.Cycle.Start),
		CycleEnd:          types.StringValue(usage.Cycle.End),
		DaysIn:            types.Int64Value(int64(usage.Cycle.DaysIn)),
		DaysLeft:          types.Int64Value(int64(usage.Cycle.DaysLeft)),
		UsageExclTax:      types.Float64Value(usage.Summary.UsageExclTax),
		Tax:               types.Float64Value(usage.Summary.Tax),
		TotalBeforeWallet: types.Float64Value(usage.Summary.TotalBeforeWallet),
		DueAmount:         types.Float64Value(usage.Summary.DueAmount),
		WalletBalance:     types.Float64Value(usage.Summary.WalletBalance),
		AvailableCredit:   types.Float64Value(usage.Summary.AvailableCredit),
	}

	for _, cat := range usage.ByCategory {
		state.ByCategory = append(state.ByCategory, BillingCategoryModel{
			Key:           types.StringValue(cat.Key),
			Label:         types.StringValue(cat.Label),
			Subtotal:      types.Float64Value(cat.Subtotal),
			ResourceCount: types.Int64Value(int64(cat.ResourceCount)),
			SharePct:      types.Float64Value(cat.SharePct),
		})
	}

	for _, tr := range usage.TopResources {
		state.TopResources = append(state.TopResources, BillingTopResourceModel{
			Name:          types.StringValue(tr.Name),
			Category:      types.StringValue(tr.Category),
			CategoryLabel: types.StringValue(tr.CategoryLabel),
			Amount:        types.Float64Value(tr.Amount),
			SharePct:      types.Float64Value(tr.SharePct),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// ══════════════════════════════════════════════════════════════
// data.utho_billing_invoices
// ══════════════════════════════════════════════════════════════

type BillingInvoicesDataSource struct{ client *client.Client }

type InvoiceModel struct {
	ID         types.Int64   `tfsdk:"id"`
	InvoiceNum types.String  `tfsdk:"invoice_num"`
	Date       types.String  `tfsdk:"date"`
	DueDate    types.String  `tfsdk:"due_date"`
	Amount     types.Float64 `tfsdk:"amount"`
	Status     types.String  `tfsdk:"status"`
}

type BillingInvoicesModel struct {
	ID       types.String   `tfsdk:"id"`
	Invoices []InvoiceModel `tfsdk:"invoices"`
}

func NewBillingInvoicesDataSource() datasource.DataSource { return &BillingInvoicesDataSource{} }

func (d *BillingInvoicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_billing_invoices"
}

func (d *BillingInvoicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve all billing invoices for your Utho account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"invoices": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of invoices.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.Int64Attribute{Computed: true},
						"invoice_num": schema.StringAttribute{Computed: true},
						"date":        schema.StringAttribute{Computed: true},
						"due_date":    schema.StringAttribute{Computed: true},
						"amount":      schema.Float64Attribute{Computed: true},
						"status":      schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *BillingInvoicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BillingInvoicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	invoices, err := d.client.GetBillingInvoices()
	if err != nil {
		resp.Diagnostics.AddError("Error reading billing invoices", fmt.Sprintf("%s", err))
		return
	}

	state := BillingInvoicesModel{ID: types.StringValue("billing_invoices")}
	for _, inv := range invoices {
		state.Invoices = append(state.Invoices, InvoiceModel{
			ID:         types.Int64Value(int64(inv.ID)),
			InvoiceNum: types.StringValue(inv.InvoiceNum),
			Date:       types.StringValue(inv.Date),
			DueDate:    types.StringValue(inv.DueDate),
			Amount:     types.Float64Value(inv.Amount),
			Status:     types.StringValue(inv.Status),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// ══════════════════════════════════════════════════════════════
// data.utho_billing_cost_by_project
// ══════════════════════════════════════════════════════════════

type BillingCostByProjectDataSource struct{ client *client.Client }

type ProjectCostModel struct {
	ProjectID     types.Int64   `tfsdk:"project_id"`
	Name          types.String  `tfsdk:"name"`
	IsDefault     types.Bool    `tfsdk:"is_default"`
	Subtotal      types.Float64 `tfsdk:"subtotal"`
	ResourceCount types.Int64   `tfsdk:"resource_count"`
}

type BillingCostByProjectModel struct {
	ID                types.String       `tfsdk:"id"`
	Projects          []ProjectCostModel `tfsdk:"projects"`
	TotalAttributed   types.Float64      `tfsdk:"total_attributed"`
	TotalUnattributed types.Float64      `tfsdk:"total_unattributed"`
	TotalMonth        types.Float64      `tfsdk:"total_month"`
	CycleStart        types.String       `tfsdk:"cycle_start"`
	CycleEnd          types.String       `tfsdk:"cycle_end"`
}

func NewBillingCostByProjectDataSource() datasource.DataSource {
	return &BillingCostByProjectDataSource{}
}

func (d *BillingCostByProjectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_billing_cost_by_project"
}

func (d *BillingCostByProjectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve current month cost breakdown by project.",
		Attributes: map[string]schema.Attribute{
			"id":                 schema.StringAttribute{Computed: true},
			"total_attributed":   schema.Float64Attribute{Computed: true, Description: "Total cost attributed to projects (INR)."},
			"total_unattributed": schema.Float64Attribute{Computed: true, Description: "Total cost not attributed to any project (INR)."},
			"total_month":        schema.Float64Attribute{Computed: true, Description: "Total month cost (INR)."},
			"cycle_start":        schema.StringAttribute{Computed: true, Description: "Billing cycle start date."},
			"cycle_end":          schema.StringAttribute{Computed: true, Description: "Billing cycle end date."},
			"projects": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Cost breakdown per project.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"project_id":     schema.Int64Attribute{Computed: true},
						"name":           schema.StringAttribute{Computed: true},
						"is_default":     schema.BoolAttribute{Computed: true},
						"subtotal":       schema.Float64Attribute{Computed: true},
						"resource_count": schema.Int64Attribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *BillingCostByProjectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BillingCostByProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data, err := d.client.GetBillingCostByProject()
	if err != nil {
		resp.Diagnostics.AddError("Error reading billing cost by project", fmt.Sprintf("%s", err))
		return
	}

	state := BillingCostByProjectModel{
		ID:                types.StringValue("billing_cost_by_project"),
		TotalAttributed:   types.Float64Value(data.TotalAttributed),
		TotalUnattributed: types.Float64Value(data.TotalUnattributed),
		TotalMonth:        types.Float64Value(data.TotalMonth),
		CycleStart:        types.StringValue(data.Cycle.Start),
		CycleEnd:          types.StringValue(data.Cycle.End),
	}

	for _, p := range data.Projects {
		state.Projects = append(state.Projects, ProjectCostModel{
			ProjectID:     types.Int64Value(int64(p.ProjectID)),
			Name:          types.StringValue(p.Name),
			IsDefault:     types.BoolValue(p.IsDefault),
			Subtotal:      types.Float64Value(p.Subtotal),
			ResourceCount: types.Int64Value(int64(p.ResourceCount)),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
