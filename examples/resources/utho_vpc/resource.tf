resource "utho_vpc" "example" {
  dcslug  = "innoida"
  name    = "vpc1"
  planid  = "1008"
  network = "10.210.100.0"
  size    = "24"
}
