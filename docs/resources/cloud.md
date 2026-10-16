---
page_title: "Utho: utho_cloud"
subcategory: "Cloud Instances"
description: |-
  Create and manage Utho Cloud instances (virtual machines).
---

# utho_cloud

Creates and manages a Utho Cloud instance — a virtual machine running in one of Utho's data centers. You can deploy from standard OS images, snapshots, backups, ISOs, or marketplace stacks.

## Before You Begin

Creating a cloud instance requires several IDs that you can't guess — plan IDs, image slugs, DC zone slugs, and VPC subnet IDs. Use these data sources to discover valid values before writing your resource block:

| What You Need | Data Source | Key Attribute |
|---------------|-------------|---------------|
| Data center slug (`dcslug`) | [utho_cloud_dczones](../data-sources/cloud_dczones) | `zones[*].slug` |
| Plan ID (`planid`) | [utho_cloud_plans](../data-sources/cloud_plans) | `plans[*].id` |
| OS image slug (`image`) | [utho_cloud_images](../data-sources/cloud_images) | `images[*].image` |
| Snapshot ID (`snapshotid`) | [utho_cloud_snapshots](../data-sources/cloud_snapshots) | `snapshots[*].id` |
| ISO name (`iso`) | [utho_cloud_isos](../data-sources/cloud_isos) | `isos[*].name` |
| VPC subnet ID (`vpc`) | [utho_vpcs](../data-sources/vpcs) | `vpcs[*].subnets[*].id` |
| SSH key ID (`sshkeys`) | [utho_ssh_keys](../data-sources/ssh_keys) | `keys[*].id` — or use `utho_ssh_key.name.id` after creating one |
| Security group ID (`firewall`) | Create first with [utho_firewall](firewall) | `utho_firewall.name.id` |

### Quick lookup example

```hcl
# Step 1 — Find active DC zones
data "utho_cloud_dczones" "all" {}

# Step 2 — Find available plans in your chosen DC
data "utho_cloud_plans" "mumbai" {
  dcslug = "inmumbaizone2"
}

# Step 3 — Find available OS images
data "utho_cloud_images" "ubuntu" {
  distro = "ubuntu"
}

# Run: terraform apply
# Then check outputs to pick the right IDs
output "zones"  { value = data.utho_cloud_dczones.all.zones[*].slug }
output "plans"  { value = [for p in data.utho_cloud_plans.mumbai.plans : { id = p.id, cpu = p.cpu, ram = p.ram, disk = p.disk } if p.disk != "0"] }
output "images" { value = data.utho_cloud_images.ubuntu.images[*].image }
```

---

## Example Usage

~> **Note on `auth` values:** `option1` means root password authentication. `option2` means SSH key authentication. Always pair `option1` with `root_password` and `option2` with `sshkeys`.

~> **Note on `enable_publicip`:** This field takes a string `"true"` or `"false"`, not a boolean — this matches the Utho API format.

### Minimal instance with password login

```hcl
resource "utho_cloud" "web" {
  hostname        = "web-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"        # 2 vCPU / 4 GB / 80 GB — use data.utho_cloud_plans to find others
  billingcycle    = "hourly"
  auth            = "option1"      # option1 = password auth, option2 = SSH key auth
  root_password   = var.root_password
  image           = "ubuntu-22.04-x86_64"   # use data.utho_cloud_images to find others
  enable_publicip = "true"
  cpumodel        = "amd"          # amd or intel — check utho_cloud_dczones for what's available in your DC
}

output "ip" { value = utho_cloud.web.ip }
output "id" { value = utho_cloud.web.id }
```

### Instance with SSH key authentication (recommended)

SSH key authentication is more secure than password auth and works seamlessly with automation tools. Create the key with `utho_ssh_key` first, then reference it here.

```hcl
resource "utho_ssh_key" "deploy" {
  name   = "deploy-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

resource "utho_cloud" "app" {
  hostname        = "app-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
}

# Connect after deployment
# ssh root@<ip>  (or ubuntu@<ip> depending on image)
output "ssh_command" {
  value = "ssh root@${utho_cloud.app.ip}"
}
```

### Scale horizontally with count

Create multiple identical servers in one block. Each gets a unique hostname via `count.index`.

```hcl
resource "utho_ssh_key" "deploy" {
  name   = "deploy-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

resource "utho_cloud" "worker" {
  count = 5

  hostname        = "worker-${count.index + 1}.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
}

output "worker_ips" {
  value = utho_cloud.worker[*].ip
}
```

### Instance inside a private VPC

Attach the instance to a VPC subnet for private networking. Use `data.utho_vpcs` to find your subnet ID.

```hcl
# Find your VPC and subnet first
data "utho_vpcs" "mumbai" {
  dcslug = "inmumbaizone2"
}

locals {
  # Get subnet ID of a specific VPC by name
  my_subnet = one([
    for vpc in data.utho_vpcs.mumbai.vpcs :
    vpc.subnets[0].id
    if vpc.name == "production" && length(vpc.subnets) > 0
  ])
}

resource "utho_cloud" "backend" {
  hostname        = "backend-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "monthly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "false"
  vpc             = local.my_subnet
}
```

### Instance with security group

Create a security group, add rules, then attach to the instance.

```hcl
resource "utho_firewall" "web" {
  name = "web-sg"
}

resource "utho_firewall_rule" "http" {
  firewall_id  = utho_firewall.web.id
  type         = "incoming"
  service      = "HTTP"
  protocol     = "tcp"
  port         = "80"
  port_range   = "80"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}

resource "utho_cloud" "web" {
  hostname        = "web-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "true"
  firewall        = utho_firewall.web.id   # attach at creation
}
```

### Deploy from a snapshot

Useful for deploying pre-configured golden images. Use `data.utho_cloud_snapshots` to find your snapshot ID.

```hcl
data "utho_cloud_snapshots" "all" {}

locals {
  golden = one([
    for s in data.utho_cloud_snapshots.all.snapshots :
    s if s.name == "golden-image-v2" && s.status == "Active"
  ])
}

resource "utho_cloud" "restored" {
  hostname        = "restored-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10308"
  billingcycle    = "hourly"
  auth            = "option1"
  root_password   = var.root_password
  snapshotid      = local.golden.id
  enable_publicip = "true"
}
```

### Instance with additional EBS storage

Attach extra block volumes at creation for databases, logs, or large datasets. Use plans with `disk = "0"` for EBS-only setups.

```hcl
resource "utho_cloud" "db" {
  hostname        = "db-01.mhc"
  dcslug          = "inmumbaizone2"
  planid          = "10336"   # 2 vCPU / 4 GB / no included disk — EBS plan
  billingcycle    = "monthly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = "ubuntu-22.04-x86_64"
  enable_publicip = "false"
  delete_ebs      = true   # delete EBS volumes when instance is destroyed

  ebs = [
    { id = "1", disk = 100, type = "nvme" },
    { id = "2", disk = 500, type = "ssd"  },
  ]
}
```

### Complete production setup — with data sources

A real-world example that uses data sources to dynamically discover all IDs rather than hardcoding them.

```hcl
# Discover available resources
data "utho_cloud_plans" "mumbai" { dcslug = "inmumbaizone2" }
data "utho_cloud_images" "ubuntu" { distro = "ubuntu" }
data "utho_vpcs" "mumbai" { dcslug = "inmumbaizone2" }

locals {
  # Pick cheapest plan with disk included
  plan = [
    for p in data.utho_cloud_plans.mumbai.plans :
    p if p.disk != "0" && p.slug == "basic"
  ][0]

  # Ubuntu 22.04 LTS image
  image = one([
    for img in data.utho_cloud_images.ubuntu.images :
    img if img.image == "ubuntu-22.04-x86_64"
  ])

  # Public subnet from production VPC
  subnet = one([
    for vpc in data.utho_vpcs.mumbai.vpcs :
    vpc.subnets[0]
    if vpc.name == "production" && length(vpc.subnets) > 0
  ])
}

resource "utho_ssh_key" "deploy" {
  name   = "deploy-key"
  sshkey = file("~/.ssh/id_ed25519.pub")
}

resource "utho_firewall" "web" { name = "web-sg" }
resource "utho_firewall_rule" "http" {
  firewall_id  = utho_firewall.web.id
  type         = "incoming"
  service      = "HTTP"
  protocol     = "tcp"
  port         = "80"
  port_range   = "80"
  addresses    = "0.0.0.0/0"
  source_range = "0.0.0.0/0"
}

resource "utho_cloud" "web" {
  count = 3

  hostname        = "web-${count.index + 1}.mhc"
  dcslug          = "inmumbaizone2"
  planid          = local.plan.id      # ← from data source
  billingcycle    = "hourly"
  auth            = "option2"
  sshkeys         = utho_ssh_key.deploy.id
  image           = local.image.image  # ← from data source
  enable_publicip = "true"
  vpc             = local.subnet.id    # ← from data source
  firewall        = utho_firewall.web.id
}

output "server_ips" { value = utho_cloud.web[*].ip }
```

## Argument Reference

### Required

| Argument       | Type   | Description |
|----------------|--------|-------------|
| `hostname`     | String | Hostname for the instance. |
| `dcslug`       | String | Data center slug. See [utho_cloud_dczones](../data-sources/cloud_dczones) for valid values. |
| `planid`       | String | Plan ID for the instance size. See [utho_cloud_plans](../data-sources/cloud_plans) for valid values. |
| `billingcycle` | String | `hourly`, `monthly`, `3month`, `6month`, `12month`, `24month`, or `36month`. Use [utho_billing_cycles](../data-sources/billing_cycles) to list valid values. |
| `auth`         | String | `option1` = root password login, `option2` = SSH key login. |

### Authentication — one required

| Argument       | Type   | When Required | Description |
|----------------|--------|---------------|-------------|
| `root_password`| String | `auth = "option1"` | Root password. **Sensitive.** |
| `sshkeys`      | String | `auth = "option2"` | SSH key ID from `utho_ssh_key`. |

### Image source — one required

| Argument     | Type   | Description |
|--------------|--------|-------------|
| `image`      | String | OS image slug. See [utho_cloud_images](../data-sources/cloud_images) for valid values. |
| `snapshotid` | String | Snapshot ID to restore from. See [utho_cloud_snapshots](../data-sources/cloud_snapshots). |
| `backupid`   | String | Backup ID to restore from. |
| `iso`        | String | ISO name for custom OS installs. See [utho_cloud_isos](../data-sources/cloud_isos). |
| `stack`      | String | Marketplace stack ID. |

### Optional

| Argument          | Type   | Description |
|-------------------|--------|-------------|
| `enable_publicip` | String | `"true"` or `"false"`. Default: `"true"`. |
| `vpc`             | String | VPC subnet ID. See [utho_vpcs](../data-sources/vpcs) for valid values. |
| `firewall`        | String | Security group ID to attach at creation. |
| `cpumodel`        | String | Required. CPU model: `amd` or `intel`. Must match what is available in the chosen DC — use [utho_cloud_dczones](../data-sources/cloud_dczones) to check `default_cpu`. |
| `enablebackup`    | String | Enable automated backups: `"true"` or `"false"`. |
| `support`         | String | `unmanaged` or `managed`. |
| `delete_ebs`      | Bool   | Delete attached EBS volumes on destroy. Default: `false`. |
| `ebs`             | List   | EBS volumes to attach at creation. See [EBS Block](#ebs-block). |

### EBS Block

```hcl
ebs = [
  { id = "1", disk = 100, type = "nvme" },
  { id = "2", disk = 500, type = "ssd"  },
]
```

| Argument | Type   | Description |
|----------|--------|-------------|
| `id`     | String | Sequential identifier: `"1"`, `"2"`, etc. |
| `disk`   | Number | Disk size in GB. |
| `type`   | String | `nvme` (faster, recommended) or `ssd`. |

## Attribute Reference

| Attribute      | Type   | Description |
|----------------|--------|-------------|
| `id`           | String | Unique instance ID. Referenced by `utho_firewall_server`, `utho_cloud_snapshot`, `utho_loadbalancer_backend`, etc. |
| `ip`           | String | Primary public IP address. Point your DNS A record here. |
| `status`       | String | Instance status (e.g. `Active`). |
| `power_status` | String | Power state (`Running`, `Shutdown`). |
| `created_at`   | String | Creation timestamp (UTC). |

## Import

Bring an existing instance under Terraform management without recreating it.

```bash
terraform import utho_cloud.web 1671990
```

After importing, fill in the required fields in your `.tf` file and run `terraform plan`:

```hcl
resource "utho_cloud" "web" {
  hostname     = "web-01.mhc"
  dcslug       = "inmumbaizone2"
  planid       = "10308"
  billingcycle = "hourly"
  auth         = "option2"
  sshkeys      = utho_ssh_key.deploy.id
}
```

~> **Note:** After import, write-only fields (`auth`, `planid`, `image`, `root_password`) will show as diffs in `terraform plan` because the API does not return them on GET requests. This is expected — fill them in manually and the diff will resolve.

## Notes

- Changing `hostname`, `dcslug`, `planid`, or `image` requires destroying and recreating the instance.
- For zero-downtime updates, create the new instance first, migrate traffic, then destroy the old one.
- Use `utho_cloud_power` to start/stop/reboot without recreating.
- Use `utho_cloud_resize` to change the plan after creation.
- Use `utho_cloud_snapshot` to take snapshots before destructive operations.

## Related Resources

Once your instance is running, use these resources to manage it:

| Resource | Purpose |
|----------|---------|
| [utho_cloud_power](cloud_power) | Start, stop, or reboot the instance |
| [utho_cloud_resize](cloud_resize) | Resize to a different plan |
| [utho_cloud_snapshot](cloud_snapshot) | Take a point-in-time snapshot |
| [utho_cloud_firewall](cloud_firewall) | Attach a security group post-deploy |
| [utho_cloud_vpc](cloud_vpc) | Attach a VPC subnet post-deploy |
| [utho_cloud_public_ip](cloud_public_ip) | Assign an additional public IP |
| [utho_cloud_storage](cloud_storage) | Add a general storage disk |
| [utho_cloud_ebs](cloud_ebs) | Attach an EBS block volume |
| [utho_cloud_iso](cloud_iso) | Mount or unmount an ISO |
