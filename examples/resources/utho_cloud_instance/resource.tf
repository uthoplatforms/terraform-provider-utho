resource "utho_cloud_instance" "example" {
  name = "example-name"
  # country slug
  dcslug       = "inbangalore"
  image        = "rocky-8.8-x86_64"
  planid       = "10045"
  firewall     = ""
  enablebackup = "false"
  billingcycle = "hourly"
  backupid     = ""
  sshkeys      = ""
  vpc_id       = ""
}
