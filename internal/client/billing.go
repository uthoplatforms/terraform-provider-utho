package client

import (
	"encoding/json"
	"fmt"
)

// ── Billing Usage structs ─────────────────────────────────────────────────

type BillingCycle struct {
	Start    string `json:"start"`
	End      string `json:"end"`
	DaysIn   int    `json:"days_in"`
	DaysLeft int    `json:"days_left"`
}

type BillingSummary struct {
	UsageExclTax      float64 `json:"usage_excl_tax"`
	FreeCreditApplied float64 `json:"free_credit_applied"`
	UsageAfterFree    float64 `json:"usage_after_free"`
	Tax               float64 `json:"tax"`
	TaxRate           float64 `json:"tax_rate"`
	TotalBeforeWallet float64 `json:"total_before_wallet"`
	PaidCreditApplied float64 `json:"paid_credit_applied"`
	DueAmount         float64 `json:"due_amount"`
	WalletBalance     float64 `json:"wallet_balance"`
	AvailableCredit   float64 `json:"availablecredit"`
}

type BillingCategory struct {
	Key           string  `json:"key"`
	Label         string  `json:"label"`
	Subtotal      float64 `json:"subtotal"`
	ResourceCount int     `json:"resource_count"`
	SharePct      float64 `json:"share_pct"`
}

type BillingTopResource struct {
	Name          string  `json:"name"`
	Category      string  `json:"category"`
	CategoryLabel string  `json:"category_label"`
	Amount        float64 `json:"amount"`
	SharePct      float64 `json:"share_pct"`
}

type BillingUsage struct {
	Cycle        BillingCycle         `json:"cycle"`
	Summary      BillingSummary       `json:"summary"`
	ByCategory   []BillingCategory    `json:"by_category"`
	TopResources []BillingTopResource `json:"top_resources"`
}

// ── Invoice structs ───────────────────────────────────────────────────────

type Invoice struct {
	ID         int     `json:"id"`
	InvoiceNum string  `json:"invoice_num"`
	Date       string  `json:"date"`
	DueDate    string  `json:"due_date"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
}

// ── Project cost structs ──────────────────────────────────────────────────

type ProjectCost struct {
	ProjectID     int     `json:"project_id"`
	Name          string  `json:"name"`
	IsDefault     bool    `json:"is_default"`
	Subtotal      float64 `json:"subtotal"`
	ResourceCount int     `json:"resource_count"`
}

type BillingCostByProject struct {
	Projects          []ProjectCost `json:"projects"`
	TotalAttributed   float64       `json:"total_attributed"`
	TotalUnattributed float64       `json:"total_unattributed"`
	TotalMonth        float64       `json:"total_month"`
	Cycle             struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"cycle"`
}

// ── Client methods ────────────────────────────────────────────────────────

func (c *Client) GetBillingUsage() (*BillingUsage, error) {
	respBytes, err := c.Get("/billing/usage/current")
	if err != nil {
		return nil, fmt.Errorf("failed to get billing usage: %w", err)
	}
	var result struct {
		Status string `json:"status"`
		BillingUsage
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse billing usage response: %w", err)
	}
	return &result.BillingUsage, nil
}

func (c *Client) GetBillingInvoices() ([]Invoice, error) {
	respBytes, err := c.Get("/billing")
	if err != nil {
		return nil, fmt.Errorf("failed to get billing invoices: %w", err)
	}
	var result struct {
		Status   string    `json:"status"`
		Invoices []Invoice `json:"invoices"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse billing invoices response: %w", err)
	}
	return result.Invoices, nil
}

func (c *Client) GetBillingCostByProject() (*BillingCostByProject, error) {
	respBytes, err := c.Get("/billing/cost-by-project")
	if err != nil {
		return nil, fmt.Errorf("failed to get billing cost by project: %w", err)
	}
	var result struct {
		Status string `json:"status"`
		BillingCostByProject
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse billing cost by project response: %w", err)
	}
	return &result.BillingCostByProject, nil
}
