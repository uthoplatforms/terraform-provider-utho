package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/uthoterraform/terraform-provider-utho/api"
)

var (
	_ datasource.DataSource              = &AccountDataSource{}
	_ datasource.DataSourceWithConfigure = &AccountDataSource{}
)

type AccountDataSource struct {
	client *api.Client
}

type (
	// AccountDataSourceModel struct {
	// 	LocationId types.String   `tfsdk:"location_id"`
	// 	Account    []AccountModel `tfsdk:"accounts"`
	// }
	// AccountModel struct {
	// 	Slug  types.String `tfsdk:"slug"`
	// 	Name  types.String `tfsdk:"name"`
	// 	RAM   types.Int64  `tfsdk:"ram"`
	// 	Disk  types.Int64  `tfsdk:"disk"`
	// 	CPU   CPU          `tfsdk:"cpu"`
	// 	Price Price        `tfsdk:"price"`
	// }
	// CPU struct {
	// 	Cores   types.Int64 `tfsdk:"cores"`
	// 	Threads types.Int64 `tfsdk:"threads"`
	// }
	// Price struct {
	// 	Amount   types.Int64  `tfsdk:"amount"`
	// 	Currency types.String `tfsdk:"currency"`
	// }

	AccountDataSourceModel struct {
		User User `tfsdk:"user"`
	}

	User struct {
		ID                      types.String  `tfsdk:"id"`
		Type                    types.String  `tfsdk:"type"`
		Fullname                types.String  `tfsdk:"fullname"`
		Company                 types.String  `tfsdk:"company"`
		Email                   types.String  `tfsdk:"email"`
		Address                 types.String  `tfsdk:"address"`
		City                    types.String  `tfsdk:"city"`
		State                   types.String  `tfsdk:"state"`
		Country                 types.String  `tfsdk:"country"`
		Postcode                types.String  `tfsdk:"postcode"`
		Mobile                  types.String  `tfsdk:"mobile"`
		Mobilecc                types.String  `tfsdk:"mobilecc"`
		Gstnumber               types.String  `tfsdk:"gstnumber"`
		SupportneedTitle        types.String  `tfsdk:"supportneed_title"`
		SupportneedUsecase      types.String  `tfsdk:"supportneed_usecase"`
		SupportneedBusinesstype types.String  `tfsdk:"supportneed_businesstype"`
		SupportneedMonthlyspend types.String  `tfsdk:"supportneed_monthlyspend"`
		SupportneedEmployeesize types.String  `tfsdk:"supportneed_employeesize"`
		SupportFieldsRequired   types.String  `tfsdk:"support_fields_required"`
		TwofaSettings           types.String  `tfsdk:"twofa_settings"`
		Currencyprefix          types.String  `tfsdk:"currencyprefix"`
		Currencyrate            types.String  `tfsdk:"currencyrate"`
		Currency                types.String  `tfsdk:"currency"`
		Credit                  types.Float64 `tfsdk:"credit"`
		Availablecredit         types.Float64 `tfsdk:"availablecredit"`
		Freecredit              types.Float64 `tfsdk:"freecredit"`
		Currentusages           types.Float64 `tfsdk:"currentusages"`
		Kyc                     types.String  `tfsdk:"kyc"`
		SmsVerified             types.String  `tfsdk:"sms_verified"`
		Verify                  types.String  `tfsdk:"verify"`
		IsPartner               types.String  `tfsdk:"is_partner"`
		Partnerid               types.String  `tfsdk:"partnerid"`
		Twofa                   types.String  `tfsdk:"twofa"`
		EmailVerified           types.String  `tfsdk:"email_verified"`
		Cloudlimit              types.String  `tfsdk:"cloudlimit"`
		K8SLimit                types.String  `tfsdk:"k8s_limit"`
		IsReseller              types.String  `tfsdk:"is_reseller"`
		Singleinvoice           types.String  `tfsdk:"singleinvoice"`
		RazorpayCustomerid      types.String  `tfsdk:"razorpay_customerid"`
		RazorpayOrderid         types.String  `tfsdk:"razorpay_orderid"`
		StripeCustomer          types.String  `tfsdk:"stripe_customer"`
		TotalCloudservers       types.String  `tfsdk:"total_cloudservers"`
		Resources               []Resources   `tfsdk:"resources"`
		Rvn                     types.String  `tfsdk:"rvn"`
		CAdded                  types.String  `tfsdk:"c_added"`
		RazorpaySub             types.String  `tfsdk:"razorpay_sub"`
		AffiliateLoginid        types.String  `tfsdk:"affiliate_loginid"`
	}
	Resources struct {
		Product types.String `tfsdk:"product"`
		Count   types.String `tfsdk:"count"`
	}
)

func NewAccountDataSource() datasource.DataSource {
	return &AccountDataSource{}
}

func (*AccountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

// Schema defines the schema for the data source.
func (d *AccountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"user": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "User model",
				Attributes: map[string]schema.Attribute{
					"id":                       schema.StringAttribute{Computed: true, Description: "Id"},
					"type":                     schema.StringAttribute{Computed: true, Description: "Type"},
					"fullname":                 schema.StringAttribute{Computed: true, Description: "Fullname"},
					"company":                  schema.StringAttribute{Computed: true, Description: "Company"},
					"email":                    schema.StringAttribute{Computed: true, Description: "Email"},
					"address":                  schema.StringAttribute{Computed: true, Description: "Address"},
					"city":                     schema.StringAttribute{Computed: true, Description: "City"},
					"state":                    schema.StringAttribute{Computed: true, Description: "State"},
					"country":                  schema.StringAttribute{Computed: true, Description: "Country"},
					"postcode":                 schema.StringAttribute{Computed: true, Description: "Postcode"},
					"mobile":                   schema.StringAttribute{Computed: true, Description: "Mobile"},
					"mobilecc":                 schema.StringAttribute{Computed: true, Description: "Mobilecc"},
					"gstnumber":                schema.StringAttribute{Computed: true, Description: "Gstnumber"},
					"supportneed_title":        schema.StringAttribute{Computed: true, Description: "Supportneed Title"},
					"supportneed_usecase":      schema.StringAttribute{Computed: true, Description: "Supportneed Usecase"},
					"supportneed_businesstype": schema.StringAttribute{Computed: true, Description: "Supportneed Businesstype"},
					"supportneed_monthlyspend": schema.StringAttribute{Computed: true, Description: "Supportneed Monthlyspend"},
					"supportneed_employeesize": schema.StringAttribute{Computed: true, Description: "Supportneed Employeesize"},
					"support_fields_required":  schema.StringAttribute{Computed: true, Description: "Support Fields Required"},
					"twofa_settings":           schema.StringAttribute{Computed: true, Description: "Twofa Settings"},
					"currencyprefix":           schema.StringAttribute{Computed: true, Description: "Currencyprefix"},
					"currencyrate":             schema.StringAttribute{Computed: true, Description: "Currencyrate"},
					"currency":                 schema.StringAttribute{Computed: true, Description: "Currency"},
					"credit":                   schema.Float64Attribute{Computed: true, Description: "Credit"},
					"availablecredit":          schema.Float64Attribute{Computed: true, Description: "Availablecredit"},
					"freecredit":               schema.Float64Attribute{Computed: true, Description: "Freecredit"},
					"currentusages":            schema.Float64Attribute{Computed: true, Description: "Currentusages"},
					"kyc":                      schema.StringAttribute{Computed: true, Description: "Kyc"},
					"sms_verified":             schema.StringAttribute{Computed: true, Description: "Sms Verified"},
					"verify":                   schema.StringAttribute{Computed: true, Description: "Verify"},
					"is_partner":               schema.StringAttribute{Computed: true, Description: "Is Partner"},
					"partnerid":                schema.StringAttribute{Computed: true, Description: "Partnerid"},
					"twofa":                    schema.StringAttribute{Computed: true, Description: "Twofa"},
					"email_verified":           schema.StringAttribute{Computed: true, Description: "Email Verified"},
					"cloudlimit":               schema.StringAttribute{Computed: true, Description: "Cloudlimit"},
					"k8s_limit":                schema.StringAttribute{Computed: true, Description: "k8s Limit"},
					"is_reseller":              schema.StringAttribute{Computed: true, Description: "Is Reseller"},
					"singleinvoice":            schema.StringAttribute{Computed: true, Description: "Singleinvoice"},
					"razorpay_customerid":      schema.StringAttribute{Computed: true, Description: "Razorpay Customerid"},
					"razorpay_orderid":         schema.StringAttribute{Computed: true, Description: "Razorpay Orderid"},
					"stripe_customer":          schema.StringAttribute{Computed: true, Description: "Stripe Customer"},
					"total_cloudservers":       schema.StringAttribute{Computed: true, Description: "Total Cloudservers"},
					"resources": schema.ListNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"product": schema.StringAttribute{Computed: true, Description: "product"},
								"count":   schema.StringAttribute{Computed: true, Description: "count"},
							},
						},
					},
					"rvn":               schema.StringAttribute{Computed: true, Description: "Rvn"},
					"c_added":           schema.StringAttribute{Computed: true, Description: "C Added"},
					"razorpay_sub":      schema.StringAttribute{Computed: true, Description: "Razorpay Sub"},
					"affiliate_loginid": schema.StringAttribute{Computed: true, Description: "Affiliate Loginid"},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *AccountDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Account Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

// Read refreshes the Terraform state with the latest data
func (d *AccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read `item` data source")
	// get account
	account, err := d.client.GetAccount(ctx)
	userInfo := account.User
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list `account`",
			err.Error(),
		)
		return
	}
	// Map response body to model
	state := AccountDataSourceModel{
		User: User{
			ID:                      types.StringValue(userInfo.ID),
			Type:                    types.StringValue(userInfo.Type),
			Fullname:                types.StringValue(userInfo.Fullname),
			Company:                 types.StringValue(userInfo.Company),
			Email:                   types.StringValue(userInfo.Email),
			Address:                 types.StringValue(userInfo.Address),
			City:                    types.StringValue(userInfo.City),
			State:                   types.StringValue(userInfo.State),
			Country:                 types.StringValue(userInfo.Country),
			Postcode:                types.StringValue(userInfo.Postcode),
			Mobile:                  types.StringValue(userInfo.Mobile),
			Mobilecc:                types.StringValue(userInfo.Mobilecc),
			Gstnumber:               types.StringValue(userInfo.Gstnumber),
			SupportneedTitle:        types.StringValue(userInfo.SupportneedTitle),
			SupportneedUsecase:      types.StringValue(userInfo.SupportneedUsecase),
			SupportneedBusinesstype: types.StringValue(userInfo.SupportneedBusinesstype),
			SupportneedMonthlyspend: types.StringValue(userInfo.SupportneedMonthlyspend),
			SupportneedEmployeesize: types.StringValue(userInfo.SupportneedEmployeesize),
			SupportFieldsRequired:   types.StringValue(userInfo.SupportFieldsRequired),
			TwofaSettings:           types.StringValue(userInfo.TwofaSettings),
			Currencyprefix:          types.StringValue(userInfo.Currencyprefix),
			Currencyrate:            types.StringValue(userInfo.Currencyrate),
			Currency:                types.StringValue(userInfo.Currency),
			Credit:                  types.Float64Value(userInfo.Credit),
			Availablecredit:         types.Float64Value(userInfo.Availablecredit),
			Freecredit:              types.Float64Value(userInfo.Freecredit),
			Currentusages:           types.Float64Value(userInfo.Currentusages),
			Kyc:                     types.StringValue(userInfo.Kyc),
			SmsVerified:             types.StringValue(userInfo.SmsVerified),
			Verify:                  types.StringValue(userInfo.Verify),
			IsPartner:               types.StringValue(userInfo.IsPartner),
			Partnerid:               types.StringValue(userInfo.Partnerid),
			Twofa:                   types.StringValue(userInfo.Twofa),
			EmailVerified:           types.StringValue(userInfo.EmailVerified),
			Cloudlimit:              types.StringValue(userInfo.Cloudlimit),
			K8SLimit:                types.StringValue(userInfo.K8SLimit),
			IsReseller:              types.StringValue(userInfo.IsReseller),
			Singleinvoice:           types.StringValue(userInfo.Singleinvoice),
			RazorpayCustomerid:      types.StringValue(userInfo.RazorpayCustomerid),
			RazorpayOrderid:         types.StringValue(userInfo.RazorpayOrderid),
			StripeCustomer:          types.StringValue(userInfo.StripeCustomer),
			TotalCloudservers:       types.StringValue(userInfo.TotalCloudservers),
			Rvn:                     types.StringValue(userInfo.Rvn),
			CAdded:                  types.StringValue(userInfo.CAdded),
			RazorpaySub:             types.StringValue(userInfo.RazorpaySub),
			AffiliateLoginid:        types.StringValue(userInfo.AffiliateLoginid),
		},
	}
	for _, resource := range userInfo.Resources {
		resourceState := Resources{
			Product: types.StringValue(resource.Product),
			Count:   types.StringValue(resource.Count),
		}
		state.User.Resources = append(state.User.Resources, resourceState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Finished reading `account` data source", map[string]any{"success": true})
}
