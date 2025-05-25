resource "utho_cloud_instance" "example" {
  name = "example-name"
  # country slug
  dcslug          = "inmumbaizone2"
  image           = "ubuntu-22.04-x86_64"
  planid          = "10045"
  enablebackup    = "false"
  billingcycle    = "hourly"
  firewall        = "23432614"
  vpc_id          = "4de5f07a-f51c-4323-b39a-ef66130e1bd9"
  cpumodel        = "amd"
  enable_publicip = "true"
  root_password   = "qwe123"
}
