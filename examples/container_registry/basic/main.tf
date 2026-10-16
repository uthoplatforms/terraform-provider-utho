# ============================================================
# Example: Basic Container Registry
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

resource "utho_container_registry" "main" {
  project_name = "my-registry"
  dcslug       = "innoida"
  planid       = 10276
  billingcycle = "monthly"
  public       = false
}

output "registry_id"   { value = utho_container_registry.main.id }
output "registry_name" { value = utho_container_registry.main.project_name }
