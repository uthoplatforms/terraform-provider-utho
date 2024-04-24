package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/uthoterraform/terraform-provider-utho/helper"
)

type (
	Account struct {
		User User `json:"user,omitempty"`
	}

	User struct {
		ID                      string      `json:"id,omitempty"`
		Type                    string      `json:"type,omitempty"`
		Fullname                string      `json:"fullname,omitempty"`
		Company                 string      `json:"company,omitempty"`
		Email                   string      `json:"email,omitempty"`
		Address                 string      `json:"address,omitempty"`
		City                    string      `json:"city,omitempty"`
		State                   string      `json:"state,omitempty"`
		Country                 string      `json:"country,omitempty"`
		Postcode                string      `json:"postcode,omitempty"`
		Mobile                  string      `json:"mobile,omitempty"`
		Mobilecc                string      `json:"mobilecc,omitempty"`
		Gstnumber               string      `json:"gstnumber,omitempty"`
		SupportneedTitle        string      `json:"supportneed_title,omitempty"`
		SupportneedUsecase      string      `json:"supportneed_usecase,omitempty"`
		SupportneedBusinesstype string      `json:"supportneed_businesstype,omitempty"`
		SupportneedMonthlyspend string      `json:"supportneed_monthlyspend,omitempty"`
		SupportneedEmployeesize string      `json:"supportneed_employeesize,omitempty"`
		SupportFieldsRequired   string      `json:"support_fields_required,omitempty"`
		TwofaSettings           string      `json:"twofa_settings,omitempty"`
		Currencyprefix          string      `json:"currencyprefix,omitempty"`
		Currencyrate            string      `json:"currencyrate,omitempty"`
		Currency                string      `json:"currency,omitempty"`
		Credit                  float64     `json:"credit,omitempty"`
		Availablecredit         float64     `json:"availablecredit,omitempty"`
		Freecredit              float64     `json:"freecredit,omitempty"`
		Currentusages           float64     `json:"currentusages,omitempty"`
		Kyc                     string      `json:"kyc,omitempty"`
		SmsVerified             string      `json:"sms_verified,omitempty"`
		Verify                  string      `json:"verify,omitempty"`
		IsPartner               string      `json:"is_partner,omitempty"`
		Partnerid               string      `json:"partnerid,omitempty"`
		Twofa                   string      `json:"twofa,omitempty"`
		EmailVerified           string      `json:"email_verified,omitempty"`
		Cloudlimit              string      `json:"cloudlimit,omitempty"`
		K8SLimit                string      `json:"k8s_limit,omitempty"`
		IsReseller              string      `json:"is_reseller,omitempty"`
		Singleinvoice           string      `json:"singleinvoice,omitempty"`
		RazorpayCustomerid      string      `json:"razorpay_customerid,omitempty"`
		RazorpayOrderid         string      `json:"razorpay_orderid,omitempty"`
		StripeCustomer          string      `json:"stripe_customer,omitempty"`
		TotalCloudservers       string      `json:"total_cloudservers,omitempty"`
		Resources               []Resources `json:"resources,omitempty"`
		Rvn                     string      `json:"rvn,omitempty"`
		CAdded                  string      `json:"c_added,omitempty"`
		RazorpaySub             string      `json:"razorpay_sub,omitempty"`
		AffiliateLoginid        string      `json:"affiliate_loginid,omitempty"`
	}
	Resources struct {
		Product string `json:"product,omitempty"`
		Count   string `json:"count,omitempty"`
	}
)

func (c *Client) GetAccount(ctx context.Context) (Account, error) {
	uri := BASE_URL + "account/info"

	resp, err := helper.NewUthoRequest(ctx, http.MethodGet, uri, nil, c.token)
	if err != nil {
		return Account{}, err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return Account{}, errors.New("unexpected http error code received for geting account data status code :" + strconv.Itoa(resp.StatusCode) + " body" + string(body))
	}
	defer resp.Body.Close()

	var account Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return Account{}, err
	}

	return account, nil
}
