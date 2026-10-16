# ============================================================
# Example: App Server + Database
# ============================================================
# Deploys a cloud instance and a managed database together.
# The app server connects to the database using the
# connection string provided in outputs.
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

variable "ssh_public_key_path" {
  type    = string
  default = "~/.ssh/id_ed25519.pub"
}

# ── SSH Key ────────────────────────────────────────────────

resource "utho_ssh_key" "deploy" {
  name   = "deploy-key"
  sshkey = file(var.ssh_public_key_path)
}

# ── App Server ─────────────────────────────────────────────

resource "utho_cloud" "app" {
  hostname        = "app-server.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10351"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
  cpumodel        = "amd"
  ebs = [
    {
      id   = "1"
      disk = "30"
      type = "nvme"
    }
  ]
}

# ── Database ───────────────────────────────────────────────

resource "utho_database" "main" {
  cluster_name  = "app-database"
  dcslug        = "inmumbaizone2"
  engine        = "mysql"
  version       = "8.0"
  size          = "10151"
  network_type  = "public"
  billing       = "hourly"
  replica_count = "0"
}

# ── Outputs ────────────────────────────────────────────────

output "app_ip"      { value = utho_cloud.app.ip }
output "db_host"     { value = utho_database.main.host }
output "db_port"     { value = utho_database.main.port }
output "db_user"     { value = utho_database.main.default_user }
output "ssh_command" { value = "ssh root@${utho_cloud.app.ip}" }
output "db_connect"  { value = "mysql -h ${utho_database.main.host} -P ${utho_database.main.port} -u ${utho_database.main.default_user} -p" }
output "db_password" {
  value     = utho_database.main.default_pass
  sensitive = true
}
