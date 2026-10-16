# ============================================================
# Example: Billing Module
# ============================================================
# Read-only data sources for account billing information.
# Useful for cost monitoring, budgeting, and automation.
# ============================================================

terraform {
  required_providers {
    utho = {
      source  = "nitinuthocloud/utho"
      version = "~> 0.7"
    }
  }
  required_version = ">= 1.0"
}

provider "utho" {
  api_key = var.utho_api_key
}

variable "utho_api_key" {
  type      = string
  sensitive = true
}

# ── Current Usage ──────────────────────────────────────────
data "utho_billing_usage" "current" {}

# ── All Invoices ───────────────────────────────────────────
data "utho_billing_invoices" "all" {}

# ── Cost by Project ────────────────────────────────────────
data "utho_billing_cost_by_project" "current" {}

# ── Outputs ────────────────────────────────────────────────
output "billing_cycle" {
  value = {
    start     = data.utho_billing_usage.current.cycle_start
    end       = data.utho_billing_usage.current.cycle_end
    days_in   = data.utho_billing_usage.current.days_in
    days_left = data.utho_billing_usage.current.days_left
  }
}

output "billing_summary" {
  value = {
    usage_excl_tax      = data.utho_billing_usage.current.usage_excl_tax
    tax                 = data.utho_billing_usage.current.tax
    total_before_wallet = data.utho_billing_usage.current.total_before_wallet
    due_amount          = data.utho_billing_usage.current.due_amount
    wallet_balance      = data.utho_billing_usage.current.wallet_balance
    available_credit    = data.utho_billing_usage.current.available_credit
  }
}

output "top_resources"   { value = data.utho_billing_usage.current.top_resources }
output "usage_by_category" { value = data.utho_billing_usage.current.by_category }
output "invoice_count"   { value = length(data.utho_billing_invoices.all.invoices) }
output "total_month_spend" { value = data.utho_billing_cost_by_project.current.total_month }
output "project_costs"   { value = data.utho_billing_cost_by_project.current.projects }
